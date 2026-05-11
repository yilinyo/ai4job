package job_posting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestJobPosting_Validate_AcceptsFrontendFields(t *testing.T) {
	posting := validPosting()
	posting.CompanyIntro = "Company intro"
	posting.JobRequirements = "Vue and TypeScript"
	posting.Base = "Shenzhen"
	posting.SalaryRange = "25k-35k x 16"
	posting.SourcePlatform = "Boss"
	posting.SourceURL = "https://example.com/job/1"
	posting.Remark = "priority"

	require.NoError(t, posting.Validate())
}

func TestJobPosting_Validate_RejectsOverlongOptionalVarcharFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*domainjob.JobPosting)
	}{
		{name: "base", mutate: func(posting *domainjob.JobPosting) { posting.Base = stringOfRunes(121) }},
		{name: "salaryRange", mutate: func(posting *domainjob.JobPosting) { posting.SalaryRange = stringOfRunes(121) }},
		{name: "sourcePlatform", mutate: func(posting *domainjob.JobPosting) { posting.SourcePlatform = stringOfRunes(121) }},
		{name: "sourceURL", mutate: func(posting *domainjob.JobPosting) { posting.SourceURL = "https://example.com/" + stringOfRunes(481) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posting := validPosting()
			tt.mutate(posting)

			err := posting.Validate()

			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}

func TestJobPosting_Validate_RejectsUnsafeSourceURL(t *testing.T) {
	tests := []string{
		"://example.com/job/1",
		"example.com/job/1",
		"https:///job/1",
		"ftp://example.com/job/1",
		"javascript:alert(1)",
	}

	for _, sourceURL := range tests {
		t.Run(sourceURL, func(t *testing.T) {
			posting := validPosting()
			posting.SourceURL = sourceURL

			err := posting.Validate()

			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}

func validPosting() *domainjob.JobPosting {
	return &domainjob.JobPosting{
		UserID:      "user-1",
		CompanyName: "Acme",
		JobTitle:    "Frontend Engineer",
		Status:      domainjob.StatusPending,
	}
}

func stringOfRunes(length int) string {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = '好'
	}
	return string(runes)
}

func TestJobPosting_Validate_RejectsMissingRequiredFrontendFields(t *testing.T) {
	tests := []struct {
		name    string
		posting domainjob.JobPosting
	}{
		{name: "missing user", posting: domainjob.JobPosting{CompanyName: "Acme", JobTitle: "Frontend", Status: domainjob.StatusPending}},
		{name: "missing companyName", posting: domainjob.JobPosting{UserID: "user-1", JobTitle: "Frontend", Status: domainjob.StatusPending}},
		{name: "missing jobTitle", posting: domainjob.JobPosting{UserID: "user-1", CompanyName: "Acme", Status: domainjob.StatusPending}},
		{name: "invalid status", posting: domainjob.JobPosting{UserID: "user-1", CompanyName: "Acme", JobTitle: "Frontend", Status: "collecting"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.posting.Validate()
			assert.ErrorIs(t, err, domainerrors.ErrBadParams)
		})
	}
}
