# job-posting-to-application-linking Specification

## Purpose
Define how job postings connect to resume applications while preserving stable job snapshots for historical review.

## Requirements
### Requirement: 岗位详情必须支持创建关联投递记录
系统 MUST 允许用户从岗位详情页创建投递记录，并自动带入该岗位相关字段。

#### Scenario: 从岗位详情创建投递记录
- **WHEN** 用户在岗位详情页点击创建投递记录
- **THEN** 系统跳转到投递记录创建页
- **AND** 自动带入 `jobPostingId` 和岗位字段

### Requirement: 投递记录创建必须支持选择已有岗位
系统 MUST 允许用户在投递记录创建流程中选择岗位库中的已有岗位，并自动回填岗位信息。

#### Scenario: 创建投递时选择岗位
- **WHEN** 用户在投递记录创建页选择已有岗位
- **THEN** 系统回填公司名称、岗位名称、岗位要求、base、薪资范围、来源平台和来源链接

### Requirement: 投递记录必须同时保存岗位关联和岗位快照
系统 MUST 在保存投递记录时写入 `jobPostingId`，并保留岗位快照字段，避免岗位库后续修改影响历史投递语义。

#### Scenario: 保存关联岗位的投递记录
- **WHEN** 用户保存从岗位库带出的投递记录
- **THEN** 投递记录保存岗位库关联 ID
- **AND** 投递记录保存当时的公司、岗位、要求、base、薪资和来源等快照字段
