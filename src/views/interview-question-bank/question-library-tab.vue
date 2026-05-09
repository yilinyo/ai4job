<template>
  <!-- eslint-disable vue/no-v-html -->
  <div class="question-library-tab">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="queryParams" class="filter-form">
        <el-form-item>
          <el-input
            v-model="queryParams.keyword"
            placeholder="搜索题目、答案或标签"
            clearable
            @clear="handleSearch"
            @keyup.enter.native="handleSearch"
          >
            <i slot="prefix" class="el-input__icon el-icon-search" />
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="selectedTags"
            multiple
            collapse-tags
            filterable
            clearable
            placeholder="标签"
            @change="handleSearch"
          >
            <el-option
              v-for="tag in questionTagSuggestions"
              :key="tag"
              :label="tag"
              :value="tag"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="queryParams.difficulty" clearable placeholder="难度" @change="handleSearch">
            <el-option
              v-for="item in questionDifficultyOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="queryParams.questionType" clearable placeholder="题型" @change="handleSearch">
            <el-option
              v-for="item in questionTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="queryParams.masteryStatus" clearable placeholder="掌握状态" @change="handleSearch">
            <el-option
              v-for="item in masteryStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="favoriteFilter" clearable placeholder="收藏状态" @change="handleSearch">
            <el-option label="已收藏" value="true" />
            <el-option label="未收藏" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item class="filter-actions">
          <el-button type="primary" icon="el-icon-search" @click="handleSearch">搜索</el-button>
          <el-button icon="el-icon-refresh-left" @click="resetSearch">重置</el-button>
          <el-button type="success" icon="el-icon-plus" @click="openCreateDialog">新增题目</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      v-loading="loading"
      :data="questions"
      class="question-table"
      border
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column label="收藏" width="70" align="center">
        <template slot-scope="{row}">
          <el-button
            type="text"
            class="favorite-button"
            :class="{'is-favorite': row.isFavorite}"
            @click.stop="toggleFavorite(row)"
          >
            <i :class="row.isFavorite ? 'el-icon-star-on' : 'el-icon-star-off'" />
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="题目" min-width="280">
        <template slot-scope="{row}">
          <div class="question-title">{{ row.title }}</div>
          <div class="question-content" v-html="renderMarkdown(row.content)" />
          <div class="tag-list">
            <el-tag v-for="tag in row.tags" :key="tag" size="mini" effect="plain">{{ tag }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="难度" width="90">
        <template slot-scope="{row}">
          <el-tag :type="difficultyTagType(row.difficulty)" size="small">
            {{ questionDifficultyLabelMap[row.difficulty] }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="题型" width="110">
        <template slot-scope="{row}">
          {{ questionTypeLabelMap[row.questionType] }}
        </template>
      </el-table-column>
      <el-table-column label="掌握" width="100">
        <template slot-scope="{row}">
          <el-tag :type="masteryTagType(row.masteryStatus)" size="small">
            {{ masteryStatusLabelMap[row.masteryStatus] }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="提问" width="80" align="center">
        <template slot-scope="{row}">
          <el-badge :value="row.occurrenceCount" class="occurrence-badge" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="120">
        <template slot-scope="{row}">
          {{ formatDate(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="最近复习" width="120">
        <template slot-scope="{row}">
          {{ row.lastReviewedAt ? formatDate(row.lastReviewedAt) : '未复习' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template slot-scope="{row}">
          <el-button type="text" size="small" @click.stop="viewQuestion(row.id)">详情</el-button>
          <el-button type="text" size="small" @click.stop="openEditDialog(row)">编辑</el-button>
          <el-button type="text" size="small" @click.stop="openOccurrenceDialog(row.id)">关联提问</el-button>
          <el-button type="text" size="small" class="danger" @click.stop="deleteQuestion(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty
      v-if="!loading && questions.length === 0"
      description="暂无面试题，点击新增题目开始沉淀"
    />

    <el-pagination
      v-if="total > 0"
      :current-page="queryParams.page"
      :page-size="queryParams.pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
    />

    <question-form-dialog
      :visible.sync="questionDialogVisible"
      :title="questionDialogTitle"
      :loading="submitting"
      :question-data="editingQuestion"
      @submit="submitQuestion"
    />

    <occurrence-dialog
      :visible.sync="occurrenceDialogVisible"
      title="新增面试提问记录"
      :loading="submitting"
      :questions="questions"
      :applications="applications"
      :progress-options="progressList"
      :question-id="occurrenceQuestionId"
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
import { InterviewQuestionBankModule } from '@/store/modules/interview-question-bank'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import { InterviewProgressModule } from '@/store/modules/interview-progress'
import {
  CreateInterviewQuestionRequest,
  CreateQuestionOccurrenceRequest,
  InterviewQuestion,
  InterviewQuestionDifficulty,
  InterviewQuestionQueryParams,
  QuestionMasteryStatus
} from '@/models'
import {
  masteryStatusLabelMap,
  masteryStatusOptions,
  questionDifficultyLabelMap,
  questionDifficultyOptions,
  questionTagSuggestions,
  questionTypeLabelMap,
  questionTypeOptions
} from './constants'
import { renderResumeMarkdown } from '@/utils/resume-markdown'

@Component({
  name: 'QuestionLibraryTab',
  components: {
    QuestionFormDialog,
    OccurrenceDialog
  }
})
export default class extends Vue {
  private questionDifficultyOptions = questionDifficultyOptions
  private questionTypeOptions = questionTypeOptions
  private masteryStatusOptions = masteryStatusOptions
  private questionTagSuggestions = questionTagSuggestions
  private questionDifficultyLabelMap = questionDifficultyLabelMap
  private questionTypeLabelMap = questionTypeLabelMap
  private masteryStatusLabelMap = masteryStatusLabelMap

  private selectedTags: string[] = []
  private favoriteFilter = ''
  private questionDialogVisible = false
  private occurrenceDialogVisible = false
  private questionDialogTitle = '新增题目'
  private editingQuestion: InterviewQuestion | null = null
  private occurrenceQuestionId = ''

  private queryParams: InterviewQuestionQueryParams = {
    keyword: '',
    tags: '',
    difficulty: undefined,
    questionType: undefined,
    masteryStatus: undefined,
    isFavorite: undefined,
    page: 1,
    pageSize: 10
  }

  get questions() {
    return InterviewQuestionBankModule.questions
  }

  get total() {
    return InterviewQuestionBankModule.total
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

  created() {
    this.loadQuestions()
    ResumeApplicationModule.GetResumeApplications({ page: 1, pageSize: 100 })
  }

  private buildQuery() {
    return {
      ...this.queryParams,
      tags: this.selectedTags.join(','),
      isFavorite: this.favoriteFilter === '' ? undefined : this.favoriteFilter === 'true'
    }
  }

  private async loadQuestions() {
    await InterviewQuestionBankModule.GetInterviewQuestions(this.buildQuery())
  }

  private handleSearch() {
    this.queryParams.page = 1
    this.loadQuestions()
  }

  private resetSearch() {
    this.selectedTags = []
    this.favoriteFilter = ''
    this.queryParams = {
      keyword: '',
      tags: '',
      difficulty: undefined,
      questionType: undefined,
      masteryStatus: undefined,
      isFavorite: undefined,
      page: 1,
      pageSize: 10
    }
    this.loadQuestions()
  }

  private handlePageChange(page: number) {
    this.queryParams.page = page
    this.loadQuestions()
  }

  private openCreateDialog() {
    this.questionDialogTitle = '新增题目'
    this.editingQuestion = null
    this.questionDialogVisible = true
  }

  private openEditDialog(question: InterviewQuestion) {
    this.questionDialogTitle = '编辑题目'
    this.editingQuestion = question
    this.questionDialogVisible = true
  }

  private async submitQuestion(payload: CreateInterviewQuestionRequest) {
    try {
      if (this.editingQuestion) {
        await InterviewQuestionBankModule.UpdateInterviewQuestion({
          id: this.editingQuestion.id,
          data: payload
        })
        this.$message.success('题目已更新')
      } else {
        await InterviewQuestionBankModule.CreateInterviewQuestion(payload)
        this.$message.success('题目已创建')
      }
      this.questionDialogVisible = false
      this.loadQuestions()
    } catch (error) {
      console.error('保存题目失败:', error)
      this.$message.error('保存失败，请稍后重试')
    }
  }

  private viewQuestion(id: string) {
    this.$router.push(`/interview-question-bank/question/${id}`)
  }

  private handleRowClick(row: InterviewQuestion) {
    this.viewQuestion(row.id)
  }

  private async toggleFavorite(question: InterviewQuestion) {
    await InterviewQuestionBankModule.ToggleFavorite({
      id: question.id,
      isFavorite: !question.isFavorite
    })
  }

  private async deleteQuestion(question: InterviewQuestion) {
    try {
      const message = question.occurrenceCount > 0
        ? `确定删除「${question.title}」吗？删除后会同步清理 ${question.occurrenceCount} 条面试提问记录和相关复习记录。`
        : `确定删除「${question.title}」吗？删除后无法恢复。`
      await this.$confirm(message, '删除题目', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
      await InterviewQuestionBankModule.DeleteInterviewQuestion(question.id)
      this.$message.success('删除成功')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除题目失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
  }

  private openOccurrenceDialog(questionId: string) {
    this.occurrenceQuestionId = questionId
    this.occurrenceDialogVisible = true
  }

  private async loadProgressForApplication(applicationId: string) {
    if (!applicationId) return
    await InterviewProgressModule.GetInterviewProgressList({ applicationId })
  }

  private async submitOccurrence(payload: CreateQuestionOccurrenceRequest) {
    try {
      await InterviewQuestionBankModule.CreateQuestionOccurrence(payload)
      this.$message.success('面试提问记录已保存')
      this.occurrenceDialogVisible = false
      this.loadQuestions()
    } catch (error) {
      console.error('保存面试提问记录失败:', error)
      this.$message.error('保存失败，请稍后重试')
    }
  }

  private renderMarkdown(content: string) {
    return renderResumeMarkdown(content || '')
  }

  private formatDate(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
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
.question-library-tab {
  .search-card {
    margin-bottom: 16px;
  }

  .filter-form {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;

    ::v-deep .el-form-item {
      margin-right: 8px;
      margin-bottom: 10px;
    }
  }

  .question-table {
    margin-bottom: 16px;

    ::v-deep .el-table__row {
      cursor: pointer;
    }
  }

  .question-title {
    font-weight: 600;
    color: #303133;
    margin-bottom: 6px;
  }

  .question-content {
    color: #606266;
    font-size: 13px;
    line-height: 1.5;
    margin-bottom: 8px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .favorite-button {
    font-size: 20px;
    color: #c0c4cc;

    &.is-favorite {
      color: #e6a23c;
    }
  }

  .danger {
    color: #f56c6c;
  }

  .el-pagination {
    text-align: center;
    margin-top: 20px;
  }
}
</style>
