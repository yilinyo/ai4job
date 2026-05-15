package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"
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
	require.NoError(t, db.Exec("CREATE INDEX IF NOT EXISTS idx_resume_repos_user_job_type ON resume_repos(user_id, job_type)").Error)
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_versions_repo_version_unique ON resume_versions(repo_id, version) WHERE deleted_at IS NULL").Error)
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_versions_repo_active_unique ON resume_versions(repo_id) WHERE is_active = TRUE AND deleted_at IS NULL").Error)
}

func testResumeRepo(id, userID, name string, jobType domainrepo.JobType) *domainrepo.ResumeRepo {
	return &domainrepo.ResumeRepo{
		ID:             id,
		UserID:         userID,
		Name:           name,
		JobType:        jobType,
		TargetPosition: name + " Engineer",
		Description:    "main",
		VersionCount:   0,
	}
}

func TestIntegrationResumeRepoRepo_CreateListAndFindByOwner(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	r := testResumeRepo(uuid.NewString(), userID, "Backend", domainrepo.JobTypeCampus)
	require.NoError(t, repo.Create(ctx, r))

	found, err := repo.FindByIDAndUserID(ctx, r.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Backend", found.Name)
	assert.Equal(t, domainrepo.JobTypeCampus, found.JobType)
	assert.Equal(t, "Backend Engineer", found.TargetPosition)
	assert.Equal(t, 0, found.VersionCount)

	rows, total, err := repo.ListByUser(ctx, resumereporepo.ListFilter{UserID: userID, Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, rows, 1)
}

func TestIntegrationResumeRepoRepo_UpdateMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()

	err := repo.Update(ctx, testResumeRepo(uuid.NewString(), uuid.NewString(), "Missing", domainrepo.JobTypeCampus))
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

	require.NoError(t, repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Backend", domainrepo.JobTypeCampus)))
	err := repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Backend", domainrepo.JobTypeSocial))
	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
}

func TestIntegrationResumeRepoRepo_ListByUserFiltersByJobType(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	require.NoError(t, repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Campus", domainrepo.JobTypeCampus)))
	require.NoError(t, repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Social", domainrepo.JobTypeSocial)))

	rows, total, err := repo.ListByUser(ctx, resumereporepo.ListFilter{UserID: userID, JobType: domainrepo.JobTypeCampus, Page: 1, PageSize: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, rows, 1)
	assert.Equal(t, domainrepo.JobTypeCampus, rows[0].JobType)
}

func TestIntegrationResumeRepoRepo_ListByUserFiltersByKeyword(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	require.NoError(t, repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Campus", domainrepo.JobTypeCampus)))
	require.NoError(t, repo.Create(ctx, testResumeRepo(uuid.NewString(), userID, "Social", domainrepo.JobTypeSocial)))

	rows, total, err := repo.ListByUser(ctx, resumereporepo.ListFilter{UserID: userID, Keyword: "camp", Page: 1, PageSize: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, rows, 1)
	assert.Equal(t, "Campus", rows[0].Name)
}

func TestIntegrationResumeRepoRepo_IncrementVersionCount(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	r := testResumeRepo(uuid.NewString(), userID, "Backend", domainrepo.JobTypeCampus)
	r.VersionCount = 1
	require.NoError(t, repo.Create(ctx, r))

	require.NoError(t, repo.IncrementVersionCount(ctx, r.ID, 2))
	found, err := repo.FindByIDAndUserID(ctx, r.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, 3, found.VersionCount)

	require.NoError(t, repo.IncrementVersionCount(ctx, r.ID, -10))
	found, err = repo.FindByIDAndUserID(ctx, r.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, found.VersionCount)
}

func TestIntegrationResumeRepoRepo_IncrementVersionCountMissingReturnsNotFound(t *testing.T) {
	db := setupTestDB(t)
	migrateResumeTables(t, db)
	repo := infrapostgres.NewResumeRepoRepo(db)
	ctx := context.Background()

	err := repo.IncrementVersionCount(ctx, uuid.NewString(), 1)

	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}
