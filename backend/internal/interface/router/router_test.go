package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	"github.com/yilin/ai-for-backend/internal/interface/router"
)

type stubUserService struct{}

func (s *stubUserService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {
	return nil, nil
}
func (s *stubUserService) Register(ctx context.Context, req dto.RegisterRequest) error { return nil }
func (s *stubUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	return nil, nil
}
func (s *stubUserService) GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error) {
	return nil, nil
}
func (s *stubUserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error) {
	return nil, nil
}

type stubResumeRepoService struct{}

func (s *stubResumeRepoService) Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error {
	return nil
}
func (s *stubResumeRepoService) List(ctx context.Context, userID string, jobType domainrepo.JobType, keyword string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
	return []domainrepo.ResumeRepo{}, 0, nil
}
func (s *stubResumeRepoService) Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error) {
	return &domainrepo.ResumeRepo{ID: id, UserID: userID, Name: "Main"}, nil
}
func (s *stubResumeRepoService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error {
	return nil
}
func (s *stubResumeRepoService) Delete(ctx context.Context, userID, id string) error { return nil }

type stubResumeVersionService struct {
	lastRepoID    string
	lastVersionID string
}

func (s *stubResumeVersionService) Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error {
	return nil
}
func (s *stubResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	s.lastRepoID = repoID
	return []domainversion.ResumeVersion{}, nil
}
func (s *stubResumeVersionService) Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	s.lastVersionID = id
	return &domainversion.ResumeVersion{ID: id, RepoID: "r1", Title: "v1", Content: "{}", Version: "v1.0.0"}, nil
}
func (s *stubResumeVersionService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error {
	return nil
}
func (s *stubResumeVersionService) Delete(ctx context.Context, userID, id string) error {
	return nil
}
func (s *stubResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	return nil
}

type stubJobPostingService struct{}

func (s *stubJobPostingService) Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) (*domainjob.JobPosting, error) {
	return &domainjob.JobPosting{ID: "job-1", UserID: userID, CompanyName: "Acme", JobTitle: "Backend", Status: domainjob.StatusPending}, nil
}
func (s *stubJobPostingService) List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	return []domainjob.JobPosting{}, 0, nil
}
func (s *stubJobPostingService) Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error) {
	return &domainjob.JobPosting{ID: id, UserID: userID, CompanyName: "Acme", JobTitle: "Backend", Status: domainjob.StatusPending}, nil
}
func (s *stubJobPostingService) Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error {
	return nil
}
func (s *stubJobPostingService) Delete(ctx context.Context, userID, id string) error { return nil }
func (s *stubJobPostingService) Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error) {
	return &dto.JobPostingResponse{CompanyName: "Acme", JobTitle: "Backend"}, nil
}

type stubResumeApplicationService struct{}

func (s *stubResumeApplicationService) Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) (*domainapp.ResumeApplication, error) {
	return &domainapp.ResumeApplication{ID: "app-1", UserID: userID, ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Backend", CurrentStatus: domainapp.StatusApplied}, nil
}
func (s *stubResumeApplicationService) List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error) {
	return []domainapp.ResumeApplication{}, 0, nil
}
func (s *stubResumeApplicationService) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	return &domainapp.ResumeApplication{ID: id, UserID: userID, ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Backend", CurrentStatus: domainapp.StatusApplied}, nil
}
func (s *stubResumeApplicationService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error {
	return nil
}
func (s *stubResumeApplicationService) Delete(ctx context.Context, userID, id string) error {
	return nil
}

type stubInterviewProgressService struct{}

func (s *stubInterviewProgressService) Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	return &domainprogress.InterviewProgress{ID: "p1", UserID: userID, ApplicationID: req.ApplicationID, Stage: req.Stage, OccurredAt: req.OccurredAt.Value(), Result: req.Result}, nil
}
func (s *stubInterviewProgressService) List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error) {
	return []domainprogress.InterviewProgress{}, nil
}
func (s *stubInterviewProgressService) Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	return &domainprogress.InterviewProgress{ID: id, UserID: userID}, nil
}
func (s *stubInterviewProgressService) Delete(ctx context.Context, userID, id string) error {
	return nil
}

type stubInterviewQuestionService struct{}

func (s *stubInterviewQuestionService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	return &domainquestion.InterviewQuestion{ID: "q1", UserID: userID, Title: req.Title, Content: req.Content, AnswerAnalysis: req.AnswerAnalysis, Tags: domainquestion.EncodeTags(req.Tags), Difficulty: req.Difficulty, QuestionType: req.QuestionType, MasteryStatus: req.MasteryStatus}, nil
}
func (s *stubInterviewQuestionService) ListQuestions(ctx context.Context, userID string, req dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	return []questionrepo.QuestionWithOccurrenceCount{}, 0, nil
}
func (s *stubInterviewQuestionService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "Question", Content: "Content", AnswerAnalysis: "Analysis", Tags: domainquestion.EncodeTags(nil), Difficulty: domainquestion.DifficultyEasy, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusUnreviewed}, 0, []domainoccurrence.QuestionOccurrence{}, domainreview.ReviewSummary{}, []domainreview.FlashcardReview{}, nil
}
func (s *stubInterviewQuestionService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "Question", Content: "Content", AnswerAnalysis: "Analysis", Tags: domainquestion.EncodeTags(nil), Difficulty: domainquestion.DifficultyEasy, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusUnreviewed}, nil
}
func (s *stubInterviewQuestionService) DeleteQuestion(ctx context.Context, userID, id string) error {
	return nil
}
func (s *stubInterviewQuestionService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	favorite := false
	if req.IsFavorite != nil {
		favorite = *req.IsFavorite
	}
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "Question", Content: "Content", AnswerAnalysis: "Analysis", Tags: domainquestion.EncodeTags(nil), Difficulty: domainquestion.DifficultyEasy, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusUnreviewed, IsFavorite: favorite}, nil
}
func (s *stubInterviewQuestionService) ListOccurrences(ctx context.Context, userID string, req dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	return []domainoccurrence.QuestionOccurrence{}, nil
}
func (s *stubInterviewQuestionService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	return &domainoccurrence.QuestionOccurrence{ID: "occ-1", UserID: userID, QuestionID: req.QuestionID, ApplicationID: req.ApplicationID, CompanyNameSnapshot: "Acme", JobTitleSnapshot: "Backend", OccurredAt: req.OccurredAt.Value()}, nil
}
func (s *stubInterviewQuestionService) DeleteOccurrence(ctx context.Context, userID, id string) error {
	return nil
}
func (s *stubInterviewQuestionService) ListFlashcards(ctx context.Context, userID string, req dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	return []questionrepo.QuestionWithOccurrenceCount{}, 0, nil
}
func (s *stubInterviewQuestionService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	return &domainreview.FlashcardReview{ID: "rev-1", UserID: userID, QuestionID: req.QuestionID, ReviewedAt: req.ReviewedAt.Value(), Result: req.Result, ReviewMode: req.ReviewMode}, questionrepo.QuestionWithOccurrenceCount{Question: domainquestion.InterviewQuestion{ID: req.QuestionID, UserID: userID, Title: "Question", Content: "Content", AnswerAnalysis: "Analysis", Tags: domainquestion.EncodeTags(nil), Difficulty: domainquestion.DifficultyEasy, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusMastered}}, domainreview.ReviewSummary{}, nil
}

type stubDashboardService struct{}

func (s *stubDashboardService) ListActivities(ctx context.Context, userID, rawRange string, page, pageSize int) (*dto.ActivityPageResponse, error) {
	return &dto.ActivityPageResponse{}, nil
}

func (s *stubDashboardService) ListTodosPage(ctx context.Context, userID string, completed *bool, page, pageSize int) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *stubDashboardService) Get(ctx context.Context, userID, rawRange string, activityPage, activityPageSize, todoPage, todoPageSize int, includeCompleted bool) (*dto.DashboardResponse, error) {
	return &dto.DashboardResponse{Range: rawRange}, nil
}
func (s *stubDashboardService) ListTodos(ctx context.Context, userID string, completed *bool) ([]dto.DashboardTodoResponse, error) {
	return []dto.DashboardTodoResponse{}, nil
}
func (s *stubDashboardService) CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	return &domaintodo.DashboardTodo{ID: "todo-1", UserID: userID, Title: req.Title}, nil
}
func (s *stubDashboardService) RouteForTodo(ctx context.Context, userID string, targetType domaintodo.TargetType, targetID *string) (string, error) {
	return "", nil
}
func (s *stubDashboardService) UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	title := "Todo"
	if req.Title != nil {
		title = *req.Title
	}
	return &domaintodo.DashboardTodo{ID: id, UserID: userID, Title: title}, nil
}
func (s *stubDashboardService) CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error) {
	return &domaintodo.DashboardTodo{ID: id, UserID: userID, Title: "Todo"}, nil
}
func (s *stubDashboardService) DeleteTodo(ctx context.Context, userID, id string) error { return nil }

func TestRouter_Healthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.NewRouter(router.Dependencies{
		UserHandler:              handler.NewUserHandler(&stubUserService{}),
		ResumeRepoHandler:        handler.NewResumeRepoHandler(&stubResumeRepoService{}),
		ResumeVersionHandler:     handler.NewResumeVersionHandler(&stubResumeVersionService{}),
		JobPostingHandler:        handler.NewJobPostingHandler(&stubJobPostingService{}),
		ResumeApplicationHandler: handler.NewResumeApplicationHandler(&stubResumeApplicationService{}),
		InterviewProgressHandler: handler.NewInterviewProgressHandler(&stubInterviewProgressService{}),
		InterviewQuestionHandler: handler.NewInterviewQuestionHandler(&stubInterviewQuestionService{}),
		DashboardHandler:         handler.NewDashboardHandler(&stubDashboardService{}),
		AuthMiddleware:           func(c *gin.Context) { c.Next() },
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRouter_ResumeRoutesRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	versionSvc := &stubResumeVersionService{}
	r := router.NewRouter(router.Dependencies{
		UserHandler:              handler.NewUserHandler(&stubUserService{}),
		ResumeRepoHandler:        handler.NewResumeRepoHandler(&stubResumeRepoService{}),
		ResumeVersionHandler:     handler.NewResumeVersionHandler(versionSvc),
		JobPostingHandler:        handler.NewJobPostingHandler(&stubJobPostingService{}),
		ResumeApplicationHandler: handler.NewResumeApplicationHandler(&stubResumeApplicationService{}),
		InterviewProgressHandler: handler.NewInterviewProgressHandler(&stubInterviewProgressService{}),
		InterviewQuestionHandler: handler.NewInterviewQuestionHandler(&stubInterviewQuestionService{}),
		DashboardHandler:         handler.NewDashboardHandler(&stubDashboardService{}),
		AuthMiddleware: func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		},
	})

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/dashboard?range=30d"},
		{http.MethodGet, "/api/v1/dashboard/todos?completed=false"},
		{http.MethodPost, "/api/v1/dashboard/todos"},
		{http.MethodPut, "/api/v1/dashboard/todos/todo-1"},
		{http.MethodPut, "/api/v1/dashboard/todos/todo-1/complete"},
		{http.MethodDelete, "/api/v1/dashboard/todos/todo-1"},
		{http.MethodGet, "/api/v1/resume-repos"},
		{http.MethodGet, "/api/v1/resume-repos/r1"},
		{http.MethodGet, "/api/v1/resume-versions?repoId=r1"},
		{http.MethodPost, "/api/v1/resume-versions"},
		{http.MethodGet, "/api/v1/resume-versions/compare"},
		{http.MethodGet, "/api/v1/resume-versions/v1"},
		{http.MethodPut, "/api/v1/resume-versions/v1"},
		{http.MethodDelete, "/api/v1/resume-versions/v1"},
		{http.MethodPost, "/api/v1/resume-versions/v1/activate"},
		{http.MethodGet, "/api/v1/job-postings"},
		{http.MethodPost, "/api/v1/job-postings"},
		{http.MethodPost, "/api/v1/job-postings/parse"},
		{http.MethodGet, "/api/v1/job-postings/job-1"},
		{http.MethodPut, "/api/v1/job-postings/job-1"},
		{http.MethodDelete, "/api/v1/job-postings/job-1"},
		{http.MethodGet, "/api/v1/resume-applications"},
		{http.MethodPost, "/api/v1/resume-applications"},
		{http.MethodGet, "/api/v1/resume-applications/app-1"},
		{http.MethodPut, "/api/v1/resume-applications/app-1"},
		{http.MethodDelete, "/api/v1/resume-applications/app-1"},
		{http.MethodGet, "/api/v1/interview-progress?applicationId=app-1"},
		{http.MethodPost, "/api/v1/interview-progress"},
		{http.MethodPut, "/api/v1/interview-progress/p1"},
		{http.MethodDelete, "/api/v1/interview-progress/p1"},
		{http.MethodGet, "/api/v1/interview-questions"},
		{http.MethodPost, "/api/v1/interview-questions"},
		{http.MethodGet, "/api/v1/interview-questions/q1"},
		{http.MethodPut, "/api/v1/interview-questions/q1"},
		{http.MethodDelete, "/api/v1/interview-questions/q1"},
		{http.MethodPut, "/api/v1/interview-questions/q1/favorite"},
		{http.MethodGet, "/api/v1/question-occurrences?questionId=q1"},
		{http.MethodPost, "/api/v1/question-occurrences"},
		{http.MethodDelete, "/api/v1/question-occurrences/occ-1"},
		{http.MethodGet, "/api/v1/flashcards"},
		{http.MethodPost, "/api/v1/flashcard-reviews"},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusNotFound, w.Code, tc.path)
	}
	assert.Equal(t, "r1", versionSvc.lastRepoID)
	assert.Equal(t, "v1", versionSvc.lastVersionID)
}

func TestRouter_ResumeVersionCompareUsesStaticRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	versionSvc := &stubResumeVersionService{}
	r := router.NewRouter(router.Dependencies{
		UserHandler:          handler.NewUserHandler(&stubUserService{}),
		ResumeRepoHandler:    handler.NewResumeRepoHandler(&stubResumeRepoService{}),
		ResumeVersionHandler: handler.NewResumeVersionHandler(versionSvc),
		AuthMiddleware: func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions/compare", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "暂不支持版本对比")
	assert.Empty(t, versionSvc.lastVersionID)
}
