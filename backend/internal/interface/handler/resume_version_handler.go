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
	Create(ctx context.Context, userID, repoID string, req dto.CreateResumeVersionRequest) error
	List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error)
	Get(ctx context.Context, userID, repoID, id string) (*domainversion.ResumeVersion, error)
	Update(ctx context.Context, userID, repoID, id string, req dto.UpdateResumeVersionRequest) error
	Delete(ctx context.Context, userID, repoID, id string) error
	SetDefault(ctx context.Context, userID, repoID, id string) error
}

type ResumeVersionHandler struct {
	svc ResumeVersionService
}

func NewResumeVersionHandler(svc ResumeVersionService) *ResumeVersionHandler {
	return &ResumeVersionHandler{svc: svc}
}

func resumeRepoParam(c *gin.Context) string {
	if repoID := c.Param("repoId"); repoID != "" {
		return repoID
	}
	return c.Param("id")
}

func resumeVersionParam(c *gin.Context) string {
	if versionID := c.Param("versionId"); versionID != "" {
		return versionID
	}
	return c.Param("id")
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
	if err := h.svc.Create(c.Request.Context(), userID, resumeRepoParam(c), req); err != nil {
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
	rows, err := h.svc.List(c.Request.Context(), userID, resumeRepoParam(c))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, rows)
}

func (h *ResumeVersionHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, resumeRepoParam(c), resumeVersionParam(c))
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
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Update(c.Request.Context(), userID, resumeRepoParam(c), resumeVersionParam(c), req); err != nil {
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
	if err := h.svc.Delete(c.Request.Context(), userID, resumeRepoParam(c), resumeVersionParam(c)); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) SetDefault(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.SetDefault(c.Request.Context(), userID, resumeRepoParam(c), resumeVersionParam(c)); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
