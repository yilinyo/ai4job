package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeVersionService interface {
	Create(ctx context.Context, repoID string, req dto.CreateResumeVersionRequest) error
	List(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error)
	Get(ctx context.Context, repoID, id string) (*domainversion.ResumeVersion, error)
	Update(ctx context.Context, repoID, id string, req dto.UpdateResumeVersionRequest) error
	Delete(ctx context.Context, repoID, id string) error
	SetDefault(ctx context.Context, repoID, id string) error
}

type ResumeVersionHandler struct {
	svc ResumeVersionService
}

func NewResumeVersionHandler(svc ResumeVersionService) *ResumeVersionHandler {
	return &ResumeVersionHandler{svc: svc}
}

func (h *ResumeVersionHandler) Create(c *gin.Context) {
	var req dto.CreateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	if err := h.svc.Create(c.Request.Context(), c.Param("repoId"), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context(), c.Param("repoId"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, rows)
}

func (h *ResumeVersionHandler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.Param("repoId"), c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}

func (h *ResumeVersionHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	if err := h.svc.Update(c.Request.Context(), c.Param("repoId"), c.Param("id"), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("repoId"), c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) SetDefault(c *gin.Context) {
	if err := h.svc.SetDefault(c.Request.Context(), c.Param("repoId"), c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
