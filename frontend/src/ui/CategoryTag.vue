<template>
  <span class="cat-tag" :class="{ empty: !category }" :title="category || '未分类'">
    <Tag v-if="category" :size="11" :stroke-width="2.4" />
    <span>{{ category || '未分类' }}</span>
  </span>
</template>

<script setup>
import { Tag } from '@lucide/vue'

defineProps({
  category: { type: String, default: '' },
})
</script>

<style scoped>
.cat-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 9px;
  border-radius: var(--r-sm);
  font-size: 11.5px;
  font-weight: 500;
  white-space: nowrap;
  line-height: 1.65;
  background: var(--primary-soft);
  color: var(--primary);
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.cat-tag.empty {
  background: transparent;
  color: var(--text-3);
  border: 1px dashed var(--border-strong);
}
</style>
