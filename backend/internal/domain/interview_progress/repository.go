package interview_progress

import "context"

type Repository interface {
	Create(ctx context.Context, p *InterviewProgress) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*InterviewProgress, error)
	Update(ctx context.Context, p *InterviewProgress) error
	Delete(ctx context.Context, id, userID string) error
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]InterviewProgress, error)
	DeleteByApplicationID(ctx context.Context, applicationID, userID string) error
}
