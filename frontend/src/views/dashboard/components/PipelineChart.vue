<template>
  <el-card shadow="never" class="pipeline-chart-card">
    <div slot="header" class="card-header">
      <span>投递管线与趋势</span>
    </div>
    <empty-state
      v-if="isEmpty"
      description="暂无投递管线和趋势数据"
    />
    <div v-else class="pipeline-chart-card__body">
      <div class="chart-section">
        <div class="chart-section__title">状态管线</div>
        <div ref="pipelineChart" class="chart" />
      </div>
      <div class="chart-section">
        <div class="chart-section__title">投递趋势</div>
        <div ref="trendChart" class="chart" />
      </div>
    </div>
  </el-card>
</template>

<script lang="ts">
import * as echarts from 'echarts'
import 'echarts/theme/macarons.js'
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { ECharts } from 'echarts'
import { DashboardStatusBucket, DashboardTrendPoint } from '@/models'

@Component({
  name: 'PipelineChart'
})
export default class extends Vue {
  @Prop({ required: true, default: () => [] }) private pipeline!: DashboardStatusBucket[]
  @Prop({ required: true, default: () => [] }) private trend!: DashboardTrendPoint[]

  private pipelineChart: ECharts | null = null
  private trendChart: ECharts | null = null

  get hasPipelineData() {
    return this.pipeline.some(item => item.count > 0)
  }

  get hasTrendData() {
    return this.trend.some(item => item.count > 0)
  }

  get isEmpty() {
    return !this.hasPipelineData && !this.hasTrendData
  }

  @Watch('pipeline', { deep: true })
  private onPipelineChange() {
    this.refreshCharts()
  }

  @Watch('trend', { deep: true })
  private onTrendChange() {
    this.refreshCharts()
  }

  mounted() {
    this.$nextTick(() => {
      this.initCharts()
      window.addEventListener('resize', this.resizeCharts)
    })
  }

  beforeDestroy() {
    window.removeEventListener('resize', this.resizeCharts)
    this.disposeCharts()
  }

  private initCharts() {
    if (this.isEmpty) return
    if (!this.pipelineChart && this.$refs.pipelineChart) {
      this.pipelineChart = echarts.init(this.$refs.pipelineChart as HTMLDivElement, 'macarons')
    }
    if (!this.trendChart && this.$refs.trendChart) {
      this.trendChart = echarts.init(this.$refs.trendChart as HTMLDivElement, 'macarons')
    }
    this.setOptions()
  }

  private refreshCharts() {
    this.$nextTick(() => {
      if (this.isEmpty) {
        this.disposeCharts()
        return
      }
      this.initCharts()
    })
  }

  private setOptions() {
    if (this.pipelineChart) {
      this.pipelineChart.setOption({
        tooltip: {
          trigger: 'axis',
          axisPointer: { type: 'shadow' }
        },
        grid: {
          top: 20,
          right: 20,
          bottom: 30,
          left: 36,
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: this.pipeline.map(item => item.label),
          axisTick: { alignWithLabel: true },
          axisLabel: { interval: 0 }
        },
        yAxis: {
          type: 'value',
          minInterval: 1
        },
        series: [{
          name: '投递数量',
          type: 'bar',
          barMaxWidth: 36,
          data: this.pipeline.map(item => item.count),
          itemStyle: {
            color: '#409eff',
            borderRadius: [6, 6, 0, 0]
          }
        }]
      }, true)
    }

    if (this.trendChart) {
      this.trendChart.setOption({
        tooltip: {
          trigger: 'axis'
        },
        grid: {
          top: 20,
          right: 20,
          bottom: 30,
          left: 36,
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: this.trend.map(item => item.date)
        },
        yAxis: {
          type: 'value',
          minInterval: 1
        },
        series: [{
          name: '新增投递',
          type: 'line',
          smooth: true,
          symbolSize: 8,
          data: this.trend.map(item => item.count),
          lineStyle: { width: 3, color: '#67c23a' },
          itemStyle: { color: '#67c23a' },
          areaStyle: { color: 'rgba(103, 194, 58, 0.14)' }
        }]
      }, true)
    }
  }

  private resizeCharts = () => {
    if (this.pipelineChart) this.pipelineChart.resize()
    if (this.trendChart) this.trendChart.resize()
  }

  private disposeCharts() {
    if (this.pipelineChart) {
      this.pipelineChart.dispose()
      this.pipelineChart = null
    }
    if (this.trendChart) {
      this.trendChart.dispose()
      this.trendChart = null
    }
  }
}
</script>

<style lang="scss" scoped>
.pipeline-chart-card {
  margin-bottom: 20px;
}

.card-header {
  font-weight: 600;
  color: #303133;
}

.pipeline-chart-card__body {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.chart-section__title {
  margin-bottom: 12px;
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

.chart {
  width: 100%;
  height: 300px;
}

@media (max-width: 992px) {
  .pipeline-chart-card__body {
    grid-template-columns: 1fr;
  }
}
</style>
