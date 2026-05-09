/**
 * 模型统一导出
 */

export * from './user'
export * from './resume-repo'
export * from './resume-version'
export * from './resume-application'
export * from './interview-progress'
export * from './job-posting'
export * from './interview-question'

// 通用响应结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页参数
export interface PaginationParams {
  page: number
  pageSize: number
}

// 分页响应
export interface PaginationResponse<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
