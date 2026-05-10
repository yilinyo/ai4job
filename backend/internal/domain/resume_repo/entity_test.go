package resume_repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
