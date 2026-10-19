package repository

import (
	"context"

	attachment "github.com/yilin/ai-for-backend/internal/domain/resume_attachment"
)

type Repository interface {
	Create(ctx context.Context, value *attachment.ResumeAttachment) error
	ListByVersionID(ctx context.Context, versionID string) ([]attachment.ResumeAttachment, error)
	FindByIDAndVersionID(ctx context.Context, id, versionID string) (*attachment.ResumeAttachment, error)
	Delete(ctx context.Context, id, versionID string) error
}
