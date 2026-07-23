<template>
  <div class="empty-state">
    <component v-if="iconComponent" :is="iconComponent" :size="40" :stroke-width="1.5" class="empty-icon" />
    <div class="empty-title">{{ title }}</div>
    <div v-if="description" class="empty-desc">{{ description }}</div>
    <slot />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import * as LucideIcons from '@lucide/vue'
import { Inbox } from '@lucide/vue'

const props = defineProps({
  icon: { type: String, default: 'Inbox' },
  title: { type: String, default: '暂无数据' },
  description: { type: String, default: '' },
})

const iconComponent = computed(() => LucideIcons[props.icon] || Inbox)
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 16px;
  text-align: center;
  color: var(--text-tertiary);
}
.empty-icon {
  color: var(--text-tertiary);
  opacity: 0.5;
  margin-bottom: 12px;
}
.empty-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 4px;
}
.empty-desc {
  font-size: 13px;
  color: var(--text-tertiary);
}
</style>
