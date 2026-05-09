# resume-repo-job-type-management Specification

## Purpose
Define the supported job-seeking types for resume repositories and ensure they remain consistent across data models, Mock data, forms, filters, cards, and detail views.

## Requirements
### Requirement: 简历仓库求职类型必须支持校招、社招和实习
系统 MUST 将“实习”作为与“校招”“社招”同级的简历仓库求职类型。

#### Scenario: 用户创建实习类型仓库
- **WHEN** 用户在创建简历仓库表单中选择“实习”
- **THEN** 系统以 `internship` 保存该仓库的求职类型

#### Scenario: 用户编辑为实习类型
- **WHEN** 用户编辑已有简历仓库并选择“实习”
- **THEN** 系统正确回显并保存 `internship`

### Requirement: 简历仓库列表必须支持实习筛选
系统 MUST 在简历仓库列表的求职类型筛选中提供“实习”选项。

#### Scenario: 用户筛选实习仓库
- **WHEN** 用户选择“实习”筛选条件
- **THEN** 系统只展示 `jobType === 'internship'` 的简历仓库

### Requirement: 简历仓库展示必须正确显示实习标签
系统 MUST 在简历仓库卡片和详情页中正确展示“实习”文案和独立标签样式。

#### Scenario: 用户查看实习仓库卡片
- **WHEN** 简历仓库类型为 `internship`
- **THEN** 卡片展示“实习”标签
- **AND** 标签样式与校招、社招可区分且风格一致

#### Scenario: 用户查看实习仓库详情
- **WHEN** 用户打开实习类型简历仓库详情页
- **THEN** 详情页求职类型字段展示“实习”
