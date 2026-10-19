package repository

import (
	"context"
	"time"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
)

type ListFilter struct {
	UserID  string
	StartAt *time.Time
	Limit   int
	Offset  int
}

type Repository interface {
	Create(ctx context.Context, activity *domainactivity.ActivityLog) error
	ListByUser(ctx context.Context, filter ListFilter) ([]domainactivity.ActivityLog, int64, error)
}
