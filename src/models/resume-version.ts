/**
 * 简历版本模型定义
 */

// 简历版本
export interface ResumeVersion {
  id: string
  repoId: string // 所属仓库ID
  version: string // 版本号(如 v1.0.0)
  title: string // 版本标题
  content: string // 简历内容(支持富文本/Markdown)
  remark?: string // 版本备注
  isActive: boolean // 是否为当前激活版本
  createdAt: string // 创建时间
  updatedAt: string // 更新时间
}

// 创建简历版本请求
export interface CreateResumeVersionRequest {
  repoId: string
  title: string
  content: string
  remark?: string
}

// 更新简历版本请求
export interface UpdateResumeVersionRequest {
  title?: string
  content?: string
  remark?: string
  isActive?: boolean
}

// 简历版本列表查询参数
export interface ResumeVersionQueryParams {
  repoId: string
  keyword?: string // 搜索关键词(版本号或标题)
  page?: number
  pageSize?: number
}

// 简历版本列表响应
export interface ResumeVersionListResponse {
  list: ResumeVersion[]
  total: number
  page: number
  pageSize: number
}

// 版本对比数据结构(预留)
export interface VersionComparison {
  oldVersion: ResumeVersion
  newVersion: ResumeVersion
  diff: string // 差异内容
}
