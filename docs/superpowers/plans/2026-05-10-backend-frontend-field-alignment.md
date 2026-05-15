# Backend Frontend Field Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the backend database schema, domain entities, DTOs, services, handlers, and routes match the existing frontend resume repo/version models and API calls.

**Architecture:** Keep the existing DDD layering: domain entities define persisted fields, application DTOs translate to frontend JSON, services enforce ownership/business rules, repositories handle GORM persistence, and handlers expose Gin routes. The backend will adopt frontend field names at the JSON boundary (`jobType`, `targetPosition`, `versionCount`, `version`, `remark`, `isActive`) and matching snake_case columns in PostgreSQL.

**Tech Stack:** Go, Gin, GORM, PostgreSQL migrations, testify, existing `pkg/response` response wrapper.

---

## File Structure

Modify these files:

- `backend/migrations/000002_create_resume_repos.up.sql` — add repo columns `job_type`, `target_position`, `version_count`.
- `backend/migrations/000003_create_resume_versions.up.sql` — replace version columns with frontend-aligned `version`, `remark`, `is_active`.
- `backend/internal/domain/resume_repo/entity.go` — add `JobType`, `TargetPosition`, `VersionCount`, validation.
- `backend/internal/domain/resume_repo/repository.go` — add version count update method and `jobType` filter support.
- `backend/internal/domain/resume_version/entity.go` — replace `VersionNum/IsDefault` with `Version/Remark/IsActive`.
- `backend/internal/domain/resume_version/repository.go` — replace numeric version/default methods with frontend-aligned methods.
- `backend/internal/application/dto/resume_dto.go` — make request/response JSON match frontend models.
- `backend/internal/application/resume_repo/resume_repo_service.go` — populate/update frontend-aligned repo fields and expose `IncrementVersionCount` for version service use.
- `backend/internal/application/resume_version/resume_version_service.go` — accept `repoId` from request, generate `v1.N.0`, maintain repo `versionCount`, activate via `isActive`.
- `backend/internal/infrastructure/persistence/postgres/resume_repo_repo.go` — persist repo fields and version count.
- `backend/internal/infrastructure/persistence/postgres/resume_version_repo.go` — persist version fields and query by version ID.
- `backend/internal/interface/handler/resume_version_handler.go` — expose `/resume-versions` query/body based API.
- `backend/internal/interface/router/router.go` — register frontend routes.
- Existing tests under `backend/internal/**` — update to encode the frontend contract.

Do not modify frontend files for this plan; the backend must adapt to the existing frontend.

---

### Task 1: Align ResumeRepo domain, DTO, migration, and repository

**Files:**
- Modify: `backend/migrations/000002_create_resume_repos.up.sql`
- Modify: `backend/internal/domain/resume_repo/entity.go`
- Modify: `backend/internal/domain/resume_repo/repository.go`
- Modify: `backend/internal/application/dto/resume_dto.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/resume_repo_repo.go`
- Test: `backend/internal/domain/resume_repo/entity_test.go`
- Test: `backend/internal/infrastructure/persistence/postgres/resume_repo_repo_test.go`

- [ ] **Step 1: Update/verify domain test for frontend repo fields**

In `backend/internal/domain/resume_repo/entity_test.go`, add these tests or update existing validation tests to include the frontend-required fields:

```go
func TestResumeRepo_Validate_RequiresFrontendFields(t *testing.T) {
	tests := []struct {
		name string
		repo resume_repo.ResumeRepo
	}{
		{
			name: "missing job type",
			repo: resume_repo.ResumeRepo{UserID: "u1", Name: "Main", TargetPosition: "Frontend"},
		},
		{
			name: "missing target position",
			repo: resume_repo.ResumeRepo{UserID: "u1", Name: "Main", JobType: resume_repo.JobTypeCampus},
		},
		{
			name: "invalid job type",
			repo: resume_repo.ResumeRepo{UserID: "u1", Name: "Main", JobType: resume_repo.JobType("unknown"), TargetPosition: "Frontend"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}

func TestResumeRepo_Validate_AcceptsFrontendFields(t *testing.T) {
	repo := resume_repo.ResumeRepo{
		UserID:         "u1",
		Name:           "Main",
		JobType:        resume_repo.JobTypeCampus,
		TargetPosition: "Frontend Engineer",
		VersionCount:   0,
	}

	require.NoError(t, repo.Validate())
}
```

Make sure imports include:

```go
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resume_repo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)
```

- [ ] **Step 2: Run the domain test to verify it fails**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_repo -run 'TestResumeRepo_Validate_(RequiresFrontendFields|AcceptsFrontendFields)' -count=1
```

Expected: FAIL because `JobType`, `JobTypeCampus`, and `TargetPosition` do not exist yet.

- [ ] **Step 3: Implement ResumeRepo domain fields and validation**

Replace `backend/internal/domain/resume_repo/entity.go` with this structure while preserving package/import style:

```go
package resume_repo

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeCampus     JobType = "campus"
	JobTypeSocial     JobType = "social"
	JobTypeInternship JobType = "internship"
)

type ResumeRepo struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;index"`
	Name           string `gorm:"type:varchar(100);not null"`
	JobType        JobType `gorm:"type:varchar(20);not null"`
	TargetPosition string `gorm:"type:varchar(100);not null"`
	Description    string `gorm:"type:text"`
	VersionCount   int    `gorm:"not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (r *ResumeRepo) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.Name) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(r.Name)) > 100 {
		return domainerrors.ErrBadParams
	}
	if r.JobType != JobTypeCampus && r.JobType != JobTypeSocial && r.JobType != JobTypeInternship {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(r.TargetPosition) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(r.TargetPosition)) > 100 {
		return domainerrors.ErrBadParams
	}
	if r.VersionCount < 0 {
		return domainerrors.ErrBadParams
	}
	return nil
}
```

- [ ] **Step 4: Update repo migration**

Change `backend/migrations/000002_create_resume_repos.up.sql` table definition to include the frontend fields:

```sql
CREATE TABLE resume_repos (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         VARCHAR(36)  NOT NULL REFERENCES users(id),
    name            VARCHAR(100) NOT NULL,
    job_type        VARCHAR(20)  NOT NULL,
    target_position VARCHAR(100) NOT NULL,
    description     TEXT,
    version_count   INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_repos_user_name_unique
ON resume_repos(user_id, name)
WHERE deleted_at IS NULL;

CREATE INDEX idx_resume_repos_user_id ON resume_repos(user_id);
CREATE INDEX idx_resume_repos_user_job_type ON resume_repos(user_id, job_type);
```

- [ ] **Step 5: Update repo DTOs**

In `backend/internal/application/dto/resume_dto.go`, change repo request/response types and converter to:

```go
type CreateResumeRepoRequest struct {
	Name           string             `json:"name" binding:"required,min=1,max=100"`
	JobType        domainrepo.JobType `json:"jobType" binding:"required"`
	TargetPosition string             `json:"targetPosition" binding:"required,min=1,max=100"`
	Description    string             `json:"description"`
}

type UpdateResumeRepoRequest struct {
	Name           string             `json:"name" binding:"required,min=1,max=100"`
	JobType        domainrepo.JobType `json:"jobType" binding:"required"`
	TargetPosition string             `json:"targetPosition" binding:"required,min=1,max=100"`
	Description    string             `json:"description"`
}

type ResumeRepoResponse struct {
	ID             string             `json:"id"`
	UserID         string             `json:"userId"`
	Name           string             `json:"name"`
	JobType        domainrepo.JobType `json:"jobType"`
	TargetPosition string             `json:"targetPosition"`
	Description    string             `json:"description"`
	VersionCount   int                `json:"versionCount"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}
```

Update `ToResumeRepoResponse` to populate:

```go
return &ResumeRepoResponse{
	ID:             repo.ID,
	UserID:         repo.UserID,
	Name:           repo.Name,
	JobType:        repo.JobType,
	TargetPosition: repo.TargetPosition,
	Description:    repo.Description,
	VersionCount:   repo.VersionCount,
	CreatedAt:      repo.CreatedAt,
	UpdatedAt:      repo.UpdatedAt,
}
```

- [ ] **Step 6: Update repo repository interface**

In `backend/internal/domain/resume_repo/repository.go`, add `JobType` to `ListFilter` and add a version-count method:

```go
type ListFilter struct {
	UserID   string
	JobType  JobType
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
	IncrementVersionCount(ctx context.Context, id string, delta int) error
}
```

- [ ] **Step 7: Update GORM repo persistence**

In `backend/internal/infrastructure/persistence/postgres/resume_repo_repo.go`:

Update `Update` map:

```go
Updates(map[string]any{
	"name":            rr.Name,
	"job_type":        rr.JobType,
	"target_position": rr.TargetPosition,
	"description":     rr.Description,
	"version_count":   rr.VersionCount,
})
```

Update `ListByUser` before count:

```go
q := r.db.WithContext(ctx).Model(&domainrepo.ResumeRepo{}).Where("user_id = ?", f.UserID)
if f.JobType != "" {
	q = q.Where("job_type = ?", f.JobType)
}
```

Add method:

```go
func (r *ResumeRepoRepo) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	result := r.db.WithContext(ctx).
		Model(&domainrepo.ResumeRepo{}).
		Where("id = ?", id).
		Update("version_count", gorm.Expr("GREATEST(version_count + ?, 0)", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}
```

This method requires `gorm` and `domainerrors`, which are already imported in the file.

- [ ] **Step 8: Run repo tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_repo ./internal/infrastructure/persistence/postgres -run 'ResumeRepo' -count=1
```

Expected: PASS after tests are updated for required fields. If existing repo persistence tests instantiate `ResumeRepo{Name: "Main"}`, update those fixtures to include:

```go
JobType: resume_repo.JobTypeCampus, TargetPosition: "Frontend Engineer"
```

---

### Task 2: Align ResumeRepo service and handler tests

**Files:**
- Modify: `backend/internal/application/resume_repo/resume_repo_service.go`
- Modify: `backend/internal/application/resume_repo/resume_repo_service_test.go`
- Modify: `backend/internal/interface/handler/resume_repo_handler_test.go`
- Modify: `backend/internal/interface/handler/resume_repo_handler.go`

- [ ] **Step 1: Update service tests for frontend request fields**

In `backend/internal/application/resume_repo/resume_repo_service_test.go`, update mock interface by adding:

```go
func (m *mockResumeRepoRepo) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return m.Called(ctx, id, delta).Error(0)
}
```

Update create success expectation to check all fields:

```go
repo.On("Create", mock.Anything, mock.MatchedBy(func(r *domainrepo.ResumeRepo) bool {
	return r.UserID == "u1" &&
		r.Name == "Main" &&
		r.JobType == domainrepo.JobTypeCampus &&
		r.TargetPosition == "Frontend Engineer" &&
		r.VersionCount == 0 &&
		r.ID != ""
})).Return(nil)
```

Update create call:

```go
err := svc.Create(context.Background(), "u1", dto.CreateResumeRepoRequest{
	Name:           "Main",
	JobType:        domainrepo.JobTypeCampus,
	TargetPosition: "Frontend Engineer",
})
```

Update update success expectation:

```go
repo.On("Update", mock.Anything, mock.MatchedBy(func(r *domainrepo.ResumeRepo) bool {
	return r.ID == "r1" &&
		r.UserID == "u1" &&
		r.Name == "Main" &&
		r.JobType == domainrepo.JobTypeSocial &&
		r.TargetPosition == "Fullstack Engineer" &&
		r.Description == "updated"
})).Return(nil)
```

Update update call:

```go
err := svc.Update(context.Background(), "u1", "r1", dto.UpdateResumeRepoRequest{
	Name:           "Main",
	JobType:        domainrepo.JobTypeSocial,
	TargetPosition: "Fullstack Engineer",
	Description:    "updated",
})
```

- [ ] **Step 2: Run service tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_repo -count=1
```

Expected: FAIL because service does not yet copy `JobType` and `TargetPosition`.

- [ ] **Step 3: Update repo service implementation**

In `backend/internal/application/resume_repo/resume_repo_service.go`, create repos with frontend fields:

```go
r := &domainrepo.ResumeRepo{
	ID:             uuid.NewString(),
	UserID:         userID,
	Name:           req.Name,
	JobType:        req.JobType,
	TargetPosition: req.TargetPosition,
	Description:    req.Description,
	VersionCount:   0,
}
```

Update repos with frontend fields:

```go
r.Name = req.Name
r.JobType = req.JobType
r.TargetPosition = req.TargetPosition
r.Description = req.Description
```

Add a pass-through method used by the version service:

```go
func (s *ResumeRepoService) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return s.repo.IncrementVersionCount(ctx, id, delta)
}
```

- [ ] **Step 4: Update handler tests for frontend JSON response**

In `backend/internal/interface/handler/resume_repo_handler_test.go`, update `TestResumeRepoHandler_Create_Success` expectation:

```go
svc.On("Create", mock.Anything, "u1", dto.CreateResumeRepoRequest{
	Name:           "Main",
	JobType:        domainrepo.JobTypeCampus,
	TargetPosition: "Frontend Engineer",
	Description:    "desc",
}).Return(nil)
```

Update request body:

```go
req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-repos", bytes.NewBufferString(`{"name":"Main","jobType":"campus","targetPosition":"Frontend Engineer","description":"desc"}`))
```

Update DTO fixture in list/get tests:

```go
{ID: "r1", UserID: "u1", Name: "Main", JobType: domainrepo.JobTypeCampus, TargetPosition: "Frontend Engineer", Description: "desc", VersionCount: 2, CreatedAt: now, UpdatedAt: now}
```

Add assertions:

```go
assert.Equal(t, "campus", repo["jobType"])
assert.Equal(t, "Frontend Engineer", repo["targetPosition"])
assert.Equal(t, float64(2), repo["versionCount"])
```

Add hidden-key assertions for Go field names:

```go
assertNoJSONKeys(t, repo, "JobType", "TargetPosition", "VersionCount")
```

- [ ] **Step 5: Run repo service/handler tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_repo ./internal/interface/handler -run 'ResumeRepo' -count=1
```

Expected: PASS.

---

### Task 3: Align ResumeVersion domain, DTO, migration, and repository

**Files:**
- Modify: `backend/migrations/000003_create_resume_versions.up.sql`
- Modify: `backend/internal/domain/resume_version/entity.go`
- Modify: `backend/internal/domain/resume_version/repository.go`
- Modify: `backend/internal/application/dto/resume_dto.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/resume_version_repo.go`
- Test: `backend/internal/domain/resume_version/entity_test.go`
- Test: `backend/internal/infrastructure/persistence/postgres/resume_version_repo_test.go`

- [ ] **Step 1: Update domain tests for frontend version fields**

In `backend/internal/domain/resume_version/entity_test.go`, update/add tests:

```go
func TestResumeVersion_ValidateForCreate_AcceptsFrontendFields(t *testing.T) {
	v := resume_version.ResumeVersion{
		RepoID:   "repo-1",
		Version:  "v1.0.0",
		Title:    "Frontend version",
		Content:  "# Resume",
		Remark:   "first version",
		IsActive: false,
	}

	require.NoError(t, v.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_RequiresVersion(t *testing.T) {
	v := resume_version.ResumeVersion{
		RepoID:  "repo-1",
		Title:   "Frontend version",
		Content: "# Resume",
	}

	err := v.ValidateForCreate()

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}
```

- [ ] **Step 2: Run domain tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_version -count=1
```

Expected: FAIL because `Version` and `Remark` do not exist yet.

- [ ] **Step 3: Implement frontend-aligned ResumeVersion entity**

Replace `backend/internal/domain/resume_version/entity.go` struct/validation with:

```go
package resume_version

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeVersion struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	RepoID    string `gorm:"type:uuid;not null;index"`
	Version   string `gorm:"type:varchar(20);not null"`
	Title     string `gorm:"type:varchar(120);not null"`
	Content   string `gorm:"type:text;not null"`
	Remark    string `gorm:"type:text"`
	IsActive  bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (v *ResumeVersion) ValidateForCreate() error {
	if strings.TrimSpace(v.RepoID) == "" || strings.TrimSpace(v.Version) == "" || strings.TrimSpace(v.Title) == "" || strings.TrimSpace(v.Content) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(v.Version)) > 20 {
		return domainerrors.ErrBadParams
	}
	if len([]rune(v.Title)) > 120 {
		return domainerrors.ErrBadParams
	}
	return nil
}
```

- [ ] **Step 4: Update version migration**

Change `backend/migrations/000003_create_resume_versions.up.sql` to:

```sql
CREATE TABLE resume_versions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id     UUID         NOT NULL REFERENCES resume_repos(id),
    version     VARCHAR(20)  NOT NULL,
    title       VARCHAR(120) NOT NULL,
    content     TEXT         NOT NULL,
    remark      TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_versions_repo_version_unique
ON resume_versions(repo_id, version)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_resume_versions_repo_active_unique
ON resume_versions(repo_id)
WHERE is_active = TRUE AND deleted_at IS NULL;

CREATE INDEX idx_resume_versions_repo_id ON resume_versions(repo_id);
```

- [ ] **Step 5: Update version DTOs**

In `backend/internal/application/dto/resume_dto.go`, replace version DTOs with:

```go
type CreateResumeVersionRequest struct {
	RepoID   string `json:"repoId" binding:"required"`
	Title    string `json:"title" binding:"required,min=1,max=120"`
	Content  string `json:"content" binding:"required"`
	Remark   string `json:"remark"`
}

type UpdateResumeVersionRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=120"`
	Content  string `json:"content" binding:"required"`
	Remark   string `json:"remark"`
	IsActive bool   `json:"isActive"`
}

type ResumeVersionResponse struct {
	ID        string    `json:"id"`
	RepoID    string    `json:"repoId"`
	Version   string    `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Remark    string    `json:"remark"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
```

Update converter:

```go
return &ResumeVersionResponse{
	ID:        version.ID,
	RepoID:    version.RepoID,
	Version:   version.Version,
	Title:     version.Title,
	Content:   version.Content,
	Remark:    version.Remark,
	IsActive:  version.IsActive,
	CreatedAt: version.CreatedAt,
	UpdatedAt: version.UpdatedAt,
}
```

- [ ] **Step 6: Update version repository interface**

In `backend/internal/domain/resume_version/repository.go`, replace methods with:

```go
type Repository interface {
	Create(ctx context.Context, v *ResumeVersion) error
	FindByID(ctx context.Context, id string) (*ResumeVersion, error)
	FindByIDAndRepoID(ctx context.Context, id, repoID string) (*ResumeVersion, error)
	ListByRepoID(ctx context.Context, repoID string) ([]ResumeVersion, error)
	Update(ctx context.Context, v *ResumeVersion) error
	Delete(ctx context.Context, id, repoID string) error
	CountByRepoID(ctx context.Context, repoID string) (int, error)
	SetActiveByID(ctx context.Context, repoID, id string) error
}
```

- [ ] **Step 7: Update GORM version persistence**

In `backend/internal/infrastructure/persistence/postgres/resume_version_repo.go`:

Add `FindByID`:

```go
func (r *ResumeVersionRepo) FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error) {
	var v domainversion.ResumeVersion
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &v, err
}
```

Change list ordering:

```go
err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Order("created_at DESC").Find(&rows).Error
```

Change update map:

```go
Updates(map[string]any{
	"title":     v.Title,
	"content":   v.Content,
	"remark":    v.Remark,
	"is_active": v.IsActive,
})
```

Replace `MaxVersionNum`, `ClearDefaultByRepoID`, and `SetDefaultByID` with:

```go
func (r *ResumeVersionRepo) CountByRepoID(ctx context.Context, repoID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Count(&count).Error
	return int(count), err
}

func (r *ResumeVersionRepo) SetActiveByID(ctx context.Context, repoID, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&domainversion.ResumeVersion{}).Where("id = ? AND repo_id = ?", id, repoID).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return domainerrors.ErrNotFound
		}
		if err := tx.Model(&domainversion.ResumeVersion{}).Where("repo_id = ?", repoID).Update("is_active", false).Error; err != nil {
			if isUniqueViolation(err) {
				return domainerrors.ErrDuplicate
			}
			return err
		}
		result := tx.Model(&domainversion.ResumeVersion{}).Where("id = ? AND repo_id = ?", id, repoID).Update("is_active", true)
		if result.Error != nil {
			if isUniqueViolation(result.Error) {
				return domainerrors.ErrDuplicate
			}
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domainerrors.ErrNotFound
		}
		return nil
	})
}
```

- [ ] **Step 8: Run version domain/persistence tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_version ./internal/infrastructure/persistence/postgres -run 'ResumeVersion' -count=1
```

Expected: PASS after updating existing fixtures from `VersionNum/IsDefault` to `Version/IsActive`.

---

### Task 4: Align ResumeVersion service with frontend API contract

**Files:**
- Modify: `backend/internal/application/resume_version/resume_version_service.go`
- Modify: `backend/internal/application/resume_version/resume_version_service_test.go`

- [ ] **Step 1: Update service interfaces and mocks in tests**

In `backend/internal/application/resume_version/resume_version_service_test.go`, add repo count mock method to `mockResumeRepoReader`:

```go
func (m *mockResumeRepoReader) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return m.Called(ctx, id, delta).Error(0)
}
```

Update version repo mock methods:

```go
func (m *mockResumeVersionRepo) FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionRepo) CountByRepoID(ctx context.Context, repoID string) (int, error) {
	args := m.Called(ctx, repoID)
	return args.Int(0), args.Error(1)
}

func (m *mockResumeVersionRepo) SetActiveByID(ctx context.Context, repoID, id string) error {
	return m.Called(ctx, repoID, id).Error(0)
}
```

Remove mock methods `MaxVersionNum`, `ClearDefaultByRepoID`, and `SetDefaultByID`.

- [ ] **Step 2: Update create test for frontend version generation**

Rename `TestResumeVersionService_Create_AssignsNextVersionNum` to `TestResumeVersionService_Create_AssignsNextFrontendVersion` and use:

```go
repo.On("CountByRepoID", mock.Anything, "repo-1").Return(2, nil)
repo.On("Create", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool {
	return v.RepoID == "repo-1" &&
		v.Title == "v3" &&
		v.Content == "{}" &&
		v.Remark == "third" &&
		v.Version == "v1.2.0" &&
		!v.IsActive &&
		v.ID != ""
})).Return(nil)
repoReader.On("IncrementVersionCount", mock.Anything, "repo-1", 1).Return(nil)

err := svc.Create(context.Background(), "user-1", dto.CreateResumeVersionRequest{
	RepoID:  "repo-1",
	Title:   "v3",
	Content: "{}",
	Remark:  "third",
})
```

- [ ] **Step 3: Update list/get/update/delete/activate tests for new signatures**

Use service method signatures:

```go
Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error
List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error)
Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error)
Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error
Delete(ctx context.Context, userID, id string) error
Activate(ctx context.Context, userID, id string) error
```

For `Get`, `Update`, `Delete`, and `Activate`, first mock `FindByID` to obtain `RepoID`, then `repoReader.Get` to verify ownership.

Example get test:

```go
existing := &domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Version: "v1.0.0", Title: "v1", Content: "{}"}
repo.On("FindByID", mock.Anything, "v1").Return(existing, nil)
allowRepo(repoReader)

got, err := svc.Get(context.Background(), "user-1", "v1")

require.NoError(t, err)
assert.Equal(t, existing, got)
```

Example update expectation:

```go
repo.On("Update", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool {
	return v.ID == "v1" &&
		v.RepoID == "repo-1" &&
		v.Version == "v1.0.0" &&
		v.Title == "new" &&
		v.Content == "{updated}" &&
		v.Remark == "updated remark" &&
		v.IsActive
})).Return(nil)
```

Example delete success should also expect decrement:

```go
repoReader.On("IncrementVersionCount", mock.Anything, "repo-1", -1).Return(nil)
```

Example activate:

```go
repo.On("FindByID", mock.Anything, "v1").Return(&domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Version: "v1.0.0", Title: "v1", Content: "{}"}, nil)
allowRepo(repoReader)
repo.On("SetActiveByID", mock.Anything, "repo-1", "v1").Return(nil)

err := svc.Activate(context.Background(), "user-1", "v1")

require.NoError(t, err)
```

- [ ] **Step 4: Run service tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_version -count=1
```

Expected: FAIL because implementation still uses old repoID positional methods, `VersionNum`, and `SetDefault`.

- [ ] **Step 5: Update service implementation**

In `backend/internal/application/resume_version/resume_version_service.go`, change `ResumeRepoReader` to:

```go
type ResumeRepoReader interface {
	Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error)
	IncrementVersionCount(ctx context.Context, id string, delta int) error
}
```

Add helper:

```go
func (s *ResumeVersionService) ensureVersionOwned(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.ensureRepoOwned(ctx, userID, v.RepoID); err != nil {
		return nil, err
	}
	return v, nil
}
```

Replace methods with:

```go
func (s *ResumeVersionService) Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error {
	if err := s.ensureRepoOwned(ctx, userID, req.RepoID); err != nil {
		return err
	}
	count, err := s.repo.CountByRepoID(ctx, req.RepoID)
	if err != nil {
		return err
	}
	v := &domainversion.ResumeVersion{
		ID:       uuid.NewString(),
		RepoID:   req.RepoID,
		Version:  fmt.Sprintf("v1.%d.0", count),
		Title:    req.Title,
		Content:  req.Content,
		Remark:   req.Remark,
		IsActive: false,
	}
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return err
	}
	return s.repoReader.IncrementVersionCount(ctx, req.RepoID, 1)
}

func (s *ResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return nil, err
	}
	return s.repo.ListByRepoID(ctx, repoID)
}

func (s *ResumeVersionService) Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	return s.ensureVersionOwned(ctx, userID, id)
}

func (s *ResumeVersionService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	v.Title = req.Title
	v.Content = req.Content
	v.Remark = req.Remark
	v.IsActive = req.IsActive
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, v)
}

func (s *ResumeVersionService) Delete(ctx context.Context, userID, id string) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	inUse, err := s.checker.ExistsByResumeVersionID(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return domainerrors.ErrVersionInUse
	}
	if err := s.repo.Delete(ctx, id, v.RepoID); err != nil {
		return err
	}
	return s.repoReader.IncrementVersionCount(ctx, v.RepoID, -1)
}

func (s *ResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	return s.repo.SetActiveByID(ctx, v.RepoID, id)
}
```

Add `fmt` to imports:

```go
import (
	"context"
	"fmt"
```

- [ ] **Step 6: Run version service tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_version -count=1
```

Expected: PASS.

---

### Task 5: Align ResumeVersion handler and routes with frontend API

**Files:**
- Modify: `backend/internal/interface/handler/resume_version_handler.go`
- Modify: `backend/internal/interface/handler/resume_version_handler_test.go`
- Modify: `backend/internal/interface/router/router.go`
- Modify: `backend/internal/interface/router/router_test.go`

- [ ] **Step 1: Update handler service interface and tests**

In `backend/internal/interface/handler/resume_version_handler_test.go`, change mock methods to:

```go
func (m *mockResumeVersionService) Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error {
	return m.Called(ctx, userID, req).Error(0)
}

func (m *mockResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, repoID)
	return args.Get(0).([]domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error {
	return m.Called(ctx, userID, id, req).Error(0)
}

func (m *mockResumeVersionService) Delete(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *mockResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}
```

Remove `SetDefault` mock.

- [ ] **Step 2: Update handler tests for frontend routes and JSON**

Use routes:

```go
r.POST("/api/v1/resume-versions", h.Create)
r.GET("/api/v1/resume-versions", h.List)
r.GET("/api/v1/resume-versions/:id", h.Get)
r.PUT("/api/v1/resume-versions/:id", h.Update)
r.DELETE("/api/v1/resume-versions/:id", h.Delete)
r.POST("/api/v1/resume-versions/:id/activate", h.Activate)
```

Update create success expectation:

```go
svc.On("Create", mock.Anything, "u1", dto.CreateResumeVersionRequest{RepoID: "r1", Title: "v1", Content: "{}", Remark: "first"}).Return(nil)
```

Update create request body:

```go
req := httptest.NewRequest(http.MethodPost, "/api/v1/resume-versions", bytes.NewBufferString(`{"repoId":"r1","title":"v1","content":"{}","remark":"first"}`))
```

Update list URL and expectation:

```go
svc.On("List", mock.Anything, "u1", "r1").Return([]domainversion.ResumeVersion{
	{ID: "v1", RepoID: "r1", Version: "v1.0.0", Title: "Main", Content: "{}", Remark: "first", IsActive: true, CreatedAt: now, UpdatedAt: now},
}, nil)
req := httptest.NewRequest(http.MethodGet, "/api/v1/resume-versions?repoId=r1", nil)
```

Update response assertions:

```go
assert.Equal(t, "v1.0.0", version["version"])
assert.Equal(t, "first", version["remark"])
assert.Equal(t, true, version["isActive"])
assertNoJSONKeys(t, version, "VersionNum", "IsDefault", "Version", "Remark", "IsActive")
```

- [ ] **Step 3: Run handler tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/handler -run 'ResumeVersion' -count=1
```

Expected: FAIL because handler still uses nested route params and old service interface.

- [ ] **Step 4: Update handler implementation**

In `backend/internal/interface/handler/resume_version_handler.go`, change interface:

```go
type ResumeVersionService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error
	List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error)
	Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error
	Delete(ctx context.Context, userID, id string) error
	Activate(ctx context.Context, userID, id string) error
}
```

Remove `resumeRepoParam` and `resumeVersionParam` helpers.

Update methods:

```go
func (h *ResumeVersionHandler) Create(c *gin.Context) {
	var req dto.CreateResumeVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Create(c.Request.Context(), userID, req); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *ResumeVersionHandler) List(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	repoID := c.Query("repoId")
	if repoID == "" {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	rows, err := h.svc.List(c.Request.Context(), userID, repoID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToResumeVersionResponses(rows))
}

func (h *ResumeVersionHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToResumeVersionResponse(row))
}
```

Update `Update`, `Delete`, and `Activate` similarly using `c.Param("id")` and `h.svc.Activate`.

- [ ] **Step 5: Update router**

In `backend/internal/interface/router/router.go`, replace version routes with:

```go
auth.GET("/resume-versions", deps.ResumeVersionHandler.List)
auth.POST("/resume-versions", deps.ResumeVersionHandler.Create)
auth.GET("/resume-versions/:id", deps.ResumeVersionHandler.Get)
auth.PUT("/resume-versions/:id", deps.ResumeVersionHandler.Update)
auth.DELETE("/resume-versions/:id", deps.ResumeVersionHandler.Delete)
auth.POST("/resume-versions/:id/activate", deps.ResumeVersionHandler.Activate)
```

Keep repo routes unchanged.

- [ ] **Step 6: Update router tests**

In `backend/internal/interface/router/router_test.go`, update expected version paths to include:

```go
"GET /api/v1/resume-versions",
"POST /api/v1/resume-versions",
"GET /api/v1/resume-versions/:id",
"PUT /api/v1/resume-versions/:id",
"DELETE /api/v1/resume-versions/:id",
"POST /api/v1/resume-versions/:id/activate",
```

Remove expected nested version paths.

- [ ] **Step 7: Run handler/router tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/handler ./internal/interface/router -count=1
```

Expected: PASS.

---

### Task 6: Final integration and full backend verification

**Files:**
- Any compile-failing backend files from previous tasks.

- [ ] **Step 1: Run full backend tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./... -count=1
```

Expected: PASS.

- [ ] **Step 2: Fix compile errors from renamed fields**

If tests fail with errors like:

```text
v.VersionNum undefined
v.IsDefault undefined
SetDefaultByID undefined
MaxVersionNum undefined
```

Replace references as follows:

```text
VersionNum -> Version
IsDefault -> IsActive
SetDefaultByID -> SetActiveByID
MaxVersionNum -> CountByRepoID
SetDefault method names -> Activate
```

For fixtures, use:

```go
domainversion.ResumeVersion{
	ID:       "v1",
	RepoID:   "repo-1",
	Version:  "v1.0.0",
	Title:    "Main",
	Content:  "{}",
	Remark:   "first",
	IsActive: true,
}
```

- [ ] **Step 3: Run full backend tests again**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./... -count=1
```

Expected: PASS.

- [ ] **Step 4: Check git diff for accidental frontend changes**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && git diff --stat && git diff --name-only
```

Expected: changed files are backend domain/application/interface/infrastructure/migrations/tests and this plan document. No frontend files should be changed.

- [ ] **Step 5: Report verification evidence**

Report the final `go test ./... -count=1` result and summarize the backend contract now matching frontend fields/routes:

```text
Backend now returns repo fields jobType/targetPosition/versionCount and version fields version/remark/isActive. Version routes are /api/v1/resume-versions..., matching frontend API wrappers.
```
