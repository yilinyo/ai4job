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

func (m *mockResumeVersionService) Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error {
	return m.Called(ctx, userID, req).Error(0)
}

func (m *mockResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, repoID)
	return args.Get(0).([]domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error {
	return m.Called(ctx, userID, id, req).Error(0)
}

func (m *mockResumeVersionService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *mockResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func TestResumeVersionHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeVersionHandler(nil)
	r.POST("/api/v1/resume-versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-versions", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
}

func TestResumeVersionHandler_Create_MissingUser(t *testing.T) {
	svc := &mockResumeVersionService{}
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-versions", bytes.NewBufferString(`{"repoId":"r1","title":"v1","content":"{}"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeUnauthorized, responseCode(t, w))
	svc.AssertNotCalled(t, "Create")
}

func TestResumeVersionHandler_Create_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	svc.On("Create", mock.Anything, "u1", dto.CreateResumeVersionRequest{RepoID: "r1", Title: "v1", Content: "{}", Remark: "first"}).Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-versions", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Create(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-versions", bytes.NewBufferString(`{"repoId":"r1","title":"v1","content":"{}","remark":"first"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_List_MissingRepoID(t *testing.T) {
	svc := &mockResumeVersionService{}
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-versions", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.List(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
	svc.AssertNotCalled(t, "List")
}

func TestResumeVersionHandler_List_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	rows := []domainversion.ResumeVersion{{ID: "v1", RepoID: "r1", Version: "v1.0.0", Title: "Main", Content: "{}", Remark: "first", IsActive: true, CreatedAt: now, UpdatedAt: now}}
	svc.On("List", mock.Anything, "u1", "r1").Return(rows, nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-versions", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.List(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions?repoId=r1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	body := w.Body.String()
	assert.Contains(t, body, `"list":[`)
	assert.Contains(t, body, `"total":1`)
	assert.Contains(t, body, `"page":1`)
	assert.Contains(t, body, `"pageSize":10`)
	assert.Contains(t, body, `"version":"v1.0.0"`)
	assert.Contains(t, body, `"remark":"first"`)
	assert.Contains(t, body, `"isActive":true`)
	assert.NotContains(t, body, "VersionNum")
	assert.NotContains(t, body, "IsDefault")
	assert.NotContains(t, body, "Version:")
	assert.NotContains(t, body, "Remark:")
	assert.NotContains(t, body, "IsActive")
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_Compare_NotImplemented(t *testing.T) {
	svc := &mockResumeVersionService{}
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-versions/compare", h.Compare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions/compare", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
	assert.Contains(t, w.Body.String(), "暂不支持版本对比")
	svc.AssertNotCalled(t, "Get")
}

func TestResumeVersionHandler_Get_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	row := &domainversion.ResumeVersion{ID: "v1", RepoID: "r1", Version: "v1.0.0", Title: "Main", Content: "{}", Remark: "first", IsActive: true, CreatedAt: now, UpdatedAt: now}
	svc.On("Get", mock.Anything, "u1", "v1").Return(row, nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.GET("/api/v1/resume-versions/:id", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Get(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions/v1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	body := w.Body.String()
	assert.Contains(t, body, `"version":"v1.0.0"`)
	assert.Contains(t, body, `"remark":"first"`)
	assert.Contains(t, body, `"isActive":true`)
	assert.NotContains(t, body, "VersionNum")
	assert.NotContains(t, body, "IsDefault")
	assert.NotContains(t, body, "Version:")
	assert.NotContains(t, body, "Remark:")
	assert.NotContains(t, body, "IsActive")
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_Update_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	active := true
	svc.On("Update", mock.Anything, "u1", "v1", dto.UpdateResumeVersionRequest{Title: "v2", Content: "{}", Remark: "second", IsActive: &active}).Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.PUT("/api/v1/resume-versions/:id", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Update(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/resume-versions/v1", bytes.NewBufferString(`{"title":"v2","content":"{}","remark":"second","isActive":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_Delete_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	svc.On("Delete", mock.Anything, "u1", "v1").Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.DELETE("/api/v1/resume-versions/:id", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Delete(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/resume-versions/v1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestResumeVersionHandler_Activate_Success(t *testing.T) {
	svc := &mockResumeVersionService{}
	svc.On("Activate", mock.Anything, "u1", "v1").Return(nil)
	r := gin.New()
	h := handler.NewResumeVersionHandler(svc)
	r.POST("/api/v1/resume-versions/:id/activate", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.Activate(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-versions/v1/activate", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}
