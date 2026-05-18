package activitylog

import (
	"context"
	"strings"

	"github.com/google/uuid"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activitylogrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
)

type Entry struct {
	UserID      string
	Action      domainactivity.Action
	TargetType  domainactivity.TargetType
	TargetID    string
	TargetLabel string
	Summary     string
	Route       string
}

type Logger interface {
	Log(ctx context.Context, entry Entry) error
}

type NoopLogger struct{}

func (NoopLogger) Log(ctx context.Context, entry Entry) error {
	return nil
}

type RepositoryLogger struct {
	repo activitylogrepo.Repository
}

func NewRepositoryLogger(repo activitylogrepo.Repository) *RepositoryLogger {
	return &RepositoryLogger{repo: repo}
}

func (l *RepositoryLogger) Log(ctx context.Context, entry Entry) error {
	if l == nil || l.repo == nil {
		return nil
	}
	activity, err := domainactivity.NewActivityLog(domainactivity.NewActivityLogParams{
		ID:          uuid.NewString(),
		UserID:      entry.UserID,
		Action:      entry.Action,
		TargetType:  entry.TargetType,
		TargetID:    entry.TargetID,
		TargetLabel: entry.TargetLabel,
		Summary:     entry.Summary,
		Route:       entry.Route,
	})
	if err != nil {
		return err
	}
	return l.repo.Create(ctx, activity)
}

func TrimLabel(label string) string {
	label = strings.TrimSpace(label)
	runes := []rune(label)
	if len(runes) > 200 {
		return string(runes[:200])
	}
	return label
}

func TrimSummary(summary string) string {
	summary = strings.TrimSpace(summary)
	runes := []rune(summary)
	if len(runes) > 500 {
		return string(runes[:500])
	}
	return summary
}

func TrimRoute(route string) string {
	route = strings.TrimSpace(route)
	runes := []rune(route)
	if len(runes) > 500 {
		return string(runes[:500])
	}
	return route
}

func OrNoop(logger Logger) Logger {
	if logger == nil {
		return NoopLogger{}
	}
	return logger
}
