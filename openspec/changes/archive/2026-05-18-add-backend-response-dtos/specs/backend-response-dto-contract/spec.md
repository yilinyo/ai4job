## ADDED Requirements

### Requirement: 后端响应必须通过 DTO 暴露前端字段格式
后端现有 `/api/v1` 接口 MUST 在 HTTP 响应边界使用 response DTO 输出数据，不得直接将领域实体或 GORM 模型作为 `data` 或 `data.list` 返回。

#### Scenario: 简历仓库列表返回 DTO 字段
- **WHEN** 前端请求 `GET /api/v1/resume-repos`
- **THEN** 响应 `data.list` 中每个仓库 MUST 使用 `id`、`userId`、`name`、`description`、`createdAt`、`updatedAt` 等小驼峰字段
- **AND** 响应项 MUST NOT 包含 `ID`、`UserID`、`Name`、`CreatedAt`、`UpdatedAt`、`DeletedAt` 等 Go 导出字段名或内部删除字段

#### Scenario: 简历版本列表返回 DTO 字段
- **WHEN** 前端请求 `GET /api/v1/resume-repos/:repoId/versions`
- **THEN** 响应 `data` 中每个版本 MUST 使用 `id`、`repoId`、`title`、`content`、`versionNum`、`isDefault`、`createdAt`、`updatedAt` 等小驼峰字段
- **AND** 响应项 MUST NOT 包含 `ID`、`RepoID`、`Title`、`Content`、`VersionNum`、`IsDefault`、`CreatedAt`、`UpdatedAt`、`DeletedAt` 等 Go 导出字段名或内部删除字段

#### Scenario: 用户资料响应保持 DTO 字段
- **WHEN** 前端登录、获取个人资料或更新个人资料
- **THEN** 用户资料响应 MUST 使用 `id`、`username`、`realName`、`email`、`educationExperiences` 等前端字段格式
- **AND** 响应 MUST NOT 暴露密码、数据库字段名或 GORM 内部字段

### Requirement: 统一响应 envelope 必须保持兼容
后端 response DTO 变更 MUST 保持现有统一响应 envelope，不得破坏前端 request 拦截器对 `code`、`message`、`data` 的处理。

#### Scenario: 普通成功响应
- **WHEN** 任一现有 `/api/v1` 接口成功返回单个对象或空结果
- **THEN** 响应 MUST 保持 `code`、`message`、`data` 三层结构
- **AND** `code` MUST 继续使用现有成功码 `20000`

#### Scenario: 分页成功响应
- **WHEN** 任一现有分页接口成功返回列表
- **THEN** 响应 MUST 保持 `data.list`、`data.total`、`data.page`、`data.pageSize`
- **AND** `data.list` 中的每个元素 MUST 已完成 response DTO 转换
