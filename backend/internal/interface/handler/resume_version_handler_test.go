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
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeVersionService struct{ mock.Mock }

func (m *mockResumeVersionService) Create(ctx context.Context, repoID string, req dto.CreateResumeVersionRequest) error {
	return m.Called(ctx, repoID, req).Error(0)
}

func (m *mockResumeVersionService) List(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) {
	args := m.Called(ctx, repoID)
	return args.Get(0).([]domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Get(ctx context.Context, repoID, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, repoID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Update(ctx context.Context, repoID, id string, req dto.UpdateResumeVersionRequest) error {
	return m.Called(ctx, repoID, id, req).Error(0)
}

func (m *mockResumeVersionService) Delete(ctx context.Context, repoID, id string) error {
	return m.Called(ctx, repoID, id).Error(0)
}

func (m *mockResumeVersionService) SetDefault(ctx context.Context, repoID, id string) error {
	return m.Called(ctx, repoID, id).Error(0)
}

func TestResumeVersionHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeVersionHandler(nil)
	r.POST("/api/v1/resume-repos/:repoId/versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos/r1/versions", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
}

func TestResumeVersionHandler_Create_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	svc.On("Create", mock.Anything, "r1", dto.CreateResumeVersionRequest{Title: "v1", Content: "{}"}).Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-repos/:repoId/versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos/r1/versions", bytes.NewBufferString(`{"title":"v1","content":"{}"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}
