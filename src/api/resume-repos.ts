import request from '@/utils/request'
import {
  CreateResumeRepoRequest,
  UpdateResumeRepoRequest,
  ResumeRepoQueryParams
} from '@/models'

/**
 * 获取简历仓库列表
 */
export const getResumeRepos = (params: ResumeRepoQueryParams) =>
  request({
    url: '/api/resume-repos',
    method: 'get',
    params
  })

/**
 * 获取简历仓库详情
 */
export const getResumeRepoById = (id: string) =>
  request({
    url: `/api/resume-repos/${id}`,
    method: 'get'
  })

/**
 * 创建简历仓库
 */
export const createResumeRepo = (data: CreateResumeRepoRequest) =>
  request({
    url: '/api/resume-repos',
    method: 'post',
    data
  })

/**
 * 更新简历仓库
 */
export const updateResumeRepo = (id: string, data: UpdateResumeRepoRequest) =>
  request({
    url: `/api/resume-repos/${id}`,
    method: 'put',
    data
  })

/**
 * 删除简历仓库
 */
export const deleteResumeRepo = (id: string) =>
  request({
    url: `/api/resume-repos/${id}`,
    method: 'delete'
  })
