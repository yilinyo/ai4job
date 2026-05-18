## ADDED Requirements
### Requirement: 用户教育经历必须以数组语义持久化
系统 MUST 在创建或更新用户资料时，将 nil 教育经历归一化为空数组，禁止向 `users.education_experiences` 写入 NULL。

#### Scenario: 注册用户未填写教育经历
- **WHEN** 后端创建新注册用户且教育经历未初始化
- **THEN** 系统将教育经历保存为空数组
- **AND** 数据库写入不触发 `education_experiences` 的 NOT NULL 约束错误

#### Scenario: 更新个人资料未传教育经历
- **WHEN** 用户更新个人资料且请求中的教育经历为 nil
- **THEN** 系统将教育经历保存为空数组
- **AND** 不向数据库写入 NULL

### Requirement: 用户查询结果必须返回数组语义
系统 MUST 在读取用户资料后，将 nil 教育经历归一化为空数组，再返回给业务层和 API 响应层。

#### Scenario: 查询用户资料时教育经历为空
- **WHEN** 后端查询用户资料且持久化结果没有教育经历内容
- **THEN** 系统返回空数组语义的教育经历
- **AND** API 响应中的 `educationExperiences` 仍为数组字段

#### Scenario: 用户已有非空教育经历
- **WHEN** 用户资料中存在一段或多段教育经历
- **THEN** 系统保留原有教育经历内容
- **AND** 不因归一化逻辑清空已有数据
