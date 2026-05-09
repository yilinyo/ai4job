<template>
  <div class="job-posting-list-container">
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="queryParams.keyword"
            placeholder="搜索公司名称或岗位名称"
            clearable
            @clear="handleSearch"
          >
            <i slot="prefix" class="el-input__icon el-icon-search" />
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="queryParams.sourcePlatform"
            placeholder="来源平台"
            clearable
            @change="handleSearch"
          >
            <el-option label="Boss" value="Boss" />
            <el-option label="官网" value="官网" />
            <el-option label="牛客" value="牛客" />
            <el-option label="内推" value="内推" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="queryParams.status"
            placeholder="岗位状态"
            clearable
            @change="handleSearch"
          >
            <el-option label="待投递" :value="JobPostingStatus.PENDING" />
            <el-option label="已投递" :value="JobPostingStatus.APPLIED" />
            <el-option label="不合适" :value="JobPostingStatus.NOT_FIT" />
            <el-option label="已关闭" :value="JobPostingStatus.CLOSED" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" icon="el-icon-search" @click="handleSearch">
            搜索
          </el-button>
        </el-col>
        <el-col :span="4" :offset="2">
          <el-button type="success" icon="el-icon-plus" @click="openCreateDialog">
            新增岗位
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <div class="job-posting-list">
      <el-row :gutter="20" v-loading="loading">
        <el-col
          v-for="jobPosting in jobPostings"
          :key="jobPosting.id"
          :span="8"
          style="margin-bottom: 20px"
        >
          <job-posting-card
            :job-posting="jobPosting"
            @view="handleView"
            @edit="handleEdit"
            @delete="handleDelete"
            @open-source="handleOpenSource"
          />
        </el-col>
      </el-row>

      <el-empty
        v-if="!loading && jobPostings.length === 0"
        description="暂无岗位记录，点击右上角新增一个吧"
      />

      <el-pagination
        v-if="total > 0"
        :current-page="queryParams.page"
        :page-size="queryParams.pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>

    <el-dialog
      :title="dialogTitle"
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
          <el-input v-model="jobPostingForm.companyName" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="岗位名称" prop="jobTitle">
          <el-input v-model="jobPostingForm.jobTitle" placeholder="请输入岗位名称" />
        </el-form-item>
        <el-form-item label="公司介绍">
          <el-input v-model="jobPostingForm.companyIntro" type="textarea" :rows="3" placeholder="请输入公司介绍" />
        </el-form-item>
        <el-form-item label="岗位要求">
          <el-input v-model="jobPostingForm.jobRequirements" type="textarea" :rows="4" placeholder="请输入岗位要求" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="base">
              <el-input v-model="jobPostingForm.base" placeholder="如：深圳 / 北京 / 上海" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="薪资范围">
              <el-input v-model="jobPostingForm.salaryRange" placeholder="如：25k-35k x 16" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="来源平台">
              <el-input v-model="jobPostingForm.sourcePlatform" placeholder="如：Boss / 官网 / 牛客" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="岗位状态" prop="status">
              <el-select v-model="jobPostingForm.status" placeholder="请选择岗位状态" style="width: 100%">
                <el-option label="待投递" :value="JobPostingStatus.PENDING" />
                <el-option label="已投递" :value="JobPostingStatus.APPLIED" />
                <el-option label="不合适" :value="JobPostingStatus.NOT_FIT" />
                <el-option label="已关闭" :value="JobPostingStatus.CLOSED" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="来源链接">
          <el-input v-model="jobPostingForm.sourceUrl" placeholder="请输入原岗位链接" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="jobPostingForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
      </el-form>

      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import JobPostingCard from '@/components/JobPostingCard/index.vue'
import { JobPostingModule } from '@/store/modules/job-posting'
import {
  JobPosting,
  JobPostingStatus,
  CreateJobPostingRequest
} from '@/models'

@Component({
  name: 'JobPostingList',
  components: {
    JobPostingCard
  }
})
export default class extends Vue {
  private JobPostingStatus = JobPostingStatus

  private queryParams = {
    keyword: '',
    sourcePlatform: undefined as string | undefined,
    status: undefined as JobPostingStatus | undefined,
    page: 1,
    pageSize: 9
  }

  private dialogVisible = false
  private dialogTitle = '新增岗位'
  private isEdit = false
  private editingJobPostingId = ''
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

  get jobPostings() {
    return JobPostingModule.jobPostings
  }

  get total() {
    return JobPostingModule.total
  }

  get loading() {
    return JobPostingModule.loading
  }

  created() {
    this.loadJobPostings()
  }

  private async loadJobPostings() {
    await JobPostingModule.GetJobPostings(this.queryParams)
  }

  private handleSearch() {
    this.queryParams.page = 1
    this.loadJobPostings()
  }

  private handlePageChange(page: number) {
    this.queryParams.page = page
    this.loadJobPostings()
  }

  private openCreateDialog() {
    this.dialogTitle = '新增岗位'
    this.isEdit = false
    this.dialogVisible = true
  }

  private handleView(jobPosting: JobPosting) {
    this.$router.push(`/job-postings/detail/${jobPosting.id}`)
  }

  private handleEdit(jobPosting: JobPosting) {
    this.dialogTitle = '编辑岗位'
    this.isEdit = true
    this.editingJobPostingId = jobPosting.id
    this.jobPostingForm = {
      companyName: jobPosting.companyName,
      companyIntro: jobPosting.companyIntro || '',
      jobTitle: jobPosting.jobTitle,
      jobRequirements: jobPosting.jobRequirements || '',
      base: jobPosting.base || '',
      salaryRange: jobPosting.salaryRange || '',
      sourcePlatform: jobPosting.sourcePlatform || '',
      sourceUrl: jobPosting.sourceUrl || '',
      status: jobPosting.status,
      remark: jobPosting.remark || ''
    }
    this.parseUrl = jobPosting.sourceUrl || ''
    this.dialogVisible = true
  }

  private handleOpenSource(jobPosting: JobPosting) {
    if (!jobPosting.sourceUrl) {
      this.$message.warning('当前岗位没有来源链接')
      return
    }
    window.open(jobPosting.sourceUrl, '_blank')
  }

  private async handleDelete(jobPosting: JobPosting) {
    try {
      await this.$confirm('确定要删除这个岗位吗？删除后无法恢复！', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await JobPostingModule.DeleteJobPosting(jobPosting.id)
      this.$message.success('删除成功！')
      this.loadJobPostings()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除岗位失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
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
        if (this.isEdit) {
          await JobPostingModule.UpdateJobPosting({
            id: this.editingJobPostingId,
            data: this.jobPostingForm
          })
          this.$message.success('更新成功！')
        } else {
          await JobPostingModule.CreateJobPosting(this.jobPostingForm)
          this.$message.success('创建成功！')
        }
        this.dialogVisible = false
        this.loadJobPostings()
      } catch (error) {
        console.error('保存岗位失败:', error)
        this.$message.error('操作失败，请稍后重试')
      } finally {
        this.submitLoading = false
      }
    })
  }

  private resetDialogState() {
    this.editingJobPostingId = ''
    this.isEdit = false
    this.parseUrl = ''
    this.jobPostingForm = {
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
    JobPostingModule.ClearParsedJobPosting()
  }
}
</script>

<style lang="scss" scoped>
.job-posting-list-container {
  .search-card {
    margin-bottom: 20px;
  }

  .job-posting-list {
    padding: 0 10px;
  }

  .parse-toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 18px;
  }
}
</style>
