package interview_progress

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Stage string

type Result string

const (
	StageApplied         Stage = "applied"
	StageViewed          Stage = "viewed"
	StageWrittenTest     Stage = "written_test"
	StageFirstInterview  Stage = "first_interview"
	StageSecondInterview Stage = "second_interview"
	StageFinalInterview  Stage = "final_interview"
	StageHRInterview     Stage = "hr_interview"
	StageOffer           Stage = "offer"
	StageRejected        Stage = "rejected"
	StageClosed          Stage = "closed"

	ResultPassed  Result = "passed"
	ResultPending Result = "pending"
	ResultFailed  Result = "failed"
)

type InterviewProgress struct {
	ID                string    `gorm:"type:uuid;primaryKey"`
	UserID            string    `gorm:"type:varchar(36);not null;index"`
	ApplicationID     string    `gorm:"type:uuid;not null;index"`
	Stage             Stage     `gorm:"type:varchar(30);not null"`
	OccurredAt        time.Time `gorm:"not null;index"`
	Result            Result    `gorm:"type:varchar(20);not null"`
	InterviewerOrTeam string    `gorm:"type:varchar(120)"`
	Note              string    `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (InterviewProgress) TableName() string {
	return "interview_progress"
}

func (p *InterviewProgress) Validate() error {
	if strings.TrimSpace(p.UserID) == "" || strings.TrimSpace(p.ApplicationID) == "" {
		return domainerrors.ErrBadParams
	}
	if p.OccurredAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if !validStage(p.Stage) || !validResult(p.Result) {
		return domainerrors.ErrBadParams
	}
	if len([]rune(p.InterviewerOrTeam)) > 120 {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validStage(stage Stage) bool {
	switch stage {
	case StageApplied, StageViewed, StageWrittenTest, StageFirstInterview, StageSecondInterview, StageFinalInterview, StageHRInterview, StageOffer, StageRejected, StageClosed:
		return true
	default:
		return false
	}
}

func validResult(result Result) bool {
	switch result {
	case ResultPassed, ResultPending, ResultFailed:
		return true
	default:
		return false
	}
}
