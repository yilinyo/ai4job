<template>
  <div class="resume-repo-list-container">
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="queryParams.keyword"
            placeholder="搜索仓库名称或岗位"
            clearable
            @clear="handleSearch"
          >
            <i
              slot="prefix"
              class="el-input__icon el-icon-search"
            />
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="queryParams.jobType"
            placeholder="求职类型"
            clearable
            @change="handleSearch"
          >
            <el-option label="校招" value="campus" />
            <el-option label="社招" value="social" />
            <el-option label="实习" value="internship" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button
            type="primary"
            icon="el-icon-search"
            @click="handleSearch"
          >
            搜索
          </el-button>
        </el-col>
        <el-col :span="4" :offset="6">
          <el-button
            type="success"
            icon="el-icon-plus"
            @click="handleCreate"
          >
            创建仓库
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <div class="repo-list">
      <el-row :gutter="20" v-loading="loading">
        <el-col
          v-for="repo in repos"
          :key="repo.id"
          :span="8"
          style="margin-bottom: 20px"
        >
          <resume-repo-card
            :repo="repo"
            @view="handleView"
            @edit="handleEdit"
            @delete="handleDelete"
          />
        </el-col>
      </el-row>

      <el-empty
        v-if="!loading && repos.length === 0"
        description="暂无简历仓库，点击右上角创建一个吧"
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      :visible.sync="dialogVisible"
      width="600px"
    >
      <el-form
        ref="repoForm"
        :model="repoForm"
        :rules="repoRules"
        label-width="100px"
      >
        <el-form-item label="仓库名称" prop="name">
          <el-input v-model="repoForm.name" placeholder="请输入仓库名称" />
        </el-form-item>

        <el-form-item label="求职类型" prop="jobType">
          <el-radio-group v-model="repoForm.jobType">
            <el-radio label="campus">校招</el-radio>
            <el-radio label="social">社招</el-radio>
            <el-radio label="internship">实习</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="目标岗位" prop="targetPosition">
          <el-input v-model="repoForm.targetPosition" placeholder="如：前端工程师" />
        </el-form-item>

        <el-form-item label="仓库描述" prop="description">
          <el-input
            v-model="repoForm.description"
            type="textarea"
            :rows="4"
            placeholder="简单描述一下这个仓库的用途"
          />
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
import { ResumeRepoModule } from '@/store/modules/resume-repo'
import { JobType, CreateResumeRepoRequest, ResumeRepo } from '@/models'
import ResumeRepoCard from '@/components/ResumeRepoCard/index.vue'

@Component({
  name: 'ResumeRepoList',
  components: {
    ResumeRepoCard
  }
})
export default class extends Vue {
  private queryParams = {
    keyword: '',
    jobType: undefined as JobType | undefined,
    page: 1,
    pageSize: 9
  }

  private dialogVisible = false
  private dialogTitle = '创建简历仓库'
  private isEdit = false
  private editingRepoId = ''
  private submitLoading = false

  private repoForm: CreateResumeRepoRequest = {
    name: '',
    jobType: JobType.CAMPUS,
    targetPosition: '',
    description: ''
  }

  private repoRules = {
    name: [{ required: true, message: '请输入仓库名称', trigger: 'blur' }],
    jobType: [{ required: true, message: '请选择求职类型', trigger: 'change' }],
    targetPosition: [{ required: true, message: '请输入目标岗位', trigger: 'blur' }]
  }

  get repos() {
    return ResumeRepoModule.repos
  }

  get total() {
    return ResumeRepoModule.total
  }

  get loading() {
    return ResumeRepoModule.loading
  }

  created() {
    this.loadRepos()
  }

  private async loadRepos() {
    await ResumeRepoModule.GetResumeRepos(this.queryParams)
  }

  private handleSearch() {
    this.queryParams.page = 1
    this.loadRepos()
  }

  private handlePageChange(page: number) {
    this.queryParams.page = page
    this.loadRepos()
  }

  private handleCreate() {
    this.dialogTitle = '创建简历仓库'
    this.isEdit = false
    this.repoForm = {
      name: '',
      jobType: JobType.CAMPUS,
      targetPosition: '',
      description: ''
    }
    this.dialogVisible = true
  }

  private handleView(repo: ResumeRepo) {
    this.$router.push(`/resume/repo/${repo.id}`)
  }

  private handleEdit(repo: ResumeRepo) {
    this.dialogTitle = '编辑简历仓库'
    this.isEdit = true
    this.editingRepoId = repo.id
    this.repoForm = {
      name: repo.name,
      jobType: repo.jobType,
      targetPosition: repo.targetPosition,
      description: repo.description
    }
    this.dialogVisible = true
  }

  private async handleDelete(repo: ResumeRepo) {
    try {
      await this.$confirm('确定要删除这个简历仓库吗？删除后无法恢复！', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await ResumeRepoModule.DeleteResumeRepo(repo.id)
      this.$message.success('删除成功！')
      this.loadRepos()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        this.$message.error('删除失败，请稍后重试')
      }
    }
  }

  private handleSubmit() {
    (this.$refs.repoForm as ElForm).validate(async(valid: boolean) => {
      if (valid) {
        this.submitLoading = true
        try {
          if (this.isEdit) {
            await ResumeRepoModule.UpdateResumeRepo({
              id: this.editingRepoId,
              data: this.repoForm
            })
            this.$message.success('更新成功！')
          } else {
            await ResumeRepoModule.CreateResumeRepo(this.repoForm)
            this.$message.success('创建成功！')
          }
          this.dialogVisible = false
          this.loadRepos()
        } catch (error) {
          console.error('提交失败:', error)
          this.$message.error('操作失败，请稍后重试')
        } finally {
          this.submitLoading = false
        }
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.resume-repo-list-container {
  padding: 20px;

  .search-card {
    margin-bottom: 20px;
  }

  .repo-list {
    .el-pagination {
      text-align: center;
      margin-top: 20px;
    }
  }
}
</style>
