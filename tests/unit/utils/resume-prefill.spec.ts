import { generateResumeDraftFromProfile } from '@/utils/resume-prefill'
import { UserProfile } from '@/models'

describe('Utils:resume-prefill', () => {
  const profile: UserProfile = {
    id: '1',
    username: 'admin',
    realName: '张三',
    age: 25,
    email: 'admin@example.com',
    phone: '13800138000',
    jobIntention: '前端工程师',
    school: '浙江大学',
    education: '本科',
    major: '计算机科学与技术',
    graduationDate: '2026-06',
    location: '杭州',
    personalAdvantage: '熟悉 Vue、TypeScript 与前端工程化。',
    createdAt: '2026-05-09T00:00:00.000Z',
    updatedAt: '2026-05-09T00:00:00.000Z'
  }

  it('generates editable resume draft from user profile', () => {
    const draft = generateResumeDraftFromProfile(profile)

    expect(draft).toContain('# 张三 | 前端工程师')
    expect(draft).toContain('## 基本信息')
    expect(draft).toContain('- 电话：13800138000')
    expect(draft).toContain('## 教育背景')
    expect(draft).toContain('- 学校：浙江大学')
    expect(draft).toContain('## 个人优势')
    expect(draft).toContain('熟悉 Vue、TypeScript 与前端工程化。')
  })

  it('skips empty fields without undefined text', () => {
    const draft = generateResumeDraftFromProfile({
      id: '2',
      username: 'user',
      createdAt: '2026-05-09T00:00:00.000Z',
      updatedAt: '2026-05-09T00:00:00.000Z'
    })

    expect(draft).toBe('# user')
    expect(draft).not.toContain('undefined')
  })
})
