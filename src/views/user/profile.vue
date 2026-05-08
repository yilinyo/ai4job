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
    jobIntention: ''
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
      // 加载用户信息到表单
      if (UserModule.userProfile) {
        this.profileForm = {
          realName: UserModule.userProfile.realName || '',
          age: UserModule.userProfile.age,
          email: UserModule.userProfile.email || '',
          phone: UserModule.userProfile.phone || '',
          jobIntention: UserModule.userProfile.jobIntention || ''
        }
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
  }
}
</style>
