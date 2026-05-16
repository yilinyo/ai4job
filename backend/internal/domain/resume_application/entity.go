package resume_application

import (
	"net/url"
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Status string

const (
	StatusApplied         Status = "applied"
	StatusViewed          Status = "viewed"
	StatusWrittenTest     Status = "written_test"
	StatusFirstInterview  Status = "first_interview"
	StatusSecondInterview Status = "second_interview"
	StatusFinalInterview  Status = "final_interview"
	StatusHRInterview     Status = "hr_interview"
	StatusOffer           Status = "offer"
	StatusRejected        Status = "rejected"
	StatusClosed          Status = "closed"
)

type ResumeApplication struct {
	ID                        string    `gorm:"type:uuid;primaryKey"`
	UserID                    string    `gorm:"type:varchar(36);not null;index"`
	JobPostingID              *string   `gorm:"type:uuid;index"`
	ResumeVersionID           string    `gorm:"type:uuid;not null;index"`
	RepoID                    string    `gorm:"type:uuid;not null;index"`
	CompanyName               string    `gorm:"type:varchar(120);not null"`
	CompanyIntro              string    `gorm:"type:text"`
	JobTitle                  string    `gorm:"type:varchar(120);not null"`
	JobRequirements           string    `gorm:"type:text"`
	Base                      string    `gorm:"type:varchar(120)"`
	SalaryRange               string    `gorm:"type:varchar(120)"`
	SourcePlatform            string    `gorm:"type:varchar(120)"`
	SourceURL                 string    `gorm:"type:varchar(500)"`
	DeliveryChannel           string    `gorm:"type:varchar(120)"`
	AppliedAt                 time.Time `gorm:"not null;index"`
	CurrentStatus             Status    `gorm:"type:varchar(30);not null;index"`
	InterviewSummary          string    `gorm:"type:text"`
	InterviewNotes            string    `gorm:"type:text"`
	ResumeMatchScore          int
	InterviewPerformanceScore int
	RoleInterestScore         int
	OverallScore              int
	Remark                    string `gorm:"type:text"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 gorm.DeletedAt `gorm:"index"`
}

type NewResumeApplicationParams struct {
	ID                        string
	UserID                    string
	JobPostingID              *string
	ResumeVersionID           string
	RepoID                    string
	CompanyName               string
	CompanyIntro              string
	JobTitle                  string
	JobRequirements           string
	Base                      string
	SalaryRange               string
	SourcePlatform            string
	SourceURL                 string
	DeliveryChannel           string
	AppliedAt                 time.Time
	CurrentStatus             Status
	InterviewSummary          string
	InterviewNotes            string
	ResumeMatchScore          int
	InterviewPerformanceScore int
	RoleInterestScore         int
	OverallScore              int
	Remark                    string
}

func NewResumeApplication(params NewResumeApplicationParams) (*ResumeApplication, error) {
	status := params.CurrentStatus
	if status == "" {
		status = StatusApplied
	}
	application := &ResumeApplication{
		ID:                        params.ID,
		UserID:                    params.UserID,
		JobPostingID:              params.JobPostingID,
		ResumeVersionID:           params.ResumeVersionID,
		RepoID:                    params.RepoID,
		CompanyName:               params.CompanyName,
		CompanyIntro:              params.CompanyIntro,
		JobTitle:                  params.JobTitle,
		JobRequirements:           params.JobRequirements,
		Base:                      params.Base,
		SalaryRange:               params.SalaryRange,
		SourcePlatform:            params.SourcePlatform,
		SourceURL:                 params.SourceURL,
		DeliveryChannel:           params.DeliveryChannel,
		AppliedAt:                 params.AppliedAt,
		CurrentStatus:             status,
		InterviewSummary:          params.InterviewSummary,
		InterviewNotes:            params.InterviewNotes,
		ResumeMatchScore:          params.ResumeMatchScore,
		InterviewPerformanceScore: params.InterviewPerformanceScore,
		RoleInterestScore:         params.RoleInterestScore,
		OverallScore:              params.OverallScore,
		Remark:                    params.Remark,
	}
	if err := application.Validate(); err != nil {
		return nil, err
	}
	return application, nil
}

type UpdateSpec struct {
	JobPostingID              **string
	ResumeVersionID           *string
	RepoID                    *string
	CompanyName               *string
	CompanyIntro              *string
	JobTitle                  *string
	JobRequirements           *string
	Base                      *string
	SalaryRange               *string
	SourcePlatform            *string
	SourceURL                 *string
	DeliveryChannel           *string
	AppliedAt                 *time.Time
	CurrentStatus             *Status
	InterviewSummary          *string
	InterviewNotes            *string
	ResumeMatchScore          *int
	InterviewPerformanceScore *int
	RoleInterestScore         *int
	OverallScore              *int
	Remark                    *string
}

func (a *ResumeApplication) ApplyUpdate(update UpdateSpec) error {
	if update.JobPostingID != nil {
		a.JobPostingID = *update.JobPostingID
	}
	if update.ResumeVersionID != nil {
		a.ResumeVersionID = *update.ResumeVersionID
	}
	if update.RepoID != nil {
		a.RepoID = *update.RepoID
	}
	if update.CompanyName != nil {
		a.CompanyName = *update.CompanyName
	}
	if update.CompanyIntro != nil {
		a.CompanyIntro = *update.CompanyIntro
	}
	if update.JobTitle != nil {
		a.JobTitle = *update.JobTitle
	}
	if update.JobRequirements != nil {
		a.JobRequirements = *update.JobRequirements
	}
	if update.Base != nil {
		a.Base = *update.Base
	}
	if update.SalaryRange != nil {
		a.SalaryRange = *update.SalaryRange
	}
	if update.SourcePlatform != nil {
		a.SourcePlatform = *update.SourcePlatform
	}
	if update.SourceURL != nil {
		a.SourceURL = *update.SourceURL
	}
	if update.DeliveryChannel != nil {
		a.DeliveryChannel = *update.DeliveryChannel
	}
	if update.AppliedAt != nil {
		a.AppliedAt = *update.AppliedAt
	}
	if update.CurrentStatus != nil {
		a.CurrentStatus = *update.CurrentStatus
	}
	if update.InterviewSummary != nil {
		a.InterviewSummary = *update.InterviewSummary
	}
	if update.InterviewNotes != nil {
		a.InterviewNotes = *update.InterviewNotes
	}
	if update.ResumeMatchScore != nil {
		a.ResumeMatchScore = *update.ResumeMatchScore
	}
	if update.InterviewPerformanceScore != nil {
		a.InterviewPerformanceScore = *update.InterviewPerformanceScore
	}
	if update.RoleInterestScore != nil {
		a.RoleInterestScore = *update.RoleInterestScore
	}
	if update.OverallScore != nil {
		a.OverallScore = *update.OverallScore
	}
	if update.Remark != nil {
		a.Remark = *update.Remark
	}
	return a.Validate()
}

func (a *ResumeApplication) UpdateCurrentStatus(status Status) error {
	oldStatus := a.CurrentStatus
	a.CurrentStatus = status
	if err := a.Validate(); err != nil {
		a.CurrentStatus = oldStatus
		return err
	}
	return nil
}

func (s Status) Valid() bool {
	return validStatus(s)
}

func (a *ResumeApplication) Validate() error {
	if strings.TrimSpace(a.UserID) == "" || strings.TrimSpace(a.ResumeVersionID) == "" || strings.TrimSpace(a.RepoID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(a.CompanyName) == "" || strings.TrimSpace(a.JobTitle) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(a.CompanyName)) > 120 || len([]rune(a.JobTitle)) > 120 {
		return domainerrors.ErrBadParams
	}
	if len([]rune(a.Base)) > 120 || len([]rune(a.SalaryRange)) > 120 || len([]rune(a.SourcePlatform)) > 120 || len([]rune(a.SourceURL)) > 500 || len([]rune(a.DeliveryChannel)) > 120 {
		return domainerrors.ErrBadParams
	}
	if !validSourceURL(a.SourceURL) {
		return domainerrors.ErrBadParams
	}
	if a.AppliedAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if !validStatus(a.CurrentStatus) {
		return domainerrors.ErrBadParams
	}
	if !validOptionalScore(a.ResumeMatchScore, 1, 5) || !validOptionalScore(a.InterviewPerformanceScore, 1, 5) || !validOptionalScore(a.RoleInterestScore, 1, 5) || !validOptionalScore(a.OverallScore, 1, 10) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validSourceURL(sourceURL string) bool {
	if sourceURL == "" {
		return true
	}

	parsed, err := url.Parse(sourceURL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
}

func validStatus(status Status) bool {
	switch status {
	case StatusApplied, StatusViewed, StatusWrittenTest, StatusFirstInterview, StatusSecondInterview, StatusFinalInterview, StatusHRInterview, StatusOffer, StatusRejected, StatusClosed:
		return true
	default:
		return false
	}
}

func validOptionalScore(score, min, max int) bool {
	return score == 0 || (score >= min && score <= max)
}
