/**
 * 岗位库模型定义
 */

export enum JobPostingStatus {
  PENDING = 'pending',
  APPLIED = 'applied',
  NOT_FIT = 'not_fit',
  CLOSED = 'closed'
}

export interface JobPosting {
  id: string
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  sourcePlatform?: string
  sourceUrl?: string
  status: JobPostingStatus
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface CreateJobPostingRequest {
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  sourcePlatform?: string
  sourceUrl?: string
  status: JobPostingStatus
  remark?: string
}

export type UpdateJobPostingRequest = Partial<CreateJobPostingRequest>

export interface JobPostingQueryParams {
  keyword?: string
  sourcePlatform?: string
  status?: JobPostingStatus
  page?: number
  pageSize?: number
}

export interface JobPostingListResponse {
  list: JobPosting[]
  total: number
  page: number
  pageSize: number
}

export interface ParseJobPostingRequest {
  url: string
}

export interface ParseJobPostingResponse {
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  sourcePlatform?: string
  sourceUrl?: string
}
