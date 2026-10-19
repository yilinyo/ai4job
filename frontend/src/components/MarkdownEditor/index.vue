<template>
  <div class="markdown-editor-wrapper" @focusout="handleEditorBlur">
    <div :id="id" />
    <div
      v-if="resizable"
      class="resize-handle"
      title="拖拽调节编辑器高度"
      @mousedown.prevent="startResize"
    >
      <i class="el-icon-rank" />
      <span>拖拽调节高度</span>
    </div>
  </div>
</template>

<script lang="ts">
import 'codemirror/lib/codemirror.css' // Editor's Dependency Style
import '@toast-ui/editor/dist/toastui-editor.css' // Editor's Style
import '@toast-ui/editor/dist/i18n/es-es'
import '@toast-ui/editor/dist/i18n/it-it'
import '@toast-ui/editor/dist/i18n/ja-jp'
import '@toast-ui/editor/dist/i18n/ko-kr'
import '@toast-ui/editor/dist/i18n/zh-cn'
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import Editor, { EditorOptions } from '@toast-ui/editor'
import defaultOptions from './default-options'

const defaultId = () => 'markdown-editor-' + +new Date() + ((Math.random() * 1000).toFixed(0) + '')

@Component({
  name: 'MarkdownEditor'
})
export default class extends Vue {
  @Prop({ default: '' }) private initialValue!: string
  @Prop({ default: defaultId }) private id!: string
  @Prop({ default: () => defaultOptions }) private options!: EditorOptions
  @Prop({ default: 'markdown' }) private mode!: string
  @Prop({ default: '300px' }) private height!: string
  @Prop({ default: 'en' }) private language!: string
  @Prop({ default: false }) private resizable!: boolean

  private markdownEditor?: Editor
  private resizeStartY = 0
  private resizeStartHeight = 0
  // Mapping for local lang to tuiEditor lang
  // https://github.com/nhn/tui.editor/blob/master/apps/editor/docs/i18n.md#supported-languages
  private languageTypeList: { [key: string]: string } = {
    en: 'en',
    zh: 'zh-CN',
    es: 'es',
    ja: 'ja',
    ko: 'ko',
    it: 'it'
  }

  get editorOptions() {
    const options = Object.assign({}, defaultOptions, this.options)
    const hooks = Object.assign({}, this.options.hooks || {})
    if (!hooks.addImageBlobHook) {
      hooks.addImageBlobHook = this.handleImageBlobHook
    }
    options.hooks = hooks
    options.initialEditType = this.mode
    options.initialValue = this.initialValue
    options.height = this.height
    options.language = this.languageTypeList[this.language]
    return options
  }

  @Watch('initialValue')
  private onInitialValueChange(value: string) {
    if (!this.markdownEditor) return
    if (value === this.markdownEditor.getMarkdown()) return
    this.markdownEditor.setMarkdown(value || '', false)
  }

  @Watch('language')
  private onLanguageChange() {
    this.destroyEditor()
    this.initEditor()
  }

  @Watch('height')
  private onHeightChange(value: string) {
    if (this.markdownEditor) {
      this.markdownEditor.height(value)
    }
  }

  @Watch('mode')
  private onModeChange(value: string) {
    if (this.markdownEditor) {
      this.markdownEditor.changeMode(value)
    }
  }

  mounted() {
    this.initEditor()
  }

  destroyed() {
    this.stopResize()
    this.destroyEditor()
  }

  private initEditor() {
    const editorElement = document.getElementById(this.id)
    if (!editorElement) return
    // eslint-disable-next-line new-cap
    this.markdownEditor = new Editor({
      ...this.editorOptions,
      el: editorElement
    })
    this.markdownEditor.on('change', this.handleEditorChange)
    this.markdownEditor.on('blur', this.handleEditorBlur)
  }

  private destroyEditor() {
    if (!this.markdownEditor) return
    this.markdownEditor.off('change')
    this.markdownEditor.off('blur')
    this.markdownEditor.remove()
    this.markdownEditor = undefined
  }

  private handleEditorChange = () => {
    if (!this.markdownEditor) return
    const markdown = this.markdownEditor.getMarkdown()
    this.$emit('change', markdown)
    this.$emit('input', markdown)
  }

  private handleEditorBlur() {
    this.$emit('blur')
  }

  private startResize(event: MouseEvent) {
    const editorElement = document.getElementById(this.id)
    if (!editorElement) return

    this.resizeStartY = event.clientY
    this.resizeStartHeight = editorElement.getBoundingClientRect().height
    document.addEventListener('mousemove', this.resizeEditor)
    document.addEventListener('mouseup', this.stopResize)
  }

  private resizeEditor = (event: MouseEvent) => {
    const nextHeight = Math.max(300, Math.round(this.resizeStartHeight + event.clientY - this.resizeStartY))
    this.$emit('height-change', `${nextHeight}px`)
  }

  private stopResize = () => {
    document.removeEventListener('mousemove', this.resizeEditor)
    document.removeEventListener('mouseup', this.stopResize)
  }

  private handleImageBlobHook(blob: File | Blob, callback: (url: string, altText: string) => void) {
    const reader = new FileReader()
    reader.onload = () => {
      callback(String(reader.result || ''), '粘贴图片')
    }
    reader.readAsDataURL(blob)
  }

  public focus() {
    if (this.markdownEditor) {
      this.markdownEditor.focus()
    }
  }

  public getMarkdown() {
    if (this.markdownEditor) {
      return this.markdownEditor.getMarkdown()
    }
    return this.initialValue
  }

  public setHtml(value: string) {
    if (this.markdownEditor) {
      this.markdownEditor.setHtml(value)
    }
  }

  public setMarkdown(value: string) {
    if (this.markdownEditor) {
      this.markdownEditor.setMarkdown(value || '', false)
    }
  }

  public getHtml() {
    if (this.markdownEditor) {
      return this.markdownEditor.getHtml()
    }
    return ''
  }
}
</script>

<style lang="scss" scoped>
.markdown-editor-wrapper {
  width: 100%;
}

.resize-handle {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 28px;
  color: #909399;
  font-size: 12px;
  cursor: ns-resize;
  user-select: none;

  &:hover {
    color: #409eff;
  }
}
</style>
