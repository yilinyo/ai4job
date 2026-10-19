# PRD：个人 API Key 与 MCP 服务

## 1. 文档信息

| 项目 | 内容 |
| --- | --- |
| 功能名称 | 个人 API Key 与 MCP 服务 |
| 目标用户 | 已登录的个人求职用户、使用 Codex 等 MCP 客户端的用户 |
| 版本 | V1.0 |
| 状态 | 已实现 |
| 关联 OpenSpec | `add-mcp-api-key-auth` |

## 2. 背景与目标

用户的简历投递、面试流程和题库数据此前只能在网页 JWT 会话中访问，无法安全地被个人 MCP 客户端使用。用户不应向外部客户端提供网页登录凭据，也不能由客户端指定其他用户身份。

本功能提供用户自助申请、轮换和撤销的个人 API Key，并在现有 Go 后端进程中提供 MCP Streamable HTTP 服务。每次 MCP 调用依据 API Key 推导当前用户，只能读取或写入该用户的数据。

## 3. 范围

### 3.1 本期范围

- 在个人中心提供“开发者访问”标签页管理 API Key。
- API Key 格式为单段 `ai4job-<random>`，支持 30 天、60 天、90 天和永不过期。
- Key 创建与轮换时仅显示完整值一次；列表仅显示掩码和元数据。
- 在 `/mcp` 提供 Streamable HTTP MCP 服务。
- 提供面试记录、面试流程、题库查询和题库批量新增四个工具。
- 支持 Key 撤销、过期校验、最近使用时间记录及用户数据隔离。

### 3.2 非本期范围

- 团队共享 Key、组织级 Key、额度计费、IP 白名单和工具级 scope。
- OAuth、动态客户端注册、使用 MCP 修改投递、面试流程、简历或既有题目。
- 将 MCP 拆分为独立服务进程。

## 4. 用户故事

1. 用户在个人中心创建一把供 Codex 使用的 API Key，并选择有效期。
2. 用户复制仅出现一次的完整 Key，写入 MCP 客户端配置。
3. MCP 客户端使用 Key 查看该用户的面试流程和题库。
4. MCP 客户端把整理好的面试题批量写入该用户题库。
5. 用户发现 Key 泄露或不再使用时，在个人中心立即撤销或轮换。

## 5. 功能需求

### 5.1 API Key 生命周期

| 项目 | 要求 |
| --- | --- |
| 格式 | `ai4job-<random>`，单段高熵不透明字符串 |
| 创建字段 | 名称、有效期（30/60/90 天或永不过期） |
| 一次性明文 | 仅创建、轮换成功响应中返回；刷新、关闭弹窗和列表查询后不可恢复 |
| 列表字段 | 名称、掩码、状态、创建时间、过期时间、最近使用时间 |
| 撤销 | Key 所属用户可撤销；撤销后 MCP 请求立即失效 |
| 轮换 | 撤销原 Key 并创建新 Key；新 Key 仅显示一次 |
| 归属 | 只能列出、撤销或轮换当前 JWT 用户名下的 Key |

服务端只持久化 Key 的 HMAC-SHA-256 摘要、显示掩码和生命周期元数据，不保存完整 Key 明文。

### 5.2 个人中心交互

入口位于个人中心的“开发者访问”标签页，与求职资料编辑页分离。页面支持：

- 查看当前用户的 Key 表格。
- 创建 Key，并以单选项选择有效期。
- 在不可误关的提示弹窗中展示完整 Key 与复制按钮。
- 确认撤销、确认轮换，操作成功后刷新列表。
- 将已撤销和已过期 Key 明确展示为不可用状态。

## 6. MCP 服务

### 6.1 接入方式

| 项目 | 说明 |
| --- | --- |
| 地址 | `https://<host>/mcp`；本地开发为 `http://localhost:8080/mcp` |
| 传输 | Streamable HTTP |
| 鉴权 | `Authorization: Bearer ai4job-...` |
| 拒绝凭据 | Cookie、`X-Access-Token`、查询参数 Key |
| 用户身份 | 仅由已验证 API Key 推导，不接受 `userId` |

Codex 示例：

```toml
[mcp_servers.ai4job]
url = "https://<host>/mcp"

[mcp_servers.ai4job.http_headers]
Authorization = "Bearer ai4job-<random>"
```

### 6.2 工具清单

| 工具 | 说明 | 输入 | 输出 |
| --- | --- | --- | --- |
| `get_resume_interviews` | 查询指定简历版本关联的投递和面试进展 | `resumeVersionId` 必填；可选分页 | 岗位快照、当前状态、流程节点 |
| `list_interview_workflows` | 查询当前用户全部投递和面试流程 | 可选 `currentStatus`、分页 | 投递记录、简历版本、岗位快照、流程节点 |
| `list_interview_questions` | 查询当前用户题库 | 关键词、标签、枚举筛选、收藏、分页 | 题库条目和分页信息 |
| `create_interview_questions` | 新增当前用户题库题目 | `questions` 数组，1 至 50 条 | 新建题目列表 |

### 6.3 题库字段约束

每条新增题目必须包含 `title`、`content`、`answerAnalysis`、`difficulty`、`questionType`、`masteryStatus`。`tags`、`isFavorite`、`remark` 可选，`source` 由服务端固定为 `mcp`。

| 字段 | 合法值 |
| --- | --- |
| `difficulty` | `easy`、`medium`、`hard` |
| `questionType` | `knowledge`、`scenario`、`coding`、`algorithm`、`project`、`behavior`、`discussion` |
| `masteryStatus` | `unreviewed`、`weak`、`normal`、`mastered` |

批量新增在单一事务中执行。任一条题目非法或数量超过 50 条时，整批不写入。

## 7. 安全与数据规则

- Key 使用随机字节生成，采用 `API_KEY_PEPPER` 计算 HMAC 摘要后查找。
- API Key 不进入数据库明文字段、列表响应或服务日志。
- 无效、撤销或过期 Key 均返回认证失败，不执行工具业务逻辑。
- 所有 MCP 查询和写入均以 Key 所属用户为边界；跨用户资源不可读取或修改。
- `last_used_at` 至多按一分钟窗口更新一次，避免每个 MCP 请求造成写放大。

## 8. 部署要求

生产环境的 self-hosted GitHub Actions Runner 从服务器本地文件 `/etc/ai-for/production.env` 读取配置，再由 Docker Compose 注入 API 容器。必须配置：

```env
API_KEY_PEPPER=<独立于 JWT_SECRET 的高熵随机秘密>
```

不得将该值提交到仓库或写入 GitHub Actions 工作流。更换 `API_KEY_PEPPER` 会使现有 Key 摘要无法匹配，所有已签发 Key 将失效并需要用户重新创建。

## 9. 验收标准

- 用户可在个人中心创建 30/60/90 天和永不过期的 Key。
- 完整 Key 仅在创建或轮换成功时出现一次，列表不返回明文。
- Key 可被撤销，撤销后 `/mcp` 立即返回认证失败。
- 使用有效 Key 可完成 MCP 初始化、工具发现和四个工具调用。
- 题库工具 schema 显式提供所有枚举值，且批量写入最多 50 条并保持原子性。
- 无 Key、错误 Key、Cookie、JWT 头和跨用户资源访问均不能访问他人数据。
- 生产部署未设置 `API_KEY_PEPPER` 时 API 容器无法启动。
