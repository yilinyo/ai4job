import request from '@/utils/request'
import {
  CreateDashboardTodoRequest,
  DashboardQueryParams,
  DashboardTodoQueryParams,
  UpdateDashboardTodoRequest
} from '@/models'

export const getDashboard = (params: DashboardQueryParams) =>
  request({
    url: '/api/dashboard',
    method: 'get',
    params
  })

export const getDashboardTodos = (params?: DashboardTodoQueryParams) =>
  request({
    url: '/api/dashboard/todos',
    method: 'get',
    params
  })

export const createDashboardTodo = (data: CreateDashboardTodoRequest) =>
  request({
    url: '/api/dashboard/todos',
    method: 'post',
    data
  })

export const updateDashboardTodo = (id: string, data: UpdateDashboardTodoRequest) =>
  request({
    url: `/api/dashboard/todos/${id}`,
    method: 'put',
    data
  })

export const completeDashboardTodo = (id: string) =>
  request({
    url: `/api/dashboard/todos/${id}/complete`,
    method: 'put'
  })

export const deleteDashboardTodo = (id: string) =>
  request({
    url: `/api/dashboard/todos/${id}`,
    method: 'delete'
  })
