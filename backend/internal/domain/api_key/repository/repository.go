package repository

import (
	"context"
	apikey "github.com/yilin/ai-for-backend/internal/domain/api_key"
	"time"
)

type Repository interface {
	Create(context.Context, *apikey.APIKey) error
	ListByUser(context.Context, string) ([]apikey.APIKey, error)
	FindByIDAndUser(context.Context, string, string) (*apikey.APIKey, error)
	FindByDigest(context.Context, string) (*apikey.APIKey, error)
	Save(context.Context, *apikey.APIKey) error
	Touch(context.Context, *apikey.APIKey, time.Time) error
}
