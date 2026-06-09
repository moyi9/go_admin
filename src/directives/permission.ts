/**
 * v-permission 按钮级权限指令。
 * 用法：`v-permission="'post.api.v1.users'"` 或 `v-permission="['perm1', 'perm2']"`
 * 无权限时直接从 DOM 中移除元素。
 */
import type { Directive } from 'vue'
import { usePermissionStore } from '@/stores/permission'

export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const store = usePermissionStore()
    const codes = Array.isArray(binding.value) ? binding.value : [binding.value]
    const has = codes.some((code) => store.hasPermission(code))
    if (!has) {
      el.parentNode?.removeChild(el)
    }
  },
}
