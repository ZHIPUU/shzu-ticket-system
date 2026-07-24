<template>
  <div class="stat-card card" :style="{ '--accent': accent, '--accent-soft': accentSoft }">
    <div class="stat-glow" aria-hidden="true" />
    <div class="stat-icon">
      <component :is="iconComponent" v-if="iconComponent" :size="19" :stroke-width="2" />
    </div>
    <div class="stat-body">
      <div class="stat-label">{{ label }}</div>
      <div class="stat-value"><CountUp :value="value" /></div>
      <div v-if="meta" class="stat-meta">{{ meta }}</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { resolveIcon } from '../ui/icons'
import CountUp from '../ui/CountUp.vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], required: true },
  meta: { type: String, default: '' },
  icon: { type: String, default: '' },
  accent: { type: String, default: '#6366F1' },
})

const iconComponent = computed(() => (props.icon ? resolveIcon(props.icon) : null))

const accentSoft = computed(() => {
  const hex = props.accent.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, 0.12)`
})
</script>

<style scoped>
.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  overflow: hidden;
  transition: transform var(--d-base) var(--ease-out), box-shadow var(--d-base) var(--ease-out);
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.stat-glow {
  position: absolute;
  top: -32px;
  right: -32px;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--accent-soft), transparent 70%);
  pointer-events: none;
}

.stat-icon {
  width: 42px;
  height: 42px;
  border-radius: 13px;
  background: var(--accent-soft);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform var(--d-base) var(--ease-spring);
}
.stat-card:hover .stat-icon { transform: scale(1.08) rotate(-4deg); }

.stat-body { flex: 1; min-width: 0; }
.stat-label { font-size: 12.5px; color: var(--text-2); font-weight: 500; }
.stat-value {
  font-size: 25px;
  font-weight: 700;
  letter-spacing: -0.6px;
  line-height: 1.25;
  color: var(--text-1);
}
.stat-meta { font-size: 11.5px; color: var(--text-3); margin-top: 1px; }

@media (max-width: 480px) {
  .stat-card { padding: 14px 15px; gap: 11px; }
  .stat-icon { width: 36px; height: 36px; border-radius: 11px; }
  .stat-value { font-size: 20px; }
  .stat-meta { display: none; }
}
</style>
