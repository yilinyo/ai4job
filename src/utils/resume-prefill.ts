import { UserProfile } from '@/models'

const appendLine = (lines: string[], label: string, value?: string | number) => {
  if (value === undefined || value === null || value === '') return
  lines.push(`- ${label}：${value}`)
}

export function generateResumeDraftFromProfile(profile?: UserProfile | null) {
  if (!profile) return ''

  const name = profile.realName || profile.username || '个人简历'
  const title = profile.jobIntention ? `${name} | ${profile.jobIntention}` : name
  const lines: string[] = [`# ${title}`, '']

  const basicInfo: string[] = []
  appendLine(basicInfo, '年龄', profile.age)
  appendLine(basicInfo, '电话', profile.phone)
  appendLine(basicInfo, '邮箱', profile.email)
  appendLine(basicInfo, '所在地', profile.location)
  if (basicInfo.length) {
    lines.push('## 基本信息', ...basicInfo, '')
  }

  const educationInfo: string[] = []
  appendLine(educationInfo, '学校', profile.school)
  appendLine(educationInfo, '学历', profile.education)
  appendLine(educationInfo, '专业', profile.major)
  appendLine(educationInfo, '毕业时间', profile.graduationDate)
  if (educationInfo.length) {
    lines.push('## 教育背景', ...educationInfo, '')
  }

  if (profile.personalAdvantage) {
    lines.push('## 个人优势', profile.personalAdvantage, '')
  }

  if (profile.jobIntention) {
    lines.push('## 求职意向', profile.jobIntention, '')
  }

  return lines.join('\n').trim()
}
