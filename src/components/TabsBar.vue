<template>
  <div class="tabs-bar">
    <div class="tabs-container">
      <div
        v-for="tab in tabsStore.tabs"
        :key="tab.path"
        class="tab-item"
        :class="{ 'is-active': tabsStore.activeTabPath === tab.path }"
        @click="switchTab(tab)"
        @contextmenu.prevent="openContextMenu($event, tab)"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <el-icon v-if="tab.closable" class="tab-close" @click.stop="closeTab(tab)">
          <Close />
        </el-icon>
      </div>
    </div>

    <teleport to="body">
      <div
        v-show="menuVisible"
        class="tab-context-menu"
        :style="{ left: menuX + 'px', top: menuY + 'px' }"
        @click.stop
      >
        <div class="menu-item" @click="closeCurrentTab">关闭当前</div>
        <div class="menu-item" @click="closeOtherTabs">关闭其他</div>
        <div class="menu-item" @click="closeAllTabs">关闭全部</div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Close } from '@element-plus/icons-vue'
import { useTabsStore, type TabItem } from '@/stores/tabs'

const router = useRouter()
const tabsStore = useTabsStore()

const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTab = ref<TabItem | null>(null)

function switchTab(tab: TabItem) {
  router.push(tab.path)
}

function closeTab(tab: TabItem) {
  const nextPath = tabsStore.removeTab(tab.path)
  if (nextPath) router.push(nextPath)
}

function openContextMenu(e: MouseEvent, tab: TabItem) {
  contextTab.value = tab
  menuX.value = e.clientX
  menuY.value = e.clientY
  menuVisible.value = true
}

function closeMenu() {
  menuVisible.value = false
  contextTab.value = null
}

function closeCurrentTab() {
  if (!contextTab.value) return
  if (contextTab.value.closable) closeTab(contextTab.value)
  closeMenu()
}

function closeOtherTabs() {
  if (!contextTab.value) return
  tabsStore.closeOthers(contextTab.value.path)
  if (tabsStore.activeTabPath !== contextTab.value.path) {
    router.push(contextTab.value.path)
  }
  closeMenu()
}

function closeAllTabs() {
  const nextPath = tabsStore.closeAll()
  closeMenu()
  if (nextPath) router.push(nextPath)
}

function handleDocumentClick() {
  closeMenu()
}

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
})

onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick)
})
</script>

<style scoped>
.tabs-bar {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  padding: 0 8px;
  user-select: none;
  flex-shrink: 0;
  position: relative;
}

.tabs-container {
  display: flex;
  align-items: center;
  gap: 2px;
  overflow-x: auto;
  scrollbar-width: none;
  padding: 6px 0 0;
}

.tabs-container::-webkit-scrollbar {
  display: none;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 13px;
  border: 1px solid transparent;
  border-bottom: none;
  border-radius: 6px 6px 0 0;
  background: transparent;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
  position: relative;
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.tab-item:hover {
  color: var(--el-text-color-primary);
  background: var(--el-fill-color-light);
}

.tab-item.is-active {
  color: var(--el-text-color-primary);
  font-weight: 500;
  background: var(--el-bg-color-page);
  border-color: var(--el-border-color-light);
  border-bottom-color: var(--el-bg-color-page);
}

.tab-item.is-active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--el-bg-color-page);
}

/* 活跃标签左侧装饰条 */
.tab-item.is-active::before {
  content: '';
  position: absolute;
  left: 10px;
  right: 10px;
  top: -1px;
  height: 2px;
  background: #0d7377;
  border-radius: 0 0 2px 2px;
  opacity: 0.8;
}

.tab-title {
  font-size: 13px;
  line-height: 1.4;
}

.tab-close {
  font-size: 11px;
  border-radius: 50%;
  padding: 2px;
  transition: all 0.2s;
  opacity: 0.5;
}

.tab-item:hover .tab-close {
  opacity: 0.8;
}

.tab-close:hover {
  background: var(--el-border-color);
  color: var(--el-text-color-primary);
  opacity: 1;
}

/* 右键菜单 */
.tab-context-menu {
  position: fixed;
  z-index: 9999;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  box-shadow: var(--el-box-shadow-light);
  padding: 4px 0;
  min-width: 120px;
  overflow: hidden;
}

.menu-item {
  padding: 8px 16px;
  font-size: 13px;
  cursor: pointer;
  color: var(--el-text-color-primary);
  transition: background 0.15s;
}

.menu-item:hover {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
</style>
