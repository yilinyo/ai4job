package apikey

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewGeneratesSinglePartKeyAndHMACDigest(t *testing.T) {
	days := 30
	key, plain, err := New("user-1", "Codex", "test-pepper", &days)

	require.NoError(t, err)
	require.True(t, strings.HasPrefix(plain, Prefix))
	require.NotContains(t, plain[len(Prefix):], ".")
	require.Equal(t, Digest(plain, "test-pepper"), key.KeyDigest)
	require.NotEqual(t, plain, key.KeyDigest)
	require.True(t, key.Matches(plain, "test-pepper"))
	require.False(t, key.Matches(plain, "other-pepper"))
	require.NotNil(t, key.ExpiresAt)
	require.WithinDuration(t, time.Now().UTC().AddDate(0, 0, 30), *key.ExpiresAt, 2*time.Second)
	require.Contains(t, key.DisplayPrefix, "...")
}

func TestAPIKeyValidityHonorsRevocationAndExpiration(t *testing.T) {
	key, _, err := New("user-1", "Codex", "test-pepper", nil)
	require.NoError(t, err)
	now := time.Now().UTC()
	require.True(t, key.ValidAt(now))

	key.Revoke(now)
	require.False(t, key.ValidAt(now))

	key.Status = StatusActive
	expired := now.Add(-time.Second)
	key.ExpiresAt = &expired
	require.False(t, key.ValidAt(now))
}
