package dto

import (
	"encoding/json"
	"strings"
	"time"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)

type FlexibleTime struct {
	time.Time
}

func (t *FlexibleTime) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parsed, err := parseFlexibleTime(value)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t FlexibleTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

func parseFlexibleTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Parse(time.RFC3339, value)
}

type CreateInterviewQuestionRequest struct {
	Title          string                       `json:"title" binding:"required,max=200"`
	Content        string                       `json:"content" binding:"required"`
	AnswerAnalysis string                       `json:"answerAnalysis" binding:"required"`
	Tags           []string                     `json:"tags"`
	Difficulty     domainquestion.Difficulty    `json:"difficulty" binding:"required"`
	QuestionType   domainquestion.QuestionType  `json:"questionType" binding:"required"`
	MasteryStatus  domainquestion.MasteryStatus `json:"masteryStatus" binding:"required"`
	IsFavorite     bool                         `json:"isFavorite"`
	Source         string                       `json:"source" binding:"max=100"`
	Remark         string                       `json:"remark"`
}

type UpdateInterviewQuestionRequest struct {
	Title          *string                       `json:"title" binding:"omitempty,max=200"`
	Content        *string                       `json:"content"`
	AnswerAnalysis *string                       `json:"answerAnalysis"`
	Tags           *[]string                     `json:"tags"`
	Difficulty     *domainquestion.Difficulty    `json:"difficulty"`
	QuestionType   *domainquestion.QuestionType  `json:"questionType"`
	MasteryStatus  *domainquestion.MasteryStatus `json:"masteryStatus"`
	IsFavorite     *bool                         `json:"isFavorite"`
	Source         *string                       `json:"source" binding:"omitempty,max=100"`
	Remark         *string                       `json:"remark"`
}

type UpdateQuestionFavoriteRequest struct {
	IsFavorite *bool `json:"isFavorite" binding:"required"`
}

type InterviewQuestionQuery struct {
	Keyword       string                       `form:"keyword" json:"keyword"`
	Tags          string                       `form:"tags" json:"tags"`
	Difficulty    domainquestion.Difficulty    `form:"difficulty" json:"difficulty"`
	QuestionType  domainquestion.QuestionType  `form:"questionType" json:"questionType"`
	MasteryStatus domainquestion.MasteryStatus `form:"masteryStatus" json:"masteryStatus"`
	IsFavorite    *bool                        `form:"isFavorite" json:"isFavorite"`
	Page          int                          `form:"page" json:"page"`
	PageSize      int                          `form:"pageSize" json:"pageSize"`
}

type CreateQuestionOccurrenceRequest struct {
	QuestionID             string                              `json:"questionId" binding:"required"`
	ApplicationID          string                              `json:"applicationId" binding:"required"`
	InterviewProgressID    *string                             `json:"interviewProgressId"`
	InterviewStageSnapshot string                              `json:"interviewStageSnapshot"`
	OccurredAt             *FlexibleTime                       `json:"occurredAt" binding:"required"`
	ActualQuestion         string                              `json:"actualQuestion"`
	AnswerPerformance      *domainoccurrence.AnswerPerformance `json:"answerPerformance"`
	Note                   string                              `json:"note"`
}

type QuestionOccurrenceQuery struct {
	QuestionID    string `form:"questionId" json:"questionId"`
	ApplicationID string `form:"applicationId" json:"applicationId"`
}

type FlashcardQuery struct {
	Mode       domainreview.ReviewMode   `form:"mode" json:"mode"`
	Tags       string                    `form:"tags" json:"tags"`
	Difficulty domainquestion.Difficulty `form:"difficulty" json:"difficulty"`
	Limit      int                       `form:"limit" json:"limit"`
}

type CreateFlashcardReviewRequest struct {
	QuestionID string                  `json:"questionId" binding:"required"`
	ReviewedAt *FlexibleTime           `json:"reviewedAt"`
	Result     domainreview.Result     `json:"result" binding:"required"`
	ReviewMode domainreview.ReviewMode `json:"reviewMode" binding:"required"`
}

type InterviewQuestionResponse struct {
	ID              string                       `json:"id"`
	Title           string                       `json:"title"`
	Content         string                       `json:"content"`
	AnswerAnalysis  string                       `json:"answerAnalysis"`
	Tags            []string                     `json:"tags"`
	Difficulty      domainquestion.Difficulty    `json:"difficulty"`
	QuestionType    domainquestion.QuestionType  `json:"questionType"`
	MasteryStatus   domainquestion.MasteryStatus `json:"masteryStatus"`
	IsFavorite      bool                         `json:"isFavorite"`
	Source          string                       `json:"source,omitempty"`
	OccurrenceCount int64                        `json:"occurrenceCount"`
	LastReviewedAt  *time.Time                   `json:"lastReviewedAt,omitempty"`
	Remark          string                       `json:"remark,omitempty"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
}

type InterviewQuestionDetailResponse struct {
	*InterviewQuestionResponse
	Occurrences   []*QuestionOccurrenceResponse   `json:"occurrences"`
	ReviewSummary *FlashcardReviewSummaryResponse `json:"reviewSummary"`
	RecentReviews []*FlashcardReviewResponse      `json:"recentReviews"`
}

type QuestionOccurrenceResponse struct {
	ID                     string                              `json:"id"`
	QuestionID             string                              `json:"questionId"`
	ApplicationID          string                              `json:"applicationId"`
	InterviewProgressID    *string                             `json:"interviewProgressId,omitempty"`
	CompanyNameSnapshot    string                              `json:"companyNameSnapshot"`
	JobTitleSnapshot       string                              `json:"jobTitleSnapshot"`
	InterviewStageSnapshot string                              `json:"interviewStageSnapshot,omitempty"`
	OccurredAt             time.Time                           `json:"occurredAt"`
	ActualQuestion         string                              `json:"actualQuestion,omitempty"`
	AnswerPerformance      *domainoccurrence.AnswerPerformance `json:"answerPerformance,omitempty"`
	Note                   string                              `json:"note,omitempty"`
	CreatedAt              time.Time                           `json:"createdAt"`
	UpdatedAt              time.Time                           `json:"updatedAt"`
}

type FlashcardReviewResponse struct {
	ID         string                  `json:"id"`
	QuestionID string                  `json:"questionId"`
	ReviewedAt time.Time               `json:"reviewedAt"`
	Result     domainreview.Result     `json:"result"`
	ReviewMode domainreview.ReviewMode `json:"reviewMode"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type FlashcardReviewSummaryResponse struct {
	TotalReviews   int64                `json:"totalReviews"`
	FailedCount    int64                `json:"failedCount"`
	VagueCount     int64                `json:"vagueCount"`
	MasteredCount  int64                `json:"masteredCount"`
	LastReviewedAt *time.Time           `json:"lastReviewedAt,omitempty"`
	LastResult     *domainreview.Result `json:"lastResult,omitempty"`
}

func ToInterviewQuestionResponse(q *domainquestion.InterviewQuestion, occurrenceCount int64) *InterviewQuestionResponse {
	if q == nil {
		return nil
	}
	tags, err := domainquestion.DecodeTags(q.Tags)
	if err != nil {
		tags = []string{}
	}
	return &InterviewQuestionResponse{
		ID:              q.ID,
		Title:           q.Title,
		Content:         q.Content,
		AnswerAnalysis:  q.AnswerAnalysis,
		Tags:            tags,
		Difficulty:      q.Difficulty,
		QuestionType:    q.QuestionType,
		MasteryStatus:   q.MasteryStatus,
		IsFavorite:      q.IsFavorite,
		Source:          q.Source,
		OccurrenceCount: occurrenceCount,
		LastReviewedAt:  q.LastReviewedAt,
		Remark:          q.Remark,
		CreatedAt:       q.CreatedAt,
		UpdatedAt:       q.UpdatedAt,
	}
}

func ToInterviewQuestionDetailResponse(q *domainquestion.InterviewQuestion, occurrenceCount int64, occurrences []*domainoccurrence.QuestionOccurrence, summary *domainreview.ReviewSummary, reviews []*domainreview.FlashcardReview) *InterviewQuestionDetailResponse {
	if q == nil {
		return nil
	}
	return &InterviewQuestionDetailResponse{
		InterviewQuestionResponse: ToInterviewQuestionResponse(q, occurrenceCount),
		Occurrences:               ToQuestionOccurrenceResponses(occurrences),
		ReviewSummary:             ToFlashcardReviewSummaryResponse(summary),
		RecentReviews:             ToFlashcardReviewResponses(reviews),
	}
}

func ToQuestionOccurrenceResponse(o *domainoccurrence.QuestionOccurrence) *QuestionOccurrenceResponse {
	if o == nil {
		return nil
	}
	return &QuestionOccurrenceResponse{
		ID:                     o.ID,
		QuestionID:             o.QuestionID,
		ApplicationID:          o.ApplicationID,
		InterviewProgressID:    o.InterviewProgressID,
		CompanyNameSnapshot:    o.CompanyNameSnapshot,
		JobTitleSnapshot:       o.JobTitleSnapshot,
		InterviewStageSnapshot: o.InterviewStageSnapshot,
		OccurredAt:             o.OccurredAt,
		ActualQuestion:         o.ActualQuestion,
		AnswerPerformance:      o.AnswerPerformance,
		Note:                   o.Note,
		CreatedAt:              o.CreatedAt,
		UpdatedAt:              o.UpdatedAt,
	}
}

func ToQuestionOccurrenceResponses(occurrences []*domainoccurrence.QuestionOccurrence) []*QuestionOccurrenceResponse {
	responses := make([]*QuestionOccurrenceResponse, 0, len(occurrences))
	for _, occurrence := range occurrences {
		responses = append(responses, ToQuestionOccurrenceResponse(occurrence))
	}
	return responses
}

func ToFlashcardReviewResponse(r *domainreview.FlashcardReview) *FlashcardReviewResponse {
	if r == nil {
		return nil
	}
	return &FlashcardReviewResponse{
		ID:         r.ID,
		QuestionID: r.QuestionID,
		ReviewedAt: r.ReviewedAt,
		Result:     r.Result,
		ReviewMode: r.ReviewMode,
		CreatedAt:  r.CreatedAt,
	}
}

func ToFlashcardReviewResponses(reviews []*domainreview.FlashcardReview) []*FlashcardReviewResponse {
	responses := make([]*FlashcardReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		responses = append(responses, ToFlashcardReviewResponse(review))
	}
	return responses
}

func ToFlashcardReviewSummaryResponse(summary *domainreview.ReviewSummary) *FlashcardReviewSummaryResponse {
	if summary == nil {
		return &FlashcardReviewSummaryResponse{}
	}
	return &FlashcardReviewSummaryResponse{
		TotalReviews:   summary.TotalReviews,
		FailedCount:    summary.FailedCount,
		VagueCount:     summary.VagueCount,
		MasteredCount:  summary.MasteredCount,
		LastReviewedAt: summary.LastReviewedAt,
		LastResult:     summary.LastResult,
	}
}
