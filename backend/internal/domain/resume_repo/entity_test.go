package resume_repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_repo"
)

func TestResumeRepo_Validate_Valid(t *testing.T) {
	repo := &resume_repo.ResumeRepo{UserID: "user-1", Name: "Backend Resume"}
	assert.NoError(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyUserID(t *testing.T) {
	repo := &resume_repo.ResumeRepo{Name: "Backend Resume"}
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_EmptyName(t *testing.T) {
	repo := &resume_repo.ResumeRepo{UserID: "user-1", Name: "   "}
	assert.Error(t, repo.Validate())
}

func TestResumeRepo_Validate_NameTooLong(t *testing.T) {
	repo := &resume_repo.ResumeRepo{UserID: "user-1", Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	assert.Error(t, repo.Validate())
}
