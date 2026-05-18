## Why

用户注册时，后端通过 GORM 创建 `users` 记录会显式写入 `education_experiences = NULL`。虽然数据库字段定义了 `JSONB NOT NULL DEFAULT '[]'`，但 PostgreSQL 默认值只在字段被省略或显式使用 `DEFAULT` 时生效，无法覆盖应用层传入的 `NULL`，最终触发 `SQLSTATE 23502`。

该问题会阻断新用户注册，也可能在用户资料更新时将教育经历写回 `NULL`。需要在后端所有用户写入与读取边界统一消除 `education_experiences` 的 nil/NULL 状态。

## What Changes

- 新增用户领域模型归一化逻辑，确保教育经历为空时使用空数组而不是 nil。
- 注册创建用户前初始化 `EducationExperiences`。
- 更新个人资料时，如果请求未传教育经历或传入 nil，写入空数组。
- 用户持久化仓储在创建、更新和查询返回时进行兜底归一化，覆盖 service 之外的直接调用场景。
- 补充单元测试和持久化集成测试，验证注册、更新、仓储创建和历史 NULL 数据读取均不会暴露 nil/NULL 教育经历。

## Capabilities

### New Capabilities

- `user-profile-json-defaults`: 用户资料 JSON 数组字段在应用层始终归一化为空数组，避免数据库 NOT NULL 约束被 nil JSON 序列化击穿。

### Modified Capabilities

- 无。

## Impact

- 影响后端用户领域模型：新增归一化方法。
- 影响注册与个人资料更新：空教育经历统一保存为 JSON 空数组。
- 影响 PostgreSQL 用户仓储：创建、更新、查询边界新增归一化兜底。
- 不改变 API 响应结构；`educationExperiences` 继续返回数组。
