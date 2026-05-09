# 求职流程管理系统后端 API 设计

版本：V0.1
日期：2026-05-09
适用范围：基于当前前端与 Mock API 反推的后端接口设计。

## 1. 设计约定

### 1.1 请求方法约束

后端首期只需要支持两类 HTTP 方法：

- `GET`：查询列表、查询详情、获取候选数据。
- `POST`：创建、更新、删除、状态变更、业务动作。

前端当前存在的 `PUT` / `DELETE` 语义在后端正式设计中统一改为 `POST` 动作接口。

### 1.2 基础地址

```text
/api
```

开发环境前端会通过 `/dev-api` 代理到后端，后端路由仍按 `/api/...` 设计。

### 1.3 通用请求格式

`GET` 使用 query 参数：

```http
GET /api/resume-repos?page=1&pageSize=10
```

`POST` 使用 JSON body：

```json
{
  "id": "1",
  "data": {
    "name": "前端校招简历"
  }
}
```

### 1.4 通用响应格式

成功：

```json
{
  "code": 20000,
  "message": "成功",
  "data": {}
}
```

失败：

```json
{
  "code": 50004,
  "message": "资源不存在",
  "data": null
}
```

建议错误码：

| code | 含义 |
| --- | --- |
| 20000 | 成功 |
| 40000 | 参数错误 |
| 40100 | 未登录或 token 无效 |
| 40300 | 无权限 |
| 50004 | 资源不存在 |
| 50008 | 业务规则冲突 |
| 50000 | 服务端错误 |

### 1.5 分页响应格式

```json
{
  "list": [],
  "total": 0,
  "page": 1,
  "pageSize": 10
}
```

### 1.6 时间格式

所有时间字段使用 ISO 8601 字符串：

```json
"2026-05-09T10:00:00.000Z"
```

### 1.7 认证约定

登录成功后返回 `token`。后续接口建议使用：

```http
Authorization: Bearer <token>
```

## 2. 枚举

### 2.1 JobType

| 值 | 说明 |
| --- | --- |
| `campus` | 校招 |
| `social` | 社招 |
| `internship` | 实习 |

### 2.2 ApplicationStatus / InterviewStage

| 值 | 说明 |
| --- | --- |
| `applied` | 已投递 |
| `viewed` | 已查看 |
| `written_test` | 笔试 |
| `first_interview` | 一面 |
| `second_interview` | 二面 |
| `final_interview` | 终面 |
| `hr_interview` | HR 面 |
| `offer` | Offer |
| `rejected` | 已拒绝 |
| `closed` | 已关闭 |

### 2.3 InterviewResult

| 值 | 说明 |
| --- | --- |
| `passed` | 通过 |
| `pending` | 待定 |
| `failed` | 未通过 |

### 2.4 JobPostingStatus

| 值 | 说明 |
| --- | --- |
| `pending` | 待投递 |
| `applied` | 已投递 |
| `not_fit` | 不合适 |
| `closed` | 已关闭 |

### 2.5 InterviewQuestionDifficulty

| 值 | 说明 |
| --- | --- |
| `easy` | 简单 |
| `medium` | 中等 |
| `hard` | 困难 |

### 2.6 InterviewQuestionType

| 值 | 说明 |
| --- | --- |
| `knowledge` | 知识点 |
| `scenario` | 场景题 |
| `coding` | 手写题 |
| `algorithm` | 算法题 |
| `project` | 项目题 |
| `behavior` | 行为面试 |
| `discussion` | 开放讨论 |

### 2.7 QuestionMasteryStatus

| 值 | 说明 |
| --- | --- |
| `unreviewed` | 未复习 |
| `weak` | 薄弱 |
| `normal` | 一般 |
| `mastered` | 已掌握 |

### 2.8 QuestionAnswerPerformance

| 值 | 说明 |
| --- | --- |
| `failed` | 没答上 |
| `vague` | 模糊 |
| `answered` | 基本答出 |
| `good` | 表现较好 |

### 2.9 FlashcardReviewMode

| 值 | 说明 |
| --- | --- |
| `random` | 随机复习 |
| `tag` | 标签专项 |
| `difficulty` | 难度专项 |
| `favorite_first` | 收藏优先 |
| `weak_first` | 薄弱优先 |
| `frequent_first` | 高频优先 |

### 2.10 FlashcardReviewResult

| 值 | 说明 |
| --- | --- |
| `failed` | 不会 |
| `vague` | 模糊 |
| `mastered` | 掌握 |

## 3. 数据模型

### 3.1 UserProfile

```json
{
  "id": "u1",
  "username": "demo",
  "realName": "张三",
  "age": 24,
  "email": "demo@example.com",
  "phone": "13800000000",
  "jobIntention": "前端开发工程师",
  "avatar": "",
  "educationExperiences": [
    {
      "school": "某大学",
      "education": "本科",
      "major": "计算机科学与技术",
      "admissionDate": "2020-09",
      "graduationDate": "2024-06"
    }
  ],
  "location": "上海",
  "personalAdvantage": "熟悉 Vue 和工程化",
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.2 ResumeRepo

```json
{
  "id": "1",
  "userId": "u1",
  "name": "前端校招简历",
  "jobType": "campus",
  "targetPosition": "前端开发工程师",
  "description": "校招投递使用",
  "versionCount": 3,
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.3 ResumeVersion

```json
{
  "id": "1",
  "repoId": "1",
  "version": "v1.0.0",
  "title": "腾讯前端定制版",
  "content": "# 张三\n\n## 项目经历",
  "remark": "突出 Vue 项目",
  "isActive": true,
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.4 JobPosting

```json
{
  "id": "1",
  "companyName": "腾讯科技",
  "companyIntro": "互联网公司",
  "jobTitle": "前端开发工程师",
  "jobRequirements": "熟悉 Vue、TypeScript",
  "base": "深圳",
  "salaryRange": "20k-35k",
  "sourcePlatform": "Boss",
  "sourceUrl": "https://example.com/jobs/1",
  "status": "pending",
  "remark": "重点岗位",
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.5 ResumeApplication

```json
{
  "id": "1",
  "jobPostingId": "1",
  "resumeVersionId": "1",
  "repoId": "1",
  "companyName": "腾讯科技",
  "companyIntro": "互联网公司",
  "jobTitle": "前端开发工程师",
  "jobRequirements": "熟悉 Vue、TypeScript",
  "base": "深圳",
  "salaryRange": "20k-35k",
  "sourcePlatform": "Boss",
  "sourceUrl": "https://example.com/jobs/1",
  "deliveryChannel": "官网",
  "appliedAt": "2026-05-09T10:00:00.000Z",
  "currentStatus": "first_interview",
  "interviewSummary": "一面通过",
  "interviewNotes": "重点问了 Vue 响应式",
  "resumeMatchScore": 4,
  "interviewPerformanceScore": 4,
  "roleInterestScore": 5,
  "overallScore": 4,
  "remark": "",
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.6 InterviewProgress

```json
{
  "id": "1",
  "applicationId": "1",
  "stage": "first_interview",
  "occurredAt": "2026-05-09T10:00:00.000Z",
  "result": "passed",
  "interviewerOrTeam": "前端团队",
  "note": "问了 Vue 响应式和性能优化",
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.7 InterviewQuestion

```json
{
  "id": "1",
  "title": "Vue 2 响应式原理",
  "content": "请说明 Vue 2 响应式实现方式。",
  "answerAnalysis": "核心是 Object.defineProperty。",
  "tags": ["Vue", "JavaScript"],
  "difficulty": "medium",
  "questionType": "knowledge",
  "masteryStatus": "normal",
  "isFavorite": true,
  "source": "面试复盘",
  "occurrenceCount": 2,
  "lastReviewedAt": "2026-05-09T10:00:00.000Z",
  "remark": "容易追问 nextTick",
  "createdAt": "2026-05-01T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.8 InterviewQuestionDetail

```json
{
  "id": "1",
  "title": "Vue 2 响应式原理",
  "content": "请说明 Vue 2 响应式实现方式。",
  "answerAnalysis": "核心是 Object.defineProperty。",
  "tags": ["Vue"],
  "difficulty": "medium",
  "questionType": "knowledge",
  "masteryStatus": "normal",
  "isFavorite": true,
  "occurrenceCount": 1,
  "createdAt": "2026-05-01T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z",
  "occurrences": [],
  "reviewSummary": {
    "totalReviews": 3,
    "failedCount": 1,
    "vagueCount": 1,
    "masteredCount": 1,
    "lastReviewedAt": "2026-05-09T10:00:00.000Z",
    "lastResult": "mastered"
  },
  "recentReviews": []
}
```

### 3.9 QuestionOccurrence

```json
{
  "id": "1",
  "questionId": "1",
  "applicationId": "1",
  "interviewProgressId": "1",
  "companyNameSnapshot": "腾讯科技",
  "jobTitleSnapshot": "前端开发工程师",
  "interviewStageSnapshot": "first_interview",
  "occurredAt": "2026-05-09T10:00:00.000Z",
  "actualQuestion": "为什么 Vue 2 对数组索引更新不敏感？",
  "answerPerformance": "answered",
  "note": "答到了 defineProperty",
  "createdAt": "2026-05-09T10:00:00.000Z",
  "updatedAt": "2026-05-09T10:00:00.000Z"
}
```

### 3.10 FlashcardReview

```json
{
  "id": "1",
  "questionId": "1",
  "reviewedAt": "2026-05-09T10:00:00.000Z",
  "result": "mastered",
  "reviewMode": "tag",
  "createdAt": "2026-05-09T10:00:00.000Z"
}
```

## 4. 用户与认证 API

### 4.1 发送邮箱验证码

```http
POST /api/user/email-code
```

请求：

```json
{
  "email": "demo@example.com"
}
```

响应 `data`：

```json
{
  "email": "demo@example.com",
  "mockCode": "123456",
  "expiresIn": 300
}
```

### 4.2 注册

```http
POST /api/user/register
```

请求：

```json
{
  "username": "demo",
  "password": "123456",
  "email": "demo@example.com",
  "emailCode": "123456"
}
```

响应 `data`：`LoginResponse`

```json
{
  "token": "token-demo",
  "userInfo": {
    "id": "u1",
    "username": "demo",
    "email": "demo@example.com",
    "createdAt": "2026-05-09T10:00:00.000Z",
    "updatedAt": "2026-05-09T10:00:00.000Z"
  }
}
```

### 4.3 登录

```http
POST /api/user/login
```

请求：

```json
{
  "username": "demo",
  "password": "123456"
}
```

响应 `data`：`LoginResponse`

### 4.4 登出

```http
POST /api/user/logout
```

请求：

```json
{}
```

响应 `data`：

```json
null
```

### 4.5 获取个人信息

```http
GET /api/user/profile
```

响应 `data`：`UserProfile`

### 4.6 更新个人信息

```http
POST /api/user/profile/update
```

请求：

```json
{
  "realName": "张三",
  "age": 24,
  "email": "demo@example.com",
  "phone": "13800000000",
  "jobIntention": "前端开发工程师",
  "avatar": "",
  "educationExperiences": [],
  "location": "上海",
  "personalAdvantage": "熟悉 Vue 和工程化"
}
```

响应 `data`：`UserProfile`

## 5. 简历仓库 API

### 5.1 查询仓库列表

```http
GET /api/resume-repos?jobType=campus&keyword=前端&page=1&pageSize=10
```

query：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `userId` | string | 否 | 用户 ID，通常由 token 推导 |
| `jobType` | JobType | 否 | 求职类型 |
| `keyword` | string | 否 | 仓库名称或目标岗位 |
| `page` | number | 否 | 页码 |
| `pageSize` | number | 否 | 每页数量 |

响应 `data`：`ResumeRepoListResponse`

### 5.2 查询仓库详情

```http
GET /api/resume-repos/detail?id=1
```

响应 `data`：`ResumeRepo`

### 5.3 创建仓库

```http
POST /api/resume-repos/create
```

请求：

```json
{
  "name": "前端校招简历",
  "jobType": "campus",
  "targetPosition": "前端开发工程师",
  "description": "校招投递使用"
}
```

响应 `data`：`ResumeRepo`

### 5.4 更新仓库

```http
POST /api/resume-repos/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "name": "前端社招简历",
    "jobType": "social",
    "targetPosition": "高级前端开发工程师",
    "description": "社招投递使用"
  }
}
```

响应 `data`：`ResumeRepo`

### 5.5 删除仓库

```http
POST /api/resume-repos/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

## 6. 简历版本 API

### 6.1 查询版本列表

```http
GET /api/resume-versions?repoId=1&keyword=腾讯&page=1&pageSize=10
```

query：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `repoId` | string | 是 | 仓库 ID |
| `keyword` | string | 否 | 版本号或标题 |
| `page` | number | 否 | 页码 |
| `pageSize` | number | 否 | 每页数量 |

响应 `data`：`ResumeVersionListResponse`

### 6.2 查询版本详情

```http
GET /api/resume-versions/detail?id=1
```

响应 `data`：`ResumeVersion`

### 6.3 创建版本

```http
POST /api/resume-versions/create
```

请求：

```json
{
  "repoId": "1",
  "title": "腾讯前端定制版",
  "content": "# 张三\n\n## 项目经历",
  "remark": "突出 Vue 项目"
}
```

响应 `data`：`ResumeVersion`

### 6.4 更新版本

```http
POST /api/resume-versions/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "title": "字节前端定制版",
    "content": "# 张三\n\n## 项目经历",
    "remark": "突出性能优化"
  }
}
```

响应 `data`：`ResumeVersion`

### 6.5 删除版本

```http
POST /api/resume-versions/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

### 6.6 激活版本

```http
POST /api/resume-versions/activate
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：`ResumeVersion`

### 6.7 版本对比

```http
GET /api/resume-versions/compare?oldVersionId=1&newVersionId=2
```

响应 `data`：

```json
{
  "oldVersion": {},
  "newVersion": {},
  "diff": "差异内容"
}
```

## 7. 岗位库 API

### 7.1 查询岗位列表

```http
GET /api/job-postings?keyword=前端&sourcePlatform=Boss&status=pending&page=1&pageSize=10
```

响应 `data`：`JobPostingListResponse`

### 7.2 查询岗位详情

```http
GET /api/job-postings/detail?id=1
```

响应 `data`：`JobPosting`

### 7.3 创建岗位

```http
POST /api/job-postings/create
```

请求：

```json
{
  "companyName": "腾讯科技",
  "companyIntro": "互联网公司",
  "jobTitle": "前端开发工程师",
  "jobRequirements": "熟悉 Vue、TypeScript",
  "base": "深圳",
  "salaryRange": "20k-35k",
  "sourcePlatform": "Boss",
  "sourceUrl": "https://example.com/jobs/1",
  "status": "pending",
  "remark": "重点岗位"
}
```

响应 `data`：`JobPosting`

### 7.4 更新岗位

```http
POST /api/job-postings/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "status": "applied",
    "remark": "已投递"
  }
}
```

响应 `data`：`JobPosting`

### 7.5 删除岗位

```http
POST /api/job-postings/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

### 7.6 解析岗位链接

```http
POST /api/job-postings/parse
```

请求：

```json
{
  "url": "https://example.com/jobs/1"
}
```

响应 `data`：`ParseJobPostingResponse`

```json
{
  "companyName": "腾讯科技",
  "companyIntro": "互联网公司",
  "jobTitle": "前端开发工程师",
  "jobRequirements": "熟悉 Vue、TypeScript",
  "base": "深圳",
  "salaryRange": "20k-35k",
  "sourcePlatform": "Boss",
  "sourceUrl": "https://example.com/jobs/1"
}
```

## 8. 投递记录 API

### 8.1 查询投递列表

```http
GET /api/resume-applications?resumeVersionId=1&repoId=1&keyword=腾讯&currentStatus=first_interview&page=1&pageSize=10
```

响应 `data`：`ResumeApplicationListResponse`

### 8.2 查询投递详情

```http
GET /api/resume-applications/detail?id=1
```

响应 `data`：`ResumeApplication`

### 8.3 创建投递

```http
POST /api/resume-applications/create
```

请求：

```json
{
  "jobPostingId": "1",
  "resumeVersionId": "1",
  "repoId": "1",
  "companyName": "腾讯科技",
  "companyIntro": "互联网公司",
  "jobTitle": "前端开发工程师",
  "jobRequirements": "熟悉 Vue、TypeScript",
  "base": "深圳",
  "salaryRange": "20k-35k",
  "sourcePlatform": "Boss",
  "sourceUrl": "https://example.com/jobs/1",
  "deliveryChannel": "官网",
  "appliedAt": "2026-05-09T10:00:00.000Z",
  "currentStatus": "applied",
  "interviewSummary": "",
  "interviewNotes": "",
  "resumeMatchScore": 4,
  "interviewPerformanceScore": 0,
  "roleInterestScore": 5,
  "overallScore": 4,
  "remark": ""
}
```

响应 `data`：`ResumeApplication`

### 8.4 更新投递

```http
POST /api/resume-applications/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "currentStatus": "first_interview",
    "interviewNotes": "问了 Vue 响应式",
    "overallScore": 4
  }
}
```

响应 `data`：`ResumeApplication`

### 8.5 删除投递

```http
POST /api/resume-applications/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

## 9. 面试进展 API

### 9.1 查询面试进展列表

```http
GET /api/interview-progress?applicationId=1
```

响应 `data`：

```json
{
  "list": [
    {
      "id": "1",
      "applicationId": "1",
      "stage": "first_interview",
      "occurredAt": "2026-05-09T10:00:00.000Z",
      "result": "passed",
      "interviewerOrTeam": "前端团队",
      "note": "问了 Vue 响应式",
      "createdAt": "2026-05-09T10:00:00.000Z",
      "updatedAt": "2026-05-09T10:00:00.000Z"
    }
  ]
}
```

### 9.2 创建面试进展

```http
POST /api/interview-progress/create
```

请求：

```json
{
  "applicationId": "1",
  "stage": "first_interview",
  "occurredAt": "2026-05-09T10:00:00.000Z",
  "result": "passed",
  "interviewerOrTeam": "前端团队",
  "note": "问了 Vue 响应式"
}
```

响应 `data`：`InterviewProgress`

### 9.3 更新面试进展

```http
POST /api/interview-progress/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "result": "pending",
    "note": "等待反馈"
  }
}
```

响应 `data`：`InterviewProgress`

### 9.4 删除面试进展

```http
POST /api/interview-progress/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

## 10. 面试题库 API

### 10.1 查询题目列表

```http
GET /api/interview-questions?keyword=Vue&tags=Vue,JavaScript&difficulty=medium&questionType=knowledge&masteryStatus=normal&isFavorite=true&page=1&pageSize=10
```

query：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `keyword` | string | 否 | 搜索标题、内容、答案、标签 |
| `tags` | string | 否 | 逗号分隔 |
| `difficulty` | InterviewQuestionDifficulty | 否 | 难度 |
| `questionType` | InterviewQuestionType | 否 | 题型 |
| `masteryStatus` | QuestionMasteryStatus | 否 | 掌握状态 |
| `isFavorite` | boolean | 否 | 是否收藏 |
| `page` | number | 否 | 页码 |
| `pageSize` | number | 否 | 每页数量 |

响应 `data`：`InterviewQuestionListResponse`

### 10.2 查询题目详情

```http
GET /api/interview-questions/detail?id=1
```

响应 `data`：`InterviewQuestionDetail`

### 10.3 创建题目

```http
POST /api/interview-questions/create
```

请求：

```json
{
  "title": "Vue 2 响应式原理",
  "content": "请说明 Vue 2 响应式实现方式。",
  "answerAnalysis": "核心是 Object.defineProperty。",
  "tags": ["Vue", "JavaScript"],
  "difficulty": "medium",
  "questionType": "knowledge",
  "masteryStatus": "unreviewed",
  "isFavorite": false,
  "source": "面试复盘",
  "remark": "容易追问 nextTick"
}
```

响应 `data`：`InterviewQuestion`

### 10.4 更新题目

```http
POST /api/interview-questions/update
```

请求：

```json
{
  "id": "1",
  "data": {
    "title": "Vue 2 响应式原理和数组更新",
    "content": "请说明 Vue 2 响应式实现方式。",
    "answerAnalysis": "核心是 Object.defineProperty。",
    "tags": ["Vue"],
    "difficulty": "medium",
    "questionType": "knowledge",
    "masteryStatus": "normal",
    "isFavorite": true,
    "source": "面试复盘",
    "remark": "重点复习"
  }
}
```

响应 `data`：`InterviewQuestion`

### 10.5 删除题目

```http
POST /api/interview-questions/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

后端应同步处理该题目的面试提问记录和闪卡复习记录，避免悬空数据。

### 10.6 更新收藏状态

```http
POST /api/interview-questions/favorite
```

请求：

```json
{
  "id": "1",
  "isFavorite": true
}
```

响应 `data`：`InterviewQuestion`

## 11. 面试提问记录 API

### 11.1 查询面试提问记录

```http
GET /api/question-occurrences?questionId=1&applicationId=1
```

query：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `questionId` | string | 否 | 题目 ID |
| `applicationId` | string | 否 | 投递记录 ID |

响应 `data`：

```json
{
  "list": []
}
```

### 11.2 创建面试提问记录

```http
POST /api/question-occurrences/create
```

请求：

```json
{
  "questionId": "1",
  "applicationId": "1",
  "interviewProgressId": "1",
  "interviewStageSnapshot": "first_interview",
  "occurredAt": "2026-05-09T10:00:00.000Z",
  "actualQuestion": "为什么 Vue 2 对数组索引更新不敏感？",
  "answerPerformance": "answered",
  "note": "答到了 defineProperty"
}
```

响应 `data`：`QuestionOccurrence`

后端创建时应根据 `applicationId` 写入 `companyNameSnapshot`、`jobTitleSnapshot`。

### 11.3 删除面试提问记录

```http
POST /api/question-occurrences/delete
```

请求：

```json
{
  "id": "1"
}
```

响应 `data`：

```json
null
```

## 12. 闪卡复习 API

### 12.1 获取闪卡队列

```http
GET /api/flashcards?mode=tag&tags=Vue&difficulty=medium&limit=10
```

query：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `mode` | FlashcardReviewMode | 是 | 抽题模式 |
| `tags` | string | 否 | 标签专项时使用，逗号分隔 |
| `difficulty` | InterviewQuestionDifficulty | 否 | 难度专项时使用 |
| `limit` | number | 否 | 抽题数量 |

响应 `data`：

```json
{
  "list": [],
  "total": 0,
  "mode": "tag"
}
```

### 12.2 保存闪卡自评

```http
POST /api/flashcard-reviews/create
```

请求：

```json
{
  "questionId": "1",
  "reviewedAt": "2026-05-09T10:00:00.000Z",
  "result": "mastered",
  "reviewMode": "tag"
}
```

响应 `data`：

```json
{
  "review": {
    "id": "1",
    "questionId": "1",
    "reviewedAt": "2026-05-09T10:00:00.000Z",
    "result": "mastered",
    "reviewMode": "tag",
    "createdAt": "2026-05-09T10:00:00.000Z"
  },
  "updatedQuestion": {},
  "reviewSummary": {
    "totalReviews": 1,
    "failedCount": 0,
    "vagueCount": 0,
    "masteredCount": 1,
    "lastReviewedAt": "2026-05-09T10:00:00.000Z",
    "lastResult": "mastered"
  }
}
```

后端保存自评后应同步更新题目的 `masteryStatus` 和 `lastReviewedAt`。

## 13. 后端实现建议

### 13.1 表结构建议

建议至少包含以下表：

- `users`
- `user_profiles`
- `resume_repos`
- `resume_versions`
- `job_postings`
- `resume_applications`
- `interview_progress`
- `interview_questions`
- `question_occurrences`
- `flashcard_reviews`

### 13.2 关键业务规则

- 删除简历版本前，应检查是否有关联投递记录；如有关联，拒绝删除或要求迁移。
- 创建投递记录时，如果传入 `jobPostingId`，应保存岗位快照字段，后续岗位库修改不影响历史投递。
- 创建或更新面试进展后，应同步投递记录 `currentStatus`。
- 删除投递记录时，应级联删除或清理关联的面试进展和面试提问记录。
- 删除面试进展时，应清理面试提问记录中的 `interviewProgressId` 或按业务需要级联删除。
- 删除面试题时，应同步删除面试提问记录和闪卡复习记录。
- 创建面试提问记录时，应保存公司、岗位和面试阶段快照。
- 保存闪卡复习结果时，应更新题目的掌握状态和最近复习时间。

### 13.3 GET/POST 与前端适配说明

当前前端代码仍使用部分 REST 风格接口：

- `PUT /api/.../:id`
- `DELETE /api/.../:id`

后端若只提供本文档中的 `GET` / `POST` 接口，前端后续需要做一层 API 适配，把更新和删除调用切换为动作型 `POST` 接口，例如：

- 更新：`POST /api/resume-repos/update`
- 删除：`POST /api/resume-repos/delete`
- 收藏：`POST /api/interview-questions/favorite`

首期也可以在后端同时兼容当前 Mock 路由和本文档动作路由，降低联调成本。
