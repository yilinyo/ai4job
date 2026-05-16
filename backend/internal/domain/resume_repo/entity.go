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

type NewResumeRepoParams struct {
	ID             string
	UserID         string
	Name           string
	JobType        JobType
	TargetPosition string
	Description    string
}

func NewResumeRepo(params NewResumeRepoParams) (*ResumeRepo, error) {
	repo := &ResumeRepo{
		ID:             params.ID,
		UserID:         params.UserID,
		Name:           params.Name,
		JobType:        params.JobType,
		TargetPosition: params.TargetPosition,
		Description:    params.Description,
		VersionCount:   0,
	}
	if err := repo.Validate(); err != nil {
		return nil, err
	}
	return repo, nil
}

type UpdateSpec struct {
	Name           string
	JobType        JobType
	TargetPosition string
	Description    string
}

func (r *ResumeRepo) ApplyUpdate(update UpdateSpec) error {
	r.Name = update.Name
	r.JobType = update.JobType
	r.TargetPosition = update.TargetPosition
	r.Description = update.Description
	return r.Validate()
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
