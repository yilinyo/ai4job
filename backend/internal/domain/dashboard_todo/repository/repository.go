package repository

import (
	"context"

	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
)

type ListFilter struct {
	UserID         string
	Completed      *bool
	TodayOnly      bool
	CompletedFirst bool
}

type Repository interface {
	Create(ctx context.Context, todo *domaintodo.DashboardTodo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domaintodo.DashboardTodo, error)
	Update(ctx context.Context, todo *domaintodo.DashboardTodo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, filter ListFilter) ([]domaintodo.DashboardTodo, error)
}
