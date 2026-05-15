package handler

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type InterviewProgressService interface {
	Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) (*domainprogress.InterviewProgress, error)
	List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) (*domainprogress.InterviewProgress, error)
	Delete(ctx context.Context, userID, id string) error
}

type InterviewProgressHandler struct {
	svc InterviewProgressService
}

func NewInterviewProgressHandler(svc InterviewProgressService) *InterviewProgressHandler {
	return &InterviewProgressHandler{svc: svc}
}

func (h *InterviewProgressHandler) List(c *gin.Context) {
	applicationID := c.Query("applicationId")
	if strings.TrimSpace(applicationID) == "" {
		response.Fail(c, pkgerrors.CodeBadParams, "applicationId 不能为空")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	rows, err := h.svc.List(c.Request.Context(), userID, applicationID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toInterviewProgressResponses(rows)})
}

func (h *InterviewProgressHandler) Create(c *gin.Context) {
	var req dto.CreateInterviewProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewProgressResponse(row))
}

func (h *InterviewProgressHandler) Update(c *gin.Context) {
	var req dto.UpdateInterviewProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	row, err := h.svc.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewProgressResponse(row))
}

func toInterviewProgressResponses(rows []domainprogress.InterviewProgress) []*dto.InterviewProgressResponse {
	responses := make([]*dto.InterviewProgressResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, dto.ToInterviewProgressResponse(&rows[i]))
	}
	return responses
}

func (h *InterviewProgressHandler) Delete(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
