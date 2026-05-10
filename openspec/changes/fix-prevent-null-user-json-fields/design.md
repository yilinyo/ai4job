## Context

`users.education_experiences` 在迁移中定义为 `JSONB NOT NULL DEFAULT '[]'`，但当前 Go 模型字段为 `[]EducationExperience`，零值是 nil slice。GORM `serializer:json` 在创建或保存模型时会把 nil slice 序列化为 SQL `NULL`，而不是省略该列，因此数据库默认值不会生效。

问题不仅存在于注册流程，也可能出现在任何直接构造 `domainuser.User` 并调用仓储创建、更新的路径。读取历史数据时，如果数据库中存在异常 NULL 或 ORM 扫描得到 nil，也应在返回业务层前归一化，避免后续保存再次写回 NULL。

## Goals / Non-Goals

**Goals:**

- 用户注册不会因 `education_experiences` 为 NULL 失败。
- 用户资料更新不会把教育经历保存为 NULL。
- 仓储创建和更新即使收到 nil slice，也会写入 JSON 空数组。
- 查询用户时返回的 `EducationExperiences` 不为 nil，API 响应保持空数组。
- 测试覆盖 service 层和 PostgreSQL 仓储边界。

**Non-Goals:**

- 不新增或修改教育经历业务字段。
- 不调整教育经历“学历唯一”等前端校验规则。
- 不引入新的数据库迁移来改变字段语义。
- 不改变注册、登录、个人资料接口路径和响应 JSON 字段名。

## Decisions

在领域模型上提供 `Normalize` 方法，并在 service 与 repository 边界都调用。
原因：service 层可以保证主要业务流程表达清晰；repository 层可以兜底直接仓储调用、测试或未来新增 service 遗漏初始化的情况。

`nil` 教育经历统一归一化为 `[]EducationExperience{}`。
原因：业务语义中“没有教育经历”应是空数组，不是未知或缺失；这也与 API 响应和数据库默认值语义一致。

查询后也进行归一化。
原因：即使当前数据库约束阻止新的 NULL，历史异常数据、测试库或迁移前数据仍可能让应用层拿到 nil。读取边界归一化能避免后续响应或保存链路再次传播 nil。

## Risks / Trade-offs

- [仓储层修改会影响所有用户写入路径] → Mitigation: 仅对 nil slice 做空数组归一化，不改变非空教育经历内容。
- [AutoMigrate 测试库与 SQL migration 默认值可能不同] → Mitigation: 测试关注应用层写入与读取行为，确保即使没有依赖数据库默认值也能成功。
