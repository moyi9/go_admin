<template>
  <el-dialog
    v-model="visible"
    title="发送通知"
    width="520px"
    :close-on-click-modal="false"
    @closed="handleClosed"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="通知类型" prop="type">
        <el-select v-model="form.type" style="width: 100%">
          <el-option label="系统通知" value="system" />
          <el-option label="安全提醒" value="security" />
        </el-select>
      </el-form-item>
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" placeholder="请输入通知标题" maxlength="255" />
      </el-form-item>
      <el-form-item label="内容" prop="content">
        <el-input
          v-model="form.content"
          type="textarea"
          :rows="4"
          placeholder="请输入通知内容"
          maxlength="2000"
        />
      </el-form-item>
      <el-form-item label="接收者" prop="target_type">
        <el-radio-group v-model="targetType">
          <el-radio value="all">全体用户</el-radio>
          <el-radio value="specific">指定用户</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="targetType === 'specific'" label="选择用户" prop="target_user_id">
        <el-select
          v-model="form.target_user_id"
          filterable
          remote
          reserve-keyword
          placeholder="搜索用户名/邮箱"
          :remote-method="searchUsers"
          :loading="userLoading"
          style="width: 100%"
        >
          <el-option
            v-for="u in userOptions"
            :key="u.id"
            :label="`${u.username} (${u.email})`"
            :value="u.id"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        发送
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { listUsersApi } from '@/api/user'
import { sendNotificationApi } from '@/api/notification'
import type { User } from '@/types/model'

const emit = defineEmits<{ sent: [] }>()

const visible = defineModel<boolean>('visible', { default: false })

const formRef = ref<FormInstance>()
const submitting = ref(false)
const targetType = ref('all')
const userLoading = ref(false)
const userOptions = ref<User[]>([])

const form = reactive({
  type: 'system',
  title: '',
  content: '',
  target_user_id: undefined as number | undefined,
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择通知类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
}

async function searchUsers(query: string) {
  if (!query) return
  userLoading.value = true
  try {
    const res = await listUsersApi(1, 20, { keyword: query })
    userOptions.value = res.data
  } catch {
    userOptions.value = []
  } finally {
    userLoading.value = false
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    await sendNotificationApi({
      type: form.type,
      title: form.title,
      content: form.content,
      target_user_id: targetType.value === 'specific' ? form.target_user_id : undefined,
    })
    ElMessage.success('通知已发送')
    visible.value = false
    emit('sent')
  } catch {
    /* handled by interceptor */
  } finally {
    submitting.value = false
  }
}

function handleClosed() {
  form.type = 'system'
  form.title = ''
  form.content = ''
  form.target_user_id = undefined
  targetType.value = 'all'
  userOptions.value = []
}
</script>
