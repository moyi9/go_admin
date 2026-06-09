/**
 * 静态路由配置：登录页、404/403/500 等无需权限的页面。
 * 动态路由（用户管理/角色管理/权限管理）由 permission store 根据权限码动态注册。
 */
import type { RouteRecordRaw } from 'vue-router'

export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'), // 主布局
    name: 'Layout',
    redirect: '/dashboard',
    children: [],
  },
]
