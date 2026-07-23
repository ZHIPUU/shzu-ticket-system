import { ref, watch, computed } from 'vue'

const STORAGE_KEY = 'ticket_theme'
const themes = ['light', 'dark']

// 读取本地存储，默认亮色
const getInitial = () => {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved && themes.includes(saved)) return saved
  } catch {}
  return 'light'
}

const theme = ref(getInitial())

// 应用到 html 根元素
const applyTheme = (val) => {
  const html = document.documentElement
  if (val === 'dark') {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
  html.setAttribute('data-theme', val)
}

applyTheme(theme.value)
watch(theme, (v) => {
  applyTheme(v)
  try { localStorage.setItem(STORAGE_KEY, v) } catch {}
})

export function useTheme() {
  const toggle = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }
  const isDark = computed(() => theme.value === 'dark')
  return { theme, toggle, isDark }
}
