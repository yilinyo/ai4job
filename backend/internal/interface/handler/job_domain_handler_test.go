package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockJobPostingService struct{ mock.Mock }

func (m *mockJobPostingService) Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) (*domainjob.JobPosting, error) {
	args := m.Called(ctx, userID, req)
	posting, _ := args.Get(0).(*domainjob.JobPosting)
	return posting, args.Error(1)
}

func (m *mockJobPostingService) List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	args := m.Called(ctx, userID, keyword, sourcePlatform, status, page, pageSize)
	return args.Get(0).([]domainjob.JobPosting), args.Get(1).(int64), args.Error(2)
}

func (m *mockJobPostingService) Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error) {
	args := m.Called(ctx, userID, id)
	posting, _ := args.Get(0).(*domainjob.JobPosting)
	return posting, args.Error(1)
}

func (m *mockJobPostingService) Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error {
	return m.Called(ctx, userID, id, req).Error(0)
}

func (m *mockJobPostingService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *mockJobPostingService) Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error) {
	args := m.Called(ctx, req)
	resp, _ := args.Get(0).(*dto.JobPostingResponse)
	return resp, args.Error(1)
}

type mockResumeApplicationJobService struct{ mock.Mock }

func (m *mockResumeApplicationJobService) Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) (*domainapp.ResumeApplication, error) {
	args := m.Called(ctx, userID, req)
	application, _ := args.Get(0).(*domainapp.ResumeApplication)
	return application, args.Error(1)
}

func (m *mockResumeApplicationJobService) List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error) {
	args := m.Called(ctx, userID, resumeVersionID, repoID, keyword, currentStatus, page, pageSize)
	return args.Get(0).([]domainapp.ResumeApplication), args.Get(1).(int64), args.Error(2)
}

func (m *mockResumeApplicationJobService) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	args := m.Called(ctx, userID, id)
	application, _ := args.Get(0).(*domainapp.ResumeApplication)
	return application, args.Error(1)
}

func (m *mockResumeApplicationJobService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error {
	return m.Called(ctx, userID, id, req).Error(0)
}

func (m *mockResumeApplicationJobService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

type mockInterviewProgressJobService struct{ mock.Mock }

func (m *mockInterviewProgressJobService) Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, userID, req)
	progress, _ := args.Get(0).(*domainprogress.InterviewProgress)
	return progress, args.Error(1)
}

func (m *mockInterviewProgressJobService) List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, userID, applicationID)
	return args.Get(0).([]domainprogress.InterviewProgress), args.Error(1)
}

func (m *mockInterviewProgressJobService) Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, userID, id, req)
	progress, _ := args.Get(0).(*domainprogress.InterviewProgress)
	return progress, args.Error(1)
}

func (m *mockInterviewProgressJobService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func TestJobPostingHandler_ListUsesDTOShapeAndSafePagination(t *testing.T) {
	svc := &mockJobPostingService{}
	createdAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	svc.On("List", mock.Anything, "u1", "backend", "Boss", string(domainjob.StatusPending), 1, 10).Return([]domainjob.JobPosting{{
		ID:          "job-1",
		UserID:      "u1",
		CompanyName: "Acme",
		JobTitle:    "Backend Engineer",
		Status:      domainjob.StatusPending,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}}, int64(1), nil)
	r := gin.New()
	h := handler.NewJobPostingHandler(svc)
	r.GET("/job-postings", func(c *gin.Context) { c.Set("userID", "u1"); h.List(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/job-postings?keyword=backend&sourcePlatform=Boss&status=pending&pageSize=999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	data := responseData(t, w)
	list := data["list"].([]any)
	first := list[0].(map[string]any)
	assert.Equal(t, "Acme", first["companyName"])
	assert.Equal(t, "Backend Engineer", first["jobTitle"])
	assert.NotContains(t, first, "CompanyName")
	assert.NotContains(t, first, "UserID")
	assert.Equal(t, float64(10), data["pageSize"])
	svc.AssertExpectations(t)
}

func TestJobPostingHandler_GetUsesDTOShape(t *testing.T) {
	svc := &mockJobPostingService{}
	svc.On("Get", mock.Anything, "u1", "job-1").Return(&domainjob.JobPosting{ID: "job-1", UserID: "u1", CompanyName: "Acme", JobTitle: "Backend", Status: domainjob.StatusPending}, nil)
	r := gin.New()
	h := handler.NewJobPostingHandler(svc)
	r.GET("/job-postings/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.Get(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/job-postings/job-1", nil)
	r.ServeHTTP(w, req)

	data := responseData(t, w)
	assert.Equal(t, "Acme", data["companyName"])
	assert.Equal(t, "Backend", data["jobTitle"])
	assert.NotContains(t, data, "CompanyName")
	assert.NotContains(t, data, "UserID")
	svc.AssertExpectations(t)
}

func TestResumeApplicationHandler_ListAndCreateUseDTOShape(t *testing.T) {
	svc := &mockResumeApplicationJobService{}
	appliedAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	row := domainapp.ResumeApplication{ID: "app-1", UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Backend", AppliedAt: appliedAt, CurrentStatus: domainapp.StatusApplied}
	svc.On("List", mock.Anything, "u1", "v1", "r1", "backend", string(domainapp.StatusApplied), 1, 10).Return([]domainapp.ResumeApplication{row}, int64(1), nil).Once()
	createReq := dto.CreateResumeApplicationRequest{ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Backend", AppliedAt: appliedAt}
	svc.On("Create", mock.Anything, "u1", createReq).Return(&row, nil).Once()
	r := gin.New()
	h := handler.NewResumeApplicationHandler(svc)
	r.GET("/applications", func(c *gin.Context) { c.Set("userID", "u1"); h.List(c) })
	r.POST("/applications", func(c *gin.Context) { c.Set("userID", "u1"); h.Create(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/applications?resumeVersionId=v1&repoId=r1&keyword=backend&currentStatus=applied&pageSize=200", nil)
	r.ServeHTTP(w, req)
	data := responseData(t, w)
	first := data["list"].([]any)[0].(map[string]any)
	assert.Equal(t, "v1", first["resumeVersionId"])
	assert.Equal(t, "applied", first["currentStatus"])
	assert.NotContains(t, first, "ResumeVersionID")
	assert.NotContains(t, first, "CurrentStatus")
	assert.Equal(t, float64(10), data["pageSize"])

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBufferString(`{"resumeVersionId":"v1","repoId":"r1","companyName":"Acme","jobTitle":"Backend","appliedAt":"2026-05-11T10:00:00Z"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	data = responseData(t, w)
	assert.Equal(t, "v1", data["resumeVersionId"])
	assert.Equal(t, "applied", data["currentStatus"])
	assert.NotContains(t, data, "ResumeVersionID")
	svc.AssertExpectations(t)
}

func TestInterviewProgressHandler_ListCreateUpdateUseDTOShapeAndMissingApplicationMessage(t *testing.T) {
	svc := &mockInterviewProgressJobService{}
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	row := domainprogress.InterviewProgress{ID: "p1", UserID: "u1", ApplicationID: "app-1", Stage: domainprogress.StageFirstInterview, OccurredAt: occurredAt, Result: domainprogress.ResultPending}
	svc.On("List", mock.Anything, "u1", "app-1").Return([]domainprogress.InterviewProgress{row}, nil).Once()
	createReq := dto.CreateInterviewProgressRequest{ApplicationID: "app-1", Stage: domainprogress.StageFirstInterview, OccurredAt: occurredAt, Result: domainprogress.ResultPending}
	svc.On("Create", mock.Anything, "u1", createReq).Return(&row, nil).Once()
	stage := domainprogress.StageSecondInterview
	updateReq := dto.UpdateInterviewProgressRequest{Stage: &stage}
	updated := row
	updated.Stage = stage
	svc.On("Update", mock.Anything, "u1", "p1", updateReq).Return(&updated, nil).Once()
	r := gin.New()
	h := handler.NewInterviewProgressHandler(svc)
	r.GET("/progress", func(c *gin.Context) { c.Set("userID", "u1"); h.List(c) })
	r.POST("/progress", func(c *gin.Context) { c.Set("userID", "u1"); h.Create(c) })
	r.PUT("/progress/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.Update(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/progress", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
	assert.Equal(t, "applicationId 不能为空", responseMessage(t, w))

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/progress?applicationId=app-1", nil)
	r.ServeHTTP(w, req)
	data := responseData(t, w)
	first := data["list"].([]any)[0].(map[string]any)
	assert.Equal(t, "app-1", first["applicationId"])
	assert.Equal(t, "first_interview", first["stage"])
	assert.NotContains(t, first, "ApplicationID")

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/progress", bytes.NewBufferString(`{"applicationId":"app-1","stage":"first_interview","occurredAt":"2026-05-11T10:00:00Z","result":"pending"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	data = responseData(t, w)
	assert.Equal(t, "app-1", data["applicationId"])
	assert.NotContains(t, data, "ApplicationID")

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/progress/p1", bytes.NewBufferString(`{"stage":"second_interview"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	data = responseData(t, w)
	assert.Equal(t, "second_interview", data["stage"])
	assert.NotContains(t, data, "ApplicationID")
	svc.AssertExpectations(t)
}

func responseData(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp.Data
}

func responseMessage(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp.Message
}
