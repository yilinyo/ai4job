package postgres

import (
	"context"
	"errors"

	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	occurrencerepo "github.com/yilin/ai-for-backend/internal/domain/question_occurrence/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type QuestionOccurrenceRepo struct {
	db *gorm.DB
}

var _ occurrencerepo.Repository = (*QuestionOccurrenceRepo)(nil)

func NewQuestionOccurrenceRepo(db *gorm.DB) *QuestionOccurrenceRepo {
	return &QuestionOccurrenceRepo{db: db}
}

func (r *QuestionOccurrenceRepo) Create(ctx context.Context, o *domainoccurrence.QuestionOccurrence) error {
	err := dbFromContext(ctx, r.db).WithContext(ctx).Create(o).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *QuestionOccurrenceRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error) {
	var o domainoccurrence.QuestionOccurrence
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &o, err
}

func (r *QuestionOccurrenceRepo) Delete(ctx context.Context, id, userID string) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainoccurrence.QuestionOccurrence{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *QuestionOccurrenceRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Delete(&domainoccurrence.QuestionOccurrence{}).Error
}

func (r *QuestionOccurrenceRepo) ListByUser(ctx context.Context, f occurrencerepo.ListFilter) ([]domainoccurrence.QuestionOccurrence, error) {
	q := dbFromContext(ctx, r.db).WithContext(ctx).Model(&domainoccurrence.QuestionOccurrence{}).Where("user_id = ?", f.UserID)
	if f.QuestionID != "" {
		q = q.Where("question_id = ?", f.QuestionID)
	}
	if f.ApplicationID != "" {
		q = q.Where("application_id = ?", f.ApplicationID)
	}

	var rows []domainoccurrence.QuestionOccurrence
	err := q.Order("occurred_at DESC").Find(&rows).Error
	return rows, err
}
