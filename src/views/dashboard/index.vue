<template>
  <div class="dashboard-page">
    <!-- 欢迎卡片 -->
    <el-card shadow="never" class="welcome-card">
      <div class="welcome-inner">
        <div class="welcome-avatar">
          <span class="avatar-text">{{ (authStore.user?.username || 'A')[0] }}</span>
        </div>
        <div class="welcome-info">
          <h2 class="welcome-title">欢迎回来，{{ authStore.user?.nickname || authStore.user?.username }}</h2>
          <p class="welcome-desc">{{ authStore.user?.email }}</p>
        </div>
      </div>
    </el-card>

    <!-- 统计指标 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card--users">
          <div class="stat-body">
            <div class="stat-value">{{ stats.users }}</div>
            <p class="stat-label">用户数</p>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card--roles">
          <div class="stat-body">
            <div class="stat-value">{{ stats.roles }}</div>
            <p class="stat-label">角色数</p>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card--perms">
          <div class="stat-body">
            <div class="stat-value">{{ stats.permissions }}</div>
            <p class="stat-label">权限数</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表卡片 -->
    <el-row :gutter="20">
      <el-col :xs="24" :md="12" :lg="8">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <span class="chart-title">用户性别分布</span>
          </template>
          <div ref="genderChartRef" style="height: 260px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12" :lg="8">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <span class="chart-title">用户状态分布</span>
          </template>
          <div ref="statusChartRef" style="height: 260px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12" :lg="8">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <span class="chart-title">角色用户分布</span>
          </template>
          <div ref="roleChartRef" style="height: 260px" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'Dashboard' })
import { onMounted, reactive, ref, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { listUsersApi } from '@/api/user'
import { listRolesApi } from '@/api/role'
import { listPermissionsApi } from '@/api/permission'
import { getDashboardStatsApi } from '@/api/dashboard'
import type { DashboardStats } from '@/api/dashboard'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const stats = reactive({ users: 0, roles: 0, permissions: 0 })

const genderChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
const roleChartRef = ref<HTMLElement>()
let genderChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null
let roleChart: echarts.ECharts | null = null

const GENDER_LABEL: Record<string, string> = { male: '男', female: '女', unknown: '未知' }

function renderGenderChart(data: DashboardStats['gender_distribution']) {
  if (!genderChartRef.value) return
  genderChart?.dispose()
  genderChart = echarts.init(genderChartRef.value)

  const entries = Object.entries(data).filter(([_, v]) => v > 0)
  genderChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} 人 ({d}%)' },
    color: ['#0d7377', '#52c0c4', '#e8e8e8'],
    series: [{
      type: 'pie',
      radius: ['40%', '65%'],
      data: entries.map(([k, v]) => ({ name: GENDER_LABEL[k] || k, value: v })),
      itemStyle: { borderRadius: 4, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{d}%' },
      emphasis: { label: { show: true, fontWeight: 'bold' } },
    }],
  })
}

function renderStatusChart(data: DashboardStats['status_distribution']) {
  if (!statusChartRef.value) return
  statusChart?.dispose()
  statusChart = echarts.init(statusChartRef.value)

  const entries = Object.entries(data).filter(([_, v]) => v > 0)
  statusChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} 人 ({d}%)' },
    color: ['#32a1ce', '#e87461'],
    series: [{
      type: 'pie',
      radius: ['40%', '65%'],
      data: entries.map(([k, v]) => ({
        name: k === 'active' ? '活跃' : k === 'disabled' ? '禁用' : k,
        value: v,
      })),
      itemStyle: { borderRadius: 4, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{d}%' },
      emphasis: { label: { show: true, fontWeight: 'bold' } },
    }],
  })
}

function renderRoleChart(data: DashboardStats['role_distribution']) {
  if (!roleChartRef.value) return
  roleChart?.dispose()
  roleChart = echarts.init(roleChartRef.value)

  roleChart.setOption({
    tooltip: { trigger: 'axis', formatter: '{b}: {c} 人' },
    color: ['#0d7377'],
    grid: { left: 80, right: 20, top: 20, bottom: 20 },
    xAxis: { type: 'value', minInterval: 1 },
    yAxis: {
      type: 'category',
      data: data.map(r => r.role_name),
      axisLabel: { fontWeight: 500 },
    },
    series: [{
      type: 'bar',
      data: data.map(r => r.user_count),
      barWidth: 20,
      itemStyle: {
        borderRadius: [0, 4, 4, 0],
        color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
          { offset: 0, color: '#52c0c4' },
          { offset: 1, color: '#0d7377' },
        ]),
      },
    }],
  })
}

function handleResize() {
  genderChart?.resize()
  statusChart?.resize()
  roleChart?.resize()
}

onMounted(async () => {
  try {
    const [u, r, p] = await Promise.all([
      listUsersApi(1, 1),
      listRolesApi(1, 1),
      listPermissionsApi(1, 1),
    ])
    stats.users = u.total
    stats.roles = r.total
    stats.permissions = p.total
  } catch { /* stats remain 0 */ }

  try {
    const dashboardStats = await getDashboardStatsApi()
    renderGenderChart(dashboardStats.gender_distribution)
    renderStatusChart(dashboardStats.status_distribution)
    renderRoleChart(dashboardStats.role_distribution)
  } catch { /* charts render empty */ }

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  genderChart?.dispose()
  statusChart?.dispose()
  roleChart?.dispose()
})
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ===== 欢迎卡片 ===== */
.welcome-card {
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.welcome-card :deep(.el-card__body) {
  padding: 0;
}

.welcome-inner {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px 32px;
  background: linear-gradient(135deg, rgba(13, 115, 119, 0.04) 0%, transparent 60%);
  position: relative;
}

.welcome-inner::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: #0d7377;
  border-radius: 0 2px 2px 0;
}

.welcome-avatar {
  width: 60px;
  height: 60px;
  border-radius: 14px;
  background: linear-gradient(135deg, #0d7377, #14a0a5);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 16px rgba(13, 115, 119, 0.25);
}

.avatar-text {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
}

.welcome-info {
  flex: 1;
}

.welcome-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 4px;
  color: var(--el-text-color-primary);
}

.welcome-desc {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

/* ===== 统计卡片行 ===== */
.stats-row {
  margin-bottom: 0 !important;
}

.stat-card {
  position: relative;
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
}

.stat-body {
  position: relative;
  padding: 24px 32px;
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1.1;
  margin-bottom: 4px;
}

.stat-label {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.stat-card--users::before,
.stat-card--roles::before,
.stat-card--perms::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
}

.stat-card--users::before { background: #0d7377; }
.stat-card--roles::before { background: #14a0a5; }
.stat-card--perms::before { background: #52c0c4; }

/* ===== 图表卡片 ===== */
.chart-card {
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.chart-card :deep(.el-card__header) {
  padding: 14px 20px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.chart-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
</style>
