<!-- 用户新增/编辑抽屉表单 Drawer，含用户字段校验和角色选择 -->
<template>
  <el-drawer
    v-model="visible"
    :title="user ? '编辑用户' : '新增用户'"
    direction="rtl"
    size="440px"
    @close="handleClose"
  >
    <div style="padding: 0 12px">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item v-if="!user" label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入昵称（选填）" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio value="male">男</el-radio>
            <el-radio value="female">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号（选填）" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="用户状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="用户角色" prop="roleIds">
          <el-select v-model="form.roleIds" multiple placeholder="请选择角色" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { createUserApi, updateUserApi } from '@/api/user'
import { listRolesApi } from '@/api/role'
import type { User, Role } from '@/types/model'
import type { FormInstance } from 'element-plus'

const props = defineProps<{ modelValue: boolean; user?: User | null }>()
const emit = defineEmits<{ 'update:modelValue': [boolean]; saved: [] }>()

const visible = ref(false)
watch(() => props.modelValue, (v) => { visible.value = v }, { immediate: true })
watch(visible, (v) => emit('update:modelValue', v))

const formRef = ref<FormInstance>()
const form = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  gender: '',
  password: '',
  status: 'active',
  roleIds: [] as number[],
})
const submitting = ref(false)
const roles = ref<Role[]>([])

onMounted(async () => {
  try {
    const res = await listRolesApi(1, 100)
    roles.value = res.data
  } catch { /* handled */ }
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度应为 2-20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  phone: [
    {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入有效的手机号',
      trigger: 'blur',
    },
  ],
  status: [
    { required: true, message: '请选择用户状态', trigger: 'change' },
  ],
  roleIds: [
    { required: true, message: '请至少选择一个角色', trigger: 'change', type: 'array', min: 1 },
  ],
}

watch(() => props.modelValue, (v) => {
  if (v && props.user) {
    form.username = props.user.username
    form.nickname = props.user.nickname || ''
    form.email = props.user.email
    form.phone = props.user.phone || ''
    form.gender = props.user.gender || ''
    form.password = ''
    form.status = props.user.status
    form.roleIds = props.user.roles?.map((r) => r.id) || []
  } else if (v) {
    form.username = ''
    form.nickname = ''
    form.email = ''
    form.phone = ''
    form.gender = ''
    form.password = ''
    form.status = 'active'
    form.roleIds = []
  }
})

function handleClose() {
  formRef.value?.resetFields()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (props.user) {
      await updateUserApi(props.user.id, {
        username: form.username,
        nickname: form.nickname || undefined,
        email: form.email,
        phone: form.phone || undefined,
        gender: form.gender || undefined,
        status: form.status,
        role_ids: form.roleIds.length > 0 ? form.roleIds : undefined,
      })
    } else {
      await createUserApi({
        username: form.username,
        nickname: form.nickname || undefined,
        email: form.email,
        phone: form.phone || undefined,
        gender: form.gender || undefined,
        password: form.password,
        status: form.status,
        role_ids: form.roleIds.length > 0 ? form.roleIds : undefined,
      })
    }
    ElMessage.success(props.user ? '更新成功' : '创建成功')
    visible.value = false
    emit('saved')
  } catch { /* handled by interceptor */ }
  finally { submitting.value = false }
}
</script>
