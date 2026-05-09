package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/infrastructure/auth"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret", 24)

	token, err := mgr.Generate("user-123", "alice")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "alice", claims.Username)
	assert.NotEmpty(t, claims.ID) // jti
}

func TestJWTManager_Parse_InvalidToken(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret", 24)
	_, err := mgr.Parse("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTManager_RemainingTTL_Positive(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret", 24)
	token, _ := mgr.Generate("user-123", "alice")
	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	ttl := mgr.RemainingTTL(claims)
	assert.Positive(t, ttl)
}
