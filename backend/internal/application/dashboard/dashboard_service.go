package dashboard

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	activityrepo "github.com/yilin/ai-for-backend/internal/domain/activity_log/repository"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
	todorepo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo/repository"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	apprepo "github.com/yilin/ai-for-backend/internal/domain/resume_application/repository"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

const followUpDays = 7
const applicationPageSize = 1000

type ApplicationRepository interface {
	ListByUser(ctx context.Context, filter apprepo.ListFilter) ([]domainapp.ResumeApplication, int64, error)
}

type InterviewProgressRepository interface {
	ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error)
}

type Service struct {
	apps       ApplicationRepository
	progresses InterviewProgressRepository
	todos      todorepo.Repository
	activities activityrepo.Repository
	now        func() time.Time
}

func NewService(apps ApplicationRepository, progresses InterviewProgressRepository, todos todorepo.Repository, activities activityrepo.Repository) *Service {
	return &Service{apps: apps, progresses: progresses, todos: todos, activities: activities, now: time.Now}
}

func (s *Service) Get(ctx context.Context, userID, rawRange string) (*dto.DashboardResponse, error) {
	rangeKey, startAt, err := parseRange(rawRange, s.now())
	if err != nil {
		return nil, err
	}

	applications, _, err := s.apps.ListByUser(ctx, apprepo.ListFilter{UserID: userID, Page: 1, PageSize: applicationPageSize})
	if err != nil {
		return nil, err
	}
	filtered := filterApplications(applications, startAt)

	manualTodos, err := s.todos.ListByUser(ctx, todorepo.ListFilter{UserID: userID, TodayOnly: true})
	if err != nil {
		return nil, err
	}

	activities, err := s.activities.ListByUser(ctx, activityrepo.ListFilter{UserID: userID, StartAt: startAt, Limit: 10})
	if err != nil {
		return nil, err
	}

	allTodos := make([]dto.DashboardTodoResponse, 0, len(manualTodos)+4)
	for _, todo := range manualTodos {
		allTodos = append(allTodos, dto.ToDashboardTodoResponse(todo, domaintodo.SourceManual, routeForTodo(todo.TargetType, todo.TargetID)))
	}
	systemTodos, err := s.systemTodos(ctx, userID, applications)
	if err != nil {
		return nil, err
	}
	allTodos = append(allTodos, systemTodos...)

	activityResponses := make([]dto.ActivityLogResponse, 0, len(activities))
	for _, activity := range activities {
		activityResponses = append(activityResponses, dto.ToActivityLogResponse(activity))
	}

	return &dto.DashboardResponse{
		Range:                rangeKey,
		Metrics:              buildMetrics(filtered, s.now()),
		StatusPipeline:       buildStatusPipeline(filtered),
		ApplicationTrend:     buildTrend(filtered),
		CompanyRanking:       buildNameCounts(filtered, func(a domainapp.ResumeApplication) string { return a.CompanyName }),
		JobTitleDistribution: buildNameCounts(filtered, func(a domainapp.ResumeApplication) string { return classifyJobTitle(a.JobTitle) }),
		SourceDistribution: buildNameCounts(filtered, func(a domainapp.ResumeApplication) string {
			return firstNonEmpty(a.SourcePlatform, a.DeliveryChannel, "其他")
		}),
		Todos:      allTodos,
		Activities: activityResponses,
	}, nil
}

func (s *Service) CreateTodo(ctx context.Context, userID string, req dto.CreateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	var dueAt *time.Time
	if req.DueAt != nil {
		value := req.DueAt.Value()
		dueAt = &value
	}
	todo, err := domaintodo.NewDashboardTodo(domaintodo.NewDashboardTodoParams{ID: uuid.NewString(), UserID: userID, Title: req.Title, Description: req.Description, DueAt: dueAt, TargetType: req.TargetType, TargetID: req.TargetID})
	if err != nil {
		return nil, err
	}
	if err := s.todos.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) UpdateTodo(ctx context.Context, userID, id string, req dto.UpdateDashboardTodoRequest) (*domaintodo.DashboardTodo, error) {
	todo, err := s.todos.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	var dueAt **time.Time
	if req.ClearDueAt {
		var cleared *time.Time
		dueAt = &cleared
	} else if req.DueAt != nil {
		value := req.DueAt.Value()
		dueValue := &value
		dueAt = &dueValue
	}

	var targetID **string
	if req.ClearTarget {
		var cleared *string
		targetID = &cleared
		empty := domaintodo.TargetType("")
		req.TargetType = &empty
	} else if req.TargetID != nil {
		targetValue := req.TargetID
		targetID = &targetValue
	}

	if err := todo.ApplyUpdate(domaintodo.UpdateSpec{Title: req.Title, Description: req.Description, DueAt: dueAt, TargetType: req.TargetType, TargetID: targetID}); err != nil {
		return nil, err
	}
	if err := s.todos.Update(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) CompleteTodo(ctx context.Context, userID, id string) (*domaintodo.DashboardTodo, error) {
	todo, err := s.todos.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	todo.Complete(s.now())
	if err := s.todos.Update(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) DeleteTodo(ctx context.Context, userID, id string) error {
	return s.todos.Delete(ctx, id, userID)
}

func (s *Service) systemTodos(ctx context.Context, userID string, applications []domainapp.ResumeApplication) ([]dto.DashboardTodoResponse, error) {
	now := s.now()
	responses := make([]dto.DashboardTodoResponse, 0)
	for _, application := range applications {
		if application.UserID != userID {
			continue
		}
		progresses, err := s.progresses.ListByApplicationAndUser(ctx, application.ID, userID)
		if err != nil {
			return nil, err
		}
		if needsFollowUp(application, progresses, now) {
			responses = append(responses, systemTodoResponse(
				"system-follow-up-"+application.ID,
				fmt.Sprintf("跟进 %s / %s", application.CompanyName, application.JobTitle),
				"投递超过 7 天未进入下一阶段",
				application.ID,
				now,
			))
		}
		for _, progress := range progresses {
			if sameDay(progress.OccurredAt, now) && isInterviewStage(progress.Stage) {
				responses = append(responses, systemTodoResponse(
					"system-interview-"+progress.ID,
					fmt.Sprintf("面试 %s / %s %s", application.CompanyName, application.JobTitle, statusLabel(domainapp.Status(progress.Stage))),
					"今天有面试安排",
					application.ID,
					progress.OccurredAt,
				))
			}
		}
	}
	return responses, nil
}

func parseRange(raw string, now time.Time) (string, *time.Time, error) {
	rangeKey := strings.TrimSpace(raw)
	if rangeKey == "" {
		rangeKey = "30d"
	}
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch rangeKey {
	case "7d":
		start := startOfToday.AddDate(0, 0, -6)
		return rangeKey, &start, nil
	case "30d":
		start := startOfToday.AddDate(0, 0, -29)
		return rangeKey, &start, nil
	case "90d":
		start := startOfToday.AddDate(0, 0, -89)
		return rangeKey, &start, nil
	case "all":
		return rangeKey, nil, nil
	default:
		return "", nil, domainerrors.ErrBadParams
	}
}

func filterApplications(rows []domainapp.ResumeApplication, startAt *time.Time) []domainapp.ResumeApplication {
	filtered := make([]domainapp.ResumeApplication, 0, len(rows))
	for _, row := range rows {
		if startAt == nil || !row.AppliedAt.Before(*startAt) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func buildMetrics(rows []domainapp.ResumeApplication, now time.Time) dto.DashboardMetrics {
	metrics := dto.DashboardMetrics{TotalApplications: len(rows)}
	for _, row := range rows {
		if isInterviewing(row.CurrentStatus) {
			metrics.Interviewing++
		}
		if row.CurrentStatus == domainapp.StatusOffer {
			metrics.Offers++
		}
		if (row.CurrentStatus == domainapp.StatusApplied || row.CurrentStatus == domainapp.StatusViewed) && now.Sub(row.AppliedAt) >= followUpDays*24*time.Hour {
			metrics.PendingFollowUp++
		}
	}
	return metrics
}

func buildStatusPipeline(rows []domainapp.ResumeApplication) []dto.DashboardStatusBucket {
	counts := make(map[domainapp.Status]int)
	for _, row := range rows {
		counts[row.CurrentStatus]++
	}
	pipeline := make([]dto.DashboardStatusBucket, 0, len(statusOrder))
	for _, item := range statusOrder {
		pipeline = append(pipeline, dto.DashboardStatusBucket{Status: string(item.status), Label: item.label, Count: counts[item.status]})
	}
	return pipeline
}

func buildTrend(rows []domainapp.ResumeApplication) []dto.DashboardTrendPoint {
	counts := make(map[string]int)
	for _, row := range rows {
		date := row.AppliedAt.Format("2006-01-02")
		counts[date]++
	}
	dates := make([]string, 0, len(counts))
	for date := range counts {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	trend := make([]dto.DashboardTrendPoint, 0, len(dates))
	for _, date := range dates {
		trend = append(trend, dto.DashboardTrendPoint{Date: date, Count: counts[date]})
	}
	return trend
}

func buildNameCounts(rows []domainapp.ResumeApplication, keyFn func(domainapp.ResumeApplication) string) []dto.DashboardNameCount {
	counts := make(map[string]int)
	for _, row := range rows {
		name := firstNonEmpty(keyFn(row), "其他")
		counts[name]++
	}
	result := make([]dto.DashboardNameCount, 0, len(counts))
	for name, count := range counts {
		result = append(result, dto.DashboardNameCount{Name: name, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Name < result[j].Name
		}
		return result[i].Count > result[j].Count
	})
	if len(result) > 10 {
		return result[:10]
	}
	return result
}

func classifyJobTitle(title string) string {
	lower := strings.ToLower(title)
	switch {
	case strings.Contains(lower, "全栈") || strings.Contains(lower, "fullstack") || strings.Contains(lower, "full-stack"):
		return "全栈"
	case strings.Contains(lower, "前端") || strings.Contains(lower, "frontend") || strings.Contains(lower, "front-end") || strings.Contains(lower, "vue") || strings.Contains(lower, "react"):
		return "前端"
	case strings.Contains(lower, "后端") || strings.Contains(lower, "backend") || strings.Contains(lower, "back-end") || strings.Contains(lower, "java") || strings.Contains(lower, "go") || strings.Contains(lower, "golang"):
		return "后端"
	case strings.Contains(lower, "ai") || strings.Contains(lower, "人工智能") || strings.Contains(lower, "机器学习") || strings.Contains(lower, "算法") || strings.Contains(lower, "大模型"):
		return "AI"
	case strings.Contains(lower, "测试") || strings.Contains(lower, "qa") || strings.Contains(lower, "test"):
		return "测试"
	case strings.Contains(lower, "产品") || strings.Contains(lower, "product") || strings.Contains(lower, "pm"):
		return "产品"
	default:
		return "其他"
	}
}

func routeForTodo(target domaintodo.TargetType, id *string) string {
	if id == nil || strings.TrimSpace(*id) == "" {
		return ""
	}
	targetID := strings.TrimSpace(*id)
	switch target {
	case domaintodo.TargetTypeApplication, domaintodo.TargetTypeInterview:
		return "/resume/application/" + targetID
	case domaintodo.TargetTypeQuestion:
		return "/interview-question-bank/question/" + targetID
	case domaintodo.TargetTypeResume:
		return "/resume/repo/" + targetID
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func needsFollowUp(application domainapp.ResumeApplication, progresses []domainprogress.InterviewProgress, now time.Time) bool {
	if application.CurrentStatus != domainapp.StatusApplied && application.CurrentStatus != domainapp.StatusViewed {
		return false
	}
	if now.Sub(application.AppliedAt) < followUpDays*24*time.Hour {
		return false
	}
	for _, progress := range progresses {
		if progress.OccurredAt.After(application.AppliedAt) && isLaterStage(progress.Stage) && progress.OccurredAt.Before(now) {
			return false
		}
	}
	return true
}

func isLaterStage(stage domainprogress.Stage) bool {
	switch stage {
	case domainprogress.StageWrittenTest, domainprogress.StageFirstInterview, domainprogress.StageSecondInterview, domainprogress.StageFinalInterview, domainprogress.StageHRInterview, domainprogress.StageOffer, domainprogress.StageRejected, domainprogress.StageClosed:
		return true
	default:
		return false
	}
}

func isInterviewStage(stage domainprogress.Stage) bool {
	switch stage {
	case domainprogress.StageWrittenTest, domainprogress.StageFirstInterview, domainprogress.StageSecondInterview, domainprogress.StageFinalInterview, domainprogress.StageHRInterview:
		return true
	default:
		return false
	}
}

func isInterviewing(status domainapp.Status) bool {
	switch status {
	case domainapp.StatusWrittenTest, domainapp.StatusFirstInterview, domainapp.StatusSecondInterview, domainapp.StatusFinalInterview, domainapp.StatusHRInterview:
		return true
	default:
		return false
	}
}

func sameDay(a, b time.Time) bool {
	a = a.In(b.Location())
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func systemTodoResponse(id, title, description, applicationID string, dueAt time.Time) dto.DashboardTodoResponse {
	targetID := applicationID
	targetType := domaintodo.TargetTypeApplication
	return dto.ToDashboardTodoResponse(domaintodo.DashboardTodo{
		ID:          id,
		Title:       title,
		Description: description,
		DueAt:       &dueAt,
		TargetType:  targetType,
		TargetID:    &targetID,
		CreatedAt:   dueAt,
		UpdatedAt:   dueAt,
	}, domaintodo.SourceSystem, routeForTodo(targetType, &targetID))
}

type statusOrderItem struct {
	status domainapp.Status
	label  string
}

var statusOrder = []statusOrderItem{
	{domainapp.StatusApplied, "已投递"},
	{domainapp.StatusViewed, "已查看"},
	{domainapp.StatusWrittenTest, "笔试"},
	{domainapp.StatusFirstInterview, "一面"},
	{domainapp.StatusSecondInterview, "二面"},
	{domainapp.StatusFinalInterview, "终面"},
	{domainapp.StatusHRInterview, "HR面"},
	{domainapp.StatusOffer, "Offer"},
	{domainapp.StatusRejected, "已拒绝"},
	{domainapp.StatusClosed, "已关闭"},
}

func statusLabel(status domainapp.Status) string {
	for _, item := range statusOrder {
		if item.status == status {
			return item.label
		}
	}
	return string(status)
}
