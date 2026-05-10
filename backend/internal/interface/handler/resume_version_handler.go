package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/pagination"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeVersionService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error
	List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error)
	Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error
	Delete(ctx context.Context, userID, id string) error
	Activate(ctx context.Context, userID, id string) error
}

type ResumeVersionHandler struct {
	svc ResumeVersionService
}

func NewResumeVersionHandler(svc ResumeVersionService) *ResumeVersionHandler {
	return &ResumeVersionHandler{svc: svc}
}

func resumeUserID(c *gin.Context) (string, bool) {
	userID := c.GetString("userID")
	if userID == "" {
		mapError(c, pkgerrors.ErrUnauthorized)
		return "", false
	}
	return userID, true
}

func (h *ResumeVersionHandler) Create(c *gin.Context) {
	var req dto.CreateResumeVersionRequest
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
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) List(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	repoID := c.Query("repoId")
	if repoID == "" {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	rows, err := h.svc.List(c.Request.Context(), userID, repoID)
	if err != nil {
		mapError(c, err)
		return
	}
	p := pagination.FromContext(c)
	resp := make([]dto.ResumeVersionResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, dto.ToResumeVersionResponse(row))
	}
	response.SuccessPage(c, resp, int64(len(rows)), p.Page, p.PageSize)
}

func (h *ResumeVersionHandler) Compare(c *gin.Context) {
	response.Fail(c, pkgerrors.CodeBadParams, "暂不支持版本对比")
}

func (h *ResumeVersionHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToResumeVersionResponse(*row))
}

func (h *ResumeVersionHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) Delete(c *gin.Context) {
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

func (h *ResumeVersionHandler) Activate(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Activate(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
