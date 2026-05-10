package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeVersionService struct{ mock.Mock }

func (m *mockResumeVersionService) Create(ctx context.Context, userID, repoID string, req dto.CreateResumeVersionRequest) error {
	return m.Called(ctx, userID, repoID, req).Error(0)
}

func (m *mockResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, repoID)
	return args.Get(0).([]domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Get(ctx context.Context, userID, repoID, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, repoID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Update(ctx context.Context, userID, repoID, id string, req dto.UpdateResumeVersionRequest) error {
	return m.Called(ctx, userID, repoID, id, req).Error(0)
}

func (m *mockResumeVersionService) Delete(ctx context.Context, userID, repoID, id string) error {
	return m.Called(ctx, userID, repoID, id).Error(0)
}

func (m *mockResumeVersionService) SetDefault(ctx context.Context, userID, repoID, id string) error {
	return m.Called(ctx, userID, repoID, id).Error(0)
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

func TestResumeVersionHandler_Create_MissingUser(t *testing.T) {
	svc := &mockResumeVersionService{}
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-repos/:repoId/versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos/r1/versions", bytes.NewBufferString(`{"title":"v1","content":"{}"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeUnauthorized, responseCode(t, w))
	svc.AssertNotCalled(t, "Create")
}

func TestResumeVersionHandler_Create_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	svc.On("Create", mock.Anything, "u1", "r1", dto.CreateResumeVersionRequest{Title: "v1", Content: "{}"}).Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-repos/:repoId/versions", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Create(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos/r1/versions", bytes.NewBufferString(`{"title":"v1","content":"{}"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_List_ReturnsDTOFields(t *testing.T) {
	now := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	svc := &mockResumeVersionService{}
	svc.On("List", mock.Anything, "u1", "r1").Return([]domainversion.ResumeVersion{
		{ID: "v1", RepoID: "r1", Title: "Main", Content: "{}", VersionNum: 1, IsDefault: true, CreatedAt: now, UpdatedAt: now},
	}, nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-repos/:repoId/versions", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.List(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-repos/r1/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	resp := responseMap(t, w)
	list := resp["data"].([]any)
	version := list[0].(map[string]any)
	assert.Equal(t, "v1", version["id"])
	assert.Equal(t, "r1", version["repoId"])
	assert.Equal(t, "Main", version["title"])
	assert.Equal(t, "{}", version["content"])
	assert.Equal(t, float64(1), version["versionNum"])
	assert.Equal(t, true, version["isDefault"])
	assert.Contains(t, version, "createdAt")
	assert.Contains(t, version, "updatedAt")
	assertNoJSONKeys(t, version, "ID", "RepoID", "Title", "Content", "VersionNum", "IsDefault", "CreatedAt", "UpdatedAt", "DeletedAt")
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_Get_ReturnsDTOFields(t *testing.T) {
	now := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	svc := &mockResumeVersionService{}
	svc.On("Get", mock.Anything, "u1", "r1", "v1").Return(
		&domainversion.ResumeVersion{ID: "v1", RepoID: "r1", Title: "Main", Content: "{}", VersionNum: 1, IsDefault: true, CreatedAt: now, UpdatedAt: now},
		nil,
	)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-repos/:repoId/versions/:versionId", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Get(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-repos/r1/versions/v1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	resp := responseMap(t, w)
	version := resp["data"].(map[string]any)
	assert.Equal(t, "v1", version["id"])
	assert.Equal(t, "r1", version["repoId"])
	assert.Equal(t, "Main", version["title"])
	assert.Equal(t, "{}", version["content"])
	assert.Equal(t, float64(1), version["versionNum"])
	assert.Equal(t, true, version["isDefault"])
	assert.Contains(t, version, "createdAt")
	assert.Contains(t, version, "updatedAt")
	assertNoJSONKeys(t, version, "ID", "RepoID", "Title", "Content", "VersionNum", "IsDefault", "CreatedAt", "UpdatedAt", "DeletedAt")
	svc.AssertExpectations(t)
}
