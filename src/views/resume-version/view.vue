<template>
  <div class="resume-version-view-container">
    <el-page-header
      @back="goBack"
      content="查看简历版本"
    />

    <el-card v-loading="loading" class="version-card">
      <div v-if="currentVersion">
        <div class="version-header">
          <div class="header-left">
            <h2>{{ currentVersion.title }}</h2>
            <div class="meta-info">
              <el-tag v-if="currentVersion.isActive" type="success" size="small">
                当前激活版本
              </el-tag>
              <span class="version-number">{{ currentVersion.version }}</span>
              <span class="create-time">创建于 {{ formatDateTime(currentVersion.createdAt) }}</span>
            </div>
          </div>
          <div class="header-right">
            <el-button
              type="primary"
              size="small"
              icon="el-icon-edit"
              @click="handleEdit"
            >
              编辑
            </el-button>
            <el-button
              v-if="!currentVersion.isActive"
              type="success"
              size="small"
              icon="el-icon-check"
              @click="handleActivate"
            >
              激活此版本
            </el-button>
            <el-button
              size="small"
              icon="el-icon-download"
              :loading="exportingPdf"
              :disabled="exportingPdf"
              @click="handleExport"
            >
              {{ exportingPdf ? '导出中...' : '导出 PDF' }}
            </el-button>
          </div>
        </div>

        <div v-if="currentVersion.remark" class="remark-section">
          <div class="section-title">版本备注</div>
          <div class="remark-content">{{ currentVersion.remark }}</div>
        </div>

        <div class="content-section">
          <div class="section-title">简历内容</div>
          <resume-markdown-preview :html="renderedContent" />
        </div>

        <div class="applications-section">
          <div class="section-header">
            <div>
              <div class="section-title">关联投递记录</div>
              <div class="section-subtitle">当前版本共关联 {{ applicationTotal }} 条投递记录</div>
            </div>
            <el-button
              type="primary"
              size="small"
              icon="el-icon-plus"
              @click="handleCreateApplication"
            >
              新增投递记录
            </el-button>
          </div>

          <el-table
            v-loading="applicationLoading"
            :data="applications"
            style="width: 100%"
          >
            <el-table-column prop="companyName" label="公司名称" min-width="160" />
            <el-table-column prop="jobTitle" label="岗位名称" min-width="160" />
            <el-table-column prop="base" label="base" width="120">
              <template slot-scope="{row}">
                {{ row.base || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="appliedAt" label="投递时间" width="140">
              <template slot-scope="{row}">
                {{ formatDate(row.appliedAt) }}
              </template>
            </el-table-column>
            <el-table-column label="当前状态" width="120">
              <template slot-scope="{row}">
                <el-tag :type="applicationStatusTagType(row.currentStatus)" size="small">
                  {{ applicationStatusLabel(row.currentStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220">
              <template slot-scope="{row}">
                <el-button type="text" size="small" @click="handleViewApplication(row.id)">
                  查看
                </el-button>
                <el-button type="text" size="small" @click="handleEditApplication(row.id)">
                  编辑
                </el-button>
                <el-button type="text" size="small" class="danger" @click="handleDeleteApplication(row.id)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-if="!applicationLoading && applications.length === 0"
            description="当前版本还没有关联投递记录"
          />
        </div>
      </div>
    </el-card>

    <div class="pdf-export-stage">
      <div
        v-if="currentVersion"
        ref="pdfContent"
        class="pdf-export-document"
      >
        <header class="pdf-header">
          <div class="pdf-title-group">
            <h1>{{ currentVersion.title }}</h1>
            <p class="pdf-subtitle">简历版本 {{ currentVersion.version }}</p>
          </div>
          <div class="pdf-meta">
            <span>导出时间 {{ pdfExportTimestamp }}</span>
            <span>创建时间 {{ formatDateTime(currentVersion.createdAt) }}</span>
          </div>
        </header>

        <section v-if="currentVersion.remark" class="pdf-section">
          <h2>版本备注</h2>
          <p class="pdf-remark">{{ currentVersion.remark }}</p>
        </section>

        <section class="pdf-section">
          <h2>简历内容</h2>
          <resume-markdown-preview
            :html="renderedContent"
            variant="pdf"
            empty-text="暂无简历内容"
          />
        </section>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import html2pdf from 'html2pdf.js'
import ResumeMarkdownPreview from '@/components/ResumeMarkdownPreview/index.vue'
import { ResumeVersionModule } from '@/store/modules/resume-version'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import { ApplicationStatus } from '@/models'
import { renderResumeMarkdown } from '@/utils/resume-markdown'

@Component({
  name: 'ResumeVersionView',
  components: {
    ResumeMarkdownPreview
  }
})
export default class extends Vue {
  private versionId = ''
  private exportingPdf = false
  private pdfExportTimestamp = ''
  private applicationStatusMap: Record<string, string> = {
    [ApplicationStatus.APPLIED]: '已投递',
    [ApplicationStatus.VIEWED]: '已查看',
    [ApplicationStatus.WRITTEN_TEST]: '笔试',
    [ApplicationStatus.FIRST_INTERVIEW]: '一面',
    [ApplicationStatus.SECOND_INTERVIEW]: '二面',
    [ApplicationStatus.FINAL_INTERVIEW]: '终面',
    [ApplicationStatus.HR_INTERVIEW]: 'HR 面',
    [ApplicationStatus.OFFER]: 'Offer',
    [ApplicationStatus.REJECTED]: '未通过',
    [ApplicationStatus.CLOSED]: '已结束'
  }

  get loading() {
    return ResumeVersionModule.loading
  }

  get applicationLoading() {
    return ResumeApplicationModule.loading
  }

  get currentVersion() {
    return ResumeVersionModule.currentVersion
  }

  get applications() {
    return ResumeApplicationModule.applications
  }

  get applicationTotal() {
    return ResumeApplicationModule.total
  }

  get renderedContent() {
    if (!this.currentVersion) return ''
    return renderResumeMarkdown(this.currentVersion.content)
  }

  created() {
    this.versionId = this.$route.params.id
    this.loadVersion()
  }

  private async loadVersion() {
    await ResumeVersionModule.GetResumeVersionById(this.versionId)
    await ResumeApplicationModule.GetResumeApplications({
      resumeVersionId: this.versionId,
      page: 1,
      pageSize: 100
    })
  }

  private formatDateTime(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }

  private formatDate(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  private handleEdit() {
    this.$router.push(`/resume/version/${this.versionId}/edit`)
  }

  private handleCreateApplication() {
    if (!this.currentVersion) return
    this.$router.push(`/resume/application/create?versionId=${this.currentVersion.id}&repoId=${this.currentVersion.repoId}`)
  }

  private handleViewApplication(id: string) {
    this.$router.push(`/resume/application/${id}`)
  }

  private handleEditApplication(id: string) {
    if (!this.currentVersion) return
    this.$router.push(`/resume/application/${id}/edit?versionId=${this.currentVersion.id}&repoId=${this.currentVersion.repoId}`)
  }

  private async handleDeleteApplication(id: string) {
    try {
      await this.$confirm('确定要删除这条投递记录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await ResumeApplicationModule.DeleteResumeApplication({
        id,
        resumeVersionId: this.versionId
      })
      this.$message.success('删除成功！')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除投递记录失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
  }

  private applicationStatusLabel(status: ApplicationStatus) {
    return this.applicationStatusMap[status] || status
  }

  private applicationStatusTagType(status: ApplicationStatus) {
    if (status === ApplicationStatus.OFFER) return 'success'
    if (status === ApplicationStatus.REJECTED) return 'danger'
    return 'primary'
  }

  private async handleActivate() {
    if (!this.currentVersion) return

    try {
      await this.$confirm('确定要激活这个版本吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      })

      await ResumeVersionModule.ActivateVersion({
        id: this.versionId,
        repoId: this.currentVersion.repoId
      })
      this.$message.success('激活成功！')
      this.loadVersion()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('激活失败:', error)
        this.$message.error('激活失败，请稍后重试')
      }
    }
  }

  private async handleExport() {
    if (!this.currentVersion || this.exportingPdf) return

    const pdfContent = this.$refs.pdfContent as HTMLElement | undefined
    if (!pdfContent) {
      this.$message.error('导出失败，未找到可导出的简历内容')
      return
    }

    const fileName = this.getPdfFileName()

    this.exportingPdf = true
    this.pdfExportTimestamp = this.formatDateTime(new Date().toISOString())

    try {
      await this.$nextTick()

      await html2pdf()
        .set({
          margin: [12, 10, 12, 10],
          filename: fileName,
          image: { type: 'jpeg', quality: 0.98 },
          html2canvas: {
            scale: 2,
            useCORS: true,
            scrollY: 0
          },
          jsPDF: {
            unit: 'mm',
            format: 'a4',
            orientation: 'portrait'
          },
          pagebreak: {
            mode: ['css', 'legacy']
          }
        })
        .from(pdfContent)
        .save()

      this.$message.success(`PDF 导出成功：${fileName}`)
    } catch (error) {
      console.error('PDF 导出失败:', error)
      this.$message.error('PDF 导出失败，请稍后重试')
    } finally {
      this.exportingPdf = false
    }
  }

  private getPdfFileName() {
    if (!this.currentVersion) return 'resume.pdf'
    const rawFileName = `${this.currentVersion.title}_${this.currentVersion.version}`
    return `${rawFileName.replace(/[\\/:*?"<>|]/g, '-').trim() || 'resume'}.pdf`
  }

  private goBack() {
    if (this.currentVersion) {
      this.$router.push(`/resume/repo/${this.currentVersion.repoId}`)
    } else {
      this.$router.back()
    }
  }
}
</script>

<style lang="scss" scoped>
.resume-version-view-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .pdf-export-stage {
    position: fixed;
    left: -99999px;
    top: 0;
    width: 794px;
    pointer-events: none;
    z-index: -1;
  }

  .version-card {
    .version-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding-bottom: 20px;
      border-bottom: 1px solid #ebeef5;
      margin-bottom: 20px;

      .header-left {
        h2 {
          margin: 0 0 10px 0;
          font-size: 24px;
          color: #303133;
        }

        .meta-info {
          display: flex;
          align-items: center;
          gap: 15px;
          color: #909399;
          font-size: 14px;

          .version-number {
            font-weight: 500;
          }
        }
      }

      .header-right {
        display: flex;
        align-items: center;
        gap: 10px;
      }
    }

    .remark-section {
      margin-bottom: 30px;
      padding: 15px;
      background-color: #f5f7fa;
      border-radius: 4px;

      .section-title {
        font-size: 14px;
        font-weight: bold;
        color: #606266;
        margin-bottom: 8px;
      }

      .remark-content {
        color: #606266;
        line-height: 1.6;
      }
    }

    .content-section {
      margin-bottom: 30px;

      .section-title {
        font-size: 16px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 15px;
      }
    }

    .applications-section {
      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
      }

      .section-title {
        font-size: 16px;
        font-weight: bold;
        color: #303133;
      }

      .section-subtitle {
        margin-top: 4px;
        color: #909399;
        font-size: 13px;
      }

      .danger {
        color: #f56c6c;
      }
    }
  }

  .pdf-export-document {
    box-sizing: border-box;
    width: 794px;
    padding: 48px 56px;
    background-color: #fff;
    color: #1f2937;
    font-family: 'Helvetica Neue', Arial, sans-serif;
    line-height: 1.75;

    .pdf-header {
      display: flex;
      justify-content: space-between;
      gap: 24px;
      padding-bottom: 24px;
      border-bottom: 2px solid #dbe3f0;
      margin-bottom: 28px;
      page-break-inside: avoid;

      .pdf-title-group {
        h1 {
          margin: 0 0 8px;
          font-size: 30px;
          font-weight: 600;
          color: #111827;
        }

        .pdf-subtitle {
          margin: 0;
          font-size: 14px;
          color: #4b5563;
          letter-spacing: 0.08em;
          text-transform: uppercase;
        }
      }

      .pdf-meta {
        display: flex;
        flex-direction: column;
        gap: 8px;
        font-size: 12px;
        color: #6b7280;
        text-align: right;
      }
    }

    .pdf-section {
      margin-bottom: 28px;
      page-break-inside: avoid;

      h2 {
        margin: 0 0 14px;
        font-size: 16px;
        font-weight: 600;
        color: #1d4ed8;
      }
    }

    .pdf-remark {
      margin: 0;
      font-size: 14px;
      color: #1f2937;
      white-space: pre-wrap;
      word-break: break-word;
    }

  }
}
</style>
