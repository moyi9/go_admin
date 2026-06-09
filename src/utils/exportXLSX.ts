/**
 * 通用数据导出工具（XLSX 格式），使用 ExcelJS 生成带样式的表格。
 */

import ExcelJS from 'exceljs'

/**
 * exportXLSX 将数据导出为 .xlsx 文件并触发下载。
 * @param data     数据数组
 * @param columns  列定义 [{ key, label }]，key 支持点号路径如 "user.name"
 * @param filename 导出文件名（不含扩展名）
 */
export async function exportXLSX<T extends Record<string, any>>(
  data: T[],
  columns: { key: string; label: string }[],
  filename: string,
): Promise<void> {
  const workbook = new ExcelJS.Workbook()
  workbook.creator = 'MyAdminSystem'
  workbook.created = new Date()

  const ws = workbook.addWorksheet('Sheet1')

  // 列定义
  ws.columns = columns.map((c) => ({
    header: c.label,
    key: c.key,
    width: Math.max(c.label.length * 2 + 2, 12), // 估算列宽
  }))

  // 表头样式
  ws.getRow(1).height = 26
  ws.getRow(1).eachCell((cell) => {
    cell.font = { name: 'Microsoft YaHei', bold: true, size: 11, color: { argb: 'FFFFFFFF' } }
    cell.fill = { type: 'pattern', pattern: 'solid', fgColor: { argb: 'FF0D7377' } }
    cell.alignment = { vertical: 'middle', horizontal: 'center' }
    cell.border = {
      top: { style: 'thin', color: { argb: 'FFB0B0B0' } },
      left: { style: 'thin', color: { argb: 'FFB0B0B0' } },
      bottom: { style: 'thin', color: { argb: 'FFB0B0B0' } },
      right: { style: 'thin', color: { argb: 'FFB0B0B0' } },
    }
  })

  // 写入数据行
  data.forEach((row, i) => {
    const values: Record<string, unknown> = {}
    for (const c of columns) {
      let val: unknown = row
      for (const part of c.key.split('.')) {
        if (val == null) break
        val = (val as Record<string, unknown>)[part]
      }
      // 若值本身是对象/数组，转 JSON 字符串
      if (val != null && typeof val === 'object') {
        values[c.key] = JSON.stringify(val)
      } else {
        values[c.key] = val ?? ''
      }
    }
    const rowRef = ws.addRow(values)

    // 数据行样式
    rowRef.height = 22
    rowRef.eachCell((cell) => {
      cell.font = { name: 'Microsoft YaHei', size: 10, color: { argb: 'FF333333' } }
      cell.alignment = { vertical: 'middle' }
      cell.border = {
        top: { style: 'thin', color: { argb: 'FFD0D0D0' } },
        left: { style: 'thin', color: { argb: 'FFD0D0D0' } },
        bottom: { style: 'thin', color: { argb: 'FFD0D0D0' } },
        right: { style: 'thin', color: { argb: 'FFD0D0D0' } },
      }
    })

    // 斑马条纹
    if (i % 2 === 1) {
      rowRef.eachCell((cell) => {
        cell.fill = { type: 'pattern', pattern: 'solid', fgColor: { argb: 'FFF5F7FA' } }
      })
    }
  })

  // 冻结表头行
  ws.views = [{ state: 'frozen', ySplit: 1 }]

  // 生成 buffer 并下载
  const buffer = await workbook.xlsx.writeBuffer()
  const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${filename}.xlsx`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
