package handler

import (
	"context"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	attachment "github.com/yilin/ai-for-backend/internal/domain/resume_attachment"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type ResumeAttachmentService interface {
	List(context.Context, string, string) ([]attachment.ResumeAttachment, error)
	Upload(context.Context, string, string, *multipart.FileHeader) (*attachment.ResumeAttachment, error)
	Download(context.Context, string, string, string) (*attachment.ResumeAttachment, error)
	Delete(context.Context, string, string, string) error
}

type ResumeAttachmentHandler struct{ svc ResumeAttachmentService }

func NewResumeAttachmentHandler(svc ResumeAttachmentService) *ResumeAttachmentHandler {
	return &ResumeAttachmentHandler{svc: svc}
}

func (h *ResumeAttachmentHandler) List(c *gin.Context) {
	h.withUser(c, func(userID string) {
		rows, err := h.svc.List(c, userID, c.Param("id"))
		if err != nil {
			mapError(c, err)
			return
		}
		out := make([]dto.ResumeAttachmentResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, dto.ToResumeAttachmentResponse(row))
		}
		response.Success(c, out)
	})
}
func (h *ResumeAttachmentHandler) Upload(c *gin.Context) {
	h.withUser(c, func(userID string) {
		file, err := c.FormFile("file")
		if err != nil {
			response.Fail(c, pkgerrors.CodeBadParams, "请选择 PDF 或 Markdown 附件")
			return
		}
		row, err := h.svc.Upload(c, userID, c.Param("id"), file)
		if err != nil {
			mapError(c, err)
			return
		}
		response.Success(c, dto.ToResumeAttachmentResponse(*row))
	})
}
func (h *ResumeAttachmentHandler) Download(c *gin.Context) {
	h.withUser(c, func(userID string) {
		row, err := h.svc.Download(c, userID, c.Param("id"), c.Param("attachmentId"))
		if err != nil {
			mapError(c, err)
			return
		}
		c.Header("Content-Type", row.ContentType)
		c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": row.FileName}))
		c.Data(http.StatusOK, row.ContentType, row.Data)
	})
}
func (h *ResumeAttachmentHandler) Delete(c *gin.Context) {
	h.withUser(c, func(userID string) {
		if err := h.svc.Delete(c, userID, c.Param("id"), c.Param("attachmentId")); err != nil {
			mapError(c, err)
			return
		}
		response.Success(c, nil)
	})
}
func (h *ResumeAttachmentHandler) withUser(c *gin.Context, fn func(string)) {
	userID, ok := resumeUserID(c)
	if ok {
		fn(userID)
	}
}
