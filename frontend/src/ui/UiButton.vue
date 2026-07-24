<template>
  <button
    class="ui-btn"
    :class="[`v-${variant}`, `s-${size}`, { block, loading, 'icon-only': iconOnly }]"
    :disabled="disabled || loading"
    :type="type"
  >
    <span v-if="loading" class="btn-spinner" aria-hidden="true" />
    <span v-else-if="$slots.icon" class="btn-icon"><slot name="icon" /></span>
    <span v-if="$slots.default" class="btn-label"><slot /></span>
  </button>
</template>

<script setup>
defineProps({
  variant: { type: String, default: 'secondary' }, // primary | secondary | ghost | danger | danger-soft | text
  size: { type: String, default: 'md' },           // sm | md | lg
  type: { type: String, default: 'button' },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  block: { type: Boolean, default: false },
  iconOnly: { type: Boolean, default: false },
})
</script>

<style scoped>
.ui-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  border: 1px solid transparent;
  border-radius: var(--r-md);
  font-weight: 500;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
  line-height: 1;
  transition: all var(--d-fast) var(--ease-out);
  position: relative;
}
.ui-btn:active:not(:disabled) { transform: scale(0.965); }
.ui-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

/* 尺寸 */
.s-sm { height: 30px; padding: 0 12px; font-size: 12.5px; border-radius: var(--r-sm); gap: 5px; }
.s-md { height: 36px; padding: 0 16px; font-size: 13.5px; }
.s-lg { height: 44px; padding: 0 22px; font-size: 15px; border-radius: var(--r-md); }
.s-sm.icon-only { width: 30px; padding: 0; }
.s-md.icon-only { width: 36px; padding: 0; }
.s-lg.icon-only { width: 44px; padding: 0; }

.block { width: 100%; }

/* 变体 */
.v-primary {
  background: var(--gradient-brand);
  color: #fff;
  box-shadow: var(--shadow-brand);
}
.v-primary:hover:not(:disabled) {
  background: var(--gradient-brand-hover);
  box-shadow: 0 8px 22px -6px rgba(99, 102, 241, 0.55);
  transform: translateY(-1px);
}
.v-primary:active:not(:disabled) { transform: scale(0.965); }

.v-secondary {
  background: var(--bg-surface);
  border-color: var(--border-strong);
  color: var(--text-1);
  box-shadow: var(--shadow-xs);
}
.v-secondary:hover:not(:disabled) {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-soft);
}

.v-ghost {
  background: transparent;
  border-color: var(--border);
  color: var(--text-2);
}
.v-ghost:hover:not(:disabled) {
  background: var(--bg-hover);
  color: var(--text-1);
}

.v-danger {
  background: var(--danger);
  color: #fff;
  box-shadow: 0 6px 16px -6px rgba(229, 72, 77, 0.5);
}
.v-danger:hover:not(:disabled) { filter: brightness(1.06); transform: translateY(-1px); }

.v-danger-soft {
  background: transparent;
  border-color: var(--border);
  color: var(--danger);
}
.v-danger-soft:hover:not(:disabled) {
  background: var(--danger-soft);
  border-color: var(--danger);
}

.v-text {
  background: transparent;
  color: var(--primary);
  padding-left: 8px;
  padding-right: 8px;
}
.v-text:hover:not(:disabled) { background: var(--primary-soft); }

.btn-icon { display: inline-flex; align-items: center; flex-shrink: 0; }
.btn-label { display: inline-flex; align-items: center; }

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}
</style>
