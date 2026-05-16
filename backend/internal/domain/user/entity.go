package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeCampus     JobType = "campus"
	JobTypeSocial     JobType = "social"
	JobTypeInternship JobType = "internship"
)

type EducationExperience struct {
	School         string `json:"school"`
	Education      string `json:"education"`
	Major          string `json:"major"`
	AdmissionDate  string `json:"admissionDate"`
	GraduationDate string `json:"graduationDate"`
}

type User struct {
	ID                   string `gorm:"type:uuid;primaryKey"`
	Username             string `gorm:"uniqueIndex;not null"`
	PasswordHash         string `gorm:"not null"`
	Email                string `gorm:"uniqueIndex;not null"`
	RealName             string
	Age                  int
	Phone                string
	JobIntention         string
	Avatar               string
	Location             string
	PersonalAdvantage    string                `gorm:"type:text"`
	EducationExperiences []EducationExperience `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

type NewUserParams struct {
	ID       string
	Username string
	Email    string
	Password string
}

func NewUser(params NewUserParams) (*User, error) {
	u := &User{
		ID:                   params.ID,
		Username:             params.Username,
		Email:                params.Email,
		EducationExperiences: []EducationExperience{},
	}
	if err := u.SetPassword(params.Password); err != nil {
		return nil, err
	}
	u.Normalize()
	return u, nil
}

type ProfileUpdate struct {
	RealName             string
	Age                  int
	Phone                string
	JobIntention         string
	Avatar               string
	Location             string
	PersonalAdvantage    string
	EducationExperiences []EducationExperience
}

func (u *User) UpdateProfile(update ProfileUpdate) {
	u.RealName = update.RealName
	u.Age = update.Age
	u.Phone = update.Phone
	u.JobIntention = update.JobIntention
	u.Avatar = update.Avatar
	u.Location = update.Location
	u.PersonalAdvantage = update.PersonalAdvantage
	u.EducationExperiences = update.EducationExperiences
	u.Normalize()
}

func (u *User) Normalize() {
	if u == nil {
		return
	}
	if u.EducationExperiences == nil {
		u.EducationExperiences = []EducationExperience{}
	}
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	u.Normalize()
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) error {
	u.Normalize()
	return nil
}

func (u *User) SetPassword(plain string) error {
	if plain == "" {
		return errors.New("password must not be empty")
	}
	if len([]byte(plain)) > 72 {
		return errors.New("password exceeds 72-byte bcrypt limit")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
	return err == nil
}
