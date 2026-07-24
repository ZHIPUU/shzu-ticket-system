<template>
  <span class="time-text tnum" :title="absolute">{{ text }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: String, default: '' },
  mode: { type: String, default: 'datetime' }, // datetime | date | relative
})

const pad = (n) => String(n).padStart(2, '0')

const absolute = computed(() => {
  if (!props.value) return ''
  const d = new Date(props.value)
  if (isNaN(d.getTime())) return props.value
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
})

const text = computed(() => {
  if (!props.value) return '—'
  if (props.mode === 'datetime') return absolute.value
  if (props.mode === 'date') return absolute.value.split(' ')[0]
  const d = new Date(props.value)
  const diff = Date.now() - d.getTime()
  if (diff < 60_000) return '刚刚'
  if (diff < 3600_000) return `${Math.floor(diff / 60_000)} 分钟前`
  if (diff < 86_400_000) return `${Math.floor(diff / 3600_000)} 小时前`
  if (diff < 30 * 86_400_000) return `${Math.floor(diff / 86_400_000)} 天前`
  return absolute.value
})
</script>

<style scoped>
.time-text {
  color: var(--text-2);
  font-size: 12.5px;
  white-space: nowrap;
}
</style>
