<!-- 权限新增/编辑弹窗，包含编码/名称/方法/路径/描述字段 -->
<template>
  <el-dialog v-model="visible" :title="permission ? '编辑权限' : '新增权限'" width="500px" @closed="handleClosed">
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="编码" prop="code">
        <el-input v-model="form.code" :disabled="!!permission" />
      </el-form-item>
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="方法" prop="method">
        <el-select v-model="form.method" style="width: 100%">
          <el-option label="GET" value="GET" />
          <el-option label="POST" value="POST" />
          <el-option label="PUT" value="PUT" />
          <el-option label="DELETE" value="DELETE" />
        </el-select>
      </el-form-item>
      <el-form-item label="路径" prop="path">
        <el-input v-model="form.path" />
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input v-model="form.description" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createPermissionApi, updatePermissionApi } from '@/api/permission'
import type { Permission } from '@/types/model'
import type { FormInstance } from 'element-plus'

const props = defineProps<{ modelValue: boolean; permission?: Permission | null }>()
const emit = defineEmits<{ 'update:modelValue': [boolean]; saved: [] }>()

const visible = ref(false)
watch(() => props.modelValue, (v) => { visible.value = v }, { immediate: true })
watch(visible, (v) => emit('update:modelValue', v))

const formRef = ref<FormInstance>()
const form = reactive({ code: '', name: '', method: 'GET', path: '', description: '' })
const submitting = ref(false)

const rules = {
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }, { min: 2, max: 128, message: '2-128 个字符' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }, { min: 2, max: 128, message: '2-128 个字符' }],
  method: [{ required: true, message: '请选择方法' }],
  path: [{ required: true, message: '请输入路径', trigger: 'blur' }],
}

watch(() => props.modelValue, (v) => {
  if (v && props.permission) {
    form.code = props.permission.code
    form.name = props.permission.name
    form.method = props.permission.method
    form.path = props.permission.path
    form.description = props.permission.description
  } else if (v) {
    form.code = ''; form.name = ''; form.method = 'GET'; form.path = ''; form.description = ''
  }
})

function handleClosed() {
  formRef.value?.resetFields()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (props.permission) {
      await updatePermissionApi(props.permission.id, { name: form.name, method: form.method, path: form.path, description: form.description })
    } else {
      await createPermissionApi({ code: form.code, name: form.name, method: form.method, path: form.path, description: form.description })
    }
    ElMessage.success(props.permission ? '更新成功' : '创建成功')
    visible.value = false
    emit('saved')
  } catch { /* handled */ }
  finally { submitting.value = false }
}
</script>
