<template>
  <div class="settings-page">
    <header class="page-head anim-fade-up">
      <div class="head-text">
        <h1>个人设置</h1>
        <p>管理账号密码与界面偏好</p>
      </div>
    </header>

    <!-- 强制改密码提示 -->
    <div v-if="mustChangePwd" class="pwd-banner anim-fade-up">
      <ShieldAlert :size="18" :stroke-width="2" />
      <div class="banner-text">
        <strong>首次登录需要修改密码</strong>
        <span>为保障账号安全，请在继续操作前修改默认密码</span>
      </div>
    </div>

    <!-- ═══ 账号密码 ═══ -->
    <section class="card block anim-stagger" style="--i: 1">
      <header class="block-head">
        <span class="block-icon"><Lock :size="16" :stroke-width="2" /></span>
        <h3>账号密码</h3>
      </header>
      <div class="block-body form-stack">
        <div class="form-item">
          <label class="form-label">当前密码</label>
          <UiInput v-model="pwdForm.old_password" type="password" placeholder="请输入当前密码" />
        </div>
        <div class="form-item">
          <label class="form-label">新密码</label>
          <UiInput v-model="pwdForm.new_password" type="password" placeholder="至少 8 位，含字母和数字" />
          <div v-if="pwdForm.new_password" class="strength">
            <div class="strength-track">
              <div class="strength-fill" :class="strengthClass" :style="{ width: strength * 25 + '%' }" />
            </div>
            <span class="strength-text" :class="strengthClass">{{ strengthText }}</span>
          </div>
        </div>
        <div class="form-item">
          <label class="form-label">确认新密码</label>
          <UiInput
            v-model="pwdForm.confirm"
            type="password"
            placeholder="再次输入新密码"
            @enter="doChangePwd"
          />
          <p v-if="pwdForm.confirm && pwdForm.confirm !== pwdForm.new_password" class="field-error">两次输入的密码不一致</p>
        </div>
        <div class="form-submit">
          <UiButton variant="primary" :loading="submitting" @click="doChangePwd">
            <template #icon><Save :size="14" :stroke-width="2" /></template>
            修改密码
          </UiButton>
        </div>
      </div>
    </section>

    <!-- ═══ 外观 ═══ -->
    <section class="card block anim-stagger" style="--i: 2">
      <header class="block-head">
        <span class="block-icon"><Sun :size="16" :stroke-width="2" /></span>
        <h3>外观主题</h3>
      </header>
      <div class="block-body">
        <div class="theme-grid">
          <button
            v-for="t in themeTiles"
            :key="t.value"
            class="theme-tile"
            :class="{ active: theme === t.value }"
            @click="theme = t.value"
          >
            <span class="tile-preview" :data-theme="t.value">
              <span class="pv-side" />
              <span class="pv-main">
                <span class="pv-bar" />
                <span class="pv-line w80" />
                <span class="pv-line w60" />
              </span>
            </span>
            <span class="tile-label">
              <component :is="t.icon" :size="14" :stroke-width="2" />
              {{ t.label }}
            </span>
            <span class="tile-check"><Check :size="12" :stroke-width="3" /></span>
          </button>
        </div>
      </div>
    </section>

    <!-- ═══ 关于 ═══ -->
    <section class="card block anim-stagger" style="--i: 3">
      <header class="block-head">
        <span class="block-icon"><Info :size="16" :stroke-width="2" /></span>
        <h3>关于系统</h3>
      </header>
      <div class="block-body info-list">
        <div class="info-row">
          <span class="info-label">当前用户</span>
          <span class="info-value">
            {{ userStore.user?.display_name || userStore.user?.username }}
            <span class="role-pill" :data-role="userStore.user?.role">{{ userStore.user?.role === 'admin' ? '管理员' : '工作人员' }}</span>
          </span>
        </div>
        <div class="info-row">
          <span class="info-label">登录到期</span>
          <span class="info-value tnum">{{ formatExp(userStore.expiresAt) }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">系统名称</span>
          <span class="info-value">石小易 AI 迎新助手 · 工单管理系统</span>
        </div>
        <div class="info-row">
          <span class="info-label">版本</span>
          <span class="info-value mono">v3.0.0</span>
        </div>
        <div class="info-row">
          <span class="info-label">后端</span>
          <span class="info-value">Go 1.22 + Gin + GORM + SQLite</span>
        </div>
        <div class="info-row">
          <span class="info-label">前端</span>
          <span class="info-value">Vue 3 + Vite + Lucide</span>
        </div>
        <div class="info-row">
          <span class="info-label">鉴权</span>
          <span class="info-value">JWT（后台） + API Key（智能体）</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Lock, Sun, Moon, Check, Info, Save, ShieldAlert } from '@lucide/vue'
import { changePassword } from '../api/auth'
import { useUserStore } from '../stores/user'
import { useTheme } from '../composables/useTheme'
import UiButton from '../ui/UiButton.vue'
import UiInput from '../ui/UiInput.vue'
import { toast } from '../ui/toast'

const router = useRouter()
const userStore = useUserStore()
const { theme } = useTheme()

const mustChangePwd = computed(() => userStore.mustChangePassword)

const submitting = ref(false)
const pwdForm = reactive({ old_password: '', new_password: '', confirm: '' })

const themeTiles = [
  { label: '亮色模式', value: 'light', icon: Sun },
  { label: '暗色模式', value: 'dark', icon: Moon },
]

// ─── 密码强度 ───
const strength = computed(() => {
  const p = pwdForm.new_password
  if (!p) return 0
  let s = 0
  if (p.length >= 8) s++
  if (p.length >= 12) s++
  if (/[a-z]/.test(p) && /[A-Z]/.test(p)) s++
  if (/\d/.test(p)) s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return Math.max(1, Math.min(s, 4))
})
const strengthClass = computed(() => ['', 'weak', 'medium', 'strong', 'very-strong'][strength.value])
const strengthText = computed(() => ['', '弱', '一般', '良好', '很强'][strength.value])

const formatExp = (iso) => {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('zh-CN', { hour12: false })
}

const doChangePwd = async () => {
  if (submitting.value) return
  if (!pwdForm.old_password) { toast.warning('请输入当前密码'); return }
  if (!/^(?=.*[A-Za-z])(?=.*\d).{8,128}$/.test(pwdForm.new_password)) {
    toast.warning('新密码至少 8 位，且必须包含字母和数字'); return
  }
  if (pwdForm.confirm !== pwdForm.new_password) {
    toast.warning('两次输入的密码不一致'); return
  }
  submitting.value = true
  try {
    await changePassword({ old_password: pwdForm.old_password, new_password: pwdForm.new_password })
    toast.success('密码修改成功')
    pwdForm.old_password = ''
    pwdForm.new_password = ''
    pwdForm.confirm = ''
    await userStore.refresh()
    if (mustChangePwd.value) router.push('/tickets')
  } catch (e) {} finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
  max-width: 720px;
}

.page-head h1 {
  margin: 0 0 3px;
  font-size: 23px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--text-1);
}
.head-text p { margin: 0; font-size: 13px; color: var(--text-3); }
@media (max-width: 767px) { .page-head h1 { font-size: 19px; } }

/* ─── 强制改密横幅 ─── */
.pwd-banner {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 14px 18px;
  border-radius: var(--r-lg);
  background: var(--warning-soft);
  border: 1px solid color-mix(in srgb, var(--warning) 30%, transparent);
  color: var(--warning);
}
.banner-text { display: flex; flex-direction: column; gap: 1px; }
.banner-text strong { font-size: 13.5px; color: var(--text-1); }
.banner-text span { font-size: 12.5px; color: var(--text-2); }

/* ─── 通用块 ─── */
.block { overflow: hidden; }
.block-head {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border-soft);
}
.block-head h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}
.block-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: var(--primary-soft);
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.block-body { padding: 20px; }

/* ─── 表单 ─── */
.form-stack { display: flex; flex-direction: column; gap: 15px; }
.form-item { display: flex; flex-direction: column; gap: 7px; }
.form-label { font-size: 13px; font-weight: 500; color: var(--text-2); }
.field-error { margin: 0; font-size: 12px; color: var(--danger); }
.form-submit { display: flex; justify-content: flex-end; }

/* ─── 强度条 ─── */
.strength { display: flex; align-items: center; gap: 10px; }
.strength-track {
  flex: 1;
  max-width: 180px;
  height: 5px;
  border-radius: 3px;
  background: var(--bg-active);
  overflow: hidden;
}
.strength-fill {
  height: 100%;
  border-radius: 3px;
  transition: width var(--d-base) var(--ease-out), background var(--d-base);
}
.strength-fill.weak { background: var(--danger); }
.strength-fill.medium { background: var(--warning); }
.strength-fill.strong { background: #84CC16; }
.strength-fill.very-strong { background: var(--success); }
.strength-text { font-size: 12px; color: var(--text-3); }
.strength-text.weak { color: var(--danger); }
.strength-text.medium { color: var(--warning); }
.strength-text.strong, .strength-text.very-strong { color: var(--success); }

/* ─── 主题瓷贴 ─── */
.theme-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}
@media (max-width: 480px) { .theme-grid { grid-template-columns: 1fr; } }

.theme-tile {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  border: 1.5px solid var(--border);
  border-radius: var(--r-lg);
  background: var(--bg-surface);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
}
.theme-tile:hover { border-color: var(--border-strong); transform: translateY(-1px); }
.theme-tile.active {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-soft-2);
}

.tile-preview {
  display: flex;
  height: 74px;
  border-radius: var(--r-md);
  overflow: hidden;
  border: 1px solid var(--border-soft);
}
.tile-preview[data-theme="light"] { background: #F3F4F9; }
.tile-preview[data-theme="dark"] { background: #0B0D15; }

.pv-side { width: 26%; }
.tile-preview[data-theme="light"] .pv-side { background: #FFFFFF; border-right: 1px solid #E8E9F0; }
.tile-preview[data-theme="dark"] .pv-side { background: #131624; border-right: 1px solid #232741; }

.pv-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 7px;
  padding: 10px;
}
.pv-bar {
  height: 8px;
  width: 55%;
  border-radius: 4px;
  background: var(--gradient-brand);
  opacity: 0.85;
}
.pv-line { height: 6px; border-radius: 3px; }
.pv-line.w80 { width: 80%; }
.pv-line.w60 { width: 60%; }
.tile-preview[data-theme="light"] .pv-line { background: #DFE1EC; }
.tile-preview[data-theme="dark"] .pv-line { background: #262B44; }

.tile-label {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-2);
}
.theme-tile.active .tile-label { color: var(--primary); font-weight: 600; }

.tile-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: scale(0.5);
  transition: all var(--d-base) var(--ease-spring);
}
.theme-tile.active .tile-check { opacity: 1; transform: scale(1); }

/* ─── 信息列表 ─── */
.info-list { display: flex; flex-direction: column; }
.info-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-soft);
}
.info-row:last-child { border-bottom: none; }
.info-label { width: 72px; flex-shrink: 0; font-size: 12.5px; color: var(--text-3); }
.info-value {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-1);
  min-width: 0;
  flex-wrap: wrap;
}

.role-pill {
  display: inline-flex;
  padding: 2px 9px;
  border-radius: var(--r-full);
  font-size: 11px;
  font-weight: 600;
  background: var(--info-soft);
  color: var(--info);
}
.role-pill[data-role="admin"] { background: var(--warning-soft); color: var(--warning); }
</style>
