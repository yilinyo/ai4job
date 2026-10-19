<template>
  <el-row :gutter="20" class="distribution-charts">
    <el-col :xs="24" :lg="8">
      <el-card shadow="never" class="distribution-card">
        <div slot="header" class="card-header">公司投递 Top N</div>
        <empty-state
          v-if="!hasCompanies"
          description="暂无公司投递数据"
        />
        <div v-else ref="companyChart" class="chart" />
      </el-card>
    </el-col>
    <el-col :xs="24" :lg="8">
      <el-card shadow="never" class="distribution-card">
        <div slot="header" class="card-header">岗位方向分布</div>
        <empty-state
          v-if="!hasJobs"
          description="暂无岗位方向数据"
        />
        <div v-else ref="jobChart" class="chart" />
      </el-card>
    </el-col>
    <el-col :xs="24" :lg="8">
      <el-card shadow="never" class="distribution-card">
        <div slot="header" class="card-header">渠道来源分布</div>
        <empty-state
          v-if="!hasSources"
          description="暂无渠道来源数据"
        />
        <div v-else ref="sourceChart" class="chart" />
      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import * as echarts from 'echarts'
import 'echarts/theme/macarons.js'
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { ECharts } from 'echarts'
import { DashboardNameCount } from '@/models'

@Component({
  name: 'DistributionCharts'
})
export default class extends Vue {
  @Prop({ required: true, default: () => [] }) private companies!: DashboardNameCount[]
  @Prop({ required: true, default: () => [] }) private jobs!: DashboardNameCount[]
  @Prop({ required: true, default: () => [] }) private sources!: DashboardNameCount[]

  private companyChart: ECharts | null = null
  private jobChart: ECharts | null = null
  private sourceChart: ECharts | null = null

  get hasCompanies() {
    return this.companies.some(item => item.count > 0)
  }

  get hasJobs() {
    return this.jobs.some(item => item.count > 0)
  }

  get hasSources() {
    return this.sources.some(item => item.count > 0)
  }

  @Watch('companies', { deep: true })
  @Watch('jobs', { deep: true })
  @Watch('sources', { deep: true })
  private onDataChange() {
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
    if (this.hasCompanies && !this.companyChart && this.$refs.companyChart) {
      this.companyChart = echarts.init(this.$refs.companyChart as HTMLDivElement, 'macarons')
    }
    if (this.hasJobs && !this.jobChart && this.$refs.jobChart) {
      this.jobChart = echarts.init(this.$refs.jobChart as HTMLDivElement, 'macarons')
    }
    if (this.hasSources && !this.sourceChart && this.$refs.sourceChart) {
      this.sourceChart = echarts.init(this.$refs.sourceChart as HTMLDivElement, 'macarons')
    }
    this.setOptions()
  }

  private refreshCharts() {
    this.$nextTick(() => {
      if (!this.hasCompanies && this.companyChart) {
        this.companyChart.dispose()
        this.companyChart = null
      }
      if (!this.hasJobs && this.jobChart) {
        this.jobChart.dispose()
        this.jobChart = null
      }
      if (!this.hasSources && this.sourceChart) {
        this.sourceChart.dispose()
        this.sourceChart = null
      }
      this.initCharts()
    })
  }

  private setOptions() {
    if (this.companyChart) {
      const sortedCompanies = [...this.companies].sort((a, b) => a.count - b.count)
      this.companyChart.setOption({
        tooltip: {
          trigger: 'axis',
          axisPointer: { type: 'shadow' }
        },
        grid: {
          top: 12,
          right: 16,
          bottom: 12,
          left: 80,
          containLabel: true
        },
        xAxis: {
          type: 'value',
          minInterval: 1
        },
        yAxis: {
          type: 'category',
          data: sortedCompanies.map(item => item.name),
          axisTick: { show: false }
        },
        series: [{
          name: '投递数量',
          type: 'bar',
          barMaxWidth: 18,
          data: sortedCompanies.map(item => item.count),
          itemStyle: {
            color: '#409eff',
            borderRadius: [0, 8, 8, 0]
          },
          label: {
            show: true,
            position: 'right'
          }
        }]
      }, true)
    }

    if (this.jobChart) {
      this.jobChart.setOption(this.createPieOptions('岗位方向', this.jobs), true)
    }

    if (this.sourceChart) {
      this.sourceChart.setOption(this.createPieOptions('渠道来源', this.sources), true)
    }
  }

  private createPieOptions(name: string, data: DashboardNameCount[]) {
    return {
      tooltip: {
        trigger: 'item'
      },
      legend: {
        bottom: 0,
        type: 'scroll'
      },
      series: [{
        name,
        type: 'pie',
        radius: ['45%', '68%'],
        center: ['50%', '44%'],
        avoidLabelOverlap: true,
        label: {
          formatter: '{b}: {c}'
        },
        data: data.map(item => ({
          name: item.name,
          value: item.count
        }))
      }]
    }
  }

  private resizeCharts = () => {
    if (this.companyChart) this.companyChart.resize()
    if (this.jobChart) this.jobChart.resize()
    if (this.sourceChart) this.sourceChart.resize()
  }

  private disposeCharts() {
    if (this.companyChart) {
      this.companyChart.dispose()
      this.companyChart = null
    }
    if (this.jobChart) {
      this.jobChart.dispose()
      this.jobChart = null
    }
    if (this.sourceChart) {
      this.sourceChart.dispose()
      this.sourceChart = null
    }
  }
}
</script>

<style lang="scss" scoped>
.distribution-charts {
  margin-top: 20px;
  margin-bottom: 20px;
}

.distribution-card {
  margin-bottom: 20px;
}

.card-header {
  font-weight: 600;
  color: #303133;
}

.chart {
  width: 100%;
  height: 300px;
}
</style>
