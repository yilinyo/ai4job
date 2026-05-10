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

func (m *mockResumeRepoRepo) Create(ctx context.Context, r *domainrepo.ResumeRepo) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockResumeRepoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainrepo.ResumeRepo), args.Error(1)
}

func (m *mockResumeRepoRepo) Update(ctx context.Context, r *domainrepo.ResumeRepo) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockResumeRepoRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockResumeRepoRepo) ListByUser(ctx context.Context, f domainrepo.ListFilter) ([]domainrepo.ResumeRepo, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domainrepo.ResumeRepo), args.Get(1).(int64), args.Error(2)
}

func (m *mockResumeRepoRepo) ExistsByUserAndName(ctx context.Context, userID, name, excludeID string) (bool, error) {
	args := m.Called(ctx, userID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *mockResumeRepoRepo) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return m.Called(ctx, id, delta).Error(0)
}

func validCreateResumeRepoRequest(name string) dto.CreateResumeRepoRequest {
	return dto.CreateResumeRepoRequest{
		Name:           name,
		JobType:        domainrepo.JobTypeCampus,
		TargetPosition: "Frontend Engineer",
		Description:    "desc",
	}
}

func validUpdateResumeRepoRequest(name string) dto.UpdateResumeRepoRequest {
	return dto.UpdateResumeRepoRequest{
		Name:           name,
		JobType:        domainrepo.JobTypeSocial,
		TargetPosition: "Fullstack Engineer",
		Description:    "updated",
	}
}

func validResumeRepo(id, userID, name string) *domainrepo.ResumeRepo {
	return &domainrepo.ResumeRepo{
		ID:             id,
		UserID:         userID,
		Name:           name,
		JobType:        domainrepo.JobTypeCampus,
		TargetPosition: "Frontend Engineer",
	}
}

func TestResumeRepoService_Create_DuplicateName(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(true, nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Create(context.Background(), "u1", validCreateResumeRepoRequest("Main"))

	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Create_Success(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(false, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(r *domainrepo.ResumeRepo) bool {
		return r.UserID == "u1" && r.Name == "Main" && r.JobType == domainrepo.JobTypeCampus && r.TargetPosition == "Frontend Engineer" && r.VersionCount == 0 && r.ID != ""
	})).Return(nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Create(context.Background(), "u1", validCreateResumeRepoRequest("Main"))

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Create_InvalidName(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "").Return(false, nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Create(context.Background(), "u1", validCreateResumeRepoRequest("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "Create")
}

func TestResumeRepoService_List(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	rows := []domainrepo.ResumeRepo{*validResumeRepo("r1", "u1", "Main")}
	repo.On("ListByUser", mock.Anything, domainrepo.ListFilter{UserID: "u1", JobType: domainrepo.JobTypeCampus, Keyword: "back", Page: 2, PageSize: 20}).Return(rows, int64(1), nil)
	svc := appsvc.NewResumeRepoService(repo)

	got, total, err := svc.List(context.Background(), "u1", domainrepo.JobTypeCampus, "back", 2, 20)

	require.NoError(t, err)
	assert.Equal(t, rows, got)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Get(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	existing := validResumeRepo("r1", "u1", "Main")
	repo.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(existing, nil)
	svc := appsvc.NewResumeRepoService(repo)

	got, err := svc.Get(context.Background(), "u1", "r1")

	require.NoError(t, err)
	assert.Equal(t, existing, got)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Update_DuplicateName(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	existing := validResumeRepo("r1", "u1", "Old")
	repo.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(existing, nil)
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "r1").Return(true, nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Update(context.Background(), "u1", "r1", validUpdateResumeRepoRequest("Main"))

	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
	repo.AssertNotCalled(t, "Update")
}

func TestResumeRepoService_Update_Success(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	existing := validResumeRepo("r1", "u1", "Old")
	repo.On("FindByIDAndUserID", mock.Anything, "r1", "u1").Return(existing, nil)
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "r1").Return(false, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *domainrepo.ResumeRepo) bool {
		return r.ID == "r1" && r.UserID == "u1" && r.Name == "Main" && r.JobType == domainrepo.JobTypeSocial && r.TargetPosition == "Fullstack Engineer" && r.Description == "updated"
	})).Return(nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Update(context.Background(), "u1", "r1", validUpdateResumeRepoRequest("Main"))

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Delete(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("Delete", mock.Anything, "r1", "u1").Return(nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Delete(context.Background(), "u1", "r1")

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
