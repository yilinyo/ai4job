# resume-application-tracking Specification

## Purpose
Define how resume versions track multiple job applications and preserve the job-search history associated with each resume version.

## Requirements
### Requirement: 简历版本必须支持关联多个投递记录
系统 MUST 允许一个简历版本关联多条投递记录，并在简历版本详情页展示这些记录。

#### Scenario: 用户查看简历版本关联投递
- **WHEN** 用户进入简历版本详情页
- **THEN** 系统展示该版本关联的投递记录数量和列表

#### Scenario: 用户新增投递记录
- **WHEN** 用户从简历版本详情页新增投递记录
- **THEN** 新投递记录关联到当前简历版本

### Requirement: 投递记录必须保留岗位级求职信息
系统 MUST 在投递记录中保存公司名称、公司介绍、岗位名称、岗位要求、base、薪资范围、投递渠道、投递时间、当前状态和备注。

#### Scenario: 用户查看投递记录详情
- **WHEN** 用户打开投递记录详情页
- **THEN** 系统展示该投递的公司、岗位、投递和状态信息

### Requirement: 已有关联投递记录的简历版本不得静默删除
系统 MUST 阻止或明确提示用户不能直接删除已有关联投递记录的简历版本。

#### Scenario: 删除已关联投递的简历版本
- **WHEN** 用户尝试删除存在投递记录的简历版本
- **THEN** 系统阻止删除并提示存在关联投递记录
