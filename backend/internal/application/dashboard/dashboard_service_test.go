package dashboard

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activityrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	todorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	apprepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestDashboardService_GetAggregatesMetricsAndDistributions(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.apps.rows = []domainapp.ResumeApplication{
		{ID: "app-1", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", SourcePlatform: "Boss", AppliedAt: now.AddDate(0, 0, -1), CurrentStatus: domainapp.StatusApplied},
		{ID: "app-2", UserID: "u1", CompanyName: "腾讯", JobTitle: "后端工程师", SourcePlatform: "官网", AppliedAt: now.AddDate(0, 0, -2), CurrentStatus: domainapp.StatusFirstInterview},
		{ID: "app-3", UserID: "u1", CompanyName: "字节", JobTitle: "AI 工程师", SourcePlatform: "Boss", AppliedAt: now.AddDate(0, 0, -40), CurrentStatus: domainapp.StatusOffer},
	}

	resp, err := svc.Service.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Equal(t, "30d", resp.Range)
	require.Equal(t, 2, resp.Metrics.TotalApplications)
	require.Equal(t, 1, resp.Metrics.Interviewing)
	require.Equal(t, 0, resp.Metrics.Offers)
	require.Equal(t, "腾讯", resp.CompanyRanking[0].Name)
	require.Equal(t, 2, resp.CompanyRanking[0].Count)
	require.ElementsMatch(t, []dto.DashboardNameCount{{Name: "后端", Count: 1}, {Name: "前端", Count: 1}}, resp.JobTitleDistribution)
	require.NotEmpty(t, resp.StatusPipeline)
	require.NotEmpty(t, resp.ApplicationTrend)
	require.Equal(t, 1000, svc.apps.lastFilter.PageSize)
}

func TestDashboardService_GetIncludesManualAndSystemTodos(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.apps.rows = []domainapp.ResumeApplication{{ID: "app-1", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", AppliedAt: now.AddDate(0, 0, -8), CurrentStatus: domainapp.StatusApplied}}
	svc.progresses.byApplication["app-1"] = []domainprogress.InterviewProgress{{ID: "p1", UserID: "u1", ApplicationID: "app-1", Stage: domainprogress.StageFirstInterview, OccurredAt: now, Result: domainprogress.ResultPending}}
	svc.todos.rows = []domaintodo.DashboardTodo{{ID: "todo-1", UserID: "u1", Title: "手动准备简历"}}

	resp, err := svc.Service.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Len(t, resp.Todos, 3)
	require.Equal(t, "manual", resp.Todos[0].Source)
	require.Contains(t, resp.Todos[1].Title, "跟进")
	require.Equal(t, "/resume/application/app-1", resp.Todos[1].Route)
	require.Contains(t, resp.Todos[2].Title, "面试")
	require.Equal(t, "system", resp.Todos[2].Source)
}

func TestDashboardService_GetIncludesRecentActivities(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.activities.rows = []domainactivity.ActivityLog{{ID: "a1", UserID: "u1", Action: domainactivity.ActionUpdated, TargetType: domainactivity.TargetTypeApplication, TargetID: "app-1", TargetLabel: "腾讯 / 前端工程师", Summary: "腾讯 / 前端工程师 投递状态更新为一面", Route: "/resume/application/app-1", CreatedAt: now}}

	resp, err := svc.Service.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Len(t, resp.Activities, 1)
	require.Equal(t, "/resume/application/app-1", resp.Activities[0].Route)
	require.NotNil(t, svc.activities.lastFilter.StartAt)
}

func TestDashboardService_GetInvalidRangeReturnsBadParams(t *testing.T) {
	svc := newTestService(time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC))

	_, err := svc.Service.Get(context.Background(), "u1", "bad")
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestDashboardService_TodoUseCases(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	ctx := context.Background()
	due := dto.FlexibleTime{Time: now.Add(24 * time.Hour)}
	targetID := "app-1"

	created, err := svc.Service.CreateTodo(ctx, "u1", dto.CreateDashboardTodoRequest{Title: "跟进腾讯", DueAt: &due, TargetType: domaintodo.TargetTypeApplication, TargetID: &targetID})
	require.NoError(t, err)
	require.Equal(t, "跟进腾讯", created.Title)
	require.NotNil(t, created.DueAt)

	updatedTitle := "准备面试"
	updated, err := svc.Service.UpdateTodo(ctx, "u1", created.ID, dto.UpdateDashboardTodoRequest{Title: &updatedTitle, ClearDueAt: true, ClearTarget: true})
	require.NoError(t, err)
	require.Equal(t, updatedTitle, updated.Title)
	require.Nil(t, updated.DueAt)
	require.Empty(t, updated.TargetType)
	require.Nil(t, updated.TargetID)

	completed, err := svc.Service.CompleteTodo(ctx, "u1", created.ID)
	require.NoError(t, err)
	require.NotNil(t, completed.CompletedAt)
	require.Equal(t, now, *completed.CompletedAt)

	require.NoError(t, svc.Service.DeleteTodo(ctx, "u1", created.ID))
	require.Empty(t, svc.todos.rows)
}

type dashboardTestService struct {
	*Service
	apps       *fakeApplicationRepo
	progresses *fakeProgressRepo
	todos      *fakeTodoRepo
	activities *fakeActivityRepo
}

func newTestService(now time.Time) *dashboardTestService {
	apps := &fakeApplicationRepo{}
	progresses := &fakeProgressRepo{byApplication: map[string][]domainprogress.InterviewProgress{}}
	todos := &fakeTodoRepo{}
	activities := &fakeActivityRepo{}
	svc := NewService(apps, progresses, todos, activities)
	svc.now = func() time.Time { return now }
	return &dashboardTestService{Service: svc, apps: apps, progresses: progresses, todos: todos, activities: activities}
}

type fakeApplicationRepo struct {
	rows       []domainapp.ResumeApplication
	lastFilter apprepo.ListFilter
}

func (r *fakeApplicationRepo) ListByUser(ctx context.Context, filter apprepo.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
	r.lastFilter = filter
	rows := make([]domainapp.ResumeApplication, 0, len(r.rows))
	for _, row := range r.rows {
		if row.UserID == filter.UserID {
			rows = append(rows, row)
		}
	}
	return rows, int64(len(rows)), nil
}

type fakeProgressRepo struct {
	byApplication map[string][]domainprogress.InterviewProgress
}

func (r *fakeProgressRepo) ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error) {
	rows := make([]domainprogress.InterviewProgress, 0)
	for _, row := range r.byApplication[applicationID] {
		if row.UserID == userID {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

type fakeTodoRepo struct {
	rows []domaintodo.DashboardTodo
}

func (r *fakeTodoRepo) Create(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	r.rows = append(r.rows, *todo)
	return nil
}

func (r *fakeTodoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domaintodo.DashboardTodo, error) {
	for i := range r.rows {
		if r.rows[i].ID == id && r.rows[i].UserID == userID {
			return &r.rows[i], nil
		}
	}
	return nil, domainerrors.ErrNotFound
}

func (r *fakeTodoRepo) Update(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	for i := range r.rows {
		if r.rows[i].ID == todo.ID && r.rows[i].UserID == todo.UserID {
			r.rows[i] = *todo
			return nil
		}
	}
	return domainerrors.ErrNotFound
}

func (r *fakeTodoRepo) Delete(ctx context.Context, id, userID string) error {
	for i := range r.rows {
		if r.rows[i].ID == id && r.rows[i].UserID == userID {
			r.rows = append(r.rows[:i], r.rows[i+1:]...)
			return nil
		}
	}
	return domainerrors.ErrNotFound
}

func (r *fakeTodoRepo) ListByUser(ctx context.Context, filter todorepo.ListFilter) ([]domaintodo.DashboardTodo, error) {
	rows := make([]domaintodo.DashboardTodo, 0, len(r.rows))
	for _, row := range r.rows {
		if row.UserID == filter.UserID {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

type fakeActivityRepo struct {
	rows       []domainactivity.ActivityLog
	lastFilter activityrepo.ListFilter
}

func (r *fakeActivityRepo) Create(ctx context.Context, activity *domainactivity.ActivityLog) error {
	r.rows = append(r.rows, *activity)
	return nil
}

func (r *fakeActivityRepo) ListByUser(ctx context.Context, filter activityrepo.ListFilter) ([]domainactivity.ActivityLog, error) {
	r.lastFilter = filter
	rows := make([]domainactivity.ActivityLog, 0, len(r.rows))
	for _, row := range r.rows {
		if row.UserID == filter.UserID && (filter.StartAt == nil || !row.CreatedAt.Before(*filter.StartAt)) {
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt.After(rows[j].CreatedAt) })
	if filter.Limit > 0 && len(rows) > filter.Limit {
		rows = rows[:filter.Limit]
	}
	return rows, nil
}
