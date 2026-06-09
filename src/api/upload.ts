/** 文件上传 API */
import request from './index'
import type { ApiResponse } from '@/types/api'

/** uploadFileApi 上传文件，返回可访问的文件 URL */
export async function uploadFileApi(file: File): Promise<string> {
  const form = new FormData()
  form.append('file', file)
  const res = await request.post<ApiResponse<{ url: string }>>('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data.data!.url
}
