package job_posting

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusApplied Status = "applied"
	StatusNotFit  Status = "not_fit"
	StatusClosed  Status = "closed"
)

type JobPosting struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	UserID          string `gorm:"type:uuid;not null;index"`
	CompanyName     string `gorm:"type:varchar(120);not null"`
	CompanyIntro    string `gorm:"type:text"`
	JobTitle        string `gorm:"type:varchar(120);not null"`
	JobRequirements string `gorm:"type:text"`
	Base            string `gorm:"type:varchar(120)"`
	SalaryRange     string `gorm:"type:varchar(120)"`
	SourcePlatform  string `gorm:"type:varchar(120);index"`
	SourceURL       string `gorm:"type:varchar(500)"`
	Status          Status `gorm:"type:varchar(20);not null;index"`
	Remark          string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (p *JobPosting) Validate() error {
	if strings.TrimSpace(p.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(p.CompanyName) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(p.JobTitle) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(p.CompanyName)) > 120 || len([]rune(p.JobTitle)) > 120 {
		return domainerrors.ErrBadParams
	}
	if !validStatus(p.Status) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusPending, StatusApplied, StatusNotFit, StatusClosed:
		return true
	default:
		return false
	}
}
