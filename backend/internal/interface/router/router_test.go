package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
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
		AuthMiddleware: func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		},
	})

	cases := []struct {
		method string
		path   string
	}{
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
