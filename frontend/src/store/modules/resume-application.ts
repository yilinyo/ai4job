import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  getResumeApplications,
  getResumeApplicationById,
  createResumeApplication,
  updateResumeApplication,
  deleteResumeApplication
} from '@/api/resume-applications'
import {
  ResumeApplication,
  ResumeApplicationQueryParams,
  CreateResumeApplicationRequest,
  UpdateResumeApplicationRequest
} from '@/models'
import store from '@/store'

export interface IResumeApplicationState {
  applications: ResumeApplication[]
  currentApplication: ResumeApplication | null
  total: number
  loading: boolean
}

@Module({ dynamic: true, store, name: 'resumeApplication' })
class ResumeApplicationStore extends VuexModule implements IResumeApplicationState {
  public applications: ResumeApplication[] = []
  public currentApplication: ResumeApplication | null = null
  public total = 0
  public loading = false

  @Mutation
  private SET_APPLICATIONS(applications: ResumeApplication[]) {
    this.applications = applications
  }

  @Mutation
  private SET_CURRENT_APPLICATION(application: ResumeApplication | null) {
    this.currentApplication = application
  }

  @Mutation
  private SET_TOTAL(total: number) {
    this.total = total
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Action
  public async GetResumeApplications(params: ResumeApplicationQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeApplications(params)
      if (data) {
        this.SET_APPLICATIONS(data.list)
        this.SET_TOTAL(data.total)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetResumeApplicationById(id: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeApplicationById(id)
      if (data) {
        this.SET_CURRENT_APPLICATION(data)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateResumeApplication(payload: CreateResumeApplicationRequest) {
    const { data } = await createResumeApplication(payload)
    if (data) {
      await this.GetResumeApplications({
        resumeVersionId: payload.resumeVersionId,
        page: 1,
        pageSize: 100
      })
      return data
    }
  }

  @Action
  public async UpdateResumeApplication(payload: { id: string, data: UpdateResumeApplicationRequest }) {
    const { data } = await updateResumeApplication(payload.id, payload.data)
    if (data) {
      if (this.currentApplication && this.currentApplication.id === payload.id) {
        this.SET_CURRENT_APPLICATION(data)
      }
      return data
    }
  }

  @Action
  public async DeleteResumeApplication(payload: { id: string, resumeVersionId: string }) {
    await deleteResumeApplication(payload.id)
    await this.GetResumeApplications({
      resumeVersionId: payload.resumeVersionId,
      page: 1,
      pageSize: 100
    })
  }

  @Action
  public ClearCurrentApplication() {
    this.SET_CURRENT_APPLICATION(null)
  }
}

export const ResumeApplicationModule = getModule(ResumeApplicationStore)
