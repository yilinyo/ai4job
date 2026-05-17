package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	applicationquestion "github.com/yilin/ai-for-backend/internal/application/interview_question"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	reviewrepo "github.com/yilin/ai-for-backend/internal/domain/flashcard_review/repository"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	occurrencerepo "github.com/yilin/ai-for-backend/internal/domain/question_occurrence/repository"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

func TestInterviewQuestionRepositoryContracts(t *testing.T) {
	var _ questionrepo.Repository = (*infrapostgres.InterviewQuestionRepo)(nil)
	var _ occurrencerepo.Repository = (*infrapostgres.QuestionOccurrenceRepo)(nil)
	var _ reviewrepo.Repository = (*infrapostgres.FlashcardReviewRepo)(nil)
	var _ applicationquestion.TransactionRunner = (*infrapostgres.TransactionRunner)(nil)
}

func migrateInterviewQuestionTables(t *testing.T, db interface {
	AutoMigrate(...interface{}) error
	Exec(string, ...interface{}) *gorm.DB
}) {
	t.Helper()
	require.NoError(t, db.AutoMigrate(&domainquestion.InterviewQuestion{}, &domainoccurrence.QuestionOccurrence{}, &domainreview.FlashcardReview{}))
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_question_occurrences_user_question_application_progress_unique ON question_occurrences(user_id, question_id, application_id, interview_progress_id) WHERE interview_progress_id IS NOT NULL AND deleted_at IS NULL").Error)
}

func testInterviewQuestion(id, userID, title string, tags []string) *domainquestion.InterviewQuestion {
	question, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             id,
		UserID:         userID,
		Title:          title,
		Content:        title + " content",
		AnswerAnalysis: title + " analysis",
		Tags:           tags,
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusUnreviewed,
	})
	if err != nil {
		panic(err)
	}
	return question
}

func testOccurrence(id, userID, questionID, applicationID string, occurredAt time.Time) *domainoccurrence.QuestionOccurrence {
	occurrence, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                  id,
		UserID:              userID,
		QuestionID:          questionID,
		ApplicationID:       applicationID,
		CompanyNameSnapshot: "Example Co",
		JobTitleSnapshot:    "Backend Engineer",
		OccurredAt:          occurredAt,
	})
	if err != nil {
		panic(err)
	}
	return occurrence
}

func testFlashcardReview(id, userID, questionID string, reviewedAt time.Time, result domainreview.Result) *domainreview.FlashcardReview {
	review, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{
		ID:         id,
		UserID:     userID,
		QuestionID: questionID,
		ReviewedAt: reviewedAt,
		Result:     result,
		ReviewMode: domainreview.ReviewModeRandom,
	})
	if err != nil {
		panic(err)
	}
	return review
}

func TestIntegrationInterviewQuestionRepo_ListByUserFiltersTagsAndCountsOccurrences(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	questions := infrapostgres.NewInterviewQuestionRepo(db)
	occurrences := infrapostgres.NewQuestionOccurrenceRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	applicationID := uuid.NewString()

	otherUserID := uuid.NewString()

	goQuestion := testInterviewQuestion(uuid.NewString(), userID, "Go channels", []string{"go", "concurrency"})
	sqlQuestion := testInterviewQuestion(uuid.NewString(), userID, "SQL joins", []string{"sql"})
	require.NoError(t, questions.Create(ctx, goQuestion))
	require.NoError(t, questions.Create(ctx, sqlQuestion))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), userID, goQuestion.ID, applicationID, time.Now())))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), userID, goQuestion.ID, uuid.NewString(), time.Now().Add(time.Second))))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), otherUserID, goQuestion.ID, uuid.NewString(), time.Now().Add(2*time.Second))))

	rows, total, err := questions.ListByUser(ctx, questionrepo.ListFilter{UserID: userID, Tags: []string{"go"}, Page: 1, PageSize: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, rows, 1)
	assert.Equal(t, goQuestion.ID, rows[0].Question.ID)
	assert.Equal(t, int64(2), rows[0].OccurrenceCount)
}

func TestIntegrationInterviewQuestionRepo_CountOccurrencesIgnoresOtherUsers(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	questions := infrapostgres.NewInterviewQuestionRepo(db)
	occurrences := infrapostgres.NewQuestionOccurrenceRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	otherUserID := uuid.NewString()
	question := testInterviewQuestion(uuid.NewString(), userID, "Scoped count", []string{"scope"})
	require.NoError(t, questions.Create(ctx, question))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), userID, question.ID, uuid.NewString(), time.Now())))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), otherUserID, question.ID, uuid.NewString(), time.Now().Add(time.Second))))

	count, err := questions.CountOccurrences(ctx, question.ID, userID)

	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestIntegrationInterviewQuestionRepo_ListFlashcardsIgnoresOtherUsersOccurrences(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	questions := infrapostgres.NewInterviewQuestionRepo(db)
	occurrences := infrapostgres.NewQuestionOccurrenceRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	otherUserID := uuid.NewString()
	question := testInterviewQuestion(uuid.NewString(), userID, "Flashcard scoped count", []string{"flashcard"})
	require.NoError(t, questions.Create(ctx, question))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), userID, question.ID, uuid.NewString(), time.Now())))
	require.NoError(t, occurrences.Create(ctx, testOccurrence(uuid.NewString(), otherUserID, question.ID, uuid.NewString(), time.Now().Add(time.Second))))

	rows, total, err := questions.ListFlashcards(ctx, questionrepo.FlashcardFilter{UserID: userID, Mode: domainreview.ReviewModeFrequentFirst, Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, rows, 1)
	assert.Equal(t, question.ID, rows[0].Question.ID)
	assert.Equal(t, int64(1), rows[0].OccurrenceCount)
}

func TestIntegrationInterviewQuestionRepo_UpdatePersistsMutableFields(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	repo := infrapostgres.NewInterviewQuestionRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	question := testInterviewQuestion(uuid.NewString(), userID, "Original", []string{"old"})
	require.NoError(t, repo.Create(ctx, question))

	reviewedAt := time.Now().UTC().Truncate(time.Microsecond)
	question.Title = "Updated"
	question.Content = "Updated content"
	question.AnswerAnalysis = "Updated analysis"
	question.Tags = domainquestion.EncodeTags([]string{"new"})
	question.Difficulty = domainquestion.DifficultyHard
	question.QuestionType = domainquestion.QuestionTypeScenario
	question.MasteryStatus = domainquestion.MasteryStatusWeak
	question.IsFavorite = true
	question.Source = "manual"
	question.LastReviewedAt = &reviewedAt
	question.Remark = "note"
	require.NoError(t, repo.Update(ctx, question))

	found, err := repo.FindByIDAndUserID(ctx, question.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", found.Title)
	assert.Equal(t, domainquestion.DifficultyHard, found.Difficulty)
	assert.Equal(t, domainquestion.QuestionTypeScenario, found.QuestionType)
	assert.Equal(t, domainquestion.MasteryStatusWeak, found.MasteryStatus)
	assert.True(t, found.IsFavorite)
	assert.Equal(t, "manual", found.Source)
	assert.Equal(t, "note", found.Remark)
	assert.Equal(t, reviewedAt, found.LastReviewedAt.UTC().Truncate(time.Microsecond))
	tags, err := domainquestion.DecodeTags(found.Tags)
	require.NoError(t, err)
	assert.Equal(t, []string{"new"}, tags)
}

func TestIntegrationFlashcardReviewRepo_SummaryExcludesSoftDeletedReviews(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	questions := infrapostgres.NewInterviewQuestionRepo(db)
	reviews := infrapostgres.NewFlashcardReviewRepo(db)
	ctx := context.Background()
	userID := uuid.NewString()
	question := testInterviewQuestion(uuid.NewString(), userID, "Review", []string{"review"})
	require.NoError(t, questions.Create(ctx, question))

	base := time.Now().UTC().Truncate(time.Microsecond)
	failed := testFlashcardReview(uuid.NewString(), userID, question.ID, base, domainreview.ResultFailed)
	vague := testFlashcardReview(uuid.NewString(), userID, question.ID, base.Add(time.Minute), domainreview.ResultVague)
	mastered := testFlashcardReview(uuid.NewString(), userID, question.ID, base.Add(2*time.Minute), domainreview.ResultMastered)
	require.NoError(t, reviews.Create(ctx, failed))
	require.NoError(t, reviews.Create(ctx, vague))
	require.NoError(t, reviews.Create(ctx, mastered))
	require.NoError(t, db.Delete(vague).Error)

	summary, err := reviews.SummaryByQuestion(ctx, question.ID, userID)

	require.NoError(t, err)
	require.NotNil(t, summary)
	assert.Equal(t, int64(2), summary.TotalReviews)
	assert.Equal(t, int64(1), summary.FailedCount)
	assert.Equal(t, int64(0), summary.VagueCount)
	assert.Equal(t, int64(1), summary.MasteredCount)
	require.NotNil(t, summary.LastReviewedAt)
	assert.Equal(t, mastered.ReviewedAt, summary.LastReviewedAt.UTC().Truncate(time.Microsecond))
	require.NotNil(t, summary.LastResult)
	assert.Equal(t, domainreview.ResultMastered, *summary.LastResult)
}

func TestIntegrationTransactionRunnerRollsBackRepositoryWrites(t *testing.T) {
	db := setupTestDB(t)
	migrateInterviewQuestionTables(t, db)
	repo := infrapostgres.NewInterviewQuestionRepo(db)
	runner := infrapostgres.NewTransactionRunner(db)
	ctx := context.Background()
	userID := uuid.NewString()
	question := testInterviewQuestion(uuid.NewString(), userID, "Rollback", []string{"tx"})

	err := runner.WithinTransaction(ctx, func(txCtx context.Context) error {
		require.NoError(t, repo.Create(txCtx, question))
		_, findErr := repo.FindByIDAndUserID(txCtx, question.ID, userID)
		require.NoError(t, findErr)
		return errors.New("rollback")
	})

	require.Error(t, err)
	_, findErr := repo.FindByIDAndUserID(ctx, question.ID, userID)
	assert.ErrorIs(t, findErr, domainerrors.ErrNotFound)
}
