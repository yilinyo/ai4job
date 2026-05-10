## Why

当前后端部分接口直接返回领域/GORM 结构体，JSON 字段名暴露为 `ID`、`UserID`、`CreatedAt` 等 Go 导出字段名，和前端模型期望的 `id`、`userId`、`createdAt` 不一致。该问题已经导致简历仓库列表点击查看时读取 `repo.id` 得到 `undefined`，真实后端接入前端流程被阻断。

需要在后端响应边界统一增加 response DTO 层，让现有 `/api/v1` 接口返回稳定、前端友好的小驼峰 JSON，同时避免数据库字段、GORM 元数据和内部领域模型泄漏到 API 契约中。

## What Changes

- 为现有后端响应增加专用 response DTO，不再直接把领域实体或 GORM 模型作为接口响应体返回。
- 统一现有 `/api/v1` 接口响应字段为前端约定的小驼峰格式，例如 `id`、`userId`、`name`、`description`、`createdAt`、`updatedAt`。
- 简历仓库列表和详情接口返回可被前端 `ResumeRepo` 模型直接消费的数据，修复 `repo.id` 为 `undefined` 的问题。
- 简历版本列表和详情接口返回可被前端版本详情、列表、编辑流程消费的数据，不暴露 Go 导出字段名。
- 用户登录、个人资料和资料更新接口继续保持前端现有 `token`、`userInfo`、`username`、`educationExperiences` 等字段格式，并通过 DTO 明确响应契约。
- 分页响应外层继续保持现有统一结构：`code`、`message`、`data.list`、`data.total`、`data.page`、`data.pageSize`。

## Capabilities

### New Capabilities

- `backend-response-dto-contract`: 后端现有接口通过 response DTO 返回前端格式的小驼峰 JSON，并隔离内部领域/数据库结构。

### Modified Capabilities

- `resume-repo-job-type-management`: 简历仓库列表、详情、创建后刷新、编辑、删除等前端流程依赖后端返回小驼峰字段，尤其是 `id` 可用于路由跳转和后续详情请求。
- `resume-version-detail-management`: 简历版本相关列表与详情接口应返回小驼峰字段，供前端版本列表、详情查看和编辑流程直接消费。
- `user-profile-management-as-top-level-entry`: 用户资料相关后端响应应继续稳定返回前端资料页使用的小驼峰字段，并由 DTO 明确而不是依赖领域实体序列化。

## Impact

- 影响后端 application/interface DTO 设计：新增或扩展 response DTO 与转换函数。
- 影响后端 handler：`user`、`resume-repo`、`resume-version` 等现有响应点需要在返回前做 DTO 转换。
- 影响后端 handler/router 测试：需要断言响应 JSON 字段名为小驼峰，且不包含 `ID`、`UserID`、`DeletedAt` 等内部字段。
- 不改变统一响应 envelope 的 `code`、`message`、`data` 结构。
- 不要求前端批量兼容大写字段；前端应继续按现有 TypeScript models 消费真实后端。
