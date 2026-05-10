package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		tcpostgres.WithDatabase("test_db"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := infrapostgres.NewDB(dsn, true)
	require.NoError(t, err)

	err = db.AutoMigrate(&domainuser.User{})
	require.NoError(t, err)

	return db
}

func TestIntegrationUserRepo_CreateAndFind(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "testuser",
		Email:    "test@example.com",
	}
	_ = u.SetPassword("password123")

	err := repo.Create(ctx, u)
	require.NoError(t, err)

	found, err := repo.FindByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
	assert.Equal(t, "test@example.com", found.Email)
	assert.NotNil(t, found.EducationExperiences)
	assert.Empty(t, found.EducationExperiences)

	var rawValue string
	err = db.Raw("SELECT education_experiences::text FROM users WHERE id = ?", u.ID).Scan(&rawValue).Error
	require.NoError(t, err)
	assert.Equal(t, "[]", rawValue)
}

func TestIntegrationUser_BeforeSaveNormalizesNilEducationExperiences(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := setupTestDB(t)

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "hookuser",
		Email:    "hook@example.com",
	}
	_ = u.SetPassword("password123")

	require.NoError(t, db.Create(u).Error)

	var rawValue string
	err := db.Raw("SELECT education_experiences::text FROM users WHERE id = ?", u.ID).Scan(&rawValue).Error
	require.NoError(t, err)
	assert.Equal(t, "[]", rawValue)
}

func TestIntegrationUserRepo_SoftDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "deleteuser",
		Email:    "delete@example.com",
	}
	_ = u.SetPassword("password123")
	require.NoError(t, repo.Create(ctx, u))

	// soft delete via GORM
	err := db.Delete(u).Error
	require.NoError(t, err)

	// should not be found after soft delete
	_, err = repo.FindByUsername(ctx, "deleteuser")
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
}

func TestIntegrationUserRepo_UpdateNormalizesNilEducationExperiences(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:                   uuid.NewString(),
		Username:             "updateuser",
		Email:                "update@example.com",
		EducationExperiences: []domainuser.EducationExperience{{School: "Old School"}},
	}
	_ = u.SetPassword("password123")
	require.NoError(t, repo.Create(ctx, u))

	u.EducationExperiences = nil
	require.NoError(t, repo.Update(ctx, u))

	found, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.NotNil(t, found.EducationExperiences)
	assert.Empty(t, found.EducationExperiences)

	var rawValue string
	err = db.Raw("SELECT education_experiences::text FROM users WHERE id = ?", u.ID).Scan(&rawValue).Error
	require.NoError(t, err)
	assert.Equal(t, "[]", rawValue)
}

func TestIntegrationUserRepo_FindNormalizesNullableLegacyEducationExperiences(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "legacyuser",
		Email:    "legacy@example.com",
	}
	_ = u.SetPassword("password123")
	require.NoError(t, db.Omit("EducationExperiences").Create(u).Error)

	found, err := repo.FindByUsername(ctx, "legacyuser")
	require.NoError(t, err)
	assert.NotNil(t, found.EducationExperiences)
	assert.Empty(t, found.EducationExperiences)
}
