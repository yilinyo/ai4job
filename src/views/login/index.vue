<template>
  <div class="login-container">
    <section class="auth-hero">
      <div class="brand-panel">
        <div class="brand-badge">
          AI For Career
        </div>
        <h1>让求职流程像项目一样可管理</h1>
        <p>
          集中维护简历、岗位、投递和面试进展，帮助你用更清晰的节奏推进每一次机会。
        </p>
        <div class="brand-stats">
          <div>
            <strong>Resume</strong>
            <span>简历版本</span>
          </div>
          <div>
            <strong>Pipeline</strong>
            <span>流程追踪</span>
          </div>
          <div>
            <strong>Profile</strong>
            <span>资料复用</span>
          </div>
        </div>
      </div>

      <div class="auth-card">
        <div class="title-container">
          <div>
            <p class="eyebrow">
              Welcome back
            </p>
            <h3 class="title">
              {{ $t('login.title') }}
            </h3>
          </div>
          <lang-select class="set-language" />
        </div>

        <el-form
          ref="loginForm"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          autocomplete="on"
          label-position="left"
        >
          <el-form-item prop="username">
            <span class="svg-container">
              <svg-icon name="user" />
            </span>
            <el-input
              ref="username"
              v-model="loginForm.username"
              :placeholder="$t('login.username')"
              name="username"
              type="text"
              tabindex="1"
              autocomplete="on"
            />
          </el-form-item>

          <el-tooltip
            v-model="capsTooltip"
            content="Caps lock is On"
            placement="right"
            manual
          >
            <el-form-item prop="password">
              <span class="svg-container">
                <svg-icon name="password" />
              </span>
              <el-input
                :key="passwordType"
                ref="password"
                v-model="loginForm.password"
                :type="passwordType"
                :placeholder="$t('login.password')"
                name="password"
                tabindex="2"
                autocomplete="on"
                @keyup.native="checkCapslock"
                @blur="capsTooltip = false"
                @keyup.enter.native="handleLogin"
              />
              <span
                class="show-pwd"
                @click="showPwd"
              >
                <svg-icon :name="passwordType === 'password' ? 'eye-off' : 'eye-on'" />
              </span>
            </el-form-item>
          </el-tooltip>

          <el-button
            :loading="loading"
            type="primary"
            class="submit-button"
            @click.native.prevent="handleLogin"
          >
            {{ $t('login.logIn') }}
          </el-button>

          <div class="auth-switch">
            <span>还没有账号？</span>
            <el-button
              type="text"
              @click="goToRegister"
            >
              立即注册
            </el-button>
          </div>

          <div class="demo-box">
            <div class="tips">
              <span>{{ $t('login.username') }} : admin </span>
              <span>{{ $t('login.password') }} : {{ $t('login.any') }} </span>
            </div>
            <div class="tips">
              <span>{{ $t('login.username') }} : editor </span>
              <span>{{ $t('login.password') }} : {{ $t('login.any') }} </span>
            </div>

            <el-button
              class="thirdparty-button"
              type="text"
              @click="showDialog=true"
            >
              {{ $t('login.thirdparty') }}
            </el-button>
          </div>
        </el-form>
      </div>
    </section>

    <el-dialog
      :title="$t('login.thirdparty')"
      :visible.sync="showDialog"
    >
      {{ $t('login.thirdpartyTips') }}
      <br>
      <br>
      <br>
      <social-sign />
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { Route } from 'vue-router'
import { Dictionary } from 'vue-router/types/router'
import { Form as ElForm, Input } from 'element-ui'
import { UserModule } from '@/store/modules/user'
import { isValidUsername } from '@/utils/validate'
import LangSelect from '@/components/LangSelect/index.vue'
import SocialSign from './components/SocialSignin.vue'

@Component({
  name: 'Login',
  components: {
    LangSelect,
    SocialSign
  }
})
export default class extends Vue {
  private validateUsername = (rule: any, value: string, callback: Function) => {
    if (!isValidUsername(value)) {
      callback(new Error('Please enter the correct user name'))
    } else {
      callback()
    }
  }

  private validatePassword = (rule: any, value: string, callback: Function) => {
    if (value.length < 6) {
      callback(new Error('The password can not be less than 6 digits'))
    } else {
      callback()
    }
  }

  private loginForm = {
    username: 'admin',
    password: '111111'
  }

  private loginRules = {
    username: [{ validator: this.validateUsername, trigger: 'blur' }],
    password: [{ validator: this.validatePassword, trigger: 'blur' }]
  }

  private passwordType = 'password'
  private loading = false
  private showDialog = false
  private capsTooltip = false
  private redirect?: string
  private otherQuery: Dictionary<string> = {}

  @Watch('$route', { immediate: true })
  private onRouteChange(route: Route) {
    // TODO: remove the "as Dictionary<string>" hack after v4 release for vue-router
    // See https://github.com/vuejs/vue-router/pull/2050 for details
    const query = route.query as Dictionary<string>
    if (query) {
      this.redirect = query.redirect
      this.otherQuery = this.getOtherQuery(query)
    }
  }

  mounted() {
    if (this.loginForm.username === '') {
      (this.$refs.username as Input).focus()
    } else if (this.loginForm.password === '') {
      (this.$refs.password as Input).focus()
    }
  }

  private checkCapslock(e: KeyboardEvent) {
    const { key } = e
    this.capsTooltip = key !== null && key.length === 1 && (key >= 'A' && key <= 'Z')
  }

  private showPwd() {
    if (this.passwordType === 'password') {
      this.passwordType = ''
    } else {
      this.passwordType = 'password'
    }
    this.$nextTick(() => {
      (this.$refs.password as Input).focus()
    })
  }

  private handleLogin() {
    (this.$refs.loginForm as ElForm).validate(async(valid: boolean) => {
      if (valid) {
        this.loading = true
        await UserModule.Login(this.loginForm)
        this.$router.push({
          path: this.redirect || '/',
          query: this.otherQuery
        }).catch(err => {
          console.warn(err)
        })
        // Just to simulate the time of the request
        setTimeout(() => {
          this.loading = false
        }, 0.5 * 1000)
      } else {
        return false
      }
    })
  }

  private goToRegister() {
    this.$router.push('/register')
  }

  private getOtherQuery(query: Dictionary<string>) {
    return Object.keys(query).reduce((acc, cur) => {
      if (cur !== 'redirect') {
        acc[cur] = query[cur]
      }
      return acc
    }, {} as Dictionary<string>)
  }
}
</script>

<style lang="scss">
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      height: 47px;
      background: #f8fafc;
      border: 1px solid #e5eaf3;
      border-radius: 10px;
      padding: 12px 5px 12px 15px;
      color: #1f2d3d;
      caret-color: #3370ff;
      -webkit-appearance: none;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px #f8fafc inset !important;
        -webkit-text-fill-color: #1f2d3d !important;
      }
    }
  }

  .el-form-item {
    border: 0;
    background: transparent;
    border-radius: 10px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  min-height: 100vh;
  overflow: auto;
  background:
    radial-gradient(circle at 18% 20%, rgba(64, 158, 255, 0.22), transparent 28%),
    linear-gradient(135deg, #eef5ff 0%, #f7fbff 45%, #ffffff 100%);
  padding: 48px;

  .auth-hero {
    display: grid;
    grid-template-columns: minmax(320px, 1fr) 460px;
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
      font-size: 46px;
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

  .brand-stats {
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
      font-size: 16px;
      margin-bottom: 8px;
    }

    span {
      color: #6b778c;
      font-size: 13px;
    }
  }

  .auth-card {
    padding: 42px;
    background: rgba(255, 255, 255, 0.94);
    border: 1px solid rgba(221, 228, 240, 0.88);
    border-radius: 26px;
    box-shadow: 0 24px 70px rgba(31, 45, 61, 0.14);
    backdrop-filter: blur(12px);
  }

  .login-form {
    position: relative;
    width: 100%;
    max-width: 100%;
    padding: 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .submit-button {
    width: 100%;
    margin: 6px 0 18px;
    height: 46px;
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
    margin-bottom: 24px;
  }

  .demo-box {
    position: relative;
    padding: 16px;
    background: #f7f9fc;
    border: 1px solid #edf1f7;
    border-radius: 14px;
  }

  .tips {
    font-size: 13px;
    color: #6b778c;
    margin-bottom: 8px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: #8a97a8;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    position: relative;
    margin-bottom: 28px;

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
      text-align: left;
      font-weight: bold;
    }

    .set-language {
      color: #3370ff;
      font-size: 18px;
      cursor: pointer;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: #8a97a8;
    cursor: pointer;
    user-select: none;
  }

  .thirdparty-button {
    position: absolute;
    right: 12px;
    bottom: 8px;
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

    .brand-stats {
      margin: 0 auto;
    }
  }

  @media only screen and (max-width: 520px) {
    .auth-card {
      padding: 28px 22px;
      border-radius: 20px;
    }

    .brand-stats {
      grid-template-columns: 1fr;
    }

    .thirdparty-button {
      position: static;
      padding-left: 0;
    }
  }
}
</style>
