# 面试题库后端设计

## 背景

前端已经完成面试题库、题目出现记录和闪卡复习页面，并通过以下 API wrapper 调用后端：

- `frontend/src/api/interview-questions.ts`
- `frontend/src/api/question-occurrences.ts`
- `frontend/src/api/flashcards.ts`

后端需要按现有 Go + Gin + GORM + PostgreSQL 的 DDD 分层实现题库后端。接口路径遵循后端现有版本化约定，只注册 `/api/v1/...`；请求字段、响应 `data` 形状、枚举值以前端模型为准。

## 目标

- 支持题目独立存在，不要求绑定投递或面试记录。
- 支持题目的创建、编辑、删除、收藏、筛选、分页和详情聚合。
- 支持把题目关联到投递或具体面试进展，记录真实面试出现情况。
- 支持闪卡队列生成和每一次复习记录。
- 支持复习后更新题目的当前掌握状态和最近复习时间。
- 所有数据按登录用户隔离。

## 非目标

- 不在本次后端设计中实现前端 UI 调整。
- 不实现 `/api/...` 兼容路由；真实联调时前端 wrapper 需要切到 `/api/v1/...`。
- 不新增单独统计表；复习统计和出现次数由明细表聚合。

## 领域模型

### InterviewQuestion

题目主资源，可以独立存在。

字段对齐前端 `InterviewQuestion`：

- `id`
- `title`
- `content`
- `answerAnalysis`
- `tags`
- `difficulty`
- `questionType`
- `masteryStatus`
- `isFavorite`
- `source`
- `occurrenceCount`
- `lastReviewedAt`
- `remark`
- `createdAt`
- `updatedAt`

后端不持久化 `occurrenceCount`，由 `question_occurrences` 聚合得到。

### QuestionOccurrence

出现记录表示“某道题在某个投递或某轮面试中出现过”。它是题目的可选关联记录。

规则：

- 题目可以有 `0..N` 条 occurrence。
- 创建题目不需要 `applicationId` 或 `interviewProgressId`。
- `questionId` 和 `applicationId` 必填。
- `interviewProgressId` 可空。
- 删除 occurrence 不删除题目。
- 删除题目时级联软删除它的 occurrence 和 review。

### FlashcardReview

复习记录表保存每一次闪卡复习，而不是只存最后一次。

规则：

- 每次提交复习结果都新增一条 review。
- 同步更新题目的 `masteryStatus` 和 `lastReviewedAt`。
- `reviewSummary` 和 `recentReviews` 由 `flashcard_reviews` 查询或聚合得到。

## 枚举

### difficulty

- `easy`
- `medium`
- `hard`

### questionType

- `knowledge`
- `scenario`
- `coding`
- `algorithm`
- `project`
- `behavior`
- `discussion`

### masteryStatus

- `unreviewed`
- `weak`
- `normal`
- `mastered`

### answerPerformance

- `failed`
- `vague`
- `answered`
- `good`

### review result

- `failed`
- `vague`
- `mastered`

### reviewMode

- `random`
- `tag`
- `difficulty`
- `favorite_first`
- `weak_first`
- `frequent_first`

复习结果更新掌握状态：

- `failed` -> `weak`
- `vague` -> `normal`
- `mastered` -> `mastered`

## 数据库设计

新增 migration：`000007_create_interview_questions.up.sql` 和 `000007_create_interview_questions.down.sql`。

### interview_questions

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `UUID` | 主键 |
| `user_id` | `VARCHAR(36)` | 所属用户 |
| `title` | `VARCHAR(200)` | 题目标题 |
| `content` | `TEXT` | 题目内容 |
| `answer_analysis` | `TEXT` | 答案解析 |
| `tags` | `JSONB` | 标签数组，默认 `[]` |
| `difficulty` | `VARCHAR(20)` | 难度 |
| `question_type` | `VARCHAR(30)` | 题目类型 |
| `mastery_status` | `VARCHAR(30)` | 掌握状态 |
| `is_favorite` | `BOOLEAN` | 是否收藏 |
| `source` | `VARCHAR(100)` | 来源，可空 |
| `last_reviewed_at` | `TIMESTAMPTZ` | 最近复习时间，可空 |
| `remark` | `TEXT` | 备注，可空 |
| `created_at` | `TIMESTAMPTZ` | 创建时间 |
| `updated_at` | `TIMESTAMPTZ` | 更新时间 |
| `deleted_at` | `TIMESTAMPTZ` | 软删除 |

索引：

- `(user_id, deleted_at)`
- `(user_id, difficulty)`
- `(user_id, question_type)`
- `(user_id, mastery_status)`
- `(user_id, is_favorite)`
- `GIN(tags)`

### question_occurrences

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `UUID` | 主键 |
| `user_id` | `VARCHAR(36)` | 所属用户 |
| `question_id` | `UUID` | 关联题目 |
| `application_id` | `UUID` | 关联投递 |
| `interview_progress_id` | `UUID` | 关联具体面试进展，可空 |
| `company_name_snapshot` | `VARCHAR(200)` | 公司名快照 |
| `job_title_snapshot` | `VARCHAR(200)` | 岗位名快照 |
| `interview_stage_snapshot` | `VARCHAR(50)` | 面试阶段快照，可空 |
| `occurred_at` | `TIMESTAMPTZ` | 出现时间 |
| `actual_question` | `TEXT` | 实际问法，可空 |
| `answer_performance` | `VARCHAR(30)` | 回答表现，可空 |
| `note` | `TEXT` | 备注，可空 |
| `created_at` | `TIMESTAMPTZ` | 创建时间 |
| `updated_at` | `TIMESTAMPTZ` | 更新时间 |
| `deleted_at` | `TIMESTAMPTZ` | 软删除 |

约束：

- `question_id` 外键到 `interview_questions.id`。
- `application_id` 外键到 `resume_applications.id`。
- `interview_progress_id` 外键到 `interview_progress.id`，可空。
- 当 `interview_progress_id` 非空时，同一用户下 `(question_id, application_id, interview_progress_id)` 排除软删除后唯一。

### flashcard_reviews

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `UUID` | 主键 |
| `user_id` | `VARCHAR(36)` | 所属用户 |
| `question_id` | `UUID` | 关联题目 |
| `reviewed_at` | `TIMESTAMPTZ` | 复习时间 |
| `result` | `VARCHAR(30)` | 复习结果 |
| `review_mode` | `VARCHAR(30)` | 复习模式 |
| `created_at` | `TIMESTAMPTZ` | 创建时间 |
| `deleted_at` | `TIMESTAMPTZ` | 软删除 |

约束：

- `question_id` 外键到 `interview_questions.id`。
- 所有查询按 `user_id` 隔离。

## API 设计

所有响应沿用后端 envelope：

```json
{
  "code": 20000,
  "message": "成功",
  "data": {}
}
```

### GET /api/v1/interview-questions

查询参数：

- `keyword`
- `tags`
- `difficulty`
- `questionType`
- `masteryStatus`
- `isFavorite`
- `page`
- `pageSize`

行为：

- `keyword` 匹配 `title`、`content`、`answer_analysis` 和 `tags`。
- `tags` 为逗号分隔字符串，命中任一标签即可。
- 默认按 `updated_at DESC` 排序。
- 返回每道题的 `occurrenceCount`。

`data`：

```json
{
  "list": [],
  "total": 0,
  "page": 1,
  "pageSize": 10
}
```

### GET /api/v1/interview-questions/:id

返回 `InterviewQuestionDetail`：

- 题目基础字段
- `occurrences`
- `reviewSummary`
- `recentReviews`

独立题目无出现记录时返回：

```json
{
  "occurrenceCount": 0,
  "occurrences": [],
  "reviewSummary": {
    "totalReviews": 0,
    "failedCount": 0,
    "vagueCount": 0,
    "masteredCount": 0,
    "lastReviewedAt": null,
    "lastResult": null
  },
  "recentReviews": []
}
```

### POST /api/v1/interview-questions

请求体对齐 `CreateInterviewQuestionRequest`：

- `title`
- `content`
- `answerAnalysis`
- `tags`
- `difficulty`
- `questionType`
- `masteryStatus`
- `isFavorite`
- `source`
- `remark`

返回创建后的 `InterviewQuestion`。

### PUT /api/v1/interview-questions/:id

请求体是创建字段的 partial。返回更新后的 `InterviewQuestion`。

### DELETE /api/v1/interview-questions/:id

软删除题目，并软删除该题目的 occurrence 和 review。返回 `null`。

### PUT /api/v1/interview-questions/:id/favorite

请求体：

```json
{
  "isFavorite": true
}
```

返回更新后的 `InterviewQuestion`。

### GET /api/v1/question-occurrences

查询参数：

- `questionId`
- `applicationId`

返回：

```json
{
  "list": []
}
```

### POST /api/v1/question-occurrences

请求体对齐 `CreateQuestionOccurrenceRequest`：

- `questionId`
- `applicationId`
- `interviewProgressId`
- `interviewStageSnapshot`
- `occurredAt`
- `actualQuestion`
- `answerPerformance`
- `note`

行为：

- 校验 question 属于当前用户。
- 校验 application 属于当前用户。
- 如果传入 `interviewProgressId`，校验 progress 属于该 application。
- `companyNameSnapshot` 和 `jobTitleSnapshot` 从 application 当前值写入。
- `interviewStageSnapshot` 优先使用请求值；缺省时从 progress 阶段补齐。

### DELETE /api/v1/question-occurrences/:id

只软删除 occurrence，不删除题目。

### GET /api/v1/flashcards

查询参数：

- `mode`
- `tags`
- `difficulty`
- `limit`

返回：

```json
{
  "list": [],
  "total": 0,
  "mode": "random"
}
```

模式：

- `random`：随机。
- `tag`：按标签过滤后随机。
- `difficulty`：按难度过滤，更新时间倒序。
- `favorite_first`：收藏优先，再按出现次数降序。
- `weak_first`：`unreviewed -> weak -> normal -> mastered` 优先。
- `frequent_first`：出现次数降序，再收藏优先。

### POST /api/v1/flashcard-reviews

请求体对齐 `CreateFlashcardReviewRequest`：

- `questionId`
- `reviewedAt`
- `result`
- `reviewMode`

行为：

- 校验 question 属于当前用户。
- `reviewedAt` 缺省时使用当前时间。
- 新增一条 `flashcard_reviews`。
- 更新 question 的 `masteryStatus` 和 `lastReviewedAt`。

返回：

```json
{
  "review": {},
  "updatedQuestion": {},
  "reviewSummary": {}
}
```

## 后端分层

新增领域目录：

- `backend/internal/domain/interview_question/`
- `backend/internal/domain/question_occurrence/`
- `backend/internal/domain/flashcard_review/`

新增仓储接口：

- `backend/internal/domain/interview_question/repository/repository.go`
- `backend/internal/domain/question_occurrence/repository/repository.go`
- `backend/internal/domain/flashcard_review/repository/repository.go`

新增应用服务：

- `backend/internal/application/interview_question/interview_question_service.go`

新增 DTO：

- `backend/internal/application/dto/interview_question_dto.go`

新增 PostgreSQL 仓储：

- `backend/internal/infrastructure/persistence/postgres/interview_question_repo.go`
- `backend/internal/infrastructure/persistence/postgres/question_occurrence_repo.go`
- `backend/internal/infrastructure/persistence/postgres/flashcard_review_repo.go`

新增 handler：

- `backend/internal/interface/handler/interview_question_handler.go`

更新：

- `backend/internal/interface/router/router.go`
- `backend/cmd/server/main.go`

应用服务统一编排题库业务，对 handler 暴露题目、出现记录和闪卡复习方法。底层仓储按三张表拆分，保持边界清晰。

## 权限与所有权

所有方法从 Gin context 获取 `userID` 并向下传递。

题目相关：

- 列表只返回当前用户题目。
- 详情、更新、删除、收藏都通过 `id + user_id` 查询。

出现记录相关：

- 创建 occurrence 前校验 question 属于当前用户。
- 创建 occurrence 前校验 application 属于当前用户。
- 如果传入 progress，校验 progress 属于该 application。
- 删除 occurrence 通过 `id + user_id` 查询。

复习相关：

- 创建 review 前校验 question 属于当前用户。
- 查询 summary 和 recent reviews 时同时限制 `question_id + user_id`。

## 错误处理

沿用现有错误码：

- 参数缺失或枚举非法：`40000`
- 未登录：`40100`
- 非本人资源或资源不存在：`50004`
- 重复 occurrence：`50008`
- 内部错误：`50000`

## 测试策略

### Domain tests

- 创建题目必填校验。
- 枚举校验。
- partial update 只更新传入字段。
- review result 到 mastery status 的映射。

### Application service tests

- 独立创建题目不需要 application。
- 创建 occurrence 校验 question/application/progress 所有权。
- occurrence 快照来自 application。
- 同一 question/application/progress 重复创建返回 conflict。
- 创建 flashcard review 新增流水并更新题目状态。
- detail 返回 `occurrences`、`reviewSummary`、`recentReviews`。

### Handler tests

- 验证所有接口注册在 `/api/v1`。
- 验证 JSON 字段为前端 camelCase。
- 验证 list/detail/review response shape 与前端模型一致。

### Repository and migration tests

- migration up/down。
- JSONB tags 读写。
- `occurrenceCount` 聚合。
- 软删除后列表、统计、详情不包含对应记录。

## 前端联调影响

由于后端只注册 `/api/v1`，真实联调时需要更新以下前端 API wrapper 的 URL：

- `frontend/src/api/interview-questions.ts`
- `frontend/src/api/question-occurrences.ts`
- `frontend/src/api/flashcards.ts`

字段名、枚举值和响应 `data` 结构保持不变。
