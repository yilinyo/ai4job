import request from '@/utils/request'
import {
  CreateJobPostingRequest,
  UpdateJobPostingRequest,
  JobPostingQueryParams,
  ParseJobPostingRequest
} from '@/models'

export const getJobPostings = (params: JobPostingQueryParams) =>
  request({
    url: '/api/job-postings',
    method: 'get',
    params
  })

export const getJobPostingById = (id: string) =>
  request({
    url: `/api/job-postings/${id}`,
    method: 'get'
  })

export const createJobPosting = (data: CreateJobPostingRequest) =>
  request({
    url: '/api/job-postings',
    method: 'post',
    data
  })

export const updateJobPosting = (id: string, data: UpdateJobPostingRequest) =>
  request({
    url: `/api/job-postings/${id}`,
    method: 'put',
    data
  })

export const deleteJobPosting = (id: string) =>
  request({
    url: `/api/job-postings/${id}`,
    method: 'delete'
  })

export const parseJobPosting = (data: ParseJobPostingRequest) =>
  request({
    url: '/api/job-postings/parse',
    method: 'post',
    data
  })
