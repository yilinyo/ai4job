package resume_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ResumeRepoService struct {
	repo           resumereporepo.Repository
	activityLogger activitylog.Logger
}

func NewResumeRepoService(repo resumereporepo.Repository) *ResumeRepoService {
	return &ResumeRepoService{repo: repo, activityLogger: activitylog.NoopLogger{}}
}

func (s *ResumeRepoService) WithActivityLogger(logger activitylog.Logger) *ResumeRepoService {
	s.activityLogger = activitylog.OrNoop(logger)
	return s
}

func (s *ResumeRepoService) Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error {
	exists, err := s.repo.ExistsByUserAndName(ctx, userID, req.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return domainerrors.ErrDuplicate
	}
	r, err := domainrepo.NewResumeRepo(domainrepo.NewResumeRepoParams{
		ID:             uuid.NewString(),
		UserID:         userID,
		Name:           req.Name,
		JobType:        req.JobType,
		TargetPosition: req.TargetPosition,
		Description:    req.Description,
	})
	if err != nil {
		return err
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return err
	}
	s.logRepoActivity(ctx, r, domainactivity.ActionCreated, "/resume/repo/"+r.ID)
	return nil
}

func (s *ResumeRepoService) List(ctx context.Context, userID string, jobType domainrepo.JobType, keyword string, page, pageSize int) ([]domainrepo.ResumeRepo, int64, error) {
	return s.repo.ListByUser(ctx, resumereporepo.ListFilter{UserID: userID, JobType: jobType, Keyword: keyword, Page: page, PageSize: pageSize})
}

func (s *ResumeRepoService) Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *ResumeRepoService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeRepoRequest) error {
	r, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	exists, err := s.repo.ExistsByUserAndName(ctx, userID, req.Name, id)
	if err != nil {
		return err
	}
	if exists {
		return domainerrors.ErrDuplicate
	}
	if err := r.ApplyUpdate(domainrepo.UpdateSpec{
		Name:           req.Name,
		JobType:        req.JobType,
		TargetPosition: req.TargetPosition,
		Description:    req.Description,
	}); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, r); err != nil {
		return err
	}
	s.logRepoActivity(ctx, r, domainactivity.ActionUpdated, "/resume/repo/"+r.ID)
	return nil
}

func (s *ResumeRepoService) Delete(ctx context.Context, userID, id string) error {
	r, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	s.logRepoActivity(ctx, r, domainactivity.ActionDeleted, "")
	return nil
}

func (s *ResumeRepoService) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return s.repo.IncrementVersionCount(ctx, id, delta)
}

func (s *ResumeRepoService) logRepoActivity(ctx context.Context, r *domainrepo.ResumeRepo, action domainactivity.Action, route string) {
	verb := "已更新"
	if action == domainactivity.ActionCreated {
		verb = "已创建"
	} else if action == domainactivity.ActionDeleted {
		verb = "已删除"
	}
	activitylog.BestEffortLog(ctx, s.activityLogger, activitylog.Entry{
		UserID:      r.UserID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeResume,
		TargetID:    r.ID,
		TargetLabel: activitylog.TrimLabel(r.Name),
		Summary:     activitylog.TrimSummary(activitylog.TrimLabel(r.Name) + " 简历仓库" + verb),
		Route:       route,
	})
}
