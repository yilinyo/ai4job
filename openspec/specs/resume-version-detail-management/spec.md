# resume-version-detail-management Specification

## Purpose
Define the resume version detail page as the hub for viewing a resume version, exporting it, and navigating to related downstream job-search records.

## Requirements
### Requirement: 简历版本详情必须展示版本内容和基础操作
系统 MUST 在简历版本详情页展示当前简历版本内容、基础信息和可用操作入口。

#### Scenario: 用户查看简历版本详情
- **WHEN** 用户打开简历版本详情页
- **THEN** 系统展示当前版本内容和版本信息

### Requirement: 简历版本详情必须展示关联投递记录
系统 MUST 在投递记录能力存在时，在简历版本详情页展示该版本关联的投递记录区域。

#### Scenario: 用户查看关联投递区域
- **WHEN** 当前简历版本存在投递记录
- **THEN** 系统展示公司名称、岗位名称、base、投递时间和当前状态

#### Scenario: 用户从详情页进入投递详情
- **WHEN** 用户点击某条关联投递记录
- **THEN** 系统跳转到该投递记录详情页

### Requirement: 简历版本详情必须提供新增投递入口
系统 MUST 在简历版本详情页提供新增投递记录入口，并将当前版本作为创建上下文。

#### Scenario: 用户新增关联投递
- **WHEN** 用户点击新增投递记录
- **THEN** 系统进入投递创建页并带入当前简历版本
