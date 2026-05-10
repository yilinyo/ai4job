package dto

import (
	"time"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

const timeFormat = time.RFC3339

type CreateResumeRepoRequest struct {
	Name           string             `json:"name" binding:"required,min=1,max=100"`
	JobType        domainrepo.JobType `json:"jobType" binding:"required"`
	TargetPosition string             `json:"targetPosition" binding:"required,min=1,max=100"`
	Description    string             `json:"description"`
}

type UpdateResumeRepoRequest struct {
	Name           string             `json:"name" binding:"required,min=1,max=100"`
	JobType        domainrepo.JobType `json:"jobType" binding:"required"`
	TargetPosition string             `json:"targetPosition" binding:"required,min=1,max=100"`
	Description    string             `json:"description"`
}

type ResumeRepoResponse struct {
	ID             string             `json:"id"`
	UserID         string             `json:"userId"`
	Name           string             `json:"name"`
	JobType        domainrepo.JobType `json:"jobType"`
	TargetPosition string             `json:"targetPosition"`
	Description    string             `json:"description"`
	VersionCount   int                `json:"versionCount"`
	CreatedAt      string             `json:"createdAt"`
	UpdatedAt      string             `json:"updatedAt"`
}

func ToResumeRepoResponse(r domainrepo.ResumeRepo) ResumeRepoResponse {
	return ResumeRepoResponse{
		ID:             r.ID,
		UserID:         r.UserID,
		Name:           r.Name,
		JobType:        r.JobType,
		TargetPosition: r.TargetPosition,
		Description:    r.Description,
		VersionCount:   r.VersionCount,
		CreatedAt:      r.CreatedAt.Format(timeFormat),
		UpdatedAt:      r.UpdatedAt.Format(timeFormat),
	}
}

type CreateResumeVersionRequest struct {
	RepoID  string `json:"repoId" binding:"required"`
	Title   string `json:"title" binding:"required,min=1,max=120"`
	Content string `json:"content" binding:"required"`
	Remark  string `json:"remark"`
}

type UpdateResumeVersionRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=120"`
	Content  string `json:"content" binding:"required"`
	Remark   string `json:"remark"`
	IsActive *bool  `json:"isActive"`
}

type ResumeVersionResponse struct {
	ID        string `json:"id"`
	RepoID    string `json:"repoId"`
	Version   string `json:"version"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Remark    string `json:"remark"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func ToResumeVersionResponse(v domainversion.ResumeVersion) ResumeVersionResponse {
	return ResumeVersionResponse{
		ID:        v.ID,
		RepoID:    v.RepoID,
		Version:   v.Version,
		Title:     v.Title,
		Content:   v.Content,
		Remark:    v.Remark,
		IsActive:  v.IsActive,
		CreatedAt: v.CreatedAt.Format(timeFormat),
		UpdatedAt: v.UpdatedAt.Format(timeFormat),
	}
}
