package resume_repo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeCampus     JobType = "campus"
	JobTypeSocial     JobType = "social"
	JobTypeInternship JobType = "internship"
)

type ResumeRepo struct {
	ID             string  `gorm:"type:uuid;primaryKey"`
	UserID         string  `gorm:"type:uuid;not null;index"`
	Name           string  `gorm:"type:varchar(100);not null"`
	JobType        JobType `gorm:"type:varchar(20);not null"`
	TargetPosition string  `gorm:"type:varchar(100);not null"`
	Description    string  `gorm:"type:text"`
	VersionCount   int     `gorm:"not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (r *ResumeRepo) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.Name) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(r.Name)) > 100 {
		return domainerrors.ErrBadParams
	}
	if !r.JobType.Valid() {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.TargetPosition) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(r.TargetPosition)) > 100 {
		return domainerrors.ErrBadParams
	}
	if r.VersionCount < 0 {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (j JobType) Valid() bool {
	switch j {
	case JobTypeCampus, JobTypeSocial, JobTypeInternship:
		return true
	default:
		return false
	}
}
