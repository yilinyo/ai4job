package activity_log

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func validNewActivityLogParams() NewActivityLogParams {
	return NewActivityLogParams{
		ID:          "log-1",
		UserID:      "user-1",
		Action:      ActionCreated,
		TargetType:  TargetTypeApplication,
		TargetID:    "11111111-1111-1111-1111-111111111111",
		TargetLabel: "Backend Engineer",
		Summary:     "Created application",
		Route:       "/applications/11111111-1111-1111-1111-111111111111",
	}
}

func TestNewActivityLogRejectsFieldsLongerThanDatabaseLimits(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*NewActivityLogParams)
	}{
		{
			name: "action longer than 40 characters",
			mutate: func(params *NewActivityLogParams) {
				params.Action = Action(strings.Repeat("a", 41))
			},
		},
		{
			name: "target type longer than 40 characters",
			mutate: func(params *NewActivityLogParams) {
				params.TargetType = TargetType(strings.Repeat("a", 41))
			},
		},
		{
			name: "target label longer than 200 characters",
			mutate: func(params *NewActivityLogParams) {
				params.TargetLabel = strings.Repeat("a", 201)
			},
		},
		{
			name: "route longer than 500 characters",
			mutate: func(params *NewActivityLogParams) {
				params.Route = "/" + strings.Repeat("a", 500)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := validNewActivityLogParams()
			tc.mutate(&params)

			log, err := NewActivityLog(params)

			assert.Nil(t, log)
			assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
		})
	}
}

func TestActivityLogValidateRejectsActorTypeLongerThanDatabaseLimit(t *testing.T) {
	log, err := NewActivityLog(validNewActivityLogParams())
	assert.NoError(t, err)
	log.ActorType = strings.Repeat("a", 41)

	err = log.validate()

	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestNewActivityLogAllowsEmptyRouteDefault(t *testing.T) {
	params := validNewActivityLogParams()
	params.Route = "   "

	log, err := NewActivityLog(params)

	assert.NoError(t, err)
	assert.Equal(t, "", log.Route)
}
