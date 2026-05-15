package repository

import (
	"context"

	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type ListFilter struct {
	UserID          string
	ResumeVersionID string
	RepoID          string
	Keyword         string
	CurrentStatus   domainapp.Status
	Page            int
	PageSize        int
}

type Repository interface {
	Create(ctx context.Context, a *domainapp.ResumeApplication) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error)
	Update(ctx context.Context, a *domainapp.ResumeApplication) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainapp.ResumeApplication, int64, error)
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
	UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error
}
