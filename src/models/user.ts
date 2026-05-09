/**
 * 用户模型定义
 */

// 用户基本信息
export interface EducationExperience {
  school?: string // 学校(选填)
  education?: string // 学历(选填)
  major?: string // 专业(选填)
  admissionDate?: string // 入学时间(选填)
  graduationDate?: string // 毕业时间(选填)
}

export interface UserProfile {
  id: string
  username: string // 用户名
  realName?: string // 真实姓名(选填)
  age?: number // 年龄(选填)
  email?: string // 邮箱(选填)
  phone?: string // 电话(选填)
  jobIntention?: string // 求职意向(选填)
  avatar?: string // 头像URL
  educationExperiences?: EducationExperience[] // 教育经历(选填)
  location?: string // 所在地(选填)
  personalAdvantage?: string // 个人优势(选填)
  createdAt: string // 创建时间
  updatedAt: string // 更新时间
}

// 用户注册请求
export interface RegisterRequest {
  username: string
  password: string
  email: string
  emailCode: string
}

export interface SendEmailCodeRequest {
  email: string
}

export interface SendEmailCodeResponse {
  email: string
  mockCode: string
  expiresIn: number
}

// 用户登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 用户登录响应
export interface LoginResponse {
  token: string
  userInfo: UserProfile
}

// 更新用户信息请求
export interface UpdateProfileRequest {
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
