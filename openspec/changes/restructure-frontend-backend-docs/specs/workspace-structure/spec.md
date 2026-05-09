## ADDED Requirements

### Requirement: 仓库必须拆分为前端、后端、文档和 OpenSpec 工作区
系统 MUST 在仓库根目录提供清晰的工作区结构，其中前端工程位于 `frontend/`，后端工程位于 `backend/`，业务文档位于 `docs/`，OpenSpec 文件位于 `openspec/`。

#### Scenario: 开发者查看仓库根目录
- **WHEN** 开发者在仓库根目录查看目录结构
- **THEN** 系统提供 `frontend/`、`backend/`、`docs/` 和 `openspec/` 目录
- **AND** 根目录保留仓库级 README、Agent 说明、CI 配置、OpenSpec 配置和 License 等跨工程文件

### Requirement: 前端工程必须完整位于 frontend 目录
系统 MUST 将现有 Vue 2 前端工程的源码、Mock、测试、静态资源、依赖清单和构建配置集中放置在 `frontend/` 目录内。

#### Scenario: 开发者运行前端本地命令
- **WHEN** 开发者需要安装依赖、启动服务、运行 lint 或执行单元测试
- **THEN** 开发者先进入 `frontend/` 目录
- **AND** 在 `frontend/` 目录内执行前端相关 npm 命令

#### Scenario: CI 构建前端项目
- **WHEN** CI 执行前端依赖安装、构建或部署
- **THEN** CI 在 `frontend/` 目录内读取依赖清单和执行构建命令
- **AND** CI 使用 `frontend/dist/` 作为前端构建产物目录

### Requirement: 后端目录必须作为后续实现占位
系统 MUST 在根目录提供 `backend/` 目录，用于承载后续后端工程，但本次变更不得引入未确认的后端框架或运行时实现。

#### Scenario: 开发者查看后端目录
- **WHEN** 开发者打开 `backend/` 目录
- **THEN** 目录存在并可被 Git 跟踪
- **AND** 当前不包含真实后端业务代码、依赖清单或启动命令

### Requirement: 文档必须集中存放在 docs 目录
系统 MUST 将文档集中存放到 `docs/` 下，并明确区分 PRD 文档和项目文档。

#### Scenario: 开发者查找业务文档
- **WHEN** 开发者需要查找 PRD、后端 API 设计或历史项目说明
- **THEN** 开发者可以在 `docs/PRD文档/` 下查找 PRD、需求方案、视觉稿和 PRD 模板
- **AND** 开发者可以在 `docs/项目文档/` 下查找项目说明、开发指南、项目总结、后端 API 设计和历史前端说明
- **AND** 根目录 README 提供文档位置说明

### Requirement: 协作说明必须反映新的目录结构
系统 MUST 更新仓库级协作说明，使 Agent 和开发者按新的 `frontend/`、`backend/`、`docs/`、`openspec/` 结构工作。

#### Scenario: Agent 根据仓库说明执行前端任务
- **WHEN** Agent 读取仓库协作说明并准备执行前端开发、验证或 Mock API 工作
- **THEN** Agent 使用 `frontend/src/`、`frontend/mock/`、`frontend/tests/` 和 `frontend/package.json` 等新路径
- **AND** Agent 在运行前端命令前进入 `frontend/` 目录
