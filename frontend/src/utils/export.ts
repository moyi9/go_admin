/**
 * 通用数据导出工具（CSV 格式），使用 UTF-8 BOM 确保 Excel 正确识别中文。
 */

/**
 * exportCSV 将数据导出为 CSV 文件并触发下载。
 * @param data     数据数组
 * @param columns  列定义 [{ key, label }]，key 支持点号路径如 "user.name"
 * @param filename 导出文件名（不含扩展名）
 */
export function exportCSV<T extends Record<string, any>>(
  data: T[],
  columns: { key: string; label: string }[],
  filename: string,
): void {
  if (!data.length) return

  const BOM = '﻿'
  const headers = columns.map((c) => c.label).join(',')

  const rows = data.map((row) =>
    columns
      .map((c) => {
        // 支持点号路径取值：如 "roles" → 取 row.roles
        let val: unknown = row
        for (const part of c.key.split('.')) {
          if (val == null) break
          val = (val as Record<string, unknown>)[part]
        }
        const str = val == null ? '' : String(val)
        // 若含逗号/引号/换行则用双引号包裹
        if (str.includes(',') || str.includes('"') || str.includes('\n')) {
          return '"' + str.replace(/"/g, '""') + '"'
        }
        return str
      })
      .join(','),
  )

  const csv = BOM + headers + '\n' + rows.join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${filename}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
