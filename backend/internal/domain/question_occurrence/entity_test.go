package question_occurrence

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func validNewQuestionOccurrenceParams() NewQuestionOccurrenceParams {
	return NewQuestionOccurrenceParams{
		ID:                     "occurrence-1",
		UserID:                 "user-1",
		QuestionID:             "question-1",
		ApplicationID:          "application-1",
		InterviewProgressID:    nil,
		CompanyNameSnapshot:    "Acme Inc",
		JobTitleSnapshot:       "Backend Engineer",
		InterviewStageSnapshot: "technical",
		OccurredAt:             time.Date(2026, 5, 16, 9, 0, 0, 0, time.UTC),
		ActualQuestion:         "Tell me about channels.",
		AnswerPerformance:      nil,
		Note:                   "asked in first round",
	}
}

func TestNewQuestionOccurrenceAllowsApplicationOnlyAssociation(t *testing.T) {
	occurrence, err := NewQuestionOccurrence(validNewQuestionOccurrenceParams())

	require.NoError(t, err)
	assert.Equal(t, "occurrence-1", occurrence.ID)
	assert.Equal(t, "user-1", occurrence.UserID)
	assert.Equal(t, "question-1", occurrence.QuestionID)
	assert.Equal(t, "application-1", occurrence.ApplicationID)
	assert.Nil(t, occurrence.InterviewProgressID)
	assert.Equal(t, "Acme Inc", occurrence.CompanyNameSnapshot)
	assert.Equal(t, "Backend Engineer", occurrence.JobTitleSnapshot)
	assert.Equal(t, "technical", occurrence.InterviewStageSnapshot)
	assert.Equal(t, "Tell me about channels.", occurrence.ActualQuestion)
	assert.Nil(t, occurrence.AnswerPerformance)
	assert.Equal(t, "asked in first round", occurrence.Note)
	assert.True(t, time.Date(2026, 5, 16, 9, 0, 0, 0, time.UTC).Equal(occurrence.OccurredAt))
}

func TestNewQuestionOccurrenceRejectsMissingQuestion(t *testing.T) {
	params := validNewQuestionOccurrenceParams()
	params.QuestionID = ""

	occurrence, err := NewQuestionOccurrence(params)

	assert.Nil(t, occurrence)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestNewQuestionOccurrenceRejectsMissingApplication(t *testing.T) {
	params := validNewQuestionOccurrenceParams()
	params.ApplicationID = ""

	occurrence, err := NewQuestionOccurrence(params)

	assert.Nil(t, occurrence)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestNewQuestionOccurrenceRejectsInvalidAnswerPerformance(t *testing.T) {
	invalid := AnswerPerformance("excellent")
	params := validNewQuestionOccurrenceParams()
	params.AnswerPerformance = &invalid

	occurrence, err := NewQuestionOccurrence(params)

	assert.Nil(t, occurrence)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}
