import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  createInterviewQuestion,
  createQuestionOccurrence,
  deleteInterviewQuestion,
  deleteQuestionOccurrence,
  getInterviewQuestionById,
  getInterviewQuestions,
  getQuestionOccurrences,
  updateInterviewQuestion,
  updateInterviewQuestionFavorite
} from '@/api'
import {
  CreateInterviewQuestionRequest,
  CreateQuestionOccurrenceRequest,
  InterviewQuestion,
  InterviewQuestionDetail,
  InterviewQuestionQueryParams,
  QuestionOccurrence,
  UpdateInterviewQuestionRequest
} from '@/models'
import store from '@/store'

export interface IInterviewQuestionBankState {
  questions: InterviewQuestion[]
  currentQuestion: InterviewQuestionDetail | null
  occurrences: QuestionOccurrence[]
  total: number
  loading: boolean
  submitting: boolean
  filters: InterviewQuestionQueryParams
}

const defaultFilters = (): InterviewQuestionQueryParams => ({
  keyword: '',
  tags: '',
  difficulty: undefined,
  questionType: undefined,
  masteryStatus: undefined,
  isFavorite: undefined,
  page: 1,
  pageSize: 10
})

@Module({ dynamic: true, store, name: 'interviewQuestionBank' })
class InterviewQuestionBankStore extends VuexModule implements IInterviewQuestionBankState {
  public questions: InterviewQuestion[] = []
  public currentQuestion: InterviewQuestionDetail | null = null
  public occurrences: QuestionOccurrence[] = []
  public total = 0
  public loading = false
  public submitting = false
  public filters: InterviewQuestionQueryParams = defaultFilters()

  @Mutation
  private SET_QUESTIONS(questions: InterviewQuestion[]) {
    this.questions = questions
  }

  @Mutation
  private SET_CURRENT_QUESTION(question: InterviewQuestionDetail | null) {
    this.currentQuestion = question
  }

  @Mutation
  private SET_OCCURRENCES(occurrences: QuestionOccurrence[]) {
    this.occurrences = occurrences
  }

  @Mutation
  private SET_TOTAL(total: number) {
    this.total = total
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Mutation
  private SET_SUBMITTING(submitting: boolean) {
    this.submitting = submitting
  }

  @Mutation
  private SET_FILTERS(filters: InterviewQuestionQueryParams) {
    this.filters = {
      ...this.filters,
      ...filters
    }
  }

  @Mutation
  private RESET_FILTERS() {
    this.filters = defaultFilters()
  }

  @Mutation
  private UPSERT_QUESTION(question: InterviewQuestion) {
    const index = this.questions.findIndex(item => item.id === question.id)
    if (index === -1) {
      this.questions = [question, ...this.questions]
      return
    }

    this.questions.splice(index, 1, question)
    this.questions = [...this.questions]
  }

  @Mutation
  private REMOVE_QUESTION(questionId: string) {
    this.questions = this.questions.filter(item => item.id !== questionId)
  }

  @Mutation
  private REPLACE_CURRENT_QUESTION_BASE(question: InterviewQuestion) {
    if (!this.currentQuestion || this.currentQuestion.id !== question.id) return

    this.currentQuestion = {
      ...this.currentQuestion,
      ...question
    }
  }

  @Action
  public async GetInterviewQuestions(params?: InterviewQuestionQueryParams) {
    const nextFilters = params ? { ...this.filters, ...params } : this.filters
    this.SET_FILTERS(nextFilters)
    this.SET_LOADING(true)
    try {
      const { data } = await getInterviewQuestions(nextFilters)
      if (data) {
        this.SET_QUESTIONS(data.list)
        this.SET_TOTAL(data.total)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetInterviewQuestionById(id: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await getInterviewQuestionById(id)
      if (data) {
        this.SET_CURRENT_QUESTION(data)
        this.SET_OCCURRENCES(data.occurrences || [])
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetQuestionOccurrences(params: { questionId?: string, applicationId?: string }) {
    const { data } = await getQuestionOccurrences(params)
    if (data) {
      this.SET_OCCURRENCES(data.list)
      return data.list
    }
    return []
  }

  @Action
  public async CreateInterviewQuestion(payload: CreateInterviewQuestionRequest) {
    this.SET_SUBMITTING(true)
    try {
      const { data } = await createInterviewQuestion(payload)
      if (data) {
        await this.GetInterviewQuestions()
        return data
      }
    } finally {
      this.SET_SUBMITTING(false)
    }
  }

  @Action
  public async UpdateInterviewQuestion(payload: { id: string, data: UpdateInterviewQuestionRequest }) {
    this.SET_SUBMITTING(true)
    try {
      const { data } = await updateInterviewQuestion(payload.id, payload.data)
      if (data) {
        this.UPSERT_QUESTION(data)
        this.REPLACE_CURRENT_QUESTION_BASE(data)
        await this.GetInterviewQuestions()
        return data
      }
    } finally {
      this.SET_SUBMITTING(false)
    }
  }

  @Action
  public async ToggleFavorite(payload: { id: string, isFavorite: boolean }) {
    const { data } = await updateInterviewQuestionFavorite(payload.id, {
      isFavorite: payload.isFavorite
    })
    if (data) {
      this.UPSERT_QUESTION(data)
      this.REPLACE_CURRENT_QUESTION_BASE(data)
      return data
    }
  }

  @Action
  public async DeleteInterviewQuestion(id: string) {
    await deleteInterviewQuestion(id)
    this.REMOVE_QUESTION(id)
    if (this.currentQuestion && this.currentQuestion.id === id) {
      this.SET_CURRENT_QUESTION(null)
      this.SET_OCCURRENCES([])
    }
    await this.GetInterviewQuestions()
  }

  @Action
  public async CreateQuestionOccurrence(payload: CreateQuestionOccurrenceRequest) {
    this.SET_SUBMITTING(true)
    try {
      const { data } = await createQuestionOccurrence(payload)
      if (data) {
        await this.GetInterviewQuestions()
        if (this.currentQuestion && this.currentQuestion.id === payload.questionId) {
          await this.GetInterviewQuestionById(payload.questionId)
        }
        return data
      }
    } finally {
      this.SET_SUBMITTING(false)
    }
  }

  @Action
  public async DeleteQuestionOccurrence(payload: { id: string, questionId?: string, applicationId?: string }) {
    await deleteQuestionOccurrence(payload.id)
    if (payload.questionId) {
      await this.GetQuestionOccurrences({ questionId: payload.questionId })
      await this.GetInterviewQuestions()
      if (this.currentQuestion && this.currentQuestion.id === payload.questionId) {
        await this.GetInterviewQuestionById(payload.questionId)
      }
    } else if (payload.applicationId) {
      await this.GetQuestionOccurrences({ applicationId: payload.applicationId })
      await this.GetInterviewQuestions()
    }
  }

  @Action
  public SetFilters(filters: InterviewQuestionQueryParams) {
    this.SET_FILTERS(filters)
  }

  @Action
  public ResetFilters() {
    this.RESET_FILTERS()
  }

  @Action
  public ClearCurrentQuestion() {
    this.SET_CURRENT_QUESTION(null)
    this.SET_OCCURRENCES([])
  }
}

export const InterviewQuestionBankModule = getModule(InterviewQuestionBankStore)
