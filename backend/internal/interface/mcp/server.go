package mcpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	apikeyapp "github.com/yilin/ai-for-backend/internal/application/api_key"
)

type contextKey string

const userIDKey contextKey = "mcpUserID"

func UserID(ctx context.Context) string { value, _ := ctx.Value(userIDKey).(string); return value }

type Handler struct {
	auth      *apikeyapp.Service
	transport http.Handler
}

func New(auth *apikeyapp.Service, server *mcp.Server) http.Handler {
	return &Handler{
		auth: auth,
		transport: mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil),
	}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Access-Token") != "" || r.Header.Get("Cookie") != "" {
		http.Error(w, "API key authentication required", http.StatusUnauthorized)
		return
	}
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		http.Error(w, "API key authentication required", http.StatusUnauthorized)
		return
	}
	key, err := h.auth.Authenticate(r.Context(), strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		http.Error(w, "API key authentication failed", http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), userIDKey, key.UserID)
	h.transport.ServeHTTP(w, r.WithContext(ctx))
}

func JSONResult(value any) (*mcp.CallToolResult, any, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(raw)}}}, nil, nil
}
