## Why

简历版本编辑页和题库题目编辑页目前使用了不同的 Markdown 输入形态，用户在两个高频入口之间切换时需要重新适应编辑方式。仓库里已经有现成的 `MarkdownEditor` 侧边栏模板，但还没有统一接入这两个编辑场景，因此 Markdown 编辑体验仍然割裂。

## What Changes

- 简历版本编辑页将改为复用统一的 MarkdownEditor 侧边栏模板，不再依赖独立的纯 textarea 编辑体验。
- 题库新增/编辑弹窗中的题目内容和答案解析将改为复用同一 MarkdownEditor 模板。
- 保留现有的 Markdown 语义、图片粘贴、保存校验和预览结果。
- 不新增后端接口，不变更现有数据结构，不做数据迁移。

## Capabilities

### New Capabilities

- 无（复用现有能力）

### Modified Capabilities

- `resume-markdown-rendering`: 简历版本编辑页的 Markdown 编辑交互将统一为侧边栏模板，并保持现有渲染效果一致。
- `interview-question-library`: 题库新增/编辑入口的 Markdown 编辑交互将统一为侧边栏模板，题目内容和答案解析的编辑体验保持一致。

## Impact

- 受影响前端页面：`frontend/src/views/resume-version/editor.vue`、`frontend/src/components/InterviewQuestion/QuestionFormDialog.vue`
- 受影响共用组件：`frontend/src/components/MarkdownEditor/index.vue`
- 现有只读预览能力保持不变，`ResumeMarkdownPreview` 继续用于详情和预览场景
- 无后端 API 影响
