# 后端 DDD 架构重构设计

**日期：** 2026-05-15  
**范围：** `backend/` 内部架构重构  
**目标：** 在不改变业务 API 行为的前提下，将当前轻量 DDD 结构重构为更清晰的领域驱动分层，并引入轻量充血模型。

## 1. 背景

当前后端已经按 `interface / application / domain / infrastructure` 分层组织，并在 `internal/domain/<aggregate>/repository.go` 中定义仓储接口。这个结构可用，但仍偏轻量：

- repository 接口散落在聚合根包内，包边界不够明确。
- domain entity 主要承担数据结构和 `Validate()` 校验职责。
- application service 中存在较多直接构造 entity、设置默认值、修改字段、再调用 `Validate()` 的逻辑。
- infrastructure repository 直接实现 domain repository，但缺少编译期接口断言。

本次重构目标不是改业务能力，而是让领域层更真实，让 application 更专注于用例编排。

## 2. 设计目标

1. 每个聚合下新增独立 repository 子包。
2. 引入轻量充血模型，把创建默认值、更新规则、状态合法性和不变量维护下沉到 domain。
3. 保持 HTTP API、请求 DTO、响应 DTO、错误码、数据库 migration 不变。
4. 保持 application service 的用例流程不变，只减少直接字段操作。
5. 保持 infrastructure repository 的持久化行为不变，只适配新的 repository 接口位置。

## 3. 非目标

本次不做以下改动：

- 不修改路由路径、HTTP 方法、请求 JSON、响应 JSON。
- 不修改数据库表结构、索引或迁移文件。
- 不引入 CQRS、Event Sourcing、Domain Event、UnitOfWork 或 DI 框架。
- 不把 entity 改成直接依赖 repository。
- 不移除 domain entity 上现有 GORM tag。
- 不引入 persistence model 与 mapper。

保留 GORM tag 是当前阶段的折中。严格 DDD 下 domain 不应感知 ORM，但移除 tag 需要新增持久化模型和映射层，改动范围较大，容易影响既有行为，适合后续单独演进。

## 4. 目标目录结构

每个聚合保持 entity 包，并在其下新增 repository 子包：

```text
backend/internal/domain/
├── user/
│   ├── entity.go
│   └── repository/
│       └── repository.go
├── resume_repo/
│   ├── entity.go
│   └── repository/
│       └── repository.go
├── resume_version/
│   ├── entity.go
│   └── repository/
│       └── repository.go
├── job_posting/
│   ├── entity.go
│   └── repository/
│       └── repository.go
├── resume_application/
│   ├── entity.go
│   └── repository/
│       └── repository.go
└── interview_progress/
    ├── entity.go
    └── repository/
        └── repository.go
```

依赖方向保持单向：

```text
interface handler
  -> application service
    -> domain aggregate + domain repository interface
      <- infrastructure postgres implementation
```

具体含义：

- `interface` 只依赖 application service interface、DTO、response。
- `application` 依赖 domain aggregate 和 `domain/<aggregate>/repository`。
- `domain` 不依赖 application、interface、infrastructure。
- `infrastructure` 依赖 domain aggregate，并实现 domain repository interface。
- `cmd/server` 继续手动组装依赖。

## 5. Repository 子包设计

现状：

```go
package resume_repo

type Repository interface { ... }
```

目标：

```go
package repository

type Repository interface { ... }
```

路径示例：

```text
internal/domain/resume_repo/repository/repository.go
```

Application 使用方式：

```go
domainrepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo"
resumereporepo "github.com/yilin/ai-for-backend/internal/domain/resume_repo/repository"

type ResumeRepoService struct {
    repo resumereporepo.Repository
}
```

Infrastructure repository 增加编译期断言：

```go
var _ repository.Repository = (*ResumeRepoRepo)(nil)
```

这样可以在接口移动后确保实现没有漂移。

## 6. 轻量充血模型设计

### 6.1 创建方法

将 application 中的实体创建、默认值和创建后校验下沉到 domain 构造方法。

示例：

```go
posting, err := job_posting.NewJobPosting(job_posting.NewJobPostingParams{
    ID:              uuid.NewString(),
    UserID:          userID,
    CompanyName:     req.CompanyName,
    CompanyIntro:    req.CompanyIntro,
    JobTitle:        req.JobTitle,
    JobRequirements: req.JobRequirements,
    Base:            req.Base,
    SalaryRange:     req.SalaryRange,
    SourcePlatform:  req.SourcePlatform,
    SourceURL:       req.SourceURL,
    Status:          req.Status,
    Remark:          req.Remark,
})
```

Domain 构造方法负责：

- 设置默认状态。
- 初始化必要默认字段。
- 校验必填字段、长度、枚举和 URL。
- 返回合法实体或领域错误。

优先覆盖：

- `user`：注册用户初始化教育经历、复用密码设置规则。
- `resume_repo`：简历仓库创建默认值和字段约束。
- `resume_version`：简历版本创建约束。
- `job_posting`：岗位创建默认 `pending`。
- `resume_application`：投递创建默认 `applied`。
- `interview_progress`：面试进展创建约束。

### 6.2 更新方法

将 application 中的直接字段赋值下沉为领域行为。

现状模式：

```go
if req.CompanyName != nil {
    posting.CompanyName = *req.CompanyName
}
if req.Status != nil {
    posting.Status = *req.Status
}
if err := posting.Validate(); err != nil {
    return err
}
```

目标模式：

```go
err := posting.ApplyUpdate(job_posting.UpdateSpec{
    CompanyName: req.CompanyName,
    Status:      req.Status,
})
```

Domain 更新方法负责：

- 只更新非 nil 字段。
- 维护状态和字段不变量。
- 更新后统一校验。
- 返回与现有行为一致的错误。

优先覆盖：

- `User.UpdateProfile(...)`
- `ResumeRepo.ApplyUpdate(...)`
- `ResumeVersion.ApplyUpdate(...)`
- `JobPosting.ApplyUpdate(...)`
- `ResumeApplication.ApplyUpdate(...)`
- `InterviewProgress.ApplyUpdate(...)`

### 6.3 状态合法性

将 application 中重复的 status switch 移入 domain：

```go
func (s Status) Valid() bool
```

构造和更新方法内部应调用 `Valid()`。Application 如需提前校验查询过滤条件，也调用 domain 暴露的方法，而不是复制 switch。

## 7. Application 层设计

Application service 保留现有职责：

- DTO 到 domain 参数的转换。
- 生成 ID。
- 调用 domain 构造/更新方法。
- 调用 repository。
- 分页默认值处理。
- 跨聚合所有权校验。
- 用例流程编排。

Application service 不再负责：

- 直接设置创建默认状态。
- 大量直接修改 entity 字段。
- 重复维护 status switch。
- 手动调用零散字段校验逻辑。

跨聚合查库校验仍留在 application。例如投递记录创建/更新时：

- 岗位必须属于当前用户。
- 简历仓库必须属于当前用户。
- 简历版本必须属于当前用户可访问的仓库。

这些规则需要 repository 查询，因此不放进 entity。

## 8. Interface 层设计

Interface 层保持当前设计：

- handler 继续通过局部 service interface 依赖 application。
- router 继续只负责路由注册。
- middleware 继续处理认证边界。
- handler 不直接依赖 repository。
- handler 不直接处理持久化细节。

API 行为保持不变：

- 路由路径不变。
- 请求 DTO 不变。
- 响应 DTO 不变。
- `response.Response` 包装结构不变。
- `Update` 后重新 `Get` 并返回最新数据的行为不变。

可选的小整理：如果多个 handler 共用分页 helper，可将 helper 移到 handler 内部公共文件。但这不是本次重构的必要条件，不能因此扩大范围。

## 9. Infrastructure 层设计

Postgres repository 保留当前持久化行为：

- 继续使用 GORM。
- 继续直接持久化 domain entity。
- 查询条件、排序、分页、错误映射保持不变。
- 继续将 `gorm.ErrRecordNotFound` 映射为 `pkg/errors.ErrNotFound`。
- 继续将唯一约束错误映射为 `pkg/errors.ErrDuplicate`。

本次只做必要适配：

- import 新的 `domain/<aggregate>/repository` 包。
- 增加编译期接口实现断言。

## 10. pkg 层设计

`pkg/errors` 和 `pkg/response` 保持不变。

当前 domain 依赖 `pkg/errors`，从严格 DDD 看可以进一步迁移到 `internal/domain/errors`。但这会影响大量 import，并且不是本次重构的关键目标，因此暂不调整。

## 11. 测试策略

重构完成后至少执行：

```bash
cd backend && go test ./...
```

测试关注点：

1. repository 接口移动后所有包能编译。
2. application service 单测仍通过。
3. domain entity 单测覆盖新增构造和更新方法。
4. handler/router 单测验证 API 行为未变。
5. postgres repository 单测验证持久化行为未变。

如发现新增 domain 方法缺少覆盖，应优先补 domain 层测试，而不是只依赖 application 测试。

## 12. 实施顺序

1. 为每个聚合创建 `repository` 子包并移动接口。
2. 更新 application、infrastructure、tests 的 import 和类型引用。
3. 在 infrastructure repository 增加编译期接口断言。
4. 为 domain aggregate 增加构造参数结构和构造方法。
5. 为 domain aggregate 增加更新参数结构和 `ApplyUpdate` / `UpdateProfile` 方法。
6. 将 application service 中的创建和更新逻辑改为调用 domain 方法。
7. 将重复 status 校验改为调用 domain `Valid()`。
8. 运行并修复后端测试。

## 13. 风险与控制

主要风险是重构时无意改变业务行为。控制方式：

- 不改 DTO、路由、response、migration。
- 不改 repository 查询条件和更新字段集合。
- 优先小步移动接口，再做 domain 方法下沉。
- 每个聚合完成后运行相关测试或全量 `go test ./...`。
- 新增 domain 方法时保持原有错误类型和默认值一致。

## 14. 验收标准

1. `backend/internal/domain/<aggregate>/repository/repository.go` 存在并承载对应仓储接口。
2. application service 不再依赖 `domain.<Repository>`，而是依赖 repository 子包接口。
3. infrastructure postgres repository 显式实现对应 repository 接口。
4. 主要创建/更新用例通过 domain 构造或更新方法表达业务意图。
5. API 路径、请求响应结构、错误码和数据库 schema 无变化。
6. `cd backend && go test ./...` 通过。
