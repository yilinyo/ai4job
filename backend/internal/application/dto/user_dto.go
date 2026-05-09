package dto

import "github.com/yilin/ai-for-backend/internal/domain/user"

type SendEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SendEmailCodeResponse struct {
	Email     string `json:"email"`
	MockCode  string `json:"mockCode,omitempty"` // only in development
	ExpiresIn int    `json:"expiresIn"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	EmailCode string `json:"emailCode" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string            `json:"token"`
	UserInfo *UserInfoResponse `json:"userInfo"`
}

type UserInfoResponse struct {
	ID                   string                     `json:"id"`
	Username             string                     `json:"username"`
	RealName             string                     `json:"realName"`
	Age                  int                        `json:"age"`
	Email                string                     `json:"email"`
	Phone                string                     `json:"phone"`
	JobIntention         string                     `json:"jobIntention"`
	Avatar               string                     `json:"avatar"`
	Location             string                     `json:"location"`
	PersonalAdvantage    string                     `json:"personalAdvantage"`
	EducationExperiences []user.EducationExperience `json:"educationExperiences"`
}

type UpdateProfileRequest struct {
	RealName             string                     `json:"realName"`
	Age                  int                        `json:"age"`
	Phone                string                     `json:"phone"`
	JobIntention         string                     `json:"jobIntention"`
	Avatar               string                     `json:"avatar"`
	Location             string                     `json:"location"`
	PersonalAdvantage    string                     `json:"personalAdvantage"`
	EducationExperiences []user.EducationExperience `json:"educationExperiences"`
}

func ToUserInfoResponse(u *user.User) *UserInfoResponse {
	exps := u.EducationExperiences
	if exps == nil {
		exps = []user.EducationExperience{}
	}
	return &UserInfoResponse{
		ID:                   u.ID,
		Username:             u.Username,
		RealName:             u.RealName,
		Age:                  u.Age,
		Email:                u.Email,
		Phone:                u.Phone,
		JobIntention:         u.JobIntention,
		Avatar:               u.Avatar,
		Location:             u.Location,
		PersonalAdvantage:    u.PersonalAdvantage,
		EducationExperiences: exps,
	}
}
