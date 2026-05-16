package flashcard_review

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func validNewFlashcardReviewParams() NewFlashcardReviewParams {
	return NewFlashcardReviewParams{
		ID:         "review-1",
		UserID:     "user-1",
		QuestionID: "question-1",
		ReviewedAt: time.Date(2026, 5, 16, 11, 0, 0, 0, time.UTC),
		Result:     ResultMastered,
		ReviewMode: ReviewModeWeakFirst,
	}
}

func TestNewFlashcardReviewCopiesFields(t *testing.T) {
	review, err := NewFlashcardReview(validNewFlashcardReviewParams())

	require.NoError(t, err)
	assert.Equal(t, "review-1", review.ID)
	assert.Equal(t, "user-1", review.UserID)
	assert.Equal(t, "question-1", review.QuestionID)
	assert.True(t, time.Date(2026, 5, 16, 11, 0, 0, 0, time.UTC).Equal(review.ReviewedAt))
	assert.Equal(t, ResultMastered, review.Result)
	assert.Equal(t, ReviewModeWeakFirst, review.ReviewMode)
}

func TestNewFlashcardReviewRejectsInvalidResult(t *testing.T) {
	params := validNewFlashcardReviewParams()
	params.Result = Result("unknown")

	review, err := NewFlashcardReview(params)

	assert.Nil(t, review)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestReviewSummaryRecordsLastReview(t *testing.T) {
	lastReviewedAt := time.Date(2026, 5, 16, 11, 0, 0, 0, time.UTC)
	lastResult := ResultMastered

	summary := ReviewSummary{
		TotalReviews:   3,
		FailedCount:    1,
		VagueCount:     1,
		MasteredCount:  1,
		LastReviewedAt: &lastReviewedAt,
		LastResult:     &lastResult,
	}

	assert.Equal(t, int64(3), summary.TotalReviews)
	assert.Equal(t, int64(1), summary.FailedCount)
	assert.Equal(t, int64(1), summary.VagueCount)
	assert.Equal(t, int64(1), summary.MasteredCount)
	require.NotNil(t, summary.LastReviewedAt)
	assert.True(t, lastReviewedAt.Equal(*summary.LastReviewedAt))
	require.NotNil(t, summary.LastResult)
	assert.Equal(t, ResultMastered, *summary.LastResult)
}
