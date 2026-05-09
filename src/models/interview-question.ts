/**
 * 面试题库与闪卡复习模型定义
 */

export enum InterviewQuestionDifficulty {
  EASY = 'easy',
  MEDIUM = 'medium',
  HARD = 'hard'
}

export enum InterviewQuestionType {
  KNOWLEDGE = 'knowledge',
  SCENARIO = 'scenario',
  CODING = 'coding',
  ALGORITHM = 'algorithm',
  PROJECT = 'project',
  BEHAVIOR = 'behavior',
  DISCUSSION = 'discussion'
}

export enum QuestionMasteryStatus {
  UNREVIEWED = 'unreviewed',
  WEAK = 'weak',
  NORMAL = 'normal',
  MASTERED = 'mastered'
}

export enum QuestionAnswerPerformance {
  FAILED = 'failed',
  VAGUE = 'vague',
  ANSWERED = 'answered',
  GOOD = 'good'
}

export enum FlashcardReviewResult {
  FAILED = 'failed',
  VAGUE = 'vague',
  MASTERED = 'mastered'
}

export enum FlashcardReviewMode {
  RANDOM = 'random',
  TAG = 'tag',
  DIFFICULTY = 'difficulty',
  FAVORITE_FIRST = 'favorite_first',
  WEAK_FIRST = 'weak_first',
  FREQUENT_FIRST = 'frequent_first'
}

export interface InterviewQuestion {
  id: string
  title: string
  content: string
  answerAnalysis: string
  tags: string[]
  difficulty: InterviewQuestionDifficulty
  questionType: InterviewQuestionType
  masteryStatus: QuestionMasteryStatus
  isFavorite: boolean
  source?: string
  occurrenceCount: number
  lastReviewedAt?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface FlashcardReviewSummary {
  totalReviews: number
  failedCount: number
  vagueCount: number
  masteredCount: number
  lastReviewedAt?: string
  lastResult?: FlashcardReviewResult
}

export interface QuestionOccurrence {
  id: string
  questionId: string
  applicationId: string
  interviewProgressId?: string
  companyNameSnapshot: string
  jobTitleSnapshot: string
  interviewStageSnapshot?: string
  occurredAt: string
  actualQuestion?: string
  answerPerformance?: QuestionAnswerPerformance
  note?: string
  createdAt: string
  updatedAt: string
}

export interface FlashcardReview {
  id: string
  questionId: string
  reviewedAt: string
  result: FlashcardReviewResult
  reviewMode: FlashcardReviewMode
  createdAt: string
}

export interface InterviewQuestionDetail extends InterviewQuestion {
  occurrences: QuestionOccurrence[]
  reviewSummary: FlashcardReviewSummary
  recentReviews: FlashcardReview[]
}

export interface CreateInterviewQuestionRequest {
  title: string
  content: string
  answerAnalysis: string
  tags?: string[]
  difficulty: InterviewQuestionDifficulty
  questionType: InterviewQuestionType
  masteryStatus: QuestionMasteryStatus
  isFavorite?: boolean
  source?: string
  remark?: string
}

export type UpdateInterviewQuestionRequest = Partial<CreateInterviewQuestionRequest>

export interface InterviewQuestionQueryParams {
  keyword?: string
  tags?: string
  difficulty?: InterviewQuestionDifficulty
  questionType?: InterviewQuestionType
  masteryStatus?: QuestionMasteryStatus
  isFavorite?: boolean
  page?: number
  pageSize?: number
}

export interface InterviewQuestionListResponse {
  list: InterviewQuestion[]
  total: number
  page: number
  pageSize: number
}

export interface UpdateQuestionFavoriteRequest {
  isFavorite: boolean
}

export interface QuestionOccurrenceQueryParams {
  questionId?: string
  applicationId?: string
}

export interface QuestionOccurrenceListResponse {
  list: QuestionOccurrence[]
}

export interface CreateQuestionOccurrenceRequest {
  questionId: string
  applicationId: string
  interviewProgressId?: string
  interviewStageSnapshot?: string
  occurredAt: string
  actualQuestion?: string
  answerPerformance?: QuestionAnswerPerformance
  note?: string
}

export interface FlashcardQueryParams {
  mode: FlashcardReviewMode
  tags?: string
  difficulty?: InterviewQuestionDifficulty
  limit?: number
}

export interface FlashcardQueueResponse {
  list: InterviewQuestion[]
  total: number
  mode: FlashcardReviewMode
}

export interface CreateFlashcardReviewRequest {
  questionId: string
  reviewedAt?: string
  result: FlashcardReviewResult
  reviewMode: FlashcardReviewMode
}

export interface CreateFlashcardReviewResponse {
  review: FlashcardReview
  updatedQuestion: InterviewQuestion
  reviewSummary: FlashcardReviewSummary
}
