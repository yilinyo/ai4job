package flashcard_review

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Result string

const (
	ResultFailed   Result = "failed"
	ResultVague    Result = "vague"
	ResultMastered Result = "mastered"
)

type ReviewMode string

const (
	ReviewModeRandom        ReviewMode = "random"
	ReviewModeTag           ReviewMode = "tag"
	ReviewModeDifficulty    ReviewMode = "difficulty"
	ReviewModeFavoriteFirst ReviewMode = "favorite_first"
	ReviewModeWeakFirst     ReviewMode = "weak_first"
	ReviewModeFrequentFirst ReviewMode = "frequent_first"
)

type FlashcardReview struct {
	ID         string     `gorm:"type:uuid;primaryKey"`
	UserID     string     `gorm:"type:varchar(36);not null;index"`
	QuestionID string     `gorm:"type:uuid;not null;index"`
	ReviewedAt time.Time  `gorm:"not null"`
	Result     Result     `gorm:"type:varchar(20);not null"`
	ReviewMode ReviewMode `gorm:"type:varchar(30);not null"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type NewFlashcardReviewParams struct {
	ID         string
	UserID     string
	QuestionID string
	ReviewedAt time.Time
	Result     Result
	ReviewMode ReviewMode
}

type ReviewSummary struct {
	TotalReviews   int64
	FailedCount    int64
	VagueCount     int64
	MasteredCount  int64
	LastReviewedAt *time.Time
	LastResult     *Result
}

func NewFlashcardReview(params NewFlashcardReviewParams) (*FlashcardReview, error) {
	review := &FlashcardReview{
		ID:         params.ID,
		UserID:     params.UserID,
		QuestionID: params.QuestionID,
		ReviewedAt: params.ReviewedAt,
		Result:     params.Result,
		ReviewMode: params.ReviewMode,
	}
	if err := review.Validate(); err != nil {
		return nil, err
	}
	return review, nil
}

func (r *FlashcardReview) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.UserID) == "" || strings.TrimSpace(r.QuestionID) == "" {
		return domainerrors.ErrBadParams
	}
	if r.ReviewedAt.IsZero() || !r.Result.Valid() || !r.ReviewMode.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (r Result) Valid() bool {
	switch r {
	case ResultFailed, ResultVague, ResultMastered:
		return true
	default:
		return false
	}
}

func (m ReviewMode) Valid() bool {
	switch m {
	case ReviewModeRandom, ReviewModeTag, ReviewModeDifficulty, ReviewModeFavoriteFirst, ReviewModeWeakFirst, ReviewModeFrequentFirst:
		return true
	default:
		return false
	}
}

func (FlashcardReview) TableName() string {
	return "flashcard_reviews"
}
