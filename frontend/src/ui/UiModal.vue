<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-overlay" @mousedown="onOverlayDown">
        <div class="modal-panel card" :style="panelStyle" role="dialog" :aria-label="title">
          <div class="sheet-handle" aria-hidden="true" />
          <header v-if="title || closable" class="modal-head">
            <h3 class="modal-title">{{ title }}</h3>
            <button v-if="closable" class="modal-close" aria-label="关闭" @click="close">
              <X :size="17" :stroke-width="2" />
            </button>
          </header>
          <div class="modal-body">
            <slot />
          </div>
          <footer v-if="$slots.footer" class="modal-foot">
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, watch, onBeforeUnmount } from 'vue'
import { X } from '@lucide/vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: '' },
  width: { type: String, default: '520px' },
  closable: { type: Boolean, default: true },
  closeOnOverlay: { type: Boolean, default: true },
})

const emit = defineEmits(['update:modelValue', 'close'])

const panelStyle = computed(() => ({ maxWidth: props.width }))

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const onOverlayDown = (e) => {
  if (props.closeOnOverlay && e.target === e.currentTarget) close()
}

const onKey = (e) => {
  if (e.key === 'Escape' && props.modelValue && props.closable) close()
}

watch(
  () => props.modelValue,
  (v) => {
    if (v) {
      document.addEventListener('keydown', onKey)
      document.body.style.overflow = 'hidden'
    } else {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }
)
onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(15, 17, 26, 0.45);
  backdrop-filter: blur(3px);
  -webkit-backdrop-filter: blur(3px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.modal-panel {
  width: 100%;
  max-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

.sheet-handle { display: none; }

.modal-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 18px 22px 0;
  flex-shrink: 0;
}
.modal-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.2px;
  color: var(--text-1);
}
.modal-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
  flex-shrink: 0;
}
.modal-close:hover { background: var(--bg-hover); color: var(--text-1); }

.modal-body {
  padding: 16px 22px;
  overflow-y: auto;
  flex: 1;
}

.modal-foot {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 0 22px 18px;
  flex-shrink: 0;
}

/* 移动端：底部抽屉 */
@media (max-width: 640px) {
  .modal-overlay {
    padding: 0;
    align-items: flex-end;
  }
  .modal-panel {
    max-width: 100% !important;
    max-height: 86vh;
    border-radius: var(--r-xl) var(--r-xl) 0 0;
    border-bottom: none;
  }
  .sheet-handle {
    display: block;
    width: 38px;
    height: 4px;
    border-radius: 2px;
    background: var(--border-strong);
    margin: 10px auto 2px;
    flex-shrink: 0;
  }
  .modal-head { padding: 12px 20px 0; }
  .modal-body { padding: 14px 20px; }
  .modal-foot {
    padding: 4px 20px calc(16px + env(safe-area-inset-bottom, 0px));
  }
}
</style>
