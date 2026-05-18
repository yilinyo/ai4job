import { createLocalVue, shallowMount } from '@vue/test-utils'
import ElementUI from 'element-ui'

const mockResumeVersionModule = {
  currentVersion: {
    id: 'version-1',
    repoId: 'repo-1',
    title: '基础版',
    content: '# 原始内容',
    remark: 'remark',
    isActive: true,
    createdAt: '2026-05-18T00:00:00.000Z',
    updatedAt: '2026-05-18T00:00:00.000Z'
  },
  loading: false,
  GetResumeVersionById: jest.fn().mockResolvedValue(undefined),
  UpdateResumeVersion: jest.fn().mockResolvedValue(undefined),
  CreateResumeVersion: jest.fn().mockResolvedValue(undefined)
}

const mockUserModule = {
  userProfile: {
    name: '测试用户'
  },
  GetUserInfo: jest.fn().mockResolvedValue(undefined)
}

jest.doMock('@/store/modules/resume-version', () => ({
  ResumeVersionModule: mockResumeVersionModule
}))

jest.doMock('@/store/modules/user', () => ({
  UserModule: mockUserModule
}))

jest.doMock('@/components/MarkdownEditor/index.vue', () => ({
  default: {
    name: 'MarkdownEditor',
    template: '<div class="markdown-editor-stub" />'
  }
}))

const ResumeVersionEditor = require('@/views/resume-version/editor.vue').default

const localVue = createLocalVue()
localVue.use(ElementUI)

const flushPromises = () => new Promise(resolve => setTimeout(resolve, 0))

describe('Resume version editor', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockResumeVersionModule.currentVersion = {
      id: 'version-1',
      repoId: 'repo-1',
      title: '基础版',
      content: '# 原始内容',
      remark: 'remark',
      isActive: true,
      createdAt: '2026-05-18T00:00:00.000Z',
      updatedAt: '2026-05-18T00:00:00.000Z'
    }
  })

  it('renders the shared markdown editor and submits the latest markdown content', async() => {
    const wrapper = shallowMount(ResumeVersionEditor, {
      localVue,
      mocks: {
        $route: {
          params: { id: 'version-1' },
          query: { repoId: 'repo-1' },
          path: '/resume/version/1/edit'
        },
        $router: { push: jest.fn(), back: jest.fn() },
        $message: {
          success: jest.fn(),
          error: jest.fn()
        }
      }
    })

    await flushPromises()

    expect(mockResumeVersionModule.GetResumeVersionById).toHaveBeenCalledWith('version-1')
    expect(wrapper.html()).toContain('markdown-editor-stub')
    expect((wrapper.vm as any).versionForm.content).toBe('# 原始内容')

    ;(wrapper.vm as any).handleContentChange('## 更新后的简历')
    ;(wrapper.vm as any).$refs = {
      versionForm: {
        validate: jest.fn((cb) => {
          const valid = true
          cb(valid)
        })
      }
    }
    ;(wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect(mockResumeVersionModule.UpdateResumeVersion).toHaveBeenCalledWith({
      id: 'version-1',
      data: {
        title: '基础版',
        content: '## 更新后的简历',
        remark: 'remark'
      }
    })
  })
})
