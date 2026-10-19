package apikey

import (
	"context"
	"time"

	"github.com/yilin/ai-for-backend/internal/application/dto"
	domain "github.com/yilin/ai-for-backend/internal/domain/api_key"
	repo "github.com/yilin/ai-for-backend/internal/domain/api_key/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type Service struct {
	repo         repo.Repository
	pepper       string
	transactions TransactionRunner
}

type TransactionRunner interface {
	WithinTransaction(context.Context, func(context.Context) error) error
}

func NewService(repo repo.Repository, pepper string, transactions TransactionRunner) *Service {
	return &Service{repo: repo, pepper: pepper, transactions: transactions}
}
func (s *Service) List(ctx context.Context, userID string) ([]dto.APIKeyResponse, error) {
	keys, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.APIKeyResponse, 0, len(keys))
	for i := range keys {
		out = append(out, response(&keys[i]))
	}
	return out, nil
}
func (s *Service) Create(ctx context.Context, userID string, req dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error) {
	key, plain, err := domain.New(userID, req.Name, s.pepper, req.ExpiresInDays)
	if err != nil {
		return nil, domainerrors.ErrBadParams
	}
	if err = s.repo.Create(ctx, key); err != nil {
		return nil, err
	}
	base := response(key)
	return &dto.CreatedAPIKeyResponse{APIKeyResponse: &base, Key: plain}, nil
}
func (s *Service) Revoke(ctx context.Context, userID, id string) error {
	key, err := s.repo.FindByIDAndUser(ctx, id, userID)
	if err != nil {
		return err
	}
	if key.Status == domain.StatusRevoked {
		return nil
	}
	key.Revoke(time.Now().UTC())
	return s.repo.Save(ctx, key)
}
func (s *Service) Rotate(ctx context.Context, userID, id string, req dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error) {
	if s.transactions == nil {
		return nil, domainerrors.ErrBadParams
	}
	var result *dto.CreatedAPIKeyResponse
	err := s.transactions.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.Revoke(txCtx, userID, id); err != nil {
			return err
		}
		created, err := s.Create(txCtx, userID, req)
		if err != nil {
			return err
		}
		result = created
		return nil
	})
	return result, err
}
func (s *Service) Authenticate(ctx context.Context, plain string) (*domain.APIKey, error) {
	if len(plain) < len(domain.Prefix)+20 || plain[:len(domain.Prefix)] != domain.Prefix {
		return nil, domainerrors.ErrUnauthorized
	}
	key, err := s.repo.FindByDigest(ctx, domain.Digest(plain, s.pepper))
	if err != nil || !key.Matches(plain, s.pepper) || !key.ValidAt(time.Now().UTC()) {
		return nil, domainerrors.ErrUnauthorized
	}
	if key.LastUsedAt == nil || time.Since(*key.LastUsedAt) > time.Minute {
		_ = s.repo.Touch(ctx, key, time.Now().UTC())
	}
	return key, nil
}
func response(key *domain.APIKey) dto.APIKeyResponse {
	return dto.APIKeyResponse{ID: key.ID, Name: key.Name, DisplayPrefix: key.DisplayPrefix, Status: string(key.Status), CreatedAt: key.CreatedAt, ExpiresAt: key.ExpiresAt, LastUsedAt: key.LastUsedAt}
}
