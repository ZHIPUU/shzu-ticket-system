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

// ─── 列表 / 详情 ───
export const listTickets = (params) => api.get('/tickets', { params }).then((r) => r.data)
export const getTicket = (id) => api.get(`/tickets/${id}`).then((r) => r.data)

// ─── 提交（HiAgent 智能体调用） ───
export const submitTicket = (data) => api.post('/tickets', data).then((r) => r.data)

// ─── 答复 / 关闭 / 重开 ───
export const answerTicket = (id, data) => api.post(`/tickets/${id}/answer`, data).then((r) => r.data)
export const closeTicket = (id, reason) => api.post(`/tickets/${id}/close`, { reason }).then((r) => r.data)
export const reopenTicket = (id) => api.post(`/tickets/${id}/reopen`, {}).then((r) => r.data)

// ─── 编辑（分类 / 备注）───
export const patchTicket = (id, data) => api.patch(`/tickets/${id}`, data).then((r) => r.data)

// ─── 归档（archive + unarchive 走同一接口，靠 body 区分） ───
export const archiveTicket = (id, archive = true) =>
  api.post(`/tickets/${id}/archive`, { archive }).then((r) => r.data)
export const unarchiveTicket = (id) => archiveTicket(id, false)

// ─── 删除（软删 / 硬删）───
export const deleteTicket = (id, hard = false) =>
  api.delete(`/tickets/${id}`, { params: hard ? { hard: true } : {} }).then((r) => r.data)

// ─── 批量删除 ───
export const batchDeleteTickets = (ticket_ids) =>
  api.post('/tickets/batch-delete', { ticket_ids }).then((r) => r.data)

// ─── 导出（按当前筛选条件） ───
export const exportTicketsUrl = (params) => {
  const base = import.meta.env.VITE_API_BASE || '/api/v1'
  const token = localStorage.getItem('ticket_token') || ''
  const search = new URLSearchParams()
  Object.entries(params || {}).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') search.append(k, v)
  })
  return {
    url: `${base}/tickets/export?${search.toString()}`,
    token,
  }
}

export const exportTickets = async (params) => {
  const { url, token } = exportTicketsUrl(params)
  const r = await axios.get(url, {
    responseType: 'blob',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
  return r
}
