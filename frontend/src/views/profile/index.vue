<template>
  <div class="profile-page">
    <!-- 顶部用户信息 -->
    <el-card shadow="never" class="profile-banner">
      <div class="banner-inner">
        <div class="banner-avatar" @click="triggerUpload">
          <img v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" class="avatar-img" alt="头像" />
          <span v-else class="avatar-text">{{ (authStore.user?.username || '用')[0] }}</span>
          <div class="avatar-overlay">
            <el-icon :size="18"><Camera /></el-icon>
            <span>更换头像</span>
          </div>
        </div>
        <input ref="fileInputRef" type="file" accept="image/*" hidden @change="handleFileChange" />
        <div class="banner-meta">
          <h2 class="banner-name">{{ authStore.user?.nickname || authStore.user?.username }}</h2>
          <p class="banner-email">{{ authStore.user?.email }}</p>
          <div class="banner-roles">
            <el-tag
              v-for="role in (authStore.user?.roles || [])"
              :key="role.id"
              size="small"
              class="role-tag"
            >{{ role.name }}</el-tag>
            <span v-if="!(authStore.user?.roles || []).length" class="no-role">暂无角色</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 左侧导航 + 右侧面板 -->
    <div class="profile-body">
      <aside class="profile-aside">
        <nav class="aside-nav">
          <a
            v-for="item in navItems"
            :key="item.key"
            class="nav-item"
            :class="{ active: activeTab === item.key }"
            @click="activeTab = item.key"
          >
            <span class="nav-icon" v-html="item.icon" />
            <span class="nav-label">{{ item.label }}</span>
          </a>
        </nav>
      </aside>

      <div class="profile-main">
        <!-- 个人资料 -->
        <transition name="fade" mode="out-in">
          <div v-if="activeTab === 'profile'" key="profile" class="main-panel">
            <div class="panel-header">
              <span class="panel-title">个人资料</span>
              <p class="panel-desc">编辑你的个人信息</p>
            </div>
            <el-form
              ref="profileFormRef"
              :model="profileForm"
              :rules="profileRules"
              label-width="80px"
              class="panel-form"
            >
              <el-form-item label="昵称" prop="nickname">
                <el-input v-model="profileForm.nickname" placeholder="请输入昵称" maxlength="64" />
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="profileForm.email" placeholder="请输入邮箱" maxlength="128" />
              </el-form-item>
              <el-form-item label="手机号" prop="phone">
                <el-input v-model="profileForm.phone" placeholder="请输入手机号" maxlength="20" />
              </el-form-item>
              <el-form-item label="性别" prop="gender">
                <el-select v-model="profileForm.gender" placeholder="请选择性别" style="width: 100%">
                  <el-option label="保密" value="unknown" />
                  <el-option label="男" value="male" />
                  <el-option label="女" value="female" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="profileLoading" @click="handleUpdateProfile">
                  保存修改
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 修改密码 -->
          <div v-else key="password" class="main-panel">
            <div class="panel-header">
              <span class="panel-title">修改密码</span>
              <p class="panel-desc">请定期更换密码以保障账户安全</p>
            </div>
            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-width="100px"
              class="panel-form"
            >
              <el-form-item label="当前密码" prop="current_password">
                <el-input
                  v-model="passwordForm.current_password"
                  type="password"
                  placeholder="请输入当前密码"
                  show-password
                />
              </el-form-item>
              <el-form-item label="新密码" prop="new_password">
                <el-input
                  v-model="passwordForm.new_password"
                  type="password"
                  placeholder="至少6位，建议包含字母和数字"
                  show-password
                />
              </el-form-item>
              <el-form-item label="确认新密码" prop="confirm_password">
                <el-input
                  v-model="passwordForm.confirm_password"
                  type="password"
                  placeholder="请再次输入新密码"
                  show-password
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="passwordLoading" @click="handleUpdatePassword">
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'Profile' })
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Camera } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { updateProfileApi, updatePasswordApi } from '@/api/auth'
import { uploadFileApi } from '@/api/upload'

const authStore = useAuthStore()

const fileInputRef = ref<HTMLInputElement | null>(null)
const uploadingAvatar = ref(false)

function triggerUpload() {
  fileInputRef.value?.click()
}

async function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  uploadingAvatar.value = true
  try {
    const url = await uploadFileApi(file)
    const user = await updateProfileApi({ avatar_url: url })
    authStore.updateUser(user)
    ElMessage.success('头像已更新')
  } catch {
    /* handled by interceptor */
  } finally {
    uploadingAvatar.value = false
    input.value = ''
  }
}

/* ===== 左侧导航 ===== */
const activeTab = ref('profile')

interface NavItem {
  key: string
  label: string
  icon: string
}

const navItems: NavItem[] = [
  { key: 'profile', label: '个人资料', icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>' },
  { key: 'password', label: '修改密码', icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>' },
]

/* ===== 个人信息表单 ===== */
const profileFormRef = ref<FormInstance>()
void profileFormRef
const profileLoading = ref(false)

const profileForm = reactive({
  nickname: '',
  email: '',
  phone: '',
  gender: 'unknown',
})

const profileRules: FormRules = {
  email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
}

/* ===== 密码表单 ===== */
const passwordFormRef = ref<FormInstance>()
void passwordFormRef
const passwordLoading = ref(false)

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: '',
})

const validateConfirm = (_rule: any, value: string, callback: (e?: Error) => void) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules: FormRules = {
  current_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' },
  ],
}

onMounted(() => {
  const u = authStore.user
  if (u) {
    profileForm.nickname = u.nickname || ''
    profileForm.email = u.email || ''
    profileForm.phone = u.phone || ''
    profileForm.gender = u.gender || 'unknown'
  }
})

async function handleUpdateProfile() {
  const valid = await profileFormRef.value?.validate().catch(() => false)
  if (!valid) return
  profileLoading.value = true
  try {
    const user = await updateProfileApi({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone,
      gender: profileForm.gender,
    })
    authStore.updateUser(user)
    ElMessage.success('个人信息已更新')
  } catch {
    /* handled by interceptor */
  } finally {
    profileLoading.value = false
  }
}

async function handleUpdatePassword() {
  const valid = await passwordFormRef.value?.validate().catch(() => false)
  if (!valid) return
  passwordLoading.value = true
  try {
    await updatePasswordApi({
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password,
    })
    ElMessage.success('密码已修改，下次登录请使用新密码')
    passwordFormRef.value?.resetFields()
  } catch {
    /* handled by interceptor */
  } finally {
    passwordLoading.value = false
  }
}
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ===== 顶部 Banner ===== */
.profile-banner {
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.profile-banner :deep(.el-card__body) {
  padding: 0;
}

.banner-inner {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 28px 32px;
  background: linear-gradient(135deg, rgba(13, 115, 119, 0.04) 0%, transparent 60%);
  position: relative;
}

.banner-inner::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: #0d7377;
  border-radius: 0 2px 2px 0;
}

.banner-avatar {
  position: relative;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: linear-gradient(135deg, #0d7377, #14a0a5);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 16px rgba(13, 115, 119, 0.25);
  cursor: pointer;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  color: #fff;
  font-size: 11px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.banner-avatar:hover .avatar-overlay {
  opacity: 1;
}

.avatar-text {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
}

.banner-meta {
  flex: 1;
}

.banner-name {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 4px;
  color: var(--el-text-color-primary);
}

.banner-email {
  margin: 0 0 10px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.banner-roles {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.role-tag {
  border-radius: 4px;
}

.no-role {
  font-size: 13px;
  color: var(--el-text-color-placeholder);
}

/* ===== 主体区域：左侧导航 + 右侧面板 ===== */
.profile-body {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

/* ===== 左侧导航 ===== */
.profile-aside {
  width: 200px;
  flex-shrink: 0;
  background: var(--el-bg-color);
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  padding: 12px 0;
  position: sticky;
  top: 80px;
}

.aside-nav {
  display: flex;
  flex-direction: column;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  font-size: 14px;
  color: var(--el-text-color-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  border-left: 3px solid transparent;
  text-decoration: none;
}

.nav-item:hover {
  color: var(--el-text-color-primary);
  background: rgba(13, 115, 119, 0.04);
}

.nav-item.active {
  color: #0d7377;
  background: rgba(13, 115, 119, 0.06);
  border-left-color: #0d7377;
  font-weight: 600;
}

.nav-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.nav-item.active .nav-icon {
  color: #0d7377;
}

/* ===== 右侧主面板 ===== */
.profile-main {
  flex: 1;
  min-width: 0;
  background: var(--el-bg-color);
  border-radius: 10px;
  border: 1px solid var(--el-border-color-light);
  padding: 28px 32px;
}

.main-panel {
  max-width: 560px;
}

.panel-header {
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.panel-title {
  display: block;
  font-size: 17px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 6px;
}

.panel-desc {
  margin: 0;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.panel-form {
  max-width: 440px;
}

/* ===== 切换动画 ===== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
