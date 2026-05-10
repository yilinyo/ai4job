package postgres

import (
	"context"
	"errors"

	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersionRepo struct {
	db *gorm.DB
}

func NewResumeVersionRepo(db *gorm.DB) *ResumeVersionRepo {
	return &ResumeVersionRepo{db: db}
}

func (r *ResumeVersionRepo) Create(ctx context.Context, v *domainversion.ResumeVersion) error {
	err := r.db.WithContext(ctx).Create(v).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *ResumeVersionRepo) FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error) {
	var v domainversion.ResumeVersion
	err := r.db.WithContext(ctx).Where("id = ? AND repo_id = ?", id, repoID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &v, err
}

func (r *ResumeVersionRepo) ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) {
	var rows []domainversion.ResumeVersion
	err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Order("version_num DESC").Find(&rows).Error
	return rows, err
}

func (r *ResumeVersionRepo) Update(ctx context.Context, v *domainversion.ResumeVersion) error {
	result := r.db.WithContext(ctx).
		Model(&domainversion.ResumeVersion{}).
		Where("id = ? AND repo_id = ?", v.ID, v.RepoID).
		Updates(map[string]any{
			"title":       v.Title,
			"content":     v.Content,
			"version_num": v.VersionNum,
			"is_default":  v.IsDefault,
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

func (r *ResumeVersionRepo) Delete(ctx context.Context, id, repoID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND repo_id = ?", id, repoID).Delete(&domainversion.ResumeVersion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *ResumeVersionRepo) MaxVersionNum(ctx context.Context, repoID string) (int, error) {
	var max int
	err := r.db.WithContext(ctx).Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Select("COALESCE(MAX(version_num), 0)").Scan(&max).Error
	return max, err
}

func (r *ResumeVersionRepo) ClearDefaultByRepoID(ctx context.Context, repoID string) error {
	err := r.db.WithContext(ctx).Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Update("is_default", false).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *ResumeVersionRepo) SetDefaultByID(ctx context.Context, repoID, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&domainversion.ResumeVersion{}).Where("id = ? AND repo_id = ?", id, repoID).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return domainerrors.ErrNotFound
		}
		if err := tx.Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Update("is_default", false).Error; err != nil {
			if isUniqueViolation(err) {
				return domainerrors.ErrDuplicate
			}
			return err
		}
		result := tx.Model(&domainversion.ResumeVersion{}).Where("id = ? AND repo_id = ?", id, repoID).Update("is_default", true)
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
	})
}
