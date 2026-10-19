package postgres

import (
	"context"
	"errors"

	attachment "github.com/yilin/ai-for-backend/internal/domain/resume_attachment"
	attachmentrepo "github.com/yilin/ai-for-backend/internal/domain/resume_attachment/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeAttachmentRepo struct{ db *gorm.DB }

var _ attachmentrepo.Repository = (*ResumeAttachmentRepo)(nil)

func NewResumeAttachmentRepo(db *gorm.DB) *ResumeAttachmentRepo { return &ResumeAttachmentRepo{db: db} }

func (r *ResumeAttachmentRepo) Create(ctx context.Context, value *attachment.ResumeAttachment) error {
	return r.db.WithContext(ctx).Create(value).Error
}

func (r *ResumeAttachmentRepo) ListByVersionID(ctx context.Context, versionID string) ([]attachment.ResumeAttachment, error) {
	var rows []attachment.ResumeAttachment
	err := r.db.WithContext(ctx).Where("version_id = ?", versionID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *ResumeAttachmentRepo) FindByIDAndVersionID(ctx context.Context, id, versionID string) (*attachment.ResumeAttachment, error) {
	var row attachment.ResumeAttachment
	err := r.db.WithContext(ctx).Where("id = ? AND version_id = ?", id, versionID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &row, err
}

func (r *ResumeAttachmentRepo) Delete(ctx context.Context, id, versionID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND version_id = ?", id, versionID).Delete(&attachment.ResumeAttachment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}
