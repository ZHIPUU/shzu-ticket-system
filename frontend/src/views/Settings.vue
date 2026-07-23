<template>
  <div class="settings">
    <PageHeader title="个人设置" description="管理账号密码与界面偏好" />

    <!-- 修改密码 -->
    <SectionCard
      title="账号密码"
      icon="Lock"
      :padded="false"
    >
      <template #header>
        <el-alert
          v-if="mustChangePwd"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 4px"
        >
          <template #title>首次登录需要修改密码</template>
        </el-alert>
      </template>
      <div class="card-pad">
        <div v-if="mustChangePwd" class="pwd-hint">
          为保障账号安全，请在继续操作前修改默认密码
        </div>
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
      </div>
    </SectionCard>

    <!-- 外观 -->
    <SectionCard title="外观" icon="Sun" :padded="false">
      <div class="card-pad">
        <el-form label-width="100px">
          <el-form-item label="主题模式">
            <el-radio-group v-model="theme" size="large">
              <el-radio-button label="light">
                <span class="radio-content">
                  <Sun :size="16" :stroke-width="2" />
                  亮色
                </span>
              </el-radio-button>
              <el-radio-button label="dark">
                <span class="radio-content">
                  <Moon :size="16" :stroke-width="2" />
                  暗色
                </span>
              </el-radio-button>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
    </SectionCard>

    <!-- 关于 -->
    <SectionCard title="关于" icon="Info" :padded="false">
      <div class="card-pad">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="当前用户">{{ userStore.user?.display_name || userStore.user?.username }}（{{ userStore.user?.role === 'admin' ? '管理员' : '工作人员' }}）</el-descriptions-item>
          <el-descriptions-item label="登录到期">{{ formatExp(userStore.expiresAt) }}</el-descriptions-item>
          <el-descriptions-item label="系统名称">石小易 AI 迎新助手 · 工单管理系统</el-descriptions-item>
          <el-descriptions-item label="版本">v3.0.0</el-descriptions-item>
          <el-descriptions-item label="后端">Go 1.22 + Gin + GORM + SQLite</el-descriptions-item>
          <el-descriptions-item label="前端">Vue 3 + Vite + Element Plus + Lucide</el-descriptions-item>
          <el-descriptions-item label="鉴权">JWT（后台） + API Key（智能体）</el-descriptions-item>
        </el-descriptions>
      </div>
    </SectionCard>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Sun, Moon } from '@lucide/vue'
import { changePassword } from '../api/auth'
import { useUserStore } from '../stores/user'
import { useTheme } from '../composables/useTheme'
import PageHeader from '../components/common/PageHeader.vue'
import SectionCard from '../components/common/SectionCard.vue'

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
      await userStore.refresh()
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
  gap: 16px;
  max-width: 800px;
}

.card-pad {
  padding: 20px;
}

.pwd-hint {
  margin-bottom: 16px;
  color: var(--text-secondary);
  font-size: 13px;
}

.radio-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
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
.strength-bar.weak::after { width: 25%; background: var(--color-danger); }
.strength-bar.medium::after { width: 50%; background: var(--color-warning); }
.strength-bar.strong::after { width: 75%; background: var(--color-success); }
.strength-bar.very-strong::after { width: 100%; background: var(--color-success); }
.strength-text { color: var(--text-tertiary); }
</style>
