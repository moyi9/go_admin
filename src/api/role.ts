/** 角色管理 API */
import request from './index'
import type { ApiResponse, PaginatedResponse } from '@/types/api'
import type { Role } from '@/types/model'

/** ListRolesParams 角色列表查询参数 */
export interface ListRolesParams {
  name?: string
  code?: string
  status?: string
}

/** listRolesApi 分页获取角色列表 */
export async function listRolesApi(page: number, pageSize: number, params?: ListRolesParams) {
  const res = await request.get<PaginatedResponse<Role>>('/roles', {
    params: { page, page_size: pageSize, ...params },
  })
  return res.data
}

/** createRoleApi 创建角色 */
export async function createRoleApi(data: Record<string, unknown>) {
  const res = await request.post<ApiResponse<Role>>('/roles', data)
  return res.data.data!
}

/** updateRoleApi 更新角色 */
export async function updateRoleApi(id: number, data: Record<string, unknown>) {
  const res = await request.put<ApiResponse<Role>>(`/roles/${id}`, data)
  return res.data.data!
}

/** deleteRoleApi 删除角色 */
export async function deleteRoleApi(id: number) {
  await request.delete(`/roles/${id}`)
}

/** batchDeleteRolesApi 批量删除角色 */
export async function batchDeleteRolesApi(ids: number[]) {
  await request.delete('/roles/batch', { data: { ids } })
}

/** assignPermissionsApi 为角色分配权限 */
export async function assignPermissionsApi(roleId: number, permissionIds: number[]) {
  await request.post(`/roles/${roleId}/permissions`, { permission_ids: permissionIds })
}
