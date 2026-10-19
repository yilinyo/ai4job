## Purpose

为本项目提供由生产服务器主动领取任务的可重复部署能力，使每次生产更新均可按改动范围发布、识别版本并验证运行状态。

## ADDED Requirements

### Requirement: 服务器 Runner 执行受控生产部署
系统 SHALL 在标记为 production 的 self-hosted runner 上执行生产部署任务。部署任务 MUST 仅由 `master` 分支推送或人工触发，且不依赖 GitHub Actions 通过 SSH 入站连接服务器。

#### Scenario: master 提交触发部署
- **WHEN** `master` 收到通过验证的提交
- **THEN** production runner SHALL 通过其主动建立的 HTTPS 连接领取并执行部署任务

#### Scenario: 拉取请求不触发生产部署
- **WHEN** 拉取请求或来自非主分支的提交触发工作流
- **THEN** 系统 MUST NOT 在 production runner 上执行部署步骤

### Requirement: 按改动范围选择性更新服务
系统 SHALL 根据本次提交中前端、后端、迁移文件和部署配置的改动范围，选择更新 Web、API、迁移任务或全部生产服务。仅文档改动 MUST NOT 触发服务重建。

#### Scenario: 仅前端改动
- **WHEN** 本次提交仅修改 frontend 目录
- **THEN** 系统 SHALL 仅重建并更新 Web 服务，不重启 API、PostgreSQL 或 Redis

#### Scenario: 后端改动
- **WHEN** 本次提交修改 backend 目录或数据库迁移文件
- **THEN** 系统 SHALL 在更新 API 前执行迁移任务，并仅更新 API 服务及必要依赖

#### Scenario: 部署配置改动
- **WHEN** 本次提交修改部署编排、Nginx 或镜像构建配置
- **THEN** 系统 SHALL 重新协调全部生产服务

### Requirement: 部署版本可被容器和镜像识别
系统 SHALL 为每次部署生成包含基础发布版本、工作流运行序号和 Git 短 SHA 的不可变版本标识，并将该标识写入 Web 与 API 镜像标签及运行容器元数据。

#### Scenario: 查询当前容器版本
- **WHEN** 运维人员列出运行中的容器或检查容器标签
- **THEN** 系统 SHALL 显示本次部署的版本标识和 Git SHA

#### Scenario: 查询历史镜像
- **WHEN** 运维人员列出 Web 或 API 镜像
- **THEN** 系统 SHALL 能区分不同已构建部署版本的镜像标签

### Requirement: 部署结果经过健康检查并可追溯
系统 SHALL 在更新服务后检查目标服务和 API 的可用性，并把版本、提交、服务范围、时间和结果写入服务器持久部署记录。健康检查失败 MUST 使工作流失败且保留失败记录。

#### Scenario: 成功部署
- **WHEN** 已更新服务通过健康检查
- **THEN** 系统 SHALL 记录成功状态并将工作流标记为成功

#### Scenario: 健康检查失败
- **WHEN** 已更新服务未在限定时间内通过健康检查
- **THEN** 系统 SHALL 记录失败状态并将工作流标记为失败
