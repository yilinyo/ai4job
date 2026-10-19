# Job Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a job-search dashboard that replaces the template dashboard with backend-aggregated metrics, todos, activities, charts, and matching frontend mock data.

**Architecture:** Add focused backend domain/repository/service/handler units for dashboard todos, activity logs, and dashboard aggregation. Wire activity logging into resume, application, interview, and question write operations while keeping job-posting creation out of the activity stream. On the frontend, add typed dashboard models, API wrappers, Vuex module, mock endpoints, and a Vue 2 dashboard composed from small components.

**Tech Stack:** Go + Gin + GORM + PostgreSQL migrations for backend; Vue 2 + TypeScript + Vuex module decorators + Element UI + ECharts for frontend; Express mock server under `frontend/mock`.

---

## File Structure

### Backend files to create

- `backend/migrations/000008_create_dashboard_todos.up.sql` — creates manual dashboard todo table and indexes.
- `backend/migrations/000008_create_dashboard_todos.down.sql` — drops dashboard todo table.
- `backend/migrations/000009_create_activity_logs.up.sql` — creates activity log table and indexes.
- `backend/migrations/000009_create_activity_logs.down.sql` — drops activity log table.
- `backend/internal/domain/dashboard_todo/dashboard_todo.go` — dashboard todo entity, source/target constants, validation, completion/update behavior.
- `backend/internal/domain/dashboard_todo/repository/repository.go` — todo repository interface and list filter.
- `backend/internal/domain/activity_log/activity_log.go` — activity log entity, action/target constants, constructor.
- `backend/internal/domain/activity_log/repository/repository.go` — activity log repository interface and list filter.
- `backend/internal/application/dashboard/dashboard_service.go` — dashboard aggregation, range parsing, automatic todo generation, manual todo use cases.
- `backend/internal/application/dashboard/dashboard_service_test.go` — service tests for aggregation, ranges, manual todos, automatic todos.
- `backend/internal/application/activitylog/activity_logger.go` — small application logger wrapper used by existing services.
- `backend/internal/application/dto/dashboard_dto.go` — dashboard response/request DTOs and converters.
- `backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo.go` — GORM todo repository.
- `backend/internal/infrastructure/persistence/postgres/activity_log_repo.go` — GORM activity repository.
- `backend/internal/interface/handler/dashboard_handler.go` — dashboard and todo HTTP handlers.
- `backend/internal/interface/handler/dashboard_handler_test.go` — handler shape and auth-scoping tests.

### Backend files to modify

- `backend/cmd/server/main.go` — instantiate new repos/services/handlers and inject activity logger.
- `backend/internal/interface/router/router.go` — register `/dashboard` and `/dashboard/todos` routes.
- `backend/internal/application/resume_repo/resume_repo_service.go` — log resume repo create/update/delete activities.
- `backend/internal/application/resume_version/resume_version_service.go` — log resume version create/update/delete/activate activities.
- `backend/internal/application/resume_application/resume_application_service.go` — log application create/update/delete/status-change activities.
- `backend/internal/application/interview_progress/interview_progress_service.go` — log interview create/update/delete activities.
- `backend/internal/application/interview_question/interview_question_service.go` — log question create/update/delete/favorite/occurrence activities.
- Existing service tests in the same packages — add logger-aware tests or no-op logger setup where constructors change.
- `backend/migrations/down_migrations_test.go` — add expected down migration statements for migrations 000008 and 000009.

### Frontend files to create

- `frontend/src/models/dashboard.ts` — dashboard TypeScript models.
- `frontend/src/api/dashboard.ts` — dashboard API wrapper.
- `frontend/src/store/modules/dashboard.ts` — dashboard Vuex module.
- `frontend/src/views/dashboard/components/MetricCards.vue` — KPI cards.
- `frontend/src/views/dashboard/components/PipelineChart.vue` — status pipeline and trend chart.
- `frontend/src/views/dashboard/components/TodoPanel.vue` — manual/system todo list and edit dialog.
- `frontend/src/views/dashboard/components/DistributionCharts.vue` — company, job-title, source charts.
- `frontend/src/views/dashboard/components/ActivityFeed.vue` — recent activity feed with optional routes.
- `frontend/mock/dashboard.ts` — mock dashboard and todo endpoints.

### Frontend files to modify

- `frontend/src/models/index.ts` — export dashboard models.
- `frontend/src/api/index.ts` — export dashboard API wrappers.
- `frontend/src/views/dashboard/index.vue` — replace template role switcher with job dashboard page.
- `frontend/mock/mock-server.ts` — register dashboard mock routes.

---

## Task 1: Add Dashboard Tables and Domain Models

**Files:**
- Create: `backend/migrations/000008_create_dashboard_todos.up.sql`
- Create: `backend/migrations/000008_create_dashboard_todos.down.sql`
- Create: `backend/migrations/000009_create_activity_logs.up.sql`
- Create: `backend/migrations/000009_create_activity_logs.down.sql`
- Create: `backend/internal/domain/dashboard_todo/dashboard_todo.go`
- Create: `backend/internal/domain/dashboard_todo/repository/repository.go`
- Create: `backend/internal/domain/activity_log/activity_log.go`
- Create: `backend/internal/domain/activity_log/repository/repository.go`
- Modify: `backend/migrations/down_migrations_test.go`

- [ ] **Step 1: Add failing down migration test expectations**

Modify `backend/migrations/down_migrations_test.go` so the expected statements include:

```go
"000008_create_dashboard_todos.down.sql": "DROP TABLE IF EXISTS dashboard_todos;",
"000009_create_activity_logs.down.sql":    "DROP TABLE IF EXISTS activity_logs;",
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./migrations
```

Expected: FAIL because the new migration files do not exist yet.

- [ ] **Step 2: Create dashboard todo migrations**

Create `backend/migrations/000008_create_dashboard_todos.up.sql`:

```sql
CREATE TABLE dashboard_todos (
    id           UUID        PRIMARY KEY,
    user_id      VARCHAR     NOT NULL REFERENCES users(id),
    title        VARCHAR(200) NOT NULL,
    description  TEXT        NOT NULL DEFAULT '',
    due_at       TIMESTAMP,
    completed_at TIMESTAMP,
    target_type  VARCHAR(40) NOT NULL DEFAULT '',
    target_id    UUID,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dashboard_todos_user_id ON dashboard_todos(user_id);
CREATE INDEX idx_dashboard_todos_user_completed ON dashboard_todos(user_id, completed_at);
CREATE INDEX idx_dashboard_todos_user_due_at ON dashboard_todos(user_id, due_at);
CREATE INDEX idx_dashboard_todos_target ON dashboard_todos(target_type, target_id);
```

Create `backend/migrations/000008_create_dashboard_todos.down.sql`:

```sql
DROP TABLE IF EXISTS dashboard_todos;
```

- [ ] **Step 3: Create activity log migrations**

Create `backend/migrations/000009_create_activity_logs.up.sql`:

```sql
CREATE TABLE activity_logs (
    id           UUID        PRIMARY KEY,
    user_id      VARCHAR     NOT NULL REFERENCES users(id),
    actor_type   VARCHAR(40) NOT NULL DEFAULT 'user',
    action       VARCHAR(40) NOT NULL,
    target_type  VARCHAR(40) NOT NULL,
    target_id    UUID        NOT NULL,
    target_label VARCHAR(200) NOT NULL,
    summary      TEXT        NOT NULL,
    route        VARCHAR(500) NOT NULL DEFAULT '',
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_created_at ON activity_logs(user_id, created_at DESC);
CREATE INDEX idx_activity_logs_target ON activity_logs(target_type, target_id);
```

Create `backend/migrations/000009_create_activity_logs.down.sql`:

```sql
DROP TABLE IF EXISTS activity_logs;
```

- [ ] **Step 4: Add dashboard todo domain entity**

Create `backend/internal/domain/dashboard_todo/dashboard_todo.go`:

```go
package dashboard_todo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type Source string

const (
	SourceManual Source = "manual"
	SourceSystem Source = "system"
)

type TargetType string

const (
	TargetTypeApplication TargetType = "application"
	TargetTypeInterview    TargetType = "interview"
	TargetTypeQuestion     TargetType = "question"
	TargetTypeResume       TargetType = "resume"
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
		ID:          params.ID,
		UserID:      strings.TrimSpace(params.UserID),
		Title:       strings.TrimSpace(params.Title),
		Description: strings.TrimSpace(params.Description),
		DueAt:       params.DueAt,
		TargetType:  params.TargetType,
		TargetID:    params.TargetID,
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
		t.TargetID = *spec.TargetID
	}
	return t.validate()
}

func (t *DashboardTodo) Complete(now time.Time) {
	t.CompletedAt = &now
}

func (t *DashboardTodo) validate() error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.UserID) == "" || strings.TrimSpace(t.Title) == "" {
		return domainerrors.ErrBadParams
	}
	if t.TargetType != "" && !t.TargetType.Valid() {
		return domainerrors.ErrBadParams
	}
	if t.TargetType != "" && t.TargetID == nil {
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
```

- [ ] **Step 5: Add dashboard todo repository interface**

Create `backend/internal/domain/dashboard_todo/repository/repository.go`:

```go
package repository

import (
	"context"

	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
)

type ListFilter struct {
	UserID    string
	Completed *bool
	TodayOnly bool
}

type Repository interface {
	Create(ctx context.Context, todo *domaintodo.DashboardTodo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domaintodo.DashboardTodo, error)
	Update(ctx context.Context, todo *domaintodo.DashboardTodo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, filter ListFilter) ([]domaintodo.DashboardTodo, error)
}
```

- [ ] **Step 6: Add activity log domain entity**

Create `backend/internal/domain/activity_log/activity_log.go`:

```go
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
	TargetTypeInterview    TargetType = "interview"
	TargetTypeQuestion     TargetType = "question"
	TargetTypeResume       TargetType = "resume"
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
		ID:          params.ID,
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
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.UserID) == "" || strings.TrimSpace(l.TargetID) == "" || strings.TrimSpace(l.TargetLabel) == "" || strings.TrimSpace(l.Summary) == "" {
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
```

- [ ] **Step 7: Add activity log repository interface**

Create `backend/internal/domain/activity_log/repository/repository.go`:

```go
package repository

import (
	"context"
	"time"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
)

type ListFilter struct {
	UserID    string
	StartAt   *time.Time
	Limit     int
}

type Repository interface {
	Create(ctx context.Context, activity *domainactivity.ActivityLog) error
	ListByUser(ctx context.Context, filter ListFilter) ([]domainactivity.ActivityLog, error)
}
```

- [ ] **Step 8: Run migration/domain tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./migrations ./internal/domain/...
```

Expected: PASS.

- [ ] **Step 9: Commit task 1**

```bash
cd /Users/yilin/project/ai-for
git add backend/migrations/000008_create_dashboard_todos.up.sql backend/migrations/000008_create_dashboard_todos.down.sql backend/migrations/000009_create_activity_logs.up.sql backend/migrations/000009_create_activity_logs.down.sql backend/migrations/down_migrations_test.go backend/internal/domain/dashboard_todo backend/internal/domain/activity_log
git commit -m "feat: add dashboard todo and activity log models"
```

---

## Task 2: Add Postgres Repositories

**Files:**
- Create: `backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/activity_log_repo.go`
- Test: `backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo_test.go`
- Test: `backend/internal/infrastructure/persistence/postgres/activity_log_repo_test.go`

- [ ] **Step 1: Write failing repository tests**

Create `backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo_test.go`:

```go
package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	todorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
)

func TestDashboardTodoRepo_ListByUserScopesAndFiltersCompletion(t *testing.T) {
	db := openTestDB(t)
	repo := NewDashboardTodoRepo(db)
	ctx := context.Background()
	due := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	completedAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)

	openTodo := &domaintodo.DashboardTodo{ID: "00000000-0000-0000-0000-000000000801", UserID: "u1", Title: "跟进腾讯", DueAt: &due}
	completedTodo := &domaintodo.DashboardTodo{ID: "00000000-0000-0000-0000-000000000802", UserID: "u1", Title: "准备面试", CompletedAt: &completedAt}
	otherUserTodo := &domaintodo.DashboardTodo{ID: "00000000-0000-0000-0000-000000000803", UserID: "u2", Title: "其他用户"}
	require.NoError(t, repo.Create(ctx, openTodo))
	require.NoError(t, repo.Create(ctx, completedTodo))
	require.NoError(t, repo.Create(ctx, otherUserTodo))

	completed := false
	rows, err := repo.ListByUser(ctx, todorepo.ListFilter{UserID: "u1", Completed: &completed})
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "跟进腾讯", rows[0].Title)
}
```

Create `backend/internal/infrastructure/persistence/postgres/activity_log_repo_test.go`:

```go
package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activityrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
)

func TestActivityLogRepo_ListByUserScopesOrdersAndLimits(t *testing.T) {
	db := openTestDB(t)
	repo := NewActivityLogRepo(db)
	ctx := context.Background()

	older := &domainactivity.ActivityLog{ID: "00000000-0000-0000-0000-000000000901", UserID: "u1", ActorType: "user", Action: domainactivity.ActionCreated, TargetType: domainactivity.TargetTypeApplication, TargetID: "00000000-0000-0000-0000-000000000501", TargetLabel: "A / Backend", Summary: "创建投递 A", CreatedAt: time.Date(2026, 5, 17, 9, 0, 0, 0, time.UTC)}
	newer := &domainactivity.ActivityLog{ID: "00000000-0000-0000-0000-000000000902", UserID: "u1", ActorType: "user", Action: domainactivity.ActionUpdated, TargetType: domainactivity.TargetTypeQuestion, TargetID: "00000000-0000-0000-0000-000000000701", TargetLabel: "Vue 生命周期", Summary: "更新题目 Vue 生命周期", CreatedAt: time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)}
	other := &domainactivity.ActivityLog{ID: "00000000-0000-0000-0000-000000000903", UserID: "u2", ActorType: "user", Action: domainactivity.ActionCreated, TargetType: domainactivity.TargetTypeApplication, TargetID: "00000000-0000-0000-0000-000000000502", TargetLabel: "B / Frontend", Summary: "创建投递 B", CreatedAt: time.Date(2026, 5, 19, 9, 0, 0, 0, time.UTC)}
	require.NoError(t, repo.Create(ctx, older))
	require.NoError(t, repo.Create(ctx, newer))
	require.NoError(t, repo.Create(ctx, other))

	rows, err := repo.ListByUser(ctx, activityrepo.ListFilter{UserID: "u1", Limit: 1})
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "更新题目 Vue 生命周期", rows[0].Summary)
}
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/infrastructure/persistence/postgres -run 'TestDashboardTodoRepo|TestActivityLogRepo'
```

Expected: FAIL because repository constructors do not exist.

- [ ] **Step 2: Implement dashboard todo repository**

Create `backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo.go`:

```go
package postgres

import (
	"context"
	"errors"
	"time"

	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	todorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type DashboardTodoRepo struct {
	db *gorm.DB
}

var _ todorepo.Repository = (*DashboardTodoRepo)(nil)

func NewDashboardTodoRepo(db *gorm.DB) *DashboardTodoRepo {
	return &DashboardTodoRepo{db: db}
}

func (r *DashboardTodoRepo) Create(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *DashboardTodoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domaintodo.DashboardTodo, error) {
	var todo domaintodo.DashboardTodo
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &todo, err
}

func (r *DashboardTodoRepo) Update(ctx context.Context, todo *domaintodo.DashboardTodo) error {
	result := r.db.WithContext(ctx).
		Model(&domaintodo.DashboardTodo{}).
		Where("id = ? AND user_id = ?", todo.ID, todo.UserID).
		Updates(map[string]any{
			"title":        todo.Title,
			"description":  todo.Description,
			"due_at":       todo.DueAt,
			"completed_at": todo.CompletedAt,
			"target_type":  todo.TargetType,
			"target_id":    todo.TargetID,
			"updated_at":   time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *DashboardTodoRepo) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domaintodo.DashboardTodo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *DashboardTodoRepo) ListByUser(ctx context.Context, filter todorepo.ListFilter) ([]domaintodo.DashboardTodo, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", filter.UserID)
	if filter.Completed != nil {
		if *filter.Completed {
			q = q.Where("completed_at IS NOT NULL")
		} else {
			q = q.Where("completed_at IS NULL")
		}
	}
	if filter.TodayOnly {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)
		q = q.Where("due_at IS NULL OR due_at < ?", end).Where("completed_at IS NULL OR completed_at >= ?", start)
	}
	var rows []domaintodo.DashboardTodo
	err := q.Order("COALESCE(due_at, created_at) ASC").Find(&rows).Error
	return rows, err
}
```

- [ ] **Step 3: Implement activity log repository**

Create `backend/internal/infrastructure/persistence/postgres/activity_log_repo.go`:

```go
package postgres

import (
	"context"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activityrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	"gorm.io/gorm"
)

type ActivityLogRepo struct {
	db *gorm.DB
}

var _ activityrepo.Repository = (*ActivityLogRepo)(nil)

func NewActivityLogRepo(db *gorm.DB) *ActivityLogRepo {
	return &ActivityLogRepo{db: db}
}

func (r *ActivityLogRepo) Create(ctx context.Context, activity *domainactivity.ActivityLog) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

func (r *ActivityLogRepo) ListByUser(ctx context.Context, filter activityrepo.ListFilter) ([]domainactivity.ActivityLog, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", filter.UserID)
	if filter.StartAt != nil {
		q = q.Where("created_at >= ?", *filter.StartAt)
	}
	limit := filter.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	var rows []domainactivity.ActivityLog
	err := q.Order("created_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}
```

- [ ] **Step 4: Run repository tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/infrastructure/persistence/postgres -run 'TestDashboardTodoRepo|TestActivityLogRepo'
```

Expected: PASS.

- [ ] **Step 5: Commit task 2**

```bash
cd /Users/yilin/project/ai-for
git add backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo.go backend/internal/infrastructure/persistence/postgres/activity_log_repo.go backend/internal/infrastructure/persistence/postgres/dashboard_todo_repo_test.go backend/internal/infrastructure/persistence/postgres/activity_log_repo_test.go
git commit -m "feat: add dashboard postgres repositories"
```

---

## Task 3: Add Dashboard DTOs, Service, and Handler

**Files:**
- Create: `backend/internal/application/dto/dashboard_dto.go`
- Create: `backend/internal/application/dashboard/dashboard_service.go`
- Create: `backend/internal/application/dashboard/dashboard_service_test.go`
- Create: `backend/internal/interface/handler/dashboard_handler.go`
- Create: `backend/internal/interface/handler/dashboard_handler_test.go`
- Modify: `backend/internal/interface/router/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Add dashboard DTOs**

Create `backend/internal/application/dto/dashboard_dto.go`:

```go
package dto

import (
	"time"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
)

type DashboardQuery struct {
	Range string `json:"range" form:"range"`
}

type CreateDashboardTodoRequest struct {
	Title       string                 `json:"title" binding:"required,min=1,max=200"`
	Description string                 `json:"description"`
	DueAt       *FlexibleTime          `json:"dueAt"`
	TargetType  domaintodo.TargetType  `json:"targetType"`
	TargetID    *string                `json:"targetId"`
}

type UpdateDashboardTodoRequest struct {
	Title       *string                `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string                `json:"description"`
	DueAt       *FlexibleTime          `json:"dueAt"`
	ClearDueAt  bool                   `json:"clearDueAt"`
	TargetType  *domaintodo.TargetType `json:"targetType"`
	TargetID    *string                `json:"targetId"`
	ClearTarget bool                   `json:"clearTarget"`
}

type DashboardResponse struct {
	Range                string                   `json:"range"`
	Metrics              DashboardMetrics         `json:"metrics"`
	StatusPipeline       []DashboardStatusBucket  `json:"statusPipeline"`
	ApplicationTrend     []DashboardTrendPoint    `json:"applicationTrend"`
	CompanyRanking       []DashboardNameCount     `json:"companyRanking"`
	JobTitleDistribution []DashboardNameCount     `json:"jobTitleDistribution"`
	SourceDistribution   []DashboardNameCount     `json:"sourceDistribution"`
	Todos                []DashboardTodoResponse  `json:"todos"`
	Activities           []ActivityLogResponse    `json:"activities"`
}

type DashboardMetrics struct {
	TotalApplications int `json:"totalApplications"`
	Interviewing      int `json:"interviewing"`
	PendingFollowUp   int `json:"pendingFollowUp"`
	Offers            int `json:"offers"`
}

type DashboardStatusBucket struct {
	Status string `json:"status"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
}

type DashboardTrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DashboardNameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DashboardTodoResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DueAt       *time.Time `json:"dueAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	TargetType  string  `json:"targetType,omitempty"`
	TargetID    *string `json:"targetId,omitempty"`
	Route       string  `json:"route,omitempty"`
	Source      string  `json:"source"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ActivityLogResponse struct {
	ID          string    `json:"id"`
	Action      string    `json:"action"`
	TargetType  string    `json:"targetType"`
	TargetID    string    `json:"targetId"`
	TargetLabel string    `json:"targetLabel"`
	Summary     string    `json:"summary"`
	Route       string    `json:"route,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func ToDashboardTodoResponse(todo domaintodo.DashboardTodo, source domaintodo.Source, route string) DashboardTodoResponse {
	return DashboardTodoResponse{ID: todo.ID, Title: todo.Title, Description: todo.Description, DueAt: todo.DueAt, CompletedAt: todo.CompletedAt, TargetType: string(todo.TargetType), TargetID: todo.TargetID, Route: route, Source: string(source), CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt}
}

func ToActivityLogResponse(activity domainactivity.ActivityLog) ActivityLogResponse {
	return ActivityLogResponse{ID: activity.ID, Action: string(activity.Action), TargetType: string(activity.TargetType), TargetID: activity.TargetID, TargetLabel: activity.TargetLabel, Summary: activity.Summary, Route: activity.Route, CreatedAt: activity.CreatedAt}
}
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/dto
```

Expected: PASS or `no test files`; compilation succeeds.

- [ ] **Step 2: Write failing dashboard service tests**

Create `backend/internal/application/dashboard/dashboard_service_test.go` with in-memory fakes. Include these tests:

```go
package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

func TestDashboardService_GetAggregatesMetricsAndDistributions(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.applications = []domainapp.ResumeApplication{
		{ID: "app-1", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", SourcePlatform: "Boss", AppliedAt: now.AddDate(0, 0, -1), CurrentStatus: domainapp.StatusApplied},
		{ID: "app-2", UserID: "u1", CompanyName: "腾讯", JobTitle: "后端工程师", SourcePlatform: "官网", AppliedAt: now.AddDate(0, 0, -2), CurrentStatus: domainapp.StatusFirstInterview},
		{ID: "app-3", UserID: "u1", CompanyName: "字节", JobTitle: "AI 工程师", SourcePlatform: "Boss", AppliedAt: now.AddDate(0, 0, -40), CurrentStatus: domainapp.StatusOffer},
	}

	resp, err := svc.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Equal(t, "30d", resp.Range)
	require.Equal(t, 2, resp.Metrics.TotalApplications)
	require.Equal(t, 1, resp.Metrics.Interviewing)
	require.Equal(t, 0, resp.Metrics.Offers)
	require.Equal(t, "腾讯", resp.CompanyRanking[0].Name)
	require.Equal(t, 2, resp.CompanyRanking[0].Count)
	require.NotEmpty(t, resp.StatusPipeline)
	require.NotEmpty(t, resp.ApplicationTrend)
}

func TestDashboardService_GetIncludesManualAndSystemTodos(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.applications = []domainapp.ResumeApplication{{ID: "app-1", UserID: "u1", CompanyName: "腾讯", JobTitle: "前端工程师", AppliedAt: now.AddDate(0, 0, -8), CurrentStatus: domainapp.StatusApplied}}
	svc.progresses = []domainprogress.InterviewProgress{{ID: "p1", UserID: "u1", ApplicationID: "app-1", Stage: domainprogress.StageFirstInterview, OccurredAt: now, Result: domainprogress.ResultPending}}
	svc.todos = []domaintodo.DashboardTodo{{ID: "todo-1", UserID: "u1", Title: "手动准备简历"}}

	resp, err := svc.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Len(t, resp.Todos, 3)
	require.Equal(t, "manual", resp.Todos[0].Source)
	require.Contains(t, resp.Todos[1].Title, "跟进")
	require.Contains(t, resp.Todos[2].Title, "面试")
}

func TestDashboardService_GetIncludesRecentActivities(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	svc := newTestService(now)
	svc.activities = []domainactivity.ActivityLog{{ID: "a1", UserID: "u1", Action: domainactivity.ActionUpdated, TargetType: domainactivity.TargetTypeApplication, TargetID: "app-1", TargetLabel: "腾讯 / 前端工程师", Summary: "腾讯 / 前端工程师 投递状态更新为一面", Route: "/resume/application/app-1", CreatedAt: now}}

	resp, err := svc.Get(context.Background(), "u1", "30d")
	require.NoError(t, err)
	require.Len(t, resp.Activities, 1)
	require.Equal(t, "/resume/application/app-1", resp.Activities[0].Route)
}
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/dashboard
```

Expected: FAIL because service does not exist.

- [ ] **Step 3: Implement dashboard service**

Create `backend/internal/application/dashboard/dashboard_service.go`. Implement these public methods and helpers:

```go
package dashboard

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
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

const followUpDays = 7

type ApplicationRepository interface {
	ListByUser(ctx context.Context, filter apprepo.ListFilter) ([]domainapp.ResumeApplication, int64, error)
}

type InterviewProgressRepository interface {
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error)
}

type Service struct {
	apps       ApplicationRepository
	progresses InterviewProgressRepository
	todos      todorepo.Repository
	activities activityrepo.Repository
	now        func() time.Time
}

func NewService(apps ApplicationRepository, progresses InterviewProgressRepository, todos todorepo.Repository, activities activityrepo.Repository) *Service {
	return &Service{apps: apps, progresses: progresses, todos: todos, activities: activities, now: time.Now}
}

func (s *Service) Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error) {
	rangeKey, startAt, err := parseRange(rawRange, s.now())
	if err != nil {
		return nil, err
	}
	applications, _, err := s.apps.ListByUser(ctx, apprepo.ListFilter{UserID: userID, Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	filtered := filterApplications(applications, startAt)
	manualTodos, err := s.todos.ListByUser(ctx, todorepo.ListFilter{UserID: userID, TodayOnly: true})
	if err != nil {
		return nil, err
	}
	activities, err := s.activities.ListByUser(ctx, activityrepo.ListFilter{UserID: userID, StartAt: startAt, Limit: 10})
	if err != nil {
		return nil, err
	}
	allTodos := make([]dto.DashboardTodoResponse, 0, len(manualTodos)+4)
	for _, todo := range manualTodos {
		allTodos = append(allTodos, dto.ToDashboardTodoResponse(todo, domaintodo.SourceManual, routeForTodo(todo.TargetType, todo.TargetID)))
	}
	systemTodos, err := s.systemTodos(ctx, userID, applications)
	if err != nil {
		return nil, err
	}
	allTodos = append(allTodos, systemTodos...)
	activityResponses := make([]dto.ActivityLogResponse, 0, len(activities))
	for _, activity := range activities {
		activityResponses = append(activityResponses, dto.ToActivityLogResponse(activity))
	}
	return &dto.DashboardResponse{Range: rangeKey, Metrics: buildMetrics(filtered), StatusPipeline: buildStatusPipeline(filtered), ApplicationTrend: buildTrend(filtered), CompanyRanking: buildNameCounts(filtered, func(a domainapp.ResumeApplication) string { return a.CompanyName }), JobTitleDistribution: buildNameCounts(filtered, func(a domainapp.ResumeApplication) string { return classifyJobTitle(a.JobTitle) }), SourceDistribution: buildNameCounts(filtered, func(a domainapp.ResumeApplication) string { return firstNonEmpty(a.SourcePlatform, a.DeliveryChannel, "其他") }), Todos: allTodos, Activities: activityResponses}, nil
}

func (s *Service) CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	var dueAt *time.Time
	if req.DueAt != nil {
		value := req.DueAt.Value()
		dueAt = &value
	}
	todo, err := domaintodo.NewDashboardTodo(domaintodo.NewDashboardTodoParams{ID: uuid.NewString(), UserID: userID, Title: req.Title, Description: req.Description, DueAt: dueAt, TargetType: req.TargetType, TargetID: req.TargetID})
	if err != nil {
		return nil, err
	}
	if err := s.todos.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	todo, err := s.todos.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	var dueAt **time.Time
	if req.ClearDueAt {
		var cleared *time.Time
		dueAt = &cleared
	} else if req.DueAt != nil {
		value := req.DueAt.Value()
		dueAt = &(&value)
	}
	var targetID **string
	if req.ClearTarget {
		var cleared *string
		targetID = &cleared
		empty := domaintodo.TargetType("")
		req.TargetType = &empty
	} else if req.TargetID != nil {
		targetID = &req.TargetID
	}
	if err := todo.ApplyUpdate(domaintodo.UpdateSpec{Title: req.Title, Description: req.Description, DueAt: dueAt, TargetType: req.TargetType, TargetID: targetID}); err != nil {
		return nil, err
	}
	if err := s.todos.Update(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error) {
	todo, err := s.todos.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	todo.Complete(s.now())
	if err := s.todos.Update(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) DeleteTodo(ctx context.Context, userID, id string) error {
	return s.todos.Delete(ctx, id, userID)
}
```

Also implement helper functions in the same file:

```go
func parseRange(raw string, now time.Time) (string, *time.Time, error) { /* accepts "", 7d, 30d, 90d, all */ }
func filterApplications(rows []domainapp.ResumeApplication, startAt *time.Time) []domainapp.ResumeApplication { /* AppliedAt >= startAt */ }
func buildMetrics(rows []domainapp.ResumeApplication) dto.DashboardMetrics { /* total, interviewing statuses, pending follow-up, offers */ }
func buildStatusPipeline(rows []domainapp.ResumeApplication) []dto.DashboardStatusBucket { /* fixed ordered statuses */ }
func buildTrend(rows []domainapp.ResumeApplication) []dto.DashboardTrendPoint { /* date yyyy-mm-dd counts */ }
func buildNameCounts(rows []domainapp.ResumeApplication, keyFn func(domainapp.ResumeApplication) string) []dto.DashboardNameCount { /* top 10 sorted desc */ }
func classifyJobTitle(title string) string { /* 前端, 后端, 全栈, AI, 测试, 产品, 其他 */ }
func routeForTodo(target domaintodo.TargetType, id *string) string { /* application/question/resume routes; interview target route uses application route when id points to app */ }
func firstNonEmpty(values ...string) string { /* trim and fallback */ }
```

Use `domainerrors.ErrBadParams` for invalid ranges.

- [ ] **Step 4: Run dashboard service tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/dashboard
```

Expected: PASS.

- [ ] **Step 5: Write failing dashboard handler tests**

Create `backend/internal/interface/handler/dashboard_handler_test.go`:

```go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

type mockDashboardService struct{ mock.Mock }

func (m *mockDashboardService) Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error) { args := m.Called(ctx, userID, rawRange); return args.Get(0).(*dto.DashboardResponse), args.Error(1) }
func (m *mockDashboardService) CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) { args := m.Called(ctx, userID, req); return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1) }
func (m *mockDashboardService) UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) { args := m.Called(ctx, userID, id, req); return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1) }
func (m *mockDashboardService) CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error) { args := m.Called(ctx, userID, id); return args.Get(0).(*domaintodo.DashboardTodo), args.Error(1) }
func (m *mockDashboardService) DeleteTodo(ctx context.Context, userID, id string) error { args := m.Called(ctx, userID, id); return args.Error(0) }

func TestDashboardHandler_GetUsesRangeAndDTOShape(t *testing.T) {
	svc := &mockDashboardService{}
	svc.On("Get", mock.Anything, "u1", "30d").Return(&dto.DashboardResponse{Range: "30d", Metrics: dto.DashboardMetrics{TotalApplications: 2}}, nil).Once()
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	r.GET("/dashboard", func(c *gin.Context) { c.Set("userID", "u1"); h.Get(c) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dashboard?range=30d", nil)
	r.ServeHTTP(w, req)
	var resp struct { Data map[string]any `json:"data"` }
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "30d", resp.Data["range"])
	require.Contains(t, resp.Data, "metrics")
	svc.AssertExpectations(t)
}

func TestDashboardHandler_CreateAndCompleteTodo(t *testing.T) {
	svc := &mockDashboardService{}
	createReq := dto.CreateDashboardTodoRequest{Title: "跟进腾讯"}
	todo := &domaintodo.DashboardTodo{ID: "todo-1", UserID: "u1", Title: "跟进腾讯"}
	svc.On("CreateTodo", mock.Anything, "u1", createReq).Return(todo, nil).Once()
	svc.On("CompleteTodo", mock.Anything, "u1", "todo-1").Return(todo, nil).Once()
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	r.POST("/dashboard/todos", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateTodo(c) })
	r.PUT("/dashboard/todos/:id/complete", func(c *gin.Context) { c.Set("userID", "u1"); h.CompleteTodo(c) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/dashboard/todos", bytes.NewBufferString(`{"title":"跟进腾讯"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/dashboard/todos/todo-1/complete", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/handler -run TestDashboardHandler
```

Expected: FAIL because handler does not exist.

- [ ] **Step 6: Implement dashboard handler**

Create `backend/internal/interface/handler/dashboard_handler.go`:

```go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type DashboardService interface {
	Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error)
	CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error)
	UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error)
	CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error)
	DeleteTodo(ctx context.Context, userID, id string) error
}

type DashboardHandler struct { svc DashboardService }

func NewDashboardHandler(svc DashboardService) *DashboardHandler { return &DashboardHandler{svc: svc} }

func (h *DashboardHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c); if !ok { return }
	row, err := h.svc.Get(c.Request.Context(), userID, c.Query("range"))
	if err != nil { mapError(c, err); return }
	response.Success(c, row)
}

func (h *DashboardHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateDashboardTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }
	userID, ok := resumeUserID(c); if !ok { return }
	todo, err := h.svc.CreateTodo(c.Request.Context(), userID, req)
	if err != nil { mapError(c, err); return }
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) UpdateTodo(c *gin.Context) {
	var req dto.UpdateDashboardTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }
	userID, ok := resumeUserID(c); if !ok { return }
	todo, err := h.svc.UpdateTodo(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil { mapError(c, err); return }
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) CompleteTodo(c *gin.Context) {
	userID, ok := resumeUserID(c); if !ok { return }
	todo, err := h.svc.CompleteTodo(c.Request.Context(), userID, c.Param("id"))
	if err != nil { mapError(c, err); return }
	response.Success(c, dto.ToDashboardTodoResponse(*todo, domaintodo.SourceManual, ""))
}

func (h *DashboardHandler) DeleteTodo(c *gin.Context) {
	userID, ok := resumeUserID(c); if !ok { return }
	if err := h.svc.DeleteTodo(c.Request.Context(), userID, c.Param("id")); err != nil { mapError(c, err); return }
	response.Success(c, nil)
}
```

- [ ] **Step 7: Wire router and server**

Modify `backend/internal/interface/router/router.go`:

```go
DashboardHandler *handler.DashboardHandler
```

Add routes inside the authenticated group before wildcard-like detail routes:

```go
auth.GET("/dashboard", deps.DashboardHandler.Get)
auth.POST("/dashboard/todos", deps.DashboardHandler.CreateTodo)
auth.PUT("/dashboard/todos/:id", deps.DashboardHandler.UpdateTodo)
auth.PUT("/dashboard/todos/:id/complete", deps.DashboardHandler.CompleteTodo)
auth.DELETE("/dashboard/todos/:id", deps.DashboardHandler.DeleteTodo)
```

Modify `backend/cmd/server/main.go` imports:

```go
dashboardapp "github.com/yilin/ai-for-backend/internal/application/dashboard"
```

Instantiate repositories:

```go
dashboardTodoRepo := infraDB.NewDashboardTodoRepo(db)
activityLogRepo := infraDB.NewActivityLogRepo(db)
```

Instantiate service and handler:

```go
dashboardSvc := dashboardapp.NewService(resumeApplicationRepo, interviewProgressRepo, dashboardTodoRepo, activityLogRepo)
dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
```

Pass handler into router dependencies:

```go
DashboardHandler: dashboardHandler,
```

- [ ] **Step 8: Run handler/router/server compilation tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/dashboard ./internal/interface/handler ./internal/interface/router ./cmd/server
```

Expected: PASS.

- [ ] **Step 9: Commit task 3**

```bash
cd /Users/yilin/project/ai-for
git add backend/internal/application/dto/dashboard_dto.go backend/internal/application/dashboard backend/internal/interface/handler/dashboard_handler.go backend/internal/interface/handler/dashboard_handler_test.go backend/internal/interface/router/router.go backend/cmd/server/main.go
git commit -m "feat: add dashboard aggregation API"
```

---

## Task 4: Add Activity Logging to Existing Write Operations

**Files:**
- Create: `backend/internal/application/activitylog/activity_logger.go`
- Modify: `backend/internal/application/resume_repo/resume_repo_service.go`
- Modify: `backend/internal/application/resume_version/resume_version_service.go`
- Modify: `backend/internal/application/resume_application/resume_application_service.go`
- Modify: `backend/internal/application/interview_progress/interview_progress_service.go`
- Modify: `backend/internal/application/interview_question/interview_question_service.go`
- Modify: `backend/cmd/server/main.go`
- Test: existing service tests in modified packages

- [ ] **Step 1: Create activity logger wrapper**

Create `backend/internal/application/activitylog/activity_logger.go`:

```go
package activitylog

import (
	"context"

	"github.com/google/uuid"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activityrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
)

type Logger interface {
	Log(ctx context.Context, userID string, action domainactivity.Action, targetType domainactivity.TargetType, targetID, targetLabel, summary, route string) error
}

type NoopLogger struct{}

func (NoopLogger) Log(context.Context, string, domainactivity.Action, domainactivity.TargetType, string, string, string, string) error { return nil }

type RepositoryLogger struct { repo activityrepo.Repository }

func NewRepositoryLogger(repo activityrepo.Repository) *RepositoryLogger { return &RepositoryLogger{repo: repo} }

func (l *RepositoryLogger) Log(ctx context.Context, userID string, action domainactivity.Action, targetType domainactivity.TargetType, targetID, targetLabel, summary, route string) error {
	activity, err := domainactivity.NewActivityLog(domainactivity.NewActivityLogParams{ID: uuid.NewString(), UserID: userID, Action: action, TargetType: targetType, TargetID: targetID, TargetLabel: targetLabel, Summary: summary, Route: route})
	if err != nil { return err }
	return l.repo.Create(ctx, activity)
}
```

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/activitylog
```

Expected: PASS or `no test files`.

- [ ] **Step 2: Add logger fields and WithActivityLogger methods**

For each service struct, add an `activityLogger activitylog.Logger` field and method:

```go
func (s *ServiceName) WithActivityLogger(logger activitylog.Logger) *ServiceName {
	if logger == nil {
		logger = activitylog.NoopLogger{}
	}
	s.activityLogger = logger
	return s
}

func (s *ServiceName) logActivity(ctx context.Context, userID string, action domainactivity.Action, targetType domainactivity.TargetType, targetID, targetLabel, summary, route string) error {
	if s.activityLogger == nil {
		return nil
	}
	return s.activityLogger.Log(ctx, userID, action, targetType, targetID, targetLabel, summary, route)
}
```

Apply to these service files with the real receiver type names used in each file:

- `backend/internal/application/resume_repo/resume_repo_service.go`
- `backend/internal/application/resume_version/resume_version_service.go`
- `backend/internal/application/resume_application/resume_application_service.go`
- `backend/internal/application/interview_progress/interview_progress_service.go`
- `backend/internal/application/interview_question/interview_question_service.go`

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/...
```

Expected: PASS because logger is no-op unless configured.

- [ ] **Step 3: Log resume activities**

In resume repo/version services, after successful create/update/delete/activate, call `logActivity`.

Use these summaries:

```go
"创建简历仓库 " + repo.Name
"更新简历仓库 " + repo.Name
"删除简历仓库 " + repo.Name
"创建简历版本 " + version.Name
"更新简历版本 " + version.Name
"删除简历版本 " + version.Name
"激活简历版本 " + version.Name
```

Use target type `domainactivity.TargetTypeResume` and routes:

```go
"/resume/repo/" + repo.ID
"/resume/version/" + version.ID
```

For delete operations, fetch the row before deletion to preserve `targetLabel` and use an empty route if the detail page no longer exists.

- [ ] **Step 4: Log application activities**

In `backend/internal/application/resume_application/resume_application_service.go`, log:

- Create: `创建投递记录 <companyName> / <jobTitle>`
- Update with status change: `<companyName> / <jobTitle> 投递状态更新为<中文状态>` with `ActionStatusChanged`
- Update without status change: `更新投递记录 <companyName> / <jobTitle>`
- Delete: `删除投递记录 <companyName> / <jobTitle>` with empty route

Add helper:

```go
func applicationLabel(a *domainapp.ResumeApplication) string { return a.CompanyName + " / " + a.JobTitle }
func applicationRoute(id string) string { return "/resume/application/" + id }
func applicationStatusLabel(status domainapp.Status) string { /* applied=已投递, viewed=已查看, written_test=笔试, first_interview=一面, second_interview=二面, final_interview=终面, hr_interview=HR面, offer=Offer, rejected=已拒绝, closed=已关闭 */ }
```

- [ ] **Step 5: Log interview activities**

In `backend/internal/application/interview_progress/interview_progress_service.go`, log:

- Create: `新增面试进展 <stage label>`
- Update: `更新面试进展为<stage label>`
- Delete: `删除面试进展 <stage label>`

Use target type `TargetTypeInterview`. If the progress has an `ApplicationID`, use route `"/resume/application/" + progress.ApplicationID`.

- [ ] **Step 6: Log question activities**

In `backend/internal/application/interview_question/interview_question_service.go`, log:

- CreateQuestion: `新增题目 <title>` with route `/interview-question-bank/question/<id>`
- UpdateQuestion: `更新题目 <title>`
- DeleteQuestion: fetch before deletion, then `删除题目 <title>` with empty route
- UpdateFavorite: `收藏题目 <title>` or `取消收藏题目 <title>`
- CreateOccurrence: `题目 <title> 关联到 <companyName> / <jobTitle>` with `ActionLinked`
- DeleteOccurrence: fetch occurrence before delete if repository supports it; if not, add repository method `FindByIDAndUserID` before logging deletion. Use summary `删除题目关联记录 <companyNameSnapshot> / <jobTitleSnapshot>`.

Do not log ListQuestions or GetQuestion.

- [ ] **Step 7: Wire repository logger in server**

Modify `backend/cmd/server/main.go` imports:

```go
activitylogapp "github.com/yilin/ai-for-backend/internal/application/activitylog"
```

After `activityLogRepo := infraDB.NewActivityLogRepo(db)`, add:

```go
activityLogger := activitylogapp.NewRepositoryLogger(activityLogRepo)
```

Apply to services:

```go
resumeRepoSvc.WithActivityLogger(activityLogger)
resumeVersionSvc.WithActivityLogger(activityLogger)
resumeApplicationSvc.WithActivityLogger(activityLogger)
interviewProgressSvc.WithActivityLogger(activityLogger)
interviewQuestionSvc.WithActivityLogger(activityLogger)
```

Do not add activity logging to `jobPostingSvc`.

- [ ] **Step 8: Add service tests for logging boundaries**

Add tests in existing service test files that use a fake logger and verify:

```go
// Resume application status update logs status_changed.
require.Equal(t, domainactivity.ActionStatusChanged, logger.entries[0].Action)
require.Contains(t, logger.entries[0].Summary, "投递状态更新为一面")

// Interview question GetQuestion does not log.
_, _, _, _, _, err := svc.GetQuestion(ctx, "u1", "q1")
require.NoError(t, err)
require.Empty(t, logger.entries)

// Job posting create has no logger field and does not write activity.
// This is covered by not modifying job_posting_service.go.
```

- [ ] **Step 9: Run application tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/...
```

Expected: PASS.

- [ ] **Step 10: Commit task 4**

```bash
cd /Users/yilin/project/ai-for
git add backend/internal/application/activitylog backend/internal/application/resume_repo backend/internal/application/resume_version backend/internal/application/resume_application backend/internal/application/interview_progress backend/internal/application/interview_question backend/cmd/server/main.go
git commit -m "feat: log dashboard activities from user operations"
```

---

## Task 5: Add Frontend Models, API, Store, and Mock Endpoints

**Files:**
- Create: `frontend/src/models/dashboard.ts`
- Create: `frontend/src/api/dashboard.ts`
- Create: `frontend/src/store/modules/dashboard.ts`
- Create: `frontend/mock/dashboard.ts`
- Modify: `frontend/src/models/index.ts`
- Modify: `frontend/src/api/index.ts`
- Modify: `frontend/mock/mock-server.ts`

- [ ] **Step 1: Add frontend dashboard models**

Create `frontend/src/models/dashboard.ts`:

```ts
export type DashboardRange = '7d' | '30d' | '90d' | 'all'
export type DashboardTargetType = 'application' | 'interview' | 'question' | 'resume'
export type DashboardTodoSource = 'manual' | 'system'

export interface DashboardMetrics {
  totalApplications: number
  interviewing: number
  pendingFollowUp: number
  offers: number
}

export interface DashboardStatusBucket {
  status: string
  label: string
  count: number
}

export interface DashboardTrendPoint {
  date: string
  count: number
}

export interface DashboardNameCount {
  name: string
  count: number
}

export interface DashboardTodo {
  id: string
  title: string
  description?: string
  dueAt?: string
  completedAt?: string
  targetType?: DashboardTargetType
  targetId?: string
  route?: string
  source: DashboardTodoSource
  createdAt: string
  updatedAt: string
}

export interface ActivityLog {
  id: string
  action: string
  targetType: DashboardTargetType
  targetId: string
  targetLabel: string
  summary: string
  route?: string
  createdAt: string
}

export interface DashboardResponse {
  range: DashboardRange
  metrics: DashboardMetrics
  statusPipeline: DashboardStatusBucket[]
  applicationTrend: DashboardTrendPoint[]
  companyRanking: DashboardNameCount[]
  jobTitleDistribution: DashboardNameCount[]
  sourceDistribution: DashboardNameCount[]
  todos: DashboardTodo[]
  activities: ActivityLog[]
}

export interface DashboardQueryParams {
  range: DashboardRange
}

export interface CreateDashboardTodoRequest {
  title: string
  description?: string
  dueAt?: string
  targetType?: DashboardTargetType
  targetId?: string
}

export interface UpdateDashboardTodoRequest extends Partial<CreateDashboardTodoRequest> {
  clearDueAt?: boolean
  clearTarget?: boolean
}
```

Modify `frontend/src/models/index.ts`:

```ts
export * from './dashboard'
```

- [ ] **Step 2: Add dashboard API wrapper**

Create `frontend/src/api/dashboard.ts`:

```ts
import request from '@/utils/request'
import {
  DashboardQueryParams,
  CreateDashboardTodoRequest,
  UpdateDashboardTodoRequest
} from '@/models'

export const getDashboard = (params: DashboardQueryParams) =>
  request({
    url: '/api/dashboard',
    method: 'get',
    params
  })

export const createDashboardTodo = (data: CreateDashboardTodoRequest) =>
  request({
    url: '/api/dashboard/todos',
    method: 'post',
    data
  })

export const updateDashboardTodo = (id: string, data: UpdateDashboardTodoRequest) =>
  request({
    url: `/api/dashboard/todos/${id}`,
    method: 'put',
    data
  })

export const completeDashboardTodo = (id: string) =>
  request({
    url: `/api/dashboard/todos/${id}/complete`,
    method: 'put'
  })

export const deleteDashboardTodo = (id: string) =>
  request({
    url: `/api/dashboard/todos/${id}`,
    method: 'delete'
  })
```

Modify `frontend/src/api/index.ts`:

```ts
export * from './dashboard'
```

- [ ] **Step 3: Add dashboard store**

Create `frontend/src/store/modules/dashboard.ts`:

```ts
import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import store from '@/store'
import {
  getDashboard,
  createDashboardTodo,
  updateDashboardTodo,
  completeDashboardTodo,
  deleteDashboardTodo
} from '@/api/dashboard'
import {
  DashboardRange,
  DashboardResponse,
  CreateDashboardTodoRequest,
  UpdateDashboardTodoRequest
} from '@/models'

export interface IDashboardState {
  dashboard: DashboardResponse | null
  range: DashboardRange
  loading: boolean
}

@Module({ dynamic: true, store, name: 'dashboard' })
class DashboardStore extends VuexModule implements IDashboardState {
  public dashboard: DashboardResponse | null = null
  public range: DashboardRange = '30d'
  public loading = false

  @Mutation
  private SET_DASHBOARD(dashboard: DashboardResponse | null) {
    this.dashboard = dashboard
  }

  @Mutation
  private SET_RANGE(range: DashboardRange) {
    this.range = range
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Action
  public async GetDashboard(payload: { range?: DashboardRange } = {}) {
    const range = payload.range || this.range
    this.SET_RANGE(range)
    this.SET_LOADING(true)
    try {
      const { data } = await getDashboard({ range })
      if (data) {
        this.SET_DASHBOARD(data)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateTodo(payload: CreateDashboardTodoRequest) {
    await createDashboardTodo(payload)
    await this.GetDashboard({ range: this.range })
  }

  @Action
  public async UpdateTodo(payload: { id: string, data: UpdateDashboardTodoRequest }) {
    await updateDashboardTodo(payload.id, payload.data)
    await this.GetDashboard({ range: this.range })
  }

  @Action
  public async CompleteTodo(id: string) {
    await completeDashboardTodo(id)
    await this.GetDashboard({ range: this.range })
  }

  @Action
  public async DeleteTodo(id: string) {
    await deleteDashboardTodo(id)
    await this.GetDashboard({ range: this.range })
  }
}

export const DashboardModule = getModule(DashboardStore)
```

- [ ] **Step 4: Add mock dashboard endpoints**

Create `frontend/mock/dashboard.ts`:

```ts
import { Request, Response } from 'express'

let manualTodos = [
  {
    id: 'todo-manual-1',
    title: '完善腾讯一面复盘',
    description: '整理面试问题和回答情况',
    dueAt: new Date().toISOString(),
    targetType: 'application',
    targetId: 'app-1',
    route: '/resume/application/app-1',
    source: 'manual',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
]

const systemTodos = [
  {
    id: 'todo-system-follow-up-1',
    title: '跟进字节 / 前端工程师',
    description: '投递超过 7 天未进入下一阶段',
    dueAt: new Date().toISOString(),
    targetType: 'application',
    targetId: 'app-2',
    route: '/resume/application/app-2',
    source: 'system',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: 'todo-system-interview-1',
    title: '准备腾讯 / 前端工程师 一面',
    description: '今天有面试安排',
    dueAt: new Date().toISOString(),
    targetType: 'application',
    targetId: 'app-1',
    route: '/resume/application/app-1',
    source: 'system',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
]

export const getDashboard = (req: Request, res: Response) => {
  const range = String(req.query.range || '30d')
  res.json({
    code: 20000,
    message: 'success',
    data: {
      range,
      metrics: { totalApplications: 42, interviewing: 8, pendingFollowUp: 5, offers: 2 },
      statusPipeline: [
        { status: 'applied', label: '已投递', count: 14 },
        { status: 'viewed', label: '已查看', count: 8 },
        { status: 'written_test', label: '笔试', count: 5 },
        { status: 'first_interview', label: '一面', count: 6 },
        { status: 'second_interview', label: '二面', count: 3 },
        { status: 'final_interview', label: '终面', count: 2 },
        { status: 'hr_interview', label: 'HR 面', count: 1 },
        { status: 'offer', label: 'Offer', count: 2 },
        { status: 'rejected', label: '拒绝', count: 4 },
        { status: 'closed', label: '关闭', count: 3 }
      ],
      applicationTrend: [
        { date: '05-12', count: 2 }, { date: '05-13', count: 5 }, { date: '05-14', count: 3 },
        { date: '05-15', count: 8 }, { date: '05-16', count: 4 }, { date: '05-17', count: 6 }, { date: '05-18', count: 3 }
      ],
      companyRanking: [
        { name: '腾讯', count: 6 }, { name: '字节', count: 5 }, { name: '阿里', count: 4 }, { name: '美团', count: 3 }
      ],
      jobTitleDistribution: [
        { name: '前端', count: 18 }, { name: '后端', count: 10 }, { name: '全栈', count: 6 }, { name: 'AI', count: 5 }, { name: '其他', count: 3 }
      ],
      sourceDistribution: [
        { name: 'Boss', count: 16 }, { name: '官网', count: 9 }, { name: '内推', count: 8 }, { name: '牛客', count: 5 }, { name: '其他', count: 4 }
      ],
      todos: [...manualTodos, ...systemTodos],
      activities: [
        { id: 'act-1', action: 'status_changed', targetType: 'application', targetId: 'app-1', targetLabel: '腾讯 / 前端工程师', summary: '腾讯 / 前端工程师 投递状态更新为一面', route: '/resume/application/app-1', createdAt: new Date().toISOString() },
        { id: 'act-2', action: 'updated', targetType: 'question', targetId: 'q-1', targetLabel: 'Vue 生命周期', summary: '更新题目 Vue 生命周期', route: '/interview-question-bank/question/q-1', createdAt: new Date().toISOString() },
        { id: 'act-3', action: 'created', targetType: 'interview', targetId: 'p-1', targetLabel: '一面', summary: '新增面试进展 一面', route: '/resume/application/app-1', createdAt: new Date().toISOString() }
      ]
    }
  })
}

export const getDashboardTodos = (_req: Request, res: Response) => {
  res.json({ code: 20000, message: 'success', data: { list: manualTodos } })
}

export const createDashboardTodo = (req: Request, res: Response) => {
  const now = new Date().toISOString()
  const todo = { id: `todo-manual-${Date.now()}`, source: 'manual', completedAt: undefined, createdAt: now, updatedAt: now, ...req.body }
  manualTodos.unshift(todo)
  res.json({ code: 20000, message: 'success', data: todo })
}

export const updateDashboardTodo = (req: Request, res: Response) => {
  const index = manualTodos.findIndex(item => item.id === req.params.id)
  if (index >= 0) manualTodos[index] = { ...manualTodos[index], ...req.body, updatedAt: new Date().toISOString() }
  res.json({ code: 20000, message: 'success', data: manualTodos[index] })
}

export const completeDashboardTodo = (req: Request, res: Response) => {
  const index = manualTodos.findIndex(item => item.id === req.params.id)
  if (index >= 0) manualTodos[index] = { ...manualTodos[index], completedAt: new Date().toISOString(), updatedAt: new Date().toISOString() }
  res.json({ code: 20000, message: 'success', data: manualTodos[index] })
}

export const deleteDashboardTodo = (req: Request, res: Response) => {
  manualTodos = manualTodos.filter(item => item.id !== req.params.id)
  res.json({ code: 20000, message: 'success', data: null })
}
```

Modify `frontend/mock/mock-server.ts` imports:

```ts
import * as dashboard from './dashboard'
```

Add routes before the 404 handler:

```ts
app.get('/api/dashboard', dashboard.getDashboard)
app.get('/api/dashboard/todos', dashboard.getDashboardTodos)
app.post('/api/dashboard/todos', dashboard.createDashboardTodo)
app.put('/api/dashboard/todos/:id', dashboard.updateDashboardTodo)
app.put('/api/dashboard/todos/:id/complete', dashboard.completeDashboardTodo)
app.delete('/api/dashboard/todos/:id', dashboard.deleteDashboardTodo)
```

- [ ] **Step 5: Run frontend type checks for new model/api/store/mock**

Run:

```bash
cd /Users/yilin/project/ai-for/frontend && npm run lint
```

Expected: PASS. If lint fails because existing unrelated files already fail, capture the failing output and only fix failures introduced by dashboard files.

- [ ] **Step 6: Commit task 5**

```bash
cd /Users/yilin/project/ai-for
git add frontend/src/models/dashboard.ts frontend/src/models/index.ts frontend/src/api/dashboard.ts frontend/src/api/index.ts frontend/src/store/modules/dashboard.ts frontend/mock/dashboard.ts frontend/mock/mock-server.ts
git commit -m "feat: add dashboard frontend data layer and mocks"
```

---

## Task 6: Build Dashboard UI Components

**Files:**
- Create: `frontend/src/views/dashboard/components/MetricCards.vue`
- Create: `frontend/src/views/dashboard/components/PipelineChart.vue`
- Create: `frontend/src/views/dashboard/components/TodoPanel.vue`
- Create: `frontend/src/views/dashboard/components/DistributionCharts.vue`
- Create: `frontend/src/views/dashboard/components/ActivityFeed.vue`
- Modify: `frontend/src/views/dashboard/index.vue`

- [ ] **Step 1: Replace dashboard page shell**

Modify `frontend/src/views/dashboard/index.vue` to render the job dashboard directly:

```vue
<template>
  <div class="job-dashboard-container">
    <div class="dashboard-header">
      <div>
        <h2>求职工作台</h2>
        <p>聚合投递进展、今日待办、统计分析和最近动态</p>
      </div>
      <el-radio-group v-model="selectedRange" size="small" @change="handleRangeChange">
        <el-radio-button label="7d">近7天</el-radio-button>
        <el-radio-button label="30d">近30天</el-radio-button>
        <el-radio-button label="90d">近90天</el-radio-button>
        <el-radio-button label="all">全部</el-radio-button>
      </el-radio-group>
    </div>

    <div v-loading="loading">
      <metric-cards :metrics="metrics" />
      <el-row :gutter="16" class="dashboard-row">
        <el-col :xs="24" :lg="16">
          <pipeline-chart :pipeline="statusPipeline" :trend="applicationTrend" />
        </el-col>
        <el-col :xs="24" :lg="8">
          <todo-panel :todos="todos" @create="handleCreateTodo" @update="handleUpdateTodo" @complete="handleCompleteTodo" @delete="handleDeleteTodo" @open="openRoute" />
        </el-col>
      </el-row>
      <distribution-charts :companies="companyRanking" :jobs="jobTitleDistribution" :sources="sourceDistribution" />
      <activity-feed :activities="activities" @open="openRoute" />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import MetricCards from './components/MetricCards.vue'
import PipelineChart from './components/PipelineChart.vue'
import TodoPanel from './components/TodoPanel.vue'
import DistributionCharts from './components/DistributionCharts.vue'
import ActivityFeed from './components/ActivityFeed.vue'
import { DashboardModule } from '@/store/modules/dashboard'
import { DashboardRange, CreateDashboardTodoRequest, UpdateDashboardTodoRequest } from '@/models'

@Component({
  name: 'Dashboard',
  components: { MetricCards, PipelineChart, TodoPanel, DistributionCharts, ActivityFeed }
})
export default class extends Vue {
  private selectedRange: DashboardRange = '30d'

  get loading() { return DashboardModule.loading }
  get dashboard() { return DashboardModule.dashboard }
  get metrics() { return this.dashboard ? this.dashboard.metrics : { totalApplications: 0, interviewing: 0, pendingFollowUp: 0, offers: 0 } }
  get statusPipeline() { return this.dashboard ? this.dashboard.statusPipeline : [] }
  get applicationTrend() { return this.dashboard ? this.dashboard.applicationTrend : [] }
  get companyRanking() { return this.dashboard ? this.dashboard.companyRanking : [] }
  get jobTitleDistribution() { return this.dashboard ? this.dashboard.jobTitleDistribution : [] }
  get sourceDistribution() { return this.dashboard ? this.dashboard.sourceDistribution : [] }
  get todos() { return this.dashboard ? this.dashboard.todos : [] }
  get activities() { return this.dashboard ? this.dashboard.activities : [] }

  created() { this.loadDashboard() }

  private async loadDashboard() { await DashboardModule.GetDashboard({ range: this.selectedRange }) }
  private handleRangeChange(range: DashboardRange) { this.selectedRange = range; this.loadDashboard() }
  private handleCreateTodo(payload: CreateDashboardTodoRequest) { return DashboardModule.CreateTodo(payload) }
  private handleUpdateTodo(payload: { id: string, data: UpdateDashboardTodoRequest }) { return DashboardModule.UpdateTodo(payload) }
  private handleCompleteTodo(id: string) { return DashboardModule.CompleteTodo(id) }
  private handleDeleteTodo(id: string) { return DashboardModule.DeleteTodo(id) }
  private openRoute(route?: string) { if (route) this.$router.push(route).catch(() => {}) }
}
</script>

<style lang="scss" scoped>
.job-dashboard-container {
  padding: 24px;
  background: #f5f7fb;
  min-height: calc(100vh - 84px);

  .dashboard-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 18px;

    h2 { margin: 0 0 6px; font-size: 22px; color: #1f2d3d; }
    p { margin: 0; color: #8c98a8; }
  }

  .dashboard-row { margin-bottom: 16px; }
}

@media (max-width: 768px) {
  .job-dashboard-container {
    padding: 12px;
    .dashboard-header { display: block; }
    .el-radio-group { margin-top: 12px; }
  }
}
</style>
```

- [ ] **Step 2: Create MetricCards component**

Create `frontend/src/views/dashboard/components/MetricCards.vue`:

```vue
<template>
  <el-row :gutter="16" class="metric-cards">
    <el-col v-for="card in cards" :key="card.label" :xs="24" :sm="12" :lg="6">
      <el-card shadow="never" class="metric-card">
        <div class="metric-icon" :class="card.type"><i :class="card.icon" /></div>
        <div>
          <div class="metric-label">{{ card.label }}</div>
          <div class="metric-value">{{ card.value }}</div>
        </div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { DashboardMetrics } from '@/models'

@Component({ name: 'MetricCards' })
export default class extends Vue {
  @Prop({ required: true }) private metrics!: DashboardMetrics

  get cards() {
    return [
      { label: '总投递', value: this.metrics.totalApplications, icon: 'el-icon-s-promotion', type: 'blue' },
      { label: '面试中', value: this.metrics.interviewing, icon: 'el-icon-user-solid', type: 'green' },
      { label: '待跟进', value: this.metrics.pendingFollowUp, icon: 'el-icon-time', type: 'orange' },
      { label: 'Offer', value: this.metrics.offers, icon: 'el-icon-suitcase', type: 'purple' }
    ]
  }
}
</script>

<style lang="scss" scoped>
.metric-cards { margin-bottom: 16px; }
.metric-card ::v-deep .el-card__body { display: flex; align-items: center; gap: 14px; }
.metric-icon { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; }
.metric-icon.blue { color: #409eff; background: #ecf5ff; }
.metric-icon.green { color: #67c23a; background: #f0f9eb; }
.metric-icon.orange { color: #e6a23c; background: #fdf6ec; }
.metric-icon.purple { color: #7c3aed; background: #f3e8ff; }
.metric-label { color: #8c98a8; font-size: 13px; }
.metric-value { margin-top: 6px; font-size: 24px; font-weight: 600; color: #1f2d3d; }
</style>
```

- [ ] **Step 3: Create chart components**

Create `PipelineChart.vue` and `DistributionCharts.vue` using ECharts. Use the existing dashboard chart components as the pattern: initialize chart in `mounted`, dispose in `beforeDestroy`, and watch props to refresh options.

For `PipelineChart.vue`, required props:

```ts
@Prop({ default: () => [] }) private pipeline!: DashboardStatusBucket[]
@Prop({ default: () => [] }) private trend!: DashboardTrendPoint[]
```

The chart options must include:

```ts
{
  tooltip: { trigger: 'axis' },
  legend: { data: ['阶段数量', '投递趋势'] },
  xAxis: [{ type: 'category', data: this.pipeline.map(item => item.label) }],
  yAxis: [{ type: 'value' }],
  series: [{ name: '阶段数量', type: 'bar', data: this.pipeline.map(item => item.count) }]
}
```

For `DistributionCharts.vue`, required props:

```ts
@Prop({ default: () => [] }) private companies!: DashboardNameCount[]
@Prop({ default: () => [] }) private jobs!: DashboardNameCount[]
@Prop({ default: () => [] }) private sources!: DashboardNameCount[]
```

Render three `el-card` columns with one chart each.

- [ ] **Step 4: Create TodoPanel component**

Create `frontend/src/views/dashboard/components/TodoPanel.vue`. It must:

- Render system and manual todos in one list.
- Show `el-tag` with `系统` or `手动`.
- Disable edit/delete for `source === 'system'`.
- Emit `create`, `update`, `complete`, `delete`, and `open` events.
- Use an `el-dialog` with title, description, dueAt, targetType, targetId fields.

Use event payload types:

```ts
this.$emit('create', this.form)
this.$emit('update', { id: this.editingId, data: this.form })
this.$emit('complete', todo.id)
this.$emit('delete', todo.id)
this.$emit('open', todo.route)
```

- [ ] **Step 5: Create ActivityFeed component**

Create `frontend/src/views/dashboard/components/ActivityFeed.vue`. It must:

- Render an `el-card` titled `最近动态`.
- Render an `el-timeline` of activities.
- Use `activity.summary` as primary text.
- Display formatted `createdAt`.
- If `activity.route` exists, make row clickable and emit `open`.
- If `route` is absent, render disabled text without click affordance.

- [ ] **Step 6: Run frontend lint and unit tests**

Run:

```bash
cd /Users/yilin/project/ai-for/frontend && npm run lint && npm run test:unit
```

Expected: PASS. If unrelated existing tests fail, record the failing test names and only fix dashboard-caused failures.

- [ ] **Step 7: Manually verify with mock server**

Run:

```bash
cd /Users/yilin/project/ai-for/frontend && npm run serve
```

Open `/dashboard` and verify:

- Range defaults to `近30天`.
- KPI cards render.
- Pipeline and distribution charts render.
- Manual and system todos render.
- New manual todo can be created.
- Manual todo can be completed.
- System todo cannot be edited/deleted.
- Activity with route jumps to target page.
- Activity without route is not clickable.

- [ ] **Step 8: Commit task 6**

```bash
cd /Users/yilin/project/ai-for
git add frontend/src/views/dashboard/index.vue frontend/src/views/dashboard/components
git commit -m "feat: replace dashboard with job workspace UI"
```

---

## Task 7: Full Verification and Cleanup

**Files:**
- Modify only files with failures found during verification.

- [ ] **Step 1: Run backend full test suite**

```bash
cd /Users/yilin/project/ai-for/backend && go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run frontend lint and unit tests**

```bash
cd /Users/yilin/project/ai-for/frontend && npm run lint && npm run test:unit
```

Expected: PASS.

- [ ] **Step 3: Run frontend production build**

```bash
cd /Users/yilin/project/ai-for/frontend && npm run build:prod
```

Expected: PASS.

- [ ] **Step 4: Manual UI verification**

With `npm run serve` running from `frontend`, verify `/dashboard`:

- Top range selector works for 7d, 30d, 90d, all.
- Cards, charts, todos, activities all update after range switch.
- Todo create/edit/complete/delete refreshes the page data.
- Dynamic click routes resolve for application and question examples.
- On a medium-width browser, cards wrap without overlapping.

- [ ] **Step 5: Check git diff for accidental unrelated changes**

```bash
cd /Users/yilin/project/ai-for && git status --short && git diff --stat
```

Expected: only dashboard/backend/frontend files listed in this plan and the reviewed spec/plan docs.

- [ ] **Step 6: Commit verification fixes if any**

If Step 1-4 required fixes, commit only those files:

```bash
cd /Users/yilin/project/ai-for
git add <fixed-files>
git commit -m "fix: stabilize job dashboard implementation"
```

If no fixes were required, do not create an empty commit.

---

## Self-Review

### Spec coverage

- Dashboard page layout B: Task 6.
- Time range 7d/30d/90d/all, default 30d: Task 3 backend and Task 6 frontend.
- Backend aggregation endpoint: Task 3.
- Manual todos CRUD and completion: Tasks 1, 2, 3, 5, 6.
- Automatic todos for overdue follow-up and interview day: Task 3.
- Activity logs for resume/question/application/interview operations: Task 4.
- Job posting creation excluded from activity logging: Task 4 explicitly does not modify job posting service.
- Question ordinary views excluded from logging: Task 4 explicitly does not log ListQuestions/GetQuestion.
- Company/job/source charts: Task 3 aggregation, Task 5 models/mock, Task 6 UI.
- Frontend mock same structure as backend: Task 5.
- Verification: Task 7.

### Placeholder scan

No `TBD`, `TODO`, `implement later`, or unspecified “handle edge cases” instructions remain. Helper functions are named with exact behavior and expected inputs/outputs.

### Type consistency

Backend DTO names, repository interfaces, handler method names, frontend model names, API wrapper names, and Vuex action names are consistent across tasks. Activity and todo target/source strings match the reviewed spec.
