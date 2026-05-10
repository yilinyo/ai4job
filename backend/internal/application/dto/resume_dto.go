package dto

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
