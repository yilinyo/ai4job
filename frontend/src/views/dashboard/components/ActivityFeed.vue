<template>
  <el-card shadow="never" class="activity-feed">
    <div slot="header" class="card-header">最近动态</div>
    <el-empty
      v-if="!activities.length"
      description="暂无动态"
      :image-size="100"
    />
    <el-timeline v-else class="activity-feed__timeline">
      <el-timeline-item
        v-for="activity in activities"
        :key="activity.id"
        :timestamp="formatDateTime(activity.createdAt)"
        placement="top"
      >
        <div
          :class="['activity-feed__item', {'is-clickable': activity.route}]"
          @click="handleOpen(activity)"
        >
          <div class="activity-feed__summary">{{ activity.summary }}</div>
          <div class="activity-feed__meta">{{ activity.targetLabel }}</div>
        </div>
      </el-timeline-item>
    </el-timeline>
  </el-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { ActivityLog } from '@/models'
import { parseTime } from '@/utils'

const formatDashboardDateTime = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) {
    const pad = (num: number) => String(num).padStart(2, '0')
    return `${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
  }
  return parseTime(value, '{m}-{d} {h}:{i}') || '-'
}

@Component({
  name: 'ActivityFeed'
})
export default class extends Vue {
  @Prop({ required: true, default: () => [] }) private activities!: ActivityLog[]

  private formatDateTime(value?: string) {
    return formatDashboardDateTime(value)
  }

  private handleOpen(activity: ActivityLog) {
    if (!activity.route) return
    this.$emit('open', activity.route)
  }
}
</script>

<style lang="scss" scoped>
.activity-feed {
  margin-bottom: 20px;
}

.card-header {
  font-weight: 600;
  color: #303133;
}

.activity-feed__timeline {
  padding-left: 0;
}

.activity-feed__item {
  padding: 10px 12px;
  border-radius: 8px;
  background: #f8fafc;

  &.is-clickable {
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      color: #409eff;
      background: #ecf5ff;
    }
  }
}

.activity-feed__summary {
  color: #303133;
  line-height: 1.5;
}

.activity-feed__meta {
  color: #909399;
  font-size: 12px;
  margin-top: 6px;
}
</style>
