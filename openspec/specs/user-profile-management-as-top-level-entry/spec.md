# user-profile-management-as-top-level-entry Specification

## Purpose
Define account-level personal profile management as a top-level navigation entry with structured form fields and account-scoped Mock persistence.

## Requirements
### Requirement: 个人信息页必须作为一级导航入口
系统 MUST 将业务个人信息页作为与 Guide 同级的一级菜单入口展示，并移除简历管理下的重复个人信息入口。

#### Scenario: 用户查看左侧菜单
- **WHEN** 用户登录后查看左侧导航
- **THEN** 系统展示一级“个人信息”入口
- **AND** “简历管理”下不再展示个人信息子入口

### Requirement: 个人信息页必须使用结构化表单维护资料
系统 MUST 通过表单维护姓名、年龄、邮箱、电话、求职意向、头像、所在地、个人优势和教育经历等账号级资料。

#### Scenario: 用户编辑基础资料
- **WHEN** 用户进入个人信息页
- **THEN** 系统展示结构化表单字段
- **AND** 用户可以保存或重置当前账号资料

### Requirement: 教育经历必须支持多段且学历唯一
系统 MUST 支持多段教育经历，每段包含学校、学历、专业、入学时间和毕业时间，并限制同一种学历只能维护一段。

#### Scenario: 用户添加教育经历
- **WHEN** 用户添加一段教育经历
- **THEN** 系统允许填写学校、学历、专业、入学时间和毕业时间

#### Scenario: 用户添加重复学历
- **WHEN** 用户尝试添加已存在学历的教育经历
- **THEN** 系统阻止保存或提示同一种学历只允许一段

#### Scenario: 用户设置错误毕业时间
- **WHEN** 用户设置毕业时间早于入学时间
- **THEN** 系统提示时间范围不合法
