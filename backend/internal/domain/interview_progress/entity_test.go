package interview_progress_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestInterviewProgress_Validate_AcceptsFrontendFields(t *testing.T) {
	progress := &domainprogress.InterviewProgress{
		UserID:            "user-1",
		ApplicationID:     "application-1",
		Stage:             domainprogress.StageFirstInterview,
		OccurredAt:        time.Now(),
		Result:            domainprogress.ResultPassed,
		InterviewerOrTeam: "Frontend team",
		Note:              "Asked Vue questions",
	}

	require.NoError(t, progress.Validate())
}

func TestInterviewProgress_Validate_RejectsInvalidFrontendFields(t *testing.T) {
	tests := []struct {
		name     string
		progress domainprogress.InterviewProgress
	}{
		{name: "missing user", progress: domainprogress.InterviewProgress{ApplicationID: "a1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "missing application", progress: domainprogress.InterviewProgress{UserID: "u1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "missing occurredAt", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: domainprogress.StageApplied, Result: domainprogress.ResultPassed}},
		{name: "invalid stage", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: "phone", OccurredAt: time.Now(), Result: domainprogress.ResultPassed}},
		{name: "invalid result", progress: domainprogress.InterviewProgress{UserID: "u1", ApplicationID: "a1", Stage: domainprogress.StageApplied, OccurredAt: time.Now(), Result: "cancelled"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.progress.Validate(), domainerrors.ErrBadParams)
		})
	}
}

func TestInterviewProgress_Validate_RejectsOverlongInterviewerOrTeam(t *testing.T) {
	progress := &domainprogress.InterviewProgress{
		UserID:            "user-1",
		ApplicationID:     "application-1",
		Stage:             domainprogress.StageFirstInterview,
		OccurredAt:        time.Now(),
		Result:            domainprogress.ResultPassed,
		InterviewerOrTeam: stringOfRunes(121),
	}

	assert.ErrorIs(t, progress.Validate(), domainerrors.ErrBadParams)
}

func stringOfRunes(length int) string {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = '好'
	}
	return string(runes)
}
