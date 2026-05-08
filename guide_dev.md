# AI-For 简历管理系统 - 开发指南

> 基于 **Vue 2 + TypeScript + Element UI + Vuex (vuex-module-decorators)** 的后台管理系统，fork 自 `vue-typescript-admin-template`，扩展了简历管理业务功能。

---

## 一、项目工程结构

### 根目录

```
ai-for/
├── mock/                  # Mock 服务器（Express 独立进程，端口 9528）
├── src/                   # 前端源码
├── public/                # 静态资源（index.html 等）
├── tests/                 # 单元测试
├── .env.development       # 开发环境变量 → VUE_APP_BASE_API=/dev-api
├── .env.production        # 生产环境变量
├── .env.staging           # 预发布环境变量
├── vue.config.js          # Webpack / DevServer 配置（含 proxy 代理到 mock）
├── tsconfig.json          # TypeScript 配置
├── babel.config.js        # Babel 配置
├── jest.config.js         # 测试配置
└── package.json           # 依赖和脚本
```

### src/ 核心目录

```
src/
├── main.ts                # 入口：初始化 Vue、Element UI、Router、Store、i18n
├── App.vue                # 根组件
├── permission.ts          # ⭐ 路由守卫：鉴权 + 动态路由注入（核心文件）
├── settings.ts            # 全局配置（标题、侧边栏固定等）
├── shims.d.ts             # TypeScript 声明（.vue 文件模块声明）
│
├── models/                # 📦 数据模型（TypeScript 接口定义）
│   ├── index.ts           #   统一导出 + 通用类型（ApiResponse、PaginationParams）
│   ├── user.ts            #   UserProfile、LoginRequest、LoginResponse 等
│   ├── resume-repo.ts     #   ResumeRepo、JobType 枚举、CRUD 请求/响应类型
│   └── resume-version.ts  #   ResumeVersion、CRUD 请求/响应类型、VersionComparison
│
├── api/                   # 📡 API 请求层（每个业务模块一个文件）
│   ├── users.ts           #   login / logout / getUserProfile / updateUserProfile
│   ├── resume-repos.ts    #   getResumeRepos / getResumeRepoById / create / update / delete
│   ├── resume-versions.ts #   getResumeVersions / create / update / delete / activate
│   ├── articles.ts        #   (原模板) 文章相关
│   ├── roles.ts           #   (原模板) 角色相关
│   └── transactions.ts    #   (原模板) 交易相关
│
├── store/                 # 🗄️ Vuex 状态管理（vuex-module-decorators 动态模块）
│   ├── index.ts           #   创建空 store，各模块通过装饰器动态注册
│   └── modules/
│       ├── user.ts        #   ⭐ 用户模块：token / roles / userProfile / Login / Logout
│       ├── permission.ts  #   ⭐ 权限模块：根据 roles 过滤 asyncRoutes 生成动态路由
│       ├── resume-repo.ts #   简历仓库模块：列表 / 详情 / CRUD actions
│       ├── resume-version.ts # 简历版本模块：列表 / 详情 / CRUD / 激活
│       ├── app.ts         #   应用状态：侧边栏开关 / 语言 / 尺寸
│       ├── settings.ts    #   设置状态：主题 / 固定 header / 显示标签栏
│       ├── tags-view.ts   #   标签页导航状态
│       └── error-log.ts   #   错误日志状态
│
├── router/                # 🔀 路由配置
│   ├── index.ts           #   constantRoutes（公开路由）+ asyncRoutes（权限路由）
│   └── modules/
│       ├── resume.ts      #   ⭐ 简历管理路由（仓库列表/详情、版本创建/编辑/查看）
│       ├── components.ts  #   (原模板) 组件示例
│       ├── charts.ts      #   (原模板) 图表
│       ├── table.ts       #   (原模板) 表格
│       └── nested.ts      #   (原模板) 嵌套路由
│
├── views/                 # 📄 页面组件（与 router 一一对应）
│   ├── login/             #   登录页（index.vue）
│   ├── user/              #   profile.vue（个人信息）、register.vue（注册）
│   ├── resume-repo/       #   list.vue（仓库列表）、detail.vue（仓库详情）
│   ├── resume-version/    #   editor.vue（版本编辑器）、view.vue（版本查看）
│   ├── dashboard/         #   仪表盘首页（admin / editor 两套）
│   ├── permission/        #   权限管理示例
│   └── ...                #   原模板其他示例页面
│
├── components/            # 🧩 公共/可复用组件
│   ├── ResumeRepoCard/    #   简历仓库卡片组件
│   ├── VersionEditor/     #   版本编辑器组件
│   ├── UserInfoForm/      #   用户信息表单组件
│   ├── Pagination/        #   分页组件
│   ├── Breadcrumb/        #   面包屑
│   ├── Hamburger/         #   侧边栏切换按钮
│   ├── HeaderSearch/      #   顶部搜索
│   ├── Charts/            #   图表组件（ECharts）
│   └── ...                #   更多通用组件
│
├── layout/                # 🖼️ 页面布局框架
│   │                      #   侧边栏(Sidebar) + 顶栏(Navbar) + 内容区(AppMain)
│   └── components/
│       ├── Sidebar/       #   侧边栏（根据权限路由动态生成菜单）
│       ├── Navbar/        #   顶部导航栏（头像、退出、面包屑）
│       └── AppMain/       #   主内容区（router-view）
│
├── utils/                 # 🔧 工具函数
│   ├── request.ts         #   ⭐ Axios 实例 + 请求/响应拦截器（最核心的文件之一）
│   ├── cookies.ts         #   Token 的 cookie 读写（getToken / setToken / removeToken）
│   ├── validate.ts        #   表单验证函数
│   ├── permission.ts      #   权限检查工具函数
│   └── ...                #   其他工具
│
├── styles/                # 🎨 全局样式（SCSS 变量、mixin、reset）
├── lang/                  # 🌐 国际化 i18n（中文/英文/日文/韩文）
├── icons/                 # SVG 图标
├── directives/            # 自定义指令（如 v-permission）
└── filters/               # 全局过滤器
```

### mock/ 目录

```
mock/
├── mock-server.ts         # Express 服务器入口（端口 9528），注册所有 API 路由
├── users.ts               # 用户 mock：login / register / getUserProfile / updateProfile
├── resume-repos.ts        # 简历仓库 mock：CRUD
├── resume-versions.ts     # 简历版本 mock：CRUD + 激活
├── articles.ts            # (原模板) 文章 mock
├── api.ts                 # 统一导出（给 swagger-routes-express 用）
├── security.ts            # Token 校验中间件
├── swagger.yml            # 原模板 Swagger API 定义
└── role/                  # 角色相关 mock 数据
```

---

## 二、核心机制详解

### 2.1 请求数据流

一次 API 请求的完整链路：

```
页面组件 (views/*.vue)
  │  调用 Store Action
  ▼
Store Module (store/modules/*.ts)
  │  调用 API 函数
  ▼
API 层 (api/*.ts)
  │  调用 request(config)
  ▼
request.ts 请求拦截器
  │  自动附加 X-Access-Token header
  ▼
vue.config.js devServer.proxy
  │  /dev-api/api/xxx → http://localhost:9528/api/xxx
  ▼
mock-server.ts (Express)
  │  处理请求，返回 { code: 20000, data: {...} }
  ▼
request.ts 响应拦截器
  │  code === 20000 → 返回 response.data（即 { code, data }）
  │  code !== 20000 → 弹出错误提示，reject
  ▼
Store Action 拿到数据
  │  const { data } = await apiCall()
  │  data 就是业务数据（不需要再 .data）
  │  调用 Mutation 更新 state
  ▼
页面组件通过 getter 响应式更新
```

### 2.2 登录鉴权流程

```
用户提交登录表单
  │
  ▼
UserModule.Login()
  │  POST /api/user/login → 获取 token
  │  setToken(token) → 写入 cookie
  │  SET_TOKEN / SET_USER_PROFILE → 更新 store
  ▼
router.push('/') → 触发路由守卫
  │
  ▼
permission.ts beforeEach 守卫
  │  有 token？
  │  ├── 是 → roles 为空？
  │  │       ├── 是 → GetUserInfo() 获取用户信息和角色
  │  │       │       → PermissionModule.GenerateRoutes(roles) 生成动态路由
  │  │       │       → router.addRoute() 注入路由
  │  │       │       → next({ ...to, replace: true })
  │  │       └── 否 → next() 直接放行
  │  └── 否 → 在白名单？
  │          ├── 是 → next()
  │          └── 否 → next('/login?redirect=xxx')
```

### 2.3 Vuex Store 模式

本项目使用 `vuex-module-decorators`，通过装饰器定义 store 模块：

```ts
// 类装饰器注册模块（dynamic: true 表示动态注册）
@Module({ dynamic: true, store, name: 'moduleName' })
class MyStore extends VuexModule {
  // state —— 直接定义为 class 属性
  public list: Item[] = []

  // getter —— 使用 get 访问器
  get itemCount() { return this.list.length }

  // mutation —— @Mutation 装饰器
  @Mutation
  private SET_LIST(list: Item[]) { this.list = list }

  // action —— @Action 装饰器（可以是 async）
  @Action
  public async FetchList() {
    const { data } = await getList()
    this.SET_LIST(data.list)  // action 中直接调用 mutation
  }
}

// 导出模块实例，页面中直接使用
export const MyModule = getModule(MyStore)

// 页面使用：
// MyModule.list        → 读 state
// MyModule.itemCount   → 读 getter
// MyModule.FetchList() → 调 action
```

### 2.4 路由与权限

- **constantRoutes**（`router/index.ts`）：所有用户可访问，如 `/login`、`/dashboard`、`/404`
- **asyncRoutes**（`router/index.ts`）：按角色过滤后动态注入，如简历管理、权限管理等
- 路由 `meta` 常用字段：
  - `title`: 显示在侧边栏和面包屑
  - `icon`: 侧边栏图标
  - `hidden`: 是否在侧边栏隐藏
  - `roles`: 允许访问的角色列表
  - `noCache`: 是否不缓存

---

## 三、快速上手：新增一个业务模块

以新增「公司管理」模块为例，完整流程 6 步：

### Step 1：定义数据模型 → `src/models/company.ts`

```ts
export interface Company {
  id: string
  name: string
  industry: string
  address?: string
  createdAt: string
  updatedAt: string
}

export interface CreateCompanyRequest {
  name: string
  industry: string
  address?: string
}

export interface UpdateCompanyRequest {
  name?: string
  industry?: string
  address?: string
}

export interface CompanyQueryParams {
  keyword?: string
  page?: number
  pageSize?: number
}
```

在 `src/models/index.ts` 添加：

```ts
export * from './company'
```

### Step 2：编写 API → `src/api/company.ts`

```ts
import request from '@/utils/request'
import { CreateCompanyRequest, UpdateCompanyRequest, CompanyQueryParams } from '@/models'

export const getCompanies = (params: CompanyQueryParams) =>
  request({ url: '/api/companies', method: 'get', params })

export const getCompanyById = (id: string) =>
  request({ url: `/api/companies/${id}`, method: 'get' })

export const createCompany = (data: CreateCompanyRequest) =>
  request({ url: '/api/companies', method: 'post', data })

export const updateCompany = (id: string, data: UpdateCompanyRequest) =>
  request({ url: `/api/companies/${id}`, method: 'put', data })

export const deleteCompany = (id: string) =>
  request({ url: `/api/companies/${id}`, method: 'delete' })
```

> **注意**：`request(...)` 不加泛型参数（如 `request<T>(...)`），AxiosInstance 直接调用不支持。

### Step 3：创建 Store → `src/store/modules/company.ts`

```ts
import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import { getCompanies, getCompanyById, createCompany, updateCompany, deleteCompany } from '@/api/company'
import { Company, CompanyQueryParams, CreateCompanyRequest, UpdateCompanyRequest } from '@/models'
import store from '@/store'

export interface ICompanyState {
  list: Company[]
  current: Company | null
  total: number
  loading: boolean
}

@Module({ dynamic: true, store, name: 'company' })
class CompanyStore extends VuexModule implements ICompanyState {
  // ⚠️ 类名必须和 model 接口 Company 不同！加 Store 后缀
  public list: Company[] = []
  public current: Company | null = null
  public total = 0
  public loading = false

  @Mutation
  private SET_LIST(list: Company[]) { this.list = list }

  @Mutation
  private SET_CURRENT(item: Company | null) { this.current = item }

  @Mutation
  private SET_TOTAL(total: number) { this.total = total }

  @Mutation
  private SET_LOADING(loading: boolean) { this.loading = loading }

  @Action
  public async GetCompanies(params: CompanyQueryParams) {
    this.SET_LOADING(true)
    try {
      const { data } = await getCompanies(params)
      if (data) {
        this.SET_LIST(data.list)     // ✅ 直接 data.list
        this.SET_TOTAL(data.total)   // ❌ 不是 data.data.list
      }
    } finally {
      this.SET_LOADING(false)
    }
  }

  @Action
  public async GetCompanyById(id: string) {
    const { data } = await getCompanyById(id)
    if (data) {
      this.SET_CURRENT(data)
    }
  }
}

export const CompanyModule = getModule(CompanyStore)
```

### Step 4：编写 Mock → `mock/company.ts`

```ts
import { Request, Response } from 'express'

const companyList = [
  { id: '1', name: '示例科技', industry: 'IT', address: '北京', createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
  { id: '2', name: '示例金融', industry: '金融', address: '上海', createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() }
]

export const getCompanies = (req: Request, res: Response) => {
  const { keyword, page = 1, pageSize = 10 } = req.query
  let filtered = companyList
  if (keyword) {
    filtered = companyList.filter(c => c.name.includes(keyword as string))
  }
  res.json({
    code: 20000,
    data: {
      list: filtered.slice((+page - 1) * +pageSize, +page * +pageSize),
      total: filtered.length
    }
  })
}

export const getCompanyById = (req: Request, res: Response) => {
  const company = companyList.find(c => c.id === req.params.id)
  res.json({ code: 20000, data: company || null })
}
```

在 `mock/mock-server.ts` 注册路由：

```ts
import * as company from './company'

app.get('/api/companies', company.getCompanies)
app.get('/api/companies/:id', company.getCompanyById)
```

### Step 5：添加路由 → `src/router/modules/company.ts`

```ts
import { RouteConfig } from 'vue-router'
import Layout from '@/layout/index.vue'

const companyRoutes: RouteConfig = {
  path: '/company',
  component: Layout,
  redirect: '/company/list',
  name: 'Company',
  meta: {
    title: '公司管理',
    icon: 'tree'
  },
  children: [
    {
      path: 'list',
      component: () => import('@/views/company/list.vue'),
      name: 'CompanyList',
      meta: { title: '公司列表' }
    },
    {
      path: ':id',
      component: () => import('@/views/company/detail.vue'),
      name: 'CompanyDetail',
      meta: { title: '公司详情', hidden: true }  // hidden 放在 meta 里！
    }
  ]
}

export default companyRoutes
```

在 `src/router/index.ts` 的 `asyncRoutes` 中引入：

```ts
import companyRoutes from './modules/company'

export const asyncRoutes: RouteConfig[] = [
  companyRoutes,
  // ...其他路由
]
```

### Step 6：编写页面 → `src/views/company/list.vue`

```vue
<template>
  <div class="company-list-container">
    <el-button type="primary" @click="handleCreate">新增公司</el-button>
    <el-table :data="list" v-loading="loading" style="width: 100%; margin-top: 20px">
      <el-table-column prop="name" label="公司名称" />
      <el-table-column prop="industry" label="行业" />
      <el-table-column prop="address" label="地址" />
    </el-table>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { CompanyModule } from '@/store/modules/company'

@Component({ name: 'CompanyList' })
export default class extends Vue {
  get list() { return CompanyModule.list }
  get loading() { return CompanyModule.loading }

  created() {
    CompanyModule.GetCompanies({ page: 1, pageSize: 10 })
  }

  private handleCreate() {
    this.$router.push('/company/create')
  }
}
</script>
```

---

## 四、开发约定与注意事项

### 命名规范

| 位置 | 规范 | 示例 |
|------|------|------|
| Model 接口 | PascalCase | `Company`、`ResumeRepo` |
| Store 类名 | PascalCase + **Store 后缀** | `CompanyStore`（避免与 model 同名） |
| Store 导出 | PascalCase + **Module 后缀** | `CompanyModule` |
| API 函数 | camelCase | `getCompanies`、`createCompany` |
| 路由文件 | kebab-case | `company.ts`、`resume-repo.ts` |
| 页面文件 | kebab-case | `list.vue`、`detail.vue` |
| Mock 文件 | kebab-case（与 API 对应） | `company.ts` |

### 已知的坑（避免踩雷）

| 坑 | 原因 | 正确做法 |
|----|------|----------|
| Store 类名和 Model 接口同名 | TS 报 `Import declaration conflicts with local declaration` | **Store 类名加 `Store` 后缀**：`CompanyStore` |
| `data.data.xxx` 取不到值 | `request.ts` 拦截器已返回 `response.data` | `const { data } = await api()` 后直接用 `data.xxx` |
| `request<T>(...)` 编译报错 | `axios.create()` 返回的实例直接调用不支持泛型 | 写 `request({...})` 不加泛型 |
| Vue 2 模板中用 `?.` 可选链 | Vue 2 模板编译器不支持 ES2020 语法 | 用 `v-if` 判断或三元表达式 |
| 路由上写 `hidden: true` | `RouteConfig` 类型没有 `hidden` 属性 | 放到 `meta: { hidden: true }` |
| `UserModule.roles` undefined | 原模板的 User store 被重写后缺少 roles | 必须在 Store 中定义 `roles` 字段，`GetUserInfo` 中赋值 |
| `npm install` 失败 | npm 7+ 严格检查 peer deps | 使用 `npm install --legacy-peer-deps` |
| webpack 5 与 vue-cli 4 冲突 | vue-cli-service 4.x 只兼容 webpack 4 | `package.json` 中 webpack 版本应为 `^4.x` |
| 修改文件后编译不生效 | babel-loader 缓存 | 删除 `node_modules/.cache` 后重启 |

### Mock 数据响应格式约定

所有 Mock 接口必须返回统一格式：

```json
{
  "code": 20000,
  "message": "操作成功",
  "data": { ... }
}
```

- `code: 20000` — 成功
- `code: 50001` — 业务错误
- `code: 50008` — 未授权（触发自动登出）
- HTTP 4xx — 请求错误（被 axios error 拦截器捕获）

---

## 五、常用命令

```bash
# 开发（同时启动前端 + mock 服务）
npm run serve

# 仅启动 mock 服务
npm run mock

# 生产构建
npm run build:prod

# ESLint 检查
npm run lint

# 单元测试
npm run test:unit

# 清除缓存（编译不生效时使用）
rm -rf node_modules/.cache
```

---

## 六、文件依赖关系图

```
main.ts
  ├── App.vue
  ├── router/index.ts ─────────── constantRoutes + asyncRoutes
  │     └── modules/resume.ts     （简历管理路由）
  ├── store/index.ts ──────────── 空 store（模块动态注册）
  │     └── modules/
  │           ├── user.ts          Login → api/users.ts → mock/users.ts
  │           ├── permission.ts    过滤 asyncRoutes
  │           ├── resume-repo.ts   CRUD → api/resume-repos.ts → mock/resume-repos.ts
  │           └── resume-version.ts CRUD → api/resume-versions.ts → mock/resume-versions.ts
  ├── permission.ts ───────────── beforeEach 守卫（调用 user + permission store）
  └── utils/request.ts ───────── Axios 实例（拦截器处理 token + code 校验）
        └── vue.config.js proxy    /dev-api → localhost:9528
```
