<template>
  <span class="status-pill" :data-status="status">
    <span class="pill-dot" />
    <span>{{ label }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: String, required: true },
})

const MAP = {
  pending: '待处理',
  processing: '处理中',
  answered: '已答复',
  closed: '已关闭',
}
const label = computed(() => MAP[props.status] || props.status)
</script>

<style scoped>
.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  border-radius: var(--r-full);
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
  line-height: 1.55;
  background: var(--bg-hover);
  color: var(--text-2);
}
.pill-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  flex-shrink: 0;
}

.status-pill[data-status="pending"] { background: var(--warning-soft); color: var(--warning); }
.status-pill[data-status="pending"] .pill-dot { animation: pulse-dot 1.8s ease-out infinite; }
.status-pill[data-status="processing"] { background: var(--primary-soft); color: var(--primary); }
.status-pill[data-status="answered"] { background: var(--success-soft); color: var(--success); }
.status-pill[data-status="closed"] { background: var(--info-soft); color: var(--info); }
</style>
