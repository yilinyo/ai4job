package resume_application_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_application"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	resumeapplicationrepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeApplicationRepo struct{ mock.Mock }

func (m *mockResumeApplicationRepo) Create(ctx context.Context, a *domainapp.ResumeApplication) error {
	return m.Called(ctx, a).Error(0)
}

func (m *mockResumeApplicationRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainapp.ResumeApplication, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainapp.ResumeApplication), args.Error(1)
}

func (m *mockResumeApplicationRepo) Update(ctx context.Context, a *domainapp.ResumeApplication) error {
	return m.Called(ctx, a).Error(0)
}

func (m *mockResumeApplicationRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockResumeApplicationRepo) ListByUser(ctx context.Context, f resumeapplicationrepo.ListFilter) ([]domainapp.ResumeApplication, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domainapp.ResumeApplication), args.Get(1).(int64), args.Error(2)
}

func (m *mockResumeApplicationRepo) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	args := m.Called(ctx, versionID)
	return args.Bool(0), args.Error(1)
}

func (m *mockResumeApplicationRepo) UpdateCurrentStatus(ctx context.Context, id, userID string, status domainapp.Status) error {
	return m.Called(ctx, id, userID, status).Error(0)
}

type mockJobPostingReader struct{ mock.Mock }

func (m *mockJobPostingReader) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) {
	args := m.Called(ctx, id, userID)
	posting, _ := args.Get(0).(*domainjob.JobPosting)
	return posting, args.Error(1)
}

type mockResumeRepoReader struct{ mock.Mock }

func (m *mockResumeRepoReader) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) {
	args := m.Called(ctx, id, userID)
	repo, _ := args.Get(0).(*domainrepo.ResumeRepo)
	return repo, args.Error(1)
}

type mockResumeVersionReader struct{ mock.Mock }

func (m *mockResumeVersionReader) FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, id)
	version, _ := args.Get(0).(*domainversion.ResumeVersion)
	return version, args.Error(1)
}

func newOwnedReferenceReaders(jobPostingID *string, repoID, versionID string) (*mockJobPostingReader, *mockResumeRepoReader, *mockResumeVersionReader) {
	jobs := &mockJobPostingReader{}
	repos := &mockResumeRepoReader{}
	versions := &mockResumeVersionReader{}
	if jobPostingID != nil {
		jobs.On("FindByIDAndUserID", mock.Anything, *jobPostingID, "u1").Return(&domainjob.JobPosting{ID: *jobPostingID, UserID: "u1"}, nil)
	}
	repos.On("FindByIDAndUserID", mock.Anything, repoID, "u1").Return(&domainrepo.ResumeRepo{ID: repoID, UserID: "u1"}, nil)
	versions.On("FindByID", mock.Anything, versionID).Return(&domainversion.ResumeVersion{ID: versionID, RepoID: repoID}, nil)
	return jobs, repos, versions
}

func TestResumeApplicationService_CreateCopiesFieldsDefaultsStatusAndSupportsOptionalJobPostingID(t *testing.T) {
	appliedAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)

	for _, tt := range []struct {
		name         string
		jobPostingID *string
	}{
		{name: "nil job posting id", jobPostingID: nil},
		{name: "provided job posting id", jobPostingID: ptr("job-1")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockResumeApplicationRepo{}
			jobs, repos, versions := newOwnedReferenceReaders(tt.jobPostingID, "r1", "v1")
			repo.On("Create", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
				return a.ID != "" &&
					a.UserID == "u1" &&
					stringPtrEqual(a.JobPostingID, tt.jobPostingID) &&
					a.ResumeVersionID == "v1" &&
					a.RepoID == "r1" &&
					a.CompanyName == "Acme" &&
					a.CompanyIntro == "intro" &&
					a.JobTitle == "Backend Engineer" &&
					a.JobRequirements == "Go" &&
					a.Base == "Shanghai" &&
					a.SalaryRange == "20k-30k" &&
					a.SourcePlatform == "Boss" &&
					a.SourceURL == "https://www.zhipin.com/job/1" &&
					a.DeliveryChannel == "manual" &&
					a.AppliedAt.Equal(appliedAt) &&
					a.CurrentStatus == domainapp.StatusApplied &&
					a.InterviewSummary == "summary" &&
					a.InterviewNotes == "notes" &&
					a.ResumeMatchScore == 5 &&
					a.InterviewPerformanceScore == 4 &&
					a.RoleInterestScore == 3 &&
					a.OverallScore == 9 &&
					a.Remark == "remark"
			})).Return(nil)
			svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

			got, err := svc.Create(context.Background(), "u1", dto.CreateResumeApplicationRequest{
				JobPostingID:              tt.jobPostingID,
				ResumeVersionID:           "v1",
				RepoID:                    "r1",
				CompanyName:               "Acme",
				CompanyIntro:              "intro",
				JobTitle:                  "Backend Engineer",
				JobRequirements:           "Go",
				Base:                      "Shanghai",
				SalaryRange:               "20k-30k",
				SourcePlatform:            "Boss",
				SourceURL:                 "https://www.zhipin.com/job/1",
				DeliveryChannel:           "manual",
				AppliedAt:                 dto.FlexibleTime{Time: appliedAt},
				InterviewSummary:          "summary",
				InterviewNotes:            "notes",
				ResumeMatchScore:          5,
				InterviewPerformanceScore: 4,
				RoleInterestScore:         3,
				OverallScore:              9,
				Remark:                    "remark",
			})

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "u1", got.UserID)
			assert.Equal(t, "Acme", got.CompanyName)
			repo.AssertExpectations(t)
			jobs.AssertExpectations(t)
			repos.AssertExpectations(t)
			versions.AssertExpectations(t)
		})
	}
}

func TestResumeApplicationService_CreateNormalizesBlankJobPostingIDToNil(t *testing.T) {
	appliedAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	blankJobPostingID := "  \t  "
	repo := &mockResumeApplicationRepo{}
	repo.On("Create", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
		return a.JobPostingID == nil
	})).Return(nil)
	jobs, repos, versions := newOwnedReferenceReaders(nil, "r1", "v1")
	svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

	got, err := svc.Create(context.Background(), "u1", dto.CreateResumeApplicationRequest{
		JobPostingID:    &blankJobPostingID,
		ResumeVersionID: "v1",
		RepoID:          "r1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       dto.FlexibleTime{Time: appliedAt},
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Nil(t, got.JobPostingID)
	repo.AssertExpectations(t)
	jobs.AssertExpectations(t)
	repos.AssertExpectations(t)
	versions.AssertExpectations(t)
}

func TestResumeApplicationService_CreateRejectsUnownedOrMissingReferencesBeforeCreate(t *testing.T) {
	appliedAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	jobPostingID := "job-1"

	for _, tt := range []struct {
		name  string
		setup func(*mockJobPostingReader, *mockResumeRepoReader, *mockResumeVersionReader)
	}{
		{
			name: "job posting",
			setup: func(jobs *mockJobPostingReader, repos *mockResumeRepoReader, versions *mockResumeVersionReader) {
				jobs.On("FindByIDAndUserID", mock.Anything, "job-1", "u1").Return(nil, domainerrors.ErrNotFound)
			},
		},
		{
			name: "repo",
			setup: func(jobs *mockJobPostingReader, repos *mockResumeRepoReader, versions *mockResumeVersionReader) {
				jobs.On("FindByIDAndUserID", mock.Anything, "job-1", "u1").Return(&domainjob.JobPosting{ID: "job-1", UserID: "u1"}, nil)
				repos.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(nil, domainerrors.ErrNotFound)
			},
		},
		{
			name: "resume version",
			setup: func(jobs *mockJobPostingReader, repos *mockResumeRepoReader, versions *mockResumeVersionReader) {
				jobs.On("FindByIDAndUserID", mock.Anything, "job-1", "u1").Return(&domainjob.JobPosting{ID: "job-1", UserID: "u1"}, nil)
				repos.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(&domainrepo.ResumeRepo{ID: "r1", UserID: "u1"}, nil)
				versions.On("FindByID", mock.Anything, "v1").Return(nil, domainerrors.ErrNotFound)
			},
		},
		{
			name: "resume version repo",
			setup: func(jobs *mockJobPostingReader, repos *mockResumeRepoReader, versions *mockResumeVersionReader) {
				jobs.On("FindByIDAndUserID", mock.Anything, "job-1", "u1").Return(&domainjob.JobPosting{ID: "job-1", UserID: "u1"}, nil)
				repos.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(&domainrepo.ResumeRepo{ID: "r1", UserID: "u1"}, nil)
				versions.On("FindByID", mock.Anything, "v1").Return(&domainversion.ResumeVersion{ID: "v1", RepoID: "other-repo"}, nil)
				repos.On("FindByIDAndUserID", mock.Anything, "other-repo", "u1").Return(nil, domainerrors.ErrNotFound)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockResumeApplicationRepo{}
			jobs := &mockJobPostingReader{}
			repos := &mockResumeRepoReader{}
			versions := &mockResumeVersionReader{}
			tt.setup(jobs, repos, versions)
			svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

			got, err := svc.Create(context.Background(), "u1", dto.CreateResumeApplicationRequest{
				JobPostingID:    &jobPostingID,
				ResumeVersionID: "v1",
				RepoID:          "r1",
				CompanyName:     "Acme",
				JobTitle:        "Backend Engineer",
				AppliedAt:       dto.FlexibleTime{Time: appliedAt},
			})

			assert.Nil(t, got)
			assert.Error(t, err)
			repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			jobs.AssertExpectations(t)
			repos.AssertExpectations(t)
			versions.AssertExpectations(t)
		})
	}
}

func TestResumeApplicationService_ListUsesFiltersAndDefaultsPagination(t *testing.T) {
	repo := &mockResumeApplicationRepo{}
	rows := []domainapp.ResumeApplication{{ID: "a1", UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Backend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}}
	repo.On("ListByUser", mock.Anything, resumeapplicationrepo.ListFilter{
		UserID:          "u1",
		ResumeVersionID: "v1",
		RepoID:          "r1",
		Keyword:         "backend",
		CurrentStatus:   domainapp.StatusApplied,
		Page:            1,
		PageSize:        10,
	}).Return(rows, int64(1), nil)
	svc := appsvc.NewResumeApplicationService(repo, nil, nil, nil)

	got, total, err := svc.List(context.Background(), "u1", "v1", "r1", "backend", string(domainapp.StatusApplied), 0, 0)

	require.NoError(t, err)
	assert.Equal(t, rows, got)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestResumeApplicationService_ListRejectsInvalidStatusWithoutCallingRepo(t *testing.T) {
	repo := &mockResumeApplicationRepo{}
	svc := appsvc.NewResumeApplicationService(repo, nil, nil, nil)

	got, total, err := svc.List(context.Background(), "u1", "v1", "r1", "backend", "invalid", 1, 10)

	assert.Nil(t, got)
	assert.Zero(t, total)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything)
}

func TestResumeApplicationService_UpdatePreservesOmittedFieldsAndAppliesProvidedFields(t *testing.T) {
	appliedAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	newAppliedAt := appliedAt.Add(24 * time.Hour)
	existingJobPostingID := "job-old"
	newJobPostingID := "job-new"
	existing := &domainapp.ResumeApplication{
		ID:                        "a1",
		UserID:                    "u1",
		JobPostingID:              &existingJobPostingID,
		ResumeVersionID:           "v1",
		RepoID:                    "r1",
		CompanyName:               "Old Co",
		CompanyIntro:              "old intro",
		JobTitle:                  "Old Title",
		JobRequirements:           "old req",
		Base:                      "Beijing",
		SalaryRange:               "10k-20k",
		SourcePlatform:            "官网",
		SourceURL:                 "https://jobs.bytedance.com/old",
		DeliveryChannel:           "old channel",
		AppliedAt:                 appliedAt,
		CurrentStatus:             domainapp.StatusApplied,
		InterviewSummary:          "old summary",
		InterviewNotes:            "old notes",
		ResumeMatchScore:          3,
		InterviewPerformanceScore: 2,
		RoleInterestScore:         1,
		OverallScore:              6,
		Remark:                    "old remark",
	}
	newTitle := "New Title"
	newURL := "https://www.zhipin.com/job/2"
	newStatus := domainapp.StatusFirstInterview
	newOverallScore := 8
	repo := &mockResumeApplicationRepo{}
	repo.On("FindByIDAndUserID", mock.Anything, "a1", "u1").Return(existing, nil)
	jobs, repos, versions := newOwnedReferenceReaders(&newJobPostingID, "r1", "v1")
	repo.On("Update", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
		return a.ID == "a1" &&
			a.UserID == "u1" &&
			a.JobPostingID != nil && *a.JobPostingID == newJobPostingID &&
			a.ResumeVersionID == "v1" &&
			a.RepoID == "r1" &&
			a.CompanyName == "Old Co" &&
			a.CompanyIntro == "old intro" &&
			a.JobTitle == newTitle &&
			a.JobRequirements == "old req" &&
			a.Base == "Beijing" &&
			a.SalaryRange == "10k-20k" &&
			a.SourcePlatform == "官网" &&
			a.SourceURL == newURL &&
			a.DeliveryChannel == "old channel" &&
			a.AppliedAt.Equal(newAppliedAt) &&
			a.CurrentStatus == newStatus &&
			a.InterviewSummary == "old summary" &&
			a.InterviewNotes == "old notes" &&
			a.ResumeMatchScore == 3 &&
			a.InterviewPerformanceScore == 2 &&
			a.RoleInterestScore == 1 &&
			a.OverallScore == newOverallScore &&
			a.Remark == "old remark"
	})).Return(nil)
	svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

	err := svc.Update(context.Background(), "u1", "a1", dto.UpdateResumeApplicationRequest{
		JobPostingID:  &newJobPostingID,
		JobTitle:      &newTitle,
		SourceURL:     &newURL,
		AppliedAt:     &dto.FlexibleTime{Time: newAppliedAt},
		CurrentStatus: &newStatus,
		OverallScore:  &newOverallScore,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	jobs.AssertExpectations(t)
	repos.AssertExpectations(t)
	versions.AssertExpectations(t)
}

func TestResumeApplicationService_UpdateNormalizesBlankJobPostingIDToNil(t *testing.T) {
	existing := &domainapp.ResumeApplication{
		ID:              "a1",
		UserID:          "u1",
		ResumeVersionID: "v1",
		RepoID:          "r1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC),
		CurrentStatus:   domainapp.StatusApplied,
	}
	blankJobPostingID := " \n\t "
	repo := &mockResumeApplicationRepo{}
	repo.On("FindByIDAndUserID", mock.Anything, "a1", "u1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
		return a.JobPostingID == nil
	})).Return(nil)
	jobs, repos, versions := newOwnedReferenceReaders(nil, "r1", "v1")
	svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

	err := svc.Update(context.Background(), "u1", "a1", dto.UpdateResumeApplicationRequest{
		JobPostingID: &blankJobPostingID,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	jobs.AssertExpectations(t)
	repos.AssertExpectations(t)
	versions.AssertExpectations(t)
}

func TestResumeApplicationService_UpdateClearsExistingJobPostingIDWhenProvidedBlank(t *testing.T) {
	existingJobPostingID := "job-old"
	existing := &domainapp.ResumeApplication{
		ID:              "a1",
		UserID:          "u1",
		JobPostingID:    &existingJobPostingID,
		ResumeVersionID: "v1",
		RepoID:          "r1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC),
		CurrentStatus:   domainapp.StatusApplied,
	}
	blankJobPostingID := ""
	repo := &mockResumeApplicationRepo{}
	repo.On("FindByIDAndUserID", mock.Anything, "a1", "u1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(a *domainapp.ResumeApplication) bool {
		return a.JobPostingID == nil
	})).Return(nil)
	jobs, repos, versions := newOwnedReferenceReaders(nil, "r1", "v1")
	svc := appsvc.NewResumeApplicationService(repo, jobs, repos, versions)

	err := svc.Update(context.Background(), "u1", "a1", dto.UpdateResumeApplicationRequest{
		JobPostingID: &blankJobPostingID,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	jobs.AssertExpectations(t)
	repos.AssertExpectations(t)
	versions.AssertExpectations(t)
}

func TestResumeApplicationService_UpdateCurrentStatusRejectsInvalidStatusAndDelegatesValidStatus(t *testing.T) {
	repo := &mockResumeApplicationRepo{}
	svc := appsvc.NewResumeApplicationService(repo, nil, nil, nil)

	err := svc.UpdateCurrentStatus(context.Background(), "u1", "a1", domainapp.Status("invalid"))
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "UpdateCurrentStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	repo.On("UpdateCurrentStatus", mock.Anything, "a1", "u1", domainapp.StatusOffer).Return(nil)
	err = svc.UpdateCurrentStatus(context.Background(), "u1", "a1", domainapp.StatusOffer)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func ptr(v string) *string {
	return &v
}

func stringPtrEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
