<template>
  <div ref="root" class="ui-seg" :class="[`s-${size}`]" role="tablist">
    <span v-if="pill.show" class="seg-pill" :style="pillStyle" aria-hidden="true" />
    <button
      v-for="opt in options"
      :key="String(opt.value)"
      ref="btns"
      type="button"
      class="seg-item"
      :class="{ active: opt.value === modelValue }"
      role="tab"
      :aria-selected="opt.value === modelValue"
      @click="pick(opt)"
    >
      <span class="seg-label">{{ opt.label }}</span>
      <span v-if="opt.count !== undefined && opt.count !== null" class="seg-count tnum">{{ opt.count }}</span>
    </button>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  options: { type: Array, default: () => [] }, // [{ label, value, count? }]
  size: { type: String, default: 'md' }, // sm | md
})
const emit = defineEmits(['update:modelValue', 'change'])

const root = ref(null)
const btns = ref([])
const pill = reactive({ show: false, left: 0, width: 0 })

const pillStyle = computed(() => ({
  transform: `translateX(${pill.left}px)`,
  width: `${pill.width}px`,
}))

const measure = () => {
  const idx = props.options.findIndex((o) => o.value === props.modelValue)
  const el = idx >= 0 ? btns.value?.[idx] : null
  if (!el) { pill.show = false; return }
  pill.left = el.offsetLeft
  pill.width = el.offsetWidth
  pill.show = true
}

const pick = (opt) => {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
}

watch(() => [props.modelValue, props.options], () => nextTick(measure), { deep: true })
onMounted(() => {
  nextTick(measure)
  window.addEventListener('resize', measure)
})
onBeforeUnmount(() => window.removeEventListener('resize', measure))
</script>

<style scoped>
.ui-seg {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 3px;
  background: var(--bg-sink);
  border: 1px solid var(--border-soft);
  border-radius: var(--r-full);
  max-width: 100%;
  overflow-x: auto;
  scrollbar-width: none;
}
.ui-seg::-webkit-scrollbar { display: none; }

.seg-pill {
  position: absolute;
  top: 3px;
  bottom: 3px;
  left: 0;
  border-radius: var(--r-full);
  background: var(--bg-surface);
  box-shadow: var(--shadow-sm);
  transition: transform var(--d-base) var(--ease-spring), width var(--d-base) var(--ease-spring);
  pointer-events: none;
}
html.dark .seg-pill { background: var(--bg-elevated); }

.seg-item {
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: none;
  background: transparent;
  border-radius: var(--r-full);
  color: var(--text-2);
  cursor: pointer;
  white-space: nowrap;
  transition: color var(--d-fast) var(--ease-out);
  line-height: 1;
}
.s-md .seg-item { height: 30px; padding: 0 14px; font-size: 13px; }
.s-sm .seg-item { height: 26px; padding: 0 11px; font-size: 12.5px; }

.seg-item:hover:not(.active) { color: var(--text-1); }
.seg-item.active { color: var(--text-1); font-weight: 500; }

.seg-count {
  font-size: 11px;
  padding: 1px 7px;
  border-radius: var(--r-full);
  background: var(--bg-active);
  color: var(--text-2);
}
.seg-item.active .seg-count {
  background: var(--primary-soft);
  color: var(--primary);
}
</style>
