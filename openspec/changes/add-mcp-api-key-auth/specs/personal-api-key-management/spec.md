## Purpose

为每位已登录用户提供可自助管理且可追溯的个人 API Key，使用户能够安全地授权外部 MCP 客户端访问其个人求职数据，而无需暴露浏览器登录凭据。

## ADDED Requirements

### Requirement: 用户必须能创建个人 API Key
系统 MUST 允许已登录用户从个人中心创建归属于自身账号的 API Key，并要求为 Key 填写名称和有效期。有效期选项 MUST 为 30 天、60 天、90 天或永不过期。新 Key 的完整值 MUST 以 `ai4job-` 作为前缀，且系统 MUST 仅在创建成功后向创建者展示完整值一次。

#### Scenario: 用户成功创建 API Key
- **WHEN** 已登录用户填写有效名称并选择 30 天、60 天、90 天或永不过期后提交创建请求
- **THEN** 系统创建一把归属于该用户的有效 API Key
- **AND** 系统仅在本次创建结果中返回完整的 `ai4job-` Key
- **AND** 系统提示用户立即保存该 Key

### Requirement: API Key 必须记录创建和过期时间
系统 MUST 为每把 API Key 保存创建时间和过期时间；选择永不过期时过期时间为空。系统 MUST 在 Key 过期后立即拒绝其 MCP 鉴权，且列表中 MUST 将其显示为已过期。

#### Scenario: 用户查看有期限 Key 的时间信息
- **WHEN** 用户查看一把选择了固定有效期的 Key
- **THEN** 系统展示该 Key 的创建时间和由有效期计算得到的过期时间

#### Scenario: 用户查看永不过期 Key 的时间信息
- **WHEN** 用户查看一把选择永不过期的 Key
- **THEN** 系统展示创建时间
- **AND** 系统将过期时间显示为永不过期

#### Scenario: 过期 Key 发起 MCP 调用
- **WHEN** MCP 客户端使用当前时间已超过过期时间的 Key 调用工具
- **THEN** 系统返回认证失败
- **AND** 系统不返回业务数据且不产生业务写入

#### Scenario: 用户关闭一次性展示后再次查看
- **WHEN** 用户在创建成功后关闭或刷新页面并再次进入 Key 管理区域
- **THEN** 系统不再返回该 Key 的完整值
- **AND** 系统仅展示其名称、前缀、状态、创建时间和最近使用时间

### Requirement: API Key 必须只以安全摘要持久化
系统 MUST 保存 API Key 的所属用户、显示前缀、名称、不可逆安全摘要、状态、创建时间、过期时间和最近使用时间；系统 MUST NOT 持久化、日志记录或在后续查询接口返回完整 Key。

#### Scenario: 服务端持久化新 Key
- **WHEN** 系统创建新的 API Key
- **THEN** 持久化记录包含该 Key 的安全摘要和可识别前缀
- **AND** 持久化记录不包含完整 Key 明文

### Requirement: 用户必须只能管理自己的 API Key
系统 MUST 仅向当前登录用户列出其名下 API Key，并且 MUST 拒绝用户查看、撤销或重新生成其他用户的 Key。

#### Scenario: 用户查看 API Key 列表
- **WHEN** 已登录用户打开 Key 管理区域
- **THEN** 系统只返回该用户创建的 Key 元数据

#### Scenario: 用户操作其他账号的 Key
- **WHEN** 用户构造请求尝试撤销或重新生成不属于自己的 Key
- **THEN** 系统拒绝该请求
- **AND** 目标 Key 的状态和完整值均不被泄露

### Requirement: 用户必须能撤销或重新生成 API Key
系统 MUST 允许 Key 所属用户撤销有效 Key；撤销后的 Key MUST 立即不能再用于 MCP 鉴权。重新生成 MUST 创建新的 Key 并撤销被替代的 Key，且新完整值仍只展示一次。

#### Scenario: 用户撤销 API Key
- **WHEN** Key 所属用户确认撤销一把有效 Key
- **THEN** 系统将该 Key 标记为已撤销
- **AND** 后续使用该 Key 的 MCP 请求被拒绝

#### Scenario: 用户重新生成 API Key
- **WHEN** Key 所属用户请求重新生成一把有效 Key
- **THEN** 系统撤销原 Key 并创建一把新的有效 Key
- **AND** 系统仅在本次结果中展示新 Key 的完整值一次
