<template>
  <div class="notification-page">
    <!-- 表格卡片 -->
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <span class="table-title">通知中心</span>
        <div class="table-actions">
          <el-button
            v-if="canSend"
            type="primary"
            @click="sendDialogVisible = true"
          >
            <el-icon><Plus /></el-icon>
            发送通知
          </el-button>
          <el-button
            v-if="notifStore.unreadCount > 0"
            @click="handleMarkAllRead"
          >
            <el-icon><Select /></el-icon>
            全部已读
          </el-button>
          <el-button @click="handleRefresh" :loading="refreshing">
            <el-icon :class="{ 'is-spinning': refreshing }"><RefreshRight /></el-icon>
            刷新
          </el-button>
        </div>
      </div>

      <el-table
        :data="notifications"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
        @row-click="handleRowClick"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tooltip :content="typeLabel(row.type)" placement="top">
              <el-icon :size="18" :color="typeIconColor(row.type)">
                <Bell v-if="row.type === 'system'" />
                <WarningFilled v-else />
              </el-icon>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'notif-unread': !row.is_read }">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="240" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'notif-unread': !row.is_read }">{{ row.content }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag
              :type="row.is_read ? 'info' : 'primary'"
              size="small"
              effect="plain"
              style="border-radius: 12px;"
            >
              {{ row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <span class="pagination-total">共 {{ total }} 条</span>
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="sizes, prev, pager, next"
          background
          @current-change="loadData"
          @size-change="loadData"
        />
      </div>
    </el-card>

    <SendNotificationDialog
      v-model:visible="sendDialogVisible"
      @sent="loadData"
    />
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'Notifications' })
import { ref, onMounted, computed } from 'vue'
import { Bell, WarningFilled, RefreshRight, Plus, Select } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { listNotificationsApi, markAsReadApi, markAllAsReadApi } from '@/api/notification'
import { useNotificationStore } from '@/stores/notification'
import { usePermissionStore } from '@/stores/permission'
import type { Notification } from '@/types/model'
import { formatTime } from '@/utils/format'
import SendNotificationDialog from './components/SendNotificationDialog.vue'

const notifStore = useNotificationStore()
const permissionStore = usePermissionStore()

const notifications = ref<Notification[]>([])
const loading = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const sendDialogVisible = ref(false)

const canSend = computed(() => permissionStore.hasPermission('post.api.v1.notifications'))

function typeLabel(type: string) {
  return type === 'system' ? '系统通知' : '安全提醒'
}

function typeIconColor(type: string) {
  return type === 'security' ? '#e6a23c' : '#0d7377'
}

async function loadData() {
  loading.value = true
  try {
    const res = await listNotificationsApi(page.value, pageSize.value)
    notifications.value = res.data
    total.value = res.total
  } catch {
    /* handled by interceptor */
  } finally {
    loading.value = false
  }
}

async function handleRowClick(row: Notification) {
  if (row.is_read) return
  try {
    await markAsReadApi(row.id)
    row.is_read = true
    notifStore.unreadCount = Math.max(0, notifStore.unreadCount - 1)
  } catch {
    /* handled by interceptor */
  }
}

async function handleMarkAllRead() {
  try {
    await markAllAsReadApi()
    notifications.value.forEach((n) => (n.is_read = true))
    notifStore.unreadCount = 0
    ElMessage.success('全部已读')
  } catch {
    /* handled by interceptor */
  }
}

async function handleRefresh() {
  refreshing.value = true
  await loadData()
  refreshing.value = false
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.notification-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
  flex-wrap: wrap;
  gap: 8px;
}

.table-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.table-actions {
  display: flex;
  gap: 8px;
}

.notif-unread {
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.pagination-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-top: 1px solid var(--el-border-color-light);
}

.pagination-total {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.is-spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
