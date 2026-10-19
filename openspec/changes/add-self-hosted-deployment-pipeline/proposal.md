## Why

当前项目只能在本地分别启动前端、后端及基础依赖，缺少可重复、可审计的生产部署路径。需要让生产服务器上的 GitHub Actions self-hosted runner 主动通过 HTTPS 获取已验证的部署任务，并按改动范围更新服务、暴露明确版本并保留可追溯的部署结果。

## What Changes

- 增加面向单台 Linux 服务器的 Docker Compose 生产部署拓扑，运行 Web、API、PostgreSQL、Redis 和一次性迁移任务。
- 为前端和 Go 后端增加可构建的生产镜像，并由 Nginx 提供静态资源、同域代理 API 请求。
- 增加由 self-hosted production runner 执行的 GitHub Actions 工作流；验证通过后，根据前端、后端或基础设施改动选择性部署。
- 增加 Shell 部署脚本，负责版本生成、服务重建、数据库迁移、健康检查和部署结果记录。
- 为镜像和容器增加发布版本、Git SHA 等可查询元数据；不再通过 SSH 从 GitHub Actions 连接服务器。

## Capabilities

### New Capabilities

- `self-hosted-deployment-pipeline`: 定义基于服务器 self-hosted runner 的选择性 Docker 部署、版本标记、迁移与健康检查行为。

### Modified Capabilities

- 无。

## Impact

- 新增前后端 Dockerfile、生产 Compose 与 Nginx 配置。
- 新增 GitHub Actions 工作流及部署 Shell 脚本。
- 调整前端生产 API 基址与资源发布路径。
- 生产服务器需安装 Docker、Compose 与 self-hosted runner，并在本机保存生产环境密钥和持久化数据卷。
