import faker from 'faker'
import { Response, Request } from 'express'

interface EducationExperience {
  school?: string
  education?: string
  major?: string
  admissionDate?: string
  graduationDate?: string
}

// 简历管理系统的用户数据
interface UserProfile {
  id: string
  username: string
  password: string
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

interface EmailCodeRecord {
  email: string
  code: string
  expiresAt: number
}

const MOCK_EMAIL_CODE = '123456'
const EMAIL_CODE_EXPIRES_IN = 5 * 60
const emailCodeMap: Record<string, EmailCodeRecord> = {}

const userList: UserProfile[] = [
  {
    id: '1',
    username: 'admin',
    password: '111111',
    realName: '张三',
    age: 25,
    email: 'admin@example.com',
    phone: '13800138000',
    jobIntention: '前端工程师',
    educationExperiences: [
      {
        school: '浙江大学',
        education: '本科',
        major: '计算机科学与技术',
        admissionDate: '2022-09',
        graduationDate: '2026-06'
      }
    ],
    location: '杭州',
    personalAdvantage: '熟悉 Vue、TypeScript 与前端工程化，有多个后台系统项目经验。',
    avatar: 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: '2',
    username: 'user',
    password: '111111',
    realName: '李四',
    age: 28,
    email: 'user@example.com',
    phone: '13900139000',
    jobIntention: '全栈工程师',
    educationExperiences: [
      {
        school: '上海交通大学',
        education: '硕士',
        major: '软件工程',
        admissionDate: '2021-09',
        graduationDate: '2024-06'
      },
      {
        school: '南京大学',
        education: '本科',
        major: '计算机科学与技术',
        admissionDate: '2017-09',
        graduationDate: '2021-06'
      }
    ],
    location: '上海',
    personalAdvantage: '具备前后端全栈开发经验，熟悉 Node.js、Vue 和接口设计。',
    avatar: 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
]

export const sendEmailCode = (req: Request, res: Response) => {
  const { email } = req.body

  if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return res.status(400).json({
      code: 50006,
      message: '邮箱格式不正确'
    })
  }

  emailCodeMap[email] = {
    email,
    code: MOCK_EMAIL_CODE,
    expiresAt: Date.now() + EMAIL_CODE_EXPIRES_IN * 1000
  }

  return res.json({
    code: 20000,
    message: '验证码已发送',
    data: {
      email,
      mockCode: MOCK_EMAIL_CODE,
      expiresIn: EMAIL_CODE_EXPIRES_IN
    }
  })
}

export const register = (req: Request, res: Response) => {
  const { username, email, emailCode } = req.body

  if (userList.find(u => u.username === username)) {
    return res.status(400).json({
      code: 50001,
      message: '用户名已存在'
    })
  }

  if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return res.status(400).json({
      code: 50006,
      message: '邮箱格式不正确'
    })
  }

  if (userList.find(u => u.email === email)) {
    return res.status(400).json({
      code: 50007,
      message: '邮箱已注册'
    })
  }

  const emailCodeRecord = emailCodeMap[email]
  if (!emailCodeRecord) {
    return res.status(400).json({
      code: 50009,
      message: '请先获取邮箱验证码'
    })
  }

  if (emailCodeRecord.expiresAt < Date.now()) {
    delete emailCodeMap[email]
    return res.status(400).json({
      code: 50010,
      message: '验证码已过期'
    })
  }

  if (emailCodeRecord.code !== emailCode) {
    return res.status(400).json({
      code: 50011,
      message: '验证码错误'
    })
  }

  const newUser: UserProfile = {
    id: String(userList.length + 1),
    username,
    password: req.body.password,
    email,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
  userList.push(newUser)
  delete emailCodeMap[email]

  return res.json({
    code: 20000,
    message: '注册成功',
    data: {
      id: newUser.id,
      username: newUser.username,
      email: newUser.email,
      createdAt: newUser.createdAt,
      updatedAt: newUser.updatedAt
    }
  })
}

export const login = (req: Request, res: Response) => {
  const { username, password } = req.body
  const user = userList.find(u => u.username === username)

  if (!user) {
    return res.status(400).json({
      code: 50004,
      message: '用户不存在'
    })
  }

  if (user.password !== password) {
    return res.status(400).json({
      code: 50005,
      message: '密码错误'
    })
  }

  const { password: pwd, ...userInfo } = user

  return res.json({
    code: 20000,
    message: '登录成功',
    data: {
      token: username + '-token',
      userInfo
    }
  })
}

export const logout = (req: Request, res: Response) => {
  return res.json({
    code: 20000,
    message: '登出成功'
  })
}

export const getUserProfile = (req: Request, res: Response) => {
  const token = req.header('X-Access-Token')
  const username = token?.replace('-token', '')
  const user = userList.find(u => u.username === username)

  if (!user) {
    return res.status(401).json({
      code: 50008,
      message: '未授权'
    })
  }

  const { password, ...userProfile } = user

  return res.json({
    code: 20000,
    message: '获取成功',
    data: {
      ...userProfile,
      roles: ['admin']
    }
  })
}

export const updateUserProfile = (req: Request, res: Response) => {
  const token = req.header('X-Access-Token')
  const username = token?.replace('-token', '')
  const user = userList.find(u => u.username === username)

  if (!user) {
    return res.status(401).json({
      code: 50008,
      message: '未授权'
    })
  }

  // 更新用户信息
  Object.assign(user, {
    ...req.body,
    updatedAt: new Date().toISOString()
  })

  const { password, ...userProfile } = user

  return res.json({
    code: 20000,
    message: '更新成功',
    data: userProfile
  })
}
