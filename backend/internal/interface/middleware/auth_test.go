package middleware_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/interface/middleware"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

// --- hand-written mocks ---

type mockParser struct {
	userID   string
	username string
	jti      string
	err      error
}

func (m *mockParser) Parse(tokenStr string) (string, string, string, error) {
	return m.userID, m.username, m.jti, m.err
}

type mockBlacklist struct {
	blacklisted bool
	err         error
}

func (m *mockBlacklist) IsBlacklisted(_ context.Context, _ string) (bool, error) {
	return m.blacklisted, m.err
}

// --- helpers ---

func newRouter(parser middleware.JWTParser, bl middleware.Blacklist, downstream gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", middleware.Auth(parser, bl), downstream)
	return r
}

func doGet(r *gin.Engine, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if token != "" {
		req.Header.Set("X-Access-Token", token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func bodyCode(t *testing.T, w *httptest.ResponseRecorder) int {
	t.Helper()
	var resp struct {
		Code int `json:"code"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return resp.Code
}

// --- tests ---

func TestAuth_MissingToken(t *testing.T) {
	r := newRouter(&mockParser{}, &mockBlacklist{}, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := doGet(r, "")
	if code := bodyCode(t, w); code != pkgerrors.CodeUnauthorized {
		t.Errorf("expected code %d, got %d", pkgerrors.CodeUnauthorized, code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	parser := &mockParser{err: errors.New("bad token")}
	r := newRouter(parser, &mockBlacklist{}, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := doGet(r, "invalid.token")
	if code := bodyCode(t, w); code != pkgerrors.CodeUnauthorized {
		t.Errorf("expected code %d, got %d", pkgerrors.CodeUnauthorized, code)
	}
}

func TestAuth_BlacklistedToken(t *testing.T) {
	parser := &mockParser{userID: "1", username: "alice", jti: "jti-1"}
	bl := &mockBlacklist{blacklisted: true}
	r := newRouter(parser, bl, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := doGet(r, "some.valid.token")
	if code := bodyCode(t, w); code != pkgerrors.CodeUnauthorized {
		t.Errorf("expected code %d, got %d", pkgerrors.CodeUnauthorized, code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	parser := &mockParser{userID: "42", username: "bob", jti: "jti-2"}
	bl := &mockBlacklist{blacklisted: false}
	downstream := func(c *gin.Context) {
		uid, _ := c.Get("userID")
		uname, _ := c.Get("username")
		c.Header("X-UserID", uid.(string))
		c.Header("X-Username", uname.(string))
		c.Status(http.StatusOK)
	}
	r := newRouter(parser, bl, downstream)
	w := doGet(r, "good.token")

	if uid := w.Header().Get("X-UserID"); uid != "42" {
		t.Errorf("expected userID '42', got %q", uid)
	}
	if uname := w.Header().Get("X-Username"); uname != "bob" {
		t.Errorf("expected username 'bob', got %q", uname)
	}
}
