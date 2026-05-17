package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	interviewprogressapp "github.com/yilin/ai-for-backend/internal/application/interview_progress"
	interviewquestionapp "github.com/yilin/ai-for-backend/internal/application/interview_question"
	jobpostingapp "github.com/yilin/ai-for-backend/internal/application/job_posting"
	resumeapplicationapp "github.com/yilin/ai-for-backend/internal/application/resume_application"
	resumerepoapp "github.com/yilin/ai-for-backend/internal/application/resume_repo"
	resumeversionapp "github.com/yilin/ai-for-backend/internal/application/resume_version"
	"github.com/yilin/ai-for-backend/internal/application/user"
	"github.com/yilin/ai-for-backend/internal/infrastructure/auth"
	infraCache "github.com/yilin/ai-for-backend/internal/infrastructure/cache"
	infraDB "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	"github.com/yilin/ai-for-backend/internal/interface/middleware"
	"github.com/yilin/ai-for-backend/internal/interface/router"
	"github.com/yilin/ai-for-backend/pkg/config"
)

// jwtParserAdapter adapts auth.JWTManager to middleware.JWTParser interface.
type jwtParserAdapter struct {
	mgr *auth.JWTManager
}

func (a *jwtParserAdapter) Parse(tokenStr string) (userID, username, jti string, err error) {
	claims, parseErr := a.mgr.Parse(tokenStr)
	if parseErr != nil {
		return "", "", "", parseErr
	}
	return claims.UserID, claims.Username, claims.ID, nil
}

// blacklistAdapter adapts infraCache.TokenCache to middleware.Blacklist interface.
type blacklistAdapter struct {
	cache *infraCache.TokenCache
}

func (a *blacklistAdapter) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	return a.cache.IsBlacklisted(ctx, jti)
}

func main() {
	cfg := config.Load()

	// --- Infrastructure: Database ---
	db, err := infraDB.NewDB(cfg.DSN(), cfg.IsDevelopment())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// --- Infrastructure: Redis ---
	redisClient, err := infraCache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	tokenCache := infraCache.NewTokenCache(redisClient)

	// --- Infrastructure: JWT ---
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpireHours)

	// --- Repositories ---
	userRepo := infraDB.NewUserRepo(db)
	resumeRepoRepo := infraDB.NewResumeRepoRepo(db)
	resumeVersionRepo := infraDB.NewResumeVersionRepo(db)
	jobPostingRepo := infraDB.NewJobPostingRepo(db)
	resumeApplicationRepo := infraDB.NewResumeApplicationRepo(db)
	interviewProgressRepo := infraDB.NewInterviewProgressRepo(db)
	interviewQuestionRepo := infraDB.NewInterviewQuestionRepo(db)
	questionOccurrenceRepo := infraDB.NewQuestionOccurrenceRepo(db)
	flashcardReviewRepo := infraDB.NewFlashcardReviewRepo(db)
	transactionRunner := infraDB.NewTransactionRunner(db)

	// --- Application Services ---
	userSvc := user.NewUserService(userRepo, tokenCache, jwtMgr, cfg.EmailMock)
	resumeRepoSvc := resumerepoapp.NewResumeRepoService(resumeRepoRepo)
	jobPostingSvc := jobpostingapp.NewJobPostingService(jobPostingRepo)
	resumeApplicationSvc := resumeapplicationapp.NewResumeApplicationService(resumeApplicationRepo, jobPostingRepo, resumeRepoRepo, resumeVersionRepo)
	interviewProgressSvc := interviewprogressapp.NewInterviewProgressService(interviewProgressRepo, resumeApplicationSvc)
	resumeVersionSvc := resumeversionapp.NewResumeVersionService(resumeRepoSvc, resumeVersionRepo, resumeApplicationRepo)
	interviewQuestionSvc := interviewquestionapp.NewInterviewQuestionServiceWithTransactions(interviewQuestionRepo, questionOccurrenceRepo, flashcardReviewRepo, resumeApplicationSvc, interviewProgressRepo, transactionRunner)

	// --- Middleware ---
	authMiddleware := middleware.Auth(
		&jwtParserAdapter{mgr: jwtMgr},
		&blacklistAdapter{cache: tokenCache},
	)

	// --- Handlers ---
	userHandler := handler.NewUserHandler(userSvc)
	resumeRepoHandler := handler.NewResumeRepoHandler(resumeRepoSvc)
	resumeVersionHandler := handler.NewResumeVersionHandler(resumeVersionSvc)
	jobPostingHandler := handler.NewJobPostingHandler(jobPostingSvc)
	resumeApplicationHandler := handler.NewResumeApplicationHandler(resumeApplicationSvc)
	interviewProgressHandler := handler.NewInterviewProgressHandler(interviewProgressSvc)
	interviewQuestionHandler := handler.NewInterviewQuestionHandler(interviewQuestionSvc)

	// --- Router ---
	r := router.NewRouter(router.Dependencies{
		UserHandler:              userHandler,
		ResumeRepoHandler:        resumeRepoHandler,
		ResumeVersionHandler:     resumeVersionHandler,
		JobPostingHandler:        jobPostingHandler,
		ResumeApplicationHandler: resumeApplicationHandler,
		InterviewProgressHandler: interviewProgressHandler,
		InterviewQuestionHandler: interviewQuestionHandler,
		AuthMiddleware:           authMiddleware,
	})

	// --- HTTP Server ---
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: r,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}
