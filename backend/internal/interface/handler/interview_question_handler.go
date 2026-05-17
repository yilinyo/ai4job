package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type InterviewQuestionService interface {
	CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error)
	ListQuestions(ctx context.Context, userID string, req dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error)
	GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error)
	UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error)
	DeleteQuestion(ctx context.Context, userID, id string) error
	UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error)
	ListOccurrences(ctx context.Context, userID string, req dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error)
	CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error)
	DeleteOccurrence(ctx context.Context, userID, id string) error
	ListFlashcards(ctx context.Context, userID string, req dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error)
	CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error)
}

type InterviewQuestionHandler struct {
	svc InterviewQuestionService
}

func NewInterviewQuestionHandler(svc InterviewQuestionService) *InterviewQuestionHandler {
	return &InterviewQuestionHandler{svc: svc}
}

func (h *InterviewQuestionHandler) ListQuestions(c *gin.Context) {
	var req dto.InterviewQuestionQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	page, pageSize := frontendPage(c)
	req.Page = page
	req.PageSize = pageSize
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	rows, total, err := h.svc.ListQuestions(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toInterviewQuestionResponses(rows), "total": total, "page": page, "pageSize": pageSize})
}

func (h *InterviewQuestionHandler) GetQuestion(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, count, occurrences, summary, reviews, err := h.svc.GetQuestion(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionDetailResponse(question, count, occurrencePointers(occurrences), &summary, reviewPointers(reviews)))
}

func (h *InterviewQuestionHandler) CreateQuestion(c *gin.Context) {
	var req dto.CreateInterviewQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, err := h.svc.CreateQuestion(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, 0))
}

func (h *InterviewQuestionHandler) UpdateQuestion(c *gin.Context) {
	var req dto.UpdateInterviewQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	question, err := h.svc.UpdateQuestion(c.Request.Context(), userID, id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	count, ok := h.questionOccurrenceCount(c, userID, id)
	if !ok {
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, count))
}

func (h *InterviewQuestionHandler) DeleteQuestion(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteQuestion(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *InterviewQuestionHandler) UpdateFavorite(c *gin.Context) {
	var req dto.UpdateQuestionFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.IsFavorite == nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	question, err := h.svc.UpdateFavorite(c.Request.Context(), userID, id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	count, ok := h.questionOccurrenceCount(c, userID, id)
	if !ok {
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, count))
}

func (h *InterviewQuestionHandler) ListOccurrences(c *gin.Context) {
	var req dto.QuestionOccurrenceQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	rows, err := h.svc.ListOccurrences(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": dto.ToQuestionOccurrenceResponses(occurrencePointers(rows))})
}

func (h *InterviewQuestionHandler) CreateOccurrence(c *gin.Context) {
	var req dto.CreateQuestionOccurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.OccurredAt == nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	occurrence, err := h.svc.CreateOccurrence(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToQuestionOccurrenceResponse(occurrence))
}

func (h *InterviewQuestionHandler) DeleteOccurrence(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteOccurrence(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *InterviewQuestionHandler) ListFlashcards(c *gin.Context) {
	var req dto.FlashcardQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	rows, total, err := h.svc.ListFlashcards(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toInterviewQuestionResponses(rows), "total": total, "mode": req.Mode})
}

func (h *InterviewQuestionHandler) CreateFlashcardReview(c *gin.Context) {
	var req dto.CreateFlashcardReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	review, updatedQuestion, summary, err := h.svc.CreateFlashcardReview(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{
		"review":          dto.ToFlashcardReviewResponse(review),
		"updatedQuestion": dto.ToInterviewQuestionResponse(&updatedQuestion.Question, updatedQuestion.OccurrenceCount),
		"reviewSummary":   dto.ToFlashcardReviewSummaryResponse(&summary),
	})
}

func (h *InterviewQuestionHandler) questionOccurrenceCount(c *gin.Context, userID, id string) (int64, bool) {
	_, count, _, _, _, err := h.svc.GetQuestion(c.Request.Context(), userID, id)
	if err != nil {
		mapError(c, err)
		return 0, false
	}
	return count, true
}

func toInterviewQuestionResponses(rows []questionrepo.QuestionWithOccurrenceCount) []*dto.InterviewQuestionResponse {
	responses := make([]*dto.InterviewQuestionResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, dto.ToInterviewQuestionResponse(&rows[i].Question, rows[i].OccurrenceCount))
	}
	return responses
}

func occurrencePointers(rows []domainoccurrence.QuestionOccurrence) []*domainoccurrence.QuestionOccurrence {
	pointers := make([]*domainoccurrence.QuestionOccurrence, 0, len(rows))
	for i := range rows {
		pointers = append(pointers, &rows[i])
	}
	return pointers
}

func reviewPointers(rows []domainreview.FlashcardReview) []*domainreview.FlashcardReview {
	pointers := make([]*domainreview.FlashcardReview, 0, len(rows))
	for i := range rows {
		pointers = append(pointers, &rows[i])
	}
	return pointers
}
