package postgres

import (
	"context"
	"errors"
	"strings"

	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	resumeapplicationrepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeApplicationRepo struct {
	db *gorm.DB
}

var _ resumeapplicationrepo.Repository = (*ResumeApplicationRepo)(nil)

func NewResumeApplicationRepo(db *gorm.DB) *ResumeApplicationRepo {
	return &ResumeApplicationRepo{db: db}
}

func (r *ResumeApplicationRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error {
	err := r.db.WithContext(ctx).Create(a).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *ResumeApplicationRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error) {
	var a domainapp.ResumeApplication
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &a, err
}

func (r *ResumeApplicationRepo) Update(ctx context.Context, a *domainapp.ResumeApplication) error {
	result := r.db.WithContext(ctx).
		Model(&domainapp.ResumeApplication{}).
		Where("id = ? AND user_id = ?", a.ID, a.UserID).
		Updates(map[string]any{
			"job_posting_id":              a.JobPostingID,
			"resume_version_id":           a.ResumeVersionID,
			"repo_id":                     a.RepoID,
			"company_name":                a.CompanyName,
			"company_intro":               a.CompanyIntro,
			"job_title":                   a.JobTitle,
			"job_requirements":            a.JobRequirements,
			"base":                        a.Base,
			"salary_range":                a.SalaryRange,
			"source_platform":             a.SourcePlatform,
			"source_url":                  a.SourceURL,
			"delivery_channel":            a.DeliveryChannel,
			"applied_at":                  a.AppliedAt,
			"current_status":              a.CurrentStatus,
			"interview_summary":           a.InterviewSummary,
			"interview_notes":             a.InterviewNotes,
			"resume_match_score":          a.ResumeMatchScore,
			"interview_performance_score": a.InterviewPerformanceScore,
			"role_interest_score":         a.RoleInterestScore,
			"overall_score":               a.OverallScore,
			"remark":                      a.Remark,
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

func (r *ResumeApplicationRepo) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainapp.ResumeApplication{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *ResumeApplicationRepo) ListByUser(ctx context.Context, f resumeapplicationrepo.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("user_id = ?", f.UserID)
	q = applyResumeApplicationFilters(q, f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []domainapp.ResumeApplication
	err := q.Order("applied_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}

func (r *ResumeApplicationRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("resume_version_id = ?", versionID).Count(&count).Error
	return count > 0, err
}

func (r *ResumeApplicationRepo) UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error {
	result := r.db.WithContext(ctx).
		Model(&domainapp.ResumeApplication{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("current_status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func applyResumeApplicationFilters(q *gorm.DB, f resumeapplicationrepo.ListFilter) *gorm.DB {
	if resumeVersionID := strings.TrimSpace(f.ResumeVersionID); resumeVersionID != "" {
		q = q.Where("resume_version_id = ?", resumeVersionID)
	}
	if repoID := strings.TrimSpace(f.RepoID); repoID != "" {
		q = q.Where("repo_id = ?", repoID)
	}
	if keyword := strings.TrimSpace(f.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("company_name ILIKE ? OR job_title ILIKE ?", like, like)
	}
	if f.CurrentStatus != "" {
		q = q.Where("current_status = ?", f.CurrentStatus)
	}
	return q
}
