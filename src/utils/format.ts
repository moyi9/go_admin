/** 格式化工具函数 */
import dayjs from 'dayjs'

/** formatTime 将 ISO 时间字符串格式化为可读日期，为空时返回占位符 "-" */
export function formatTime(iso: string | undefined | null, template = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!iso) return '-'
  return dayjs(iso).format(template)
}
