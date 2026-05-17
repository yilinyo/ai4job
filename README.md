# AI For Job Frontend/Backend Workspace

本仓库已重构为前后端分层工作区：

```text
.
├── frontend/   # Vue 2 + TypeScript 前端工程，包含 Mock 服务
├── backend/    # Go 后端工程
├── docs/       # PRD 文档与项目文档
├── openspec/   # OpenSpec changes/specs/archive
└── .codex/     # 项目级 Agent skills
```

## 前端开发

```bash
cd frontend
npm install --legacy-peer-deps
npm run serve
```

默认 `npm run serve` 会同时启动前端和 Mock API。请求真实后端时，请先确保后端运行在 `127.0.0.1:8080`，再使用：

```bash
cd frontend
VUE_APP_USE_MOCK=false npx vue-cli-service serve
```

常用命令：

```bash
cd frontend
npm run mock
npm run lint
npm run test:unit
npm run build:prod
```

端口：

- 前端开发服务：`9527`
- Mock API：`9528`
- 真实后端 API：`8080`
- 开发 API 前缀：`/dev-api`

## 后端开发

后端工程位于 `backend/`。API 设计可参考：

- [后端 API 设计](docs/项目文档/后端API设计/API_求职流程管理系统_GET_POST_JSON.md)

## 文档

- `docs/PRD文档/`：PRD、需求方案、视觉稿和 PRD 模板。
- `docs/项目文档/`：项目说明、开发指南、项目总结、后端 API 设计和历史前端说明。

## OpenSpec

OpenSpec 仍位于根目录 `openspec/`，不放入 `frontend/` 或 `backend/`。
