# 后端设计文档 — 求职流程管理系统

**日期：** 2026-05-10  
**技术栈：** Go（最新版）+ Gin + GORM + PostgreSQL + Redis  
**架构：** 经典四层 DDD  
**开发模式：** TDD

---

## 1. 整体架构

### 1.1 分层结构

```
Interface 层 → Application 层 → Domain 层 ← Infrastructure 层
```

- **Interface 层**：HTTP Handler、Middleware、Router，负责请求解析与响应格式化
- **Application 层**：用例编排（Use Case），调用 Domain 层，不含业务规则
- **Domain 层**：实体、值对象、仓储接口、领域服务，零框架依赖
- **Infrastructure 层**：GORM 仓储实现、Redis 客户端、JWT、邮件

依赖方向单向，Domain 层通过接口与 Infrastructure 层解耦（依赖倒置）。

### 1.2 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 入口，依赖注入，启动
├── internal/
│   ├── interface/
│   │   ├── handler/             # HTTP handlers（user, resume, job, question 等）
│   │   ├── middleware/          # JWT 认证、日志、Recovery、CORS
│   │   └── router/              # 路由注册
│   ├── application/
│   │   ├── user/                # 用户用例
│   │   ├── resume/              # 简历用例
│   │   ├── job/                 # 岗位/投递用例
│   │   ├── question/            # 题库/闪卡用例
│   │   └── dto/                 # 请求/响应 DTO
│   ├── domain/
│   │   ├── user/                # User 聚合：entity + repo 接口 + 领域服务
│   │   ├── resume/              # ResumeRepo + ResumeVersion 聚合
│   │   ├── job/                 # JobPosting + ResumeApplication + InterviewProgress 聚合
│   │   └── question/            # InterviewQuestion + QuestionOccurrence + FlashcardReview 聚合
│   └── infrastructure/
│       ├── persistence/postgres/ # GORM 仓储实现
│       ├── cache/               # Redis 封装（黑名单、验证码）
│       ├── auth/                # JWT 生成/验证
│       └── email/               # 验证码（开发环境 mock）
├── pkg/
│   ├── response/                # 统一响应结构
│   ├── errors/                  # 业务错误码
│   ├── config/                  # Viper 配置加载
│   └── pagination/              # 分页结构
├── tests/
│   └── integration/             # 集成测试（testcontainers）
├── migrations/                  # SQL 迁移文件（golang-migrate）
├── .env.example
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## 2. 领域模型

### 2.1 聚合划分

| 聚合 | 聚合根 | 子实体 | 值对象 |
|------|--------|--------|--------|
| User | User | EducationExperience（JSONB） | Email, Phone |
| Resume | ResumeRepo | ResumeVersion | JobType, VersionString |
| Job | ResumeApplication | JobPosting, InterviewProgress | ApplicationStatus, JobSnapshot, SalaryRange |
| Question | InterviewQuestion | QuestionOccurrence, FlashcardReview | MasteryStatus, Difficulty, ReviewMode |

### 2.2 数据库表设计（PostgreSQL）

所有表：
- 主键使用 **UUID**
- 时间字段使用 **TIMESTAMPTZ**
- 统一包含 **`deleted_at TIMESTAMPTZ NULL`** 软删除字段（GORM `gorm.DeletedAt`）
- 涉及唯一约束的字段使用 PostgreSQL **部分索引**（`WHERE deleted_at IS NULL`）

```sql
-- users
id UUID PK | username VARCHAR(50) UNIQUE | password_hash VARCHAR
email VARCHAR(100) UNIQUE | real_name VARCHAR | age INT
phone VARCHAR(20) | job_intention VARCHAR | avatar VARCHAR
location VARCHAR | personal_advantage TEXT
education_experiences JSONB        -- 内嵌存储，避免过度拆表
created_at TIMESTAMPTZ | updated_at TIMESTAMPTZ | deleted_at TIMESTAMPTZ

-- resume_repos
id UUID PK | user_id UUID FK→users | name VARCHAR(100)
job_type VARCHAR(20)               -- campus/social/internship
target_position VARCHAR(100) | description TEXT
version_count INT DEFAULT 0
created_at | updated_at | deleted_at

-- resume_versions
id UUID PK | repo_id UUID FK→resume_repos
version VARCHAR(20)                -- e.g. v1.0.0
title VARCHAR(200) | content TEXT | remark TEXT
is_active BOOLEAN DEFAULT false
created_at | updated_at | deleted_at
UNIQUE INDEX ON (repo_id, version) WHERE deleted_at IS NULL

-- job_postings
id UUID PK | user_id UUID FK→users
company_name VARCHAR | company_intro TEXT | job_title VARCHAR
job_requirements TEXT | base VARCHAR | salary_range VARCHAR
source_platform VARCHAR | source_url VARCHAR
status VARCHAR(20)                 -- pending/applied/not_fit/closed
remark TEXT
created_at | updated_at | deleted_at

-- resume_applications
id UUID PK | job_posting_id UUID FK→job_postings (nullable)
resume_version_id UUID FK→resume_versions | repo_id UUID FK→resume_repos
-- 快照字段（创建时从 JobPosting 复制，此后独立）
company_name VARCHAR | company_intro TEXT | job_title VARCHAR
job_requirements TEXT | base VARCHAR | salary_range VARCHAR
source_platform VARCHAR | source_url VARCHAR | delivery_channel VARCHAR
applied_at TIMESTAMPTZ | current_status VARCHAR(30)
interview_summary TEXT | interview_notes TEXT
resume_match_score INT | interview_performance_score INT
role_interest_score INT | overall_score INT | remark TEXT
created_at | updated_at | deleted_at

-- interview_progress
id UUID PK | application_id UUID FK→resume_applications
stage VARCHAR(30) | occurred_at TIMESTAMPTZ | result VARCHAR(20)
interviewer_or_team VARCHAR | note TEXT
created_at | updated_at | deleted_at

-- interview_questions
id UUID PK | user_id UUID FK→users
title VARCHAR | content TEXT | answer_analysis TEXT
tags VARCHAR[]                     -- PostgreSQL 原生数组，支持 @> 查询
difficulty VARCHAR(10) | question_type VARCHAR(20)
mastery_status VARCHAR(20) DEFAULT 'unreviewed'
is_favorite BOOLEAN DEFAULT false | source VARCHAR
last_reviewed_at TIMESTAMPTZ | remark TEXT
created_at | updated_at | deleted_at

-- question_occurrences
id UUID PK | question_id UUID FK→interview_questions
application_id UUID FK→resume_applications
interview_progress_id UUID FK→interview_progress (nullable)
company_name_snapshot VARCHAR | job_title_snapshot VARCHAR
interview_stage_snapshot VARCHAR | occurred_at TIMESTAMPTZ
actual_question TEXT | answer_performance VARCHAR(20) | note TEXT
created_at | updated_at | deleted_at
UNIQUE INDEX ON (question_id, application_id, interview_progress_id) WHERE deleted_at IS NULL

-- flashcard_reviews
id UUID PK | question_id UUID FK→interview_questions
reviewed_at TIMESTAMPTZ | result VARCHAR(20) | review_mode VARCHAR(30)
created_at TIMESTAMPTZ            -- 无 updated_at，review 不可修改
deleted_at TIMESTAMPTZ
```

### 2.3 关键业务规则实现位置

| 规则 | 实现位置 |
|------|---------|
| `version_count` 维护 | Domain Service（创建/删除 version 时同步） |
| 删除 version 前检查 application 引用 | Application Service（查询存在则返回 50008） |
| interview_progress 变更同步 `currentStatus` | Domain Service：`RecalculateApplicationStatus`（从全量 progress 重算） |
| delete progress 时清空 occurrence 的 `interviewProgressId` | Application Service 编排 |
| flashcard review 同步 `masteryStatus` | Domain Entity 方法：`question.ApplyReview(result)` |
| `occurrenceCount`（计算字段） | Repository 层 COUNT 子查询，不持久化 |
| 快照字段（application 创建时） | Application Service：从 JobPosting 复制到 Application |
| application 的 currentStatus 枚举映射 | `failed→rejected`；其余 stage 1:1 映射到 status |

### 2.4 枚举值

| 枚举 | 值 |
|------|----|
| JobType | `campus`, `social`, `internship` |
| ApplicationStatus | `applied`, `viewed`, `written_test`, `first_interview`, `second_interview`, `final_interview`, `hr_interview`, `offer`, `rejected`, `closed` |
| InterviewResult | `passed`, `pending`, `failed` |
| JobPostingStatus | `pending`, `applied`, `not_fit`, `closed` |
| Difficulty | `easy`, `medium`, `hard` |
| QuestionType | `knowledge`, `scenario`, `coding`, `algorithm`, `project`, `behavior`, `discussion` |
| MasteryStatus | `unreviewed`, `weak`, `normal`, `mastered` |
| AnswerPerformance | `failed`, `vague`, `answered`, `good` |
| FlashcardReviewMode | `random`, `tag`, `difficulty`, `favorite_first`, `weak_first`, `frequent_first` |
| FlashcardReviewResult | `failed`, `vague`, `mastered` |

---

## 3. 认证机制与 API 设计

### 3.1 认证流程（JWT + Redis 黑名单）

**登录：**  
`POST /api/v1/user/login` → 验证密码（bcrypt）→ 生成 JWT（payload: userID, username, jti, exp）→ 返回 token

**请求认证：**  
Header `X-Access-Token: <token>` → AuthMiddleware 验证签名+过期 → 检查 Redis 黑名单（key: `blacklist:<jti>`）→ 注入 userID 到 Gin Context

**登出：**  
`POST /api/v1/user/logout` → 解析 token 取 jti → 写入 Redis 黑名单，TTL = token 剩余有效期

**验证码：**  
`POST /api/v1/user/email-code` → 生成6位码 → 存 Redis（key: `email_code:<email>`，TTL=300s）→ 开发环境响应中返回 mockCode  
`POST /api/v1/user/register` → 从 Redis 取码比对 → 创建用户 → 删除 Redis key

### 3.2 统一响应结构

```json
// 成功
{ "code": 20000, "message": "成功", "data": {} }

// 分页
{ "code": 20000, "message": "成功", "data": { "list": [], "total": 0, "page": 1, "pageSize": 10 } }

// 错误
{ "code": 50004, "message": "资源不存在", "data": null }
```

错误码：`20000` 成功 · `40000` 参数错误 · `40100` 未登录 · `40300` 无权限 · `50004` 不存在 · `50008` 业务冲突 · `50000` 服务器错误

### 3.3 路由列表（/api/v1 前缀）

**无需认证：**
```
POST /api/v1/user/login
POST /api/v1/user/register
POST /api/v1/user/email-code
```

**需要 JWT 认证（以下所有路由）：**
```
GET  /api/v1/user/profile
PUT  /api/v1/user/profile
POST /api/v1/user/logout

GET    /api/v1/resume-repos
POST   /api/v1/resume-repos
GET    /api/v1/resume-repos/:id
PUT    /api/v1/resume-repos/:id
DELETE /api/v1/resume-repos/:id

GET    /api/v1/resume-versions
POST   /api/v1/resume-versions
GET    /api/v1/resume-versions/:id
PUT    /api/v1/resume-versions/:id
DELETE /api/v1/resume-versions/:id
POST   /api/v1/resume-versions/:id/activate

GET    /api/v1/job-postings
POST   /api/v1/job-postings
GET    /api/v1/job-postings/:id
PUT    /api/v1/job-postings/:id
DELETE /api/v1/job-postings/:id
POST   /api/v1/job-postings/parse          # mock 实现，返回预填数据

GET    /api/v1/resume-applications
POST   /api/v1/resume-applications
GET    /api/v1/resume-applications/:id
PUT    /api/v1/resume-applications/:id
DELETE /api/v1/resume-applications/:id

GET    /api/v1/interview-progress
POST   /api/v1/interview-progress
PUT    /api/v1/interview-progress/:id
DELETE /api/v1/interview-progress/:id

GET    /api/v1/interview-questions
POST   /api/v1/interview-questions
GET    /api/v1/interview-questions/:id
PUT    /api/v1/interview-questions/:id
DELETE /api/v1/interview-questions/:id
PUT    /api/v1/interview-questions/:id/favorite

GET    /api/v1/question-occurrences
POST   /api/v1/question-occurrences
DELETE /api/v1/question-occurrences/:id

GET    /api/v1/flashcards
POST   /api/v1/flashcard-reviews
```

### 3.4 中间件栈（按执行顺序）

| 中间件 | 职责 |
|--------|------|
| Recovery | panic 恢复，返回 50000 |
| Logger (zap) | 请求日志（方法、路径、耗时、状态码） |
| CORS | 允许前端开发服务器（9527端口）跨域 |
| AuthMiddleware | 验证 JWT + 检查 Redis 黑名单 + 注入 userID |

---

## 4. 技术栈与开发规范

### 4.1 依赖库

| 用途 | 库 |
|------|----|
| HTTP 框架 | `gin-gonic/gin` |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` |
| Redis | `redis/go-redis/v9` |
| JWT | `golang-jwt/jwt/v5` |
| 配置 | `spf13/viper` |
| 密码哈希 | `golang.org/x/crypto/bcrypt`（cost=12） |
| UUID | `google/uuid` |
| DB 迁移 | `golang-migrate/migrate` |
| 测试断言 | `stretchr/testify` |
| Mock 生成 | `vektra/mockery v2` |
| 集成测试 DB | `testcontainers/testcontainers-go` |
| 日志 | `uber-go/zap` |

### 4.2 TDD 工作流

每个功能按以下循环开发：

1. **写测试**（先写失败的测试）
2. **跑测试**（确认 RED）
3. **写最小实现**（让测试通过）
4. **跑测试**（确认 GREEN）
5. **重构**（整理代码，测试保持绿）

各层测试策略：

| 层 | 测试类型 | 工具 | 重点 |
|----|---------|------|------|
| Domain | 纯单元测试 | Go testing | 实体方法、值对象、领域服务（零依赖） |
| Application | 单元测试 + Mock | testify + mockery | Mock Repository 接口，验证用例逻辑 |
| Infrastructure | 集成测试 | testcontainers-go | 真实 PostgreSQL 容器，验证 GORM 仓储 |
| Interface | Handler 测试 | httptest + testify | Mock Application Service，验证请求响应 |

### 4.3 本地开发环境

**docker-compose.yml：**
```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: ai_for
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports: ["5432:5432"]
    volumes: [postgres_data:/var/lib/postgresql/data]

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
```

**环境变量（.env）：**
```
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=ai_for
DB_USER=postgres
DB_PASSWORD=postgres

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

EMAIL_MOCK=true     # true=开发环境返回 mockCode，false=发真实邮件
```

**Makefile 常用命令：**
```
make run          # 启动后端服务
make test         # 跑所有单元测试
make test-int     # 跑集成测试（需要 Docker）
make migrate-up   # 执行数据库迁移
make migrate-down # 回滚迁移
make mock         # mockery 生成所有 Repository mock
make lint         # golangci-lint 检查
```

### 4.4 前端对接说明

前端需要做以下调整：
- `vue.config.js` 代理目标从 `http://127.0.0.1:9528` 改为 `http://127.0.0.1:8080`
- 所有 API 调用路径从 `/api/*` 改为 `/api/v1/*`（或在 devServer proxy 中 rewrite 路径前缀）

---

## 5. 资源隔离

所有资源操作通过 JWT 中注入的 `userID` 过滤，用户只能访问自己的数据。Repository 层的所有查询均带 `WHERE user_id = ?` 条件（或通过关联关系间接隔离）。
