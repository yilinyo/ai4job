import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  createFlashcardReview,
  getFlashcards
} from '@/api'
import {
  CreateFlashcardReviewRequest,
  FlashcardQueryParams,
  FlashcardReviewMode,
  FlashcardReviewResult,
  InterviewQuestion,
  InterviewQuestionDifficulty
} from '@/models'
import store from '@/store'
import { InterviewQuestionBankModule } from './interview-question-bank'

interface FlashcardReviewStats {
  total: number
  failedCount: number
  vagueCount: number
  masteredCount: number
}

export interface IFlashcardReviewState {
  queue: InterviewQuestion[]
  currentIndex: number
  loading: boolean
  saving: boolean
  reviewMode: FlashcardReviewMode
  selectedTags: string[]
  selectedDifficulty?: InterviewQuestionDifficulty
  totalAvailable: number
  results: Record<string, FlashcardReviewResult>
  summary: FlashcardReviewStats
}

const createEmptySummary = (): FlashcardReviewStats => ({
  total: 0,
  failedCount: 0,
  vagueCount: 0,
  masteredCount: 0
})

@Module({ dynamic: true, store, name: 'flashcardReview' })
class FlashcardReviewStore extends VuexModule implements IFlashcardReviewState {
  public queue: InterviewQuestion[] = []
  public currentIndex = 0
  public loading = false
  public saving = false
  public reviewMode = FlashcardReviewMode.RANDOM
  public selectedTags: string[] = []
  public selectedDifficulty: InterviewQuestionDifficulty | undefined = undefined
  public totalAvailable = 0
  public results: Record<string, FlashcardReviewResult> = {}
  public summary: FlashcardReviewStats = createEmptySummary()

  get currentQuestion() {
    return this.queue[this.currentIndex] || null
  }

  get hasPrevious() {
    return this.currentIndex > 0
  }

  get hasNext() {
    return this.currentIndex < this.queue.length - 1
  }

  get isFinished() {
    return this.queue.length > 0 && Object.keys(this.results).length >= this.queue.length
  }

  @Mutation
  private SET_QUEUE(payload: { list: InterviewQuestion[], total: number, mode: FlashcardReviewMode }) {
    this.queue = payload.list
    this.totalAvailable = payload.total
    this.reviewMode = payload.mode
    this.currentIndex = 0
    this.results = {}
    this.summary = createEmptySummary()
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Mutation
  private SET_SAVING(saving: boolean) {
    this.saving = saving
  }

  @Mutation
  private SET_CURRENT_INDEX(index: number) {
    this.currentIndex = index
  }

  @Mutation
  private SET_REVIEW_MODE(mode: FlashcardReviewMode) {
    this.reviewMode = mode
  }

  @Mutation
  private SET_SELECTED_TAGS(tags: string[]) {
    this.selectedTags = tags
  }

  @Mutation
  private SET_SELECTED_DIFFICULTY(difficulty?: InterviewQuestionDifficulty) {
    this.selectedDifficulty = difficulty
  }

  @Mutation
  private ADD_RESULT(payload: { questionId: string, result: FlashcardReviewResult }) {
    this.results = {
      ...this.results,
      [payload.questionId]: payload.result
    }

    const values = Object.values(this.results)
    this.summary = {
      total: values.length,
      failedCount: values.filter(item => item === FlashcardReviewResult.FAILED).length,
      vagueCount: values.filter(item => item === FlashcardReviewResult.VAGUE).length,
      masteredCount: values.filter(item => item === FlashcardReviewResult.MASTERED).length
    }
  }

  @Mutation
  private UPDATE_QUESTION(question: InterviewQuestion) {
    const index = this.queue.findIndex(item => item.id === question.id)
    if (index === -1) return

    this.queue.splice(index, 1, question)
    this.queue = [...this.queue]
  }

  @Action
  public async LoadFlashcards(payload?: Partial<FlashcardQueryParams> & { limit?: number }) {
    const params: FlashcardQueryParams = {
      mode: payload?.mode || this.reviewMode,
      tags: payload?.tags ?? (this.selectedTags.length ? this.selectedTags.join(',') : undefined),
      difficulty: payload?.difficulty ?? this.selectedDifficulty,
      limit: payload?.limit || 8
    }

    this.SET_REVIEW_MODE(params.mode)
    this.SET_LOADING(true)
    try {
      const { data } = await getFlashcards(params)
      if (data) {
        this.SET_QUEUE({
          list: data.list,
          total: data.total,
          mode: data.mode
        })
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async SaveReview(payload: CreateFlashcardReviewRequest) {
    this.SET_SAVING(true)
    try {
      const { data } = await createFlashcardReview(payload)
      if (data) {
        this.ADD_RESULT({
          questionId: payload.questionId,
          result: payload.result
        })
        this.UPDATE_QUESTION(data.updatedQuestion)
        InterviewQuestionBankModule.GetInterviewQuestions()
        InterviewQuestionBankModule.GetInterviewQuestionById(payload.questionId).catch(error => {
          console.warn('刷新题目详情失败:', error)
        })
        return data
      }
    } finally {
      this.SET_SAVING(false)
    }
  }

  @Action
  public SetCurrentIndex(index: number) {
    this.SET_CURRENT_INDEX(index)
  }

  @Action
  public SetReviewMode(mode: FlashcardReviewMode) {
    this.SET_REVIEW_MODE(mode)
  }

  @Action
  public SetSelectedTags(tags: string[]) {
    this.SET_SELECTED_TAGS(tags)
  }

  @Action
  public SetSelectedDifficulty(difficulty?: InterviewQuestionDifficulty) {
    this.SET_SELECTED_DIFFICULTY(difficulty)
  }

  @Action
  public ResetSession() {
    this.SET_QUEUE({
      list: [],
      total: 0,
      mode: this.reviewMode
    })
  }
}

export const FlashcardReviewModule = getModule(FlashcardReviewStore)
