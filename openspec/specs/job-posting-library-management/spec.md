# job-posting-library-management Specification

## Purpose
Define the standalone job posting library used to collect, manage, filter, view, and maintain job postings independently from resume applications.

## Requirements
### Requirement: 岗位库必须作为独立一级模块管理岗位
系统 MUST 提供独立岗位库模块，用于创建、查看、编辑、删除和筛选岗位信息，并在左侧导航中作为与简历管理同级的一级入口展示。

#### Scenario: 用户查看岗位列表
- **WHEN** 用户进入岗位列表页
- **THEN** 系统展示岗位卡片列表、筛选区、分页和新增岗位入口
- **AND** 岗位列表作为一级导航模块展示

#### Scenario: 用户筛选岗位
- **WHEN** 用户按关键词、来源平台或岗位状态筛选岗位
- **THEN** 系统只展示匹配条件的岗位数据

### Requirement: 岗位卡片必须展示核心岗位信息和操作
系统 MUST 在岗位卡片中展示公司名称、岗位名称、base、薪资范围、来源平台、岗位状态和创建时间，并提供查看、编辑、删除、查看原岗位操作。

#### Scenario: 用户查看岗位卡片
- **WHEN** 岗位列表加载成功
- **THEN** 每张岗位卡片展示岗位核心信息和操作入口
- **AND** 卡片尺寸、内容区高度、留白和 footer 位置与简历仓库卡片风格保持一致

### Requirement: 岗位详情页必须集中展示岗位信息和关联投递
系统 MUST 提供岗位详情页，展示岗位摘要、岗位基础信息与操作区域，以及该岗位关联的投递记录。

#### Scenario: 用户进入岗位详情
- **WHEN** 用户点击岗位卡片查看操作
- **THEN** 系统展示公司、岗位、状态、来源链接、岗位要求、备注和关联投递记录
- **AND** 创建投递记录、编辑岗位、查看原岗位操作位于岗位基础信息与操作区域

### Requirement: 查看原岗位必须打开来源链接
系统 MUST 在岗位存在来源链接时提供查看原岗位操作，并新开页面访问该来源地址。

#### Scenario: 用户查看原岗位
- **WHEN** 用户点击查看原岗位
- **THEN** 系统打开该岗位的 `sourceUrl`
