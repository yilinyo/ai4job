package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/pagination"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type JobPostingService interface {
	Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) (*domainjob.JobPosting, error)
	List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error)
	Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error)
	Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error
	Delete(ctx context.Context, userID, id string) error
	Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error)
}

type JobPostingHandler struct {
	svc JobPostingService
}

func NewJobPostingHandler(svc JobPostingService) *JobPostingHandler {
	return &JobPostingHandler{svc: svc}
}

func frontendPage(c *gin.Context) (int, int) {
	p := pagination.FromContext(c)
	return p.Page, p.PageSize
}

func (h *JobPostingHandler) List(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	page, pageSize := frontendPage(c)
	rows, total, err := h.svc.List(c.Request.Context(), userID, c.Query("keyword"), c.Query("sourcePlatform"), c.Query("status"), page, pageSize)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toJobPostingResponses(rows), "total": total, "page": page, "pageSize": pageSize})
}

func (h *JobPostingHandler) Get(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToJobPostingResponse(row))
}

func (h *JobPostingHandler) Create(c *gin.Context) {
	var req dto.CreateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	row, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToJobPostingResponse(row))
}

func (h *JobPostingHandler) Update(c *gin.Context) {
	var req dto.UpdateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := h.svc.Update(c.Request.Context(), userID, id, req); err != nil {
		mapError(c, err)
		return
	}
	row, err := h.svc.Get(c.Request.Context(), userID, id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToJobPostingResponse(row))
}

func (h *JobPostingHandler) Delete(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func toJobPostingResponses(rows []domainjob.JobPosting) []*dto.JobPostingResponse {
	responses := make([]*dto.JobPostingResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, dto.ToJobPostingResponse(&rows[i]))
	}
	return responses
}

func (h *JobPostingHandler) Parse(c *gin.Context) {
	var req dto.ParseJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	row, err := h.svc.Parse(c.Request.Context(), req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, row)
}
