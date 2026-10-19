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
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	apprepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
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

	resp, err := svc.Service.Get(context.Background(), "u1", "30d", 1, 5, 1, 5, false)
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
	require.Equal(t, 100, svc.apps.lastFilter.PageSize)
}

func TestDashboardService_GetIncludesManualAndSystemTodos(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.apps.rows = []domainapp.ResumeApplication{{ID: "app-1", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", AppliedAt: now.AddDate(0, 0, -8), CurrentStatus: domainapp.StatusApplied}}
	svc.progresses.byApplication["app-1"] = []domainprogress.InterviewProgress{{ID: "p1", UserID: "u1", ApplicationID: "app-1", Stage: domainprogress.StageFirstInterview, OccurredAt: now, Result: domainprogress.ResultPending}}
	svc.todos.rows = []domaintodo.DashboardTodo{{ID: "todo-1", UserID: "u1", Title: "手动准备简历"}}

	resp, err := svc.Service.Get(context.Background(), "u1", "30d", 1, 5, 1, 5, false)
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

	resp, err := svc.Service.Get(context.Background(), "u1", "30d", 1, 5, 1, 5, false)
	require.NoError(t, err)
	require.Len(t, resp.Activities, 1)
	require.Equal(t, "/resume/application/app-1", resp.Activities[0].Route)
	require.NotNil(t, svc.activities.lastFilter.StartAt)
}

func TestDashboardService_GetFetchesAllApplicationPages(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	for i := 0; i < 101; i++ {
		svc.apps.rows = append(svc.apps.rows, domainapp.ResumeApplication{ID: "app", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", AppliedAt: now, CurrentStatus: domainapp.StatusApplied})
	}

	resp, err := svc.Service.Get(context.Background(), "u1", "30d", 1, 5, 1, 5, false)
	require.NoError(t, err)
	require.Equal(t, 101, resp.Metrics.TotalApplications)
	require.Equal(t, []int{1, 2}, svc.apps.pages)
}

func TestDashboardService_ListTodosFiltersManualTodosAndComputesInterviewRoute(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	completedAt := now.Add(-time.Hour)
	targetID := "progress-1"
	svc.todos.rows = []domaintodo.DashboardTodo{
		{ID: "open", UserID: "u1", Title: "准备面试", TargetType: domaintodo.TargetTypeInterview, TargetID: &targetID},
		{ID: "done", UserID: "u1", Title: "完成", CompletedAt: &completedAt},
	}
	svc.progresses.byID = map[string]domainprogress.InterviewProgress{"progress-1": {ID: "progress-1", UserID: "u1", ApplicationID: "app-1"}}
	completed := false

	rows, err := svc.Service.ListTodos(context.Background(), "u1", &completed)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "open", rows[0].ID)
	require.Equal(t, "/resume/application/app-1", rows[0].Route)
}

func TestDashboardService_GetInvalidRangeReturnsBadParams(t *testing.T) {
	svc := newTestService(time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC))

	_, err := svc.Service.Get(context.Background(), "u1", "bad", 1, 5, 1, 5, false)
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestDashboardService_TodoUseCases(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	ctx := context.Background()
	due := dto.FlexibleTime{Time: now.Add(24 * time.Hour)}
	targetID := "app-1"
	svc.apps.rows = []domainapp.ResumeApplication{{ID: targetID, UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师"}}

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

func TestDashboardService_CreateTodoRejectsTargetsNotOwnedByUser(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	targetID := "app-1"

	_, err := svc.Service.CreateTodo(context.Background(), "u1", dto.CreateDashboardTodoRequest{Title: "跟进腾讯", TargetType: domaintodo.TargetTypeApplication, TargetID: &targetID})
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Empty(t, svc.todos.rows)
}

func TestDashboardService_UpdateTodoValidatesNewTargetOwnership(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	todo := domaintodo.DashboardTodo{ID: "todo-1", UserID: "u1", Title: "跟进腾讯"}
	svc.todos.rows = []domaintodo.DashboardTodo{todo}
	targetID := "question-1"
	targetType := domaintodo.TargetTypeQuestion

	_, err := svc.Service.UpdateTodo(context.Background(), "u1", "todo-1", dto.UpdateDashboardTodoRequest{TargetType: &targetType, TargetID: &targetID})
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Empty(t, svc.todos.rows[0].TargetType)
	require.Nil(t, svc.todos.rows[0].TargetID)
}

func TestDashboardService_CreateTodoRequiresConfiguredTargetReader(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.Service.WithTargetReaders(TargetReaders{})
	targetID := "app-1"

	_, err := svc.Service.CreateTodo(context.Background(), "u1", dto.CreateDashboardTodoRequest{Title: "跟进腾讯", TargetType: domaintodo.TargetTypeApplication, TargetID: &targetID})
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestDashboardService_RouteForTodoResolvesInterviewProgressApplication(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	targetID := "progress-1"
	svc.progresses.byID = map[string]domainprogress.InterviewProgress{"progress-1": {ID: "progress-1", UserID: "u1", ApplicationID: "app-1"}}

	route, err := svc.Service.RouteForTodo(context.Background(), "u1", domaintodo.TargetTypeInterview, &targetID)
	require.NoError(t, err)
	require.Equal(t, "/resume/application/app-1", route)
}

type dashboardTestService struct {
	*Service
	apps       *fakeApplicationRepo
	progresses *fakeProgressRepo
	questions  *fakeQuestionReader
	resumes    *fakeResumeReader
	todos      *fakeTodoRepo
	activities *fakeActivityRepo
}

func newTestService(now time.Time) *dashboardTestService {
	apps := &fakeApplicationRepo{}
	progresses := &fakeProgressRepo{byApplication: map[string][]domainprogress.InterviewProgress{}}
	questions := &fakeQuestionReader{idsByUser: map[string]map[string]bool{}}
	resumes := &fakeResumeReader{idsByUser: map[string]map[string]bool{}}
	todos := &fakeTodoRepo{}
	activities := &fakeActivityRepo{}
	svc := NewService(apps, progresses, todos, activities)
	svc.WithTargetReaders(TargetReaders{Applications: apps, Progresses: progresses, Questions: questions, Resumes: resumes})
	svc.now = func() time.Time { return now }
	return &dashboardTestService{Service: svc, apps: apps, progresses: progresses, questions: questions, resumes: resumes, todos: todos, activities: activities}
}

type fakeApplicationRepo struct {
	rows       []domainapp.ResumeApplication
	lastFilter apprepo.ListFilter
	pages      []int
}

func (r *fakeApplicationRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error) {
	for i := range r.rows {
		if r.rows[i].ID == id && r.rows[i].UserID == userID {
			return &r.rows[i], nil
		}
	}
	return nil, domainerrors.ErrNotFound
}

func (r *fakeApplicationRepo) ListByUser(ctx context.Context, filter apprepo.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
	r.lastFilter = filter
	r.pages = append(r.pages, filter.Page)
	rows := make([]domainapp.ResumeApplication, 0, len(r.rows))
	for _, row := range r.rows {
		if row.UserID == filter.UserID {
			rows = append(rows, row)
		}
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = len(rows)
	}
	start := (page - 1) * pageSize
	if start >= len(rows) {
		return []domainapp.ResumeApplication{}, int64(len(rows)), nil
	}
	end := start + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end], int64(len(rows)), nil
}

type fakeProgressRepo struct {
	byApplication map[string][]domainprogress.InterviewProgress
	byID          map[string]domainprogress.InterviewProgress
}

func (r *fakeProgressRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) {
	if r.byID != nil {
		if row, ok := r.byID[id]; ok && row.UserID == userID {
			return &row, nil
		}
	}
	for _, rows := range r.byApplication {
		for i := range rows {
			if rows[i].ID == id && rows[i].UserID == userID {
				return &rows[i], nil
			}
		}
	}
	return nil, domainerrors.ErrNotFound
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

type fakeQuestionReader struct {
	idsByUser map[string]map[string]bool
}

func (r *fakeQuestionReader) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error) {
	if r.idsByUser[userID] != nil && r.idsByUser[userID][id] {
		return &domainquestion.InterviewQuestion{ID: id, UserID: userID}, nil
	}
	return nil, domainerrors.ErrNotFound
}

type fakeResumeReader struct {
	idsByUser map[string]map[string]bool
}

func (r *fakeResumeReader) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) {
	if r.idsByUser[userID] != nil && r.idsByUser[userID][id] {
		return &domainrepo.ResumeRepo{ID: id, UserID: userID}, nil
	}
	return nil, domainerrors.ErrNotFound
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
		if row.UserID != filter.UserID {
			continue
		}
		if filter.Completed != nil {
			completed := row.CompletedAt != nil
			if completed != *filter.Completed {
				continue
			}
		}
		rows = append(rows, row)
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

func (r *fakeActivityRepo) ListByUser(ctx context.Context, filter activityrepo.ListFilter) ([]domainactivity.ActivityLog, int64, error) {
	r.lastFilter = filter
	rows := make([]domainactivity.ActivityLog, 0, len(r.rows))
	for _, row := range r.rows {
		if row.UserID == filter.UserID && (filter.StartAt == nil || !row.CreatedAt.Before(*filter.StartAt)) {
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt.After(rows[j].CreatedAt) })
	total := int64(len(rows))
	if filter.Offset < len(rows) {
		rows = rows[filter.Offset:]
	} else {
		rows = []domainactivity.ActivityLog{}
	}
	if filter.Limit > 0 && len(rows) > filter.Limit {
		rows = rows[:filter.Limit]
	}
	return rows, total, nil
}
