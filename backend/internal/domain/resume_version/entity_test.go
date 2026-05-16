package resume_version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestNewResumeVersionBuildsValidVersion(t *testing.T) {
	version, err := resume_version.NewResumeVersion(resume_version.NewResumeVersionParams{
		ID:       "version-1",
		RepoID:   "repo-1",
		Version:  "v1.0.0",
		Title:    "Backend Resume",
		Content:  "# Resume",
		Remark:   "first",
		IsActive: false,
	})

	require.NoError(t, err)
	assert.Equal(t, "version-1", version.ID)
	assert.Equal(t, "repo-1", version.RepoID)
	assert.False(t, version.IsActive)
}

func TestResumeVersionApplyUpdateHandlesActivationFlag(t *testing.T) {
	version := &resume_version.ResumeVersion{ID: "version-1", RepoID: "repo-1", Version: "v1.0.0", Title: "old", Content: "old", IsActive: true}
	active := false

	err := version.ApplyUpdate(resume_version.UpdateSpec{Title: "new", Content: "# New", Remark: "updated", IsActive: &active})

	require.NoError(t, err)
	assert.Equal(t, "new", version.Title)
	assert.Equal(t, "# New", version.Content)
	assert.Equal(t, "updated", version.Remark)
	assert.False(t, version.IsActive)
}

func TestResumeVersionApplyUpdateLeavesActivationForTrueFlagToRepository(t *testing.T) {
	version := &resume_version.ResumeVersion{ID: "version-1", RepoID: "repo-1", Version: "v1.0.0", Title: "old", Content: "old", IsActive: false}
	active := true

	err := version.ApplyUpdate(resume_version.UpdateSpec{Title: "new", Content: "# New", IsActive: &active})

	require.NoError(t, err)
	assert.False(t, version.IsActive)
}

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
