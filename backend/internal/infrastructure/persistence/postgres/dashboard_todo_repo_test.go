package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	dashboardtodorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func migrateDashboardTodoTables(t *testing.T, db interface{ AutoMigrate(...interface{}) error }) {
	t.Helper()
	require.NoError(t, db.AutoMigrate(&domaintodo.DashboardTodo{}))
}

func testDashboardTodo(id, userID, title string, dueAt *time.Time) *domaintodo.DashboardTodo {
	todo, err := domaintodo.NewDashboardTodo(domaintodo.NewDashboardTodoParams{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: title + " description",
		DueAt:       dueAt,
	})
	if err != nil {
		panic(err)
	}
	return todo
}

func TestDashboardTodoRepoContracts(t *testing.T) {
	var _ dashboardtodorepo.Repository = (*infrapostgres.DashboardTodoRepo)(nil)
}

func TestDashboardTodoRepo_CreateFindUpdateDelete(t *testing.T) {
	db := setupTestDB(t)
	migrateDashboardTodoTables(t, db)
	repo := infrapostgres.NewDashboardTodoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	dueAt := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Microsecond)
	todo := testDashboardTodo(uuid.NewString(), userID, "Follow up", &dueAt)

	require.NoError(t, repo.Create(ctx, todo))

	found, err := repo.FindByIDAndUserID(ctx, todo.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, todo.ID, found.ID)
	assert.Equal(t, "Follow up", found.Title)
	require.NotNil(t, found.DueAt)
	assert.Equal(t, dueAt, found.DueAt.UTC().Truncate(time.Microsecond))

	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	found.Title = "Follow up again"
	found.Description = "updated description"
	found.CompletedAt = &completedAt
	require.NoError(t, repo.Update(ctx, found))

	updated, err := repo.FindByIDAndUserID(ctx, todo.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Follow up again", updated.Title)
	assert.Equal(t, "updated description", updated.Description)
	require.NotNil(t, updated.CompletedAt)
	assert.Equal(t, completedAt, updated.CompletedAt.UTC().Truncate(time.Microsecond))

	require.NoError(t, repo.Delete(ctx, todo.ID, userID))
	_, err = repo.FindByIDAndUserID(ctx, todo.ID, userID)
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestDashboardTodoRepo_ListByUserFiltersCompleted(t *testing.T) {
	db := setupTestDB(t)
	migrateDashboardTodoTables(t, db)
	repo := infrapostgres.NewDashboardTodoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	otherUserID := uuid.NewString()

	openTodo := testDashboardTodo(uuid.NewString(), userID, "Open", nil)
	completedTodo := testDashboardTodo(uuid.NewString(), userID, "Done", nil)
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	completedTodo.CompletedAt = &completedAt
	require.NoError(t, repo.Create(ctx, openTodo))
	require.NoError(t, repo.Create(ctx, completedTodo))
	require.NoError(t, repo.Create(ctx, testDashboardTodo(uuid.NewString(), otherUserID, "Other", nil)))

	completed := false
	rows, err := repo.ListByUser(ctx, dashboardtodorepo.ListFilter{UserID: userID, Completed: &completed})

	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, openTodo.ID, rows[0].ID)

	completed = true
	rows, err = repo.ListByUser(ctx, dashboardtodorepo.ListFilter{UserID: userID, Completed: &completed})

	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, completedTodo.ID, rows[0].ID)
}

func TestDashboardTodoRepo_ListByUserTodayOnlyIncludesRelevantTodos(t *testing.T) {
	db := setupTestDB(t)
	migrateDashboardTodoTables(t, db)
	repo := infrapostgres.NewDashboardTodoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Truncate(time.Microsecond)
	todayDue := startOfToday.Add(12 * time.Hour)
	overdueDue := startOfToday.Add(-24 * time.Hour)
	futureDue := startOfToday.Add(48 * time.Hour)
	completedBeforeTodayAt := startOfToday.Add(-time.Hour)

	undated := testDashboardTodo(uuid.NewString(), userID, "Undated", nil)
	overdue := testDashboardTodo(uuid.NewString(), userID, "Overdue", &overdueDue)
	today := testDashboardTodo(uuid.NewString(), userID, "Today", &todayDue)
	future := testDashboardTodo(uuid.NewString(), userID, "Future", &futureDue)
	completedBeforeToday := testDashboardTodo(uuid.NewString(), userID, "Completed Yesterday", &todayDue)
	completedBeforeToday.CompletedAt = &completedBeforeTodayAt
	require.NoError(t, repo.Create(ctx, undated))
	require.NoError(t, repo.Create(ctx, overdue))
	require.NoError(t, repo.Create(ctx, today))
	require.NoError(t, repo.Create(ctx, future))
	require.NoError(t, repo.Create(ctx, completedBeforeToday))

	rows, err := repo.ListByUser(ctx, dashboardtodorepo.ListFilter{UserID: userID, TodayOnly: true})

	require.NoError(t, err)
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	assert.ElementsMatch(t, []string{undated.ID, overdue.ID, today.ID}, ids)
}

func TestDashboardTodoRepo_UpdateAndDeleteMissingReturnNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateDashboardTodoTables(t, db)
	repo := infrapostgres.NewDashboardTodoRepo(db)
	ctx := context.Background()

	err := repo.Update(ctx, testDashboardTodo(uuid.NewString(), uuid.NewString(), "Missing", nil))
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)

	err = repo.Delete(ctx, uuid.NewString(), uuid.NewString())
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}
