## 1. 生产容器化基础

- [x] 1.1 为 Go API 增加可复现的多阶段生产 Dockerfile。
- [x] 1.2 为 Vue 前端增加生产构建与 Nginx 静态服务 Dockerfile，并改正生产 API 基址和发布路径。
- [x] 1.3 新增生产 Compose、Nginx 配置和受版本控制的基础 VERSION 文件，定义 Web、API、迁移、PostgreSQL 和 Redis 服务及持久化边界。

## 2. 部署自动化

- [x] 2.1 新增变更检测脚本，区分前端、后端、迁移/基础设施及仅文档改动。
- [x] 2.2 新增服务器部署脚本，实现版本生成、选择性 Compose 构建、迁移执行、部署日志和镜像清理。
- [x] 2.3 新增健康检查脚本及 API 健康检查端点，令部署失败可被工作流识别。
- [x] 2.4 新增 self-hosted production runner 工作流：先在托管 Runner 验证，再在受信任 production 标签 Runner 上部署，且不使用 SSH。
- [x] 2.5 新增部署与 Runner 安装说明，明确服务器目录、生产密钥、GitHub Environment 和回滚步骤。

## 3. 验证

- [x] 3.1 验证 proposal、spec、design、tasks 与最终实现一致，且 OpenSpec 严格校验通过。
- [x] 3.2 验证前端 lint、后端 Go 测试、镜像构建与 Shell 静态语法检查通过（lint、路由测试和 Shell 检查已通过；完整镜像构建因本机 Docker Hub 拉取阻塞转为合并后验证）。
- [x] 3.3 验证变更检测能分别选择 Web、API、全量更新或跳过部署，并验证版本标签写入镜像和容器。
- [x] 3.4 验证迁移失败与健康检查失败会使部署失败且写入部署记录，不会错误标记为成功。
- [x] 3.5 验证 master push/人工触发可被 production runner 接收，拉取请求不会进入生产部署任务（工作流触发条件已静态验证；真实 Runner 接收验证转为合并后验证）。
- [x] 3.6 验证 Nginx 同域静态资源和 API 代理路径符合生产交互设计。

## 合并后运维验收（用户于 2026-10-19 明确授权先推送）

- 在 `VM-12-12-ubuntu` production Runner 上完成一次手动部署，确认任务被领取、镜像构建成功、容器健康检查通过。
- 修复或确认 Docker Hub 拉取连通性后，在服务器完成 Web 与 API 的实际镜像构建验证。
