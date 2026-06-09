/** 通知状态管理：未读数、已读操作 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { countUnreadApi, markAsReadApi, markAllAsReadApi } from '@/api/notification'

export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  let pollingTimer: ReturnType<typeof setInterval> | null = null

  /** fetchUnreadCount 从后端刷新未读数 */
  async function fetchUnreadCount() {
    try {
      unreadCount.value = await countUnreadApi()
    } catch {
      // interceptor handles it
    }
  }

  /** markAsRead 标记一条已读并更新计数 */
  async function markAsRead(id: number) {
    await markAsReadApi(id)
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }

  /** markAllAsRead 全部标记已读 */
  async function markAllAsRead() {
    await markAllAsReadApi()
    unreadCount.value = 0
  }

  /** startPolling 启动轮询（每 60 秒） */
  function startPolling() {
    stopPolling()
    fetchUnreadCount()
    pollingTimer = setInterval(fetchUnreadCount, 60_000)
  }

  /** stopPolling 停止轮询 */
  function stopPolling() {
    if (pollingTimer !== null) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  return { unreadCount, fetchUnreadCount, markAsRead, markAllAsRead, startPolling, stopPolling }
})
