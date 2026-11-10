import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'

const TOKEN_KEY = 'blog_token'
const USER_KEY = 'blog_user'

// 安全读取本地缓存中的用户信息
function readCachedUser() {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    user: readCachedUser()
  }),

  getters: {
    isLogin: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    displayName: (state) => state.user?.nickname || state.user?.username || '游客'
  },

  actions: {
    setToken(token) {
      this.token = token || ''
      if (token) localStorage.setItem(TOKEN_KEY, token)
      else localStorage.removeItem(TOKEN_KEY)
    },

    setUser(user) {
      this.user = user || null
      if (user) localStorage.setItem(USER_KEY, JSON.stringify(user))
      else localStorage.removeItem(USER_KEY)
    },

    // 登录：写入 token 与用户信息
    async login(payload) {
      const data = await authApi.login(payload)
      this.setToken(data.token)
      this.setUser(data.user)
      return data
    },

    // 拉取当前用户信息
    async fetchMe() {
      const data = await authApi.me()
      this.setUser(data)
      return data
    },

    // 退出登录：清空本地状态
    logout() {
      this.setToken('')
      this.setUser(null)
    }
  }
})
