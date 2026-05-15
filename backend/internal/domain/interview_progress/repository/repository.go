package repository

import (
	"context"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
)

type Repository interface {
	Create(ctx context.Context, p *domainprogress.InterviewProgress) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error)
	Update(ctx context.Context, p *domainprogress.InterviewProgress) error
	Delete(ctx context.Context, id, userID string) error
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error)
	DeleteByApplicationID(ctx context.Context, applicationID, userID string) error
}
