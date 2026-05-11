package job_posting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestJobPosting_Validate_AcceptsFrontendFields(t *testing.T) {
	posting := &domainjob.JobPosting{
		UserID:          "user-1",
		CompanyName:     "Acme",
		CompanyIntro:    "Company intro",
		JobTitle:        "Frontend Engineer",
		JobRequirements: "Vue and TypeScript",
		Base:            "Shenzhen",
		SalaryRange:     "25k-35k x 16",
		SourcePlatform:  "Boss",
		SourceURL:       "https://example.com/job/1",
		Status:          domainjob.StatusPending,
		Remark:          "priority",
	}

	require.NoError(t, posting.Validate())
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
