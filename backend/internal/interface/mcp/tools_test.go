package mcpserver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuestionToolSchemasExposeEnumsAndBatchLimits(t *testing.T) {
	var listSchema map[string]any
	require.NoError(t, json.Unmarshal(questionListInputSchema(), &listSchema))
	properties := listSchema["properties"].(map[string]any)
	require.ElementsMatch(t, []any{"easy", "medium", "hard"}, properties["difficulty"].(map[string]any)["enum"])
	require.ElementsMatch(t, []any{"knowledge", "scenario", "coding", "algorithm", "project", "behavior", "discussion"}, properties["questionType"].(map[string]any)["enum"])
	require.ElementsMatch(t, []any{"unreviewed", "weak", "normal", "mastered"}, properties["masteryStatus"].(map[string]any)["enum"])
	require.Contains(t, properties, "isFavorite")

	var createSchema map[string]any
	require.NoError(t, json.Unmarshal(questionCreateInputSchema(), &createSchema))
	questions := createSchema["properties"].(map[string]any)["questions"].(map[string]any)
	require.Equal(t, float64(1), questions["minItems"])
	require.Equal(t, float64(50), questions["maxItems"])
	itemProperties := questions["items"].(map[string]any)["properties"].(map[string]any)
	require.ElementsMatch(t, []any{"easy", "medium", "hard"}, itemProperties["difficulty"].(map[string]any)["enum"])
}
