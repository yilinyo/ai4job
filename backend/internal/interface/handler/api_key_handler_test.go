package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockAPIKeyService struct{ mock.Mock }

func (m *mockAPIKeyService) List(ctx context.Context, userID string) ([]dto.APIKeyResponse, error) {
	args := m.Called(ctx, userID)
	rows, _ := args.Get(0).([]dto.APIKeyResponse)
	return rows, args.Error(1)
}
func (m *mockAPIKeyService) Create(ctx context.Context, userID string, req dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error) {
	args := m.Called(ctx, userID, req)
	row, _ := args.Get(0).(*dto.CreatedAPIKeyResponse)
	return row, args.Error(1)
}
func (m *mockAPIKeyService) Revoke(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}
func (m *mockAPIKeyService) Rotate(ctx context.Context, userID, id string, req dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error) {
	args := m.Called(ctx, userID, id, req)
	row, _ := args.Get(0).(*dto.CreatedAPIKeyResponse)
	return row, args.Error(1)
}

func newAPIKeyRouter(service *mockAPIKeyService) *gin.Engine {
	r := gin.New()
	h := handler.NewAPIKeyHandler(service)
	r.GET("/api-keys", func(c *gin.Context) { c.Set("userID", "user-1"); h.List(c) })
	r.POST("/api-keys", func(c *gin.Context) { c.Set("userID", "user-1"); h.Create(c) })
	r.DELETE("/api-keys/:id", func(c *gin.Context) { c.Set("userID", "user-1"); h.Revoke(c) })
	r.POST("/api-keys/:id/rotate", func(c *gin.Context) { c.Set("userID", "user-1"); h.Rotate(c) })
	return r
}

func TestAPIKeyHandlerReturnsPlaintextOnlyForCreateAndRotate(t *testing.T) {
	service := &mockAPIKeyService{}
	createdAt := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	row := dto.APIKeyResponse{ID: "key-1", Name: "Codex", DisplayPrefix: "ai4job-AbCd...", Status: "active", CreatedAt: createdAt}
	service.On("List", mock.Anything, "user-1").Return([]dto.APIKeyResponse{row}, nil).Once()
	service.On("Create", mock.Anything, "user-1", dto.CreateAPIKeyRequest{Name: "Codex", ExpiresInDays: nil}).Return(&dto.CreatedAPIKeyResponse{APIKeyResponse: &row, Key: "ai4job-secret"}, nil).Once()
	service.On("Rotate", mock.Anything, "user-1", "key-1", dto.CreateAPIKeyRequest{Name: "Codex", ExpiresInDays: nil}).Return(&dto.CreatedAPIKeyResponse{APIKeyResponse: &row, Key: "ai4job-rotated"}, nil).Once()
	r := newAPIKeyRouter(service)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api-keys", nil))
	list := responseMap(t, w)["data"].(map[string]any)["list"].([]any)[0].(map[string]any)
	require.NotContains(t, list, "key")
	require.Equal(t, "ai4job-AbCd...", list["displayPrefix"])

	for _, tc := range []struct{ path, want string }{{"/api-keys", "ai4job-secret"}, {"/api-keys/key-1/rotate", "ai4job-rotated"}} {
		w = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBufferString(`{"name":"Codex","expiresInDays":null}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		data := responseMap(t, w)["data"].(map[string]any)
		require.Equal(t, tc.want, data["key"])
	}
	service.AssertExpectations(t)
}

func TestAPIKeyHandlerDoesNotAllowCrossUserRevocation(t *testing.T) {
	service := &mockAPIKeyService{}
	service.On("Revoke", mock.Anything, "user-1", "other-users-key").Return(pkgerrors.ErrNotFound).Once()
	r := newAPIKeyRouter(service)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api-keys/other-users-key", nil))
	require.Equal(t, pkgerrors.CodeNotFound, responseCode(t, w))
	service.AssertExpectations(t)
}
