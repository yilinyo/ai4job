package resume_version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

func TestResumeVersion_ValidateForCreate_Valid(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Title: "v1", Content: "{}"}
	assert.NoError(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_EmptyRepoID(t *testing.T) {
	version := &resume_version.ResumeVersion{Title: "v1", Content: "{}"}
	assert.Error(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_EmptyContent(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Title: "v1", Content: ""}
	assert.Error(t, version.ValidateForCreate())
}
