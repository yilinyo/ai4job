# job-posting-link-parse Specification

## Purpose
Define the frontend flow for importing job posting information from a source URL using a backend parsing API, with Mock-first support.

## Requirements
### Requirement: 岗位创建编辑必须支持链接解析导入
系统 MUST 在岗位创建或编辑流程中允许用户输入岗位来源链接，并调用解析接口获取结构化岗位数据回填表单。

#### Scenario: 用户解析岗位链接
- **WHEN** 用户输入岗位来源链接并点击解析导入
- **THEN** 系统调用岗位链接解析接口
- **AND** 解析成功后将公司名称、岗位名称、岗位要求、base、薪资范围、来源平台和来源链接等字段回填到岗位表单

### Requirement: 解析后用户必须可以继续编辑字段
系统 MUST 允许用户在解析结果回填后继续手动补充或修正岗位字段。

#### Scenario: 用户修改解析结果
- **WHEN** 岗位解析结果回填到表单后
- **THEN** 用户仍可编辑表单字段并保存最终岗位信息

### Requirement: 前端不得承担真实网页抓取职责
系统 MUST 将真实网页抓取和解析职责保留在后端边界内，前端只消费解析接口返回的结构化 JSON。

#### Scenario: 前端发起解析请求
- **WHEN** 用户触发链接解析
- **THEN** 前端只提交来源链接并消费接口响应
- **AND** 当前阶段可通过 Mock 返回结构化岗位数据
