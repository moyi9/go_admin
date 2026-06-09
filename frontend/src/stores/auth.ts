/** 认证状态管理：Token 存取、用户信息获取、登录/登出 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/model'
import { loginApi, getMeApi } from '@/api/auth'
import { setToken, removeToken, getToken } from '@/utils/token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken()) // JWT Token
  const user = ref<User | null>(null)           // 当前登录用户信息

  /** login 登录并存储 Token，remember 控制是否持久化 */
  async function login(username: string, password: string, remember: boolean) {
    const res = await loginApi(username, password)
    setToken(res.access_token, remember)
    token.value = res.access_token
    user.value = res.user
  }

  /** getUserInfo 从后端获取当前用户信息，失败时清除 Token */
  async function getUserInfo() {
    if (!token.value) return null
    try {
      const u = await getMeApi()
      user.value = u
      return u
    } catch {
      removeToken()
      token.value = null
      return null
    }
  }

  /** updateUser 更新本地用户信息（编辑个人资料后同步 store） */
  function updateUser(partial: Partial<User>) {
    if (user.value) {
      user.value = { ...user.value, ...partial }
    }
  }

  /** logout 清除 Token 和用户信息 */
  function logout() {
    removeToken()
    token.value = null
    user.value = null
  }

  return { token, user, login, getUserInfo, updateUser, logout }
})
