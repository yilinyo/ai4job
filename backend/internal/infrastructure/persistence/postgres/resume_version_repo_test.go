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

func TestIntegrationResumeVersionRepo_CreateMaxAndDefault(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()
	repoID := uuid.NewString()

	v1 := &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Title: "v1", Content: "{}", VersionNum: 1, IsDefault: true}
	require.NoError(t, repo.Create(ctx, v1))

	max, err := repo.MaxVersionNum(ctx, repoID)
	require.NoError(t, err)
	assert.Equal(t, 1, max)

	require.NoError(t, repo.ClearDefaultByRepoID(ctx, repoID))
	found, err := repo.FindByIDAndRepoID(ctx, v1.ID, repoID)
	require.NoError(t, err)
	assert.False(t, found.IsDefault)
}

func TestIntegrationResumeVersionRepo_UpdateMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()

	err := repo.Update(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: uuid.NewString(), Title: "missing", Content: "{}", VersionNum: 1})
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationResumeVersionRepo_DuplicateVersionNumReturnsDuplicate(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeVersionRepo(db)
	ctx := context.Background()
	repoID := uuid.NewString()

	require.NoError(t, db.Create(&domainrepo.ResumeRepo{ID: repoID, UserID: uuid.NewString(), Name: "Backend"}).Error)
	require.NoError(t, repo.Create(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Title: "v1", Content: "{}", VersionNum: 1}))
	err := repo.Create(ctx, &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Title: "v1-copy", Content: "{}", VersionNum: 1})
	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
}
