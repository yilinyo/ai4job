package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

type mockDashboardService struct{ mock.Mock }

func (m *mockDashboardService) Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error) {
	args := m.Called(ctx, userID, rawRange)
	return args.Get(0).(*dto.DashboardResponse), args.Error(1)
}

func (m *mockDashboardService) CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1)
}

func (m *mockDashboardService) UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	args := m.Called(ctx, userID, id, req)
	return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1)
}

func (m *mockDashboardService) CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error) {
	args := m.Called(ctx, userID, id)
	return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1)
}

func (m *mockDashboardService) DeleteTodo(ctx context.Context, userID, id string) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

func TestDashboardHandler_GetUsesRangeAndDTOShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockDashboardService{}
	svc.On("Get", mock.Anything, "u1", "30d").Return(&dto.DashboardResponse{Range: "30d", Metrics: dto.DashboardMetrics{TotalApplications: 2}}, nil).Once()
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	r.GET("/dashboard", func(c *gin.Context) { c.Set("userID", "u1"); h.Get(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dashboard?range=30d", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "30d", resp.Data["range"])
	require.Contains(t, resp.Data, "metrics")
	svc.AssertExpectations(t)
}

func TestDashboardHandler_CreateCompleteUpdateAndDeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockDashboardService{}
	createReq := dto.CreateDashboardTodoRequest{Title: "跟进腾讯"}
	updateTitle := "准备面试"
	updateReq := dto.UpdateDashboardTodoRequest{Title: &updateTitle, ClearDueAt: true}
	todo := &domaintodo.DashboardTodo{ID: "todo-1", UserID: "u1", Title: "跟进腾讯"}
	updatedTodo := &domaintodo.DashboardTodo{ID: "todo-1", UserID: "u1", Title: updateTitle}
	svc.On("CreateTodo", mock.Anything, "u1", createReq).Return(todo, nil).Once()
	svc.On("CompleteTodo", mock.Anything, "u1", "todo-1").Return(todo, nil).Once()
	svc.On("UpdateTodo", mock.Anything, "u1", "todo-1", updateReq).Return(updatedTodo, nil).Once()
	svc.On("DeleteTodo", mock.Anything, "u1", "todo-1").Return(nil).Once()
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	r.POST("/dashboard/todos", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateTodo(c) })
	r.PUT("/dashboard/todos/:id/complete", func(c *gin.Context) { c.Set("userID", "u1"); h.CompleteTodo(c) })
	r.PUT("/dashboard/todos/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.UpdateTodo(c) })
	r.DELETE("/dashboard/todos/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.DeleteTodo(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/dashboard/todos", bytes.NewBufferString(`{"title":"跟进腾讯"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/dashboard/todos/todo-1/complete", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/dashboard/todos/todo-1", bytes.NewBufferString(`{"title":"准备面试","clearDueAt":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/dashboard/todos/todo-1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
