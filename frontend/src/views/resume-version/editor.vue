<template>
  <div class="resume-version-editor-container">
    <el-page-header
      @back="goBack"
      :content="pageTitle"
    />
    <div v-if="isEdit" class="auto-save-banner">
      <i class="el-icon-document-checked" />
      {{ autoSaveStatus }}
    </div>

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
            @input="markAsChanged"
            @blur="scheduleAutoSave"
          />
        </el-form-item>

        <el-form-item label="版本备注" prop="remark">
          <el-input
            v-model="versionForm.remark"
            type="textarea"
            :rows="2"
            placeholder="简单描述一下这个版本的特点或修改内容"
            @input="markAsChanged"
            @blur="scheduleAutoSave"
          />
        </el-form-item>

        <el-form-item label="简历内容" prop="content">
          <markdown-editor
            ref="markdownEditor"
            class="resume-markdown-editor"
            :initial-value="versionForm.content"
            :height="editorHeight"
            :resizable="true"
            @change="handleContentChange"
            @blur="scheduleAutoSave"
            @height-change="handleEditorHeightChange"
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
  private autoSaving = false
  private autoSavePending = false
  private autoSaveStatus = '失去焦点后自动保存'
  private lastSavedSnapshot = ''
  private autoSaveTimer?: number
  // 默认高度为原编辑区的两倍；用户可通过编辑器底部的拖拽手柄继续调整。
  private editorHeight = 'calc(200vh - 720px)'

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
      this.lastSavedSnapshot = this.formSnapshot()
    }
  }

  private handleSubmit() {
    this.cancelAutoSave()
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
    this.markAsChanged()
  }

  private markAsChanged() {
    if (this.isEdit && this.formSnapshot() !== this.lastSavedSnapshot && !this.autoSaving) {
      this.autoSaveStatus = '有未保存改动，失去焦点后自动保存'
    }
  }

  private scheduleAutoSave() {
    if (!this.isEdit || this.submitLoading) return
    if (this.autoSaving) {
      this.autoSavePending = true
      this.autoSaveStatus = '修改待自动保存…'
      return
    }
    this.cancelAutoSave()
    this.autoSaveStatus = '将在失焦后自动保存…'
    this.autoSaveTimer = window.setTimeout(() => this.autoSave(), 300)
  }

  private cancelAutoSave() {
    if (this.autoSaveTimer !== undefined) {
      window.clearTimeout(this.autoSaveTimer)
      this.autoSaveTimer = undefined
    }
  }

  private formSnapshot() {
    return JSON.stringify(this.versionForm)
  }

  private async autoSave() {
    this.autoSaveTimer = undefined
    const markdownEditor = this.$refs.markdownEditor as MarkdownEditor
    if (markdownEditor) this.versionForm.content = markdownEditor.getMarkdown()

    const snapshot = this.formSnapshot()
    if (snapshot === this.lastSavedSnapshot) {
      this.autoSaveStatus = '已保存'
      return
    }
    if (!this.versionForm.title.trim() || !this.versionForm.content.trim()) {
      this.autoSaveStatus = '标题和简历内容不能为空，未自动保存'
      return
    }

    this.autoSaving = true
    this.autoSaveStatus = '正在自动保存…'
    try {
      const data: UpdateResumeVersionRequest = {
        title: this.versionForm.title,
        content: this.versionForm.content,
        remark: this.versionForm.remark
      }
      await ResumeVersionModule.UpdateResumeVersion({
        id: this.versionId,
        data
      })
      this.lastSavedSnapshot = snapshot
      this.autoSaveStatus = '已自动保存'
    } catch (error) {
      console.error('自动保存失败:', error)
      this.autoSaveStatus = '自动保存失败，请点击“保存修改”重试'
    } finally {
      this.autoSaving = false
      if (this.autoSavePending || this.formSnapshot() !== this.lastSavedSnapshot) {
        this.autoSavePending = false
        this.scheduleAutoSave()
      }
    }
  }

  private handleEditorHeightChange(height: string) {
    this.editorHeight = height
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
    this.cancelAutoSave()
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

  .auto-save-banner {
    margin: -8px 0 16px;
    color: #409eff;
    font-size: 13px;

    i {
      margin-right: 4px;
    }
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
