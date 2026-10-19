package dto

import (
	"time"

	attachment "github.com/yilin/ai-for-backend/internal/domain/resume_attachment"
)

type ResumeAttachmentResponse struct {
	ID          string `json:"id"`
	VersionID   string `json:"versionId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int64  `json:"fileSize"`
	CreatedAt   string `json:"createdAt"`
}

func ToResumeAttachmentResponse(value attachment.ResumeAttachment) ResumeAttachmentResponse {
	return ResumeAttachmentResponse{ID: value.ID, VersionID: value.VersionID, FileName: value.FileName, ContentType: value.ContentType, FileSize: value.FileSize, CreatedAt: value.CreatedAt.Format(time.RFC3339)}
}
