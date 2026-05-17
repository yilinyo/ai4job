package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

// Dependencies holds all handler and middleware dependencies for the router.
type Dependencies struct {
	UserHandler              *handler.UserHandler
	ResumeRepoHandler        *handler.ResumeRepoHandler
	ResumeVersionHandler     *handler.ResumeVersionHandler
	JobPostingHandler        *handler.JobPostingHandler
	ResumeApplicationHandler *handler.ResumeApplicationHandler
	InterviewProgressHandler *handler.InterviewProgressHandler
	InterviewQuestionHandler *handler.InterviewQuestionHandler
	AuthMiddleware           gin.HandlerFunc
}

// NewRouter creates and configures the Gin engine with all routes.
func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")

	user := v1.Group("/user")
	{
		user.POST("/email-code", deps.UserHandler.SendEmailCode)
		user.POST("/register", deps.UserHandler.Register)
		user.POST("/login", deps.UserHandler.Login)

		authed := user.Group("/", deps.AuthMiddleware)
		{
			authed.GET("/profile", deps.UserHandler.GetProfile)
			authed.PUT("/profile", deps.UserHandler.UpdateProfile)
			authed.POST("/logout", deps.UserHandler.Logout)
		}
	}

	auth := v1.Group("/", deps.AuthMiddleware)
	{
		auth.GET("/resume-repos", deps.ResumeRepoHandler.List)
		auth.POST("/resume-repos", deps.ResumeRepoHandler.Create)

		auth.GET("/resume-versions", deps.ResumeVersionHandler.List)
		auth.POST("/resume-versions", deps.ResumeVersionHandler.Create)
		auth.GET("/resume-versions/compare", deps.ResumeVersionHandler.Compare)
		auth.GET("/resume-versions/:id", deps.ResumeVersionHandler.Get)
		auth.PUT("/resume-versions/:id", deps.ResumeVersionHandler.Update)
		auth.DELETE("/resume-versions/:id", deps.ResumeVersionHandler.Delete)
		auth.POST("/resume-versions/:id/activate", deps.ResumeVersionHandler.Activate)

		auth.GET("/job-postings", deps.JobPostingHandler.List)
		auth.POST("/job-postings", deps.JobPostingHandler.Create)
		auth.POST("/job-postings/parse", deps.JobPostingHandler.Parse)
		auth.GET("/job-postings/:id", deps.JobPostingHandler.Get)
		auth.PUT("/job-postings/:id", deps.JobPostingHandler.Update)
		auth.DELETE("/job-postings/:id", deps.JobPostingHandler.Delete)

		auth.GET("/resume-applications", deps.ResumeApplicationHandler.List)
		auth.POST("/resume-applications", deps.ResumeApplicationHandler.Create)
		auth.GET("/resume-applications/:id", deps.ResumeApplicationHandler.Get)
		auth.PUT("/resume-applications/:id", deps.ResumeApplicationHandler.Update)
		auth.DELETE("/resume-applications/:id", deps.ResumeApplicationHandler.Delete)

		auth.GET("/interview-progress", deps.InterviewProgressHandler.List)
		auth.POST("/interview-progress", deps.InterviewProgressHandler.Create)
		auth.PUT("/interview-progress/:id", deps.InterviewProgressHandler.Update)
		auth.DELETE("/interview-progress/:id", deps.InterviewProgressHandler.Delete)

		auth.GET("/interview-questions", deps.InterviewQuestionHandler.ListQuestions)
		auth.POST("/interview-questions", deps.InterviewQuestionHandler.CreateQuestion)
		auth.GET("/interview-questions/:id", deps.InterviewQuestionHandler.GetQuestion)
		auth.PUT("/interview-questions/:id", deps.InterviewQuestionHandler.UpdateQuestion)
		auth.DELETE("/interview-questions/:id", deps.InterviewQuestionHandler.DeleteQuestion)
		auth.PUT("/interview-questions/:id/favorite", deps.InterviewQuestionHandler.UpdateFavorite)

		auth.GET("/question-occurrences", deps.InterviewQuestionHandler.ListOccurrences)
		auth.POST("/question-occurrences", deps.InterviewQuestionHandler.CreateOccurrence)
		auth.DELETE("/question-occurrences/:id", deps.InterviewQuestionHandler.DeleteOccurrence)

		auth.GET("/flashcards", deps.InterviewQuestionHandler.ListFlashcards)
		auth.POST("/flashcard-reviews", deps.InterviewQuestionHandler.CreateFlashcardReview)

		auth.GET("/resume-repos/:id", deps.ResumeRepoHandler.Get)
		auth.PUT("/resume-repos/:id", deps.ResumeRepoHandler.Update)
		auth.DELETE("/resume-repos/:id", deps.ResumeRepoHandler.Delete)
	}

	return r
}
