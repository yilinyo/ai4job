## Context

当前真实后端已提供 `/api/v1` 用户、简历仓库、简历版本等接口，并通过 `pkg/response` 保持统一 envelope。用户相关接口已有 `UserInfoResponse` 等 DTO，但简历仓库和简历版本接口仍直接返回领域/GORM 实体，例如 `domainrepo.ResumeRepo` 和 `domainversion.ResumeVersion`。

Go JSON 序列化会使用导出字段名，导致响应出现 `ID`、`UserID`、`CreatedAt`、`DeletedAt` 等字段。前端 Vue/TypeScript 模型期望小驼峰字段，例如 `id`、`userId`、`createdAt`，因此真实后端返回无法被现有前端直接消费。

本变更以“后端响应边界 DTO 化”为核心，不改变数据库表结构，不改变统一响应 envelope，不要求前端兼容后端内部字段名。

## Goals / Non-Goals

**Goals:**

- 为现有后端响应增加明确的 response DTO 层，隔离领域实体、GORM 模型和 API 契约。
- 让 `/api/v1` 已实现接口返回前端可直接消费的小驼峰 JSON。
- 修复简历仓库列表点击查看时 `repo.id` 为 `undefined` 的问题。
- 避免响应中暴露 `DeletedAt`、GORM 软删除结构、Go 导出字段名等内部实现细节。
- 用 handler 层测试锁定响应字段名，防止后续回退到直接返回实体。

**Non-Goals:**

- 不重构数据库 schema 或迁移历史数据。
- 不新增新的业务模块或未路由接口。
- 不要求前端批量改成兼容 `ID`、`UserID` 等大写字段。
- 不在本变更中补齐前端模型中尚未由后端持久化的业务字段，除非现有接口已经具备对应数据来源。
- 不改变 `code`、`message`、`data`、分页 `list/total/page/pageSize` 的统一响应格式。

## Decisions

### 1. 响应 DTO 放在 application/dto 包

在 `backend/internal/application/dto` 中新增或扩展 response DTO，例如：

- `ResumeRepoResponse`
- `ResumeVersionResponse`
- `ToResumeRepoResponse`
- `ToResumeRepoResponses`
- `ToResumeVersionResponse`
- `ToResumeVersionResponses`

理由：该包已经承载请求 DTO 和用户响应 DTO，handler 已依赖该包；把转换函数放在这里可以减少新包引入和循环依赖风险。

替代方案是在 `internal/interface/response` 或各 handler 文件内定义响应结构。前者会增加新的分层包和迁移成本；后者容易让转换逻辑散落在 handler 中，不利于复用和测试。

### 2. Handler 返回 DTO，不直接返回 domain entity

简历仓库和简历版本 handler 在调用 service 后，将 domain entity 转为 response DTO，再调用 `response.Success` 或 `response.SuccessPage`。

示例目标：

```go
response.SuccessPage(c, dto.ToResumeRepoResponses(rows), total, p.Page, p.PageSize)
response.Success(c, dto.ToResumeRepoResponse(row))
```

理由：handler 是 HTTP 边界，最适合完成“内部模型到 API 契约”的转换。service 继续返回领域对象，避免把前端字段命名侵入业务层。

### 3. 时间字段保持字符串，由标准 JSON 时间格式输出

DTO 的 `createdAt`、`updatedAt` 可以使用 `time.Time` 并通过 `json` tag 输出小驼峰字段；Go 默认 RFC3339/RFC3339Nano 时间字符串可继续被前端 `new Date(...)` 解析。

替代方案是转换成自定义格式字符串。该方案会引入额外格式约定，当前前端已能处理 ISO/RFC3339 时间，不需要增加复杂度。

### 4. 删除和内部字段默认不出现在响应中

`DeletedAt`、密码、数据库 serializer 内部状态、GORM metadata 等字段不进入 response DTO。对于可选字段，按前端模型需要返回空字符串、空数组或省略字段，避免返回内部零值对象。

### 5. 前端格式以现有 TypeScript models 为准

响应字段以 `frontend/src/models` 中已存在的模型命名为主要依据，例如：

- `ResumeRepo.id`
- `ResumeRepo.userId`
- `ResumeRepo.name`
- `ResumeRepo.description`
- `ResumeRepo.createdAt`
- `ResumeRepo.updatedAt`
- `ResumeVersion.id`
- `ResumeVersion.repoId`
- `ResumeVersion.title`
- `ResumeVersion.content`
- `ResumeVersion.versionNum`
- `ResumeVersion.isDefault`

如果前端模型包含当前后端实体没有的数据字段，应在设计和任务中记录缺口，不通过假数据伪造业务含义。

## Risks / Trade-offs

- [Risk] 前端模型包含后端当前没有持久化的字段，例如简历仓库历史前端模型中的 `jobType`、`targetPosition`、`versionCount`。→ Mitigation：本变更只保证已有后端数据字段命名正确；缺失业务字段需要另开 change 或在本 change 实施时明确补充数据来源后再纳入。
- [Risk] 只给 domain entity 加 `json` tag 看起来更快。→ Mitigation：这会继续把数据库/领域结构绑定到 API 响应，且容易暴露 `DeletedAt`；本变更采用 DTO 层作为稳定边界。
- [Risk] 不同 handler 可能遗漏 DTO 转换。→ Mitigation：通过 `rg "response.Success.*row|rows"` 排查直接返回 entity 的点，并在 handler 测试中断言响应不包含大写内部字段。
- [Risk] DTO 转换函数增加少量重复代码。→ Mitigation：先保持显式转换，避免过早引入反射或泛型 mapper；字段变化时可读性更高。
