package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type DashboardService interface {
	Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error)
	CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error)
	UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error)
	CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error)
	DeleteTodo(ctx context.Context, userID, id string) error
}

type DashboardHandler struct {
	svc DashboardService
}

func NewDashboardHandler(svc DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, c.Query("range"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}

func (h *DashboardHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateDashboardTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	todo, err := h.svc.CreateTodo(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) UpdateTodo(c *gin.Context) {
	var req dto.UpdateDashboardTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	todo, err := h.svc.UpdateTodo(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) CompleteTodo(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	todo, err := h.svc.CompleteTodo(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) DeleteTodo(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteTodo(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}
