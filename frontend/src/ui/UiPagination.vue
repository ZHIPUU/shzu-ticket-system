<template>
  <div class="ui-pager" :class="{ compact }">
    <span class="pager-total tnum">共 {{ total }} 条</span>

    <UiSelect
      v-if="!compact && pageSizes.length"
      class="pager-size"
      size="sm"
      :model-value="String(pageSize)"
      :options="pageSizes.map((n) => ({ label: `${n} 条/页`, value: String(n) }))"
      @update:model-value="changeSize"
    />

    <div class="pager-btns">
      <button class="pager-btn" :disabled="page <= 1" aria-label="上一页" @click="go(page - 1)">
        <ChevronLeft :size="15" :stroke-width="2.2" />
      </button>

      <template v-if="!compact">
        <button
          v-for="(p, i) in pageItems"
          :key="`${p}-${i}`"
          class="pager-btn num tnum"
          :class="{ active: p === page, dots: p === '…' }"
          :disabled="p === '…'"
          @click="go(p)"
        >{{ p }}</button>
      </template>
      <span v-else class="pager-info tnum">{{ page }} / {{ pageCount }}</span>

      <button class="pager-btn" :disabled="page >= pageCount" aria-label="下一页" @click="go(page + 1)">
        <ChevronRight :size="15" :stroke-width="2.2" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from '@lucide/vue'
import UiSelect from './UiSelect.vue'

const props = defineProps({
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 20 },
  total: { type: Number, default: 0 },
  pageSizes: { type: Array, default: () => [10, 20, 50, 100] },
  compact: { type: Boolean, default: false },
})
const emit = defineEmits(['update:page', 'update:pageSize', 'change'])

const pageCount = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const pageItems = computed(() => {
  const n = pageCount.value
  const c = props.page
  if (n <= 7) return Array.from({ length: n }, (_, i) => i + 1)
  const items = [1]
  const lo = Math.max(2, c - 1)
  const hi = Math.min(n - 1, c + 1)
  if (lo > 2) items.push('…')
  for (let i = lo; i <= hi; i++) items.push(i)
  if (hi < n - 1) items.push('…')
  items.push(n)
  return items
})

const go = (p) => {
  if (p === '…' || p < 1 || p > pageCount.value || p === props.page) return
  emit('update:page', p)
  emit('change')
}
const changeSize = (v) => {
  emit('update:pageSize', Number(v))
  emit('update:page', 1)
  emit('change')
}
</script>

<style scoped>
.ui-pager {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}
.ui-pager.compact { justify-content: center; gap: 10px; }

.pager-total { font-size: 12.5px; color: var(--text-3); white-space: nowrap; }
.pager-size { width: 104px; }

.pager-btns { display: inline-flex; align-items: center; gap: 5px; }

.pager-btn {
  min-width: 30px;
  height: 30px;
  padding: 0 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border);
  border-radius: var(--r-sm);
  background: var(--bg-surface);
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
}
.pager-btn:hover:not(:disabled):not(.active) {
  border-color: var(--primary);
  color: var(--primary);
}
.pager-btn.active {
  background: var(--gradient-brand);
  border-color: transparent;
  color: #fff;
  font-weight: 600;
  box-shadow: var(--shadow-brand);
}
.pager-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.pager-btn.dots { border: none; background: transparent; cursor: default; }

.pager-info { font-size: 13px; color: var(--text-2); padding: 0 4px; }
</style>
