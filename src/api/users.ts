import request from '@/utils/request'
import {
  RegisterRequest,
  LoginRequest,
  UpdateProfileRequest
} from '@/models'

/**
 * 用户注册
 */
export const register = (data: RegisterRequest) =>
  request({
    url: '/api/user/register',
    method: 'post',
    data
  })

/**
 * 用户登录
 */
export const login = (data: LoginRequest) =>
  request({
    url: '/api/user/login',
    method: 'post',
    data
  })

/**
 * 用户登出
 */
export const logout = () =>
  request({
    url: '/api/user/logout',
    method: 'post'
  })

/**
 * 获取用户个人信息
 */
export const getUserProfile = () =>
  request({
    url: '/api/user/profile',
    method: 'get'
  })

/**
 * 更新用户个人信息
 */
export const updateUserProfile = (data: UpdateProfileRequest) =>
  request({
    url: '/api/user/profile',
    method: 'put',
    data
  })

/**
 * 获取用户列表（兼容原模板 ArticleDetail 组件）
 */
export const getUsers = (params: any) =>
  request({
    url: '/users',
    method: 'get',
    params
  })
