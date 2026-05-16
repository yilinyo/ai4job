package repository

import (
	"context"

	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

type Repository interface {
	Create(ctx context.Context, v *domainversion.ResumeVersion) error
	FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error)
	FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error)
	ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error)
	Update(ctx context.Context, v *domainversion.ResumeVersion) error
	Delete(ctx context.Context, id, repoID string) error
	CountByRepoID(ctx context.Context, repoID string) (int64, error)
	SetActiveByID(ctx context.Context, repoID, id string) error
}
