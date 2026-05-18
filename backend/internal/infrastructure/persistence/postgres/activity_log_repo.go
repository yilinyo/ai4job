package postgres

import (
	"context"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activitylogrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	"gorm.io/gorm"
)

type ActivityLogRepo struct {
	db *gorm.DB
}

var _ activitylogrepo.Repository = (*ActivityLogRepo)(nil)

func NewActivityLogRepo(db *gorm.DB) *ActivityLogRepo {
	return &ActivityLogRepo{db: db}
}

func (r *ActivityLogRepo) Create(ctx context.Context, activity *domainactivity.ActivityLog) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

func (r *ActivityLogRepo) ListByUser(ctx context.Context, filter activitylogrepo.ListFilter) ([]domainactivity.ActivityLog, error) {
	q := r.db.WithContext(ctx).Model(&domainactivity.ActivityLog{}).Where("user_id = ?", filter.UserID)
	if filter.StartAt != nil {
		q = q.Where("created_at >= ?", *filter.StartAt)
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	}

	var rows []domainactivity.ActivityLog
	err := q.Order("created_at DESC").Find(&rows).Error
	return rows, err
}
