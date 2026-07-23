import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 10000,
})

// 注入 JWT
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
    const status = err.response?.status
    const detail = err.response?.data?.detail || err.response?.data?.error_message || err.message
    if (status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      localStorage.removeItem('ticket_token')
      localStorage.removeItem('ticket_user')
      // 路由守卫会处理跳转
    } else if (status === 403) {
      ElMessage.error('权限不足：' + detail)
    } else if (status === 404) {
      ElMessage.warning('工单不存在')
    } else if (status === 422) {
      ElMessage.warning('参数错误：' + detail)
    } else if (status >= 500) {
      ElMessage.error('服务器内部错误：' + detail)
    } else if (err.code === 'ERR_NETWORK') {
      ElMessage.error('网络异常，无法连接后端服务')
    }
    return Promise.reject(err)
  }
)

export const listTickets = (params) => api.get('/tickets', { params }).then((r) => r.data)
export const getTicket = (id) => api.get(`/tickets/${id}`).then((r) => r.data)
export const answerTicket = (id, data) => api.post(`/tickets/${id}/answer`, data).then((r) => r.data)
export const closeTicket = (id, reason) => api.post(`/tickets/${id}/close`, { reason }).then((r) => r.data)
export const submitTicket = (data) => api.post('/tickets', data).then((r) => r.data)
