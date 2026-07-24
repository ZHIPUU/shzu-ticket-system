<template>
  <UiModal
    :model-value="confirmState.open"
    :title="confirmState.title"
    width="420px"
    @update:model-value="onCancel"
  >
    <div class="confirm-body">
      <div class="confirm-icon" :class="{ danger: confirmState.danger }">
        <TriangleAlert v-if="confirmState.danger" :size="22" :stroke-width="2" />
        <CircleQuestionMark v-else :size="22" :stroke-width="2" />
      </div>
      <p class="confirm-msg">{{ confirmState.message }}</p>

      <div v-if="confirmState.mode === 'prompt'" class="confirm-input">
        <UiInput
          ref="inputRef"
          v-model="confirmState.inputValue"
          :placeholder="confirmState.placeholder"
          clearable
          @enter="onOk"
        />
        <p v-if="errorMsg" class="confirm-error anim-shake" :key="errorMsg">{{ errorMsg }}</p>
      </div>
    </div>

    <template #footer>
      <UiButton variant="ghost" @click="onCancel">{{ confirmState.cancelText }}</UiButton>
      <UiButton :variant="confirmState.danger ? 'danger' : 'primary'" @click="onOk">
        {{ confirmState.confirmText }}
      </UiButton>
    </template>
  </UiModal>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { TriangleAlert, CircleQuestionMark } from '@lucide/vue'
import UiModal from './UiModal.vue'
import UiButton from './UiButton.vue'
import UiInput from './UiInput.vue'
import { confirmState, resolveDialog } from './confirm'

const errorMsg = ref('')
const inputRef = ref(null)

watch(
  () => confirmState.open,
  async (v) => {
    errorMsg.value = ''
    if (v && confirmState.mode === 'prompt') {
      await nextTick()
      // 等弹窗过渡后聚焦
      setTimeout(() => inputRef.value?.focus(), 240)
    }
  }
)

const onOk = () => {
  if (confirmState.mode === 'prompt') {
    const val = (confirmState.inputValue || '').trim()
    if (confirmState.pattern && !confirmState.pattern.test(val)) {
      errorMsg.value = confirmState.patternMessage || '输入格式不正确'
      return
    }
    resolveDialog(val)
    return
  }
  resolveDialog(true)
}

const onCancel = () => {
  if (!confirmState.open) return
  resolveDialog(confirmState.mode === 'prompt' ? null : false)
}
</script>

<style scoped>
.confirm-body { display: flex; flex-direction: column; align-items: center; text-align: center; }

.confirm-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--r-full);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-soft);
  color: var(--primary);
  margin-bottom: 14px;
}
.confirm-icon.danger { background: var(--danger-soft); color: var(--danger); }

.confirm-msg {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-2);
  white-space: pre-wrap;
}

.confirm-input { width: 100%; margin-top: 16px; }
.confirm-error {
  margin: 8px 0 0;
  font-size: 12.5px;
  color: var(--danger);
  text-align: left;
}
</style>
