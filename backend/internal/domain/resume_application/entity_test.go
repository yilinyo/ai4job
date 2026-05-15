package resume_application_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestResumeApplication_Validate_AcceptsFrontendFields(t *testing.T) {
	jobPostingID := "job-1"
	app := validResumeApplication()
	app.JobPostingID = &jobPostingID
	app.ResumeMatchScore = 5
	app.InterviewPerformanceScore = 4
	app.RoleInterestScore = 5
	app.OverallScore = 9

	require.NoError(t, app.Validate())
}

func TestResumeApplication_Validate_AcceptsNilJobPostingID(t *testing.T) {
	app := validResumeApplication()
	app.JobPostingID = nil

	require.NoError(t, app.Validate())
}

func TestResumeApplication_Validate_RejectsInvalidFrontendFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*domainapp.ResumeApplication)
	}{
		{name: "missing user", mutate: func(app *domainapp.ResumeApplication) { app.UserID = "" }},
		{name: "missing version", mutate: func(app *domainapp.ResumeApplication) { app.ResumeVersionID = "" }},
		{name: "missing repo", mutate: func(app *domainapp.ResumeApplication) { app.RepoID = "" }},
		{name: "missing company name", mutate: func(app *domainapp.ResumeApplication) { app.CompanyName = "" }},
		{name: "missing job title", mutate: func(app *domainapp.ResumeApplication) { app.JobTitle = "" }},
		{name: "missing applied at", mutate: func(app *domainapp.ResumeApplication) { app.AppliedAt = time.Time{} }},
		{name: "invalid status", mutate: func(app *domainapp.ResumeApplication) { app.CurrentStatus = "pending" }},
		{name: "invalid score", mutate: func(app *domainapp.ResumeApplication) { app.OverallScore = 11 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := validResumeApplication()
			tt.mutate(app)

			assert.ErrorIs(t, app.Validate(), domainerrors.ErrBadParams)
		})
	}
}

func TestResumeApplication_Validate_RejectsOverlongOptionalVarcharFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*domainapp.ResumeApplication)
	}{
		{name: "base", mutate: func(app *domainapp.ResumeApplication) { app.Base = strings.Repeat("a", 121) }},
		{name: "salary range", mutate: func(app *domainapp.ResumeApplication) { app.SalaryRange = strings.Repeat("a", 121) }},
		{name: "source platform", mutate: func(app *domainapp.ResumeApplication) { app.SourcePlatform = strings.Repeat("a", 121) }},
		{name: "source url", mutate: func(app *domainapp.ResumeApplication) {
			app.SourceURL = "https://example.com/" + strings.Repeat("a", 501)
		}},
		{name: "delivery channel", mutate: func(app *domainapp.ResumeApplication) { app.DeliveryChannel = strings.Repeat("a", 121) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := validResumeApplication()
			tt.mutate(app)

			assert.ErrorIs(t, app.Validate(), domainerrors.ErrBadParams)
		})
	}
}

func TestResumeApplication_Validate_RejectsUnsafeSourceURL(t *testing.T) {
	tests := []struct {
		name      string
		sourceURL string
	}{
		{name: "javascript scheme", sourceURL: "javascript:alert(1)"},
		{name: "non http scheme", sourceURL: "ftp://example.com/job"},
		{name: "missing host", sourceURL: "https:///job"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := validResumeApplication()
			app.SourceURL = tt.sourceURL

			assert.ErrorIs(t, app.Validate(), domainerrors.ErrBadParams)
		})
	}
}

func TestResumeApplication_Validate_RejectsIndividualScoreBounds(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*domainapp.ResumeApplication)
	}{
		{name: "resume match score below minimum", mutate: func(app *domainapp.ResumeApplication) { app.ResumeMatchScore = -1 }},
		{name: "resume match score above maximum", mutate: func(app *domainapp.ResumeApplication) { app.ResumeMatchScore = 6 }},
		{name: "interview performance score below minimum", mutate: func(app *domainapp.ResumeApplication) { app.InterviewPerformanceScore = -1 }},
		{name: "interview performance score above maximum", mutate: func(app *domainapp.ResumeApplication) { app.InterviewPerformanceScore = 6 }},
		{name: "role interest score below minimum", mutate: func(app *domainapp.ResumeApplication) { app.RoleInterestScore = -1 }},
		{name: "role interest score above maximum", mutate: func(app *domainapp.ResumeApplication) { app.RoleInterestScore = 6 }},
		{name: "overall score below minimum", mutate: func(app *domainapp.ResumeApplication) { app.OverallScore = -1 }},
		{name: "overall score above maximum", mutate: func(app *domainapp.ResumeApplication) { app.OverallScore = 11 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := validResumeApplication()
			tt.mutate(app)

			assert.ErrorIs(t, app.Validate(), domainerrors.ErrBadParams)
		})
	}
}

func TestNewResumeApplicationDefaultsAppliedStatus(t *testing.T) {
	appliedAt := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	application, err := domainapp.NewResumeApplication(domainapp.NewResumeApplicationParams{
		ID:              "app-1",
		UserID:          "user-1",
		ResumeVersionID: "version-1",
		RepoID:          "repo-1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       appliedAt,
	})

	require.NoError(t, err)
	assert.Equal(t, domainapp.StatusApplied, application.CurrentStatus)
}

func TestResumeApplicationStatusValid(t *testing.T) {
	assert.True(t, domainapp.StatusApplied.Valid())
	assert.True(t, domainapp.StatusOffer.Valid())
	assert.False(t, domainapp.Status("unknown").Valid())
}

func TestResumeApplicationApplyUpdateValidatesResult(t *testing.T) {
	application := validResumeApplication()
	company := "Beta"
	status := domainapp.StatusViewed

	err := application.ApplyUpdate(domainapp.UpdateSpec{CompanyName: &company, CurrentStatus: &status})

	require.NoError(t, err)
	assert.Equal(t, "Beta", application.CompanyName)
	assert.Equal(t, domainapp.StatusViewed, application.CurrentStatus)
}

func TestResumeApplicationUpdateCurrentStatusValidatesStatus(t *testing.T) {
	application := validResumeApplication()

	err := application.UpdateCurrentStatus(domainapp.StatusOffer)

	require.NoError(t, err)
	assert.Equal(t, domainapp.StatusOffer, application.CurrentStatus)
	assert.Error(t, application.UpdateCurrentStatus(domainapp.Status("unknown")))
}

func validResumeApplication() *domainapp.ResumeApplication {
	return &domainapp.ResumeApplication{
		ID:              "app-1",
		UserID:          "user-1",
		ResumeVersionID: "version-1",
		RepoID:          "repo-1",
		CompanyName:     "Acme",
		JobTitle:        "Backend Engineer",
		AppliedAt:       time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		CurrentStatus:   domainapp.StatusApplied,
	}
}
