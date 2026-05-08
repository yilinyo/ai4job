## 背景

现有项目中已经存在以下基础：

- `src/models/user.ts` 定义了 `UserProfile` 和 `UpdateProfileRequest`。
- `src/views/user/profile.vue` 已有个人信息表单雏形，支持姓名、年龄、邮箱、电话、求职意向等字段。
- `mock/users.ts` 中的 `getUserProfile` 和 `updateUserProfile` 已通过 `X-Access-Token` 解析当前 username，并按账号读取 / 更新 Mock 用户数据。
- `src/views/resume-version/editor.vue` 在创建简历版本时默认 `content` 为空，尚未根据用户资料生成草稿。
- 当前路由里既有顶层模板个人中心 `/profile/index`，也有业务个人信息页 `/resume/profile`，信息架构上容易混淆。

本次变更应把业务个人信息明确为账号级资料入口，并服务于简历创建自动填充。

## 目标 / 非目标

**目标：**

- 个人信息页以结构化表单维护，不使用 Markdown。
- 个人信息页作为一级菜单入口，与“引导页 / Guide”同级显示。
- 移除“简历管理”下的个人信息入口。
- 扩展个人信息字段，覆盖简历草稿所需的基础资料，并支持多段教育经历。
- 使用 Mock 数据模拟“每个账号独立保存个人资料”。
- 创建简历版本时自动生成默认简历内容草稿，且用户可以继续编辑修改。
- 编辑已有简历版本时不自动覆盖已有内容。

**非目标：**

- 不改造简历编辑器为完整结构化简历表单。
- 不新增真实头像上传服务。
- 不新增后端数据库或真实账号系统。
- 不支持多个个人信息版本。
- 不自动根据岗位 JD 生成定制简历。

## 关键设计决策

个人信息页使用表单，而不是 Markdown。
原因：个人信息是结构化账号资料，需要可校验、可复用、可按字段写入后端。Markdown 不利于后续自动填充、字段校验和账号资料复用。
备选方案：
- 用 Markdown 维护个人信息：灵活但难以结构化复用，不符合本次诉求。

个人信息作为顶层导航入口。
原因：个人信息属于账号级资产，不依附于某个简历仓库。放在“简历管理”下会弱化账号资料的独立性，也会导致用户误以为它只服务于简历模块。
备选方案：
- 保留在“简历管理”下：改动小，但信息架构不清晰。

复用现有业务个人信息页面并调整入口。
原因：`src/views/user/profile.vue` 已经具备表单、加载、保存、校验等基础能力，适合在原有基础上扩展字段和样式。
备选方案：
- 新建另一个页面：会造成重复页面和后续维护成本。

创建简历版本时自动填充一次。
原因：用户希望创建时减少重复输入，但生成后仍要可编辑。自动填充应只在新建态且内容为空时发生，避免覆盖用户已输入内容或编辑已有版本。
备选方案：
- 提供“从个人信息填充”按钮：可控性更强，但用户明确希望自动填入。
- 每次进入创建页都覆盖：风险高，会丢失用户临时编辑内容。

Mock 按 token 对应账号隔离。
原因：当前 Mock 已通过 token 推导 username，继续沿用该机制即可模拟后端“当前账号只能读写自己的资料”。
备选方案：
- 前端传 `userId`：不符合真实后端安全模型，容易被篡改。

## 数据模型设计

### 用户个人信息

建议扩展 `UserProfile`：

```ts
interface UserProfile {
  id: string
  username: string
  realName?: string
  age?: number
  email?: string
  phone?: string
  jobIntention?: string
  avatar?: string
  educationExperiences?: EducationExperience[]
  location?: string
  personalAdvantage?: string
  createdAt: string
  updatedAt: string
}

interface EducationExperience {
  school?: string
  education?: string
  major?: string
  admissionDate?: string
  graduationDate?: string
}
```

### 更新个人信息请求

`UpdateProfileRequest` 同步支持上述可编辑字段：

```ts
interface UpdateProfileRequest {
  realName?: string
  age?: number
  email?: string
  phone?: string
  jobIntention?: string
  avatar?: string
  educationExperiences?: EducationExperience[]
  location?: string
  personalAdvantage?: string
}
```

## 页面与交互设计

### 顶层个人信息入口

导航要求：

- 左侧一级菜单显示“个人信息”。
- 与“引导页 / Guide”同级。
- 不再作为“简历管理”的二级菜单。
- 建议业务路由使用 `/user/profile`，避免与模板自带 `/profile/index` 混淆。

页面结构建议：

1. 基础信息区：头像、姓名、年龄、邮箱、电话、所在地。
2. 求职意向区：求职意向、个人优势。
3. 教育经历区：支持添加多段学校经历，每段包含学校、学历、专业、入学时间、毕业时间。

字段建议：

- 姓名：文本输入
- 年龄：数字输入
- 邮箱：邮箱格式校验
- 电话：手机号格式校验
- 求职意向：文本或 textarea
- 头像：URL 输入或沿用现有头像字段展示
- 教育经历：可添加多段，每段包含学校、学历、专业、入学时间、毕业时间
- 学历唯一约束：同一种学历只允许添加一段经历，例如只能有一段“本科”和一段“硕士”
- 入学时间 / 毕业时间：使用月份选择，毕业时间不能早于入学时间
- 所在地：文本输入
- 个人优势：textarea

### 简历创建自动填充

触发条件：

- 当前页面是创建简历版本，不是编辑已有版本。
- 当前 `versionForm.content` 为空。
- 成功获取当前账号个人信息。

生成规则：

- 使用个人信息字段生成默认简历内容草稿。
- 空字段不强行展示或展示为可填写占位项，建议优先跳过空字段，避免生成大量空内容。
- 草稿生成后写入 `versionForm.content`，用户可继续在编辑器中修改。
- 如果用户已经手动输入内容，不自动覆盖。

默认草稿建议结构：

```md
# 姓名 | 求职意向

## 基本信息
- 年龄：年龄
- 电话：电话
- 邮箱：邮箱
- 所在地：所在地

## 教育背景
- 学历：学校 / 专业 / 入学时间 - 毕业时间

## 个人优势
个人优势
```

说明：个人信息页本身不使用 Markdown；该模板只是写入当前简历版本编辑器的默认内容草稿，符合现有简历内容模型。

## API 与 Mock 设计

沿用现有接口：

- `GET /api/user/profile`
- `PUT /api/user/profile`

Mock 要求：

- `getUserProfile` 必须基于 `X-Access-Token` 解析当前账号。
- `updateUserProfile` 只能更新当前 token 对应账号的数据。
- `admin` 和 `user` 两个 Mock 账号应有不同的个人资料样例，便于验证账号隔离。

## 实现方案

建议按以下顺序实现：

1. 扩展 `src/models/user.ts` 中的用户资料类型。
2. 扩展 `mock/users.ts` 中的用户样例和更新字段。
3. 扩展 `src/views/user/profile.vue` 表单字段、初始化逻辑和校验。
4. 调整路由：新增顶层个人信息入口，移除“简历管理”下的个人信息子路由。
5. 更新语言包或路由标题，确保左侧菜单显示“个人信息”。
6. 在简历版本创建页加载当前用户资料，并在创建态自动生成默认草稿。
7. 抽取简历草稿生成函数，避免模板拼接散落在页面方法中。

## 风险与权衡

- [导航重复或混淆] -> 缓解：移除 `/resume/profile` 入口，并避免与模板 `/profile/index` 语义冲突。
- [自动填充覆盖用户输入] -> 缓解：仅在创建态且内容为空时填充。
- [账号隔离验证不足] -> 缓解：Mock 中使用不同账号资料，并验证切换 token 后读写互不影响。
- [草稿内容格式与用户期望不一致] -> 缓解：生成结果只是初始草稿，用户可继续编辑；后续如需结构化简历编辑器应另起 change。
- [头像上传范围膨胀] -> 缓解：本次只保留头像 URL 字段或现有头像展示，不实现上传服务。

## 迁移计划

该变更不涉及真实数据迁移。Mock 中为已有账号补充新增字段。已有前端用户资料字段继续兼容，新增字段为空时不影响页面加载和简历创建。

## Open Questions

- 顶层个人信息路由最终使用 `/user/profile`，并保留模板 `/profile/index` 为隐藏页，是否符合后续信息架构？本次建议按此实现。
