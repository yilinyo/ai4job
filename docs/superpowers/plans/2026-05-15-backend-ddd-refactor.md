# Backend DDD Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the Go backend toward stricter DDD by moving repository interfaces into per-aggregate repository subpackages and introducing lightweight rich-domain creation/update methods without changing API behavior.

**Architecture:** Keep the existing `interface -> application -> domain <- infrastructure` dependency direction. Domain aggregate packages own entity behavior; `domain/<aggregate>/repository` packages own persistence contracts; application services orchestrate use cases and cross-aggregate ownership checks; postgres repositories keep existing GORM persistence behavior.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, Redis, `go test`, existing `pkg/errors` and `pkg/response` contracts.

---

## File Structure

### Repository interface relocation

- Create: `backend/internal/domain/user/repository/repository.go` — user repository contract.
- Create: `backend/internal/domain/resume_repo/repository/repository.go` — resume repository contract plus `ListFilter`.
- Create: `backend/internal/domain/resume_version/repository/repository.go` — resume version repository contract.
- Create: `backend/internal/domain/job_posting/repository/repository.go` — job posting repository contract plus `ListFilter`.
- Create: `backend/internal/domain/resume_application/repository/repository.go` — resume application repository contract plus `ListFilter`.
- Create: `backend/internal/domain/interview_progress/repository/repository.go` — interview progress repository contract.
- Delete after migration: existing `backend/internal/domain/*/repository.go` files.

### Domain behavior additions

- Modify: `backend/internal/domain/user/entity.go` — add `NewUser`, `NewUserParams`, `ProfileUpdate`, `UpdateProfile`.
- Modify: `backend/internal/domain/resume_repo/entity.go` — add `NewResumeRepo`, `NewResumeRepoParams`, `UpdateSpec`, `ApplyUpdate`.
- Modify: `backend/internal/domain/resume_version/entity.go` — add `NewResumeVersion`, `NewResumeVersionParams`, `UpdateSpec`, `ApplyUpdate`, keep `ValidateForCreate` behavior.
- Modify: `backend/internal/domain/job_posting/entity.go` — add `NewJobPosting`, `NewJobPostingParams`, `UpdateSpec`, `ApplyUpdate`, export `Status.Valid()`.
- Modify: `backend/internal/domain/resume_application/entity.go` — add `NewResumeApplication`, `NewResumeApplicationParams`, `UpdateSpec`, `ApplyUpdate`, `UpdateCurrentStatus`, export `Status.Valid()`.
- Modify: `backend/internal/domain/interview_progress/entity.go` — add `NewInterviewProgress`, `NewInterviewProgressParams`, `UpdateSpec`, `ApplyUpdate`, export `Stage.Valid()` and `Result.Valid()`.

### Domain tests

- Modify: `backend/internal/domain/user/entity_test.go` — cover `NewUser` and `UpdateProfile`.
- Modify: `backend/internal/domain/resume_repo/entity_test.go` — cover `NewResumeRepo` and `ApplyUpdate`.
- Modify: `backend/internal/domain/resume_version/entity_test.go` — cover `NewResumeVersion` and `ApplyUpdate`.
- Modify: `backend/internal/domain/job_posting/entity_test.go` — cover `NewJobPosting`, default status, `ApplyUpdate`, and `Status.Valid()`.
- Modify: `backend/internal/domain/resume_application/entity_test.go` — cover `NewResumeApplication`, default status, `ApplyUpdate`, `UpdateCurrentStatus`, and `Status.Valid()`.
- Modify: `backend/internal/domain/interview_progress/entity_test.go` — cover `NewInterviewProgress`, `ApplyUpdate`, `Stage.Valid()`, and `Result.Valid()`.

### Application and infrastructure adaptation

- Modify: all files under `backend/internal/application/*/*.go` and their tests that refer to `domainxxx.Repository` or `domainxxx.ListFilter`.
- Modify: postgres repository files under `backend/internal/infrastructure/persistence/postgres/*repo*.go` to import the new repository packages and add compile-time assertions.
- Modify: tests under `backend/internal/interface/**` only if imports break due to repository type relocation.
- Modify: `backend/cmd/server/main.go` only if constructor type changes require import cleanup; no DI behavior changes.

---

### Task 1: Move repository contracts into per-aggregate subpackages

**Files:**
- Create: `backend/internal/domain/user/repository/repository.go`
- Create: `backend/internal/domain/resume_repo/repository/repository.go`
- Create: `backend/internal/domain/resume_version/repository/repository.go`
- Create: `backend/internal/domain/job_posting/repository/repository.go`
- Create: `backend/internal/domain/resume_application/repository/repository.go`
- Create: `backend/internal/domain/interview_progress/repository/repository.go`
- Delete: `backend/internal/domain/user/repository.go`
- Delete: `backend/internal/domain/resume_repo/repository.go`
- Delete: `backend/internal/domain/resume_version/repository.go`
- Delete: `backend/internal/domain/job_posting/repository.go`
- Delete: `backend/internal/domain/resume_application/repository.go`
- Delete: `backend/internal/domain/interview_progress/repository.go`

- [ ] **Step 1: Create user repository subpackage**

Create `backend/internal/domain/user/repository/repository.go`:

```go
package repository

import (
	"context"

	"github.com/yilin/ai-for-backend/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
```

- [ ] **Step 2: Create resume repo repository subpackage**

Create `backend/internal/domain/resume_repo/repository/repository.go`:

```go
package repository

import (
	"context"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
)

type ListFilter struct {
	UserID   string
	JobType  domainrepo.JobType
	Keyword  string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, r *domainrepo.ResumeRepo) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error)
	Update(ctx context.Context, r *domainrepo.ResumeRepo) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainrepo.ResumeRepo, int64, error)
	ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error)
	IncrementVersionCount(ctx context.Context, id string, delta int) error
}
```

- [ ] **Step 3: Create resume version repository subpackage**

Create `backend/internal/domain/resume_version/repository/repository.go`:

```go
package repository

import (
	"context"

	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

type Repository interface {
	Create(ctx context.Context, v *domainversion.ResumeVersion) error
	FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error)
	FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error)
	ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error)
	Update(ctx context.Context, v *domainversion.ResumeVersion) error
	Delete(ctx context.Context, id, repoID string) error
	CountByRepoID(ctx context.Context, repoID string) (int64, error)
	SetActiveByID(ctx context.Context, repoID, id string) error
}
```

- [ ] **Step 4: Create job posting repository subpackage**

Create `backend/internal/domain/job_posting/repository/repository.go`:

```go
package repository

import (
	"context"

	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

type ListFilter struct {
	UserID         string
	Keyword        string
	SourcePlatform string
	Status         domainjob.Status
	Page           int
	PageSize       int
}

type Repository interface {
	Create(ctx context.Context, p *domainjob.JobPosting) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error)
	Update(ctx context.Context, p *domainjob.JobPosting) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainjob.JobPosting, int64, error)
}
```

- [ ] **Step 5: Create resume application repository subpackage**

Create `backend/internal/domain/resume_application/repository/repository.go`:

```go
package repository

import (
	"context"

	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type ListFilter struct {
	UserID          string
	ResumeVersionID string
	RepoID          string
	Keyword         string
	CurrentStatus   domainapp.Status
	Page            int
	PageSize        int
}

type Repository interface {
	Create(ctx context.Context, a *domainapp.ResumeApplication) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error)
	Update(ctx context.Context, a *domainapp.ResumeApplication) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainapp.ResumeApplication, int64, error)
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
	UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error
}
```

- [ ] **Step 6: Create interview progress repository subpackage**

Create `backend/internal/domain/interview_progress/repository/repository.go`:

```go
package repository

import (
	"context"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
)

type Repository interface {
	Create(ctx context.Context, p *domainprogress.InterviewProgress) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error)
	Update(ctx context.Context, p *domainprogress.InterviewProgress) error
	Delete(ctx context.Context, id, userID string) error
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error)
	DeleteByApplicationID(ctx context.Context, applicationID, userID string) error
}
```

- [ ] **Step 7: Remove old repository files**

Delete these files after the new files exist:

```text
backend/internal/domain/user/repository.go
backend/internal/domain/resume_repo/repository.go
backend/internal/domain/resume_version/repository.go
backend/internal/domain/job_posting/repository.go
backend/internal/domain/resume_application/repository.go
backend/internal/domain/interview_progress/repository.go
```

- [ ] **Step 8: Run tests to capture expected compile failures**

Run:

```bash
cd backend && go test ./...
```

Expected: FAIL. Errors should mention undefined `domainxxx.Repository` and `domainxxx.ListFilter` references. Do not fix by re-adding old repository files.

- [ ] **Step 9: Commit repository package skeleton**

```bash
git add backend/internal/domain/*/repository/repository.go backend/internal/domain/*/repository.go
git commit -m "refactor: move domain repository contracts into subpackages"
```

---

### Task 2: Update application services to use repository subpackages

**Files:**
- Modify: `backend/internal/application/user/user_service.go`
- Modify: `backend/internal/application/resume_repo/resume_repo_service.go`
- Modify: `backend/internal/application/resume_version/resume_version_service.go`
- Modify: `backend/internal/application/job_posting/job_posting_service.go`
- Modify: `backend/internal/application/resume_application/resume_application_service.go`
- Modify: `backend/internal/application/interview_progress/interview_progress_service.go`
- Modify: matching `*_test.go` files only if they directly mention repository types.

- [ ] **Step 1: Update user service repository type**

In `backend/internal/application/user/user_service.go`, add import:

```go
userrepo "github.com/yilin/ai-for-backend/internal/domain/user/repository"
```

Change fields and constructor:

```go
type UserService struct {
	repo      userrepo.Repository
	cache     TokenCache
	jwt       JWTManager
	emailMock bool
}

func NewUserService(repo userrepo.Repository, cache TokenCache, jwt JWTManager, emailMock bool) *UserService {
	return &UserService{repo: repo, cache: cache, jwt: jwt, emailMock: emailMock}
}
```

- [ ] **Step 2: Update resume repo service repository type and list filter**

In `backend/internal/application/resume_repo/resume_repo_service.go`, add import:

```go
resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"
```

Change service type and constructor:

```go
type ResumeRepoService struct {
	repo resumereporepo.Repository
}

func NewResumeRepoService(repo resumereporepo.Repository) *ResumeRepoService {
	return &ResumeRepoService{repo: repo}
}
```

Change `List` to use the repository package filter:

```go
func (s *ResumeRepoService) List(ctx context.Context, userID string, jobType domainrepo.JobType, keyword string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
	return s.repo.ListByUser(ctx, resumereporepo.ListFilter{UserID: userID, JobType: jobType, Keyword: keyword, Page: page, PageSize: pageSize})
}
```

- [ ] **Step 3: Update resume version service repository type**

In `backend/internal/application/resume_version/resume_version_service.go`, add import:

```go
resumeversionrepo "github.com/yilin/ai-for-backend/internal/domain/resume_version/repository"
```

Change service type and constructor:

```go
type ResumeVersionService struct {
	repoReader ResumeRepoReader
	repo       resumeversionrepo.Repository
	checker    VersionUseChecker
}

func NewResumeVersionService(repoReader ResumeRepoReader, repo resumeversionrepo.Repository, checker VersionUseChecker) *ResumeVersionService {
	return &ResumeVersionService{repoReader: repoReader, repo: repo, checker: checker}
}
```

- [ ] **Step 4: Update job posting service repository type and list filter**

In `backend/internal/application/job_posting/job_posting_service.go`, add import:

```go
jobpostingrepo "github.com/yilin/ai-for-backend/internal/domain/job_posting/repository"
```

Change service type and constructor:

```go
type JobPostingService struct {
	repo jobpostingrepo.Repository
}

func NewJobPostingService(repo jobpostingrepo.Repository) *JobPostingService {
	return &JobPostingService{repo: repo}
}
```

Change `List` filter construction:

```go
return s.repo.ListByUser(ctx, jobpostingrepo.ListFilter{
	UserID:         userID,
	Keyword:        keyword,
	SourcePlatform: sourcePlatform,
	Status:         filterStatus,
	Page:           page,
	PageSize:       pageSize,
})
```

- [ ] **Step 5: Update resume application service repository type and list filter**

In `backend/internal/application/resume_application/resume_application_service.go`, add import:

```go
resumeapplicationrepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
```

Change service type and constructor:

```go
type ResumeApplicationService struct {
	repo           resumeapplicationrepo.Repository
	jobPostings    JobPostingReader
	resumeRepos    ResumeRepoReader
	resumeVersions ResumeVersionReader
}

func NewResumeApplicationService(repo resumeapplicationrepo.Repository, jobPostings JobPostingReader, resumeRepos ResumeRepoReader, resumeVersions ResumeVersionReader) *ResumeApplicationService {
	return &ResumeApplicationService{repo: repo, jobPostings: jobPostings, resumeRepos: resumeRepos, resumeVersions: resumeVersions}
}
```

Change `List` filter construction:

```go
return s.repo.ListByUser(ctx, resumeapplicationrepo.ListFilter{
	UserID:          userID,
	ResumeVersionID: resumeVersionID,
	RepoID:          repoID,
	Keyword:         keyword,
	CurrentStatus:   filterStatus,
	Page:            page,
	PageSize:        pageSize,
})
```

- [ ] **Step 6: Update interview progress service repository type**

In `backend/internal/application/interview_progress/interview_progress_service.go`, add import:

```go
interviewprogressrepo "github.com/yilin/ai-for-backend/internal/domain/interview_progress/repository"
```

Change service type and constructor:

```go
type InterviewProgressService struct {
	repo interviewprogressrepo.Repository
	apps ApplicationService
}

func NewInterviewProgressService(repo interviewprogressrepo.Repository, apps ApplicationService) *InterviewProgressService {
	return &InterviewProgressService{repo: repo, apps: apps}
}
```

- [ ] **Step 7: Run application compile tests**

Run:

```bash
cd backend && go test ./internal/application/...
```

Expected: FAIL only if infrastructure or tests still import old repository symbols. No application service should refer to `domainxxx.Repository` or `domainxxx.ListFilter`.

- [ ] **Step 8: Commit application repository import migration**

```bash
git add backend/internal/application
git commit -m "refactor: point application services at domain repository subpackages"
```

---

### Task 3: Update postgres repositories to implement new contracts

**Files:**
- Modify: `backend/internal/infrastructure/persistence/postgres/user_repo.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/resume_repo_repo.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/resume_version_repo.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/job_posting_repo.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/resume_application_repo.go`
- Modify: `backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go`

- [ ] **Step 1: Add compile-time assertion for user repo**

In `backend/internal/infrastructure/persistence/postgres/user_repo.go`, add import:

```go
userrepo "github.com/yilin/ai-for-backend/internal/domain/user/repository"
```

After `type UserRepo struct`, add:

```go
var _ userrepo.Repository = (*UserRepo)(nil)
```

- [ ] **Step 2: Update resume repo filter and assertion**

In `backend/internal/infrastructure/persistence/postgres/resume_repo_repo.go`, add import:

```go
resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"
```

After `type ResumeRepoRepo struct`, add:

```go
var _ resumereporepo.Repository = (*ResumeRepoRepo)(nil)
```

Change signatures/helper types from `domainrepo.ListFilter` to `resumereporepo.ListFilter`:

```go
func (r *ResumeRepoRepo) ListByUser(ctx context.Context, f resumereporepo.ListFilter) ([]domainrepo.ResumeRepo, int64, error) {
```

and:

```go
func applyResumeRepoFilters(q *gorm.DB, f resumereporepo.ListFilter) *gorm.DB {
```

- [ ] **Step 3: Add compile-time assertion for resume version repo**

In `backend/internal/infrastructure/persistence/postgres/resume_version_repo.go`, add import:

```go
resumeversionrepo "github.com/yilin/ai-for-backend/internal/domain/resume_version/repository"
```

After `type ResumeVersionRepo struct`, add:

```go
var _ resumeversionrepo.Repository = (*ResumeVersionRepo)(nil)
```

- [ ] **Step 4: Update job posting filter and assertion**

In `backend/internal/infrastructure/persistence/postgres/job_posting_repo.go`, add import:

```go
jobpostingrepo "github.com/yilin/ai-for-backend/internal/domain/job_posting/repository"
```

After `type JobPostingRepo struct`, add:

```go
var _ jobpostingrepo.Repository = (*JobPostingRepo)(nil)
```

Change signatures/helper types from `domainjob.ListFilter` to `jobpostingrepo.ListFilter`:

```go
func (r *JobPostingRepo) ListByUser(ctx context.Context, f jobpostingrepo.ListFilter) ([]domainjob.JobPosting, int64, error) {
```

and:

```go
func applyJobPostingFilters(q *gorm.DB, f jobpostingrepo.ListFilter) *gorm.DB {
```

- [ ] **Step 5: Update resume application filter and assertion**

In `backend/internal/infrastructure/persistence/postgres/resume_application_repo.go`, add import:

```go
resumeapplicationrepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
```

After `type ResumeApplicationRepo struct`, add:

```go
var _ resumeapplicationrepo.Repository = (*ResumeApplicationRepo)(nil)
```

Change signatures/helper types from `domainapp.ListFilter` to `resumeapplicationrepo.ListFilter`:

```go
func (r *ResumeApplicationRepo) ListByUser(ctx context.Context, f resumeapplicationrepo.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
```

and:

```go
func applyResumeApplicationFilters(q *gorm.DB, f resumeapplicationrepo.ListFilter) *gorm.DB {
```

- [ ] **Step 6: Add compile-time assertion for interview progress repo**

In `backend/internal/infrastructure/persistence/postgres/interview_progress_repo.go`, add import:

```go
interviewprogressrepo "github.com/yilin/ai-for-backend/internal/domain/interview_progress/repository"
```

After `type InterviewProgressRepo struct`, add:

```go
var _ interviewprogressrepo.Repository = (*InterviewProgressRepo)(nil)
```

- [ ] **Step 7: Run infrastructure tests**

Run:

```bash
cd backend && go test ./internal/infrastructure/persistence/postgres/...
```

Expected: PASS or the same environment-dependent database failures that existed before. There should be no compile errors related to repository interfaces.

- [ ] **Step 8: Run full backend compile/test suite**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS or only pre-existing external service/database setup failures. There should be no undefined symbol or import-cycle errors.

- [ ] **Step 9: Commit infrastructure repository adaptation**

```bash
git add backend/internal/infrastructure/persistence/postgres backend/internal/application backend/internal/domain
git commit -m "refactor: adapt postgres repositories to domain repository packages"
```

---

### Task 4: Add rich-domain behavior for user and resume repo

**Files:**
- Modify: `backend/internal/domain/user/entity.go`
- Modify: `backend/internal/domain/user/entity_test.go`
- Modify: `backend/internal/domain/resume_repo/entity.go`
- Modify: `backend/internal/domain/resume_repo/entity_test.go`
- Modify: `backend/internal/application/user/user_service.go`
- Modify: `backend/internal/application/resume_repo/resume_repo_service.go`

- [ ] **Step 1: Add failing user domain tests**

Append to `backend/internal/domain/user/entity_test.go`:

```go
func TestNewUserInitializesAndHashesPassword(t *testing.T) {
	u, err := user.NewUser(user.NewUserParams{
		ID:       "user-1",
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", u.ID)
	assert.Equal(t, "alice", u.Username)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.NotEmpty(t, u.PasswordHash)
	assert.NotEqual(t, "secret123", u.PasswordHash)
	assert.NotNil(t, u.EducationExperiences)
	assert.True(t, u.CheckPassword("secret123"))
}

func TestNewUserRejectsInvalidPassword(t *testing.T) {
	u, err := user.NewUser(user.NewUserParams{ID: "user-1", Username: "alice", Email: "alice@example.com", Password: ""})

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestUserUpdateProfileUpdatesFieldsAndNormalizesEducation(t *testing.T) {
	u := &user.User{ID: "user-1", Username: "alice", Email: "alice@example.com"}

	u.UpdateProfile(user.ProfileUpdate{
		RealName:             "Alice",
		Age:                  22,
		Phone:                "13800000000",
		JobIntention:         "Backend",
		Avatar:               "avatar.png",
		Location:             "Shanghai",
		PersonalAdvantage:    "DDD",
		EducationExperiences: nil,
	})

	assert.Equal(t, "Alice", u.RealName)
	assert.Equal(t, 22, u.Age)
	assert.Equal(t, "Backend", u.JobIntention)
	assert.NotNil(t, u.EducationExperiences)
}
```

- [ ] **Step 2: Run user domain tests to verify failure**

Run:

```bash
cd backend && go test ./internal/domain/user -run 'TestNewUser|TestUserUpdateProfile' -v
```

Expected: FAIL with undefined `NewUser`, `NewUserParams`, `ProfileUpdate`, or `UpdateProfile`.

- [ ] **Step 3: Implement user domain methods**

Add to `backend/internal/domain/user/entity.go` after the `User` struct:

```go
type NewUserParams struct {
	ID       string
	Username string
	Email    string
	Password string
}

func NewUser(params NewUserParams) (*User, error) {
	u := &User{
		ID:                   params.ID,
		Username:             params.Username,
		Email:                params.Email,
		EducationExperiences: []EducationExperience{},
	}
	if err := u.SetPassword(params.Password); err != nil {
		return nil, err
	}
	u.Normalize()
	return u, nil
}

type ProfileUpdate struct {
	RealName             string
	Age                  int
	Phone                string
	JobIntention         string
	Avatar               string
	Location             string
	PersonalAdvantage    string
	EducationExperiences []EducationExperience
}

func (u *User) UpdateProfile(update ProfileUpdate) {
	u.RealName = update.RealName
	u.Age = update.Age
	u.Phone = update.Phone
	u.JobIntention = update.JobIntention
	u.Avatar = update.Avatar
	u.Location = update.Location
	u.PersonalAdvantage = update.PersonalAdvantage
	u.EducationExperiences = update.EducationExperiences
	u.Normalize()
}
```

- [ ] **Step 4: Refactor user application service to call domain methods**

In `backend/internal/application/user/user_service.go`, replace the manual user construction in `Register` with:

```go
u, err := domainuser.NewUser(domainuser.NewUserParams{
	ID:       uuid.NewString(),
	Username: req.Username,
	Email:    req.Email,
	Password: req.Password,
})
if err != nil {
	return err
}
if err := s.repo.Create(ctx, u); err != nil {
	return err
}
```

In `UpdateProfile`, replace direct field assignments with:

```go
u.UpdateProfile(domainuser.ProfileUpdate{
	RealName:             req.RealName,
	Age:                  req.Age,
	Phone:                req.Phone,
	JobIntention:         req.JobIntention,
	Avatar:               req.Avatar,
	Location:             req.Location,
	PersonalAdvantage:    req.PersonalAdvantage,
	EducationExperiences: req.EducationExperiences,
})
```

- [ ] **Step 5: Add failing resume repo domain tests**

Append to `backend/internal/domain/resume_repo/entity_test.go`:

```go
func TestNewResumeRepoDefaultsVersionCount(t *testing.T) {
	repo, err := resume_repo.NewResumeRepo(resume_repo.NewResumeRepoParams{
		ID:             "repo-1",
		UserID:         "user-1",
		Name:           "Backend",
		JobType:        resume_repo.JobTypeCampus,
		TargetPosition: "Go Engineer",
		Description:    "desc",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, repo.VersionCount)
	assert.Equal(t, "Backend", repo.Name)
}

func TestResumeRepoApplyUpdateValidatesResult(t *testing.T) {
	repo := validResumeRepo()

	err := repo.ApplyUpdate(resume_repo.UpdateSpec{
		Name:           "Backend v2",
		JobType:        resume_repo.JobTypeSocial,
		TargetPosition: "Senior Go Engineer",
		Description:    "updated",
	})

	require.NoError(t, err)
	assert.Equal(t, "Backend v2", repo.Name)
	assert.Equal(t, resume_repo.JobTypeSocial, repo.JobType)
	assert.Equal(t, "Senior Go Engineer", repo.TargetPosition)
	assert.Equal(t, "updated", repo.Description)
}

func TestResumeRepoApplyUpdateRejectsInvalidResult(t *testing.T) {
	repo := validResumeRepo()

	err := repo.ApplyUpdate(resume_repo.UpdateSpec{Name: "", JobType: resume_repo.JobTypeCampus, TargetPosition: "Go Engineer"})

	assert.Error(t, err)
}
```

- [ ] **Step 6: Run resume repo domain tests to verify failure**

Run:

```bash
cd backend && go test ./internal/domain/resume_repo -run 'TestNewResumeRepo|TestResumeRepoApplyUpdate' -v
```

Expected: FAIL with undefined `NewResumeRepo`, `NewResumeRepoParams`, `UpdateSpec`, or `ApplyUpdate`.

- [ ] **Step 7: Implement resume repo domain methods**

Add to `backend/internal/domain/resume_repo/entity.go` after the `ResumeRepo` struct:

```go
type NewResumeRepoParams struct {
	ID             string
	UserID         string
	Name           string
	JobType        JobType
	TargetPosition string
	Description    string
}

func NewResumeRepo(params NewResumeRepoParams) (*ResumeRepo, error) {
	repo := &ResumeRepo{
		ID:             params.ID,
		UserID:         params.UserID,
		Name:           params.Name,
		JobType:        params.JobType,
		TargetPosition: params.TargetPosition,
		Description:    params.Description,
		VersionCount:   0,
	}
	if err := repo.Validate(); err != nil {
		return nil, err
	}
	return repo, nil
}

type UpdateSpec struct {
	Name           string
	JobType        JobType
	TargetPosition string
	Description    string
}

func (r *ResumeRepo) ApplyUpdate(update UpdateSpec) error {
	r.Name = update.Name
	r.JobType = update.JobType
	r.TargetPosition = update.TargetPosition
	r.Description = update.Description
	return r.Validate()
}
```

- [ ] **Step 8: Refactor resume repo application service to call domain methods**

In `backend/internal/application/resume_repo/resume_repo_service.go`, replace manual construction in `Create` with:

```go
r, err := domainrepo.NewResumeRepo(domainrepo.NewResumeRepoParams{
	ID:             uuid.NewString(),
	UserID:         userID,
	Name:           req.Name,
	JobType:        req.JobType,
	TargetPosition: req.TargetPosition,
	Description:    req.Description,
})
if err != nil {
	return err
}
return s.repo.Create(ctx, r)
```

In `Update`, replace direct field assignments and validation with:

```go
if err := r.ApplyUpdate(domainrepo.UpdateSpec{
	Name:           req.Name,
	JobType:        req.JobType,
	TargetPosition: req.TargetPosition,
	Description:    req.Description,
}); err != nil {
	return err
}
return s.repo.Update(ctx, r)
```

- [ ] **Step 9: Run user and resume repo tests**

Run:

```bash
cd backend && go test ./internal/domain/user ./internal/domain/resume_repo ./internal/application/user ./internal/application/resume_repo
```

Expected: PASS.

- [ ] **Step 10: Commit user and resume repo domain behavior**

```bash
git add backend/internal/domain/user backend/internal/domain/resume_repo backend/internal/application/user backend/internal/application/resume_repo
git commit -m "refactor: move user and resume repo behavior into domain"
```

---

### Task 5: Add rich-domain behavior for resume versions and job postings

**Files:**
- Modify: `backend/internal/domain/resume_version/entity.go`
- Modify: `backend/internal/domain/resume_version/entity_test.go`
- Modify: `backend/internal/domain/job_posting/entity.go`
- Modify: `backend/internal/domain/job_posting/entity_test.go`
- Modify: `backend/internal/application/resume_version/resume_version_service.go`
- Modify: `backend/internal/application/job_posting/job_posting_service.go`

- [ ] **Step 1: Add failing resume version domain tests**

Append to `backend/internal/domain/resume_version/entity_test.go`:

```go
func TestNewResumeVersionBuildsValidVersion(t *testing.T) {
	version, err := resume_version.NewResumeVersion(resume_version.NewResumeVersionParams{
		ID:       "version-1",
		RepoID:   "repo-1",
		Version:  "v1.0.0",
		Title:    "Backend Resume",
		Content:  "# Resume",
		Remark:   "first",
		IsActive: false,
	})

	require.NoError(t, err)
	assert.Equal(t, "version-1", version.ID)
	assert.Equal(t, "repo-1", version.RepoID)
	assert.False(t, version.IsActive)
}

func TestResumeVersionApplyUpdateHandlesActivationFlag(t *testing.T) {
	version := &resume_version.ResumeVersion{ID: "version-1", RepoID: "repo-1", Version: "v1.0.0", Title: "old", Content: "old", IsActive: true}
	active := false

	err := version.ApplyUpdate(resume_version.UpdateSpec{Title: "new", Content: "# New", Remark: "updated", IsActive: &active})

	require.NoError(t, err)
	assert.Equal(t, "new", version.Title)
	assert.Equal(t, "# New", version.Content)
	assert.Equal(t, "updated", version.Remark)
	assert.False(t, version.IsActive)
}

func TestResumeVersionApplyUpdateLeavesActivationForTrueFlagToRepository(t *testing.T) {
	version := &resume_version.ResumeVersion{ID: "version-1", RepoID: "repo-1", Version: "v1.0.0", Title: "old", Content: "old", IsActive: false}
	active := true

	err := version.ApplyUpdate(resume_version.UpdateSpec{Title: "new", Content: "# New", IsActive: &active})

	require.NoError(t, err)
	assert.False(t, version.IsActive)
}
```

- [ ] **Step 2: Implement resume version domain methods**

Add to `backend/internal/domain/resume_version/entity.go` after the `ResumeVersion` struct:

```go
type NewResumeVersionParams struct {
	ID       string
	RepoID   string
	Version  string
	Title    string
	Content  string
	Remark   string
	IsActive bool
}

func NewResumeVersion(params NewResumeVersionParams) (*ResumeVersion, error) {
	version := &ResumeVersion{
		ID:       params.ID,
		RepoID:   params.RepoID,
		Version:  params.Version,
		Title:    params.Title,
		Content:  params.Content,
		Remark:   params.Remark,
		IsActive: params.IsActive,
	}
	if err := version.ValidateForCreate(); err != nil {
		return nil, err
	}
	return version, nil
}

type UpdateSpec struct {
	Title    string
	Content  string
	Remark   string
	IsActive *bool
}

func (v *ResumeVersion) ApplyUpdate(update UpdateSpec) error {
	v.Title = update.Title
	v.Content = update.Content
	v.Remark = update.Remark
	if update.IsActive != nil && !*update.IsActive {
		v.IsActive = false
	}
	return v.ValidateForCreate()
}
```

- [ ] **Step 3: Refactor resume version service**

In `backend/internal/application/resume_version/resume_version_service.go`, replace manual construction in `Create` with:

```go
v, err := domainversion.NewResumeVersion(domainversion.NewResumeVersionParams{
	ID:       uuid.NewString(),
	RepoID:   req.RepoID,
	Version:  fmt.Sprintf("v1.%d.0", count),
	Title:    req.Title,
	Content:  req.Content,
	Remark:   req.Remark,
	IsActive: false,
})
if err != nil {
	return err
}
```

In `Update`, replace field assignments and validation with:

```go
activate := req.IsActive != nil && *req.IsActive
if err := v.ApplyUpdate(domainversion.UpdateSpec{
	Title:    req.Title,
	Content:  req.Content,
	Remark:   req.Remark,
	IsActive: req.IsActive,
}); err != nil {
	return err
}
```

Keep the existing repository update and `SetActiveByID` flow after this block.

- [ ] **Step 4: Add failing job posting domain tests**

Append to `backend/internal/domain/job_posting/entity_test.go`:

```go
func TestNewJobPostingDefaultsPendingStatus(t *testing.T) {
	posting, err := job_posting.NewJobPosting(job_posting.NewJobPostingParams{
		ID:          "job-1",
		UserID:      "user-1",
		CompanyName: "Acme",
		JobTitle:    "Backend Engineer",
	})

	require.NoError(t, err)
	assert.Equal(t, job_posting.StatusPending, posting.Status)
}

func TestJobPostingStatusValid(t *testing.T) {
	assert.True(t, job_posting.StatusPending.Valid())
	assert.True(t, job_posting.StatusApplied.Valid())
	assert.False(t, job_posting.Status("unknown").Valid())
}

func TestJobPostingApplyUpdateValidatesResult(t *testing.T) {
	posting := &job_posting.JobPosting{ID: "job-1", UserID: "user-1", CompanyName: "Acme", JobTitle: "Backend", Status: job_posting.StatusPending}
	company := "Beta"
	status := job_posting.StatusApplied

	err := posting.ApplyUpdate(job_posting.UpdateSpec{CompanyName: &company, Status: &status})

	require.NoError(t, err)
	assert.Equal(t, "Beta", posting.CompanyName)
	assert.Equal(t, job_posting.StatusApplied, posting.Status)
}

func TestJobPostingApplyUpdateRejectsInvalidStatus(t *testing.T) {
	posting := &job_posting.JobPosting{ID: "job-1", UserID: "user-1", CompanyName: "Acme", JobTitle: "Backend", Status: job_posting.StatusPending}
	status := job_posting.Status("unknown")

	err := posting.ApplyUpdate(job_posting.UpdateSpec{Status: &status})

	assert.Error(t, err)
}
```

- [ ] **Step 5: Implement job posting domain methods**

Add to `backend/internal/domain/job_posting/entity.go` after the `JobPosting` struct:

```go
type NewJobPostingParams struct {
	ID              string
	UserID          string
	CompanyName     string
	CompanyIntro    string
	JobTitle        string
	JobRequirements string
	Base            string
	SalaryRange     string
	SourcePlatform  string
	SourceURL       string
	Status          Status
	Remark          string
}

func NewJobPosting(params NewJobPostingParams) (*JobPosting, error) {
	status := params.Status
	if status == "" {
		status = StatusPending
	}
	posting := &JobPosting{
		ID:              params.ID,
		UserID:          params.UserID,
		CompanyName:     params.CompanyName,
		CompanyIntro:    params.CompanyIntro,
		JobTitle:        params.JobTitle,
		JobRequirements: params.JobRequirements,
		Base:            params.Base,
		SalaryRange:     params.SalaryRange,
		SourcePlatform:  params.SourcePlatform,
		SourceURL:       params.SourceURL,
		Status:          status,
		Remark:          params.Remark,
	}
	if err := posting.Validate(); err != nil {
		return nil, err
	}
	return posting, nil
}

type UpdateSpec struct {
	CompanyName     *string
	CompanyIntro    *string
	JobTitle        *string
	JobRequirements *string
	Base            *string
	SalaryRange     *string
	SourcePlatform  *string
	SourceURL       *string
	Status          *Status
	Remark          *string
}

func (p *JobPosting) ApplyUpdate(update UpdateSpec) error {
	if update.CompanyName != nil {
		p.CompanyName = *update.CompanyName
	}
	if update.CompanyIntro != nil {
		p.CompanyIntro = *update.CompanyIntro
	}
	if update.JobTitle != nil {
		p.JobTitle = *update.JobTitle
	}
	if update.JobRequirements != nil {
		p.JobRequirements = *update.JobRequirements
	}
	if update.Base != nil {
		p.Base = *update.Base
	}
	if update.SalaryRange != nil {
		p.SalaryRange = *update.SalaryRange
	}
	if update.SourcePlatform != nil {
		p.SourcePlatform = *update.SourcePlatform
	}
	if update.SourceURL != nil {
		p.SourceURL = *update.SourceURL
	}
	if update.Status != nil {
		p.Status = *update.Status
	}
	if update.Remark != nil {
		p.Remark = *update.Remark
	}
	return p.Validate()
}

func (s Status) Valid() bool {
	return validStatus(s)
}
```

- [ ] **Step 6: Refactor job posting service**

In `backend/internal/application/job_posting/job_posting_service.go`, replace manual construction in `Create` with:

```go
posting, err := domainjob.NewJobPosting(domainjob.NewJobPostingParams{
	ID:              uuid.NewString(),
	UserID:          userID,
	CompanyName:     req.CompanyName,
	CompanyIntro:    req.CompanyIntro,
	JobTitle:        req.JobTitle,
	JobRequirements: req.JobRequirements,
	Base:            req.Base,
	SalaryRange:     req.SalaryRange,
	SourcePlatform:  req.SourcePlatform,
	SourceURL:       req.SourceURL,
	Status:          req.Status,
	Remark:          req.Remark,
})
if err != nil {
	return nil, err
}
```

Change `List` status validation to:

```go
if status != "" && !filterStatus.Valid() {
	return nil, 0, domainerrors.ErrBadParams
}
```

In `Update`, replace direct field assignments and validation with:

```go
if err := posting.ApplyUpdate(domainjob.UpdateSpec{
	CompanyName:     req.CompanyName,
	CompanyIntro:    req.CompanyIntro,
	JobTitle:        req.JobTitle,
	JobRequirements: req.JobRequirements,
	Base:            req.Base,
	SalaryRange:     req.SalaryRange,
	SourcePlatform:  req.SourcePlatform,
	SourceURL:       req.SourceURL,
	Status:          req.Status,
	Remark:          req.Remark,
}); err != nil {
	return err
}
return s.repo.Update(ctx, posting)
```

Delete the service-level `validStatus` function from `job_posting_service.go` after replacing all uses.

- [ ] **Step 7: Run focused tests**

Run:

```bash
cd backend && go test ./internal/domain/resume_version ./internal/domain/job_posting ./internal/application/resume_version ./internal/application/job_posting
```

Expected: PASS.

- [ ] **Step 8: Commit resume version and job posting domain behavior**

```bash
git add backend/internal/domain/resume_version backend/internal/domain/job_posting backend/internal/application/resume_version backend/internal/application/job_posting
git commit -m "refactor: move resume version and job posting behavior into domain"
```

---

### Task 6: Add rich-domain behavior for applications and interview progress

**Files:**
- Modify: `backend/internal/domain/resume_application/entity.go`
- Modify: `backend/internal/domain/resume_application/entity_test.go`
- Modify: `backend/internal/domain/interview_progress/entity.go`
- Modify: `backend/internal/domain/interview_progress/entity_test.go`
- Modify: `backend/internal/application/resume_application/resume_application_service.go`
- Modify: `backend/internal/application/interview_progress/interview_progress_service.go`

- [ ] **Step 1: Add failing resume application domain tests**

Append to `backend/internal/domain/resume_application/entity_test.go`:

```go
func TestNewResumeApplicationDefaultsAppliedStatus(t *testing.T) {
	appliedAt := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	application, err := resume_application.NewResumeApplication(resume_application.NewResumeApplicationParams{
		ID:              "app-1",
		UserID:          "user-1",
		ResumeVersionID: "version-1",
		RepoID:          "repo-1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       appliedAt,
	})

	require.NoError(t, err)
	assert.Equal(t, resume_application.StatusApplied, application.CurrentStatus)
}

func TestResumeApplicationStatusValid(t *testing.T) {
	assert.True(t, resume_application.StatusApplied.Valid())
	assert.True(t, resume_application.StatusOffer.Valid())
	assert.False(t, resume_application.Status("unknown").Valid())
}

func TestResumeApplicationApplyUpdateValidatesResult(t *testing.T) {
	application := validResumeApplication()
	company := "Beta"
	status := resume_application.StatusViewed

	err := application.ApplyUpdate(resume_application.UpdateSpec{CompanyName: &company, CurrentStatus: &status})

	require.NoError(t, err)
	assert.Equal(t, "Beta", application.CompanyName)
	assert.Equal(t, resume_application.StatusViewed, application.CurrentStatus)
}

func TestResumeApplicationUpdateCurrentStatusValidatesStatus(t *testing.T) {
	application := validResumeApplication()

	err := application.UpdateCurrentStatus(resume_application.StatusOffer)

	require.NoError(t, err)
	assert.Equal(t, resume_application.StatusOffer, application.CurrentStatus)
	assert.Error(t, application.UpdateCurrentStatus(resume_application.Status("unknown")))
}
```

If `validResumeApplication()` does not exist, add this helper in the test file:

```go
func validResumeApplication() *resume_application.ResumeApplication {
	return &resume_application.ResumeApplication{
		ID:              "app-1",
		UserID:          "user-1",
		ResumeVersionID: "version-1",
		RepoID:          "repo-1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		CurrentStatus:   resume_application.StatusApplied,
	}
}
```

- [ ] **Step 2: Implement resume application domain methods**

Add to `backend/internal/domain/resume_application/entity.go` after the `ResumeApplication` struct:

```go
type NewResumeApplicationParams struct {
	ID                        string
	UserID                    string
	JobPostingID              *string
	ResumeVersionID           string
	RepoID                    string
	CompanyName               string
	CompanyIntro              string
	JobTitle                  string
	JobRequirements           string
	Base                      string
	SalaryRange               string
	SourcePlatform            string
	SourceURL                 string
	DeliveryChannel           string
	AppliedAt                 time.Time
	CurrentStatus             Status
	InterviewSummary          string
	InterviewNotes            string
	ResumeMatchScore          int
	InterviewPerformanceScore int
	RoleInterestScore         int
	OverallScore              int
	Remark                    string
}

func NewResumeApplication(params NewResumeApplicationParams) (*ResumeApplication, error) {
	status := params.CurrentStatus
	if status == "" {
		status = StatusApplied
	}
	application := &ResumeApplication{
		ID:                        params.ID,
		UserID:                    params.UserID,
		JobPostingID:              params.JobPostingID,
		ResumeVersionID:           params.ResumeVersionID,
		RepoID:                    params.RepoID,
		CompanyName:               params.CompanyName,
		CompanyIntro:              params.CompanyIntro,
		JobTitle:                  params.JobTitle,
		JobRequirements:           params.JobRequirements,
		Base:                      params.Base,
		SalaryRange:               params.SalaryRange,
		SourcePlatform:            params.SourcePlatform,
		SourceURL:                 params.SourceURL,
		DeliveryChannel:           params.DeliveryChannel,
		AppliedAt:                 params.AppliedAt,
		CurrentStatus:             status,
		InterviewSummary:          params.InterviewSummary,
		InterviewNotes:            params.InterviewNotes,
		ResumeMatchScore:          params.ResumeMatchScore,
		InterviewPerformanceScore: params.InterviewPerformanceScore,
		RoleInterestScore:         params.RoleInterestScore,
		OverallScore:              params.OverallScore,
		Remark:                    params.Remark,
	}
	if err := application.Validate(); err != nil {
		return nil, err
	}
	return application, nil
}

type UpdateSpec struct {
	JobPostingID              **string
	ResumeVersionID           *string
	RepoID                    *string
	CompanyName               *string
	CompanyIntro              *string
	JobTitle                  *string
	JobRequirements           *string
	Base                      *string
	SalaryRange               *string
	SourcePlatform            *string
	SourceURL                 *string
	DeliveryChannel           *string
	AppliedAt                 *time.Time
	CurrentStatus             *Status
	InterviewSummary          *string
	InterviewNotes            *string
	ResumeMatchScore          *int
	InterviewPerformanceScore *int
	RoleInterestScore         *int
	OverallScore              *int
	Remark                    *string
}

func (a *ResumeApplication) ApplyUpdate(update UpdateSpec) error {
	if update.JobPostingID != nil {
		a.JobPostingID = *update.JobPostingID
	}
	if update.ResumeVersionID != nil {
		a.ResumeVersionID = *update.ResumeVersionID
	}
	if update.RepoID != nil {
		a.RepoID = *update.RepoID
	}
	if update.CompanyName != nil {
		a.CompanyName = *update.CompanyName
	}
	if update.CompanyIntro != nil {
		a.CompanyIntro = *update.CompanyIntro
	}
	if update.JobTitle != nil {
		a.JobTitle = *update.JobTitle
	}
	if update.JobRequirements != nil {
		a.JobRequirements = *update.JobRequirements
	}
	if update.Base != nil {
		a.Base = *update.Base
	}
	if update.SalaryRange != nil {
		a.SalaryRange = *update.SalaryRange
	}
	if update.SourcePlatform != nil {
		a.SourcePlatform = *update.SourcePlatform
	}
	if update.SourceURL != nil {
		a.SourceURL = *update.SourceURL
	}
	if update.DeliveryChannel != nil {
		a.DeliveryChannel = *update.DeliveryChannel
	}
	if update.AppliedAt != nil {
		a.AppliedAt = *update.AppliedAt
	}
	if update.CurrentStatus != nil {
		a.CurrentStatus = *update.CurrentStatus
	}
	if update.InterviewSummary != nil {
		a.InterviewSummary = *update.InterviewSummary
	}
	if update.InterviewNotes != nil {
		a.InterviewNotes = *update.InterviewNotes
	}
	if update.ResumeMatchScore != nil {
		a.ResumeMatchScore = *update.ResumeMatchScore
	}
	if update.InterviewPerformanceScore != nil {
		a.InterviewPerformanceScore = *update.InterviewPerformanceScore
	}
	if update.RoleInterestScore != nil {
		a.RoleInterestScore = *update.RoleInterestScore
	}
	if update.OverallScore != nil {
		a.OverallScore = *update.OverallScore
	}
	if update.Remark != nil {
		a.Remark = *update.Remark
	}
	return a.Validate()
}

func (a *ResumeApplication) UpdateCurrentStatus(status Status) error {
	oldStatus := a.CurrentStatus
	a.CurrentStatus = status
	if err := a.Validate(); err != nil {
		a.CurrentStatus = oldStatus
		return err
	}
	return nil
}

func (s Status) Valid() bool {
	return validStatus(s)
}
```

- [ ] **Step 3: Refactor resume application service**

In `backend/internal/application/resume_application/resume_application_service.go`, replace manual construction in `Create` with:

```go
application, err := domainapp.NewResumeApplication(domainapp.NewResumeApplicationParams{
	ID:                        uuid.NewString(),
	UserID:                    userID,
	JobPostingID:              normalizeOptionalJobPostingID(req.JobPostingID),
	ResumeVersionID:           req.ResumeVersionID,
	RepoID:                    req.RepoID,
	CompanyName:               req.CompanyName,
	CompanyIntro:              req.CompanyIntro,
	JobTitle:                  req.JobTitle,
	JobRequirements:           req.JobRequirements,
	Base:                      req.Base,
	SalaryRange:               req.SalaryRange,
	SourcePlatform:            req.SourcePlatform,
	SourceURL:                 req.SourceURL,
	DeliveryChannel:           req.DeliveryChannel,
	AppliedAt:                 req.AppliedAt.Value(),
	CurrentStatus:             req.CurrentStatus,
	InterviewSummary:          req.InterviewSummary,
	InterviewNotes:            req.InterviewNotes,
	ResumeMatchScore:          req.ResumeMatchScore,
	InterviewPerformanceScore: req.InterviewPerformanceScore,
	RoleInterestScore:         req.RoleInterestScore,
	OverallScore:              req.OverallScore,
	Remark:                    req.Remark,
})
if err != nil {
	return nil, err
}
```

Change `List` status validation to:

```go
if currentStatus != "" && !filterStatus.Valid() {
	return nil, 0, domainerrors.ErrBadParams
}
```

In `Update`, convert the optional applied date and normalized job posting pointer before calling domain:

```go
var appliedAt *time.Time
if req.AppliedAt != nil {
	value := req.AppliedAt.Value()
	appliedAt = &value
}
var jobPostingID **string
if req.JobPostingID != nil {
	value := normalizeOptionalJobPostingID(req.JobPostingID)
	jobPostingID = &value
}
if err := application.ApplyUpdate(domainapp.UpdateSpec{
	JobPostingID:              jobPostingID,
	ResumeVersionID:           req.ResumeVersionID,
	RepoID:                    req.RepoID,
	CompanyName:               req.CompanyName,
	CompanyIntro:              req.CompanyIntro,
	JobTitle:                  req.JobTitle,
	JobRequirements:           req.JobRequirements,
	Base:                      req.Base,
	SalaryRange:               req.SalaryRange,
	SourcePlatform:            req.SourcePlatform,
	SourceURL:                 req.SourceURL,
	DeliveryChannel:           req.DeliveryChannel,
	AppliedAt:                 appliedAt,
	CurrentStatus:             req.CurrentStatus,
	InterviewSummary:          req.InterviewSummary,
	InterviewNotes:            req.InterviewNotes,
	ResumeMatchScore:          req.ResumeMatchScore,
	InterviewPerformanceScore: req.InterviewPerformanceScore,
	RoleInterestScore:         req.RoleInterestScore,
	OverallScore:              req.OverallScore,
	Remark:                    req.Remark,
}); err != nil {
	return err
}
```

In `UpdateCurrentStatus`, replace service-level validation with:

```go
application, err := s.repo.FindByIDAndUserID(ctx, id, userID)
if err != nil {
	return err
}
if err := application.UpdateCurrentStatus(status); err != nil {
	return err
}
return s.repo.UpdateCurrentStatus(ctx, id, userID, application.CurrentStatus)
```

Delete the service-level `validStatus` function from `resume_application_service.go` after replacing all uses.

Add `time` to the imports if not already present.

- [ ] **Step 4: Add failing interview progress domain tests**

Append to `backend/internal/domain/interview_progress/entity_test.go`:

```go
func TestNewInterviewProgressBuildsValidProgress(t *testing.T) {
	occurredAt := time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC)
	progress, err := interview_progress.NewInterviewProgress(interview_progress.NewInterviewProgressParams{
		ID:                "progress-1",
		UserID:            "user-1",
		ApplicationID:     "app-1",
		Stage:             interview_progress.StageFirstInterview,
		OccurredAt:        occurredAt,
		Result:            interview_progress.ResultPending,
		InterviewerOrTeam: "Backend Team",
		Note:              "note",
	})

	require.NoError(t, err)
	assert.Equal(t, "progress-1", progress.ID)
	assert.Equal(t, interview_progress.StageFirstInterview, progress.Stage)
}

func TestInterviewProgressEnumsValid(t *testing.T) {
	assert.True(t, interview_progress.StageApplied.Valid())
	assert.True(t, interview_progress.ResultPassed.Valid())
	assert.False(t, interview_progress.Stage("unknown").Valid())
	assert.False(t, interview_progress.Result("unknown").Valid())
}

func TestInterviewProgressApplyUpdateValidatesResult(t *testing.T) {
	progress := validInterviewProgress()
	stage := interview_progress.StageOffer
	result := interview_progress.ResultPassed

	err := progress.ApplyUpdate(interview_progress.UpdateSpec{Stage: &stage, Result: &result})

	require.NoError(t, err)
	assert.Equal(t, interview_progress.StageOffer, progress.Stage)
	assert.Equal(t, interview_progress.ResultPassed, progress.Result)
}
```

If `validInterviewProgress()` does not exist, add this helper in the test file:

```go
func validInterviewProgress() *interview_progress.InterviewProgress {
	return &interview_progress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         interview_progress.StageFirstInterview,
		OccurredAt:    time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC),
		Result:        interview_progress.ResultPending,
	}
}
```

- [ ] **Step 5: Implement interview progress domain methods**

Add to `backend/internal/domain/interview_progress/entity.go` after the `InterviewProgress` struct:

```go
type NewInterviewProgressParams struct {
	ID                string
	UserID            string
	ApplicationID     string
	Stage             Stage
	OccurredAt        time.Time
	Result            Result
	InterviewerOrTeam string
	Note              string
}

func NewInterviewProgress(params NewInterviewProgressParams) (*InterviewProgress, error) {
	progress := &InterviewProgress{
		ID:                params.ID,
		UserID:            params.UserID,
		ApplicationID:     params.ApplicationID,
		Stage:             params.Stage,
		OccurredAt:        params.OccurredAt,
		Result:            params.Result,
		InterviewerOrTeam: params.InterviewerOrTeam,
		Note:              params.Note,
	}
	if err := progress.Validate(); err != nil {
		return nil, err
	}
	return progress, nil
}

type UpdateSpec struct {
	ApplicationID     *string
	Stage             *Stage
	OccurredAt        *time.Time
	Result            *Result
	InterviewerOrTeam *string
	Note              *string
}

func (p *InterviewProgress) ApplyUpdate(update UpdateSpec) error {
	if update.ApplicationID != nil {
		p.ApplicationID = *update.ApplicationID
	}
	if update.Stage != nil {
		p.Stage = *update.Stage
	}
	if update.OccurredAt != nil {
		p.OccurredAt = *update.OccurredAt
	}
	if update.Result != nil {
		p.Result = *update.Result
	}
	if update.InterviewerOrTeam != nil {
		p.InterviewerOrTeam = *update.InterviewerOrTeam
	}
	if update.Note != nil {
		p.Note = *update.Note
	}
	return p.Validate()
}

func (s Stage) Valid() bool {
	return validStage(s)
}

func (r Result) Valid() bool {
	return validResult(r)
}
```

- [ ] **Step 6: Refactor interview progress service**

In `backend/internal/application/interview_progress/interview_progress_service.go`, replace manual construction in `Create` with:

```go
progress, err := domainprogress.NewInterviewProgress(domainprogress.NewInterviewProgressParams{
	ID:                uuid.NewString(),
	UserID:            userID,
	ApplicationID:     req.ApplicationID,
	Stage:             req.Stage,
	OccurredAt:        req.OccurredAt.Value(),
	Result:            req.Result,
	InterviewerOrTeam: req.InterviewerOrTeam,
	Note:              req.Note,
})
if err != nil {
	return nil, err
}
```

In `Update`, convert optional occurred date before calling domain:

```go
var occurredAt *time.Time
if req.OccurredAt != nil {
	value := req.OccurredAt.Value()
	occurredAt = &value
}
if err := progress.ApplyUpdate(domainprogress.UpdateSpec{
	ApplicationID:     req.ApplicationID,
	Stage:             req.Stage,
	OccurredAt:        occurredAt,
	Result:            req.Result,
	InterviewerOrTeam: req.InterviewerOrTeam,
	Note:              req.Note,
}); err != nil {
	return nil, err
}
```

Keep the existing `oldApplicationID`, ownership check, repository update, and sync flow unchanged.

Add `time` to imports if not already present.

- [ ] **Step 7: Run focused tests**

Run:

```bash
cd backend && go test ./internal/domain/resume_application ./internal/domain/interview_progress ./internal/application/resume_application ./internal/application/interview_progress
```

Expected: PASS.

- [ ] **Step 8: Commit application and interview progress domain behavior**

```bash
git add backend/internal/domain/resume_application backend/internal/domain/interview_progress backend/internal/application/resume_application backend/internal/application/interview_progress
git commit -m "refactor: move application and interview progress behavior into domain"
```

---

### Task 7: Final API-behavior verification and cleanup

**Files:**
- Modify only files required by failing tests or formatting.
- Verify: `backend/internal/interface/handler/*_test.go`
- Verify: `backend/internal/interface/router/router_test.go`
- Verify: `backend/internal/application/dto/*_test.go`

- [ ] **Step 1: Format backend Go code**

Run:

```bash
cd backend && gofmt -w cmd internal pkg
```

Expected: no output.

- [ ] **Step 2: Run full backend tests**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS. If external database tests fail because local PostgreSQL/Redis is unavailable, capture the exact failing package and error. Compile errors must be fixed.

- [ ] **Step 3: Verify no old repository package references remain**

Run:

```bash
cd backend && grep -R "domain.*\.Repository\|domain.*\.ListFilter" -n internal --include='*.go'
```

Expected: no matches for old `domainxxx.Repository` or `domainxxx.ListFilter` patterns. Matches inside comments should be removed or updated.

- [ ] **Step 4: Verify API layer tests still pass**

Run:

```bash
cd backend && go test ./internal/interface/... ./internal/application/dto/...
```

Expected: PASS.

- [ ] **Step 5: Check git diff for forbidden scope changes**

Run:

```bash
git diff --stat HEAD
```

Expected: changes limited to backend domain/application/infrastructure code and tests. No migration, router path, response contract, or DTO JSON tag changes.

- [ ] **Step 6: Commit final cleanup**

```bash
git add backend
git commit -m "test: verify backend DDD refactor behavior"
```

---

## Self-Review

- Spec coverage: Tasks cover repository subpackages, application import migration, infrastructure assertions, rich-domain create/update methods, status validity methods, tests, and final API behavior verification.
- Placeholder scan: The plan contains no TBD/TODO/fill-in placeholders. Where existing helper existence is uncertain, the plan provides exact helper code to add.
- Type consistency: Repository package names and domain method names are consistent across tasks. Optional nullable job posting updates use `**string` to preserve the distinction between omitted `jobPostingId` and clearing it to nil.
