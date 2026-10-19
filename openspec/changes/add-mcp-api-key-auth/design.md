## Context

详见 [proposal.md](proposal.md)。当前后端是单个 Go/Gin 进程：浏览器 API 使用 `X-Access-Token` JWT 中间件，路由层把 `userID` 传给应用服务；业务仓储也以 `user_id` 作为查询边界。前端个人信息页是一个 Vue 2 + Element UI 表单页面。系统尚无 API Key 持久化、服务凭据鉴权或 MCP 协议入口。

## Goals / Non-Goals

**Goals:**

- 在不改变既有 JWT 浏览器 API 的前提下，建立可撤销、可追溯的个人 API Key 生命周期。
- 以同一份服务端用户边界向 MCP 暴露首批四个工具，并杜绝客户端指定用户身份。
- 复用现有投递、面试进展和题库领域校验，避免为 MCP 创建与网页端不一致的数据规则。
- 让用户在个人中心完成 Key 管理，并使完整密钥只在创建或重新生成时出现。

**Non-Goals:**

- 首版不提供组织级 Key、团队共享、计费、额度、IP 白名单或细粒度工具 scope。
- 首版不提供 OAuth、动态客户端注册或替代浏览器 JWT 的登录能力。
- 首版不通过 MCP 修改面试流程、投递记录、简历或已有题目，也不抓取外部招聘链接。
- 不将 MCP 独立部署成第二个进程；部署拆分仅在连接数、隔离或扩缩容需求出现后评估。

## Decisions

### 1. API Key 采用单段高熵格式，并仅保存带 pepper 的摘要

新 Key 格式为 `ai4job-<random>`：`random` 使用 `crypto/rand` 生成至少 32 字节随机值并采用 URL-safe 编码。完整 Key 与主流 API Key 一样是单段不透明字符串，用户只需复制一整串值；它不携带可见的用户标识、Key 标识或分隔段。

新增 `api_keys` 表，包含 `id`、`user_id`、`name`、`display_prefix`、`key_digest`、`status`（`active` / `revoked`）、`created_at`、`expires_at`、`last_used_at`、`revoked_at`。`key_digest` 为 `HMAC-SHA-256(API_KEY_PEPPER, full_key)`，具有唯一索引；`user_id` 与有效状态建立查询索引。`expires_at` 对应创建时选择的 30、60 或 90 天；永不过期存为 `NULL`。请求认证时服务端对整串 Key 计算相同摘要并按索引定位记录，随后使用常量时间比较确认摘要，并拒绝已撤销或已过期的记录。`API_KEY_PEPPER` 仅来自服务端环境配置。创建响应中的完整 Key 不写入数据库、应用日志、活动日志或后续列表响应；列表仅返回固定长度掩码和元数据。

选择单段随机 Key 与 HMAC 摘要，是为了符合主流 API Key 使用体验，同时保持 O(1) 查找和数据库泄露后的离线猜测防护。备选方案是公开标识加密钥的双段格式：查询直接但增加了用户无须理解的格式复杂度；只存 bcrypt/Argon2 哈希：安全但无法高效定位 Key。

### 2. 浏览器 Key 管理 API 继续使用既有 JWT 身份

在现有 `/api/v1` JWT 路由组下增加以下端点：

- `GET /api-keys`：列出当前用户 Key 元数据和掩码。
- `POST /api-keys`：以 `expiresInDays`（`30`、`60`、`90` 或 `null`）创建 Key，响应一次性携带完整 Key 与过期时间。
- `POST /api-keys/:id/rotate`：在事务中撤销旧 Key 并创建新 Key，响应一次性携带新完整 Key。
- `DELETE /api-keys/:id`：撤销当前用户的一把 Key，不执行物理删除。

所有端点从 JWT 上下文取 `userID`，仓储查询必须同时匹配 `id` 与 `user_id`。前端在 `frontend/src/models/`、`api/`、`store/modules/` 增加对应类型与请求；个人信息页在资料表单之后加入“开发者访问”分区，使用表格展示名称、掩码、状态、创建/过期/最近使用时间，并通过确认对话框执行撤销或重新生成。创建和重新生成表单以单选控件提供 30 天、60 天、90 天和永不过期；创建和重新生成使用不可关闭前的结果弹窗显示完整 Key，关闭后只刷新元数据。

选择扩展个人信息页而非新增导航项，是为了保持“个人凭据属于账号设置”的信息层级。备选的独立一级菜单会扩大导航面并使用户资料与开发者凭据分散。

### 3. MCP 使用同一服务进程的 Streamable HTTP 传输，并只接受 Bearer API Key

后端引入经过维护的 Go MCP SDK，在现有 HTTP Server 中挂载独立 `/mcp` Streamable HTTP 端点。该端点仅接受 `Authorization: Bearer ai4job-...`，明确拒绝 Cookie、`X-Access-Token` 和查询字符串 Key。API Key 中间件验证前缀、计算整串 Key 的摘要、按摘要查询有效记录、使用常量时间比较确认，并在执行工具前拒绝已撤销或当前时间已超过 `expires_at` 的记录；验证成功后将 `userID`、`apiKeyID` 写入请求上下文，并异步或受控地更新 `last_used_at`。

MCP 请求在 HTTPS 终止后的反向代理或应用 TLS 中传输；CORS 不作为 MCP 安全边界。认证失败统一映射为 MCP 认证错误，不区分 Key 不存在、密钥错误或已撤销。访问日志只记录 `apiKeyID` 或掩码前缀，禁止记录 Authorization 头、完整参数中可能重复的 Key 或工具返回的敏感内容。

选择 Streamable HTTP 而非 stdio，是因为服务需要对外暴露并与现有 HTTP 部署共存；选择独立鉴权中间件而非复用 JWT 中间件，是因为凭据语义、轮换和撤销机制不同。stdio 仅适合本地单用户进程；将 Key 放进 URL 或 MCP 参数会造成代理、日志和会话记录泄露风险。

### 4. MCP 工具调用应用服务适配层，而不复用 Gin Handler

新增 MCP application adapter，以已验证上下文中的 `userID` 调用领域服务或用户范围仓储；MCP tool handler 不调用 Gin HTTP handler，也不接受 `userId` 参数。适配层负责 MCP JSON schema 到现有 DTO/领域类型的转换，以及将领域错误转换为稳定的 MCP 工具错误。

首批工具及参数边界如下：

| 工具 | 输入 | 输出 |
| --- | --- | --- |
| `get_resume_interviews` | `resumeVersionId`、可选分页 | 该简历版本关联的投递岗位快照、当前状态与面试进展 |
| `list_interview_workflows` | 可选状态、分页 | 当前用户全部投递的简历版本、岗位快照、当前状态和流程节点 |
| `list_interview_questions` | 可选关键词、标签、难度、题型、掌握状态、收藏状态、分页 | 当前用户题库条目及分页信息 |
| `create_interview_questions` | `questions` 数组，长度 1-50 | 每条新题目的标识及创建结果 |

`get_resume_interviews` 先确认简历版本归属，再查询投递和进展；其余查询直接带入已验证的 `userID`。批量创建新增 `CreateQuestions` 应用服务，在单一数据库事务内依次复用现有题目 DTO 与领域校验；任一条失败即回滚全部，并返回带数组索引的输入错误。创建结果将 `source` 固定为 `mcp` 或由服务端补充，绝不接收客户端 `userId`。

选择适配层而非让 MCP 直接读数据库，是为了保留现有所有权校验、枚举校验、活动记录和事务边界。备选的“调用现有 REST API”会要求服务端回调自身、耦合 JWT 传输并增加失败链路。

### 5. 可测试的分层边界与错误模型

- API Key domain/repository：覆盖生成、摘要验证、按用户列出、撤销和轮换。
- API Key middleware：覆盖缺失、格式非法、未知、已撤销和有效 Key；只在有效验证后建立身份上下文。
- MCP adapter：覆盖四个 tool schema、用户数据隔离、资源归属校验和批量事务回滚。
- HTTP handler 与 Vue 页面：覆盖创建只显示一次、掩码列表、撤销/轮换刷新，以及 JWT API 与 MCP Key API 不可互换。

浏览器 API 继续使用项目标准响应包络；MCP 使用协议规定的初始化、tool list 和 tool call 响应，同时将可预期业务错误作为结构化工具错误返回。内部数据库、认证和解析错误不返回堆栈、SQL 或 Key 信息。

## Risks / Trade-offs

- [API Key 被用户复制到不安全的客户端] → 完整值只显示一次，支持即时撤销和轮换；上线文档明确要求用户使用客户端的密钥存储。
- [Key 明文进入日志或错误追踪] → 在鉴权、请求日志和错误序列化处显式脱敏 Authorization 头；测试中断言日志不包含完整 Key。
- [MCP 增加写入入口扩大误操作面] → 首版仅允许新增题库，批量上限 50，并要求单事务全有或全无。
- [现有服务的用户范围校验被适配层绕过] → 适配层仅调用带 `userID` 的服务/仓储接口；跨用户资源测试作为必需回归项。
- [每次调用更新 `last_used_at` 带来写放大] → 仅在距上次记录达到设定时间窗口后更新，或将更新时间写入轻量异步队列；认证正确性不依赖该字段。
- [MCP SDK 或传输版本演进] → 将 SDK 限制在接口层，工具逻辑保持在 application adapter，升级时只替换协议适配器并复跑互操作测试。

## Migration Plan

1. 增加配置项 `API_KEY_PEPPER`，生产环境缺失时服务启动失败；为开发和测试环境提供独立安全值。
2. 执行 `api_keys` 上行迁移，部署 domain、repository、JWT 管理 API 与前端页面，但先不公开 MCP 地址。
3. 部署 MCP 端点并在受控环境使用新 Key 完成初始化、四个工具和跨账号隔离验证后，再在网关公开 `/mcp`。
4. 监控认证失败、工具错误和 Key 轮换事件，确认不含明文后逐步开放给用户。
5. 回滚时先从网关下线 `/mcp`；保留 `api_keys` 数据和浏览器管理页以便用户撤销凭据。若需回退数据库，先确认无后续表依赖后执行下行迁移。
