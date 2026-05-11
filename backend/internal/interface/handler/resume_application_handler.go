package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeApplicationService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) (*domainapp.ResumeApplication, error)
	List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error)
	Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type ResumeApplicationHandler struct {
	svc ResumeApplicationService
}

func NewResumeApplicationHandler(svc ResumeApplicationService) *ResumeApplicationHandler {
	return &ResumeApplicationHandler{svc: svc}
}

func (h *ResumeApplicationHandler) List(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	page, pageSize := frontendPage(c)
	rows, total, err := h.svc.List(c.Request.Context(), userID, c.Query("resumeVersionId"), c.Query("repoId"), c.Query("keyword"), c.Query("currentStatus"), page, pageSize)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toResumeApplicationResponses(rows), "total": total, "page": page, "pageSize": pageSize})
}

func (h *ResumeApplicationHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToResumeApplicationResponse(row))
}

func (h *ResumeApplicationHandler) Create(c *gin.Context) {
	var req dto.CreateResumeApplicationRequest
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
	response.Success(c, dto.ToResumeApplicationResponse(row))
}

func (h *ResumeApplicationHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeApplicationRequest
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
	row, err := h.svc.Get(c.Request.Context(), userID, id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToResumeApplicationResponse(row))
}

func toResumeApplicationResponses(rows []domainapp.ResumeApplication) []*dto.ResumeApplicationResponse {
	responses := make([]*dto.ResumeApplicationResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, dto.ToResumeApplicationResponse(&rows[i]))
	}
	return responses
}

func (h *ResumeApplicationHandler) Delete(c *gin.Context) {
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
