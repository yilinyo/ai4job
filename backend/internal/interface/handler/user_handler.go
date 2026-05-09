package handler

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

// UserService defines the use cases needed by the user handler.
type UserService interface {
	SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error)
}

// UserHandler handles HTTP requests for the user resource.
type UserHandler struct {
	svc UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// SendEmailCode handles POST /user/email-code
func (h *UserHandler) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	resp, err := h.svc.SendEmailCode(c.Request.Context(), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, resp)
}

// Register handles POST /user/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	if err := h.svc.Register(c.Request.Context(), req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

// Login handles POST /user/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, resp)
}

// GetProfile handles GET /user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	resp, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, resp)
}

// UpdateProfile handles PUT /user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	resp, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, resp)
}

// Logout handles POST /user/logout
// Full blacklist logic is wired in main.go where JWTParser is available.
func (h *UserHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// mapError translates domain errors into HTTP responses.
func mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pkgerrors.ErrNotFound):
		response.Fail(c, pkgerrors.CodeNotFound, "资源不存在")
	case errors.Is(err, pkgerrors.ErrUnauthorized):
		response.Fail(c, pkgerrors.CodeUnauthorized, "未登录")
	case errors.Is(err, pkgerrors.ErrDuplicate), errors.Is(err, pkgerrors.ErrConflict):
		response.Fail(c, pkgerrors.CodeConflict, "资源已存在")
	case errors.Is(err, pkgerrors.ErrEmailCodeWrong):
		response.Fail(c, pkgerrors.CodeConflict, "验证码错误")
	default:
		response.Fail(c, pkgerrors.CodeInternalError, "服务器错误")
	}
}
