# resume-application-creation-flow Specification

## Purpose
Define the creation flow for resume applications, including support for direct resume-version context, optional job posting library linkage, and resume version activation checks.

## Requirements
### Requirement: 投递记录创建必须从简历版本上下文带入关联
系统 MUST 在从简历版本详情创建投递记录时自动带入当前简历版本信息；若当前版本未激活，则不得允许用户继续发起新增投递。

#### Scenario: 从简历版本创建投递
- **WHEN** 用户在简历版本详情页点击新增投递记录
- **THEN** 系统进入投递记录创建页
- **AND** 自动带入当前 `resumeVersionId`

#### Scenario: 从未激活版本查看详情
- **WHEN** 用户进入未激活的简历版本详情页
- **THEN** 系统不展示“新增投递记录”入口
- **AND** 用户不能从该页面直接发起新的投递创建

### Requirement: 投递记录创建必须拒绝未激活简历版本
系统 MUST 在创建投递记录时校验 `resumeVersionId` 对应版本处于激活状态；未激活版本不得创建新的投递记录。

#### Scenario: 用户提交未激活版本创建投递
- **WHEN** 用户通过手工构造请求、绕过前端入口或其它链路提交未激活版本的 `resumeVersionId`
- **THEN** 系统拒绝创建投递记录
- **AND** 不写入新的投递记录数据

#### Scenario: 用户提交激活版本创建投递
- **WHEN** 用户提交激活版本的 `resumeVersionId` 创建投递记录
- **THEN** 系统允许继续保存结构化公司、岗位和投递信息

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
