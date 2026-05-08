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
          placeholder="请输入邮箱(选填)"
        />
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
import { register } from '@/api'
import { RegisterRequest } from '@/models'

@Component({
  name: 'Register'
})
export default class extends Vue {
  private registerForm = {
    username: '',
    password: '',
    confirmPassword: '',
    email: ''
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
    if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      callback(new Error('邮箱格式不正确'))
    } else {
      callback()
    }
  }

  private registerRules = {
    username: [{ validator: this.validateUsername, trigger: 'blur' }],
    password: [{ validator: this.validatePassword, trigger: 'blur' }],
    confirmPassword: [{ validator: this.validateConfirmPassword, trigger: 'blur' }],
    email: [{ validator: this.validateEmail, trigger: 'blur' }]
  }

  private loading = false

  private handleRegister() {
    (this.$refs.registerForm as ElForm).validate(async(valid: boolean) => {
      if (valid) {
        this.loading = true
        try {
          const data: RegisterRequest = {
            username: this.registerForm.username,
            password: this.registerForm.password,
            email: this.registerForm.email || undefined
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
  }
}
</style>
