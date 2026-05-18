package interview_question

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	reviewrepo "github.com/yilin/ai-for-backend/internal/domain/flashcard_review/repository"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	occurrencerepo "github.com/yilin/ai-for-backend/internal/domain/question_occurrence/repository"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

const (
	defaultPage        = 1
	defaultPageSize    = 10
	maxPageSize        = 100
	defaultFlashLimit  = 20
	maxFlashLimit      = 20
	recentReviewsLimit = 10
)

type ApplicationReader interface {
	Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error)
}

type ProgressReader interface {
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error)
}

type TransactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}

type InterviewQuestionService struct {
	questions      questionrepo.Repository
	occurrences    occurrencerepo.Repository
	reviews        reviewrepo.Repository
	apps           ApplicationReader
	progresses     ProgressReader
	transactions   TransactionRunner
	activityLogger activitylog.Logger
}

func NewInterviewQuestionService(questions questionrepo.Repository, occurrences occurrencerepo.Repository, reviews reviewrepo.Repository, apps ApplicationReader, progresses ProgressReader) *InterviewQuestionService {
	return &InterviewQuestionService{
		questions:      questions,
		occurrences:    occurrences,
		reviews:        reviews,
		apps:           apps,
		progresses:     progresses,
		activityLogger: activitylog.NoopLogger{},
	}
}

func NewInterviewQuestionServiceWithTransactions(questions questionrepo.Repository, occurrences occurrencerepo.Repository, reviews reviewrepo.Repository, apps ApplicationReader, progresses ProgressReader, transactions TransactionRunner) *InterviewQuestionService {
	service := NewInterviewQuestionService(questions, occurrences, reviews, apps, progresses)
	service.transactions = transactions
	return service
}

func (s *InterviewQuestionService) WithTransactions(transactions TransactionRunner) *InterviewQuestionService {
	s.transactions = transactions
	return s
}

func (s *InterviewQuestionService) WithActivityLogger(logger activitylog.Logger) *InterviewQuestionService {
	s.activityLogger = activitylog.OrNoop(logger)
	return s
}

func (s *InterviewQuestionService) withinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if s.transactions == nil {
		return domainerrors.ErrBadParams
	}
	return s.transactions.WithinTransaction(ctx, fn)
}

func (s *InterviewQuestionService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	question, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             uuid.NewString(),
		UserID:         userID,
		Title:          req.Title,
		Content:        req.Content,
		AnswerAnalysis: req.AnswerAnalysis,
		Tags:           req.Tags,
		Difficulty:     req.Difficulty,
		QuestionType:   req.QuestionType,
		MasteryStatus:  req.MasteryStatus,
		IsFavorite:     req.IsFavorite,
		Source:         req.Source,
		Remark:         req.Remark,
	})
	if err != nil {
		return nil, err
	}
	if s.questions == nil {
		return nil, domainerrors.ErrBadParams
	}
	if err := s.questions.Create(ctx, question); err != nil {
		return nil, err
	}
	s.logQuestionActivity(ctx, question, domainactivity.ActionCreated)
	return question, nil
}

func (s *InterviewQuestionService) ListQuestions(ctx context.Context, userID string, req dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	filter, err := buildQuestionListFilter(userID, req)
	if err != nil {
		return nil, 0, err
	}
	if s.questions == nil {
		return nil, 0, domainerrors.ErrBadParams
	}
	return s.questions.ListByUser(ctx, filter)
}

func (s *InterviewQuestionService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
	if s.questions == nil || s.occurrences == nil || s.reviews == nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	count, err := s.questions.CountOccurrences(ctx, id, userID)
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	occurrences, err := s.occurrences.ListByUser(ctx, occurrencerepo.ListFilter{UserID: userID, QuestionID: id})
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	summary, err := s.reviews.SummaryByQuestion(ctx, id, userID)
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	recentReviews, err := s.reviews.ListRecentByQuestion(ctx, id, userID, recentReviewsLimit)
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	if summary == nil {
		summary = &domainreview.ReviewSummary{}
	}
	return question, count, occurrences, *summary, recentReviews, nil
}

func (s *InterviewQuestionService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	if s.questions == nil {
		return nil, domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := question.ApplyUpdate(domainquestion.UpdateSpec{
		Title:          req.Title,
		Content:        req.Content,
		AnswerAnalysis: req.AnswerAnalysis,
		Tags:           req.Tags,
		Difficulty:     req.Difficulty,
		QuestionType:   req.QuestionType,
		MasteryStatus:  req.MasteryStatus,
		IsFavorite:     req.IsFavorite,
		Source:         req.Source,
		Remark:         req.Remark,
	}); err != nil {
		return nil, err
	}
	if err := s.questions.Update(ctx, question); err != nil {
		return nil, err
	}
	s.logQuestionActivity(ctx, question, domainactivity.ActionUpdated)
	return question, nil
}

func (s *InterviewQuestionService) DeleteQuestion(ctx context.Context, userID, id string) error {
	if s.questions == nil || s.occurrences == nil || s.reviews == nil {
		return domainerrors.ErrBadParams
	}
	if s.transactions == nil {
		return domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := s.withinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.occurrences.DeleteByQuestionID(txCtx, id, userID); err != nil {
			return err
		}
		if err := s.reviews.DeleteByQuestionID(txCtx, id, userID); err != nil {
			return err
		}
		return s.questions.Delete(txCtx, id, userID)
	}); err != nil {
		return err
	}
	s.logQuestionActivity(ctx, question, domainactivity.ActionDeleted)
	return nil
}

func (s *InterviewQuestionService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	if req.IsFavorite == nil {
		return nil, domainerrors.ErrBadParams
	}
	if s.questions == nil {
		return nil, domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := question.ApplyUpdate(domainquestion.UpdateSpec{IsFavorite: req.IsFavorite}); err != nil {
		return nil, err
	}
	if err := s.questions.Update(ctx, question); err != nil {
		return nil, err
	}
	action := domainactivity.ActionUnfavorited
	if question.IsFavorite {
		action = domainactivity.ActionFavorited
	}
	s.logQuestionActivity(ctx, question, action)
	return question, nil
}

func (s *InterviewQuestionService) ListOccurrences(ctx context.Context, userID string, req dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	if s.occurrences == nil {
		return nil, domainerrors.ErrBadParams
	}
	return s.occurrences.ListByUser(ctx, occurrencerepo.ListFilter{
		UserID:        userID,
		QuestionID:    req.QuestionID,
		ApplicationID: req.ApplicationID,
	})
}

func (s *InterviewQuestionService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	if req.OccurredAt == nil {
		return nil, domainerrors.ErrBadParams
	}
	if s.questions == nil || s.occurrences == nil || s.apps == nil {
		return nil, domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, req.QuestionID, userID)
	if err != nil {
		return nil, err
	}
	application, err := s.apps.Get(ctx, userID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	stage := strings.TrimSpace(req.InterviewStageSnapshot)
	if req.InterviewProgressID != nil {
		if s.progresses == nil {
			return nil, domainerrors.ErrBadParams
		}
		progress, err := s.progresses.FindByIDAndUserID(ctx, *req.InterviewProgressID, userID)
		if err != nil {
			return nil, err
		}
		if progress.ApplicationID != req.ApplicationID {
			return nil, domainerrors.ErrBadParams
		}
		if stage == "" {
			stage = string(progress.Stage)
		}
	}

	occurrence, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                     uuid.NewString(),
		UserID:                 userID,
		QuestionID:             req.QuestionID,
		ApplicationID:          req.ApplicationID,
		InterviewProgressID:    req.InterviewProgressID,
		CompanyNameSnapshot:    application.CompanyName,
		JobTitleSnapshot:       application.JobTitle,
		InterviewStageSnapshot: stage,
		OccurredAt:             req.OccurredAt.Value(),
		ActualQuestion:         req.ActualQuestion,
		AnswerPerformance:      req.AnswerPerformance,
		Note:                   req.Note,
	})
	if err != nil {
		return nil, err
	}
	if err := s.occurrences.Create(ctx, occurrence); err != nil {
		return nil, err
	}
	s.logOccurrenceActivity(ctx, occurrence, question, domainactivity.ActionLinked)
	return occurrence, nil
}

func (s *InterviewQuestionService) DeleteOccurrence(ctx context.Context, userID, id string) error {
	if s.occurrences == nil || s.questions == nil {
		return domainerrors.ErrBadParams
	}
	occurrence, err := s.occurrences.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	question, err := s.questions.FindByIDAndUserID(ctx, occurrence.QuestionID, userID)
	if err != nil {
		return err
	}
	if err := s.occurrences.Delete(ctx, id, userID); err != nil {
		return err
	}
	s.logOccurrenceActivity(ctx, occurrence, question, domainactivity.ActionDeleted)
	return nil
}

func (s *InterviewQuestionService) ListFlashcards(ctx context.Context, userID string, req dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	filter, err := buildFlashcardFilter(userID, req)
	if err != nil {
		return nil, 0, err
	}
	if s.questions == nil {
		return nil, 0, domainerrors.ErrBadParams
	}
	return s.questions.ListFlashcards(ctx, filter)
}

func (s *InterviewQuestionService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	if s.questions == nil || s.reviews == nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, domainerrors.ErrBadParams
	}
	question, err := s.questions.FindByIDAndUserID(ctx, req.QuestionID, userID)
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	reviewedAt := time.Now()
	if req.ReviewedAt != nil {
		reviewedAt = req.ReviewedAt.Value()
	}
	review, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{
		ID:         uuid.NewString(),
		UserID:     userID,
		QuestionID: req.QuestionID,
		ReviewedAt: reviewedAt,
		Result:     req.Result,
		ReviewMode: req.ReviewMode,
	})
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	if err := s.withinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.reviews.Create(txCtx, review); err != nil {
			return err
		}
		if err := question.ApplyReviewResult(domainquestion.ReviewResult(req.Result), reviewedAt); err != nil {
			return err
		}
		return s.questions.Update(txCtx, question)
	}); err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	summary, err := s.reviews.SummaryByQuestion(ctx, req.QuestionID, userID)
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	count, err := s.questions.CountOccurrences(ctx, req.QuestionID, userID)
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	if summary == nil {
		summary = &domainreview.ReviewSummary{}
	}
	return review, questionrepo.QuestionWithOccurrenceCount{Question: *question, OccurrenceCount: count}, *summary, nil
}

func (s *InterviewQuestionService) logQuestionActivity(ctx context.Context, question *domainquestion.InterviewQuestion, action domainactivity.Action) {
	activitylog.BestEffortLog(ctx, s.activityLogger, questionActivityEntry(question, action))
}

func questionActivityEntry(question *domainquestion.InterviewQuestion, action domainactivity.Action) activitylog.Entry {
	label := activitylog.TrimLabel(question.Title)
	route := "/interview-question-bank/question/" + question.ID
	if action == domainactivity.ActionDeleted {
		route = ""
	}
	return activitylog.Entry{
		UserID:      question.UserID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeQuestion,
		TargetID:    question.ID,
		TargetLabel: label,
		Summary:     activitylog.TrimSummary(questionSummary(label, action)),
		Route:       route,
	}
}

func (s *InterviewQuestionService) logOccurrenceActivity(ctx context.Context, occurrence *domainoccurrence.QuestionOccurrence, question *domainquestion.InterviewQuestion, action domainactivity.Action) {
	companyJobLabel := activitylog.TrimLabel(occurrence.CompanyNameSnapshot + " / " + occurrence.JobTitleSnapshot)
	questionLabel := activitylog.TrimLabel(question.Title)
	activitylog.BestEffortLog(ctx, s.activityLogger, activitylog.Entry{
		UserID:      occurrence.UserID,
		Action:      action,
		TargetType:  domainactivity.TargetTypeQuestion,
		TargetID:    question.ID,
		TargetLabel: questionLabel,
		Summary:     activitylog.TrimSummary(occurrenceSummary(questionLabel, companyJobLabel, action)),
		Route:       "/interview-question-bank/question/" + question.ID,
	})
}

func questionSummary(label string, action domainactivity.Action) string {
	switch action {
	case domainactivity.ActionCreated:
		return label + " 面试题已创建"
	case domainactivity.ActionDeleted:
		return label + " 面试题已删除"
	case domainactivity.ActionFavorited:
		return label + " 面试题已收藏"
	case domainactivity.ActionUnfavorited:
		return label + " 面试题已取消收藏"
	default:
		return label + " 面试题已更新"
	}
}

func occurrenceSummary(questionLabel, companyJobLabel string, action domainactivity.Action) string {
	if action == domainactivity.ActionDeleted {
		return "删除题目关联记录 " + questionLabel + " - " + companyJobLabel
	}
	return "题目 " + questionLabel + " 关联到 " + companyJobLabel
}

func buildQuestionListFilter(userID string, req dto.InterviewQuestionQuery) (questionrepo.ListFilter, error) {
	if req.Difficulty != "" && !req.Difficulty.Valid() {
		return questionrepo.ListFilter{}, domainerrors.ErrBadParams
	}
	if req.QuestionType != "" && !req.QuestionType.Valid() {
		return questionrepo.ListFilter{}, domainerrors.ErrBadParams
	}
	if req.MasteryStatus != "" && !req.MasteryStatus.Valid() {
		return questionrepo.ListFilter{}, domainerrors.ErrBadParams
	}
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	return questionrepo.ListFilter{
		UserID:        userID,
		Keyword:       req.Keyword,
		Tags:          splitTags(req.Tags),
		Difficulty:    req.Difficulty,
		QuestionType:  req.QuestionType,
		MasteryStatus: req.MasteryStatus,
		IsFavorite:    req.IsFavorite,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func buildFlashcardFilter(userID string, req dto.FlashcardQuery) (questionrepo.FlashcardFilter, error) {
	mode := req.Mode
	if mode == "" {
		mode = domainreview.ReviewModeRandom
	}
	if !mode.Valid() {
		return questionrepo.FlashcardFilter{}, domainerrors.ErrBadParams
	}
	if req.Difficulty != "" && !req.Difficulty.Valid() {
		return questionrepo.FlashcardFilter{}, domainerrors.ErrBadParams
	}
	limit := req.Limit
	if limit <= 0 || limit > maxFlashLimit {
		limit = defaultFlashLimit
	}
	return questionrepo.FlashcardFilter{
		UserID:     userID,
		Mode:       mode,
		Tags:       splitTags(req.Tags),
		Difficulty: req.Difficulty,
		Limit:      limit,
	}, nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}
	return page, pageSize
}

func splitTags(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
