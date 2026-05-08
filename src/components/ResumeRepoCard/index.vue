<template>
  <el-card class="repo-card" :body-style="{padding: '20px'}">
    <div class="card-header">
      <h3 class="repo-name">{{ repo.name }}</h3>
      <el-tag :type="jobTypeTag">{{ jobTypeText }}</el-tag>
    </div>

    <div class="card-content">
      <div class="info-item">
        <i class="el-icon-suitcase" />
        <span>{{ repo.targetPosition }}</span>
      </div>
      <div class="info-item">
        <i class="el-icon-document" />
        <span>{{ repo.versionCount }} 个版本</span>
      </div>
      <div v-if="repo.description" class="description">
        {{ repo.description }}
      </div>
    </div>

    <div class="card-footer">
      <div class="time-info">
        <span>创建于 {{ formatDate(repo.createdAt) }}</span>
      </div>
      <div class="actions">
        <el-button type="text" size="small" @click="handleView">
          查看
        </el-button>
        <el-button type="text" size="small" @click="handleEdit">
          编辑
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
import { ResumeRepo } from '@/models'
import { getJobTypeTag, getJobTypeText } from '@/utils/resume-repo'

@Component({
  name: 'ResumeRepoCard'
})
export default class extends Vue {
  @Prop({ required: true }) private repo!: ResumeRepo

  get jobTypeText() {
    return getJobTypeText(this.repo.jobType)
  }

  get jobTypeTag() {
    return getJobTypeTag(this.repo.jobType)
  }

  private formatDate(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  private handleView() {
    this.$emit('view', this.repo)
  }

  private handleEdit() {
    this.$emit('edit', this.repo)
  }

  private handleDelete() {
    this.$emit('delete', this.repo)
  }
}
</script>

<style lang="scss" scoped>
.repo-card {
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

    .repo-name {
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
    margin-bottom: 15px;

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
