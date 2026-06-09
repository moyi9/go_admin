<template>
  <div class="role-page">
    <!-- 搜索卡片：3 字段多条件搜索（角色名称/编码/状态） -->
    <el-card shadow="never" class="search-card">
      <el-form ref="searchFormRef" :model="searchForm" label-width="auto" size="default">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="searchForm.name" placeholder="请输入角色名称" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="角色编码" prop="code">
              <el-input v-model="searchForm.code" placeholder="请输入角色编码" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="角色状态" prop="status">
              <el-select v-model="searchForm.status" placeholder="请选择角色状态" clearable style="width: 100%">
                <el-option label="启用" value="active" />
                <el-option label="禁用" value="disabled" />
              </el-select>
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
      </el-form>
    </el-card>

    <!-- 底部表格卡片 -->
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <span class="table-title">角色列表</span>
        <div class="table-actions">
          <el-button v-permission="'post.api.v1.roles'" type="primary" @click="openCreate">
            <el-icon><Plus /></el-icon>
            新增
          </el-button>
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
          <el-button @click="handleRefresh" :loading="refreshing">
            <el-icon :class="{ 'is-spinning': refreshing }"><RefreshRight /></el-icon>
            刷新
          </el-button>
          <el-popover trigger="click" placement="bottom-end" :width="160">
            <template #reference>
              <el-button>
                <el-icon><Setting /></el-icon>
                列设置
              </el-button>
            </template>
            <div style="padding: 4px 0">
              <el-checkbox
                v-for="col in columnDefs"
                :key="col.key"
                v-model="col.visible"
                :label="col.label"
                style="display: flex; margin-bottom: 8px"
                @change="forceUpdate = Date.now()"
              />
            </div>
          </el-popover>
        </div>
      </div>

      <el-table
        ref="tableRef"
        :data="roles"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column v-if="getCol('selection')?.visible" type="selection" width="50" />
        <el-table-column v-if="getCol('index')?.visible" type="index" label="序号" width="60" />
        <el-table-column v-if="getCol('name')?.visible" prop="name" label="角色名称" min-width="140" />
        <el-table-column v-if="getCol('code')?.visible" prop="code" label="角色编码" min-width="140" />
        <el-table-column v-if="getCol('description')?.visible" prop="description" label="角色描述" min-width="180" />
        <el-table-column v-if="getCol('status')?.visible" prop="status" label="角色状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag
              v-if="row.status === 'active'"
              size="small"
              effect="plain"
              style="border: 1px solid #67c23a; color: #67c23a; background: transparent; border-radius: 12px; padding: 0 12px"
            >
              启用
            </el-tag>
            <el-tag
              v-else
              size="small"
              effect="plain"
              style="border: 1px solid #e6a23c; color: #e6a23c; background: transparent; border-radius: 12px; padding: 0 12px"
            >
              禁用
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="getCol('permissions_count')?.visible" label="权限数" width="80" align="center">
          <template #default="{ row }">{{ row.permissions?.length || 0 }}</template>
        </el-table-column>
        <el-table-column v-if="getCol('created_at')?.visible" label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column v-if="getCol('actions')?.visible" label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'put.api.v1.roles.id'" size="small" type="primary" link @click="openEdit(row)">编辑</el-button>
            <el-button v-permission="'post.api.v1.roles.id.permissions'" size="small" type="primary" link @click="openAssignPermission(row)">分配权限</el-button>
            <el-button v-permission="'delete.api.v1.roles.id'" size="small" type="danger" link @click="handleDelete(row)">删除</el-button>
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

    <RoleFormDialog v-model="dialogVisible" :role="editingRole" @saved="onSaved" />
    <PermissionAssignDialog v-model="permDialogVisible" :role="assigningRole" @saved="loadData" />
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'Roles' })
import { ref, reactive, onMounted } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus, Delete, RefreshRight, Setting, Search, Refresh, Download, ArrowDown } from '@element-plus/icons-vue'
import { listRolesApi, deleteRoleApi, batchDeleteRolesApi } from '@/api/role'
import type { Role } from '@/types/model'
import type { FormInstance } from 'element-plus'
import RoleFormDialog from './components/RoleFormDialog.vue'
import PermissionAssignDialog from './components/PermissionAssignDialog.vue'
import { formatTime } from '@/utils/format'
import { exportCSV } from '@/utils/export'
import { exportXLSX } from '@/utils/exportXLSX'

const roles = ref<Role[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedIds = ref<number[]>([])

const searchFormRef = ref<FormInstance>()
void searchFormRef
const searchForm = reactive({
  name: '',
  code: '',
  status: '',
})

const refreshing = ref(false)
const forceUpdate = ref(0)

interface ColumnDef {
  key: string
  label: string
  visible: boolean
}

const columnDefs = ref<ColumnDef[]>([
  { key: 'selection', label: '多选', visible: true },
  { key: 'index', label: '序号', visible: true },
  { key: 'name', label: '角色名称', visible: true },
  { key: 'code', label: '角色编码', visible: true },
  { key: 'description', label: '角色描述', visible: true },
  { key: 'status', label: '角色状态', visible: true },
  { key: 'permissions_count', label: '权限数', visible: true },
  { key: 'created_at', label: '创建时间', visible: true },
  { key: 'actions', label: '操作', visible: true },
])

function getCol(key: string) {
  return columnDefs.value.find((c) => c.key === key)
}

const dialogVisible = ref(false)
const permDialogVisible = ref(false)
const editingRole = ref<Role | null>(null)
const assigningRole = ref<Role | null>(null)

onMounted(() => loadData())

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    for (const [key, value] of Object.entries(searchForm)) {
      if (value) params[key] = value
    }
    const res = await listRolesApi(page.value, pageSize.value, Object.keys(params).length ? params : undefined)
    roles.value = res.data
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
  searchForm.code = ''
  searchForm.status = ''
  page.value = 1
  loadData()
}

async function handleRefresh() {
  refreshing.value = true
  await loadData()
  refreshing.value = false
}

async function handleExport(format: string) {
  const columns = [
    { key: 'id', label: '编号' },
    { key: 'name', label: '角色名称' },
    { key: 'code', label: '角色编码' },
    { key: 'description', label: '角色描述' },
    { key: 'status', label: '状态' },
    { key: 'created_at', label: '创建时间' },
  ]
  try {
    const res = await listRolesApi(1, 10000)
    if (format === 'xlsx') {
      await exportXLSX(res.data, columns, '角色列表')
    } else {
      exportCSV(res.data, columns, '角色列表')
    }
  } catch { /* handled */ }
}

function handleSelectionChange(rows: Role[]) {
  selectedIds.value = rows.map((r) => r.id)
}

async function handleBatchDelete() {
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 个角色？`, '提示')
    await batchDeleteRolesApi(selectedIds.value)
    ElMessage.success('批量删除成功')
    selectedIds.value = []
    loadData()
  } catch { /* cancelled or failed */ }
}

function openCreate() {
  editingRole.value = null
  dialogVisible.value = true
}

function openEdit(role: Role) {
  editingRole.value = role
  dialogVisible.value = true
}

function openAssignPermission(role: Role) {
  assigningRole.value = role
  permDialogVisible.value = true
}

async function handleDelete(role: Role) {
  try {
    await ElMessageBox.confirm(`确定删除角色 "${role.name}"？`, '提示')
    await deleteRoleApi(role.id)
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
.role-page {
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

/* ===== 表格卡片 ===== */
.table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.table-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
}

/* 表头区域 + 青瓷色装饰条 */
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
  position: relative;
  background: linear-gradient(90deg, rgba(13, 115, 119, 0.03) 0%, transparent 40%);
}

.table-header::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: #0d7377;
  border-radius: 0 2px 2px 0;
}

.table-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  padding-left: 4px;
}

.table-actions {
  display: flex;
  gap: 8px;
}

/* ===== 表格主体 ===== */
.table-card :deep(.el-table) {
  flex: 1;
  border: none;
  border-radius: 0;
}

.table-card :deep(.el-table th.el-table__cell) {
  background-color: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 13px;
  letter-spacing: 0.3px;
  border-bottom: 2px solid var(--el-border-color-light);
}

.table-card :deep(.el-table th.el-table__cell .cell) {
  padding: 8px 12px;
}

/* 行悬停 — 青瓷色微 tint */
.table-card :deep(.el-table__body tr:hover > td) {
  background-color: var(--el-color-primary-light-9);
}

/* ===== 分页 ===== */
.pagination-wrapper {
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

/* ===== 工具类 ===== */
.is-spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
