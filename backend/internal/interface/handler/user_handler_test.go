package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- mock ---

type mockUserService struct{ mock.Mock }

func (m *mockUserService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SendEmailCodeResponse), args.Error(1)
}
func (m *mockUserService) Register(ctx context.Context, req dto.RegisterRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *mockUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}
func (m *mockUserService) GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserInfoResponse), args.Error(1)
}
func (m *mockUserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserInfoResponse), args.Error(1)
}

// --- helpers ---

func responseCode(t *testing.T, w *httptest.ResponseRecorder) int {
	t.Helper()
	var resp struct {
		Code int `json:"code"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp.Code
}

func newRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.POST("/email-code", h.SendEmailCode)
	r.POST("/login", h.Login)
	r.GET("/profile", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.GetProfile(c)
	})
	return r
}

// --- tests ---

func TestUserHandler_SendEmailCode_Success(t *testing.T) {
	svc := &mockUserService{}
	svc.On("SendEmailCode", mock.Anything, dto.SendEmailCodeRequest{Email: "a@b.com"}).
		Return(&dto.SendEmailCodeResponse{Email: "a@b.com", MockCode: "123456", ExpiresIn: 300}, nil)

	h := handler.NewUserHandler(svc)
	r := newRouter(h)

	body := `{"email":"a@b.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/email-code", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestUserHandler_SendEmailCode_BadBody(t *testing.T) {
	svc := &mockUserService{}
	h := handler.NewUserHandler(svc)
	r := newRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/email-code", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
}

func TestUserHandler_Login_Success(t *testing.T) {
	svc := &mockUserService{}
	svc.On("Login", mock.Anything, dto.LoginRequest{Username: "alice", Password: "pass123"}).
		Return(&dto.LoginResponse{Token: "tok"}, nil)

	h := handler.NewUserHandler(svc)
	r := newRouter(h)

	body := `{"username":"alice","password":"pass123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}

func TestUserHandler_Login_Unauthorized(t *testing.T) {
	svc := &mockUserService{}
	svc.On("Login", mock.Anything, mock.Anything).
		Return(nil, pkgerrors.ErrUnauthorized)

	h := handler.NewUserHandler(svc)
	r := newRouter(h)

	body := `{"username":"alice","password":"wrong"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeUnauthorized, responseCode(t, w))
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	svc := &mockUserService{}
	svc.On("GetProfile", mock.Anything, "u1").
		Return(&dto.UserInfoResponse{ID: "u1", Username: "alice"}, nil)

	h := handler.NewUserHandler(svc)
	r := newRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	svc.AssertExpectations(t)
}
