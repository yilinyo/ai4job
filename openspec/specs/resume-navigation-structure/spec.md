# resume-navigation-structure Specification

## Purpose
Define resume-related navigation boundaries so account-level profile management and resume management are not duplicated or confused.

## Requirements
### Requirement: 简历管理下不得重复展示个人信息入口
系统 MUST 将账号级个人信息从简历管理二级菜单中移除，避免与顶层个人信息入口重复。

#### Scenario: 用户查看简历管理菜单
- **WHEN** 用户展开简历管理菜单
- **THEN** 系统不展示个人信息子入口

### Requirement: 个人信息入口必须与 Guide 同级
系统 MUST 将业务个人信息入口放在顶层导航，与 Guide 同级展示。

#### Scenario: 用户查看顶层菜单
- **WHEN** 用户登录后查看左侧顶层导航
- **THEN** 系统展示“个人信息”入口
- **AND** 用户点击后进入账号级个人信息表单页
