# mock-user-profile-account-isolation Specification

## Purpose
Define Mock-first account isolation for user profile reads and writes so local development simulates backend per-account profile storage.

## Requirements
### Requirement: Mock 个人资料读取必须按当前账号隔离
系统 MUST 在 Mock 环境中基于当前登录 token 识别账号，并只返回该账号的个人资料。

#### Scenario: 不同账号读取个人资料
- **WHEN** `admin` 和 `user` 分别请求个人信息
- **THEN** 系统返回各自账号对应的个人资料

### Requirement: Mock 个人资料更新必须只影响当前账号
系统 MUST 在 Mock 环境中只更新当前 token 对应账号的个人资料，不影响其他账号。

#### Scenario: 当前账号保存个人资料
- **WHEN** 当前账号提交个人信息更新
- **THEN** 系统只更新当前账号资料
- **AND** 其他账号资料保持不变

### Requirement: Mock 用户资料样例必须支持验证隔离
系统 MUST 为至少两个 Mock 账号提供不同个人资料样例，方便验证账号隔离行为。

#### Scenario: 验证账号资料差异
- **WHEN** 开发者分别以不同 Mock 账号登录
- **THEN** 个人信息页展示不同的资料内容
