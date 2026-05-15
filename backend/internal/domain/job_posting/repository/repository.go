package repository

import (
	"context"

	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

type ListFilter struct {
	UserID         string
	Keyword        string
	SourcePlatform string
	Status         domainjob.Status
	Page           int
	PageSize       int
}

type Repository interface {
	Create(ctx context.Context, p *domainjob.JobPosting) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error)
	Update(ctx context.Context, p *domainjob.JobPosting) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainjob.JobPosting, int64, error)
}
