/** 用户模型 */
export interface User {
  id: number
  username: string
  nickname?: string
  email: string
  phone?: string
  gender?: string
  avatar_url?: string
  status: 'active' | 'disabled'
  roles?: Role[]
  created_at: string
  updated_at: string
}

/** 角色模型 */
export interface Role {
  id: number
  code: string
  name: string
  description: string
  status?: string
  permissions?: Permission[]
  created_at: string
  updated_at: string
}

/** 操作日志模型 */
export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  resource: string
  resource_id: string
  detail: string
  ip: string
  user_agent: string
  created_at: string
}

/** 通知模型 */
export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  link?: string
  is_read: boolean
  created_at: string
}

/** 权限模型 */
export interface Permission {
  id: number
  code: string
  name: string
  method: string
  path: string
  description: string
  created_at: string
  updated_at: string
}
