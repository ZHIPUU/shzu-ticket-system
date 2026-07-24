<template>
  <div ref="root" class="ui-select" :class="[`s-${size}`, { open, disabled }]">
    <button type="button" class="select-trigger" :disabled="disabled" @click="toggle">
      <span class="trigger-text" :class="{ placeholder: !current }">{{ current?.label || placeholder }}</span>
      <ChevronDown :size="15" :stroke-width="2" class="trigger-chev" />
    </button>

    <Transition name="pop">
      <div v-if="open" class="select-pop card">
        <div
          v-if="clearable && current"
          class="select-option clear-opt"
          @click="choose(null)"
        >
          <RotateCcw :size="14" :stroke-width="2" />
          <span>清除选择</span>
        </div>
        <div
          v-for="opt in options"
          :key="String(opt.value)"
          class="select-option"
          :class="{ active: opt.value === modelValue }"
          @click="choose(opt)"
        >
          <span class="opt-label">{{ opt.label }}</span>
          <Check v-if="opt.value === modelValue" :size="15" :stroke-width="2.4" class="opt-check" />
        </div>
        <div v-if="!options.length" class="select-empty muted">暂无可选项</div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { ChevronDown, Check, RotateCcw } from '@lucide/vue'

const props = defineProps({
  modelValue: { type: [String, Number, Boolean, null], default: '' },
  options: { type: Array, default: () => [] }, // [{ label, value }]
  placeholder: { type: String, default: '请选择' },
  clearable: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  size: { type: String, default: 'md' },
})

const emit = defineEmits(['update:modelValue', 'change'])

const open = ref(false)
const root = ref(null)

const current = computed(() => props.options.find((o) => o.value === props.modelValue) || null)

const toggle = () => { open.value = !open.value }
const close = () => { open.value = false }

const choose = (opt) => {
  const val = opt ? opt.value : ''
  emit('update:modelValue', val)
  emit('change', val)
  close()
}

const onDocDown = (e) => {
  if (root.value && !root.value.contains(e.target)) close()
}
const onKey = (e) => { if (e.key === 'Escape') close() }

watch(open, (v) => {
  if (v) {
    document.addEventListener('mousedown', onDocDown)
    document.addEventListener('keydown', onKey)
  } else {
    document.removeEventListener('mousedown', onDocDown)
    document.removeEventListener('keydown', onKey)
  }
})
onBeforeUnmount(() => {
  document.removeEventListener('mousedown', onDocDown)
  document.removeEventListener('keydown', onKey)
})
</script>

<style scoped>
.ui-select { position: relative; display: inline-block; min-width: 0; }

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  width: 100%;
  background: var(--bg-surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--r-md);
  color: var(--text-1);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
  line-height: 1;
}
.s-sm .select-trigger { height: 32px; padding: 0 10px; font-size: 12.5px; }
.s-md .select-trigger { height: 38px; padding: 0 12px; font-size: 13.5px; }

.select-trigger:hover:not(:disabled) { border-color: var(--text-3); }
.ui-select.open .select-trigger {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-soft-2);
}
.select-trigger:disabled { opacity: 0.55; cursor: not-allowed; }

.trigger-text { overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
.trigger-text.placeholder { color: var(--text-3); }

.trigger-chev {
  color: var(--text-3);
  flex-shrink: 0;
  transition: transform var(--d-base) var(--ease-out);
}
.ui-select.open .trigger-chev { transform: rotate(180deg); color: var(--primary); }

.select-pop {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  min-width: 140px;
  z-index: 60;
  padding: 5px;
  max-height: 260px;
  overflow-y: auto;
  box-shadow: var(--shadow-lg);
  transform-origin: top center;
}

.select-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 10px;
  border-radius: var(--r-sm);
  font-size: 13px;
  color: var(--text-1);
  cursor: pointer;
  transition: background var(--d-fast) var(--ease-out);
  white-space: nowrap;
}
.select-option:hover { background: var(--bg-hover); }
.select-option.active { color: var(--primary); font-weight: 500; background: var(--primary-soft); }
.select-option.clear-opt {
  color: var(--text-3);
  border-bottom: 1px solid var(--border-soft);
  border-radius: var(--r-sm) var(--r-sm) 0 0;
  justify-content: flex-start;
}
.opt-label { overflow: hidden; text-overflow: ellipsis; }
.opt-check { flex-shrink: 0; }
.select-empty { padding: 14px; text-align: center; font-size: 12.5px; }
</style>
