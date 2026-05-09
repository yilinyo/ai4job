import { Request, Response } from 'express'

export enum JobPostingStatus {
  PENDING = 'pending',
  APPLIED = 'applied',
  NOT_FIT = 'not_fit',
  CLOSED = 'closed'
}

interface JobPosting {
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

const jobPostingList: JobPosting[] = [
  {
    id: '1',
    companyName: '腾讯科技',
    companyIntro: '互联网科技公司，业务覆盖社交、内容、游戏与云服务。',
    jobTitle: '前端开发工程师',
    jobRequirements: '熟悉 Vue / React，具备工程化与性能优化经验，有中大型项目经历更佳。',
    base: '深圳',
    salaryRange: '25k-35k x 16',
    sourcePlatform: 'Boss',
    sourceUrl: 'https://www.zhipin.com/job_detail/tencent-frontend',
    status: JobPostingStatus.PENDING,
    remark: '优先级高，适合大厂定制版简历',
    createdAt: new Date('2026-05-01').toISOString(),
    updatedAt: new Date('2026-05-01').toISOString()
  },
  {
    id: '2',
    companyName: '字节跳动',
    companyIntro: '内容平台与企业服务并行发展的互联网公司。',
    jobTitle: '资深前端工程师',
    jobRequirements: '要求扎实的 JavaScript 基础，熟悉 React 生态，具备系统设计能力。',
    base: '北京',
    salaryRange: '30k-45k x 15',
    sourcePlatform: '内推',
    sourceUrl: 'https://jobs.bytedance.com/frontend-senior',
    status: JobPostingStatus.APPLIED,
    remark: '已经投递过一次，后续可能考虑新版简历再次投递',
    createdAt: new Date('2026-04-28').toISOString(),
    updatedAt: new Date('2026-05-02').toISOString()
  },
  {
    id: '3',
    companyName: '某创业公司',
    companyIntro: '高速增长的 SaaS 创业团队，强调业务交付与 owner 意识。',
    jobTitle: '全栈工程师',
    jobRequirements: '具备 Node.js、Vue 和 MySQL 经验，能独立推进需求落地。',
    base: '上海',
    salaryRange: '20k-28k x 14',
    sourcePlatform: '官网',
    sourceUrl: 'https://company.example.com/jobs/fullstack',
    status: JobPostingStatus.CLOSED,
    remark: '岗位已下线',
    createdAt: new Date('2026-04-20').toISOString(),
    updatedAt: new Date('2026-04-30').toISOString()
  }
]

export const getJobPostingByIdInternal = (id: string) => jobPostingList.find(item => item.id === id)

export const getJobPostings = (req: Request, res: Response) => {
  const { keyword, sourcePlatform, status, page = 1, pageSize = 10 } = req.query

  let filteredList = [...jobPostingList]

  if (keyword) {
    const searchTerm = String(keyword).toLowerCase()
    filteredList = filteredList.filter(item =>
      item.companyName.toLowerCase().includes(searchTerm) ||
      item.jobTitle.toLowerCase().includes(searchTerm)
    )
  }

  if (sourcePlatform) {
    filteredList = filteredList.filter(item => item.sourcePlatform === sourcePlatform)
  }

  if (status) {
    filteredList = filteredList.filter(item => item.status === status)
  }

  filteredList.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

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

export const getJobPostingById = (req: Request, res: Response) => {
  const item = getJobPostingByIdInternal(req.params.id)

  if (!item) {
    return res.status(404).json({
      code: 50004,
      message: '岗位不存在'
    })
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: item
  })
}

export const createJobPosting = (req: Request, res: Response) => {
  const now = new Date().toISOString()
  const newItem: JobPosting = {
    id: String(jobPostingList.length + 1),
    companyName: req.body.companyName,
    companyIntro: req.body.companyIntro,
    jobTitle: req.body.jobTitle,
    jobRequirements: req.body.jobRequirements,
    base: req.body.base,
    salaryRange: req.body.salaryRange,
    sourcePlatform: req.body.sourcePlatform,
    sourceUrl: req.body.sourceUrl,
    status: req.body.status,
    remark: req.body.remark,
    createdAt: now,
    updatedAt: now
  }

  jobPostingList.push(newItem)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newItem
  })
}

export const updateJobPosting = (req: Request, res: Response) => {
  const item = getJobPostingByIdInternal(req.params.id)

  if (!item) {
    return res.status(404).json({
      code: 50004,
      message: '岗位不存在'
    })
  }

  Object.assign(item, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  return res.json({
    code: 20000,
    message: '更新成功',
    data: item
  })
}

export const deleteJobPosting = (req: Request, res: Response) => {
  const index = jobPostingList.findIndex(item => item.id === req.params.id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '岗位不存在'
    })
  }

  jobPostingList.splice(index, 1)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}

export const parseJobPosting = (req: Request, res: Response) => {
  const url = String(req.body.url || '')

  if (!url) {
    return res.status(400).json({
      code: 50001,
      message: '岗位链接不能为空'
    })
  }

  const sourcePlatform = url.includes('zhipin') ? 'Boss' : url.includes('jobs.bytedance') ? '官网' : '其他'

  return res.json({
    code: 20000,
    message: '解析成功',
    data: {
      companyName: sourcePlatform === 'Boss' ? '腾讯科技' : '示例公司',
      companyIntro: '通过解析接口返回的公司介绍示例',
      jobTitle: sourcePlatform === 'Boss' ? '前端开发工程师' : '待确认岗位',
      jobRequirements: '通过解析接口返回的岗位要求示例，用户可继续修改',
      base: sourcePlatform === 'Boss' ? '深圳' : '北京',
      salaryRange: '25k-35k x 16',
      sourcePlatform,
      sourceUrl: url
    }
  })
}
