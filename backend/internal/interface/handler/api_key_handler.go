package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type APIKeyService interface {
	List(context.Context, string) ([]dto.APIKeyResponse, error)
	Create(context.Context, string, dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error)
	Revoke(context.Context, string, string) error
	Rotate(context.Context, string, string, dto.CreateAPIKeyRequest) (*dto.CreatedAPIKeyResponse, error)
}
type APIKeyHandler struct{ svc APIKeyService }

func NewAPIKeyHandler(s APIKeyService) *APIKeyHandler { return &APIKeyHandler{svc: s} }
func (h *APIKeyHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": rows})
}
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req dto.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	row, err := h.svc.Create(c.Request.Context(), c.GetString("userID"), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	if err := h.svc.Revoke(c.Request.Context(), c.GetString("userID"), c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
func (h *APIKeyHandler) Rotate(c *gin.Context) {
	var req dto.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	row, err := h.svc.Rotate(c.Request.Context(), c.GetString("userID"), c.Param("id"), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}
