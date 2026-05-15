package resume_version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestResumeVersion_ValidateForCreate_ValidFrontendFields(t *testing.T) {
	version := &resume_version.ResumeVersion{
		RepoID:   "repo-1",
		Version:  "v1.0.0",
		Title:    "Frontend version",
		Content:  "# Resume",
		Remark:   "first version",
		IsActive: false,
	}
	assert.NoError(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_MissingVersion(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Title: "Frontend version", Content: "# Resume"}
	assert.ErrorIs(t, version.ValidateForCreate(), domainerrors.ErrBadParams)
}

func TestResumeVersion_ValidateForCreate_EmptyRepoID(t *testing.T) {
	version := &resume_version.ResumeVersion{Version: "v1.0.0", Title: "v1", Content: "{}"}
	assert.Error(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_EmptyContent(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Version: "v1.0.0", Title: "v1", Content: ""}
	assert.Error(t, version.ValidateForCreate())
}

func TestResumeVersion_ValidateForCreate_TitleTooLong(t *testing.T) {
	version := &resume_version.ResumeVersion{RepoID: "repo-1", Version: "v1.0.0", Title: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Content: "{}"}
	assert.Error(t, version.ValidateForCreate())
}
