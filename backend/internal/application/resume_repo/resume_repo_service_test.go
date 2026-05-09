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

func TestResumeRepoService_Create_DuplicateName(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(true, nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Create(context.Background(), "u1", dto.CreateResumeRepoRequest{Name: "Main"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicate)
	repo.AssertExpectations(t)
}

func TestResumeRepoService_Create_Success(t *testing.T) {
	repo := &mockResumeRepoRepo{}
	repo.On("ExistsByUserAndName", mock.Anything, "u1", "Main", "").Return(false, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(r *domainrepo.ResumeRepo) bool {
		return r.UserID == "u1" && r.Name == "Main" && r.ID != ""
	})).Return(nil)
	svc := appsvc.NewResumeRepoService(repo)

	err := svc.Create(context.Background(), "u1", dto.CreateResumeRepoRequest{Name: "Main"})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
