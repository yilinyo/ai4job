package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewUserInitializesAndHashesPassword(t *testing.T) {
	u, err := user.NewUser(user.NewUserParams{
		ID:       "user-1",
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", u.ID)
	assert.Equal(t, "alice", u.Username)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.NotEmpty(t, u.PasswordHash)
	assert.NotEqual(t, "secret123", u.PasswordHash)
	assert.NotNil(t, u.EducationExperiences)
	assert.True(t, u.CheckPassword("secret123"))
}

func TestNewUserRejectsInvalidPassword(t *testing.T) {
	u, err := user.NewUser(user.NewUserParams{ID: "user-1", Username: "alice", Email: "alice@example.com", Password: ""})

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestUserUpdateProfileUpdatesFieldsAndNormalizesEducation(t *testing.T) {
	u := &user.User{ID: "user-1", Username: "alice", Email: "alice@example.com"}

	u.UpdateProfile(user.ProfileUpdate{
		RealName:             "Alice",
		Age:                  22,
		Phone:                "13800000000",
		JobIntention:         "Backend",
		Avatar:               "avatar.png",
		Location:             "Shanghai",
		PersonalAdvantage:    "DDD",
		EducationExperiences: nil,
	})

	assert.Equal(t, "Alice", u.RealName)
	assert.Equal(t, 22, u.Age)
	assert.Equal(t, "Backend", u.JobIntention)
	assert.NotNil(t, u.EducationExperiences)
}

func TestUser_Normalize_InitializesNilEducationExperiences(t *testing.T) {
	u := &user.User{}

	u.Normalize()

	assert.NotNil(t, u.EducationExperiences)
	assert.Empty(t, u.EducationExperiences)
}

func TestUser_Normalize_PreservesExistingEducationExperiences(t *testing.T) {
	exps := []user.EducationExperience{
		{School: "Test University", Education: "本科", Major: "Computer Science"},
	}
	u := &user.User{EducationExperiences: exps}

	u.Normalize()

	assert.Equal(t, exps, u.EducationExperiences)
}
