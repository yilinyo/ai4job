package resume_version_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/resume_version"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type mockResumeRepoReader struct{ mock.Mock }

func (m *mockResumeRepoReader) Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainrepo.ResumeRepo), args.Error(1)
}

func allowRepo(repoReader *mockResumeRepoReader) {
	repoReader.On("Get", mock.Anything, "user-1", "repo-1").Return(&domainrepo.ResumeRepo{ID: "repo-1", UserID: "user-1"}, nil)
}

type mockResumeVersionRepo struct{ mock.Mock }

func (m *mockResumeVersionRepo) Create(ctx context.Context, v *domainversion.ResumeVersion) error {
	return m.Called(ctx, v).Error(0)
}

func (m *mockResumeVersionRepo) FindByIDAndRepoID(ctx context.Context, id, repoID string) (*domainversion.ResumeVersion, error) {
	args := m.Called(ctx, id, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionRepo) ListByRepoID(ctx context.Context, repoID string) ([]domainversion.ResumeVersion, error) {
	args := m.Called(ctx, repoID)
	return args.Get(0).([]domainversion.ResumeVersion), args.Error(1)
}

func (m *mockResumeVersionRepo) Update(ctx context.Context, v *domainversion.ResumeVersion) error {
	return m.Called(ctx, v).Error(0)
}

func (m *mockResumeVersionRepo) Delete(ctx context.Context, id, repoID string) error {
	return m.Called(ctx, id, repoID).Error(0)
}

func (m *mockResumeVersionRepo) MaxVersionNum(ctx context.Context, repoID string) (int, error) {
	args := m.Called(ctx, repoID)
	return args.Int(0), args.Error(1)
}

func (m *mockResumeVersionRepo) ClearDefaultByRepoID(ctx context.Context, repoID string) error {
	return m.Called(ctx, repoID).Error(0)
}

func (m *mockResumeVersionRepo) SetDefaultByID(ctx context.Context, repoID, id string) error {
	return m.Called(ctx, repoID, id).Error(0)
}

type mockVersionUseChecker struct{ mock.Mock }

func (m *mockVersionUseChecker) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	args := m.Called(ctx, versionID)
	return args.Bool(0), args.Error(1)
}

func TestResumeVersionService_Create_AssignsNextVersionNum(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("MaxVersionNum", mock.Anything, "repo-1").Return(2, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool {
		return v.RepoID == "repo-1" && v.Title == "v3" && v.Content == "{}" && v.VersionNum == 3 && v.ID != ""
	})).Return(nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Create(context.Background(), "user-1", "repo-1", dto.CreateResumeVersionRequest{Title: "v3", Content: "{}"})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestResumeVersionService_Create_DefaultPersistsBeforeSwitchingDefault(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("MaxVersionNum", mock.Anything, "repo-1").Return(0, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool {
		return !v.IsDefault && v.VersionNum == 1
	})).Return(nil)
	repo.On("SetDefaultByID", mock.Anything, "repo-1", mock.AnythingOfType("string")).Return(nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Create(context.Background(), "user-1", "repo-1", dto.CreateResumeVersionRequest{Title: "v1", Content: "{}", IsDefault: true})

	require.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "ClearDefaultByRepoID")
}

func TestResumeVersionService_List_ChecksRepoOwnership(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	rows := []domainversion.ResumeVersion{{ID: "v1", RepoID: "repo-1", Title: "v1", Content: "{}", VersionNum: 1}}
	repoReader.On("Get", mock.Anything, "user-1", "repo-1").Return(&domainrepo.ResumeRepo{ID: "repo-1", UserID: "user-1"}, nil)
	repo.On("ListByRepoID", mock.Anything, "repo-1").Return(rows, nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	got, err := svc.List(context.Background(), "user-1", "repo-1")

	require.NoError(t, err)
	assert.Equal(t, rows, got)
	repoReader.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestResumeVersionService_List_RejectsOtherUsersRepo(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	repoReader.On("Get", mock.Anything, "user-1", "repo-1").Return(nil, domainerrors.ErrNotFound)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	got, err := svc.List(context.Background(), "user-1", "repo-1")

	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	assert.Nil(t, got)
	repo.AssertNotCalled(t, "ListByRepoID")
}

func TestResumeVersionService_Get(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	existing := &domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Title: "v1", Content: "{}", VersionNum: 1}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(existing, nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	got, err := svc.Get(context.Background(), "user-1", "repo-1", "v1")

	require.NoError(t, err)
	assert.Equal(t, existing, got)
	repo.AssertExpectations(t)
}

func TestResumeVersionService_Update_Success(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	existing := &domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Title: "old", Content: "{}", VersionNum: 1}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(v *domainversion.ResumeVersion) bool {
		return v.ID == "v1" && v.RepoID == "repo-1" && v.Title == "new" && v.Content == "{updated}" && v.VersionNum == 1
	})).Return(nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Update(context.Background(), "user-1", "repo-1", "v1", dto.UpdateResumeVersionRequest{Title: "new", Content: "{updated}"})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestResumeVersionService_Update_InvalidTitle(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	existing := &domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Title: "old", Content: "{}", VersionNum: 1}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(existing, nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Update(context.Background(), "user-1", "repo-1", "v1", dto.UpdateResumeVersionRequest{Title: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Content: "{}"})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "Update")
}

func TestResumeVersionService_Delete_InUse(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(&domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Title: "v1", Content: "{}", VersionNum: 1}, nil)
	checker.On("ExistsByResumeVersionID", mock.Anything, "v1").Return(true, nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Delete(context.Background(), "user-1", "repo-1", "v1")

	assert.ErrorIs(t, err, domainerrors.ErrVersionInUse)
	repo.AssertNotCalled(t, "Delete")
}

func TestResumeVersionService_Delete_NotFoundBeforeUseCheck(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(nil, domainerrors.ErrNotFound)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Delete(context.Background(), "user-1", "repo-1", "v1")

	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	checker.AssertNotCalled(t, "ExistsByResumeVersionID")
	repo.AssertNotCalled(t, "Delete")
}

func TestResumeVersionService_Delete_Success(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("FindByIDAndRepoID", mock.Anything, "v1", "repo-1").Return(&domainversion.ResumeVersion{ID: "v1", RepoID: "repo-1", Title: "v1", Content: "{}", VersionNum: 1}, nil)
	checker.On("ExistsByResumeVersionID", mock.Anything, "v1").Return(false, nil)
	repo.On("Delete", mock.Anything, "v1", "repo-1").Return(nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.Delete(context.Background(), "user-1", "repo-1", "v1")

	require.NoError(t, err)
	repo.AssertExpectations(t)
	checker.AssertExpectations(t)
}

func TestResumeVersionService_SetDefault(t *testing.T) {
	repoReader := &mockResumeRepoReader{}
	repo := &mockResumeVersionRepo{}
	checker := &mockVersionUseChecker{}
	allowRepo(repoReader)
	repo.On("SetDefaultByID", mock.Anything, "repo-1", "v1").Return(nil)
	svc := appsvc.NewResumeVersionService(repoReader, repo, checker)

	err := svc.SetDefault(context.Background(), "user-1", "repo-1", "v1")

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
