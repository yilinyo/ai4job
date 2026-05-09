import { shallowMount } from '@vue/test-utils'
import { readFileSync } from 'fs'
import {
  FlashcardReviewMode,
  FlashcardReviewResult,
  InterviewQuestionDifficulty,
  InterviewQuestionType,
  InterviewResult,
  InterviewStage,
  QuestionMasteryStatus
} from '@/models'

const mockQuestionBankModule = {
  questions: [
    {
      id: 'q1',
      title: 'Vue 响应式',
      content: '解释 Vue 2 响应式',
      answerAnalysis: 'Object.defineProperty',
      tags: ['Vue'],
      difficulty: InterviewQuestionDifficulty.MEDIUM,
      questionType: InterviewQuestionType.KNOWLEDGE,
      masteryStatus: QuestionMasteryStatus.NORMAL,
      isFavorite: false,
      occurrenceCount: 2,
      createdAt: '2026-05-09T00:00:00.000Z',
      updatedAt: '2026-05-09T00:00:00.000Z'
    }
  ],
  total: 1,
  loading: false,
  submitting: false,
  GetInterviewQuestions: jest.fn().mockResolvedValue(undefined),
  CreateInterviewQuestion: jest.fn().mockResolvedValue({
    id: 'created-question',
    title: '新题目'
  }),
  UpdateInterviewQuestion: jest.fn().mockResolvedValue(undefined),
  DeleteInterviewQuestion: jest.fn().mockResolvedValue(undefined),
  ToggleFavorite: jest.fn().mockResolvedValue(undefined),
  CreateQuestionOccurrence: jest.fn().mockResolvedValue(undefined)
}

const mockFlashcardModule = {
  queue: [
    {
      id: 'q1',
      title: 'Vue 响应式',
      content: '解释 Vue 2 响应式',
      answerAnalysis: 'Object.defineProperty',
      tags: ['Vue'],
      difficulty: InterviewQuestionDifficulty.MEDIUM,
      questionType: InterviewQuestionType.KNOWLEDGE,
      masteryStatus: QuestionMasteryStatus.NORMAL,
      isFavorite: false,
      occurrenceCount: 2,
      createdAt: '2026-05-09T00:00:00.000Z',
      updatedAt: '2026-05-09T00:00:00.000Z'
    },
    {
      id: 'q2',
      title: '性能优化',
      content: '解释性能优化体系',
      answerAnalysis: '网络、构建、运行时',
      tags: ['工程化'],
      difficulty: InterviewQuestionDifficulty.HARD,
      questionType: InterviewQuestionType.SCENARIO,
      masteryStatus: QuestionMasteryStatus.WEAK,
      isFavorite: true,
      occurrenceCount: 1,
      createdAt: '2026-05-09T00:00:00.000Z',
      updatedAt: '2026-05-09T00:00:00.000Z'
    }
  ],
  currentIndex: 0,
  loading: false,
  saving: false,
  summary: {
    total: 0,
    failedCount: 0,
    vagueCount: 0,
    masteredCount: 0
  },
  get currentQuestion() {
    return this.queue[this.currentIndex]
  },
  get hasPrevious() {
    return this.currentIndex > 0
  },
  get hasNext() {
    return this.currentIndex < this.queue.length - 1
  },
  SetReviewMode: jest.fn().mockResolvedValue(undefined),
  SetSelectedTags: jest.fn().mockResolvedValue(undefined),
  SetSelectedDifficulty: jest.fn().mockResolvedValue(undefined),
  LoadFlashcards: jest.fn().mockResolvedValue(undefined),
  SetCurrentIndex: jest.fn((index: number) => {
    mockFlashcardModule.currentIndex = index
    return Promise.resolve()
  }),
  SaveReview: jest.fn().mockImplementation(() => {
    mockFlashcardModule.summary.total += 1
    mockFlashcardModule.summary.masteredCount += 1
    return Promise.resolve()
  })
}

const mockResumeApplicationModule = {
  currentApplication: {
    id: 'app1',
    companyName: '腾讯科技',
    jobTitle: '前端开发工程师',
    appliedAt: '2026-05-09T00:00:00.000Z',
    currentStatus: 'interviewing',
    interviewNotes: '请解释 Vue 响应式'
  },
  loading: false,
  GetResumeApplicationById: jest.fn().mockResolvedValue(undefined),
  GetResumeApplications: jest.fn().mockResolvedValue(undefined)
}

const mockInterviewProgressModule = {
  progressList: [
    {
      id: 'p1',
      applicationId: 'app1',
      stage: InterviewStage.FIRST_INTERVIEW,
      occurredAt: '2026-05-09T10:00:00.000Z',
      result: InterviewResult.PASSED,
      note: '一面问了 Vue 响应式'
    }
  ],
  loading: false,
  GetInterviewProgressList: jest.fn().mockResolvedValue(undefined)
}

jest.doMock('@/store/modules/interview-question-bank', () => ({
  InterviewQuestionBankModule: mockQuestionBankModule
}))

jest.doMock('@/store/modules/flashcard-review', () => ({
  FlashcardReviewModule: mockFlashcardModule
}))

jest.doMock('@/store/modules/resume-application', () => ({
  ResumeApplicationModule: mockResumeApplicationModule
}))

jest.doMock('@/store/modules/interview-progress', () => ({
  InterviewProgressModule: mockInterviewProgressModule
}))

const QuestionLibraryTab = require('@/views/interview-question-bank/question-library-tab.vue').default
const FlashcardReviewTab = require('@/views/interview-question-bank/flashcard-review-tab.vue').default
const ResumeApplicationView = require('@/views/resume-application/view.vue').default

describe('Interview question bank validation', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockFlashcardModule.currentIndex = 0
    mockFlashcardModule.summary.total = 0
    mockFlashcardModule.summary.failedCount = 0
    mockFlashcardModule.summary.vagueCount = 0
    mockFlashcardModule.summary.masteredCount = 0
  })

  it('runs the question library flow hooks for filtering, favorite, detail, edit, delete and occurrence creation', async() => {
    const confirm = jest.fn().mockResolvedValue(undefined)
    const push = jest.fn()
    const wrapper = shallowMount(QuestionLibraryTab, {
      mocks: {
        $confirm: confirm,
        $message: {
          success: jest.fn(),
          error: jest.fn()
        },
        $router: { push }
      }
    })
    const vm = wrapper.vm as any

    vm.selectedTags = ['Vue']
    vm.favoriteFilter = 'true'
    vm.queryParams.keyword = '响应式'
    await vm.handleSearch()
    expect(mockQuestionBankModule.GetInterviewQuestions).toHaveBeenLastCalledWith(expect.objectContaining({
      keyword: '响应式',
      tags: 'Vue',
      isFavorite: true,
      page: 1
    }))

    vm.openCreateDialog()
    expect(vm.questionDialogVisible).toBe(true)
    await vm.submitQuestion({
      title: '新题目',
      content: '题目内容',
      answerAnalysis: '答案',
      tags: ['Vue'],
      difficulty: InterviewQuestionDifficulty.MEDIUM,
      questionType: InterviewQuestionType.KNOWLEDGE,
      masteryStatus: QuestionMasteryStatus.UNREVIEWED,
      isFavorite: false
    })
    expect(mockQuestionBankModule.CreateInterviewQuestion).toHaveBeenCalled()

    await vm.toggleFavorite(mockQuestionBankModule.questions[0])
    expect(mockQuestionBankModule.ToggleFavorite).toHaveBeenCalledWith({
      id: 'q1',
      isFavorite: true
    })

    vm.viewQuestion('q1')
    expect(push).toHaveBeenCalledWith('/interview-question-bank/question/q1')

    vm.handleRowClick(mockQuestionBankModule.questions[0])
    expect(push).toHaveBeenLastCalledWith('/interview-question-bank/question/q1')

    vm.openEditDialog(mockQuestionBankModule.questions[0])
    expect(vm.questionDialogTitle).toBe('编辑题目')

    vm.openOccurrenceDialog('q1')
    expect(vm.occurrenceDialogVisible).toBe(true)
    await vm.submitOccurrence({
      questionId: 'q1',
      applicationId: 'app1',
      occurredAt: '2026-05-09T10:00:00'
    })
    expect(mockQuestionBankModule.CreateQuestionOccurrence).toHaveBeenCalledWith(expect.objectContaining({
      questionId: 'q1',
      applicationId: 'app1'
    }))

    await vm.deleteQuestion(mockQuestionBankModule.questions[0])
    expect(confirm).toHaveBeenCalledWith(expect.stringContaining('同步清理 2 条面试提问记录'), '删除题目', expect.any(Object))
    expect(mockQuestionBankModule.DeleteInterviewQuestion).toHaveBeenCalledWith('q1')
  })

  it('runs the resume review sediment flow for creating a new question and linking an existing question', async() => {
    const wrapper = shallowMount(ResumeApplicationView, {
      mocks: {
        $route: { params: { id: 'app1' } },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          error: jest.fn()
        },
        $confirm: jest.fn()
      }
    })
    const vm = wrapper.vm as any
    const progress = mockInterviewProgressModule.progressList[0]

    vm.openQuestionCreateDialog(progress)
    expect(vm.questionDialogVisible).toBe(true)
    expect(vm.questionDraft.content).toBe(progress.note)
    expect(vm.linkProgressId).toBe('p1')

    await vm.submitQuestionAndOccurrence({
      title: '沉淀题目',
      content: '真实问法',
      answerAnalysis: '答案解析',
      tags: ['面试复盘'],
      difficulty: InterviewQuestionDifficulty.MEDIUM,
      questionType: InterviewQuestionType.SCENARIO,
      masteryStatus: QuestionMasteryStatus.UNREVIEWED,
      isFavorite: false
    })
    expect(mockQuestionBankModule.CreateInterviewQuestion).toHaveBeenCalled()
    expect(mockQuestionBankModule.CreateQuestionOccurrence).toHaveBeenCalledWith(expect.objectContaining({
      questionId: 'created-question',
      applicationId: 'app1',
      interviewProgressId: 'p1',
      interviewStageSnapshot: InterviewStage.FIRST_INTERVIEW,
      actualQuestion: '真实问法'
    }))

    vm.openQuestionLinkDialog(progress)
    expect(vm.questionLinkDialogVisible).toBe(true)
    await vm.submitQuestionOccurrence({
      questionId: 'q1',
      applicationId: 'app1',
      interviewProgressId: 'p1',
      occurredAt: '2026-05-09T10:00:00'
    })
    expect(mockQuestionBankModule.CreateQuestionOccurrence).toHaveBeenLastCalledWith(expect.objectContaining({
      questionId: 'q1',
      applicationId: 'app1',
      interviewProgressId: 'p1'
    }))
  })

  it('runs the flashcard review flow and guards repeated animation triggers', async() => {
    jest.useFakeTimers()
    const wrapper = shallowMount(FlashcardReviewTab, {
      mocks: {
        $message: {
          success: jest.fn()
        }
      }
    })
    const vm = wrapper.vm as any

    vm.reviewMode = FlashcardReviewMode.TAG
    vm.selectedTags = ['Vue']
    await vm.loadFlashcards()
    expect(mockFlashcardModule.LoadFlashcards).toHaveBeenCalledWith(expect.objectContaining({
      mode: FlashcardReviewMode.TAG,
      tags: 'Vue'
    }))

    vm.flipCard()
    expect(vm.isFlipped).toBe(true)

    vm.moveCard(1)
    expect(vm.isAnimating).toBe(true)
    vm.moveCard(1)
    expect(mockFlashcardModule.SetCurrentIndex).not.toHaveBeenCalled()

    jest.advanceTimersByTime(260)
    expect(mockFlashcardModule.SetCurrentIndex).toHaveBeenCalledWith(1)
    expect(vm.isFlipped).toBe(false)
    jest.advanceTimersByTime(60)
    expect(vm.isAnimating).toBe(false)

    await vm.saveReview(FlashcardReviewResult.MASTERED)
    expect(mockFlashcardModule.SaveReview).toHaveBeenCalledWith({
      questionId: 'q1',
      result: FlashcardReviewResult.MASTERED,
      reviewMode: FlashcardReviewMode.TAG
    })
    expect(mockFlashcardModule.summary.total).toBe(1)
    jest.useRealTimers()
  })

  it('keeps the visual and interaction contract visible in rendered markup', () => {
    const library = shallowMount(QuestionLibraryTab, {
      mocks: {
        $router: { push: jest.fn() }
      }
    })
    const flashcard = shallowMount(FlashcardReviewTab)

    const librarySource = readFileSync('src/views/interview-question-bank/question-library-tab.vue', 'utf8')
    const detailSource = readFileSync('src/views/interview-question-bank/detail.vue', 'utf8')
    const formSource = readFileSync('src/components/InterviewQuestion/QuestionFormDialog.vue', 'utf8')
    const indexSource = readFileSync('src/views/interview-question-bank/index.vue', 'utf8')
    const routerSource = readFileSync('src/router/modules/resume.ts', 'utf8')
    const mockSource = readFileSync('mock/interview-questions.ts', 'utf8')

    expect(library.text()).toContain('新增题目')
    expect(library.html()).toContain('label="提问"')
    expect(library.html()).toContain('label="创建时间"')
    expect(library.html()).not.toContain('label="修改时间"')
    expect(indexSource).not.toContain('<el-tabs')
    expect(routerSource).toContain("redirect: '/interview-question-bank/list'")
    expect(routerSource).toContain('alwaysShow: true')
    expect(routerSource).toContain("path: 'list'")
    expect(routerSource).toContain("path: 'flashcard'")
    expect(librarySource).toContain('@row-click="handleRowClick"')
    expect(librarySource).toContain('@click.stop="viewQuestion')
    expect(librarySource).toContain('el-icon-star')
    expect(librarySource).toContain('删除')
    expect(librarySource).toContain('关联提问')
    expect(detailSource).toContain('@dblclick.native="openEditDialog"')
    expect(detailSource).toContain('创建时间：{{ formatDateTime(question.createdAt) }}')
    expect(detailSource).toContain('修改时间：{{ formatDateTime(question.updatedAt) }}')
    expect(formSource).toContain('@paste.native="handleMarkdownPaste')
    expect(mockSource).toContain('createdAt: now')
    expect(mockSource).toContain('updatedAt: now')
    expect(mockSource).toContain('updatedAt: new Date().toISOString()')

    expect(flashcard.text()).toContain('开始复习')
    expect(readFileSync('src/views/interview-question-bank/flashcard-review-tab.vue', 'utf8')).toContain('label="复习模式"')
    expect(flashcard.html()).toContain('flashcard-shell')
    expect(flashcard.html()).toContain('flashcard-front')
    expect(flashcard.html()).toContain('flashcard-back')
    expect(flashcard.html()).toContain('card-controls')
  })
})
