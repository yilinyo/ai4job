<template>
  <div class="register-container">
    <section class="auth-hero">
      <div class="brand-panel">
        <div class="brand-badge">
          AI For Career
        </div>
        <h1>创建账号，沉淀你的求职资产</h1>
        <p>
          注册后邮箱会进入个人信息页，后续创建简历时可自动复用基础资料，让每次投递更轻一些。
        </p>
        <div class="brand-steps">
          <div>
            <strong>01</strong>
            <span>邮箱验证</span>
          </div>
          <div>
            <strong>02</strong>
            <span>资料沉淀</span>
          </div>
          <div>
            <strong>03</strong>
            <span>简历复用</span>
          </div>
        </div>
      </div>

      <div class="auth-card">
        <div class="title-container">
          <p class="eyebrow">
            Start here
          </p>
          <h3 class="title">
            用户注册
          </h3>
        </div>

        <el-form
          ref="registerForm"
          :model="registerForm"
          :rules="registerRules"
          class="register-form"
          label-position="top"
        >
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

          <el-button
            :loading="loading"
            type="primary"
            class="submit-button"
            @click="handleRegister"
          >
            注册
          </el-button>

          <div class="auth-switch">
            <span>已有账号？</span>
            <el-button
              type="text"
              @click="goToLogin"
            >
              去登录
            </el-button>
          </div>
        </el-form>
      </div>
    </section>
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
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  overflow: auto;
  padding: 48px;
  background:
    radial-gradient(circle at 78% 16%, rgba(64, 158, 255, 0.2), transparent 26%),
    linear-gradient(135deg, #f7fbff 0%, #eef5ff 48%, #ffffff 100%);

  .auth-hero {
    display: grid;
    grid-template-columns: minmax(320px, 1fr) 480px;
    gap: 56px;
    align-items: center;
    width: 100%;
    max-width: 1120px;
  }

  .brand-panel {
    color: #17233d;

    .brand-badge {
      display: inline-flex;
      align-items: center;
      height: 36px;
      padding: 0 16px;
      margin-bottom: 28px;
      color: #1d63ed;
      font-weight: 600;
      background: rgba(51, 112, 255, 0.1);
      border: 1px solid rgba(51, 112, 255, 0.18);
      border-radius: 999px;
    }

    h1 {
      max-width: 560px;
      margin: 0;
      font-size: 44px;
      line-height: 1.18;
      letter-spacing: -1px;
    }

    p {
      max-width: 540px;
      margin: 24px 0 36px;
      color: #5c6b82;
      font-size: 17px;
      line-height: 1.8;
    }
  }

  .brand-steps {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 14px;
    max-width: 520px;

    div {
      padding: 18px;
      background: rgba(255, 255, 255, 0.72);
      border: 1px solid rgba(220, 229, 244, 0.9);
      border-radius: 18px;
      box-shadow: 0 16px 36px rgba(31, 45, 61, 0.08);
    }

    strong {
      display: block;
      color: #1d63ed;
      font-size: 18px;
      margin-bottom: 8px;
    }

    span {
      color: #6b778c;
      font-size: 13px;
    }
  }

  .auth-card {
    padding: 40px 42px;
    background: rgba(255, 255, 255, 0.95);
    border: 1px solid rgba(221, 228, 240, 0.88);
    border-radius: 26px;
    box-shadow: 0 24px 70px rgba(31, 45, 61, 0.14);
    backdrop-filter: blur(12px);
  }

  .title-container {
    margin-bottom: 24px;

    .eyebrow {
      margin: 0 0 8px;
      color: #3370ff;
      font-size: 13px;
      font-weight: 600;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    .title {
      font-size: 28px;
      color: #17233d;
      margin: 0;
      font-weight: bold;
    }
  }

  .register-form {
    width: 100%;

    ::v-deep .el-form-item {
      margin-bottom: 18px;
    }

    ::v-deep .el-form-item__label {
      padding-bottom: 6px;
      color: #42526e;
      font-weight: 600;
      line-height: 1.3;
    }

    ::v-deep .el-input__inner {
      height: 44px;
      background: #f8fafc;
      border-color: #e5eaf3;
      border-radius: 10px;
      color: #1f2d3d;
    }

    .email-code-row {
      display: flex;
      gap: 10px;

      .el-input {
        flex: 1;
      }

      .el-button {
        width: 122px;
        border-radius: 10px;
      }
    }

    .submit-button {
      width: 100%;
      height: 46px;
      margin-top: 4px;
      border-radius: 10px;
      font-size: 16px;
      font-weight: 600;
    }

    .auth-switch {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 4px;
      color: #6b778c;
      margin-top: 18px;
    }
  }

  @media only screen and (max-width: 960px) {
    padding: 32px 20px;

    .auth-hero {
      grid-template-columns: 1fr;
      gap: 28px;
    }

    .brand-panel {
      text-align: center;

      h1,
      p {
        max-width: none;
      }

      h1 {
        font-size: 34px;
      }
    }

    .brand-steps {
      margin: 0 auto;
    }
  }

  @media only screen and (max-width: 520px) {
    .auth-card {
      padding: 28px 22px;
      border-radius: 20px;
    }

    .brand-steps {
      grid-template-columns: 1fr;
    }

    .email-code-row {
      flex-direction: column;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
