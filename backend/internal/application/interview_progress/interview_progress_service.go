package interview_progress

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ApplicationStatusUpdater interface {
	UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error
}

type InterviewProgressService struct {
	repo domainprogress.Repository
	apps ApplicationStatusUpdater
}

func NewInterviewProgressService(repo domainprogress.Repository, apps ApplicationStatusUpdater) *InterviewProgressService {
	return &InterviewProgressService{repo: repo, apps: apps}
}

func (s *InterviewProgressService) Create(ctx context.Context, userID string, req dto.CreateInterviewProgressRequest) error {
	progress := &domainprogress.InterviewProgress{
		ID:                uuid.NewString(),
		UserID:            userID,
		ApplicationID:     req.ApplicationID,
		Stage:             req.Stage,
		OccurredAt:        req.OccurredAt,
		Result:            req.Result,
		InterviewerOrTeam: req.InterviewerOrTeam,
		Note:              req.Note,
	}
	if err := progress.Validate(); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, progress); err != nil {
		return err
	}
	return s.syncApplicationStatus(ctx, userID, progress.ApplicationID)
}

func (s *InterviewProgressService) List(ctx context.Context, userID, applicationID string) ([]domainprogress.InterviewProgress, error) {
	if strings.TrimSpace(applicationID) == "" {
		return nil, domainerrors.ErrBadParams
	}
	return s.repo.ListByApplicationAndUser(ctx, applicationID, userID)
}

func (s *InterviewProgressService) Update(ctx context.Context, userID, id string, req dto.UpdateInterviewProgressRequest) error {
	if strings.TrimSpace(id) == "" {
		return domainerrors.ErrBadParams
	}
	progress, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	oldApplicationID := progress.ApplicationID

	if req.ApplicationID != nil {
		progress.ApplicationID = *req.ApplicationID
	}
	if req.Stage != nil {
		progress.Stage = *req.Stage
	}
	if req.OccurredAt != nil {
		progress.OccurredAt = *req.OccurredAt
	}
	if req.Result != nil {
		progress.Result = *req.Result
	}
	if req.InterviewerOrTeam != nil {
		progress.InterviewerOrTeam = *req.InterviewerOrTeam
	}
	if req.Note != nil {
		progress.Note = *req.Note
	}

	if err := progress.Validate(); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, progress); err != nil {
		return err
	}
	if oldApplicationID != progress.ApplicationID {
		if err := s.syncApplicationStatus(ctx, userID, oldApplicationID); err != nil {
			return err
		}
	}
	return s.syncApplicationStatus(ctx, userID, progress.ApplicationID)
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
	return s.syncApplicationStatus(ctx, userID, progress.ApplicationID)
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
