<template>
  <!-- eslint-disable vue/no-v-html -->
  <div class="flashcard-review-tab">
    <el-card shadow="never" class="review-toolbar">
      <el-form :inline="true" class="review-form">
        <el-form-item label="复习模式">
          <el-select v-model="reviewMode" placeholder="选择模式" @change="handleModeChange">
            <el-option
              v-for="item in flashcardModeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="reviewMode === FlashcardReviewMode.TAG" label="标签">
          <el-select v-model="selectedTags" multiple collapse-tags filterable clearable placeholder="选择标签">
            <el-option
              v-for="tag in questionTagSuggestions"
              :key="tag"
              :label="tag"
              :value="tag"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="reviewMode === FlashcardReviewMode.DIFFICULTY" label="难度">
          <el-select v-model="selectedDifficulty" clearable placeholder="选择难度">
            <el-option
              v-for="item in questionDifficultyOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="limit" :min="1" :max="20" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-video-play" :loading="loading" @click="loadFlashcards">
            开始复习
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div v-if="currentQuestion" class="flashcard-workspace">
      <div class="progress-line">
        <span>{{ currentIndex + 1 }} / {{ queue.length }}</span>
        <el-progress :percentage="progressPercentage" :show-text="false" />
      </div>

      <div
        class="flashcard-shell"
        :class="[
          {'is-flipped': isFlipped},
          swipeDirection ? `is-swipe-${swipeDirection}` : ''
        ]"
        @click="flipCard"
        @touchstart.passive="handleTouchStart"
        @touchend.passive="handleTouchEnd"
      >
        <div class="flashcard-inner">
          <div class="flashcard-face flashcard-front">
            <div class="card-meta">
              <el-tag :type="difficultyTagType(currentQuestion.difficulty)" size="small">
                {{ questionDifficultyLabelMap[currentQuestion.difficulty] }}
              </el-tag>
              <el-tag size="small" type="info">{{ questionTypeLabelMap[currentQuestion.questionType] }}</el-tag>
              <el-tag size="small" :type="masteryTagType(currentQuestion.masteryStatus)">
                {{ masteryStatusLabelMap[currentQuestion.masteryStatus] }}
              </el-tag>
            </div>
            <h3>{{ currentQuestion.title }}</h3>
            <div class="question-text" v-html="renderMarkdown(currentQuestion.content)" />
            <div class="tag-list">
              <el-tag v-for="tag in currentQuestion.tags" :key="tag" size="mini" effect="plain">{{ tag }}</el-tag>
            </div>
          </div>
          <div class="flashcard-face flashcard-back">
            <h3>答案解析</h3>
            <div class="answer-text" v-html="renderMarkdown(currentQuestion.answerAnalysis)" />
            <div class="review-actions" @click.stop>
              <el-button
                v-for="item in flashcardResultActions"
                :key="item.value"
                :type="item.type"
                :loading="saving && pendingResult === item.value"
                @click="saveReview(item.value)"
              >
                {{ item.label }}
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <div class="card-controls">
        <el-button icon="el-icon-arrow-left" :disabled="!hasPrevious || isAnimating" @click="moveCard(-1)">上一题</el-button>
        <el-button icon="el-icon-refresh" :disabled="isAnimating" @click="flipCard">翻转</el-button>
        <el-button type="primary" :disabled="!hasNext || isAnimating" @click="moveCard(1)">
          下一题
          <i class="el-icon-arrow-right el-icon--right" />
        </el-button>
      </div>

      <el-card v-if="summary.total > 0" shadow="never" class="summary-card">
        <div class="summary-item">
          <span>已复习</span>
          <strong>{{ summary.total }}</strong>
        </div>
        <div class="summary-item danger">
          <span>不会</span>
          <strong>{{ summary.failedCount }}</strong>
        </div>
        <div class="summary-item warning">
          <span>模糊</span>
          <strong>{{ summary.vagueCount }}</strong>
        </div>
        <div class="summary-item success">
          <span>掌握</span>
          <strong>{{ summary.masteredCount }}</strong>
        </div>
      </el-card>
    </div>

    <el-empty
      v-else
      :description="loading ? '正在抽题' : '选择复习模式后开始抽题'"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { FlashcardReviewModule } from '@/store/modules/flashcard-review'
import {
  FlashcardReviewMode,
  FlashcardReviewResult,
  InterviewQuestionDifficulty,
  QuestionMasteryStatus
} from '@/models'
import {
  flashcardModeOptions,
  flashcardResultActions,
  masteryStatusLabelMap,
  questionDifficultyLabelMap,
  questionDifficultyOptions,
  questionTagSuggestions,
  questionTypeLabelMap
} from './constants'
import { renderResumeMarkdown } from '@/utils/resume-markdown'

@Component({
  name: 'FlashcardReviewTab'
})
export default class extends Vue {
  private FlashcardReviewMode = FlashcardReviewMode
  private flashcardModeOptions = flashcardModeOptions
  private flashcardResultActions = flashcardResultActions
  private questionTagSuggestions = questionTagSuggestions
  private questionDifficultyOptions = questionDifficultyOptions
  private questionDifficultyLabelMap = questionDifficultyLabelMap
  private questionTypeLabelMap = questionTypeLabelMap
  private masteryStatusLabelMap = masteryStatusLabelMap

  private reviewMode = FlashcardReviewMode.RANDOM
  private selectedTags: string[] = []
  private selectedDifficulty: InterviewQuestionDifficulty | undefined = undefined
  private limit = 8
  private isFlipped = false
  private isAnimating = false
  private swipeDirection = ''
  private touchStartX = 0
  private pendingResult: FlashcardReviewResult | '' = ''

  get queue() {
    return FlashcardReviewModule.queue
  }

  get currentQuestion() {
    return FlashcardReviewModule.currentQuestion
  }

  get currentIndex() {
    return FlashcardReviewModule.currentIndex
  }

  get hasPrevious() {
    return FlashcardReviewModule.hasPrevious
  }

  get hasNext() {
    return FlashcardReviewModule.hasNext
  }

  get loading() {
    return FlashcardReviewModule.loading
  }

  get saving() {
    return FlashcardReviewModule.saving
  }

  get summary() {
    return FlashcardReviewModule.summary
  }

  get progressPercentage() {
    if (!this.queue.length) return 0
    return Math.round(((this.currentIndex + 1) / this.queue.length) * 100)
  }

  created() {
    this.loadFlashcards()
  }

  private async loadFlashcards() {
    await FlashcardReviewModule.SetReviewMode(this.reviewMode)
    await FlashcardReviewModule.SetSelectedTags(this.selectedTags)
    await FlashcardReviewModule.SetSelectedDifficulty(this.selectedDifficulty)
    await FlashcardReviewModule.LoadFlashcards({
      mode: this.reviewMode,
      tags: this.selectedTags.join(','),
      difficulty: this.selectedDifficulty,
      limit: this.limit
    })
    this.isFlipped = false
    this.swipeDirection = ''
  }

  private handleModeChange() {
    if (this.reviewMode !== FlashcardReviewMode.TAG) {
      this.selectedTags = []
    }
    if (this.reviewMode !== FlashcardReviewMode.DIFFICULTY) {
      this.selectedDifficulty = undefined
    }
  }

  private flipCard() {
    if (this.isAnimating || !this.currentQuestion) return
    this.isFlipped = !this.isFlipped
  }

  private moveCard(direction: number) {
    if (this.isAnimating) return
    const nextIndex = this.currentIndex + direction
    if (nextIndex < 0 || nextIndex >= this.queue.length) return

    this.isAnimating = true
    this.swipeDirection = direction > 0 ? 'left' : 'right'
    window.setTimeout(() => {
      FlashcardReviewModule.SetCurrentIndex(nextIndex)
      this.isFlipped = false
      this.swipeDirection = ''
      window.setTimeout(() => {
        this.isAnimating = false
      }, 60)
    }, 260)
  }

  private handleTouchStart(event: TouchEvent) {
    this.touchStartX = event.changedTouches[0].clientX
  }

  private handleTouchEnd(event: TouchEvent) {
    const deltaX = event.changedTouches[0].clientX - this.touchStartX
    if (Math.abs(deltaX) < 45) return
    this.moveCard(deltaX < 0 ? 1 : -1)
  }

  private async saveReview(result: FlashcardReviewResult) {
    if (!this.currentQuestion || this.saving) return
    this.pendingResult = result
    await FlashcardReviewModule.SaveReview({
      questionId: this.currentQuestion.id,
      result,
      reviewMode: this.reviewMode
    })
    this.pendingResult = ''

    if (this.hasNext) {
      this.moveCard(1)
    } else {
      this.isFlipped = false
      this.$message.success('本组闪卡已完成')
    }
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
.flashcard-review-tab {
  .review-toolbar {
    margin-bottom: 20px;
  }

  .review-form {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    ::v-deep .el-form-item {
      margin-bottom: 8px;
      margin-right: 8px;
    }
  }

  .flashcard-workspace {
    max-width: 820px;
    margin: 0 auto;
  }

  .progress-line {
    display: grid;
    grid-template-columns: 80px 1fr;
    align-items: center;
    gap: 12px;
    margin-bottom: 18px;
    color: #606266;
  }

  .flashcard-shell {
    width: min(100%, 760px);
    height: 420px;
    margin: 0 auto;
    perspective: 1200px;
    cursor: pointer;
    transition: transform 260ms ease, opacity 260ms ease;

    &.is-swipe-left {
      transform: translateX(-90px) scale(0.96);
      opacity: 0.25;
    }

    &.is-swipe-right {
      transform: translateX(90px) scale(0.96);
      opacity: 0.25;
    }

    &.is-flipped .flashcard-inner {
      transform: rotateY(180deg);
    }
  }

  .flashcard-inner {
    position: relative;
    width: 100%;
    height: 100%;
    transform-style: preserve-3d;
    transition: transform 340ms ease;
  }

  .flashcard-face {
    position: absolute;
    inset: 0;
    padding: 28px;
    border: 1px solid #dcdfe6;
    border-radius: 8px;
    background: #fff;
    box-shadow: 0 12px 28px rgba(0, 0, 0, 0.08);
    backface-visibility: hidden;
    overflow: hidden;
    display: flex;
    flex-direction: column;

    h3 {
      margin: 16px 0;
      color: #303133;
      font-size: 22px;
    }
  }

  .flashcard-back {
    transform: rotateY(180deg);
  }

  .card-meta,
  .tag-list,
  .review-actions,
  .card-controls {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
  }

  .question-text,
  .answer-text {
    flex: 1;
    color: #303133;
    font-size: 16px;
    line-height: 1.8;
    overflow-y: auto;
    white-space: pre-wrap;
  }

  .answer-text {
    margin-bottom: 18px;
  }

  .card-controls {
    justify-content: center;
    margin: 20px 0;
  }

  .summary-card {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;

    ::v-deep .el-card__body {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 12px;
    }
  }

  .summary-item {
    border: 1px solid #ebeef5;
    border-radius: 4px;
    padding: 12px;
    text-align: center;

    span {
      display: block;
      color: #909399;
      margin-bottom: 6px;
    }

    strong {
      color: #303133;
      font-size: 22px;
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
}

@media (max-width: 768px) {
  .flashcard-review-tab {
    .flashcard-shell {
      height: 480px;
    }

    .flashcard-face {
      padding: 20px;
    }

    .summary-card ::v-deep .el-card__body {
      grid-template-columns: repeat(2, 1fr);
    }
  }
}
</style>
