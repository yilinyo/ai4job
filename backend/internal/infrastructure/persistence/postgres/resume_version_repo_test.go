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
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestIntegrationResumeVersionRepo_CreateCountAndSetActive(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()
	repoID := uuid.NewString()
	require.NoError(t, db.Create(testResumeRepo(repoID, uuid.NewString(), "Backend", domainrepo.JobTypeCampus)).Error)

	v1 := &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Version: "v1.0.0", Title: "v1", Content: "{}", Remark: "first", IsActive: false}
	require.NoError(t, repo.Create(ctx, v1))

	count, err := repo.CountByRepoID(ctx, repoID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	require.NoError(t, repo.SetActiveByID(ctx, repoID, v1.ID))
	found, err := repo.FindByID(ctx, v1.ID)
	require.NoError(t, err)
	assert.True(t, found.IsActive)
	assert.Equal(t, "first", found.Remark)
}

func TestIntegrationResumeVersionRepo_UpdateMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()

	err := repo.Update(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: uuid.NewString(), Version: "v1.0.0", Title: "missing", Content: "{}"})
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationResumeVersionRepo_DeleteMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.NewString(), uuid.NewString())
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationResumeVersionRepo_DuplicateVersionReturnsDuplicate(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()
	repoID := uuid.NewString()

	require.NoError(t, db.Create(testResumeRepo(repoID, uuid.NewString(), "Backend", domainrepo.JobTypeCampus)).Error)
	require.NoError(t, repo.Create(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Version: "v1.0.0", Title: "v1", Content: "{}"}))
	err := repo.Create(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Version: "v1.0.0", Title: "v1-copy", Content: "{}"})
	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
}
