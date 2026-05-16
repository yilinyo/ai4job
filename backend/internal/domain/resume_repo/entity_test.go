package resume_repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/domain/resume_repo"
)

func validResumeRepo() *resume_repo.ResumeRepo {
	return &resume_repo.ResumeRepo{
		UserID:         "user-1",
		Name:           "Backend Resume",
		JobType:        resume_repo.JobTypeCampus,
		TargetPosition: "Backend Engineer",
		VersionCount:   0,
	}
}

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

func TestResumeRepo_Validate_Valid(t *testing.T) {
	repo := validResumeRepo()
	assert.NoError(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyUserID(t *testing.T) {
	repo := validResumeRepo()
	repo.UserID = ""
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyName(t *testing.T) {
	repo := validResumeRepo()
	repo.Name = "   "
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_NameTooLong(t *testing.T) {
	repo := validResumeRepo()
	repo.Name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyJobType(t *testing.T) {
	repo := validResumeRepo()
	repo.JobType = ""
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_InvalidJobType(t *testing.T) {
	repo := validResumeRepo()
	repo.JobType = "contract"
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyTargetPosition(t *testing.T) {
	repo := validResumeRepo()
	repo.TargetPosition = "   "
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_TargetPositionTooLong(t *testing.T) {
	repo := validResumeRepo()
	repo.TargetPosition = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_NegativeVersionCount(t *testing.T) {
	repo := validResumeRepo()
	repo.VersionCount = -1
	assert.Error(t, repo.Validate())
}
