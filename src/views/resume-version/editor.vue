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
          <div class="editor-mode-toolbar">
            <el-tooltip content="编辑模式" placement="top">
              <el-button
                :type="editorMode === 'edit' ? 'primary' : 'default'"
                icon="el-icon-edit"
                circle
                @click="setEditorMode('edit')"
              />
            </el-tooltip>
            <el-tooltip content="预览模式" placement="top">
              <el-button
                :type="editorMode === 'preview' ? 'primary' : 'default'"
                icon="el-icon-notebook-2"
                circle
                @click="setEditorMode('preview')"
              />
            </el-tooltip>
          </div>

          <el-input
            v-if="editorMode === 'edit'"
            v-model="versionForm.content"
            type="textarea"
            :rows="20"
            placeholder="请输入简历内容（支持Markdown格式）"
          />
          <resume-markdown-preview
            v-else
            :html="renderedContent"
            empty-text="切换到预览模式后，这里会显示 Markdown 渲染结果"
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
import ResumeMarkdownPreview from '@/components/ResumeMarkdownPreview/index.vue'
import { ResumeVersionModule } from '@/store/modules/resume-version'
import { CreateResumeVersionRequest, UpdateResumeVersionRequest } from '@/models'
import { renderResumeMarkdown } from '@/utils/resume-markdown'

@Component({
  name: 'ResumeVersionEditor',
  components: {
    ResumeMarkdownPreview
  }
})
export default class extends Vue {
  private versionId = ''
  private repoId = ''
  private isEdit = false
  private submitLoading = false
  private editorMode: 'edit' | 'preview' = 'edit'

  private versionForm = {
    title: '',
    content: '',
    remark: ''
  }

  private versionRules = {
    title: [{ required: true, message: '请输入版本标题', trigger: 'blur' }],
    content: [{ required: true, message: '请输入简历内容', trigger: 'blur' }]
  }

  get loading() {
    return ResumeVersionModule.loading
  }

  get pageTitle() {
    return this.isEdit ? '编辑简历版本' : '创建简历版本'
  }

  get renderedContent() {
    return renderResumeMarkdown(this.versionForm.content)
  }

  created() {
    this.repoId = this.$route.query.repoId as string || ''
    this.versionId = this.$route.params.id || ''
    this.isEdit = !!this.versionId && this.$route.path.includes('/edit')

    if (this.isEdit) {
      this.loadVersion()
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
    (this.$refs.versionForm as ElForm).validate(async(valid: boolean) => {
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

  private setEditorMode(mode: 'edit' | 'preview') {
    this.editorMode = mode
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
    .editor-mode-toolbar {
      display: flex;
      justify-content: flex-end;
      gap: 10px;
      margin-bottom: 12px;
    }

    .el-textarea {
      font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
    }
  }
}
</style>
