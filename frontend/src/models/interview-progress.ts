/**
 * 面试进展模型定义
 */

export enum InterviewStage {
  APPLIED = 'applied',
  VIEWED = 'viewed',
  WRITTEN_TEST = 'written_test',
  FIRST_INTERVIEW = 'first_interview',
  SECOND_INTERVIEW = 'second_interview',
  FINAL_INTERVIEW = 'final_interview',
  HR_INTERVIEW = 'hr_interview',
  OFFER = 'offer',
  REJECTED = 'rejected',
  CLOSED = 'closed'
}

export enum InterviewResult {
  PASSED = 'passed',
  PENDING = 'pending',
  FAILED = 'failed'
}

export interface InterviewProgress {
  id: string
  applicationId: string
  stage: InterviewStage
  occurredAt: string
  result: InterviewResult
  interviewerOrTeam?: string
  note?: string
  createdAt: string
  updatedAt: string
}

export interface CreateInterviewProgressRequest {
  applicationId: string
  stage: InterviewStage
  occurredAt: string
  result: InterviewResult
  interviewerOrTeam?: string
  note?: string
}

export type UpdateInterviewProgressRequest = Partial<CreateInterviewProgressRequest>

export interface InterviewProgressQueryParams {
  applicationId: string
}
