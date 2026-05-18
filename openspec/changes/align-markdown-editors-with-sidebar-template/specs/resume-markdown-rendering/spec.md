## MODIFIED Requirements

### Requirement: 编辑页必须支持编辑模式与预览模式切换
系统 MUST 在简历版本编辑页提供统一的 `MarkdownEditor` 侧边栏模板来编辑和预览简历内容，且编辑区域内应同时呈现 Markdown 输入与预览结果，避免再依赖独立 textarea 与独立预览切换。

#### Scenario: 用户进入简历版本编辑页
- **WHEN** 用户打开简历版本创建或编辑页面
- **THEN** 系统展示统一的 Markdown 编辑模板
- **AND** 用户可以在同一编辑区域内查看 Markdown 输入和预览结果

#### Scenario: 用户保存简历内容
- **WHEN** 用户在编辑器中修改简历内容并点击保存
- **THEN** 系统保存与编辑器当前内容一致的 Markdown 文本
- **AND** 简历详情页和 PDF 导出继续按既有 Markdown 渲染规则展示内容
