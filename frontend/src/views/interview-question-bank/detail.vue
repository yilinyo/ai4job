<template>
  <!-- eslint-disable vue/no-v-html -->
  <div class="question-detail-container">
    <el-page-header @back="goBack" content="题目详情" />

    <el-card v-loading="loading" class="detail-card">
      <template v-if="question">
        <div class="detail-header">
          <div>
            <h2>{{ question.title }}</h2>
            <div class="meta-line">
              <el-tag :type="difficultyTagType(question.difficulty)">{{ questionDifficultyLabelMap[question.difficulty] }}</el-tag>
              <el-tag type="info">{{ questionTypeLabelMap[question.questionType] }}</el-tag>
              <el-tag :type="masteryTagType(question.masteryStatus)">{{ masteryStatusLabelMap[question.masteryStatus] }}</el-tag>
              <el-tag v-if="question.isFavorite" type="warning">
                <i class="el-icon-star-on" /> 已收藏
              </el-tag>
              <span>提问 {{ question.occurrenceCount }} 次</span>
              <span>最近复习：{{ question.lastReviewedAt ? formatDateTime(question.lastReviewedAt) : '未复习' }}</span>
              <span>创建时间：{{ formatDateTime(question.createdAt) }}</span>
              <span>修改时间：{{ formatDateTime(question.updatedAt) }}</span>
            </div>
          </div>
          <div class="header-actions">
            <el-button icon="el-icon-plus" @click="openOccurrenceDialog">新增提问记录</el-button>
            <el-button type="primary" icon="el-icon-edit" @click="openEditDialog">编辑题目</el-button>
            <el-button
              :type="question.isFavorite ? 'warning' : 'default'"
              :icon="question.isFavorite ? 'el-icon-star-on' : 'el-icon-star-off'"
              @click="toggleFavorite"
            >
              {{ question.isFavorite ? '取消收藏' : '收藏' }}
            </el-button>
          </div>
        </div>

        <div class="tag-list">
          <el-tag v-for="tag in question.tags" :key="tag" size="small" effect="plain">{{ tag }}</el-tag>
        </div>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-card shadow="never" class="content-card">
              <div slot="header" class="card-header">
                <span>题目内容</span>
                <span class="edit-hint">双击编辑</span>
              </div>
              <resume-markdown-preview
                class="question-markdown-preview"
                :html="renderMarkdown(question.content)"
                empty-text="暂无题目内容"
                @dblclick.native="openEditDialog"
              />
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="never" class="content-card">
              <div slot="header" class="card-header">
                <span>答案解析</span>
                <span class="edit-hint">双击编辑</span>
              </div>
              <resume-markdown-preview
                class="question-markdown-preview"
                :html="renderMarkdown(question.answerAnalysis)"
                empty-text="暂无答案解析"
                @dblclick.native="openEditDialog"
              />
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="16">
            <el-card shadow="never" class="content-card">
              <div slot="header" class="card-header">
                <span>真实面试提问记录</span>
                <el-button type="text" icon="el-icon-plus" @click="openOccurrenceDialog">新增</el-button>
              </div>
              <el-empty v-if="question.occurrences.length === 0" description="暂无真实面试提问记录" />
              <el-table v-else :data="question.occurrences" style="width: 100%">
                <el-table-column label="公司 / 岗位" min-width="180">
                  <template slot-scope="{row}">
                    <div>{{ row.companyNameSnapshot }}</div>
                    <div class="sub-text">{{ row.jobTitleSnapshot }}</div>
                  </template>
                </el-table-column>
                <el-table-column label="阶段" width="100">
                  <template slot-scope="{row}">
                    {{ row.interviewStageSnapshot ? interviewStageLabelMap[row.interviewStageSnapshot] || row.interviewStageSnapshot : '未填写' }}
                  </template>
                </el-table-column>
                <el-table-column label="提问时间" width="150">
                  <template slot-scope="{row}">{{ formatDateTime(row.occurredAt) }}</template>
                </el-table-column>
                <el-table-column label="当时问法" min-width="220">
                  <template slot-scope="{row}">{{ row.actualQuestion || '未记录' }}</template>
                </el-table-column>
                <el-table-column label="表现" width="100">
                  <template slot-scope="{row}">
                    {{ row.answerPerformance ? answerPerformanceLabelMap[row.answerPerformance] : '未记录' }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="130">
                  <template slot-scope="{row}">
                    <el-button type="text" size="small" @click="goApplication(row.applicationId)">投递详情</el-button>
                    <el-button type="text" size="small" class="danger" @click="deleteOccurrence(row.id)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="never" class="content-card">
              <div slot="header">复习摘要</div>
              <div class="summary-grid">
                <div class="summary-item">
                  <span>总复习</span>
                  <strong>{{ question.reviewSummary.totalReviews }}</strong>
                </div>
                <div class="summary-item danger">
                  <span>不会</span>
                  <strong>{{ question.reviewSummary.failedCount }}</strong>
                </div>
                <div class="summary-item warning">
                  <span>模糊</span>
                  <strong>{{ question.reviewSummary.vagueCount }}</strong>
                </div>
                <div class="summary-item success">
                  <span>掌握</span>
                  <strong>{{ question.reviewSummary.masteredCount }}</strong>
                </div>
              </div>
              <el-divider />
              <div v-if="question.recentReviews.length">
                <div
                  v-for="review in question.recentReviews"
                  :key="review.id"
                  class="review-row"
                >
                  <span>{{ formatDateTime(review.reviewedAt) }}</span>
                  <el-tag size="mini">{{ reviewResultLabelMap[review.result] }}</el-tag>
                </div>
              </div>
              <el-empty v-else description="暂无复习记录" />
            </el-card>
          </el-col>
        </el-row>
      </template>
    </el-card>

    <question-form-dialog
      :visible.sync="questionDialogVisible"
      title="编辑题目"
      :loading="submitting"
      :question-data="question"
      @submit="submitQuestion"
    />

    <occurrence-dialog
      :visible.sync="occurrenceDialogVisible"
      title="新增面试提问记录"
      :loading="submitting"
      :questions="question ? [question] : []"
      :applications="applications"
      :progress-options="progressList"
      :question-id="questionId"
      :lock-question="true"
      @application-change="loadProgressForApplication"
      @submit="submitOccurrence"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import QuestionFormDialog from '@/components/InterviewQuestion/QuestionFormDialog.vue'
import OccurrenceDialog from '@/components/InterviewQuestion/OccurrenceDialog.vue'
import ResumeMarkdownPreview from '@/components/ResumeMarkdownPreview/index.vue'
import { InterviewQuestionBankModule } from '@/store/modules/interview-question-bank'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import { InterviewProgressModule } from '@/store/modules/interview-progress'
import {
  CreateInterviewQuestionRequest,
  CreateQuestionOccurrenceRequest,
  InterviewQuestionDifficulty,
  QuestionMasteryStatus
} from '@/models'
import {
  answerPerformanceLabelMap,
  interviewStageLabelMap,
  masteryStatusLabelMap,
  questionDifficultyLabelMap,
  questionTypeLabelMap,
  reviewResultLabelMap
} from './constants'
import { renderResumeMarkdown } from '@/utils/resume-markdown'

@Component({
  name: 'InterviewQuestionDetail',
  components: {
    QuestionFormDialog,
    OccurrenceDialog,
    ResumeMarkdownPreview
  }
})
export default class extends Vue {
  private questionId = ''
  private questionDialogVisible = false
  private occurrenceDialogVisible = false
  private questionDifficultyLabelMap = questionDifficultyLabelMap
  private questionTypeLabelMap = questionTypeLabelMap
  private masteryStatusLabelMap = masteryStatusLabelMap
  private answerPerformanceLabelMap = answerPerformanceLabelMap
  private interviewStageLabelMap = interviewStageLabelMap
  private reviewResultLabelMap = reviewResultLabelMap

  get question() {
    return InterviewQuestionBankModule.currentQuestion
  }

  get loading() {
    return InterviewQuestionBankModule.loading
  }

  get submitting() {
    return InterviewQuestionBankModule.submitting
  }

  get applications() {
    return ResumeApplicationModule.applications
  }

  get progressList() {
    return InterviewProgressModule.progressList
  }

  async created() {
    this.questionId = this.$route.params.id
    await this.loadData()
  }

  private async loadData() {
    await InterviewQuestionBankModule.GetInterviewQuestionById(this.questionId)
    await ResumeApplicationModule.GetResumeApplications({ page: 1, pageSize: 100 })
  }

  private goBack() {
    this.$router.push('/interview-question-bank/list')
  }

  private openEditDialog() {
    this.questionDialogVisible = true
  }

  private openOccurrenceDialog() {
    this.occurrenceDialogVisible = true
  }

  private async loadProgressForApplication(applicationId: string) {
    if (!applicationId) return
    await InterviewProgressModule.GetInterviewProgressList({ applicationId })
  }

  private async submitQuestion(payload: CreateInterviewQuestionRequest) {
    await InterviewQuestionBankModule.UpdateInterviewQuestion({
      id: this.questionId,
      data: payload
    })
    await this.loadData()
    this.questionDialogVisible = false
    this.$message.success('题目已更新')
  }

  private async submitOccurrence(payload: CreateQuestionOccurrenceRequest) {
    await InterviewQuestionBankModule.CreateQuestionOccurrence(payload)
    await this.loadData()
    this.occurrenceDialogVisible = false
    this.$message.success('面试提问记录已保存')
  }

  private async toggleFavorite() {
    if (!this.question) return
    await InterviewQuestionBankModule.ToggleFavorite({
      id: this.question.id,
      isFavorite: !this.question.isFavorite
    })
    await this.loadData()
  }

  private async deleteOccurrence(id: string) {
    await this.$confirm('确定删除这条面试提问记录吗？', '删除面试提问记录', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await InterviewQuestionBankModule.DeleteQuestionOccurrence({
      id,
      questionId: this.questionId
    })
    await this.loadData()
    this.$message.success('删除成功')
  }

  private goApplication(applicationId: string) {
    this.$router.push(`/resume/application/${applicationId}`)
  }

  private formatDateTime(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }

  private renderMarkdown(content: string) {
    return renderResumeMarkdown(content || '')
  }

  private difficultyTagType(difficulty: InterviewQuestionDifficulty) {
    if (difficulty === InterviewQuestionDifficulty.EASY) return 'success'
    if (difficulty === InterviewQuestionDifficulty.HARD) return 'danger'
    return 'warning'
  }

  private masteryTagType(status: QuestionMasteryStatus) {
    if (status === QuestionMasteryStatus.MASTERED) return 'success'
    if (status === QuestionMasteryStatus.WEAK) return 'danger'
    if (status === QuestionMasteryStatus.NORMAL) return 'warning'
    return 'info'
  }
}
</script>

<style lang="scss" scoped>
.question-detail-container {
  padding: 20px;

  .detail-card {
    margin-top: 20px;
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 20px;
    margin-bottom: 12px;

    h2 {
      margin: 0 0 10px;
      color: #303133;
    }
  }

  .meta-line,
  .tag-list,
  .header-actions,
  .card-header {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
  }

  .tag-list {
    margin-bottom: 20px;
  }

  .content-card {
    margin-bottom: 20px;
    min-height: 220px;
  }

  .edit-hint {
    color: #909399;
    font-size: 12px;
  }

  .question-markdown-preview {
    cursor: text;
  }

  .sub-text {
    color: #909399;
    font-size: 12px;
    margin-top: 4px;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .summary-item {
    padding: 12px;
    border: 1px solid #ebeef5;
    border-radius: 4px;

    span {
      display: block;
      color: #909399;
      margin-bottom: 6px;
    }

    strong {
      font-size: 22px;
      color: #303133;
    }

    &.danger strong {
      color: #f56c6c;
    }

    &.warning strong {
      color: #e6a23c;
    }

    &.success strong {
      color: #67c23a;
    }
  }

  .review-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    color: #606266;
  }

  .danger {
    color: #f56c6c;
  }
}
</style>
