<template>
  <span v-if="archived" class="archived-badge" :title="`归档于 ${archivedAt || ''}`">
    <Archive :size="12" :stroke-width="2.4" />
    <span>已归档</span>
  </span>
</template>

<script setup>
import { Archive } from '@lucide/vue'

defineProps({
  archived: { type: Boolean, default: false },
  archivedAt: { type: String, default: '' },
})
</script>

<style scoped>
.archived-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11.5px;
  font-weight: 500;
  background: var(--color-info-soft);
  color: var(--color-info);
  white-space: nowrap;
  line-height: 1.7;
}
</style>
