<template>
  <button
    type="button"
    class="ui-check"
    :class="{ checked: modelValue, partial: indeterminate && !modelValue, disabled }"
    role="checkbox"
    :aria-checked="indeterminate ? 'mixed' : modelValue"
    :disabled="disabled"
    @click.stop="$emit('update:modelValue', !modelValue); $emit('change', !modelValue)"
  >
    <Check v-if="modelValue" :size="12" :stroke-width="3.4" />
    <Minus v-else-if="indeterminate" :size="12" :stroke-width="3.4" />
  </button>
</template>

<script setup>
import { Check, Minus } from '@lucide/vue'

defineProps({
  modelValue: { type: Boolean, default: false },
  indeterminate: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue', 'change'])
</script>

<style scoped>
.ui-check {
  width: 18px;
  height: 18px;
  border-radius: 5px;
  border: 1.5px solid var(--border-strong);
  background: var(--bg-surface);
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0;
  transition: all var(--d-fast) var(--ease-out);
  flex-shrink: 0;
}
.ui-check:hover:not(.disabled):not(.checked) { border-color: var(--primary); }
.ui-check.checked, .ui-check.partial {
  background: var(--gradient-brand);
  border-color: transparent;
  animation: pop-in var(--d-base) var(--ease-spring);
}
.ui-check.disabled { opacity: 0.45; cursor: not-allowed; }
</style>
