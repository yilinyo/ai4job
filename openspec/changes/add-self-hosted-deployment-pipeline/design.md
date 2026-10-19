## Context

当前仓库可在本地运行 Vue 前端、Go API、PostgreSQL 和 Redis，但仅有基础依赖的 Compose 文件，且生产前端配置仍指向模板 Mock 服务。参见 proposal.md 了解此次变更动机。

## Goals / Non-Goals

**Goals:**

- 使用同域名 Nginx 代理提供生产 Web 和 `/api/` 请求，避免浏览器跨域配置。
- 由服务器本机的 self-hosted runner 主动通过 HTTPS 领取部署任务，无需 GitHub 对服务器开放 SSH。
- 按改动范围构建 Web 或 API，并给镜像与容器注入可查询的不可变版本。
- 迁移、健康检查与部署历史作为发布闭环的一部分。

**Non-Goals:**

- 不提供多服务器编排、自动扩缩容或零停机蓝绿发布。
- 不提供用户界面中的部署历史或一键回滚页面。
- 不将生产密钥、证书或数据库数据纳入 Git。

## Decisions

### 单机 Docker Compose 作为生产编排

新增 `web`、`api`、`migrate`、`postgres`、`redis` 服务。`web` 通过多阶段构建生成 Vue 静态资源，并以 Nginx 对外提供页面和 API 反向代理；数据库和缓存仅加入内部网络及持久化卷。

选择单机 Compose 而非 Kubernetes：当前目标为一台服务器，Compose 的复杂度、资源占用和运维门槛更低。

### 生产 Runner 仅执行受信任分支部署

工作流分为 GitHub 托管 Runner 上的验证 job 和标签为 `production` 的 self-hosted 部署 job。部署 job 只接受 `main` push 或人工触发，并使用 GitHub Environment 控制生产审批。Runner 以专用 deploy 用户运行，通过 Docker 执行部署。

选择 self-hosted runner 而非 SSH action：服务器保持仅出站 HTTPS 连接，减少入站凭据和 SSH Secrets 的维护。

### Shell 脚本承载部署决策

工作流 YAML 只负责编排；`scripts/ci/detect-changes.sh` 把路径差异写入 GitHub outputs，`deploy/scripts/deploy.sh` 负责版本、构建、迁移和服务更新，`deploy/scripts/healthcheck.sh` 负责轮询服务。脚本必须使用严格错误处理并接受显式参数。

### 使用镜像标签和容器标签记录版本

根目录 `VERSION` 提供人工维护的基础版本。流水线生成 `v<base>-run.<run>-<short_sha>`，作为 `web`、`api` 镜像 tag，并写入 `com.ai-for.version`、`com.ai-for.git-sha` 标签。容器名保持稳定，版本通过 `docker ps` 的镜像列和 `docker inspect` 查询。

选择镜像标签与 Label 而非版本化容器名：稳定容器名便于 Nginx、日志、巡检和运维命令引用。

### 迁移优先于 API 更新

后端或迁移文件变动时，先运行一次性 `migrate` 服务，再更新 API。迁移必须向后兼容；若迁移失败，API 不更新。成功与失败部署记录写入服务器挂载的持久日志文件，避免部署记录依赖待部署 API。

## Risks / Trade-offs

- [self-hosted runner 执行不受信任代码可危及服务器] → 仅私有仓库、仅受信任的 main/workflow_dispatch 事件使用 production 标签，生产 Environment 需审批。
- [迁移后新 API 健康检查失败] → 记录失败版本；本次范围不承诺自动回滚数据库，迁移必须采用向后兼容策略。
- [服务器构建消耗 CPU、磁盘] → 选择性构建，设置 Docker 镜像保留与清理策略，并在部署前检查磁盘空间。
- [Runner 工作区含旧构建产物] → 每次部署 checkout 精确 Git SHA，构建使用版本化镜像标签，不依赖 `latest`。

## Migration Plan

1. 在服务器安装 Docker Compose 与 self-hosted runner，并配置仅出站 HTTPS 和生产 Environment。
2. 将生产环境变量、TLS 证书和数据卷放置于服务器受限路径。
3. 合并部署文件后由人工触发首次部署，初始化数据卷和迁移。
4. 验证容器版本、部署日志、API 与页面；后续 main 推送走自动选择性部署。
5. 回滚时使用上一个成功版本的镜像 tag 重新启动对应服务；数据库回滚按独立迁移策略执行。
