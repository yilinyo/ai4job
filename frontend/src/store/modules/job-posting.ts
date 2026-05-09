import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import {
  getJobPostings,
  getJobPostingById,
  createJobPosting,
  updateJobPosting,
  deleteJobPosting,
  parseJobPosting
} from '@/api/job-postings'
import {
  JobPosting,
  JobPostingQueryParams,
  CreateJobPostingRequest,
  UpdateJobPostingRequest,
  ParseJobPostingResponse
} from '@/models'
import store from '@/store'

export interface IJobPostingState {
  jobPostings: JobPosting[]
  currentJobPosting: JobPosting | null
  parsedJobPosting: ParseJobPostingResponse | null
  total: number
  loading: boolean
}

@Module({ dynamic: true, store, name: 'jobPosting' })
class JobPostingStore extends VuexModule implements IJobPostingState {
  public jobPostings: JobPosting[] = []
  public currentJobPosting: JobPosting | null = null
  public parsedJobPosting: ParseJobPostingResponse | null = null
  public total = 0
  public loading = false

  @Mutation
  private SET_JOB_POSTINGS(jobPostings: JobPosting[]) {
    this.jobPostings = jobPostings
  }

  @Mutation
  private SET_CURRENT_JOB_POSTING(jobPosting: JobPosting | null) {
    this.currentJobPosting = jobPosting
  }

  @Mutation
  private SET_PARSED_JOB_POSTING(jobPosting: ParseJobPostingResponse | null) {
    this.parsedJobPosting = jobPosting
  }

  @Mutation
  private SET_TOTAL(total: number) {
    this.total = total
  }

  @Mutation
  private SET_LOADING(loading: boolean) {
    this.loading = loading
  }

  @Action
  public async GetJobPostings(params: JobPostingQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getJobPostings(params)
      if (data) {
        this.SET_JOB_POSTINGS(data.list)
        this.SET_TOTAL(data.total)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetJobPostingById(id: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await getJobPostingById(id)
      if (data) {
        this.SET_CURRENT_JOB_POSTING(data)
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async CreateJobPosting(payload: CreateJobPostingRequest) {
    const { data } = await createJobPosting(payload)
    if (data) {
      await this.GetJobPostings({ page: 1, pageSize: 9 })
      return data
    }
  }

  @Action
  public async UpdateJobPosting(payload: { id: string, data: UpdateJobPostingRequest }) {
    const { data } = await updateJobPosting(payload.id, payload.data)
    if (data) {
      if (this.currentJobPosting && this.currentJobPosting.id === payload.id) {
        this.SET_CURRENT_JOB_POSTING(data)
      }
      await this.GetJobPostings({ page: 1, pageSize: 9 })
      return data
    }
  }

  @Action
  public async DeleteJobPosting(id: string) {
    await deleteJobPosting(id)
    await this.GetJobPostings({ page: 1, pageSize: 9 })
  }

  @Action
  public async ParseJobPosting(url: string) {
    this.SET_LOADING(true)
    try {
      const { data } = await parseJobPosting({ url })
      if (data) {
        this.SET_PARSED_JOB_POSTING(data)
        return data
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public ClearCurrentJobPosting() {
    this.SET_CURRENT_JOB_POSTING(null)
  }

  @Action
  public ClearParsedJobPosting() {
    this.SET_PARSED_JOB_POSTING(null)
  }
}

export const JobPostingModule = getModule(JobPostingStore)
