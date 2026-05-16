# Backend Plan 3: Frontend-Aligned Job Domain Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the backend job-domain APIs so they match the existing frontend job posting, resume application, and interview progress interfaces exactly.

**Architecture:** Keep the current backend DDD layering: domain entities define persisted state and validation, application DTOs translate between Go and frontend JSON, services enforce user ownership and cross-aggregate rules, GORM repositories persist snake_case PostgreSQL columns, and Gin handlers expose the frontend API contract. The frontend is the source of truth for route paths, JSON names, query parameters, status enums, and response shapes; do not change frontend files for this plan.

**Tech Stack:** Go, Gin, GORM, PostgreSQL migrations, google/uuid, testify/mock, existing `pkg/response`, existing pagination conventions.

---

## Frontend Contract to Preserve

The backend must adapt to these existing frontend files:

- `frontend/src/models/job-posting.ts`
- `frontend/src/api/job-postings.ts`
- `frontend/mock/job-postings.ts`
- `frontend/src/models/resume-application.ts`
- `frontend/src/api/resume-applications.ts`
- `frontend/mock/resume-applications.ts`
- `frontend/src/models/interview-progress.ts`
- `frontend/src/api/interview-progress.ts`
- `frontend/mock/interview-progress.ts`

### JobPosting API

Routes:

```text
GET    /api/job-postings
GET    /api/job-postings/:id
POST   /api/job-postings
PUT    /api/job-postings/:id
DELETE /api/job-postings/:id
POST   /api/job-postings/parse
```

JSON fields:

```text
id, companyName, companyIntro, jobTitle, jobRequirements, base, salaryRange,
sourcePlatform, sourceUrl, status, remark, createdAt, updatedAt
```

Query fields:

```text
keyword, sourcePlatform, status, page, pageSize
```

Status values:

```text
pending, applied, not_fit, closed
```

### ResumeApplication API

Routes:

```text
GET    /api/resume-applications
GET    /api/resume-applications/:id
POST   /api/resume-applications
PUT    /api/resume-applications/:id
DELETE /api/resume-applications/:id
```

JSON fields:

```text
id, jobPostingId, resumeVersionId, repoId, companyName, companyIntro, jobTitle,
jobRequirements, base, salaryRange, sourcePlatform, sourceUrl, deliveryChannel,
appliedAt, currentStatus, interviewSummary, interviewNotes, resumeMatchScore,
interviewPerformanceScore, roleInterestScore, overallScore, remark, createdAt, updatedAt
```

Query fields:

```text
resumeVersionId, repoId, keyword, currentStatus, page, pageSize
```

Status values:

```text
applied, viewed, written_test, first_interview, second_interview,
final_interview, hr_interview, offer, rejected, closed
```

### InterviewProgress API

Routes:

```text
GET    /api/interview-progress?applicationId=...
POST   /api/interview-progress
PUT    /api/interview-progress/:id
DELETE /api/interview-progress/:id
```

JSON fields:

```text
id, applicationId, stage, occurredAt, result, interviewerOrTeam, note, createdAt, updatedAt
```

Stage values match application statuses. Result values:

```text
passed, pending, failed
```

Behavior:

- Creating, updating, or deleting interview progress recalculates the parent application's `currentStatus` from the latest progress by `occurredAt`.
- If the latest progress has `result === failed`, the application status becomes `rejected`.
- Otherwise the application status becomes the latest progress `stage`.

---

## File Map

Create these files:

```text
backend/internal/domain/job_posting/entity.go
backend/internal/domain/job_posting/repository.go
backend/internal/domain/job_posting/entity_test.go
backend/internal/domain/resume_application/entity.go
backend/internal/domain/resume_application/repository.go
backend/internal/domain/resume_application/entity_test.go
backend/internal/domain/interview_progress/entity.go
backend/internal/domain/interview_progress/repository.go
backend/internal/domain/interview_progress/entity_test.go
backend/internal/application/dto/job_dto.go
backend/internal/application/job_posting/job_posting_service.go
backend/internal/application/job_posting/job_posting_service_test.go
backend/internal/application/resume_application/resume_application_service.go
backend/internal/application/resume_application/resume_application_service_test.go
backend/internal/application/interview_progress/interview_progress_service.go
backend/internal/application/interview_progress/interview_progress_service_test.go
backend/internal/infrastructure/persistence/postgres/job_posting_repo.go
backend/internal/infrastructure/persistence/postgres/job_posting_repo_test.go
backend/internal/infrastructure/persistence/postgres/resume_application_repo.go
backend/internal/infrastructure/persistence/postgres/resume_application_repo_test.go
backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go
backend/internal/infrastructure/persistence/postgres/interview_progress_repo_test.go
backend/internal/interface/handler/job_posting_handler.go
backend/internal/interface/handler/job_posting_handler_test.go
backend/internal/interface/handler/resume_application_handler.go
backend/internal/interface/handler/resume_application_handler_test.go
backend/internal/interface/handler/interview_progress_handler.go
backend/internal/interface/handler/interview_progress_handler_test.go
backend/migrations/000004_create_job_postings.up.sql
backend/migrations/000004_create_job_postings.down.sql
backend/migrations/000005_create_resume_applications.up.sql
backend/migrations/000005_create_resume_applications.down.sql
backend/migrations/000006_create_interview_progress.up.sql
backend/migrations/000006_create_interview_progress.down.sql
```

Modify these files:

```text
backend/internal/interface/router/router.go
backend/internal/interface/router/router_test.go
backend/cmd/server/main.go
```

Do not modify frontend files.

---

### Task 1: Add frontend-aligned JobPosting domain

**Files:**
- Create: `backend/internal/domain/job_posting/entity.go`
- Create: `backend/internal/domain/job_posting/repository.go`
- Create: `backend/internal/domain/job_posting/entity_test.go`

- [ ] **Step 1: Write failing JobPosting entity tests**

Create `backend/internal/domain/job_posting/entity_test.go`:

```go
package job_posting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestJobPosting_Validate_AcceptsFrontendFields(t *testing.T) {
	posting := &domainjob.JobPosting{
		UserID:          "user-1",
		CompanyName:     "Acme",
		CompanyIntro:    "Company intro",
		JobTitle:        "Frontend Engineer",
		JobRequirements: "Vue and TypeScript",
		Base:            "Shenzhen",
		SalaryRange:     "25k-35k x 16",
		SourcePlatform:  "Boss",
		SourceURL:       "https://example.com/job/1",
		Status:          domainjob.StatusPending,
		Remark:          "priority",
	}

	require.NoError(t, posting.Validate())
}

func TestJobPosting_Validate_RejectsMissingRequiredFrontendFields(t *testing.T) {
	tests := []struct {
		name    string
		posting domainjob.JobPosting
	}{
		{name: "missing user", posting: domainjob.JobPosting{CompanyName: "Acme", JobTitle: "Frontend", Status: domainjob.StatusPending}},
		{name: "missing companyName", posting: domainjob.JobPosting{UserID: "user-1", JobTitle: "Frontend", Status: domainjob.StatusPending}},
		{name: "missing jobTitle", posting: domainjob.JobPosting{UserID: "user-1", CompanyName: "Acme", Status: domainjob.StatusPending}},
		{name: "invalid status", posting: domainjob.JobPosting{UserID: "user-1", CompanyName: "Acme", JobTitle: "Frontend", Status: "collecting"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.posting.Validate()
			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/job_posting -count=1
```

Expected: FAIL because package `internal/domain/job_posting` does not exist.

- [ ] **Step 3: Implement JobPosting entity**

Create `backend/internal/domain/job_posting/entity.go`:

```go
package job_posting

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusApplied Status = "applied"
	StatusNotFit  Status = "not_fit"
	StatusClosed  Status = "closed"
)

type JobPosting struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	UserID          string `gorm:"type:uuid;not null;index"`
	CompanyName     string `gorm:"type:varchar(120);not null"`
	CompanyIntro    string `gorm:"type:text"`
	JobTitle        string `gorm:"type:varchar(120);not null"`
	JobRequirements string `gorm:"type:text"`
	Base            string `gorm:"type:varchar(120)"`
	SalaryRange     string `gorm:"type:varchar(120)"`
	SourcePlatform  string `gorm:"type:varchar(120);index"`
	SourceURL       string `gorm:"type:varchar(500)"`
	Status          Status `gorm:"type:varchar(20);not null;index"`
	Remark          string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (p *JobPosting) Validate() error {
	if strings.TrimSpace(p.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(p.CompanyName) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(p.JobTitle) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(p.CompanyName)) > 120 || len([]rune(p.JobTitle)) > 120 {
		return domainerrors.ErrBadParams
	}
	if !validStatus(p.Status) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusPending, StatusApplied, StatusNotFit, StatusClosed:
		return true
	default:
		return false
	}
}
```

- [ ] **Step 4: Implement JobPosting repository interface**

Create `backend/internal/domain/job_posting/repository.go`:

```go
package job_posting

import "context"

type ListFilter struct {
	UserID         string
	Keyword        string
	SourcePlatform string
	Status         Status
	Page           int
	PageSize       int
}

type Repository interface {
	Create(ctx context.Context, p *JobPosting) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*JobPosting, error)
	Update(ctx context.Context, p *JobPosting) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]JobPosting, int64, error)
}
```

- [ ] **Step 5: Run JobPosting domain tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/job_posting -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/domain/job_posting && git commit -m "feat: add frontend-aligned job posting domain"
```

---

### Task 2: Add frontend-aligned ResumeApplication domain

**Files:**
- Create: `backend/internal/domain/resume_application/entity.go`
- Create: `backend/internal/domain/resume_application/repository.go`
- Create: `backend/internal/domain/resume_application/entity_test.go`

- [ ] **Step 1: Write failing ResumeApplication entity tests**

Create `backend/internal/domain/resume_application/entity_test.go`:

```go
package resume_application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestResumeApplication_Validate_AcceptsFrontendFields(t *testing.T) {
	app := &domainapp.ResumeApplication{
		UserID:                    "user-1",
		JobPostingID:              "job-1",
		ResumeVersionID:           "version-1",
		RepoID:                    "repo-1",
		CompanyName:               "Acme",
		JobTitle:                  "Frontend Engineer",
		AppliedAt:                 time.Now(),
		CurrentStatus:             domainapp.StatusApplied,
		ResumeMatchScore:          5,
		InterviewPerformanceScore: 4,
		RoleInterestScore:         5,
		OverallScore:              9,
	}

	require.NoError(t, app.Validate())
}

func TestResumeApplication_Validate_RejectsInvalidFrontendFields(t *testing.T) {
	tests := []struct {
		name string
		app  domainapp.ResumeApplication
	}{
		{name: "missing user", app: domainapp.ResumeApplication{ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "missing version", app: domainapp.ResumeApplication{UserID: "u1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "missing repo", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "invalid status", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: "pending"}},
		{name: "invalid score", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied, OverallScore: 11}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.app.Validate(), domainerrors.ErrBadParams)
		})
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_application -count=1
```

Expected: FAIL because package `internal/domain/resume_application` does not exist.

- [ ] **Step 3: Implement ResumeApplication entity**

Create `backend/internal/domain/resume_application/entity.go`:

```go
package resume_application

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Status string

const (
	StatusApplied         Status = "applied"
	StatusViewed          Status = "viewed"
	StatusWrittenTest     Status = "written_test"
	StatusFirstInterview  Status = "first_interview"
	StatusSecondInterview Status = "second_interview"
	StatusFinalInterview  Status = "final_interview"
	StatusHRInterview     Status = "hr_interview"
	StatusOffer           Status = "offer"
	StatusRejected        Status = "rejected"
	StatusClosed          Status = "closed"
)

type ResumeApplication struct {
	ID                        string `gorm:"type:uuid;primaryKey"`
	UserID                    string `gorm:"type:uuid;not null;index"`
	JobPostingID              string `gorm:"type:uuid;index"`
	ResumeVersionID           string `gorm:"type:uuid;not null;index"`
	RepoID                    string `gorm:"type:uuid;not null;index"`
	CompanyName               string `gorm:"type:varchar(120);not null"`
	CompanyIntro              string `gorm:"type:text"`
	JobTitle                  string `gorm:"type:varchar(120);not null"`
	JobRequirements           string `gorm:"type:text"`
	Base                      string `gorm:"type:varchar(120)"`
	SalaryRange               string `gorm:"type:varchar(120)"`
	SourcePlatform            string `gorm:"type:varchar(120)"`
	SourceURL                 string `gorm:"type:varchar(500)"`
	DeliveryChannel           string `gorm:"type:varchar(120)"`
	AppliedAt                 time.Time `gorm:"not null;index"`
	CurrentStatus             Status    `gorm:"type:varchar(30);not null;index"`
	InterviewSummary          string    `gorm:"type:text"`
	InterviewNotes            string    `gorm:"type:text"`
	ResumeMatchScore          int
	InterviewPerformanceScore int
	RoleInterestScore         int
	OverallScore              int
	Remark                    string `gorm:"type:text"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 gorm.DeletedAt `gorm:"index"`
}

func (a *ResumeApplication) Validate() error {
	if strings.TrimSpace(a.UserID) == "" || strings.TrimSpace(a.ResumeVersionID) == "" || strings.TrimSpace(a.RepoID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(a.CompanyName) == "" || strings.TrimSpace(a.JobTitle) == "" {
		return domainerrors.ErrBadParams
	}
	if a.AppliedAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if !validStatus(a.CurrentStatus) {
		return domainerrors.ErrBadParams
	}
	if !validOptionalScore(a.ResumeMatchScore, 1, 5) || !validOptionalScore(a.InterviewPerformanceScore, 1, 5) || !validOptionalScore(a.RoleInterestScore, 1, 5) || !validOptionalScore(a.OverallScore, 1, 10) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusApplied, StatusViewed, StatusWrittenTest, StatusFirstInterview, StatusSecondInterview, StatusFinalInterview, StatusHRInterview, StatusOffer, StatusRejected, StatusClosed:
		return true
	default:
		return false
	}
}

func validOptionalScore(score, min, max int) bool {
	return score == 0 || (score >= min && score <= max)
}
```

- [ ] **Step 4: Implement ResumeApplication repository interface**

Create `backend/internal/domain/resume_application/repository.go`:

```go
package resume_application

import "context"

type ListFilter struct {
	UserID          string
	ResumeVersionID string
	RepoID          string
	Keyword         string
	CurrentStatus   Status
	Page            int
	PageSize        int
}

type Repository interface {
	Create(ctx context.Context, a *ResumeApplication) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*ResumeApplication, error)
	Update(ctx context.Context, a *ResumeApplication) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]ResumeApplication, int64, error)
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
	UpdateCurrentStatus(ctx context.Context, id, userID string, status Status) error
}
```

- [ ] **Step 5: Run ResumeApplication domain tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/resume_application -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/domain/resume_application && git commit -m "feat: add frontend-aligned resume application domain"
```

---

### Task 3: Add frontend-aligned InterviewProgress domain

**Files:**
- Create: `backend/internal/domain/interview_progress/entity.go`
- Create: `backend/internal/domain/interview_progress/repository.go`
- Create: `backend/internal/domain/interview_progress/entity_test.go`

- [ ] **Step 1: Write failing InterviewProgress entity tests**

Create `backend/internal/domain/interview_progress/entity_test.go`:

```go
package interview_progress_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestInterviewProgress_Validate_AcceptsFrontendFields(t *testing.T) {
	progress := &domainprogress.InterviewProgress{
		UserID:            "user-1",
		ApplicationID:     "application-1",
		Stage:             domainprogress.StageFirstInterview,
		OccurredAt:        time.Now(),
		Result:            domainprogress.ResultPassed,
		InterviewerOrTeam: "Frontend team",
		Note:              "Asked Vue questions",
	}

	require.NoError(t, progress.Validate())
}

func TestInterviewProgress_Validate_RejectsInvalidFrontendFields(t *testing.T) {
	tests := []struct {
		name     string
		progress domainprogress.InterviewProgress
	}{
		{name: "missing user", progress: domainprogress.InterviewProgress{ApplicationID: "a1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "missing application", progress: domainprogress.InterviewProgress{UserID: "u1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "missing occurredAt", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: domainprogress.StageApplied, Result: domainprogress.ResultPassed}},
		{name: "invalid stage", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: "phone", OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "invalid result", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: "cancelled"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.progress.Validate(), domainerrors.ErrBadParams)
		})
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/interview_progress -count=1
```

Expected: FAIL because package `internal/domain/interview_progress` does not exist.

- [ ] **Step 3: Implement InterviewProgress entity**

Create `backend/internal/domain/interview_progress/entity.go`:

```go
package interview_progress

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Stage string

type Result string

const (
	StageApplied         Stage = "applied"
	StageViewed          Stage = "viewed"
	StageWrittenTest     Stage = "written_test"
	StageFirstInterview  Stage = "first_interview"
	StageSecondInterview Stage = "second_interview"
	StageFinalInterview  Stage = "final_interview"
	StageHRInterview     Stage = "hr_interview"
	StageOffer           Stage = "offer"
	StageRejected        Stage = "rejected"
	StageClosed          Stage = "closed"

	ResultPassed  Result = "passed"
	ResultPending Result = "pending"
	ResultFailed  Result = "failed"
)

type InterviewProgress struct {
	ID                string `gorm:"type:uuid;primaryKey"`
	UserID            string `gorm:"type:uuid;not null;index"`
	ApplicationID     string `gorm:"type:uuid;not null;index"`
	Stage             Stage  `gorm:"type:varchar(30);not null"`
	OccurredAt        time.Time `gorm:"not null;index"`
	Result            Result    `gorm:"type:varchar(20);not null"`
	InterviewerOrTeam string    `gorm:"type:varchar(120)"`
	Note              string    `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (p *InterviewProgress) Validate() error {
	if strings.TrimSpace(p.UserID) == "" || strings.TrimSpace(p.ApplicationID) == "" {
		return domainerrors.ErrBadParams
	}
	if p.OccurredAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if !validStage(p.Stage) || !validResult(p.Result) {
		return domainerrors.ErrBadParams
	}
	return nil
}

func validStage(stage Stage) bool {
	switch stage {
	case StageApplied, StageViewed, StageWrittenTest, StageFirstInterview, StageSecondInterview, StageFinalInterview, StageHRInterview, StageOffer, StageRejected, StageClosed:
		return true
	default:
		return false
	}
}

func validResult(result Result) bool {
	switch result {
	case ResultPassed, ResultPending, ResultFailed:
		return true
	default:
		return false
	}
}
```

- [ ] **Step 4: Implement InterviewProgress repository interface**

Create `backend/internal/domain/interview_progress/repository.go`:

```go
package interview_progress

import "context"

type Repository interface {
	Create(ctx context.Context, p *InterviewProgress) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*InterviewProgress, error)
	Update(ctx context.Context, p *InterviewProgress) error
	Delete(ctx context.Context, id, userID string) error
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]InterviewProgress, error)
	DeleteByApplicationID(ctx context.Context, applicationID, userID string) error
}
```

- [ ] **Step 5: Run InterviewProgress domain tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/interview_progress -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/domain/interview_progress && git commit -m "feat: add frontend-aligned interview progress domain"
```

---

### Task 4: Add job-domain migrations with frontend-aligned columns

**Files:**
- Create: `backend/migrations/000004_create_job_postings.up.sql`
- Create: `backend/migrations/000004_create_job_postings.down.sql`
- Create: `backend/migrations/000005_create_resume_applications.up.sql`
- Create: `backend/migrations/000005_create_resume_applications.down.sql`
- Create: `backend/migrations/000006_create_interview_progress.up.sql`
- Create: `backend/migrations/000006_create_interview_progress.down.sql`

- [ ] **Step 1: Create job posting migration**

Create `backend/migrations/000004_create_job_postings.up.sql`:

```sql
CREATE TABLE job_postings (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID         NOT NULL REFERENCES users(id),
    company_name     VARCHAR(120) NOT NULL,
    company_intro    TEXT,
    job_title        VARCHAR(120) NOT NULL,
    job_requirements TEXT,
    base             VARCHAR(120),
    salary_range     VARCHAR(120),
    source_platform  VARCHAR(120),
    source_url       VARCHAR(500),
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending',
    remark           TEXT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_job_postings_user_id ON job_postings(user_id);
CREATE INDEX idx_job_postings_user_status ON job_postings(user_id, status);
CREATE INDEX idx_job_postings_user_source_platform ON job_postings(user_id, source_platform);
```

Create `backend/migrations/000004_create_job_postings.down.sql`:

```sql
DROP TABLE IF EXISTS job_postings;
```

- [ ] **Step 2: Create resume application migration**

Create `backend/migrations/000005_create_resume_applications.up.sql`:

```sql
CREATE TABLE resume_applications (
    id                          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                     UUID         NOT NULL REFERENCES users(id),
    job_posting_id              UUID         REFERENCES job_postings(id),
    resume_version_id           UUID         NOT NULL REFERENCES resume_versions(id),
    repo_id                     UUID         NOT NULL REFERENCES resume_repos(id),
    company_name                VARCHAR(120) NOT NULL,
    company_intro               TEXT,
    job_title                   VARCHAR(120) NOT NULL,
    job_requirements            TEXT,
    base                        VARCHAR(120),
    salary_range                VARCHAR(120),
    source_platform             VARCHAR(120),
    source_url                  VARCHAR(500),
    delivery_channel            VARCHAR(120),
    applied_at                  TIMESTAMPTZ  NOT NULL,
    current_status              VARCHAR(30)  NOT NULL DEFAULT 'applied',
    interview_summary           TEXT,
    interview_notes             TEXT,
    resume_match_score          INT          NOT NULL DEFAULT 0,
    interview_performance_score INT          NOT NULL DEFAULT 0,
    role_interest_score         INT          NOT NULL DEFAULT 0,
    overall_score               INT          NOT NULL DEFAULT 0,
    remark                      TEXT,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ
);

CREATE INDEX idx_resume_applications_user_id ON resume_applications(user_id);
CREATE INDEX idx_resume_applications_resume_version_id ON resume_applications(resume_version_id);
CREATE INDEX idx_resume_applications_repo_id ON resume_applications(repo_id);
CREATE INDEX idx_resume_applications_current_status ON resume_applications(current_status);
CREATE INDEX idx_resume_applications_applied_at ON resume_applications(applied_at);
```

Create `backend/migrations/000005_create_resume_applications.down.sql`:

```sql
DROP TABLE IF EXISTS resume_applications;
```

- [ ] **Step 3: Create interview progress migration**

Create `backend/migrations/000006_create_interview_progress.up.sql`:

```sql
CREATE TABLE interview_progress (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID        NOT NULL REFERENCES users(id),
    application_id      UUID        NOT NULL REFERENCES resume_applications(id),
    stage               VARCHAR(30) NOT NULL,
    occurred_at         TIMESTAMPTZ NOT NULL,
    result              VARCHAR(20) NOT NULL,
    interviewer_or_team VARCHAR(120),
    note                TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_interview_progress_user_id ON interview_progress(user_id);
CREATE INDEX idx_interview_progress_application_id ON interview_progress(application_id);
CREATE INDEX idx_interview_progress_occurred_at ON interview_progress(occurred_at);
```

Create `backend/migrations/000006_create_interview_progress.down.sql`:

```sql
DROP TABLE IF EXISTS interview_progress;
```

- [ ] **Step 4: Run migration-related tests or compile check**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/job_posting ./internal/domain/resume_application ./internal/domain/interview_progress -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add migrations/000004_* migrations/000005_* migrations/000006_* && git commit -m "feat: add frontend-aligned job domain migrations"
```

---

### Task 5: Add application DTOs matching frontend JSON

**Files:**
- Create: `backend/internal/application/dto/job_dto.go`

- [ ] **Step 1: Create DTO file**

Create `backend/internal/application/dto/job_dto.go`:

```go
package dto

import (
	"time"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type CreateJobPostingRequest struct {
	CompanyName     string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro    string           `json:"companyIntro"`
	JobTitle        string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements string           `json:"jobRequirements"`
	Base            string           `json:"base"`
	SalaryRange     string           `json:"salaryRange"`
	SourcePlatform  string           `json:"sourcePlatform"`
	SourceURL       string           `json:"sourceUrl"`
	Status          domainjob.Status `json:"status" binding:"required"`
	Remark          string           `json:"remark"`
}

type UpdateJobPostingRequest = CreateJobPostingRequest

type ParseJobPostingRequest struct {
	URL string `json:"url" binding:"required"`
}

type JobPostingResponse struct {
	ID              string           `json:"id"`
	CompanyName     string           `json:"companyName"`
	CompanyIntro    string           `json:"companyIntro"`
	JobTitle        string           `json:"jobTitle"`
	JobRequirements string           `json:"jobRequirements"`
	Base            string           `json:"base"`
	SalaryRange     string           `json:"salaryRange"`
	SourcePlatform  string           `json:"sourcePlatform"`
	SourceURL       string           `json:"sourceUrl"`
	Status          domainjob.Status `json:"status"`
	Remark          string           `json:"remark"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type CreateResumeApplicationRequest struct {
	JobPostingID              string           `json:"jobPostingId"`
	ResumeVersionID           string           `json:"resumeVersionId" binding:"required"`
	RepoID                    string           `json:"repoId" binding:"required"`
	CompanyName               string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro              string           `json:"companyIntro"`
	JobTitle                  string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements           string           `json:"jobRequirements"`
	Base                      string           `json:"base"`
	SalaryRange               string           `json:"salaryRange"`
	SourcePlatform            string           `json:"sourcePlatform"`
	SourceURL                 string           `json:"sourceUrl"`
	DeliveryChannel           string           `json:"deliveryChannel"`
	AppliedAt                 time.Time        `json:"appliedAt" binding:"required"`
	CurrentStatus             domainapp.Status `json:"currentStatus" binding:"required"`
	InterviewSummary          string           `json:"interviewSummary"`
	InterviewNotes            string           `json:"interviewNotes"`
	ResumeMatchScore          int              `json:"resumeMatchScore"`
	InterviewPerformanceScore int              `json:"interviewPerformanceScore"`
	RoleInterestScore         int              `json:"roleInterestScore"`
	OverallScore              int              `json:"overallScore"`
	Remark                    string           `json:"remark"`
}

type UpdateResumeApplicationRequest = CreateResumeApplicationRequest

type ResumeApplicationResponse struct {
	ID                        string           `json:"id"`
	JobPostingID              string           `json:"jobPostingId"`
	ResumeVersionID           string           `json:"resumeVersionId"`
	RepoID                    string           `json:"repoId"`
	CompanyName               string           `json:"companyName"`
	CompanyIntro              string           `json:"companyIntro"`
	JobTitle                  string           `json:"jobTitle"`
	JobRequirements           string           `json:"jobRequirements"`
	Base                      string           `json:"base"`
	SalaryRange               string           `json:"salaryRange"`
	SourcePlatform            string           `json:"sourcePlatform"`
	SourceURL                 string           `json:"sourceUrl"`
	DeliveryChannel           string           `json:"deliveryChannel"`
	AppliedAt                 time.Time        `json:"appliedAt"`
	CurrentStatus             domainapp.Status `json:"currentStatus"`
	InterviewSummary          string           `json:"interviewSummary"`
	InterviewNotes            string           `json:"interviewNotes"`
	ResumeMatchScore          int              `json:"resumeMatchScore"`
	InterviewPerformanceScore int              `json:"interviewPerformanceScore"`
	RoleInterestScore         int              `json:"roleInterestScore"`
	OverallScore              int              `json:"overallScore"`
	Remark                    string           `json:"remark"`
	CreatedAt                 time.Time        `json:"createdAt"`
	UpdatedAt                 time.Time        `json:"updatedAt"`
}

type CreateInterviewProgressRequest struct {
	ApplicationID     string                `json:"applicationId" binding:"required"`
	Stage             domainprogress.Stage  `json:"stage" binding:"required"`
	OccurredAt        time.Time             `json:"occurredAt" binding:"required"`
	Result            domainprogress.Result `json:"result" binding:"required"`
	InterviewerOrTeam string                `json:"interviewerOrTeam"`
	Note              string                `json:"note"`
}

type UpdateInterviewProgressRequest = CreateInterviewProgressRequest

type InterviewProgressResponse struct {
	ID                string                `json:"id"`
	ApplicationID     string                `json:"applicationId"`
	Stage             domainprogress.Stage  `json:"stage"`
	OccurredAt        time.Time             `json:"occurredAt"`
	Result            domainprogress.Result `json:"result"`
	InterviewerOrTeam string                `json:"interviewerOrTeam"`
	Note              string                `json:"note"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

func ToJobPostingResponse(p *domainjob.JobPosting) *JobPostingResponse {
	if p == nil {
		return nil
	}
	return &JobPostingResponse{ID: p.ID, CompanyName: p.CompanyName, CompanyIntro: p.CompanyIntro, JobTitle: p.JobTitle, JobRequirements: p.JobRequirements, Base: p.Base, SalaryRange: p.SalaryRange, SourcePlatform: p.SourcePlatform, SourceURL: p.SourceURL, Status: p.Status, Remark: p.Remark, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func ToJobPostingResponses(rows []domainjob.JobPosting) []JobPostingResponse {
	out := make([]JobPostingResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *ToJobPostingResponse(&rows[i]))
	}
	return out
}

func ToResumeApplicationResponse(a *domainapp.ResumeApplication) *ResumeApplicationResponse {
	if a == nil {
		return nil
	}
	return &ResumeApplicationResponse{ID: a.ID, JobPostingID: a.JobPostingID, ResumeVersionID: a.ResumeVersionID, RepoID: a.RepoID, CompanyName: a.CompanyName, CompanyIntro: a.CompanyIntro, JobTitle: a.JobTitle, JobRequirements: a.JobRequirements, Base: a.Base, SalaryRange: a.SalaryRange, SourcePlatform: a.SourcePlatform, SourceURL: a.SourceURL, DeliveryChannel: a.DeliveryChannel, AppliedAt: a.AppliedAt, CurrentStatus: a.CurrentStatus, InterviewSummary: a.InterviewSummary, InterviewNotes: a.InterviewNotes, ResumeMatchScore: a.ResumeMatchScore, InterviewPerformanceScore: a.InterviewPerformanceScore, RoleInterestScore: a.RoleInterestScore, OverallScore: a.OverallScore, Remark: a.Remark, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}

func ToResumeApplicationResponses(rows []domainapp.ResumeApplication) []ResumeApplicationResponse {
	out := make([]ResumeApplicationResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *ToResumeApplicationResponse(&rows[i]))
	}
	return out
}

func ToInterviewProgressResponse(p *domainprogress.InterviewProgress) *InterviewProgressResponse {
	if p == nil {
		return nil
	}
	return &InterviewProgressResponse{ID: p.ID, ApplicationID: p.ApplicationID, Stage: p.Stage, OccurredAt: p.OccurredAt, Result: p.Result, InterviewerOrTeam: p.InterviewerOrTeam, Note: p.Note, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func ToInterviewProgressResponses(rows []domainprogress.InterviewProgress) []InterviewProgressResponse {
	out := make([]InterviewProgressResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *ToInterviewProgressResponse(&rows[i]))
	}
	return out
}
```

- [ ] **Step 2: Run compile check**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/dto -run '^$' -count=1
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/application/dto/job_dto.go && git commit -m "feat: add frontend-aligned job domain DTOs"
```

---

### Task 6: Implement JobPosting service

**Files:**
- Create: `backend/internal/application/job_posting/job_posting_service.go`
- Create: `backend/internal/application/job_posting/job_posting_service_test.go`

- [ ] **Step 1: Write failing service tests**

Create `backend/internal/application/job_posting/job_posting_service_test.go`:

```go
package job_posting_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/job_posting"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

type mockJobPostingRepo struct{ mock.Mock }

func (m *mockJobPostingRepo) Create(ctx context.Context, p *domainjob.JobPosting) error { return m.Called(ctx, p).Error(0) }
func (m *mockJobPostingRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) { args := m.Called(ctx, id, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).(*domainjob.JobPosting), args.Error(1) }
func (m *mockJobPostingRepo) Update(ctx context.Context, p *domainjob.JobPosting) error { return m.Called(ctx, p).Error(0) }
func (m *mockJobPostingRepo) Delete(ctx context.Context, id, userID string) error { return m.Called(ctx, id, userID).Error(0) }
func (m *mockJobPostingRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) { args := m.Called(ctx, f); return args.Get(0).([]domainjob.JobPosting), args.Get(1).(int64), args.Error(2) }

func TestJobPostingService_Create_CopiesFrontendFields(t *testing.T) {
	repo := &mockJobPostingRepo{}
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domainjob.JobPosting) bool {
		return p.ID != "" && p.UserID == "u1" && p.CompanyName == "Acme" && p.JobTitle == "Frontend" && p.SourcePlatform == "Boss" && p.SourceURL == "https://example.com" && p.Status == domainjob.StatusPending && p.Remark == "priority"
	})).Return(nil)
	svc := appsvc.NewJobPostingService(repo)

	err := svc.Create(context.Background(), "u1", dto.CreateJobPostingRequest{CompanyName: "Acme", JobTitle: "Frontend", SourcePlatform: "Boss", SourceURL: "https://example.com", Status: domainjob.StatusPending, Remark: "priority"})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestJobPostingService_List_UsesFrontendFilters(t *testing.T) {
	repo := &mockJobPostingRepo{}
	repo.On("ListByUser", mock.Anything, domainjob.ListFilter{UserID: "u1", Keyword: "acme", SourcePlatform: "Boss", Status: domainjob.StatusApplied, Page: 2, PageSize: 9}).Return([]domainjob.JobPosting{}, int64(0), nil)
	svc := appsvc.NewJobPostingService(repo)

	_, _, err := svc.List(context.Background(), "u1", "acme", "Boss", "applied", 2, 9)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
```

- [ ] **Step 2: Run service tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/job_posting -count=1
```

Expected: FAIL because service package does not exist.

- [ ] **Step 3: Implement JobPosting service**

Create `backend/internal/application/job_posting/job_posting_service.go`:

```go
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
	p := &domainjob.JobPosting{ID: uuid.NewString(), UserID: userID, CompanyName: req.CompanyName, CompanyIntro: req.CompanyIntro, JobTitle: req.JobTitle, JobRequirements: req.JobRequirements, Base: req.Base, SalaryRange: req.SalaryRange, SourcePlatform: req.SourcePlatform, SourceURL: req.SourceURL, Status: req.Status, Remark: req.Remark}
	if p.Status == "" { p.Status = domainjob.StatusPending }
	if err := p.Validate(); err != nil { return err }
	return s.repo.Create(ctx, p)
}

func (s *JobPostingService) List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
	return s.repo.ListByUser(ctx, domainjob.ListFilter{UserID: userID, Keyword: keyword, SourcePlatform: sourcePlatform, Status: domainjob.Status(status), Page: page, PageSize: pageSize})
}

func (s *JobPostingService) Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *JobPostingService) Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error {
	p, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	p.CompanyName, p.CompanyIntro, p.JobTitle, p.JobRequirements = req.CompanyName, req.CompanyIntro, req.JobTitle, req.JobRequirements
	p.Base, p.SalaryRange, p.SourcePlatform, p.SourceURL = req.Base, req.SalaryRange, req.SourcePlatform, req.SourceURL
	p.Status, p.Remark = req.Status, req.Remark
	if err := p.Validate(); err != nil { return err }
	return s.repo.Update(ctx, p)
}

func (s *JobPostingService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *JobPostingService) Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error) {
	platform := "其他"
	if contains(req.URL, "zhipin") { platform = "Boss" }
	return &dto.JobPostingResponse{CompanyName: "示例公司", JobTitle: "待确认岗位", JobRequirements: "通过解析接口返回的岗位要求示例，用户可继续修改", Base: "北京", SalaryRange: "25k-35k x 16", SourcePlatform: platform, SourceURL: req.URL, Status: domainjob.StatusPending}, nil
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && index(s, substr) >= 0)
}

func index(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr { return i }
	}
	return -1
}
```

- [ ] **Step 4: Run service tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/job_posting -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/application/job_posting && git commit -m "feat: add job posting service"
```

---

### Task 7: Implement ResumeApplication service

**Files:**
- Create: `backend/internal/application/resume_application/resume_application_service.go`
- Create: `backend/internal/application/resume_application/resume_application_service_test.go`

- [ ] **Step 1: Write failing service tests**

Create `backend/internal/application/resume_application/resume_application_service_test.go`:

```go
package resume_application_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_application"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type mockApplicationRepo struct{ mock.Mock }

func (m *mockApplicationRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error { return m.Called(ctx, a).Error(0) }
func (m *mockApplicationRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error) { args := m.Called(ctx, id, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).(*domainapp.ResumeApplication), args.Error(1) }
func (m *mockApplicationRepo) Update(ctx context.Context, a *domainapp.ResumeApplication) error { return m.Called(ctx, a).Error(0) }
func (m *mockApplicationRepo) Delete(ctx context.Context, id, userID string) error { return m.Called(ctx, id, userID).Error(0) }
func (m *mockApplicationRepo) ListByUser(ctx context.Context, f domainapp.ListFilter) ([]domainapp.ResumeApplication, int64, error) { args := m.Called(ctx, f); return args.Get(0).([]domainapp.ResumeApplication), args.Get(1).(int64), args.Error(2) }
func (m *mockApplicationRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) { args := m.Called(ctx, versionID); return args.Bool(0), args.Error(1) }
func (m *mockApplicationRepo) UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error { return m.Called(ctx, id, userID, status).Error(0) }

func TestResumeApplicationService_Create_CopiesFrontendFields(t *testing.T) {
	repo := &mockApplicationRepo{}
	appliedAt := time.Now()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
		return a.ID != "" && a.UserID == "u1" && a.JobPostingID == "j1" && a.ResumeVersionID == "v1" && a.RepoID == "r1" && a.CompanyName == "Acme" && a.CurrentStatus == domainapp.StatusApplied && a.OverallScore == 8
	})).Return(nil)
	svc := appsvc.NewResumeApplicationService(repo)

	err := svc.Create(context.Background(), "u1", dto.CreateResumeApplicationRequest{JobPostingID: "j1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: appliedAt, CurrentStatus: domainapp.StatusApplied, OverallScore: 8})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
```

- [ ] **Step 2: Run service tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_application -count=1
```

Expected: FAIL because service package does not exist.

- [ ] **Step 3: Implement ResumeApplication service**

Create `backend/internal/application/resume_application/resume_application_service.go`:

```go
package resume_application

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type ResumeApplicationService struct{ repo domainapp.Repository }

func NewResumeApplicationService(repo domainapp.Repository) *ResumeApplicationService { return &ResumeApplicationService{repo: repo} }

func (s *ResumeApplicationService) Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) error {
	a := &domainapp.ResumeApplication{ID: uuid.NewString(), UserID: userID, JobPostingID: req.JobPostingID, ResumeVersionID: req.ResumeVersionID, RepoID: req.RepoID, CompanyName: req.CompanyName, CompanyIntro: req.CompanyIntro, JobTitle: req.JobTitle, JobRequirements: req.JobRequirements, Base: req.Base, SalaryRange: req.SalaryRange, SourcePlatform: req.SourcePlatform, SourceURL: req.SourceURL, DeliveryChannel: req.DeliveryChannel, AppliedAt: req.AppliedAt, CurrentStatus: req.CurrentStatus, InterviewSummary: req.InterviewSummary, InterviewNotes: req.InterviewNotes, ResumeMatchScore: req.ResumeMatchScore, InterviewPerformanceScore: req.InterviewPerformanceScore, RoleInterestScore: req.RoleInterestScore, OverallScore: req.OverallScore, Remark: req.Remark}
	if a.CurrentStatus == "" { a.CurrentStatus = domainapp.StatusApplied }
	if err := a.Validate(); err != nil { return err }
	return s.repo.Create(ctx, a)
}

func (s *ResumeApplicationService) List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
	return s.repo.ListByUser(ctx, domainapp.ListFilter{UserID: userID, ResumeVersionID: resumeVersionID, RepoID: repoID, Keyword: keyword, CurrentStatus: domainapp.Status(currentStatus), Page: page, PageSize: pageSize})
}

func (s *ResumeApplicationService) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *ResumeApplicationService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error {
	a, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	a.JobPostingID, a.ResumeVersionID, a.RepoID = req.JobPostingID, req.ResumeVersionID, req.RepoID
	a.CompanyName, a.CompanyIntro, a.JobTitle, a.JobRequirements = req.CompanyName, req.CompanyIntro, req.JobTitle, req.JobRequirements
	a.Base, a.SalaryRange, a.SourcePlatform, a.SourceURL, a.DeliveryChannel = req.Base, req.SalaryRange, req.SourcePlatform, req.SourceURL, req.DeliveryChannel
	a.AppliedAt, a.CurrentStatus = req.AppliedAt, req.CurrentStatus
	a.InterviewSummary, a.InterviewNotes = req.InterviewSummary, req.InterviewNotes
	a.ResumeMatchScore, a.InterviewPerformanceScore, a.RoleInterestScore, a.OverallScore = req.ResumeMatchScore, req.InterviewPerformanceScore, req.RoleInterestScore, req.OverallScore
	a.Remark = req.Remark
	if err := a.Validate(); err != nil { return err }
	return s.repo.Update(ctx, a)
}

func (s *ResumeApplicationService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *ResumeApplicationService) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	return s.repo.ExistsByResumeVersionID(ctx, versionID)
}

func (s *ResumeApplicationService) UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error {
	return s.repo.UpdateCurrentStatus(ctx, id, userID, status)
}
```

- [ ] **Step 4: Run service tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/resume_application -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/application/resume_application && git commit -m "feat: add resume application service"
```

---

### Task 8: Implement InterviewProgress service with application status sync

**Files:**
- Create: `backend/internal/application/interview_progress/interview_progress_service.go`
- Create: `backend/internal/application/interview_progress/interview_progress_service_test.go`

- [ ] **Step 1: Write failing service tests**

Create `backend/internal/application/interview_progress/interview_progress_service_test.go`:

```go
package interview_progress_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/interview_progress"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type mockProgressRepo struct{ mock.Mock }
func (m *mockProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error { return m.Called(ctx, p).Error(0) }
func (m *mockProgressRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) { args := m.Called(ctx, id, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).(*domainprogress.InterviewProgress), args.Error(1) }
func (m *mockProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error { return m.Called(ctx, p).Error(0) }
func (m *mockProgressRepo) Delete(ctx context.Context, id, userID string) error { return m.Called(ctx, id, userID).Error(0) }
func (m *mockProgressRepo) ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error) { args := m.Called(ctx, applicationID, userID); return args.Get(0).([]domainprogress.InterviewProgress), args.Error(1) }
func (m *mockProgressRepo) DeleteByApplicationID(ctx context.Context, applicationID, userID string) error { return m.Called(ctx, applicationID, userID).Error(0) }

type mockStatusUpdater struct{ mock.Mock }
func (m *mockStatusUpdater) UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error { return m.Called(ctx, userID, id, status).Error(0) }

func TestInterviewProgressService_Create_SyncsFailedToRejected(t *testing.T) {
	repo := &mockProgressRepo{}
	apps := &mockStatusUpdater{}
	now := time.Now()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*interview_progress.InterviewProgress")).Return(nil)
	repo.On("ListByApplicationAndUser", mock.Anything, "a1", "u1").Return([]domainprogress.InterviewProgress{{ApplicationID: "a1", Stage: domainprogress.StageFirstInterview, Result: domainprogress.ResultFailed, OccurredAt: now}}, nil)
	apps.On("UpdateCurrentStatus", mock.Anything, "u1", "a1", domainapp.StatusRejected).Return(nil)
	svc := appsvc.NewInterviewProgressService(repo, apps)

	err := svc.Create(context.Background(), "u1", dto.CreateInterviewProgressRequest{ApplicationID: "a1", Stage: domainprogress.StageFirstInterview, Result: domainprogress.ResultFailed, OccurredAt: now})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}
```

- [ ] **Step 2: Run service tests to verify failure**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/interview_progress -count=1
```

Expected: FAIL because service package does not exist.

- [ ] **Step 3: Implement InterviewProgress service**

Create `backend/internal/application/interview_progress/interview_progress_service.go`:

```go
package interview_progress

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type ApplicationStatusUpdater interface {
	UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error
}

type InterviewProgressService struct {
	repo domainprogress.Repository
	apps ApplicationStatusUpdater
}

func NewInterviewProgressService(repo domainprogress.Repository, apps ApplicationStatusUpdater) *InterviewProgressService {
	return &InterviewProgressService{repo: repo, apps: apps}
}

func (s *InterviewProgressService) Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) error {
	p := &domainprogress.InterviewProgress{ID: uuid.NewString(), UserID: userID, ApplicationID: req.ApplicationID, Stage: req.Stage, OccurredAt: req.OccurredAt, Result: req.Result, InterviewerOrTeam: req.InterviewerOrTeam, Note: req.Note}
	if err := p.Validate(); err != nil { return err }
	if err := s.repo.Create(ctx, p); err != nil { return err }
	return s.syncApplicationStatus(ctx, userID, req.ApplicationID)
}

func (s *InterviewProgressService) List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error) {
	return s.repo.ListByApplicationAndUser(ctx, applicationID, userID)
}

func (s *InterviewProgressService) Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) error {
	p, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	p.ApplicationID, p.Stage, p.OccurredAt, p.Result = req.ApplicationID, req.Stage, req.OccurredAt, req.Result
	p.InterviewerOrTeam, p.Note = req.InterviewerOrTeam, req.Note
	if err := p.Validate(); err != nil { return err }
	if err := s.repo.Update(ctx, p); err != nil { return err }
	return s.syncApplicationStatus(ctx, userID, p.ApplicationID)
}

func (s *InterviewProgressService) Delete(ctx context.Context, userID, id string) error {
	p, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return err }
	if err := s.repo.Delete(ctx, id, userID); err != nil { return err }
	return s.syncApplicationStatus(ctx, userID, p.ApplicationID)
}

func (s *InterviewProgressService) syncApplicationStatus(ctx context.Context, userID, applicationID string) error {
	rows, err := s.repo.ListByApplicationAndUser(ctx, applicationID, userID)
	if err != nil { return err }
	if len(rows) == 0 { return nil }
	sort.Slice(rows, func(i, j int) bool { return rows[i].OccurredAt.Before(rows[j].OccurredAt) })
	latest := rows[len(rows)-1]
	if latest.Result == domainprogress.ResultFailed {
		return s.apps.UpdateCurrentStatus(ctx, userID, applicationID, domainapp.StatusRejected)
	}
	return s.apps.UpdateCurrentStatus(ctx, userID, applicationID, domainapp.Status(latest.Stage))
}
```

- [ ] **Step 4: Run service tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/interview_progress -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/application/interview_progress && git commit -m "feat: add interview progress service"
```

---

### Task 9: Implement PostgreSQL repositories

**Files:**
- Create: `backend/internal/infrastructure/persistence/postgres/job_posting_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/resume_application_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go`

- [ ] **Step 1: Implement JobPosting repository**

Create `backend/internal/infrastructure/persistence/postgres/job_posting_repo.go`:

```go
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
func (r *JobPostingRepo) Create(ctx context.Context, p *domainjob.JobPosting) error { return r.db.WithContext(ctx).Create(p).Error }
func (r *JobPostingRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) { var p domainjob.JobPosting; err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error; if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }; return &p, err }
func (r *JobPostingRepo) Update(ctx context.Context, p *domainjob.JobPosting) error { return r.db.WithContext(ctx).Save(p).Error }
func (r *JobPostingRepo) Delete(ctx context.Context, id, userID string) error { result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainjob.JobPosting{}); if result.Error != nil { return result.Error }; if result.RowsAffected == 0 { return domainerrors.ErrNotFound }; return nil }
func (r *JobPostingRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainjob.JobPosting{}).Where("user_id = ?", f.UserID)
	if f.Keyword != "" { like := "%" + f.Keyword + "%"; q = q.Where("company_name ILIKE ? OR job_title ILIKE ?", like, like) }
	if f.SourcePlatform != "" { q = q.Where("source_platform = ?", f.SourcePlatform) }
	if f.Status != "" { q = q.Where("status = ?", f.Status) }
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	var rows []domainjob.JobPosting
	err := q.Order("created_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}
```

- [ ] **Step 2: Implement ResumeApplication repository**

Create `backend/internal/infrastructure/persistence/postgres/resume_application_repo.go`:

```go
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
func (r *ResumeApplicationRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error { return r.db.WithContext(ctx).Create(a).Error }
func (r *ResumeApplicationRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error) { var a domainapp.ResumeApplication; err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error; if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }; return &a, err }
func (r *ResumeApplicationRepo) Update(ctx context.Context, a *domainapp.ResumeApplication) error { return r.db.WithContext(ctx).Save(a).Error }
func (r *ResumeApplicationRepo) Delete(ctx context.Context, id, userID string) error { result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainapp.ResumeApplication{}); if result.Error != nil { return result.Error }; if result.RowsAffected == 0 { return domainerrors.ErrNotFound }; return nil }
func (r *ResumeApplicationRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) { var count int64; err := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("resume_version_id = ?", versionID).Count(&count).Error; return count > 0, err }
func (r *ResumeApplicationRepo) UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error { result := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("id = ? AND user_id = ?", id, userID).Update("current_status", status); if result.Error != nil { return result.Error }; if result.RowsAffected == 0 { return domainerrors.ErrNotFound }; return nil }
func (r *ResumeApplicationRepo) ListByUser(ctx context.Context, f domainapp.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainapp.ResumeApplication{}).Where("user_id = ?", f.UserID)
	if f.ResumeVersionID != "" { q = q.Where("resume_version_id = ?", f.ResumeVersionID) }
	if f.RepoID != "" { q = q.Where("repo_id = ?", f.RepoID) }
	if f.Keyword != "" { like := "%" + f.Keyword + "%"; q = q.Where("company_name ILIKE ? OR job_title ILIKE ?", like, like) }
	if f.CurrentStatus != "" { q = q.Where("current_status = ?", f.CurrentStatus) }
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	var rows []domainapp.ResumeApplication
	err := q.Order("applied_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&rows).Error
	return rows, total, err
}
```

- [ ] **Step 3: Implement InterviewProgress repository**

Create `backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go`:

```go
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
func (r *InterviewProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error { return r.db.WithContext(ctx).Create(p).Error }
func (r *InterviewProgressRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) { var p domainprogress.InterviewProgress; err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error; if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrNotFound }; return &p, err }
func (r *InterviewProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error { return r.db.WithContext(ctx).Save(p).Error }
func (r *InterviewProgressRepo) Delete(ctx context.Context, id, userID string) error { result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainprogress.InterviewProgress{}); if result.Error != nil { return result.Error }; if result.RowsAffected == 0 { return domainerrors.ErrNotFound }; return nil }
func (r *InterviewProgressRepo) ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error) { var rows []domainprogress.InterviewProgress; err := r.db.WithContext(ctx).Where("application_id = ? AND user_id = ?", applicationID, userID).Order("occurred_at ASC").Find(&rows).Error; return rows, err }
func (r *InterviewProgressRepo) DeleteByApplicationID(ctx context.Context, applicationID, userID string) error { return r.db.WithContext(ctx).Where("application_id = ? AND user_id = ?", applicationID, userID).Delete(&domainprogress.InterviewProgress{}).Error }
```

- [ ] **Step 4: Run repository compile check**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/infrastructure/persistence/postgres -run '^$' -count=1
```

Expected: PASS or only unrelated existing integration setup failures. Fix compile errors before continuing.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/infrastructure/persistence/postgres/job_posting_repo.go internal/infrastructure/persistence/postgres/resume_application_repo.go internal/infrastructure/persistence/postgres/interview_progress_repo.go && git commit -m "feat: add postgres repositories for job domain"
```

---

### Task 10: Implement frontend route handlers

**Files:**
- Create: `backend/internal/interface/handler/job_posting_handler.go`
- Create: `backend/internal/interface/handler/resume_application_handler.go`
- Create: `backend/internal/interface/handler/interview_progress_handler.go`

- [ ] **Step 1: Implement JobPosting handler**

Create `backend/internal/interface/handler/job_posting_handler.go`:

```go
package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type JobPostingService interface {
	Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) error
	List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error)
	Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error
	Delete(ctx context.Context, userID, id string) error
	Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error)
}

type JobPostingHandler struct{ svc JobPostingService }
func NewJobPostingHandler(svc JobPostingService) *JobPostingHandler { return &JobPostingHandler{svc: svc} }

func (h *JobPostingHandler) List(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; page, pageSize := parsePage(c); rows, total, err := h.svc.List(c.Request.Context(), userID, c.Query("keyword"), c.Query("sourcePlatform"), c.Query("status"), page, pageSize); if err != nil { mapError(c, err); return }; response.Success(c, gin.H{"list": dto.ToJobPostingResponses(rows), "total": total, "page": page, "pageSize": pageSize}) }
func (h *JobPostingHandler) Get(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id")); if err != nil { mapError(c, err); return }; response.Success(c, dto.ToJobPostingResponse(row)) }
func (h *JobPostingHandler) Create(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.CreateJobPostingRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Create(c.Request.Context(), userID, req); err != nil { mapError(c, err); return }; rows, _, err := h.svc.List(c.Request.Context(), userID, "", "", "", 1, 1); if err == nil && len(rows) > 0 { response.Success(c, dto.ToJobPostingResponse(&rows[0])); return }; response.Success(c, nil) }
func (h *JobPostingHandler) Update(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.UpdateJobPostingRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req); err != nil { mapError(c, err); return }; row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id")); if err != nil { mapError(c, err); return }; response.Success(c, dto.ToJobPostingResponse(row)) }
func (h *JobPostingHandler) Delete(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil { mapError(c, err); return }; response.Success(c, nil) }
func (h *JobPostingHandler) Parse(c *gin.Context) { var req dto.ParseJobPostingRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "岗位链接不能为空"); return }; parsed, err := h.svc.Parse(c.Request.Context(), req); if err != nil { mapError(c, err); return }; response.Success(c, parsed) }

func parsePage(c *gin.Context) (int, int) { page, _ := strconv.Atoi(c.DefaultQuery("page", "1")); pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10")); if page <= 0 { page = 1 }; if pageSize <= 0 { pageSize = 10 }; return page, pageSize }
```

- [ ] **Step 2: Implement ResumeApplication handler**

Create `backend/internal/interface/handler/resume_application_handler.go`:

```go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeApplicationService interface {
	Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) error
	List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error)
	Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type ResumeApplicationHandler struct{ svc ResumeApplicationService }
func NewResumeApplicationHandler(svc ResumeApplicationService) *ResumeApplicationHandler { return &ResumeApplicationHandler{svc: svc} }

func (h *ResumeApplicationHandler) List(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; page, pageSize := parsePage(c); rows, total, err := h.svc.List(c.Request.Context(), userID, c.Query("resumeVersionId"), c.Query("repoId"), c.Query("keyword"), c.Query("currentStatus"), page, pageSize); if err != nil { mapError(c, err); return }; response.Success(c, gin.H{"list": dto.ToResumeApplicationResponses(rows), "total": total, "page": page, "pageSize": pageSize}) }
func (h *ResumeApplicationHandler) Get(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id")); if err != nil { mapError(c, err); return }; response.Success(c, dto.ToResumeApplicationResponse(row)) }
func (h *ResumeApplicationHandler) Create(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.CreateResumeApplicationRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Create(c.Request.Context(), userID, req); err != nil { mapError(c, err); return }; rows, _, err := h.svc.List(c.Request.Context(), userID, req.ResumeVersionID, "", "", "", 1, 1); if err == nil && len(rows) > 0 { response.Success(c, dto.ToResumeApplicationResponse(&rows[0])); return }; response.Success(c, nil) }
func (h *ResumeApplicationHandler) Update(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.UpdateResumeApplicationRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req); err != nil { mapError(c, err); return }; row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id")); if err != nil { mapError(c, err); return }; response.Success(c, dto.ToResumeApplicationResponse(row)) }
func (h *ResumeApplicationHandler) Delete(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil { mapError(c, err); return }; response.Success(c, nil) }
```

- [ ] **Step 3: Implement InterviewProgress handler**

Create `backend/internal/interface/handler/interview_progress_handler.go`:

```go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type InterviewProgressService interface {
	Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) error
	List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) error
	Delete(ctx context.Context, userID, id string) error
}

type InterviewProgressHandler struct{ svc InterviewProgressService }
func NewInterviewProgressHandler(svc InterviewProgressService) *InterviewProgressHandler { return &InterviewProgressHandler{svc: svc} }

func (h *InterviewProgressHandler) List(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; applicationID := c.Query("applicationId"); if applicationID == "" { response.Fail(c, pkgerrors.CodeBadParams, "applicationId 不能为空"); return }; rows, err := h.svc.List(c.Request.Context(), userID, applicationID); if err != nil { mapError(c, err); return }; response.Success(c, gin.H{"list": dto.ToInterviewProgressResponses(rows)}) }
func (h *InterviewProgressHandler) Create(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.CreateInterviewProgressRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Create(c.Request.Context(), userID, req); err != nil { mapError(c, err); return }; rows, err := h.svc.List(c.Request.Context(), userID, req.ApplicationID); if err == nil && len(rows) > 0 { response.Success(c, dto.ToInterviewProgressResponse(&rows[len(rows)-1])); return }; response.Success(c, nil) }
func (h *InterviewProgressHandler) Update(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; var req dto.UpdateInterviewProgressRequest; if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, pkgerrors.CodeBadParams, "参数错误"); return }; if err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req); err != nil { mapError(c, err); return }; response.Success(c, nil) }
func (h *InterviewProgressHandler) Delete(c *gin.Context) { userID, ok := resumeUserID(c); if !ok { return }; if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil { mapError(c, err); return }; response.Success(c, nil) }
```

- [ ] **Step 4: Run handler compile check**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/handler -run '^$' -count=1
```

Expected: PASS. If `resumeUserID` or `mapError` names differ in the current handler package, use the existing helper names already used by resume repo/version handlers.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/interface/handler/job_posting_handler.go internal/interface/handler/resume_application_handler.go internal/interface/handler/interview_progress_handler.go && git commit -m "feat: add frontend-aligned job domain handlers"
```

---

### Task 11: Wire router and server

**Files:**
- Modify: `backend/internal/interface/router/router.go`
- Modify: `backend/internal/interface/router/router_test.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Update router dependencies and routes**

In `backend/internal/interface/router/router.go`, add dependencies for the three new handlers using the same pattern as existing resume handlers, then register these authenticated routes under the existing `/api` or `/api/v1` group that is proxied by the frontend:

```go
auth.GET("/job-postings", deps.JobPostingHandler.List)
auth.GET("/job-postings/:id", deps.JobPostingHandler.Get)
auth.POST("/job-postings", deps.JobPostingHandler.Create)
auth.PUT("/job-postings/:id", deps.JobPostingHandler.Update)
auth.DELETE("/job-postings/:id", deps.JobPostingHandler.Delete)
auth.POST("/job-postings/parse", deps.JobPostingHandler.Parse)

auth.GET("/resume-applications", deps.ResumeApplicationHandler.List)
auth.GET("/resume-applications/:id", deps.ResumeApplicationHandler.Get)
auth.POST("/resume-applications", deps.ResumeApplicationHandler.Create)
auth.PUT("/resume-applications/:id", deps.ResumeApplicationHandler.Update)
auth.DELETE("/resume-applications/:id", deps.ResumeApplicationHandler.Delete)

auth.GET("/interview-progress", deps.InterviewProgressHandler.List)
auth.POST("/interview-progress", deps.InterviewProgressHandler.Create)
auth.PUT("/interview-progress/:id", deps.InterviewProgressHandler.Update)
auth.DELETE("/interview-progress/:id", deps.InterviewProgressHandler.Delete)
```

Important route-order rule: register `POST /job-postings/parse` before `GET /job-postings/:id` if this router matches static paths after params. Gin usually handles exact static routes correctly, but registering static paths first prevents ambiguity.

- [ ] **Step 2: Update router tests**

In `backend/internal/interface/router/router_test.go`, add these expected routes to the existing route list/assertions:

```go
"GET /api/job-postings",
"GET /api/job-postings/:id",
"POST /api/job-postings",
"PUT /api/job-postings/:id",
"DELETE /api/job-postings/:id",
"POST /api/job-postings/parse",
"GET /api/resume-applications",
"GET /api/resume-applications/:id",
"POST /api/resume-applications",
"PUT /api/resume-applications/:id",
"DELETE /api/resume-applications/:id",
"GET /api/interview-progress",
"POST /api/interview-progress",
"PUT /api/interview-progress/:id",
"DELETE /api/interview-progress/:id",
```

If the existing router test expects `/api/v1/...`, use `/api/v1/...` consistently. Match the current router prefix, not a new prefix.

- [ ] **Step 3: Wire repositories, services, and handlers in main**

In `backend/cmd/server/main.go`, add imports for the new packages:

```go
jobpostingapp "github.com/yilin/ai-for-backend/internal/application/job_posting"
interviewprogressapp "github.com/yilin/ai-for-backend/internal/application/interview_progress"
resumeapplicationapp "github.com/yilin/ai-for-backend/internal/application/resume_application"
domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
```

Add models to AutoMigrate only if `main.go` currently uses AutoMigrate:

```go
&domainjob.JobPosting{},
&domainapp.ResumeApplication{},
&domainprogress.InterviewProgress{},
```

Create repositories:

```go
jobPostingRepo := infrapostgres.NewJobPostingRepo(db)
resumeApplicationRepo := infrapostgres.NewResumeApplicationRepo(db)
interviewProgressRepo := infrapostgres.NewInterviewProgressRepo(db)
```

Create services:

```go
jobPostingSvc := jobpostingapp.NewJobPostingService(jobPostingRepo)
resumeApplicationSvc := resumeapplicationapp.NewResumeApplicationService(resumeApplicationRepo)
interviewProgressSvc := interviewprogressapp.NewInterviewProgressService(interviewProgressRepo, resumeApplicationSvc)
```

Create handlers:

```go
jobPostingHandler := handler.NewJobPostingHandler(jobPostingSvc)
resumeApplicationHandler := handler.NewResumeApplicationHandler(resumeApplicationSvc)
interviewProgressHandler := handler.NewInterviewProgressHandler(interviewProgressSvc)
```

Pass these handlers to router setup using the current router dependency style.

- [ ] **Step 4: Run build and router tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/router ./cmd/server/... -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
cd /Users/yilin/project/ai-for/backend && git add internal/interface/router/router.go internal/interface/router/router_test.go cmd/server/main.go && git commit -m "feat: wire frontend-aligned job domain routes"
```

---

### Task 12: Final verification

**Files:**
- No new files. Fix compile/test failures in files touched by earlier tasks.

- [ ] **Step 1: Run domain tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/domain/job_posting ./internal/domain/resume_application ./internal/domain/interview_progress -count=1
```

Expected: PASS.

- [ ] **Step 2: Run application tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/application/job_posting ./internal/application/resume_application ./internal/application/interview_progress -count=1
```

Expected: PASS.

- [ ] **Step 3: Run handler/router tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./internal/interface/handler ./internal/interface/router -count=1
```

Expected: PASS.

- [ ] **Step 4: Run full backend tests**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && go test ./... -count=1
```

Expected: PASS.

- [ ] **Step 5: Verify no frontend files changed**

Run:

```bash
cd /Users/yilin/project/ai-for/backend && git diff --name-only
```

Expected: changed files are backend files and this plan document only. No `frontend/` files should be listed.

- [ ] **Step 6: Commit final fixes if needed**

If Step 4 required fixes after the previous commits, commit only those fixes:

```bash
cd /Users/yilin/project/ai-for/backend && git add internal cmd migrations ../docs/superpowers/plans/2026-05-10-backend-plan3-job-domain.md && git commit -m "test: finalize frontend-aligned job domain"
```

---

## Self-Review

**Spec coverage:**

- JobPosting fields match `frontend/src/models/job-posting.ts`: covered in Tasks 1, 4, 5, 6, 9, 10, 11.
- JobPosting API routes match `frontend/src/api/job-postings.ts`: covered in Tasks 10 and 11.
- JobPosting filters `keyword`, `sourcePlatform`, `status`, `page`, `pageSize`: covered in Tasks 1, 6, 9, 10.
- ResumeApplication fields match `frontend/src/models/resume-application.ts`: covered in Tasks 2, 4, 5, 7, 9, 10, 11.
- ResumeApplication API routes match `frontend/src/api/resume-applications.ts`: covered in Tasks 10 and 11.
- ResumeApplication filters `resumeVersionId`, `repoId`, `keyword`, `currentStatus`, `page`, `pageSize`: covered in Tasks 2, 7, 9, 10.
- InterviewProgress fields match `frontend/src/models/interview-progress.ts`: covered in Tasks 3, 4, 5, 8, 9, 10, 11.
- InterviewProgress API routes match `frontend/src/api/interview-progress.ts`: covered in Tasks 10 and 11.
- Interview progress syncs parent application status from latest progress and failed result maps to rejected: covered in Task 8.
- Backend architecture remains DDD with domain/application/infrastructure/interface layers: covered across Tasks 1-11.
- Database columns remain snake_case while JSON remains frontend camelCase: covered in Tasks 4 and 5.

**Placeholder scan:** No TBD/TODO/implement-later placeholders remain. All new files have concrete code content.

**Type consistency:**

- `domainjob.Status` values are `pending/applied/not_fit/closed` and are reused by DTO/service/repository.
- `domainapp.Status` values match frontend `ApplicationStatus` and are reused by DTO/service/status sync.
- `domainprogress.Stage` values match application status names and convert to `domainapp.Status` during sync.
- Repository interfaces match service constructors and method calls.
- Handler interfaces match service methods.
- Route paths match frontend API wrappers, not the previous nested backend-only plan.
