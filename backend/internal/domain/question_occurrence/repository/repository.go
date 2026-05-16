package repository

import (
	"context"

	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)

type ListFilter struct {
	UserID        string
	QuestionID    string
	ApplicationID string
}

type Repository interface {
	Create(ctx context.Context, o *domainoccurrence.QuestionOccurrence) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error)
	Delete(ctx context.Context, id, userID string) error
	DeleteByQuestionID(ctx context.Context, questionID, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainoccurrence.QuestionOccurrence, error)
}
