## Why

当前仓库仍保持早期前端项目的根目录结构，`src/`、`mock/`、`package.json`、前端构建配置、项目文档和 OpenSpec 文件混放在仓库根目录。后续需要基于现有前端继续设计和实现后端能力，如果不先拆分目录，前端依赖、后端工程、业务文档和 OpenSpec 变更会继续互相挤占根目录，增加维护成本。

本次变更将仓库重构为更适合前后端协作的工作区结构：前端工程集中到 `frontend/`，后端预留到 `backend/`，PRD 文档和项目文档在 `docs/` 下分区存放，OpenSpec 继续保留在根目录。

## What Changes

- 将现有 Vue 2 + TypeScript 前端工程相关文件整体移动到 `frontend/`，包括源码、Mock、测试、构建配置、依赖清单和前端静态资源。
- 在根目录新增 `backend/` 目录并保留占位文件，为后续后端实现预留独立工程位置。
- 将 PRD、需求方案、视觉稿和 PRD 模板归档到 `docs/PRD文档/`。
- 将项目说明、开发指南、历史前端说明和后端 API 设计文档归档到 `docs/项目文档/`。
- 保持 `openspec/`、`.codex/`、`AGENTS.md`、`CLAUDE.md`、CI 配置和仓库级 README 等协作文件在根目录。
- 更新根目录 README、Agent 说明、Claude 说明、CI workflow、CircleCI 配置和 `.gitignore`，确保命令和路径指向 `frontend/`。
- 验证前端移动后仍可执行 lint 和相关单元测试。

## Capabilities

### New Capabilities

- `workspace-structure`: 管理仓库根目录、前端工程、后端占位、文档目录和 OpenSpec 目录的结构约束。

### Modified Capabilities

- 无。

## Impact

- 影响仓库目录结构和所有前端文件路径。
- 影响本地开发命令，需要先进入 `frontend/` 再执行前端命令。
- 影响 CI/CD 配置，需要在 `frontend/` 内安装依赖、执行构建和读取产物。
- 影响 Agent 协作说明和项目文档引用路径。
- 不改变前端业务功能、路由、Mock API 行为或用户界面。
- 不新增真实后端功能，`backend/` 当前仅作为后续实现占位。
