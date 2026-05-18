import { Request, Response } from 'express'
import { getApplicationByIdInternal } from './resume-applications'
import { getInterviewProgressByIdInternal } from './interview-progress'
import { getInterviewQuestionByIdInternal } from './interview-questions'
import { getResumeRepoByIdInternal } from './resume-repos'

type DashboardRange = '7d' | '30d' | '90d' | 'all'
type DashboardTargetType = 'application' | 'interview' | 'question' | 'resume'
type DashboardTodoSource = 'manual' | 'system'

interface DashboardMetrics {
  totalApplications: number
  interviewing: number
  pendingFollowUp: number
  offers: number
}

interface DashboardStatusBucket {
  status: string
  label: string
  count: number
}

interface DashboardTrendPoint {
  date: string
  count: number
}

interface DashboardNameCount {
  name: string
  count: number
}

interface DashboardTodo {
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

interface ActivityLog {
  id: string
  action: string
  targetType: DashboardTargetType
  targetId: string
  targetLabel: string
  summary: string
  route?: string
  createdAt: string
}

interface DashboardResponse {
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

interface UpdateDashboardTodoRequest {
  title?: string
  description?: string
  dueAt?: string
  clearDueAt?: boolean
  targetType?: DashboardTargetType
  targetId?: string
  clearTarget?: boolean
}

const success = (res: Response, data: any) => res.json({
  code: 20000,
  message: 'success',
  data
})

const badRequest = (res: Response, message = '参数错误') => res.status(400).json({
  code: 50001,
  message,
  data: null
})

const notFound = (res: Response) => res.status(404).json({
  code: 50004,
  message: '待办事项不存在',
  data: null
})

const now = new Date('2026-05-18T09:30:00.000Z')
const iso = (value: string) => new Date(value).toISOString()
const daysAgo = (days: number) => new Date(now.getTime() - days * 24 * 60 * 60 * 1000).toISOString()

let nextTodoId = 4

const manualTodos: DashboardTodo[] = [
  {
    id: '1',
    title: '跟进腾讯 HR 面结果',
    description: '确认 HR 面反馈和下一步时间安排。',
    dueAt: daysAgo(-1),
    targetType: 'application',
    targetId: '1',
    route: '/resume/application/1',
    source: 'manual',
    createdAt: iso('2026-05-14T10:00:00.000Z'),
    updatedAt: iso('2026-05-14T10:00:00.000Z')
  },
  {
    id: '2',
    title: '补充全栈项目 STAR 案例',
    description: '整理 SaaS 项目中的复杂问题拆解和业务结果。',
    dueAt: iso('2026-05-20T12:00:00.000Z'),
    targetType: 'resume',
    targetId: '2',
    route: '/resume/repo/2',
    source: 'manual',
    createdAt: iso('2026-05-15T08:30:00.000Z'),
    updatedAt: iso('2026-05-15T08:30:00.000Z')
  },
  {
    id: '3',
    title: '复盘字节一面算法题',
    description: '把未完成的双指针题加入面试题库复习。',
    dueAt: iso('2026-05-10T12:00:00.000Z'),
    completedAt: iso('2026-05-12T16:00:00.000Z'),
    targetType: 'question',
    targetId: '1',
    route: '/interview-question-bank/question/1',
    source: 'manual',
    createdAt: iso('2026-05-09T09:00:00.000Z'),
    updatedAt: iso('2026-05-12T16:00:00.000Z')
  }
]

const systemTodos: DashboardTodo[] = [
  {
    id: 'system-follow-up-1',
    title: '腾讯科技投递超过 7 天未更新',
    description: '建议联系招聘方确认当前流程状态。',
    dueAt: daysAgo(0),
    targetType: 'application',
    targetId: '1',
    route: '/resume/application/1',
    source: 'system',
    createdAt: daysAgo(0),
    updatedAt: daysAgo(0)
  },
  {
    id: 'system-interview-1',
    title: '准备某创业公司终面复盘',
    description: '整理终面追问和 Offer 沟通要点。',
    dueAt: daysAgo(-2),
    targetType: 'interview',
    targetId: '2',
    route: '/resume/application/3',
    source: 'system',
    createdAt: daysAgo(0),
    updatedAt: daysAgo(0)
  }
]

const allTrend: DashboardTrendPoint[] = [
  { date: '2026-04-08', count: 1 },
  { date: '2026-04-10', count: 1 },
  { date: '2026-04-12', count: 1 },
  { date: '2026-04-22', count: 2 },
  { date: '2026-05-02', count: 1 },
  { date: '2026-05-09', count: 2 },
  { date: '2026-05-15', count: 1 }
]

const activities: ActivityLog[] = [
  {
    id: '1',
    action: 'updated',
    targetType: 'application',
    targetId: '1',
    targetLabel: '腾讯科技 / 前端开发工程师',
    summary: '腾讯科技 / 前端开发工程师 投递状态更新为 HR 面',
    route: '/resume/application/1',
    createdAt: iso('2026-05-17T14:30:00.000Z')
  },
  {
    id: '2',
    action: 'created',
    targetType: 'interview',
    targetId: '2',
    targetLabel: '终面 / 通过',
    summary: '终面 / 通过 面试进展创建',
    route: '/resume/application/3',
    createdAt: iso('2026-05-16T11:00:00.000Z')
  },
  {
    id: '3',
    action: 'status_changed',
    targetType: 'resume',
    targetId: '4',
    targetLabel: '全栈岗位定制版',
    summary: '全栈岗位定制版 简历版本状态变更',
    route: '/resume/version/4',
    createdAt: iso('2026-05-14T09:00:00.000Z')
  },
  {
    id: '4',
    action: 'created',
    targetType: 'question',
    targetId: '1',
    targetLabel: 'Vue 响应式原理',
    summary: 'Vue 响应式原理 面试题创建',
    route: '/interview-question-bank/question/1',
    createdAt: iso('2026-05-12T10:20:00.000Z')
  }
]

const validTargetTypes: DashboardTargetType[] = ['application', 'interview', 'question', 'resume']

const toRange = (raw: any): DashboardRange => {
  const range = String(raw || '30d')
  return ['7d', '30d', '90d', 'all'].includes(range) ? range as DashboardRange : '30d'
}

const trendLimit: Record<DashboardRange, number> = {
  '7d': 2,
  '30d': 5,
  '90d': 7,
  all: 7
}

const buildDashboard = (range: DashboardRange): DashboardResponse => {
  const trend = allTrend.slice(-trendLimit[range])
  const totalApplications = range === '7d' ? 3 : range === '30d' ? 7 : 9

  return {
    range,
    metrics: {
      totalApplications,
      interviewing: 3,
      pendingFollowUp: manualTodos.filter(item => !item.completedAt).length + systemTodos.length,
      offers: 1
    },
    statusPipeline: [
      { status: 'applied', label: '已投递', count: 2 },
      { status: 'viewed', label: '已查看', count: 1 },
      { status: 'first_interview', label: '一面', count: 2 },
      { status: 'hr_interview', label: 'HR 面', count: 1 },
      { status: 'offer', label: 'Offer', count: 1 },
      { status: 'rejected', label: '已拒绝', count: 2 }
    ],
    applicationTrend: trend,
    companyRanking: [
      { name: '腾讯科技', count: 2 },
      { name: '字节跳动', count: 2 },
      { name: '某创业公司', count: 1 },
      { name: '阿里云', count: 1 },
      { name: '美团', count: 1 }
    ],
    jobTitleDistribution: [
      { name: '前端开发工程师', count: 4 },
      { name: '全栈工程师', count: 2 },
      { name: 'Node.js 工程师', count: 1 }
    ],
    sourceDistribution: [
      { name: 'Boss', count: 3 },
      { name: '内推', count: 2 },
      { name: '官网', count: 2 }
    ],
    todos: [
      ...manualTodos.filter(item => !item.completedAt),
      ...systemTodos
    ],
    activities
  }
}

const getManualTodo = (id: string) => manualTodos.find(item => item.id === id)

const trimOptional = (value: unknown) => typeof value === 'string' ? value.trim() : undefined

const validateTitle = (title: unknown, required: boolean) => {
  if (title === undefined && !required) return true
  if (typeof title !== 'string') return false
  const trimmed = title.trim()
  return trimmed.length > 0 && trimmed.length <= 200
}

const validateDueAt = (dueAt: unknown) => {
  if (dueAt === undefined) return true
  if (typeof dueAt !== 'string') return false
  return !Number.isNaN(new Date(dueAt).getTime())
}

const validateTargetPair = (targetType: unknown, targetId: unknown) => {
  if (targetType === undefined && targetId === undefined) return true
  if (typeof targetType !== 'string' || !validTargetTypes.includes(targetType as DashboardTargetType)) return false
  if (typeof targetId !== 'string' || targetId.trim() === '') return false
  return Boolean(buildRoute(targetType as DashboardTargetType, targetId))
}

const validateCreateTodo = (body: any) =>
  validateTitle(body.title, true) &&
  validateDueAt(body.dueAt) &&
  validateTargetPair(body.targetType, body.targetId)

const validateUpdateTodo = (body: UpdateDashboardTodoRequest) => {
  if (!validateTitle(body.title, false) || !validateDueAt(body.dueAt)) return false
  if (body.clearTarget) return true
  if (body.targetType === undefined && body.targetId === undefined) return true
  return validateTargetPair(body.targetType, body.targetId)
}

const ensureTodoMutable = (req: Request, res: Response) => {
  const todo = getManualTodo(req.params.id)
  if (!todo) {
    notFound(res)
    return null
  }
  return todo
}

export const getDashboard = (req: Request, res: Response) => success(res, buildDashboard(toRange(req.query.range)))

export const getDashboardTodos = (req: Request, res: Response) => {
  const { completed } = req.query
  let filteredTodos = [...manualTodos]

  if (completed !== undefined) {
    const rawCompleted = String(completed)
    if (!['true', 'false'].includes(rawCompleted)) {
      return badRequest(res)
    }
    const isCompleted = rawCompleted === 'true'
    filteredTodos = filteredTodos.filter(item => Boolean(item.completedAt) === isCompleted)
  }

  return success(res, { list: filteredTodos })
}

export const createDashboardTodo = (req: Request, res: Response) => {
  if (!validateCreateTodo(req.body)) {
    return badRequest(res)
  }

  const current = new Date().toISOString()
  const targetType = req.body.targetType as DashboardTargetType | undefined
  const targetId = trimOptional(req.body.targetId)
  const todo: DashboardTodo = {
    id: String(nextTodoId++),
    title: req.body.title.trim(),
    description: trimOptional(req.body.description) || '',
    dueAt: req.body.dueAt,
    targetType,
    targetId,
    route: buildRoute(targetType, targetId),
    source: 'manual',
    createdAt: current,
    updatedAt: current
  }

  manualTodos.unshift(todo)
  return success(res, todo)
}

export const updateDashboardTodo = (req: Request, res: Response) => {
  const todo = ensureTodoMutable(req, res)
  if (!todo) return

  const body = req.body as UpdateDashboardTodoRequest
  if (!validateUpdateTodo(body)) {
    return badRequest(res)
  }

  if (body.title !== undefined) todo.title = body.title.trim()
  if (body.description !== undefined) todo.description = trimOptional(body.description) || ''
  if (body.clearDueAt) {
    delete todo.dueAt
  } else if (body.dueAt !== undefined) {
    todo.dueAt = body.dueAt
  }
  if (body.clearTarget) {
    delete todo.targetType
    delete todo.targetId
    delete todo.route
  } else if (body.targetType !== undefined || body.targetId !== undefined) {
    todo.targetType = body.targetType
    todo.targetId = trimOptional(body.targetId)
    todo.route = buildRoute(todo.targetType, todo.targetId)
  }
  todo.updatedAt = new Date().toISOString()

  return success(res, todo)
}

export const completeDashboardTodo = (req: Request, res: Response) => {
  const todo = ensureTodoMutable(req, res)
  if (!todo) return

  const current = new Date().toISOString()
  todo.completedAt = current
  todo.updatedAt = current

  return success(res, todo)
}

export const deleteDashboardTodo = (req: Request, res: Response) => {
  const index = manualTodos.findIndex(item => item.id === req.params.id)
  if (index === -1) {
    return notFound(res)
  }

  manualTodos.splice(index, 1)
  return success(res, null)
}

const buildRoute = (targetType?: DashboardTargetType, targetId?: string) => {
  const id = trimOptional(targetId)
  if (!targetType || !id) return undefined

  switch (targetType) {
    case 'application':
      return getApplicationByIdInternal(id) ? `/resume/application/${id}` : undefined
    case 'interview': {
      const progress = getInterviewProgressByIdInternal(id)
      return progress ? `/resume/application/${progress.applicationId}` : undefined
    }
    case 'question':
      return getInterviewQuestionByIdInternal(id) ? `/interview-question-bank/question/${id}` : undefined
    case 'resume':
      return getResumeRepoByIdInternal(id) ? `/resume/repo/${id}` : undefined
    default:
      return undefined
  }
}
