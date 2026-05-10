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
	"gorm.io/gorm"
)

func migrateResumeTables(t *testing.T, db interface {
	AutoMigrate(...interface{}) error
	Exec(string, ...interface{}) *gorm.DB
}) {
	t.Helper()
	require.NoError(t, db.AutoMigrate(&domainrepo.ResumeRepo{}, &domainversion.ResumeVersion{}))
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_repos_user_name_unique ON resume_repos(user_id, name) WHERE deleted_at IS NULL").Error)
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_versions_repo_version_unique ON resume_versions(repo_id, version_num) WHERE deleted_at IS NULL").Error)
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_versions_repo_default_unique ON resume_versions(repo_id) WHERE is_default = TRUE AND deleted_at IS NULL").Error)
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

func TestIntegrationResumeRepoRepo_UpdateMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()

	err := repo.Update(ctx, &domainrepo.ResumeRepo{ID: uuid.NewString(), UserID: uuid.NewString(), Name: "Missing"})
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationResumeRepoRepo_DeleteMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.NewString(), uuid.NewString())
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationResumeRepoRepo_DuplicateNameReturnsDuplicate(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	require.NoError(t, repo.Create(ctx, &domainrepo.ResumeRepo{ID: uuid.NewString(), UserID: userID, Name: "Backend"}))
	err := repo.Create(ctx, &domainrepo.ResumeRepo{ID: uuid.NewString(), UserID: userID, Name: "Backend"})
	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
}
