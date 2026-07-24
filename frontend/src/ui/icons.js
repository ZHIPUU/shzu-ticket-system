// 按名字解析的图标注册表（避免全量引入 lucide 图标，控制包体）
import {
  Layers, Timer, MessageSquareCheck, Database,
  Inbox, MessageSquare, Users,
} from '@lucide/vue'

const registry = {
  Layers, Timer, MessageSquareCheck, Database,
  Inbox, MessageSquare, Users,
}

export const resolveIcon = (name) => registry[name] || null
