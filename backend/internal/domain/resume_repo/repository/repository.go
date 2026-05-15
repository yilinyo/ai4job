package repository

import (
	"context"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
)

type ListFilter struct {
	UserID   string
	JobType  domainrepo.JobType
	Keyword  string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, r *domainrepo.ResumeRepo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error)
	Update(ctx context.Context, r *domainrepo.ResumeRepo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainrepo.ResumeRepo, int64, error)
	ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error)
	IncrementVersionCount(ctx context.Context, id string, delta int) error
}
