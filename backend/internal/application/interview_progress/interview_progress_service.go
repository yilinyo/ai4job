package interview_progress

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	interviewprogressrepo "github.com/yilin/ai-for-backend/internal/domain/interview_progress/repository"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ApplicationReader interface {
	Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error)
}

type ApplicationStatusUpdater interface {
	UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error
}

type ApplicationService interface {
	ApplicationReader
	ApplicationStatusUpdater
}

type InterviewProgressService struct {
	repo           interviewprogressrepo.Repository
	apps           ApplicationService
	activityLogger activitylog.Logger
}

func NewInterviewProgressService(repo interviewprogressrepo.Repository, apps ApplicationService) *InterviewProgressService {
	return &InterviewProgressService{repo: repo, apps: apps, activityLogger: activitylog.NoopLogger{}}
}

func (s *InterviewProgressService) WithActivityLogger(logger activitylog.Logger) *InterviewProgressService {
	s.activityLogger = activitylog.OrNoop(logger)
	return s
}

func (s *InterviewProgressService) Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	progress, err := domainprogress.NewInterviewProgress(domainprogress.NewInterviewProgressParams{
		ID:                uuid.NewString(),
		UserID:            userID,
		ApplicationID:     req.ApplicationID,
		Stage:             req.Stage,
		OccurredAt:        req.OccurredAt.Value(),
		Result:            req.Result,
		InterviewerOrTeam: req.InterviewerOrTeam,
		Note:              req.Note,
	})
	if err != nil {
		return nil, err
	}
	if err := s.ensureApplicationOwned(ctx, userID, progress.ApplicationID); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, progress); err != nil {
		return nil, err
	}
	if err := s.syncApplicationStatus(ctx, userID, progress.ApplicationID); err != nil {
		return nil, err
	}
	if err := s.logProgressActivity(ctx, progress, domainactivity.ActionCreated); err != nil {
		return nil, err
	}
	return progress, nil
}

func (s *InterviewProgressService) List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error) {
	if strings.TrimSpace(applicationID) == "" {
		return nil, domainerrors.ErrBadParams
	}
	return s.repo.ListByApplicationAndUser(ctx, applicationID, userID)
}

func (s *InterviewProgressService) Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) (*domainprogress.InterviewProgress, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domainerrors.ErrBadParams
	}
	progress, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	oldApplicationID := progress.ApplicationID

	var occurredAt *time.Time
	if req.OccurredAt != nil {
		value := req.OccurredAt.Value()
		occurredAt = &value
	}

	if err := progress.ApplyUpdate(domainprogress.UpdateSpec{
		ApplicationID:     req.ApplicationID,
		Stage:             req.Stage,
		OccurredAt:        occurredAt,
		Result:            req.Result,
		InterviewerOrTeam: req.InterviewerOrTeam,
		Note:              req.Note,
	}); err != nil {
		return nil, err
	}
	if req.ApplicationID != nil && oldApplicationID != progress.ApplicationID {
		if err := s.ensureApplicationOwned(ctx, userID, progress.ApplicationID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Update(ctx, progress); err != nil {
		return nil, err
	}
	if oldApplicationID != progress.ApplicationID {
		if err := s.syncApplicationStatus(ctx, userID, oldApplicationID); err != nil {
			return nil, err
		}
	}
	if err := s.syncApplicationStatus(ctx, userID, progress.ApplicationID); err != nil {
		return nil, err
	}
	if err := s.logProgressActivity(ctx, progress, domainactivity.ActionUpdated); err != nil {
		return nil, err
	}
	return progress, nil
}

func (s *InterviewProgressService) Delete(ctx context.Context, userID, id string) error {
	if strings.TrimSpace(id) == "" {
		return domainerrors.ErrBadParams
	}
	progress, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	if err := s.syncApplicationStatus(ctx, userID, progress.ApplicationID); err != nil {
		return err
	}
	return s.logProgressActivity(ctx, progress, domainactivity.ActionDeleted)
}

func (s *InterviewProgressService) logProgressActivity(ctx context.Context, progress *domainprogress.InterviewProgress, action domainactivity.Action) error {
	label := activitylog.TrimLabel(stageLabel(progress.Stage) + " / " + resultLabel(progress.Result))
	verb := "已更新"
	if action == domainactivity.ActionCreated {
		verb = "已创建"
	} else if action == domainactivity.ActionDeleted {
		verb = "已删除"
	}
	return activitylog.OrNoop(s.activityLogger).Log(ctx, activitylog.Entry{
		UserID:      progress.UserID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeInterview,
		TargetID:    progress.ID,
		TargetLabel: label,
		Summary:     activitylog.TrimSummary(label + " 面试进展" + verb),
		Route:       "/resume/application/" + progress.ApplicationID,
	})
}

func stageLabel(stage domainprogress.Stage) string {
	switch stage {
	case domainprogress.StageApplied:
		return "已投递"
	case domainprogress.StageViewed:
		return "已查看"
	case domainprogress.StageWrittenTest:
		return "笔试"
	case domainprogress.StageFirstInterview:
		return "一面"
	case domainprogress.StageSecondInterview:
		return "二面"
	case domainprogress.StageFinalInterview:
		return "终面"
	case domainprogress.StageHRInterview:
		return "HR面"
	case domainprogress.StageOffer:
		return "Offer"
	case domainprogress.StageRejected:
		return "已拒绝"
	case domainprogress.StageClosed:
		return "已关闭"
	default:
		return string(stage)
	}
}

func resultLabel(result domainprogress.Result) string {
	switch result {
	case domainprogress.ResultPassed:
		return "通过"
	case domainprogress.ResultPending:
		return "待定"
	case domainprogress.ResultFailed:
		return "未通过"
	default:
		return string(result)
	}
}

func (s *InterviewProgressService) ensureApplicationOwned(ctx context.Context, userID, applicationID string) error {
	if s.apps == nil {
		return domainerrors.ErrBadParams
	}
	_, err := s.apps.Get(ctx, userID, applicationID)
	return err
}

func (s *InterviewProgressService) syncApplicationStatus(ctx context.Context, userID, applicationID string) error {
	if s.apps == nil {
		return domainerrors.ErrBadParams
	}
	progresses, err := s.repo.ListByApplicationAndUser(ctx, applicationID, userID)
	if err != nil {
		return err
	}
	if len(progresses) == 0 {
		return s.apps.UpdateCurrentStatus(ctx, userID, applicationID, domainapp.StatusApplied)
	}

	latest := progresses[0]
	for _, progress := range progresses[1:] {
		if progress.OccurredAt.After(latest.OccurredAt) {
			latest = progress
		}
	}

	status := domainapp.Status(latest.Stage)
	if latest.Result == domainprogress.ResultFailed {
		status = domainapp.StatusRejected
	}
	return s.apps.UpdateCurrentStatus(ctx, userID, applicationID, status)
}
