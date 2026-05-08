# AGENTS.md

本文件为本仓库中的通用 Agent 协作说明，适用于在 `/Users/yilin/project/ai-for` 中工作的编码代理。

## 项目概览

- 本项目基于 `vue-typescript-admin-template`
- 当前已扩展为一个求职流程管理前端系统
- 主要业务包括：
  - 简历仓库管理
  - 简历版本管理
  - 投递记录管理
  - 面试进展管理
  - 岗位库管理

技术栈：

- Vue 2
- TypeScript
- Element UI
- Vuex + `vuex-module-decorators`
- Vue Router
- Express Mock API

## 常用命令

安装依赖：

```bash
npm install --legacy-peer-deps
```

启动前端 + Mock：

```bash
npm run serve
```

仅启动 Mock：

```bash
npm run mock
```

静态检查：

```bash
npm run lint
```

单元测试：

```bash
npm run test:unit
```

## 运行与端口

- 前端开发端口：`9527`
- Mock API 端口：`9528`
- 开发环境 API 前缀：`/dev-api`
- 代理配置见：
  - [vue.config.js](/Users/yilin/project/ai-for/vue.config.js)
  - [mock/mock-server.ts](/Users/yilin/project/ai-for/mock/mock-server.ts)

## 代码分层约定

新增业务模块时，优先遵循以下分层：

1. `src/models/`
2. `src/api/`
3. `mock/`
4. `src/store/modules/`
5. `src/router/modules/`
6. `src/views/`
7. `src/components/`

不要只改页面而跳过 model / api / store 层。

## 当前重要业务约定

### 简历与投递

- 一个简历版本可关联多条投递记录
- 投递记录保留岗位快照字段
- 已有关联投递记录的简历版本不可直接删除

### 面试进展

- 面试进展挂在投递记录下
- `currentStatus` 需要与面试进展同步
- 节点显示规则：
  - 已完成：绿色
  - 当前节点：绿色高亮
  - 失败节点：红色
  - 未开始：灰色

### 岗位库

- 岗位库是独立一级模块，不挂在“简历管理”二级菜单下
- 岗位支持手动录入和链接解析导入
- 前端只消费解析接口返回的 JSON，不负责抓取解析逻辑
- 岗位卡片尺寸、留白和简历列表风格保持一致

## OpenSpec 工作约定

本仓库使用 OpenSpec 做需求与实现跟踪。

### 基本规则

- 新需求优先通过 OpenSpec change 管理
- 中途补充需求时，必须同步更新对应的：
  - `proposal.md`
  - `design.md`
  - `tasks.md`
- 不允许只改代码、不改 OpenSpec 文档

### 什么时候需要提出一个 change

以下情况应优先提出新的 OpenSpec change：

1. 提出一个全新的模块
2. 提出一个已有模块的新功能
3. 提出某个功能的 bug 并修复
4. 删除某个单独功能

约束：

- 拒绝一次性在一个 change 中实现多个一级模块
- 拒绝一次性在一个 change 中删除多个功能
- 拒绝通过一个 change 删除整个模块
- 删除功能时必须分析影响范围，并说明可能受影响的相关功能

### change 粒度约定

一个 change 应只处理一件相对独立的事情。

推荐粒度：

- 一个新模块
- 一个已有模块中的一个新功能
- 一个明确的 bug 修复
- 一个单独功能删除

禁止把以下内容混在同一个 change 中：

- 新模块 + 多个旧模块联动重构
- 新功能 + 无关 bug 修复
- 删除功能 + 大范围结构改造

### 文档要求

每个 change 默认应有需求文档支撑。

建议分级：

- 轻量变更：
  - 小 bug
  - 小交互修复
  - 文案修正
  - 可使用简版需求说明

- 标准变更：
  - 已有模块新功能
  - 页面结构调整
  - 流程调整
  - 应有 PRD

- 模块级变更：
  - 新模块
  - 应有完整 PRD
  - 必要时补视觉稿

### OpenSpec 与 Git 分支工作流

本项目采用“文档在 `master`，代码在功能分支”的工作模式。

规则如下：

1. 提出 change、撰写 PRD、维护 OpenSpec 文档：
   - 默认在 `master` 分支进行

2. 进入实现阶段（apply）前：
   - 必须从最新 `master` 切出新分支
   - 禁止直接在 `master` 上改代码

3. 分支命名建议：
   - 新功能：`feat/<change-name>`
   - bug 修复：`fix/<change-name>`
   - 删除功能：`chore/remove-<feature-name>`
   - 纯文档：`docs/<change-name>`

4. apply 阶段要求：
   - 如果当前分支是 `master`，不得直接开始实现
   - 必须先创建并切换到新分支
   - 当当前分支不符合当前操作权限时，必须先明确检查分支状态，并向用户说明，再决定下一步操作
   - 不允许在分支状态未确认的情况下直接继续执行 change apply、代码修改、归档前提交等操作
   - 如果当前 change 的 `tasks.md` 仍有未完成项，不允许提交实现代码，也不允许合并回 `master`

5. 提交与合并要求：
   - 在功能分支提交前，必须确认对应 OpenSpec change 的 `tasks.md` 已全部完成
   - 在功能分支合并回 `master` 前，必须再次确认 `tasks.md` 已无未完成项
   - 如果存在未完成任务，只能继续实现或补充验证，不得通过提交或合并绕过任务闭环

6. 归档前要求：
   - 当前分支必须是 `master`
   - OpenSpec tasks 完成
   - 验证完成
   - 代码已提交
   - 功能分支已合并回 `master`
   - 归档动作完成后，必须立即提交 archive 变更并推送到远程 `master`

也就是说，本项目中的 OpenSpec archive 不仅是文档动作，也是 Git 流程闭环动作。

### 分支状态检查规则

以下场景必须先检查当前 Git 分支状态，再决定是否继续：

1. 准备 apply 一个 change
2. 准备开始代码实现
3. 准备执行提交、合并、归档闭环
4. 用户要求直接在当前仓库继续某个 OpenSpec 工作流

执行要求：

- 如果当前分支不满足当前操作权限，必须先明确指出
- 必须显式提示用户当前分支状态与预期流程不一致
- 在未完成分支确认前，不能直接继续后续操作

典型示例：

- 在 `master` 上尝试 apply 改代码：必须先提示并要求切新分支
- 在功能分支上尝试提出应当落在 `master` 的正式 change 文档：必须先提示并确认是否回到 `master`
- 在未确认分支是否已合并回 `master` 的情况下尝试 archive：必须先检查再继续
- 在非 `master` 分支尝试 archive：必须停止，并先切回 `master`
- archive 完成但尚未提交或推送：不得视为流程闭环完成

### 删除功能的额外要求

删除功能必须单独开 change，并且至少补充以下内容：

- 删除原因
- 受影响页面
- 受影响接口
- 受影响数据结构
- 是否影响已有用户数据
- 是否有回滚方案

如果删除带来替代路径变化，必须在 `design.md` 中写清楚。

### 项目级 skills

以下 skills 已在项目中沉淀，应优先遵循：

- [openspec-change-sync](/Users/yilin/project/ai-for/.codex/skills/openspec-change-sync/SKILL.md)
  - 用于保证需求中途补充时及时回写 OpenSpec 文档

- [openspec-validation-template](/Users/yilin/project/ai-for/.codex/skills/openspec-validation-template/SKILL.md)
  - 用于统一每个 OpenSpec change 的验证逻辑与任务写法

- [openspec-branch-workflow](/Users/yilin/project/ai-for/.codex/skills/openspec-branch-workflow/SKILL.md)
  - 用于约束 change 粒度、分支创建、apply 时机与 archive 前 Git 闭环

### 验证默认要求

每个 OpenSpec change 默认应覆盖：

1. 文档一致性验证
2. 静态检查验证
3. 数据 / API 验证
4. 主用户流程验证
5. 约束 / 异常验证
6. 视觉 / 交互一致性验证

如果没有真实验证，不要把 `tasks.md` 对应项勾成完成。

### Archive 前附加检查

除了验证完成外，还要确认：

1. 中途补充需求已同步进 OpenSpec 文档
2. 当前 change 没有 scope 漂移
3. 功能分支不是悬空状态
4. change 对应代码已准备合并回 `master`

## 文档存放约定

业务 PRD 和视觉稿统一归档到 `docs/` 下，按主题拆目录。

当前已有：

- [docs/岗位库方案](/Users/yilin/project/ai-for/docs/岗位库方案)
- [docs/投递与面试方案](/Users/yilin/project/ai-for/docs/投递与面试方案)

如果新增同类文档，优先归到 `docs/`，不要长期散落在根目录。

## Git 约定

当前远程仓库：

```bash
origin https://github.com/yilinyo/ai-for-frontend.git
```

不要随意修改 remote，除非用户明确要求。

## 修改前检查

在执行较大改动前，建议先看：

- [CLAUDE.md](/Users/yilin/project/ai-for/CLAUDE.md)
- [PROJECT_SUMMARY.md](/Users/yilin/project/ai-for/PROJECT_SUMMARY.md)
- [guide_dev.md](/Users/yilin/project/ai-for/guide_dev.md)

## 完成前检查

结束一个较完整任务前，至少确认：

- 代码层是否跑通
- lint 是否通过
- Mock / API 是否符合预期
- OpenSpec 文档是否同步
- 文档是否需要归档到 `docs/`
