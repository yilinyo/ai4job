import { Request, Response } from 'express'

export enum InterviewQuestionDifficulty {
  EASY = 'easy',
  MEDIUM = 'medium',
  HARD = 'hard'
}

export enum InterviewQuestionType {
  KNOWLEDGE = 'knowledge',
  SCENARIO = 'scenario',
  CODING = 'coding',
  ALGORITHM = 'algorithm',
  PROJECT = 'project',
  BEHAVIOR = 'behavior',
  DISCUSSION = 'discussion'
}

export enum QuestionMasteryStatus {
  UNREVIEWED = 'unreviewed',
  WEAK = 'weak',
  NORMAL = 'normal',
  MASTERED = 'mastered'
}

export enum QuestionAnswerPerformance {
  FAILED = 'failed',
  VAGUE = 'vague',
  ANSWERED = 'answered',
  GOOD = 'good'
}

export enum FlashcardReviewResult {
  FAILED = 'failed',
  VAGUE = 'vague',
  MASTERED = 'mastered'
}

export enum FlashcardReviewMode {
  RANDOM = 'random',
  TAG = 'tag',
  DIFFICULTY = 'difficulty',
  FAVORITE_FIRST = 'favorite_first',
  WEAK_FIRST = 'weak_first',
  FREQUENT_FIRST = 'frequent_first'
}

interface InterviewQuestionRecord {
  id: string
  title: string
  content: string
  answerAnalysis: string
  tags: string[]
  difficulty: InterviewQuestionDifficulty
  questionType: InterviewQuestionType
  masteryStatus: QuestionMasteryStatus
  isFavorite: boolean
  source?: string
  lastReviewedAt?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

interface QuestionOccurrence {
  id: string
  questionId: string
  applicationId: string
  interviewProgressId?: string
  companyNameSnapshot: string
  jobTitleSnapshot: string
  interviewStageSnapshot?: string
  occurredAt: string
  actualQuestion?: string
  answerPerformance?: QuestionAnswerPerformance
  note?: string
  createdAt: string
  updatedAt: string
}

interface FlashcardReview {
  id: string
  questionId: string
  reviewedAt: string
  result: FlashcardReviewResult
  reviewMode: FlashcardReviewMode
  createdAt: string
}

const questionList: InterviewQuestionRecord[] = [
  {
    id: '1',
    title: 'Vue 2 响应式原理和数组更新陷阱',
    content: '请说明 Vue 2 的响应式实现方式，以及为什么直接通过下标修改数组可能不触发更新。',
    answerAnalysis: '核心是 Object.defineProperty 劫持属性访问。数组通过重写变异方法实现依赖通知，直接设置索引或修改 length 无法被拦截，需要使用 Vue.set 或 splice。',
    tags: ['Vue', 'JavaScript', '浏览器'],
    difficulty: InterviewQuestionDifficulty.MEDIUM,
    questionType: InterviewQuestionType.KNOWLEDGE,
    masteryStatus: QuestionMasteryStatus.NORMAL,
    isFavorite: true,
    source: '面试复盘',
    lastReviewedAt: new Date('2026-05-06T20:30:00').toISOString(),
    remark: '经常会追问到 nextTick 和异步更新队列。',
    createdAt: new Date('2026-04-14T10:00:00').toISOString(),
    updatedAt: new Date('2026-05-06T20:30:00').toISOString()
  },
  {
    id: '2',
    title: '前端性能优化体系怎么拆',
    content: '如果让你系统性梳理一个前端项目的性能优化方案，你会从哪些维度展开？',
    answerAnalysis: '可以从网络传输、构建产物、运行时渲染、缓存策略、监控和体验指标拆分。回答时要把指标、措施和验证闭环说清楚。',
    tags: ['性能优化', '工程化', '项目经历'],
    difficulty: InterviewQuestionDifficulty.HARD,
    questionType: InterviewQuestionType.SCENARIO,
    masteryStatus: QuestionMasteryStatus.WEAK,
    isFavorite: true,
    source: '自己整理',
    lastReviewedAt: new Date('2026-05-04T08:15:00').toISOString(),
    remark: '适合结合真实项目讲 bundle、懒加载和监控。',
    createdAt: new Date('2026-04-18T14:00:00').toISOString(),
    updatedAt: new Date('2026-05-04T08:15:00').toISOString()
  },
  {
    id: '3',
    title: '实现防抖和节流有什么边界差异',
    content: '请分别实现 debounce 和 throttle，并说明在输入搜索和滚动监听中的使用差异。',
    answerAnalysis: '防抖更适合用户停止输入后再触发，节流适合持续事件中的固定频率执行。实现上要说明 timer、leading / trailing 和 this / 参数透传。',
    tags: ['JavaScript', '浏览器'],
    difficulty: InterviewQuestionDifficulty.EASY,
    questionType: InterviewQuestionType.CODING,
    masteryStatus: QuestionMasteryStatus.MASTERED,
    isFavorite: false,
    source: '牛客',
    lastReviewedAt: new Date('2026-05-01T18:00:00').toISOString(),
    createdAt: new Date('2026-04-20T09:00:00').toISOString(),
    updatedAt: new Date('2026-05-01T18:00:00').toISOString()
  },
  {
    id: '4',
    title: '如何设计前端错误监控与告警闭环',
    content: '如果线上出现大量前端异常，你会如何设计采集、聚合、告警与排查机制？',
    answerAnalysis: '重点是错误分类、上下文采集、Sourcemap 还原、告警阈值和回溯链路。答案要体现发现问题、止损和复盘闭环。',
    tags: ['工程化', '系统设计', '性能优化'],
    difficulty: InterviewQuestionDifficulty.HARD,
    questionType: InterviewQuestionType.PROJECT,
    masteryStatus: QuestionMasteryStatus.UNREVIEWED,
    isFavorite: false,
    source: '文章',
    createdAt: new Date('2026-05-02T11:00:00').toISOString(),
    updatedAt: new Date('2026-05-02T11:00:00').toISOString()
  },
  {
    id: '5',
    title: '自我介绍里项目亮点怎么讲',
    content: '请结合真实项目说明你在团队里的核心贡献，以及为什么它能体现你的成长。',
    answerAnalysis: '行为面要用场景、动作、结果和反思组织表达，避免流水账。可以用 STAR 结构，并强调定量结果与角色责任。',
    tags: ['HR 面', '项目经历', '行为面试'],
    difficulty: InterviewQuestionDifficulty.MEDIUM,
    questionType: InterviewQuestionType.BEHAVIOR,
    masteryStatus: QuestionMasteryStatus.NORMAL,
    isFavorite: false,
    source: '朋友分享',
    createdAt: new Date('2026-04-28T16:00:00').toISOString(),
    updatedAt: new Date('2026-04-28T16:00:00').toISOString()
  }
]

let occurrenceList: QuestionOccurrence[] = [
  {
    id: '1',
    questionId: '1',
    applicationId: '1',
    interviewProgressId: '3',
    companyNameSnapshot: '腾讯科技',
    jobTitleSnapshot: '前端开发工程师',
    interviewStageSnapshot: 'first_interview',
    occurredAt: new Date('2026-04-14T15:00:00').toISOString(),
    actualQuestion: '为什么 Vue 2 对数组索引更新不敏感？如果要兼容怎么处理？',
    answerPerformance: QuestionAnswerPerformance.ANSWERED,
    note: '答到了 defineProperty 和 splice，但没展开依赖收集。',
    createdAt: new Date('2026-04-14T15:30:00').toISOString(),
    updatedAt: new Date('2026-04-14T15:30:00').toISOString()
  },
  {
    id: '2',
    questionId: '2',
    applicationId: '1',
    interviewProgressId: '4',
    companyNameSnapshot: '腾讯科技',
    jobTitleSnapshot: '前端开发工程师',
    interviewStageSnapshot: 'second_interview',
    occurredAt: new Date('2026-04-16T20:00:00').toISOString(),
    actualQuestion: '如果让你优化一个首屏白屏严重的中后台项目，你会怎么拆？',
    answerPerformance: QuestionAnswerPerformance.GOOD,
    note: '回答时结合拆包、缓存和监控，反馈不错。',
    createdAt: new Date('2026-04-16T20:15:00').toISOString(),
    updatedAt: new Date('2026-04-16T20:15:00').toISOString()
  },
  {
    id: '3',
    questionId: '2',
    applicationId: '2',
    interviewProgressId: '7',
    companyNameSnapshot: '字节跳动',
    jobTitleSnapshot: '资深前端工程师',
    interviewStageSnapshot: 'first_interview',
    occurredAt: new Date('2026-04-16T10:00:00').toISOString(),
    actualQuestion: '描述一下你做性能优化时的指标体系。',
    answerPerformance: QuestionAnswerPerformance.VAGUE,
    note: '没有把性能指标和业务收益讲清楚。',
    createdAt: new Date('2026-04-16T10:20:00').toISOString(),
    updatedAt: new Date('2026-04-16T10:20:00').toISOString()
  },
  {
    id: '4',
    questionId: '5',
    applicationId: '3',
    interviewProgressId: '10',
    companyNameSnapshot: '某创业公司',
    jobTitleSnapshot: '全栈工程师',
    interviewStageSnapshot: 'hr_interview',
    occurredAt: new Date('2026-04-15T14:00:00').toISOString(),
    actualQuestion: '你最有成就感的一次项目经历是什么？',
    answerPerformance: QuestionAnswerPerformance.GOOD,
    note: '重点讲了 owner 意识和项目结果。',
    createdAt: new Date('2026-04-15T14:10:00').toISOString(),
    updatedAt: new Date('2026-04-15T14:10:00').toISOString()
  }
]

let reviewList: FlashcardReview[] = [
  {
    id: '1',
    questionId: '1',
    reviewedAt: new Date('2026-05-06T20:30:00').toISOString(),
    result: FlashcardReviewResult.VAGUE,
    reviewMode: FlashcardReviewMode.TAG,
    createdAt: new Date('2026-05-06T20:30:00').toISOString()
  },
  {
    id: '2',
    questionId: '2',
    reviewedAt: new Date('2026-05-04T08:15:00').toISOString(),
    result: FlashcardReviewResult.FAILED,
    reviewMode: FlashcardReviewMode.WEAK_FIRST,
    createdAt: new Date('2026-05-04T08:15:00').toISOString()
  },
  {
    id: '3',
    questionId: '3',
    reviewedAt: new Date('2026-05-01T18:00:00').toISOString(),
    result: FlashcardReviewResult.MASTERED,
    reviewMode: FlashcardReviewMode.RANDOM,
    createdAt: new Date('2026-05-01T18:00:00').toISOString()
  }
]

const getNextId = (list: Array<{ id: string }>) => {
  const nextId = list.reduce((maxId, item) => {
    const currentId = Number(item.id)
    return Number.isNaN(currentId) ? maxId : Math.max(maxId, currentId)
  }, 0) + 1

  return String(nextId)
}

const parseTags = (value: unknown) => {
  if (!value) return []
  if (Array.isArray(value)) {
    return value
      .join(',')
      .split(',')
      .map(item => item.trim())
      .filter(Boolean)
  }

  return String(value)
    .split(',')
    .map(item => item.trim())
    .filter(Boolean)
}

const shuffleList = <T>(list: T[]) => {
  const copied = [...list]
  for (let index = copied.length - 1; index > 0; index -= 1) {
    const randomIndex = Math.floor(Math.random() * (index + 1))
    const temp = copied[index]
    copied[index] = copied[randomIndex]
    copied[randomIndex] = temp
  }
  return copied
}

const getOccurrenceCount = (questionId: string) =>
  occurrenceList.filter(item => item.questionId === questionId).length

const buildReviewSummary = (questionId: string) => {
  const reviews = reviewList
    .filter(item => item.questionId === questionId)
    .sort((a, b) => new Date(b.reviewedAt).getTime() - new Date(a.reviewedAt).getTime())

  return {
    totalReviews: reviews.length,
    failedCount: reviews.filter(item => item.result === FlashcardReviewResult.FAILED).length,
    vagueCount: reviews.filter(item => item.result === FlashcardReviewResult.VAGUE).length,
    masteredCount: reviews.filter(item => item.result === FlashcardReviewResult.MASTERED).length,
    lastReviewedAt: reviews[0]?.reviewedAt,
    lastResult: reviews[0]?.result
  }
}

const buildQuestion = (question: InterviewQuestionRecord) => ({
  ...question,
  occurrenceCount: getOccurrenceCount(question.id)
})

const buildQuestionDetail = (question: InterviewQuestionRecord) => ({
  ...buildQuestion(question),
  occurrences: occurrenceList
    .filter(item => item.questionId === question.id)
    .sort((a, b) => new Date(b.occurredAt).getTime() - new Date(a.occurredAt).getTime()),
  reviewSummary: buildReviewSummary(question.id),
  recentReviews: reviewList
    .filter(item => item.questionId === question.id)
    .sort((a, b) => new Date(b.reviewedAt).getTime() - new Date(a.reviewedAt).getTime())
    .slice(0, 5)
})

const syncQuestionReviewState = (
  question: InterviewQuestionRecord,
  result: FlashcardReviewResult,
  reviewedAt: string
) => {
  const masteryMap: Record<FlashcardReviewResult, QuestionMasteryStatus> = {
    [FlashcardReviewResult.FAILED]: QuestionMasteryStatus.WEAK,
    [FlashcardReviewResult.VAGUE]: QuestionMasteryStatus.NORMAL,
    [FlashcardReviewResult.MASTERED]: QuestionMasteryStatus.MASTERED
  }

  question.masteryStatus = masteryMap[result]
  question.lastReviewedAt = reviewedAt
  question.updatedAt = new Date().toISOString()
}

export const getInterviewQuestionByIdInternal = (id: string) =>
  questionList.find(item => item.id === id)

export const cascadeDeleteQuestionOccurrencesByApplicationId = (applicationId: string) => {
  occurrenceList = occurrenceList.filter(item => item.applicationId !== applicationId)
}

export const clearQuestionOccurrenceProgressLinkByProgressId = (progressId: string) => {
  occurrenceList = occurrenceList.map(item => {
    if (item.interviewProgressId !== progressId) return item

    return {
      ...item,
      interviewProgressId: undefined,
      updatedAt: new Date().toISOString()
    }
  })
}

const matchTags = (questionTags: string[], selectedTags: string[]) => {
  if (!selectedTags.length) return true
  return selectedTags.some(tag => questionTags.includes(tag))
}

export const getInterviewQuestions = (req: Request, res: Response) => {
  const {
    keyword,
    tags,
    difficulty,
    questionType,
    masteryStatus,
    isFavorite,
    page = 1,
    pageSize = 10
  } = req.query

  const selectedTags = parseTags(tags)

  let filteredList = questionList
    .map(buildQuestion)
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())

  if (keyword) {
    const searchTerm = String(keyword).toLowerCase()
    filteredList = filteredList.filter(item =>
      item.title.toLowerCase().includes(searchTerm) ||
      item.content.toLowerCase().includes(searchTerm) ||
      item.answerAnalysis.toLowerCase().includes(searchTerm) ||
      item.tags.some(tag => tag.toLowerCase().includes(searchTerm))
    )
  }

  if (selectedTags.length) {
    filteredList = filteredList.filter(item => matchTags(item.tags, selectedTags))
  }

  if (difficulty) {
    filteredList = filteredList.filter(item => item.difficulty === difficulty)
  }

  if (questionType) {
    filteredList = filteredList.filter(item => item.questionType === questionType)
  }

  if (masteryStatus) {
    filteredList = filteredList.filter(item => item.masteryStatus === masteryStatus)
  }

  if (typeof isFavorite !== 'undefined' && isFavorite !== '') {
    const favoriteFlag = String(isFavorite) === 'true'
    filteredList = filteredList.filter(item => item.isFavorite === favoriteFlag)
  }

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

export const getInterviewQuestionById = (req: Request, res: Response) => {
  const question = getInterviewQuestionByIdInternal(req.params.id)

  if (!question) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: buildQuestionDetail(question)
  })
}

export const createInterviewQuestion = (req: Request, res: Response) => {
  const now = new Date().toISOString()
  const newQuestion: InterviewQuestionRecord = {
    id: getNextId(questionList),
    title: req.body.title,
    content: req.body.content,
    answerAnalysis: req.body.answerAnalysis,
    tags: Array.isArray(req.body.tags) ? req.body.tags : [],
    difficulty: req.body.difficulty,
    questionType: req.body.questionType,
    masteryStatus: req.body.masteryStatus,
    isFavorite: Boolean(req.body.isFavorite),
    source: req.body.source,
    remark: req.body.remark,
    createdAt: now,
    updatedAt: now
  }

  questionList.unshift(newQuestion)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: buildQuestion(newQuestion)
  })
}

export const updateInterviewQuestion = (req: Request, res: Response) => {
  const question = getInterviewQuestionByIdInternal(req.params.id)

  if (!question) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  Object.assign(question, {
    ...req.body,
    tags: Array.isArray(req.body.tags) ? req.body.tags : question.tags,
    updatedAt: new Date().toISOString()
  })

  return res.json({
    code: 20000,
    message: '更新成功',
    data: buildQuestion(question)
  })
}

export const deleteInterviewQuestion = (req: Request, res: Response) => {
  const index = questionList.findIndex(item => item.id === req.params.id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  const questionId = questionList[index].id
  questionList.splice(index, 1)
  occurrenceList = occurrenceList.filter(item => item.questionId !== questionId)
  reviewList = reviewList.filter(item => item.questionId !== questionId)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}

export const updateInterviewQuestionFavorite = (req: Request, res: Response) => {
  const question = getInterviewQuestionByIdInternal(req.params.id)

  if (!question) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  question.isFavorite = Boolean(req.body.isFavorite)
  question.updatedAt = new Date().toISOString()

  return res.json({
    code: 20000,
    message: '收藏状态更新成功',
    data: buildQuestion(question)
  })
}

export const getQuestionOccurrences = (req: Request, res: Response) => {
  const { questionId, applicationId } = req.query

  let filteredList = [...occurrenceList]

  if (questionId) {
    filteredList = filteredList.filter(item => item.questionId === questionId)
  }

  if (applicationId) {
    filteredList = filteredList.filter(item => item.applicationId === applicationId)
  }

  filteredList.sort((a, b) => new Date(b.occurredAt).getTime() - new Date(a.occurredAt).getTime())

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list: filteredList
    }
  })
}

export const createQuestionOccurrence = async(req: Request, res: Response) => {
  const question = getInterviewQuestionByIdInternal(req.body.questionId)
  const { getApplicationByIdInternal } = await import('./resume-applications')
  const application = getApplicationByIdInternal(req.body.applicationId)

  if (!question) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  if (!application) {
    return res.status(404).json({
      code: 50004,
      message: '投递记录不存在'
    })
  }

  const duplicatedOccurrence = req.body.interviewProgressId
    ? occurrenceList.find(item =>
      item.questionId === req.body.questionId &&
      item.applicationId === req.body.applicationId &&
      item.interviewProgressId === req.body.interviewProgressId
    )
    : undefined

  if (duplicatedOccurrence) {
    return res.status(400).json({
      code: 50008,
      message: '当前面试进展已经关联过这道题了'
    })
  }

  const now = new Date().toISOString()
  const newOccurrence: QuestionOccurrence = {
    id: getNextId(occurrenceList),
    questionId: req.body.questionId,
    applicationId: req.body.applicationId,
    interviewProgressId: req.body.interviewProgressId,
    companyNameSnapshot: application.companyName,
    jobTitleSnapshot: application.jobTitle,
    interviewStageSnapshot: req.body.interviewStageSnapshot,
    occurredAt: req.body.occurredAt,
    actualQuestion: req.body.actualQuestion,
    answerPerformance: req.body.answerPerformance,
    note: req.body.note,
    createdAt: now,
    updatedAt: now
  }

  occurrenceList.unshift(newOccurrence)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newOccurrence
  })
}

export const deleteQuestionOccurrence = (req: Request, res: Response) => {
  const index = occurrenceList.findIndex(item => item.id === req.params.id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '题目出现记录不存在'
    })
  }

  occurrenceList.splice(index, 1)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}

const getWeakPriority = (status: QuestionMasteryStatus) => {
  const priorityMap: Record<QuestionMasteryStatus, number> = {
    [QuestionMasteryStatus.UNREVIEWED]: 0,
    [QuestionMasteryStatus.WEAK]: 1,
    [QuestionMasteryStatus.NORMAL]: 2,
    [QuestionMasteryStatus.MASTERED]: 3
  }

  return priorityMap[status]
}

export const getFlashcards = (req: Request, res: Response) => {
  const { mode, tags, difficulty, limit = 8 } = req.query
  const selectedTags = parseTags(tags)
  const limitSize = Math.max(1, Number(limit))

  let resultList = questionList.map(buildQuestion)

  if (selectedTags.length) {
    resultList = resultList.filter(item => matchTags(item.tags, selectedTags))
  }

  if (difficulty) {
    resultList = resultList.filter(item => item.difficulty === difficulty)
  }

  switch (mode) {
    case FlashcardReviewMode.TAG:
      resultList = shuffleList(resultList)
      break
    case FlashcardReviewMode.DIFFICULTY:
      resultList = resultList.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
      break
    case FlashcardReviewMode.FAVORITE_FIRST:
      resultList = resultList.sort((a, b) => {
        if (a.isFavorite === b.isFavorite) return b.occurrenceCount - a.occurrenceCount
        return a.isFavorite ? -1 : 1
      })
      break
    case FlashcardReviewMode.WEAK_FIRST:
      resultList = resultList.sort((a, b) => {
        const priorityDiff = getWeakPriority(a.masteryStatus) - getWeakPriority(b.masteryStatus)
        if (priorityDiff !== 0) return priorityDiff
        return b.occurrenceCount - a.occurrenceCount
      })
      break
    case FlashcardReviewMode.FREQUENT_FIRST:
      resultList = resultList.sort((a, b) => {
        if (b.occurrenceCount !== a.occurrenceCount) return b.occurrenceCount - a.occurrenceCount
        return Number(b.isFavorite) - Number(a.isFavorite)
      })
      break
    case FlashcardReviewMode.RANDOM:
    default:
      resultList = shuffleList(resultList)
      break
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list: resultList.slice(0, limitSize),
      total: resultList.length,
      mode
    }
  })
}

export const createFlashcardReview = (req: Request, res: Response) => {
  const question = getInterviewQuestionByIdInternal(req.body.questionId)

  if (!question) {
    return res.status(404).json({
      code: 50004,
      message: '题目不存在'
    })
  }

  const reviewedAt = req.body.reviewedAt || new Date().toISOString()
  const now = new Date().toISOString()
  const review: FlashcardReview = {
    id: getNextId(reviewList),
    questionId: req.body.questionId,
    reviewedAt,
    result: req.body.result,
    reviewMode: req.body.reviewMode,
    createdAt: now
  }

  reviewList.unshift(review)
  syncQuestionReviewState(question, req.body.result, reviewedAt)

  return res.json({
    code: 20000,
    message: '保存成功',
    data: {
      review,
      updatedQuestion: buildQuestion(question),
      reviewSummary: buildReviewSummary(question.id)
    }
  })
}
