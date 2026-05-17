package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockInterviewQuestionService struct{ mock.Mock }

func (m *mockInterviewQuestionService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, req)
	question, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return question, args.Error(1)
}

func (m *mockInterviewQuestionService) ListQuestions(ctx context.Context, userID string, req dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}

func (m *mockInterviewQuestionService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
	args := m.Called(ctx, userID, id)
	question, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return question, args.Get(1).(int64), args.Get(2).([]domainoccurrence.QuestionOccurrence), args.Get(3).(domainreview.ReviewSummary), args.Get(4).([]domainreview.FlashcardReview), args.Error(5)
}

func (m *mockInterviewQuestionService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, id, req)
	question, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return question, args.Error(1)
}

func (m *mockInterviewQuestionService) DeleteQuestion(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *mockInterviewQuestionService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, id, req)
	question, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return question, args.Error(1)
}

func (m *mockInterviewQuestionService) ListOccurrences(ctx context.Context, userID string, req dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).([]domainoccurrence.QuestionOccurrence), args.Error(1)
}

func (m *mockInterviewQuestionService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, userID, req)
	occurrence, _ := args.Get(0).(*domainoccurrence.QuestionOccurrence)
	return occurrence, args.Error(1)
}

func (m *mockInterviewQuestionService) DeleteOccurrence(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *mockInterviewQuestionService) ListFlashcards(ctx context.Context, userID string, req dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}

func (m *mockInterviewQuestionService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	args := m.Called(ctx, userID, req)
	review, _ := args.Get(0).(*domainreview.FlashcardReview)
	return review, args.Get(1).(questionrepo.QuestionWithOccurrenceCount), args.Get(2).(domainreview.ReviewSummary), args.Error(3)
}

func TestInterviewQuestionHandler_ListCreateDetailUseFrontendShape(t *testing.T) {
	svc := &mockInterviewQuestionService{}
	createdAt := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)
	question := domainquestion.InterviewQuestion{
		ID:             "q1",
		UserID:         "u1",
		Title:          "Explain indexes",
		Content:        "What is a DB index?",
		AnswerAnalysis: "Speeds reads",
		Tags:           domainquestion.EncodeTags([]string{"db", "backend"}),
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusUnreviewed,
		IsFavorite:     true,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}
	svc.On("ListQuestions", mock.Anything, "u1", dto.InterviewQuestionQuery{Keyword: "index", Page: 1, PageSize: 10}).Return([]questionrepo.QuestionWithOccurrenceCount{{Question: question, OccurrenceCount: 2}}, int64(1), nil).Once()
	svc.On("CreateQuestion", mock.Anything, "u1", dto.CreateInterviewQuestionRequest{Title: "Explain indexes", Content: "What is a DB index?", AnswerAnalysis: "Speeds reads", Tags: []string{"db"}, Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusUnreviewed, IsFavorite: true}).Return(&question, nil).Once()
	svc.On("GetQuestion", mock.Anything, "u1", "q1").Return(&question, int64(3), []domainoccurrence.QuestionOccurrence{{ID: "occ-1", UserID: "u1", QuestionID: "q1", ApplicationID: "app-1", CompanyNameSnapshot: "Acme", JobTitleSnapshot: "Backend", OccurredAt: createdAt, CreatedAt: createdAt, UpdatedAt: createdAt}}, domainreview.ReviewSummary{TotalReviews: 1}, []domainreview.FlashcardReview{}, nil).Once()
	r := gin.New()
	h := handler.NewInterviewQuestionHandler(svc)
	r.GET("/questions", func(c *gin.Context) { c.Set("userID", "u1"); h.ListQuestions(c) })
	r.POST("/questions", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateQuestion(c) })
	r.GET("/questions/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.GetQuestion(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/questions?keyword=index", nil)
	r.ServeHTTP(w, req)
	data := responseData(t, w)
	first := data["list"].([]any)[0].(map[string]any)
	assert.Equal(t, "Explain indexes", first["title"])
	assert.Equal(t, "Speeds reads", first["answerAnalysis"])
	assert.Equal(t, float64(2), first["occurrenceCount"])
	assert.Equal(t, true, first["isFavorite"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["pageSize"])
	assert.NotContains(t, first, "AnswerAnalysis")
	assert.NotContains(t, first, "OccurrenceCount")
	assert.NotContains(t, first, "UserID")

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBufferString(`{"title":"Explain indexes","content":"What is a DB index?","answerAnalysis":"Speeds reads","tags":["db"],"difficulty":"medium","questionType":"knowledge","masteryStatus":"unreviewed","isFavorite":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	data = responseData(t, w)
	assert.Equal(t, "Explain indexes", data["title"])
	assert.Equal(t, float64(0), data["occurrenceCount"])
	assert.NotContains(t, data, "AnswerAnalysis")

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/questions/q1", nil)
	r.ServeHTTP(w, req)
	data = responseData(t, w)
	assert.Equal(t, "Explain indexes", data["title"])
	assert.Equal(t, float64(3), data["occurrenceCount"])
	assert.Contains(t, data, "reviewSummary")
	assert.Contains(t, data, "recentReviews")
	assert.NotContains(t, data, "AnswerAnalysis")
	svc.AssertExpectations(t)
}

func TestInterviewQuestionHandler_FavoriteValidationAndExplicitFalse(t *testing.T) {
	svc := &mockInterviewQuestionService{}
	question := domainquestion.InterviewQuestion{ID: "q1", UserID: "u1", Title: "T", Content: "C", AnswerAnalysis: "A", Tags: domainquestion.EncodeTags(nil), Difficulty: domainquestion.DifficultyEasy, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusWeak}
	isFavorite := false
	svc.On("UpdateFavorite", mock.Anything, "u1", "q1", dto.UpdateQuestionFavoriteRequest{IsFavorite: &isFavorite}).Return(&question, nil).Once()
	svc.On("GetQuestion", mock.Anything, "u1", "q1").Return(&question, int64(4), []domainoccurrence.QuestionOccurrence{}, domainreview.ReviewSummary{}, []domainreview.FlashcardReview{}, nil).Once()
	r := gin.New()
	h := handler.NewInterviewQuestionHandler(svc)
	r.PUT("/questions/:id/favorite", func(c *gin.Context) { c.Set("userID", "u1"); h.UpdateFavorite(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/questions/q1/favorite", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
	svc.AssertNotCalled(t, "UpdateFavorite", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/questions/q1/favorite", bytes.NewBufferString(`{"isFavorite":false}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	data := responseData(t, w)
	assert.Equal(t, false, data["isFavorite"])
	assert.Equal(t, float64(4), data["occurrenceCount"])
	svc.AssertExpectations(t)
}

func TestInterviewQuestionHandler_OccurrenceValidationAndCreateShape(t *testing.T) {
	svc := &mockInterviewQuestionService{}
	occurredAt := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	performance := domainoccurrence.AnswerPerformanceGood
	createReq := dto.CreateQuestionOccurrenceRequest{QuestionID: "q1", ApplicationID: "app-1", OccurredAt: &dto.FlexibleTime{Time: occurredAt}, AnswerPerformance: &performance, ActualQuestion: "Actual", Note: "note"}
	row := &domainoccurrence.QuestionOccurrence{ID: "occ-1", UserID: "u1", QuestionID: "q1", ApplicationID: "app-1", CompanyNameSnapshot: "Acme", JobTitleSnapshot: "Backend", OccurredAt: occurredAt, AnswerPerformance: &performance, ActualQuestion: "Actual", Note: "note", CreatedAt: occurredAt, UpdatedAt: occurredAt}
	svc.On("CreateOccurrence", mock.Anything, "u1", createReq).Return(row, nil).Once()
	r := gin.New()
	h := handler.NewInterviewQuestionHandler(svc)
	r.POST("/occurrences", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateOccurrence(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/occurrences", bytes.NewBufferString(`{"questionId":"q1","applicationId":"app-1"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeBadParams, responseCode(t, w))
	svc.AssertNotCalled(t, "CreateOccurrence", mock.Anything, mock.Anything, mock.Anything)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/occurrences", bytes.NewBufferString(`{"questionId":"q1","applicationId":"app-1","occurredAt":"2026-05-15T09:30:00Z","answerPerformance":"good","actualQuestion":"Actual","note":"note"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	data := responseData(t, w)
	assert.Equal(t, "q1", data["questionId"])
	assert.Equal(t, "app-1", data["applicationId"])
	assert.Equal(t, "Acme", data["companyNameSnapshot"])
	assert.Equal(t, "good", data["answerPerformance"])
	assert.NotContains(t, data, "QuestionID")
	assert.NotContains(t, data, "ApplicationID")
	svc.AssertExpectations(t)
}

func TestInterviewQuestionHandler_FlashcardReviewShape(t *testing.T) {
	svc := &mockInterviewQuestionService{}
	reviewedAt := time.Date(2026, 5, 15, 11, 0, 0, 0, time.UTC)
	reqDTO := dto.CreateFlashcardReviewRequest{QuestionID: "q1", ReviewedAt: &dto.FlexibleTime{Time: reviewedAt}, Result: domainreview.ResultMastered, ReviewMode: domainreview.ReviewModeRandom}
	review := &domainreview.FlashcardReview{ID: "rev-1", UserID: "u1", QuestionID: "q1", ReviewedAt: reviewedAt, Result: domainreview.ResultMastered, ReviewMode: domainreview.ReviewModeRandom, CreatedAt: reviewedAt}
	question := domainquestion.InterviewQuestion{ID: "q1", UserID: "u1", Title: "Explain indexes", Content: "C", AnswerAnalysis: "A", Tags: domainquestion.EncodeTags([]string{"db"}), Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.QuestionTypeKnowledge, MasteryStatus: domainquestion.MasteryStatusMastered, LastReviewedAt: &reviewedAt}
	lastResult := domainreview.ResultMastered
	svc.On("CreateFlashcardReview", mock.Anything, "u1", reqDTO).Return(review, questionrepo.QuestionWithOccurrenceCount{Question: question, OccurrenceCount: 5}, domainreview.ReviewSummary{TotalReviews: 2, MasteredCount: 2, LastReviewedAt: &reviewedAt, LastResult: &lastResult}, nil).Once()
	r := gin.New()
	h := handler.NewInterviewQuestionHandler(svc)
	r.POST("/reviews", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateFlashcardReview(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewBufferString(`{"questionId":"q1","reviewedAt":"2026-05-15T11:00:00Z","result":"mastered","reviewMode":"random"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, pkgerrors.CodeSuccess, responseCode(t, w))
	data := responseData(t, w)
	reviewData := data["review"].(map[string]any)
	updatedQuestion := data["updatedQuestion"].(map[string]any)
	summary := data["reviewSummary"].(map[string]any)
	assert.Equal(t, "rev-1", reviewData["id"])
	assert.Equal(t, "q1", reviewData["questionId"])
	assert.Equal(t, "Explain indexes", updatedQuestion["title"])
	assert.Equal(t, float64(5), updatedQuestion["occurrenceCount"])
	assert.Equal(t, float64(2), summary["totalReviews"])
	assert.Equal(t, "mastered", summary["lastResult"])
	assert.NotContains(t, reviewData, "QuestionID")
	assert.NotContains(t, updatedQuestion, "OccurrenceCount")
	svc.AssertExpectations(t)
}
