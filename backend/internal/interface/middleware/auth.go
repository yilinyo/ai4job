package middleware

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

// JWTParser parses a raw JWT string and returns its claims.
type JWTParser interface {
	Parse(tokenStr string) (userID, username, jti string, err error)
}

// Blacklist checks whether a JWT ID has been revoked.
type Blacklist interface {
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

// Auth returns a gin middleware that validates the JWT from X-Access-Token.
func Auth(parser JWTParser, bl Blacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Access-Token")
		if token == "" {
			response.Fail(c, pkgerrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		userID, username, jti, err := parser.Parse(token)
		if err != nil {
			response.Fail(c, pkgerrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		blacklisted, err := bl.IsBlacklisted(c.Request.Context(), jti)
		if err != nil {
			log.Printf("auth middleware: blacklist check failed for jti %s: %v", jti, err)
			response.Fail(c, pkgerrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}
		if blacklisted {
			response.Fail(c, pkgerrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}
