package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appuser "github.com/yilin/ai-for-backend/internal/application/user"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

// --- hand-written mocks ---

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, u *domainuser.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, u *domainuser.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
func (m *mockUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

type mockTokenCache struct{ mock.Mock }

func (m *mockTokenCache) SetEmailCode(ctx context.Context, email, code string) error {
	return m.Called(ctx, email, code).Error(0)
}
func (m *mockTokenCache) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	args := m.Called(ctx, email, code)
	return args.Bool(0), args.Error(1)
}
func (m *mockTokenCache) DeleteEmailCode(ctx context.Context, email string) error {
	return m.Called(ctx, email).Error(0)
}
func (m *mockTokenCache) AddToBlacklist(ctx context.Context, jti string, ttl time.Duration) error {
	return m.Called(ctx, jti, ttl).Error(0)
}
func (m *mockTokenCache) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

type mockJWT struct{ mock.Mock }

func (m *mockJWT) Generate(userID, username string) (string, error) {
	args := m.Called(userID, username)
	return args.String(0), args.Error(1)
}

// --- tests ---

func TestUserService_Login_Success(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	existingUser := &domainuser.User{ID: "uid-1", Username: "alice", Email: "a@b.com"}
	_ = existingUser.SetPassword("password123")

	repo.On("FindByUsername", mock.Anything, "alice").Return(existingUser, nil)
	jwtMgr.On("Generate", "uid-1", "alice").Return("signed.jwt.token", nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	resp, err := svc.Login(context.Background(), dto.LoginRequest{
		Username: "alice",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "signed.jwt.token", resp.Token)
	assert.Equal(t, "alice", resp.UserInfo.Username)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	existingUser := &domainuser.User{ID: "uid-1", Username: "alice"}
	_ = existingUser.SetPassword("correct")

	repo.On("FindByUsername", mock.Anything, "alice").Return(existingUser, nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Username: "alice",
		Password: "wrong",
	})

	assert.ErrorIs(t, err, domainerrors.ErrUnauthorized)
}

func TestUserService_Register_InvalidEmailCode(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	cache.On("VerifyEmailCode", mock.Anything, "test@example.com", "000000").Return(false, nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	err := svc.Register(context.Background(), dto.RegisterRequest{
		Username:  "alice",
		Password:  "pass1234",
		Email:     "test@example.com",
		EmailCode: "000000",
	})

	assert.ErrorIs(t, err, domainerrors.ErrEmailCodeWrong)
}

func TestUserService_Register_InitializesEducationExperiences(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	cache.On("VerifyEmailCode", mock.Anything, "test@example.com", "123456").Return(true, nil)
	repo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	repo.On("ExistsByUsername", mock.Anything, "alice").Return(false, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(u *domainuser.User) bool {
		return u.Username == "alice" &&
			u.Email == "test@example.com" &&
			u.PasswordHash != "" &&
			u.EducationExperiences != nil &&
			len(u.EducationExperiences) == 0
	})).Return(nil)
	cache.On("DeleteEmailCode", mock.Anything, "test@example.com").Return(nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	err := svc.Register(context.Background(), dto.RegisterRequest{
		Username:  "alice",
		Password:  "pass1234",
		Email:     "test@example.com",
		EmailCode: "123456",
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestUserService_UpdateProfile_NormalizesNilEducationExperiences(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	existingUser := &domainuser.User{ID: "uid-1", Username: "alice", Email: "test@example.com"}

	repo.On("FindByID", mock.Anything, "uid-1").Return(existingUser, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domainuser.User) bool {
		return u.ID == "uid-1" &&
			u.RealName == "Alice" &&
			u.EducationExperiences != nil &&
			len(u.EducationExperiences) == 0
	})).Return(nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	resp, err := svc.UpdateProfile(context.Background(), "uid-1", dto.UpdateProfileRequest{
		RealName: "Alice",
	})

	require.NoError(t, err)
	assert.NotNil(t, resp.EducationExperiences)
	assert.Empty(t, resp.EducationExperiences)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_PreservesNonEmptyEducationExperiences(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	exps := []domainuser.EducationExperience{
		{School: "Test University", Education: "本科", Major: "Computer Science"},
	}
	existingUser := &domainuser.User{ID: "uid-1", Username: "alice", Email: "test@example.com"}

	repo.On("FindByID", mock.Anything, "uid-1").Return(existingUser, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domainuser.User) bool {
		return assert.ObjectsAreEqual(exps, u.EducationExperiences)
	})).Return(nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	resp, err := svc.UpdateProfile(context.Background(), "uid-1", dto.UpdateProfileRequest{
		EducationExperiences: exps,
	})

	require.NoError(t, err)
	assert.Equal(t, exps, resp.EducationExperiences)
	repo.AssertExpectations(t)
}
