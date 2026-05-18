package dto

import (
	"time"

	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domaintodo "github.com/yilin/ai-for-backend/internal/domain/dashboard_todo"
)

type DashboardQuery struct {
	Range string `json:"range" form:"range"`
}

type CreateDashboardTodoRequest struct {
	Title       string                `json:"title" binding:"required,min=1,max=200"`
	Description string                `json:"description"`
	DueAt       *FlexibleTime         `json:"dueAt"`
	TargetType  domaintodo.TargetType `json:"targetType"`
	TargetID    *string               `json:"targetId"`
}

type UpdateDashboardTodoRequest struct {
	Title       *string                `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string                `json:"description"`
	DueAt       *FlexibleTime          `json:"dueAt"`
	ClearDueAt  bool                   `json:"clearDueAt"`
	TargetType  *domaintodo.TargetType `json:"targetType"`
	TargetID    *string                `json:"targetId"`
	ClearTarget bool                   `json:"clearTarget"`
}

type DashboardResponse struct {
	Range                string                  `json:"range"`
	Metrics              DashboardMetrics        `json:"metrics"`
	StatusPipeline       []DashboardStatusBucket `json:"statusPipeline"`
	ApplicationTrend     []DashboardTrendPoint   `json:"applicationTrend"`
	CompanyRanking       []DashboardNameCount    `json:"companyRanking"`
	JobTitleDistribution []DashboardNameCount    `json:"jobTitleDistribution"`
	SourceDistribution   []DashboardNameCount    `json:"sourceDistribution"`
	Todos                []DashboardTodoResponse `json:"todos"`
	Activities           []ActivityLogResponse   `json:"activities"`
}

type DashboardMetrics struct {
	TotalApplications int `json:"totalApplications"`
	Interviewing      int `json:"interviewing"`
	PendingFollowUp   int `json:"pendingFollowUp"`
	Offers            int `json:"offers"`
}

type DashboardStatusBucket struct {
	Status string `json:"status"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
}

type DashboardTrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DashboardNameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DashboardTodoResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueAt       *time.Time `json:"dueAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	TargetType  string     `json:"targetType,omitempty"`
	TargetID    *string    `json:"targetId,omitempty"`
	Route       string     `json:"route,omitempty"`
	Source      string     `json:"source"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ActivityLogResponse struct {
	ID          string    `json:"id"`
	Action      string    `json:"action"`
	TargetType  string    `json:"targetType"`
	TargetID    string    `json:"targetId"`
	TargetLabel string    `json:"targetLabel"`
	Summary     string    `json:"summary"`
	Route       string    `json:"route,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func ToDashboardTodoResponse(todo domaintodo.DashboardTodo, source domaintodo.Source, route string) DashboardTodoResponse {
	return DashboardTodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		DueAt:       todo.DueAt,
		CompletedAt: todo.CompletedAt,
		TargetType:  string(todo.TargetType),
		TargetID:    todo.TargetID,
		Route:       route,
		Source:      string(source),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func ToActivityLogResponse(activity domainactivity.ActivityLog) ActivityLogResponse {
	return ActivityLogResponse{
		ID:          activity.ID,
		Action:      string(activity.Action),
		TargetType:  string(activity.TargetType),
		TargetID:    activity.TargetID,
		TargetLabel: activity.TargetLabel,
		Summary:     activity.Summary,
		Route:       activity.Route,
		CreatedAt:   activity.CreatedAt,
	}
}
