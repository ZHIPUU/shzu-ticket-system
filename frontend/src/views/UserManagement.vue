<template>
  <div class="users-page">
    <!-- ═══ 页头 ═══ -->
    <header class="page-head anim-fade-up">
      <div class="head-text">
        <h1>用户管理</h1>
        <p>管理系统用户、角色和权限</p>
      </div>
      <UiButton variant="primary" :icon-only="isMobile" @click="openCreate">
        <template #icon><UserRoundPlus :size="16" :stroke-width="2.2" /></template>
        <span v-if="!isMobile">新建用户</span>
      </UiButton>
    </header>

    <!-- ═══ 列表 ═══ -->
    <div class="users-card card anim-stagger" style="--i: 1">
      <div v-if="loading" class="list-loading"><UiSkeleton type="rows" :count="4" /></div>

      <UiEmpty v-else-if="!users.length" icon="Users" title="暂无用户" description="点击右上角新建第一个用户" />

      <!-- 桌面表格 -->
      <template v-else-if="!isMobile">
        <div class="t-head">
          <div class="t-cell c-user">用户</div>
          <div class="t-cell c-email">邮箱</div>
          <div class="t-cell c-role">角色</div>
          <div class="t-cell c-status">状态</div>
          <div class="t-cell c-login">最后登录</div>
          <div class="t-cell c-ops">操作</div>
        </div>
        <div v-for="(u, i) in users" :key="u.id" class="t-row anim-stagger" :style="{ '--i': i }">
          <div class="t-cell c-user">
            <div class="u-avatar" :data-role="u.role">{{ letter(u) }}</div>
            <div class="u-names">
              <span class="u-display">{{ u.display_name || u.username }}</span>
              <span class="u-username mono">@{{ u.username }}</span>
            </div>
            <span v-if="u.must_change_password" class="pwd-warn" title="下次登录需修改密码">
              <ShieldAlert :size="13" :stroke-width="2.2" />
            </span>
          </div>
          <div class="t-cell c-email"><span class="ellipsis">{{ u.email || '—' }}</span></div>
          <div class="t-cell c-role"><span class="role-pill" :data-role="u.role">{{ u.role === 'admin' ? '管理员' : '工作人员' }}</span></div>
          <div class="t-cell c-status"><span class="state-pill" :class="{ off: !u.active }">{{ u.active ? '启用' : '已禁用' }}</span></div>
          <div class="t-cell c-login">
            <TimeText v-if="u.last_login_at" :value="u.last_login_at" />
            <span v-else class="muted">从未登录</span>
          </div>
          <div class="t-cell c-ops">
            <UiButton variant="text" size="sm" @click="openEdit(u)">
              <template #icon><Pencil :size="13" :stroke-width="2" /></template>编辑
            </UiButton>
            <UiButton variant="text" size="sm" @click="onResetPwd(u)">
              <template #icon><KeyRound :size="13" :stroke-width="2" /></template>重置密码
            </UiButton>
            <UiButton
              variant="text"
              size="sm"
              :class="u.active ? 'op-danger' : 'op-success'"
              :disabled="u.id === currentUserId"
              @click="onToggleActive(u)"
            >
              <template #icon>
                <Ban v-if="u.active" :size="13" :stroke-width="2" />
                <CircleCheck v-else :size="13" :stroke-width="2" />
              </template>
              {{ u.active ? '禁用' : '启用' }}
            </UiButton>
          </div>
        </div>
      </template>

      <!-- 移动卡片 -->
      <template v-else>
        <div v-for="(u, i) in users" :key="u.id" class="m-user anim-stagger" :style="{ '--i': i }">
          <div class="u-avatar" :data-role="u.role">{{ letter(u) }}</div>
          <div class="m-user-body">
            <div class="m-user-top">
              <span class="u-display">{{ u.display_name || u.username }}</span>
              <span v-if="u.must_change_password" class="pwd-warn" title="下次登录需修改密码">
                <ShieldAlert :size="13" :stroke-width="2.2" />
              </span>
              <span class="state-pill sm" :class="{ off: !u.active }">{{ u.active ? '启用' : '禁用' }}</span>
            </div>
            <div class="u-username mono">@{{ u.username }}<template v-if="u.email"> · {{ u.email }}</template></div>
            <div class="m-user-meta">
              <span class="role-pill" :data-role="u.role">{{ u.role === 'admin' ? '管理员' : '工作人员' }}</span>
              <span class="muted last-login">
                <TimeText v-if="u.last_login_at" :value="u.last_login_at" mode="relative" />
                <template v-else>从未登录</template>
              </span>
            </div>
          </div>
          <UiDropdown :items="userActions(u)" placement="end" @select="(cmd) => onUserAction(cmd, u)">
            <template #trigger>
              <button class="m-kebab" aria-label="更多操作"><EllipsisVertical :size="17" :stroke-width="2" /></button>
            </template>
          </UiDropdown>
        </div>
      </template>
    </div>

    <!-- ═══ 新建用户 ═══ -->
    <UiModal v-model="createVisible" title="新建用户" width="480px">
      <div class="form-stack">
        <div class="form-item">
          <label class="form-label">用户名 <em>*</em></label>
          <UiInput v-model="createForm.username" :maxlength="64" placeholder="3-64 个字符" clearable />
        </div>
        <div class="form-item">
          <label class="form-label">初始密码 <em>*</em></label>
          <UiInput v-model="createForm.password" type="password" placeholder="至少 8 位，含字母和数字" />
        </div>
        <div class="form-item">
          <label class="form-label">角色 <em>*</em></label>
          <UiSegmented v-model="createForm.role" :options="roleOptions" />
        </div>
        <div class="form-item">
          <label class="form-label">显示名</label>
          <UiInput v-model="createForm.display_name" :maxlength="64" placeholder="选填" clearable />
        </div>
        <div class="form-item">
          <label class="form-label">邮箱</label>
          <UiInput v-model="createForm.email" :maxlength="128" placeholder="选填" clearable />
        </div>
        <p class="form-note">
          <Info :size="13" :stroke-width="2" />
          新建用户首次登录需强制修改密码
        </p>
      </div>
      <template #footer>
        <UiButton variant="ghost" @click="createVisible = false">取消</UiButton>
        <UiButton variant="primary" :loading="submitting" @click="doCreate">创建</UiButton>
      </template>
    </UiModal>

    <!-- ═══ 编辑用户 ═══ -->
    <UiModal v-model="editVisible" title="编辑用户" width="480px">
      <div v-if="editing" class="form-stack">
        <div class="form-item">
          <label class="form-label">用户名</label>
          <UiInput v-model="editing.username" disabled />
        </div>
        <div class="form-item">
          <label class="form-label">显示名</label>
          <UiInput v-model="editing.display_name" :maxlength="64" clearable />
        </div>
        <div class="form-item">
          <label class="form-label">邮箱</label>
          <UiInput v-model="editing.email" :maxlength="128" clearable />
        </div>
        <div class="form-item">
          <label class="form-label">角色</label>
          <UiSegmented v-model="editing.role" :options="roleOptions" />
        </div>
        <div class="switch-row">
          <div class="switch-text">
            <span class="switch-title">启用账号</span>
            <span class="switch-desc">禁用后该用户无法登录系统</span>
          </div>
          <UiSwitch v-model="editing.active" :disabled="editing.id === currentUserId" />
        </div>
        <div class="switch-row">
          <div class="switch-text">
            <span class="switch-title">强制修改密码</span>
            <span class="switch-desc">开启后将生成临时密码，用户下次登录需重设</span>
          </div>
          <UiSwitch v-model="editing.must_change_password" />
        </div>
      </div>
      <template #footer>
        <UiButton variant="ghost" @click="editVisible = false">取消</UiButton>
        <UiButton variant="primary" :loading="submitting" @click="doEdit">保存</UiButton>
      </template>
    </UiModal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import {
  UserRoundPlus, Pencil, KeyRound, Ban, CircleCheck, ShieldAlert,
  EllipsisVertical, Info,
} from '@lucide/vue'
import { listUsers, createUser, updateUser } from '../api/auth'
import { useUserStore } from '../stores/user'
import UiButton from '../ui/UiButton.vue'
import UiInput from '../ui/UiInput.vue'
import UiModal from '../ui/UiModal.vue'
import UiSegmented from '../ui/UiSegmented.vue'
import UiSwitch from '../ui/UiSwitch.vue'
import UiDropdown from '../ui/UiDropdown.vue'
import UiSkeleton from '../ui/UiSkeleton.vue'
import UiEmpty from '../ui/UiEmpty.vue'
import TimeText from '../ui/TimeText.vue'
import { toast } from '../ui/toast'
import { confirmDialog, promptDialog } from '../ui/confirm'
import { useIsMobile } from '../composables/useMediaQuery'

const userStore = useUserStore()
const isMobile = useIsMobile()
const currentUserId = computed(() => userStore.user?.id)

const loading = ref(false)
const users = ref([])
const createVisible = ref(false)
const editVisible = ref(false)
const editing = ref(null)
const submitting = ref(false)

const createForm = reactive({ username: '', password: '', role: 'staff', display_name: '', email: '' })

const roleOptions = [
  { label: '工作人员', value: 'staff' },
  { label: '管理员', value: 'admin' },
]

const PWD_PATTERN = /^(?=.*[A-Za-z])(?=.*\d).{8,128}$/
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const letter = (u) => (u.display_name?.[0] || u.username?.[0] || '?').toUpperCase()

const userActions = (u) => [
  { label: '编辑资料', value: 'edit', icon: Pencil },
  { label: '重置密码', value: 'reset', icon: KeyRound },
  {
    label: u.active ? '禁用账号' : '启用账号',
    value: 'toggle',
    icon: u.active ? Ban : CircleCheck,
    danger: u.active,
    disabled: u.id === currentUserId.value,
    divided: true,
  },
]

const onUserAction = (cmd, u) => {
  if (cmd === 'edit') openEdit(u)
  else if (cmd === 'reset') onResetPwd(u)
  else if (cmd === 'toggle') onToggleActive(u)
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

// ─── 新建 ───
const openCreate = () => {
  Object.assign(createForm, { username: '', password: '', role: 'staff', display_name: '', email: '' })
  createVisible.value = true
}

const doCreate = async () => {
  const f = createForm
  if (f.username.trim().length < 3 || f.username.trim().length > 64) {
    toast.warning('用户名需为 3-64 个字符'); return
  }
  if (!PWD_PATTERN.test(f.password)) {
    toast.warning('密码至少 8 位，且必须包含字母和数字'); return
  }
  if (f.email && !EMAIL_PATTERN.test(f.email)) {
    toast.warning('邮箱格式不正确'); return
  }
  submitting.value = true
  try {
    await createUser({ ...f, username: f.username.trim() })
    toast.success('用户创建成功')
    createVisible.value = false
    await fetch()
  } catch (e) {} finally {
    submitting.value = false
  }
}

// ─── 编辑 ───
const openEdit = (row) => {
  editing.value = { ...row }
  editVisible.value = true
}

const doEdit = async () => {
  if (!editing.value) return
  if (editing.value.email && !EMAIL_PATTERN.test(editing.value.email)) {
    toast.warning('邮箱格式不正确'); return
  }
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
    toast.success('已更新')
    editVisible.value = false
    await fetch()
  } catch (e) {} finally {
    submitting.value = false
  }
}

// ─── 重置密码 ───
const onResetPwd = async (row) => {
  const val = await promptDialog({
    title: `重置密码`,
    message: `为用户 ${row.username} 设置新密码（至少 8 位，含字母和数字）`,
    placeholder: '输入新密码…',
    pattern: PWD_PATTERN,
    patternMessage: '至少 8 位，必须含字母和数字',
    confirmText: '确认重置',
  })
  if (val === null) return
  try {
    await updateUser(row.id, { password: val })
    toast.success('密码已重置，用户下次登录需修改')
  } catch (e) {}
}

// ─── 启用/禁用 ───
const onToggleActive = async (row) => {
  const action = row.active ? '禁用' : '启用'
  const ok = await confirmDialog({
    title: `${action}用户`,
    message: `确认${action}用户 ${row.username}？${row.active ? '禁用后其将无法登录。' : ''}`,
    confirmText: `确认${action}`,
    danger: row.active,
  })
  if (!ok) return
  try {
    await updateUser(row.id, { active: !row.active })
    toast.success(`已${action}`)
    await fetch()
  } catch (e) {}
}

onMounted(fetch)
</script>

<style scoped>
.users-page { display: flex; flex-direction: column; gap: 18px; }

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.head-text h1 {
  margin: 0 0 3px;
  font-size: 23px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--text-1);
}
.head-text p { margin: 0; font-size: 13px; color: var(--text-3); }
@media (max-width: 767px) { .head-text h1 { font-size: 19px; } }

.users-card { overflow: hidden; }
.list-loading { padding: 8px 20px; }

/* ─── 桌面表格 ─── */
.t-head, .t-row {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1.6fr) 100px 90px 140px 250px;
  align-items: center;
  gap: 14px;
  padding: 0 22px;
}
.t-head {
  height: 42px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-3);
  border-bottom: 1px solid var(--border-soft);
  background: var(--bg-sink);
}
.t-row {
  min-height: 66px;
  border-bottom: 1px solid var(--border-soft);
  transition: background var(--d-fast) var(--ease-out);
}
.t-row:last-child { border-bottom: none; }
.t-row:hover { background: var(--bg-hover); }
.t-cell { min-width: 0; display: flex; align-items: center; }

.u-avatar {
  width: 36px;
  height: 36px;
  border-radius: 11px;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}
.u-avatar[data-role="admin"] { background: linear-gradient(135deg, #F59E0B, #F97316); }
.c-user { gap: 11px; }
.u-names { display: flex; flex-direction: column; min-width: 0; line-height: 1.35; }
.u-display { font-size: 13.5px; font-weight: 600; color: var(--text-1); }
.u-username { font-size: 11.5px; color: var(--text-3); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.c-email { font-size: 13px; color: var(--text-2); }

.pwd-warn { color: var(--warning); display: inline-flex; flex-shrink: 0; }

.role-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: var(--r-full);
  font-size: 11.5px;
  font-weight: 600;
  background: var(--info-soft);
  color: var(--info);
  white-space: nowrap;
}
.role-pill[data-role="admin"] { background: var(--warning-soft); color: var(--warning); }

.state-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: var(--r-full);
  font-size: 11.5px;
  font-weight: 600;
  background: var(--success-soft);
  color: var(--success);
  white-space: nowrap;
}
.state-pill.off { background: var(--danger-soft); color: var(--danger); }
.state-pill.sm { padding: 2px 8px; font-size: 10.5px; }

.c-ops { gap: 2px; }
.op-danger { color: var(--danger) !important; }
.op-danger:hover { background: var(--danger-soft) !important; }
.op-success { color: var(--success) !important; }
.op-success:hover { background: var(--success-soft) !important; }

/* ─── 移动卡片 ─── */
.m-user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-soft);
}
.m-user:last-child { border-bottom: none; }
.m-user-body { flex: 1; min-width: 0; }
.m-user-top { display: flex; align-items: center; gap: 8px; margin-bottom: 2px; }
.m-user-meta { display: flex; align-items: center; gap: 10px; margin-top: 7px; }
.last-login { font-size: 11.5px; }
.m-kebab {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  flex-shrink: 0;
}
.m-kebab:hover { background: var(--bg-hover); color: var(--text-1); }

/* ─── 表单 ─── */
.form-stack { display: flex; flex-direction: column; gap: 15px; }
.form-item { display: flex; flex-direction: column; gap: 7px; }
.form-label { font-size: 13px; font-weight: 500; color: var(--text-2); }
.form-label em { color: var(--danger); font-style: normal; }

.form-note {
  display: flex;
  align-items: center;
  gap: 7px;
  margin: 0;
  padding: 10px 13px;
  border-radius: var(--r-md);
  background: var(--primary-soft);
  color: var(--text-2);
  font-size: 12.5px;
}
.form-note svg { color: var(--primary); flex-shrink: 0; }

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 12px 14px;
  border: 1px solid var(--border-soft);
  border-radius: var(--r-md);
  background: var(--bg-sink);
}
.switch-text { display: flex; flex-direction: column; gap: 2px; }
.switch-title { font-size: 13.5px; font-weight: 500; color: var(--text-1); }
.switch-desc { font-size: 12px; color: var(--text-3); }
</style>
