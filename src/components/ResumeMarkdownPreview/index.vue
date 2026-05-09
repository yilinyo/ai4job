<template>
  <div
    class="resume-markdown-preview"
    :class="[`resume-markdown-preview--${variant}`, {'is-empty': !html}]"
  >
    <div
      v-if="html"
      class="resume-markdown-preview__content"
      v-html="html"
    />
    <div v-else class="resume-markdown-preview__empty">
      {{ emptyText }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component({
  name: 'ResumeMarkdownPreview'
})
export default class extends Vue {
  @Prop({ default: '' }) private html!: string
  @Prop({ default: 'screen' }) private variant!: 'screen' | 'pdf'
  @Prop({ default: '暂无内容可预览' }) private emptyText!: string
}
</script>

<style lang="scss" scoped>
.resume-markdown-preview {
  color: #1f2937;

  &.is-empty {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 220px;
    color: #9ca3af;
    background-color: #f8fafc;
  }

  &__content {
    line-height: 1.8;
    word-break: break-word;

    ::v-deep h1,
    ::v-deep h2,
    ::v-deep h3,
    ::v-deep h4,
    ::v-deep h5,
    ::v-deep h6 {
      margin: 1.2em 0 0.6em;
      color: #111827;
      font-weight: 600;
      line-height: 1.35;
      page-break-after: avoid;
    }

    ::v-deep h1 {
      font-size: 28px;
      border-bottom: 2px solid #dbe3f0;
      padding-bottom: 0.35em;
    }

    ::v-deep h2 {
      font-size: 22px;
    }

    ::v-deep h3 {
      font-size: 18px;
    }

    ::v-deep p,
    ::v-deep ul,
    ::v-deep ol,
    ::v-deep blockquote,
    ::v-deep pre,
    ::v-deep table {
      margin: 0 0 1em;
    }

    ::v-deep ul,
    ::v-deep ol {
      padding-left: 1.5em;
    }

    ::v-deep li + li {
      margin-top: 0.35em;
    }

    ::v-deep blockquote {
      padding: 0.8em 1em;
      border-left: 4px solid #93c5fd;
      background-color: #eff6ff;
      color: #334155;
    }

    ::v-deep pre {
      overflow-x: auto;
      padding: 12px 14px;
      border-radius: 8px;
      background-color: #0f172a;
      color: #e2e8f0;
      white-space: pre-wrap;
    }

    ::v-deep code {
      padding: 0.15em 0.4em;
      border-radius: 4px;
      background-color: #eef2ff;
      color: #1d4ed8;
      font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
      font-size: 0.92em;
    }

    ::v-deep pre code {
      padding: 0;
      background-color: transparent;
      color: inherit;
    }

    ::v-deep table {
      width: 100%;
      border-collapse: collapse;
      table-layout: fixed;
      page-break-inside: avoid;
    }

    ::v-deep th,
    ::v-deep td {
      border: 1px solid #dbe3f0;
      padding: 10px 12px;
      text-align: left;
      vertical-align: top;
    }

    ::v-deep th {
      background-color: #eff6ff;
      color: #1e3a8a;
      font-weight: 600;
    }

    ::v-deep a {
      color: #2563eb;
      text-decoration: none;
    }

    ::v-deep img {
      display: block;
      max-width: 100%;
      max-height: 360px;
      margin: 12px 0;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      object-fit: contain;
    }

    ::v-deep hr {
      border: none;
      border-top: 1px solid #dbe3f0;
      margin: 1.5em 0;
    }

    ::v-deep > :first-child {
      margin-top: 0;
    }

    ::v-deep > :last-child {
      margin-bottom: 0;
    }
  }

  &--screen {
    padding: 20px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    background-color: #fff;
    min-height: 360px;
  }

  &--pdf {
    padding: 0;
    background-color: transparent;

    .resume-markdown-preview__content {
      font-size: 14px;
      line-height: 1.75;
    }

    .resume-markdown-preview__content ::v-deep h1 {
      font-size: 26px;
    }

    .resume-markdown-preview__content ::v-deep h2 {
      font-size: 18px;
      color: #1d4ed8;
    }

    .resume-markdown-preview__content ::v-deep table {
      font-size: 13px;
    }
  }
}
</style>
