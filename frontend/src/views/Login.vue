<template>
  <div class="login-page">
    <div class="login-bg" :class="{ dark: isDark }">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <div class="login-container">
      <div class="brand-area">
        <div class="logo">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 7l9-4 9 4-9 4-9-4z" />
            <path d="M3 12l9 4 9-4" />
            <path d="M3 17l9 4 9-4" />
          </svg>
        </div>
        <h1 class="brand-name">石小易工单</h1>
        <p class="brand-sub">迎新智能体后台管理系统</p>
      </div>

      <el-card class="login-card" shadow="never">
        <h2 class="form-title">登录</h2>
        <p class="form-sub">使用您的账号继续</p>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="onSubmit" size="large">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" clearable autofocus />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="密码" :prefix-icon="Lock" show-password
              @keyup.enter="onSubmit" />
          </el-form-item>
          <el-form-item v-if="errorMsg">
            <el-alert :title="errorMsg" type="error" :closable="false" show-icon />
          </el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="onSubmit" class="submit-btn">
            登录
          </el-button>
        </el-form>

        <div class="login-tips">
          <p>首次部署默认账号：<code>admin / admin123</code></p>
          <p>登录后会强制要求修改密码</p>
        </div>
      </el-card>

      <div class="footer">
        © 2026 石河子大学信息化中心 · 石小易 AI 迎新助手
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import { useTheme } from '../composables/useTheme'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { isDark } = useTheme()

const formRef = ref(null)
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const onSubmit = async () => {
  errorMsg.value = ''
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const r = await userStore.login(form.username, form.password)
      ElMessage.success(`欢迎回来，${r.user.display_name || r.user.username}`)
      const redirect = route.query.redirect || '/tickets'
      // 强制改密码时先跳设置
      if (r.must_change_password) {
        router.push('/settings')
      } else {
        router.push(redirect)
      }
    } catch (e) {
      const detail = e.response?.data?.error_message || '登录失败'
      errorMsg.value = detail
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f0fbf8 0%, #e6f9f5 50%, #f7f9fc 100%);
  overflow: hidden;
}
html.dark .login-page {
  background: linear-gradient(135deg, #0a1413 0%, #0f1115 50%, #142826 100%);
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: float 20s ease-in-out infinite;
}
.blob-1 { width: 480px; height: 480px; background: #00b894; top: -100px; left: -100px; }
.blob-2 { width: 380px; height: 380px; background: #00cec9; bottom: -80px; right: -80px; animation-delay: -7s; }
.blob-3 { width: 320px; height: 320px; background: #74b9ff; top: 40%; left: 50%; animation-delay: -14s; opacity: 0.3; }
@media (max-width: 768px) {
  .blob-1 { width: 260px; height: 260px; top: -60px; left: -60px; }
  .blob-2 { width: 200px; height: 200px; bottom: -40px; right: -40px; }
  .blob-3 { display: none; }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(60px, -40px) scale(1.1); }
}

.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.brand-area {
  text-align: center;
  margin-bottom: 32px;
  color: var(--text-primary);
}
.logo {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  background: var(--gradient-header);
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: 0 8px 24px rgba(0, 184, 148, 0.3);
}
.logo svg { width: 32px; height: 32px; }
.brand-name {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 4px;
  letter-spacing: -0.5px;
}
.brand-sub {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0;
}

.login-card {
  width: 100%;
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
  padding: 32px;
}
@media (max-width: 480px) { .login-card { padding: 24px 20px; } }
html.dark .login-card {
  background: var(--bg-surface);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.form-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
  letter-spacing: -0.3px;
}
.form-sub {
  color: var(--text-secondary);
  font-size: 13px;
  margin: 0 0 24px;
}

.submit-btn {
  width: 100%;
  background: var(--gradient-header);
  border: none;
  height: 44px;
  font-size: 15px;
  font-weight: 500;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 12px rgba(0, 184, 148, 0.25);
  transition: all var(--transition-base);
}
.submit-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(0, 184, 148, 0.35);
  background: var(--gradient-header) !important;
}

.login-tips {
  margin-top: 20px;
  padding: 12px 14px;
  background: var(--bg-base);
  border-radius: 8px;
  border: 1px dashed var(--border-color);
  font-size: 12px;
  color: var(--text-tertiary);
  line-height: 1.7;
}
.login-tips code {
  background: var(--bg-elevated);
  padding: 1px 6px;
  border-radius: 3px;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  color: var(--color-primary);
  font-size: 11.5px;
}

.footer {
  margin-top: 24px;
  color: var(--text-tertiary);
  font-size: 12px;
  text-align: center;
}
</style>
