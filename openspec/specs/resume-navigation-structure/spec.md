# resume-navigation-structure Specification

## Purpose
Define resume-related navigation boundaries so account-level profile management and resume management are not duplicated or confused.

## Requirements
### Requirement: 简历相关导航标题必须使用稳定且一致的展示文案
系统 MUST 为简历相关侧边栏菜单使用统一的国际化 key 或统一可翻译文案，不得在导航中直接展示 `route.*` 原始 key。

#### Scenario: 用户查看简历仓库菜单
- **WHEN** 用户登录后查看左侧导航中的简历仓库入口
- **THEN** 系统展示正常的“简历仓库”文案
- **AND** 菜单标题不得显示 `route.简历仓库`、缺失 key 或其它原始国际化占位内容

#### Scenario: 用户切换语言或读取菜单标题
- **WHEN** 侧边栏组件根据路由 `meta.title` 渲染菜单标题
- **THEN** 系统使用稳定的标题 key 或可翻译文案进行展示
- **AND** 简历相关菜单与其它菜单采用相同的标题配置策略

### Requirement: 简历仓库入口必须使用贴合业务语义的导航图标
系统 MUST 为简历仓库入口使用与“简历 / 文档资产”语义匹配的图标，并与岗位仓库等相邻导航入口保持视觉风格一致。

#### Scenario: 用户查看简历仓库图标
- **WHEN** 用户查看左侧侧边栏中的简历仓库入口
- **THEN** 系统展示可表达简历、文档或资料仓库语义的图标
- **AND** 图标风格与侧边栏其它业务入口保持一致

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
