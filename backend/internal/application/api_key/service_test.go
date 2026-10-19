package apikey

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domain "github.com/yilin/ai-for-backend/internal/domain/api_key"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type memoryRepo struct{ keys map[string]*domain.APIKey }

func newMemoryRepo() *memoryRepo { return &memoryRepo{keys: map[string]*domain.APIKey{}} }
func (r *memoryRepo) Create(_ context.Context, key *domain.APIKey) error {
	if _, exists := r.keys[key.ID]; exists {
		return errors.New("duplicate")
	}
	key.CreatedAt = time.Now().UTC()
	r.keys[key.ID] = key
	return nil
}
func (r *memoryRepo) ListByUser(_ context.Context, userID string) ([]domain.APIKey, error) {
	keys := []domain.APIKey{}
	for _, key := range r.keys {
		if key.UserID == userID {
			keys = append(keys, *key)
		}
	}
	return keys, nil
}
func (r *memoryRepo) FindByIDAndUser(_ context.Context, id, userID string) (*domain.APIKey, error) {
	key, ok := r.keys[id]
	if !ok || key.UserID != userID {
		return nil, domainerrors.ErrNotFound
	}
	return key, nil
}
func (r *memoryRepo) FindByDigest(_ context.Context, digest string) (*domain.APIKey, error) {
	for _, key := range r.keys {
		if key.KeyDigest == digest {
			return key, nil
		}
	}
	return nil, domainerrors.ErrNotFound
}
func (r *memoryRepo) Save(_ context.Context, key *domain.APIKey) error {
	r.keys[key.ID] = key
	return nil
}
func (r *memoryRepo) Touch(_ context.Context, key *domain.APIKey, now time.Time) error {
	key.LastUsedAt = &now
	return nil
}

type immediateTransactionRunner struct{}

func (immediateTransactionRunner) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestServiceCreateListRevokeAndAuthenticate(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, "test-pepper", immediateTransactionRunner{})
	days := 30
	created, err := service.Create(context.Background(), "user-1", dto.CreateAPIKeyRequest{Name: "Codex", ExpiresInDays: &days})
	require.NoError(t, err)
	require.NotEmpty(t, created.Key)
	require.NotEqual(t, created.Key, repo.keys[created.ID].KeyDigest)

	listed, err := service.List(context.Background(), "user-1")
	require.NoError(t, err)
	require.Len(t, listed, 1)
	require.NotEmpty(t, listed[0].DisplayPrefix)

	authenticated, err := service.Authenticate(context.Background(), created.Key)
	require.NoError(t, err)
	require.Equal(t, "user-1", authenticated.UserID)
	require.NotNil(t, authenticated.LastUsedAt)

	require.NoError(t, service.Revoke(context.Background(), "user-1", created.ID))
	_, err = service.Authenticate(context.Background(), created.Key)
	require.ErrorIs(t, err, domainerrors.ErrUnauthorized)

	require.ErrorIs(t, service.Revoke(context.Background(), "another-user", created.ID), domainerrors.ErrNotFound)
}

func TestServiceRotateOnlyReturnsNewPlaintext(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, "test-pepper", immediateTransactionRunner{})
	first, err := service.Create(context.Background(), "user-1", dto.CreateAPIKeyRequest{Name: "old"})
	require.NoError(t, err)
	rotated, err := service.Rotate(context.Background(), "user-1", first.ID, dto.CreateAPIKeyRequest{Name: "new"})
	require.NoError(t, err)
	require.NotEqual(t, first.Key, rotated.Key)
	require.Equal(t, domain.StatusRevoked, repo.keys[first.ID].Status)
	_, err = service.Authenticate(context.Background(), first.Key)
	require.ErrorIs(t, err, domainerrors.ErrUnauthorized)
	_, err = service.Authenticate(context.Background(), rotated.Key)
	require.NoError(t, err)
}
