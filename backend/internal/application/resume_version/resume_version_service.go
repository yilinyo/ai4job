package resume_version

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	resumeversionrepo "github.com/yilin/ai-for-backend/internal/domain/resume_version/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type VersionUseChecker interface {
	ExistsByResumeVersionID(ctx context.Context, versionID string) (bool, error)
}

type ResumeRepoReader interface {
	Get(ctx context.Context, userID, id string) (*domainrepo.ResumeRepo, error)
	IncrementVersionCount(ctx context.Context, id string, delta int) error
}

type ResumeVersionService struct {
	repoReader     ResumeRepoReader
	repo           resumeversionrepo.Repository
	checker        VersionUseChecker
	activityLogger activitylog.Logger
}

func NewResumeVersionService(repoReader ResumeRepoReader, repo resumeversionrepo.Repository, checker VersionUseChecker) *ResumeVersionService {
	return &ResumeVersionService{repoReader: repoReader, repo: repo, checker: checker, activityLogger: activitylog.NoopLogger{}}
}

func (s *ResumeVersionService) WithActivityLogger(logger activitylog.Logger) *ResumeVersionService {
	s.activityLogger = activitylog.OrNoop(logger)
	return s
}

func (s *ResumeVersionService) ensureRepoOwned(ctx context.Context, userID, repoID string) error {
	_, err := s.repoReader.Get(ctx, userID, repoID)
	return err
}

func (s *ResumeVersionService) ensureVersionOwned(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.ensureRepoOwned(ctx, userID, v.RepoID); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *ResumeVersionService) Create(ctx context.Context, userID string, req dto.CreateResumeVersionRequest) error {
	if err := s.ensureRepoOwned(ctx, userID, req.RepoID); err != nil {
		return err
	}
	count, err := s.repo.CountByRepoID(ctx, req.RepoID)
	if err != nil {
		return err
	}
	v, err := domainversion.NewResumeVersion(domainversion.NewResumeVersionParams{
		ID:       uuid.NewString(),
		RepoID:   req.RepoID,
		Version:  fmt.Sprintf("v1.%d.0", count),
		Title:    req.Title,
		Content:  req.Content,
		Remark:   req.Remark,
		IsActive: false,
	})
	if err != nil {
		return err
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return err
	}
	if err := s.repoReader.IncrementVersionCount(ctx, req.RepoID, 1); err != nil {
		return err
	}
	return s.logVersionActivity(ctx, userID, v, domainactivity.ActionCreated, "/resume/version/"+v.ID)
}

func (s *ResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return nil, err
	}
	return s.repo.ListByRepoID(ctx, repoID)
}

func (s *ResumeVersionService) Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error) {
	return s.ensureVersionOwned(ctx, userID, id)
}

func (s *ResumeVersionService) Update(ctx context.Context, userID, id string, req dto.UpdateResumeVersionRequest) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	activate := req.IsActive != nil && *req.IsActive
	if err := v.ApplyUpdate(domainversion.UpdateSpec{Title: req.Title, Content: req.Content, Remark: req.Remark, IsActive: req.IsActive}); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, v); err != nil {
		return err
	}
	if activate {
		if err := s.repo.SetActiveByID(ctx, v.RepoID, id); err != nil {
			return err
		}
		return s.logVersionActivity(ctx, userID, v, domainactivity.ActionStatusChanged, "/resume/version/"+v.ID)
	}
	return s.logVersionActivity(ctx, userID, v, domainactivity.ActionUpdated, "/resume/version/"+v.ID)
}

func (s *ResumeVersionService) Delete(ctx context.Context, userID, id string) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	inUse, err := s.checker.ExistsByResumeVersionID(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return domainerrors.ErrVersionInUse
	}
	if err := s.repo.Delete(ctx, id, v.RepoID); err != nil {
		return err
	}
	if err := s.repoReader.IncrementVersionCount(ctx, v.RepoID, -1); err != nil {
		return err
	}
	return s.logVersionActivity(ctx, userID, v, domainactivity.ActionDeleted, "")
}

func (s *ResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	if err := s.repo.SetActiveByID(ctx, v.RepoID, id); err != nil {
		return err
	}
	return s.logVersionActivity(ctx, userID, v, domainactivity.ActionStatusChanged, "/resume/version/"+v.ID)
}

func (s *ResumeVersionService) logVersionActivity(ctx context.Context, userID string, v *domainversion.ResumeVersion, action domainactivity.Action, route string) error {
	verb := "已更新"
	if action == domainactivity.ActionCreated {
		verb = "已创建"
	} else if action == domainactivity.ActionDeleted {
		verb = "已删除"
	} else if action == domainactivity.ActionStatusChanged {
		verb = "已激活"
	}
	label := activitylog.TrimLabel(v.Title)
	return activitylog.OrNoop(s.activityLogger).Log(ctx, activitylog.Entry{
		UserID:      userID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeResume,
		TargetID:    v.ID,
		TargetLabel: label,
		Summary:     activitylog.TrimSummary(label + " 简历版本" + verb),
		Route:       route,
	})
}
