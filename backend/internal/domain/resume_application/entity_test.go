package resume_application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestResumeApplication_Validate_AcceptsFrontendFields(t *testing.T) {
	app := &domainapp.ResumeApplication{
		UserID:                    "user-1",
		JobPostingID:              "job-1",
		ResumeVersionID:           "version-1",
		RepoID:                    "repo-1",
		CompanyName:               "Acme",
		JobTitle:                  "Frontend Engineer",
		AppliedAt:                 time.Now(),
		CurrentStatus:             domainapp.StatusApplied,
		ResumeMatchScore:          5,
		InterviewPerformanceScore: 4,
		RoleInterestScore:         5,
		OverallScore:              9,
	}

	require.NoError(t, app.Validate())
}

func TestResumeApplication_Validate_RejectsInvalidFrontendFields(t *testing.T) {
	tests := []struct {
		name string
		app  domainapp.ResumeApplication
	}{
		{name: "missing user", app: domainapp.ResumeApplication{ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "missing version", app: domainapp.ResumeApplication{UserID: "u1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "missing repo", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied}},
		{name: "invalid status", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: "pending"}},
		{name: "invalid score", app: domainapp.ResumeApplication{UserID: "u1", ResumeVersionID: "v1", RepoID: "r1", CompanyName: "Acme", JobTitle: "Frontend", AppliedAt: time.Now(), CurrentStatus: domainapp.StatusApplied, OverallScore: 11}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.app.Validate(), domainerrors.ErrBadParams)
		})
	}
}
