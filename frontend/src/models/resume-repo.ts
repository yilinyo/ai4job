/**
 * 简历仓库模型定义
 */

// 求职类型枚举
export enum JobType {
  CAMPUS = 'campus', // 校招
  SOCIAL = 'social', // 社招
  INTERNSHIP = 'internship' // 实习
}

// 简历仓库
export interface ResumeRepo {
  id: string
  userId: string // 所属用户ID
  name: string // 仓库名称
  jobType: JobType // 求职类型: 校招/社招/实习
  targetPosition: string // 目标岗位
  description?: string // 仓库描述
  versionCount: number // 版本数量
  createdAt: string // 创建时间
  updatedAt: string // 更新时间
}

// 创建简历仓库请求
export interface CreateResumeRepoRequest {
  name: string
  jobType: JobType
  targetPosition: string
  description?: string
}

// 更新简历仓库请求
export interface UpdateResumeRepoRequest {
  name?: string
  jobType?: JobType
  targetPosition?: string
  description?: string
}

// 简历仓库列表查询参数
export interface ResumeRepoQueryParams {
  userId?: string
  jobType?: JobType
  keyword?: string // 搜索关键词(仓库名称或目标岗位)
  page?: number
  pageSize?: number
}

// 简历仓库列表响应
export interface ResumeRepoListResponse {
  list: ResumeRepo[]
  total: number
  page: number
  pageSize: number
}
