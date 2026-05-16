package postgres

import (
	"context"
	"errors"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	interviewprogressrepo "github.com/yilin/ai-for-backend/internal/domain/interview_progress/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type InterviewProgressRepo struct {
	db *gorm.DB
}

var _ interviewprogressrepo.Repository = (*InterviewProgressRepo)(nil)

func NewInterviewProgressRepo(db *gorm.DB) *InterviewProgressRepo {
	return &InterviewProgressRepo{db: db}
}

func (r *InterviewProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error {
	err := r.db.WithContext(ctx).Create(p).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *InterviewProgressRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) {
	var p domainprogress.InterviewProgress
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &p, err
}

func (r *InterviewProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error {
	result := r.db.WithContext(ctx).
		Model(&domainprogress.InterviewProgress{}).
		Where("id = ? AND user_id = ?", p.ID, p.UserID).
		Updates(map[string]any{
			"application_id":      p.ApplicationID,
			"stage":               p.Stage,
			"occurred_at":         p.OccurredAt,
			"result":              p.Result,
			"interviewer_or_team": p.InterviewerOrTeam,
			"note":                p.Note,
		})
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return domainerrors.ErrDuplicate
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *InterviewProgressRepo) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainprogress.InterviewProgress{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *InterviewProgressRepo) ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error) {
	var rows []domainprogress.InterviewProgress
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND user_id = ?", applicationID, userID).
		Order("occurred_at ASC").
		Find(&rows).Error
	return rows, err
}

func (r *InterviewProgressRepo) DeleteByApplicationID(ctx context.Context, applicationID, userID string) error {
	return r.db.WithContext(ctx).
		Where("application_id = ? AND user_id = ?", applicationID, userID).
		Delete(&domainprogress.InterviewProgress{}).Error
}
