<template>
  <el-container class="app-container">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="app-aside">
      <SidebarMenu :is-collapse="isCollapse" />
    </el-aside>
    <el-container class="app-right">
      <el-header class="app-header">
        <Navbar @toggle-collapse="isCollapse = !isCollapse" />
      </el-header>
      <TabsBar />
      <div class="loading-bar" :class="{ active: loadingStore.count > 0 }" />
      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <keep-alive :include="keepAliveNames">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<!-- 主布局组件：左-侧边栏 + 右-顶栏/标签栏/主内容区，主内容支持 keep-alive 缓存 -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import SidebarMenu from '@/components/SidebarMenu.vue'
import Navbar from '@/components/Navbar.vue'
import TabsBar from '@/components/TabsBar.vue'
import { useTabsStore } from '@/stores/tabs'
import { useLoadingStore } from '@/stores/loading'

const isCollapse = ref(false) // 侧边栏折叠状态

const route = useRoute()
const tabsStore = useTabsStore()
const loadingStore = useLoadingStore()

/** keepAliveNames 根据当前打开的标签页列表决定哪些组件需要缓存 */
const keepAliveNames = computed(() => {
  const names = new Set(tabsStore.tabs.map((t) => t.name))
  if (route.name) names.add(route.name as string)
  return Array.from(names)
})
</script>

<style scoped>
.app-container {
  height: 100vh;
  background: var(--el-bg-color-page);
}

.app-aside {
  background-color: var(--aside-bg, #1a1d23);
  overflow: hidden;
  transition: width 0.28s ease;
  border-right: 1px solid var(--aside-border, rgba(255, 255, 255, 0.06));
}

.app-right {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-header {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  flex-shrink: 0;
}

.app-main {
  padding: 16px 24px;
  overflow-y: auto;
  flex: 1;
  background: var(--el-bg-color-page);
}

.loading-bar {
  height: 2px;
  background: linear-gradient(90deg, #0d7377, #14a0a5, #0d7377);
  background-size: 200% 100%;
  transition: opacity 0.2s;
  opacity: 0;
  flex-shrink: 0;
}

.loading-bar.active {
  opacity: 1;
  animation: loading-slide 1.2s ease-in-out infinite;
}

@keyframes loading-slide {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 路由切换淡入淡出 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
