<template>
  <div class="resume-application-editor-container">
    <el-page-header
      @back="goBack"
      :content="pageTitle"
    />

    <el-card v-loading="loading" class="editor-card">
      <el-form
        ref="applicationForm"
        :model="applicationForm"
        :rules="applicationRules"
        label-width="120px"
      >
        <el-divider content-position="left">基础岗位信息</el-divider>

        <el-form-item label="公司名称" prop="companyName">
          <el-input v-model="applicationForm.companyName" placeholder="请输入公司名称" />
        </el-form-item>

        <el-form-item label="关联岗位">
          <el-select
            v-model="selectedJobPostingId"
            filterable
            clearable
            placeholder="可选择岗位库中的已有岗位"
            style="width: 100%"
            @change="handleJobPostingChange"
          >
            <el-option
              v-for="item in jobPostings"
              :key="item.id"
              :label="`${item.companyName} / ${item.jobTitle}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="岗位名称" prop="jobTitle">
          <el-input v-model="applicationForm.jobTitle" placeholder="请输入岗位名称" />
        </el-form-item>

        <el-form-item label="投递日期" prop="appliedAt">
          <el-date-picker
            v-model="applicationForm.appliedAt"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="选择投递日期"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="当前状态" prop="currentStatus">
          <el-select v-model="applicationForm.currentStatus" placeholder="请选择当前状态" style="width: 100%">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="base">
          <el-input v-model="applicationForm.base" placeholder="如：深圳 / 北京 / 上海" />
        </el-form-item>

        <el-form-item label="薪资范围">
          <el-input v-model="applicationForm.salaryRange" placeholder="如：25k-35k x 16" />
        </el-form-item>

        <el-form-item label="投递渠道">
          <el-input v-model="applicationForm.deliveryChannel" placeholder="如：官网 / 内推 / Boss 直聘" />
        </el-form-item>

        <el-divider content-position="left">公司与 JD 信息</el-divider>

        <el-form-item label="公司介绍">
          <el-input
            v-model="applicationForm.companyIntro"
            type="textarea"
            :rows="3"
            placeholder="记录公司介绍、业务方向、团队情况等"
          />
        </el-form-item>

        <el-form-item label="岗位要求">
          <el-input
            v-model="applicationForm.jobRequirements"
            type="textarea"
            :rows="4"
            placeholder="记录 JD 摘要、核心要求和加分项"
          />
        </el-form-item>

        <el-divider content-position="left">面试总结与笔记</el-divider>

        <el-form-item label="面试总结">
          <el-input
            v-model="applicationForm.interviewSummary"
            type="textarea"
            :rows="3"
            placeholder="总结这次投递和面试推进情况"
          />
        </el-form-item>

        <el-form-item label="面试小记">
          <el-input
            v-model="applicationForm.interviewNotes"
            type="textarea"
            :rows="4"
            placeholder="记录问题、亮点、失误点和后续补强方向"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="applicationForm.remark"
            type="textarea"
            :rows="2"
            placeholder="其他补充信息"
          />
        </el-form-item>

        <el-divider content-position="left">个人评分</el-divider>

        <div class="score-grid">
          <el-form-item label="简历匹配度">
            <el-input-number v-model="applicationForm.resumeMatchScore" :min="1" :max="5" />
          </el-form-item>
          <el-form-item label="面试表现">
            <el-input-number v-model="applicationForm.interviewPerformanceScore" :min="1" :max="5" />
          </el-form-item>
          <el-form-item label="岗位意向度">
            <el-input-number v-model="applicationForm.roleInterestScore" :min="1" :max="5" />
          </el-form-item>
          <el-form-item label="综合评分">
            <el-input-number v-model="applicationForm.overallScore" :min="1" :max="10" />
          </el-form-item>
        </div>

        <el-form-item>
          <el-button
            type="primary"
            :loading="submitLoading"
            @click="handleSubmit"
          >
            {{ isEdit ? '保存修改' : '创建投递记录' }}
          </el-button>
          <el-button @click="goBack">
            取消
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import { ResumeApplicationModule } from '@/store/modules/resume-application'
import { JobPostingModule } from '@/store/modules/job-posting'
import {
  ApplicationStatus,
  CreateResumeApplicationRequest,
  UpdateResumeApplicationRequest
} from '@/models'

@Component({
  name: 'ResumeApplicationEditor'
})
export default class extends Vue {
  private applicationId = ''
  private versionId = ''
  private repoId = ''
  private selectedJobPostingId = ''
  private isEdit = false
  private submitLoading = false

  private statusOptions = [
    { label: '已投递', value: ApplicationStatus.APPLIED },
    { label: '已查看', value: ApplicationStatus.VIEWED },
    { label: '笔试', value: ApplicationStatus.WRITTEN_TEST },
    { label: '一面', value: ApplicationStatus.FIRST_INTERVIEW },
    { label: '二面', value: ApplicationStatus.SECOND_INTERVIEW },
    { label: '终面', value: ApplicationStatus.FINAL_INTERVIEW },
    { label: 'HR 面', value: ApplicationStatus.HR_INTERVIEW },
    { label: 'Offer', value: ApplicationStatus.OFFER },
    { label: '未通过', value: ApplicationStatus.REJECTED },
    { label: '已结束', value: ApplicationStatus.CLOSED }
  ]

  private applicationForm = {
    companyName: '',
    jobTitle: '',
    companyIntro: '',
    jobRequirements: '',
    base: '',
    salaryRange: '',
    sourcePlatform: '',
    sourceUrl: '',
    deliveryChannel: '',
    appliedAt: this.formatDateOnly(new Date().toISOString()),
    currentStatus: ApplicationStatus.APPLIED,
    interviewSummary: '',
    interviewNotes: '',
    resumeMatchScore: 3,
    interviewPerformanceScore: 3,
    roleInterestScore: 3,
    overallScore: 6,
    remark: ''
  }

  private applicationRules = {
    companyName: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
    jobTitle: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
    appliedAt: [{ required: true, message: '请选择投递日期', trigger: 'change' }],
    currentStatus: [{ required: true, message: '请选择当前状态', trigger: 'change' }]
  }

  get loading() {
    return ResumeApplicationModule.loading
  }

  get jobPostings() {
    return JobPostingModule.jobPostings
  }

  get pageTitle() {
    return this.isEdit ? '编辑投递记录' : '创建投递记录'
  }

  async created() {
    this.applicationId = this.$route.params.id || ''
    this.versionId = this.$route.query.versionId as string || ''
    this.repoId = this.$route.query.repoId as string || ''
    this.selectedJobPostingId = this.$route.query.jobPostingId as string || ''
    this.isEdit = !!this.applicationId && this.$route.path.includes('/edit')

    await this.loadJobPostings()

    if (this.isEdit) {
      await this.loadApplication()
    } else if (this.selectedJobPostingId) {
      await this.fillFromJobPosting(this.selectedJobPostingId)
    }
  }

  private async loadJobPostings() {
    await JobPostingModule.GetJobPostings({
      page: 1,
      pageSize: 100
    })
  }

  private async loadApplication() {
    await ResumeApplicationModule.GetResumeApplicationById(this.applicationId)
    const application = ResumeApplicationModule.currentApplication
    if (application) {
      this.selectedJobPostingId = application.jobPostingId || ''
      this.versionId = application.resumeVersionId
      this.repoId = application.repoId
      this.applicationForm = {
        companyName: application.companyName,
        jobTitle: application.jobTitle,
        companyIntro: application.companyIntro || '',
        jobRequirements: application.jobRequirements || '',
        base: application.base || '',
        salaryRange: application.salaryRange || '',
        sourcePlatform: application.sourcePlatform || '',
        sourceUrl: application.sourceUrl || '',
        deliveryChannel: application.deliveryChannel || '',
        appliedAt: this.formatDateOnly(application.appliedAt),
        currentStatus: application.currentStatus,
        interviewSummary: application.interviewSummary || '',
        interviewNotes: application.interviewNotes || '',
        resumeMatchScore: application.resumeMatchScore || 3,
        interviewPerformanceScore: application.interviewPerformanceScore || 3,
        roleInterestScore: application.roleInterestScore || 3,
        overallScore: application.overallScore || 6,
        remark: application.remark || ''
      }
    }
  }

  private async fillFromJobPosting(jobPostingId: string) {
    if (!jobPostingId) return

    await JobPostingModule.GetJobPostingById(jobPostingId)
    const jobPosting = JobPostingModule.currentJobPosting
    if (!jobPosting) return

    this.applicationForm = {
      ...this.applicationForm,
      companyName: jobPosting.companyName,
      companyIntro: jobPosting.companyIntro || '',
      jobTitle: jobPosting.jobTitle,
      jobRequirements: jobPosting.jobRequirements || '',
      base: jobPosting.base || '',
      salaryRange: jobPosting.salaryRange || '',
      sourcePlatform: jobPosting.sourcePlatform || '',
      sourceUrl: jobPosting.sourceUrl || ''
    }
  }

  private async handleJobPostingChange(jobPostingId: string) {
    this.selectedJobPostingId = jobPostingId || ''
    if (this.selectedJobPostingId) {
      await this.fillFromJobPosting(this.selectedJobPostingId)
    }
  }

  private formatDateOnly(dateString: string) {
    const date = new Date(dateString)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  private handleSubmit() {
    (this.$refs.applicationForm as ElForm).validate(async(valid: boolean) => {
      if (!valid) return

      this.submitLoading = true
      try {
        if (this.isEdit) {
          const payload: UpdateResumeApplicationRequest = {
            jobPostingId: this.selectedJobPostingId || undefined,
            ...this.applicationForm
          }
          await ResumeApplicationModule.UpdateResumeApplication({
            id: this.applicationId,
            data: payload
          })
          this.$message.success('更新成功！')
          this.$router.push(`/resume/application/${this.applicationId}`)
        } else {
          const payload: CreateResumeApplicationRequest = {
            jobPostingId: this.selectedJobPostingId || undefined,
            resumeVersionId: this.versionId,
            repoId: this.repoId,
            ...this.applicationForm
          }
          const created = await ResumeApplicationModule.CreateResumeApplication(payload)
          this.$message.success('创建成功！')
          if (created && created.id) {
            this.$router.push(`/resume/application/${created.id}`)
          } else {
            this.goBack()
          }
        }
      } catch (error) {
        console.error('提交失败:', error)
        this.$message.error('操作失败，请稍后重试')
      } finally {
        this.submitLoading = false
      }
    })
  }

  private goBack() {
    if (this.isEdit && this.applicationId) {
      this.$router.push(`/resume/application/${this.applicationId}`)
      return
    }

    if (this.versionId) {
      this.$router.push(`/resume/version/${this.versionId}`)
      return
    }

    this.$router.back()
  }
}
</script>

<style lang="scss" scoped>
.resume-application-editor-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .score-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 8px 16px;
  }
}
</style>
