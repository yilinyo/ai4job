package repository

import (
	"context"

	"github.com/yilin/ai-for-backend/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
