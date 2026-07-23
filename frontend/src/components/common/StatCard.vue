<template>
  <div class="stat-card" :style="cardStyle">
    <div class="stat-icon" :style="{ background: tintedBg, color: accent }">
      <component :is="iconComponent" v-if="iconComponent" :size="18" :stroke-width="2" />
    </div>
    <div class="stat-content">
      <div class="stat-label">{{ label }}</div>
      <div class="stat-value">{{ value }}</div>
      <div v-if="meta" class="stat-meta">{{ meta }}</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import * as LucideIcons from '@lucide/vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], required: true },
  meta: { type: String, default: '' },
  icon: { type: String, default: '' },
  accent: { type: String, default: '#2563EB' },
  accentAlpha: { type: Number, default: 0.10 },
})

const iconComponent = computed(() => {
  if (!props.icon) return null
  return LucideIcons[props.icon] || null
})

const cardStyle = computed(() => ({
  '--accent': props.accent,
}))

const tintedBg = computed(() => {
  // 简单把 hex 转 rgba
  const hex = props.accent.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, ${props.accentAlpha})`
})
</script>

<style scoped>
.stat-card {
  background: var(--bg-surface);
  border-radius: var(--radius-lg);
  padding: 18px 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-soft);
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.stat-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--accent);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-icon {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content { flex: 1; min-width: 0; }

.stat-label {
  color: var(--text-secondary);
  font-size: 13px;
  margin-bottom: 4px;
  font-weight: 500;
}

.stat-value {
  color: var(--text-primary);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.5px;
  line-height: 1.2;
  font-variant-numeric: tabular-nums;
}

.stat-meta {
  color: var(--text-tertiary);
  font-size: 12px;
  margin-top: 4px;
}

@media (max-width: 480px) {
  .stat-card { padding: 14px 16px; gap: 10px; }
  .stat-icon { width: 32px; height: 32px; }
  .stat-value { font-size: 20px; }
  .stat-label { font-size: 12px; }
  .stat-meta { font-size: 11px; }
}
</style>
