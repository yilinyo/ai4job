# 求职工作台 Dashboard 设计

## 背景

当前 `/dashboard` 仍使用模板自带的 admin/editor dashboard。求职流程系统已经具备简历、岗位、投递、面试进展、题库等核心数据，但缺少一个面向日常求职管理的首页总览。

本设计将 dashboard 改造为求职工作台，提供数据聚合、今日待办、最近动态和投递分析能力。方案采用后端聚合接口，同时提供同结构前端 mock 数据，支持前后端并行开发。

## 目标

- 用求职工作台替换模板 dashboard。
- 支持时间范围筛选：近 7 天、近 30 天、近 90 天、全部，默认近 30 天。
- 展示核心 KPI、Pipeline 状态、投递趋势、公司/岗位/渠道统计。
- 支持系统自动待办和用户手动待办。
- 展示用户操作相关的最近动态，并支持点击跳转详情。
- 后端提供聚合接口，前端 mock server 提供同结构数据。

## 非目标

- 不做待办通知、重复提醒或日历集成。
- 不做用户自定义自动待办规则天数。
- 不做独立动态列表页。
- 不记录岗位创建动态。
- 不记录普通题目列表/详情查看行为。
- 不做多用户协作动态。

## 页面设计

Dashboard 采用分区清晰的工作台布局。

### 顶部区域

- 标题：`求职工作台`
- 时间范围筛选：`近7天 / 近30天 / 近90天 / 全部`
- 默认选择：近 30 天

### KPI 卡片

展示四个核心指标：

- 总投递数
- 面试中
- 待跟进
- Offer

KPI 数据来自后端 dashboard 聚合接口。可保留上一周期对比字段；若后端第一版未计算对比，前端隐藏对比展示。

### 核心进展区

左侧大卡展示 Pipeline 状态与投递趋势，支持图表切换或分上下展示：

- Pipeline 状态：已投递、已查看、笔试、一面、二面、终面、HR 面、Offer、拒绝、关闭。
- 投递趋势：按日或按周展示范围内投递数量。

右侧展示今日待办：

- 合并系统自动待办和手动待办。
- 支持新增、完成、编辑、删除手动待办。
- 待办可关联投递、面试、题目、简历。
- 点击有关联对象的待办时跳转目标详情。

### 统计分析区

展示三张并列小图：

- 公司投递 Top N
- 岗位方向分布
- 渠道来源分布

岗位方向第一版可基于 `jobTitle` 关键词归类；无法归类的进入“其他”。

### 最近动态

最近动态展示用户操作相关事件：

- 简历相关操作
- 题目相关操作
- 投递相关操作
- 面试相关操作

不记录岗位创建。题目动态记录会改变数据或建立关联的行为，包括新增、更新、删除、收藏/取消收藏、关联投递；不记录普通列表/详情查看。

动态以摘要形式展示，例如：`腾讯 / 前端工程师 投递状态更新为一面`。动态项带目标类型、目标 ID 和可选路由信息。可跳转动态点击后进入对应投递、面试、题目或简历详情；删除类动态如果目标详情已不存在，可以不提供跳转或跳转到父级列表。

## 后端设计

### Dashboard 聚合接口

新增接口：

```http
GET /api/dashboard?range=7d|30d|90d|all
```

响应结构：

```ts
interface DashboardResponse {
  range: '7d' | '30d' | '90d' | 'all'
  metrics: DashboardMetrics
  statusPipeline: DashboardStatusBucket[]
  applicationTrend: DashboardTrendPoint[]
  companyRanking: DashboardNameCount[]
  jobTitleDistribution: DashboardNameCount[]
  sourceDistribution: DashboardNameCount[]
  todos: DashboardTodo[]
  activities: ActivityLog[]
}

interface DashboardMetrics {
  totalApplications: number
  interviewing: number
  pendingFollowUp: number
  offers: number
}

interface DashboardStatusBucket {
  status: string
  label: string
  count: number
}

interface DashboardTrendPoint {
  date: string
  count: number
}

interface DashboardNameCount {
  name: string
  count: number
}

interface DashboardTodo {
  id: string
  title: string
  description?: string
  dueAt?: string
  completedAt?: string
  targetType?: 'application' | 'interview' | 'question' | 'resume'
  targetId?: string
  route?: string
  source: 'manual' | 'system'
  createdAt: string
  updatedAt: string
}

interface ActivityLog {
  id: string
  action: string
  targetType: 'application' | 'interview' | 'question' | 'resume'
  targetId: string
  targetLabel: string
  summary: string
  route?: string
  createdAt: string
}
```

The backend should compute all response sections for the authenticated user only.

### 时间范围

`range` 取值：

- `7d`: 最近 7 天
- `30d`: 最近 30 天
- `90d`: 最近 90 天
- `all`: 不按时间过滤

时间范围用于统计投递数量、状态分布、趋势、公司/岗位/渠道分布和最近动态筛选。今日待办始终以当天和未完成待办为主，不受统计范围严格过滤。

### 待办数据模型

新增表 `dashboard_todos`，保存手动待办。系统自动待办查询时生成，不落库。

字段：

- `id UUID PRIMARY KEY`
- `user_id VARCHAR NOT NULL`
- `title VARCHAR NOT NULL`
- `description TEXT`
- `due_at TIMESTAMP`
- `completed_at TIMESTAMP`
- `target_type VARCHAR`
- `target_id UUID`
- `created_at TIMESTAMP NOT NULL`
- `updated_at TIMESTAMP NOT NULL`

手动待办的 source 在 DTO 层返回为 `manual`。系统待办的 source 返回为 `system`。

手动待办 API：

```http
GET /api/dashboard/todos?completed=false
POST /api/dashboard/todos
PUT /api/dashboard/todos/:id
PUT /api/dashboard/todos/:id/complete
DELETE /api/dashboard/todos/:id
```

手动待办支持：

- 新增
- 编辑
- 删除
- 完成
- 截止时间
- 关联投递、面试、题目、简历

### 自动待办规则

第一版包含两类自动待办：

1. 超时未跟进
   - 投递后 7 天仍未进入下一阶段时生成跟进待办。
   - 目标关联到投递记录。
2. 面试当天
   - 当天存在面试进展或面试安排时生成准备/参加面试待办。
   - 目标关联到投递或面试记录。

自动待办不持久化，因此无法编辑或删除。前端可以展示为系统来源，并只允许跳转，不提供编辑删除。

### 动态日志数据模型

新增表 `activity_logs`，记录用户操作相关事件。

字段：

- `id UUID PRIMARY KEY`
- `user_id VARCHAR NOT NULL`
- `actor_type VARCHAR NOT NULL DEFAULT 'user'`
- `action VARCHAR NOT NULL`
- `target_type VARCHAR NOT NULL`
- `target_id UUID NOT NULL`
- `target_label VARCHAR NOT NULL`
- `summary TEXT NOT NULL`
- `created_at TIMESTAMP NOT NULL`

记录范围：

- 简历：创建、更新、删除、激活版本等用户操作。
- 题目：新增、更新、删除、收藏/取消收藏、关联投递。
- 投递：创建、更新、删除、状态变更。
- 面试：创建、更新、删除、阶段/结果更新。

不记录：

- 岗位创建。
- 普通查看题目列表或详情。

Dashboard 聚合接口返回最近 N 条动态，第一版 N 可固定为 10。

### 跳转路由

后端可以在 DTO 中返回可选 `route`，前端点击时优先使用。删除类动态如果无法跳到详情，后端可省略 `route`。建议路由：

- 投递：`/resume/application/:id`
- 面试：若无独立面试详情页，跳转到关联投递详情 `/resume/application/:applicationId`
- 题目：以现有题库详情路由为准
- 简历：以简历仓库或版本详情路由为准

## 前端设计

### 文件结构

新增或改造：

- `frontend/src/views/dashboard/index.vue`
- `frontend/src/views/dashboard/components/MetricCards.vue`
- `frontend/src/views/dashboard/components/PipelineChart.vue`
- `frontend/src/views/dashboard/components/TodoPanel.vue`
- `frontend/src/views/dashboard/components/DistributionCharts.vue`
- `frontend/src/views/dashboard/components/ActivityFeed.vue`
- `frontend/src/models/dashboard.ts`
- `frontend/src/api/dashboard.ts`
- `frontend/src/store/modules/dashboard.ts`

现有模板 dashboard admin/editor 组件不再作为首页内容使用。

### 页面数据流

1. 页面默认 `range = '30d'`。
2. 调用 `DashboardModule.GetDashboard({ range })`。
3. 渲染 KPI、图表、待办、最近动态。
4. 用户切换时间范围后重新请求 dashboard。
5. 用户新增、编辑、完成、删除手动待办后刷新 dashboard 或刷新 todo 列表。

### 图表

复用当前项目已有 ECharts 依赖，按求职业务重新封装图表组件：

- 状态 Pipeline：柱状图。
- 投递趋势：折线图。
- 公司 Top N：横向柱状图。
- 岗位方向分布：环形图或饼图。
- 渠道来源分布：环形图或饼图。

## Mock 设计

新增：

- `frontend/mock/dashboard.ts`

在 `frontend/mock/mock-server.ts` 注册：

```ts
app.get('/api/dashboard', dashboard.getDashboard)
app.get('/api/dashboard/todos', dashboard.getDashboardTodos)
app.post('/api/dashboard/todos', dashboard.createDashboardTodo)
app.put('/api/dashboard/todos/:id', dashboard.updateDashboardTodo)
app.put('/api/dashboard/todos/:id/complete', dashboard.completeDashboardTodo)
app.delete('/api/dashboard/todos/:id', dashboard.deleteDashboardTodo)
```

Mock 响应结构必须与后端真实接口一致，支持前端本地开发。Mock 数据应覆盖 KPI、Pipeline、趋势、三类统计图、系统待办、手动待办、最近动态。

## 测试策略

### 后端测试

- Dashboard 聚合：
  - `range` 过滤正确。
  - KPI 统计正确。
  - Pipeline 状态统计正确。
  - 投递趋势统计正确。
  - 公司、岗位、渠道聚合正确。
  - 自动待办规则正确。
- 手动待办：
  - 创建、编辑、完成、删除。
  - 只能访问当前用户数据。
- 动态日志：
  - 简历、题目、投递、面试变更写入日志。
  - 岗位创建不写入日志。
  - 普通查看题目列表/详情不写入日志。

### 前端测试和验证

- Dashboard store/API 能消费 mock 数据。
- 页面能渲染 KPI、核心进展、今日待办、统计分析和最近动态。
- 切换 range 会重新加载数据。
- 手动待办新增、编辑、完成、删除可用。
- 动态和待办点击能跳转到目标详情。
- 手动启动前端 + mock server 验证 `/dashboard`。

## 实施顺序

1. 后端新增 migration：`dashboard_todos`、`activity_logs`。
2. 后端新增 dashboard DTO、repository、service、handler、router。
3. 后端接入 activity log 写入点：简历、题目、投递、面试。
4. 前端新增 dashboard model、api、store。
5. 前端新增 mock dashboard 数据与 todo API。
6. 重写 dashboard 页面与组件。
7. 补充测试并手动验证。