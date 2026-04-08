import { defineStore } from 'pinia'
import { login as loginApi, register as registerApi } from '@/api/auth'
import { getUser } from '@/api/user'
import { storage } from '@/utils/storage'
import type { User } from '@/types'

interface AuthState {
  token: string | null
  userInfo: User | null
  isLoggedIn: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: storage.get<string>('token'),
    userInfo: storage.get<User>('userInfo'),
    isLoggedIn: !!storage.get<string>('token')
  }),

  getters: {
    userEmail: (state) => state.userInfo?.email || ''
  },

  actions: {
    async login(email: string, password: string) {
      try {
        const res = await loginApi({ email, password })
        const token = (res.data as any).token
        this.token = token
        this.isLoggedIn = true
        storage.set('token', token)
        await this.fetchUserInfo()
        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '登录失败' }
      }
    },

    async register(email: string, password: string) {
      try {
        const res = await registerApi({ email, password })
        const token = (res.data as any).token
        this.token = token
        this.isLoggedIn = true
        storage.set('token', token)
        await this.fetchUserInfo()
        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '注册失败' }
      }
    },

    async fetchUserInfo() {
      try {
        const res = await getUser()
        this.userInfo = res.data as User
        storage.set('userInfo', this.userInfo)
      } catch (error) {
        console.error('Failed to fetch user info:', error)
      }
    },

    logout() {
      this.token = null
      this.userInfo = null
      this.isLoggedIn = false
      storage.remove('token')
      storage.remove('userInfo')
      uni.reLaunch({ url: '/pages/auth/login' })
    },

    updateApiKeyStatus(set: boolean) {
      if (this.userInfo) {
        this.userInfo.api_key_set = set
        storage.set('userInfo', this.userInfo)
      }
    }
  }
})
