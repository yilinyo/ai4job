# 简历管理系统 (Resume Management System)

> 基于 Vue + TypeScript + Element UI 的简历版本管理系统
>
> 目录结构已调整：前端工程位于仓库根目录的 `frontend/` 下。本文为历史前端 README，命令请在 `frontend/` 目录执行；仓库总览请查看根目录 `README.md`。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue](https://img.shields.io/badge/vue-2.6.12-brightgreen.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/typescript-4.2.3-blue.svg)](https://www.typescriptlang.org/)
[![Element UI](https://img.shields.io/badge/element--ui-2.15.1-blue.svg)](https://element.eleme.io/)

## 📖 项目介绍

简历管理系统是一个帮助求职者高效管理简历版本的工具,支持:

- 📂 **多仓库管理** - 为不同岗位/公司创建独立的简历仓库
- 📝 **版本控制** - 像Git一样管理简历的不同版本
- 🎯 **分类管理** - 区分校招和社招简历
- 💾 **导出功能** - 一键导出简历文件
- 👤 **个人信息** - 集中管理个人求职信息

## ✨ 核心功能

### 1. 简历仓库管理
- 创建多个简历仓库,针对不同求职方向
- 支持校招/社招分类
- 搜索、筛选、分页功能
- 查看每个仓库的版本数量

### 2. 简历版本管理
- 为每个仓库创建多个简历版本
- 版本自动编号(v1.0.0, v1.1.0...)
- 激活/切换当前使用的版本
- 版本备注功能,记录修改说明
- 简历内容支持Markdown格式

### 3. 个人信息管理
- 填写基本信息(姓名、年龄、邮箱、电话)
- 设置求职意向
- 信息可随时更新

### 4. 简历查看与导出
- 预览简历内容
- 导出为文本文件
- 未来可扩展为PDF导出

## 🚀 快速开始

### 环境要求
- Node.js >= 14.x
- npm >= 6.x

### 安装依赖
```bash
npm install --legacy-peer-deps
```

### 启动开发服务器
```bash
npm run serve
```

这会同时启动:
- 前端开发服务器: http://localhost:8080
- Mock API服务器: http://localhost:9528

### 登录系统
```
用户名: admin
密码: 111111
```

> 📘 详细使用说明请查看 [QUICK_START.md](./QUICK_START.md)

## 📁 项目结构

```
src/
├── models/              # 数据模型定义
│   ├── user.ts
│   ├── resume-repo.ts
│   └── resume-version.ts
├── api/                 # API接口层
│   ├── users.ts
│   ├── resume-repos.ts
│   └── resume-versions.ts
├── store/modules/       # Vuex状态管理
│   ├── user.ts
│   ├── resume-repo.ts
│   └── resume-version.ts
├── views/               # 页面组件
│   ├── user/           # 用户相关页面
│   ├── resume-repo/    # 简历仓库页面
│   └── resume-version/ # 简历版本页面
├── components/          # 通用组件
│   └── ResumeRepoCard/ # 仓库卡片组件
└── router/modules/      # 路由配置
    └── resume.ts

mock/                    # Mock数据服务
├── users.ts            # 用户Mock
├── resume-repos.ts     # 仓库Mock
├── resume-versions.ts  # 版本Mock
└── mock-server.ts      # Mock服务器
```

## 🛠️ 技术栈

### 前端核心
- **Vue 2.6** - 渐进式JavaScript框架
- **TypeScript 4.2** - JavaScript超集,提供类型安全
- **Element UI 2.15** - 优秀的Vue UI组件库
- **Vue Router 3.5** - 官方路由管理器
- **Vuex 3.6** - 状态管理模式
- **Axios** - HTTP客户端

### 工程化
- **Vue CLI 4.5** - Vue项目脚手架
- **ESLint** - 代码质量检查
- **Sass** - CSS预处理器
- **TypeScript** - 类型检查

### Mock服务
- **Express** - Node.js Web框架
- **CORS** - 跨域资源共享
- **Morgan** - HTTP请求日志

## 📊 数据模型

### 用户模型
```typescript
interface UserProfile {
  id: string
  username: string
  realName?: string
  age?: number
  email?: string
  phone?: string
  jobIntention?: string
  avatar?: string
  createdAt: string
  updatedAt: string
}
```

### 简历仓库模型
```typescript
interface ResumeRepo {
  id: string
  userId: string
  name: string
  jobType: 'campus' | 'social'  // 校招/社招
  targetPosition: string
  description?: string
  versionCount: number
  createdAt: string
  updatedAt: string
}
```

### 简历版本模型
```typescript
interface ResumeVersion {
  id: string
  repoId: string
  version: string
  title: string
  content: string
  remark?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}
```

## 🔌 API接口

### 用户相关
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `POST /api/user/logout` - 用户登出
- `GET /api/user/profile` - 获取个人信息
- `PUT /api/user/profile` - 更新个人信息

### 简历仓库
- `GET /api/resume-repos` - 获取仓库列表
- `GET /api/resume-repos/:id` - 获取仓库详情
- `POST /api/resume-repos` - 创建仓库
- `PUT /api/resume-repos/:id` - 更新仓库
- `DELETE /api/resume-repos/:id` - 删除仓库

### 简历版本
- `GET /api/resume-versions` - 获取版本列表
- `GET /api/resume-versions/:id` - 获取版本详情
- `POST /api/resume-versions` - 创建版本
- `PUT /api/resume-versions/:id` - 更新版本
- `DELETE /api/resume-versions/:id` - 删除版本
- `POST /api/resume-versions/:id/activate` - 激活版本

> 📘 完整API文档请查看 [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)

## 🎯 测试数据

系统预置了完整的测试数据:

- **2个测试账号**: admin/111111 和 user/111111
- **3个简历仓库**: 涵盖前端、全栈、后端等岗位
- **5个简历版本**: 包含不同场景的简历示例

## 📝 开发命令

```bash
# 安装依赖
npm install --legacy-peer-deps

# 启动开发服务器(前端+Mock)
npm run serve

# 代码检查
npm run lint

# 构建生产版本
npm run build:prod

# 构建预发布版本
npm run build:stage
```

## 🔄 接入真实后端

当后端API开发完成后,只需两步即可切换:

1. **修改API基础地址**
   ```typescript
   // src/utils/request.ts
   const service = axios.create({
     baseURL: 'https://your-api-domain.com',
     timeout: 5000
   })
   ```

2. **停止Mock服务,启动前端**
   ```bash
   npx vue-cli-service serve
   ```

前端代码无需任何修改! 🎉

## 🌟 功能亮点

1. **完整的TypeScript类型系统** - 所有接口和模型都有完整类型定义
2. **模块化架构** - 按功能模块清晰分包,易于维护
3. **Mock与真实API无缝切换** - 开发阶段使用Mock,生产环境切换真实API
4. **响应式设计** - 界面美观,支持不同屏幕尺寸
5. **完善的错误处理** - 友好的错误提示和异常处理

## 🚧 未来规划

- [ ] AI辅助简历优化功能
- [ ] 简历版本对比功能
- [ ] 简历导出为PDF
- [ ] 岗位投递记录管理
- [ ] 面试记录管理
- [ ] 简历模板系统
- [ ] 简历数据统计分析

## 📄 文档

- [快速开始指南](./QUICK_START.md) - 详细的启动和使用说明
- [项目进度总结](./PROJECT_SUMMARY.md) - 完整的功能清单和API文档

## 🤝 贡献

欢迎提交Issue和Pull Request!

## 📜 开源协议

[MIT](LICENSE)

---

## 🙏 致谢

本项目基于 [vue-typescript-admin-template](https://github.com/Armour/vue-typescript-admin-template) 模板开发。

---

⭐ 如果这个项目对你有帮助,请给一个Star支持! ⭐
