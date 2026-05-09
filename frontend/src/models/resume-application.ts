/**
 * 投递记录模型定义
 */

export enum ApplicationStatus {
  APPLIED = 'applied',
  VIEWED = 'viewed',
  WRITTEN_TEST = 'written_test',
  FIRST_INTERVIEW = 'first_interview',
  SECOND_INTERVIEW = 'second_interview',
  FINAL_INTERVIEW = 'final_interview',
  HR_INTERVIEW = 'hr_interview',
  OFFER = 'offer',
  REJECTED = 'rejected',
  CLOSED = 'closed'
}

export interface ResumeApplication {
  id: string
  jobPostingId?: string
  resumeVersionId: string
  repoId: string
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  sourcePlatform?: string
  sourceUrl?: string
  deliveryChannel?: string
  appliedAt: string
  currentStatus: ApplicationStatus
  interviewSummary?: string
  interviewNotes?: string
  resumeMatchScore?: number
  interviewPerformanceScore?: number
  roleInterestScore?: number
  overallScore?: number
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface CreateResumeApplicationRequest {
  jobPostingId?: string
  resumeVersionId: string
  repoId: string
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  sourcePlatform?: string
  sourceUrl?: string
  deliveryChannel?: string
  appliedAt: string
  currentStatus: ApplicationStatus
  interviewSummary?: string
  interviewNotes?: string
  resumeMatchScore?: number
  interviewPerformanceScore?: number
  roleInterestScore?: number
  overallScore?: number
  remark?: string
}

export type UpdateResumeApplicationRequest = Partial<CreateResumeApplicationRequest>

export interface ResumeApplicationQueryParams {
  resumeVersionId?: string
  repoId?: string
  keyword?: string
  currentStatus?: ApplicationStatus
  page?: number
  pageSize?: number
}

export interface ResumeApplicationListResponse {
  list: ResumeApplication[]
  total: number
  page: number
  pageSize: number
}
