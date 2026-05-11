package job_posting

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

type JobPostingService struct {
	repo domainjob.Repository
}

func NewJobPostingService(repo domainjob.Repository) *JobPostingService {
	return &JobPostingService{repo: repo}
}

func (s *JobPostingService) Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) error {
	status := req.Status
	if status == "" {
		status = domainjob.StatusPending
	}

	posting := &domainjob.JobPosting{
		ID:              uuid.NewString(),
		UserID:          userID,
		CompanyName:     req.CompanyName,
		CompanyIntro:    req.CompanyIntro,
		JobTitle:        req.JobTitle,
		JobRequirements: req.JobRequirements,
		Base:            req.Base,
		SalaryRange:     req.SalaryRange,
		SourcePlatform:  req.SourcePlatform,
		SourceURL:       req.SourceURL,
		Status:          status,
		Remark:          req.Remark,
	}
	if err := posting.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, posting)
}

func (s *JobPostingService) List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return s.repo.ListByUser(ctx, domainjob.ListFilter{
		UserID:         userID,
		Keyword:        keyword,
		SourcePlatform: sourcePlatform,
		Status:         domainjob.Status(status),
		Page:           page,
		PageSize:       pageSize,
	})
}

func (s *JobPostingService) Get(ctx context.Context, userID, id string) (*domainjob.JobPosting, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *JobPostingService) Update(ctx context.Context, userID, id string, req dto.UpdateJobPostingRequest) error {
	posting, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	if req.CompanyName != nil {
		posting.CompanyName = *req.CompanyName
	}
	if req.CompanyIntro != nil {
		posting.CompanyIntro = *req.CompanyIntro
	}
	if req.JobTitle != nil {
		posting.JobTitle = *req.JobTitle
	}
	if req.JobRequirements != nil {
		posting.JobRequirements = *req.JobRequirements
	}
	if req.Base != nil {
		posting.Base = *req.Base
	}
	if req.SalaryRange != nil {
		posting.SalaryRange = *req.SalaryRange
	}
	if req.SourcePlatform != nil {
		posting.SourcePlatform = *req.SourcePlatform
	}
	if req.SourceURL != nil {
		posting.SourceURL = *req.SourceURL
	}
	if req.Status != nil {
		posting.Status = *req.Status
	}
	if req.Remark != nil {
		posting.Remark = *req.Remark
	}

	if err := posting.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, posting)
}

func (s *JobPostingService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *JobPostingService) Parse(ctx context.Context, req dto.ParseJobPostingRequest) (*dto.JobPostingResponse, error) {
	_ = ctx
	return &dto.JobPostingResponse{
		SourceURL:      req.URL,
		SourcePlatform: guessSourcePlatform(req.URL),
	}, nil
}

func guessSourcePlatform(rawURL string) string {
	lowerURL := strings.ToLower(rawURL)
	if strings.Contains(lowerURL, "zhipin") {
		return "Boss"
	}
	if strings.Contains(lowerURL, "jobs.bytedance") {
		return "官网"
	}
	return "其他"
}
