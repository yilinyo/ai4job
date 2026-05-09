package resume_version

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersion struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	RepoID     string         `gorm:"type:uuid;not null;index"`
	Title      string         `gorm:"type:varchar(120);not null"`
	Content    string         `gorm:"type:text;not null"`
	VersionNum int            `gorm:"not null"`
	IsDefault  bool           `gorm:"not null;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (v *ResumeVersion) ValidateForCreate() error {
	if strings.TrimSpace(v.RepoID) == "" || strings.TrimSpace(v.Title) == "" || strings.TrimSpace(v.Content) == "" {
		return domainerrors.ErrBadParams
	}
	return nil
}
