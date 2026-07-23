<template>
  <div class="section-card" :class="{ padded }">
    <div v-if="title || $slots.header" class="section-header">
      <div class="section-header-left">
        <component v-if="iconComponent" :is="iconComponent" :size="16" :stroke-width="2" class="section-icon" />
        <span class="section-title">{{ title }}</span>
      </div>
      <div v-if="$slots.header" class="section-header-right">
        <slot name="header" />
      </div>
    </div>
    <div class="section-body" :class="{ 'no-header': !title && !$slots.header }">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import * as LucideIcons from '@lucide/vue'

const props = defineProps({
  title: { type: String, default: '' },
  icon: { type: String, default: '' },
  padded: { type: Boolean, default: true },
})

const iconComponent = computed(() => props.icon ? (LucideIcons[props.icon] || null) : null)
</script>

<style scoped>
.section-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border-soft);
  gap: 12px;
}
.section-header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.section-icon {
  color: var(--color-primary);
  flex-shrink: 0;
}
.section-title {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
  letter-spacing: -0.2px;
}
.section-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.section-body { padding: 20px; }
.section-body.no-header { padding: 0; }
</style>
