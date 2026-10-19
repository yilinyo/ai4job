<template>
  <div class="resume-version-editor-container">
    <el-page-header
      @back="goBack"
      :content="pageTitle"
    />

    <el-card v-loading="loading" class="editor-card">
      <el-form
        ref="versionForm"
        :model="versionForm"
        :rules="versionRules"
        label-width="100px"
      >
        <el-form-item label="版本标题" prop="title">
          <el-input
            v-model="versionForm.title"
            placeholder="请输入版本标题，如：互联网大厂通用版"
          />
        </el-form-item>

        <el-form-item label="版本备注" prop="remark">
          <el-input
            v-model="versionForm.remark"
            type="textarea"
            :rows="2"
            placeholder="简单描述一下这个版本的特点或修改内容"
          />
        </el-form-item>

        <el-form-item label="简历内容" prop="content">
          <markdown-editor
            ref="markdownEditor"
            class="resume-markdown-editor"
            :initial-value="versionForm.content"
            :height="editorHeight"
            @change="handleContentChange"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="submitLoading"
            @click="handleSubmit"
          >
            {{ isEdit ? '保存修改' : '创建版本' }}
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
import MarkdownEditor from '@/components/MarkdownEditor/index.vue'
import { ResumeRepoModule } from '@/store/modules/resume-repo'
import { ResumeVersionModule } from '@/store/modules/resume-version'
import { UserModule } from '@/store/modules/user'
import { CreateResumeVersionRequest, UpdateResumeVersionRequest } from '@/models'
import { generateResumeDraftFromProfile } from '@/utils/resume-prefill'

@Component({
  name: 'ResumeVersionEditor',
  components: {
    MarkdownEditor
  }
})
export default class extends Vue {
  private versionId = ''
  private repoId = ''
  private isEdit = false
  private submitLoading = false
  private editorHeight = 'calc(100vh - 360px)'

  private versionForm = {
    title: '',
    content: '',
    remark: ''
  }

  private versionRules = {
    title: [{ required: true, message: '请输入版本标题', trigger: 'blur' }],
    content: [{ required: true, message: '请输入简历内容', trigger: 'change' }]
  }

  get loading() {
    return ResumeVersionModule.loading
  }

  get pageTitle() {
    return this.isEdit ? '编辑简历版本' : '创建简历版本'
  }

  created() {
    this.repoId = this.$route.query.repoId as string || ''
    this.versionId = this.$route.params.id || ''
    this.isEdit = !!this.versionId && this.$route.path.includes('/edit')

    if (this.isEdit) {
      this.loadVersion()
    } else {
      this.prefillContentFromUserProfile()
    }
  }

  private async loadVersion() {
    await ResumeVersionModule.GetResumeVersionById(this.versionId)
    const version = ResumeVersionModule.currentVersion
    if (version) {
      this.versionForm = {
        title: version.title,
        content: version.content,
        remark: version.remark || ''
      }
      this.repoId = version.repoId
    }
  }

  private handleSubmit() {
    const markdownEditor = this.$refs.markdownEditor as MarkdownEditor
    this.versionForm.content = markdownEditor
      ? markdownEditor.getMarkdown()
      : this.versionForm.content

    const versionForm = this.$refs.versionForm as ElForm
    versionForm.validate(async(valid: boolean) => {
      if (valid) {
        this.submitLoading = true
        try {
          if (this.isEdit) {
            const data: UpdateResumeVersionRequest = {
              title: this.versionForm.title,
              content: this.versionForm.content,
              remark: this.versionForm.remark
            }
            await ResumeVersionModule.UpdateResumeVersion({
              id: this.versionId,
              data
            })
            this.$message.success('更新成功！')
          } else {
            const data: CreateResumeVersionRequest = {
              repoId: this.repoId,
              title: this.versionForm.title,
              content: this.versionForm.content,
              remark: this.versionForm.remark
            }
            await ResumeVersionModule.CreateResumeVersion(data)
            await ResumeRepoModule.GetResumeRepoById(this.repoId)
            this.$message.success('创建成功！')
          }
          this.goBack()
        } catch (error) {
          console.error('提交失败:', error)
          this.$message.error('操作失败，请稍后重试')
        } finally {
          this.submitLoading = false
        }
      }
    })
  }

  private handleContentChange(value: string) {
    this.versionForm.content = value
  }

  private async prefillContentFromUserProfile() {
    if (this.versionForm.content.trim()) return

    try {
      await UserModule.GetUserInfo()
      if (this.versionForm.content.trim()) return

      const draft = generateResumeDraftFromProfile(UserModule.userProfile)
      if (draft) {
        this.versionForm.content = draft
      }
    } catch (error) {
      console.error('生成默认简历草稿失败:', error)
    }
  }

  private goBack() {
    if (this.repoId) {
      this.$router.push(`/resume/repo/${this.repoId}`)
    } else {
      this.$router.back()
    }
  }
}
</script>

<style lang="scss" scoped>
.resume-version-editor-container {
  padding: 20px;

  .el-page-header {
    margin-bottom: 20px;
  }

  .editor-card {
    max-width: 1280px;
    margin: 0 auto;

    .resume-markdown-editor {
      display: block;
      width: 100%;
    }
  }
}
</style>
