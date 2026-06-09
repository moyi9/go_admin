/** 仪表盘统计 API */
import request from './index'
import type { ApiResponse } from '@/types/api'

/** DashboardStats 仪表盘聚合统计 */
export interface DashboardStats {
  gender_distribution: Record<string, number>
  status_distribution: Record<string, number>
  role_distribution: { role_name: string; user_count: number }[]
}

/** getDashboardStatsApi 获取仪表盘统计数据 */
export async function getDashboardStatsApi() {
  const res = await request.get<ApiResponse<DashboardStats>>('/dashboard/stats')
  return res.data.data!
}
