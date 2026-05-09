package postgres

import (
	"context"
	"errors"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeRepoRepo struct {
	db *gorm.DB
}

func NewResumeRepoRepo(db *gorm.DB) *ResumeRepoRepo {
	return &ResumeRepoRepo{db: db}
}

func (r *ResumeRepoRepo) Create(ctx context.Context, rr *domainrepo.ResumeRepo) error {
	return r.db.WithContext(ctx).Create(rr).Error
}

func (r *ResumeRepoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) {
	var rr domainrepo.ResumeRepo
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&rr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &rr, err
}

func (r *ResumeRepoRepo) Update(ctx context.Context, rr *domainrepo.ResumeRepo) error {
	return r.db.WithContext(ctx).Save(rr).Error
}

func (r *ResumeRepoRepo) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainrepo.ResumeRepo{}).Error
}

func (r *ResumeRepoRepo) ListByUser(ctx context.Context, f domainrepo.ListFilter) ([]domainrepo.ResumeRepo, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainrepo.ResumeRepo{}).Where("user_id = ?", f.UserID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []domainrepo.ResumeRepo
	err := q.Order("updated_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}

func (r *ResumeRepoRepo) ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&domainrepo.ResumeRepo{}).Where("user_id = ? AND name = ?", userID, name)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	err := q.Count(&count).Error
	return count > 0, err
}
