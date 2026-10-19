import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  getResumeVersions,
  getResumeVersionById,
  createResumeVersion,
  updateResumeVersion,
  deleteResumeVersion,
  activateResumeVersion
} from '@/api/resume-versions'
import {
  ResumeVersion,
  CreateResumeVersionRequest,
  UpdateResumeVersionRequest,
  ResumeVersionQueryParams
} from '@/models'
import store from '@/store'

export interface IResumeVersionState {
  versions: ResumeVersion[]
  currentVersion: ResumeVersion | null
  total: number
  loading: boolean
}

@Module({ dynamic: true, store, name: 'resumeVersion' })
class ResumeVersionStore extends VuexModule implements IResumeVersionState {
  public versions: ResumeVersion[] = []
  public currentVersion: ResumeVersion | null = null
  public total = 0
  public loading = false

  @Mutation
  private SET_VERSIONS(versions: ResumeVersion[]) {
    this.versions = versions
  }

  @Mutation
  private SET_CURRENT_VERSION(version: ResumeVersion | null) {
    this.currentVersion = version
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
  public async GetResumeVersions(params: ResumeVersionQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeVersions(params)
      // 后端分页响应为 { list, total }；无论响应是否为空，都必须覆盖旧状态，
      // 否则从创建页返回时会继续显示此前缓存的版本列表。
      const payload = (data || {}) as { list?: ResumeVersion[], total?: number }
      const versions = Array.isArray(payload.list) ? payload.list : []
      const total = typeof payload.total === 'number' ? payload.total : versions.length
      this.SET_VERSIONS(versions)
      this.SET_TOTAL(total)
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetResumeVersionById(id: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeVersionById(id)
      if (data) {
        this.SET_CURRENT_VERSION(data)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateResumeVersion(versionData: CreateResumeVersionRequest) {
    // 创建接口按约定返回 data: null；请求成功后仍必须刷新列表。
    const { data } = await createResumeVersion(versionData)
    await this.GetResumeVersions({
      repoId: versionData.repoId,
      page: 1,
      pageSize: 100
    })
    return data
  }

  @Action
  public async UpdateResumeVersion(payload: { id: string, data: UpdateResumeVersionRequest }) {
    const { data } = await updateResumeVersion(payload.id, payload.data)
    if (data) {
      if (this.currentVersion && (this.currentVersion as any).id === payload.id) {
        this.SET_CURRENT_VERSION(data)
      }
      return data
    }
  }

  @Action
  public async DeleteResumeVersion(payload: { id: string, repoId: string }) {
    await deleteResumeVersion(payload.id)
    // 删除成功后刷新列表
    await this.GetResumeVersions({
      repoId: payload.repoId,
      page: 1,
      pageSize: 10
    })
  }

  @Action
  public async ActivateVersion(payload: { id: string, repoId: string }) {
    await activateResumeVersion(payload.id)
    // 激活成功后刷新列表
    await this.GetResumeVersions({
      repoId: payload.repoId,
      page: 1,
      pageSize: 10
    })
  }

  @Action
  public ClearCurrentVersion() {
    this.SET_CURRENT_VERSION(null)
  }
}

export const ResumeVersionModule = getModule(ResumeVersionStore)
