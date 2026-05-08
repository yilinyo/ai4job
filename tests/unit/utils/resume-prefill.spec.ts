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
    educationExperiences: [
      {
        school: '浙江大学',
        education: '本科',
        major: '计算机科学与技术',
        admissionDate: '2022-09',
        graduationDate: '2026-06'
      },
      {
        school: '清华大学',
        education: '硕士',
        major: '软件工程',
        admissionDate: '2026-09',
        graduationDate: '2029-06'
      }
    ],
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
    expect(draft).toContain('- 本科：浙江大学：计算机科学与技术 / 2022-09 - 2026-06')
    expect(draft).toContain('- 硕士：清华大学：软件工程 / 2026-09 - 2029-06')
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
