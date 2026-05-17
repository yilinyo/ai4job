import request from '@/utils/request'
import {
  CreateInterviewQuestionRequest,
  InterviewQuestionQueryParams,
  UpdateInterviewQuestionRequest,
  UpdateQuestionFavoriteRequest
} from '@/models'

export const getInterviewQuestions = (params: InterviewQuestionQueryParams) =>
  request({
    url: '/api/interview-questions',
    method: 'get',
    params
  })

export const getInterviewQuestionById = (id: string) =>
  request({
    url: `/api/interview-questions/${id}`,
    method: 'get'
  })

export const createInterviewQuestion = (data: CreateInterviewQuestionRequest) =>
  request({
    url: '/api/interview-questions',
    method: 'post',
    data
  })

export const updateInterviewQuestion = (id: string, data: UpdateInterviewQuestionRequest) =>
  request({
    url: `/api/interview-questions/${id}`,
    method: 'put',
    data
  })

export const deleteInterviewQuestion = (id: string) =>
  request({
    url: `/api/interview-questions/${id}`,
    method: 'delete'
  })

export const updateInterviewQuestionFavorite = (id: string, data: UpdateQuestionFavoriteRequest) =>
  request({
    url: `/api/interview-questions/${id}/favorite`,
    method: 'put',
    data
  })
