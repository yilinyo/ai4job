import request from '@/utils/request'
import {
  CreateResumeVersionRequest,
  UpdateResumeVersionRequest,
  ResumeVersionQueryParams
} from '@/models'

/**
 * 获取简历版本列表
 */
export const getResumeVersions = (params: ResumeVersionQueryParams) =>
  request({
    url: '/api/resume-versions',
    method: 'get',
    params
  })

/**
 * 获取简历版本详情
 */
export const getResumeVersionById = (id: string) =>
  request({
    url: `/api/resume-versions/${id}`,
    method: 'get'
  })

/**
 * 创建简历版本
 */
export const createResumeVersion = (data: CreateResumeVersionRequest) =>
  request({
    url: '/api/resume-versions',
    method: 'post',
    data
  })

/**
 * 更新简历版本
 */
export const updateResumeVersion = (id: string, data: UpdateResumeVersionRequest) =>
  request({
    url: `/api/resume-versions/${id}`,
    method: 'put',
    data
  })

/**
 * 删除简历版本
 */
export const deleteResumeVersion = (id: string) =>
  request({
    url: `/api/resume-versions/${id}`,
    method: 'delete'
  })

/**
 * 激活指定版本(同时取消同仓库下其他版本的激活状态)
 */
export const activateResumeVersion = (id: string) =>
  request({
    url: `/api/resume-versions/${id}/activate`,
    method: 'post'
  })

/**
 * 版本对比(预留功能)
 */
export const compareVersions = (oldVersionId: string, newVersionId: string) =>
  request({
    url: '/api/resume-versions/compare',
    method: 'get',
    params: {
      oldVersionId,
      newVersionId
    }
  })
