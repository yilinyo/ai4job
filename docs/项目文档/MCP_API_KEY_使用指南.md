# MCP 与 API Key 使用指南

## 部署配置

服务端必须通过环境变量配置 `API_KEY_PEPPER`。该值用于计算 API Key 的 HMAC-SHA-256 摘要，必须是随机的高熵秘密，不能提交到仓库，也不能与 JWT 密钥复用。

```env
APP_ENV=production
API_KEY_PEPPER=replace-with-a-long-random-secret
```

本仓库的生产流水线运行在服务器上的 self-hosted GitHub Actions Runner。请在该服务器编辑 `/etc/ai-for/production.env`，部署脚本会以 `docker compose --env-file` 读取并注入 API 容器；不需要把该值配置为 GitHub Actions Secret。首次配置可执行：

```bash
sudo install -o deploy -g deploy -m 600 deploy/production.env.example /etc/ai-for/production.env
sudoedit /etc/ai-for/production.env
```

填写 `API_KEY_PEPPER=<独立随机值>` 后，推送到 `master` 或手动触发 **Verify and deploy production**。请保持该值稳定；更换它会使数据库中已有 Key 的摘要无法匹配，所有已签发 Key 都需要重新创建。

生产环境必须由 HTTPS 或可信反向代理保护服务端连接。MCP 端点为 `https://<host>/mcp`，与既有业务 API 共用后端服务进程和端口；不要通过 Cookie、`X-Access-Token` 或 URL 查询参数传递 API Key。

个人中心的“开发者访问”标签页可创建、轮换和撤销 Key。完整 Key 只在创建或轮换成功时显示一次；应保存到 MCP 客户端的安全凭据存储中。泄露后立即撤销或轮换。

## 客户端配置

MCP 使用 Streamable HTTP 和 Bearer 鉴权：

```toml
[mcp_servers.ai4job]
url = "https://<host>/mcp"

[mcp_servers.ai4job.http_headers]
Authorization = "Bearer ai4job-<random>"
```

本地开发可使用 `http://localhost:8080/mcp`。配置变更后重启 MCP 客户端或其扩展，使其重新初始化会话。

## 工具

所有工具从 API Key 推导调用用户，不接受 `userId`，也不会返回其他用户数据。`page` 默认 1，`pageSize` 默认 20，最大 100。

| 工具 | 参数 | 返回 |
| --- | --- | --- |
| `get_resume_interviews` | `resumeVersionId` 必填；可选 `page`、`pageSize` | 指定简历版本的投递记录和面试进展 |
| `list_interview_workflows` | 可选 `currentStatus`、`page`、`pageSize` | 当前用户全部投递的状态和流程节点 |
| `list_interview_questions` | 可选 `keyword`、`tags`、`difficulty`、`questionType`、`masteryStatus`、`isFavorite`、分页 | 当前用户题库及分页信息 |
| `create_interview_questions` | `questions`，长度 1 至 50 | 创建成功的题目列表；任一条无效则整批不写入 |

题库枚举：

- `difficulty`：`easy`、`medium`、`hard`。
- `questionType`：`knowledge`、`scenario`、`coding`、`algorithm`、`project`、`behavior`、`discussion`。
- `masteryStatus`：`unreviewed`、`weak`、`normal`、`mastered`。

创建题目时每条至少包含 `title`、`content`、`answerAnalysis`、`difficulty`、`questionType` 和 `masteryStatus`。`tags`、`isFavorite`、`remark` 可选；`source` 固定由服务端写为 `mcp`。

## 错误处理

缺少、格式错误、无效、已撤销或已过期的 Key 都返回 HTTP 401，且不会执行工具逻辑。MCP 会话先调用 `initialize`，服务端返回 `Mcp-Session-Id` 后，后续 `tools/list` 与 `tools/call` 必须带上该会话标识。业务校验失败、跨用户资源和超过批量上限的请求不会写入数据。
