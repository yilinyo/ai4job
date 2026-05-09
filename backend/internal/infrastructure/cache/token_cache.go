package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	blacklistPrefix = "blacklist:"
	emailCodePrefix = "email_code:"
	emailCodeExpiry = 5 * time.Minute
)

type TokenCache struct {
	client *redis.Client
}

func NewTokenCache(client *redis.Client) *TokenCache {
	return &TokenCache{client: client}
}

// AddToBlacklist marks a JWT jti as revoked until its expiry time.
func (c *TokenCache) AddToBlacklist(ctx context.Context, jti string, ttl time.Duration) error {
	key := blacklistPrefix + jti
	return c.client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted returns true if the jti has been revoked.
func (c *TokenCache) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := blacklistPrefix + jti
	val, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// SetEmailCode stores a 6-digit code for the given email with a 5-minute TTL.
func (c *TokenCache) SetEmailCode(ctx context.Context, email, code string) error {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	return c.client.Set(ctx, key, code, emailCodeExpiry).Err()
}

// VerifyEmailCode checks the stored code. Returns true if it matches.
func (c *TokenCache) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	stored, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

// DeleteEmailCode removes the email code after successful registration.
func (c *TokenCache) DeleteEmailCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	return c.client.Del(ctx, key).Err()
}
