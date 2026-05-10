package resume_version

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
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
	repoReader ResumeRepoReader
	repo       domainversion.Repository
	checker    VersionUseChecker
}

func NewResumeVersionService(repoReader ResumeRepoReader, repo domainversion.Repository, checker VersionUseChecker) *ResumeVersionService {
	return &ResumeVersionService{repoReader: repoReader, repo: repo, checker: checker}
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
	v := &domainversion.ResumeVersion{
		ID:       uuid.NewString(),
		RepoID:   req.RepoID,
		Version:  fmt.Sprintf("v1.%d.0", count),
		Title:    req.Title,
		Content:  req.Content,
		Remark:   req.Remark,
		IsActive: false,
	}
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return err
	}
	return s.repoReader.IncrementVersionCount(ctx, req.RepoID, 1)
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
	v.Title = req.Title
	v.Content = req.Content
	v.Remark = req.Remark
	if req.IsActive != nil && !*req.IsActive {
		v.IsActive = false
	}
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, v); err != nil {
		return err
	}
	if activate {
		return s.repo.SetActiveByID(ctx, v.RepoID, id)
	}
	return nil
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
	return s.repoReader.IncrementVersionCount(ctx, v.RepoID, -1)
}

func (s *ResumeVersionService) Activate(ctx context.Context, userID, id string) error {
	v, err := s.ensureVersionOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	return s.repo.SetActiveByID(ctx, v.RepoID, id)
}
