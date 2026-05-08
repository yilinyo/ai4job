## ADDED Requirements

### Requirement: 简历内容必须按 Markdown 语义渲染
系统必须在简历版本编辑页和详情页中将简历内容按 Markdown 语义解析并渲染，而不是直接展示原始 Markdown 标记文本。

#### Scenario: 标题与列表按 Markdown 显示
- **WHEN** 用户查看包含 `## 教育经历` 和 `- 熟悉 Vue` 等 Markdown 语法的简历内容
- **THEN** 系统以标题和列表形式渲染这些内容，而不是显示原始 `##` 和 `-` 标记

#### Scenario: 段落与强调语法正确显示
- **WHEN** 用户查看包含段落、加粗或引用语法的简历内容
- **THEN** 系统按照对应的 Markdown 结构和样式展示内容

#### Scenario: Markdown 表格正确显示
- **WHEN** 用户查看包含标准 Markdown 表格语法的简历内容
- **THEN** 系统将其渲染为可读的表格结构，而不是显示原始管道符和分隔线文本

### Requirement: 编辑页必须支持编辑模式与预览模式切换
系统必须在简历版本编辑页提供清晰的编辑模式和预览模式切换能力，交互上接近 VSCode 的 Markdown 使用方式。

#### Scenario: 用户切换到编辑模式
- **WHEN** 用户点击编辑模式图标（如小铅笔）
- **THEN** 系统显示可编辑的 Markdown 输入区域，并允许用户继续修改简历内容

#### Scenario: 用户切换到预览模式
- **WHEN** 用户点击预览模式图标（如小书本）
- **THEN** 系统在当前编辑区域内显示该份简历的 Markdown 渲染结果，而不是弹出独立预览窗口

### Requirement: 页面预览与 PDF 导出必须使用一致的 Markdown 内容结构
系统必须让简历编辑页预览、简历详情页预览和 PDF 导出基于一致的 Markdown 解析结果，避免不同场景展示不同内容结构。

#### Scenario: 页面与导出结构一致
- **WHEN** 用户预览一份 Markdown 简历并立即导出 PDF
- **THEN** PDF 中的标题、列表、段落和强调结构与页面预览保持一致

#### Scenario: 编辑页预览与详情页一致
- **WHEN** 用户在编辑页切换到预览模式后再进入简历详情页
- **THEN** 两处看到的 Markdown 结构和主要样式层级保持一致

### Requirement: Markdown 渲染结果必须具备基础安全控制
系统必须在将 Markdown 渲染结果注入页面或导出容器前进行安全控制，并按纯 Markdown 策略处理内容，不保留原生 HTML 标签。

#### Scenario: 不安全标签不会直接执行
- **WHEN** 简历内容中包含潜在不安全的原始 HTML 片段
- **THEN** 系统不会将其作为可执行内容直接注入预览或导出文档

#### Scenario: 原生 HTML 不作为格式化能力保留
- **WHEN** 用户在简历内容中写入原生 HTML 标签
- **THEN** 系统不会依赖这些标签提供格式化展示能力，而是仅以受支持的 Markdown 语法作为渲染来源

### Requirement: Markdown 渲染必须适配简历阅读与导出排版
系统必须为常见 Markdown 元素提供适合简历阅读和 PDF 导出的排版样式，并覆盖 Markdown 表格。

#### Scenario: 常见 Markdown 元素具备可读样式
- **WHEN** 简历内容包含标题、列表、引用和代码块等常见 Markdown 元素
- **THEN** 系统为这些元素应用适合简历阅读的间距、字号和层级样式

#### Scenario: 表格具备可读样式
- **WHEN** 简历内容包含 Markdown 表格
- **THEN** 系统为表头、边框、单元格间距和内容对齐提供可读的预览与导出样式

#### Scenario: 长内容导出具备可接受分页
- **WHEN** 用户导出包含多个 Markdown 分节和列表的长简历
- **THEN** 导出的 PDF 在主要分节和内容块上具备可接受的分页与阅读体验
