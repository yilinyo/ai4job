import request from '@/utils/request'
import {
  CreateInterviewProgressRequest,
  UpdateInterviewProgressRequest,
  InterviewProgressQueryParams
} from '@/models'

export const getInterviewProgressList = (params: InterviewProgressQueryParams) =>
  request({
    url: '/api/interview-progress',
    method: 'get',
    params
  })

export const createInterviewProgress = (data: CreateInterviewProgressRequest) =>
  request({
    url: '/api/interview-progress',
    method: 'post',
    data
  })

export const updateInterviewProgress = (id: string, data: UpdateInterviewProgressRequest) =>
  request({
    url: `/api/interview-progress/${id}`,
    method: 'put',
    data
  })

export const deleteInterviewProgress = (id: string) =>
  request({
    url: `/api/interview-progress/${id}`,
    method: 'delete'
  })
