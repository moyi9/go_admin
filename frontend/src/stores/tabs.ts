/** 多标签页状态管理：标签页的添加、关闭、右键菜单操作 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export interface TabItem {
  path: string
  name: string
  title: string
  closable: boolean // Dashboard 不允许关闭
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<TabItem[]>([])
  const activeTabPath = ref<string>('')

  /** addTab 添加标签页，重复路径激活已有标签 */
  function addTab(route: RouteLocationNormalized) {
    const title = route.meta?.title as string | undefined
    if (!title || route.meta?.hidden) return

    const existing = tabs.value.find((t) => t.path === route.path)
    if (existing) {
      activeTabPath.value = route.path
      return
    }

    tabs.value.push({
      path: route.path,
      name: (route.name as string) || '',
      title,
      closable: route.name !== 'Dashboard',
    })
    activeTabPath.value = route.path
  }

  /** removeTab 关闭标签页，返回下一个激活的路径 */
  function removeTab(path: string): string | null {
    const idx = tabs.value.findIndex((t) => t.path === path)
    if (idx === -1) return null
    tabs.value.splice(idx, 1)

    if (activeTabPath.value === path) {
      const next = tabs.value[Math.min(idx, tabs.value.length - 1)]
      return next?.path || null
    }
    return null
  }

  /** closeOthers 关闭除指定路径外的所有可关闭标签 */
  function closeOthers(path: string) {
    tabs.value = tabs.value.filter((t) => t.path === path || !t.closable)
    activeTabPath.value = path
  }

  /** closeAll 关闭所有可关闭标签，返回 Dashboard 路径 */
  function closeAll() {
    const keep = tabs.value.filter((t) => !t.closable)
    tabs.value = keep
    return keep[0]?.path || null
  }

  return { tabs, activeTabPath, addTab, removeTab, closeOthers, closeAll }
})
