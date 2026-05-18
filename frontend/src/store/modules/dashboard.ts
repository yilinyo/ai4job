import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  completeDashboardTodo,
  createDashboardTodo,
  deleteDashboardTodo,
  getDashboard,
  getDashboardTodos,
  updateDashboardTodo
} from '@/api/dashboard'
import {
  CreateDashboardTodoRequest,
  DashboardRange,
  DashboardResponse,
  DashboardTodo,
  DashboardTodoQueryParams,
  UpdateDashboardTodoRequest
} from '@/models'
import store from '@/store'

export interface IDashboardState {
  dashboard: DashboardResponse | null
  todos: DashboardTodo[]
  range: DashboardRange
  loading: boolean
}

@Module({ dynamic: true, store, name: 'dashboard' })
class DashboardStore extends VuexModule implements IDashboardState {
  public dashboard: DashboardResponse | null = null
  public todos: DashboardTodo[] = []
  public range: DashboardRange = '30d'
  public loading = false

  @Mutation
  private SET_DASHBOARD(dashboard: DashboardResponse | null) {
    this.dashboard = dashboard
  }

  @Mutation
  private SET_TODOS(todos: DashboardTodo[]) {
    this.todos = todos
  }

  @Mutation
  private SET_RANGE(range: DashboardRange) {
    this.range = range
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Action
  public async GetDashboard(range?: DashboardRange) {
    const currentRange = range || this.range
    this.SET_LOADING(true)
    try {
      const { data } = await getDashboard({ range: currentRange })
      if (data) {
        this.SET_RANGE(data.range || currentRange)
        this.SET_DASHBOARD(data)
        this.SET_TODOS(data.todos || [])
      }
      return data
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetTodos(params?: DashboardTodoQueryParams) {
    const { data } = await getDashboardTodos(params)
    const list = data && data.list ? data.list : []
    this.SET_TODOS(list)
    return list
  }

  @Action
  public async CreateTodo(payload: CreateDashboardTodoRequest) {
    const { data } = await createDashboardTodo(payload)
    await this.GetDashboard(this.range)
    return data
  }

  @Action
  public async UpdateTodo(payload: { id: string, data: UpdateDashboardTodoRequest }) {
    const { data } = await updateDashboardTodo(payload.id, payload.data)
    await this.GetDashboard(this.range)
    return data
  }

  @Action
  public async CompleteTodo(id: string) {
    const { data } = await completeDashboardTodo(id)
    await this.GetDashboard(this.range)
    return data
  }

  @Action
  public async DeleteTodo(id: string) {
    await deleteDashboardTodo(id)
    await this.GetDashboard(this.range)
  }
}

export const DashboardModule = getModule(DashboardStore)
