/** 权限管理 API */
import request from './index'
import type { ApiResponse, PaginatedResponse } from '@/types/api'
import type { Permission } from '@/types/model'

/** listPermissionsApi 分页获取权限列表 */
export async function listPermissionsApi(page: number, pageSize: number, params?: Record<string, string>) {
  const res = await request.get<PaginatedResponse<Permission>>('/permissions', { params: { page, page_size: pageSize, ...params } })
  return res.data
}

/** createPermissionApi 创建权限 */
export async function createPermissionApi(data: Record<string, unknown>) {
  const res = await request.post<ApiResponse<Permission>>('/permissions', data)
  return res.data.data!
}

/** updatePermissionApi 更新权限 */
export async function updatePermissionApi(id: number, data: Record<string, unknown>) {
  const res = await request.put<ApiResponse<Permission>>(`/permissions/${id}`, data)
  return res.data.data!
}

/** deletePermissionApi 删除权限 */
export async function deletePermissionApi(id: number) {
  await request.delete(`/permissions/${id}`)
}

/** batchDeletePermissionsApi 批量删除权限 */
export async function batchDeletePermissionsApi(ids: number[]) {
  await request.delete('/permissions/batch', { data: { ids } })
}
