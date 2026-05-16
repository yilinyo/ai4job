package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexibleTime_UnmarshalJSON_SupportsDateOnly(t *testing.T) {
	var got FlexibleTime

	err := json.Unmarshal([]byte(`"2026-05-11"`), &got)

	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC), got.Value())
}

func TestFlexibleTime_UnmarshalJSON_SupportsRFC3339(t *testing.T) {
	var got FlexibleTime

	err := json.Unmarshal([]byte(`"2026-05-11T10:00:00Z"`), &got)

	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC), got.Value())
}

func TestFlexibleTime_UnmarshalJSON_SupportsLocalDateTimeWithoutTimezone(t *testing.T) {
	var got FlexibleTime

	err := json.Unmarshal([]byte(`"2026-05-11T23:28:04"`), &got)

	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 5, 11, 23, 28, 4, 0, time.Local), got.Value())
}

func TestFlexibleTime_UnmarshalJSON_RejectsInvalidFormat(t *testing.T) {
	var got FlexibleTime

	err := json.Unmarshal([]byte(`"2026/05/11"`), &got)

	require.Error(t, err)
}
