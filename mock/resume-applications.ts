import { Request, Response } from 'express'

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

const applicationList: ResumeApplication[] = [
  {
    id: '1',
    jobPostingId: '1',
    resumeVersionId: '3',
    repoId: '1',
    companyName: '腾讯科技',
    companyIntro: '面向海量用户的互联网科技公司，业务覆盖社交、内容、游戏与云服务。',
    jobTitle: '前端开发工程师',
    jobRequirements: '熟悉 Vue/React，具备工程化能力与性能优化经验，有中大型项目经历更佳。',
    base: '深圳',
    salaryRange: '25k-35k x 16',
    sourcePlatform: 'Boss',
    sourceUrl: 'https://www.zhipin.com/job_detail/tencent-frontend',
    deliveryChannel: '官网',
    appliedAt: new Date('2026-04-10').toISOString(),
    currentStatus: ApplicationStatus.HR_INTERVIEW,
    interviewSummary: '整体流程推进顺利，技术面表现较稳定，HR 面重点关注稳定性与意向度。',
    interviewNotes: '一面追问性能优化和工程化；二面关注项目复杂度与团队协作；HR 面关注职业规划。',
    resumeMatchScore: 5,
    interviewPerformanceScore: 4,
    roleInterestScore: 5,
    overallScore: 9,
    remark: '大厂重点投递',
    createdAt: new Date('2026-04-10').toISOString(),
    updatedAt: new Date('2026-04-18').toISOString()
  },
  {
    id: '2',
    jobPostingId: '2',
    resumeVersionId: '3',
    repoId: '1',
    companyName: '字节跳动',
    companyIntro: '以内容平台和企业服务为核心的互联网公司。',
    jobTitle: '资深前端工程师',
    jobRequirements: '要求扎实的 JavaScript 基础，熟悉 React 生态，具备性能优化和工程实践经验。',
    base: '北京',
    salaryRange: '30k-45k x 15',
    sourcePlatform: '内推',
    sourceUrl: 'https://jobs.bytedance.com/frontend-senior',
    deliveryChannel: '内推',
    appliedAt: new Date('2026-04-12').toISOString(),
    currentStatus: ApplicationStatus.REJECTED,
    interviewSummary: '一面技术题完成度一般，系统设计表达不够聚焦。',
    interviewNotes: '需要强化复杂场景拆解和高并发前端方案表达。',
    resumeMatchScore: 4,
    interviewPerformanceScore: 2,
    roleInterestScore: 4,
    overallScore: 6,
    remark: '一面后挂',
    createdAt: new Date('2026-04-12').toISOString(),
    updatedAt: new Date('2026-04-16').toISOString()
  },
  {
    id: '3',
    jobPostingId: '3',
    resumeVersionId: '4',
    repoId: '2',
    companyName: '某创业公司',
    companyIntro: '高速增长的 SaaS 创业团队，强调全栈和业务交付能力。',
    jobTitle: '全栈工程师',
    jobRequirements: '具备 Node.js、Vue 和 MySQL 经验，能独立推进需求落地。',
    base: '上海',
    salaryRange: '20k-28k x 14',
    sourcePlatform: '官网',
    sourceUrl: 'https://company.example.com/jobs/fullstack',
    deliveryChannel: 'Boss 直聘',
    appliedAt: new Date('2026-04-08').toISOString(),
    currentStatus: ApplicationStatus.OFFER,
    interviewSummary: '业务匹配度较高，整体沟通顺畅，已进入 Offer 阶段。',
    interviewNotes: '技术问题较贴近日常业务，重点在交付和 owner 意识。',
    resumeMatchScore: 5,
    interviewPerformanceScore: 5,
    roleInterestScore: 4,
    overallScore: 9,
    remark: '已发口头 offer',
    createdAt: new Date('2026-04-08').toISOString(),
    updatedAt: new Date('2026-04-20').toISOString()
  }
]

export const getApplicationList = () => applicationList

export const getApplicationByIdInternal = (id: string) => applicationList.find(item => item.id === id)

export const updateApplicationStatusInternal = (applicationId: string, status: ApplicationStatus) => {
  const application = getApplicationByIdInternal(applicationId)
  if (!application) return

  application.currentStatus = status
  application.updatedAt = new Date().toISOString()
}

export const hasApplicationsForResumeVersion = (resumeVersionId: string) =>
  applicationList.some(item => item.resumeVersionId === resumeVersionId)

export const getResumeApplications = (req: Request, res: Response) => {
  const { resumeVersionId, repoId, keyword, currentStatus, page = 1, pageSize = 10 } = req.query

  let filteredList = [...applicationList]

  if (resumeVersionId) {
    filteredList = filteredList.filter(item => item.resumeVersionId === resumeVersionId)
  }

  if (repoId) {
    filteredList = filteredList.filter(item => item.repoId === repoId)
  }

  if (keyword) {
    const searchTerm = String(keyword).toLowerCase()
    filteredList = filteredList.filter(item =>
      item.companyName.toLowerCase().includes(searchTerm) ||
      item.jobTitle.toLowerCase().includes(searchTerm)
    )
  }

  if (currentStatus) {
    filteredList = filteredList.filter(item => item.currentStatus === currentStatus)
  }

  filteredList.sort((a, b) => new Date(b.appliedAt).getTime() - new Date(a.appliedAt).getTime())

  const pageNum = Number(page)
  const size = Number(pageSize)
  const start = (pageNum - 1) * size
  const end = start + size

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list: filteredList.slice(start, end),
      total: filteredList.length,
      page: pageNum,
      pageSize: size
    }
  })
}

export const getResumeApplicationById = (req: Request, res: Response) => {
  const application = getApplicationByIdInternal(req.params.id)

  if (!application) {
    return res.status(404).json({
      code: 50004,
      message: '投递记录不存在'
    })
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: application
  })
}

export const createResumeApplication = (req: Request, res: Response) => {
  const now = new Date().toISOString()
  const newApplication: ResumeApplication = {
    id: String(applicationList.length + 1),
    resumeVersionId: req.body.resumeVersionId,
    jobPostingId: req.body.jobPostingId,
    repoId: req.body.repoId,
    companyName: req.body.companyName,
    companyIntro: req.body.companyIntro,
    jobTitle: req.body.jobTitle,
    jobRequirements: req.body.jobRequirements,
    base: req.body.base,
    salaryRange: req.body.salaryRange,
    sourcePlatform: req.body.sourcePlatform,
    sourceUrl: req.body.sourceUrl,
    deliveryChannel: req.body.deliveryChannel,
    appliedAt: req.body.appliedAt,
    currentStatus: req.body.currentStatus,
    interviewSummary: req.body.interviewSummary,
    interviewNotes: req.body.interviewNotes,
    resumeMatchScore: req.body.resumeMatchScore,
    interviewPerformanceScore: req.body.interviewPerformanceScore,
    roleInterestScore: req.body.roleInterestScore,
    overallScore: req.body.overallScore,
    remark: req.body.remark,
    createdAt: now,
    updatedAt: now
  }

  applicationList.push(newApplication)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newApplication
  })
}

export const updateResumeApplication = (req: Request, res: Response) => {
  const application = getApplicationByIdInternal(req.params.id)

  if (!application) {
    return res.status(404).json({
      code: 50004,
      message: '投递记录不存在'
    })
  }

  Object.assign(application, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  return res.json({
    code: 20000,
    message: '更新成功',
    data: application
  })
}

export const deleteResumeApplication = (req: Request, res: Response) => {
  const index = applicationList.findIndex(item => item.id === req.params.id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '投递记录不存在'
    })
  }

  applicationList.splice(index, 1)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}
