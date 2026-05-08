<template>
  <el-card class="job-posting-card" :body-style="{padding: '20px'}">
    <div class="card-header">
      <h3 class="company-name">{{ jobPosting.companyName }}</h3>
      <el-tag :type="statusTagType">{{ statusText }}</el-tag>
    </div>

    <div class="card-content">
      <div class="info-item">
        <i class="el-icon-suitcase" />
        <span>{{ jobPosting.jobTitle }}</span>
      </div>
      <div class="info-item">
        <i class="el-icon-location-outline" />
        <span>{{ jobPosting.base || '未填写 base' }}</span>
      </div>
      <div class="info-item">
        <i class="el-icon-money" />
        <span>{{ jobPosting.salaryRange || '未填写薪资范围' }}</span>
      </div>
      <div class="info-item">
        <i class="el-icon-link" />
        <span>{{ jobPosting.sourcePlatform || '未知来源' }}</span>
      </div>
      <div v-if="jobPosting.jobRequirements" class="description">
        {{ jobPosting.jobRequirements }}
      </div>
    </div>

    <div class="card-footer">
      <div class="time-info">
        <span>创建于 {{ formatDate(jobPosting.createdAt) }}</span>
      </div>
      <div class="actions">
        <el-button type="text" size="small" @click="handleView">
          查看
        </el-button>
        <el-button type="text" size="small" @click="handleEdit">
          编辑
        </el-button>
        <el-button type="text" size="small" @click="handleOpenSource">
          原岗位
        </el-button>
        <el-button type="text" size="small" class="danger" @click="handleDelete">
          删除
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import { JobPosting, JobPostingStatus } from '@/models'

@Component({
  name: 'JobPostingCard'
})
export default class extends Vue {
  @Prop({ required: true }) private jobPosting!: JobPosting

  get statusText() {
    const statusMap: Record<string, string> = {
      [JobPostingStatus.PENDING]: '待投递',
      [JobPostingStatus.APPLIED]: '已投递',
      [JobPostingStatus.NOT_FIT]: '不合适',
      [JobPostingStatus.CLOSED]: '已关闭'
    }
    return statusMap[this.jobPosting.status] || this.jobPosting.status
  }

  get statusTagType() {
    if (this.jobPosting.status === JobPostingStatus.PENDING) return 'success'
    if (this.jobPosting.status === JobPostingStatus.APPLIED) return 'primary'
    if (this.jobPosting.status === JobPostingStatus.CLOSED) return 'danger'
    return 'info'
  }

  private formatDate(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  private handleView() {
    this.$emit('view', this.jobPosting)
  }

  private handleEdit() {
    this.$emit('edit', this.jobPosting)
  }

  private handleOpenSource() {
    this.$emit('open-source', this.jobPosting)
  }

  private handleDelete() {
    this.$emit('delete', this.jobPosting)
  }
}
</script>

<style lang="scss" scoped>
.job-posting-card {
  height: 100%;
  transition: all 0.3s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;

    .company-name {
      margin: 0;
      font-size: 18px;
      font-weight: bold;
      color: #303133;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      flex: 1;
      margin-right: 10px;
    }
  }

  .card-content {
    min-height: 140px;
    margin-bottom: 15px;
    display: flex;
    flex-direction: column;

    .info-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      color: #606266;
      font-size: 14px;

      i {
        margin-right: 8px;
        color: #909399;
      }
    }

    .description {
      margin-top: 12px;
      min-height: 40px;
      color: #909399;
      font-size: 13px;
      line-height: 1.5;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }

  .card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid #ebeef5;
    margin-top: auto;

    .time-info {
      font-size: 12px;
      color: #c0c4cc;
    }

    .actions {
      .el-button {
        padding: 0;
        margin-left: 10px;

        &.danger {
          color: #f56c6c;
        }
      }
    }
  }
}
</style>
