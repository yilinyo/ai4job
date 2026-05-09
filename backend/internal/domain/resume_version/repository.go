package resume_version

import "context"

type Repository interface {
	Create(ctx context.Context, v *ResumeVersion) error
	FindByIDAndRepoID(ctx context.Context, id, repoID string) (*ResumeVersion, error)
	ListByRepoID(ctx context.Context, repoID string) ([]ResumeVersion, error)
	Update(ctx context.Context, v *ResumeVersion) error
	Delete(ctx context.Context, id, repoID string) error
	MaxVersionNum(ctx context.Context, repoID string) (int, error)
	ClearDefaultByRepoID(ctx context.Context, repoID string) error
}
