<template>
  <span class="status-badge" :data-status="status">
    <span class="status-dot" />
    <span class="status-text">{{ label }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: String, required: true },
})

const STATUS_MAP = {
  pending: { label: '待处理', color: '#F59E0B' },
  processing: { label: '处理中', color: '#3B82F6' },
  answered: { label: '已答复', color: '#10B981' },
  closed: { label: '已关闭', color: '#6B7280' },
}

const label = computed(() => STATUS_MAP[props.status]?.label || props.status)
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: var(--bg-hover);
  color: var(--text-secondary);
  white-space: nowrap;
  line-height: 1.6;
}
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-tertiary);
  flex-shrink: 0;
}
.status-badge[data-status="pending"] {
  background: var(--color-warning-soft);
  color: var(--color-warning);
}
.status-badge[data-status="pending"] .status-dot { background: var(--color-warning); }
.status-badge[data-status="processing"] {
  background: var(--color-primary-soft);
  color: var(--color-primary);
}
.status-badge[data-status="processing"] .status-dot { background: var(--color-primary); }
.status-badge[data-status="answered"] {
  background: var(--color-success-soft);
  color: var(--color-success);
}
.status-badge[data-status="answered"] .status-dot { background: var(--color-success); }
.status-badge[data-status="closed"] {
  background: var(--color-info-soft);
  color: var(--color-info);
}
.status-badge[data-status="closed"] .status-dot { background: var(--color-info); }
</style>
