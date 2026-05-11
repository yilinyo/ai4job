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
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockJobPostingRepo struct{ mock.Mock }

func (m *mockJobPostingRepo) Create(ctx context.Context, p *domainjob.JobPosting) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockJobPostingRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainjob.JobPosting), args.Error(1)
}

func (m *mockJobPostingRepo) Update(ctx context.Context, p *domainjob.JobPosting) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockJobPostingRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockJobPostingRepo) ListByUser(ctx context.Context, f domainjob.ListFilter) ([]domainjob.JobPosting, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domainjob.JobPosting), args.Get(1).(int64), args.Error(2)
}

func TestJobPostingService_CreateCopiesFieldsAndDefaultsStatus(t *testing.T) {
	repo := &mockJobPostingRepo{}
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domainjob.JobPosting) bool {
		return p.ID != "" &&
			p.UserID == "u1" &&
			p.CompanyName == "Acme" &&
			p.CompanyIntro == "intro" &&
			p.JobTitle == "Backend Engineer" &&
			p.JobRequirements == "Go" &&
			p.Base == "Shanghai" &&
			p.SalaryRange == "20k-30k" &&
			p.SourcePlatform == "Boss" &&
			p.SourceURL == "https://www.zhipin.com/job/1" &&
			p.Status == domainjob.StatusPending &&
			p.Remark == "note"
	})).Return(nil)
	svc := appsvc.NewJobPostingService(repo)

	err := svc.Create(context.Background(), "u1", dto.CreateJobPostingRequest{
		CompanyName:     "Acme",
		CompanyIntro:    "intro",
		JobTitle:        "Backend Engineer",
		JobRequirements: "Go",
		Base:            "Shanghai",
		SalaryRange:     "20k-30k",
		SourcePlatform:  "Boss",
		SourceURL:       "https://www.zhipin.com/job/1",
		Remark:          "note",
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestJobPostingService_ListUsesFiltersAndDefaultsPagination(t *testing.T) {
	repo := &mockJobPostingRepo{}
	rows := []domainjob.JobPosting{{ID: "j1", UserID: "u1", CompanyName: "Acme", JobTitle: "Backend", Status: domainjob.StatusPending}}
	repo.On("ListByUser", mock.Anything, domainjob.ListFilter{
		UserID:         "u1",
		Keyword:        "backend",
		SourcePlatform: "Boss",
		Status:         domainjob.StatusPending,
		Page:           1,
		PageSize:       10,
	}).Return(rows, int64(1), nil)
	svc := appsvc.NewJobPostingService(repo)

	got, total, err := svc.List(context.Background(), "u1", "backend", "Boss", string(domainjob.StatusPending), 0, 0)

	require.NoError(t, err)
	assert.Equal(t, rows, got)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestJobPostingService_ListRejectsInvalidStatusWithoutCallingRepo(t *testing.T) {
	repo := &mockJobPostingRepo{}
	svc := appsvc.NewJobPostingService(repo)

	got, total, err := svc.List(context.Background(), "u1", "backend", "Boss", "invalid", 1, 10)

	assert.Nil(t, got)
	assert.Zero(t, total)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything)
}

func TestJobPostingService_UpdatePreservesOmittedFieldsAndAppliesProvidedFields(t *testing.T) {
	repo := &mockJobPostingRepo{}
	existing := &domainjob.JobPosting{
		ID:              "j1",
		UserID:          "u1",
		CompanyName:     "Old Co",
		CompanyIntro:    "old intro",
		JobTitle:        "Old Title",
		JobRequirements: "old req",
		Base:            "Beijing",
		SalaryRange:     "10k-20k",
		SourcePlatform:  "官网",
		SourceURL:       "https://jobs.bytedance.com/old",
		Status:          domainjob.StatusPending,
		Remark:          "old remark",
	}
	newTitle := "New Title"
	newURL := "https://www.zhipin.com/job/2"
	newStatus := domainjob.StatusApplied
	repo.On("FindByIDAndUserID", mock.Anything, "j1", "u1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domainjob.JobPosting) bool {
		return p.ID == "j1" &&
			p.UserID == "u1" &&
			p.CompanyName == "Old Co" &&
			p.CompanyIntro == "old intro" &&
			p.JobTitle == newTitle &&
			p.JobRequirements == "old req" &&
			p.Base == "Beijing" &&
			p.SalaryRange == "10k-20k" &&
			p.SourcePlatform == "官网" &&
			p.SourceURL == newURL &&
			p.Status == newStatus &&
			p.Remark == "old remark"
	})).Return(nil)
	svc := appsvc.NewJobPostingService(repo)

	err := svc.Update(context.Background(), "u1", "j1", dto.UpdateJobPostingRequest{
		JobTitle:  &newTitle,
		SourceURL: &newURL,
		Status:    &newStatus,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestJobPostingService_ParseReturnsPrefilledFields(t *testing.T) {
	svc := appsvc.NewJobPostingService(&mockJobPostingRepo{})

	cases := []struct {
		name           string
		rawURL         string
		sourcePlatform string
		companyName    string
		jobTitle       string
		base           string
	}{
		{
			name:           "boss",
			rawURL:         "https://www.zhipin.com/job/1",
			sourcePlatform: "Boss",
			companyName:    "腾讯科技",
			jobTitle:       "前端开发工程师",
			base:           "深圳",
		},
		{
			name:           "bytedance jobs",
			rawURL:         "https://jobs.bytedance.com/referral/pc/position/1",
			sourcePlatform: "官网",
			companyName:    "示例公司",
			jobTitle:       "待确认岗位",
			base:           "北京",
		},
		{
			name:           "other",
			rawURL:         "https://example.com/jobs/1",
			sourcePlatform: "其他",
			companyName:    "示例公司",
			jobTitle:       "待确认岗位",
			base:           "北京",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Parse(context.Background(), dto.ParseJobPostingRequest{URL: tt.rawURL})
			require.NoError(t, err)

			assert.Equal(t, tt.companyName, got.CompanyName)
			assert.Equal(t, "通过解析接口返回的公司介绍示例", got.CompanyIntro)
			assert.Equal(t, tt.jobTitle, got.JobTitle)
			assert.Equal(t, "通过解析接口返回的岗位要求示例，用户可继续修改", got.JobRequirements)
			assert.Equal(t, tt.base, got.Base)
			assert.Equal(t, "25k-35k x 16", got.SalaryRange)
			assert.Equal(t, tt.sourcePlatform, got.SourcePlatform)
			assert.Equal(t, tt.rawURL, got.SourceURL)
		})
	}
}

func TestJobPostingService_ParseRejectsInvalidURLs(t *testing.T) {
	svc := appsvc.NewJobPostingService(&mockJobPostingRepo{})

	for _, rawURL := range []string{"", "javascript:alert(1)", "not-a-url", "https://"} {
		t.Run(rawURL, func(t *testing.T) {
			got, err := svc.Parse(context.Background(), dto.ParseJobPostingRequest{URL: rawURL})

			assert.Nil(t, got)
			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}
