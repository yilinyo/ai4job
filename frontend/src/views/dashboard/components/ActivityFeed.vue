<template>
  <el-card shadow="never" class="activity-feed">
    <div slot="header" class="card-header">
      <span>最近动态</span>
      <div class="activity-feed__pager">
        <button type="button" class="pager-button" :disabled="currentPage === 1" @click="handlePageChange(currentPage - 1)">‹</button>
        <button
          v-for="page in pageCount"
          :key="page"
          type="button"
          class="pager-button"
          :class="{'is-active': page === currentPage}"
          @click="handlePageChange(page)"
        >{{ page }}</button>
        <button type="button" class="pager-button" :disabled="currentPage === pageCount" @click="handlePageChange(currentPage + 1)">›</button>
      </div>
    </div>
    <empty-state
      v-if="!activities.length"
      description="暂无动态"
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
  @Prop({ required: true, default: 0 }) private total!: number
  @Prop({ required: true, default: 1 }) private currentPage!: number

  private readonly pageSize = 5

  get pageCount() {
    return Math.max(1, Math.ceil(this.total / this.pageSize))
  }

  private formatDateTime(value?: string) {
    return formatDashboardDateTime(value)
  }

  private handleOpen(activity: ActivityLog) {
    if (!activity.route) return
    this.$emit('open', activity.route)
  }

  private handlePageChange(page: number) {
    this.$emit('page-change', page)
  }
}
</script>

<style lang="scss" scoped>
.activity-feed {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #303133;
}

.activity-feed__pager {
  display: flex;
  align-items: center;
}

.pager-button {
  min-width: 28px;
  height: 24px;
  padding: 0 6px;
  border: 1px solid #dcdfe6;
  border-radius: 3px;
  background: #fff;
  color: #606266;
  cursor: pointer;

  & + .pager-button {
    margin-left: 4px;
  }

  &.is-active {
    border-color: #409eff;
    background: #409eff;
    color: #fff;
  }

  &:disabled {
    color: #c0c4cc;
    cursor: not-allowed;
  }
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
