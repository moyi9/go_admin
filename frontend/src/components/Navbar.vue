<template>
  <div class="navbar">
    <div class="navbar-left">
      <div class="fold-btn" @click="$emit('toggle-collapse')">
        <el-icon :size="18"><Fold /></el-icon>
      </div>
      <el-breadcrumb separator="/" class="navbar-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item v-for="item in breadcrumbItems" :key="item.path" :to="item.to">
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="navbar-right">
      <NotificationBell />
      <button class="theme-btn" @click="themeStore.toggle()" :title="themeStore.isDark ? '切换亮色模式' : '切换暗黑模式'">
        <el-icon :size="18">
          <Moon v-if="themeStore.isDark" />
          <Sunny v-else />
        </el-icon>
      </button>
      <el-dropdown trigger="click" class="user-dropdown">
        <span class="user-trigger">
          <span class="user-avatar">{{ (authStore.user?.username || '用')[0] }}</span>
          <span class="user-name">{{ authStore.user?.username || '用户' }}</span>
          <el-icon class="user-arrow"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="$router.push('/profile')">
              <el-icon><User /></el-icon>个人中心
            </el-dropdown-item>
            <el-dropdown-item divided @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Fold, ArrowDown, SwitchButton, Moon, Sunny, User } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import NotificationBell from '@/components/NotificationBell.vue'
import { resetRouter } from '@/router'

defineEmits<{ 'toggle-collapse': [] }>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const breadcrumbItems = computed(() => {
  return route.matched
    .filter((r) => r.meta?.title && !r.meta?.hidden && r.name !== 'Layout')
    .map((r) => ({
      path: r.path,
      title: r.meta.title as string,
      to: r.path === '' || r.path === '/' ? undefined : { path: r.path },
    }))
})

function handleLogout() {
  authStore.logout()
  resetRouter()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 100%;
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.fold-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  color: var(--el-text-color-secondary);
  transition: all 0.2s;
}

.fold-btn:hover {
  background: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
}

/* 面包屑 */
.navbar-breadcrumb :deep(.el-breadcrumb__inner) {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.navbar-breadcrumb :deep(.el-breadcrumb__inner.is-link) {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.navbar-breadcrumb :deep(.el-breadcrumb__inner.is-link:hover) {
  color: var(--el-color-primary);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.theme-btn {
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

.theme-btn:hover {
  background: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
}

/* 用户下拉 */
.user-dropdown {
  cursor: pointer;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px 4px 4px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-trigger:hover {
  background: var(--el-fill-color-light);
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: linear-gradient(135deg, #0d7377, #14a0a5);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-name {
  font-size: 14px;
  color: var(--el-text-color-primary);
  font-weight: 500;
}

.user-arrow {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
