<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <div slot="header">
        <span class="card-title">个人信息管理</span>
      </div>

      <el-form
        ref="profileForm"
        :model="profileForm"
        :rules="profileRules"
        label-width="120px"
      >
        <el-form-item label="用户名">
          <el-input
            v-model="userProfile.username"
            disabled
          />
        </el-form-item>

        <el-divider content-position="left">基础信息</el-divider>

        <el-form-item
          label="头像"
          prop="avatar"
        >
          <div class="avatar-form-item">
            <el-avatar
              :size="56"
              :src="profileForm.avatar"
              icon="el-icon-user-solid"
            />
            <el-input
              v-model="profileForm.avatar"
              placeholder="请输入头像 URL"
            />
          </div>
        </el-form-item>

        <el-form-item
          label="真实姓名"
          prop="realName"
        >
          <el-input
            v-model="profileForm.realName"
            placeholder="请输入真实姓名"
          />
        </el-form-item>

        <el-form-item
          label="年龄"
          prop="age"
        >
          <el-input-number
            v-model="profileForm.age"
            :min="18"
            :max="100"
            placeholder="请输入年龄"
          />
        </el-form-item>

        <el-form-item
          label="邮箱"
          prop="email"
        >
          <el-input
            v-model="profileForm.email"
            placeholder="请输入邮箱"
          />
        </el-form-item>

        <el-form-item
          label="电话"
          prop="phone"
        >
          <el-input
            v-model="profileForm.phone"
            placeholder="请输入手机号"
          />
        </el-form-item>

        <el-form-item
          label="所在地"
          prop="location"
        >
          <el-input
            v-model="profileForm.location"
            placeholder="请输入所在地，如：上海"
          />
        </el-form-item>

        <el-divider content-position="left">求职意向</el-divider>

        <el-form-item
          label="求职意向"
          prop="jobIntention"
        >
          <el-input
            v-model="profileForm.jobIntention"
            type="textarea"
            :rows="3"
            placeholder="请输入求职意向，如：前端工程师、全栈开发等"
          />
        </el-form-item>

        <el-form-item
          label="个人优势"
          prop="personalAdvantage"
        >
          <el-input
            v-model="profileForm.personalAdvantage"
            type="textarea"
            :rows="4"
            placeholder="请输入个人优势，如：熟悉 Vue、TypeScript，有后台系统项目经验"
          />
        </el-form-item>

        <el-divider content-position="left">教育经历</el-divider>

        <div
          v-for="(experience, index) in profileForm.educationExperiences"
          :key="index"
          class="education-experience-card"
        >
          <div class="education-experience-card__header">
            <span>教育经历 {{ index + 1 }}</span>
            <el-button
              type="text"
              class="danger"
              :disabled="profileForm.educationExperiences.length === 1"
              @click="removeEducationExperience(index)"
            >
              删除
            </el-button>
          </div>

          <el-form-item label="学校">
            <el-input
              v-model="experience.school"
              placeholder="请输入学校名称"
            />
          </el-form-item>

          <el-form-item label="学历">
            <el-select
              v-model="experience.education"
              placeholder="请选择学历，同一种学历只能添加一段"
              clearable
            >
              <el-option
                v-for="option in educationOptions"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="专业">
            <el-input
              v-model="experience.major"
              placeholder="请输入专业名称"
            />
          </el-form-item>

          <el-form-item label="时间">
            <div class="education-date-range">
              <el-date-picker
                v-model="experience.admissionDate"
                type="month"
                value-format="yyyy-MM"
                placeholder="入学时间"
              />
              <span class="date-separator">至</span>
              <el-date-picker
                v-model="experience.graduationDate"
                type="month"
                value-format="yyyy-MM"
                placeholder="毕业时间"
              />
            </div>
          </el-form-item>
        </div>

        <el-form-item>
          <el-button
            type="primary"
            plain
            icon="el-icon-plus"
            @click="addEducationExperience"
          >
            添加教育经历
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleUpdate"
          >
            保存修改
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import { UserModule } from '@/store/modules/user'
import { EducationExperience, UpdateProfileRequest } from '@/models'

@Component({
  name: 'UserProfile'
})
export default class extends Vue {
  private profileForm: UpdateProfileRequest = {
    realName: '',
    age: undefined,
    email: '',
    phone: '',
    jobIntention: '',
    avatar: '',
    educationExperiences: [this.createEducationExperience()],
    location: '',
    personalAdvantage: ''
  }

  private loading = false
  private educationOptions = ['专科', '本科', '硕士', '博士', '其他']

  get userProfile() {
    return UserModule.userProfile || {
      username: '',
      id: '',
      createdAt: '',
      updatedAt: ''
    }
  }

  private validateEmail = (rule: any, value: string, callback: Function) => {
    if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      callback(new Error('邮箱格式不正确'))
    } else {
      callback()
    }
  }

  private validatePhone = (rule: any, value: string, callback: Function) => {
    if (value && !/^1[3-9]\d{9}$/.test(value)) {
      callback(new Error('手机号格式不正确'))
    } else {
      callback()
    }
  }

  private profileRules = {
    email: [{ validator: this.validateEmail, trigger: 'blur' }],
    phone: [{ validator: this.validatePhone, trigger: 'blur' }],
    age: [
      { type: 'number', message: '年龄必须为数字', trigger: 'blur' }
    ]
  }

  created() {
    this.loadUserProfile()
  }

  private async loadUserProfile() {
    try {
      await UserModule.GetUserInfo()
      if (UserModule.userProfile) {
        this.profileForm = this.getProfileForm()
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
    }
  }

  private handleUpdate() {
    (this.$refs.profileForm as ElForm).validate(async(valid: boolean) => {
      if (valid && this.validateEducationExperiences()) {
        this.loading = true
        try {
          await UserModule.UpdateProfile(this.getSubmitProfileForm())
          this.$message.success('保存成功！')
        } catch (error) {
          console.error('更新失败:', error)
          this.$message.error('保存失败，请稍后重试')
        } finally {
          this.loading = false
        }
      }
    })
  }

  private handleReset() {
    this.loadUserProfile()
    this.$message.info('已重置为当前保存的信息')
  }

  private getProfileForm(): UpdateProfileRequest {
    const profile = UserModule.userProfile
    return {
      realName: profile?.realName || '',
      age: profile?.age,
      email: profile?.email || '',
      phone: profile?.phone || '',
      jobIntention: profile?.jobIntention || '',
      avatar: profile?.avatar || '',
      educationExperiences: profile?.educationExperiences?.length
        ? profile.educationExperiences.map(experience => ({ ...experience }))
        : [this.createEducationExperience()],
      location: profile?.location || '',
      personalAdvantage: profile?.personalAdvantage || ''
    }
  }

  private createEducationExperience(): EducationExperience {
    return {
      school: '',
      education: '',
      major: '',
      admissionDate: '',
      graduationDate: ''
    }
  }

  private addEducationExperience() {
    this.profileForm.educationExperiences = [
      ...(this.profileForm.educationExperiences || []),
      this.createEducationExperience()
    ]
  }

  private removeEducationExperience(index: number) {
    const experiences = this.profileForm.educationExperiences || []
    if (experiences.length <= 1) return
    this.profileForm.educationExperiences = experiences.filter((_, currentIndex) => currentIndex !== index)
  }

  private getSubmitProfileForm(): UpdateProfileRequest {
    return {
      ...this.profileForm,
      educationExperiences: (this.profileForm.educationExperiences || []).filter(experience =>
        this.hasEducationContent(experience)
      )
    }
  }

  private validateEducationExperiences() {
    const experiences = (this.profileForm.educationExperiences || []).filter(experience =>
      this.hasEducationContent(experience)
    )
    const educationSet = new Set<string>()

    for (const experience of experiences) {
      if (!experience.education) {
        this.$message.error('请为已填写的教育经历选择学历')
        return false
      }

      if (educationSet.has(experience.education)) {
        this.$message.error(`学历“${experience.education}”只能添加一段教育经历`)
        return false
      }

      if (
        experience.admissionDate &&
        experience.graduationDate &&
        experience.graduationDate < experience.admissionDate
      ) {
        this.$message.error('毕业时间不能早于入学时间')
        return false
      }

      educationSet.add(experience.education)
    }

    return true
  }

  private hasEducationContent(experience: EducationExperience) {
    return Boolean(
      experience.school ||
      experience.education ||
      experience.major ||
      experience.admissionDate ||
      experience.graduationDate
    )
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  padding: 20px;

  .profile-card {
    max-width: 800px;
    margin: 0 auto;

    .card-title {
      font-size: 18px;
      font-weight: bold;
    }

    .avatar-form-item {
      display: flex;
      align-items: center;
      gap: 16px;

      .el-input {
        flex: 1;
      }
    }

    .el-select,
    .el-date-editor {
      width: 100%;
    }

    .education-experience-card {
      margin-bottom: 18px;
      padding: 16px 16px 1px;
      border: 1px solid #ebeef5;
      border-radius: 4px;
      background: #fafafa;

      &__header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 16px;
        font-weight: 600;
        color: #303133;

        .danger {
          color: #f56c6c;
        }
      }
    }

    .education-date-range {
      display: flex;
      align-items: center;
      gap: 10px;

      .date-separator {
        color: #909399;
      }
    }
  }
}
</style>
