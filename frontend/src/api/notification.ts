/** 通知中心 API */
import request from './index'
import type { PaginatedResponse } from '@/types/api'
import type { Notification } from '@/types/model'

/** SendNotificationParams 发送通知参数 */
export interface SendNotificationParams {
  type: string
  title: string
  content: string
  target_user_id?: number
}

/** listNotificationsApi 分页获取当前用户通知列表 */
export async function listNotificationsApi(page: number, pageSize: number) {
  const res = await request.get<PaginatedResponse<Notification>>('/notifications', {
    params: { page, page_size: pageSize },
  })
  return res.data
}

/** countUnreadApi 获取未读通知数 */
export async function countUnreadApi() {
  const res = await request.get<{ code: number; data: { count: number } }>('/notifications/unread')
  return res.data.data.count
}

/** listUnreadNotificationsApi 获取最新 N 条未读通知 */
export async function listUnreadNotificationsApi() {
  const res = await request.get<{ code: number; data: Notification[] }>('/notifications/unread-list')
  return res.data.data
}

/** sendNotificationApi 发送通知 */
export async function sendNotificationApi(data: SendNotificationParams) {
  await request.post('/notifications', data)
}

/** markAsReadApi 标记单条通知已读 */
export async function markAsReadApi(id: number) {
  await request.put(`/notifications/${id}/read`)
}

/** markAllAsReadApi 全部标记已读 */
export async function markAllAsReadApi() {
  await request.put('/notifications/read-all')
}
