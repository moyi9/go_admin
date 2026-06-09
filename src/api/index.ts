/**
 * Axios HTTP 请求封装。
 * - 请求拦截器：自动附加 Token 和全局 Loading
 * - 响应拦截器：统一业务错误提示、401 跳登录、403 无权限提示、429 限流提示
 */
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from '@/utils/token'
import { useLoadingStore } from '@/stores/loading'
import type { ApiResponse } from '@/types/api'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL, // 由 .env.development / .env.production 配置
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const loading = useLoadingStore()
  loading.show()
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    const loading = useLoadingStore()
    loading.hide()
    const body = response.data as ApiResponse
    if (body.code !== 0) {
      ElMessage.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message))
    }
    return response
  },
  (error) => {
    const loading = useLoadingStore()
    loading.hide()
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        removeToken()
        window.location.href = '/#/login'
        ElMessage.error('登录已过期，请重新登录')
      } else if (status === 403) {
        ElMessage.error('无权限访问')
      } else if (status === 429) {
        ElMessage.error('请求太频繁，请稍后再试')
      } else {
        ElMessage.error(error.response.data?.message || '服务器错误')
      }
    } else {
      ElMessage.error('网络错误，请检查连接')
    }
    return Promise.reject(error)
  },
)

export default request
