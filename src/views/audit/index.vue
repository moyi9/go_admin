<template>
  <div class="audit-page">
    <!-- 搜索卡片：时间范围 + 操作类型 + 资源类型 + 关键词 -->
    <el-card shadow="never" class="search-card">
      <el-form ref="searchFormRef" :model="searchForm" label-width="auto" size="default">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="操作类型" prop="action">
              <el-select v-model="searchForm.action" placeholder="全部操作" clearable style="width: 100%">
                <el-option label="登录" value="LOGIN" />
                <el-option label="创建" value="CREATE" />
                <el-option label="更新" value="UPDATE" />
                <el-option label="删除" value="DELETE" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="资源类型" prop="resource">
              <el-select v-model="searchForm.resource" placeholder="全部资源" clearable style="width: 100%">
                <el-option label="用户" value="user" />
                <el-option label="角色" value="role" />
                <el-option label="权限" value="permission" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="关键词" prop="keyword">
              <el-input v-model="searchForm.keyword" placeholder="用户名 / 操作描述" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="6" style="display: flex; align-items: flex-start; justify-content: flex-end; gap: 8px">
            <el-form-item label=" " label-width="auto">
              <div style="display: flex; gap: 8px">
                <el-button @click="handleReset">
                  <el-icon><Refresh /></el-icon>
                  重置
                </el-button>
                <el-button type="primary" @click="handleSearch">
                  <el-icon><Search /></el-icon>
                  搜索
                </el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="起始时间" prop="since">
              <el-date-picker
                v-model="searchForm.since"
                type="datetime"
                placeholder="选择起始时间"
                format="YYYY-MM-DD HH:mm"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="结束时间" prop="until">
              <el-date-picker
                v-model="searchForm.until"
                type="datetime"
                placeholder="选择结束时间"
                format="YYYY-MM-DD HH:mm"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 表格卡片 -->
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <span class="table-title">操作日志</span>
        <div class="table-actions">
          <el-button @click="handleRefresh" :loading="refreshing">
            <el-icon :class="{ 'is-spinning': refreshing }"><RefreshRight /></el-icon>
            刷新
          </el-button>
        </div>
      </div>

      <el-table
        :data="auditLogs"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="created_at" label="操作时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="username" label="操作人" width="100" />
        <el-table-column prop="action" label="操作类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag
              :type="actionTagType(row.action)"
              size="small"
              effect="plain"
              style="border-radius: 12px; padding: 0 12px; border-width: 1px;"
            >{{ actionLabel(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源类型" width="90" align="center">
          <template #default="{ row }">{{ resourceLabel(row.resource) }}</template>
        </el-table-column>
        <el-table-column prop="resource_id" label="资源ID" width="80" align="center" />
        <el-table-column prop="detail" label="操作描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP 地址" width="140" />
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
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'AuditLogs' })
import { ref, reactive, onMounted } from 'vue'
import { RefreshRight, Search, Refresh } from '@element-plus/icons-vue'
import { listAuditLogsApi } from '@/api/audit'
import type { AuditLog } from '@/types/model'
import type { FormInstance } from 'element-plus'
import { formatTime } from '@/utils/format'

const auditLogs = ref<AuditLog[]>([])
const loading = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const searchFormRef = ref<FormInstance>()
void searchFormRef
const searchForm = reactive({
  action: '',
  resource: '',
  keyword: '',
  since: '',
  until: '',
})

function actionLabel(action: string) {
  const map: Record<string, string> = { LOGIN: '登录', CREATE: '创建', UPDATE: '更新', DELETE: '删除' }
  return map[action] || action
}

function actionTagType(action: string) {
  const map: Record<string, string> = { LOGIN: 'info', CREATE: 'success', UPDATE: 'warning', DELETE: 'danger' }
  return map[action] || 'info'
}

function resourceLabel(resource: string) {
  const map: Record<string, string> = { user: '用户', role: '角色', permission: '权限' }
  return map[resource] || resource
}

async function loadData() {
  loading.value = true
  try {
    const res = await listAuditLogsApi(page.value, pageSize.value, {
      action: searchForm.action || undefined,
      resource: searchForm.resource || undefined,
      keyword: searchForm.keyword || undefined,
      since: searchForm.since || undefined,
      until: searchForm.until || undefined,
    })
    auditLogs.value = res.data
    total.value = res.total
  } catch {
    /* handled by interceptor */
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handleReset() {
  searchForm.action = ''
  searchForm.resource = ''
  searchForm.keyword = ''
  searchForm.since = ''
  searchForm.until = ''
  page.value = 1
  loadData()
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
.audit-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-card :deep(.el-card__body) {
  padding: 20px 24px 8px;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
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
</style>
