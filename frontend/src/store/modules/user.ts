import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import { login, logout, getUserProfile, updateUserProfile } from '@/api/users'
import { getToken, setToken, removeToken } from '@/utils/cookies'
import router, { resetRouter } from '@/router'
import { TagsViewModule } from './tags-view'
import { PermissionModule } from './permission'
import { UserProfile, UpdateProfileRequest } from '@/models'
import store from '@/store'

export interface IUserState {
  token: string
  roles: string[]
  userProfile: UserProfile | null
}

@Module({ dynamic: true, store, name: 'user' })
class User extends VuexModule implements IUserState {
  public token = getToken() || ''
  public roles: string[] = []
  public userProfile: UserProfile | null = null

  // 兼容原模板组件的属性访问
  get name() {
    return this.userProfile?.realName || this.userProfile?.username || ''
  }

  get avatar() {
    return this.userProfile?.avatar || ''
  }

  get email() {
    return this.userProfile?.email || ''
  }

  @Mutation
  private SET_TOKEN(token: string) {
    this.token = token
  }

  @Mutation
  private SET_USER_PROFILE(profile: UserProfile | null) {
    this.userProfile = profile
  }

  @Mutation
  private SET_ROLES(roles: string[]) {
    this.roles = roles
  }

  @Action
  public async Login(userInfo: { username: string, password: string}) {
    let { username, password } = userInfo
    username = username.trim()
    const { data } = await login({ username, password })
    setToken(data.token)
    this.SET_TOKEN(data.token)
    this.SET_USER_PROFILE(data.userInfo)
  }

  @Action
  public ResetToken() {
    removeToken()
    this.SET_TOKEN('')
    this.SET_ROLES([])
    this.SET_USER_PROFILE(null)
  }

  @Action
  public async GetUserInfo() {
    if (this.token === '') {
      throw Error('GetUserInfo: token is undefined!')
    }
    const { data } = await getUserProfile()
    if (!data) {
      throw Error('Verification failed, please Login again.')
    }
    this.SET_USER_PROFILE(data)
    this.SET_ROLES(data.roles || ['admin'])
  }

  @Action
  public async UpdateProfile(profileData: UpdateProfileRequest) {
    const { data } = await updateUserProfile(profileData)
    if (data) {
      this.SET_USER_PROFILE(data)
    }
  }

  @Action
  public async ChangeRoles(role: string) {
    const token = role + '-token'
    this.SET_TOKEN(token)
    setToken(token)
    await this.GetUserInfo()
    resetRouter()
    PermissionModule.GenerateRoutes(this.roles)
    PermissionModule.dynamicRoutes.forEach(route => {
      router.addRoute(route)
    })
  }

  @Action
  public async LogOut() {
    if (this.token === '') {
      throw Error('LogOut: token is undefined!')
    }
    await logout()
    removeToken()
    resetRouter()

    // Reset visited views and cached views
    TagsViewModule.delAllViews()
    this.SET_TOKEN('')
    this.SET_ROLES([])
    this.SET_USER_PROFILE(null)
  }
}

export const UserModule = getModule(User)
