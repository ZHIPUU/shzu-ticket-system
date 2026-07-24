<template>
  <Teleport to="body">
    <div class="toast-wrap" aria-live="polite">
      <TransitionGroup name="toast">
        <div v-for="t in toasts" :key="t.id" class="toast card" :data-type="t.type">
          <span class="toast-icon">
            <CircleCheckBig v-if="t.type === 'success'" :size="16" :stroke-width="2.2" />
            <CircleAlert v-else-if="t.type === 'error'" :size="16" :stroke-width="2.2" />
            <TriangleAlert v-else-if="t.type === 'warning'" :size="16" :stroke-width="2.2" />
            <Info v-else :size="16" :stroke-width="2.2" />
          </span>
          <span class="toast-msg">{{ t.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { CircleCheckBig, CircleAlert, TriangleAlert, Info } from '@lucide/vue'
import { toasts } from './toast'
</script>

<style scoped>
.toast-wrap {
  position: fixed;
  top: 18px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 300;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  pointer-events: none;
  width: max-content;
  max-width: calc(100vw - 32px);
}

.toast {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 16px;
  border-radius: var(--r-full);
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-1);
  box-shadow: var(--shadow-lg);
  pointer-events: auto;
  max-width: 100%;
}

.toast-icon { display: inline-flex; flex-shrink: 0; }
.toast[data-type="success"] .toast-icon { color: var(--success); }
.toast[data-type="error"] .toast-icon { color: var(--danger); }
.toast[data-type="warning"] .toast-icon { color: var(--warning); }
.toast[data-type="info"] .toast-icon { color: var(--primary); }

.toast-msg {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
