package resume_repo

import "context"

type ListFilter struct {
	UserID   string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, r *ResumeRepo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*ResumeRepo, error)
	Update(ctx context.Context, r *ResumeRepo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]ResumeRepo, int64, error)
	ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error)
}
