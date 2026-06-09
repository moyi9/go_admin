<template>
  <el-popover
    placement="bottom-end"
    :width="360"
    trigger="click"
    :visible="popoverVisible"
    @show="handleShow"
    @hide="popoverVisible = false"
  >
    <template #reference>
      <button class="bell-btn" @click="popoverVisible = !popoverVisible" title="通知">
        <el-badge :value="notifStore.unreadCount" :hidden="notifStore.unreadCount === 0" :max="99">
          <el-icon :size="18"><Bell /></el-icon>
        </el-badge>
      </button>
    </template>

    <div class="notif-popover">
      <div class="notif-popover-header">
        <span class="notif-popover-title">通知</span>
        <span class="notif-popover-subtitle" v-if="notifStore.unreadCount > 0">
          {{ notifStore.unreadCount }} 条未读
        </span>
      </div>

      <div v-if="loading" class="notif-loading">加载中...</div>
      <div v-else-if="unreadList.length === 0" class="notif-empty">
        <el-icon :size="32" color="#c0c4cc"><Bell /></el-icon>
        <p>暂无通知</p>
      </div>
      <div v-else class="notif-list">
        <div
          v-for="item in unreadList"
          :key="item.id"
          class="notif-item"
          @click="handleMarkRead(item)"
        >
          <div class="notif-item-icon">
            <el-icon :size="16" :color="iconColor(item.type)">
              <Bell v-if="item.type === 'system'" />
              <WarningFilled v-else />
            </el-icon>
          </div>
          <div class="notif-item-body">
            <div class="notif-item-title">{{ item.title }}</div>
            <div class="notif-item-time">{{ formatTime(item.created_at) }}</div>
          </div>
        </div>
      </div>

      <div class="notif-popover-footer">
        <el-button
          v-if="notifStore.unreadCount > 0"
          text
          size="small"
          @click="handleMarkAllRead"
        >
          全部已读
        </el-button>
        <el-button text size="small" type="primary" @click="goToCenter">
          查看全部 →
        </el-button>
      </div>
    </div>
  </el-popover>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Bell, WarningFilled } from '@element-plus/icons-vue'
import { useNotificationStore } from '@/stores/notification'
import { listUnreadNotificationsApi } from '@/api/notification'
import type { Notification } from '@/types/model'
import { formatTime } from '@/utils/format'

const router = useRouter()
const notifStore = useNotificationStore()

const popoverVisible = ref(false)
const loading = ref(false)
const unreadList = ref<Notification[]>([])

function iconColor(type: string) {
  return type === 'security' ? '#e6a23c' : '#0d7377'
}

async function handleShow() {
  loading.value = true
  try {
    unreadList.value = await listUnreadNotificationsApi()
  } catch {
    unreadList.value = []
  } finally {
    loading.value = false
  }
}

async function handleMarkRead(item: Notification) {
  try {
    await notifStore.markAsRead(item.id)
    unreadList.value = unreadList.value.filter((n) => n.id !== item.id)
    ElMessage.success('已标记已读')
  } catch {
    // interceptor handles it
  }
}

async function handleMarkAllRead() {
  try {
    await notifStore.markAllAsRead()
    unreadList.value = []
    ElMessage.success('全部已读')
  } catch {
    // interceptor handles it
  }
}

function goToCenter() {
  popoverVisible.value = false
  router.push('/notifications')
}

onMounted(() => {
  notifStore.startPolling()
})

onUnmounted(() => {
  notifStore.stopPolling()
})
</script>

<style scoped>
.bell-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--el-text-color-secondary);
  transition: all 0.2s;
}

.bell-btn:hover {
  background: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
}

.notif-popover {
  display: flex;
  flex-direction: column;
}

.notif-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 12px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.notif-popover-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.notif-popover-subtitle {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.notif-loading,
.notif-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.notif-list {
  max-height: 320px;
  overflow-y: auto;
  margin: 0 -12px;
}

.notif-item {
  display: flex;
  gap: 12px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s;
}

.notif-item:hover {
  background: var(--el-fill-color-light);
}

.notif-item-icon {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color);
  border-radius: 8px;
  margin-top: 2px;
}

.notif-item-body {
  flex: 1;
  min-width: 0;
}

.notif-item-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  line-height: 1.4;
}

.notif-item-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
}

.notif-popover-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0 0;
  border-top: 1px solid var(--el-border-color-light);
  margin-top: 4px;
}
</style>
