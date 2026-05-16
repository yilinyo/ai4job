package repository

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
)

type ListFilter struct {
	UserID        string
	Keyword       string
	Tags          []string
	Difficulty    domainquestion.Difficulty
	QuestionType  domainquestion.QuestionType
	MasteryStatus domainquestion.MasteryStatus
	IsFavorite    *bool
	Page          int
	PageSize      int
}

type FlashcardFilter struct {
	UserID     string
	Mode       domainreview.ReviewMode
	Tags       []string
	Difficulty domainquestion.Difficulty
	Limit      int
}

type QuestionWithOccurrenceCount struct {
	Question        domainquestion.InterviewQuestion
	OccurrenceCount int64
}

type Repository interface {
	Create(ctx context.Context, q *domainquestion.InterviewQuestion) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error)
	Update(ctx context.Context, q *domainquestion.InterviewQuestion) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]QuestionWithOccurrenceCount, int64, error)
	CountOccurrences(ctx context.Context, questionID, userID string) (int64, error)
	ListFlashcards(ctx context.Context, f FlashcardFilter) ([]QuestionWithOccurrenceCount, int64, error)
}
