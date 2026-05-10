package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/pagination"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeRepoService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error
	List(ctx context.Context, userID string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error)
	Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type ResumeRepoHandler struct {
	svc ResumeRepoService
}

func NewResumeRepoHandler(svc ResumeRepoService) *ResumeRepoHandler {
	return &ResumeRepoHandler{svc: svc}
}

func (h *ResumeRepoHandler) Create(c *gin.Context) {
	var req dto.CreateResumeRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	if err := h.svc.Create(c.Request.Context(), c.GetString("userID"), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeRepoHandler) List(c *gin.Context) {
	p := pagination.FromContext(c)
	rows, total, err := h.svc.List(c.Request.Context(), c.GetString("userID"), p.Page, p.PageSize)
	if err != nil {
		mapError(c, err)
		return
	}
	response.SuccessPage(c, rows, total, p.Page, p.PageSize)
}

func (h *ResumeRepoHandler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.GetString("userID"), c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}

func (h *ResumeRepoHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	if err := h.svc.Update(c.Request.Context(), c.GetString("userID"), c.Param("id"), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeRepoHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.GetString("userID"), c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
