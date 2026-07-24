<template>
  <div class="ui-empty anim-fade-up">
    <div class="empty-art">
      <component :is="iconComponent" :size="30" :stroke-width="1.6" />
    </div>
    <div class="empty-title">{{ title }}</div>
    <div v-if="description" class="empty-desc">{{ description }}</div>
    <div v-if="$slots.default" class="empty-action"><slot /></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { resolveIcon } from './icons'
import { Inbox } from '@lucide/vue'

const props = defineProps({
  icon: { type: String, default: 'Inbox' },
  title: { type: String, default: '暂无数据' },
  description: { type: String, default: '' },
})

const iconComponent = computed(() => resolveIcon(props.icon) || Inbox)
</script>

<style scoped>
.ui-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 56px 24px;
  text-align: center;
}
.empty-art {
  width: 68px;
  height: 68px;
  border-radius: var(--r-xl);
  background: var(--primary-soft);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  animation: float-y 3.2s ease-in-out infinite;
}
.empty-title {
  font-size: 14.5px;
  font-weight: 600;
  color: var(--text-1);
  margin-bottom: 4px;
}
.empty-desc { font-size: 13px; color: var(--text-3); max-width: 320px; }
.empty-action { margin-top: 18px; }
</style>
