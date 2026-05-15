package resume_version

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersion struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	RepoID    string `gorm:"type:uuid;not null;index"`
	Version   string `gorm:"type:varchar(20);not null"`
	Title     string `gorm:"type:varchar(120);not null"`
	Content   string `gorm:"type:text;not null"`
	Remark    string `gorm:"type:text"`
	IsActive  bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NewResumeVersionParams struct {
	ID       string
	RepoID   string
	Version  string
	Title    string
	Content  string
	Remark   string
	IsActive bool
}

func NewResumeVersion(params NewResumeVersionParams) (*ResumeVersion, error) {
	version := &ResumeVersion{
		ID:       params.ID,
		RepoID:   params.RepoID,
		Version:  params.Version,
		Title:    params.Title,
		Content:  params.Content,
		Remark:   params.Remark,
		IsActive: params.IsActive,
	}
	if err := version.ValidateForCreate(); err != nil {
		return nil, err
	}
	return version, nil
}

type UpdateSpec struct {
	Title    string
	Content  string
	Remark   string
	IsActive *bool
}

func (v *ResumeVersion) ApplyUpdate(update UpdateSpec) error {
	v.Title = update.Title
	v.Content = update.Content
	v.Remark = update.Remark
	if update.IsActive != nil && !*update.IsActive {
		v.IsActive = false
	}
	return v.ValidateForCreate()
}

func (v *ResumeVersion) ValidateForCreate() error {
	if strings.TrimSpace(v.RepoID) == "" || strings.TrimSpace(v.Version) == "" || strings.TrimSpace(v.Title) == "" || strings.TrimSpace(v.Content) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(v.Version)) > 20 || len([]rune(v.Title)) > 120 {
		return domainerrors.ErrBadParams
	}
	return nil
}
