<template>
  <el-dialog
    :title="title"
    :visible.sync="dialogVisible"
    width="960px"
    top="4vh"
    custom-class="question-form-dialog"
    @closed="resetForm"
  >
    <el-form
      ref="questionForm"
      :model="form"
      :rules="rules"
      label-width="110px"
    >
      <el-form-item label="题目标题" prop="title">
        <el-input v-model="form.title" placeholder="请输入题目标题" maxlength="80" show-word-limit />
      </el-form-item>
      <el-form-item label="题目内容" prop="content">
        <markdown-editor
          class="question-markdown-editor"
          :initial-value="form.content"
          :height="editorHeight"
          @change="handleContentChange"
        />
      </el-form-item>
      <el-form-item label="答案解析" prop="answerAnalysis">
        <markdown-editor
          class="question-markdown-editor"
          :initial-value="form.answerAnalysis"
          :height="editorHeight"
          @change="handleAnswerAnalysisChange"
        />
      </el-form-item>
      <el-form-item label="标签">
        <el-select
          v-model="form.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="可输入多个标签"
          style="width: 100%"
        >
          <el-option
            v-for="tag in questionTagSuggestions"
            :key="tag"
            :label="tag"
            :value="tag"
          />
        </el-select>
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="难度" prop="difficulty">
            <el-select v-model="form.difficulty" placeholder="请选择难度" style="width: 100%">
              <el-option
                v-for="item in questionDifficultyOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="题型" prop="questionType">
            <el-select v-model="form.questionType" placeholder="请选择题型" style="width: 100%">
              <el-option
                v-for="item in questionTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="掌握状态" prop="masteryStatus">
            <el-select v-model="form.masteryStatus" placeholder="请选择掌握状态" style="width: 100%">
              <el-option
                v-for="item in masteryStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="来源">
            <el-input v-model="form.source" placeholder="如：面试复盘 / 牛客 / 自己整理" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="重点收藏">
            <el-switch
              v-model="form.isFavorite"
              active-text="收藏"
              inactive-text="普通"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="备注">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="3"
          placeholder="补充说明这道题的背景、坑点或后续行动"
        />
      </el-form-item>
    </el-form>

    <span slot="footer">
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        保存题目
      </el-button>
    </span>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import MarkdownEditor from '@/components/MarkdownEditor/index.vue'
import {
  CreateInterviewQuestionRequest,
  InterviewQuestionDifficulty,
  InterviewQuestionType,
  QuestionMasteryStatus
} from '@/models'
import {
  masteryStatusOptions,
  questionDifficultyOptions,
  questionTagSuggestions,
  questionTypeOptions
} from '@/views/interview-question-bank/constants'

type QuestionFormValue = Partial<CreateInterviewQuestionRequest> & {
  isFavorite?: boolean
}

@Component({
  name: 'QuestionFormDialog',
  components: {
    MarkdownEditor
  }
})
export default class extends Vue {
  @Prop({ default: false }) private visible!: boolean
  @Prop({ default: false }) private loading!: boolean
  @Prop({ default: '新增题目' }) private title!: string
  @Prop({ default: () => null }) private questionData!: QuestionFormValue | null

  private dialogVisible = false
  private questionDifficultyOptions = questionDifficultyOptions
  private questionTypeOptions = questionTypeOptions
  private masteryStatusOptions = masteryStatusOptions
  private questionTagSuggestions = questionTagSuggestions
  private editorHeight = '240px'

  private form: Required<CreateInterviewQuestionRequest> & { isFavorite: boolean } = this.getDefaultForm()

  private rules = {
    title: [{ required: true, message: '请输入题目标题', trigger: 'blur' }],
    content: [{ required: true, message: '请输入题目内容', trigger: 'change' }],
    answerAnalysis: [{ required: true, message: '请输入答案解析', trigger: 'change' }],
    difficulty: [{ required: true, message: '请选择难度', trigger: 'change' }],
    questionType: [{ required: true, message: '请选择题型', trigger: 'change' }],
    masteryStatus: [{ required: true, message: '请选择掌握状态', trigger: 'change' }]
  }

  @Watch('visible', { immediate: true })
  private onVisibleChange(value: boolean) {
    this.dialogVisible = value
    if (value) {
      this.syncForm()
    }
  }

  @Watch('questionData', { deep: true })
  private onQuestionDataChange() {
    if (this.visible) {
      this.syncForm()
    }
  }

  @Watch('dialogVisible')
  private onDialogVisibleChange(value: boolean) {
    if (value !== this.visible) {
      this.$emit('update:visible', value)
    }
  }

  private getDefaultForm() {
    return {
      title: '',
      content: '',
      answerAnalysis: '',
      tags: [],
      difficulty: InterviewQuestionDifficulty.MEDIUM,
      questionType: InterviewQuestionType.KNOWLEDGE,
      masteryStatus: QuestionMasteryStatus.UNREVIEWED,
      isFavorite: false,
      source: '',
      remark: ''
    }
  }

  private syncForm() {
    const base = this.getDefaultForm()
    const nextData = this.questionData || {}
    this.form = {
      ...base,
      ...nextData,
      tags: Array.isArray(nextData.tags) ? [...nextData.tags] : []
    }
  }

  private resetForm() {
    this.form = this.getDefaultForm()
    const formRef = this.$refs.questionForm as ElForm | undefined
    if (formRef) {
      formRef.clearValidate()
    }
  }

  private handleSubmit() {
    (this.$refs.questionForm as ElForm).validate((valid: boolean) => {
      if (!valid) return

      const payload: CreateInterviewQuestionRequest = {
        title: this.form.title.trim(),
        content: this.form.content.trim(),
        answerAnalysis: this.form.answerAnalysis.trim(),
        tags: this.form.tags,
        difficulty: this.form.difficulty,
        questionType: this.form.questionType,
        masteryStatus: this.form.masteryStatus,
        isFavorite: this.form.isFavorite,
        source: this.form.source?.trim() || '',
        remark: this.form.remark?.trim() || ''
      }

      this.$emit('submit', payload)
    })
  }

  private handleContentChange(value: string) {
    this.form.content = value
  }

  private handleAnswerAnalysisChange(value: string) {
    this.form.answerAnalysis = value
  }
}
</script>

<style lang="scss" scoped>
::v-deep .question-form-dialog {
  width: min(960px, calc(100vw - 32px)) !important;
}

::v-deep .question-form-dialog .el-dialog__body {
  max-height: calc(100vh - 180px);
  overflow-y: auto;
}

.question-markdown-editor {
  display: block;
  width: 100%;
}

::v-deep .question-markdown-editor .toastui-editor-defaultUI {
  min-height: 240px;
}
</style>
