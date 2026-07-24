import { reactive } from 'vue'

let seq = 0
export const toasts = reactive([])

const push = (type, message, duration = 2600) => {
  const id = ++seq
  toasts.push({ id, type, message })
  setTimeout(() => {
    const i = toasts.findIndex((t) => t.id === id)
    if (i !== -1) toasts.splice(i, 1)
  }, duration)
}

export const toast = {
  success: (m, d) => push('success', m, d),
  error: (m, d) => push('error', m, d ?? 3200),
  warning: (m, d) => push('warning', m, d),
  info: (m, d) => push('info', m, d),
}
