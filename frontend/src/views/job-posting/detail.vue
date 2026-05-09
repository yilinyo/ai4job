<template>
  <div class="job-posting-detail-container">
    <el-page-header
      @back="goBack"
      content="岗位详情"
    />

    <el-card v-loading="loading" class="detail-card">
      <div v-if="currentJobPosting">
        <div class="header-block">
          <div>
            <h2>{{ currentJobPosting.companyName }} / {{ currentJobPosting.jobTitle }}</h2>
            <div class="meta-info">
              <el-tag :type="statusTagType">{{ statusText }}</el-tag>
              <span>{{ currentJobPosting.base || '未填写 base' }}</span>
              <span>{{ currentJobPosting.salaryRange || '未填写薪资范围' }}</span>
              <span>{{ currentJobPosting.sourcePlatform || '未知来源' }}</span>
            </div>
          </div>
        </div>

        <el-card shadow="never" class="info-operation-card">
          <div slot="header">岗位基础信息与操作</div>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="公司介绍">
              {{ currentJobPosting.companyIntro || '暂无' }}
            </el-descriptions-item>
            <el-descriptions-item label="岗位要求">
              {{ currentJobPosting.jobRequirements || '暂无' }}
            </el-descriptions-item>
            <el-descriptions-item label="来源链接">
              <el-link
                v-if="currentJobPosting.sourceUrl"
                type="primary"
                :href="currentJobPosting.sourceUrl"
                target="_blank"
              >
                {{ currentJobPosting.sourceUrl }}
              </el-link>
              <span v-else>暂无</span>
            </el-descriptions-item>
            <el-descriptions-item label="备注">
              {{ currentJobPosting.remark || '暂无' }}
            </el-descriptions-item>
          </el-descriptions>

          <div class="action-row">
            <el-button type="success" icon="el-icon-plus" @click="handleCreateApplication">
              创建投递记录
            </el-button>
            <el-button type="primary" icon="el-icon-edit" @click="openEditDialog">
              编辑岗位
            </el-button>
            <el-button plain icon="el-icon-link" @click="handleOpenSource">
              查看原岗位
            </el-button>
          </div>
        </el-card>

        <el-card shadow="never" class="applications-card">
          <div slot="header">
            关联投递记录（{{ applications.length }}）
          </div>

          <el-empty
            v-if="!applicationLoading && applications.length === 0"
            description="当前岗位还没有关联投递记录"
          />

          <el-table v-else v-loading="applicationLoading" :data="applications" style="width: 100%">
            <el-table-column prop="companyName" label="公司名称" min-width="150" />
            <el-table-column prop="jobTitle" label="岗位名称" min-width="150" />
            <el-table-column label="简历版本" min-width="140">
              <template slot-scope="{row}">
                {{ getVersionDisplayName(row.resumeVersionId) }}
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
                  {{ applicationStatusText(row.currentStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template slot-scope="{row}">
                <el-button type="text" size="small" @click="handleViewApplication(row.id)">
                  查看投递
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </el-card>

    <el-dialog
      title="编辑岗位"
      :visible.sync="dialogVisible"
      width="720px"
      @closed="resetDialogState"
    >
      <div class="parse-toolbar">
        <el-input
          v-model="parseUrl"
          placeholder="输入岗位来源链接，点击解析导入"
          clearable
        />
        <el-button type="primary" plain :loading="parsing" @click="handleParse">
          解析导入
        </el-button>
      </div>

      <el-form
        ref="jobPostingForm"
        :model="jobPostingForm"
        :rules="jobPostingRules"
        label-width="100px"
      >
        <el-form-item label="公司名称" prop="companyName">
          <el-input v-model="jobPostingForm.companyName" />
        </el-form-item>
        <el-form-item label="岗位名称" prop="jobTitle">
          <el-input v-model="jobPostingForm.jobTitle" />
        </el-form-item>
        <el-form-item label="公司介绍">
          <el-input v-model="jobPostingForm.companyIntro" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="岗位要求">
          <el-input v-model="jobPostingForm.jobRequirements" type="textarea" :rows="4" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="base">
              <el-input v-model="jobPostingForm.base" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="薪资范围">
              <el-input v-model="jobPostingForm.salaryRange" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="来源平台">
              <el-input v-model="jobPostingForm.sourcePlatform" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="岗位状态" prop="status">
              <el-select v-model="jobPostingForm.status" style="width: 100%">
                <el-option label="待投递" :value="JobPostingStatus.PENDING" />
                <el-option label="已投递" :value="JobPostingStatus.APPLIED" />
                <el-option label="不合适" :value="JobPostingStatus.NOT_FIT" />
                <el-option label="已关闭" :value="JobPostingStatus.CLOSED" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="来源链接">
          <el-input v-model="jobPostingForm.sourceUrl" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="jobPostingForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>

      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          保存
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import { JobPostingModule } from '@/store/modules/job-posting'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import {
  ApplicationStatus,
  CreateJobPostingRequest,
  JobPostingStatus
} from '@/models'

@Component({
  name: 'JobPostingDetail'
})
export default class extends Vue {
  private JobPostingStatus = JobPostingStatus
  private jobPostingId = ''
  private dialogVisible = false
  private submitLoading = false
  private parsing = false
  private parseUrl = ''

  private jobPostingForm: CreateJobPostingRequest = {
    companyName: '',
    companyIntro: '',
    jobTitle: '',
    jobRequirements: '',
    base: '',
    salaryRange: '',
    sourcePlatform: '',
    sourceUrl: '',
    status: JobPostingStatus.PENDING,
    remark: ''
  }

  private jobPostingRules = {
    companyName: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
    jobTitle: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
    status: [{ required: true, message: '请选择岗位状态', trigger: 'change' }]
  }

  get loading() {
    return JobPostingModule.loading
  }

  get applicationLoading() {
    return ResumeApplicationModule.loading
  }

  get currentJobPosting() {
    return JobPostingModule.currentJobPosting
  }

  get applications() {
    return ResumeApplicationModule.applications.filter(item => item.jobPostingId === this.jobPostingId)
  }

  get statusText() {
    const statusMap: Record<string, string> = {
      [JobPostingStatus.PENDING]: '待投递',
      [JobPostingStatus.APPLIED]: '已投递',
      [JobPostingStatus.NOT_FIT]: '不合适',
      [JobPostingStatus.CLOSED]: '已关闭'
    }
    return this.currentJobPosting ? statusMap[this.currentJobPosting.status] || this.currentJobPosting.status : ''
  }

  get statusTagType() {
    if (!this.currentJobPosting) return 'info'
    if (this.currentJobPosting.status === JobPostingStatus.PENDING) return 'success'
    if (this.currentJobPosting.status === JobPostingStatus.APPLIED) return 'primary'
    if (this.currentJobPosting.status === JobPostingStatus.CLOSED) return 'danger'
    return 'info'
  }

  created() {
    this.jobPostingId = this.$route.params.id
    this.loadData()
  }

  private async loadData() {
    await JobPostingModule.GetJobPostingById(this.jobPostingId)
    await ResumeApplicationModule.GetResumeApplications({
      page: 1,
      pageSize: 100
    })
  }

  private formatDate(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  private applicationStatusText(status: ApplicationStatus) {
    const statusMap: Record<string, string> = {
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
    return statusMap[status] || status
  }

  private applicationStatusTagType(status: ApplicationStatus) {
    if (status === ApplicationStatus.OFFER) return 'success'
    if (status === ApplicationStatus.REJECTED) return 'danger'
    return 'primary'
  }

  private getVersionDisplayName(resumeVersionId: string) {
    return `简历版本 #${resumeVersionId}`
  }

  private handleCreateApplication() {
    if (!this.currentJobPosting) return
    this.$router.push(`/resume/application/create?jobPostingId=${this.currentJobPosting.id}`)
  }

  private openEditDialog() {
    if (!this.currentJobPosting) return

    this.jobPostingForm = {
      companyName: this.currentJobPosting.companyName,
      companyIntro: this.currentJobPosting.companyIntro || '',
      jobTitle: this.currentJobPosting.jobTitle,
      jobRequirements: this.currentJobPosting.jobRequirements || '',
      base: this.currentJobPosting.base || '',
      salaryRange: this.currentJobPosting.salaryRange || '',
      sourcePlatform: this.currentJobPosting.sourcePlatform || '',
      sourceUrl: this.currentJobPosting.sourceUrl || '',
      status: this.currentJobPosting.status,
      remark: this.currentJobPosting.remark || ''
    }
    this.parseUrl = this.currentJobPosting.sourceUrl || ''
    this.dialogVisible = true
  }

  private handleOpenSource() {
    if (!this.currentJobPosting?.sourceUrl) {
      this.$message.warning('当前岗位没有来源链接')
      return
    }
    window.open(this.currentJobPosting.sourceUrl, '_blank')
  }

  private async handleParse() {
    if (!this.parseUrl) {
      this.$message.warning('请先输入岗位链接')
      return
    }

    this.parsing = true
    try {
      const parsed = await JobPostingModule.ParseJobPosting(this.parseUrl)
      if (parsed) {
        this.jobPostingForm = {
          ...this.jobPostingForm,
          companyName: parsed.companyName,
          companyIntro: parsed.companyIntro || '',
          jobTitle: parsed.jobTitle,
          jobRequirements: parsed.jobRequirements || '',
          base: parsed.base || '',
          salaryRange: parsed.salaryRange || '',
          sourcePlatform: parsed.sourcePlatform || '',
          sourceUrl: parsed.sourceUrl || this.parseUrl
        }
        this.$message.success('解析成功，已回填岗位信息')
      }
    } catch (error) {
      console.error('解析岗位失败:', error)
      this.$message.error('解析失败，请稍后重试')
    } finally {
      this.parsing = false
    }
  }

  private handleSubmit() {
    (this.$refs.jobPostingForm as ElForm).validate(async(valid: boolean) => {
      if (!valid) return

      this.submitLoading = true
      try {
        await JobPostingModule.UpdateJobPosting({
          id: this.jobPostingId,
          data: this.jobPostingForm
        })
        await this.loadData()
        this.dialogVisible = false
        this.$message.success('更新成功！')
      } catch (error) {
        console.error('更新岗位失败:', error)
        this.$message.error('更新失败，请稍后重试')
      } finally {
        this.submitLoading = false
      }
    })
  }

  private handleViewApplication(id: string) {
    this.$router.push(`/resume/application/${id}`)
  }

  private resetDialogState() {
    this.parseUrl = ''
    JobPostingModule.ClearParsedJobPosting()
  }

  private goBack() {
    this.$router.push('/job-postings/list')
  }
}
</script>

<style lang="scss" scoped>
.job-posting-detail-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .header-block {
    margin-bottom: 20px;

    h2 {
      margin: 0 0 10px;
      font-size: 24px;
      color: #303133;
    }

    .meta-info {
      display: flex;
      align-items: center;
      gap: 15px;
      color: #909399;
      font-size: 14px;
      flex-wrap: wrap;
    }
  }

  .info-operation-card,
  .applications-card {
    margin-bottom: 20px;
  }

  .action-row {
    display: flex;
    gap: 12px;
    margin-top: 18px;
  }

  .parse-toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 18px;
  }
}
</style>
