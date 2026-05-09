# resume-application-creation-flow Specification

## Purpose
Define the creation flow for resume applications, including support for direct resume-version context and optional job posting library linkage.

## Requirements
### Requirement: 投递记录创建必须从简历版本上下文带入关联
系统 MUST 在从简历版本详情创建投递记录时自动带入当前简历版本信息。

#### Scenario: 从简历版本创建投递
- **WHEN** 用户在简历版本详情页点击新增投递记录
- **THEN** 系统进入投递记录创建页
- **AND** 自动带入当前 `resumeVersionId`

### Requirement: 投递记录创建必须支持结构化岗位信息
系统 MUST 在投递记录创建表单中支持公司名称、公司介绍、岗位名称、岗位要求、base、薪资范围、投递渠道、投递时间、当前状态和备注等字段。

#### Scenario: 用户填写投递记录
- **WHEN** 用户创建投递记录
- **THEN** 系统允许用户填写并保存结构化公司、岗位和投递信息

### Requirement: 投递记录创建必须兼容岗位库选择
系统 MUST 在岗位库能力存在时支持选择已有岗位并回填岗位字段。

#### Scenario: 用户选择岗位库岗位
- **WHEN** 用户在投递记录创建页选择已有岗位
- **THEN** 系统自动填充岗位相关字段
- **AND** 用户仍可继续补充简历版本、投递日期、当前状态和备注
