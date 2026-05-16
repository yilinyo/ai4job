package postgres

import (
	"context"
	"errors"
	"strings"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type InterviewQuestionRepo struct {
	db *gorm.DB
}

var _ questionrepo.Repository = (*InterviewQuestionRepo)(nil)

func NewInterviewQuestionRepo(db *gorm.DB) *InterviewQuestionRepo {
	return &InterviewQuestionRepo{db: db}
}

func (r *InterviewQuestionRepo) Create(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	err := dbFromContext(ctx, r.db).WithContext(ctx).Create(q).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *InterviewQuestionRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error) {
	var q domainquestion.InterviewQuestion
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &q, err
}

func (r *InterviewQuestionRepo) Update(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&domainquestion.InterviewQuestion{}).
		Where("id = ? AND user_id = ?", q.ID, q.UserID).
		Updates(map[string]any{
			"title":            q.Title,
			"content":          q.Content,
			"answer_analysis":  q.AnswerAnalysis,
			"tags":             q.Tags,
			"difficulty":       q.Difficulty,
			"question_type":    q.QuestionType,
			"mastery_status":   q.MasteryStatus,
			"is_favorite":      q.IsFavorite,
			"source":           q.Source,
			"last_reviewed_at": q.LastReviewedAt,
			"remark":           q.Remark,
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

func (r *InterviewQuestionRepo) Delete(ctx context.Context, id, userID string) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainquestion.InterviewQuestion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *InterviewQuestionRepo) ListByUser(ctx context.Context, f questionrepo.ListFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	db := dbFromContext(ctx, r.db).WithContext(ctx)
	q := applyQuestionListFilters(db.Model(&domainquestion.InterviewQuestion{}), f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var questions []domainquestion.InterviewQuestion
	err := q.Order("updated_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&questions).Error
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.attachOccurrenceCounts(ctx, questions)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *InterviewQuestionRepo) CountOccurrences(ctx context.Context, questionID, userID string) (int64, error) {
	var count int64
	err := dbFromContext(ctx, r.db).WithContext(ctx).
		Table("question_occurrences").
		Where("question_id = ? AND user_id = ? AND deleted_at IS NULL", questionID, userID).
		Count(&count).Error
	return count, err
}

func (r *InterviewQuestionRepo) ListFlashcards(ctx context.Context, f questionrepo.FlashcardFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	db := dbFromContext(ctx, r.db).WithContext(ctx)
	q := applyFlashcardFilters(db.Model(&domainquestion.InterviewQuestion{}), f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	selectSQL := "interview_questions.*, COALESCE(occurrence_counts.occurrence_count, 0) AS occurrence_count"
	q = q.Select(selectSQL).
		Joins("LEFT JOIN (SELECT question_id, COUNT(*) AS occurrence_count FROM question_occurrences WHERE user_id = ? AND deleted_at IS NULL GROUP BY question_id) occurrence_counts ON occurrence_counts.question_id = interview_questions.id", f.UserID)

	q = applyFlashcardOrder(q, f.Mode)
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}

	var rows []questionWithOccurrenceRow
	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]questionrepo.QuestionWithOccurrenceCount, 0, len(rows))
	for _, row := range rows {
		items = append(items, questionrepo.QuestionWithOccurrenceCount{Question: row.InterviewQuestion, OccurrenceCount: row.OccurrenceCount})
	}
	return items, total, nil
}

type questionWithOccurrenceRow struct {
	domainquestion.InterviewQuestion
	OccurrenceCount int64 `gorm:"column:occurrence_count"`
}

func applyQuestionListFilters(q *gorm.DB, f questionrepo.ListFilter) *gorm.DB {
	q = q.Where("user_id = ?", f.UserID)
	if f.Keyword != "" {
		keyword := "%" + f.Keyword + "%"
		q = q.Where("title ILIKE ? OR content ILIKE ? OR answer_analysis ILIKE ? OR tags::text ILIKE ?", keyword, keyword, keyword, keyword)
	}
	q = applyTagsAny(q, f.Tags)
	if f.Difficulty != "" {
		q = q.Where("difficulty = ?", f.Difficulty)
	}
	if f.QuestionType != "" {
		q = q.Where("question_type = ?", f.QuestionType)
	}
	if f.MasteryStatus != "" {
		q = q.Where("mastery_status = ?", f.MasteryStatus)
	}
	if f.IsFavorite != nil {
		q = q.Where("is_favorite = ?", *f.IsFavorite)
	}
	return q
}

func applyFlashcardFilters(q *gorm.DB, f questionrepo.FlashcardFilter) *gorm.DB {
	q = q.Where("user_id = ?", f.UserID)
	q = applyTagsAny(q, f.Tags)
	if f.Difficulty != "" {
		q = q.Where("difficulty = ?", f.Difficulty)
	}
	return q
}

func applyTagsAny(q *gorm.DB, tags []string) *gorm.DB {
	if len(tags) == 0 {
		return q
	}
	conditions := make([]string, 0, len(tags))
	args := make([]any, 0, len(tags))
	for _, tag := range tags {
		conditions = append(conditions, "tags @> ?::jsonb")
		args = append(args, domainquestion.EncodeTags([]string{tag}).String())
	}
	return q.Where(strings.Join(conditions, " OR "), args...)
}

func applyFlashcardOrder(q *gorm.DB, mode domainreview.ReviewMode) *gorm.DB {
	switch mode {
	case domainreview.ReviewModeTag, domainreview.ReviewModeRandom:
		return q.Order("RANDOM()")
	case domainreview.ReviewModeDifficulty:
		return q.Order("interview_questions.updated_at DESC")
	case domainreview.ReviewModeFavoriteFirst:
		return q.Order("is_favorite DESC").Order("occurrence_count DESC")
	case domainreview.ReviewModeWeakFirst:
		return q.Order("CASE mastery_status WHEN 'unreviewed' THEN 1 WHEN 'weak' THEN 2 WHEN 'normal' THEN 3 WHEN 'mastered' THEN 4 ELSE 5 END ASC").Order("occurrence_count DESC")
	case domainreview.ReviewModeFrequentFirst:
		return q.Order("occurrence_count DESC").Order("is_favorite DESC")
	default:
		return q.Order("RANDOM()")
	}
}

func (r *InterviewQuestionRepo) attachOccurrenceCounts(ctx context.Context, questions []domainquestion.InterviewQuestion) ([]questionrepo.QuestionWithOccurrenceCount, error) {
	items := make([]questionrepo.QuestionWithOccurrenceCount, 0, len(questions))
	if len(questions) == 0 {
		return items, nil
	}

	ids := make([]string, 0, len(questions))
	for _, question := range questions {
		ids = append(ids, question.ID)
	}

	type countRow struct {
		QuestionID string
		Count      int64
	}
	var counts []countRow
	err := dbFromContext(ctx, r.db).WithContext(ctx).
		Table("question_occurrences").
		Select("question_id, COUNT(*) AS count").
		Where("question_id IN ? AND deleted_at IS NULL", ids).
		Group("question_id").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	countsByQuestionID := make(map[string]int64, len(counts))
	for _, row := range counts {
		countsByQuestionID[row.QuestionID] = row.Count
	}
	for _, question := range questions {
		items = append(items, questionrepo.QuestionWithOccurrenceCount{Question: question, OccurrenceCount: countsByQuestionID[question.ID]})
	}
	return items, nil
}
