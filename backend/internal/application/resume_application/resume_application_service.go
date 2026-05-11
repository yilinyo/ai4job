package resume_application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

type JobPostingReader interface {
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainjob.JobPosting, error)
}

type ResumeRepoReader interface {
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainrepo.ResumeRepo, error)
}

type ResumeVersionReader interface {
	FindByID(ctx context.Context, id string) (*domainversion.ResumeVersion, error)
}

type ResumeApplicationService struct {
	repo           domainapp.Repository
	jobPostings    JobPostingReader
	resumeRepos    ResumeRepoReader
	resumeVersions ResumeVersionReader
}

func NewResumeApplicationService(repo domainapp.Repository, jobPostings JobPostingReader, resumeRepos ResumeRepoReader, resumeVersions ResumeVersionReader) *ResumeApplicationService {
	return &ResumeApplicationService{repo: repo, jobPostings: jobPostings, resumeRepos: resumeRepos, resumeVersions: resumeVersions}
}

func (s *ResumeApplicationService) Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) (*domainapp.ResumeApplication, error) {
	status := req.CurrentStatus
	if status == "" {
		status = domainapp.StatusApplied
	}

	application := &domainapp.ResumeApplication{
		ID:                        uuid.NewString(),
		UserID:                    userID,
		JobPostingID:              normalizeOptionalJobPostingID(req.JobPostingID),
		ResumeVersionID:           req.ResumeVersionID,
		RepoID:                    req.RepoID,
		CompanyName:               req.CompanyName,
		CompanyIntro:              req.CompanyIntro,
		JobTitle:                  req.JobTitle,
		JobRequirements:           req.JobRequirements,
		Base:                      req.Base,
		SalaryRange:               req.SalaryRange,
		SourcePlatform:            req.SourcePlatform,
		SourceURL:                 req.SourceURL,
		DeliveryChannel:           req.DeliveryChannel,
		AppliedAt:                 req.AppliedAt,
		CurrentStatus:             status,
		InterviewSummary:          req.InterviewSummary,
		InterviewNotes:            req.InterviewNotes,
		ResumeMatchScore:          req.ResumeMatchScore,
		InterviewPerformanceScore: req.InterviewPerformanceScore,
		RoleInterestScore:         req.RoleInterestScore,
		OverallScore:              req.OverallScore,
		Remark:                    req.Remark,
	}
	if err := application.Validate(); err != nil {
		return nil, err
	}
	if err := s.ensureReferencesOwned(ctx, userID, application.JobPostingID, application.RepoID, application.ResumeVersionID); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, application); err != nil {
		return nil, err
	}
	return application, nil
}

func (s *ResumeApplicationService) List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error) {
	filterStatus := domainapp.Status(currentStatus)
	if currentStatus != "" && !validStatus(filterStatus) {
		return nil, 0, domainerrors.ErrBadParams
	}

	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return s.repo.ListByUser(ctx, domainapp.ListFilter{
		UserID:          userID,
		ResumeVersionID: resumeVersionID,
		RepoID:          repoID,
		Keyword:         keyword,
		CurrentStatus:   filterStatus,
		Page:            page,
		PageSize:        pageSize,
	})
}

func (s *ResumeApplicationService) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *ResumeApplicationService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeApplicationRequest) error {
	application, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	if req.JobPostingID != nil {
		application.JobPostingID = normalizeOptionalJobPostingID(req.JobPostingID)
	}
	if req.ResumeVersionID != nil {
		application.ResumeVersionID = *req.ResumeVersionID
	}
	if req.RepoID != nil {
		application.RepoID = *req.RepoID
	}
	if req.CompanyName != nil {
		application.CompanyName = *req.CompanyName
	}
	if req.CompanyIntro != nil {
		application.CompanyIntro = *req.CompanyIntro
	}
	if req.JobTitle != nil {
		application.JobTitle = *req.JobTitle
	}
	if req.JobRequirements != nil {
		application.JobRequirements = *req.JobRequirements
	}
	if req.Base != nil {
		application.Base = *req.Base
	}
	if req.SalaryRange != nil {
		application.SalaryRange = *req.SalaryRange
	}
	if req.SourcePlatform != nil {
		application.SourcePlatform = *req.SourcePlatform
	}
	if req.SourceURL != nil {
		application.SourceURL = *req.SourceURL
	}
	if req.DeliveryChannel != nil {
		application.DeliveryChannel = *req.DeliveryChannel
	}
	if req.AppliedAt != nil {
		application.AppliedAt = *req.AppliedAt
	}
	if req.CurrentStatus != nil {
		application.CurrentStatus = *req.CurrentStatus
	}
	if req.InterviewSummary != nil {
		application.InterviewSummary = *req.InterviewSummary
	}
	if req.InterviewNotes != nil {
		application.InterviewNotes = *req.InterviewNotes
	}
	if req.ResumeMatchScore != nil {
		application.ResumeMatchScore = *req.ResumeMatchScore
	}
	if req.InterviewPerformanceScore != nil {
		application.InterviewPerformanceScore = *req.InterviewPerformanceScore
	}
	if req.RoleInterestScore != nil {
		application.RoleInterestScore = *req.RoleInterestScore
	}
	if req.OverallScore != nil {
		application.OverallScore = *req.OverallScore
	}
	if req.Remark != nil {
		application.Remark = *req.Remark
	}

	if err := application.Validate(); err != nil {
		return err
	}
	if err := s.ensureReferencesOwned(ctx, userID, application.JobPostingID, application.RepoID, application.ResumeVersionID); err != nil {
		return err
	}
	return s.repo.Update(ctx, application)
}

func (s *ResumeApplicationService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *ResumeApplicationService) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	return s.repo.ExistsByResumeVersionID(ctx, versionID)
}

func (s *ResumeApplicationService) UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error {
	if !validStatus(status) {
		return domainerrors.ErrBadParams
	}
	return s.repo.UpdateCurrentStatus(ctx, id, userID, status)
}

func (s *ResumeApplicationService) ensureReferencesOwned(ctx context.Context, userID string, jobPostingID *string, repoID, resumeVersionID string) error {
	if jobPostingID != nil {
		if s.jobPostings == nil {
			return domainerrors.ErrBadParams
		}
		if _, err := s.jobPostings.FindByIDAndUserID(ctx, *jobPostingID, userID); err != nil {
			return err
		}
	}

	if s.resumeRepos == nil {
		return domainerrors.ErrBadParams
	}
	if _, err := s.resumeRepos.FindByIDAndUserID(ctx, repoID, userID); err != nil {
		return err
	}

	if s.resumeVersions == nil {
		return domainerrors.ErrBadParams
	}
	version, err := s.resumeVersions.FindByID(ctx, resumeVersionID)
	if err != nil {
		return err
	}
	_, err = s.resumeRepos.FindByIDAndUserID(ctx, version.RepoID, userID)
	return err
}

func normalizeOptionalJobPostingID(jobPostingID *string) *string {
	if jobPostingID == nil || strings.TrimSpace(*jobPostingID) == "" {
		return nil
	}
	return jobPostingID
}

func validStatus(status domainapp.Status) bool {
	switch status {
	case domainapp.StatusApplied,
		domainapp.StatusViewed,
		domainapp.StatusWrittenTest,
		domainapp.StatusFirstInterview,
		domainapp.StatusSecondInterview,
		domainapp.StatusFinalInterview,
		domainapp.StatusHRInterview,
		domainapp.StatusOffer,
		domainapp.StatusRejected,
		domainapp.StatusClosed:
		return true
	default:
		return false
	}
}
