# resume-version-prefill-from-user-profile Specification

## Purpose
Define how creating a resume version uses the current account profile to generate an editable default resume draft.

## Requirements
### Requirement: 创建简历版本必须可从个人信息自动生成草稿
系统 MUST 在新建简历版本且内容为空时，根据当前账号个人信息自动生成默认简历内容草稿。

#### Scenario: 新建简历版本自动填充
- **WHEN** 用户创建新的简历版本且内容为空
- **THEN** 系统读取当前账号个人信息
- **AND** 生成包含基本信息、教育背景和个人优势的默认简历草稿

### Requirement: 自动草稿必须保持可编辑
系统 MUST 将自动生成的草稿写入现有简历编辑器，使用户可以继续编辑修改。

#### Scenario: 用户修改自动草稿
- **WHEN** 系统生成默认简历草稿后
- **THEN** 用户可以在简历编辑器中继续修改并保存

### Requirement: 自动填充不得覆盖已有内容
系统 MUST 只在新建模式且内容为空时自动填充，不得覆盖用户已输入内容或已有简历版本内容。

#### Scenario: 编辑已有简历版本
- **WHEN** 用户编辑已有简历版本
- **THEN** 系统不使用个人信息自动覆盖已有内容

#### Scenario: 新建时用户已有输入
- **WHEN** 用户已经在内容字段输入文本
- **THEN** 系统不覆盖用户输入

### Requirement: 空字段不得生成错误占位内容
系统 MUST 在生成草稿时跳过或合理处理空字段，避免出现 `undefined` 等错误文本。

#### Scenario: 个人信息字段为空
- **WHEN** 当前账号部分个人信息为空
- **THEN** 生成的简历草稿不包含明显错误内容或 `undefined` 文本
