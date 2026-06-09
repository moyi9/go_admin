<!-- 用户角色分配弹窗：复选框列表选择角色，保存时覆盖已有角色 -->
<template>
  <el-dialog v-model="visible" title="分配角色" width="500px">
    <el-checkbox-group v-model="selectedIds">
      <div v-for="role in roles" :key="role.id" style="margin-bottom: 8px">
        <el-checkbox :value="role.id">{{ role.name }} ({{ role.code }})</el-checkbox>
      </div>
    </el-checkbox-group>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { listRolesApi } from '@/api/role'
import { assignRolesApi } from '@/api/user'
import type { User, Role } from '@/types/model'

const props = defineProps<{ modelValue: boolean; user?: User | null }>()
const emit = defineEmits<{ 'update:modelValue': [boolean]; saved: [] }>()

const visible = ref(false)
watch(() => props.modelValue, (v) => { visible.value = v }, { immediate: true })
watch(visible, (v) => emit('update:modelValue', v))

const roles = ref<Role[]>([])
const selectedIds = ref<number[]>([])
const submitting = ref(false)

watch(() => props.modelValue, async (v) => {
  if (v) {
    const res = await listRolesApi(1, 100)
    roles.value = res.data
    if (props.user) {
      selectedIds.value = props.user.roles?.map((r) => r.id) || []
    } else {
      selectedIds.value = []
    }
  }
})

async function handleSubmit() {
  if (!props.user) return
  submitting.value = true
  try {
    await assignRolesApi(props.user.id, selectedIds.value)
    ElMessage.success('角色分配成功')
    visible.value = false
    emit('saved')
  } catch { /* handled */ }
  finally { submitting.value = false }
}
</script>
