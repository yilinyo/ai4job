package dashboard_todo

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func validNewDashboardTodoParams() NewDashboardTodoParams {
	targetID := "11111111-1111-1111-1111-111111111111"
	return NewDashboardTodoParams{
		ID:          "todo-1",
		UserID:      "user-1",
		Title:       "Follow up with company",
		Description: "Send interview availability",
		DueAt:       nil,
		TargetType:  TargetTypeApplication,
		TargetID:    &targetID,
	}
}

func TestNewDashboardTodoTrimsFields(t *testing.T) {
	todo, err := NewDashboardTodo(NewDashboardTodoParams{
		ID:     " todo-1 ",
		UserID: " user-1 ",
		Title:  " Follow up ",
	})

	require.NoError(t, err)
	assert.Equal(t, "todo-1", todo.ID)
	assert.Equal(t, "user-1", todo.UserID)
	assert.Equal(t, "Follow up", todo.Title)
}

func TestNewDashboardTodoRejectsTitleLongerThanDatabaseLimit(t *testing.T) {
	params := validNewDashboardTodoParams()
	params.Title = strings.Repeat("a", 201)

	todo, err := NewDashboardTodo(params)

	assert.Nil(t, todo)
	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
}

func TestDashboardTodoApplyUpdateRestoresOriginalValuesWhenValidationFails(t *testing.T) {
	todo, err := NewDashboardTodo(validNewDashboardTodoParams())
	require.NoError(t, err)
	original := *todo

	title := "   "
	description := "updated description"
	dueAt := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	dueAtPtr := &dueAt
	invalidTargetType := TargetType("unknown")
	targetID := "22222222-2222-2222-2222-222222222222"
	targetIDPtr := &targetID

	err = todo.ApplyUpdate(UpdateSpec{
		Title:       &title,
		Description: &description,
		DueAt:       &dueAtPtr,
		TargetType:  &invalidTargetType,
		TargetID:    &targetIDPtr,
	})

	assert.True(t, errors.Is(err, domainerrors.ErrBadParams))
	assert.Equal(t, original, *todo)
}
