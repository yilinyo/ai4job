# 生产部署指南

本项目使用一台服务器上的 Docker Compose 运行 `web`、`api`、PostgreSQL 和 Redis。GitHub Actions Runner 由服务器主动通过 HTTPS 连接 GitHub 并领取任务；GitHub 不会通过 SSH 连接服务器。

## 一次性服务器准备

以专用的非 root 用户（示例为 `deploy`）安装 Docker Engine、Docker Compose plugin、Git、curl 和 GitHub Actions Runner。将该用户加入 `docker` 组后重新登录：

```bash
sudo useradd --create-home --shell /bin/bash deploy
sudo usermod -aG docker deploy
sudo mkdir -p /srv/ai-for /etc/ai-for /var/lib/ai-for
sudo chown -R deploy:deploy /srv/ai-for /var/lib/ai-for
sudo chmod 700 /etc/ai-for
```

在 GitHub 仓库的 **Settings → Actions → Runners** 新建 Linux self-hosted runner，使用页面提供的下载与注册命令在 `/srv/actions-runner` 中注册。注册标签必须包含：

```text
self-hosted, linux, x64, production
```

将 Runner 作为 `deploy` 用户的 systemd 服务运行。Runner 只需访问 GitHub 的 HTTPS（443）出站网络；不需要开放 SSH 入站端口。生产 runner 只能用于此私有仓库受信任的 `master` 推送和人工触发任务。

## 生产变量与首次部署

创建服务器本地密钥文件，禁止将它提交到仓库：

```bash
sudo install -o deploy -g deploy -m 600 deploy/production.env.example /etc/ai-for/production.env
sudoedit /etc/ai-for/production.env
```

必须替换 `DB_PASSWORD`、`REDIS_PASSWORD`、`JWT_SECRET`。密码若包含 URL 保留字符（尤其 `@`、`:`、`/`、`?`、`#`），请进行 URL 编码，因为迁移命令使用 PostgreSQL URL。

创建 GitHub 的 `production` Environment，并按需要配置必需审批者。首次上线前，从 Actions 页面手动运行 **Verify and deploy production**；它会构建镜像、启动依赖、执行迁移并检查服务健康。

生产服务器上的持久化边界：

- PostgreSQL：Compose 卷 `postgres_data`
- Redis：Compose 卷 `redis_data`
- 部署历史：`/var/lib/ai-for/deployments.jsonl`
- 生产密钥：`/etc/ai-for/production.env`

## 发布行为与可观测性

每次发布的镜像版本为 `v<VERSION>-run.<workflow-run>-<git-short-sha>`。可通过下列命令检查版本和部署历史：

```bash
docker compose --env-file /etc/ai-for/production.env -f deploy/docker-compose.prod.yml ps
docker image ls 'ai-for-*'
docker inspect --format '{{ index .Config.Labels "com.ai-for.version" }} {{ index .Config.Labels "com.ai-for.git-sha" }}' "$(docker compose --env-file /etc/ai-for/production.env -f deploy/docker-compose.prod.yml ps -q api)"
tail -n 20 /var/lib/ai-for/deployments.jsonl
```

仅修改 `frontend/` 时仅更新 Web，不会重启 API、PostgreSQL 或 Redis；修改 `backend/`（包含迁移）时先迁移再更新 API；修改 `deploy/`、构建文件、部署脚本或工作流时更新 Web 和 API。仅文档改动不部署。

## 回滚

先从部署记录或 `docker image ls` 找到一个已验证版本。通过临时覆盖 `APP_VERSION` 启动相应服务：

```bash
APP_VERSION=v1.0.0-run.42-abcdef0 \
docker compose --env-file /etc/ai-for/production.env -f deploy/docker-compose.prod.yml up -d --no-build web api
```

回滚应用镜像不会回滚数据库。生产迁移必须向后兼容；数据库回滚需要独立的、经审核的迁移方案。

## 故障处理

部署失败时，工作流会失败并在部署历史写入 `failed`。检查服务和日志：

```bash
docker compose --env-file /etc/ai-for/production.env -f deploy/docker-compose.prod.yml ps
docker compose --env-file /etc/ai-for/production.env -f deploy/docker-compose.prod.yml logs --tail=200 api web postgres redis
```

不要使用 `docker compose down -v` 处理生产故障，它会删除数据库和 Redis 持久化卷。
