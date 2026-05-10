package resume_version

import (
	"context"

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

func (s *ResumeVersionService) Create(ctx context.Context, userID, repoID string, req dto.CreateResumeVersionRequest) error {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return err
	}
	maxNum, err := s.repo.MaxVersionNum(ctx, repoID)
	if err != nil {
		return err
	}
	v := &domainversion.ResumeVersion{ID: uuid.NewString(), RepoID: repoID, Title: req.Title, Content: req.Content, VersionNum: maxNum + 1}
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return err
	}
	if req.IsDefault {
		return s.repo.SetDefaultByID(ctx, repoID, v.ID)
	}
	return nil
}

func (s *ResumeVersionService) List(ctx context.Context, userID, repoID string) ([]domainversion.ResumeVersion, error) {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return nil, err
	}
	return s.repo.ListByRepoID(ctx, repoID)
}

func (s *ResumeVersionService) Get(ctx context.Context, userID, repoID, id string) (*domainversion.ResumeVersion, error) {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return nil, err
	}
	return s.repo.FindByIDAndRepoID(ctx, id, repoID)
}

func (s *ResumeVersionService) Update(ctx context.Context, userID, repoID, id string, req dto.UpdateResumeVersionRequest) error {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return err
	}
	v, err := s.repo.FindByIDAndRepoID(ctx, id, repoID)
	if err != nil {
		return err
	}
	v.Title = req.Title
	v.Content = req.Content
	if err := v.ValidateForCreate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, v)
}

func (s *ResumeVersionService) Delete(ctx context.Context, userID, repoID, id string) error {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return err
	}
	if _, err := s.repo.FindByIDAndRepoID(ctx, id, repoID); err != nil {
		return err
	}
	inUse, err := s.checker.ExistsByResumeVersionID(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return domainerrors.ErrVersionInUse
	}
	return s.repo.Delete(ctx, id, repoID)
}

func (s *ResumeVersionService) SetDefault(ctx context.Context, userID, repoID, id string) error {
	if err := s.ensureRepoOwned(ctx, userID, repoID); err != nil {
		return err
	}
	return s.repo.SetDefaultByID(ctx, repoID, id)
}
