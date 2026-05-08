import { Response, Request } from 'express'

type JobType = 'campus' | 'social' | 'internship'

interface ResumeRepo {
  id: string
  userId: string
  name: string
  jobType: JobType
  targetPosition: string
  description?: string
  versionCount: number
  createdAt: string
  updatedAt: string
}

const repoList: ResumeRepo[] = [
  {
    id: '1',
    userId: '1',
    name: '前端工程师-互联网大厂',
    jobType: 'campus',
    targetPosition: '前端工程师',
    description: '针对互联网大厂的前端岗位',
    versionCount: 3,
    createdAt: new Date('2024-01-15').toISOString(),
    updatedAt: new Date('2024-03-10').toISOString()
  },
  {
    id: '2',
    userId: '1',
    name: '全栈开发-创业公司',
    jobType: 'social',
    targetPosition: '全栈工程师',
    description: '适合创业公司的全栈岗位',
    versionCount: 2,
    createdAt: new Date('2024-02-01').toISOString(),
    updatedAt: new Date('2024-02-20').toISOString()
  },
  {
    id: '3',
    userId: '1',
    name: 'Node.js后端',
    jobType: 'social',
    targetPosition: 'Node.js工程师',
    description: '专注于Node.js后端开发',
    versionCount: 1,
    createdAt: new Date('2024-02-15').toISOString(),
    updatedAt: new Date('2024-02-15').toISOString()
  },
  {
    id: '4',
    userId: '1',
    name: '前端实习-暑期项目',
    jobType: 'internship',
    targetPosition: '前端实习生',
    description: '面向暑期实习和日常实习的前端岗位',
    versionCount: 0,
    createdAt: new Date('2024-03-01').toISOString(),
    updatedAt: new Date('2024-03-01').toISOString()
  }
]

export const incrementRepoVersionCount = (repoId: string) => {
  const repo = repoList.find(item => item.id === repoId)
  if (!repo) return
  repo.versionCount += 1
  repo.updatedAt = new Date().toISOString()
}

export const decrementRepoVersionCount = (repoId: string) => {
  const repo = repoList.find(item => item.id === repoId)
  if (!repo) return
  repo.versionCount = Math.max(0, repo.versionCount - 1)
  repo.updatedAt = new Date().toISOString()
}

export const getResumeRepos = (req: Request, res: Response) => {
  const { keyword, jobType, page = 1, pageSize = 10 } = req.query

  let filteredRepos = [...repoList]

  // 关键词搜索
  if (keyword) {
    const searchTerm = (keyword as string).toLowerCase()
    filteredRepos = filteredRepos.filter(repo =>
      repo.name.toLowerCase().includes(searchTerm) ||
      repo.targetPosition.toLowerCase().includes(searchTerm)
    )
  }

  // 求职类型筛选
  if (jobType) {
    filteredRepos = filteredRepos.filter(repo => repo.jobType === jobType)
  }

  // 分页
  const pageNum = Number(page)
  const size = Number(pageSize)
  const start = (pageNum - 1) * size
  const end = start + size
  const pagedRepos = filteredRepos.slice(start, end)

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list: pagedRepos,
      total: filteredRepos.length,
      page: pageNum,
      pageSize: size
    }
  })
}

export const getResumeRepoById = (req: Request, res: Response) => {
  const { id } = req.params
  const repo = repoList.find(r => r.id === id)

  if (!repo) {
    return res.status(404).json({
      code: 50004,
      message: '仓库不存在'
    })
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: repo
  })
}

export const createResumeRepo = (req: Request, res: Response) => {
  const { name, jobType, targetPosition, description } = req.body

  const newRepo: ResumeRepo = {
    id: String(repoList.length + 1),
    userId: '1', // Mock userId
    name,
    jobType,
    targetPosition,
    description,
    versionCount: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }

  repoList.push(newRepo)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newRepo
  })
}

export const updateResumeRepo = (req: Request, res: Response) => {
  const { id } = req.params
  const repo = repoList.find(r => r.id === id)

  if (!repo) {
    return res.status(404).json({
      code: 50004,
      message: '仓库不存在'
    })
  }

  Object.assign(repo, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  return res.json({
    code: 20000,
    message: '更新成功',
    data: repo
  })
}

export const deleteResumeRepo = (req: Request, res: Response) => {
  const { id } = req.params
  const index = repoList.findIndex(r => r.id === id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '仓库不存在'
    })
  }

  repoList.splice(index, 1)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}
