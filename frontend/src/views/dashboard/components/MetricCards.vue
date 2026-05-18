<template>
  <el-row :gutter="20" class="metric-cards">
    <el-col
      v-for="card in cards"
      :key="card.key"
      :xs="24"
      :sm="12"
      :lg="6"
    >
      <el-card shadow="hover" class="metric-card">
        <div class="metric-card__content">
          <div :class="['metric-card__icon', card.iconClass]">
            <i :class="card.icon" />
          </div>
          <div class="metric-card__body">
            <div class="metric-card__label">{{ card.label }}</div>
            <div class="metric-card__value">{{ card.value }}</div>
          </div>
        </div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { DashboardMetrics } from '@/models'

interface MetricCardItem {
  key: string
  label: string
  value: number
  icon: string
  iconClass: string
}

@Component({
  name: 'MetricCards'
})
export default class extends Vue {
  @Prop({
    required: true,
    default: () => ({
      totalApplications: 0,
      interviewing: 0,
      pendingFollowUp: 0,
      offers: 0
    })
  }) private metrics!: DashboardMetrics

  get cards(): MetricCardItem[] {
    return [
      {
        key: 'totalApplications',
        label: '总投递',
        value: this.metrics.totalApplications || 0,
        icon: 'el-icon-s-promotion',
        iconClass: 'is-blue'
      },
      {
        key: 'interviewing',
        label: '面试中',
        value: this.metrics.interviewing || 0,
        icon: 'el-icon-chat-dot-round',
        iconClass: 'is-purple'
      },
      {
        key: 'pendingFollowUp',
        label: '待跟进',
        value: this.metrics.pendingFollowUp || 0,
        icon: 'el-icon-alarm-clock',
        iconClass: 'is-orange'
      },
      {
        key: 'offers',
        label: 'Offer',
        value: this.metrics.offers || 0,
        icon: 'el-icon-medal',
        iconClass: 'is-green'
      }
    ]
  }
}
</script>

<style lang="scss" scoped>
.metric-cards {
  margin-bottom: 20px;
}

.metric-card {
  margin-bottom: 20px;

  ::v-deep .el-card__body {
    padding: 22px;
  }
}

.metric-card__content {
  display: flex;
  align-items: center;
}

.metric-card__icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 26px;

  &.is-blue {
    color: #409eff;
    background: #ecf5ff;
  }

  &.is-purple {
    color: #7c3aed;
    background: #f3e8ff;
  }

  &.is-orange {
    color: #e6a23c;
    background: #fdf6ec;
  }

  &.is-green {
    color: #67c23a;
    background: #f0f9eb;
  }
}

.metric-card__body {
  min-width: 0;
}

.metric-card__label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 8px;
}

.metric-card__value {
  color: #303133;
  font-size: 30px;
  line-height: 1;
  font-weight: 600;
}
</style>
