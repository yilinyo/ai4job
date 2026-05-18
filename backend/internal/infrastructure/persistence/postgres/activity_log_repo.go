package postgres

import (
	"context"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activitylogrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
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
	err := dbFromContext(ctx, r.db).WithContext(ctx).Create(activity).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *ActivityLogRepo) ListByUser(ctx context.Context, filter activitylogrepo.ListFilter) ([]domainactivity.ActivityLog, error) {
	q := dbFromContext(ctx, r.db).WithContext(ctx).Model(&domainactivity.ActivityLog{}).Where("user_id = ?", filter.UserID)
	if filter.StartAt != nil {
		q = q.Where("created_at >= ?", *filter.StartAt)
	}
	limit := normalizeActivityLogLimit(filter.Limit)

	var rows []domainactivity.ActivityLog
	err := q.Order("created_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func normalizeActivityLogLimit(limit int) int {
	if limit <= 0 || limit > 50 {
		return 10
	}
	return limit
}
