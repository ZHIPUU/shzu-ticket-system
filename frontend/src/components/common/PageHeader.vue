<template>
  <header class="page-header">
    <div class="page-header-text">
      <h1 class="page-title">{{ title }}</h1>
      <p v-if="description" class="page-desc">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="page-header-actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 4px;
  flex-wrap: wrap;
}
.page-header-text { min-width: 0; }
.page-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
  letter-spacing: -0.4px;
}
.page-desc {
  color: var(--text-secondary);
  margin: 0;
  font-size: 13px;
}
.page-header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-shrink: 0;
}
@media (max-width: 768px) {
  .page-title { font-size: 18px; }
  .page-desc { font-size: 12px; }
}
</style>
