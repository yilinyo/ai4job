import { Response, Request } from 'express'
import { hasApplicationsForResumeVersion } from './resume-applications'
import { incrementRepoVersionCount, decrementRepoVersionCount } from './resume-repos'

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

const versionList: ResumeVersion[] = [
  {
    id: '1',
    repoId: '1',
    version: 'v1.0.0',
    title: '初始版本',
    content: `# 个人简历

## 基本信息
- 姓名：张三
- 年龄：25岁
- 电话：13800138000
- 邮箱：admin@example.com

## 教育经历
2017-2021 某某大学 计算机科学与技术 本科

## 工作经历
### 前端工程师 | ABC科技公司 | 2021.07-至今
- 负责公司主要产品的前端开发
- 使用Vue.js/React开发企业级应用
- 优化页面性能，提升用户体验

## 技能清单
- 熟练掌握 HTML/CSS/JavaScript
- 熟悉 Vue.js、React 等前端框架
- 了解 Node.js 后端开发
- 熟悉 Git 版本控制`,
    remark: '第一版简历',
    isActive: false,
    createdAt: new Date('2024-01-15').toISOString(),
    updatedAt: new Date('2024-01-15').toISOString()
  },
  {
    id: '2',
    repoId: '1',
    version: 'v1.1.0',
    title: '增加项目经验',
    content: `# 个人简历

## 基本信息
- 姓名：张三
- 年龄：25岁
- 电话：13800138000
- 邮箱：admin@example.com

## 教育经历
2017-2021 某某大学 计算机科学与技术 本科

## 工作经历
### 前端工程师 | ABC科技公司 | 2021.07-至今
- 负责公司主要产品的前端开发
- 使用Vue.js/React开发企业级应用
- 优化页面性能，提升用户体验

## 项目经验
### 企业管理系统
- 技术栈：Vue3 + TypeScript + Element Plus
- 负责整体架构设计和核心模块开发
- 实现了权限管理、数据可视化等功能

### 电商平台移动端
- 技术栈：React + Redux + Ant Design Mobile
- 开发商品列表、购物车、订单管理等模块
- 优化首屏加载速度，提升30%性能

## 技能清单
- 熟练掌握 HTML/CSS/JavaScript
- 熟悉 Vue.js、React 等前端框架
- 了解 Node.js 后端开发
- 熟悉 Git 版本控制`,
    remark: '添加了详细的项目经验',
    isActive: false,
    createdAt: new Date('2024-02-01').toISOString(),
    updatedAt: new Date('2024-02-01').toISOString()
  },
  {
    id: '3',
    repoId: '1',
    version: 'v2.0.0',
    title: '大厂定制版',
    content: `# 张三 | 前端工程师

📧 admin@example.com | 📱 13800138000 | 💼 3年工作经验

## 💼 工作经历

### 高级前端工程师 | ABC科技公司 | 2021.07-至今
**核心职责：**
- 负责公司核心产品的前端架构设计与开发，服务用户数 100万+
- 主导前端工程化建设，搭建 CI/CD 流程，提升团队开发效率 40%
- 优化首屏加载速度，FCP 从 3.2s 降低到 1.1s
- 带领 5 人团队完成多个重点项目

## 🚀 项目经验

### 企业级 SaaS 管理平台
**技术栈：** Vue3 + TypeScript + Vite + Pinia + Element Plus
**项目描述：** 面向中大型企业的一站式管理平台，包含权限管理、数据分析、工作流等核心模块
**个人贡献：**
- 设计并实现微前端架构，支持子应用独立开发部署
- 封装了 20+ 业务组件，提升开发效率 50%
- 实现了动态权限路由系统，支持角色级权限控制
- 使用 Web Worker 优化大数据量渲染，处理 10万+ 条数据不卡顿

### 高性能电商移动端
**技术栈：** React18 + TypeScript + Redux Toolkit + TailwindCSS
**项目描述：** 面向 C 端用户的电商平台，日活 50万+
**个人贡献：**
- 实现虚拟滚动和图片懒加载，优化长列表性能
- 使用 SSR 技术提升 SEO 和首屏加载速度
- 接入第三方支付、物流系统，保证交易稳定性
- 完善单元测试覆盖率达 85%

## 🎓 教育背景
**某某大学** | 计算机科学与技术 | 本科 | 2017-2021

## 💡 技能清单
- **前端框架：** Vue2/3、React、Angular（熟练）
- **工程化：** Webpack、Vite、Rollup、pnpm（熟练）
- **编程语言：** TypeScript、JavaScript、Node.js（熟练）
- **UI 框架：** Element Plus、Ant Design、TailwindCSS
- **状态管理：** Vuex、Pinia、Redux、Zustand
- **测试：** Jest、Vitest、Cypress
- **其他：** Git、Docker、Linux、Nginx

## 🏆 个人优势
- 3年大厂前端开发经验，参与多个千万级用户产品
- 具备良好的工程化思维和架构设计能力
- 注重代码质量和性能优化
- 学习能力强，善于团队协作`,
    remark: '针对大厂优化的版本，突出量化指标和技术深度',
    isActive: true,
    createdAt: new Date('2024-03-10').toISOString(),
    updatedAt: new Date('2024-03-10').toISOString()
  },
  {
    id: '4',
    repoId: '2',
    version: 'v1.0.0',
    title: '全栈通用版',
    content: `# 个人简历 - 全栈工程师

## 基本信息
姓名：张三 | 年龄：25岁 | 电话：13800138000 | 邮箱：admin@example.com

## 技能栈
**前端：** Vue、React、TypeScript
**后端：** Node.js、Express、Koa、NestJS
**数据库：** MySQL、MongoDB、Redis
**运维：** Docker、Nginx、Linux

## 工作经历
### 全栈工程师 | ABC科技 | 2021.07-至今
独立负责项目的前后端开发和部署

## 项目经验
### CMS内容管理系统
全栈开发，包含用户管理、内容发布、数据统计等功能`,
    remark: '全栈版本',
    isActive: true,
    createdAt: new Date('2024-02-01').toISOString(),
    updatedAt: new Date('2024-02-01').toISOString()
  },
  {
    id: '5',
    repoId: '2',
    version: 'v1.1.0',
    title: '创业公司版',
    content: `# 张三 - 全栈工程师

💻 一人可独当一面，快速迭代产品

## 个人优势
- 3年全栈开发经验，能独立完成项目从 0 到 1
- 快速学习新技术，适应创业公司节奏
- 有创业公司工作经验，能承受高压

## 核心技能
前端、后端、数据库、运维全栈技能

## 项目案例
独立完成多个 MVP 产品的开发和上线`,
    remark: '突出独立开发能力和创业经验',
    isActive: false,
    createdAt: new Date('2024-02-20').toISOString(),
    updatedAt: new Date('2024-02-20').toISOString()
  }
]

export const getResumeVersions = (req: Request, res: Response) => {
  const { repoId, keyword, page = 1, pageSize = 10 } = req.query

  if (!repoId) {
    return res.status(400).json({
      code: 50001,
      message: 'repoId 不能为空'
    })
  }

  let filteredVersions = versionList.filter(v => v.repoId === repoId)

  // 关键词搜索
  if (keyword) {
    const searchTerm = (keyword as string).toLowerCase()
    filteredVersions = filteredVersions.filter(version =>
      version.version.toLowerCase().includes(searchTerm) ||
      version.title.toLowerCase().includes(searchTerm)
    )
  }

  // 按创建时间倒序
  filteredVersions.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

  // 分页
  const pageNum = Number(page)
  const size = Number(pageSize)
  const start = (pageNum - 1) * size
  const end = start + size
  const pagedVersions = filteredVersions.slice(start, end)

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      list: pagedVersions,
      total: filteredVersions.length,
      page: pageNum,
      pageSize: size
    }
  })
}

export const getResumeVersionById = (req: Request, res: Response) => {
  const { id } = req.params
  const version = versionList.find(v => v.id === id)

  if (!version) {
    return res.status(404).json({
      code: 50004,
      message: '版本不存在'
    })
  }

  return res.json({
    code: 20000,
    message: '获取成功',
    data: version
  })
}

export const createResumeVersion = (req: Request, res: Response) => {
  const { repoId, title, content, remark } = req.body

  // 生成版本号
  const repoVersions = versionList.filter(v => v.repoId === repoId)
  const versionNumber = `v1.${repoVersions.length}.0`

  const newVersion: ResumeVersion = {
    id: String(versionList.length + 1),
    repoId,
    version: versionNumber,
    title,
    content,
    remark,
    isActive: false,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }

  versionList.push(newVersion)
  incrementRepoVersionCount(repoId)

  return res.json({
    code: 20000,
    message: '创建成功',
    data: newVersion
  })
}

export const updateResumeVersion = (req: Request, res: Response) => {
  const { id } = req.params
  const version = versionList.find(v => v.id === id)

  if (!version) {
    return res.status(404).json({
      code: 50004,
      message: '版本不存在'
    })
  }

  Object.assign(version, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  return res.json({
    code: 20000,
    message: '更新成功',
    data: version
  })
}

export const deleteResumeVersion = (req: Request, res: Response) => {
  const { id } = req.params
  const index = versionList.findIndex(v => v.id === id)

  if (index === -1) {
    return res.status(404).json({
      code: 50004,
      message: '版本不存在'
    })
  }

  if (hasApplicationsForResumeVersion(id)) {
    return res.status(400).json({
      code: 50001,
      message: '该简历版本已关联投递记录，暂不能删除，请先处理关联投递记录'
    })
  }

  decrementRepoVersionCount(versionList[index].repoId)
  versionList.splice(index, 1)

  return res.json({
    code: 20000,
    message: '删除成功',
    data: null
  })
}

export const activateResumeVersion = (req: Request, res: Response) => {
  const { id } = req.params
  const version = versionList.find(v => v.id === id)

  if (!version) {
    return res.status(404).json({
      code: 50004,
      message: '版本不存在'
    })
  }

  // 取消同仓库下其他版本的激活状态
  versionList.forEach(v => {
    if (v.repoId === version.repoId) {
      v.isActive = false
    }
  })

  // 激活当前版本
  version.isActive = true
  version.updatedAt = new Date().toISOString()

  return res.json({
    code: 20000,
    message: '激活成功',
    data: null
  })
}
