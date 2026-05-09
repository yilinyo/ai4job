import request from '@/utils/request'
import {
  CreateQuestionOccurrenceRequest,
  QuestionOccurrenceQueryParams
} from '@/models'

export const getQuestionOccurrences = (params: QuestionOccurrenceQueryParams) =>
  request({
    url: '/api/question-occurrences',
    method: 'get',
    params
  })

export const createQuestionOccurrence = (data: CreateQuestionOccurrenceRequest) =>
  request({
    url: '/api/question-occurrences',
    method: 'post',
    data
  })

export const deleteQuestionOccurrence = (id: string) =>
  request({
    url: `/api/question-occurrences/${id}`,
    method: 'delete'
  })
