/** 操作日志 API */
import request from './index'
import type { PaginatedResponse } from '@/types/api'
import type { AuditLog } from '@/types/model'

/** ListAuditLogsParams 日志列表查询参数 */
export interface ListAuditLogsParams {
  action?: string
  resource?: string
  keyword?: string
  since?: string
  until?: string
}

/** listAuditLogsApi 分页获取操作日志列表 */
export async function listAuditLogsApi(page: number, pageSize: number, params?: ListAuditLogsParams) {
  const res = await request.get<PaginatedResponse<AuditLog>>('/audit-logs', {
    params: { page, page_size: pageSize, ...params },
  })
  return res.data
}
