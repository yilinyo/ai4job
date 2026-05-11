package resume_application

import "context"

type ListFilter struct {
	UserID          string
	ResumeVersionID string
	RepoID          string
	Keyword         string
	CurrentStatus   Status
	Page            int
	PageSize        int
}

type Repository interface {
	Create(ctx context.Context, a *ResumeApplication) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*ResumeApplication, error)
	Update(ctx context.Context, a *ResumeApplication) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]ResumeApplication, int64, error)
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
	UpdateCurrentStatus(ctx context.Context, id, userID string, status Status) error
}
