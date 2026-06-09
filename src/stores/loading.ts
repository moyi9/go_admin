/** 全局请求计数器状态管理：用于控制顶部加载条的显示 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const count = ref(0) // 当前正在进行的请求数

  /** show 增加一个请求计数 */
  function show() {
    count.value++
  }

  /** hide 减少一个请求计数，最小为 0 */
  function hide() {
    count.value--
    if (count.value < 0) count.value = 0
  }

  return { count, show, hide }
})
