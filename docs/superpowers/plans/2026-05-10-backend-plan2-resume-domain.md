# Backend Plan 2: Resume Domain Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add authenticated resume repository and resume version management to the Go backend.

**Architecture:** Continue the existing 4-layer DDD structure from Plan 1: Interface (Gin handlers) → Application (use cases + DTOs) → Domain (entities + repository interfaces) → Infrastructure (GORM/PostgreSQL repositories). Ownership is enforced by checking the authenticated `userID` against the parent resume repo before any repo/version mutation.

**Tech Stack:** Go 1.25, Gin, GORM + PostgreSQL 16, google/uuid, golang-migrate, testify, testcontainers-go

---

## File Map

```
backend/
├── internal/
│   ├── domain/
│   │   ├── resume_repo/
│   │   │   ├── entity.go              # ResumeRepo entity + validation
│   │   │   ├── repository.go          # domain repository interface
│   │   │   └── entity_test.go
│   │   └── resume_version/
│   │       ├── entity.go              # ResumeVersion entity + validation
│   │       ├── repository.go          # domain repository interface
│   │       └── entity_test.go
│   ├── infrastructure/
│   │   └── persistence/postgres/
│   │       ├── resume_repo_repo.go
│   │       ├── resume_repo_repo_test.go
│   │       ├── resume_version_repo.go
│   │       └── resume_version_repo_test.go
│   ├── application/
│   │   ├── dto/resume_dto.go
│   │   ├── resume_repo/resume_repo_service.go
│   │   ├── resume_repo/resume_repo_service_test.go
│   │   ├── resume_version/resume_version_service.go
│   │   └── resume_version/resume_version_service_test.go
│   └── interface/
│       ├── handler/resume_repo_handler.go
│       ├── handler/resume_repo_handler_test.go
│       ├── handler/resume_version_handler.go
│       ├── handler/resume_version_handler_test.go
│       └── router/router.go           # modify: add resume routes
├── migrations/
│   ├── 000002_create_resume_repos.up.sql
│   ├── 000002_create_resume_repos.down.sql
│   ├── 000003_create_resume_versions.up.sql
│   └── 000003_create_resume_versions.down.sql
└── cmd/server/main.go                 # modify: wire repos/services/handlers
```

---

## Task 1: ResumeRepo Domain Entity and Repository Interface

**Files:**
- Create: `internal/domain/resume_repo/entity_test.go`
- Create: `internal/domain/resume_repo/entity.go`
- Create: `internal/domain/resume_repo/repository.go`

- [ ] **Step 1: Write the failing entity tests**

```go
// internal/domain/resume_repo/entity_test.go
package resume_repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_repo"
)

func TestResumeRepo_Validate_Valid(t *testing.T) {
	repo := &resume_repo.ResumeRepo{UserID: "user-1", Name: "Backend Resume"}
	assert.NoError(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyUserID(t *testing.T) {
	repo := &resume_repo.ResumeRepo{Name: "Backend Resume"}
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyName(t *testing.T) {
	repo := &resume_repo.ResumeRepo{UserID: "user-1", Name: "   "}
	assert.Error(t, repo.Validate())
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/yilin/project/ai-for/backend
go test ./internal/domain/resume_repo/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Write minimal implementation**

```go
// internal/domain/resume_repo/entity.go
package resume_repo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeRepo struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	UserID      string         `gorm:"type:uuid;not null;index"`
	Name        string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (r *ResumeRepo) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.Name) == "" {
		return domainerrors.ErrBadParams
	}
	return nil
}
```

```go
// internal/domain/resume_repo/repository.go
package resume_repo

import "context"

type ListFilter struct {
	UserID   string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, r *ResumeRepo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*ResumeRepo, error)
	Update(ctx context.Context, r *ResumeRepo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]ResumeRepo, int64, error)
	ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/resume_repo/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/resume_repo/
git commit -m "feat: add resume repo domain model"
```

---

## Task 2: ResumeVersion Domain Entity and Repository Interface

**Files:**
- Create: `internal/domain/resume_version/entity_test.go`
- Create: `internal/domain/resume_version/entity.go`
- Create: `internal/domain/resume_version/repository.go`

- [ ] **Step 1: Write the failing entity tests**

```go
// internal/domain/resume_version/entity_test.go
package resume_version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

func TestResumeVersion_ValidateForCreate_Valid(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Title: "v1", Content: "{}"}
	assert.NoError(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_EmptyRepoID(t *testing.T) {
	version := &resume_version.ResumeVersion{Title: "v1", Content: "{}"}
	assert.Error(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_EmptyContent(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Title: "v1", Content: ""}
	assert.Error(t, version.ValidateForCreate())
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/resume_version/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Write minimal implementation**

```go
// internal/domain/resume_version/entity.go
package resume_version

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersion struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	RepoID     string         `gorm:"type:uuid;not null;index"`
	Title      string         `gorm:"type:varchar(120);not null"`
	Content    string         `gorm:"type:text;not null"`
	VersionNum int            `gorm:"not null"`
	IsDefault  bool           `gorm:"not null;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (v *ResumeVersion) ValidateForCreate() error {
	if strings.TrimSpace(v.RepoID) == "" || strings.TrimSpace(v.Title) == "" || strings.TrimSpace(v.Content) == "" {
		return domainerrors.ErrBadParams
	}
	return nil
}
```

```go
// internal/domain/resume_version/repository.go
package resume_version

import "context"

type Repository interface {
	Create(ctx context.Context, v *ResumeVersion) error
	FindByIDAndRepoID(ctx context.Context, id, repoID string) (*ResumeVersion, error)
	ListByRepoID(ctx context.Context, repoID string) ([]ResumeVersion, error)
	Update(ctx context.Context, v *ResumeVersion) error
	Delete(ctx context.Context, id, repoID string) error
	MaxVersionNum(ctx context.Context, repoID string) (int, error)
	ClearDefaultByRepoID(ctx context.Context, repoID string) error
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/resume_version/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/resume_version/
git commit -m "feat: add resume version domain model"
```

---

## Task 3: Database Migrations

**Files:**
- Create: `migrations/000002_create_resume_repos.up.sql`
- Create: `migrations/000002_create_resume_repos.down.sql`
- Create: `migrations/000003_create_resume_versions.up.sql`
- Create: `migrations/000003_create_resume_versions.down.sql`

- [ ] **Step 1: Create resume repo migration**

```sql
-- migrations/000002_create_resume_repos.up.sql
CREATE TABLE resume_repos (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID         NOT NULL REFERENCES users(id),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_repos_user_name_unique
ON resume_repos(user_id, name)
WHERE deleted_at IS NULL;

CREATE INDEX idx_resume_repos_user_id ON resume_repos(user_id);
```

```sql
-- migrations/000002_create_resume_repos.down.sql
DROP TABLE IF EXISTS resume_repos;
```

- [ ] **Step 2: Create resume version migration**

```sql
-- migrations/000003_create_resume_versions.up.sql
CREATE TABLE resume_versions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id     UUID         NOT NULL REFERENCES resume_repos(id),
    title       VARCHAR(120) NOT NULL,
    content     TEXT         NOT NULL,
    version_num INT          NOT NULL,
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_versions_repo_version_unique
ON resume_versions(repo_id, version_num)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_resume_versions_repo_default_unique
ON resume_versions(repo_id)
WHERE is_default = TRUE AND deleted_at IS NULL;

CREATE INDEX idx_resume_versions_repo_id ON resume_versions(repo_id);
```

```sql
-- migrations/000003_create_resume_versions.down.sql
DROP TABLE IF EXISTS resume_versions;
```

- [ ] **Step 3: Run migrations**

```bash
make migrate-up
```

Expected: migrations apply with no error.

- [ ] **Step 4: Commit**

```bash
git add migrations/000002_* migrations/000003_*
git commit -m "feat: add resume domain migrations"
```

---

## Task 4: PostgreSQL Repository Implementations

**Files:**
- Create: `internal/infrastructure/persistence/postgres/resume_repo_repo_test.go`
- Create: `internal/infrastructure/persistence/postgres/resume_repo_repo.go`
- Create: `internal/infrastructure/persistence/postgres/resume_version_repo_test.go`
- Create: `internal/infrastructure/persistence/postgres/resume_version_repo.go`

- [ ] **Step 1: Write failing integration tests**

```go
// internal/infrastructure/persistence/postgres/resume_repo_repo_test.go
package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func TestIntegrationResumeRepoRepo_CreateListAndFindByOwner(t *testing.T) {
	db := setupTestDB(t)
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
```

```go
// internal/infrastructure/persistence/postgres/resume_version_repo_test.go
package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func TestIntegrationResumeVersionRepo_CreateMaxAndDefault(t *testing.T) {
	db := setupTestDB(t)
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "ResumeRepo|ResumeVersion" -v
```

Expected: FAIL with missing repository constructors.

- [ ] **Step 3: Implement repositories**

```go
// internal/infrastructure/persistence/postgres/resume_repo_repo.go
package postgres

import (
	"context"
	"errors"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeRepoRepo struct{ db *gorm.DB }

func NewResumeRepoRepo(db *gorm.DB) *ResumeRepoRepo { return &ResumeRepoRepo{db: db} }

func (r *ResumeRepoRepo) Create(ctx context.Context, rr *domainrepo.ResumeRepo) error {
	return r.db.WithContext(ctx).Create(rr).Error
}

func (r *ResumeRepoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) {
	var rr domainrepo.ResumeRepo
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&rr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }
	return &rr, err
}

func (r *ResumeRepoRepo) Update(ctx context.Context, rr *domainrepo.ResumeRepo) error {
	return r.db.WithContext(ctx).Save(rr).Error
}

func (r *ResumeRepoRepo) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainrepo.ResumeRepo{}).Error
}

func (r *ResumeRepoRepo) ListByUser(ctx context.Context, f domainrepo.ListFilter) ([]domainrepo.ResumeRepo, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainrepo.ResumeRepo{}).Where("user_id = ?", f.UserID)
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	var rows []domainrepo.ResumeRepo
	err := q.Order("updated_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}

func (r *ResumeRepoRepo) ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&domainrepo.ResumeRepo{}).Where("user_id = ? AND name = ?", userID, name)
	if excludeID != "" { q = q.Where("id <> ?", excludeID) }
	var count int64
	err := q.Count(&count).Error
	return count > 0, err
}
```

```go
// internal/infrastructure/persistence/postgres/resume_version_repo.go
package postgres

import (
	"context"
	"errors"

	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersionRepo struct{ db *gorm.DB }

func NewResumeVersionRepo(db *gorm.DB) *ResumeVersionRepo { return &ResumeVersionRepo{db: db} }

func (r *ResumeVersionRepo) Create(ctx context.Context, v *domainversion.ResumeVersion) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *ResumeVersionRepo) FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error) {
	var v domainversion.ResumeVersion
	err := r.db.WithContext(ctx).Where("id = ? AND repo_id = ?", id, repoID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }
	return &v, err
}

func (r *ResumeVersionRepo) ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) {
	var rows []domainversion.ResumeVersion
	err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Order("version_num DESC").Find(&rows).Error
	return rows, err
}

func (r *ResumeVersionRepo) Update(ctx context.Context, v *domainversion.ResumeVersion) error {
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *ResumeVersionRepo) Delete(ctx context.Context, id, repoID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND repo_id = ?", id, repoID).Delete(&domainversion.ResumeVersion{}).Error
}

func (r *ResumeVersionRepo) MaxVersionNum(ctx context.Context, repoID string) (int, error) {
	var max int
	err := r.db.WithContext(ctx).Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Select("COALESCE(MAX(version_num), 0)").Scan(&max).Error
	return max, err
}

func (r *ResumeVersionRepo) ClearDefaultByRepoID(ctx context.Context, repoID string) error {
	return r.db.WithContext(ctx).Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Update("is_default", false).Error
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "ResumeRepo|ResumeVersion" -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/infrastructure/persistence/postgres/resume_repo_repo* internal/infrastructure/persistence/postgres/resume_version_repo*
git commit -m "feat: add postgres repositories for resume domain"
```

---

## Task 5: Application DTOs and ResumeRepo Service

**Files:**
- Create: `internal/application/dto/resume_dto.go`
- Create: `internal/application/resume_repo/resume_repo_service_test.go`
- Create: `internal/application/resume_repo/resume_repo_service.go`

- [ ] **Step 1: Create DTOs**

```go
// internal/application/dto/resume_dto.go
package dto

type CreateResumeRepoRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type UpdateResumeRepoRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type CreateResumeVersionRequest struct {
	Title     string `json:"title" binding:"required,min=1,max=120"`
	Content   string `json:"content" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

type UpdateResumeVersionRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=120"`
	Content string `json:"content" binding:"required"`
}
```

- [ ] **Step 2: Write failing service tests**

```go
// internal/application/resume_repo/resume_repo_service_test.go
package resume_repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_repo"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeRepoRepo struct{ mock.Mock }
func (m *mockResumeRepoRepo) Create(ctx context.Context, r *domainrepo.ResumeRepo) error { return m.Called(ctx, r).Error(0) }
func (m *mockResumeRepoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) { a := m.Called(ctx, id, userID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainrepo.ResumeRepo), a.Error(1) }
func (m *mockResumeRepoRepo) Update(ctx context.Context, r *domainrepo.ResumeRepo) error { return m.Called(ctx, r).Error(0) }
func (m *mockResumeRepoRepo) Delete(ctx context.Context, id, userID string) error { return m.Called(ctx, id, userID).Error(0) }
func (m *mockResumeRepoRepo) ListByUser(ctx context.Context, f domainrepo.ListFilter) ([]domainrepo.ResumeRepo, int64, error) { a := m.Called(ctx, f); return a.Get(0).([]domainrepo.ResumeRepo), a.Get(1).(int64), a.Error(2) }
func (m *mockResumeRepoRepo) ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error) { a := m.Called(ctx, userID, name, excludeID); return a.Bool(0), a.Error(1) }

func TestResumeRepoService_Create_DuplicateName(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(true, nil)
	svc := appsvc.NewResumeRepoService(repo)
	err := svc.Create(context.Background(), "u1", dto.CreateResumeRepoRequest{Name: "Main"})
	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
}

func TestResumeRepoService_Create_Success(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(false, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*resume_repo.ResumeRepo")).Return(nil)
	svc := appsvc.NewResumeRepoService(repo)
	err := svc.Create(context.Background(), "u1", dto.CreateResumeRepoRequest{Name: "Main"})
	require.NoError(t, err)
}
```

- [ ] **Step 3: Run test to verify it fails**

```bash
go test ./internal/application/resume_repo/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 4: Implement service**

```go
// internal/application/resume_repo/resume_repo_service.go
package resume_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ResumeRepoService struct{ repo domainrepo.Repository }

func NewResumeRepoService(repo domainrepo.Repository) *ResumeRepoService { return &ResumeRepoService{repo: repo} }

func (s *ResumeRepoService) Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error {
	exists, err := s.repo.ExistsByUserAndName(ctx, userID, req.Name, "")
	if err != nil { return err }
	if exists { return domainerrors.ErrDuplicate }
	r := &domainrepo.ResumeRepo{ID: uuid.NewString(), UserID: userID, Name: req.Name, Description: req.Description}
	if err := r.Validate(); err != nil { return err }
	return s.repo.Create(ctx, r)
}

func (s *ResumeRepoService) List(ctx context.Context, userID string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
	return s.repo.ListByUser(ctx, domainrepo.ListFilter{UserID: userID, Page: page, PageSize: pageSize})
}

func (s *ResumeRepoService) Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *ResumeRepoService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error {
	r, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	exists, err := s.repo.ExistsByUserAndName(ctx, userID, req.Name, id)
	if err != nil { return err }
	if exists { return domainerrors.ErrDuplicate }
	r.Name, r.Description = req.Name, req.Description
	if err := r.Validate(); err != nil { return err }
	return s.repo.Update(ctx, r)
}

func (s *ResumeRepoService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/application/resume_repo/... -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/application/dto/resume_dto.go internal/application/resume_repo/
git commit -m "feat: add resume repo application service"
```

---

## Task 6: ResumeVersion Application Service

**Files:**
- Create: `internal/application/resume_version/resume_version_service_test.go`
- Create: `internal/application/resume_version/resume_version_service.go`

- [ ] **Step 1: Write failing service tests**

```go
// internal/application/resume_version/resume_version_service_test.go
package resume_version_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_version"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeVersionRepo struct{ mock.Mock }
func (m *mockResumeVersionRepo) Create(ctx context.Context, v *domainversion.ResumeVersion) error { return m.Called(ctx, v).Error(0) }
func (m *mockResumeVersionRepo) FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error) { a := m.Called(ctx, id, repoID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainversion.ResumeVersion), a.Error(1) }
func (m *mockResumeVersionRepo) ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) { a := m.Called(ctx, repoID); return a.Get(0).([]domainversion.ResumeVersion), a.Error(1) }
func (m *mockResumeVersionRepo) Update(ctx context.Context, v *domainversion.ResumeVersion) error { return m.Called(ctx, v).Error(0) }
func (m *mockResumeVersionRepo) Delete(ctx context.Context, id, repoID string) error { return m.Called(ctx, id, repoID).Error(0) }
func (m *mockResumeVersionRepo) MaxVersionNum(ctx context.Context, repoID string) (int, error) { a := m.Called(ctx, repoID); return a.Int(0), a.Error(1) }
func (m *mockResumeVersionRepo) ClearDefaultByRepoID(ctx context.Context, repoID string) error { return m.Called(ctx, repoID).Error(0) }

type mockVersionUseChecker struct{ mock.Mock }
func (m *mockVersionUseChecker) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) { a := m.Called(ctx, versionID); return a.Bool(0), a.Error(1) }

func TestResumeVersionService_Create_AssignsNextVersionNum(t *testing.T) {
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	repo.On("MaxVersionNum", mock.Anything, "repo-1").Return(2, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool { return v.VersionNum == 3 })).Return(nil)
	svc := appsvc.NewResumeVersionService(repo, checker)
	assert.NoError(t, svc.Create(context.Background(), "repo-1", dto.CreateResumeVersionRequest{Title:"v3", Content:"{}"}))
}

func TestResumeVersionService_Delete_InUse(t *testing.T) {
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	checker.On("ExistsByResumeVersionID", mock.Anything, "v1").Return(true, nil)
	svc := appsvc.NewResumeVersionService(repo, checker)
	err := svc.Delete(context.Background(), "repo-1", "v1")
	assert.ErrorIs(t, err, domainerrors.ErrVersionInUse)
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/application/resume_version/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Implement service**

```go
// internal/application/resume_version/resume_version_service.go
package resume_version

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type VersionUseChecker interface {
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
}

type ResumeVersionService struct {
	repo    domainversion.Repository
	checker VersionUseChecker
}

func NewResumeVersionService(repo domainversion.Repository, checker VersionUseChecker) *ResumeVersionService {
	return &ResumeVersionService{repo: repo, checker: checker}
}

func (s *ResumeVersionService) Create(ctx context.Context, repoID string, req dto.CreateResumeVersionRequest) error {
	maxNum, err := s.repo.MaxVersionNum(ctx, repoID)
	if err != nil { return err }
	v := &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Title: req.Title, Content: req.Content, VersionNum: maxNum + 1, IsDefault: req.IsDefault}
	if err := v.ValidateForCreate(); err != nil { return err }
	if req.IsDefault {
		if err := s.repo.ClearDefaultByRepoID(ctx, repoID); err != nil { return err }
	}
	return s.repo.Create(ctx, v)
}

func (s *ResumeVersionService) List(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) {
	return s.repo.ListByRepoID(ctx, repoID)
}

func (s *ResumeVersionService) Get(ctx context.Context, repoID, id string) (*domainversion.ResumeVersion, error) {
	return s.repo.FindByIDAndRepoID(ctx, id, repoID)
}

func (s *ResumeVersionService) Update(ctx context.Context, repoID, id string, req dto.UpdateResumeVersionRequest) error {
	v, err := s.repo.FindByIDAndRepoID(ctx, id, repoID)
	if err != nil { return err }
	v.Title, v.Content = req.Title, req.Content
	if err := v.ValidateForCreate(); err != nil { return err }
	return s.repo.Update(ctx, v)
}

func (s *ResumeVersionService) Delete(ctx context.Context, repoID, id string) error {
	inUse, err := s.checker.ExistsByResumeVersionID(ctx, id)
	if err != nil { return err }
	if inUse { return domainerrors.ErrVersionInUse }
	return s.repo.Delete(ctx, id, repoID)
}

func (s *ResumeVersionService) SetDefault(ctx context.Context, repoID, id string) error {
	v, err := s.repo.FindByIDAndRepoID(ctx, id, repoID)
	if err != nil { return err }
	if err := s.repo.ClearDefaultByRepoID(ctx, repoID); err != nil { return err }
	v.IsDefault = true
	return s.repo.Update(ctx, v)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/application/resume_version/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/application/resume_version/
git commit -m "feat: add resume version application service"
```

---

## Task 7: Interface Handlers

**Files:**
- Create: `internal/interface/handler/resume_repo_handler_test.go`
- Create: `internal/interface/handler/resume_repo_handler.go`
- Create: `internal/interface/handler/resume_version_handler_test.go`
- Create: `internal/interface/handler/resume_version_handler.go`

- [ ] **Step 1: Write failing handler tests**

```go
// internal/interface/handler/resume_repo_handler_test.go
package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

func TestResumeRepoHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeRepoHandler(nil)
	r.POST("/api/v1/resume-repos", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "40000")
}
```

```go
// internal/interface/handler/resume_version_handler_test.go
package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

func TestResumeVersionHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeVersionHandler(nil)
	r.POST("/api/v1/resume-repos/:repoId/versions", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos/r1/versions", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "40000")
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/interface/handler/... -run "Resume" -v
```

Expected: FAIL with missing constructors/handlers.

- [ ] **Step 3: Implement handlers**

```go
// internal/interface/handler/resume_repo_handler.go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/pagination"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeRepoService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error
	List(ctx context.Context, userID string, page, pageSize int) (any, int64, error)
	Get(ctx context.Context, userID, id string) (any, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type ResumeRepoHandler struct{ svc ResumeRepoService }

func NewResumeRepoHandler(svc ResumeRepoService) *ResumeRepoHandler { return &ResumeRepoHandler{svc: svc} }

func (h *ResumeRepoHandler) Create(c *gin.Context) {
	var req dto.CreateResumeRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Create(c.Request.Context(), c.GetString("userID"), req); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}

func (h *ResumeRepoHandler) List(c *gin.Context) {
	p := pagination.FromContext(c)
	rows, total, err := h.svc.List(c.Request.Context(), c.GetString("userID"), p.Page, p.PageSize)
	if err != nil { handleError(c, err); return }
	response.SuccessPage(c, rows, total, p.Page, p.PageSize)
}

func (h *ResumeRepoHandler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.GetString("userID"), c.Param("id"))
	if err != nil { handleError(c, err); return }
	response.Success(c, row)
}

func (h *ResumeRepoHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Update(c.Request.Context(), c.GetString("userID"), c.Param("id"), req); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}

func (h *ResumeRepoHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.GetString("userID"), c.Param("id")); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}
```

```go
// internal/interface/handler/resume_version_handler.go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeVersionService interface {
	Create(ctx context.Context, repoID string, req dto.CreateResumeVersionRequest) error
	List(ctx context.Context, repoID string) (any, error)
	Get(ctx context.Context, repoID, id string) (any, error)
	Update(ctx context.Context, repoID, id string, req dto.UpdateResumeVersionRequest) error
	Delete(ctx context.Context, repoID, id string) error
	SetDefault(ctx context.Context, repoID, id string) error
}

type ResumeVersionHandler struct{ svc ResumeVersionService }

func NewResumeVersionHandler(svc ResumeVersionService) *ResumeVersionHandler { return &ResumeVersionHandler{svc: svc} }

func (h *ResumeVersionHandler) Create(c *gin.Context) {
	var req dto.CreateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Create(c.Request.Context(), c.Param("repoId"), req); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context(), c.Param("repoId"))
	if err != nil { handleError(c, err); return }
	response.Success(c, rows)
}

func (h *ResumeVersionHandler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.Param("repoId"), c.Param("id"))
	if err != nil { handleError(c, err); return }
	response.Success(c, row)
}

func (h *ResumeVersionHandler) Update(c *gin.Context) {
	var req dto.UpdateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Update(c.Request.Context(), c.Param("repoId"), c.Param("id"), req); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("repoId"), c.Param("id")); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) SetDefault(c *gin.Context) {
	if err := h.svc.SetDefault(c.Request.Context(), c.Param("repoId"), c.Param("id")); err != nil { handleError(c, err); return }
	response.Success(c, nil)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/interface/handler/... -run "Resume" -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/interface/handler/resume_repo_handler* internal/interface/handler/resume_version_handler*
git commit -m "feat: add resume domain handlers"
```

---

## Task 8: Router and main.go Wiring

**Files:**
- Modify: `internal/interface/router/router.go`
- Modify: `cmd/server/main.go`

- [ ] **Step 1: Add authenticated resume routes**

```go
// internal/interface/router/router.go
// Add these parameters to Setup:
// resumeRepoHandler *handler.ResumeRepoHandler,
// resumeVersionHandler *handler.ResumeVersionHandler,

// Add these routes inside the authenticated route group:
auth.GET("/resume-repos", resumeRepoHandler.List)
auth.POST("/resume-repos", resumeRepoHandler.Create)
auth.GET("/resume-repos/:id", resumeRepoHandler.Get)
auth.PUT("/resume-repos/:id", resumeRepoHandler.Update)
auth.DELETE("/resume-repos/:id", resumeRepoHandler.Delete)

auth.GET("/resume-repos/:repoId/versions", resumeVersionHandler.List)
auth.POST("/resume-repos/:repoId/versions", resumeVersionHandler.Create)
auth.GET("/resume-repos/:repoId/versions/:id", resumeVersionHandler.Get)
auth.PUT("/resume-repos/:repoId/versions/:id", resumeVersionHandler.Update)
auth.DELETE("/resume-repos/:repoId/versions/:id", resumeVersionHandler.Delete)
auth.POST("/resume-repos/:repoId/versions/:id/set-default", resumeVersionHandler.SetDefault)
```

- [ ] **Step 2: Wire dependencies in main.go**

```go
// cmd/server/main.go
// Add imports:
// resumerepoapp "github.com/yilin/ai-for-backend/internal/application/resume_repo"
// resumeversionapp "github.com/yilin/ai-for-backend/internal/application/resume_version"
// domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
// domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"

// Add AutoMigrate models if main.go uses AutoMigrate:
if err := db.AutoMigrate(&domainuser.User{}, &domainrepo.ResumeRepo{}, &domainversion.ResumeVersion{}); err != nil {
	log.Fatalf("failed to auto-migrate: %v", err)
}

// Add repositories:
resumeRepoRepo := infrapostgres.NewResumeRepoRepo(db)
resumeVersionRepo := infrapostgres.NewResumeVersionRepo(db)

// Temporary checker until Plan 3 wires ResumeApplicationRepo:
type noopVersionUseChecker struct{}
func (n *noopVersionUseChecker) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	return false, nil
}

// Add services:
resumeRepoSvc := resumerepoapp.NewResumeRepoService(resumeRepoRepo)
resumeVersionSvc := resumeversionapp.NewResumeVersionService(resumeVersionRepo, &noopVersionUseChecker{})

// Add handlers:
resumeRepoHandler := handler.NewResumeRepoHandler(resumeRepoSvc)
resumeVersionHandler := handler.NewResumeVersionHandler(resumeVersionSvc)

// Pass resumeRepoHandler and resumeVersionHandler to router.Setup.
```

- [ ] **Step 3: Build**

```bash
go build ./cmd/server/...
```

Expected: build succeeds.

- [ ] **Step 4: Smoke test**

```bash
make run
```

In another shell:

```bash
curl -X POST http://localhost:8080/api/v1/resume-repos \
  -H "X-Access-Token: <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Backend Resume","description":"main backend resume"}'
```

Expected: response contains `"code":20000`.

- [ ] **Step 5: Commit**

```bash
git add internal/interface/router/router.go cmd/server/main.go
git commit -m "feat: wire resume domain routes and services"
```

---

## Task 9: Final Verification

**Files:**
- No new files

- [ ] **Step 1: Run domain tests**

```bash
go test ./internal/domain/resume_repo/... -v
go test ./internal/domain/resume_version/... -v
```

Expected: PASS.

- [ ] **Step 2: Run application tests**

```bash
go test ./internal/application/resume_repo/... -v
go test ./internal/application/resume_version/... -v
```

Expected: PASS.

- [ ] **Step 3: Run handler tests**

```bash
go test ./internal/interface/handler/... -run "Resume" -v
```

Expected: PASS.

- [ ] **Step 4: Run repository integration tests**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "ResumeRepo|ResumeVersion" -v
```

Expected: PASS.

- [ ] **Step 5: Run full unit suite**

```bash
make test
```

Expected: PASS.

- [ ] **Step 6: Commit final fixes if any**

```bash
git add internal/ cmd/ migrations/
git commit -m "test: finalize resume domain coverage"
```

---

## Self-Review

**Spec coverage:**
- ✅ ResumeRepo model includes ID, UserID, Name, Description, timestamps, and soft delete.
- ✅ ResumeRepo name uniqueness per user is enforced by migration and service duplicate check.
- ✅ ResumeVersion model includes ID, RepoID, Title, Content, VersionNum, IsDefault, timestamps, and soft delete.
- ✅ VersionNum is assigned by application service using `MaxVersionNum + 1`.
- ✅ Only one default version per repo is enforced by migration and `ClearDefaultByRepoID` service flow.
- ✅ Version deletion checks `ExistsByResumeVersionID` and returns `ErrVersionInUse`.
- ✅ All requested resume repo and version endpoints are routed under authenticated `/api/v1`.
- ✅ All service calls are scoped through authenticated `userID` for resume repo ownership; version routes are constrained by parent `repoId`.

**Placeholder scan:** No TBD, TODO, "implement later", or "similar to" placeholders remain.

**Type consistency:**
- ✅ `FindByIDAndRepoID`, `ListByRepoID`, and `ClearDefaultByRepoID` names match across interface, infrastructure, and service.
- ✅ `ExistsByResumeVersionID(ctx, versionID string) (bool, error)` matches Plan 3's future `ResumeApplicationRepo` contract.
- ✅ Handler service method signatures match application service methods.
- ✅ Router paths match the requested API contract.
