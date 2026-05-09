<template>
  <div class="register-container">
    <el-form
      ref="registerForm"
      :model="registerForm"
      :rules="registerRules"
      class="register-form"
      label-width="100px"
    >
      <div class="title-container">
        <h3 class="title">
          用户注册
        </h3>
      </div>

      <el-form-item
        label="用户名"
        prop="username"
      >
        <el-input
          v-model="registerForm.username"
          placeholder="请输入用户名"
        />
      </el-form-item>

      <el-form-item
        label="密码"
        prop="password"
      >
        <el-input
          v-model="registerForm.password"
          type="password"
          placeholder="请输入密码(至少6位)"
        />
      </el-form-item>

      <el-form-item
        label="确认密码"
        prop="confirmPassword"
      >
        <el-input
          v-model="registerForm.confirmPassword"
          type="password"
          placeholder="请再次输入密码"
        />
      </el-form-item>

      <el-form-item
        label="邮箱"
        prop="email"
      >
        <el-input
          v-model="registerForm.email"
          placeholder="请输入邮箱"
          @input="handleEmailChange"
        />
      </el-form-item>

      <el-form-item
        label="验证码"
        prop="emailCode"
      >
        <div class="email-code-row">
          <el-input
            v-model="registerForm.emailCode"
            placeholder="请输入邮箱验证码"
          />
          <el-button
            :loading="emailCodeLoading"
            :disabled="emailCodeCountdown > 0"
            @click="handleSendEmailCode"
          >
            {{ emailCodeButtonText }}
          </el-button>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button
          :loading="loading"
          type="primary"
          style="width:100%"
          @click="handleRegister"
        >
          注册
        </el-button>
      </el-form-item>

      <el-form-item>
        <el-button
          style="width:100%"
          @click="goToLogin"
        >
          已有账号？去登录
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Form as ElForm } from 'element-ui'
import { register, sendEmailCode } from '@/api'
import { RegisterRequest } from '@/models'

@Component({
  name: 'Register'
})
export default class extends Vue {
  private registerForm = {
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    emailCode: ''
  }

  private validateUsername = (rule: any, value: string, callback: Function) => {
    if (!value || value.length < 3) {
      callback(new Error('用户名至少3个字符'))
    } else {
      callback()
    }
  }

  private validatePassword = (rule: any, value: string, callback: Function) => {
    if (value.length < 6) {
      callback(new Error('密码至少6位'))
    } else {
      callback()
    }
  }

  private validateConfirmPassword = (rule: any, value: string, callback: Function) => {
    if (value !== this.registerForm.password) {
      callback(new Error('两次输入的密码不一致'))
    } else {
      callback()
    }
  }

  private validateEmail = (rule: any, value: string, callback: Function) => {
    if (!value) {
      callback(new Error('请输入邮箱'))
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      callback(new Error('邮箱格式不正确'))
    } else {
      callback()
    }
  }

  private validateEmailCode = (rule: any, value: string, callback: Function) => {
    if (!value) {
      callback(new Error('请输入邮箱验证码'))
    } else {
      callback()
    }
  }

  private registerRules = {
    username: [{ validator: this.validateUsername, trigger: 'blur' }],
    password: [{ validator: this.validatePassword, trigger: 'blur' }],
    confirmPassword: [{ validator: this.validateConfirmPassword, trigger: 'blur' }],
    email: [{ validator: this.validateEmail, trigger: 'blur' }],
    emailCode: [{ validator: this.validateEmailCode, trigger: 'blur' }]
  }

  private loading = false
  private emailCodeLoading = false
  private emailCodeCountdown = 0
  private emailCodeTimer: number | null = null

  get emailCodeButtonText() {
    return this.emailCodeCountdown > 0 ? `${this.emailCodeCountdown}s 后重发` : '发送验证码'
  }

  beforeDestroy() {
    this.clearEmailCodeTimer()
  }

  private handleEmailChange() {
    this.registerForm.emailCode = ''
  }

  private handleSendEmailCode() {
    (this.$refs.registerForm as ElForm).validateField('email', async(errorMessage: string) => {
      if (errorMessage) return

      this.emailCodeLoading = true
      try {
        const { data } = await sendEmailCode({ email: this.registerForm.email })
        this.$message.success(`验证码已发送${data.mockCode ? `，Mock 验证码：${data.mockCode}` : ''}`)
        this.startEmailCodeCountdown(data.expiresIn || 60)
      } catch (error) {
        console.error('发送验证码失败:', error)
        this.$message.error('发送验证码失败，请稍后重试')
      } finally {
        this.emailCodeLoading = false
      }
    })
  }

  private startEmailCodeCountdown(seconds: number) {
    this.clearEmailCodeTimer()
    this.emailCodeCountdown = Math.min(seconds, 60)
    this.emailCodeTimer = window.setInterval(() => {
      this.emailCodeCountdown -= 1
      if (this.emailCodeCountdown <= 0) {
        this.clearEmailCodeTimer()
      }
    }, 1000)
  }

  private clearEmailCodeTimer() {
    if (this.emailCodeTimer !== null) {
      window.clearInterval(this.emailCodeTimer)
      this.emailCodeTimer = null
    }
    if (this.emailCodeCountdown < 0) {
      this.emailCodeCountdown = 0
    }
  }

  private handleRegister() {
    (this.$refs.registerForm as ElForm).validate(async(valid: boolean) => {
      if (valid) {
        this.loading = true
        try {
          const data: RegisterRequest = {
            username: this.registerForm.username,
            password: this.registerForm.password,
            email: this.registerForm.email,
            emailCode: this.registerForm.emailCode
          }
          await register(data)
          this.$message.success('注册成功！')
          // 注册成功后跳转到登录页
          setTimeout(() => {
            this.$router.push('/login')
          }, 1000)
        } catch (error) {
          console.error('注册失败:', error)
          this.$message.error('注册失败，请稍后重试')
        } finally {
          this.loading = false
        }
      }
    })
  }

  private goToLogin() {
    this.$router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;

  .register-form {
    width: 450px;
    padding: 40px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

    .title-container {
      .title {
        font-size: 26px;
        color: #333;
        margin: 0 auto 30px;
        text-align: center;
        font-weight: bold;
      }
    }

    .email-code-row {
      display: flex;
      gap: 10px;

      .el-input {
        flex: 1;
      }

      .el-button {
        width: 120px;
      }
    }
  }
}
</style>
