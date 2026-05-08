import { Request, Response } from 'express'
import { ApplicationStatus, getApplicationByIdInternal, updateApplicationStatusInternal } from './resume-applications'

export enum InterviewStage {
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

export enum InterviewResult {
  PASSED = 'passed',
  PENDING = 'pending',
  FAILED = 'failed'
}

interface InterviewProgress {
  id: string
  applicationId: string
  stage: InterviewStage
  occurredAt: string
  result: InterviewResult
  interviewerOrTeam?: string
  note?: string
  createdAt: string
  updatedAt: string
}

let progressList: InterviewProgress[] = [
  {
    id: '1',
    applicationId: '1',
    stage: InterviewStage.APPLIED,
    occurredAt: new Date('2026-04-10').toISOString(),
    result: InterviewResult.PASSED,
    note: '官网投递成功',
    createdAt: new Date('2026-04-10').toISOString(),
    updatedAt: new Date('2026-04-10').toISOString()
  },
  {
    id: '2',
    applicationId: '1',
    stage: InterviewStage.VIEWED,
    occurredAt: new Date('2026-04-11').toISOString(),
    result: InterviewResult.PASSED,
    note: '简历已被查看',
    createdAt: new Date('2026-04-11').toISOString(),
    updatedAt: new Date('2026-04-11').toISOString()
  },
  {
    id: '3',
    applicationId: '1',
    stage: InterviewStage.FIRST_INTERVIEW,
    occurredAt: new Date('2026-04-14').toISOString(),
    result: InterviewResult.PASSED,
    interviewerOrTeam: '前端基础架构组',
    note: '重点问性能优化和工程化体系',
    createdAt: new Date('2026-04-14').toISOString(),
    updatedAt: new Date('2026-04-14').toISOString()
  },
  {
    id: '4',
    applicationId: '1',
    stage: InterviewStage.SECOND_INTERVIEW,
    occurredAt: new Date('2026-04-16').toISOString(),
    result: InterviewResult.PASSED,
    interviewerOrTeam: '业务技术负责人',
    note: '项目复杂度和跨团队协作表现良好',
    createdAt: new Date('2026-04-16').toISOString(),
    updatedAt: new Date('2026-04-16').toISOString()
  },
  {
    id: '5',
    applicationId: '1',
    stage: InterviewStage.HR_INTERVIEW,
    occurredAt: new Date('2026-04-18').toISOString(),
    result: InterviewResult.PENDING,
    interviewerOrTeam: 'HR',
    note: '待最终结果',
    createdAt: new Date('2026-04-18').toISOString(),
    updatedAt: new Date('2026-04-18').toISOString()
  },
  {
    id: '6',
    applicationId: '2',
    stage: InterviewStage.APPLIED,
    occurredAt: new Date('2026-04-12').toISOString(),
    result: InterviewResult.PASSED,
    note: '内推完成',
    createdAt: new Date('2026-04-12').toISOString(),
    updatedAt: new Date('2026-04-12').toISOString()
  },
  {
    id: '7',
    applicationId: '2',
    stage: InterviewStage.FIRST_INTERVIEW,
    occurredAt: new Date('2026-04-16').toISOString(),
    result: InterviewResult.FAILED,
    interviewerOrTeam: '业务团队',
    note: '系统设计问题回答不完整',
    createdAt: new Date('2026-04-16').toISOString(),
    updatedAt: new Date('2026-04-16').toISOString()
  },
  {
    id: '8',
    applicationId: '3',
    stage: InterviewStage.APPLIED,
    occurredAt: new Date('2026-04-08').toISOString(),
    result: InterviewResult.PASSED,
    note: 'Boss 直聘投递',
    createdAt: new Date('2026-04-08').toISOString(),
    updatedAt: new Date('2026-04-08').toISOString()
  },
  {
    id: '9',
    applicationId: '3',
    stage: InterviewStage.FIRST_INTERVIEW,
    occurredAt: new Date('2026-04-11').toISOString(),
    result: InterviewResult.PASSED,
    note: '偏业务场景，沟通顺畅',
    createdAt: new Date('2026-04-11').toISOString(),
    updatedAt: new Date('2026-04-11').toISOString()
  },
  {
    id: '10',
    applicationId: '3',
    stage: InterviewStage.HR_INTERVIEW,
    occurredAt: new Date('2026-04-15').toISOString(),
    result: InterviewResult.PASSED,
    note: '进入 offer 流程',
    createdAt: new Date('2026-04-15').toISOString(),
    updatedAt: new Date('2026-04-15').toISOString()
  },
  {
    id: '11',
    applicationId: '3',
    stage: InterviewStage.OFFER,
    occurredAt: new Date('2026-04-20').toISOString(),
    result: InterviewResult.PASSED,
    note: '已收到口头 offer',
    createdAt: new Date('2026-04-20').toISOString(),
    updatedAt: new Date('2026-04-20').toISOString()
  }
]

export const deleteInterviewProgressByApplicationId = (applicationId: string) => {
  progressList = progressList.filter(item => item.applicationId !== applicationId)
}

const stageToStatusMap: Record<InterviewStage, ApplicationStatus> = {
  [InterviewStage.APPLIED]: ApplicationStatus.APPLIED,
  [InterviewStage.VIEWED]: ApplicationStatus.VIEWED,
  [InterviewStage.WRITTEN_TEST]: ApplicationStatus.WRITTEN_TEST,
  [InterviewStage.FIRST_INTERVIEW]: ApplicationStatus.FIRST_INTERVIEW,
  [InterviewStage.SECOND_INTERVIEW]: ApplicationStatus.SECOND_INTERVIEW,
  [InterviewStage.FINAL_INTERVIEW]: ApplicationStatus.FINAL_INTERVIEW,
  [InterviewStage.HR_INTERVIEW]: ApplicationStatus.HR_INTERVIEW,
  [InterviewStage.OFFER]: ApplicationStatus.OFFER,
  [InterviewStage.REJECTED]: ApplicationStatus.REJECTED,
  [InterviewStage.CLOSED]: ApplicationStatus.CLOSED
}

const syncApplicationStatus = (applicationId: string) => {
  const application = getApplicationByIdInternal(applicationId)
  if (!application) return

  const sortedProgress = progressList
    .filter(item => item.applicationId === applicationId)
    .sort((a, b) => new Date(a.occurredAt).getTime() - new Date(b.occurredAt).getTime())

  if (!sortedProgress.length) return

  const latest = sortedProgress[sortedProgress.length - 1]
  const nextStatus = latest.result === InterviewResult.FAILED
    ? ApplicationStatus.REJECTED
    : stageToStatusMap[latest.stage]

  updateApplicationStatusInternal(applicationId, nextStatus)
}

export const getInterviewProgressList = (req: Request, res: Response) => {
  const { applicationId } = req.query

  if (!applicationId) {
    return res.status(400).json({
      code: 50001,
      message: 'applicationId 不能为空'
    })
  }

  const list = progressList
    .filter(item => item.applicationId === applicationId)
    .sort((a, b) => new Date(a.occurredAt).getTime() - new Date(b.occurredAt).getTime())

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list
    }
  })
}

export const createInterviewProgress = (req: Request, res: Response) => {
  const now = new Date().toISOString()
  const newItem: InterviewProgress = {
    id: String(progressList.length + 1),
    applicationId: req.body.applicationId,
    stage: req.body.stage,
    occurredAt: req.body.occurredAt,
    result: req.body.result,
    interviewerOrTeam: req.body.interviewerOrTeam,
    note: req.body.note,
    createdAt: now,
    updatedAt: now
  }

  progressList.push(newItem)
  syncApplicationStatus(newItem.applicationId)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newItem
  })
}

export const updateInterviewProgress = (req: Request, res: Response) => {
  const progress = progressList.find(item => item.id === req.params.id)

  if (!progress) {
    return res.status(404).json({
      code: 50004,
      message: '面试进展不存在'
    })
  }

  Object.assign(progress, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  syncApplicationStatus(progress.applicationId)

  return res.json({
    code: 20000,
    message: '更新成功',
    data: progress
  })
}

export const deleteInterviewProgress = (req: Request, res: Response) => {
  const index = progressList.findIndex(item => item.id === req.params.id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '面试进展不存在'
    })
  }

  const { applicationId } = progressList[index]
  progressList.splice(index, 1)
  syncApplicationStatus(applicationId)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}
