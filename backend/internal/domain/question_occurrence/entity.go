package question_occurrence

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type AnswerPerformance string

const (
	AnswerPerformanceFailed   AnswerPerformance = "failed"
	AnswerPerformanceVague    AnswerPerformance = "vague"
	AnswerPerformanceAnswered AnswerPerformance = "answered"
	AnswerPerformanceGood     AnswerPerformance = "good"
)

type QuestionOccurrence struct {
	ID                     string             `gorm:"type:uuid;primaryKey"`
	UserID                 string             `gorm:"type:varchar(36);not null;index"`
	QuestionID             string             `gorm:"type:uuid;not null;index"`
	ApplicationID          string             `gorm:"type:uuid;not null;index"`
	InterviewProgressID    *string            `gorm:"type:uuid;index"`
	CompanyNameSnapshot    string             `gorm:"type:varchar(200);not null"`
	JobTitleSnapshot       string             `gorm:"type:varchar(200);not null"`
	InterviewStageSnapshot string             `gorm:"type:varchar(50)"`
	OccurredAt             time.Time          `gorm:"not null"`
	ActualQuestion         string             `gorm:"type:text"`
	AnswerPerformance      *AnswerPerformance `gorm:"type:varchar(20)"`
	Note                   string             `gorm:"type:text"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type NewQuestionOccurrenceParams struct {
	ID                     string
	UserID                 string
	QuestionID             string
	ApplicationID          string
	InterviewProgressID    *string
	CompanyNameSnapshot    string
	JobTitleSnapshot       string
	InterviewStageSnapshot string
	OccurredAt             time.Time
	ActualQuestion         string
	AnswerPerformance      *AnswerPerformance
	Note                   string
}

func NewQuestionOccurrence(params NewQuestionOccurrenceParams) (*QuestionOccurrence, error) {
	occurrence := &QuestionOccurrence{
		ID:                     params.ID,
		UserID:                 params.UserID,
		QuestionID:             params.QuestionID,
		ApplicationID:          params.ApplicationID,
		InterviewProgressID:    params.InterviewProgressID,
		CompanyNameSnapshot:    params.CompanyNameSnapshot,
		JobTitleSnapshot:       params.JobTitleSnapshot,
		InterviewStageSnapshot: params.InterviewStageSnapshot,
		OccurredAt:             params.OccurredAt,
		ActualQuestion:         params.ActualQuestion,
		AnswerPerformance:      params.AnswerPerformance,
		Note:                   params.Note,
	}
	if err := occurrence.Validate(); err != nil {
		return nil, err
	}
	return occurrence, nil
}

func (o *QuestionOccurrence) Validate() error {
	if strings.TrimSpace(o.ID) == "" || strings.TrimSpace(o.UserID) == "" || strings.TrimSpace(o.QuestionID) == "" || strings.TrimSpace(o.ApplicationID) == "" || strings.TrimSpace(o.CompanyNameSnapshot) == "" || strings.TrimSpace(o.JobTitleSnapshot) == "" {
		return domainerrors.ErrBadParams
	}
	if o.OccurredAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if len([]rune(o.CompanyNameSnapshot)) > 200 || len([]rune(o.JobTitleSnapshot)) > 200 || len([]rune(o.InterviewStageSnapshot)) > 50 {
		return domainerrors.ErrBadParams
	}
	if o.AnswerPerformance != nil && !o.AnswerPerformance.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (p AnswerPerformance) Valid() bool {
	switch p {
	case AnswerPerformanceFailed, AnswerPerformanceVague, AnswerPerformanceAnswered, AnswerPerformanceGood:
		return true
	default:
		return false
	}
}

func (QuestionOccurrence) TableName() string {
	return "question_occurrences"
}
