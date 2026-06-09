<!-- 角色新增/编辑抽屉表单 Drawer，包含角色编码/名称/描述/状态字段 -->
<template>
  <el-drawer
    v-model="visible"
    :title="role ? '编辑角色' : '新增角色'"
    direction="rtl"
    size="440px"
    @close="handleClose"
  >
    <div style="padding: 0 12px">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" :disabled="!!role" placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="角色状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
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
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createRoleApi, updateRoleApi } from '@/api/role'
import type { Role } from '@/types/model'
import type { FormInstance } from 'element-plus'

const props = defineProps<{ modelValue: boolean; role?: Role | null }>()
const emit = defineEmits<{ 'update:modelValue': [boolean]; saved: [] }>()

const visible = ref(false)
watch(() => props.modelValue, (v) => { visible.value = v }, { immediate: true })
watch(visible, (v) => emit('update:modelValue', v))

const formRef = ref<FormInstance>()
const form = reactive({
  code: '',
  name: '',
  description: '',
  status: 'active',
})
const submitting = ref(false)

const rules = {
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { min: 2, max: 64, message: '角色编码长度应为 2-64 个字符', trigger: 'blur' },
  ],
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 128, message: '角色名称长度应为 2-128 个字符', trigger: 'blur' },
  ],
  status: [
    { required: true, message: '请选择角色状态', trigger: 'change' },
  ],
}

watch(() => props.modelValue, (v) => {
  if (v && props.role) {
    form.code = props.role.code
    form.name = props.role.name
    form.description = props.role.description
    form.status = props.role.status || 'active'
  } else if (v) {
    form.code = ''
    form.name = ''
    form.description = ''
    form.status = 'active'
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
    if (props.role) {
      await updateRoleApi(props.role.id, {
        name: form.name,
        description: form.description,
        status: form.status,
      })
    } else {
      await createRoleApi({
        code: form.code,
        name: form.name,
        description: form.description,
        status: form.status,
      })
    }
    ElMessage.success(props.role ? '更新成功' : '创建成功')
    visible.value = false
    emit('saved')
  } catch { /* handled */ }
  finally { submitting.value = false }
}
</script>
