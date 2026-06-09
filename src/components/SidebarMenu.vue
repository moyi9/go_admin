<template>
  <div class="sidebar-container">
    <div class="sidebar-logo">
      <div class="logo-icon">⬡</div>
      <h2 v-show="!isCollapse" class="logo-text">后台管理控制</h2>
    </div>
    <el-menu
      :default-active="route.path"
      router
      :collapse="isCollapse"
      :collapse-transition="false"
      class="sidebar-menu"
    >
      <template v-for="item in visibleRoutes" :key="item.path">
        <el-sub-menu v-if="item.children?.length" :index="'/' + item.path" class="sub-menu-custom">
          <template #title>
            <el-icon class="menu-icon"><component :is="item.meta?.icon as string" /></el-icon>
            <span class="menu-title">{{ item.meta?.title }}</span>
          </template>
          <el-menu-item
            v-for="child in item.children"
            :key="child.path"
            :index="'/' + item.path + '/' + child.path"
            class="menu-item-custom menu-item-custom--child"
          >
            <el-icon class="menu-icon"><component :is="child.meta?.icon as string" /></el-icon>
            <template #title>
              <span class="menu-title">{{ child.meta?.title }}</span>
            </template>
          </el-menu-item>
        </el-sub-menu>
        <el-menu-item v-else :index="'/' + item.path" class="menu-item-custom">
          <el-icon class="menu-icon"><component :is="item.meta?.icon as string" /></el-icon>
          <template #title>
            <span class="menu-title">{{ item.meta?.title }}</span>
          </template>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { usePermissionStore } from '@/stores/permission'

defineProps<{ isCollapse?: boolean }>()

const route = useRoute()
const store = usePermissionStore()
const visibleRoutes = computed(() => store.routes.filter((r) => !r.meta?.hidden))
</script>

<style scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--aside-bg, #1a1d23);
}

.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logo-icon {
  font-size: 22px;
  color: #14a0a5;
  line-height: 1;
  flex-shrink: 0;
}

.logo-text {
  color: rgba(255, 255, 255, 0.9);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 1px;
  white-space: nowrap;
  font-family: "PingFang SC", "Noto Sans SC", "Microsoft YaHei", sans-serif;
}

/* ===== Menu 样式覆盖 ===== */
.sidebar-menu {
  flex: 1;
  overflow-y: auto;
  border-right: none !important;
  background: transparent !important;
  padding: 8px 0;
}

.sidebar-menu:not(.el-menu--collapse) {
  padding: 8px 8px;
}

/* 每个菜单项 */
.menu-item-custom {
  height: 44px;
  line-height: 44px;
  margin: 2px 0;
  border-radius: 8px;
  padding: 0 12px !important;
  color: rgba(255, 255, 255, 0.55) !important;
  background: transparent !important;
  transition: all 0.2s ease;
  position: relative;
}

.menu-item-custom:hover {
  color: rgba(255, 255, 255, 0.85) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}

/* 活跃状态 */
.menu-item-custom.is-active {
  color: #fff !important;
  background: linear-gradient(135deg, rgba(13, 115, 119, 0.5) 0%, rgba(20, 160, 165, 0.25) 100%) !important;
  box-shadow: 0 2px 8px rgba(13, 115, 119, 0.2);
}

.menu-item-custom.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 18px;
  background: #14a0a5;
  border-radius: 0 3px 3px 0;
}

/* 图标 */
.menu-icon {
  font-size: 18px;
  margin-right: 4px;
  color: inherit;
}

/* 标题 */
.menu-title {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
}

/* 折叠时效果 */
.sidebar-menu.el-menu--collapse .menu-item-custom {
  border-radius: 0;
  margin: 2px 0;
  padding: 0 !important;
  justify-content: center;
}

.sidebar-menu.el-menu--collapse .menu-item-custom.is-active {
  background: rgba(13, 115, 119, 0.4) !important;
  box-shadow: none;
}

.sidebar-menu.el-menu--collapse .menu-item-custom.is-active::before {
  display: none;
}

/* ===== 子菜单样式 ===== */
.sub-menu-custom {
  background: transparent !important;
}
.sub-menu-custom :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  margin: 2px 0;
  border-radius: 8px;
  padding: 0 12px !important;
  color: rgba(255, 255, 255, 0.55) !important;
  background: transparent !important;
  transition: all 0.2s ease;
}
.sub-menu-custom :deep(.el-sub-menu__title:hover) {
  color: rgba(255, 255, 255, 0.85) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}
.sub-menu-custom.is-active :deep(.el-sub-menu__title) {
  color: #fff !important;
  background: linear-gradient(135deg, rgba(13, 115, 119, 0.5) 0%, rgba(20, 160, 165, 0.25) 100%) !important;
  box-shadow: 0 2px 8px rgba(13, 115, 119, 0.2);
}
.sub-menu-custom :deep(.el-menu) {
  background: transparent !important;
}
.menu-item-custom--child {
  padding-left: 44px !important;
  height: 38px !important;
  line-height: 38px !important;
  font-size: 13px !important;
}
.menu-item-custom--child .menu-icon {
  font-size: 14px !important;
}

/* 滚动条 */
.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}
.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}
.sidebar-menu::-webkit-scrollbar-track {
  background: transparent;
}
</style>
