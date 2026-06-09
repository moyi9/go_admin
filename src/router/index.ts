/**
 * Vue Router 配置与导航守卫。
 *
 * 导航流程：
 * 1. /login 免鉴权直接放行
 * 2. 无 Token → 跳转登录页并携带 redirect 参数
 * 3. 用户信息为空 → 从后端获取 /auth/me
 * 4. 首次进入动态路由未注册 → 从权限 store 构建并注册
 * 5. 每次路由切换后 → 更新标签页状态
 */
import { createRouter, createWebHashHistory } from 'vue-router'
import { staticRoutes } from './routes'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { useTabsStore } from '@/stores/tabs'
import type { User } from '@/types/model'

const router = createRouter({
  history: createWebHashHistory(),
  routes: staticRoutes,
})

const dynamicRouteNames = ['Dashboard', 'Profile', 'Users', 'Roles', 'Permissions', 'AuditLogs', 'Forbidden', 'NotFound', 'ServerError']

/** resetRouter 移除动态注册的路由和 CatchAll，用于登出清理 */
export function resetRouter() {
  dynamicRouteNames.forEach((name) => {
    if (router.hasRoute(name)) router.removeRoute(name)
  })
  if (router.hasRoute('CatchAll')) router.removeRoute('CatchAll')
}

/** initDynamicRoutes 根据当前用户权限构建并注册动态路由 */
export function initDynamicRoutes(user: User) {
  const perm = usePermissionStore()
  perm.setPermissions(user)
  const routes = perm.buildRoutes()
  for (const r of routes) router.addRoute('Layout', r)
  router.addRoute({
    path: '/:pathMatch(.*)*',
    name: 'CatchAll',
    component: () => import('@/views/error/NotFound.vue'),
    meta: { hidden: true },
  })
  return routes
}

router.beforeEach(async (to, _from, next) => {
  if (to.path === '/login') return next()

  const auth = useAuthStore()
  if (!auth.token) return next(`/login?redirect=${to.fullPath}`)

  // 尝试获取用户信息，失败则跳转登录
  if (!auth.user) {
    try {
      const user = await auth.getUserInfo()
      if (!user) {
        auth.logout()
        return next('/login')
      }
    } catch {
      return next('/login')
    }
  }

  // 动态路由尚未注册时，根据权限码构建并注册
  const perm = usePermissionStore()
  if (perm.routes.length === 0 && auth.user) {
    initDynamicRoutes(auth.user)
    return next({ ...to, replace: true })
  }

  next()
})

router.afterEach((to) => {
  const tabs = useTabsStore()
  tabs.addTab(to)
})

export default router
