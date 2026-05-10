package dto

import (
	"time"

	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
)

type CreateResumeRepoRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type UpdateResumeRepoRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type CreateResumeVersionRequest struct {
	Title     string `json:"title" binding:"required,min=1,max=120"`
	Content   string `json:"content" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

type UpdateResumeVersionRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=120"`
	Content string `json:"content" binding:"required"`
}

type ResumeRepoResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ResumeVersionResponse struct {
	ID         string    `json:"id"`
	RepoID     string    `json:"repoId"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	VersionNum int       `json:"versionNum"`
	IsDefault  bool      `json:"isDefault"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func ToResumeRepoResponse(repo *domainrepo.ResumeRepo) *ResumeRepoResponse {
	if repo == nil {
		return nil
	}
	return &ResumeRepoResponse{
		ID:          repo.ID,
		UserID:      repo.UserID,
		Name:        repo.Name,
		Description: repo.Description,
		CreatedAt:   repo.CreatedAt,
		UpdatedAt:   repo.UpdatedAt,
	}
}

func ToResumeRepoResponses(repos []domainrepo.ResumeRepo) []ResumeRepoResponse {
	responses := make([]ResumeRepoResponse, 0, len(repos))
	for i := range repos {
		responses = append(responses, *ToResumeRepoResponse(&repos[i]))
	}
	return responses
}

func ToResumeVersionResponse(version *domainversion.ResumeVersion) *ResumeVersionResponse {
	if version == nil {
		return nil
	}
	return &ResumeVersionResponse{
		ID:         version.ID,
		RepoID:     version.RepoID,
		Title:      version.Title,
		Content:    version.Content,
		VersionNum: version.VersionNum,
		IsDefault:  version.IsDefault,
		CreatedAt:  version.CreatedAt,
		UpdatedAt:  version.UpdatedAt,
	}
}

func ToResumeVersionResponses(versions []domainversion.ResumeVersion) []ResumeVersionResponse {
	responses := make([]ResumeVersionResponse, 0, len(versions))
	for i := range versions {
		responses = append(responses, *ToResumeVersionResponse(&versions[i]))
	}
	return responses
}
