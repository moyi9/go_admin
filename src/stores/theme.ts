/** 暗黑主题状态管理：切换 Element Plus 暗黑模式并持久化偏好 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(localStorage.getItem('theme_dark') === 'true')

  /** applyTheme 应用指定主题：切换 <html>.dark 类并持久化 */
  function applyTheme(dark: boolean) {
    isDark.value = dark
    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
    localStorage.setItem('theme_dark', String(dark))
  }

  /** toggle 切换暗黑模式 */
  function toggle() {
    applyTheme(!isDark.value)
  }

  // 初始化时从 localStorage 恢复主题
  applyTheme(isDark.value)

  return { isDark, toggle, applyTheme }
})
