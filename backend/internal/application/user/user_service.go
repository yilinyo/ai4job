package user

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

// TokenCache defines the cache operations needed by the user service.
type TokenCache interface {
	SetEmailCode(ctx context.Context, email, code string) error
	VerifyEmailCode(ctx context.Context, email, code string) (bool, error)
	DeleteEmailCode(ctx context.Context, email string) error
	AddToBlacklist(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

// JWTManager defines the JWT operations needed by the user service.
type JWTManager interface {
	Generate(userID, username string) (string, error)
}

type UserService struct {
	repo      domainuser.Repository
	cache     TokenCache
	jwt       JWTManager
	emailMock bool
}

func NewUserService(repo domainuser.Repository, cache TokenCache, jwt JWTManager, emailMock bool) *UserService {
	return &UserService{repo: repo, cache: cache, jwt: jwt, emailMock: emailMock}
}

func (s *UserService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {
	code := fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
	if err := s.cache.SetEmailCode(ctx, req.Email, code); err != nil {
		return nil, err
	}
	resp := &dto.SendEmailCodeResponse{Email: req.Email, ExpiresIn: 300}
	if s.emailMock {
		resp.MockCode = code
	}
	return resp, nil
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) error {
	ok, err := s.cache.VerifyEmailCode(ctx, req.Email, req.EmailCode)
	if err != nil {
		return err
	}
	if !ok {
		return domainerrors.ErrEmailCodeWrong
	}
	existsEmail, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existsEmail {
		return domainerrors.ErrDuplicate
	}
	existsUsername, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existsUsername {
		return domainerrors.ErrDuplicate
	}
	u := &domainuser.User{
		ID:                   uuid.NewString(),
		Username:             req.Username,
		Email:                req.Email,
		EducationExperiences: []domainuser.EducationExperience{},
	}
	if err := u.SetPassword(req.Password); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}
	return s.cache.DeleteEmailCode(ctx, req.Email)
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, domainerrors.ErrUnauthorized
	}
	if !u.CheckPassword(req.Password) {
		return nil, domainerrors.ErrUnauthorized
	}
	token, err := s.jwt.Generate(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Token: token, UserInfo: dto.ToUserInfoResponse(u)}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToUserInfoResponse(u), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.RealName = req.RealName
	u.Age = req.Age
	u.Phone = req.Phone
	u.JobIntention = req.JobIntention
	u.Avatar = req.Avatar
	u.Location = req.Location
	u.PersonalAdvantage = req.PersonalAdvantage
	u.EducationExperiences = req.EducationExperiences
	u.Normalize()
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return dto.ToUserInfoResponse(u), nil
}
