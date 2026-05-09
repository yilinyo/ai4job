package resume_repo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeRepo struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	UserID      string         `gorm:"type:uuid;not null;index"`
	Name        string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (r *ResumeRepo) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.Name) == "" {
		return domainerrors.ErrBadParams
	}
	return nil
}
