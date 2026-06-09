<!-- 角色权限分配弹窗：使用 ElTransfer 穿梭框选择权限，保存时覆盖已有权限 -->
<template>
  <el-dialog v-model="visible" title="分配权限" width="600px">
    <el-transfer
      v-model="selectedIds"
      :data="permissionOptions"
      :titles="['可分配权限', '已分配权限']"
    />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { listPermissionsApi } from '@/api/permission'
import { assignPermissionsApi } from '@/api/role'
import type { Role, Permission } from '@/types/model'

const props = defineProps<{ modelValue: boolean; role?: Role | null }>()
const emit = defineEmits<{ 'update:modelValue': [boolean]; saved: [] }>()

const visible = ref(false)
watch(() => props.modelValue, (v) => { visible.value = v }, { immediate: true })
watch(visible, (v) => emit('update:modelValue', v))

const allPermissions = ref<Permission[]>([])
const selectedIds = ref<number[]>([])
const submitting = ref(false)

const permissionOptions = computed(() =>
  allPermissions.value.map((p) => ({
    key: p.id,
    label: p.name,
    disabled: false,
  })),
)

watch(() => props.modelValue, async (v) => {
  if (v) {
    const res = await listPermissionsApi(1, 100)
    allPermissions.value = res.data
    if (props.role) {
      selectedIds.value = props.role.permissions?.map((p) => p.id) || []
    } else {
      selectedIds.value = []
    }
  }
})

async function handleSubmit() {
  if (!props.role) return
  submitting.value = true
  try {
    await assignPermissionsApi(props.role.id, selectedIds.value)
    ElMessage.success('权限分配成功')
    visible.value = false
    emit('saved')
  } catch { /* handled */ }
  finally { submitting.value = false }
}
</script>
