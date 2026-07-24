<template>
  <div class="ui-field" :class="[`s-${size}`, { focused, disabled, 'has-prefix': !!$slots.prefix, textarea: isTextarea }]">
    <span v-if="$slots.prefix" class="field-prefix"><slot name="prefix" /></span>

    <textarea
      v-if="isTextarea"
      ref="inputRef"
      class="field-control"
      :value="modelValue"
      :placeholder="placeholder"
      :maxlength="maxlength"
      :rows="rows"
      :disabled="disabled"
      @input="onInput"
      @focus="focused = true"
      @blur="focused = false"
      @keydown="$emit('keydown', $event)"
    />
    <input
      v-else
      ref="inputRef"
      class="field-control"
      :value="modelValue"
      :type="actualType"
      :placeholder="placeholder"
      :maxlength="maxlength"
      :disabled="disabled"
      :autofocus="autofocus"
      @input="onInput"
      @focus="focused = true"
      @blur="focused = false"
      @keyup.enter="$emit('enter')"
      @keydown="$emit('keydown', $event)"
    />

    <span v-if="isPassword" class="field-suffix eye" @click="showPwd = !showPwd">
      <Eye v-if="!showPwd" :size="15" :stroke-width="2" />
      <EyeOff v-else :size="15" :stroke-width="2" />
    </span>
    <span v-else-if="clearable && modelValue && !disabled" class="field-suffix clear" @click="clear">
      <X :size="14" :stroke-width="2.2" />
    </span>
    <span v-if="showWordLimit && maxlength" class="word-count tnum">{{ String(modelValue || '').length }}/{{ maxlength }}</span>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { X, Eye, EyeOff } from '@lucide/vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  type: { type: String, default: 'text' }, // text | password | textarea | date | email | number
  placeholder: { type: String, default: '' },
  maxlength: { type: [String, Number], default: undefined },
  showWordLimit: { type: Boolean, default: false },
  rows: { type: [String, Number], default: 3 },
  clearable: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  autofocus: { type: Boolean, default: false },
  size: { type: String, default: 'md' }, // sm | md | lg
})

const emit = defineEmits(['update:modelValue', 'enter', 'keydown', 'clear'])

const focused = ref(false)
const showPwd = ref(false)
const inputRef = ref(null)

const isTextarea = computed(() => props.type === 'textarea')
const isPassword = computed(() => props.type === 'password')
const actualType = computed(() => (isPassword.value ? (showPwd.value ? 'text' : 'password') : props.type))

const onInput = (e) => emit('update:modelValue', e.target.value)
const clear = () => {
  emit('update:modelValue', '')
  emit('clear')
}

const focus = () => inputRef.value?.focus()
defineExpose({ focus })
</script>

<style scoped>
.ui-field {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  background: var(--bg-surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--r-md);
  transition: border-color var(--d-fast) var(--ease-out), box-shadow var(--d-fast) var(--ease-out),
    background var(--d-fast) var(--ease-out);
  position: relative;
}
.ui-field:hover:not(.disabled):not(.focused) { border-color: var(--text-3); }
.ui-field.focused {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-soft-2);
}
.ui-field.disabled { opacity: 0.55; cursor: not-allowed; background: var(--bg-sink); }

.s-sm { min-height: 32px; padding: 0 10px; font-size: 12.5px; }
.s-md { min-height: 38px; padding: 0 12px; font-size: 13.5px; }
.s-lg { min-height: 46px; padding: 0 14px; font-size: 15px; border-radius: var(--r-md); }

.ui-field.textarea { align-items: flex-start; padding-top: 9px; padding-bottom: 9px; }

.field-control {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  color: var(--text-1);
  font-size: inherit;
  font-family: inherit;
  line-height: 1.55;
  padding: 0;
  resize: vertical;
}
.field-control::placeholder { color: var(--text-3); }
.field-control:disabled { cursor: not-allowed; }

.field-prefix {
  display: inline-flex;
  align-items: center;
  color: var(--text-3);
  flex-shrink: 0;
}
.ui-field.focused .field-prefix { color: var(--primary); }

.field-suffix {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-3);
  cursor: pointer;
  border-radius: 5px;
  padding: 2px;
  transition: all var(--d-fast) var(--ease-out);
  flex-shrink: 0;
}
.field-suffix:hover { color: var(--text-1); background: var(--bg-hover); }

.word-count {
  font-size: 11px;
  color: var(--text-3);
  flex-shrink: 0;
  align-self: flex-end;
}

/* 原生 date 输入美化 */
input[type="date"].field-control { color-scheme: light; }
html.dark input[type="date"].field-control { color-scheme: dark; }
input[type="date"].field-control::-webkit-calendar-picker-indicator {
  cursor: pointer;
  opacity: 0.55;
}
</style>
