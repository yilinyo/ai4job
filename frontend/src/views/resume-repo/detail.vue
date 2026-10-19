<template>
  <div class="resume-repo-detail-container">
    <el-page-header
      @back="goBack"
      content="简历仓库详情"
    />

    <el-card v-loading="loading" class="detail-card">
      <div slot="header" class="card-header">
        <span class="title">仓库概览</span>
        <el-button
          type="primary"
          size="small"
          icon="el-icon-edit"
          @click="goToVersionList"
        >
          管理版本
        </el-button>
      </div>

      <div v-if="currentRepo" class="repo-info">
        <div class="repo-name-block">
          <span class="field-label">仓库名称</span>
          <h2>{{ currentRepo.name }}</h2>
        </div>

        <div class="info-grid">
          <div class="info-item">
            <span class="field-label">求职类型</span>
            <el-tag :type="jobTypeTag">{{ jobTypeText }}</el-tag>
          </div>
          <div class="info-item">
            <span class="field-label">目标岗位</span>
            <span class="field-value">{{ currentRepo.targetPosition }}</span>
          </div>
          <div class="info-item">
            <span class="field-label">版本数量</span>
            <span class="field-value">{{ currentRepo.versionCount }} 个</span>
          </div>
          <div class="info-item">
            <span class="field-label">创建时间</span>
            <span class="field-value">{{ formatDateTime(currentRepo.createdAt) }}</span>
          </div>
          <div class="info-item">
            <span class="field-label">更新时间</span>
            <span class="field-value">{{ formatDateTime(currentRepo.updatedAt) }}</span>
          </div>
        </div>

        <div class="description-block">
          <span class="field-label">仓库描述</span>
          <p>{{ currentRepo.description || '暂无描述' }}</p>
        </div>
      </div>
    </el-card>

    <el-card class="version-list-card">
      <div slot="header">
        <span>简历版本列表</span>
        <el-button
          type="success"
          size="small"
          icon="el-icon-plus"
          style="float: right"
          @click="handleCreateVersion"
        >
          创建新版本
        </el-button>
      </div>

      <el-table
        :data="versions"
        v-loading="versionLoading"
        style="width: 100%"
      >
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="title" label="版本标题" />
        <el-table-column label="状态" width="100">
          <template slot-scope="{row}">
            <el-tag v-if="row.isActive" type="success" size="small">
              当前版本
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template slot-scope="{row}">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template slot-scope="{row}">
            <el-button type="text" size="small" @click="handleViewVersion(row)">
              查看
            </el-button>
            <el-button type="text" size="small" @click="handleEditVersion(row)">
              编辑
            </el-button>
            <el-button
              v-if="!row.isActive"
              type="text"
              size="small"
              @click="handleActivateVersion(row)"
            >
              激活
            </el-button>
            <el-button
              type="text"
              size="small"
              class="danger"
              @click="handleDeleteVersion(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <empty-state
        v-if="!versionLoading && versions.length === 0"
        description="暂无简历版本"
      />
    </el-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { ResumeRepoModule } from '@/store/modules/resume-repo'
import { ResumeVersionModule } from '@/store/modules/resume-version'
import { ResumeVersion } from '@/models'
import { getJobTypeTag, getJobTypeText } from '@/utils/resume-repo'

@Component({
  name: 'ResumeRepoDetail'
})
export default class extends Vue {
  private repoId = ''

  get currentRepo() {
    return ResumeRepoModule.currentRepo
  }

  get loading() {
    return ResumeRepoModule.loading
  }

  get versions() {
    return ResumeVersionModule.versions
  }

  get versionLoading() {
    return ResumeVersionModule.loading
  }

  get jobTypeText() {
    return getJobTypeText(this.currentRepo?.jobType)
  }

  get jobTypeTag() {
    return getJobTypeTag(this.currentRepo?.jobType)
  }

  created() {
    this.repoId = this.$route.params.id
    this.loadRepoDetail()
    this.loadVersions()
  }

  activated() {
    // 详情页会被 keep-alive 缓存。从仓库列表编辑后再次进入时，重新同步信息。
    this.repoId = this.$route.params.id
    this.loadRepoDetail()
    this.loadVersions()
  }

  private async loadRepoDetail() {
    await ResumeRepoModule.GetResumeRepoById(this.repoId)
  }

  private async loadVersions() {
    await ResumeVersionModule.GetResumeVersions({
      repoId: this.repoId,
      page: 1,
      pageSize: 100
    })
  }

  private formatDateTime(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }

  private goBack() {
    this.$router.push('/resume/repos')
  }

  private goToVersionList() {
    this.$router.push(`/resume/version?repoId=${this.repoId}`)
  }

  private handleCreateVersion() {
    this.$router.push(`/resume/version/create?repoId=${this.repoId}`)
  }

  private handleViewVersion(version: ResumeVersion) {
    this.$router.push(`/resume/version/${version.id}`)
  }

  private handleEditVersion(version: ResumeVersion) {
    this.$router.push(`/resume/version/${version.id}/edit`)
  }

  private async handleActivateVersion(version: ResumeVersion) {
    try {
      await this.$confirm('确定要激活这个版本吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      })

      await ResumeVersionModule.ActivateVersion({
        id: version.id,
        repoId: this.repoId
      })
      this.$message.success('激活成功！')
      this.loadVersions()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('激活失败:', error)
        this.$message.error('激活失败，请稍后重试')
      }
    }
  }

  private async handleDeleteVersion(version: ResumeVersion) {
    try {
      await this.$confirm('确定要删除这个版本吗？删除后无法恢复！', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await ResumeVersionModule.DeleteResumeVersion({
        id: version.id,
        repoId: this.repoId
      })
      this.$message.success('删除成功！')
      this.loadVersions()
      this.loadRepoDetail()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.resume-repo-detail-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .detail-card {
    margin-bottom: 20px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 18px;
        font-weight: bold;
      }
    }

    .repo-info {
      .field-label {
        display: block;
        margin-bottom: 8px;
        color: #909399;
        font-size: 13px;
        line-height: 1;
      }

      .repo-name-block {
        padding-bottom: 20px;
        border-bottom: 1px solid #ebeef5;

        h2 {
          margin: 0;
          color: #303133;
          font-size: 24px;
          line-height: 1.35;
        }
      }

      .info-grid {
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 20px 32px;
        padding: 22px 0;
      }

      .info-item {
        min-width: 0;
      }

      .field-value {
        display: block;
        overflow: hidden;
        color: #303133;
        font-size: 15px;
        line-height: 1.5;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .description-block {
        padding: 18px 20px;
        background: #f8fafc;
        border-radius: 4px;

        p {
          margin: 0;
          color: #606266;
          font-size: 14px;
          line-height: 1.7;
          white-space: pre-wrap;
        }
      }
    }
  }

  @media (max-width: 900px) {
    .detail-card .repo-info .info-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }

  @media (max-width: 640px) {
    .detail-card .repo-info .info-grid {
      grid-template-columns: 1fr;
      gap: 18px;
    }
  }

  .version-list-card {
    .danger {
      color: #f56c6c;
    }
  }
}
</style>
