/** 用户管理 API */
import request from './index'
import type { ApiResponse, PaginatedResponse } from '@/types/api'
import type { User } from '@/types/model'

/** ListUsersParams 用户列表查询参数 */
export interface ListUsersParams {
  username?: string
  nickname?: string
  phone?: string
  email?: string
  gender?: string
  status?: string
  keyword?: string
}

/** listUsersApi 分页获取用户列表 */
export async function listUsersApi(page: number, pageSize: number, params?: ListUsersParams) {
  const res = await request.get<PaginatedResponse<User>>('/users', {
    params: { page, page_size: pageSize, ...params },
  })
  return res.data
}

/** createUserApi 创建用户 */
export async function createUserApi(data: {
  username: string
  nickname?: string
  email: string
  phone?: string
  gender?: string
  password: string
  status: string
  role_ids?: number[]
}) {
  const res = await request.post<ApiResponse<User>>('/users', data)
  return res.data.data!
}

/** updateUserApi 更新用户信息 */
export async function updateUserApi(id: number, data: {
  username: string
  nickname?: string
  email: string
  phone?: string
  gender?: string
  status: string
  role_ids?: number[]
}) {
  const res = await request.put<ApiResponse<User>>(`/users/${id}`, data)
  return res.data.data!
}

/** deleteUserApi 删除用户 */
export async function deleteUserApi(id: number) {
  await request.delete(`/users/${id}`)
}

/** batchDeleteUsersApi 批量删除用户 */
export async function batchDeleteUsersApi(ids: number[]) {
  await request.delete('/users/batch', { data: { ids } })
}

/** assignRolesApi 为用户分配角色 */
export async function assignRolesApi(userId: number, roleIds: number[]) {
  await request.post(`/users/${userId}/roles`, { role_ids: roleIds })
}
