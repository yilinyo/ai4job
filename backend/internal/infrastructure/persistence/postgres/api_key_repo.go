package postgres

import (
	"context"
	"errors"
	"time"

	apikey "github.com/yilin/ai-for-backend/internal/domain/api_key"
	repo "github.com/yilin/ai-for-backend/internal/domain/api_key/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type APIKeyRepo struct{ db *gorm.DB }

var _ repo.Repository = (*APIKeyRepo)(nil)

func NewAPIKeyRepo(db *gorm.DB) *APIKeyRepo { return &APIKeyRepo{db: db} }
func (r *APIKeyRepo) Create(ctx context.Context, key *apikey.APIKey) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(key).Error
}
func (r *APIKeyRepo) ListByUser(ctx context.Context, userID string) ([]apikey.APIKey, error) {
	var keys []apikey.APIKey
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error
	return keys, err
}
func (r *APIKeyRepo) FindByIDAndUser(ctx context.Context, id, userID string) (*apikey.APIKey, error) {
	var key apikey.APIKey
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &key, err
}
func (r *APIKeyRepo) FindByDigest(ctx context.Context, digest string) (*apikey.APIKey, error) {
	var key apikey.APIKey
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("key_digest = ?", digest).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &key, err
}
func (r *APIKeyRepo) Save(ctx context.Context, key *apikey.APIKey) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Save(key).Error
}
func (r *APIKeyRepo) Touch(ctx context.Context, key *apikey.APIKey, now time.Time) error {
	key.LastUsedAt = &now
	return r.Save(ctx, key)
}
