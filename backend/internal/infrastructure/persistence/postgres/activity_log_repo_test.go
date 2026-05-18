package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	activitylogrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func migrateActivityLogTables(t *testing.T, db interface{ AutoMigrate(...interface{}) error }) {
	t.Helper()
	require.NoError(t, db.AutoMigrate(&domainactivity.ActivityLog{}))
}

func testActivityLog(id, userID string, createdAt time.Time) *domainactivity.ActivityLog {
	log, err := domainactivity.NewActivityLog(domainactivity.NewActivityLogParams{
		ID:          id,
		UserID:      userID,
		Action:      domainactivity.ActionCreated,
		TargetType:  domainactivity.TargetTypeApplication,
		TargetID:    uuid.NewString(),
		TargetLabel: "Example Application",
		Summary:     "Created Example Application",
		Route:       "/applications/" + uuid.NewString(),
	})
	if err != nil {
		panic(err)
	}
	log.CreatedAt = createdAt
	return log
}

func TestActivityLogRepoContracts(t *testing.T) {
	var _ activitylogrepo.Repository = (*infrapostgres.ActivityLogRepo)(nil)
}

func TestActivityLogRepo_CreateAndListByUser(t *testing.T) {
	db := setupTestDB(t)
	migrateActivityLogTables(t, db)
	repo := infrapostgres.NewActivityLogRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	otherUserID := uuid.NewString()
	base := time.Now().UTC().Truncate(time.Microsecond)
	older := testActivityLog(uuid.NewString(), userID, base.Add(-time.Hour))
	newer := testActivityLog(uuid.NewString(), userID, base)
	other := testActivityLog(uuid.NewString(), otherUserID, base.Add(time.Hour))

	require.NoError(t, repo.Create(ctx, older))
	require.NoError(t, repo.Create(ctx, newer))
	require.NoError(t, repo.Create(ctx, other))

	rows, err := repo.ListByUser(ctx, activitylogrepo.ListFilter{UserID: userID, Limit: 10})

	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, newer.ID, rows[0].ID)
	assert.Equal(t, older.ID, rows[1].ID)
	assert.Equal(t, "Created Example Application", rows[0].Summary)
}

func TestActivityLogRepo_ListByUserFiltersStartAtAndLimit(t *testing.T) {
	db := setupTestDB(t)
	migrateActivityLogTables(t, db)
	repo := infrapostgres.NewActivityLogRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	base := time.Now().UTC().Truncate(time.Microsecond)
	oldest := testActivityLog(uuid.NewString(), userID, base.Add(-2*time.Hour))
	middle := testActivityLog(uuid.NewString(), userID, base.Add(-time.Hour))
	newest := testActivityLog(uuid.NewString(), userID, base)

	require.NoError(t, repo.Create(ctx, oldest))
	require.NoError(t, repo.Create(ctx, middle))
	require.NoError(t, repo.Create(ctx, newest))

	startAt := middle.CreatedAt
	rows, err := repo.ListByUser(ctx, activitylogrepo.ListFilter{UserID: userID, StartAt: &startAt, Limit: 1})

	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, newest.ID, rows[0].ID)
}

func TestActivityLogRepo_ListByUserZeroLimitDefaultsToTen(t *testing.T) {
	db := setupTestDB(t)
	migrateActivityLogTables(t, db)
	repo := infrapostgres.NewActivityLogRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	base := time.Now().UTC().Truncate(time.Microsecond)

	for i := 0; i < 11; i++ {
		require.NoError(t, repo.Create(ctx, testActivityLog(uuid.NewString(), userID, base.Add(time.Duration(i)*time.Minute))))
	}

	rows, err := repo.ListByUser(ctx, activitylogrepo.ListFilter{UserID: userID, Limit: 0})

	require.NoError(t, err)
	require.Len(t, rows, 10)
}
