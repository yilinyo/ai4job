package apikey

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const Prefix = "ai4job-"

type Status string

const (
	StatusActive  Status = "active"
	StatusRevoked Status = "revoked"
)

type APIKey struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	UserID        string `gorm:"type:varchar(36);not null;index"`
	Name          string `gorm:"type:varchar(100);not null"`
	DisplayPrefix string `gorm:"type:varchar(32);not null"`
	KeyDigest     string `gorm:"type:varchar(64);not null;uniqueIndex"`
	Status        Status `gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time
	ExpiresAt     *time.Time
	LastUsedAt    *time.Time
	RevokedAt     *time.Time
}

func New(userID, name, pepper string, expiresInDays *int) (*APIKey, string, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(name) == "" || len([]rune(name)) > 100 || pepper == "" {
		return nil, "", errors.New("invalid api key parameters")
	}
	if expiresInDays != nil && *expiresInDays != 30 && *expiresInDays != 60 && *expiresInDays != 90 {
		return nil, "", errors.New("invalid api key expiry")
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, "", err
	}
	plain := Prefix + base64.RawURLEncoding.EncodeToString(raw)
	var expiresAt *time.Time
	if expiresInDays != nil {
		value := time.Now().UTC().AddDate(0, 0, *expiresInDays)
		expiresAt = &value
	}
	return &APIKey{ID: uuid.NewString(), UserID: userID, Name: strings.TrimSpace(name), DisplayPrefix: displayPrefix(plain), KeyDigest: Digest(plain, pepper), Status: StatusActive, ExpiresAt: expiresAt}, plain, nil
}

func Digest(plain, pepper string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte(plain))
	return hex.EncodeToString(mac.Sum(nil))
}
func (k *APIKey) Matches(plain, pepper string) bool {
	return hmac.Equal([]byte(k.KeyDigest), []byte(Digest(plain, pepper)))
}
func (k *APIKey) ValidAt(now time.Time) bool {
	return k != nil && k.Status == StatusActive && (k.ExpiresAt == nil || now.Before(*k.ExpiresAt))
}
func (k *APIKey) Revoke(now time.Time) { k.Status = StatusRevoked; k.RevokedAt = &now }
func displayPrefix(plain string) string {
	if len(plain) <= len(Prefix)+8 {
		return plain + "..."
	}
	return plain[:len(Prefix)+8] + "..."
}
