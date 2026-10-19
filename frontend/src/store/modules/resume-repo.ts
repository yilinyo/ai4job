import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  getResumeRepos,
  getResumeRepoById,
  createResumeRepo,
  updateResumeRepo,
  deleteResumeRepo
} from '@/api/resume-repos'
import {
  ResumeRepo,
  CreateResumeRepoRequest,
  UpdateResumeRepoRequest,
  ResumeRepoQueryParams
} from '@/models'
import store from '@/store'

export interface IResumeRepoState {
  repos: ResumeRepo[]
  currentRepo: ResumeRepo | null
  total: number
  loading: boolean
}

@Module({ dynamic: true, store, name: 'resumeRepo' })
class ResumeRepoStore extends VuexModule implements IResumeRepoState {
  public repos: ResumeRepo[] = []
  public currentRepo: ResumeRepo | null = null
  public total = 0
  public loading = false

  @Mutation
  private SET_REPOS(repos: ResumeRepo[]) {
    this.repos = repos
  }

  @Mutation
  private SET_CURRENT_REPO(repo: ResumeRepo | null) {
    this.currentRepo = repo
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
  public async GetResumeRepos(params: ResumeRepoQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeRepos(params)
      // 兼容后端返回空数组或缺少 list 的响应，避免将 undefined 写入列表状态，
      // 从而导致页面的空态条件无法渲染。
      const payload = (data || {}) as unknown as { list?: ResumeRepo[], total?: number }
      const repos: ResumeRepo[] = Array.isArray(data)
        ? (data as unknown as ResumeRepo[])
        : Array.isArray(payload.list)
          ? payload.list
          : []
      const total: number = typeof payload.total === 'number' ? payload.total : repos.length

      this.SET_REPOS(repos)
      this.SET_TOTAL(total)
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetResumeRepoById(id: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await getResumeRepoById(id)
      if (data) {
        this.SET_CURRENT_REPO(data)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateResumeRepo(repoData: CreateResumeRepoRequest) {
    const { data } = await createResumeRepo(repoData)
    if (data) {
      await this.GetResumeRepos({ page: 1, pageSize: 10 })
      return data
    }
  }

  @Action
  public async UpdateResumeRepo(payload: { id: string, data: UpdateResumeRepoRequest }) {
    const { data } = await updateResumeRepo(payload.id, payload.data)
    if (data) {
      await this.GetResumeRepos({ page: 1, pageSize: 10 })
      return data
    }
  }

  @Action
  public async DeleteResumeRepo(id: string) {
    await deleteResumeRepo(id)
    // 删除成功后刷新列表
    await this.GetResumeRepos({ page: 1, pageSize: 10 })
  }

  @Action
  public ClearCurrentRepo() {
    this.SET_CURRENT_REPO(null)
  }
}

export const ResumeRepoModule = getModule(ResumeRepoStore)
