package job_posting

import "context"

type ListFilter struct {
	UserID         string
	Keyword        string
	SourcePlatform string
	Status         Status
	Page           int
	PageSize       int
}

type Repository interface {
	Create(ctx context.Context, p *JobPosting) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*JobPosting, error)
	Update(ctx context.Context, p *JobPosting) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]JobPosting, int64, error)
}
