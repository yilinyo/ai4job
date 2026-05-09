# register-email-verification Specification

## Purpose
TBD - created by archiving change enhance-register-email-verification. Update Purpose after archive.
## Requirements
### Requirement: 注册页必须要求填写邮箱和验证码
系统 MUST 在注册页要求用户填写邮箱和邮箱验证码，并在表单层面对邮箱格式、验证码必填、密码确认进行校验。

#### Scenario: 注册表单展示邮箱验证码字段
- **WHEN** 用户打开注册页
- **THEN** 系统展示用户名、密码、确认密码、邮箱、邮箱验证码和发送验证码操作

#### Scenario: 邮箱或验证码缺失时提交注册
- **WHEN** 用户未填写邮箱或邮箱验证码就提交注册
- **THEN** 系统在注册页提示对应字段必填，不提交注册请求

#### Scenario: 邮箱格式错误时提交注册
- **WHEN** 用户填写格式错误的邮箱并提交注册
- **THEN** 系统在注册页提示邮箱格式错误，不提交注册请求

### Requirement: 注册页必须支持发送邮箱验证码交互
系统 MUST 允许用户在填写有效邮箱后触发发送验证码操作，并在发送成功后给出反馈和倒计时状态。

#### Scenario: 有效邮箱发送验证码
- **WHEN** 用户填写有效邮箱并点击发送验证码
- **THEN** 系统调用发送验证码接口，发送成功后提示用户并启动倒计时

#### Scenario: 倒计时期间重复发送验证码
- **WHEN** 验证码发送成功且倒计时尚未结束
- **THEN** 系统禁用发送验证码操作或阻止重复发送

#### Scenario: 修改邮箱后重新填写验证码
- **WHEN** 用户修改邮箱字段
- **THEN** 系统清空已输入验证码，要求用户基于新邮箱重新获取或填写验证码

### Requirement: 注册成功后邮箱必须展示在个人信息页
系统 MUST 在邮箱验证码注册成功后，将注册邮箱作为账号个人资料邮箱，并在用户登录后通过个人信息页展示。

#### Scenario: 注册后查看个人信息邮箱
- **WHEN** 用户使用邮箱验证码注册成功并登录新账号后进入个人信息页
- **THEN** 个人信息页邮箱字段显示注册时填写的邮箱

