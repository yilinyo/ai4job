package repository

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
)

type Repository interface {
	Create(ctx context.Context, r *domainreview.FlashcardReview) error
	DeleteByQuestionID(ctx context.Context, questionID, userID string) error
	ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error)
	SummaryByQuestion(ctx context.Context, questionID, userID string) (*domainreview.ReviewSummary, error)
}
