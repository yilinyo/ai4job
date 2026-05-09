<template>
  <div class="interview-question-bank-container">
    <div class="module-header">
      <div>
        <h2>面试题库</h2>
        <p>把真实面试问题沉淀为可筛选、可复习、可追踪的个人题库。</p>
      </div>
    </div>

    <el-card class="module-card">
      <question-library-tab />
    </el-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import QuestionLibraryTab from './question-library-tab.vue'

@Component({
  name: 'InterviewQuestionBankIndex',
  components: {
    QuestionLibraryTab
  }
})
export default class extends Vue {}
</script>

<style lang="scss" scoped>
.interview-question-bank-container {
  padding: 20px;

  .module-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0 0 8px;
      font-size: 24px;
      color: #303133;
    }

    p {
      margin: 0;
      color: #909399;
      font-size: 14px;
    }
  }

  .module-card {
    min-height: 640px;
  }
}
</style>
