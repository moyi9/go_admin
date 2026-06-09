/** 权限状态管理：用户权限码存储、动态路由生成 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import type { User } from '@/types/model'
import { generateRoutes } from '@/router/dynamic-routes'

export const usePermissionStore = defineStore('permission', () => {
  const routes = ref<RouteRecordRaw[]>([]) // 根据权限动态生成的路由表
  let codes: string[] = []                  // 当前用户的权限码列表

  /** setPermissions 从用户角色中提取所有权限码 */
  function setPermissions(user: User) {
    const s = new Set<string>()
    for (const r of user.roles || []) {
      for (const p of r.permissions || []) s.add(p.code)
    }
    codes = Array.from(s)
  }

  /** hasPermission 检查是否拥有指定权限码 */
  function hasPermission(code: string) {
    return codes.includes(code)
  }

  /** buildRoutes 根据权限码生成可访问的动态路由 */
  function buildRoutes() {
    routes.value = generateRoutes(codes)
    return routes.value
  }

  return { routes, setPermissions, hasPermission, buildRoutes }
})
