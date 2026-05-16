package interview_question

import (
	"encoding/json"
	"strings"
	"time"

	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type QuestionType string

const (
	QuestionTypeKnowledge  QuestionType = "knowledge"
	QuestionTypeScenario   QuestionType = "scenario"
	QuestionTypeCoding     QuestionType = "coding"
	QuestionTypeAlgorithm  QuestionType = "algorithm"
	QuestionTypeProject    QuestionType = "project"
	QuestionTypeBehavior   QuestionType = "behavior"
	QuestionTypeDiscussion QuestionType = "discussion"
)

type MasteryStatus string

const (
	MasteryStatusUnreviewed MasteryStatus = "unreviewed"
	MasteryStatusWeak       MasteryStatus = "weak"
	MasteryStatusNormal     MasteryStatus = "normal"
	MasteryStatusMastered   MasteryStatus = "mastered"
)

type ReviewResult string

const (
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
	Difficulty     Difficulty     `gorm:"type:varchar(20);not null"`
	QuestionType   QuestionType   `gorm:"type:varchar(20);not null"`
	MasteryStatus  MasteryStatus  `gorm:"type:varchar(20);not null"`
	IsFavorite     bool           `gorm:"not null;default:false"`
	Source         string         `gorm:"type:varchar(100)"`
	LastReviewedAt *time.Time
	Remark         string `gorm:"type:text"`
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
	LastReviewedAt *time.Time
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
	LastReviewedAt *time.Time
	Remark         *string
}

func NewInterviewQuestion(params NewInterviewQuestionParams) (*InterviewQuestion, error) {
	masteryStatus := params.MasteryStatus
	if masteryStatus == "" {
		masteryStatus = MasteryStatusUnreviewed
	}

	question := &InterviewQuestion{
		ID:             params.ID,
		UserID:         params.UserID,
		Title:          params.Title,
		Content:        params.Content,
		AnswerAnalysis: params.AnswerAnalysis,
		Tags:           EncodeTags(params.Tags),
		Difficulty:     params.Difficulty,
		QuestionType:   params.QuestionType,
		MasteryStatus:  masteryStatus,
		IsFavorite:     params.IsFavorite,
		Source:         params.Source,
		LastReviewedAt: params.LastReviewedAt,
		Remark:         params.Remark,
	}
	if err := question.Validate(); err != nil {
		return nil, err
	}
	return question, nil
}

func EncodeTags(tags []string) datatypes.JSON {
	if tags == nil {
		tags = []string{}
	}
	encoded, err := json.Marshal(tags)
	if err != nil {
		return datatypes.JSON([]byte("[]"))
	}
	return datatypes.JSON(encoded)
}

func DecodeTags(tags datatypes.JSON) ([]string, error) {
	if len(tags) == 0 {
		return []string{}, nil
	}
	decoded := []string{}
	if err := json.Unmarshal(tags, &decoded); err != nil {
		return nil, err
	}
	if decoded == nil {
		return []string{}, nil
	}
	return decoded, nil
}

func (q *InterviewQuestion) ApplyUpdate(spec UpdateSpec) error {
	original := *q

	if spec.Title != nil {
		q.Title = *spec.Title
	}
	if spec.Content != nil {
		q.Content = *spec.Content
	}
	if spec.AnswerAnalysis != nil {
		q.AnswerAnalysis = *spec.AnswerAnalysis
	}
	if spec.Tags != nil {
		q.Tags = EncodeTags(*spec.Tags)
	}
	if spec.Difficulty != nil {
		q.Difficulty = *spec.Difficulty
	}
	if spec.QuestionType != nil {
		q.QuestionType = *spec.QuestionType
	}
	if spec.MasteryStatus != nil {
		q.MasteryStatus = *spec.MasteryStatus
	}
	if spec.IsFavorite != nil {
		q.IsFavorite = *spec.IsFavorite
	}
	if spec.Source != nil {
		q.Source = *spec.Source
	}
	if spec.LastReviewedAt != nil {
		q.LastReviewedAt = spec.LastReviewedAt
	}
	if spec.Remark != nil {
		q.Remark = *spec.Remark
	}

	if err := q.Validate(); err != nil {
		*q = original
		return err
	}
	return nil
}

func (q *InterviewQuestion) ApplyReviewResult(result ReviewResult, reviewedAt time.Time) error {
	switch result {
	case ReviewResultFailed:
		q.MasteryStatus = MasteryStatusWeak
	case ReviewResultVague:
		q.MasteryStatus = MasteryStatusNormal
	case ReviewResultMastered:
		q.MasteryStatus = MasteryStatusMastered
	default:
		return domainerrors.ErrBadParams
	}
	q.LastReviewedAt = &reviewedAt
	return q.Validate()
}

func (q *InterviewQuestion) Validate() error {
	if strings.TrimSpace(q.ID) == "" || strings.TrimSpace(q.UserID) == "" || strings.TrimSpace(q.Title) == "" || strings.TrimSpace(q.Content) == "" || strings.TrimSpace(q.AnswerAnalysis) == "" {
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
	case QuestionTypeKnowledge, QuestionTypeScenario, QuestionTypeCoding, QuestionTypeAlgorithm, QuestionTypeProject, QuestionTypeBehavior, QuestionTypeDiscussion:
		return true
	default:
		return false
	}
}

func (s MasteryStatus) Valid() bool {
	switch s {
	case MasteryStatusUnreviewed, MasteryStatusWeak, MasteryStatusNormal, MasteryStatusMastered:
		return true
	default:
		return false
	}
}

func (InterviewQuestion) TableName() string {
	return "interview_questions"
}
