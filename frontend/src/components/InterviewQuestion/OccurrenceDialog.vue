<template>
  <el-dialog
    :title="title"
    :visible.sync="dialogVisible"
    width="700px"
    @closed="resetForm"
  >
    <el-form
      ref="occurrenceForm"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item v-if="!lockQuestion" label="关联题目" prop="questionId">
        <el-select
          v-model="form.questionId"
          filterable
          placeholder="请选择要关联的题目"
          style="width: 100%"
        >
          <el-option
            v-for="item in questions"
            :key="item.id"
            :label="item.title"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="投递记录" prop="applicationId">
        <el-select
          v-model="form.applicationId"
          :disabled="lockApplication"
          filterable
          placeholder="请选择投递记录"
          style="width: 100%"
          @change="handleApplicationChange"
        >
          <el-option
            v-for="item in applications"
            :key="item.id"
            :label="`${item.companyName} / ${item.jobTitle}`"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="关联进展">
        <el-select
          v-model="form.interviewProgressId"
          clearable
          :disabled="lockProgress || !form.applicationId"
          placeholder="可选关联具体面试进展"
          style="width: 100%"
          @change="handleProgressChange"
        >
          <el-option
            v-for="item in progressOptions"
            :key="item.id"
            :label="buildProgressLabel(item)"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="面试阶段">
            <el-select
              v-model="form.interviewStageSnapshot"
              :disabled="Boolean(form.interviewProgressId)"
              clearable
              placeholder="未绑定进展时可手动选择"
              style="width: 100%"
            >
              <el-option
                v-for="item in interviewStageOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="提问时间" prop="occurredAt">
            <el-date-picker
              v-model="form.occurredAt"
              type="datetime"
              value-format="yyyy-MM-ddTHH:mm:ss"
              placeholder="请选择提问时间"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="当时问法">
        <el-input
          v-model="form.actualQuestion"
          type="textarea"
          :rows="3"
          placeholder="记录面试官真实提问方式或追问点"
        />
      </el-form-item>

      <el-form-item label="回答表现">
        <el-select
          v-model="form.answerPerformance"
          clearable
          placeholder="请选择回答表现"
          style="width: 100%"
        >
          <el-option
            v-for="item in answerPerformanceOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="复盘备注">
        <el-input
          v-model="form.note"
          type="textarea"
          :rows="3"
          placeholder="记录这道题在真实面试中的上下文"
        />
      </el-form-item>
    </el-form>

    <span slot="footer">
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        保存提问记录
      </el-button>
    </span>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import {
  InterviewProgress,
  QuestionAnswerPerformance,
  ResumeApplication
} from '@/models'
import {
  answerPerformanceOptions,
  interviewStageLabelMap,
  interviewStageOptions
} from '@/views/interview-question-bank/constants'

interface QuestionLite {
  id: string
  title: string
}

interface OccurrenceFormPayload {
  questionId: string
  applicationId: string
  interviewProgressId?: string
  interviewStageSnapshot?: string
  occurredAt: string
  actualQuestion?: string
  answerPerformance?: QuestionAnswerPerformance
  note?: string
}

@Component({
  name: 'OccurrenceDialog'
})
export default class extends Vue {
  @Prop({ default: false }) private visible!: boolean
  @Prop({ default: false }) private loading!: boolean
  @Prop({ default: '新增面试提问记录' }) private title!: string
  @Prop({ default: () => [] }) private questions!: QuestionLite[]
  @Prop({ default: () => [] }) private applications!: ResumeApplication[]
  @Prop({ default: () => [] }) private progressOptions!: InterviewProgress[]
  @Prop({ default: false }) private lockQuestion!: boolean
  @Prop({ default: false }) private lockApplication!: boolean
  @Prop({ default: false }) private lockProgress!: boolean
  @Prop({ default: '' }) private questionId!: string
  @Prop({ default: '' }) private applicationId!: string
  @Prop({ default: '' }) private progressId!: string
  @Prop({ default: '' }) private occurredAt!: string
  @Prop({ default: '' }) private stageSnapshot!: string

  private dialogVisible = false
  private answerPerformanceOptions = answerPerformanceOptions
  private interviewStageOptions = interviewStageOptions
  private interviewStageLabelMap = interviewStageLabelMap

  private form: OccurrenceFormPayload = this.getDefaultForm()

  private rules = {
    questionId: [{ required: true, message: '请选择题目', trigger: 'change' }],
    applicationId: [{ required: true, message: '请选择投递记录', trigger: 'change' }],
    occurredAt: [{ required: true, message: '请选择提问时间', trigger: 'change' }]
  }

  @Watch('visible', { immediate: true })
  private onVisibleChange(value: boolean) {
    this.dialogVisible = value
    if (value) {
      this.syncForm()
    }
  }

  @Watch('dialogVisible')
  private onDialogVisibleChange(value: boolean) {
    if (value !== this.visible) {
      this.$emit('update:visible', value)
    }
  }

  @Watch('questionId')
  @Watch('applicationId')
  @Watch('progressId')
  @Watch('occurredAt')
  @Watch('stageSnapshot')
  private onDefaultsChange() {
    if (this.visible) {
      this.syncForm()
    }
  }

  private getDefaultForm(): OccurrenceFormPayload {
    return {
      questionId: '',
      applicationId: '',
      interviewProgressId: '',
      interviewStageSnapshot: '',
      occurredAt: this.formatDateTimeInput(new Date().toISOString()),
      actualQuestion: '',
      answerPerformance: undefined,
      note: ''
    }
  }

  private syncForm() {
    this.form = {
      ...this.getDefaultForm(),
      questionId: this.questionId || '',
      applicationId: this.applicationId || '',
      interviewProgressId: this.progressId || '',
      interviewStageSnapshot: this.stageSnapshot || '',
      occurredAt: this.occurredAt || this.formatDateTimeInput(new Date().toISOString())
    }
  }

  private resetForm() {
    this.form = this.getDefaultForm()
    const formRef = this.$refs.occurrenceForm as ElForm | undefined
    if (formRef) {
      formRef.clearValidate()
    }
  }

  private formatDateTimeInput(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}T${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
  }

  private buildProgressLabel(item: InterviewProgress) {
    const stageLabel = this.interviewStageLabelMap[item.stage] || item.stage
    const date = new Date(item.occurredAt)
    const dateText = `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    return `${stageLabel} · ${dateText}`
  }

  private handleApplicationChange(value: string) {
    this.form.interviewProgressId = ''
    if (!this.lockProgress) {
      this.form.interviewStageSnapshot = ''
    }
    this.$emit('application-change', value)
  }

  private handleProgressChange(value: string) {
    const targetProgress = this.progressOptions.find(item => item.id === value)
    this.form.interviewStageSnapshot = targetProgress ? targetProgress.stage : this.form.interviewStageSnapshot
  }

  private handleSubmit() {
    (this.$refs.occurrenceForm as ElForm).validate((valid: boolean) => {
      if (!valid) return

      const payload: OccurrenceFormPayload = {
        questionId: this.lockQuestion ? this.questionId : this.form.questionId,
        applicationId: this.form.applicationId,
        interviewProgressId: this.form.interviewProgressId || undefined,
        interviewStageSnapshot: this.form.interviewStageSnapshot || undefined,
        occurredAt: this.form.occurredAt,
        actualQuestion: this.form.actualQuestion?.trim() || undefined,
        answerPerformance: this.form.answerPerformance,
        note: this.form.note?.trim() || undefined
      }

      if (!this.lockQuestion && !payload.questionId) {
        this.$message.warning('请选择题目')
        return
      }

      this.$emit('submit', payload)
    })
  }
}
</script>
