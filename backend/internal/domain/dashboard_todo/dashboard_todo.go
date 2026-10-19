package dashboard_todo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type Source string

const (
	// Source is response metadata only: manual todos are persisted, system todos are generated.
	// It is intentionally not stored on DashboardTodo.
	SourceManual Source = "manual"
	SourceSystem Source = "system"
)

const maxTitleLength = 200

type TargetType string

const (
	TargetTypeApplication TargetType = "application"
	TargetTypeInterview   TargetType = "interview"
	TargetTypeQuestion    TargetType = "question"
	TargetTypeResume      TargetType = "resume"
)

type DashboardTodo struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	UserID      string     `gorm:"column:user_id"`
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	DueAt       *time.Time `gorm:"column:due_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	TargetType  TargetType `gorm:"column:target_type"`
	TargetID    *string    `gorm:"column:target_id"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

type NewDashboardTodoParams struct {
	ID          string
	UserID      string
	Title       string
	Description string
	DueAt       *time.Time
	TargetType  TargetType
	TargetID    *string
}

type UpdateSpec struct {
	Title       *string
	Description *string
	DueAt       **time.Time
	TargetType  *TargetType
	TargetID    **string
}

func NewDashboardTodo(params NewDashboardTodoParams) (*DashboardTodo, error) {
	todo := &DashboardTodo{
		ID:          strings.TrimSpace(params.ID),
		UserID:      strings.TrimSpace(params.UserID),
		Title:       strings.TrimSpace(params.Title),
		Description: strings.TrimSpace(params.Description),
		DueAt:       params.DueAt,
		TargetType:  params.TargetType,
		TargetID:    normalizeTargetID(params.TargetID),
	}
	if err := todo.validate(); err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *DashboardTodo) TableName() string {
	return "dashboard_todos"
}

func (t *DashboardTodo) ApplyUpdate(spec UpdateSpec) error {
	original := *t

	if spec.Title != nil {
		t.Title = strings.TrimSpace(*spec.Title)
	}
	if spec.Description != nil {
		t.Description = strings.TrimSpace(*spec.Description)
	}
	if spec.DueAt != nil {
		t.DueAt = *spec.DueAt
	}
	if spec.TargetType != nil {
		t.TargetType = *spec.TargetType
	}
	if spec.TargetID != nil {
		t.TargetID = normalizeTargetID(*spec.TargetID)
	}
	if err := t.validate(); err != nil {
		*t = original
		return err
	}
	return nil
}

func (t *DashboardTodo) Complete(now time.Time) {
	t.CompletedAt = &now
}

func (t *DashboardTodo) validate() error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.UserID) == "" || strings.TrimSpace(t.Title) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(t.Title)) > maxTitleLength {
		return domainerrors.ErrBadParams
	}
	if t.TargetType != "" && !t.TargetType.Valid() {
		return domainerrors.ErrBadParams
	}
	if t.TargetType == "" && t.TargetID != nil {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (t TargetType) Valid() bool {
	switch t {
	case TargetTypeApplication, TargetTypeInterview, TargetTypeQuestion, TargetTypeResume:
		return true
	default:
		return false
	}
}

func normalizeTargetID(targetID *string) *string {
	if targetID == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*targetID)
	return &trimmed
}
