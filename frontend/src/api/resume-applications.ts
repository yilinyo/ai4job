import request from '@/utils/request'
import {
  CreateResumeApplicationRequest,
  UpdateResumeApplicationRequest,
  ResumeApplicationQueryParams
} from '@/models'

export const getResumeApplications = (params: ResumeApplicationQueryParams) =>
  request({
    url: '/api/resume-applications',
    method: 'get',
    params
  })

export const getResumeApplicationById = (id: string) =>
  request({
    url: `/api/resume-applications/${id}`,
    method: 'get'
  })

export const createResumeApplication = (data: CreateResumeApplicationRequest) =>
  request({
    url: '/api/resume-applications',
    method: 'post',
    data
  })

export const updateResumeApplication = (id: string, data: UpdateResumeApplicationRequest) =>
  request({
    url: `/api/resume-applications/${id}`,
    method: 'put',
    data
  })

export const deleteResumeApplication = (id: string) =>
  request({
    url: `/api/resume-applications/${id}`,
    method: 'delete'
  })
