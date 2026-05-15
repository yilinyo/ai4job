package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeRepoService struct{ mock.Mock }

func (m *mockResumeRepoService) Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error {
	return m.Called(ctx, userID, req).Error(0)
}

func (m *mockResumeRepoService) List(ctx context.Context, userID string, jobType domainrepo.JobType, keyword string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
	args := m.Called(ctx, userID, jobType, keyword, page, pageSize)
	return args.Get(0).([]domainrepo.ResumeRepo), args.Get(1).(int64), args.Error(2)
}

func (m *mockResumeRepoService) Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainrepo.ResumeRepo), args.Error(1)
}

func (m *mockResumeRepoService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error {
	return m.Called(ctx, userID, id, req).Error(0)
}

func (m *mockResumeRepoService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func TestResumeRepoHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeRepoHandler(nil)
	r.POST("/api/v1/resume-repos", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
}

func TestResumeRepoHandler_Create_Success(t *testing.T) {
	svc := &mockResumeRepoService{}
	svc.On("Create", mock.Anything, "u1", dto.CreateResumeRepoRequest{Name: "Main", JobType: domainrepo.JobTypeCampus, TargetPosition: "Frontend Engineer", Description: "desc"}).Return(nil)
	r := gin.New()
	h := handler.NewResumeRepoHandler(svc)
	r.POST("/api/v1/resume-repos", func(c *gin.Context) { c.Set("userID", "u1"); h.Create(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos", bytes.NewBufferString(`{"name":"Main","jobType":"campus","targetPosition":"Frontend Engineer","description":"desc"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestResumeRepoHandler_List_Success(t *testing.T) {
	svc := &mockResumeRepoService{}
	rows := []domainrepo.ResumeRepo{{ID: "r1", UserID: "u1", Name: "Main", JobType: domainrepo.JobTypeCampus, TargetPosition: "Frontend Engineer", VersionCount: 2}}
	svc.On("List", mock.Anything, "u1", domainrepo.JobTypeCampus, "front", 1, 10).Return(rows, int64(1), nil)
	r := gin.New()
	h := handler.NewResumeRepoHandler(svc)
	r.GET("/api/v1/resume-repos", func(c *gin.Context) { c.Set("userID", "u1"); h.List(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-repos?jobType=campus&keyword=front", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	body := w.Body.String()
	assert.Contains(t, body, `"jobType":"campus"`)
	assert.Contains(t, body, `"targetPosition":"Frontend Engineer"`)
	assert.Contains(t, body, `"versionCount":2`)
	assert.NotContains(t, body, "JobType")
	assert.NotContains(t, body, "TargetPosition")
	assert.NotContains(t, body, "VersionCount")
	svc.AssertExpectations(t)
}

func TestResumeRepoHandler_Get_Success(t *testing.T) {
	svc := &mockResumeRepoService{}
	row := &domainrepo.ResumeRepo{ID: "r1", UserID: "u1", Name: "Main", JobType: domainrepo.JobTypeCampus, TargetPosition: "Frontend Engineer", VersionCount: 2}
	svc.On("Get", mock.Anything, "u1", "r1").Return(row, nil)
	r := gin.New()
	h := handler.NewResumeRepoHandler(svc)
	r.GET("/api/v1/resume-repos/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.Get(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-repos/r1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	body := w.Body.String()
	assert.Contains(t, body, `"jobType":"campus"`)
	assert.Contains(t, body, `"targetPosition":"Frontend Engineer"`)
	assert.Contains(t, body, `"versionCount":2`)
	assert.NotContains(t, body, "JobType")
	assert.NotContains(t, body, "TargetPosition")
	assert.NotContains(t, body, "VersionCount")
	svc.AssertExpectations(t)
}
