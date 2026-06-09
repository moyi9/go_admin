<template>
  <div class="perm-page">
    <!-- 搜索卡片 -->
    <el-card shadow="never" class="search-card">
      <el-form ref="searchFormRef" :model="searchForm" label-width="auto" size="default">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="权限名称" prop="name">
              <el-input v-model="searchForm.name" placeholder="请输入权限名称" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="请求方法" prop="method">
              <el-select v-model="searchForm.method" placeholder="全部方法" clearable style="width: 100%">
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6" style="display: flex; align-items: flex-start; justify-content: flex-end">
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
      </el-form>
    </el-card>

    <el-card shadow="never" class="perm-card">
      <div class="perm-header">
        <span class="perm-title">权限列表</span>
        <div class="table-actions">
          <el-button
            type="danger"
            :disabled="!selectedIds.length"
            @click="handleBatchDelete"
          >
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
          <el-dropdown @command="handleExport">
            <el-button>
              <el-icon><Download /></el-icon>
              导出
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="xlsx">Excel (.xlsx)</el-dropdown-item>
                <el-dropdown-item command="csv">CSV</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button v-permission="'post.api.v1.permissions'" type="primary" @click="openCreate">
            <el-icon><Plus /></el-icon>
            新增权限
          </el-button>
        </div>
      </div>
      <el-table :data="permissions" v-loading="loading" border stripe style="width: 100%" class="perm-table" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="method" label="方法" width="100">
          <template #default="{ row }">
            <el-tag :type="methodTagType(row.method)">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'put.api.v1.permissions.id'" size="small" type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-permission="'delete.api.v1.permissions.id'" size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="perm-pagination">
        <span class="pagination-total">共 {{ total }} 条</span>
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          background
          @current-change="loadData"
        />
      </div>
    </el-card>
    <PermissionFormDialog v-model="dialogVisible" :permission="editingPermission" @saved="onSaved" />
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'Permissions' })
import { ref, reactive, onMounted } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Refresh, Delete, Plus, Download, ArrowDown } from '@element-plus/icons-vue'
import { listPermissionsApi, deletePermissionApi, batchDeletePermissionsApi } from '@/api/permission'
import type { Permission } from '@/types/model'
import type { FormInstance } from 'element-plus'
import PermissionFormDialog from './components/PermissionFormDialog.vue'
import { formatTime } from '@/utils/format'
import { exportCSV } from '@/utils/export'
import { exportXLSX } from '@/utils/exportXLSX'

const permissions = ref<Permission[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)

const searchFormRef = ref<FormInstance>()
void searchFormRef
const searchForm = reactive({
  name: '',
  method: '',
})

const selectedIds = ref<number[]>([])
const dialogVisible = ref(false)
const editingPermission = ref<Permission | null>(null)

onMounted(() => loadData())

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    for (const [key, value] of Object.entries(searchForm)) {
      if (value) params[key] = value
    }
    const res = await listPermissionsApi(page.value, pageSize, Object.keys(params).length ? params : undefined)
    permissions.value = res.data
    total.value = res.total
  } catch { /* handled */ }
  finally { loading.value = false }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handleReset() {
  searchForm.name = ''
  searchForm.method = ''
  page.value = 1
  loadData()
}

function handleSelectionChange(rows: Permission[]) {
  selectedIds.value = rows.map((r) => r.id)
}

async function handleBatchDelete() {
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 个权限？`, '提示')
    await batchDeletePermissionsApi(selectedIds.value)
    ElMessage.success('批量删除成功')
    selectedIds.value = []
    loadData()
  } catch { /* cancelled or failed */ }
}

async function handleExport(format: string) {
  const columns = [
    { key: 'id', label: '编号' },
    { key: 'name', label: '名称' },
    { key: 'method', label: '方法' },
    { key: 'path', label: '路径' },
    { key: 'description', label: '描述' },
    { key: 'created_at', label: '创建时间' },
  ]
  try {
    const res = await listPermissionsApi(1, 10000)
    if (format === 'xlsx') {
      await exportXLSX(res.data, columns, '权限列表')
    } else {
      exportCSV(res.data, columns, '权限列表')
    }
  } catch { /* handled */ }
}

/** methodTagType 根据 HTTP 方法返回 Element Plus Tag 类型：GET→success, POST→warning, PUT→primary, DELETE→danger */
function methodTagType(method: string) {
  const map: Record<string, string> = { GET: 'success', POST: 'warning', PUT: 'primary', DELETE: 'danger' }
  return map[method] || 'info'
}

function openCreate() {
  editingPermission.value = null
  dialogVisible.value = true
}

function openEdit(perm: Permission) {
  editingPermission.value = perm
  dialogVisible.value = true
}

async function handleDelete(perm: Permission) {
  try {
    await ElMessageBox.confirm(`确定删除权限 "${perm.name}"？`, '提示')
    await deletePermissionApi(perm.id)
    ElMessage.success('删除成功')
    loadData()
  } catch { /* cancelled or failed */ }
}

function onSaved() {
  dialogVisible.value = false
  loadData()
}
</script>

<style scoped>
.perm-page {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  gap: 16px;
}

/* ===== 搜索卡片 ===== */
.search-card {
  flex-shrink: 0;
  border: 1px solid var(--el-border-color-light);
  border-radius: 10px;
  border-left: 3px solid #0d7377;
  transition: box-shadow 0.25s ease;
}

.search-card:hover {
  box-shadow: 0 2px 12px rgba(13, 115, 119, 0.06);
}

.search-card :deep(.el-card__body) {
  padding: 20px 24px 4px 24px;
}

.perm-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.perm-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
}

.perm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
  position: relative;
  background: linear-gradient(90deg, rgba(13, 115, 119, 0.03) 0%, transparent 40%);
}

.perm-header::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: #0d7377;
  border-radius: 0 2px 2px 0;
}

.perm-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  padding-left: 4px;
}

.table-actions {
  display: flex;
  gap: 8px;
}

/* 表格 */
.perm-table {
  flex: 1;
  border: none;
  border-radius: 0;
}

.perm-table :deep(.el-table) {
  border: none;
}

.perm-table :deep(.el-table th.el-table__cell) {
  background-color: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 13px;
  letter-spacing: 0.3px;
  border-bottom: 2px solid var(--el-border-color-light);
}

.perm-table :deep(.el-table__body tr:hover > td) {
  background-color: var(--el-color-primary-light-9);
}

/* 分页 */
.perm-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  border-top: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
  background: var(--el-bg-color);
}

.pagination-total {
  font-size: 14px;
  color: #909399;
}
</style>
