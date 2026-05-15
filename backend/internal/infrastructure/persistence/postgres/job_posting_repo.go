package postgres

import (
	"context"
	"errors"
	"strings"

	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type JobPostingRepo struct {
	db *gorm.DB
}

func NewJobPostingRepo(db *gorm.DB) *JobPostingRepo {
	return &JobPostingRepo{db: db}
}

func (r *JobPostingRepo) Create(ctx context.Context, p *domainjob.JobPosting) error {
	err := r.db.WithContext(ctx).Create(p).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *JobPostingRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) {
	var p domainjob.JobPosting
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &p, err
}

func (r *JobPostingRepo) Update(ctx context.Context, p *domainjob.JobPosting) error {
	result := r.db.WithContext(ctx).
		Model(&domainjob.JobPosting{}).
		Where("id = ? AND user_id = ?", p.ID, p.UserID).
		Updates(map[string]any{
			"company_name":     p.CompanyName,
			"company_intro":    p.CompanyIntro,
			"job_title":        p.JobTitle,
			"job_requirements": p.JobRequirements,
			"base":             p.Base,
			"salary_range":     p.SalaryRange,
			"source_platform":  p.SourcePlatform,
			"source_url":       p.SourceURL,
			"status":           p.Status,
			"remark":           p.Remark,
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

func (r *JobPostingRepo) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainjob.JobPosting{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *JobPostingRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainjob.JobPosting{}).Where("user_id = ?", f.UserID)
	q = applyJobPostingFilters(q, f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []domainjob.JobPosting
	err := q.Order("created_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}

func applyJobPostingFilters(q *gorm.DB, f domainjob.ListFilter) *gorm.DB {
	if keyword := strings.TrimSpace(f.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("company_name ILIKE ? OR job_title ILIKE ?", like, like)
	}
	if sourcePlatform := strings.TrimSpace(f.SourcePlatform); sourcePlatform != "" {
		q = q.Where("source_platform = ?", sourcePlatform)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	return q
}
