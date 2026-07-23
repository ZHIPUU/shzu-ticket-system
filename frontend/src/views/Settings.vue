<template>
  <div class="settings">
    <header class="page-header">
      <h1 class="page-title">个人设置</h1>
      <p class="page-desc">管理账号密码与界面偏好</p>
    </header>

    <!-- 修改密码（如果首登强制改密码时高亮提示） -->
    <el-card class="settings-card" :class="{ alert: mustChangePwd }" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="card-header-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" />
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
          </div>
          <div>
            <div class="card-header-title">账号密码</div>
            <div class="card-header-sub">定期更换密码可提升账号安全性</div>
          </div>
        </div>
      </template>

      <el-alert
        v-if="mustChangePwd"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #title>首次登录需要修改密码</template>
        为保障账号安全，请在继续操作前修改默认密码
      </el-alert>

      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="100px">
        <el-form-item label="当前密码" prop="old_password">
          <el-input v-model="pwdForm.old_password" type="password" show-password placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="pwdForm.new_password" type="password" show-password placeholder="至少 8 位，含字母和数字" />
          <div class="pwd-strength" v-if="pwdForm.new_password">
            <div class="strength-bar" :class="strengthClass"></div>
            <span class="strength-text">{{ strengthText }}</span>
          </div>
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm">
          <el-input v-model="pwdForm.confirm" type="password" show-password placeholder="再次输入新密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="doChangePwd">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 外观 -->
    <el-card class="settings-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="card-header-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5" />
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" />
            </svg>
          </div>
          <div>
            <div class="card-header-title">外观</div>
            <div class="card-header-sub">界面显示偏好</div>
          </div>
        </div>
      </template>

      <el-form label-width="100px">
        <el-form-item label="主题模式">
          <el-radio-group v-model="theme" size="large">
            <el-radio-button label="light">
              <span style="display: inline-flex; align-items: center; gap: 6px">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px">
                  <circle cx="12" cy="12" r="4" />
                  <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
                </svg>
                亮色
              </span>
            </el-radio-button>
            <el-radio-button label="dark">
              <span style="display: inline-flex; align-items: center; gap: 6px">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px">
                  <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
                </svg>
                暗色
              </span>
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 关于 -->
    <el-card class="settings-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="card-header-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <path d="M12 16v-4M12 8h.01" />
            </svg>
          </div>
          <div>
            <div class="card-header-title">关于</div>
            <div class="card-header-sub">系统信息</div>
          </div>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="当前用户">{{ userStore.user?.display_name || userStore.user?.username }}（{{ userStore.user?.role === 'admin' ? '管理员' : '工作人员' }}）</el-descriptions-item>
        <el-descriptions-item label="登录到期">{{ formatExp(userStore.expiresAt) }}</el-descriptions-item>
        <el-descriptions-item label="系统名称">石小易 AI 迎新助手 · 工单管理系统</el-descriptions-item>
        <el-descriptions-item label="版本">v1.0.0</el-descriptions-item>
        <el-descriptions-item label="后端">Go 1.22 + Gin + GORM + SQLite</el-descriptions-item>
        <el-descriptions-item label="前端">Vue 3 + Vite + Element Plus</el-descriptions-item>
        <el-descriptions-item label="鉴权">JWT（后台） + API Key（智能体）</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { changePassword } from '../api/auth'
import { useUserStore } from '../stores/user'
import { useTheme } from '../composables/useTheme'

const router = useRouter()
const userStore = useUserStore()
const { theme } = useTheme()

const mustChangePwd = computed(() => userStore.mustChangePassword)

const pwdFormRef = ref(null)
const submitting = ref(false)
const pwdForm = reactive({
  old_password: '',
  new_password: '',
  confirm: '',
})

const validateConfirm = (_, value, callback) => {
  if (value !== pwdForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const pwdRules = {
  old_password: [{ required: true, min: 6, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, min: 8, max: 128, message: '至少 8 位', trigger: 'blur' },
    {
      pattern: /^(?=.*[A-Za-z])(?=.*\d).{8,128}$/,
      message: '必须同时包含字母和数字',
      trigger: 'blur',
    },
  ],
  confirm: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' },
  ],
}

const strength = computed(() => {
  const p = pwdForm.new_password
  if (!p) return 0
  let s = 0
  if (p.length >= 8) s++
  if (p.length >= 12) s++
  if (/[a-z]/.test(p) && /[A-Z]/.test(p)) s++
  if (/\d/.test(p)) s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return Math.min(s, 4)
})

const strengthClass = computed(() => ['', 'weak', 'medium', 'strong', 'very-strong'][strength.value])
const strengthText = computed(() => ['', '弱', '一般', '良好', '很强'][strength.value])

const formatExp = (iso) => {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('zh-CN', { hour12: false })
}

const doChangePwd = async () => {
  await pwdFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await changePassword({
        old_password: pwdForm.old_password,
        new_password: pwdForm.new_password,
      })
      ElMessage.success('密码修改成功')
      pwdForm.old_password = ''
      pwdForm.new_password = ''
      pwdForm.confirm = ''
      // 更新本地 user 状态
      await userStore.refresh()
      // 强制改密码完成后跳到工单列表
      if (mustChangePwd.value) {
        router.push('/tickets')
      }
    } catch (e) {} finally {
      submitting.value = false
    }
  })
}
</script>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
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
  margin: 0 0 16px;
  font-size: 13px;
}

.settings-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}
.settings-card:hover { box-shadow: var(--shadow-md); }
.settings-card.alert {
  border-color: #fdcb6e;
  box-shadow: 0 0 0 3px rgba(253, 203, 110, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 14px;
}
.card-header-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--color-primary-soft);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.card-header-icon svg { width: 20px; height: 20px; }
.card-header-title {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}
.card-header-sub {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.pwd-strength {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}
.strength-bar {
  height: 4px;
  width: 80px;
  border-radius: 2px;
  background: var(--bg-base);
  position: relative;
  overflow: hidden;
}
.strength-bar::after {
  content: '';
  position: absolute;
  inset: 0;
  width: 0;
  transition: all var(--transition-base);
}
.strength-bar.weak::after { width: 25%; background: #d63031; }
.strength-bar.medium::after { width: 50%; background: #fdcb6e; }
.strength-bar.strong::after { width: 75%; background: #00b894; }
.strength-bar.very-strong::after { width: 100%; background: #00b894; }
.strength-text { color: var(--text-tertiary); }
</style>
