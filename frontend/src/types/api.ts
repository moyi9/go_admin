/** 标准 API 响应格式 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  request_id?: string
}

/** 分页 API 响应格式 */
export interface PaginatedResponse<T> {
  code: number
  message: string
  data: T[]
  total: number
  page: number
  page_size: number
  request_id?: string
}

/** 登录接口返回结果 */
export interface LoginResult {
  access_token: string
  token_type: string
  user: import('./model').User
}
