package dto

import "time"

type CreateAPIKeyRequest struct {
	Name          string `json:"name" binding:"required,max=100"`
	ExpiresInDays *int   `json:"expiresInDays"`
}

type APIKeyResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	DisplayPrefix string     `json:"displayPrefix"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	LastUsedAt    *time.Time `json:"lastUsedAt"`
}

type CreatedAPIKeyResponse struct {
	*APIKeyResponse
	Key string `json:"key"`
}
