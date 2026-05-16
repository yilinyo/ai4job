package postgres

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	reviewrepo "github.com/yilin/ai-for-backend/internal/domain/flashcard_review/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type FlashcardReviewRepo struct {
	db *gorm.DB
}

var _ reviewrepo.Repository = (*FlashcardReviewRepo)(nil)

func NewFlashcardReviewRepo(db *gorm.DB) *FlashcardReviewRepo {
	return &FlashcardReviewRepo{db: db}
}

func (r *FlashcardReviewRepo) Create(ctx context.Context, review *domainreview.FlashcardReview) error {
	err := dbFromContext(ctx, r.db).WithContext(ctx).Create(review).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *FlashcardReviewRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Delete(&domainreview.FlashcardReview{}).Error
}

func (r *FlashcardReviewRepo) ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error) {
	q := dbFromContext(ctx, r.db).WithContext(ctx).
		Where("question_id = ? AND user_id = ?", questionID, userID).
		Order("reviewed_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}

	var rows []domainreview.FlashcardReview
	err := q.Find(&rows).Error
	return rows, err
}

func (r *FlashcardReviewRepo) SummaryByQuestion(ctx context.Context, questionID, userID string) (*domainreview.ReviewSummary, error) {
	db := dbFromContext(ctx, r.db).WithContext(ctx)

	type aggregateRow struct {
		TotalReviews  int64
		FailedCount   int64
		VagueCount    int64
		MasteredCount int64
	}
	var aggregate aggregateRow
	if err := db.Model(&domainreview.FlashcardReview{}).
		Select("COUNT(*) AS total_reviews, COUNT(*) FILTER (WHERE result = ?) AS failed_count, COUNT(*) FILTER (WHERE result = ?) AS vague_count, COUNT(*) FILTER (WHERE result = ?) AS mastered_count", domainreview.ResultFailed, domainreview.ResultVague, domainreview.ResultMastered).
		Where("question_id = ? AND user_id = ?", questionID, userID).
		Scan(&aggregate).Error; err != nil {
		return nil, err
	}

	summary := &domainreview.ReviewSummary{
		TotalReviews:  aggregate.TotalReviews,
		FailedCount:   aggregate.FailedCount,
		VagueCount:    aggregate.VagueCount,
		MasteredCount: aggregate.MasteredCount,
	}
	if aggregate.TotalReviews == 0 {
		return summary, nil
	}

	var latest domainreview.FlashcardReview
	if err := db.Where("question_id = ? AND user_id = ?", questionID, userID).Order("reviewed_at DESC").First(&latest).Error; err != nil {
		return nil, err
	}
	summary.LastReviewedAt = &latest.ReviewedAt
	summary.LastResult = &latest.Result
	return summary, nil
}
