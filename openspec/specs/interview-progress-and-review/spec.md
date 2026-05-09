# interview-progress-and-review Specification

## Purpose
Define interview progress tracking, timeline visualization, notes, summaries, and subjective review scores for each resume application.

## Requirements
### Requirement: 投递记录必须支持面试进展历史
系统 MUST 允许用户围绕每条投递记录维护多个面试进展节点，并保留阶段、时间、结果和备注。

#### Scenario: 用户新增面试进展
- **WHEN** 用户在投递记录详情页新增面试进展
- **THEN** 系统保存该进展的阶段、发生时间、结果和备注
- **AND** 进展归属于当前投递记录

### Requirement: 当前状态必须与面试进展同步
系统 MUST 在新增或更新面试进展时同步投递记录的当前状态。

#### Scenario: 面试进展更新当前状态
- **WHEN** 用户新增或编辑一个面试阶段进展
- **THEN** 系统根据进展阶段更新投递记录 `currentStatus`

### Requirement: 面试流程节点必须按状态可视化
系统 MUST 在投递详情页以流程节点展示简历投递、简历评估、笔试、一面、二面、终面、HR 面和 Offer 等阶段状态。

#### Scenario: 用户查看面试流程节点
- **WHEN** 用户打开投递详情页
- **THEN** 已完成节点显示绿色
- **AND** 当前节点绿色高亮
- **AND** 失败节点显示红色
- **AND** 未开始节点显示灰色

### Requirement: 投递详情必须支持复盘与评分
系统 MUST 支持在投递详情中维护面试小记、总结和个人评分字段。

#### Scenario: 用户填写复盘信息
- **WHEN** 用户编辑投递详情
- **THEN** 系统允许保存面试小记、总结、简历匹配度、面试表现、岗位意向度和综合评价
