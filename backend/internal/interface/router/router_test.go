package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/application/dto"
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
func (s *stubResumeRepoService) List(ctx context.Context, userID string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
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

func (s *stubResumeVersionService) Create(ctx context.Context, userID, repoID string, req dto.CreateResumeVersionRequest) error {
	return nil
}
func (s *stubResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	s.lastRepoID = repoID
	return []domainversion.ResumeVersion{}, nil
}
func (s *stubResumeVersionService) Get(ctx context.Context, userID, repoID, id string) (*domainversion.ResumeVersion, error) {
	s.lastRepoID = repoID
	s.lastVersionID = id
	return &domainversion.ResumeVersion{ID: id, RepoID: repoID, Title: "v1", Content: "{}", VersionNum: 1}, nil
}
func (s *stubResumeVersionService) Update(ctx context.Context, userID, repoID, id string, req dto.UpdateResumeVersionRequest) error {
	return nil
}
func (s *stubResumeVersionService) Delete(ctx context.Context, userID, repoID, id string) error {
	return nil
}
func (s *stubResumeVersionService) SetDefault(ctx context.Context, userID, repoID, id string) error {
	return nil
}

func TestRouter_ResumeRoutesRequireAuth(t *testing.T) {
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

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/resume-repos"},
		{http.MethodGet, "/api/v1/resume-repos/r1"},
		{http.MethodGet, "/api/v1/resume-repos/r1/versions"},
		{http.MethodGet, "/api/v1/resume-repos/r1/versions/v1"},
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
