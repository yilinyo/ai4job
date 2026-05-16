package interview_question

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func validNewInterviewQuestionParams() NewInterviewQuestionParams {
	return NewInterviewQuestionParams{
		ID:             "question-1",
		UserID:         "user-1",
		Title:          "Explain goroutines",
		Content:        "What is a goroutine?",
		AnswerAnalysis: "A lightweight thread managed by Go runtime.",
		Tags:           []string{"go", "concurrency"},
		Difficulty:     DifficultyMedium,
		QuestionType:   QuestionTypeKnowledge,
		Source:         "manual",
		Remark:         "important",
	}
}

func TestNewInterviewQuestionCopiesFieldsAndDefaults(t *testing.T) {
	question, err := NewInterviewQuestion(validNewInterviewQuestionParams())

	require.NoError(t, err)
	assert.Equal(t, "question-1", question.ID)
	assert.Equal(t, "user-1", question.UserID)
	assert.Equal(t, "Explain goroutines", question.Title)
	assert.Equal(t, "What is a goroutine?", question.Content)
	assert.Equal(t, "A lightweight thread managed by Go runtime.", question.AnswerAnalysis)
	assert.Equal(t, DifficultyMedium, question.Difficulty)
	assert.Equal(t, QuestionTypeKnowledge, question.QuestionType)
	assert.Equal(t, MasteryStatusUnreviewed, question.MasteryStatus)
	assert.False(t, question.IsFavorite)
	assert.Equal(t, "manual", question.Source)
	assert.Equal(t, "important", question.Remark)
	assert.Nil(t, question.LastReviewedAt)

	tags, err := DecodeTags(question.Tags)
	require.NoError(t, err)
	assert.Equal(t, []string{"go", "concurrency"}, tags)
}

func TestNewInterviewQuestionEncodesNilTagsAsEmptyArray(t *testing.T) {
	params := validNewInterviewQuestionParams()
	params.Tags = nil

	question, err := NewInterviewQuestion(params)

	require.NoError(t, err)
	tags, err := DecodeTags(question.Tags)
	require.NoError(t, err)
	assert.Equal(t, []string{}, tags)
}

func TestNewInterviewQuestionMissingRequiredTitleReturnsBadParams(t *testing.T) {
	params := validNewInterviewQuestionParams()
	params.Title = "   "

	question, err := NewInterviewQuestion(params)

	assert.Nil(t, question)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestNewInterviewQuestionInvalidDifficultyReturnsBadParams(t *testing.T) {
	params := validNewInterviewQuestionParams()
	params.Difficulty = Difficulty("impossible")

	question, err := NewInterviewQuestion(params)

	assert.Nil(t, question)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestApplyUpdatePreservesOmittedFieldsAndUpdatesProvidedFields(t *testing.T) {
	question, err := NewInterviewQuestion(validNewInterviewQuestionParams())
	require.NoError(t, err)

	newTitle := "Explain Go channels"
	newTags := []string{"go", "channels"}
	favorite := true
	update := UpdateSpec{
		Title:      &newTitle,
		Tags:       &newTags,
		IsFavorite: &favorite,
	}

	err = question.ApplyUpdate(update)

	require.NoError(t, err)
	assert.Equal(t, "Explain Go channels", question.Title)
	assert.Equal(t, "What is a goroutine?", question.Content)
	assert.Equal(t, "A lightweight thread managed by Go runtime.", question.AnswerAnalysis)
	assert.Equal(t, DifficultyMedium, question.Difficulty)
	assert.Equal(t, QuestionTypeKnowledge, question.QuestionType)
	assert.Equal(t, MasteryStatusUnreviewed, question.MasteryStatus)
	assert.True(t, question.IsFavorite)

	tags, err := DecodeTags(question.Tags)
	require.NoError(t, err)
	assert.Equal(t, []string{"go", "channels"}, tags)
}

func TestApplyReviewResultMapsMasteryStatusAndSetsLastReviewedAt(t *testing.T) {
	cases := []struct {
		name           string
		result         ReviewResult
		expectedStatus MasteryStatus
	}{
		{name: "failed maps to weak", result: ReviewResultFailed, expectedStatus: MasteryStatusWeak},
		{name: "vague maps to normal", result: ReviewResultVague, expectedStatus: MasteryStatusNormal},
		{name: "mastered maps to mastered", result: ReviewResultMastered, expectedStatus: MasteryStatusMastered},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			question, err := NewInterviewQuestion(validNewInterviewQuestionParams())
			require.NoError(t, err)
			reviewedAt := time.Date(2026, 5, 16, 10, 30, 0, 0, time.UTC)

			err = question.ApplyReviewResult(tc.result, reviewedAt)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, question.MasteryStatus)
			require.NotNil(t, question.LastReviewedAt)
			assert.True(t, reviewedAt.Equal(*question.LastReviewedAt))
		})
	}
}
