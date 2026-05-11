package dto

import (
	"time"

	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainjob "github.com/yilin/ai-for-backend/internal/domain/job_posting"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
)

type CreateJobPostingRequest struct {
	CompanyName     string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro    string           `json:"companyIntro"`
	JobTitle        string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements string           `json:"jobRequirements"`
	Base            string           `json:"base" binding:"max=120"`
	SalaryRange     string           `json:"salaryRange" binding:"max=120"`
	SourcePlatform  string           `json:"sourcePlatform" binding:"max=120"`
	SourceURL       string           `json:"sourceUrl" binding:"omitempty,url,max=500"`
	Status          domainjob.Status `json:"status"`
	Remark          string           `json:"remark"`
}

type UpdateJobPostingRequest struct {
	CompanyName     string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro    string           `json:"companyIntro"`
	JobTitle        string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements string           `json:"jobRequirements"`
	Base            string           `json:"base" binding:"max=120"`
	SalaryRange     string           `json:"salaryRange" binding:"max=120"`
	SourcePlatform  string           `json:"sourcePlatform" binding:"max=120"`
	SourceURL       string           `json:"sourceUrl" binding:"omitempty,url,max=500"`
	Status          domainjob.Status `json:"status"`
	Remark          string           `json:"remark"`
}

type ParseJobPostingRequest struct {
	Content        string `json:"content" binding:"required"`
	SourcePlatform string `json:"sourcePlatform" binding:"max=120"`
	SourceURL      string `json:"sourceUrl" binding:"omitempty,url,max=500"`
}

type JobPostingResponse struct {
	ID              string           `json:"id"`
	CompanyName     string           `json:"companyName"`
	CompanyIntro    string           `json:"companyIntro"`
	JobTitle        string           `json:"jobTitle"`
	JobRequirements string           `json:"jobRequirements"`
	Base            string           `json:"base"`
	SalaryRange     string           `json:"salaryRange"`
	SourcePlatform  string           `json:"sourcePlatform"`
	SourceURL       string           `json:"sourceUrl"`
	Status          domainjob.Status `json:"status"`
	Remark          string           `json:"remark"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

func ToJobPostingResponse(p *domainjob.JobPosting) *JobPostingResponse {
	if p == nil {
		return nil
	}
	return &JobPostingResponse{
		ID:              p.ID,
		CompanyName:     p.CompanyName,
		CompanyIntro:    p.CompanyIntro,
		JobTitle:        p.JobTitle,
		JobRequirements: p.JobRequirements,
		Base:            p.Base,
		SalaryRange:     p.SalaryRange,
		SourcePlatform:  p.SourcePlatform,
		SourceURL:       p.SourceURL,
		Status:          p.Status,
		Remark:          p.Remark,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func ToJobPostingResponses(postings []*domainjob.JobPosting) []*JobPostingResponse {
	responses := make([]*JobPostingResponse, 0, len(postings))
	for _, posting := range postings {
		responses = append(responses, ToJobPostingResponse(posting))
	}
	return responses
}

type CreateResumeApplicationRequest struct {
	JobPostingID              *string          `json:"jobPostingId,omitempty"`
	ResumeVersionID           string           `json:"resumeVersionId" binding:"required"`
	RepoID                    string           `json:"repoId" binding:"required"`
	CompanyName               string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro              string           `json:"companyIntro"`
	JobTitle                  string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements           string           `json:"jobRequirements"`
	Base                      string           `json:"base" binding:"max=120"`
	SalaryRange               string           `json:"salaryRange" binding:"max=120"`
	SourcePlatform            string           `json:"sourcePlatform" binding:"max=120"`
	SourceURL                 string           `json:"sourceUrl" binding:"omitempty,url,max=500"`
	DeliveryChannel           string           `json:"deliveryChannel" binding:"max=120"`
	AppliedAt                 time.Time        `json:"appliedAt" binding:"required"`
	CurrentStatus             domainapp.Status `json:"currentStatus"`
	InterviewSummary          string           `json:"interviewSummary"`
	InterviewNotes            string           `json:"interviewNotes"`
	ResumeMatchScore          int              `json:"resumeMatchScore" binding:"omitempty,min=1,max=5"`
	InterviewPerformanceScore int              `json:"interviewPerformanceScore" binding:"omitempty,min=1,max=5"`
	RoleInterestScore         int              `json:"roleInterestScore" binding:"omitempty,min=1,max=5"`
	OverallScore              int              `json:"overallScore" binding:"omitempty,min=1,max=10"`
	Remark                    string           `json:"remark"`
}

type UpdateResumeApplicationRequest struct {
	JobPostingID              *string          `json:"jobPostingId,omitempty"`
	ResumeVersionID           string           `json:"resumeVersionId" binding:"required"`
	RepoID                    string           `json:"repoId" binding:"required"`
	CompanyName               string           `json:"companyName" binding:"required,min=1,max=120"`
	CompanyIntro              string           `json:"companyIntro"`
	JobTitle                  string           `json:"jobTitle" binding:"required,min=1,max=120"`
	JobRequirements           string           `json:"jobRequirements"`
	Base                      string           `json:"base" binding:"max=120"`
	SalaryRange               string           `json:"salaryRange" binding:"max=120"`
	SourcePlatform            string           `json:"sourcePlatform" binding:"max=120"`
	SourceURL                 string           `json:"sourceUrl" binding:"omitempty,url,max=500"`
	DeliveryChannel           string           `json:"deliveryChannel" binding:"max=120"`
	AppliedAt                 time.Time        `json:"appliedAt" binding:"required"`
	CurrentStatus             domainapp.Status `json:"currentStatus"`
	InterviewSummary          string           `json:"interviewSummary"`
	InterviewNotes            string           `json:"interviewNotes"`
	ResumeMatchScore          int              `json:"resumeMatchScore" binding:"omitempty,min=1,max=5"`
	InterviewPerformanceScore int              `json:"interviewPerformanceScore" binding:"omitempty,min=1,max=5"`
	RoleInterestScore         int              `json:"roleInterestScore" binding:"omitempty,min=1,max=5"`
	OverallScore              int              `json:"overallScore" binding:"omitempty,min=1,max=10"`
	Remark                    string           `json:"remark"`
}

type ResumeApplicationResponse struct {
	ID                        string           `json:"id"`
	JobPostingID              *string          `json:"jobPostingId,omitempty"`
	ResumeVersionID           string           `json:"resumeVersionId"`
	RepoID                    string           `json:"repoId"`
	CompanyName               string           `json:"companyName"`
	CompanyIntro              string           `json:"companyIntro"`
	JobTitle                  string           `json:"jobTitle"`
	JobRequirements           string           `json:"jobRequirements"`
	Base                      string           `json:"base"`
	SalaryRange               string           `json:"salaryRange"`
	SourcePlatform            string           `json:"sourcePlatform"`
	SourceURL                 string           `json:"sourceUrl"`
	DeliveryChannel           string           `json:"deliveryChannel"`
	AppliedAt                 time.Time        `json:"appliedAt"`
	CurrentStatus             domainapp.Status `json:"currentStatus"`
	InterviewSummary          string           `json:"interviewSummary"`
	InterviewNotes            string           `json:"interviewNotes"`
	ResumeMatchScore          int              `json:"resumeMatchScore"`
	InterviewPerformanceScore int              `json:"interviewPerformanceScore"`
	RoleInterestScore         int              `json:"roleInterestScore"`
	OverallScore              int              `json:"overallScore"`
	Remark                    string           `json:"remark"`
	CreatedAt                 time.Time        `json:"createdAt"`
	UpdatedAt                 time.Time        `json:"updatedAt"`
}

func ToResumeApplicationResponse(a *domainapp.ResumeApplication) *ResumeApplicationResponse {
	if a == nil {
		return nil
	}
	return &ResumeApplicationResponse{
		ID:                        a.ID,
		JobPostingID:              a.JobPostingID,
		ResumeVersionID:           a.ResumeVersionID,
		RepoID:                    a.RepoID,
		CompanyName:               a.CompanyName,
		CompanyIntro:              a.CompanyIntro,
		JobTitle:                  a.JobTitle,
		JobRequirements:           a.JobRequirements,
		Base:                      a.Base,
		SalaryRange:               a.SalaryRange,
		SourcePlatform:            a.SourcePlatform,
		SourceURL:                 a.SourceURL,
		DeliveryChannel:           a.DeliveryChannel,
		AppliedAt:                 a.AppliedAt,
		CurrentStatus:             a.CurrentStatus,
		InterviewSummary:          a.InterviewSummary,
		InterviewNotes:            a.InterviewNotes,
		ResumeMatchScore:          a.ResumeMatchScore,
		InterviewPerformanceScore: a.InterviewPerformanceScore,
		RoleInterestScore:         a.RoleInterestScore,
		OverallScore:              a.OverallScore,
		Remark:                    a.Remark,
		CreatedAt:                 a.CreatedAt,
		UpdatedAt:                 a.UpdatedAt,
	}
}

func ToResumeApplicationResponses(applications []*domainapp.ResumeApplication) []*ResumeApplicationResponse {
	responses := make([]*ResumeApplicationResponse, 0, len(applications))
	for _, application := range applications {
		responses = append(responses, ToResumeApplicationResponse(application))
	}
	return responses
}

type CreateInterviewProgressRequest struct {
	ApplicationID     string                `json:"applicationId" binding:"required"`
	Stage             domainprogress.Stage  `json:"stage" binding:"required"`
	OccurredAt        time.Time             `json:"occurredAt" binding:"required"`
	Result            domainprogress.Result `json:"result" binding:"required"`
	InterviewerOrTeam string                `json:"interviewerOrTeam" binding:"max=120"`
	Note              string                `json:"note"`
}

type UpdateInterviewProgressRequest struct {
	ApplicationID     string                `json:"applicationId" binding:"required"`
	Stage             domainprogress.Stage  `json:"stage" binding:"required"`
	OccurredAt        time.Time             `json:"occurredAt" binding:"required"`
	Result            domainprogress.Result `json:"result" binding:"required"`
	InterviewerOrTeam string                `json:"interviewerOrTeam" binding:"max=120"`
	Note              string                `json:"note"`
}

type InterviewProgressResponse struct {
	ID                string                `json:"id"`
	ApplicationID     string                `json:"applicationId"`
	Stage             domainprogress.Stage  `json:"stage"`
	OccurredAt        time.Time             `json:"occurredAt"`
	Result            domainprogress.Result `json:"result"`
	InterviewerOrTeam string                `json:"interviewerOrTeam"`
	Note              string                `json:"note"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

func ToInterviewProgressResponse(p *domainprogress.InterviewProgress) *InterviewProgressResponse {
	if p == nil {
		return nil
	}
	return &InterviewProgressResponse{
		ID:                p.ID,
		ApplicationID:     p.ApplicationID,
		Stage:             p.Stage,
		OccurredAt:        p.OccurredAt,
		Result:            p.Result,
		InterviewerOrTeam: p.InterviewerOrTeam,
		Note:              p.Note,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func ToInterviewProgressResponses(progresses []*domainprogress.InterviewProgress) []*InterviewProgressResponse {
	responses := make([]*InterviewProgressResponse, 0, len(progresses))
	for _, progress := range progresses {
		responses = append(responses, ToInterviewProgressResponse(progress))
	}
	return responses
}
