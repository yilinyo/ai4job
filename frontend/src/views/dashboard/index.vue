<template>
  <div class="job-dashboard-container">
    <div class="dashboard-hero">
      <div>
        <h1>求职工作台</h1>
        <p>集中查看投递进展、跟进事项、岗位分布和最近动态。</p>
      </div>
      <el-radio-group v-model="selectedRange" size="small" @change="handleRangeChange">
        <el-radio-button label="7d">近 7 天</el-radio-button>
        <el-radio-button label="30d">近 30 天</el-radio-button>
        <el-radio-button label="90d">近 90 天</el-radio-button>
        <el-radio-button label="all">全部</el-radio-button>
      </el-radio-group>
    </div>

    <div v-loading="loading">
      <metric-cards :metrics="metrics" />

      <pipeline-chart
        :pipeline="statusPipeline"
        :trend="applicationTrend"
      />

      <el-row :gutter="20" class="dashboard-main-row">
        <el-col :xs="24" :lg="15">
          <todo-panel
            :todos="todos"
            :submitting="todoSubmitting"
            :submit-success-key="todoSubmitSuccessKey"
            @create="handleCreateTodo"
            @update="handleUpdateTodo"
            @complete="handleCompleteTodo"
            @delete="handleDeleteTodo"
            @open="handleOpenRoute"
          />
        </el-col>
        <el-col :xs="24" :lg="9">
          <activity-feed
            :activities="activities"
            @open="handleOpenRoute"
          />
        </el-col>
      </el-row>

      <distribution-charts
        :companies="companyRanking"
        :jobs="jobTitleDistribution"
        :sources="sourceDistribution"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import MetricCards from './components/MetricCards.vue'
import PipelineChart from './components/PipelineChart.vue'
import TodoPanel from './components/TodoPanel.vue'
import DistributionCharts from './components/DistributionCharts.vue'
import ActivityFeed from './components/ActivityFeed.vue'
import { DashboardModule } from '@/store/modules/dashboard'
import {
  ActivityLog,
  CreateDashboardTodoRequest,
  DashboardMetrics,
  DashboardNameCount,
  DashboardRange,
  DashboardResponse,
  DashboardStatusBucket,
  DashboardTodo,
  DashboardTrendPoint,
  UpdateDashboardTodoRequest
} from '@/models'

const emptyMetrics: DashboardMetrics = {
  totalApplications: 0,
  interviewing: 0,
  pendingFollowUp: 0,
  offers: 0
}

@Component({
  name: 'Dashboard',
  components: {
    MetricCards,
    PipelineChart,
    TodoPanel,
    DistributionCharts,
    ActivityFeed
  }
})
export default class extends Vue {
  private selectedRange: DashboardRange = '30d'
  private todoSubmitting = false
  private todoSubmitSuccessKey = 0

  get loading() {
    return DashboardModule.loading
  }

  get dashboard(): DashboardResponse | null {
    return DashboardModule.dashboard
  }

  get metrics(): DashboardMetrics {
    return this.dashboard ? this.dashboard.metrics : emptyMetrics
  }

  get statusPipeline(): DashboardStatusBucket[] {
    return this.dashboard ? this.dashboard.statusPipeline : []
  }

  get applicationTrend(): DashboardTrendPoint[] {
    return this.dashboard ? this.dashboard.applicationTrend : []
  }

  get companyRanking(): DashboardNameCount[] {
    return this.dashboard ? this.dashboard.companyRanking : []
  }

  get jobTitleDistribution(): DashboardNameCount[] {
    return this.dashboard ? this.dashboard.jobTitleDistribution : []
  }

  get sourceDistribution(): DashboardNameCount[] {
    return this.dashboard ? this.dashboard.sourceDistribution : []
  }

  get todos(): DashboardTodo[] {
    return DashboardModule.todos
  }

  get activities(): ActivityLog[] {
    return this.dashboard ? this.dashboard.activities : []
  }

  created() {
    this.loadDashboard()
  }

  private async loadDashboard() {
    try {
      await DashboardModule.GetDashboard(this.selectedRange)
    } catch (error) {
      console.error('加载工作台失败:', error)
      this.$message.error('加载工作台失败，请稍后重试')
    }
  }

  private handleRangeChange() {
    this.loadDashboard()
  }

  private async handleCreateTodo(payload: CreateDashboardTodoRequest) {
    this.todoSubmitting = true
    try {
      await DashboardModule.CreateTodo(payload)
      this.todoSubmitSuccessKey += 1
      this.$message.success('待办创建成功')
    } catch (error) {
      console.error('创建待办失败:', error)
      this.$message.error('创建待办失败，请检查目标信息')
    } finally {
      this.todoSubmitting = false
    }
  }

  private async handleUpdateTodo(payload: { id: string, data: UpdateDashboardTodoRequest }) {
    this.todoSubmitting = true
    try {
      await DashboardModule.UpdateTodo(payload)
      this.todoSubmitSuccessKey += 1
      this.$message.success('待办更新成功')
    } catch (error) {
      console.error('更新待办失败:', error)
      this.$message.error('更新待办失败，请检查目标信息')
    } finally {
      this.todoSubmitting = false
    }
  }

  private async handleCompleteTodo(id: string) {
    try {
      await DashboardModule.CompleteTodo(id)
      this.$message.success('待办已完成')
    } catch (error) {
      console.error('完成待办失败:', error)
      this.$message.error('完成待办失败，请稍后重试')
    }
  }

  private async handleDeleteTodo(id: string) {
    try {
      await DashboardModule.DeleteTodo(id)
      this.$message.success('待办已删除')
    } catch (error) {
      console.error('删除待办失败:', error)
      this.$message.error('删除待办失败，请稍后重试')
    }
  }

  private handleOpenRoute(route: string) {
    if (!route) return
    this.$router.push(route).catch((error: Error) => {
      if (error.name !== 'NavigationDuplicated') {
        console.error('打开详情失败:', error)
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.job-dashboard-container {
  min-height: calc(100vh - 84px);
  padding: 24px;
  background: #f5f7fa;
}

.dashboard-hero {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  margin-bottom: 20px;
  padding: 24px;
  border-radius: 12px;
  background: linear-gradient(135deg, #1f2937 0%, #2563eb 100%);
  color: #fff;

  h1 {
    margin: 0 0 10px;
    font-size: 28px;
    font-weight: 600;
  }

  p {
    margin: 0;
    color: rgba(255, 255, 255, 0.82);
  }

  ::v-deep .el-radio-button__inner {
    border-color: rgba(255, 255, 255, 0.4);
  }
}

.dashboard-main-row {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .job-dashboard-container {
    padding: 16px;
  }

  .dashboard-hero {
    flex-direction: column;

    h1 {
      font-size: 24px;
    }
  }
}
</style>
