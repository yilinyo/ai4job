## MODIFIED Requirements

### Requirement: 题库列表必须支持题目管理和筛选
系统 MUST 允许用户在题库列表中新增、编辑、查看、删除面试题，并支持按关键词、标签、难度、题型、掌握状态和收藏状态筛选题目。新增和编辑题目时，题目内容与答案解析 MUST 使用统一的 `MarkdownEditor` 侧边栏模板进行编辑，以保持两个字段的 Markdown 体验一致。

#### Scenario: 用户新增面试题
- **WHEN** 用户创建一条面试题并填写题目内容和答案解析
- **THEN** 系统使用统一的 Markdown 编辑模板承载两个 Markdown 字段
- **AND** 题目内容和答案解析都可以按 Markdown 语义编辑和预览
- **AND** 用户粘贴图片时系统继续以 Markdown 图片形式保存

#### Scenario: 用户编辑面试题
- **WHEN** 用户打开已有题目的编辑弹窗
- **THEN** 系统将当前题目内容和答案解析预填到统一编辑模板中
- **AND** 用户保存后题库列表和详情页继续按既有规则展示 Markdown 内容
