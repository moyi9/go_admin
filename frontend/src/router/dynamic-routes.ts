/**
 * 动态路由配置：根据后端权限码过滤可访问的菜单页面。
 * generateRoutes() 会根据当前用户权限码列表，只返回有权限的路由。
 */
import type { RouteRecordRaw } from 'vue-router'
import DashboardPage from '@/views/dashboard/index.vue'
import ProfilePage from '@/views/profile/index.vue'
import UserPage from '@/views/user/index.vue'
import RolePage from '@/views/role/index.vue'
import PermissionPage from '@/views/permission/index.vue'
import AuditPage from '@/views/audit/index.vue'
import NotificationPage from '@/views/notification/index.vue'
import ForbiddenPage from '@/views/error/Forbidden.vue'
import NotFoundPage from '@/views/error/NotFound.vue'
import ServerErrorPage from '@/views/error/ServerError.vue'

interface RouteConfig {
  path: string
  name: string
  component?: any
  meta: { title: string; icon: string; permissionCode?: string; hidden?: boolean }
  children?: RouteConfig[]
}

const configs: RouteConfig[] = [
  {
    path: 'dashboard',
    name: 'Dashboard',
    component: DashboardPage,
    meta: { title: '仪表盘', icon: 'Odometer' }, // 首页，无权限要求
  },
  {
    path: 'profile',
    name: 'Profile',
    component: ProfilePage,
    meta: { title: '个人中心', icon: 'User', hidden: true }, // hidden 不在侧边栏显示
  },
  {
    path: 'users',
    name: 'Users',
    component: UserPage,
    meta: { title: '用户管理', icon: 'User', permissionCode: 'get.api.v1.users' },
  },
  {
    path: 'roles',
    name: 'Roles',
    component: RolePage,
    meta: { title: '角色管理', icon: 'Setting', permissionCode: 'get.api.v1.roles' },
  },
  {
    path: 'permissions',
    name: 'Permissions',
    component: PermissionPage,
    meta: { title: '权限管理', icon: 'Lock', permissionCode: 'get.api.v1.permissions' },
  },
  {
    path: 'audit-logs',
    name: 'AuditLogs',
    component: AuditPage,
    meta: { title: '操作日志', icon: 'List', permissionCode: 'get.api.v1.audit-logs' },
  },
  {
    path: 'notifications',
    name: 'Notifications',
    component: NotificationPage,
    meta: { title: '通知中心', icon: 'Bell', permissionCode: 'get.api.v1.notifications' },
  },
  {
    path: 'exception',
    name: 'Exception',
    meta: { title: '异常页面', icon: 'WarningFilled' },
    children: [
      {
        path: '403',
        name: 'Forbidden',
        component: ForbiddenPage,
        meta: { title: '403 禁止访问', icon: 'RemoveFilled' },
      },
      {
        path: '404',
        name: 'NotFound',
        component: NotFoundPage,
        meta: { title: '404 页面未找到', icon: 'QuestionFilled' },
      },
      {
        path: '500',
        name: 'ServerError',
        component: ServerErrorPage,
        meta: { title: '500 服务器错误', icon: 'CircleCloseFilled' },
      },
    ],
  },
]

export function generateRoutes(codes: string[]): RouteRecordRaw[] {
  const result: RouteRecordRaw[] = []
  for (const cfg of configs) {
    if (cfg.meta.permissionCode && !codes.includes(cfg.meta.permissionCode)) continue
    if (cfg.children) {
      const filteredChildren = cfg.children
        .filter((child) => !child.meta.permissionCode || codes.includes(child.meta.permissionCode))
        .map((child) => ({
          path: child.path,
          name: child.name,
          component: child.component,
          meta: { ...child.meta },
        }))
      result.push({ path: cfg.path, name: cfg.name, meta: { ...cfg.meta }, children: filteredChildren })
    } else {
      result.push({ path: cfg.path, name: cfg.name, component: cfg.component, meta: { ...cfg.meta } })
    }
  }
  return result
}

export { configs as routeConfigs }
