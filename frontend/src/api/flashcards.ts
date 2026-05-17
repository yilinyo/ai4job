import request from '@/utils/request'
import {
  CreateFlashcardReviewRequest,
  FlashcardQueryParams
} from '@/models'

export const getFlashcards = (params: FlashcardQueryParams) =>
  request({
    url: '/api/flashcards',
    method: 'get',
    params
  })

export const createFlashcardReview = (data: CreateFlashcardReviewRequest) =>
  request({
    url: '/api/flashcard-reviews',
    method: 'post',
    data
  })
