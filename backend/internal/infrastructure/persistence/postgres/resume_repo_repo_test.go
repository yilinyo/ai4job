package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func migrateResumeTables(t *testing.T, db interface{ AutoMigrate(...interface{}) error }) {
	t.Helper()
	require.NoError(t, db.AutoMigrate(&domainrepo.ResumeRepo{}, &domainversion.ResumeVersion{}))
}

func TestIntegrationResumeRepoRepo_CreateListAndFindByOwner(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	r := &domainrepo.ResumeRepo{ID: uuid.NewString(), UserID: userID, Name: "Backend", Description: "main"}
	require.NoError(t, repo.Create(ctx, r))

	found, err := repo.FindByIDAndUserID(ctx, r.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Backend", found.Name)

	rows, total, err := repo.ListByUser(ctx, domainrepo.ListFilter{UserID: userID, Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, rows, 1)
}
