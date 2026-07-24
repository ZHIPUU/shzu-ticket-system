<template>
  <div class="login-page">
    <!-- 极光背景 -->
    <div class="aurora" aria-hidden="true">
      <div class="blob b1" />
      <div class="blob b2" />
      <div class="blob b3" />
      <div class="grid-mask" />
    </div>

    <div class="login-wrap">
      <div class="login-card anim-fade-up">
        <div class="brand">
          <div class="brand-logo">
            <TicketsPlane :size="26" :stroke-width="2" />
          </div>
          <h1 class="brand-name">石小易工单</h1>
          <p class="brand-sub">迎新智能体 · 工单管理后台</p>
        </div>

        <form class="login-form" @submit.prevent="onSubmit">
          <UiInput
            v-model="form.username"
            size="lg"
            placeholder="用户名"
            clearable
            autofocus
            @enter="onSubmit"
          >
            <template #prefix><User :size="16" :stroke-width="2" /></template>
          </UiInput>

          <UiInput
            v-model="form.password"
            type="password"
            size="lg"
            placeholder="密码"
            @enter="onSubmit"
          >
            <template #prefix><Lock :size="16" :stroke-width="2" /></template>
          </UiInput>

          <p v-if="errorMsg" :key="errorMsg" class="login-error anim-shake">
            <CircleAlert :size="14" :stroke-width="2.2" />
            <span>{{ errorMsg }}</span>
          </p>

          <UiButton variant="primary" size="lg" block :loading="loading" type="submit">
            登 录
          </UiButton>
        </form>

        <div class="login-tips">
          <p>首次部署默认账号 <code>admin / admin123</code>，登录后需修改密码</p>
        </div>
      </div>

      <p class="login-footer anim-fade-up" style="--i: 3">© 2026 石河子大学信息化中心 · 石小易 AI 迎新助手</p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { TicketsPlane, User, Lock, CircleAlert } from '@lucide/vue'
import UiInput from '../ui/UiInput.vue'
import UiButton from '../ui/UiButton.vue'
import { toast } from '../ui/toast'
import { useUserStore } from '../stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: '',
})

const onSubmit = async () => {
  if (loading.value) return
  errorMsg.value = ''
  if (!form.username.trim()) { errorMsg.value = '请输入用户名'; return }
  if (!form.password) { errorMsg.value = '请输入密码'; return }

  loading.value = true
  try {
    const r = await userStore.login(form.username.trim(), form.password)
    toast.success(`欢迎回来，${r.user.display_name || r.user.username}`)
    if (r.must_change_password) {
      router.push('/settings')
    } else {
      router.push(route.query.redirect || '/tickets')
    }
  } catch (e) {
    errorMsg.value = e.response?.data?.error_message || '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
  background: var(--bg-app);
}

/* ── 极光背景 ── */
.aurora { position: absolute; inset: 0; overflow: hidden; pointer-events: none; }
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  will-change: transform;
}
.b1 {
  width: 520px; height: 520px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.55), transparent 65%);
  top: -140px; left: -120px;
  animation: drift-1 22s ease-in-out infinite;
}
.b2 {
  width: 460px; height: 460px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.45), transparent 65%);
  bottom: -120px; right: -100px;
  animation: drift-2 26s ease-in-out infinite;
}
.b3 {
  width: 380px; height: 380px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.35), transparent 65%);
  top: 42%; left: 52%;
  animation: drift-3 30s ease-in-out infinite;
}
html.dark .b1 { background: radial-gradient(circle, rgba(99, 102, 241, 0.35), transparent 65%); }
html.dark .b2 { background: radial-gradient(circle, rgba(139, 92, 246, 0.3), transparent 65%); }
html.dark .b3 { background: radial-gradient(circle, rgba(56, 189, 248, 0.18), transparent 65%); }

@keyframes drift-1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(70px, 50px) scale(1.12); }
}
@keyframes drift-2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-60px, -40px) scale(1.08); }
}
@keyframes drift-3 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-80px, 60px) scale(0.92); }
}

.grid-mask {
  position: absolute;
  inset: 0;
  background-image: radial-gradient(rgba(99, 102, 241, 0.14) 1px, transparent 1px);
  background-size: 26px 26px;
  mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, black, transparent);
  -webkit-mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, black, transparent);
}

/* ── 内容 ── */
.login-wrap {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-card {
  width: 100%;
  padding: 36px 32px 28px;
  border-radius: var(--r-xl);
  background: color-mix(in srgb, var(--bg-surface) 78%, transparent);
  -webkit-backdrop-filter: blur(18px) saturate(1.4);
  backdrop-filter: blur(18px) saturate(1.4);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-lg);
}
@supports not (background: color-mix(in srgb, red 50%, blue)) {
  .login-card { background: var(--bg-surface); }
}

.brand { text-align: center; margin-bottom: 28px; }
.brand-logo {
  width: 58px;
  height: 58px;
  margin: 0 auto 16px;
  border-radius: 17px;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-brand);
  animation: float-y 3.6s ease-in-out infinite;
}
.brand-name {
  margin: 0 0 5px;
  font-size: 23px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--text-1);
}
.brand-sub { margin: 0; font-size: 13px; color: var(--text-3); }

.login-form { display: flex; flex-direction: column; gap: 14px; }

.login-error {
  display: flex;
  align-items: center;
  gap: 7px;
  margin: 0;
  padding: 9px 13px;
  border-radius: var(--r-md);
  background: var(--danger-soft);
  color: var(--danger);
  font-size: 13px;
  font-weight: 500;
}

.login-tips {
  margin-top: 22px;
  padding: 11px 14px;
  border-radius: var(--r-md);
  border: 1px dashed var(--border-strong);
  font-size: 12px;
  color: var(--text-3);
  line-height: 1.7;
  text-align: center;
}
.login-tips p { margin: 0; }
.login-tips code {
  padding: 1px 7px;
  border-radius: 5px;
  background: var(--primary-soft);
  color: var(--primary);
  font-family: "SF Mono", Menlo, Consolas, monospace;
  font-size: 11.5px;
}

.login-footer {
  margin: 22px 0 0;
  font-size: 12px;
  color: var(--text-3);
  text-align: center;
  animation-delay: 260ms;
}

@media (max-width: 480px) {
  .login-card { padding: 30px 22px 24px; }
  .b3 { display: none; }
}
</style>
