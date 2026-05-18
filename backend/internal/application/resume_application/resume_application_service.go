package resume_application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	resumeapplicationrepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
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
	repo           resumeapplicationrepo.Repository
	jobPostings    JobPostingReader
	resumeRepos    ResumeRepoReader
	resumeVersions ResumeVersionReader
	activityLogger activitylog.Logger
}

func NewResumeApplicationService(repo resumeapplicationrepo.Repository, jobPostings JobPostingReader, resumeRepos ResumeRepoReader, resumeVersions ResumeVersionReader) *ResumeApplicationService {
	return &ResumeApplicationService{repo: repo, jobPostings: jobPostings, resumeRepos: resumeRepos, resumeVersions: resumeVersions, activityLogger: activitylog.NoopLogger{}}
}

func (s *ResumeApplicationService) WithActivityLogger(logger activitylog.Logger) *ResumeApplicationService {
	s.activityLogger = activitylog.OrNoop(logger)
	return s
}

func (s *ResumeApplicationService) Create(ctx context.Context, userID string, req dto.CreateResumeApplicationRequest) (*domainapp.ResumeApplication, error) {
	application, err := domainapp.NewResumeApplication(domainapp.NewResumeApplicationParams{
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
		AppliedAt:                 req.AppliedAt.Value(),
		CurrentStatus:             req.CurrentStatus,
		InterviewSummary:          req.InterviewSummary,
		InterviewNotes:            req.InterviewNotes,
		ResumeMatchScore:          req.ResumeMatchScore,
		InterviewPerformanceScore: req.InterviewPerformanceScore,
		RoleInterestScore:         req.RoleInterestScore,
		OverallScore:              req.OverallScore,
		Remark:                    req.Remark,
	})
	if err != nil {
		return nil, err
	}
	if err := s.ensureReferencesOwned(ctx, userID, application.JobPostingID, application.RepoID, application.ResumeVersionID); err != nil {
		return nil, err
	}
	if err := s.ensureResumeVersionActive(ctx, application.ResumeVersionID); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, application); err != nil {
		return nil, err
	}
	if err := s.logApplicationActivity(ctx, application, domainactivity.ActionCreated, applicationActivitySummary(application, domainactivity.ActionCreated)); err != nil {
		return nil, err
	}
	return application, nil
}

func (s *ResumeApplicationService) List(ctx context.Context, userID, resumeVersionID, repoID, keyword, currentStatus string, page, pageSize int) ([]domainapp.ResumeApplication, int64, error) {
	filterStatus := domainapp.Status(currentStatus)
	if currentStatus != "" && !filterStatus.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}

	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return s.repo.ListByUser(ctx, resumeapplicationrepo.ListFilter{
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
	oldStatus := application.CurrentStatus

	var appliedAt *time.Time
	if req.AppliedAt != nil {
		value := req.AppliedAt.Value()
		appliedAt = &value
	}
	var jobPostingID **string
	if req.JobPostingID != nil {
		value := normalizeOptionalJobPostingID(req.JobPostingID)
		jobPostingID = &value
	}

	if err := application.ApplyUpdate(domainapp.UpdateSpec{
		JobPostingID:              jobPostingID,
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
		AppliedAt:                 appliedAt,
		CurrentStatus:             req.CurrentStatus,
		InterviewSummary:          req.InterviewSummary,
		InterviewNotes:            req.InterviewNotes,
		ResumeMatchScore:          req.ResumeMatchScore,
		InterviewPerformanceScore: req.InterviewPerformanceScore,
		RoleInterestScore:         req.RoleInterestScore,
		OverallScore:              req.OverallScore,
		Remark:                    req.Remark,
	}); err != nil {
		return err
	}
	if err := s.ensureReferencesOwned(ctx, userID, application.JobPostingID, application.RepoID, application.ResumeVersionID); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, application); err != nil {
		return err
	}
	action := domainactivity.ActionUpdated
	summary := applicationActivitySummary(application, action)
	if oldStatus != application.CurrentStatus {
		action = domainactivity.ActionStatusChanged
		summary = applicationStatusChangedSummary(application)
	}
	return s.logApplicationActivity(ctx, application, action, summary)
}

func (s *ResumeApplicationService) Delete(ctx context.Context, userID, id string) error {
	application, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	return s.logApplicationActivity(ctx, application, domainactivity.ActionDeleted, applicationActivitySummary(application, domainactivity.ActionDeleted))
}

func (s *ResumeApplicationService) ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error) {
	return s.repo.ExistsByResumeVersionID(ctx, versionID)
}

func (s *ResumeApplicationService) UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error {
	application, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	oldStatus := application.CurrentStatus
	if err := application.UpdateCurrentStatus(status); err != nil {
		return err
	}
	if err := s.repo.UpdateCurrentStatus(ctx, id, userID, application.CurrentStatus); err != nil {
		return err
	}
	if oldStatus == application.CurrentStatus {
		return s.logApplicationActivity(ctx, application, domainactivity.ActionUpdated, applicationActivitySummary(application, domainactivity.ActionUpdated))
	}
	return s.logApplicationActivity(ctx, application, domainactivity.ActionStatusChanged, applicationStatusChangedSummary(application))
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

func (s *ResumeApplicationService) ensureResumeVersionActive(ctx context.Context, resumeVersionID string) error {
	if s.resumeVersions == nil {
		return domainerrors.ErrBadParams
	}
	version, err := s.resumeVersions.FindByID(ctx, resumeVersionID)
	if err != nil {
		return err
	}
	if !version.IsActive {
		return domainerrors.ErrInactiveVersion
	}
	return nil
}

func (s *ResumeApplicationService) logApplicationActivity(ctx context.Context, application *domainapp.ResumeApplication, action domainactivity.Action, summary string) error {
	return activitylog.OrNoop(s.activityLogger).Log(ctx, activitylog.Entry{
		UserID:      application.UserID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeApplication,
		TargetID:    application.ID,
		TargetLabel: applicationLabel(application),
		Summary:     activitylog.TrimSummary(summary),
		Route:       "/resume/application/" + application.ID,
	})
}

func applicationLabel(application *domainapp.ResumeApplication) string {
	return activitylog.TrimLabel(application.CompanyName + " / " + application.JobTitle)
}

func applicationActivitySummary(application *domainapp.ResumeApplication, action domainactivity.Action) string {
	switch action {
	case domainactivity.ActionCreated:
		return applicationLabel(application) + " 投递已创建"
	case domainactivity.ActionDeleted:
		return applicationLabel(application) + " 投递已删除"
	default:
		return applicationLabel(application) + " 投递已更新"
	}
}

func applicationStatusChangedSummary(application *domainapp.ResumeApplication) string {
	return applicationLabel(application) + " 投递状态更新为" + statusLabel(application.CurrentStatus)
}

func statusLabel(status domainapp.Status) string {
	switch status {
	case domainapp.StatusApplied:
		return "已投递"
	case domainapp.StatusViewed:
		return "已查看"
	case domainapp.StatusWrittenTest:
		return "笔试"
	case domainapp.StatusFirstInterview:
		return "一面"
	case domainapp.StatusSecondInterview:
		return "二面"
	case domainapp.StatusFinalInterview:
		return "终面"
	case domainapp.StatusHRInterview:
		return "HR面"
	case domainapp.StatusOffer:
		return "Offer"
	case domainapp.StatusRejected:
		return "已拒绝"
	case domainapp.StatusClosed:
		return "已关闭"
	default:
		return string(status)
	}
}

func normalizeOptionalJobPostingID(jobPostingID *string) *string {
	if jobPostingID == nil || strings.TrimSpace(*jobPostingID) == "" {
		return nil
	}
	return jobPostingID
}
