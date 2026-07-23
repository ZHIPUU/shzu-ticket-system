import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 10000,
})

// 注入 JWT Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('ticket_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (resp) => resp,
  (err) => {
    // 401 时清 token，让路由守卫跳登录
    if (err.response?.status === 401) {
      localStorage.removeItem('ticket_token')
      localStorage.removeItem('ticket_user')
      localStorage.removeItem('ticket_token_exp')
    }
    return Promise.reject(err)
  }
)

export const login = (data) => api.post('/auth/login', data).then((r) => r.data)
export const getMe = () => api.get('/auth/me').then((r) => r.data)
export const changePassword = (data) => api.post('/auth/change-password', data).then((r) => r.data)

// 用户管理（admin）
export const listUsers = () => api.get('/users').then((r) => r.data)
export const createUser = (data) => api.post('/users', data).then((r) => r.data)
export const updateUser = (id, data) => api.patch(`/users/${id}`, data).then((r) => r.data)
