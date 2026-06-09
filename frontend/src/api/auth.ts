/** 认证相关 API */
import request from './index'
import type { ApiResponse, LoginResult } from '@/types/api'
import type { User } from '@/types/model'

/** loginApi 用户登录，返回 JWT Token 和用户信息 */
export async function loginApi(username: string, password: string) {
  const res = await request.post<ApiResponse<LoginResult>>('/auth/login', { username, password })
  return res.data.data!
}

/** getMeApi 获取当前登录用户信息 */
export async function getMeApi() {
  const res = await request.get<ApiResponse<User>>('/auth/me')
  return res.data.data!
}

/** updateProfileApi 更新当前用户个人信息 */
export async function updateProfileApi(data: {
  nickname?: string
  email?: string
  phone?: string
  gender?: string
  avatar_url?: string
}) {
  const res = await request.put<ApiResponse<User>>('/auth/profile', data)
  return res.data.data!
}

/** updatePasswordApi 修改当前用户密码 */
export async function updatePasswordApi(data: {
  current_password: string
  new_password: string
}) {
  await request.put('/auth/password', data)
}
