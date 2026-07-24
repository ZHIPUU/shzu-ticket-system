<template>
  <button
    type="button"
    class="ui-switch"
    :class="{ on: modelValue, disabled }"
    role="switch"
    :aria-checked="modelValue"
    :disabled="disabled"
    @click="toggle"
  >
    <span class="switch-knob" />
  </button>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue', 'change'])
const toggle = () => {
  if (props.disabled) return
  emit('update:modelValue', !props.modelValue)
  emit('change', !props.modelValue)
}
</script>

<style scoped>
.ui-switch {
  width: 42px;
  height: 24px;
  border-radius: var(--r-full);
  border: none;
  background: var(--bg-active);
  cursor: pointer;
  position: relative;
  padding: 2px;
  transition: background var(--d-base) var(--ease-out);
  flex-shrink: 0;
}
.ui-switch.on { background: var(--gradient-brand); }
.ui-switch.disabled { opacity: 0.5; cursor: not-allowed; }

.switch-knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.22);
  transition: transform var(--d-base) var(--ease-spring);
}
.ui-switch.on .switch-knob { transform: translateX(18px); }
</style>
