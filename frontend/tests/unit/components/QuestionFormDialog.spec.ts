import { createLocalVue, shallowMount } from '@vue/test-utils'
import ElementUI from 'element-ui'

jest.doMock('@/components/MarkdownEditor/index.vue', () => ({
  default: {
    name: 'MarkdownEditor',
    template: '<div class="markdown-editor-stub" />'
  }
}))

const QuestionFormDialog = require('@/components/InterviewQuestion/QuestionFormDialog.vue').default

const localVue = createLocalVue()
localVue.use(ElementUI)

const flushPromises = () => new Promise(resolve => setTimeout(resolve, 0))

describe('QuestionFormDialog', () => {
  it('syncs data, validates submission and resets cleanly on close', async() => {
    const wrapper = shallowMount(QuestionFormDialog, {
      localVue,
      propsData: {
        visible: true,
        loading: false,
        questionData: {
          title: 'Vue 响应式',
          content: '# 题目内容',
          answerAnalysis: '答案解析',
          tags: ['Vue'],
          difficulty: 'medium',
          questionType: 'knowledge',
          masteryStatus: 'unreviewed',
          isFavorite: false,
          source: '复盘',
          remark: '备注'
        }
      }
    })

    await flushPromises()

    expect((wrapper.vm as any).form.content).toBe('# 题目内容')
    expect((wrapper.vm as any).form.answerAnalysis).toBe('答案解析')
    expect(wrapper.html()).toContain('markdown-editor-stub')

    ;(wrapper.vm as any).handleContentChange('## 新内容')
    ;(wrapper.vm as any).handleAnswerAnalysisChange('新的答案')
    ;(wrapper.vm as any).$refs = {
      questionForm: {
        validate: jest.fn((cb) => {
          const valid = true
          cb(valid)
        }),
        clearValidate: jest.fn()
      }
    }
    ;(wrapper.vm as any).handleSubmit()

    expect(wrapper.emitted().submit?.[0]?.[0]).toMatchObject({
      title: 'Vue 响应式',
      content: '## 新内容',
      answerAnalysis: '新的答案',
      tags: ['Vue']
    })

    ;(wrapper.vm as any).$refs.questionForm.validate = jest.fn((cb) => {
      const valid = false
      cb(valid)
    })
    ;(wrapper.vm as any).handleSubmit()
    expect(wrapper.emitted().submit).toHaveLength(1)

    ;(wrapper.vm as any).$refs.questionForm.clearValidate = jest.fn()
    ;(wrapper.vm as any).resetForm()
    expect((wrapper.vm as any).form.content).toBe('')
    expect((wrapper.vm as any).form.answerAnalysis).toBe('')
    expect((wrapper.vm as any).$refs.questionForm.clearValidate).toHaveBeenCalled()
  })
})
