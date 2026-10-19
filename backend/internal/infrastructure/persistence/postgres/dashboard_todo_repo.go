package postgres

import (
	"context"
	"errors"
	"time"

	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	dashboardtodorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type DashboardTodoRepo struct {
	db *gorm.DB
}

var _ dashboardtodorepo.Repository = (*DashboardTodoRepo)(nil)

func NewDashboardTodoRepo(db *gorm.DB) *DashboardTodoRepo {
	return &DashboardTodoRepo{db: db}
}

func (r *DashboardTodoRepo) Create(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	err := dbFromContext(ctx, r.db).WithContext(ctx).Create(todo).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *DashboardTodoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domaintodo.DashboardTodo, error) {
	var todo domaintodo.DashboardTodo
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &todo, err
}

func (r *DashboardTodoRepo) Update(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&domaintodo.DashboardTodo{}).
		Where("id = ? AND user_id = ?", todo.ID, todo.UserID).
		Updates(map[string]any{
			"title":        todo.Title,
			"description":  todo.Description,
			"due_at":       todo.DueAt,
			"completed_at": todo.CompletedAt,
			"target_type":  todo.TargetType,
			"target_id":    todo.TargetID,
		})
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return domainerrors.ErrDuplicate
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *DashboardTodoRepo) Delete(ctx context.Context, id, userID string) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domaintodo.DashboardTodo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *DashboardTodoRepo) ListByUser(ctx context.Context, filter dashboardtodorepo.ListFilter) ([]domaintodo.DashboardTodo, error) {
	q := dbFromContext(ctx, r.db).WithContext(ctx).Model(&domaintodo.DashboardTodo{}).Where("user_id = ?", filter.UserID)
	if filter.Completed != nil {
		if *filter.Completed {
			q = q.Where("completed_at IS NOT NULL")
		} else {
			q = q.Where("completed_at IS NULL")
		}
	}
	if filter.TodayOnly {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)
		q = q.Where("(due_at IS NULL OR due_at < ?) AND (completed_at IS NULL OR completed_at >= ?)", end, start)
	}

	var rows []domaintodo.DashboardTodo
	order := "completed_at ASC NULLS FIRST, due_at ASC NULLS LAST, created_at DESC"
	if filter.CompletedFirst {
		order = "completed_at DESC NULLS LAST, due_at ASC NULLS LAST, created_at DESC"
	}
	err := q.Order(order).Find(&rows).Error
	return rows, err
}
