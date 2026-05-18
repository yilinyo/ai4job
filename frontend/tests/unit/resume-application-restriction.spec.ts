import { createLocalVue, shallowMount } from '@vue/test-utils'
import ElementUI from 'element-ui'

jest.mock('html2pdf.js', () => jest.fn())

const mockResumeVersionModule = {
  currentVersion: {
    id: 'v1',
    repoId: 'r1',
    title: 'v1',
    version: 'v1.0.0',
    content: '# Resume',
    remark: '',
    isActive: false,
    createdAt: '2026-05-18T00:00:00.000Z',
    updatedAt: '2026-05-18T00:00:00.000Z'
  },
  loading: false,
  GetResumeVersionById: jest.fn().mockResolvedValue(undefined),
  GetResumeVersions: jest.fn().mockResolvedValue(undefined)
}

const historyApplication = {
  id: 'app1',
  companyName: '腾讯科技',
  jobTitle: '前端开发工程师',
  base: '深圳',
  appliedAt: '2026-04-10T00:00:00.000Z',
  currentStatus: 'applied'
}

const mockResumeApplicationModule = {
  applications: [historyApplication],
  total: 1,
  loading: false,
  currentApplication: null,
  GetResumeApplications: jest.fn().mockResolvedValue(undefined),
  GetResumeApplicationById: jest.fn().mockResolvedValue(undefined),
  CreateResumeApplication: jest.fn().mockResolvedValue(undefined),
  UpdateResumeApplication: jest.fn().mockResolvedValue(undefined),
  DeleteResumeApplication: jest.fn().mockResolvedValue(undefined)
}

const mockJobPostingModule = {
  jobPostings: [],
  loading: false,
  currentJobPosting: null,
  GetJobPostings: jest.fn().mockResolvedValue(undefined),
  GetJobPostingById: jest.fn().mockResolvedValue(undefined)
}

jest.doMock('@/store/modules/resume-version', () => ({
  ResumeVersionModule: mockResumeVersionModule
}))

jest.doMock('@/store/modules/resume-application', () => ({
  ResumeApplicationModule: mockResumeApplicationModule
}))

jest.doMock('@/store/modules/job-posting', () => ({
  JobPostingModule: mockJobPostingModule
}))

const ResumeVersionView = require('@/views/resume-version/view.vue').default
const ResumeApplicationEditor = require('@/views/resume-application/editor.vue').default

const localVue = createLocalVue()
localVue.use(ElementUI)

const flushPromises = () => new Promise(resolve => setTimeout(resolve, 0))

describe('Resume application restriction flow', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockResumeApplicationModule.applications = [historyApplication]
    mockResumeApplicationModule.total = 1
    mockResumeVersionModule.currentVersion = {
      id: 'v1',
      repoId: 'r1',
      title: 'v1',
      version: 'v1.0.0',
      content: '# Resume',
      remark: '',
      isActive: false,
      createdAt: '2026-05-18T00:00:00.000Z',
      updatedAt: '2026-05-18T00:00:00.000Z'
    }
  })

  it('hides the create-application entry on inactive resume version detail pages', async() => {
    const wrapper = shallowMount(ResumeVersionView, {
      localVue,
      mocks: {
        $route: { params: { id: 'v1' } },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          warning: jest.fn(),
          error: jest.fn()
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).not.toContain('新增投递记录')
    expect(wrapper.text()).toContain('当前版本共关联 1 条投递记录')
    expect(mockResumeVersionModule.GetResumeVersionById).toHaveBeenCalledWith('v1')
    expect(mockResumeApplicationModule.GetResumeApplications).toHaveBeenCalledWith({
      resumeVersionId: 'v1',
      page: 1,
      pageSize: 100
    })
  })

  it('shows the create-application entry on active resume version detail pages', async() => {
    mockResumeVersionModule.currentVersion.isActive = true

    const wrapper = shallowMount(ResumeVersionView, {
      localVue,
      mocks: {
        $route: { params: { id: 'v1' } },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          warning: jest.fn(),
          error: jest.fn()
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('新增投递记录')
  })

  it('blocks application creation when the editor is opened without an active version context', async() => {
    const wrapper = shallowMount(ResumeApplicationEditor, {
      localVue,
      mocks: {
        $route: {
          params: {},
          query: {
            jobPostingId: 'job-1'
          },
          path: '/resume/application/create'
        },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          warning: jest.fn(),
          error: jest.fn()
        },
        $confirm: jest.fn()
      }
    })

    await flushPromises()

    expect((wrapper.vm as any).creationBlockedMessage).toBe('请先从激活的简历版本详情页进入，再创建投递记录。')
    expect(mockJobPostingModule.GetJobPostings).toHaveBeenCalled()
    expect(mockResumeVersionModule.GetResumeVersionById).not.toHaveBeenCalled()

    ;(wrapper.vm as any).$refs = {
      applicationForm: {
        validate: jest.fn()
      }
    }
    ;(wrapper.vm as any).handleSubmit()

    expect((wrapper.vm as any).$message.warning).toHaveBeenCalledWith('请先从激活的简历版本详情页进入，再创建投递记录。')
    expect(mockResumeApplicationModule.CreateResumeApplication).not.toHaveBeenCalled()
  })

  it('allows application creation when the editor is opened with an active version context', async() => {
    mockResumeVersionModule.currentVersion.isActive = true

    const wrapper = shallowMount(ResumeApplicationEditor, {
      localVue,
      mocks: {
        $route: {
          params: {},
          query: {
            versionId: 'v1',
            repoId: 'r1'
          },
          path: '/resume/application/create'
        },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          warning: jest.fn(),
          error: jest.fn()
        },
        $confirm: jest.fn()
      }
    })

    await flushPromises()

    expect((wrapper.vm as any).creationBlockedMessage).toBe('')
    ;(wrapper.vm as any).$refs = {
      applicationForm: {
        validate: jest.fn(cb => cb(true))
      }
    }
    ;(wrapper.vm as any).handleSubmit()

    expect(mockResumeApplicationModule.CreateResumeApplication).toHaveBeenCalled()
  })
})
