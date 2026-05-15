package job_posting

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	jobpostingrepo "github.com/yilin/ai-for-backend/internal/domain/job_posting/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

type JobPostingService struct {
	repo jobpostingrepo.Repository
}

func NewJobPostingService(repo jobpostingrepo.Repository) *JobPostingService {
	return &JobPostingService{repo: repo}
}

func (s *JobPostingService) Create(ctx context.Context, userID string, req dto.CreateJobPostingRequest) (*domainjob.JobPosting, error) {
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
		return nil, err
	}
	if err := s.repo.Create(ctx, posting); err != nil {
		return nil, err
	}
	return posting, nil
}

func (s *JobPostingService) List(ctx context.Context, userID, keyword, sourcePlatform, status string, page, pageSize int) ([]domainjob.JobPosting, int64, error) {
	filterStatus := domainjob.Status(status)
	if status != "" && !validStatus(filterStatus) {
		return nil, 0, domainerrors.ErrBadParams
	}

	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return s.repo.ListByUser(ctx, jobpostingrepo.ListFilter{
		UserID:         userID,
		Keyword:        keyword,
		SourcePlatform: sourcePlatform,
		Status:         filterStatus,
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
	if !validParseURL(req.URL) {
		return nil, domainerrors.ErrBadParams
	}

	sourcePlatform := guessSourcePlatform(req.URL)
	return &dto.JobPostingResponse{
		CompanyName:     parsedCompanyName(sourcePlatform),
		CompanyIntro:    "通过解析接口返回的公司介绍示例",
		JobTitle:        parsedJobTitle(sourcePlatform),
		JobRequirements: "通过解析接口返回的岗位要求示例，用户可继续修改",
		Base:            parsedBase(sourcePlatform),
		SalaryRange:     "25k-35k x 16",
		SourcePlatform:  sourcePlatform,
		SourceURL:       req.URL,
	}, nil
}

func validStatus(status domainjob.Status) bool {
	switch status {
	case domainjob.StatusPending, domainjob.StatusApplied, domainjob.StatusNotFit, domainjob.StatusClosed:
		return true
	default:
		return false
	}
}

func validParseURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
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

func parsedCompanyName(sourcePlatform string) string {
	if sourcePlatform == "Boss" {
		return "腾讯科技"
	}
	return "示例公司"
}

func parsedJobTitle(sourcePlatform string) string {
	if sourcePlatform == "Boss" {
		return "前端开发工程师"
	}
	return "待确认岗位"
}

func parsedBase(sourcePlatform string) string {
	if sourcePlatform == "Boss" {
		return "深圳"
	}
	return "北京"
}
