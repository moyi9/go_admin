<!-- 根组件：仅渲染路由视图，全局暗黑模式样式覆盖 -->
<template>
  <router-view />
</template>

<style>
/* ========== 暗黑模式变量覆盖 ========== */
html.dark {
  --aside-bg: #121212;
  --aside-border: rgba(255, 255, 255, 0.06);
}

/* 侧边栏背景 */
html.dark .app-aside {
  background-color: var(--aside-bg, #121212) !important;
}

/* 侧边栏容器 */
html.dark .sidebar-container {
  background: var(--aside-bg, #121212) !important;
}

/* 暗黑模式 logo 文字 */
html.dark .sidebar-logo .logo-text {
  color: rgba(255, 255, 255, 0.85);
}

/* 暗黑模式头部 */
html.dark .app-header {
  background: var(--el-bg-color) !important;
  border-bottom-color: var(--el-border-color-light) !important;
}

/* 暗黑模式标签栏 */
html.dark .tabs-bar {
  background: var(--el-bg-color);
  border-bottom-color: var(--el-border-color-light);
}

/* 暗黑模式标签项 */
html.dark .tab-item.is-active {
  border-color: var(--el-border-color-light);
}

/* 历史兼容：原 class 名保持 */
html.dark .sidebar-logo h2 {
  color: rgba(255, 255, 255, 0.85);
}

html.dark .el-aside {
  background-color: var(--aside-bg, #121212) !important;
}

html.dark .el-menu {
  background-color: transparent !important;
  border-right: none !important;
}

html.dark .el-header {
  background: var(--el-bg-color) !important;
  border-bottom-color: var(--el-border-color-light) !important;
}

html.dark .tabs-bar {
  background: var(--el-bg-color);
  border-bottom-color: var(--el-border-color-light);
}

html.dark .tab-item {
  background: transparent;
  border-color: transparent;
  color: var(--el-text-color-secondary);
}

html.dark .tab-item.is-active {
  background: var(--el-bg-color-page);
  border-color: var(--el-border-color-light);
  color: var(--el-text-color-primary);
}

html.dark .tab-item.is-active::after {
  background: var(--el-bg-color-page);
}

html.dark .tab-context-menu {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color-light);
}

html.dark .menu-item {
  color: var(--el-text-color-primary);
}

html.dark .fold-btn {
  color: var(--el-text-color-primary);
}
</style>
