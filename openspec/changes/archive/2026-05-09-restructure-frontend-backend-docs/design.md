## Context

仓库当前已经从单纯的前端模板项目演进为求职流程管理系统，并开始沉淀后端 API 设计文档。为了支持后续后端开发，需要先将仓库从“前端项目根目录”调整为“前后端工作区根目录”。

重构后的目录约定为：

- `frontend/`：现有 Vue 2、TypeScript、Element UI、Vuex、Vue Router、Express Mock API 前端工程。
- `backend/`：后端工程预留目录，当前仅保留占位文件。
- `docs/`：文档根目录，下分 `docs/PRD文档/` 和 `docs/项目文档/`。
- `openspec/`：OpenSpec 变更与规格，继续位于根目录。
- 根目录：仓库级协作配置、CI 配置、Agent 说明、README、License 等。

## Goals / Non-Goals

**Goals:**

- 将前端工程文件完整移动到 `frontend/`，并保持前端 lint、测试和构建入口可用。
- 创建 `backend/` 占位目录，为后续后端工程提供稳定位置。
- 将 PRD 和项目文档分别归档到 `docs/PRD文档/` 与 `docs/项目文档/`，并保留后端 API 设计文档。
- 更新根目录说明、Agent 协作说明和 CI 配置，避免继续引用旧的根目录前端路径。
- 保持 OpenSpec 文件在根目录，继续作为跨前后端的需求与实现跟踪机制。

**Non-Goals:**

- 不实现后端服务、数据库、认证或真实 API。
- 不修改前端业务功能、视觉样式、Mock 数据语义或用户流程。
- 不引入 monorepo 包管理工具或 workspace 依赖管理。
- 不归档既有 OpenSpec change。

## Decisions

### 1. 前端工程整体移动到 `frontend/`

现有 `src/`、`mock/`、`public/`、`tests/`、`package.json`、`package-lock.json`、`vue.config.js`、`tsconfig.json`、Jest/Babel/PostCSS/ESLint 配置等均归属前端工程，统一移动到 `frontend/`。

理由：

- 后续后端工程会拥有自己的源码、依赖和启动命令，继续放在根目录会混淆工程边界。
- 前端自身的相对路径在整体搬迁后大多仍保持一致，迁移风险低。
- 本地开发命令只需切换为 `cd frontend` 后执行。

### 2. 根目录保留跨工程协作文件

`AGENTS.md`、`CLAUDE.md`、`.codex/`、`.github/`、`.circleci/`、`openspec/`、`README.md`、`.gitignore` 和 `LICENSE` 保留在根目录。

理由：

- 这些文件约束的是整个仓库，而不是单一前端工程。
- OpenSpec 需要继续跨前端、后端和文档工作流。
- CI 和 Agent 说明需要从根目录协调不同子工程。

### 3. `backend/` 当前只保留占位

本次只创建 `backend/.gitkeep`，不提前引入后端框架、依赖或目录模板。

理由：

- 用户当前目标是仓库结构重构，不是后端实现。
- 后端技术栈和接口实现仍需基于后续 API 文档和需求确认。
- 保持占位可以让目录结构立即稳定，又避免提前做错误架构选择。

### 4. 文档在 `docs/` 下区分 PRD 文档和项目文档

`docs/PRD文档/` 用于保存 PRD、需求方案、视觉稿和 PRD 模板；`docs/项目文档/` 用于保存项目说明、开发指南、项目总结、后端 API 设计和历史前端说明。历史前端文档保留迁移说明，提示前端工程已移动到 `frontend/`。

理由：

- 文档不应散落在根目录，尤其是在后续前后端并行开发后。
- PRD 和项目文档的生命周期不同，分开后更容易查找、评审和归档。
- 历史文档仍有参考价值，不应直接删除。
- 新根目录 README 只保留工作区级入口信息。

### 5. CI 和忽略规则显式适配 `frontend/`

GitHub Actions、CircleCI 和 `.gitignore` 均显式引用 `frontend/` 下的依赖、构建产物和测试覆盖率目录。

理由：

- 避免 CI 继续在根目录查找 `package.json`。
- 避免 `frontend/dist/` 和 `frontend/node_modules/` 被误提交。
- 保持后续新增后端工程时规则更清晰。

## Risks / Trade-offs

- [路径引用遗漏] -> 更新根目录协作文件、CI 文件和主要项目文档，并通过 `rg` 检查常见旧路径引用。
- [CI 命令不兼容] -> 将 GitHub Actions 的前端命令改为在 `frontend/` 内执行，将 CircleCI 的工作目录切到 `frontend/`。
- [历史文档仍包含旧命令] -> 保留迁移说明，避免强行重写所有历史内容造成文档语义漂移。
- [后端目录过空] -> 使用 `.gitkeep` 固化目录，后续在专门 change 中实现后端结构。

## Migration Plan

1. 创建 `frontend/` 和 `backend/` 目录。
2. 将前端工程文件移动到 `frontend/`。
3. 将根目录 PRD、需求方案、视觉稿和 PRD 模板移动到 `docs/PRD文档/`。
4. 将项目说明、开发指南、后端 API 设计和历史前端说明移动到 `docs/项目文档/`。
5. 新增工作区级 README，说明目录结构和前端命令。
6. 更新 Agent 说明、Claude 说明、CI 配置和忽略规则。
7. 运行前端 lint 和相关单元测试确认移动后仍可执行。

## Open Questions

- 后端工程具体技术栈、目录分层、运行命令和 API 实现将在后续后端设计或实现 change 中确定。
