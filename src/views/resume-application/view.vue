<template>
  <div class="resume-application-view-container">
    <el-page-header
      @back="goBack"
      content="投递记录详情"
    />

    <el-card v-loading="loading" class="application-card">
      <div v-if="currentApplication">
        <div class="application-header">
          <div>
            <h2>{{ currentApplication.companyName }} / {{ currentApplication.jobTitle }}</h2>
            <div class="header-meta">
              <el-tag :type="statusTagType">{{ statusLabel }}</el-tag>
              <span>{{ currentApplication.base || '未填写 base' }}</span>
              <span>投递于 {{ formatDate(currentApplication.appliedAt) }}</span>
            </div>
          </div>
          <div class="header-actions">
            <el-button type="primary" size="small" icon="el-icon-edit" @click="handleEdit">
              编辑投递
            </el-button>
            <el-button size="small" icon="el-icon-plus" @click="openProgressDialog()">
              新增进展
            </el-button>
          </div>
        </div>

        <div class="process-strip">
          <div
            v-for="node in processNodes"
            :key="node.stage"
            class="process-node"
            :class="`is-${node.state}`"
          >
            <div class="process-node__dot" />
            <div class="process-node__label">{{ node.label }}</div>
          </div>
        </div>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-card shadow="never" class="info-block">
              <div slot="header">基础信息</div>
              <el-descriptions :column="1" size="small" border>
                <el-descriptions-item label="公司介绍">
                  {{ currentApplication.companyIntro || '暂无' }}
                </el-descriptions-item>
                <el-descriptions-item label="岗位要求">
                  {{ currentApplication.jobRequirements || '暂无' }}
                </el-descriptions-item>
                <el-descriptions-item label="薪资范围">
                  {{ currentApplication.salaryRange || '暂无' }}
                </el-descriptions-item>
                <el-descriptions-item label="投递渠道">
                  {{ currentApplication.deliveryChannel || '暂无' }}
                </el-descriptions-item>
                <el-descriptions-item label="备注">
                  {{ currentApplication.remark || '暂无' }}
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>

          <el-col :span="12">
            <el-card shadow="never" class="info-block">
              <div slot="header">个人评分</div>
              <div class="score-list">
                <div class="score-item">
                  <span>简历匹配度</span>
                  <strong>{{ currentApplication.resumeMatchScore || '-' }}/5</strong>
                </div>
                <div class="score-item">
                  <span>面试表现</span>
                  <strong>{{ currentApplication.interviewPerformanceScore || '-' }}/5</strong>
                </div>
                <div class="score-item">
                  <span>岗位意向度</span>
                  <strong>{{ currentApplication.roleInterestScore || '-' }}/5</strong>
                </div>
                <div class="score-item">
                  <span>综合评分</span>
                  <strong>{{ currentApplication.overallScore || '-' }}/10</strong>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-card shadow="never" class="info-block">
              <div slot="header">面试总结</div>
              <div class="text-block">{{ currentApplication.interviewSummary || '暂无总结' }}</div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="never" class="info-block">
              <div slot="header">面试小记</div>
              <div class="text-block">{{ currentApplication.interviewNotes || '暂无面试小记' }}</div>
            </el-card>
          </el-col>
        </el-row>

        <el-card shadow="never" class="timeline-card">
          <div slot="header" class="timeline-header">
            <span>面试进展历史</span>
            <el-button size="small" type="primary" plain @click="openProgressDialog()">
              新增进展
            </el-button>
          </div>

          <el-empty v-if="!progressLoading && progressList.length === 0" description="暂无面试进展" />

          <el-timeline v-else>
            <el-timeline-item
              v-for="item in progressList"
              :key="item.id"
              :timestamp="formatDateTime(item.occurredAt)"
              :type="timelineType(item)"
            >
              <div class="timeline-item">
                <div class="timeline-item__header">
                  <strong>{{ stageLabelMap[item.stage] }}</strong>
                  <span>{{ resultLabelMap[item.result] }}</span>
                </div>
                <div class="timeline-item__meta">{{ item.interviewerOrTeam || '未填写面试官/团队' }}</div>
                <div class="timeline-item__note">{{ item.note || '暂无备注' }}</div>
                <div class="timeline-item__actions">
                  <el-button type="text" size="mini" @click="openProgressDialog(item)">
                    编辑
                  </el-button>
                  <el-button type="text" size="mini" class="danger" @click="handleDeleteProgress(item.id)">
                    删除
                  </el-button>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </div>
    </el-card>

    <el-dialog
      :title="progressDialogTitle"
      :visible.sync="progressDialogVisible"
      width="520px"
    >
      <el-form
        ref="progressForm"
        :model="progressForm"
        :rules="progressRules"
        label-width="110px"
      >
        <el-form-item label="流程阶段" prop="stage">
          <el-select v-model="progressForm.stage" placeholder="请选择流程阶段" style="width: 100%">
            <el-option
              v-for="item in stageOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发生时间" prop="occurredAt">
          <el-date-picker
            v-model="progressForm.occurredAt"
            type="datetime"
            value-format="yyyy-MM-ddTHH:mm:ss"
            placeholder="请选择时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结果" prop="result">
          <el-select v-model="progressForm.result" placeholder="请选择结果" style="width: 100%">
            <el-option
              v-for="item in resultOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试官/团队">
          <el-input v-model="progressForm.interviewerOrTeam" placeholder="可选填写" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="progressForm.note" type="textarea" :rows="3" placeholder="记录本轮进展情况" />
        </el-form-item>
      </el-form>

      <span slot="footer">
        <el-button @click="progressDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="progressSubmitting" @click="submitProgress">
          确定
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import { InterviewProgressModule } from '@/store/modules/interview-progress'
import {
  ApplicationStatus,
  InterviewProgress,
  InterviewStage,
  InterviewResult
} from '@/models'

interface ProcessNodeView {
  stage: InterviewStage
  label: string
  state: 'done' | 'current' | 'failed' | 'pending'
}

@Component({
  name: 'ResumeApplicationView'
})
export default class extends Vue {
  private applicationId = ''
  private progressDialogVisible = false
  private editingProgressId = ''
  private progressSubmitting = false

  private stageLabelMap: Record<string, string> = {
    [InterviewStage.APPLIED]: '简历投递',
    [InterviewStage.VIEWED]: '简历评估',
    [InterviewStage.WRITTEN_TEST]: '笔试',
    [InterviewStage.FIRST_INTERVIEW]: '一面',
    [InterviewStage.SECOND_INTERVIEW]: '二面',
    [InterviewStage.FINAL_INTERVIEW]: '三面/终面',
    [InterviewStage.HR_INTERVIEW]: 'HR 面',
    [InterviewStage.OFFER]: 'Offer',
    [InterviewStage.REJECTED]: '未通过',
    [InterviewStage.CLOSED]: '已结束'
  }

  private resultLabelMap: Record<string, string> = {
    [InterviewResult.PASSED]: '通过',
    [InterviewResult.PENDING]: '待定',
    [InterviewResult.FAILED]: '未通过'
  }

  private statusLabelMap: Record<string, string> = {
    [ApplicationStatus.APPLIED]: '已投递',
    [ApplicationStatus.VIEWED]: '已查看',
    [ApplicationStatus.WRITTEN_TEST]: '笔试',
    [ApplicationStatus.FIRST_INTERVIEW]: '一面',
    [ApplicationStatus.SECOND_INTERVIEW]: '二面',
    [ApplicationStatus.FINAL_INTERVIEW]: '终面',
    [ApplicationStatus.HR_INTERVIEW]: 'HR 面',
    [ApplicationStatus.OFFER]: 'Offer',
    [ApplicationStatus.REJECTED]: '未通过',
    [ApplicationStatus.CLOSED]: '已结束'
  }

  private stageOptions = [
    { label: '简历投递', value: InterviewStage.APPLIED },
    { label: '简历评估', value: InterviewStage.VIEWED },
    { label: '笔试', value: InterviewStage.WRITTEN_TEST },
    { label: '一面', value: InterviewStage.FIRST_INTERVIEW },
    { label: '二面', value: InterviewStage.SECOND_INTERVIEW },
    { label: '三面/终面', value: InterviewStage.FINAL_INTERVIEW },
    { label: 'HR 面', value: InterviewStage.HR_INTERVIEW },
    { label: 'Offer', value: InterviewStage.OFFER },
    { label: '未通过', value: InterviewStage.REJECTED },
    { label: '已结束', value: InterviewStage.CLOSED }
  ]

  private resultOptions = [
    { label: '通过', value: InterviewResult.PASSED },
    { label: '待定', value: InterviewResult.PENDING },
    { label: '未通过', value: InterviewResult.FAILED }
  ]

  private progressForm = {
    stage: InterviewStage.APPLIED,
    occurredAt: this.formatDateTimeInput(new Date().toISOString()),
    result: InterviewResult.PENDING,
    interviewerOrTeam: '',
    note: ''
  }

  private progressRules = {
    stage: [{ required: true, message: '请选择流程阶段', trigger: 'change' }],
    occurredAt: [{ required: true, message: '请选择时间', trigger: 'change' }],
    result: [{ required: true, message: '请选择结果', trigger: 'change' }]
  }

  get loading() {
    return ResumeApplicationModule.loading
  }

  get progressLoading() {
    return InterviewProgressModule.loading
  }

  get currentApplication() {
    return ResumeApplicationModule.currentApplication
  }

  get progressList() {
    return InterviewProgressModule.progressList
  }

  get progressDialogTitle() {
    return this.editingProgressId ? '编辑面试进展' : '新增面试进展'
  }

  get statusLabel() {
    if (!this.currentApplication) return ''
    return this.statusLabelMap[this.currentApplication.currentStatus] || this.currentApplication.currentStatus
  }

  get statusTagType() {
    if (!this.currentApplication) return 'info'
    if (this.currentApplication.currentStatus === ApplicationStatus.OFFER) return 'success'
    if (this.currentApplication.currentStatus === ApplicationStatus.REJECTED) return 'danger'
    return 'primary'
  }

  get processNodes(): ProcessNodeView[] {
    const orderedStages: InterviewStage[] = [
      InterviewStage.APPLIED,
      InterviewStage.VIEWED,
      InterviewStage.WRITTEN_TEST,
      InterviewStage.FIRST_INTERVIEW,
      InterviewStage.SECOND_INTERVIEW,
      InterviewStage.FINAL_INTERVIEW,
      InterviewStage.HR_INTERVIEW,
      InterviewStage.OFFER
    ]

    const latestProgress = this.progressList.length ? this.progressList[this.progressList.length - 1] : null
    const failedStage = latestProgress && latestProgress.result === InterviewResult.FAILED ? latestProgress.stage : ''
    const currentStage = this.currentApplication?.currentStatus === ApplicationStatus.REJECTED
      ? failedStage
      : (latestProgress ? latestProgress.stage : '')

    const reachedStages = new Set(this.progressList.map(item => item.stage))

    return orderedStages.map((stage, index) => {
      const currentIndex = orderedStages.indexOf(currentStage as InterviewStage)

      if (failedStage === stage) {
        return { stage, label: this.stageLabelMap[stage], state: 'failed' }
      }

      if (reachedStages.has(stage)) {
        if (stage === currentStage && this.currentApplication?.currentStatus !== ApplicationStatus.OFFER) {
          return { stage, label: this.stageLabelMap[stage], state: 'current' }
        }

        if (this.currentApplication?.currentStatus === ApplicationStatus.OFFER && stage === InterviewStage.OFFER) {
          return { stage, label: this.stageLabelMap[stage], state: 'current' }
        }

        return { stage, label: this.stageLabelMap[stage], state: 'done' }
      }

      if (currentIndex !== -1 && index < currentIndex) {
        return { stage, label: this.stageLabelMap[stage], state: 'done' }
      }

      return { stage, label: this.stageLabelMap[stage], state: 'pending' }
    })
  }

  created() {
    this.applicationId = this.$route.params.id
    this.loadData()
  }

  private async loadData() {
    await ResumeApplicationModule.GetResumeApplicationById(this.applicationId)
    await InterviewProgressModule.GetInterviewProgressList({ applicationId: this.applicationId })
  }

  private formatDate(dateString: string) {
    return this.formatDateTime(dateString, false)
  }

  private formatDateTimeInput(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}T${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
  }

  private formatDateTime(dateString: string, includeTime = true) {
    const date = new Date(dateString)
    const datePart = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    if (!includeTime) return datePart
    return `${datePart} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }

  private timelineType(item: InterviewProgress) {
    return item.result === InterviewResult.FAILED ? 'danger' : item.result === InterviewResult.PENDING ? 'primary' : 'success'
  }

  private handleEdit() {
    if (!this.currentApplication) return
    this.$router.push(`/resume/application/${this.currentApplication.id}/edit?versionId=${this.currentApplication.resumeVersionId}&repoId=${this.currentApplication.repoId}`)
  }

  private openProgressDialog(progress?: InterviewProgress) {
    if (progress) {
      this.editingProgressId = progress.id
      this.progressForm = {
        stage: progress.stage,
        occurredAt: this.formatDateTimeInput(progress.occurredAt),
        result: progress.result,
        interviewerOrTeam: progress.interviewerOrTeam || '',
        note: progress.note || ''
      }
    } else {
      this.editingProgressId = ''
      this.progressForm = {
        stage: InterviewStage.APPLIED,
        occurredAt: this.formatDateTimeInput(new Date().toISOString()),
        result: InterviewResult.PENDING,
        interviewerOrTeam: '',
        note: ''
      }
    }
    this.progressDialogVisible = true
  }

  private submitProgress() {
    (this.$refs.progressForm as ElForm).validate(async(valid: boolean) => {
      if (!valid) return

      this.progressSubmitting = true
      try {
        if (this.editingProgressId) {
          await InterviewProgressModule.UpdateInterviewProgress({
            id: this.editingProgressId,
            applicationId: this.applicationId,
            data: {
              applicationId: this.applicationId,
              ...this.progressForm
            }
          })
          this.$message.success('更新成功！')
        } else {
          await InterviewProgressModule.CreateInterviewProgress({
            applicationId: this.applicationId,
            ...this.progressForm
          })
          this.$message.success('创建成功！')
        }
        await ResumeApplicationModule.GetResumeApplicationById(this.applicationId)
        this.progressDialogVisible = false
      } catch (error) {
        console.error('保存面试进展失败:', error)
        this.$message.error('保存失败，请稍后重试')
      } finally {
        this.progressSubmitting = false
      }
    })
  }

  private async handleDeleteProgress(progressId: string) {
    try {
      await this.$confirm('确定要删除这条面试进展吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await InterviewProgressModule.DeleteInterviewProgress({
        id: progressId,
        applicationId: this.applicationId
      })
      await ResumeApplicationModule.GetResumeApplicationById(this.applicationId)
      this.$message.success('删除成功！')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
  }

  private goBack() {
    if (this.currentApplication) {
      this.$router.push(`/resume/version/${this.currentApplication.resumeVersionId}`)
    } else {
      this.$router.back()
    }
  }
}
</script>

<style lang="scss" scoped>
.resume-application-view-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .application-card {
    .application-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 20px;
      margin-bottom: 24px;

      h2 {
        margin: 0 0 8px;
      }

      .header-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
        color: #606266;
      }

      .header-actions {
        display: flex;
        gap: 10px;
      }
    }

    .process-strip {
      display: grid;
      grid-template-columns: repeat(8, minmax(0, 1fr));
      gap: 12px;
      margin-bottom: 24px;
    }

    .process-node {
      padding: 12px 8px;
      border-radius: 12px;
      background: #f5f7fa;
      border: 1px solid #ebeef5;
      text-align: center;
      transition: all 0.2s ease;

      &__dot {
        width: 12px;
        height: 12px;
        margin: 0 auto 8px;
        border-radius: 50%;
        background: #c0c4cc;
      }

      &__label {
        font-size: 12px;
        color: #606266;
        line-height: 1.4;
      }

      &.is-done,
      &.is-current {
        background: #f0f9eb;
        border-color: #b7eb8f;

        .process-node__dot {
          background: #67c23a;
        }

        .process-node__label {
          color: #389e0d;
        }
      }

      &.is-current {
        box-shadow: 0 0 0 2px rgba(103, 194, 58, 0.18);
      }

      &.is-failed {
        background: #fef0f0;
        border-color: #f56c6c;

        .process-node__dot {
          background: #f56c6c;
        }

        .process-node__label {
          color: #c45656;
        }
      }
    }

    .info-block {
      margin-bottom: 20px;
    }

    .score-list {
      display: grid;
      gap: 12px;
    }

    .score-item {
      display: flex;
      justify-content: space-between;
      padding: 10px 12px;
      border-radius: 8px;
      background: #f5f7fa;
    }

    .text-block {
      white-space: pre-wrap;
      line-height: 1.8;
      color: #606266;
      min-height: 90px;
    }

    .timeline-card {
      margin-top: 4px;

      .timeline-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .timeline-item {
        &__header {
          display: flex;
          justify-content: space-between;
          margin-bottom: 6px;
        }

        &__meta,
        &__note {
          color: #606266;
          margin-bottom: 4px;
          white-space: pre-wrap;
        }

        &__actions {
          .danger {
            color: #f56c6c;
          }
        }
      }
    }
  }
}
</style>
