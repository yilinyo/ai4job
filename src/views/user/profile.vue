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

        <el-divider content-position="left">教育信息</el-divider>

        <el-form-item
          label="学校"
          prop="school"
        >
          <el-input
            v-model="profileForm.school"
            placeholder="请输入学校名称"
          />
        </el-form-item>

        <el-form-item
          label="学历"
          prop="education"
        >
          <el-select
            v-model="profileForm.education"
            placeholder="请选择学历"
            clearable
          >
            <el-option label="专科" value="专科" />
            <el-option label="本科" value="本科" />
            <el-option label="硕士" value="硕士" />
            <el-option label="博士" value="博士" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>

        <el-form-item
          label="专业"
          prop="major"
        >
          <el-input
            v-model="profileForm.major"
            placeholder="请输入专业名称"
          />
        </el-form-item>

        <el-form-item
          label="毕业时间"
          prop="graduationDate"
        >
          <el-date-picker
            v-model="profileForm.graduationDate"
            type="month"
            value-format="yyyy-MM"
            placeholder="请选择毕业年月"
          />
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
import { UpdateProfileRequest } from '@/models'

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
    school: '',
    education: '',
    major: '',
    graduationDate: '',
    location: '',
    personalAdvantage: ''
  }

  private loading = false

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
      if (valid) {
        this.loading = true
        try {
          await UserModule.UpdateProfile(this.profileForm)
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
      school: profile?.school || '',
      education: profile?.education || '',
      major: profile?.major || '',
      graduationDate: profile?.graduationDate || '',
      location: profile?.location || '',
      personalAdvantage: profile?.personalAdvantage || ''
    }
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
  }
}
</style>
