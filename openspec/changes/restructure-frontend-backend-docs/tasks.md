## 1. 目录重构

- [x] 1.1 创建 `frontend/`、`backend/`、`docs/` 工作区目录
- [x] 1.2 将前端源码、Mock、测试、静态资源、依赖清单和构建配置移动到 `frontend/`
- [x] 1.3 在 `backend/` 中保留占位文件，暂不引入后端框架或业务代码
- [x] 1.4 将 PRD、需求方案、视觉稿和 PRD 模板归档到 `docs/PRD文档/`
- [x] 1.5 将项目说明、开发指南、后端 API 设计和历史前端说明归档到 `docs/项目文档/`
- [x] 1.6 保持 `openspec/`、`.codex/`、`AGENTS.md`、`CLAUDE.md`、CI 配置和仓库级 README 在根目录

## 2. 协作与构建配置

- [x] 2.1 新增根目录工作区 README，说明 `frontend/`、`backend/`、`docs/` 和 `openspec/` 的职责
- [x] 2.2 更新 `AGENTS.md` 和 `CLAUDE.md`，将前端命令和路径改为 `frontend/` 下执行
- [x] 2.3 更新 GitHub Actions 和 CircleCI 配置，使依赖安装、构建和产物路径指向 `frontend/`
- [x] 2.4 更新 `.gitignore`，忽略 `frontend/node_modules/`、`frontend/dist/` 和前端测试覆盖率目录
- [x] 2.5 更新 OpenSpec 验证模板中的前端 lint 和测试命令

## 3. 验证

- [x] 3.1 验证文档、设计、specs 与实际目录结构保持一致
- [x] 3.2 验证静态检查通过：在 `frontend/` 下运行 `npm run lint`
- [x] 3.3 验证相关单元测试通过：在 `frontend/` 下运行面试题库相关单元测试
- [x] 3.4 验证 `frontend/dist/` 与 `frontend/node_modules/` 被忽略且不会进入提交
- [x] 3.5 验证主要协作文件和 CI 配置不再依赖旧的根目录前端路径
