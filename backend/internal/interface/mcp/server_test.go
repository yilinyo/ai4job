package mcpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
	apikeyapp "github.com/yilin/ai-for-backend/internal/application/api_key"
	domain "github.com/yilin/ai-for-backend/internal/domain/api_key"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type apiKeyMemoryRepo struct{ keys map[string]*domain.APIKey }

func (r *apiKeyMemoryRepo) Create(_ context.Context, key *domain.APIKey) error {
	r.keys[key.ID] = key
	return nil
}
func (r *apiKeyMemoryRepo) ListByUser(_ context.Context, userID string) ([]domain.APIKey, error) {
	return nil, nil
}
func (r *apiKeyMemoryRepo) FindByIDAndUser(_ context.Context, id, userID string) (*domain.APIKey, error) {
	return nil, domainerrors.ErrNotFound
}
func (r *apiKeyMemoryRepo) FindByDigest(_ context.Context, digest string) (*domain.APIKey, error) {
	for _, key := range r.keys {
		if key.KeyDigest == digest {
			return key, nil
		}
	}
	return nil, domainerrors.ErrNotFound
}
func (r *apiKeyMemoryRepo) Save(_ context.Context, key *domain.APIKey) error {
	r.keys[key.ID] = key
	return nil
}
func (r *apiKeyMemoryRepo) Touch(_ context.Context, key *domain.APIKey, now time.Time) error {
	key.LastUsedAt = &now
	return nil
}

func mcpInitializeRequest(key string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/mcp", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Body = ioNopCloser(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`)
	return req
}

type stringReadCloser struct{ *strings.Reader }

func (stringReadCloser) Close() error           { return nil }
func ioNopCloser(value string) stringReadCloser { return stringReadCloser{strings.NewReader(value)} }

func TestHandlerAcceptsOnlyActiveBearerAPIKeys(t *testing.T) {
	const pepper = "test-pepper"
	active, plain, err := domain.New("user-1", "active", pepper, nil)
	require.NoError(t, err)
	revoked, revokedPlain, err := domain.New("user-1", "revoked", pepper, nil)
	require.NoError(t, err)
	revoked.Revoke(time.Now().UTC())
	expired, expiredPlain, err := domain.New("user-1", "expired", pepper, nil)
	require.NoError(t, err)
	expiresAt := time.Now().UTC().Add(-time.Second)
	expired.ExpiresAt = &expiresAt

	repo := &apiKeyMemoryRepo{keys: map[string]*domain.APIKey{active.ID: active, revoked.ID: revoked, expired.ID: expired}}
	h := New(apikeyapp.NewService(repo, pepper, nil), mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1"}, nil))

	for _, tc := range []struct {
		name, key string
		mutate    func(*http.Request)
	}{
		{name: "missing bearer", key: "", mutate: func(req *http.Request) { req.Header.Del("Authorization") }},
		{name: "invalid key", key: "ai4job-not-a-real-key-material"},
		{name: "revoked key", key: revokedPlain},
		{name: "expired key", key: expiredPlain},
		{name: "jwt header", key: plain, mutate: func(req *http.Request) { req.Header.Set("X-Access-Token", "jwt") }},
		{name: "cookie", key: plain, mutate: func(req *http.Request) { req.Header.Set("Cookie", "session=value") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := mcpInitializeRequest(tc.key)
			if tc.mutate != nil {
				tc.mutate(req)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, mcpInitializeRequest(plain))
	require.Equal(t, http.StatusOK, w.Code)
	require.NotEmpty(t, w.Header().Get("Mcp-Session-Id"))
}
