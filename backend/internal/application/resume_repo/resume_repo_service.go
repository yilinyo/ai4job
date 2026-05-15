package resume_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
	resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type ResumeRepoService struct {
	repo resumereporepo.Repository
}

func NewResumeRepoService(repo resumereporepo.Repository) *ResumeRepoService {
	return &ResumeRepoService{repo: repo}
}

func (s *ResumeRepoService) Create(ctx context.Context, userID string, req dto.CreateResumeRepoRequest) error {
	exists, err := s.repo.ExistsByUserAndName(ctx, userID, req.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return domainerrors.ErrDuplicate
	}
	r := &domainrepo.ResumeRepo{
		ID:             uuid.NewString(),
		UserID:         userID,
		Name:           req.Name,
		JobType:        req.JobType,
		TargetPosition: req.TargetPosition,
		Description:    req.Description,
		VersionCount:   0,
	}
	if err := r.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, r)
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
	r.Name = req.Name
	r.JobType = req.JobType
	r.TargetPosition = req.TargetPosition
	r.Description = req.Description
	if err := r.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, r)
}

func (s *ResumeRepoService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *ResumeRepoService) IncrementVersionCount(ctx context.Context, id string, delta int) error {
	return s.repo.IncrementVersionCount(ctx, id, delta)
}
