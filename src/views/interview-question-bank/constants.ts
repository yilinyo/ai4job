import {
  FlashcardReviewMode,
  FlashcardReviewResult,
  InterviewStage,
  InterviewQuestionDifficulty,
  InterviewQuestionType,
  QuestionAnswerPerformance,
  QuestionMasteryStatus
} from '@/models'

export const questionDifficultyOptions = [
  { label: '简单', value: InterviewQuestionDifficulty.EASY },
  { label: '中等', value: InterviewQuestionDifficulty.MEDIUM },
  { label: '困难', value: InterviewQuestionDifficulty.HARD }
]

export const questionTypeOptions = [
  { label: '八股题', value: InterviewQuestionType.KNOWLEDGE },
  { label: '场景题', value: InterviewQuestionType.SCENARIO },
  { label: '手写题', value: InterviewQuestionType.CODING },
  { label: '算法题', value: InterviewQuestionType.ALGORITHM },
  { label: '项目追问', value: InterviewQuestionType.PROJECT },
  { label: '行为面试', value: InterviewQuestionType.BEHAVIOR },
  { label: '开放讨论题', value: InterviewQuestionType.DISCUSSION }
]

export const masteryStatusOptions = [
  { label: '未复习', value: QuestionMasteryStatus.UNREVIEWED },
  { label: '生疏', value: QuestionMasteryStatus.WEAK },
  { label: '一般', value: QuestionMasteryStatus.NORMAL },
  { label: '熟练', value: QuestionMasteryStatus.MASTERED }
]

export const answerPerformanceOptions = [
  { label: '不会', value: QuestionAnswerPerformance.FAILED },
  { label: '模糊', value: QuestionAnswerPerformance.VAGUE },
  { label: '基本答出', value: QuestionAnswerPerformance.ANSWERED },
  { label: '回答较好', value: QuestionAnswerPerformance.GOOD }
]

export const flashcardModeOptions = [
  { label: '随机复习', value: FlashcardReviewMode.RANDOM },
  { label: '标签专项', value: FlashcardReviewMode.TAG },
  { label: '难度专项', value: FlashcardReviewMode.DIFFICULTY },
  { label: '收藏优先', value: FlashcardReviewMode.FAVORITE_FIRST },
  { label: '薄弱优先', value: FlashcardReviewMode.WEAK_FIRST },
  { label: '高频优先', value: FlashcardReviewMode.FREQUENT_FIRST }
]

export const flashcardResultActions = [
  { label: '不会', value: FlashcardReviewResult.FAILED, type: 'danger' },
  { label: '模糊', value: FlashcardReviewResult.VAGUE, type: 'warning' },
  { label: '掌握', value: FlashcardReviewResult.MASTERED, type: 'success' }
]

export const questionDifficultyLabelMap: Record<string, string> = {
  [InterviewQuestionDifficulty.EASY]: '简单',
  [InterviewQuestionDifficulty.MEDIUM]: '中等',
  [InterviewQuestionDifficulty.HARD]: '困难'
}

export const questionTypeLabelMap: Record<string, string> = {
  [InterviewQuestionType.KNOWLEDGE]: '八股题',
  [InterviewQuestionType.SCENARIO]: '场景题',
  [InterviewQuestionType.CODING]: '手写题',
  [InterviewQuestionType.ALGORITHM]: '算法题',
  [InterviewQuestionType.PROJECT]: '项目追问',
  [InterviewQuestionType.BEHAVIOR]: '行为面试',
  [InterviewQuestionType.DISCUSSION]: '开放讨论题'
}

export const masteryStatusLabelMap: Record<string, string> = {
  [QuestionMasteryStatus.UNREVIEWED]: '未复习',
  [QuestionMasteryStatus.WEAK]: '生疏',
  [QuestionMasteryStatus.NORMAL]: '一般',
  [QuestionMasteryStatus.MASTERED]: '熟练'
}

export const answerPerformanceLabelMap: Record<string, string> = {
  [QuestionAnswerPerformance.FAILED]: '不会',
  [QuestionAnswerPerformance.VAGUE]: '模糊',
  [QuestionAnswerPerformance.ANSWERED]: '基本答出',
  [QuestionAnswerPerformance.GOOD]: '回答较好'
}

export const flashcardModeLabelMap: Record<string, string> = {
  [FlashcardReviewMode.RANDOM]: '随机复习',
  [FlashcardReviewMode.TAG]: '标签专项',
  [FlashcardReviewMode.DIFFICULTY]: '难度专项',
  [FlashcardReviewMode.FAVORITE_FIRST]: '收藏优先',
  [FlashcardReviewMode.WEAK_FIRST]: '薄弱优先',
  [FlashcardReviewMode.FREQUENT_FIRST]: '高频优先'
}

export const reviewResultLabelMap: Record<string, string> = {
  [FlashcardReviewResult.FAILED]: '不会',
  [FlashcardReviewResult.VAGUE]: '模糊',
  [FlashcardReviewResult.MASTERED]: '掌握'
}

export const interviewStageOptions = [
  { label: '简历投递', value: InterviewStage.APPLIED },
  { label: '简历评估', value: InterviewStage.VIEWED },
  { label: '笔试', value: InterviewStage.WRITTEN_TEST },
  { label: '一面', value: InterviewStage.FIRST_INTERVIEW },
  { label: '二面', value: InterviewStage.SECOND_INTERVIEW },
  { label: '终面', value: InterviewStage.FINAL_INTERVIEW },
  { label: 'HR 面', value: InterviewStage.HR_INTERVIEW },
  { label: 'Offer', value: InterviewStage.OFFER },
  { label: '未通过', value: InterviewStage.REJECTED },
  { label: '已结束', value: InterviewStage.CLOSED }
]

export const interviewStageLabelMap: Record<string, string> = {
  [InterviewStage.APPLIED]: '简历投递',
  [InterviewStage.VIEWED]: '简历评估',
  [InterviewStage.WRITTEN_TEST]: '笔试',
  [InterviewStage.FIRST_INTERVIEW]: '一面',
  [InterviewStage.SECOND_INTERVIEW]: '二面',
  [InterviewStage.FINAL_INTERVIEW]: '终面',
  [InterviewStage.HR_INTERVIEW]: 'HR 面',
  [InterviewStage.OFFER]: 'Offer',
  [InterviewStage.REJECTED]: '未通过',
  [InterviewStage.CLOSED]: '已结束'
}

export const questionTagSuggestions = [
  'JavaScript',
  'TypeScript',
  'Vue',
  'React',
  'CSS',
  '浏览器',
  '网络',
  '工程化',
  '性能优化',
  '算法',
  '项目经历',
  '系统设计',
  'HR 面'
]
