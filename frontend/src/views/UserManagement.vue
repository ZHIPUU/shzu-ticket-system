<template>
  <div class="user-mgmt">
    <header class="page-header">
      <div>
        <h1 class="page-title">用户管理</h1>
        <p class="page-desc">管理系统用户、角色和权限</p>
      </div>
      <el-button type="primary" @click="openCreate">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px;margin-right:6px">
          <path d="M12 5v14M5 12h14" />
        </svg>
        新建用户
      </el-button>
    </header>

    <el-card shadow="never" class="user-card">
      <div class="table-wrapper">
      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column label="显示名" width="140">
          <template #default="{ row }">{{ row.display_name || '—' }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="{ row }">{{ row.email || '—' }}</template>
        </el-table-column>
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
              {{ row.role === 'admin' ? '管理员' : '工作人员' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.active ? 'success' : 'warning'" size="small">
              {{ row.active ? '启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="需改密码" width="100">
          <template #default="{ row }">
            <span v-if="row.must_change_password" style="color: #e17055">⚠ 是</span>
            <span v-else style="color: var(--text-tertiary)">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最后登录" width="160">
          <template #default="{ row }">{{ row.last_login_at ? formatTime(row.last_login_at) : '从未' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" link type="warning" @click="onResetPwd(row)">重置密码</el-button>
            <el-button
              size="small"
              link
              :type="row.active ? 'danger' : 'success'"
              :disabled="row.id === currentUserId"
              @click="onToggleActive(row)"
            >
              {{ row.active ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      </div>
    </el-card>

    <!-- 新建用户 -->
    <el-dialog v-model="createVisible" title="新建用户" :width="dialogWidth" :destroy-on-close="true">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="84px">
        <el-form-item label="用户名" required prop="username">
          <el-input v-model="createForm.username" placeholder="3-64 字符" maxlength="64" />
        </el-form-item>
        <el-form-item label="初始密码" required prop="password">
          <el-input v-model="createForm.password" type="password" show-password placeholder="至少 8 位，含字母和数字" />
        </el-form-item>
        <el-form-item label="角色" required prop="role">
          <el-radio-group v-model="createForm.role">
            <el-radio-button label="staff">工作人员</el-radio-button>
            <el-radio-button label="admin">管理员</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="显示名" prop="display_name">
          <el-input v-model="createForm.display_name" maxlength="64" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="createForm.email" maxlength="128" />
        </el-form-item>
      </el-form>
      <el-alert type="info" :closable="false" show-icon style="margin: 0 0 16px">
        <template #title>提示</template>
        新建用户首次登录需强制修改密码
      </el-alert>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="doCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- 编辑用户 -->
    <el-dialog v-model="editVisible" title="编辑用户" :width="dialogWidth" :destroy-on-close="true">
      <el-form v-if="editing" :model="editing" label-width="84px">
        <el-form-item label="用户名">
          <el-input v-model="editing.username" disabled />
        </el-form-item>
        <el-form-item label="显示名">
          <el-input v-model="editing.display_name" maxlength="64" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editing.email" maxlength="128" />
        </el-form-item>
        <el-form-item label="角色">
          <el-radio-group v-model="editing.role">
            <el-radio-button label="staff">工作人员</el-radio-button>
            <el-radio-button label="admin">管理员</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="editing.active" />
        </el-form-item>
        <el-form-item label="强制改密码">
          <el-switch v-model="editing.must_change_password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="doEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listUsers, createUser, updateUser } from '../api/auth'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const currentUserId = computed(() => userStore.user?.id)

const isMobile = computed(() => window.innerWidth < 768)
const dialogWidth = computed(() => isMobile.value ? '95vw' : '480px')

const loading = ref(false)
const users = ref([])
const createVisible = ref(false)
const editVisible = ref(false)
const editing = ref(null)
const submitting = ref(false)

const createFormRef = ref(null)
const createForm = reactive({
  username: '',
  password: '',
  role: 'staff',
  display_name: '',
  email: '',
})

const createRules = {
  username: [{ required: true, min: 3, max: 64, message: '3-64 字符', trigger: 'blur' }],
  password: [{ required: true, min: 8, message: '至少 8 位', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
}

const formatTime = (iso) => {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const fetch = async () => {
  loading.value = true
  try {
    const r = await listUsers()
    users.value = r.items
  } catch (e) {} finally {
    loading.value = false
  }
}

const openCreate = () => {
  createForm.username = ''
  createForm.password = ''
  createForm.role = 'staff'
  createForm.display_name = ''
  createForm.email = ''
  createVisible.value = true
}

const doCreate = async () => {
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await createUser(createForm)
      ElMessage.success('用户创建成功')
      createVisible.value = false
      await fetch()
    } catch (e) {} finally {
      submitting.value = false
    }
  })
}

const openEdit = (row) => {
  editing.value = { ...row }
  editVisible.value = true
}

const doEdit = async () => {
  if (!editing.value) return
  submitting.value = true
  try {
    const payload = {
      display_name: editing.value.display_name,
      email: editing.value.email,
      role: editing.value.role,
      active: editing.value.active,
    }
    if (editing.value.must_change_password) {
      payload.password = 'temp-' + Math.random().toString(36).slice(2, 10)
    }
    await updateUser(editing.value.id, payload)
    ElMessage.success('已更新')
    editVisible.value = false
    await fetch()
  } catch (e) {} finally {
    submitting.value = false
  }
}

const onResetPwd = (row) => {
  ElMessageBox.prompt(
    '请输入新密码（至少 8 位，含字母和数字）',
    `重置 ${row.username} 的密码`,
    {
      inputPattern: /^(?=.*[A-Za-z])(?=.*\d).{8,128}$/,
      inputErrorMessage: '至少 8 位，必须含字母和数字',
      inputValue: '',
      confirmButtonText: '确认重置',
    }
  )
    .then(async ({ value }) => {
      await updateUser(row.id, { password: value })
      ElMessage.success('密码已重置，用户下次登录需修改')
    })
    .catch(() => {})
}

const onToggleActive = async (row) => {
  const action = row.active ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确认${action}用户 ${row.username}？`, `${action}用户`, {
      type: row.active ? 'warning' : 'success',
    })
    await updateUser(row.id, { active: !row.active })
    ElMessage.success(`已${action}`)
    await fetch()
  } catch (e) {
    if (e !== 'cancel') {
      // 错误已拦截
    }
  }
}

onMounted(fetch)
</script>

<style scoped>
.user-mgmt {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}
.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
  letter-spacing: -0.5px;
}
.page-desc {
  color: var(--text-secondary);
  margin: 0;
  font-size: 13px;
}

.user-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

@media (max-width: 768px) {
  .page-title { font-size: 18px; }
}
</style>
