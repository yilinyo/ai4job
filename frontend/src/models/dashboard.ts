/**
 * 仪表盘模型定义
 */

export type DashboardRange = '7d' | '30d' | '90d' | 'all'

export type DashboardTargetType = 'application' | 'interview' | 'question' | 'resume'

export type DashboardTodoSource = 'manual' | 'system'

export interface DashboardMetrics {
  totalApplications: number
  interviewing: number
  pendingFollowUp: number
  offers: number
}

export interface DashboardStatusBucket {
  status: string
  label: string
  count: number
}

export interface DashboardTrendPoint {
  date: string
  count: number
}

export interface DashboardNameCount {
  name: string
  count: number
}

export interface DashboardTodo {
  id: string
  title: string
  description: string
  dueAt?: string
  completedAt?: string
  targetType?: DashboardTargetType
  targetId?: string
  route?: string
  source: DashboardTodoSource
  createdAt: string
  updatedAt: string
}

export interface ActivityLog {
  id: string
  action: string
  targetType: DashboardTargetType
  targetId: string
  targetLabel: string
  summary: string
  route?: string
  createdAt: string
}

export interface DashboardResponse {
  range: DashboardRange
  metrics: DashboardMetrics
  statusPipeline: DashboardStatusBucket[]
  applicationTrend: DashboardTrendPoint[]
  companyRanking: DashboardNameCount[]
  jobTitleDistribution: DashboardNameCount[]
  sourceDistribution: DashboardNameCount[]
  todos: DashboardTodo[]
  activities: ActivityLog[]
}

export interface DashboardQueryParams {
  range?: DashboardRange
}

export interface DashboardTodoQueryParams {
  completed?: boolean
}

export interface DashboardTodoListResponse {
  list: DashboardTodo[]
}

export interface CreateDashboardTodoRequest {
  title: string
  description?: string
  dueAt?: string
  targetType?: DashboardTargetType
  targetId?: string
}

export interface UpdateDashboardTodoRequest {
  title?: string
  description?: string
  dueAt?: string
  clearDueAt?: boolean
  targetType?: DashboardTargetType
  targetId?: string
  clearTarget?: boolean
}
