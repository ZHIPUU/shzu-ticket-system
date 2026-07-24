import { reactive } from 'vue'

// 全局确认/输入对话框状态（同一时间只展示一个）
export const confirmState = reactive({
  open: false,
  mode: 'confirm', // confirm | prompt
  title: '',
  message: '',
  confirmText: '确认',
  cancelText: '取消',
  danger: false,
  // prompt 专用
  placeholder: '',
  inputValue: '',
  pattern: null,
  patternMessage: '',
  // 内部
  _resolve: null,
})

const openDialog = (opts) =>
  new Promise((resolve) => {
    Object.assign(confirmState, {
      open: true,
      mode: 'confirm',
      title: '',
      message: '',
      confirmText: '确认',
      cancelText: '取消',
      danger: false,
      placeholder: '',
      inputValue: '',
      pattern: null,
      patternMessage: '',
      ...opts,
      _resolve: resolve,
    })
  })

/** 确认对话框 → Promise<boolean> */
export const confirmDialog = (opts) => openDialog({ ...opts, mode: 'confirm' })

/** 输入对话框 → Promise<string|null>（取消为 null） */
export const promptDialog = (opts) => openDialog({ ...opts, mode: 'prompt' })

export const resolveDialog = (result) => {
  confirmState._resolve?.(result)
  confirmState.open = false
  confirmState._resolve = null
}
