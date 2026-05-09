import { JobType } from '@/models'

export const jobTypeTextMap: Record<JobType, string> = {
  [JobType.CAMPUS]: '校招',
  [JobType.SOCIAL]: '社招',
  [JobType.INTERNSHIP]: '实习'
}

export const jobTypeTagMap: Record<JobType, string> = {
  [JobType.CAMPUS]: 'success',
  [JobType.SOCIAL]: 'primary',
  [JobType.INTERNSHIP]: 'warning'
}

export const getJobTypeText = (jobType?: JobType) => {
  return jobType ? jobTypeTextMap[jobType] : ''
}

export const getJobTypeTag = (jobType?: JobType) => {
  return jobType ? jobTypeTagMap[jobType] : ''
}
