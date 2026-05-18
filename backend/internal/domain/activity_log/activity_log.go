package activity_log

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type Action string

const (
	ActionCreated       Action = "created"
	ActionUpdated       Action = "updated"
	ActionDeleted       Action = "deleted"
	ActionStatusChanged Action = "status_changed"
	ActionCompleted     Action = "completed"
	ActionFavorited     Action = "favorited"
	ActionUnfavorited   Action = "unfavorited"
	ActionLinked        Action = "linked"
)

type TargetType string

const (
	TargetTypeApplication TargetType = "application"
	TargetTypeInterview   TargetType = "interview"
	TargetTypeQuestion    TargetType = "question"
	TargetTypeResume      TargetType = "resume"
)

type ActivityLog struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	UserID      string     `gorm:"column:user_id"`
	ActorType   string     `gorm:"column:actor_type"`
	Action      Action     `gorm:"column:action"`
	TargetType  TargetType `gorm:"column:target_type"`
	TargetID    string     `gorm:"column:target_id"`
	TargetLabel string     `gorm:"column:target_label"`
	Summary     string     `gorm:"column:summary"`
	Route       string     `gorm:"column:route"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
}

type NewActivityLogParams struct {
	ID          string
	UserID      string
	Action      Action
	TargetType  TargetType
	TargetID    string
	TargetLabel string
	Summary     string
	Route       string
}

func NewActivityLog(params NewActivityLogParams) (*ActivityLog, error) {
	log := &ActivityLog{
		ID:          strings.TrimSpace(params.ID),
		UserID:      strings.TrimSpace(params.UserID),
		ActorType:   "user",
		Action:      params.Action,
		TargetType:  params.TargetType,
		TargetID:    strings.TrimSpace(params.TargetID),
		TargetLabel: strings.TrimSpace(params.TargetLabel),
		Summary:     strings.TrimSpace(params.Summary),
		Route:       strings.TrimSpace(params.Route),
	}
	if err := log.validate(); err != nil {
		return nil, err
	}
	return log, nil
}

func (l *ActivityLog) TableName() string {
	return "activity_logs"
}

func (l *ActivityLog) validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.UserID) == "" || strings.TrimSpace(l.ActorType) == "" || strings.TrimSpace(l.TargetID) == "" || strings.TrimSpace(l.TargetLabel) == "" || strings.TrimSpace(l.Summary) == "" {
		return domainerrors.ErrBadParams
	}
	if !l.Action.Valid() || !l.TargetType.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (a Action) Valid() bool {
	switch a {
	case ActionCreated, ActionUpdated, ActionDeleted, ActionStatusChanged, ActionCompleted, ActionFavorited, ActionUnfavorited, ActionLinked:
		return true
	default:
		return false
	}
}

func (t TargetType) Valid() bool {
	switch t {
	case TargetTypeApplication, TargetTypeInterview, TargetTypeQuestion, TargetTypeResume:
		return true
	default:
		return false
	}
}
