package interview_question

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	occurrencerepo "github.com/yilin/ai-for-backend/internal/domain/question_occurrence/repository"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type spyActivityLogger struct {
	Entries []activitylog.Entry
}

func (l *spyActivityLogger) Log(ctx context.Context, entry activitylog.Entry) error {
	l.Entries = append(l.Entries, entry)
	return nil
}

type mockQuestionRepo struct{ mock.Mock }

func (m *mockQuestionRepo) Create(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockQuestionRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, id, userID)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Error(1)
}

func (m *mockQuestionRepo) Update(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockQuestionRepo) Delete(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockQuestionRepo) ListByUser(ctx context.Context, f questionrepo.ListFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, f)
	items, _ := args.Get(0).([]questionrepo.QuestionWithOccurrenceCount)
	return items, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) CountOccurrences(ctx context.Context, questionID, userID string) (int64, error) {
	args := m.Called(ctx, questionID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockQuestionRepo) ListFlashcards(ctx context.Context, f questionrepo.FlashcardFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, f)
	items, _ := args.Get(0).([]questionrepo.QuestionWithOccurrenceCount)
	return items, args.Get(1).(int64), args.Error(2)
}

type mockOccurrenceRepo struct{ mock.Mock }

func (m *mockOccurrenceRepo) Create(ctx context.Context, o *domainoccurrence.QuestionOccurrence) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOccurrenceRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, id, userID)
	o, _ := args.Get(0).(*domainoccurrence.QuestionOccurrence)
	return o, args.Error(1)
}

func (m *mockOccurrenceRepo) Delete(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockOccurrenceRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	args := m.Called(ctx, questionID, userID)
	return args.Error(0)
}

func (m *mockOccurrenceRepo) ListByUser(ctx context.Context, f occurrencerepo.ListFilter) ([]domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, f)
	items, _ := args.Get(0).([]domainoccurrence.QuestionOccurrence)
	return items, args.Error(1)
}

type mockReviewRepo struct{ mock.Mock }

func (m *mockReviewRepo) Create(ctx context.Context, r *domainreview.FlashcardReview) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *mockReviewRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	args := m.Called(ctx, questionID, userID)
	return args.Error(0)
}

func (m *mockReviewRepo) ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error) {
	args := m.Called(ctx, questionID, userID, limit)
	items, _ := args.Get(0).([]domainreview.FlashcardReview)
	return items, args.Error(1)
}

func (m *mockReviewRepo) SummaryByQuestion(ctx context.Context, questionID, userID string) (*domainreview.ReviewSummary, error) {
	args := m.Called(ctx, questionID, userID)
	summary, _ := args.Get(0).(*domainreview.ReviewSummary)
	return summary, args.Error(1)
}

type mockApplicationReader struct{ mock.Mock }

func (m *mockApplicationReader) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	args := m.Called(ctx, userID, id)
	app, _ := args.Get(0).(*domainapp.ResumeApplication)
	return app, args.Error(1)
}

type mockProgressReader struct{ mock.Mock }

func (m *mockProgressReader) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, id, userID)
	progress, _ := args.Get(0).(*domainprogress.InterviewProgress)
	return progress, args.Error(1)
}

type transactionContextKey struct{}

type fakeTransactionRunner struct {
	calls int
}

func (r *fakeTransactionRunner) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	r.calls++
	return fn(context.WithValue(ctx, transactionContextKey{}, true))
}

func transactionContextMatcher() interface{} {
	return mock.MatchedBy(func(ctx context.Context) bool {
		inTransaction, _ := ctx.Value(transactionContextKey{}).(bool)
		return inTransaction
	})
}

func TestInterviewQuestionService_GetQuestionDoesNotLogActivity(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	reviews := new(mockReviewRepo)
	logger := &spyActivityLogger{}
	service := NewInterviewQuestionService(questions, occurrences, reviews, nil, nil).WithActivityLogger(logger)
	question := mustQuestion(t, "question-1", "user-1")

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	questions.On("CountOccurrences", ctx, "question-1", "user-1").Return(int64(0), nil).Once()
	occurrences.On("ListByUser", ctx, occurrencerepo.ListFilter{UserID: "user-1", QuestionID: "question-1"}).Return([]domainoccurrence.QuestionOccurrence{}, nil).Once()
	reviews.On("SummaryByQuestion", ctx, "question-1", "user-1").Return(&domainreview.ReviewSummary{}, nil).Once()
	reviews.On("ListRecentByQuestion", ctx, "question-1", "user-1", 10).Return([]domainreview.FlashcardReview{}, nil).Once()

	got, count, gotOccurrences, _, gotReviews, err := service.GetQuestion(ctx, "user-1", "question-1")

	require.NoError(t, err)
	require.Equal(t, question, got)
	require.Zero(t, count)
	require.Empty(t, gotOccurrences)
	require.Empty(t, gotReviews)
	require.Empty(t, logger.Entries)
	questions.AssertExpectations(t)
	occurrences.AssertExpectations(t)
	reviews.AssertExpectations(t)
}

func TestInterviewQuestionService_CreateUpdateFavoriteAndOccurrenceLogExpectedActions(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	apps := new(mockApplicationReader)
	logger := &spyActivityLogger{}
	service := NewInterviewQuestionService(questions, occurrences, new(mockReviewRepo), apps, new(mockProgressReader)).WithActivityLogger(logger)
	question := mustQuestion(t, "question-1", "user-1")
	question.IsFavorite = false
	newTitle := "Updated question title"
	favorite := true
	occurredAt := dto.RequiredFlexibleTime{Time: time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)}

	questions.On("Create", ctx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.UserID == "user-1" && q.Title == "Created question"
	})).Return(nil).Once()
	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Twice()
	questions.On("Update", ctx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.ID == "question-1" && q.Title == newTitle
	})).Return(nil).Once()
	questions.On("Update", ctx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.ID == "question-1" && q.IsFavorite
	})).Return(nil).Once()
	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	apps.On("Get", ctx, "user-1", "app-1").Return(mustApplication(t, "app-1", "user-1", "Acme", "Backend Engineer"), nil).Once()
	occurrences.On("Create", ctx, mock.MatchedBy(func(o *domainoccurrence.QuestionOccurrence) bool {
		return o.QuestionID == "question-1" && o.ApplicationID == "app-1"
	})).Return(nil).Once()

	created, err := service.CreateQuestion(ctx, "user-1", dto.CreateInterviewQuestionRequest{
		Title:          "Created question",
		Content:        "Question content",
		AnswerAnalysis: "Answer analysis",
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusUnreviewed,
	})
	require.NoError(t, err)
	_, err = service.UpdateQuestion(ctx, "user-1", "question-1", dto.UpdateInterviewQuestionRequest{Title: &newTitle})
	require.NoError(t, err)
	_, err = service.UpdateFavorite(ctx, "user-1", "question-1", dto.UpdateQuestionFavoriteRequest{IsFavorite: &favorite})
	require.NoError(t, err)
	occurrence, err := service.CreateOccurrence(ctx, "user-1", dto.CreateQuestionOccurrenceRequest{QuestionID: "question-1", ApplicationID: "app-1", OccurredAt: &occurredAt})
	require.NoError(t, err)

	require.Len(t, logger.Entries, 4)
	assertActivity := func(index int, action domainactivity.Action, targetID, route string) {
		t.Helper()
		assert.Equal(t, action, logger.Entries[index].Action)
		assert.Equal(t, domainactivity.TargetTypeQuestion, logger.Entries[index].TargetType)
		assert.Equal(t, targetID, logger.Entries[index].TargetID)
		assert.Equal(t, route, logger.Entries[index].Route)
	}
	assertActivity(0, domainactivity.ActionCreated, created.ID, "/interview-question-bank/question/"+created.ID)
	assertActivity(1, domainactivity.ActionUpdated, "question-1", "/interview-question-bank/question/question-1")
	assertActivity(2, domainactivity.ActionFavorited, "question-1", "/interview-question-bank/question/question-1")
	assertActivity(3, domainactivity.ActionLinked, occurrence.ID, "/resume/application/app-1")
	questions.AssertExpectations(t)
	occurrences.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestCreateQuestionDoesNotRequireApplicationAndDoesNotCreateOccurrence(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	apps := new(mockApplicationReader)
	service := NewInterviewQuestionService(questions, occurrences, new(mockReviewRepo), apps, new(mockProgressReader))

	questions.On("Create", ctx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.UserID == "user-1" && q.Title == "Explain indexes" && q.MasteryStatus == domainquestion.MasteryStatusUnreviewed
	})).Return(nil).Once()

	question, err := service.CreateQuestion(ctx, "user-1", dto.CreateInterviewQuestionRequest{
		Title:          "Explain indexes",
		Content:        "What are DB indexes?",
		AnswerAnalysis: "Indexes speed up reads.",
		Tags:           []string{"db"},
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusUnreviewed,
	})

	require.NoError(t, err)
	require.NotEmpty(t, question.ID)
	questions.AssertExpectations(t)
	occurrences.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	apps.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestListQuestionsDefaultsPaginationAndSplitsTags(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	service := NewInterviewQuestionService(questions, new(mockOccurrenceRepo), new(mockReviewRepo), nil, nil)

	questions.On("ListByUser", ctx, questionrepo.ListFilter{
		UserID:   "user-1",
		Keyword:  "redis",
		Tags:     []string{"backend", "cache"},
		Page:     1,
		PageSize: 10,
	}).Return([]questionrepo.QuestionWithOccurrenceCount{}, int64(0), nil).Once()

	items, total, err := service.ListQuestions(ctx, "user-1", dto.InterviewQuestionQuery{
		Keyword:  "redis",
		Tags:     " backend, cache ,, ",
		Page:     0,
		PageSize: 150,
	})

	require.NoError(t, err)
	require.Empty(t, items)
	require.Equal(t, int64(0), total)
	questions.AssertExpectations(t)
}

func TestListQuestionsRejectsInvalidEnumsWithoutRepoCall(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name  string
		query dto.InterviewQuestionQuery
	}{
		{name: "difficulty", query: dto.InterviewQuestionQuery{Difficulty: domainquestion.Difficulty("impossible")}},
		{name: "type", query: dto.InterviewQuestionQuery{QuestionType: domainquestion.QuestionType("essay")}},
		{name: "mastery", query: dto.InterviewQuestionQuery{MasteryStatus: domainquestion.MasteryStatus("perfect")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			questions := new(mockQuestionRepo)
			service := NewInterviewQuestionService(questions, nil, nil, nil, nil)

			items, total, err := service.ListQuestions(ctx, "user-1", tt.query)

			require.ErrorIs(t, err, domainerrors.ErrBadParams)
			require.Nil(t, items)
			require.Zero(t, total)
			questions.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything)
		})
	}
}

func TestUpdateFavoriteRejectsNilIsFavoriteAndSupportsExplicitFalse(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	service := NewInterviewQuestionService(questions, nil, nil, nil, nil)

	question := mustQuestion(t, "question-1", "user-1")
	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	questions.On("Update", ctx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.ID == "question-1" && !q.IsFavorite
	})).Return(nil).Once()

	got, err := service.UpdateFavorite(ctx, "user-1", "question-1", dto.UpdateQuestionFavoriteRequest{})
	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Nil(t, got)

	falseValue := false
	got, err = service.UpdateFavorite(ctx, "user-1", "question-1", dto.UpdateQuestionFavoriteRequest{IsFavorite: &falseValue})
	require.NoError(t, err)
	require.False(t, got.IsFavorite)
	questions.AssertExpectations(t)
}

func TestCreateOccurrenceCopiesAppSnapshotsAndFillsProgressStageWhenMissing(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	apps := new(mockApplicationReader)
	progresses := new(mockProgressReader)
	service := NewInterviewQuestionService(questions, occurrences, new(mockReviewRepo), apps, progresses)
	occurredAt := dto.RequiredFlexibleTime{Time: time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)}
	progressID := "progress-1"
	app := mustApplication(t, "app-1", "user-1", "Acme", "Backend Engineer")
	progress := mustProgress(t, "progress-1", "user-1", "app-1", domainprogress.StageSecondInterview)

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(mustQuestion(t, "question-1", "user-1"), nil).Once()
	apps.On("Get", ctx, "user-1", "app-1").Return(app, nil).Once()
	progresses.On("FindByIDAndUserID", ctx, "progress-1", "user-1").Return(progress, nil).Once()
	occurrences.On("Create", ctx, mock.MatchedBy(func(o *domainoccurrence.QuestionOccurrence) bool {
		return o.UserID == "user-1" &&
			o.QuestionID == "question-1" &&
			o.ApplicationID == "app-1" &&
			o.InterviewProgressID != nil && *o.InterviewProgressID == "progress-1" &&
			o.CompanyNameSnapshot == "Acme" &&
			o.JobTitleSnapshot == "Backend Engineer" &&
			o.InterviewStageSnapshot == string(domainprogress.StageSecondInterview) &&
			o.OccurredAt.Equal(occurredAt.Value())
	})).Return(nil).Once()

	occurrence, err := service.CreateOccurrence(ctx, "user-1", dto.CreateQuestionOccurrenceRequest{
		QuestionID:          "question-1",
		ApplicationID:       "app-1",
		InterviewProgressID: &progressID,
		OccurredAt:          &occurredAt,
	})

	require.NoError(t, err)
	require.Equal(t, "Acme", occurrence.CompanyNameSnapshot)
	require.Equal(t, string(domainprogress.StageSecondInterview), occurrence.InterviewStageSnapshot)
	questions.AssertExpectations(t)
	apps.AssertExpectations(t)
	progresses.AssertExpectations(t)
	occurrences.AssertExpectations(t)
}

func TestCreateOccurrenceRejectsNilOccurredAt(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	service := NewInterviewQuestionService(questions, new(mockOccurrenceRepo), nil, nil, nil)

	occurrence, err := service.CreateOccurrence(ctx, "user-1", dto.CreateQuestionOccurrenceRequest{QuestionID: "question-1", ApplicationID: "app-1"})

	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Nil(t, occurrence)
	questions.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateOccurrenceRejectsProgressFromDifferentApplication(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	apps := new(mockApplicationReader)
	progresses := new(mockProgressReader)
	occurrences := new(mockOccurrenceRepo)
	service := NewInterviewQuestionService(questions, occurrences, nil, apps, progresses)
	occurredAt := dto.RequiredFlexibleTime{Time: time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)}
	progressID := "progress-1"

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(mustQuestion(t, "question-1", "user-1"), nil).Once()
	apps.On("Get", ctx, "user-1", "app-1").Return(mustApplication(t, "app-1", "user-1", "Acme", "Backend Engineer"), nil).Once()
	progresses.On("FindByIDAndUserID", ctx, "progress-1", "user-1").Return(mustProgress(t, "progress-1", "user-1", "other-app", domainprogress.StageFirstInterview), nil).Once()

	occurrence, err := service.CreateOccurrence(ctx, "user-1", dto.CreateQuestionOccurrenceRequest{
		QuestionID:          "question-1",
		ApplicationID:       "app-1",
		InterviewProgressID: &progressID,
		OccurredAt:          &occurredAt,
	})

	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Nil(t, occurrence)
	occurrences.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateFlashcardReviewRequiresTransactionRunnerWithoutWriting(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	reviews := new(mockReviewRepo)
	service := NewInterviewQuestionService(questions, new(mockOccurrenceRepo), reviews, nil, nil)
	question := mustQuestion(t, "question-1", "user-1")

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Maybe()

	review, gotQuestion, gotSummary, err := service.CreateFlashcardReview(ctx, "user-1", dto.CreateFlashcardReviewRequest{
		QuestionID: "question-1",
		Result:     domainreview.ResultMastered,
		ReviewMode: domainreview.ReviewModeRandom,
	})

	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	require.Nil(t, review)
	require.Empty(t, gotQuestion)
	require.Empty(t, gotSummary)
	questions.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	reviews.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateFlashcardReviewCreatesHistoryAndUpdatesQuestionMasteryAndLastReviewedAt(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	reviews := new(mockReviewRepo)
	transactions := new(fakeTransactionRunner)
	service := NewInterviewQuestionServiceWithTransactions(questions, new(mockOccurrenceRepo), reviews, nil, nil, transactions)
	reviewedAt := dto.FlexibleTime{Time: time.Date(2026, 5, 15, 11, 0, 0, 0, time.UTC)}
	question := mustQuestion(t, "question-1", "user-1")
	summary := &domainreview.ReviewSummary{TotalReviews: 1, MasteredCount: 1, LastReviewedAt: &reviewedAt.Time}
	updatedWithCount := questionrepo.QuestionWithOccurrenceCount{Question: *question, OccurrenceCount: 3}
	txCtx := transactionContextMatcher()

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	reviews.On("Create", txCtx, mock.MatchedBy(func(r *domainreview.FlashcardReview) bool {
		return r.UserID == "user-1" && r.QuestionID == "question-1" && r.Result == domainreview.ResultMastered && r.ReviewMode == domainreview.ReviewModeRandom && r.ReviewedAt.Equal(reviewedAt.Value())
	})).Return(nil).Once()
	questions.On("Update", txCtx, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.MasteryStatus == domainquestion.MasteryStatusMastered && q.LastReviewedAt != nil && q.LastReviewedAt.Equal(reviewedAt.Value())
	})).Return(nil).Once()
	reviews.On("SummaryByQuestion", ctx, "question-1", "user-1").Return(summary, nil).Once()
	questions.On("CountOccurrences", ctx, "question-1", "user-1").Return(int64(3), nil).Once()

	review, gotQuestion, gotSummary, err := service.CreateFlashcardReview(ctx, "user-1", dto.CreateFlashcardReviewRequest{
		QuestionID: "question-1",
		ReviewedAt: &reviewedAt,
		Result:     domainreview.ResultMastered,
		ReviewMode: domainreview.ReviewModeRandom,
	})

	require.NoError(t, err)
	require.Equal(t, 1, transactions.calls)
	require.Equal(t, domainreview.ResultMastered, review.Result)
	require.Equal(t, domainquestion.MasteryStatusMastered, gotQuestion.Question.MasteryStatus)
	require.Equal(t, updatedWithCount.OccurrenceCount, gotQuestion.OccurrenceCount)
	require.Equal(t, summary, &gotSummary)
	questions.AssertExpectations(t)
	reviews.AssertExpectations(t)
}

func TestCreateFlashcardReviewReturnsUpdateErrorWithoutSummaryOrCountQueries(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	reviews := new(mockReviewRepo)
	transactions := new(fakeTransactionRunner)
	service := NewInterviewQuestionServiceWithTransactions(questions, new(mockOccurrenceRepo), reviews, nil, nil, transactions)
	reviewedAt := dto.FlexibleTime{Time: time.Date(2026, 5, 15, 11, 0, 0, 0, time.UTC)}
	question := mustQuestion(t, "question-1", "user-1")
	expectedErr := errors.New("question update failed")
	txCtx := transactionContextMatcher()

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	reviews.On("Create", txCtx, mock.AnythingOfType("*flashcard_review.FlashcardReview")).Return(nil).Once()
	questions.On("Update", txCtx, mock.AnythingOfType("*interview_question.InterviewQuestion")).Return(expectedErr).Once()

	review, gotQuestion, gotSummary, err := service.CreateFlashcardReview(ctx, "user-1", dto.CreateFlashcardReviewRequest{
		QuestionID: "question-1",
		ReviewedAt: &reviewedAt,
		Result:     domainreview.ResultMastered,
		ReviewMode: domainreview.ReviewModeRandom,
	})

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, review)
	require.Empty(t, gotQuestion)
	require.Empty(t, gotSummary)
	require.Equal(t, 1, transactions.calls)
	reviews.AssertNotCalled(t, "SummaryByQuestion", mock.Anything, mock.Anything, mock.Anything)
	questions.AssertNotCalled(t, "CountOccurrences", mock.Anything, mock.Anything, mock.Anything)
	questions.AssertExpectations(t)
	reviews.AssertExpectations(t)
}

func TestGetQuestionAggregatesCountOccurrencesSummaryAndRecentReviews(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	reviews := new(mockReviewRepo)
	service := NewInterviewQuestionService(questions, occurrences, reviews, nil, nil)
	question := mustQuestion(t, "question-1", "user-1")
	occurrence := mustOccurrence(t, "occurrence-1", "user-1", "question-1", "app-1")
	review := mustReview(t, "review-1", "user-1", "question-1")
	summary := &domainreview.ReviewSummary{TotalReviews: 5}

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(question, nil).Once()
	questions.On("CountOccurrences", ctx, "question-1", "user-1").Return(int64(2), nil).Once()
	occurrences.On("ListByUser", ctx, occurrencerepo.ListFilter{UserID: "user-1", QuestionID: "question-1"}).Return([]domainoccurrence.QuestionOccurrence{*occurrence}, nil).Once()
	reviews.On("SummaryByQuestion", ctx, "question-1", "user-1").Return(summary, nil).Once()
	reviews.On("ListRecentByQuestion", ctx, "question-1", "user-1", 10).Return([]domainreview.FlashcardReview{*review}, nil).Once()

	gotQuestion, count, gotOccurrences, gotSummary, gotReviews, err := service.GetQuestion(ctx, "user-1", "question-1")

	require.NoError(t, err)
	require.Equal(t, question, gotQuestion)
	require.Equal(t, int64(2), count)
	require.Len(t, gotOccurrences, 1)
	require.Equal(t, *summary, gotSummary)
	require.Len(t, gotReviews, 1)
	questions.AssertExpectations(t)
	occurrences.AssertExpectations(t)
	reviews.AssertExpectations(t)
}

func TestDeleteQuestionRequiresTransactionRunnerWithoutWriting(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	reviews := new(mockReviewRepo)
	service := NewInterviewQuestionService(questions, occurrences, reviews, nil, nil)

	err := service.DeleteQuestion(ctx, "user-1", "question-1")

	require.ErrorIs(t, err, domainerrors.ErrBadParams)
	occurrences.AssertNotCalled(t, "DeleteByQuestionID", mock.Anything, mock.Anything, mock.Anything)
	reviews.AssertNotCalled(t, "DeleteByQuestionID", mock.Anything, mock.Anything, mock.Anything)
	questions.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteQuestionDeletesChildrenBeforeQuestionInsideTransaction(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	reviews := new(mockReviewRepo)
	transactions := new(fakeTransactionRunner)
	service := NewInterviewQuestionServiceWithTransactions(questions, occurrences, reviews, nil, nil, transactions)
	txCtx := transactionContextMatcher()

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(mustQuestion(t, "question-1", "user-1"), nil).Once()
	occurrences.On("DeleteByQuestionID", txCtx, "question-1", "user-1").Return(nil).Once()
	reviews.On("DeleteByQuestionID", txCtx, "question-1", "user-1").Return(nil).Once()
	questions.On("Delete", txCtx, "question-1", "user-1").Return(nil).Once()

	err := service.DeleteQuestion(ctx, "user-1", "question-1")

	require.NoError(t, err)
	require.Equal(t, 1, transactions.calls)
	occurrences.AssertExpectations(t)
	reviews.AssertExpectations(t)
	questions.AssertExpectations(t)
}

func TestDeleteQuestionReturnsCascadeErrorWithoutDeletingQuestion(t *testing.T) {
	ctx := context.Background()
	questions := new(mockQuestionRepo)
	occurrences := new(mockOccurrenceRepo)
	reviews := new(mockReviewRepo)
	transactions := new(fakeTransactionRunner)
	service := NewInterviewQuestionServiceWithTransactions(questions, occurrences, reviews, nil, nil, transactions)
	expectedErr := errors.New("cascade failed")

	questions.On("FindByIDAndUserID", ctx, "question-1", "user-1").Return(mustQuestion(t, "question-1", "user-1"), nil).Once()
	occurrences.On("DeleteByQuestionID", transactionContextMatcher(), "question-1", "user-1").Return(expectedErr).Once()

	err := service.DeleteQuestion(ctx, "user-1", "question-1")

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, transactions.calls)
	reviews.AssertNotCalled(t, "DeleteByQuestionID", mock.Anything, mock.Anything, mock.Anything)
	questions.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	occurrences.AssertExpectations(t)
}

func mustQuestion(t *testing.T, id, userID string) *domainquestion.InterviewQuestion {
	t.Helper()
	q, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             id,
		UserID:         userID,
		Title:          "Question title",
		Content:        "Question content",
		AnswerAnalysis: "Answer analysis",
		Tags:           []string{"tag"},
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.QuestionTypeKnowledge,
		MasteryStatus:  domainquestion.MasteryStatusUnreviewed,
		IsFavorite:     true,
	})
	require.NoError(t, err)
	return q
}

func mustApplication(t *testing.T, id, userID, company, title string) *domainapp.ResumeApplication {
	t.Helper()
	app, err := domainapp.NewResumeApplication(domainapp.NewResumeApplicationParams{
		ID:              id,
		UserID:          userID,
		ResumeVersionID: "version-1",
		RepoID:          "repo-1",
		CompanyName:     company,
		JobTitle:        title,
		AppliedAt:       time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		CurrentStatus:   domainapp.StatusApplied,
	})
	require.NoError(t, err)
	return app
}

func mustProgress(t *testing.T, id, userID, applicationID string, stage domainprogress.Stage) *domainprogress.InterviewProgress {
	t.Helper()
	progress, err := domainprogress.NewInterviewProgress(domainprogress.NewInterviewProgressParams{
		ID:            id,
		UserID:        userID,
		ApplicationID: applicationID,
		Stage:         stage,
		OccurredAt:    time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
		Result:        domainprogress.ResultPending,
	})
	require.NoError(t, err)
	return progress
}

func mustOccurrence(t *testing.T, id, userID, questionID, applicationID string) *domainoccurrence.QuestionOccurrence {
	t.Helper()
	occurrence, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                  id,
		UserID:              userID,
		QuestionID:          questionID,
		ApplicationID:       applicationID,
		CompanyNameSnapshot: "Acme",
		JobTitleSnapshot:    "Backend Engineer",
		OccurredAt:          time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	return occurrence
}

func mustReview(t *testing.T, id, userID, questionID string) *domainreview.FlashcardReview {
	t.Helper()
	review, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{
		ID:         id,
		UserID:     userID,
		QuestionID: questionID,
		ReviewedAt: time.Date(2026, 5, 15, 11, 0, 0, 0, time.UTC),
		Result:     domainreview.ResultVague,
		ReviewMode: domainreview.ReviewModeRandom,
	})
	require.NoError(t, err)
	return review
}
