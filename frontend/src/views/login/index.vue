<template>
  <div class="login-page">
    <!-- 背景装饰：浮动光晕 + 几何纹理 -->
    <div class="bg-orb bg-orb--1" />
    <div class="bg-orb bg-orb--2" />
    <div class="bg-grid" />
    <div class="bg-noise" />

    <div class="login-container">
      <div class="login-card-wrapper">
        <!-- 左侧品牌区 -->
        <div class="login-brand">
          <div class="brand-embleme">
            <svg viewBox="0 0 64 64" class="brand-svg">
              <rect x="4" y="4" width="56" height="56" rx="8" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.6" />
              <rect x="14" y="14" width="36" height="36" rx="4" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4" />
              <path d="M32 8 L32 56 M8 32 L56 32" stroke="currentColor" stroke-width="1.2" opacity="0.25" />
              <circle cx="32" cy="32" r="8" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.5" />
            </svg>
          </div>
          <h1 class="brand-title">后台管理控制</h1>
          <p class="brand-subtitle">企业级权限管理平台</p>
          <div class="brand-divider" />
          <p class="brand-desc">
            安全 · 高效 · 可扩展
          </p>
        </div>

        <!-- 右侧登录表单 -->
        <div class="login-form-area">
          <div class="form-header">
            <h2 class="form-title">欢迎登录</h2>
            <p class="form-desc">请输入您的账号信息</p>
          </div>
          <el-form ref="formRef" :model="form" :rules="rules" @keyup.enter="handleLogin" class="login-form">
            <el-form-item prop="username">
              <el-input
                v-model="form.username"
                placeholder="用户名"
                :prefix-icon="User"
                size="large"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                placeholder="密码"
                :prefix-icon="Lock"
                show-password
                size="large"
              />
            </el-form-item>
            <el-form-item>
              <el-checkbox v-model="form.remember" class="login-checkbox">记住我</el-checkbox>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" class="login-btn" size="large" @click="handleLogin">
                {{ loading ? '登录中...' : '登 录' }}
              </el-button>
            </el-form-item>
            <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { initDynamicRoutes } from '@/router'
import type { FormInstance } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = reactive({ username: '', password: '', remember: false })
const loading = ref(false)
const errorMsg = ref('')
const formRef = ref<FormInstance>()

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

/** handleLogin 校验表单 -> 登录 -> 跳转 redirect 或首页 */
async function handleLogin() {
  errorMsg.value = ''
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authStore.login(form.username, form.password, form.remember)
    // 在导航前先构建并注册动态路由，避免路由守卫中的 addRoute 时序问题
    if (authStore.user) {
      initDynamicRoutes(authStore.user)
    }
    const redirect = (route.query.redirect as string) || '/'
    // 防止开放重定向：只允许相对路径
    const safeRedirect = redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/'
    router.push(safeRedirect)
  } catch (e: any) {
    errorMsg.value = e?.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* ========== 全局变量 ========== */
.login-page {
  --login-primary: #0d7377;
  --login-primary-light: #14a0a5;
  --login-primary-dark: #0a5c5f;
  --login-warm: #c44536;
  --login-bg-dark: #1a1d23;
  --login-bg-mid: #2a2e38;
  --login-text-light: rgba(255, 255, 255, 0.9);
  --login-text-muted: rgba(255, 255, 255, 0.55);
  --login-card-bg: rgba(26, 29, 35, 0.72);
  --login-card-border: rgba(255, 255, 255, 0.07);
  --login-glass-shadow: 0 24px 80px rgba(0, 0, 0, 0.45);
  --login-form-bg: rgba(255, 255, 255, 0.04);

  position: relative;
  min-height: 100vh;
  background: var(--login-bg-dark);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "PingFang SC", "Noto Sans SC", "Microsoft YaHei", -apple-system, sans-serif;
}

/* ========== 背景装饰 ========== */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  pointer-events: none;
  z-index: 0;
}
.bg-orb--1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(13, 115, 119, 0.3), transparent 70%);
  top: -200px;
  right: -150px;
  animation: orbFloat1 20s ease-in-out infinite;
}
.bg-orb--2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(196, 69, 54, 0.15), transparent 70%);
  bottom: -150px;
  left: -150px;
  animation: orbFloat2 25s ease-in-out infinite;
}

@keyframes orbFloat1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(60px, 40px) scale(1.05); }
  66% { transform: translate(-30px, -20px) scale(0.95); }
}
@keyframes orbFloat2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(-40px, -30px) scale(1.08); }
  66% { transform: translate(50px, 20px) scale(0.92); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  z-index: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.025) 1px, transparent 1px);
  background-size: 60px 60px;
  mask-image: radial-gradient(ellipse at 50% 50%, black 30%, transparent 70%);
  -webkit-mask-image: radial-gradient(ellipse at 50% 50%, black 30%, transparent 70%);
}

.bg-noise {
  position: absolute;
  inset: 0;
  z-index: 0;
  opacity: 0.03;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
  background-size: 256px 256px;
  pointer-events: none;
}

/* ========== 容器 ========== */
.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 920px;
  padding: 24px;
}

.login-card-wrapper {
  display: flex;
  border-radius: 16px;
  overflow: hidden;
  background: var(--login-card-bg);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid var(--login-card-border);
  box-shadow: var(--login-glass-shadow);
  animation: cardEntrance 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes cardEntrance {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* ========== 左侧品牌区 ========== */
.login-brand {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  color: var(--login-text-light);
  position: relative;
  overflow: hidden;
}

.login-brand::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(160deg, rgba(13, 115, 119, 0.2) 0%, transparent 60%);
  pointer-events: none;
}

.brand-embleme {
  width: 72px;
  height: 72px;
  margin-bottom: 24px;
  color: var(--login-primary-light);
  animation: emblemeEntrance 0.8s 0.15s cubic-bezier(0.16, 1, 0.3, 1) both;
}

.brand-svg {
  width: 100%;
  height: 100%;
}

.brand-title {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 2px;
  margin: 0 0 8px;
  animation: fadeUp 0.6s 0.3s both;
}

.brand-subtitle {
  font-size: 14px;
  color: var(--login-text-muted);
  margin: 0 0 24px;
  letter-spacing: 1px;
  animation: fadeUp 0.6s 0.4s both;
}

.brand-divider {
  width: 40px;
  height: 2px;
  background: linear-gradient(90deg, var(--login-primary-light), transparent);
  margin-bottom: 20px;
  animation: fadeUp 0.6s 0.5s both;
}

.brand-desc {
  font-size: 13px;
  color: var(--login-text-muted);
  letter-spacing: 4px;
  margin: 0;
  animation: fadeUp 0.6s 0.6s both;
}

@keyframes emblemeEntrance {
  from { opacity: 0; transform: scale(0.5) rotate(-20deg); }
  to { opacity: 1; transform: scale(1) rotate(0); }
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* ========== 右侧表单区 ========== */
.login-form-area {
  width: 400px;
  padding: 48px 40px;
  background: var(--login-form-bg);
  backdrop-filter: blur(8px);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-header {
  margin-bottom: 32px;
}

.form-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--login-text-light);
  margin: 0 0 6px;
}

.form-desc {
  font-size: 14px;
  color: var(--login-text-muted);
  margin: 0;
}

/* ========== 表单样式覆盖 ========== */
.login-form {
  width: 100%;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.login-form :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  box-shadow: none !important;
  padding: 4px 16px;
  transition: all 0.25s ease;
}

.login-form :deep(.el-input__wrapper:hover) {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.08);
}

.login-form :deep(.el-input__wrapper.is-focus) {
  border-color: var(--login-primary-light);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(13, 115, 119, 0.15) !important;
}

.login-form :deep(.el-input__inner) {
  color: var(--login-text-light);
  height: 44px;
  font-size: 15px;
}

.login-form :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.3);
}

.login-form :deep(.el-input__prefix) {
  margin-right: 8px;
}

.login-form :deep(.el-input__prefix-inner) {
  color: rgba(255, 255, 255, 0.35);
}

/* 密码可见切换 */
.login-form :deep(.el-input__suffix-inner) {
  color: rgba(255, 255, 255, 0.35);
}

.login-checkbox {
  color: var(--login-text-muted);
  font-size: 14px;
}

.login-checkbox :deep(.el-checkbox__label) {
  color: var(--login-text-muted);
  font-size: 14px;
}

.login-checkbox :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--login-primary-light);
  border-color: var(--login-primary-light);
}

.login-btn {
  width: 100%;
  height: 46px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 4px;
  border: none;
  background: linear-gradient(135deg, var(--login-primary) 0%, var(--login-primary-light) 100%);
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(13, 115, 119, 0.35);
}

.login-btn:active {
  transform: translateY(0);
}

.error-msg {
  text-align: center;
  color: #f56c6c;
  font-size: 14px;
  padding: 8px 12px;
  background: rgba(245, 108, 108, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(245, 108, 108, 0.2);
  animation: shakeX 0.4s ease;
}

@keyframes shakeX {
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-6px); }
  40% { transform: translateX(6px); }
  60% { transform: translateX(-4px); }
  80% { transform: translateX(4px); }
}
</style>
