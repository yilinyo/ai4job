## Why

当前简历详情页和 PDF 导出都把简历内容当作纯文本处理，导致 `##`、列表、加粗等 Markdown 语法原样显示，用户看到的预览和下载结果都不符合简历展示预期。与此同时，编辑页仍然只有纯文本输入框和独立弹窗预览，不支持像 VSCode 那样在编辑模式与预览模式之间顺滑切换。既然系统已经允许用户以 Markdown 方式维护简历内容，就需要提供真实、一致的 Markdown 渲染与编辑体验。

## What Changes

- 将简历详情页的内容展示从“换行替换”升级为真实 Markdown 渲染。
- 让 PDF 导出复用同一份 Markdown 渲染结果，确保页面预览与导出结果一致。
- 将简历编辑页升级为支持“编辑模式 / 预览模式”切换的 Markdown 编辑体验，交互上接近 VSCode 的编辑预览方式。
- 为 Markdown 内容补充基础排版样式，包括标题、列表、段落、强调、引用和代码块等常见结构。
- 增加对空内容、非法 HTML 注入风险和导出场景版式差异的处理，避免渲染结果不可控。

## Capabilities

### New Capabilities
- `resume-markdown-rendering`: 将简历版本内容按真实 Markdown 语义渲染，并在编辑页、详情页与 PDF 导出中保持一致的展示结果。

### Modified Capabilities

## Impact

- 受影响页面：`src/views/resume-version/view.vue`、`src/views/resume-version/editor.vue`
- 可能新增前端依赖：Markdown 解析器与 HTML 清洗相关工具
- 可能新增共享渲染/格式化工具、编辑页模式切换逻辑与导出样式
- 不涉及后端接口、数据结构或 Mock API 变更
