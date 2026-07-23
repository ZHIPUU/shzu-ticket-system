import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, getMe } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('ticket_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('ticket_user') || 'null'))
  const expiresAt = ref(localStorage.getItem('ticket_token_exp') || '')

  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const mustChangePassword = computed(() => user.value?.must_change_password === true)

  function setAuth({ token: t, user: u, expires_at }) {
    token.value = t
    user.value = u
    expiresAt.value = expires_at
    localStorage.setItem('ticket_token', t)
    localStorage.setItem('ticket_user', JSON.stringify(u))
    localStorage.setItem('ticket_token_exp', expires_at)
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    expiresAt.value = ''
    localStorage.removeItem('ticket_token')
    localStorage.removeItem('ticket_user')
    localStorage.removeItem('ticket_token_exp')
    // 清理老的 API Key 存储
    localStorage.removeItem('ticket_api_key')
  }

  async function login(username, password) {
    const r = await apiLogin({ username, password })
    setAuth(r)
    return r
  }

  async function refresh() {
    if (!token.value) return null
    try {
      const r = await getMe()
      user.value = r.user
      localStorage.setItem('ticket_user', JSON.stringify(r.user))
      return r.user
    } catch (e) {
      clearAuth()
      return null
    }
  }

  function logout() {
    clearAuth()
  }

  return {
    token,
    user,
    expiresAt,
    isLoggedIn,
    isAdmin,
    mustChangePassword,
    setAuth,
    clearAuth,
    login,
    refresh,
    logout,
  }
})
