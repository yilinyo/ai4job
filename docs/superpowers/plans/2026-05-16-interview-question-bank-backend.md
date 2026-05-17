# Interview Question Bank Backend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the backend for the interview question bank, question occurrence tracking, and flashcard review flow, aligned with the approved frontend field contract.

**Architecture:** Add three focused domain models (`interview_question`, `question_occurrence`, `flashcard_review`) and one application service that orchestrates the frontend-facing question-bank workflows. Keep persistence split by table, expose one handler for all question-bank endpoints, and register only `/api/v1/...` routes.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, testify, existing response envelope in `backend/pkg/response`, existing error package in `backend/pkg/errors`.

---

## File map

### New files

- `backend/internal/domain/interview_question/entity.go` — interview question entity, enums, validation, update behavior, review-to-mastery mapping.
- `backend/internal/domain/interview_question/entity_test.go` — domain tests for question validation and update behavior.
- `backend/internal/domain/interview_question/repository/repository.go` — repository interface and list filter for questions.
- `backend/internal/domain/question_occurrence/entity.go` — occurrence entity, enum, validation.
- `backend/internal/domain/question_occurrence/entity_test.go` — occurrence validation tests.
- `backend/internal/domain/question_occurrence/repository/repository.go` — occurrence repository interface and list filter.
- `backend/internal/domain/flashcard_review/entity.go` — flashcard review entity, enums, validation, summary type.
- `backend/internal/domain/flashcard_review/entity_test.go` — review validation and summary tests.
- `backend/internal/domain/flashcard_review/repository/repository.go` — review repository interface and queue filter.
- `backend/internal/application/dto/interview_question_dto.go` — frontend-aligned request/response DTOs and converters.
- `backend/internal/application/interview_question/interview_question_service.go` — application service for question CRUD, occurrences, flashcards, and reviews.
- `backend/internal/application/interview_question/interview_question_service_test.go` — service tests with mocked repositories/readers.
- `backend/internal/infrastructure/persistence/postgres/interview_question_repo.go` — GORM repository for questions and question list aggregation.
- `backend/internal/infrastructure/persistence/postgres/question_occurrence_repo.go` — GORM repository for occurrences.
- `backend/internal/infrastructure/persistence/postgres/flashcard_review_repo.go` — GORM repository for flashcard reviews and summaries.
- `backend/internal/interface/handler/interview_question_handler.go` — Gin handlers for all `/api/v1` question-bank endpoints.
- `backend/internal/interface/handler/interview_question_handler_test.go` — handler tests for response shape and request parsing.
- `backend/migrations/000007_create_interview_questions.up.sql` — migration creating all three tables and indexes.
- `backend/migrations/000007_create_interview_questions.down.sql` — migration dropping all three tables.

### Modified files

- `backend/internal/domain/interview_progress/repository/repository.go` — add `FindByIDAndUserID` if current interface lacks it at execution time; current repo already has it.
- `backend/internal/interface/router/router.go` — add `InterviewQuestionHandler` dependency and routes.
- `backend/internal/interface/router/router_test.go` — add stub service and route coverage for new endpoints.
- `backend/cmd/server/main.go` — wire repositories, service, handler.
- `frontend/src/api/interview-questions.ts` — change URL prefix from `/api` to `/api/v1`.
- `frontend/src/api/question-occurrences.ts` — change URL prefix from `/api` to `/api/v1`.
- `frontend/src/api/flashcards.ts` — change URL prefix from `/api` to `/api/v1`.

---

## Task 1: Add domain entities and domain tests

**Files:**
- Create: `backend/internal/domain/interview_question/entity.go`
- Create: `backend/internal/domain/interview_question/entity_test.go`
- Create: `backend/internal/domain/question_occurrence/entity.go`
- Create: `backend/internal/domain/question_occurrence/entity_test.go`
- Create: `backend/internal/domain/flashcard_review/entity.go`
- Create: `backend/internal/domain/flashcard_review/entity_test.go`

- [ ] **Step 1: Write failing interview question domain tests**

Create `backend/internal/domain/interview_question/entity_test.go`:

```go
package interview_question_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestInterviewQuestion_NewDefaultsAndValidates(t *testing.T) {
	q, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             "q1",
		UserID:         "u1",
		Title:          "Vue 2 响应式原理",
		Content:        "请说明 Vue 2 的响应式实现方式。",
		AnswerAnalysis: "Object.defineProperty 劫持 getter/setter。",
		Tags:           []string{"Vue", "JavaScript"},
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
	})

	require.NoError(t, err)
	assert.Equal(t, "q1", q.ID)
	assert.Equal(t, "u1", q.UserID)
	assert.Equal(t, []string{"Vue", "JavaScript"}, q.Tags)
	assert.False(t, q.IsFavorite)
	assert.Nil(t, q.LastReviewedAt)
}

func TestInterviewQuestion_NewRejectsMissingRequiredFields(t *testing.T) {
	_, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             "q1",
		UserID:         "u1",
		Title:          " ",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Difficulty:     domainquestion.DifficultyEasy,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
	})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestInterviewQuestion_NewRejectsInvalidEnums(t *testing.T) {
	_, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             "q1",
		UserID:         "u1",
		Title:          "title",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Difficulty:     domainquestion.Difficulty("invalid"),
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
	})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestInterviewQuestion_ApplyUpdatePreservesOmittedFields(t *testing.T) {
	q, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             "q1",
		UserID:         "u1",
		Title:          "old title",
		Content:        "old content",
		AnswerAnalysis: "old analysis",
		Tags:           []string{"Vue"},
		Difficulty:     domainquestion.DifficultyEasy,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
		IsFavorite:     false,
		Source:         "old source",
		Remark:         "old remark",
	})
	require.NoError(t, err)

	newTitle := "new title"
	newTags := []string{"Vue", "浏览器"}
	newFavorite := true
	err = q.ApplyUpdate(domainquestion.UpdateSpec{
		Title:      &newTitle,
		Tags:       &newTags,
		IsFavorite: &newFavorite,
	})

	require.NoError(t, err)
	assert.Equal(t, "new title", q.Title)
	assert.Equal(t, "old content", q.Content)
	assert.Equal(t, []string{"Vue", "浏览器"}, q.Tags)
	assert.True(t, q.IsFavorite)
	assert.Equal(t, "old source", q.Source)
}

func TestInterviewQuestion_ApplyReviewResultUpdatesMasteryAndLastReviewedAt(t *testing.T) {
	q, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{
		ID:             "q1",
		UserID:         "u1",
		Title:          "title",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
	})
	require.NoError(t, err)
	reviewedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)

	err = q.ApplyReviewResult(domainquestion.ReviewResultFailed, reviewedAt)

	require.NoError(t, err)
	assert.Equal(t, domainquestion.MasteryWeak, q.MasteryStatus)
	require.NotNil(t, q.LastReviewedAt)
	assert.True(t, q.LastReviewedAt.Equal(reviewedAt))

	err = q.ApplyReviewResult(domainquestion.ReviewResultVague, reviewedAt.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, domainquestion.MasteryNormal, q.MasteryStatus)

	err = q.ApplyReviewResult(domainquestion.ReviewResultMastered, reviewedAt.Add(2*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, domainquestion.MasteryMastered, q.MasteryStatus)
}
```

- [ ] **Step 2: Write failing occurrence and review domain tests**

Create `backend/internal/domain/question_occurrence/entity_test.go`:

```go
package question_occurrence_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestQuestionOccurrence_NewAllowsApplicationOnlyAssociation(t *testing.T) {
	occurredAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	occurrence, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                  "o1",
		UserID:              "u1",
		QuestionID:          "q1",
		ApplicationID:       "app1",
		CompanyNameSnapshot: "腾讯科技",
		JobTitleSnapshot:    "前端开发工程师",
		OccurredAt:          occurredAt,
	})

	require.NoError(t, err)
	assert.Equal(t, "q1", occurrence.QuestionID)
	assert.Nil(t, occurrence.InterviewProgressID)
	assert.Equal(t, "腾讯科技", occurrence.CompanyNameSnapshot)
}

func TestQuestionOccurrence_NewRejectsMissingQuestionOrApplication(t *testing.T) {
	_, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                  "o1",
		UserID:              "u1",
		QuestionID:          "",
		ApplicationID:       "app1",
		CompanyNameSnapshot: "Acme",
		JobTitleSnapshot:    "Frontend",
		OccurredAt:          time.Now(),
	})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestQuestionOccurrence_NewRejectsInvalidAnswerPerformance(t *testing.T) {
	performance := domainoccurrence.AnswerPerformance("unknown")
	_, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{
		ID:                  "o1",
		UserID:              "u1",
		QuestionID:          "q1",
		ApplicationID:       "app1",
		CompanyNameSnapshot: "Acme",
		JobTitleSnapshot:    "Frontend",
		OccurredAt:          time.Now(),
		AnswerPerformance:   &performance,
	})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}
```

Create `backend/internal/domain/flashcard_review/entity_test.go`:

```go
package flashcard_review_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

func TestFlashcardReview_NewValidatesAndCopiesFields(t *testing.T) {
	reviewedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	review, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{
		ID:         "r1",
		UserID:     "u1",
		QuestionID: "q1",
		ReviewedAt: reviewedAt,
		Result:     domainreview.ResultMastered,
		ReviewMode: domainreview.ModeRandom,
	})

	require.NoError(t, err)
	assert.Equal(t, "r1", review.ID)
	assert.Equal(t, domainreview.ResultMastered, review.Result)
	assert.Equal(t, domainreview.ModeRandom, review.ReviewMode)
}

func TestFlashcardReview_NewRejectsInvalidResult(t *testing.T) {
	_, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{
		ID:         "r1",
		UserID:     "u1",
		QuestionID: "q1",
		ReviewedAt: time.Now(),
		Result:     domainreview.Result("answered"),
		ReviewMode: domainreview.ModeRandom,
	})

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
}

func TestReviewSummary_RecordsLastReview(t *testing.T) {
	lastReviewedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	summary := domainreview.ReviewSummary{
		TotalReviews:   3,
		FailedCount:    1,
		VagueCount:     1,
		MasteredCount:  1,
		LastReviewedAt: &lastReviewedAt,
		LastResult:     ptr(domainreview.ResultMastered),
	}

	assert.Equal(t, int64(3), summary.TotalReviews)
	assert.Equal(t, domainreview.ResultMastered, *summary.LastResult)
}

func ptr[T any](v T) *T {
	return &v
}
```

- [ ] **Step 3: Run domain tests and verify they fail**

Run:

```bash
cd backend && go test ./internal/domain/interview_question ./internal/domain/question_occurrence ./internal/domain/flashcard_review
```

Expected: FAIL because the new packages and symbols do not exist.

- [ ] **Step 4: Implement interview question domain entity**

Create `backend/internal/domain/interview_question/entity.go`:

```go
package interview_question

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Difficulty string
type QuestionType string
type MasteryStatus string
type ReviewResult string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"

	TypeKnowledge  QuestionType = "knowledge"
	TypeScenario   QuestionType = "scenario"
	TypeCoding     QuestionType = "coding"
	TypeAlgorithm  QuestionType = "algorithm"
	TypeProject    QuestionType = "project"
	TypeBehavior   QuestionType = "behavior"
	TypeDiscussion QuestionType = "discussion"

	MasteryUnreviewed MasteryStatus = "unreviewed"
	MasteryWeak       MasteryStatus = "weak"
	MasteryNormal     MasteryStatus = "normal"
	MasteryMastered   MasteryStatus = "mastered"

	ReviewResultFailed   ReviewResult = "failed"
	ReviewResultVague    ReviewResult = "vague"
	ReviewResultMastered ReviewResult = "mastered"
)

type InterviewQuestion struct {
	ID             string         `gorm:"type:uuid;primaryKey"`
	UserID         string         `gorm:"type:varchar(36);not null;index"`
	Title          string         `gorm:"type:varchar(200);not null"`
	Content        string         `gorm:"type:text;not null"`
	AnswerAnalysis string         `gorm:"type:text;not null"`
	Tags           datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`
	Difficulty     Difficulty     `gorm:"type:varchar(20);not null;index"`
	QuestionType   QuestionType   `gorm:"type:varchar(30);not null;index"`
	MasteryStatus  MasteryStatus  `gorm:"type:varchar(30);not null;index"`
	IsFavorite     bool           `gorm:"not null;default:false;index"`
	Source         string         `gorm:"type:varchar(100)"`
	LastReviewedAt *time.Time     `gorm:"index"`
	Remark         string         `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type NewInterviewQuestionParams struct {
	ID             string
	UserID         string
	Title          string
	Content        string
	AnswerAnalysis string
	Tags           []string
	Difficulty     Difficulty
	QuestionType   QuestionType
	MasteryStatus  MasteryStatus
	IsFavorite     bool
	Source         string
	Remark         string
}

type UpdateSpec struct {
	Title          *string
	Content        *string
	AnswerAnalysis *string
	Tags           *[]string
	Difficulty     *Difficulty
	QuestionType   *QuestionType
	MasteryStatus  *MasteryStatus
	IsFavorite     *bool
	Source         *string
	Remark         *string
}

func NewInterviewQuestion(params NewInterviewQuestionParams) (*InterviewQuestion, error) {
	masteryStatus := params.MasteryStatus
	if masteryStatus == "" {
		masteryStatus = MasteryUnreviewed
	}
	tags, err := EncodeTags(params.Tags)
	if err != nil {
		return nil, err
	}
	q := &InterviewQuestion{
		ID:             params.ID,
		UserID:         params.UserID,
		Title:          params.Title,
		Content:        params.Content,
		AnswerAnalysis: params.AnswerAnalysis,
		Tags:           tags,
		Difficulty:     params.Difficulty,
		QuestionType:   params.QuestionType,
		MasteryStatus:  masteryStatus,
		IsFavorite:     params.IsFavorite,
		Source:         params.Source,
		Remark:         params.Remark,
	}
	if err := q.Validate(); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *InterviewQuestion) ApplyUpdate(update UpdateSpec) error {
	if update.Title != nil {
		q.Title = *update.Title
	}
	if update.Content != nil {
		q.Content = *update.Content
	}
	if update.AnswerAnalysis != nil {
		q.AnswerAnalysis = *update.AnswerAnalysis
	}
	if update.Tags != nil {
		tags, err := EncodeTags(*update.Tags)
		if err != nil {
			return err
		}
		q.Tags = tags
	}
	if update.Difficulty != nil {
		q.Difficulty = *update.Difficulty
	}
	if update.QuestionType != nil {
		q.QuestionType = *update.QuestionType
	}
	if update.MasteryStatus != nil {
		q.MasteryStatus = *update.MasteryStatus
	}
	if update.IsFavorite != nil {
		q.IsFavorite = *update.IsFavorite
	}
	if update.Source != nil {
		q.Source = *update.Source
	}
	if update.Remark != nil {
		q.Remark = *update.Remark
	}
	return q.Validate()
}

func (q *InterviewQuestion) ApplyReviewResult(result ReviewResult, reviewedAt time.Time) error {
	switch result {
	case ReviewResultFailed:
		q.MasteryStatus = MasteryWeak
	case ReviewResultVague:
		q.MasteryStatus = MasteryNormal
	case ReviewResultMastered:
		q.MasteryStatus = MasteryMastered
	default:
		return domainerrors.ErrBadParams
	}
	q.LastReviewedAt = &reviewedAt
	return q.Validate()
}

func (q *InterviewQuestion) Validate() error {
	if strings.TrimSpace(q.ID) == "" || strings.TrimSpace(q.UserID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(q.Title) == "" || strings.TrimSpace(q.Content) == "" || strings.TrimSpace(q.AnswerAnalysis) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(q.Title)) > 200 || len([]rune(q.Source)) > 100 {
		return domainerrors.ErrBadParams
	}
	if !q.Difficulty.Valid() || !q.QuestionType.Valid() || !q.MasteryStatus.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (d Difficulty) Valid() bool {
	switch d {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	default:
		return false
	}
}

func (t QuestionType) Valid() bool {
	switch t {
	case TypeKnowledge, TypeScenario, TypeCoding, TypeAlgorithm, TypeProject, TypeBehavior, TypeDiscussion:
		return true
	default:
		return false
	}
}

func (s MasteryStatus) Valid() bool {
	switch s {
	case MasteryUnreviewed, MasteryWeak, MasteryNormal, MasteryMastered:
		return true
	default:
		return false
	}
}

func (InterviewQuestion) TableName() string {
	return "interview_questions"
}
```

- [ ] **Step 5: Implement occurrence domain entity**

Create `backend/internal/domain/question_occurrence/entity.go`:

```go
package question_occurrence

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type AnswerPerformance string

const (
	PerformanceFailed   AnswerPerformance = "failed"
	PerformanceVague    AnswerPerformance = "vague"
	PerformanceAnswered AnswerPerformance = "answered"
	PerformanceGood     AnswerPerformance = "good"
)

type QuestionOccurrence struct {
	ID                     string             `gorm:"type:uuid;primaryKey"`
	UserID                 string             `gorm:"type:varchar(36);not null;index"`
	QuestionID             string             `gorm:"type:uuid;not null;index"`
	ApplicationID          string             `gorm:"type:uuid;not null;index"`
	InterviewProgressID    *string            `gorm:"type:uuid;index"`
	CompanyNameSnapshot    string             `gorm:"type:varchar(200);not null"`
	JobTitleSnapshot       string             `gorm:"type:varchar(200);not null"`
	InterviewStageSnapshot string             `gorm:"type:varchar(50)"`
	OccurredAt             time.Time          `gorm:"not null;index"`
	ActualQuestion         string             `gorm:"type:text"`
	AnswerPerformance      *AnswerPerformance `gorm:"type:varchar(30)"`
	Note                   string             `gorm:"type:text"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type NewQuestionOccurrenceParams struct {
	ID                     string
	UserID                 string
	QuestionID             string
	ApplicationID          string
	InterviewProgressID    *string
	CompanyNameSnapshot    string
	JobTitleSnapshot       string
	InterviewStageSnapshot string
	OccurredAt             time.Time
	ActualQuestion         string
	AnswerPerformance      *AnswerPerformance
	Note                   string
}

func NewQuestionOccurrence(params NewQuestionOccurrenceParams) (*QuestionOccurrence, error) {
	occurrence := &QuestionOccurrence{
		ID:                     params.ID,
		UserID:                 params.UserID,
		QuestionID:             params.QuestionID,
		ApplicationID:          params.ApplicationID,
		InterviewProgressID:    params.InterviewProgressID,
		CompanyNameSnapshot:    params.CompanyNameSnapshot,
		JobTitleSnapshot:       params.JobTitleSnapshot,
		InterviewStageSnapshot: params.InterviewStageSnapshot,
		OccurredAt:             params.OccurredAt,
		ActualQuestion:         params.ActualQuestion,
		AnswerPerformance:      params.AnswerPerformance,
		Note:                   params.Note,
	}
	if err := occurrence.Validate(); err != nil {
		return nil, err
	}
	return occurrence, nil
}

func (o *QuestionOccurrence) Validate() error {
	if strings.TrimSpace(o.ID) == "" || strings.TrimSpace(o.UserID) == "" || strings.TrimSpace(o.QuestionID) == "" || strings.TrimSpace(o.ApplicationID) == "" {
		return domainerrors.ErrBadParams
	}
	if strings.TrimSpace(o.CompanyNameSnapshot) == "" || strings.TrimSpace(o.JobTitleSnapshot) == "" {
		return domainerrors.ErrBadParams
	}
	if len([]rune(o.CompanyNameSnapshot)) > 200 || len([]rune(o.JobTitleSnapshot)) > 200 || len([]rune(o.InterviewStageSnapshot)) > 50 {
		return domainerrors.ErrBadParams
	}
	if o.OccurredAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if o.AnswerPerformance != nil && !o.AnswerPerformance.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (p AnswerPerformance) Valid() bool {
	switch p {
	case PerformanceFailed, PerformanceVague, PerformanceAnswered, PerformanceGood:
		return true
	default:
		return false
	}
}

func (QuestionOccurrence) TableName() string {
	return "question_occurrences"
}
```

- [ ] **Step 6: Implement flashcard review domain entity**

Create `backend/internal/domain/flashcard_review/entity.go`:

```go
package flashcard_review

import (
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type Result string
type ReviewMode string

const (
	ResultFailed   Result = "failed"
	ResultVague    Result = "vague"
	ResultMastered Result = "mastered"

	ModeRandom        ReviewMode = "random"
	ModeTag           ReviewMode = "tag"
	ModeDifficulty    ReviewMode = "difficulty"
	ModeFavoriteFirst ReviewMode = "favorite_first"
	ModeWeakFirst     ReviewMode = "weak_first"
	ModeFrequentFirst ReviewMode = "frequent_first"
)

type FlashcardReview struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	UserID     string         `gorm:"type:varchar(36);not null;index"`
	QuestionID string         `gorm:"type:uuid;not null;index"`
	ReviewedAt time.Time      `gorm:"not null;index"`
	Result     Result         `gorm:"type:varchar(30);not null"`
	ReviewMode ReviewMode     `gorm:"type:varchar(30);not null"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type NewFlashcardReviewParams struct {
	ID         string
	UserID     string
	QuestionID string
	ReviewedAt time.Time
	Result     Result
	ReviewMode ReviewMode
}

type ReviewSummary struct {
	TotalReviews   int64
	FailedCount    int64
	VagueCount     int64
	MasteredCount  int64
	LastReviewedAt *time.Time
	LastResult     *Result
}

func NewFlashcardReview(params NewFlashcardReviewParams) (*FlashcardReview, error) {
	review := &FlashcardReview{
		ID:         params.ID,
		UserID:     params.UserID,
		QuestionID: params.QuestionID,
		ReviewedAt: params.ReviewedAt,
		Result:     params.Result,
		ReviewMode: params.ReviewMode,
	}
	if err := review.Validate(); err != nil {
		return nil, err
	}
	return review, nil
}

func (r *FlashcardReview) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.UserID) == "" || strings.TrimSpace(r.QuestionID) == "" {
		return domainerrors.ErrBadParams
	}
	if r.ReviewedAt.IsZero() {
		return domainerrors.ErrBadParams
	}
	if !r.Result.Valid() || !r.ReviewMode.Valid() {
		return domainerrors.ErrBadParams
	}
	return nil
}

func (r Result) Valid() bool {
	switch r {
	case ResultFailed, ResultVague, ResultMastered:
		return true
	default:
		return false
	}
}

func (m ReviewMode) Valid() bool {
	switch m {
	case ModeRandom, ModeTag, ModeDifficulty, ModeFavoriteFirst, ModeWeakFirst, ModeFrequentFirst:
		return true
	default:
		return false
	}
}

func (FlashcardReview) TableName() string {
	return "flashcard_reviews"
}
```

- [ ] **Step 7: Run domain tests and verify they pass**

Run:

```bash
cd backend && go test ./internal/domain/interview_question ./internal/domain/question_occurrence ./internal/domain/flashcard_review
```

Expected: PASS.

- [ ] **Step 8: Commit domain entities**

Run:

```bash
git add backend/internal/domain/interview_question backend/internal/domain/question_occurrence backend/internal/domain/flashcard_review
git commit -m "feat: add interview question bank domain models"
```

---

## Task 2: Add DTOs and frontend-aligned converters

**Files:**
- Create: `backend/internal/application/dto/interview_question_dto.go`
- Test: `backend/internal/application/dto/job_dto_test.go`

- [ ] **Step 1: Write failing DTO tests**

Append to `backend/internal/application/dto/job_dto_test.go`:

```go
func TestInterviewQuestionResponse_UsesFrontendShape(t *testing.T) {
	lastReviewedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	question := &domainquestion.InterviewQuestion{
		ID:             "q1",
		UserID:         "u1",
		Title:          "Vue 2 响应式原理",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Tags:           mustTags(t, []string{"Vue", "JavaScript"}),
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryNormal,
		IsFavorite:     true,
		Source:         "面试复盘",
		LastReviewedAt: &lastReviewedAt,
		Remark:         "remark",
		CreatedAt:      lastReviewedAt.Add(-time.Hour),
		UpdatedAt:      lastReviewedAt,
	}

	resp := dto.ToInterviewQuestionResponse(question, 2)
	body, err := json.Marshal(resp)
	require.NoError(t, err)

	assert.Contains(t, string(body), "answerAnalysis")
	assert.Contains(t, string(body), "questionType")
	assert.Contains(t, string(body), "masteryStatus")
	assert.Contains(t, string(body), "isFavorite")
	assert.Contains(t, string(body), "occurrenceCount")
	assert.NotContains(t, string(body), "AnswerAnalysis")
	assert.Equal(t, int64(2), resp.OccurrenceCount)
	assert.Equal(t, []string{"Vue", "JavaScript"}, resp.Tags)
}

func TestInterviewQuestionDetailResponse_IncludesOccurrencesAndReviewSummary(t *testing.T) {
	now := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	question := &domainquestion.InterviewQuestion{
		ID:             "q1",
		UserID:         "u1",
		Title:          "title",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Tags:           mustTags(t, []string{}),
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	performance := domainoccurrence.PerformanceGood
	occurrence := domainoccurrence.QuestionOccurrence{
		ID:                  "o1",
		UserID:              "u1",
		QuestionID:          "q1",
		ApplicationID:       "app1",
		CompanyNameSnapshot: "Acme",
		JobTitleSnapshot:    "Frontend",
		OccurredAt:          now,
		AnswerPerformance:   &performance,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	review := domainreview.FlashcardReview{ID: "r1", UserID: "u1", QuestionID: "q1", ReviewedAt: now, Result: domainreview.ResultMastered, ReviewMode: domainreview.ModeRandom, CreatedAt: now}
	summary := domainreview.ReviewSummary{TotalReviews: 1, MasteredCount: 1, LastReviewedAt: &now, LastResult: &review.Result}

	resp := dto.ToInterviewQuestionDetailResponse(question, 1, []domainoccurrence.QuestionOccurrence{occurrence}, summary, []domainreview.FlashcardReview{review})

	assert.Len(t, resp.Occurrences, 1)
	assert.Equal(t, "good", string(*resp.Occurrences[0].AnswerPerformance))
	assert.Equal(t, int64(1), resp.ReviewSummary.TotalReviews)
	assert.Len(t, resp.RecentReviews, 1)
}
```

Add imports to `backend/internal/application/dto/job_dto_test.go`:

```go
import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)
```

If the file already has imports, merge these imports instead of creating a second import block.

Add helper to the same test file:

```go
func mustTags(t *testing.T, tags []string) datatypes.JSON {
	t.Helper()
	encoded, err := domainquestion.EncodeTags(tags)
	require.NoError(t, err)
	return encoded
}
```

Also add import:

```go
"gorm.io/datatypes"
```

- [ ] **Step 2: Run DTO tests and verify they fail**

Run:

```bash
cd backend && go test ./internal/application/dto
```

Expected: FAIL because DTO types and `EncodeTags` helpers do not exist.

- [ ] **Step 3: Add tag helpers to question domain**

Add to `backend/internal/domain/interview_question/entity.go` after imports and before entity constructors:

```go
func EncodeTags(tags []string) (datatypes.JSON, error) {
	if tags == nil {
		tags = []string{}
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return nil, domainerrors.ErrBadParams
	}
	return datatypes.JSON(data), nil
}

func DecodeTags(data datatypes.JSON) []string {
	if len(data) == 0 {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil {
		return []string{}
	}
	if tags == nil {
		return []string{}
	}
	return tags
}
```

Update the import block in `backend/internal/domain/interview_question/entity.go`:

```go
import (
	"encoding/json"
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)
```

- [ ] **Step 4: Implement DTOs**

Create `backend/internal/application/dto/interview_question_dto.go`:

```go
package dto

import (
	"time"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)

type CreateInterviewQuestionRequest struct {
	Title          string                       `json:"title" binding:"required,min=1,max=200"`
	Content        string                       `json:"content" binding:"required"`
	AnswerAnalysis string                       `json:"answerAnalysis" binding:"required"`
	Tags           []string                     `json:"tags"`
	Difficulty     domainquestion.Difficulty   `json:"difficulty" binding:"required"`
	QuestionType   domainquestion.QuestionType `json:"questionType" binding:"required"`
	MasteryStatus  domainquestion.MasteryStatus `json:"masteryStatus" binding:"required"`
	IsFavorite     bool                         `json:"isFavorite"`
	Source         string                       `json:"source" binding:"max=100"`
	Remark         string                       `json:"remark"`
}

type UpdateInterviewQuestionRequest struct {
	Title          *string                       `json:"title" binding:"omitempty,min=1,max=200"`
	Content        *string                       `json:"content"`
	AnswerAnalysis *string                       `json:"answerAnalysis"`
	Tags           *[]string                     `json:"tags"`
	Difficulty     *domainquestion.Difficulty   `json:"difficulty"`
	QuestionType   *domainquestion.QuestionType `json:"questionType"`
	MasteryStatus  *domainquestion.MasteryStatus `json:"masteryStatus"`
	IsFavorite     *bool                         `json:"isFavorite"`
	Source         *string                       `json:"source" binding:"omitempty,max=100"`
	Remark         *string                       `json:"remark"`
}

type UpdateQuestionFavoriteRequest struct {
	IsFavorite bool `json:"isFavorite"`
}

type InterviewQuestionQuery struct {
	Keyword       string                       `json:"keyword" form:"keyword"`
	Tags          string                       `json:"tags" form:"tags"`
	Difficulty    domainquestion.Difficulty   `json:"difficulty" form:"difficulty"`
	QuestionType  domainquestion.QuestionType `json:"questionType" form:"questionType"`
	MasteryStatus domainquestion.MasteryStatus `json:"masteryStatus" form:"masteryStatus"`
	IsFavorite    *bool                        `json:"isFavorite" form:"isFavorite"`
	Page          int                          `json:"page" form:"page"`
	PageSize      int                          `json:"pageSize" form:"pageSize"`
}

type InterviewQuestionResponse struct {
	ID              string                       `json:"id"`
	Title           string                       `json:"title"`
	Content         string                       `json:"content"`
	AnswerAnalysis  string                       `json:"answerAnalysis"`
	Tags            []string                     `json:"tags"`
	Difficulty      domainquestion.Difficulty   `json:"difficulty"`
	QuestionType    domainquestion.QuestionType `json:"questionType"`
	MasteryStatus   domainquestion.MasteryStatus `json:"masteryStatus"`
	IsFavorite      bool                         `json:"isFavorite"`
	Source          string                       `json:"source,omitempty"`
	OccurrenceCount int64                        `json:"occurrenceCount"`
	LastReviewedAt  *time.Time                   `json:"lastReviewedAt,omitempty"`
	Remark          string                       `json:"remark,omitempty"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
}

type InterviewQuestionDetailResponse struct {
	*InterviewQuestionResponse
	Occurrences   []*QuestionOccurrenceResponse `json:"occurrences"`
	ReviewSummary *FlashcardReviewSummaryResponse `json:"reviewSummary"`
	RecentReviews []*FlashcardReviewResponse `json:"recentReviews"`
}

type CreateQuestionOccurrenceRequest struct {
	QuestionID             string                              `json:"questionId" binding:"required"`
	ApplicationID          string                              `json:"applicationId" binding:"required"`
	InterviewProgressID    *string                             `json:"interviewProgressId"`
	InterviewStageSnapshot string                              `json:"interviewStageSnapshot"`
	OccurredAt             FlexibleTime                        `json:"occurredAt" binding:"required"`
	ActualQuestion         string                              `json:"actualQuestion"`
	AnswerPerformance      *domainoccurrence.AnswerPerformance `json:"answerPerformance"`
	Note                   string                              `json:"note"`
}

type QuestionOccurrenceQuery struct {
	QuestionID    string `json:"questionId" form:"questionId"`
	ApplicationID string `json:"applicationId" form:"applicationId"`
}

type QuestionOccurrenceResponse struct {
	ID                     string                              `json:"id"`
	QuestionID             string                              `json:"questionId"`
	ApplicationID          string                              `json:"applicationId"`
	InterviewProgressID    *string                             `json:"interviewProgressId,omitempty"`
	CompanyNameSnapshot    string                              `json:"companyNameSnapshot"`
	JobTitleSnapshot       string                              `json:"jobTitleSnapshot"`
	InterviewStageSnapshot string                              `json:"interviewStageSnapshot,omitempty"`
	OccurredAt             time.Time                           `json:"occurredAt"`
	ActualQuestion         string                              `json:"actualQuestion,omitempty"`
	AnswerPerformance      *domainoccurrence.AnswerPerformance `json:"answerPerformance,omitempty"`
	Note                   string                              `json:"note,omitempty"`
	CreatedAt              time.Time                           `json:"createdAt"`
	UpdatedAt              time.Time                           `json:"updatedAt"`
}

type FlashcardQuery struct {
	Mode       domainreview.ReviewMode      `json:"mode" form:"mode" binding:"required"`
	Tags       string                       `json:"tags" form:"tags"`
	Difficulty domainquestion.Difficulty   `json:"difficulty" form:"difficulty"`
	Limit      int                          `json:"limit" form:"limit"`
}

type CreateFlashcardReviewRequest struct {
	QuestionID string                  `json:"questionId" binding:"required"`
	ReviewedAt *FlexibleTime           `json:"reviewedAt"`
	Result     domainreview.Result     `json:"result" binding:"required"`
	ReviewMode domainreview.ReviewMode `json:"reviewMode" binding:"required"`
}

type FlashcardReviewResponse struct {
	ID         string                  `json:"id"`
	QuestionID string                  `json:"questionId"`
	ReviewedAt time.Time               `json:"reviewedAt"`
	Result     domainreview.Result     `json:"result"`
	ReviewMode domainreview.ReviewMode `json:"reviewMode"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type FlashcardReviewSummaryResponse struct {
	TotalReviews   int64                `json:"totalReviews"`
	FailedCount    int64                `json:"failedCount"`
	VagueCount     int64                `json:"vagueCount"`
	MasteredCount  int64                `json:"masteredCount"`
	LastReviewedAt *time.Time           `json:"lastReviewedAt,omitempty"`
	LastResult     *domainreview.Result `json:"lastResult,omitempty"`
}

func ToInterviewQuestionResponse(q *domainquestion.InterviewQuestion, occurrenceCount int64) *InterviewQuestionResponse {
	if q == nil {
		return nil
	}
	return &InterviewQuestionResponse{
		ID:              q.ID,
		Title:           q.Title,
		Content:         q.Content,
		AnswerAnalysis:  q.AnswerAnalysis,
		Tags:            domainquestion.DecodeTags(q.Tags),
		Difficulty:      q.Difficulty,
		QuestionType:    q.QuestionType,
		MasteryStatus:   q.MasteryStatus,
		IsFavorite:      q.IsFavorite,
		Source:          q.Source,
		OccurrenceCount: occurrenceCount,
		LastReviewedAt:  q.LastReviewedAt,
		Remark:          q.Remark,
		CreatedAt:       q.CreatedAt,
		UpdatedAt:       q.UpdatedAt,
	}
}

func ToInterviewQuestionDetailResponse(q *domainquestion.InterviewQuestion, occurrenceCount int64, occurrences []domainoccurrence.QuestionOccurrence, summary domainreview.ReviewSummary, reviews []domainreview.FlashcardReview) *InterviewQuestionDetailResponse {
	return &InterviewQuestionDetailResponse{
		InterviewQuestionResponse: ToInterviewQuestionResponse(q, occurrenceCount),
		Occurrences:              ToQuestionOccurrenceResponses(occurrences),
		ReviewSummary:            ToFlashcardReviewSummaryResponse(summary),
		RecentReviews:            ToFlashcardReviewResponses(reviews),
	}
}

func ToQuestionOccurrenceResponse(o *domainoccurrence.QuestionOccurrence) *QuestionOccurrenceResponse {
	if o == nil {
		return nil
	}
	return &QuestionOccurrenceResponse{
		ID:                     o.ID,
		QuestionID:             o.QuestionID,
		ApplicationID:          o.ApplicationID,
		InterviewProgressID:    o.InterviewProgressID,
		CompanyNameSnapshot:    o.CompanyNameSnapshot,
		JobTitleSnapshot:       o.JobTitleSnapshot,
		InterviewStageSnapshot: o.InterviewStageSnapshot,
		OccurredAt:             o.OccurredAt,
		ActualQuestion:         o.ActualQuestion,
		AnswerPerformance:      o.AnswerPerformance,
		Note:                   o.Note,
		CreatedAt:              o.CreatedAt,
		UpdatedAt:              o.UpdatedAt,
	}
}

func ToQuestionOccurrenceResponses(rows []domainoccurrence.QuestionOccurrence) []*QuestionOccurrenceResponse {
	responses := make([]*QuestionOccurrenceResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, ToQuestionOccurrenceResponse(&rows[i]))
	}
	return responses
}

func ToFlashcardReviewResponse(r *domainreview.FlashcardReview) *FlashcardReviewResponse {
	if r == nil {
		return nil
	}
	return &FlashcardReviewResponse{
		ID:         r.ID,
		QuestionID: r.QuestionID,
		ReviewedAt: r.ReviewedAt,
		Result:     r.Result,
		ReviewMode: r.ReviewMode,
		CreatedAt:  r.CreatedAt,
	}
}

func ToFlashcardReviewResponses(rows []domainreview.FlashcardReview) []*FlashcardReviewResponse {
	responses := make([]*FlashcardReviewResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, ToFlashcardReviewResponse(&rows[i]))
	}
	return responses
}

func ToFlashcardReviewSummaryResponse(summary domainreview.ReviewSummary) *FlashcardReviewSummaryResponse {
	return &FlashcardReviewSummaryResponse{
		TotalReviews:   summary.TotalReviews,
		FailedCount:    summary.FailedCount,
		VagueCount:     summary.VagueCount,
		MasteredCount:  summary.MasteredCount,
		LastReviewedAt: summary.LastReviewedAt,
		LastResult:     summary.LastResult,
	}
}
```

- [ ] **Step 5: Run DTO tests and verify they pass**

Run:

```bash
cd backend && go test ./internal/application/dto
```

Expected: PASS.

- [ ] **Step 6: Commit DTOs**

Run:

```bash
git add backend/internal/application/dto/interview_question_dto.go backend/internal/application/dto/job_dto_test.go backend/internal/domain/interview_question/entity.go
git commit -m "feat: add interview question DTOs"
```

---

## Task 3: Add repository interfaces and migrations

**Files:**
- Create: `backend/internal/domain/interview_question/repository/repository.go`
- Create: `backend/internal/domain/question_occurrence/repository/repository.go`
- Create: `backend/internal/domain/flashcard_review/repository/repository.go`
- Create: `backend/migrations/000007_create_interview_questions.up.sql`
- Create: `backend/migrations/000007_create_interview_questions.down.sql`

- [ ] **Step 1: Add repository interfaces**

Create `backend/internal/domain/interview_question/repository/repository.go`:

```go
package repository

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
)

type ListFilter struct {
	UserID        string
	Keyword       string
	Tags          []string
	Difficulty    domainquestion.Difficulty
	QuestionType  domainquestion.QuestionType
	MasteryStatus domainquestion.MasteryStatus
	IsFavorite    *bool
	Page          int
	PageSize      int
}

type FlashcardFilter struct {
	UserID     string
	Mode       domainreview.ReviewMode
	Tags       []string
	Difficulty domainquestion.Difficulty
	Limit      int
}

type QuestionWithOccurrenceCount struct {
	Question        domainquestion.InterviewQuestion
	OccurrenceCount int64
}

type Repository interface {
	Create(ctx context.Context, q *domainquestion.InterviewQuestion) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error)
	Update(ctx context.Context, q *domainquestion.InterviewQuestion) error
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]QuestionWithOccurrenceCount, int64, error)
	CountOccurrences(ctx context.Context, questionID, userID string) (int64, error)
	ListFlashcards(ctx context.Context, f FlashcardFilter) ([]QuestionWithOccurrenceCount, int64, error)
}
```

Create `backend/internal/domain/question_occurrence/repository/repository.go`:

```go
package repository

import (
	"context"

	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
)

type ListFilter struct {
	UserID        string
	QuestionID    string
	ApplicationID string
}

type Repository interface {
	Create(ctx context.Context, occurrence *domainoccurrence.QuestionOccurrence) error
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error)
	Delete(ctx context.Context, id, userID string) error
	DeleteByQuestionID(ctx context.Context, questionID, userID string) error
	ListByUser(ctx context.Context, f ListFilter) ([]domainoccurrence.QuestionOccurrence, error)
}
```

Create `backend/internal/domain/flashcard_review/repository/repository.go`:

```go
package repository

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
)

type Repository interface {
	Create(ctx context.Context, review *domainreview.FlashcardReview) error
	DeleteByQuestionID(ctx context.Context, questionID, userID string) error
	ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error)
	SummaryByQuestion(ctx context.Context, questionID, userID string) (domainreview.ReviewSummary, error)
}
```

- [ ] **Step 2: Write migration up file**

Create `backend/migrations/000007_create_interview_questions.up.sql`:

```sql
CREATE TABLE interview_questions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          VARCHAR(36) NOT NULL REFERENCES users(id),
    title            VARCHAR(200) NOT NULL,
    content          TEXT NOT NULL,
    answer_analysis  TEXT NOT NULL,
    tags             JSONB NOT NULL DEFAULT '[]',
    difficulty       VARCHAR(20) NOT NULL,
    question_type    VARCHAR(30) NOT NULL,
    mastery_status   VARCHAR(30) NOT NULL DEFAULT 'unreviewed',
    is_favorite      BOOLEAN NOT NULL DEFAULT FALSE,
    source           VARCHAR(100),
    last_reviewed_at TIMESTAMPTZ,
    remark           TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE TABLE question_occurrences (
    id                       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                  VARCHAR(36) NOT NULL REFERENCES users(id),
    question_id              UUID NOT NULL REFERENCES interview_questions(id),
    application_id           UUID NOT NULL REFERENCES resume_applications(id),
    interview_progress_id    UUID REFERENCES interview_progress(id),
    company_name_snapshot    VARCHAR(200) NOT NULL,
    job_title_snapshot       VARCHAR(200) NOT NULL,
    interview_stage_snapshot VARCHAR(50),
    occurred_at              TIMESTAMPTZ NOT NULL,
    actual_question          TEXT,
    answer_performance       VARCHAR(30),
    note                     TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at               TIMESTAMPTZ
);

CREATE TABLE flashcard_reviews (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     VARCHAR(36) NOT NULL REFERENCES users(id),
    question_id UUID NOT NULL REFERENCES interview_questions(id),
    reviewed_at TIMESTAMPTZ NOT NULL,
    result      VARCHAR(30) NOT NULL,
    review_mode VARCHAR(30) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_interview_questions_user_id ON interview_questions(user_id);
CREATE INDEX idx_interview_questions_user_difficulty ON interview_questions(user_id, difficulty);
CREATE INDEX idx_interview_questions_user_question_type ON interview_questions(user_id, question_type);
CREATE INDEX idx_interview_questions_user_mastery_status ON interview_questions(user_id, mastery_status);
CREATE INDEX idx_interview_questions_user_is_favorite ON interview_questions(user_id, is_favorite);
CREATE INDEX idx_interview_questions_tags ON interview_questions USING GIN(tags);
CREATE INDEX idx_interview_questions_deleted_at ON interview_questions(deleted_at);

CREATE INDEX idx_question_occurrences_user_id ON question_occurrences(user_id);
CREATE INDEX idx_question_occurrences_question_id ON question_occurrences(question_id);
CREATE INDEX idx_question_occurrences_application_id ON question_occurrences(application_id);
CREATE INDEX idx_question_occurrences_progress_id ON question_occurrences(interview_progress_id);
CREATE INDEX idx_question_occurrences_occurred_at ON question_occurrences(occurred_at);
CREATE UNIQUE INDEX idx_question_occurrences_unique_progress
    ON question_occurrences(user_id, question_id, application_id, interview_progress_id)
    WHERE interview_progress_id IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_flashcard_reviews_user_id ON flashcard_reviews(user_id);
CREATE INDEX idx_flashcard_reviews_question_id ON flashcard_reviews(question_id);
CREATE INDEX idx_flashcard_reviews_reviewed_at ON flashcard_reviews(reviewed_at);
```

- [ ] **Step 3: Write migration down file**

Create `backend/migrations/000007_create_interview_questions.down.sql`:

```sql
DROP TABLE IF EXISTS flashcard_reviews;
DROP TABLE IF EXISTS question_occurrences;
DROP TABLE IF EXISTS interview_questions;
```

- [ ] **Step 4: Run backend tests to catch migration naming issues**

Run:

```bash
cd backend && go test ./migrations ./internal/domain/...
```

Expected: PASS.

- [ ] **Step 5: Commit interfaces and migrations**

Run:

```bash
git add backend/internal/domain/interview_question/repository backend/internal/domain/question_occurrence/repository backend/internal/domain/flashcard_review/repository backend/migrations/000007_create_interview_questions.up.sql backend/migrations/000007_create_interview_questions.down.sql
git commit -m "feat: add interview question migrations and repository contracts"
```

---

## Task 4: Implement application service with mocked tests

**Files:**
- Create: `backend/internal/application/interview_question/interview_question_service.go`
- Create: `backend/internal/application/interview_question/interview_question_service_test.go`

- [ ] **Step 1: Write failing service tests for independent questions and list filters**

Create `backend/internal/application/interview_question/interview_question_service_test.go` with package, imports, mocks, and first tests:

```go
package interview_question_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	questionsvc "github.com/yilin/ai-for-backend/internal/application/interview_question"
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

type mockQuestionRepo struct{ mock.Mock }

func (m *mockQuestionRepo) Create(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	return m.Called(ctx, q).Error(0)
}
func (m *mockQuestionRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, id, userID)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Error(1)
}
func (m *mockQuestionRepo) Update(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	return m.Called(ctx, q).Error(0)
}
func (m *mockQuestionRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *mockQuestionRepo) ListByUser(ctx context.Context, f questionrepo.ListFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) CountOccurrences(ctx context.Context, questionID, userID string) (int64, error) {
	args := m.Called(ctx, questionID, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockQuestionRepo) ListFlashcards(ctx context.Context, f questionrepo.FlashcardFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}

type mockOccurrenceRepo struct{ mock.Mock }

func (m *mockOccurrenceRepo) Create(ctx context.Context, occurrence *domainoccurrence.QuestionOccurrence) error {
	return m.Called(ctx, occurrence).Error(0)
}
func (m *mockOccurrenceRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, id, userID)
	occurrence, _ := args.Get(0).(*domainoccurrence.QuestionOccurrence)
	return occurrence, args.Error(1)
}
func (m *mockOccurrenceRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *mockOccurrenceRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return m.Called(ctx, questionID, userID).Error(0)
}
func (m *mockOccurrenceRepo) ListByUser(ctx context.Context, f occurrencerepo.ListFilter) ([]domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domainoccurrence.QuestionOccurrence), args.Error(1)
}

type mockReviewRepo struct{ mock.Mock }

func (m *mockReviewRepo) Create(ctx context.Context, review *domainreview.FlashcardReview) error {
	return m.Called(ctx, review).Error(0)
}
func (m *mockReviewRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return m.Called(ctx, questionID, userID).Error(0)
}
func (m *mockReviewRepo) ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error) {
	args := m.Called(ctx, questionID, userID, limit)
	return args.Get(0).([]domainreview.FlashcardReview), args.Error(1)
}
func (m *mockReviewRepo) SummaryByQuestion(ctx context.Context, questionID, userID string) (domainreview.ReviewSummary, error) {
	args := m.Called(ctx, questionID, userID)
	return args.Get(0).(domainreview.ReviewSummary), args.Error(1)
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

func TestInterviewQuestionService_CreateQuestionDoesNotRequireApplication(t *testing.T) {
	questions := &mockQuestionRepo{}
	occurrences := &mockOccurrenceRepo{}
	reviews := &mockReviewRepo{}
	questions.On("Create", mock.Anything, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.ID != "" && q.UserID == "u1" && q.Title == "Vue 2 响应式原理" && q.IsFavorite
	})).Return(nil)
	svc := questionsvc.NewInterviewQuestionService(questions, occurrences, reviews, nil, nil)

	got, err := svc.CreateQuestion(context.Background(), "u1", dto.CreateInterviewQuestionRequest{
		Title:          "Vue 2 响应式原理",
		Content:        "content",
		AnswerAnalysis: "analysis",
		Tags:           []string{"Vue"},
		Difficulty:     domainquestion.DifficultyMedium,
		QuestionType:   domainquestion.TypeKnowledge,
		MasteryStatus:  domainquestion.MasteryUnreviewed,
		IsFavorite:     true,
	})

	require.NoError(t, err)
	assert.Equal(t, "u1", got.UserID)
	questions.AssertExpectations(t)
	occurrences.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInterviewQuestionService_ListDefaultsPaginationAndSplitsTags(t *testing.T) {
	questions := &mockQuestionRepo{}
	rows := []questionrepo.QuestionWithOccurrenceCount{{Question: domainquestion.InterviewQuestion{ID: "q1", UserID: "u1"}, OccurrenceCount: 2}}
	questions.On("ListByUser", mock.Anything, questionrepo.ListFilter{
		UserID:        "u1",
		Keyword:       "Vue",
		Tags:          []string{"Vue", "浏览器"},
		Difficulty:    domainquestion.DifficultyMedium,
		QuestionType:  domainquestion.TypeKnowledge,
		MasteryStatus: domainquestion.MasteryWeak,
		Page:          1,
		PageSize:      10,
	}).Return(rows, int64(1), nil)
	svc := questionsvc.NewInterviewQuestionService(questions, nil, nil, nil, nil)

	got, total, err := svc.ListQuestions(context.Background(), "u1", dto.InterviewQuestionQuery{
		Keyword:       "Vue",
		Tags:          "Vue, 浏览器",
		Difficulty:    domainquestion.DifficultyMedium,
		QuestionType:  domainquestion.TypeKnowledge,
		MasteryStatus: domainquestion.MasteryWeak,
		Page:          0,
		PageSize:      0,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, got, 1)
	questions.AssertExpectations(t)
}
```

- [ ] **Step 2: Write failing service tests for occurrences and reviews**

Append to `backend/internal/application/interview_question/interview_question_service_test.go`:

```go
func TestInterviewQuestionService_CreateOccurrenceCopiesApplicationSnapshotAndProgressStage(t *testing.T) {
	questions := &mockQuestionRepo{}
	occurrences := &mockOccurrenceRepo{}
	reviews := &mockReviewRepo{}
	apps := &mockApplicationReader{}
	progresses := &mockProgressReader{}
	occurredAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	progressID := "p1"
	questions.On("FindByIDAndUserID", mock.Anything, "q1", "u1").Return(&domainquestion.InterviewQuestion{ID: "q1", UserID: "u1"}, nil)
	apps.On("Get", mock.Anything, "u1", "app1").Return(&domainapp.ResumeApplication{ID: "app1", UserID: "u1", CompanyName: "腾讯科技", JobTitle: "前端开发工程师"}, nil)
	progresses.On("FindByIDAndUserID", mock.Anything, "p1", "u1").Return(&domainprogress.InterviewProgress{ID: "p1", UserID: "u1", ApplicationID: "app1", Stage: domainprogress.StageFirstInterview}, nil)
	occurrences.On("Create", mock.Anything, mock.MatchedBy(func(o *domainoccurrence.QuestionOccurrence) bool {
		return o.ID != "" && o.UserID == "u1" && o.QuestionID == "q1" && o.ApplicationID == "app1" && o.CompanyNameSnapshot == "腾讯科技" && o.JobTitleSnapshot == "前端开发工程师" && o.InterviewStageSnapshot == "first_interview"
	})).Return(nil)
	svc := questionsvc.NewInterviewQuestionService(questions, occurrences, reviews, apps, progresses)

	got, err := svc.CreateOccurrence(context.Background(), "u1", dto.CreateQuestionOccurrenceRequest{
		QuestionID:          "q1",
		ApplicationID:       "app1",
		InterviewProgressID: &progressID,
		OccurredAt:          dto.FlexibleTime{Time: occurredAt},
		ActualQuestion:      "实际问法",
	})

	require.NoError(t, err)
	assert.Equal(t, "腾讯科技", got.CompanyNameSnapshot)
	assert.Equal(t, "first_interview", got.InterviewStageSnapshot)
	questions.AssertExpectations(t)
	apps.AssertExpectations(t)
	progresses.AssertExpectations(t)
	occurrences.AssertExpectations(t)
}

func TestInterviewQuestionService_CreateOccurrenceRejectsProgressFromDifferentApplication(t *testing.T) {
	questions := &mockQuestionRepo{}
	occurrences := &mockOccurrenceRepo{}
	apps := &mockApplicationReader{}
	progresses := &mockProgressReader{}
	progressID := "p1"
	questions.On("FindByIDAndUserID", mock.Anything, "q1", "u1").Return(&domainquestion.InterviewQuestion{ID: "q1", UserID: "u1"}, nil)
	apps.On("Get", mock.Anything, "u1", "app1").Return(&domainapp.ResumeApplication{ID: "app1", UserID: "u1", CompanyName: "Acme", JobTitle: "Frontend"}, nil)
	progresses.On("FindByIDAndUserID", mock.Anything, "p1", "u1").Return(&domainprogress.InterviewProgress{ID: "p1", UserID: "u1", ApplicationID: "other-app", Stage: domainprogress.StageFirstInterview}, nil)
	svc := questionsvc.NewInterviewQuestionService(questions, occurrences, nil, apps, progresses)

	got, err := svc.CreateOccurrence(context.Background(), "u1", dto.CreateQuestionOccurrenceRequest{QuestionID: "q1", ApplicationID: "app1", InterviewProgressID: &progressID, OccurredAt: dto.FlexibleTime{Time: time.Now()}})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	occurrences.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInterviewQuestionService_CreateFlashcardReviewCreatesHistoryAndUpdatesQuestion(t *testing.T) {
	questions := &mockQuestionRepo{}
	reviews := &mockReviewRepo{}
	reviewedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	question := &domainquestion.InterviewQuestion{ID: "q1", UserID: "u1", Title: "title", Content: "content", AnswerAnalysis: "analysis", Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryUnreviewed}
	questions.On("FindByIDAndUserID", mock.Anything, "q1", "u1").Return(question, nil)
	reviews.On("Create", mock.Anything, mock.MatchedBy(func(r *domainreview.FlashcardReview) bool {
		return r.ID != "" && r.QuestionID == "q1" && r.UserID == "u1" && r.Result == domainreview.ResultMastered && r.ReviewMode == domainreview.ModeRandom
	})).Return(nil)
	questions.On("Update", mock.Anything, mock.MatchedBy(func(q *domainquestion.InterviewQuestion) bool {
		return q.ID == "q1" && q.MasteryStatus == domainquestion.MasteryMastered && q.LastReviewedAt != nil && q.LastReviewedAt.Equal(reviewedAt)
	})).Return(nil)
	summary := domainreview.ReviewSummary{TotalReviews: 1, MasteredCount: 1, LastReviewedAt: &reviewedAt, LastResult: ptr(domainreview.ResultMastered)}
	reviews.On("SummaryByQuestion", mock.Anything, "q1", "u1").Return(summary, nil)
	questions.On("CountOccurrences", mock.Anything, "q1", "u1").Return(int64(0), nil)
	svc := questionsvc.NewInterviewQuestionService(questions, nil, reviews, nil, nil)

	review, updated, gotSummary, err := svc.CreateFlashcardReview(context.Background(), "u1", dto.CreateFlashcardReviewRequest{QuestionID: "q1", ReviewedAt: &dto.FlexibleTime{Time: reviewedAt}, Result: domainreview.ResultMastered, ReviewMode: domainreview.ModeRandom})

	require.NoError(t, err)
	assert.Equal(t, domainreview.ResultMastered, review.Result)
	assert.Equal(t, domainquestion.MasteryMastered, updated.Question.MasteryStatus)
	assert.Equal(t, int64(1), gotSummary.TotalReviews)
	reviews.AssertExpectations(t)
	questions.AssertExpectations(t)
}

func ptr[T any](v T) *T {
	return &v
}
```

- [ ] **Step 3: Run service tests and verify they fail**

Run:

```bash
cd backend && go test ./internal/application/interview_question
```

Expected: FAIL because service package does not exist.

- [ ] **Step 4: Implement service**

Create `backend/internal/application/interview_question/interview_question_service.go`:

```go
package interview_question

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
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
	defaultReviewLimit = 20
	recentReviewLimit  = 10
)

type ApplicationReader interface {
	Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error)
}

type ProgressReader interface {
	FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error)
}

type InterviewQuestionService struct {
	questions   questionrepo.Repository
	occurrences occurrencerepo.Repository
	reviews     reviewrepo.Repository
	apps        ApplicationReader
	progresses  ProgressReader
}

func NewInterviewQuestionService(questions questionrepo.Repository, occurrences occurrencerepo.Repository, reviews reviewrepo.Repository, apps ApplicationReader, progresses ProgressReader) *InterviewQuestionService {
	return &InterviewQuestionService{questions: questions, occurrences: occurrences, reviews: reviews, apps: apps, progresses: progresses}
}

func (s *InterviewQuestionService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	question, err := domainquestion.NewInterviewQuestion(domainquestion.NewInterviewQuestionParams{ID: uuid.NewString(), UserID: userID, Title: req.Title, Content: req.Content, AnswerAnalysis: req.AnswerAnalysis, Tags: req.Tags, Difficulty: req.Difficulty, QuestionType: req.QuestionType, MasteryStatus: req.MasteryStatus, IsFavorite: req.IsFavorite, Source: req.Source, Remark: req.Remark})
	if err != nil {
		return nil, err
	}
	if err := s.questions.Create(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *InterviewQuestionService) ListQuestions(ctx context.Context, userID string, query dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	if query.Difficulty != "" && !query.Difficulty.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}
	if query.QuestionType != "" && !query.QuestionType.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}
	if query.MasteryStatus != "" && !query.MasteryStatus.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}
	page, pageSize := normalizePage(query.Page, query.PageSize)
	return s.questions.ListByUser(ctx, questionrepo.ListFilter{UserID: userID, Keyword: query.Keyword, Tags: splitCSV(query.Tags), Difficulty: query.Difficulty, QuestionType: query.QuestionType, MasteryStatus: query.MasteryStatus, IsFavorite: query.IsFavorite, Page: page, PageSize: pageSize})
}

func (s *InterviewQuestionService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
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
	recentReviews, err := s.reviews.ListRecentByQuestion(ctx, id, userID, recentReviewLimit)
	if err != nil {
		return nil, 0, nil, domainreview.ReviewSummary{}, nil, err
	}
	return question, count, occurrences, summary, recentReviews, nil
}

func (s *InterviewQuestionService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	question, err := s.questions.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := question.ApplyUpdate(domainquestion.UpdateSpec{Title: req.Title, Content: req.Content, AnswerAnalysis: req.AnswerAnalysis, Tags: req.Tags, Difficulty: req.Difficulty, QuestionType: req.QuestionType, MasteryStatus: req.MasteryStatus, IsFavorite: req.IsFavorite, Source: req.Source, Remark: req.Remark}); err != nil {
		return nil, err
	}
	if err := s.questions.Update(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *InterviewQuestionService) DeleteQuestion(ctx context.Context, userID, id string) error {
	if err := s.questions.Delete(ctx, id, userID); err != nil {
		return err
	}
	if s.occurrences != nil {
		if err := s.occurrences.DeleteByQuestionID(ctx, id, userID); err != nil {
			return err
		}
	}
	if s.reviews != nil {
		return s.reviews.DeleteByQuestionID(ctx, id, userID)
	}
	return nil
}

func (s *InterviewQuestionService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	return s.UpdateQuestion(ctx, userID, id, dto.UpdateInterviewQuestionRequest{IsFavorite: &req.IsFavorite})
}

func (s *InterviewQuestionService) ListOccurrences(ctx context.Context, userID string, query dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	return s.occurrences.ListByUser(ctx, occurrencerepo.ListFilter{UserID: userID, QuestionID: query.QuestionID, ApplicationID: query.ApplicationID})
}

func (s *InterviewQuestionService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	if _, err := s.questions.FindByIDAndUserID(ctx, req.QuestionID, userID); err != nil {
		return nil, err
	}
	if s.apps == nil {
		return nil, domainerrors.ErrBadParams
	}
	application, err := s.apps.Get(ctx, userID, req.ApplicationID)
	if err != nil {
		return nil, err
	}
	stage := req.InterviewStageSnapshot
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
		if strings.TrimSpace(stage) == "" {
			stage = string(progress.Stage)
		}
	}
	occurrence, err := domainoccurrence.NewQuestionOccurrence(domainoccurrence.NewQuestionOccurrenceParams{ID: uuid.NewString(), UserID: userID, QuestionID: req.QuestionID, ApplicationID: req.ApplicationID, InterviewProgressID: req.InterviewProgressID, CompanyNameSnapshot: application.CompanyName, JobTitleSnapshot: application.JobTitle, InterviewStageSnapshot: stage, OccurredAt: req.OccurredAt.Value(), ActualQuestion: req.ActualQuestion, AnswerPerformance: req.AnswerPerformance, Note: req.Note})
	if err != nil {
		return nil, err
	}
	if err := s.occurrences.Create(ctx, occurrence); err != nil {
		return nil, err
	}
	return occurrence, nil
}

func (s *InterviewQuestionService) DeleteOccurrence(ctx context.Context, userID, id string) error {
	return s.occurrences.Delete(ctx, id, userID)
}

func (s *InterviewQuestionService) ListFlashcards(ctx context.Context, userID string, query dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	if !query.Mode.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}
	if query.Difficulty != "" && !query.Difficulty.Valid() {
		return nil, 0, domainerrors.ErrBadParams
	}
	limit := query.Limit
	if limit <= 0 || limit > defaultReviewLimit {
		limit = defaultReviewLimit
	}
	return s.questions.ListFlashcards(ctx, questionrepo.FlashcardFilter{UserID: userID, Mode: query.Mode, Tags: splitCSV(query.Tags), Difficulty: query.Difficulty, Limit: limit})
}

func (s *InterviewQuestionService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	question, err := s.questions.FindByIDAndUserID(ctx, req.QuestionID, userID)
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	reviewedAt := time.Now()
	if req.ReviewedAt != nil {
		reviewedAt = req.ReviewedAt.Value()
	}
	review, err := domainreview.NewFlashcardReview(domainreview.NewFlashcardReviewParams{ID: uuid.NewString(), UserID: userID, QuestionID: req.QuestionID, ReviewedAt: reviewedAt, Result: req.Result, ReviewMode: req.ReviewMode})
	if err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	if err := reviewToQuestion(question, review); err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	if err := s.reviews.Create(ctx, review); err != nil {
		return nil, questionrepo.QuestionWithOccurrenceCount{}, domainreview.ReviewSummary{}, err
	}
	if err := s.questions.Update(ctx, question); err != nil {
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
	return review, questionrepo.QuestionWithOccurrenceCount{Question: *question, OccurrenceCount: count}, summary, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = defaultPageSize
	}
	return page, pageSize
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func reviewToQuestion(question *domainquestion.InterviewQuestion, review *domainreview.FlashcardReview) error {
	return question.ApplyReviewResult(domainquestion.ReviewResult(review.Result), review.ReviewedAt)
}
```

- [ ] **Step 5: Run service tests and fix compile errors**

Run:

```bash
cd backend && go test ./internal/application/interview_question
```

Expected: PASS after resolving import formatting and type alignment issues caused by exact project state.

- [ ] **Step 6: Commit service**

Run:

```bash
git add backend/internal/application/interview_question
git commit -m "feat: add interview question application service"
```

---

## Task 5: Implement PostgreSQL repositories

**Files:**
- Create: `backend/internal/infrastructure/persistence/postgres/interview_question_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/question_occurrence_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/flashcard_review_repo.go`
- Test: existing package tests or add focused tests if repository test helpers are available.

- [ ] **Step 1: Add compile-time repository tests**

Create `backend/internal/infrastructure/persistence/postgres/interview_question_repo_test.go`:

```go
package postgres

import (
	"testing"

	reviewrepo "github.com/yilin/ai-for-backend/internal/domain/flashcard_review/repository"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	occurrencerepo "github.com/yilin/ai-for-backend/internal/domain/question_occurrence/repository"
)

func TestInterviewQuestionRepositories_ImplementContracts(t *testing.T) {
	var _ questionrepo.Repository = (*InterviewQuestionRepo)(nil)
	var _ occurrencerepo.Repository = (*QuestionOccurrenceRepo)(nil)
	var _ reviewrepo.Repository = (*FlashcardReviewRepo)(nil)
}
```

- [ ] **Step 2: Run repository tests and verify they fail**

Run:

```bash
cd backend && go test ./internal/infrastructure/persistence/postgres
```

Expected: FAIL because repository structs do not exist.

- [ ] **Step 3: Implement question repository**

Create `backend/internal/infrastructure/persistence/postgres/interview_question_repo.go`:

```go
package postgres

import (
	"context"
	"errors"
	"fmt"
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
	err := r.db.WithContext(ctx).Create(q).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *InterviewQuestionRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainquestion.InterviewQuestion, error) {
	var q domainquestion.InterviewQuestion
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &q, err
}

func (r *InterviewQuestionRepo) Update(ctx context.Context, q *domainquestion.InterviewQuestion) error {
	result := r.db.WithContext(ctx).Model(&domainquestion.InterviewQuestion{}).Where("id = ? AND user_id = ?", q.ID, q.UserID).Updates(map[string]any{"title": q.Title, "content": q.Content, "answer_analysis": q.AnswerAnalysis, "tags": q.Tags, "difficulty": q.Difficulty, "question_type": q.QuestionType, "mastery_status": q.MasteryStatus, "is_favorite": q.IsFavorite, "source": q.Source, "last_reviewed_at": q.LastReviewedAt, "remark": q.Remark})
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
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainquestion.InterviewQuestion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *InterviewQuestionRepo) ListByUser(ctx context.Context, f questionrepo.ListFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainquestion.InterviewQuestion{}).Where("interview_questions.user_id = ?", f.UserID)
	q = applyInterviewQuestionFilters(q, f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var questions []domainquestion.InterviewQuestion
	if err := q.Order("updated_at DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&questions).Error; err != nil {
		return nil, 0, err
	}
	rows, err := r.attachOccurrenceCounts(ctx, questions, f.UserID)
	return rows, total, err
}

func (r *InterviewQuestionRepo) CountOccurrences(ctx context.Context, questionID, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("question_occurrences").Where("question_id = ? AND user_id = ? AND deleted_at IS NULL", questionID, userID).Count(&count).Error
	return count, err
}

func (r *InterviewQuestionRepo) ListFlashcards(ctx context.Context, f questionrepo.FlashcardFilter) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	q := r.db.WithContext(ctx).Model(&domainquestion.InterviewQuestion{}).Where("interview_questions.user_id = ?", f.UserID)
	q = applyFlashcardFilters(q, f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q = applyFlashcardOrder(q, f.Mode)
	var questions []domainquestion.InterviewQuestion
	if err := q.Limit(f.Limit).Find(&questions).Error; err != nil {
		return nil, 0, err
	}
	rows, err := r.attachOccurrenceCounts(ctx, questions, f.UserID)
	return rows, total, err
}

func applyInterviewQuestionFilters(q *gorm.DB, f questionrepo.ListFilter) *gorm.DB {
	if keyword := strings.TrimSpace(f.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title ILIKE ? OR content ILIKE ? OR answer_analysis ILIKE ? OR tags::text ILIKE ?", like, like, like, like)
	}
	if len(f.Tags) > 0 {
		q = q.Where("tags ?| array"+sqlPlaceholders(len(f.Tags)), stringAnySlice(f.Tags)...)
	}
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
	if f.Mode == domainreview.ModeTag && len(f.Tags) > 0 {
		q = q.Where("tags ?| array"+sqlPlaceholders(len(f.Tags)), stringAnySlice(f.Tags)...)
	}
	if f.Mode == domainreview.ModeDifficulty && f.Difficulty != "" {
		q = q.Where("difficulty = ?", f.Difficulty)
	}
	return q
}

func applyFlashcardOrder(q *gorm.DB, mode domainreview.ReviewMode) *gorm.DB {
	subquery := "(SELECT COUNT(*) FROM question_occurrences qo WHERE qo.question_id = interview_questions.id AND qo.deleted_at IS NULL)"
	switch mode {
	case domainreview.ModeRandom, domainreview.ModeTag:
		return q.Order("RANDOM()")
	case domainreview.ModeFavoriteFirst:
		return q.Order("is_favorite DESC").Order(subquery + " DESC")
	case domainreview.ModeWeakFirst:
		return q.Order("CASE mastery_status WHEN 'unreviewed' THEN 1 WHEN 'weak' THEN 2 WHEN 'normal' THEN 3 WHEN 'mastered' THEN 4 ELSE 5 END ASC").Order(subquery + " DESC")
	case domainreview.ModeFrequentFirst:
		return q.Order(subquery + " DESC").Order("is_favorite DESC")
	default:
		return q.Order("updated_at DESC")
	}
}

func (r *InterviewQuestionRepo) attachOccurrenceCounts(ctx context.Context, questions []domainquestion.InterviewQuestion, userID string) ([]questionrepo.QuestionWithOccurrenceCount, error) {
	rows := make([]questionrepo.QuestionWithOccurrenceCount, 0, len(questions))
	for i := range questions {
		count, err := r.CountOccurrences(ctx, questions[i].ID, userID)
		if err != nil {
			return nil, err
		}
		rows = append(rows, questionrepo.QuestionWithOccurrenceCount{Question: questions[i], OccurrenceCount: count})
	}
	return rows, nil
}

func sqlPlaceholders(n int) string {
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, "?")
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func stringAnySlice(values []string) []any {
	args := make([]any, 0, len(values))
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

func ensureFmtUsed() string {
	return fmt.Sprint("")
}
```

After writing this, immediately remove the `fmt` import and `ensureFmtUsed` function if the compiler reports `fmt` is unused. The final file should not keep dead code.

- [ ] **Step 4: Implement occurrence repository**

Create `backend/internal/infrastructure/persistence/postgres/question_occurrence_repo.go`:

```go
package postgres

import (
	"context"
	"errors"
	"strings"

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

func (r *QuestionOccurrenceRepo) Create(ctx context.Context, occurrence *domainoccurrence.QuestionOccurrence) error {
	err := r.db.WithContext(ctx).Create(occurrence).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *QuestionOccurrenceRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainoccurrence.QuestionOccurrence, error) {
	var occurrence domainoccurrence.QuestionOccurrence
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&occurrence).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &occurrence, err
}

func (r *QuestionOccurrenceRepo) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domainoccurrence.QuestionOccurrence{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *QuestionOccurrenceRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return r.db.WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Delete(&domainoccurrence.QuestionOccurrence{}).Error
}

func (r *QuestionOccurrenceRepo) ListByUser(ctx context.Context, f occurrencerepo.ListFilter) ([]domainoccurrence.QuestionOccurrence, error) {
	q := r.db.WithContext(ctx).Model(&domainoccurrence.QuestionOccurrence{}).Where("user_id = ?", f.UserID)
	if questionID := strings.TrimSpace(f.QuestionID); questionID != "" {
		q = q.Where("question_id = ?", questionID)
	}
	if applicationID := strings.TrimSpace(f.ApplicationID); applicationID != "" {
		q = q.Where("application_id = ?", applicationID)
	}
	var rows []domainoccurrence.QuestionOccurrence
	err := q.Order("occurred_at DESC").Find(&rows).Error
	return rows, err
}
```

- [ ] **Step 5: Implement review repository**

Create `backend/internal/infrastructure/persistence/postgres/flashcard_review_repo.go`:

```go
package postgres

import (
	"context"

	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	reviewrepo "github.com/yilin/ai-for-backend/internal/domain/flashcard_review/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
)

type FlashcardReviewRepo struct {
	db *gorm.DB
}

var _ reviewrepo.Repository = (*FlashcardReviewRepo)(nil)

func NewFlashcardReviewRepo(db *gorm.DB) *FlashcardReviewRepo {
	return &FlashcardReviewRepo{db: db}
}

func (r *FlashcardReviewRepo) Create(ctx context.Context, review *domainreview.FlashcardReview) error {
	err := r.db.WithContext(ctx).Create(review).Error
	if err != nil && isUniqueViolation(err) {
		return domainerrors.ErrDuplicate
	}
	return err
}

func (r *FlashcardReviewRepo) DeleteByQuestionID(ctx context.Context, questionID, userID string) error {
	return r.db.WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Delete(&domainreview.FlashcardReview{}).Error
}

func (r *FlashcardReviewRepo) ListRecentByQuestion(ctx context.Context, questionID, userID string, limit int) ([]domainreview.FlashcardReview, error) {
	var rows []domainreview.FlashcardReview
	err := r.db.WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Order("reviewed_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *FlashcardReviewRepo) SummaryByQuestion(ctx context.Context, questionID, userID string) (domainreview.ReviewSummary, error) {
	var reviews []domainreview.FlashcardReview
	if err := r.db.WithContext(ctx).Where("question_id = ? AND user_id = ?", questionID, userID).Order("reviewed_at DESC").Find(&reviews).Error; err != nil {
		return domainreview.ReviewSummary{}, err
	}
	summary := domainreview.ReviewSummary{}
	for i := range reviews {
		summary.TotalReviews++
		switch reviews[i].Result {
		case domainreview.ResultFailed:
			summary.FailedCount++
		case domainreview.ResultVague:
			summary.VagueCount++
		case domainreview.ResultMastered:
			summary.MasteredCount++
		}
		if i == 0 {
			summary.LastReviewedAt = &reviews[i].ReviewedAt
			summary.LastResult = &reviews[i].Result
		}
	}
	return summary, nil
}
```

- [ ] **Step 6: Run repository package tests**

Run:

```bash
cd backend && go test ./internal/infrastructure/persistence/postgres
```

Expected: PASS after correcting SQL placeholder syntax if PostgreSQL JSONB tag filter compilation tests expose issues.

- [ ] **Step 7: Commit repositories**

Run:

```bash
git add backend/internal/infrastructure/persistence/postgres/interview_question_repo.go backend/internal/infrastructure/persistence/postgres/question_occurrence_repo.go backend/internal/infrastructure/persistence/postgres/flashcard_review_repo.go backend/internal/infrastructure/persistence/postgres/interview_question_repo_test.go
git commit -m "feat: add interview question postgres repositories"
```

---

## Task 6: Add HTTP handler and route tests

**Files:**
- Create: `backend/internal/interface/handler/interview_question_handler.go`
- Create: `backend/internal/interface/handler/interview_question_handler_test.go`
- Modify: `backend/internal/interface/router/router.go`
- Modify: `backend/internal/interface/router/router_test.go`

- [ ] **Step 1: Write failing handler tests**

Create `backend/internal/interface/handler/interview_question_handler_test.go`:

```go
package handler_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

type mockInterviewQuestionHandlerService struct{ mock.Mock }

func (m *mockInterviewQuestionHandlerService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, req)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Error(1)
}
func (m *mockInterviewQuestionHandlerService) ListQuestions(ctx context.Context, userID string, query dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, userID, query)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}
func (m *mockInterviewQuestionHandlerService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
	args := m.Called(ctx, userID, id)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Get(1).(int64), args.Get(2).([]domainoccurrence.QuestionOccurrence), args.Get(3).(domainreview.ReviewSummary), args.Get(4).([]domainreview.FlashcardReview), args.Error(5)
}
func (m *mockInterviewQuestionHandlerService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, id, req)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Error(1)
}
func (m *mockInterviewQuestionHandlerService) DeleteQuestion(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}
func (m *mockInterviewQuestionHandlerService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	args := m.Called(ctx, userID, id, req)
	q, _ := args.Get(0).(*domainquestion.InterviewQuestion)
	return q, args.Error(1)
}
func (m *mockInterviewQuestionHandlerService) ListOccurrences(ctx context.Context, userID string, query dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, userID, query)
	return args.Get(0).([]domainoccurrence.QuestionOccurrence), args.Error(1)
}
func (m *mockInterviewQuestionHandlerService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	args := m.Called(ctx, userID, req)
	occurrence, _ := args.Get(0).(*domainoccurrence.QuestionOccurrence)
	return occurrence, args.Error(1)
}
func (m *mockInterviewQuestionHandlerService) DeleteOccurrence(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}
func (m *mockInterviewQuestionHandlerService) ListFlashcards(ctx context.Context, userID string, query dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	args := m.Called(ctx, userID, query)
	return args.Get(0).([]questionrepo.QuestionWithOccurrenceCount), args.Get(1).(int64), args.Error(2)
}
func (m *mockInterviewQuestionHandlerService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	args := m.Called(ctx, userID, req)
	review, _ := args.Get(0).(*domainreview.FlashcardReview)
	return review, args.Get(1).(questionrepo.QuestionWithOccurrenceCount), args.Get(2).(domainreview.ReviewSummary), args.Error(3)
}

func TestInterviewQuestionHandler_ListCreateDetailUseFrontendShape(t *testing.T) {
	svc := &mockInterviewQuestionHandlerService{}
	now := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	question := domainquestion.InterviewQuestion{ID: "q1", UserID: "u1", Title: "Vue", Content: "content", AnswerAnalysis: "analysis", Tags: mustTags(t, []string{"Vue"}), Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryUnreviewed, CreatedAt: now, UpdatedAt: now}
	svc.On("ListQuestions", mock.Anything, "u1", mock.MatchedBy(func(q dto.InterviewQuestionQuery) bool { return q.Keyword == "Vue" && q.Page == 1 })).Return([]questionrepo.QuestionWithOccurrenceCount{{Question: question, OccurrenceCount: 1}}, int64(1), nil).Once()
	svc.On("CreateQuestion", mock.Anything, "u1", mock.MatchedBy(func(req dto.CreateInterviewQuestionRequest) bool { return req.Title == "Vue" && req.QuestionType == domainquestion.TypeKnowledge })).Return(&question, nil).Once()
	svc.On("GetQuestion", mock.Anything, "u1", "q1").Return(&question, int64(1), []domainoccurrence.QuestionOccurrence{}, domainreview.ReviewSummary{}, []domainreview.FlashcardReview{}, nil).Once()
	r := gin.New()
	h := handler.NewInterviewQuestionHandler(svc)
	r.GET("/interview-questions", func(c *gin.Context) { c.Set("userID", "u1"); h.ListQuestions(c) })
	r.POST("/interview-questions", func(c *gin.Context) { c.Set("userID", "u1"); h.CreateQuestion(c) })
	r.GET("/interview-questions/:id", func(c *gin.Context) { c.Set("userID", "u1"); h.GetQuestion(c) })

	w := performRequest(r, "GET", "/interview-questions?keyword=Vue&page=1&pageSize=10", nil)
	data := responseData(t, w)
	first := data["list"].([]any)[0].(map[string]any)
	assert.Equal(t, "analysis", first["answerAnalysis"])
	assert.Equal(t, "knowledge", first["questionType"])
	assert.NotContains(t, first, "AnswerAnalysis")

	w = performRequest(r, "POST", "/interview-questions", bytes.NewBufferString(`{"title":"Vue","content":"content","answerAnalysis":"analysis","difficulty":"medium","questionType":"knowledge","masteryStatus":"unreviewed"}`))
	data = responseData(t, w)
	assert.Equal(t, "Vue", data["title"])
	assert.Equal(t, "unreviewed", data["masteryStatus"])

	w = performRequest(r, "GET", "/interview-questions/q1", nil)
	data = responseData(t, w)
	assert.Contains(t, data, "occurrences")
	assert.Contains(t, data, "reviewSummary")
	assert.Contains(t, data, "recentReviews")
	svc.AssertExpectations(t)
}
```

This test uses `performRequest`, `responseData`, and `mustTags`. If `performRequest` is absent in the package, add this helper to the same file:

```go
func performRequest(r *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	r.ServeHTTP(w, req)
	return w
}
```

Add imports for `io`, `net/http/httptest`, and `gorm.io/datatypes` if those helpers are added.

- [ ] **Step 2: Run handler tests and verify they fail**

Run:

```bash
cd backend && go test ./internal/interface/handler
```

Expected: FAIL because handler does not exist.

- [ ] **Step 3: Implement handler**

Create `backend/internal/interface/handler/interview_question_handler.go`:

```go
package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
	questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
	domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
	pkgerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type InterviewQuestionService interface {
	CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error)
	ListQuestions(ctx context.Context, userID string, query dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error)
	GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error)
	UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error)
	DeleteQuestion(ctx context.Context, userID, id string) error
	UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error)
	ListOccurrences(ctx context.Context, userID string, query dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error)
	CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error)
	DeleteOccurrence(ctx context.Context, userID, id string) error
	ListFlashcards(ctx context.Context, userID string, query dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error)
	CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error)
}

type InterviewQuestionHandler struct {
	svc InterviewQuestionService
}

func NewInterviewQuestionHandler(svc InterviewQuestionService) *InterviewQuestionHandler {
	return &InterviewQuestionHandler{svc: svc}
}

func (h *InterviewQuestionHandler) ListQuestions(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	var query dto.InterviewQuestionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	page, pageSize := frontendPage(c)
	query.Page = page
	query.PageSize = pageSize
	rows, total, err := h.svc.ListQuestions(c.Request.Context(), userID, query)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toInterviewQuestionResponses(rows), "total": total, "page": page, "pageSize": pageSize})
}

func (h *InterviewQuestionHandler) GetQuestion(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, count, occurrences, summary, reviews, err := h.svc.GetQuestion(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionDetailResponse(question, count, occurrences, summary, reviews))
}

func (h *InterviewQuestionHandler) CreateQuestion(c *gin.Context) {
	var req dto.CreateInterviewQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, err := h.svc.CreateQuestion(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, 0))
}

func (h *InterviewQuestionHandler) UpdateQuestion(c *gin.Context) {
	var req dto.UpdateInterviewQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, err := h.svc.UpdateQuestion(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		mapError(c, err)
		return
	}
	count, err := h.svcCount(c, userID, question.ID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, count))
}

func (h *InterviewQuestionHandler) DeleteQuestion(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteQuestion(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *InterviewQuestionHandler) UpdateFavorite(c *gin.Context) {
	var req dto.UpdateQuestionFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	question, err := h.svc.UpdateFavorite(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		mapError(c, err)
		return
	}
	count, err := h.svcCount(c, userID, question.ID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToInterviewQuestionResponse(question, count))
}

func (h *InterviewQuestionHandler) ListOccurrences(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	rows, err := h.svc.ListOccurrences(c.Request.Context(), userID, dto.QuestionOccurrenceQuery{QuestionID: c.Query("questionId"), ApplicationID: c.Query("applicationId")})
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": dto.ToQuestionOccurrenceResponses(rows)})
}

func (h *InterviewQuestionHandler) CreateOccurrence(c *gin.Context) {
	var req dto.CreateQuestionOccurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	occurrence, err := h.svc.CreateOccurrence(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, dto.ToQuestionOccurrenceResponse(occurrence))
}

func (h *InterviewQuestionHandler) DeleteOccurrence(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteOccurrence(c.Request.Context(), userID, c.Param("id")); err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *InterviewQuestionHandler) ListFlashcards(c *gin.Context) {
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	var query dto.FlashcardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	rows, total, err := h.svc.ListFlashcards(c.Request.Context(), userID, query)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"list": toInterviewQuestionResponses(rows), "total": total, "mode": query.Mode})
}

func (h *InterviewQuestionHandler) CreateFlashcardReview(c *gin.Context) {
	var req dto.CreateFlashcardReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, pkgerrors.CodeBadParams, "参数错误")
		return
	}
	userID, ok := resumeUserID(c)
	if !ok {
		return
	}
	review, question, summary, err := h.svc.CreateFlashcardReview(c.Request.Context(), userID, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Success(c, gin.H{"review": dto.ToFlashcardReviewResponse(review), "updatedQuestion": dto.ToInterviewQuestionResponse(&question.Question, question.OccurrenceCount), "reviewSummary": dto.ToFlashcardReviewSummaryResponse(summary)})
}

func toInterviewQuestionResponses(rows []questionrepo.QuestionWithOccurrenceCount) []*dto.InterviewQuestionResponse {
	responses := make([]*dto.InterviewQuestionResponse, 0, len(rows))
	for i := range rows {
		responses = append(responses, dto.ToInterviewQuestionResponse(&rows[i].Question, rows[i].OccurrenceCount))
	}
	return responses
}

func (h *InterviewQuestionHandler) svcCount(c *gin.Context, userID, questionID string) (int64, error) {
	question, count, _, _, _, err := h.svc.GetQuestion(c.Request.Context(), userID, questionID)
	if err != nil {
		return 0, err
	}
	if question == nil {
		return 0, pkgerrors.ErrNotFound
	}
	return count, nil
}
```

If this creates unnecessary detail queries after update/favorite, replace `svcCount` by changing the service interface to return `questionrepo.QuestionWithOccurrenceCount` from `UpdateQuestion` and `UpdateFavorite`. Keep the final handler response identical.

- [ ] **Step 4: Add routes**

Modify `backend/internal/interface/router/router.go`:

Add field to `Dependencies`:

```go
InterviewQuestionHandler *handler.InterviewQuestionHandler
```

Add routes inside the authenticated group, before routes with generic `/:id` patterns when possible:

```go
auth.GET("/interview-questions", deps.InterviewQuestionHandler.ListQuestions)
auth.POST("/interview-questions", deps.InterviewQuestionHandler.CreateQuestion)
auth.GET("/interview-questions/:id", deps.InterviewQuestionHandler.GetQuestion)
auth.PUT("/interview-questions/:id", deps.InterviewQuestionHandler.UpdateQuestion)
auth.DELETE("/interview-questions/:id", deps.InterviewQuestionHandler.DeleteQuestion)
auth.PUT("/interview-questions/:id/favorite", deps.InterviewQuestionHandler.UpdateFavorite)

auth.GET("/question-occurrences", deps.InterviewQuestionHandler.ListOccurrences)
auth.POST("/question-occurrences", deps.InterviewQuestionHandler.CreateOccurrence)
auth.DELETE("/question-occurrences/:id", deps.InterviewQuestionHandler.DeleteOccurrence)

auth.GET("/flashcards", deps.InterviewQuestionHandler.ListFlashcards)
auth.POST("/flashcard-reviews", deps.InterviewQuestionHandler.CreateFlashcardReview)
```

- [ ] **Step 5: Update router test stubs**

In `backend/internal/interface/router/router_test.go`, add imports for new domains and repository package if not already present:

```go
domainreview "github.com/yilin/ai-for-backend/internal/domain/flashcard_review"
domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
questionrepo "github.com/yilin/ai-for-backend/internal/domain/interview_question/repository"
domainoccurrence "github.com/yilin/ai-for-backend/internal/domain/question_occurrence"
```

Add stub:

```go
type stubInterviewQuestionService struct{}

func (s *stubInterviewQuestionService) CreateQuestion(ctx context.Context, userID string, req dto.CreateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	return &domainquestion.InterviewQuestion{ID: "q1", UserID: userID, Title: req.Title, Content: req.Content, AnswerAnalysis: req.AnswerAnalysis, Difficulty: req.Difficulty, QuestionType: req.QuestionType, MasteryStatus: req.MasteryStatus}, nil
}
func (s *stubInterviewQuestionService) ListQuestions(ctx context.Context, userID string, query dto.InterviewQuestionQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	return []questionrepo.QuestionWithOccurrenceCount{}, 0, nil
}
func (s *stubInterviewQuestionService) GetQuestion(ctx context.Context, userID, id string) (*domainquestion.InterviewQuestion, int64, []domainoccurrence.QuestionOccurrence, domainreview.ReviewSummary, []domainreview.FlashcardReview, error) {
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "title", Content: "content", AnswerAnalysis: "analysis", Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryUnreviewed}, 0, []domainoccurrence.QuestionOccurrence{}, domainreview.ReviewSummary{}, []domainreview.FlashcardReview{}, nil
}
func (s *stubInterviewQuestionService) UpdateQuestion(ctx context.Context, userID, id string, req dto.UpdateInterviewQuestionRequest) (*domainquestion.InterviewQuestion, error) {
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "title", Content: "content", AnswerAnalysis: "analysis", Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryUnreviewed}, nil
}
func (s *stubInterviewQuestionService) DeleteQuestion(ctx context.Context, userID, id string) error { return nil }
func (s *stubInterviewQuestionService) UpdateFavorite(ctx context.Context, userID, id string, req dto.UpdateQuestionFavoriteRequest) (*domainquestion.InterviewQuestion, error) {
	return &domainquestion.InterviewQuestion{ID: id, UserID: userID, Title: "title", Content: "content", AnswerAnalysis: "analysis", Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryUnreviewed, IsFavorite: req.IsFavorite}, nil
}
func (s *stubInterviewQuestionService) ListOccurrences(ctx context.Context, userID string, query dto.QuestionOccurrenceQuery) ([]domainoccurrence.QuestionOccurrence, error) {
	return []domainoccurrence.QuestionOccurrence{}, nil
}
func (s *stubInterviewQuestionService) CreateOccurrence(ctx context.Context, userID string, req dto.CreateQuestionOccurrenceRequest) (*domainoccurrence.QuestionOccurrence, error) {
	return &domainoccurrence.QuestionOccurrence{ID: "o1", UserID: userID, QuestionID: req.QuestionID, ApplicationID: req.ApplicationID, CompanyNameSnapshot: "Acme", JobTitleSnapshot: "Frontend", OccurredAt: req.OccurredAt.Value()}, nil
}
func (s *stubInterviewQuestionService) DeleteOccurrence(ctx context.Context, userID, id string) error { return nil }
func (s *stubInterviewQuestionService) ListFlashcards(ctx context.Context, userID string, query dto.FlashcardQuery) ([]questionrepo.QuestionWithOccurrenceCount, int64, error) {
	return []questionrepo.QuestionWithOccurrenceCount{}, 0, nil
}
func (s *stubInterviewQuestionService) CreateFlashcardReview(ctx context.Context, userID string, req dto.CreateFlashcardReviewRequest) (*domainreview.FlashcardReview, questionrepo.QuestionWithOccurrenceCount, domainreview.ReviewSummary, error) {
	return &domainreview.FlashcardReview{ID: "r1", UserID: userID, QuestionID: req.QuestionID, Result: req.Result, ReviewMode: req.ReviewMode, ReviewedAt: time.Now()}, questionrepo.QuestionWithOccurrenceCount{Question: domainquestion.InterviewQuestion{ID: req.QuestionID, UserID: userID, Title: "title", Content: "content", AnswerAnalysis: "analysis", Difficulty: domainquestion.DifficultyMedium, QuestionType: domainquestion.TypeKnowledge, MasteryStatus: domainquestion.MasteryMastered}}, domainreview.ReviewSummary{}, nil
}
```

Add handler to `router.NewRouter` dependencies in tests:

```go
InterviewQuestionHandler: handler.NewInterviewQuestionHandler(&stubInterviewQuestionService{}),
```

Add paths to the `cases` list:

```go
{http.MethodGet, "/api/v1/interview-questions"},
{http.MethodPost, "/api/v1/interview-questions"},
{http.MethodGet, "/api/v1/interview-questions/q1"},
{http.MethodPut, "/api/v1/interview-questions/q1"},
{http.MethodDelete, "/api/v1/interview-questions/q1"},
{http.MethodPut, "/api/v1/interview-questions/q1/favorite"},
{http.MethodGet, "/api/v1/question-occurrences"},
{http.MethodPost, "/api/v1/question-occurrences"},
{http.MethodDelete, "/api/v1/question-occurrences/o1"},
{http.MethodGet, "/api/v1/flashcards?mode=random"},
{http.MethodPost, "/api/v1/flashcard-reviews"},
```

- [ ] **Step 6: Run handler and router tests**

Run:

```bash
cd backend && go test ./internal/interface/handler ./internal/interface/router
```

Expected: PASS after aligning helper imports and handler method names.

- [ ] **Step 7: Commit handler and routes**

Run:

```bash
git add backend/internal/interface/handler/interview_question_handler.go backend/internal/interface/handler/interview_question_handler_test.go backend/internal/interface/router/router.go backend/internal/interface/router/router_test.go
git commit -m "feat: add interview question API handlers"
```

---

## Task 7: Wire server dependencies and frontend API paths

**Files:**
- Modify: `backend/cmd/server/main.go`
- Modify: `frontend/src/api/interview-questions.ts`
- Modify: `frontend/src/api/question-occurrences.ts`
- Modify: `frontend/src/api/flashcards.ts`

- [ ] **Step 1: Update server wiring**

Modify `backend/cmd/server/main.go` imports:

```go
interviewquestionapp "github.com/yilin/ai-for-backend/internal/application/interview_question"
```

Add repositories after `interviewProgressRepo := ...`:

```go
interviewQuestionRepo := infraDB.NewInterviewQuestionRepo(db)
questionOccurrenceRepo := infraDB.NewQuestionOccurrenceRepo(db)
flashcardReviewRepo := infraDB.NewFlashcardReviewRepo(db)
```

Add service after `interviewProgressSvc := ...`:

```go
interviewQuestionSvc := interviewquestionapp.NewInterviewQuestionService(interviewQuestionRepo, questionOccurrenceRepo, flashcardReviewRepo, resumeApplicationSvc, interviewProgressRepo)
```

Add handler after `interviewProgressHandler := ...`:

```go
interviewQuestionHandler := handler.NewInterviewQuestionHandler(interviewQuestionSvc)
```

Add dependency in `router.NewRouter`:

```go
InterviewQuestionHandler: interviewQuestionHandler,
```

- [ ] **Step 2: Update frontend interview question API paths**

Modify `frontend/src/api/interview-questions.ts` so every URL has `/api/v1`:

```ts
url: '/api/v1/interview-questions'
url: `/api/v1/interview-questions/${id}`
url: `/api/v1/interview-questions/${id}/favorite`
```

- [ ] **Step 3: Update frontend occurrence API paths**

Modify `frontend/src/api/question-occurrences.ts`:

```ts
url: '/api/v1/question-occurrences'
url: `/api/v1/question-occurrences/${id}`
```

- [ ] **Step 4: Update frontend flashcard API paths**

Modify `frontend/src/api/flashcards.ts`:

```ts
url: '/api/v1/flashcards'
url: '/api/v1/flashcard-reviews'
```

- [ ] **Step 5: Run backend compile tests**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 6: Run frontend lint or targeted type check**

Run:

```bash
cd frontend && npm run lint
```

Expected: PASS or only pre-existing unrelated lint failures. If there are failures in the three touched API files, fix them before committing.

- [ ] **Step 7: Commit wiring and frontend path changes**

Run:

```bash
git add backend/cmd/server/main.go frontend/src/api/interview-questions.ts frontend/src/api/question-occurrences.ts frontend/src/api/flashcards.ts
git commit -m "feat: wire interview question bank backend"
```

---

## Task 8: Final verification and cleanup

**Files:**
- Review all files touched by Tasks 1-7.

- [ ] **Step 1: Run full backend tests**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run frontend lint**

Run:

```bash
cd frontend && npm run lint
```

Expected: PASS or report pre-existing unrelated failures with exact file names.

- [ ] **Step 3: Inspect git status**

Run:

```bash
git status --short
```

Expected: only intended changes remain. Do not stage unrelated pre-existing files such as root `README.md` or old untracked plan files unless the user explicitly asks.

- [ ] **Step 4: Commit remaining fixes if any**

If Step 1 or Step 2 required fixes in files from this plan, run:

```bash
git add backend frontend/src/api/interview-questions.ts frontend/src/api/question-occurrences.ts frontend/src/api/flashcards.ts
git commit -m "fix: complete interview question bank verification"
```

Expected: a new commit only when verification fixes were needed. If no fixes were needed, skip this commit.

---

## Self-review notes

Spec coverage:

- Independent question creation is covered in Task 1 and Task 4.
- Three-table database design is covered in Task 3.
- Frontend-aligned request and response DTOs are covered in Task 2 and Task 6.
- `/api/v1` route registration is covered in Task 6.
- Occurrence ownership, snapshots, and optional progress are covered in Task 4.
- Flashcard review history and mastery updates are covered in Task 1 and Task 4.
- Frontend API path alignment is covered in Task 7.
- Verification is covered in Task 8.

Placeholder scan:

- No incomplete requirement markers are intentionally present.
- Any compile alignment called out in execution steps must result in concrete code before the associated task is committed.

Type consistency:

- Service methods in Task 4 match the handler interface in Task 6.
- Repository interfaces in Task 3 match mocked tests in Task 4 and repository structs in Task 5.
- DTO field names follow the frontend model names from the approved spec.
