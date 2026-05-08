## 背景

本项目是一个基于 Vue 2、TypeScript 和 Element UI 的后台应用，当前通过 Mock API 驱动，已经具备用户信息、简历仓库、简历版本、Markdown 预览和 PDF 导出能力。简历版本是产品中的核心资产，但当前流程在“创建或导出简历”之后就结束了。

本次需求希望以“简历版本”为中心，向下游扩展出完整的求职流程记录能力，让用户能够围绕某个版本继续管理投递记录、面试进展和复盘信息。实现方案需要保持与当前项目架构一致，继续采用 typed model、API 封装、Vuex 模块和基于路由的 CRUD 页面，并先通过 Mock 数据完成闭环。

## 目标 / 非目标

**目标：**

- 支持一个简历版本关联多条投递记录。
- 支持在投递记录中存储结构化的公司和岗位信息。
- 支持为每条投递记录维护按时间排序的面试进展历史。
- 支持在投递详情中维护面试小记、总结和个人评价打分。
- 在不依赖真实后端的前提下，将新能力接入现有简历版本流程。

**非目标：**

- 本次变更不包含跨版本的数据分析大盘。
- 本次变更不包含提醒、日历同步、通知或 AI 自动反馈。
- 本次变更不包含附件能力，例如上传截图、Offer 文件或 PDF 资料。
- 本次变更不追求企业级 ATS 的高度规范化建模，当前优先满足个人求职管理场景。

## 关键设计决策

以简历版本作为投递记录的直接父级。
原因：用户需求明确要求“每个简历版本可以关联多个投递记录”，这样才能准确保留“哪一版简历投给了哪一个岗位”的历史链路。
备选方案：
- 只将投递记录关联到简历仓库：实现更简单，但会丢失版本级可追溯性。
- 允许投递记录脱离简历版本独立存在：更灵活，但不符合本产品以简历版本为核心的流程设计。

将面试进展设计为投递记录下的子集合。
原因：一条投递记录可能经历多个阶段，用户需要看到时间线历史，而不是只看到一个被覆盖后的当前状态。保留“当前状态 + 历史进展”可以同时满足列表展示和详情复盘。
备选方案：
- 只保留一个当前状态字段：实现简单，但会丢失进展历史。
- 将每一轮面试建模为独立顶层实体：扩展性更强，但对当前版本来说过重。

在首版中将面试小记、总结和评分保留在投递详情层级。
原因：用户当前的核心诉求是把一次投递作为一个整体来管理。将时间线、复盘和评分统一放在一个详情页中，录入成本更低，也更符合个人使用习惯。
备选方案：
- 强制每一轮面试都独立填写笔记和评分：粒度更细，但复杂度更高，不适合作为首版落地方案。

首版继续采用 Mock 驱动的 CRUD 流程。
原因：项目当前已经是 Mock-first 架构，这次需求会涉及多个页面和多层数据结构，继续沿用现有模式风险最低，也更容易快速验证流程。
备选方案：
- 等真实后端就绪后再做：会阻塞当前需求推进，不符合项目现有迭代方式。

对存在关联投递记录的简历版本增加删除保护。
原因：投递记录本质上是求职历史，不应因为误删简历版本而丢失链路。
备选方案：
- 删除简历版本时级联删除投递记录：技术上简单，但对用户风险过高。
- 只给出警告但仍允许删除：仍然容易造成历史数据意外丢失。

## 数据模型

### 投递记录

建议结构：

```ts
interface ResumeApplication {
  id: string
  resumeVersionId: string
  repoId: string
  companyName: string
  companyIntro?: string
  jobTitle: string
  jobRequirements?: string
  base?: string
  salaryRange?: string
  deliveryChannel?: string
  appliedAt: string
  currentStatus: ApplicationStatus
  interviewSummary?: string
  interviewNotes?: string
  resumeMatchScore?: number
  interviewPerformanceScore?: number
  roleInterestScore?: number
  overallScore?: number
  remark?: string
  createdAt: string
  updatedAt: string
}
```

### 面试进展记录

建议结构：

```ts
interface InterviewProgress {
  id: string
  applicationId: string
  stage: InterviewStage
  occurredAt: string
  result: InterviewResult
  interviewerOrTeam?: string
  note?: string
  createdAt: string
  updatedAt: string
}
```

### 枚举建议

- `ApplicationStatus`: applied, viewed, written_test, interviewing, offer, rejected, closed
- `InterviewStage`: applied, viewed, written_test, first_interview, second_interview, final_interview, hr_interview, offer, rejected, closed
- `InterviewResult`: passed, pending, failed

其中投递记录保留一个冗余的 `currentStatus` 用于列表快速展示，面试进展列表负责保存完整历史。

## 页面与流程设计

### 简历版本详情页

扩展 `src/views/resume-version/view.vue`，增加“关联投递记录”区域，包含：

- 投递记录数量汇总
- 以表格或卡片形式展示公司名称、岗位名称、base、投递时间和当前状态
- 新增投递记录入口
- 查看、编辑、删除操作入口

这样可以让新流程继续挂在现有的简历版本页面下，符合用户当前查看和导出简历的操作路径。

### 投递记录创建/编辑页

新增表单页，建议分为以下区域：

- 基础岗位信息
- 公司介绍与 JD 信息
- 投递元信息
- 面试总结与笔记
- 个人评分

创建流程中，`resumeVersionId` 应从来源简历版本上下文自动带入，而不是让用户手动选择。

### 投递记录详情页

新增详情页，包含：

- 顶部摘要区：公司名称、岗位名称、当前状态
- 基础信息面板
- 面试进展时间线，支持新增、编辑、删除进展记录
- 面试小记 / 总结面板
- 个人评分面板

该页面将作为用户在投递创建后持续更新求职进展的核心操作页。

### 面试流程节点展示规则

投递记录详情页中的面试进展展示，参考腾讯常见招聘流程图的节点表达方式，以“流程节点”而不是纯文本列表来呈现，建议节点顺序如下：

1. 简历投递
2. 简历评估 / 简历筛选
3. 笔试（如有）
4. 一面
5. 二面
6. 三面 / 终面
7. HR 面
8. Offer

考虑到不同公司流程并不完全一致，首版实现时可允许部分节点缺失，但默认展示结构仍按照上述顺序组织。

节点状态展示规则：

- 已经过的节点：显示绿色，表示该阶段已完成
- 当前进行中的节点：显示绿色高亮，表示当前正在该阶段推进
- 尚未开始的节点：显示默认灰色
- 面试失败节点：显示红色，表示流程在该节点终止

交互和展示建议：

- 在投递详情页使用步骤条、流程图或时间线节点的形式展示，而不是只显示一个状态字段
- 当前状态应同时体现在顶部摘要区和流程节点区，避免用户只在一个位置感知进度
- 当某个节点被标记为失败后，其后的节点不再显示为可进行状态，可保持灰色或直接视为未达成
- 如果流程已进入 Offer，则 Offer 节点显示绿色，其前序已完成节点也全部保持绿色

数据映射建议：

- `currentStatus` 用于列表页快速展示当前整体状态
- `InterviewProgress` 用于驱动详情页的节点状态和时间线内容
- 当新增或更新面试进展时，需要同步计算每个流程节点的展示状态

这样既能满足用户对招聘流程可视化的预期，也能与当前“投递记录 + 面试进展历史”的数据结构保持一致。

## API 与状态管理设计

按照当前项目约定新增以下模块：

- `src/models/resume-application.ts`
- `src/models/interview-progress.ts`
- `src/api/resume-applications.ts`
- `src/api/interview-progress.ts`
- `src/store/modules/resume-application.ts`
- `src/store/modules/interview-progress.ts`
- `mock/resume-applications.ts`
- `mock/interview-progress.ts`

Mock 接口建议覆盖：

- `GET /api/resume-applications?resumeVersionId=...`
- `GET /api/resume-applications/:id`
- `POST /api/resume-applications`
- `PUT /api/resume-applications/:id`
- `DELETE /api/resume-applications/:id`
- `GET /api/interview-progress?applicationId=...`
- `POST /api/interview-progress`
- `PUT /api/interview-progress/:id`
- `DELETE /api/interview-progress/:id`

Store 职责划分：

- 投递记录模块负责列表、详情、加载态和 CRUD
- 面试进展模块负责单条投递下的时间线数据管理

## 风险与权衡

- [该功能会同时影响 model、api、store、mock、route 和 view 多层结构] -> 缓解方式：严格沿用现有 resume 模块的分层模式，并按任务拆分逐步落地。
- [当前状态与进展历史存在重复，可能出现不一致] -> 缓解方式：在新增或编辑面试进展时，同步更新父级投递记录的 `currentStatus`。
- [用户可能更偏好轻量记录，而不是复杂表单] -> 缓解方式：尽量减少必填项，其余字段保持可选。
- [删除保护会改变当前简历版本管理行为] -> 缓解方式：在删除确认和错误提示中明确说明限制原因。

## 迁移方案

当前系统基于内存型 Mock 数据运行，因此不涉及持久化迁移。实现时只需要新增投递记录与面试进展的 mock 数据，并注册到现有 mock server 中，同时让简历版本详情页加载关联投递数据，并在简历版本删除逻辑中增加关联校验。

如果后续需要回滚，可移除新增路由和页面区域、注销对应 mock 接口，并恢复简历版本删除逻辑。现有简历仓库和简历版本能力与本次功能耦合较低，回滚风险可控。

## 待确认问题

- V1 的个人评分是否只保留在投递记录层级，后续是否需要细化到每一轮面试？
- 面试小记在首版中是否只保留为一个可编辑的总结字段，还是从第一版开始就按进展节点分别记录？
- 首版是否只提供“从简历版本详情进入投递记录”的入口，还是需要额外提供一个全局投递列表页？
