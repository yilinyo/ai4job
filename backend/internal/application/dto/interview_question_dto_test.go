package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)

func TestInterviewQuestionResponseMarshalsFrontendCamelCaseFields(t *testing.T) {
	lastReviewedAt := time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC)
	question := &domainquestion.InterviewQuestion{
		ID:             "question-1",
		Title:          "Vue 2 响应式原理",
		Content:        "请说明 Vue 2 响应式实现方式。",
		AnswerAnalysis: "核心是 Object.defineProperty。",
		Tags:           domainquestion.EncodeTags([]string{"Vue", "JavaScript"}),
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusNormal,
		IsFavorite:     true,
		Source:         "面试复盘",
		LastReviewedAt: &lastReviewedAt,
		Remark:         "容易追问 nextTick",
		CreatedAt:      time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
	}

	response := ToInterviewQuestionResponse(question, 2)
	payload, err := json.Marshal(response)
	require.NoError(t, err)

	var fields map[string]any
	require.NoError(t, json.Unmarshal(payload, &fields))

	for _, field := range []string{"answerAnalysis", "questionType", "masteryStatus", "isFavorite", "occurrenceCount"} {
		assert.Contains(t, fields, field)
	}
	for _, field := range []string{"AnswerAnalysis", "QuestionType", "MasteryStatus", "IsFavorite", "OccurrenceCount"} {
		assert.NotContains(t, fields, field)
	}
	assert.Equal(t, []string{"Vue", "JavaScript"}, response.Tags)
	assert.Equal(t, int64(2), response.OccurrenceCount)
}

func TestInterviewQuestionDetailResponseIncludesNestedData(t *testing.T) {
	answerPerformance := domainoccurrence.AnswerPerformanceAnswered
	lastReviewedAt := time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC)
	lastResult := domainreview.ResultMastered
	question := &domainquestion.InterviewQuestion{
		ID:             "question-1",
		Title:          "Vue 2 响应式原理",
		Content:        "请说明 Vue 2 响应式实现方式。",
		AnswerAnalysis: "核心是 Object.defineProperty。",
		Tags:           domainquestion.EncodeTags([]string{"Vue"}),
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusNormal,
		CreatedAt:      time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
	}
	occurrences := []*domainoccurrence.QuestionOccurrence{
		{
			ID:                     "occurrence-1",
			QuestionID:             "question-1",
			ApplicationID:          "application-1",
			CompanyNameSnapshot:    "腾讯科技",
			JobTitleSnapshot:       "前端开发工程师",
			InterviewStageSnapshot: "first_interview",
			OccurredAt:             time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			AnswerPerformance:      &answerPerformance,
			CreatedAt:              time.Date(2026, 5, 8, 11, 0, 0, 0, time.UTC),
			UpdatedAt:              time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
		},
	}
	summary := &domainreview.ReviewSummary{
		TotalReviews:   3,
		FailedCount:    1,
		VagueCount:     1,
		MasteredCount:  1,
		LastReviewedAt: &lastReviewedAt,
		LastResult:     &lastResult,
	}
	reviews := []*domainreview.FlashcardReview{
		{
			ID:         "review-1",
			QuestionID: "question-1",
			ReviewedAt: lastReviewedAt,
			Result:     domainreview.ResultMastered,
			ReviewMode: domainreview.ReviewModeTag,
			CreatedAt:  lastReviewedAt,
		},
	}

	response := ToInterviewQuestionDetailResponse(question, 1, occurrences, summary, reviews)
	payload, err := json.Marshal(response)
	require.NoError(t, err)

	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(payload, &fields))
	assert.Contains(t, fields, "occurrences")
	assert.Contains(t, fields, "reviewSummary")
	assert.Contains(t, fields, "recentReviews")
	assert.Len(t, response.Occurrences, 1)
	assert.Equal(t, int64(3), response.ReviewSummary.TotalReviews)
	assert.Len(t, response.RecentReviews, 1)
}

func TestQuestionOccurrenceResponseIncludesAnswerPerformanceAsStringWhenSet(t *testing.T) {
	answerPerformance := domainoccurrence.AnswerPerformanceAnswered
	occurrence := &domainoccurrence.QuestionOccurrence{
		ID:                "occurrence-1",
		QuestionID:        "question-1",
		ApplicationID:     "application-1",
		OccurredAt:        time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
		AnswerPerformance: &answerPerformance,
		CreatedAt:         time.Date(2026, 5, 8, 11, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	}

	response := ToQuestionOccurrenceResponse(occurrence)
	payload, err := json.Marshal(response)
	require.NoError(t, err)

	var fields map[string]any
	require.NoError(t, json.Unmarshal(payload, &fields))
	assert.Equal(t, "answered", fields["answerPerformance"])
	assert.Contains(t, fields, "answerPerformance")
	assert.NotContains(t, fields, "AnswerPerformance")
}

func TestQuestionOccurrenceResponseOmitsEmptyOptionalFrontendFields(t *testing.T) {
	response := &QuestionOccurrenceResponse{
		ID:                  "occurrence-1",
		QuestionID:          "question-1",
		ApplicationID:       "application-1",
		CompanyNameSnapshot: "腾讯科技",
		JobTitleSnapshot:    "前端开发工程师",
		OccurredAt:          time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
		CreatedAt:           time.Date(2026, 5, 8, 11, 0, 0, 0, time.UTC),
		UpdatedAt:           time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)

	var fields map[string]any
	require.NoError(t, json.Unmarshal(payload, &fields))
	assert.NotContains(t, fields, "interviewStageSnapshot")
	assert.NotContains(t, fields, "actualQuestion")
	assert.NotContains(t, fields, "note")
}

func TestToInterviewQuestionDetailResponseReturnsNilForNilQuestion(t *testing.T) {
	response := ToInterviewQuestionDetailResponse(nil, 0, nil, nil, nil)

	assert.Nil(t, response)
}
