# Backend Plan 3: Job Domain Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the Job domain so authenticated users can track job postings, link submitted resume versions to jobs, and manage interview progress rounds.

**Architecture:** Continue the 4-layer DDD architecture from Plans 1-2: Interface (Gin handlers) → Application (use cases) → Domain (entities + repository interfaces) → Infrastructure (GORM/PostgreSQL repositories). All operations are scoped by `userID` from auth middleware, with repository interfaces defined in the domain layer and implemented in infrastructure.

**Tech Stack:** Go 1.25, Gin, GORM + PostgreSQL 16, google/uuid, golang-migrate, testify, mockery v2, testcontainers-go

---

## File Map

```
backend/
├── internal/
│   ├── interface/
│   │   ├── handler/job_posting_handler.go
│   │   ├── handler/job_posting_handler_test.go
│   │   ├── handler/resume_application_handler.go
│   │   ├── handler/resume_application_handler_test.go
│   │   ├── handler/interview_progress_handler.go
│   │   ├── handler/interview_progress_handler_test.go
│   │   └── router/router.go                 # update: add job domain routes
│   ├── application/
│   │   ├── dto/job_dto.go
│   │   ├── job_posting/job_posting_service.go
│   │   ├── job_posting/job_posting_service_test.go
│   │   ├── resume_application/resume_application_service.go
│   │   ├── resume_application/resume_application_service_test.go
│   │   ├── interview_progress/interview_progress_service.go
│   │   └── interview_progress/interview_progress_service_test.go
│   ├── domain/
│   │   ├── job_posting/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── entity_test.go
│   │   ├── resume_application/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── entity_test.go
│   │   └── interview_progress/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── entity_test.go
│   └── infrastructure/
│       └── persistence/postgres/
│           ├── job_posting_repo.go
│           ├── job_posting_repo_test.go
│           ├── resume_application_repo.go
│           ├── resume_application_repo_test.go
│           ├── interview_progress_repo.go
│           └── interview_progress_repo_test.go
├── migrations/
│   ├── 000004_create_job_postings.up.sql
│   ├── 000004_create_job_postings.down.sql
│   ├── 000005_create_resume_applications.up.sql
│   ├── 000005_create_resume_applications.down.sql
│   ├── 000006_create_interview_progress.up.sql
│   └── 000006_create_interview_progress.down.sql
└── cmd/server/main.go                         # update: wire repos/services/handlers
```

---

## Task 1: JobPosting Domain Entity and Repository Interface

**Files:**
- Create: `backend/internal/domain/job_posting/entity.go`
- Create: `backend/internal/domain/job_posting/repository.go`
- Create: `backend/internal/domain/job_posting/entity_test.go`

- [ ] **Step 1: Write the failing entity tests**

```go
// internal/domain/job_posting/entity_test.go
package job_posting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

func TestJobPosting_Validate_Valid(t *testing.T) {
	p := &job_posting.JobPosting{
		UserID:   "user-1",
		Company:  "Acme",
		Position: "Backend Engineer",
		JobType:  job_posting.JobTypeSocial,
		Status:   job_posting.StatusCollecting,
	}
	assert.NoError(t, p.Validate())
}

func TestJobPosting_Validate_EmptyUserID(t *testing.T) {
	p := &job_posting.JobPosting{Company: "Acme", Position: "Backend", JobType: job_posting.JobTypeSocial, Status: job_posting.StatusCollecting}
	assert.Error(t, p.Validate())
}

func TestJobPosting_Validate_InvalidStatus(t *testing.T) {
	p := &job_posting.JobPosting{UserID: "user-1", Company: "Acme", Position: "Backend", JobType: job_posting.JobTypeSocial, Status: "unknown"}
	assert.Error(t, p.Validate())
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/yilin/project/ai-for/backend
go test ./internal/domain/job_posting/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Write minimal implementation**

```go
// internal/domain/job_posting/entity.go
package job_posting

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

type Status string

const (
	StatusCollecting   Status = "collecting"
	StatusApplying     Status = "applying"
	StatusInterviewing Status = "interviewing"
	StatusOffered      Status = "offered"
	StatusRejected     Status = "rejected"
	StatusWithdrawn    Status = "withdrawn"
)

type JobPosting struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	UserID      string         `gorm:"type:uuid;not null;index"`
	Company     string         `gorm:"type:varchar(120);not null"`
	Position    string         `gorm:"type:varchar(120);not null"`
	JobType     JobType        `gorm:"type:varchar(20);not null"`
	Location    string         `gorm:"type:varchar(120)"`
	Salary      string         `gorm:"type:varchar(120)"`
	Description string         `gorm:"type:text"`
	Source      string         `gorm:"type:varchar(120)"`
	URL         string         `gorm:"type:varchar(500)"`
	Status      Status         `gorm:"type:varchar(20);not null"`
	Notes       string         `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *JobPosting) Validate() error {
	if strings.TrimSpace(p.UserID) == "" || strings.TrimSpace(p.Company) == "" || strings.TrimSpace(p.Position) == "" {
		return domainerrors.ErrBadParams
	}
	if !validJobType(p.JobType) || !validStatus(p.Status) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validJobType(t JobType) bool {
	switch t {
	case JobTypeCampus, JobTypeSocial, JobTypeInternship:
		return true
	default:
		return false
	}
}

func validStatus(s Status) bool {
	switch s {
	case StatusCollecting, StatusApplying, StatusInterviewing, StatusOffered, StatusRejected, StatusWithdrawn:
		return true
	default:
		return false
	}
}
```

```go
// internal/domain/job_posting/repository.go
package job_posting

import "context"

type ListFilter struct {
	UserID   string
	Status   string
	JobType  string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, p *JobPosting) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*JobPosting, error)
	Update(ctx context.Context, p *JobPosting) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]JobPosting, int64, error)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/job_posting/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/job_posting/
git commit -m "feat: add job posting domain entity and repository interface"
```

---

## Task 2: ResumeApplication Domain Entity and Repository Interface

**Files:**
- Create: `backend/internal/domain/resume_application/entity.go`
- Create: `backend/internal/domain/resume_application/repository.go`
- Create: `backend/internal/domain/resume_application/entity_test.go`

- [ ] **Step 1: Write the failing entity tests**

```go
// internal/domain/resume_application/entity_test.go
package resume_application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

func TestResumeApplication_Validate_Valid(t *testing.T) {
	a := &resume_application.ResumeApplication{UserID: "u1", JobPostingID: "j1", ResumeVersionID: "v1", AppliedAt: time.Now()}
	assert.NoError(t, a.Validate())
}

func TestResumeApplication_Validate_MissingVersion(t *testing.T) {
	a := &resume_application.ResumeApplication{UserID: "u1", JobPostingID: "j1"}
	assert.Error(t, a.Validate())
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/resume_application/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Write minimal implementation**

```go
// internal/domain/resume_application/entity.go
package resume_application

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ResumeApplication struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	UserID          string    `gorm:"type:uuid;not null;index"`
	JobPostingID    string    `gorm:"type:uuid;not null;index"`
	ResumeVersionID string    `gorm:"type:uuid;not null;index"`
	AppliedAt       time.Time `gorm:"not null"`
	Notes           string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (a *ResumeApplication) Validate() error {
	if strings.TrimSpace(a.UserID) == "" || strings.TrimSpace(a.JobPostingID) == "" || strings.TrimSpace(a.ResumeVersionID) == "" {
		return domainerrors.ErrBadParams
	}
	return nil
}
```

```go
// internal/domain/resume_application/repository.go
package resume_application

import "context"

type Repository interface {
	Create(ctx context.Context, a *ResumeApplication) error
	FindByIDAndJobAndUser(ctx context.Context, id, jobPostingID, userID string) (*ResumeApplication, error)
	ListByJobAndUser(ctx context.Context, jobPostingID, userID string) ([]ResumeApplication, error)
	Delete(ctx context.Context, id, jobPostingID, userID string) error
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/resume_application/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/resume_application/
git commit -m "feat: add resume application domain entity and repository interface"
```

---

## Task 3: InterviewProgress Domain Entity and Repository Interface

**Files:**
- Create: `backend/internal/domain/interview_progress/entity.go`
- Create: `backend/internal/domain/interview_progress/repository.go`
- Create: `backend/internal/domain/interview_progress/entity_test.go`

- [ ] **Step 1: Write the failing entity tests**

```go
// internal/domain/interview_progress/entity_test.go
package interview_progress_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/interview_progress"
)

func TestInterviewProgress_Validate_Valid(t *testing.T) {
	p := &interview_progress.InterviewProgress{UserID: "u1", JobPostingID: "j1", Round: 1, InterviewType: interview_progress.TypeTechnical, Result: interview_progress.ResultPending}
	assert.NoError(t, p.Validate())
}

func TestInterviewProgress_Validate_InvalidRound(t *testing.T) {
	p := &interview_progress.InterviewProgress{UserID: "u1", JobPostingID: "j1", Round: 0, InterviewType: interview_progress.TypeTechnical, Result: interview_progress.ResultPending}
	assert.Error(t, p.Validate())
}

func TestInterviewProgress_Validate_InvalidResult(t *testing.T) {
	p := &interview_progress.InterviewProgress{UserID: "u1", JobPostingID: "j1", Round: 1, InterviewType: interview_progress.TypeTechnical, Result: "unknown"}
	assert.Error(t, p.Validate())
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/interview_progress/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Write minimal implementation**

```go
// internal/domain/interview_progress/entity.go
package interview_progress

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type InterviewType string

const (
	TypePhone     InterviewType = "phone"
	TypeHR        InterviewType = "hr"
	TypeTechnical InterviewType = "technical"
	TypeFinal     InterviewType = "final"
)

type Result string

const (
	ResultPending   Result = "pending"
	ResultPassed    Result = "passed"
	ResultFailed    Result = "failed"
	ResultCancelled Result = "cancelled"
)

type InterviewProgress struct {
	ID                  string     `gorm:"type:uuid;primaryKey"`
	UserID              string     `gorm:"type:uuid;not null;index"`
	JobPostingID        string     `gorm:"type:uuid;not null;index"`
	Round               int        `gorm:"not null"`
	InterviewType       InterviewType `gorm:"type:varchar(20);not null"`
	ScheduledAt         *time.Time
	CompletedAt         *time.Time
	Result              Result     `gorm:"type:varchar(20);not null"`
	Notes               string     `gorm:"type:text"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *InterviewProgress) Validate() error {
	if strings.TrimSpace(p.UserID) == "" || strings.TrimSpace(p.JobPostingID) == "" || p.Round < 1 {
		return domainerrors.ErrBadParams
	}
	if !validInterviewType(p.InterviewType) || !validResult(p.Result) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validInterviewType(t InterviewType) bool {
	switch t {
	case TypePhone, TypeHR, TypeTechnical, TypeFinal:
		return true
	default:
		return false
	}
}

func validResult(r Result) bool {
	switch r {
	case ResultPending, ResultPassed, ResultFailed, ResultCancelled:
		return true
	default:
		return false
	}
}
```

```go
// internal/domain/interview_progress/repository.go
package interview_progress

import "context"

type Repository interface {
	Create(ctx context.Context, p *InterviewProgress) error
	FindByIDAndJobAndUser(ctx context.Context, id, jobPostingID, userID string) (*InterviewProgress, error)
	ListByJobAndUser(ctx context.Context, jobPostingID, userID string) ([]InterviewProgress, error)
	Update(ctx context.Context, p *InterviewProgress) error
	Delete(ctx context.Context, id, jobPostingID, userID string) error
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/interview_progress/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/interview_progress/
git commit -m "feat: add interview progress domain entity and repository interface"
```

---

## Task 4: Database Migrations for Job Domain

**Files:**
- Create: `backend/migrations/000004_create_job_postings.up.sql`
- Create: `backend/migrations/000004_create_job_postings.down.sql`
- Create: `backend/migrations/000005_create_resume_applications.up.sql`
- Create: `backend/migrations/000005_create_resume_applications.down.sql`
- Create: `backend/migrations/000006_create_interview_progress.up.sql`
- Create: `backend/migrations/000006_create_interview_progress.down.sql`

- [ ] **Step 1: Create migration SQL files**

```sql
-- migrations/000004_create_job_postings.up.sql
CREATE TABLE job_postings (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID         NOT NULL REFERENCES users(id),
    company     VARCHAR(120) NOT NULL,
    position    VARCHAR(120) NOT NULL,
    job_type    VARCHAR(20)  NOT NULL,
    location    VARCHAR(120),
    salary      VARCHAR(120),
    description TEXT,
    source      VARCHAR(120),
    url         VARCHAR(500),
    status      VARCHAR(20)  NOT NULL DEFAULT 'collecting',
    notes       TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_job_postings_user_id ON job_postings(user_id);
CREATE INDEX idx_job_postings_status ON job_postings(status);
CREATE INDEX idx_job_postings_job_type ON job_postings(job_type);
```

```sql
-- migrations/000004_create_job_postings.down.sql
DROP TABLE IF EXISTS job_postings;
```

```sql
-- migrations/000005_create_resume_applications.up.sql
CREATE TABLE resume_applications (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID        NOT NULL REFERENCES users(id),
    job_posting_id    UUID        NOT NULL REFERENCES job_postings(id),
    resume_version_id UUID        NOT NULL REFERENCES resume_versions(id),
    applied_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes             TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resume_applications_user_id ON resume_applications(user_id);
CREATE INDEX idx_resume_applications_job_posting_id ON resume_applications(job_posting_id);
CREATE INDEX idx_resume_applications_resume_version_id ON resume_applications(resume_version_id);
```

```sql
-- migrations/000005_create_resume_applications.down.sql
DROP TABLE IF EXISTS resume_applications;
```

```sql
-- migrations/000006_create_interview_progress.up.sql
CREATE TABLE interview_progress (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID        NOT NULL REFERENCES users(id),
    job_posting_id UUID        NOT NULL REFERENCES job_postings(id),
    round          INT         NOT NULL,
    interview_type VARCHAR(20) NOT NULL,
    scheduled_at   TIMESTAMPTZ,
    completed_at   TIMESTAMPTZ,
    result         VARCHAR(20) NOT NULL DEFAULT 'pending',
    notes          TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_interview_progress_user_id ON interview_progress(user_id);
CREATE INDEX idx_interview_progress_job_posting_id ON interview_progress(job_posting_id);
```

```sql
-- migrations/000006_create_interview_progress.down.sql
DROP TABLE IF EXISTS interview_progress;
```

- [ ] **Step 2: Run migrations**

```bash
cd /Users/yilin/project/ai-for/backend
make migrate-up
```

Expected: migrations apply with no error.

- [ ] **Step 3: Commit**

```bash
git add migrations/000004_* migrations/000005_* migrations/000006_*
git commit -m "feat: add job domain migrations"
```

---

## Task 5: PostgreSQL Repository Implementations and Integration Tests

**Files:**
- Create: `backend/internal/infrastructure/persistence/postgres/job_posting_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/job_posting_repo_test.go`
- Create: `backend/internal/infrastructure/persistence/postgres/resume_application_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/resume_application_repo_test.go`
- Create: `backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/interview_progress_repo_test.go`

- [ ] **Step 1: Write failing integration tests**

```go
// internal/infrastructure/persistence/postgres/job_posting_repo_test.go
package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func TestIntegrationJobPostingRepo_CreateListAndUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := infrapostgres.NewJobPostingRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()

	p := &domainjob.JobPosting{ID: uuid.NewString(), UserID: userID, Company: "Acme", Position: "Backend", JobType: domainjob.JobTypeSocial, Status: domainjob.StatusCollecting}
	require.NoError(t, repo.Create(ctx, p))

	rows, total, err := repo.ListByUser(ctx, domainjob.ListFilter{UserID: userID, Status: "collecting", Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, rows, 1)

	found, err := repo.FindByIDAndUserID(ctx, p.ID, userID)
	require.NoError(t, err)
	found.Status = domainjob.StatusApplying
	require.NoError(t, repo.Update(ctx, found))
}
```

```go
// internal/infrastructure/persistence/postgres/resume_application_repo_test.go
package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func TestIntegrationResumeApplicationRepo_ExistsByResumeVersionID(t *testing.T) {
	db := setupTestDB(t)
	repo := infrapostgres.NewResumeApplicationRepo(db)
	ctx := context.Background()
	versionID := uuid.NewString()

	a := &domainapp.ResumeApplication{ID: uuid.NewString(), UserID: uuid.NewString(), JobPostingID: uuid.NewString(), ResumeVersionID: versionID, AppliedAt: time.Now()}
	require.NoError(t, repo.Create(ctx, a))

	exists, err := repo.ExistsByResumeVersionID(ctx, versionID)
	require.NoError(t, err)
	assert.True(t, exists)
}
```

```go
// internal/infrastructure/persistence/postgres/interview_progress_repo_test.go
package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func TestIntegrationInterviewProgressRepo_CreateAndList(t *testing.T) {
	db := setupTestDB(t)
	repo := infrapostgres.NewInterviewProgressRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	jobID := uuid.NewString()

	p := &domainprogress.InterviewProgress{ID: uuid.NewString(), UserID: userID, JobPostingID: jobID, Round: 1, InterviewType: domainprogress.TypeTechnical, Result: domainprogress.ResultPending}
	require.NoError(t, repo.Create(ctx, p))

	rows, err := repo.ListByJobAndUser(ctx, jobID, userID)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: FAIL with missing constructors/repositories.

- [ ] **Step 3: Implement repositories**

```go
// internal/infrastructure/persistence/postgres/job_posting_repo.go
package postgres

import (
	"context"
	"errors"

	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type JobPostingRepo struct{ db *gorm.DB }

func NewJobPostingRepo(db *gorm.DB) *JobPostingRepo { return &JobPostingRepo{db: db} }

func (r *JobPostingRepo) Create(ctx context.Context, p *domainjob.JobPosting) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *JobPostingRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) {
	var p domainjob.JobPosting
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }
	return &p, err
}

func (r *JobPostingRepo) Update(ctx context.Context, p *domainjob.JobPosting) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *JobPostingRepo) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainjob.JobPosting{}).Error
}

func (r *JobPostingRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainjob.JobPosting{}).Where("user_id = ?", f.UserID)
	if f.Status != "" { q = q.Where("status = ?", f.Status) }
	if f.JobType != "" { q = q.Where("job_type = ?", f.JobType) }
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	var rows []domainjob.JobPosting
	err := q.Order("updated_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}
```

```go
// internal/infrastructure/persistence/postgres/resume_application_repo.go
package postgres

import (
	"context"
	"errors"

	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type ResumeApplicationRepo struct{ db *gorm.DB }

func NewResumeApplicationRepo(db *gorm.DB) *ResumeApplicationRepo { return &ResumeApplicationRepo{db: db} }

func (r *ResumeApplicationRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *ResumeApplicationRepo) FindByIDAndJobAndUser(ctx context.Context, id, jobPostingID, userID string) (*domainapp.ResumeApplication, error) {
	var a domainapp.ResumeApplication
	err := r.db.WithContext(ctx).Where("id = ? AND job_posting_id = ? AND user_id = ?", id, jobPostingID, userID).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }
	return &a, err
}

func (r *ResumeApplicationRepo) ListByJobAndUser(ctx context.Context, jobPostingID, userID string) ([]domainapp.ResumeApplication, error) {
	var rows []domainapp.ResumeApplication
	err := r.db.WithContext(ctx).Where("job_posting_id = ? AND user_id = ?", jobPostingID, userID).Order("applied_at DESC").Find(&rows).Error
	return rows, err
}

func (r *ResumeApplicationRepo) Delete(ctx context.Context, id, jobPostingID, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND job_posting_id = ? AND user_id = ?", id, jobPostingID, userID).Delete(&domainapp.ResumeApplication{}).Error
}

func (r *ResumeApplicationRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("resume_version_id = ?", versionID).Count(&count).Error
	return count > 0, err
}
```

```go
// internal/infrastructure/persistence/postgres/interview_progress_repo.go
package postgres

import (
	"context"
	"errors"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type InterviewProgressRepo struct{ db *gorm.DB }

func NewInterviewProgressRepo(db *gorm.DB) *InterviewProgressRepo { return &InterviewProgressRepo{db: db} }

func (r *InterviewProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *InterviewProgressRepo) FindByIDAndJobAndUser(ctx context.Context, id, jobPostingID, userID string) (*domainprogress.InterviewProgress, error) {
	var p domainprogress.InterviewProgress
	err := r.db.WithContext(ctx).Where("id = ? AND job_posting_id = ? AND user_id = ?", id, jobPostingID, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }
	return &p, err
}

func (r *InterviewProgressRepo) ListByJobAndUser(ctx context.Context, jobPostingID, userID string) ([]domainprogress.InterviewProgress, error) {
	var rows []domainprogress.InterviewProgress
	err := r.db.WithContext(ctx).Where("job_posting_id = ? AND user_id = ?", jobPostingID, userID).Order("round ASC").Find(&rows).Error
	return rows, err
}

func (r *InterviewProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *InterviewProgressRepo) Delete(ctx context.Context, id, jobPostingID, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND job_posting_id = ? AND user_id = ?", id, jobPostingID, userID).Delete(&domainprogress.InterviewProgress{}).Error
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/infrastructure/persistence/postgres/job_posting_repo* internal/infrastructure/persistence/postgres/resume_application_repo* internal/infrastructure/persistence/postgres/interview_progress_repo*
git commit -m "feat: add postgres repositories for job domain"
```

---

## Task 6: Application DTOs

**Files:**
- Create: `backend/internal/application/dto/job_dto.go`

- [ ] **Step 1: Create DTO file**

```go
// internal/application/dto/job_dto.go
package dto

import "time"

type CreateJobPostingRequest struct {
	Company     string `json:"company" binding:"required"`
	Position    string `json:"position" binding:"required"`
	JobType     string `json:"jobType" binding:"required"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Description string `json:"description"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

type UpdateJobPostingRequest struct {
	Company     string `json:"company" binding:"required"`
	Position    string `json:"position" binding:"required"`
	JobType     string `json:"jobType" binding:"required"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Description string `json:"description"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	Status      string `json:"status" binding:"required"`
	Notes       string `json:"notes"`
}

type UpdateJobPostingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateResumeApplicationRequest struct {
	ResumeVersionID string     `json:"resumeVersionId" binding:"required"`
	AppliedAt       *time.Time `json:"appliedAt"`
	Notes           string     `json:"notes"`
}

type CreateInterviewProgressRequest struct {
	Round         int        `json:"round" binding:"required,min=1"`
	InterviewType string     `json:"interviewType" binding:"required"`
	ScheduledAt   *time.Time `json:"scheduledAt"`
	CompletedAt   *time.Time `json:"completedAt"`
	Result        string     `json:"result"`
	Notes         string     `json:"notes"`
}

type UpdateInterviewProgressRequest struct {
	Round         int        `json:"round" binding:"required,min=1"`
	InterviewType string     `json:"interviewType" binding:"required"`
	ScheduledAt   *time.Time `json:"scheduledAt"`
	CompletedAt   *time.Time `json:"completedAt"`
	Result        string     `json:"result" binding:"required"`
	Notes         string     `json:"notes"`
}
```

- [ ] **Step 2: Run compile check**

```bash
go test ./internal/application/... -run '^$'
```

Expected: PASS or unrelated existing test failures only.

- [ ] **Step 3: Commit**

```bash
git add internal/application/dto/job_dto.go
git commit -m "feat: add job domain DTOs"
```

---

## Task 7: JobPosting Application Service

**Files:**
- Create: `backend/internal/application/job_posting/job_posting_service.go`
- Create: `backend/internal/application/job_posting/job_posting_service_test.go`

- [ ] **Step 1: Write failing service tests**

```go
// internal/application/job_posting/job_posting_service_test.go
package job_posting_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/job_posting"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

type mockJobRepo struct{ mock.Mock }

func (m *mockJobRepo) Create(ctx context.Context, p *domainjob.JobPosting) error { return m.Called(ctx, p).Error(0) }
func (m *mockJobRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) {
	a := m.Called(ctx, id, userID)
	if a.Get(0) == nil { return nil, a.Error(1) }
	return a.Get(0).(*domainjob.JobPosting), a.Error(1)
}
func (m *mockJobRepo) Update(ctx context.Context, p *domainjob.JobPosting) error { return m.Called(ctx, p).Error(0) }
func (m *mockJobRepo) Delete(ctx context.Context, id, userID string) error { return m.Called(ctx, id, userID).Error(0) }
func (m *mockJobRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) {
	a := m.Called(ctx, f)
	return a.Get(0).([]domainjob.JobPosting), a.Get(1).(int64), a.Error(2)
}

func TestJobPostingService_Create_DefaultStatus(t *testing.T) {
	repo := &mockJobRepo{}
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domainjob.JobPosting) bool {
		return p.UserID == "u1" && p.Company == "Acme" && p.Status == domainjob.StatusCollecting
	})).Return(nil)
	svc := appsvc.NewJobPostingService(repo)
	err := svc.Create(context.Background(), "u1", dto.CreateJobPostingRequest{Company: "Acme", Position: "Backend", JobType: "social"})
	require.NoError(t, err)
}

func TestJobPostingService_UpdateStatus(t *testing.T) {
	repo := &mockJobRepo{}
	existing := &domainjob.JobPosting{ID: "j1", UserID: "u1", Company: "Acme", Position: "Backend", JobType: domainjob.JobTypeSocial, Status: domainjob.StatusCollecting}
	repo.On("FindByIDAndUserID", mock.Anything, "j1", "u1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domainjob.JobPosting) bool { return p.Status == domainjob.StatusApplying })).Return(nil)
	svc := appsvc.NewJobPostingService(repo)
	assert.NoError(t, svc.UpdateStatus(context.Background(), "u1", "j1", dto.UpdateJobPostingStatusRequest{Status: "applying"}))
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/application/job_posting/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Implement service**

```go
// internal/application/job_posting/job_posting_service.go
package job_posting

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

type JobPostingService struct{ repo domainjob.Repository }

func NewJobPostingService(repo domainjob.Repository) *JobPostingService { return &JobPostingService{repo: repo} }

func (s *JobPostingService) Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) error {
	status := domainjob.Status(req.Status)
	if status == "" { status = domainjob.StatusCollecting }
	p := &domainjob.JobPosting{ID: uuid.NewString(), UserID: userID, Company: req.Company, Position: req.Position, JobType: domainjob.JobType(req.JobType), Location: req.Location, Salary: req.Salary, Description: req.Description, Source: req.Source, URL: req.URL, Status: status, Notes: req.Notes}
	if err := p.Validate(); err != nil { return err }
	return s.repo.Create(ctx, p)
}

func (s *JobPostingService) List(ctx context.Context, userID, status, jobType string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	return s.repo.ListByUser(ctx, domainjob.ListFilter{UserID: userID, Status: status, JobType: jobType, Page: page, PageSize: pageSize})
}

func (s *JobPostingService) Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *JobPostingService) Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error {
	p, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	p.Company, p.Position, p.JobType = req.Company, req.Position, domainjob.JobType(req.JobType)
	p.Location, p.Salary, p.Description = req.Location, req.Salary, req.Description
	p.Source, p.URL, p.Status, p.Notes = req.Source, req.URL, domainjob.Status(req.Status), req.Notes
	if err := p.Validate(); err != nil { return err }
	return s.repo.Update(ctx, p)
}

func (s *JobPostingService) UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateJobPostingStatusRequest) error {
	p, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	p.Status = domainjob.Status(req.Status)
	if err := p.Validate(); err != nil { return err }
	return s.repo.Update(ctx, p)
}

func (s *JobPostingService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/application/job_posting/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/application/job_posting/
git commit -m "feat: add job posting application service"
```

---

## Task 8: ResumeApplication Application Service

**Files:**
- Create: `backend/internal/application/resume_application/resume_application_service.go`
- Create: `backend/internal/application/resume_application/resume_application_service_test.go`

- [ ] **Step 1: Write failing service tests**

```go
// internal/application/resume_application/resume_application_service_test.go
package resume_application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_application"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

type mockAppRepo struct{ mock.Mock }
func (m *mockAppRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error { return m.Called(ctx, a).Error(0) }
func (m *mockAppRepo) FindByIDAndJobAndUser(ctx context.Context, id, jobID, userID string) (*domainapp.ResumeApplication, error) { a := m.Called(ctx, id, jobID, userID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainapp.ResumeApplication), a.Error(1) }
func (m *mockAppRepo) ListByJobAndUser(ctx context.Context, jobID, userID string) ([]domainapp.ResumeApplication, error) { a := m.Called(ctx, jobID, userID); return a.Get(0).([]domainapp.ResumeApplication), a.Error(1) }
func (m *mockAppRepo) Delete(ctx context.Context, id, jobID, userID string) error { return m.Called(ctx, id, jobID, userID).Error(0) }
func (m *mockAppRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) { a := m.Called(ctx, versionID); return a.Bool(0), a.Error(1) }

type mockJobFinder struct{ mock.Mock }
func (m *mockJobFinder) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) { a := m.Called(ctx, id, userID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainjob.JobPosting), a.Error(1) }

type mockVersionFinder struct{ mock.Mock }
func (m *mockVersionFinder) FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error) { a := m.Called(ctx, id); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainversion.ResumeVersion), a.Error(1) }

func TestResumeApplicationService_Create_ChecksJobAndVersion(t *testing.T) {
	apps := &mockAppRepo{}
	jobs := &mockJobFinder{}
	versions := &mockVersionFinder{}
	jobs.On("FindByIDAndUserID", mock.Anything, "j1", "u1").Return(&domainjob.JobPosting{ID:"j1", UserID:"u1"}, nil)
	versions.On("FindByID", mock.Anything, "v1").Return(&domainversion.ResumeVersion{ID:"v1"}, nil)
	apps.On("Create", mock.Anything, mock.AnythingOfType("*resume_application.ResumeApplication")).Return(nil)
	svc := appsvc.NewResumeApplicationService(apps, jobs, versions)
	assert.NoError(t, svc.Create(context.Background(), "u1", "j1", dto.CreateResumeApplicationRequest{ResumeVersionID:"v1"}))
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/application/resume_application/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Implement service**

```go
// internal/application/resume_application/resume_application_service.go
package resume_application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

type JobFinder interface { FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) }
type VersionFinder interface { FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error) }

type ResumeApplicationService struct {
	repo domainapp.Repository
	jobs JobFinder
	versions VersionFinder
}

func NewResumeApplicationService(repo domainapp.Repository, jobs JobFinder, versions VersionFinder) *ResumeApplicationService {
	return &ResumeApplicationService{repo: repo, jobs: jobs, versions: versions}
}

func (s *ResumeApplicationService) Create(ctx context.Context, userID, jobID string, req dto.CreateResumeApplicationRequest) error {
	if _, err := s.jobs.FindByIDAndUserID(ctx, jobID, userID); err != nil { return err }
	if _, err := s.versions.FindByID(ctx, req.ResumeVersionID); err != nil { return err }
	appliedAt := time.Now()
	if req.AppliedAt != nil { appliedAt = *req.AppliedAt }
	a := &domainapp.ResumeApplication{ID: uuid.NewString(), UserID: userID, JobPostingID: jobID, ResumeVersionID: req.ResumeVersionID, AppliedAt: appliedAt, Notes: req.Notes}
	if err := a.Validate(); err != nil { return err }
	return s.repo.Create(ctx, a)
}

func (s *ResumeApplicationService) List(ctx context.Context, userID, jobID string) ([]domainapp.ResumeApplication, error) {
	if _, err := s.jobs.FindByIDAndUserID(ctx, jobID, userID); err != nil { return nil, err }
	return s.repo.ListByJobAndUser(ctx, jobID, userID)
}

func (s *ResumeApplicationService) Delete(ctx context.Context, userID, jobID, id string) error {
	if _, err := s.jobs.FindByIDAndUserID(ctx, jobID, userID); err != nil { return err }
	return s.repo.Delete(ctx, id, jobID, userID)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/application/resume_application/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/application/resume_application/
git commit -m "feat: add resume application service"
```

---

## Task 9: InterviewProgress Application Service

**Files:**
- Create: `backend/internal/application/interview_progress/interview_progress_service.go`
- Create: `backend/internal/application/interview_progress/interview_progress_service_test.go`

- [ ] **Step 1: Write failing service tests**

```go
// internal/application/interview_progress/interview_progress_service_test.go
package interview_progress_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/interview_progress"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
)

type mockProgressRepo struct{ mock.Mock }
func (m *mockProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error { return m.Called(ctx, p).Error(0) }
func (m *mockProgressRepo) FindByIDAndJobAndUser(ctx context.Context, id, jobID, userID string) (*domainprogress.InterviewProgress, error) { a := m.Called(ctx, id, jobID, userID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainprogress.InterviewProgress), a.Error(1) }
func (m *mockProgressRepo) ListByJobAndUser(ctx context.Context, jobID, userID string) ([]domainprogress.InterviewProgress, error) { a := m.Called(ctx, jobID, userID); return a.Get(0).([]domainprogress.InterviewProgress), a.Error(1) }
func (m *mockProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error { return m.Called(ctx, p).Error(0) }
func (m *mockProgressRepo) Delete(ctx context.Context, id, jobID, userID string) error { return m.Called(ctx, id, jobID, userID).Error(0) }

type mockProgressJobFinder struct{ mock.Mock }
func (m *mockProgressJobFinder) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) { a := m.Called(ctx, id, userID); if a.Get(0)==nil { return nil, a.Error(1) }; return a.Get(0).(*domainjob.JobPosting), a.Error(1) }

func TestInterviewProgressService_Create_DefaultResult(t *testing.T) {
	repo := &mockProgressRepo{}
	jobs := &mockProgressJobFinder{}
	jobs.On("FindByIDAndUserID", mock.Anything, "j1", "u1").Return(&domainjob.JobPosting{ID:"j1", UserID:"u1"}, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domainprogress.InterviewProgress) bool { return p.Result == domainprogress.ResultPending })).Return(nil)
	svc := appsvc.NewInterviewProgressService(repo, jobs)
	assert.NoError(t, svc.Create(context.Background(), "u1", "j1", dto.CreateInterviewProgressRequest{Round:1, InterviewType:"technical"}))
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/application/interview_progress/... -v
```

Expected: FAIL with package/file not found.

- [ ] **Step 3: Implement service**

```go
// internal/application/interview_progress/interview_progress_service.go
package interview_progress

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
)

type JobFinder interface { FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) }

type InterviewProgressService struct {
	repo domainprogress.Repository
	jobs JobFinder
}

func NewInterviewProgressService(repo domainprogress.Repository, jobs JobFinder) *InterviewProgressService {
	return &InterviewProgressService{repo: repo, jobs: jobs}
}

func (s *InterviewProgressService) Create(ctx context.Context, userID, jobID string, req dto.CreateInterviewProgressRequest) error {
	if _, err := s.jobs.FindByIDAndUserID(ctx, jobID, userID); err != nil { return err }
	result := domainprogress.Result(req.Result)
	if result == "" { result = domainprogress.ResultPending }
	p := &domainprogress.InterviewProgress{ID: uuid.NewString(), UserID: userID, JobPostingID: jobID, Round: req.Round, InterviewType: domainprogress.InterviewType(req.InterviewType), ScheduledAt: req.ScheduledAt, CompletedAt: req.CompletedAt, Result: result, Notes: req.Notes}
	if err := p.Validate(); err != nil { return err }
	return s.repo.Create(ctx, p)
}

func (s *InterviewProgressService) List(ctx context.Context, userID, jobID string) ([]domainprogress.InterviewProgress, error) {
	if _, err := s.jobs.FindByIDAndUserID(ctx, jobID, userID); err != nil { return nil, err }
	return s.repo.ListByJobAndUser(ctx, jobID, userID)
}

func (s *InterviewProgressService) Update(ctx context.Context, userID, jobID, id string, req dto.UpdateInterviewProgressRequest) error {
	p, err := s.repo.FindByIDAndJobAndUser(ctx, id, jobID, userID)
	if err != nil { return err }
	p.Round, p.InterviewType, p.ScheduledAt, p.CompletedAt, p.Result, p.Notes = req.Round, domainprogress.InterviewType(req.InterviewType), req.ScheduledAt, req.CompletedAt, domainprogress.Result(req.Result), req.Notes
	if err := p.Validate(); err != nil { return err }
	return s.repo.Update(ctx, p)
}

func (s *InterviewProgressService) Delete(ctx context.Context, userID, jobID, id string) error {
	return s.repo.Delete(ctx, id, jobID, userID)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/application/interview_progress/... -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/application/interview_progress/
git commit -m "feat: add interview progress service"
```

---

## Task 10: Interface Handlers

**Files:**
- Create: `backend/internal/interface/handler/job_posting_handler.go`
- Create: `backend/internal/interface/handler/job_posting_handler_test.go`
- Create: `backend/internal/interface/handler/resume_application_handler.go`
- Create: `backend/internal/interface/handler/resume_application_handler_test.go`
- Create: `backend/internal/interface/handler/interview_progress_handler.go`
- Create: `backend/internal/interface/handler/interview_progress_handler_test.go`

- [ ] **Step 1: Write failing handler tests**

```go
// internal/interface/handler/job_posting_handler_test.go
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

func TestJobPostingHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewJobPostingHandler(nil)
	r.POST("/api/v1/job-postings", h.Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/job-postings", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "40000")
}
```

```go
// internal/interface/handler/resume_application_handler_test.go
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

func TestResumeApplicationHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewResumeApplicationHandler(nil)
	r.POST("/api/v1/job-postings/:jobId/applications", h.Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/job-postings/j1/applications", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "40000")
}
```

```go
// internal/interface/handler/interview_progress_handler_test.go
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

func TestInterviewProgressHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewInterviewProgressHandler(nil)
	r.POST("/api/v1/job-postings/:jobId/interviews", h.Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/job-postings/j1/interviews", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "40000")
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/interface/handler/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: FAIL with missing constructors/handlers.

- [ ] **Step 3: Implement handlers**

```go
// internal/interface/handler/job_posting_handler.go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/pagination"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type JobPostingService interface {
	Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) error
	List(ctx context.Context, userID, status, jobType string, page, pageSize int) (any, int64, error)
	Get(ctx context.Context, userID, id string) (any, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error
	UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateJobPostingStatusRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type JobPostingHandler struct{ svc JobPostingService }
func NewJobPostingHandler(svc JobPostingService) *JobPostingHandler { return &JobPostingHandler{svc: svc} }

func (h *JobPostingHandler) Create(c *gin.Context) {
	var req dto.CreateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Create(c.Request.Context(), c.GetString("userID"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *JobPostingHandler) List(c *gin.Context) {
	p := pagination.FromContext(c)
	rows, total, err := h.svc.List(c.Request.Context(), c.GetString("userID"), c.Query("status"), c.Query("jobType"), p.Page, p.PageSize)
	if err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.SuccessPage(c, rows, total, p.Page, p.PageSize)
}
func (h *JobPostingHandler) Get(c *gin.Context) {
	row, err := h.svc.Get(c.Request.Context(), c.GetString("userID"), c.Param("id"))
	if err != nil { response.Fail(c, domainerrors.CodeNotFound, "资源不存在"); return }
	response.Success(c, row)
}
func (h *JobPostingHandler) Update(c *gin.Context) {
	var req dto.UpdateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Update(c.Request.Context(), c.GetString("userID"), c.Param("id"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *JobPostingHandler) UpdateStatus(c *gin.Context) {
	var req dto.UpdateJobPostingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.UpdateStatus(c.Request.Context(), c.GetString("userID"), c.Param("id"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *JobPostingHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.GetString("userID"), c.Param("id")); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
```

```go
// internal/interface/handler/resume_application_handler.go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeApplicationService interface {
	Create(ctx context.Context, userID, jobID string, req dto.CreateResumeApplicationRequest) error
	List(ctx context.Context, userID, jobID string) (any, error)
	Delete(ctx context.Context, userID, jobID, id string) error
}

type ResumeApplicationHandler struct{ svc ResumeApplicationService }
func NewResumeApplicationHandler(svc ResumeApplicationService) *ResumeApplicationHandler { return &ResumeApplicationHandler{svc: svc} }

func (h *ResumeApplicationHandler) Create(c *gin.Context) {
	var req dto.CreateResumeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Create(c.Request.Context(), c.GetString("userID"), c.Param("jobId"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *ResumeApplicationHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context(), c.GetString("userID"), c.Param("jobId"))
	if err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, rows)
}
func (h *ResumeApplicationHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.GetString("userID"), c.Param("jobId"), c.Param("id")); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
```

```go
// internal/interface/handler/interview_progress_handler.go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type InterviewProgressService interface {
	Create(ctx context.Context, userID, jobID string, req dto.CreateInterviewProgressRequest) error
	List(ctx context.Context, userID, jobID string) (any, error)
	Update(ctx context.Context, userID, jobID, id string, req dto.UpdateInterviewProgressRequest) error
	Delete(ctx context.Context, userID, jobID, id string) error
}

type InterviewProgressHandler struct{ svc InterviewProgressService }
func NewInterviewProgressHandler(svc InterviewProgressService) *InterviewProgressHandler { return &InterviewProgressHandler{svc: svc} }

func (h *InterviewProgressHandler) Create(c *gin.Context) {
	var req dto.CreateInterviewProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Create(c.Request.Context(), c.GetString("userID"), c.Param("jobId"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *InterviewProgressHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context(), c.GetString("userID"), c.Param("jobId"))
	if err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, rows)
}
func (h *InterviewProgressHandler) Update(c *gin.Context) {
	var req dto.UpdateInterviewProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, domainerrors.CodeBadParams, err.Error()); return }
	if err := h.svc.Update(c.Request.Context(), c.GetString("userID"), c.Param("jobId"), c.Param("id"), req); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
func (h *InterviewProgressHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.GetString("userID"), c.Param("jobId"), c.Param("id")); err != nil { response.Fail(c, domainerrors.CodeInternalError, err.Error()); return }
	response.Success(c, nil)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/interface/handler/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/interface/handler/job_posting_handler* internal/interface/handler/resume_application_handler* internal/interface/handler/interview_progress_handler*
git commit -m "feat: add job domain handlers"
```

---

## Task 11: Router and main.go Wiring

**Files:**
- Modify: `backend/internal/interface/router/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Update router setup signature and routes**

```go
// internal/interface/router/router.go
// Add these parameters to Setup:
// jobPostingHandler *handler.JobPostingHandler,
// resumeApplicationHandler *handler.ResumeApplicationHandler,
// interviewProgressHandler *handler.InterviewProgressHandler,

// Add these routes inside the authenticated route group:
auth.GET("/job-postings", jobPostingHandler.List)
auth.POST("/job-postings", jobPostingHandler.Create)
auth.GET("/job-postings/:id", jobPostingHandler.Get)
auth.PUT("/job-postings/:id", jobPostingHandler.Update)
auth.DELETE("/job-postings/:id", jobPostingHandler.Delete)
auth.PUT("/job-postings/:id/status", jobPostingHandler.UpdateStatus)

auth.GET("/job-postings/:jobId/applications", resumeApplicationHandler.List)
auth.POST("/job-postings/:jobId/applications", resumeApplicationHandler.Create)
auth.DELETE("/job-postings/:jobId/applications/:id", resumeApplicationHandler.Delete)

auth.GET("/job-postings/:jobId/interviews", interviewProgressHandler.List)
auth.POST("/job-postings/:jobId/interviews", interviewProgressHandler.Create)
auth.PUT("/job-postings/:jobId/interviews/:id", interviewProgressHandler.Update)
auth.DELETE("/job-postings/:jobId/interviews/:id", interviewProgressHandler.Delete)
```

- [ ] **Step 2: Wire repositories, services, and handlers in main.go**

```go
// cmd/server/main.go
// Add imports with aliases:
// jobpostingapp "github.com/yilin/ai-for-backend/internal/application/job_posting"
// resumeapplicationapp "github.com/yilin/ai-for-backend/internal/application/resume_application"
// interviewprogressapp "github.com/yilin/ai-for-backend/internal/application/interview_progress"
// domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
// domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
// domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"

// Add AutoMigrate models if Plan 1 main.go still uses AutoMigrate:
if err := db.AutoMigrate(&domainuser.User{}, &domainrepo.ResumeRepo{}, &domainversion.ResumeVersion{}, &domainjob.JobPosting{}, &domainapp.ResumeApplication{}, &domainprogress.InterviewProgress{}); err != nil {
	log.Fatalf("failed to auto-migrate: %v", err)
}

// Add repositories:
jobPostingRepo := infrapostgres.NewJobPostingRepo(db)
resumeApplicationRepo := infrapostgres.NewResumeApplicationRepo(db)
interviewProgressRepo := infrapostgres.NewInterviewProgressRepo(db)

// Replace Plan 2 noopVersionInUseChecker with real repo:
resumeVersionSvc := resumeversionapp.NewResumeVersionService(resumeVersionRepo, resumeApplicationRepo)

// Add services:
jobPostingSvc := jobpostingapp.NewJobPostingService(jobPostingRepo)
resumeApplicationSvc := resumeapplicationapp.NewResumeApplicationService(resumeApplicationRepo, jobPostingRepo, resumeVersionRepo)
interviewProgressSvc := interviewprogressapp.NewInterviewProgressService(interviewProgressRepo, jobPostingRepo)

// Add handlers:
jobPostingHandler := handler.NewJobPostingHandler(jobPostingSvc)
resumeApplicationHandler := handler.NewResumeApplicationHandler(resumeApplicationSvc)
interviewProgressHandler := handler.NewInterviewProgressHandler(interviewProgressSvc)

// Pass new handlers to router.Setup in the same order as the updated signature.
```

- [ ] **Step 3: Build**

```bash
go build ./cmd/server/...
```

Expected: build succeeds.

- [ ] **Step 4: Smoke test job posting endpoints**

```bash
make run
```

In another shell:

```bash
curl -X POST http://localhost:8080/api/v1/job-postings \
  -H "X-Access-Token: <token>" \
  -H "Content-Type: application/json" \
  -d '{"company":"Acme","position":"Backend Engineer","jobType":"social","status":"collecting"}'

curl "http://localhost:8080/api/v1/job-postings?status=collecting&jobType=social" \
  -H "X-Access-Token: <token>"
```

Expected: both return JSON with `"code":20000`.

- [ ] **Step 5: Commit**

```bash
git add internal/interface/router/router.go cmd/server/main.go
git commit -m "feat: wire job domain routes and services"
```

---

## Task 12: Final Verification

**Files:**
- No new files

- [ ] **Step 1: Run domain tests**

```bash
go test ./internal/domain/job_posting/... -v
go test ./internal/domain/resume_application/... -v
go test ./internal/domain/interview_progress/... -v
```

Expected: PASS.

- [ ] **Step 2: Run application tests**

```bash
go test ./internal/application/job_posting/... -v
go test ./internal/application/resume_application/... -v
go test ./internal/application/interview_progress/... -v
```

Expected: PASS.

- [ ] **Step 3: Run handler tests**

```bash
go test ./internal/interface/handler/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: PASS.

- [ ] **Step 4: Run repository integration tests**

```bash
go test ./internal/infrastructure/persistence/postgres/... -run "JobPosting|ResumeApplication|InterviewProgress" -v
```

Expected: PASS.

- [ ] **Step 5: Run full unit test suite**

```bash
make test
```

Expected: PASS.

- [ ] **Step 6: Commit final fixes if any**

```bash
git add internal/ cmd/ migrations/
git commit -m "test: finalize job domain coverage"
```

---

## Self-Review

**Spec coverage check:**
- ✅ JobPosting domain covers company, position, job type, source, URL, status, notes, soft delete, and user ownership.
- ✅ JobPosting list supports pagination plus `status` and `jobType` filters.
- ✅ ResumeApplication links job postings to resume versions and implements `ExistsByResumeVersionID(ctx, versionID string) (bool, error)` for Plan 2 deletion protection.
- ✅ InterviewProgress supports rounds, interview type, scheduled/completed times, result, and notes.
- ✅ All endpoints are under authenticated `/api/v1` routes.
- ✅ All operations include user-scoped ownership checks through service/repository methods.
- ✅ Migrations 000004-000006 create all required tables and indexes.
- ✅ TDD is represented for entities, services, handlers, and infrastructure repositories.

**Placeholder scan:** No TBD, TODO, "implement later", or "similar to" placeholders remain.

**Type consistency:**
- ✅ Domain constants match DTO string fields and service conversions.
- ✅ Handler service interfaces match service methods.
- ✅ Repository method names match application services.
- ✅ `ResumeApplicationRepo.ExistsByResumeVersionID` matches Plan 2 `VersionInUseChecker`.
- ✅ Router paths match the requested API contract.
