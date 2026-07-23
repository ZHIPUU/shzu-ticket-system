<template>
  <span v-if="category" class="category-badge" :title="category">
    <Tag :size="12" :stroke-width="2.4" />
    <span>{{ category }}</span>
  </span>
  <span v-else class="category-badge empty" :title="'未分类'">
    <Minus :size="12" :stroke-width="2.4" />
    <span>未分类</span>
  </span>
</template>

<script setup>
import { Tag, Minus } from '@lucide/vue'

defineProps({
  category: { type: String, default: '' },
})
</script>

<style scoped>
.category-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11.5px;
  font-weight: 500;
  background: var(--color-primary-soft);
  color: var(--color-primary);
  white-space: nowrap;
  line-height: 1.7;
  border: 1px solid transparent;
}
.category-badge.empty {
  background: transparent;
  color: var(--text-tertiary);
  border-color: var(--border-color);
}
</style>
