import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  getInterviewProgressList,
  createInterviewProgress,
  updateInterviewProgress,
  deleteInterviewProgress
} from '@/api/interview-progress'
import {
  InterviewProgress,
  InterviewProgressQueryParams,
  CreateInterviewProgressRequest,
  UpdateInterviewProgressRequest
} from '@/models'
import store from '@/store'

export interface IInterviewProgressState {
  progressList: InterviewProgress[]
  loading: boolean
}

@Module({ dynamic: true, store, name: 'interviewProgress' })
class InterviewProgressStore extends VuexModule implements IInterviewProgressState {
  public progressList: InterviewProgress[] = []
  public loading = false

  @Mutation
  private SET_PROGRESS_LIST(progressList: InterviewProgress[]) {
    this.progressList = progressList
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Action
  public async GetInterviewProgressList(params: InterviewProgressQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getInterviewProgressList(params)
      if (data) {
        this.SET_PROGRESS_LIST(data.list)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateInterviewProgress(payload: CreateInterviewProgressRequest) {
    const { data } = await createInterviewProgress(payload)
    if (data) {
      await this.GetInterviewProgressList({ applicationId: payload.applicationId })
      return data
    }
  }

  @Action
  public async UpdateInterviewProgress(payload: { id: string, data: UpdateInterviewProgressRequest, applicationId: string }) {
    const { data } = await updateInterviewProgress(payload.id, payload.data)
    if (data) {
      await this.GetInterviewProgressList({ applicationId: payload.applicationId })
      return data
    }
  }

  @Action
  public async DeleteInterviewProgress(payload: { id: string, applicationId: string }) {
    await deleteInterviewProgress(payload.id)
    await this.GetInterviewProgressList({ applicationId: payload.applicationId })
  }

  @Action
  public ClearInterviewProgress() {
    this.SET_PROGRESS_LIST([])
  }
}

export const InterviewProgressModule = getModule(InterviewProgressStore)
