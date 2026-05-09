package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/user"
)

func TestUser_SetPassword_HashesPassword(t *testing.T) {
	u := &user.User{}
	err := u.SetPassword("plaintext123")
	assert.NoError(t, err)
	assert.NotEqual(t, "plaintext123", u.PasswordHash)
	assert.NotEmpty(t, u.PasswordHash)
}

func TestUser_CheckPassword_CorrectPassword(t *testing.T) {
	u := &user.User{}
	_ = u.SetPassword("plaintext123")
	assert.True(t, u.CheckPassword("plaintext123"))
}

func TestUser_CheckPassword_WrongPassword(t *testing.T) {
	u := &user.User{}
	_ = u.SetPassword("plaintext123")
	assert.False(t, u.CheckPassword("wrongpassword"))
}

func TestUser_SetPassword_EmptyPassword_ReturnsError(t *testing.T) {
	u := &user.User{}
	err := u.SetPassword("")
	assert.Error(t, err)
}

func TestUser_SetPassword_TooLong_ReturnsError(t *testing.T) {
	u := &user.User{}
	longPassword := string(make([]byte, 73)) // 73 bytes > 72 byte limit
	err := u.SetPassword(longPassword)
	assert.Error(t, err)
}
