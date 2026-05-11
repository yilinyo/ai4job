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
	Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) error
	List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) error
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
	response.Success(c, gin.H{"list": rows})
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
	if err := h.svc.Create(c.Request.Context(), userID, req); err != nil {
		mapError(c, err)
		return
	}
	rows, err := h.svc.List(c.Request.Context(), userID, req.ApplicationID)
	if err != nil {
		mapError(c, err)
		return
	}
	if len(rows) == 0 {
		response.Success(c, gin.H{})
		return
	}
	response.Success(c, rows[len(rows)-1])
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
	if err := h.svc.Update(c.Request.Context(), userID, id, req); err != nil {
		mapError(c, err)
		return
	}
	applicationID := c.Query("applicationId")
	if req.ApplicationID != nil {
		applicationID = *req.ApplicationID
	}
	if strings.TrimSpace(applicationID) == "" {
		response.Success(c, gin.H{})
		return
	}
	rows, err := h.svc.List(c.Request.Context(), userID, applicationID)
	if err != nil {
		mapError(c, err)
		return
	}
	for i := range rows {
		if rows[i].ID == id {
			response.Success(c, rows[i])
			return
		}
	}
	response.Success(c, gin.H{})
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
