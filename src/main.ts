/** 应用入口：注册 Pinia 状态管理、Vue Router、Element Plus 组件库和自定义指令 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { vPermission } from '@/directives/permission'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 全局注册 Element Plus 图标组件，便于在模板中直接使用 <el-icon><icon-name /></el-icon>
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.directive('permission', vPermission)

app.mount('#app')
