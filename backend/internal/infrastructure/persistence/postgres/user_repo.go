package postgres

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domainuser.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

// isUniqueViolation detects PostgreSQL unique constraint violation (error code 23505).
func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) Update(ctx context.Context, u *domainuser.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainuser.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainuser.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
