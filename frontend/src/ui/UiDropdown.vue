<template>
  <div ref="root" class="ui-dropdown">
    <div class="dropdown-trigger" @click="toggle">
      <slot name="trigger" :open="open" />
    </div>
    <Transition name="pop">
      <div v-if="open" class="dropdown-pop card" :class="`place-${placement}`">
        <template v-for="(item, i) in items" :key="i">
          <div v-if="item.divided && i > 0" class="dropdown-divider" />
          <button
            type="button"
            class="dropdown-item"
            :class="{ danger: item.danger }"
            :disabled="item.disabled"
            @click="pick(item)"
          >
            <component v-if="item.icon" :is="item.icon" :size="15" :stroke-width="2" class="item-icon" />
            <span>{{ item.label }}</span>
          </button>
        </template>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, watch, onBeforeUnmount } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] }, // [{ label, value, icon?, danger?, divided?, disabled? }]
  placement: { type: String, default: 'end' }, // start | end
})
const emit = defineEmits(['select'])

const open = ref(false)
const root = ref(null)

const toggle = () => { open.value = !open.value }
const close = () => { open.value = false }

const pick = (item) => {
  if (item.disabled) return
  close()
  emit('select', item.value)
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
.ui-dropdown { position: relative; display: inline-block; }
.dropdown-trigger { display: inline-flex; cursor: pointer; }

.dropdown-pop {
  position: absolute;
  top: calc(100% + 8px);
  z-index: 70;
  min-width: 168px;
  padding: 6px;
  box-shadow: var(--shadow-lg);
  transform-origin: top right;
}
.place-end { right: 0; }
.place-start { left: 0; transform-origin: top left; }

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
  padding: 9px 11px;
  border: none;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-1);
  font-size: 13.5px;
  cursor: pointer;
  text-align: left;
  transition: background var(--d-fast) var(--ease-out);
}
.dropdown-item:hover:not(:disabled) { background: var(--bg-hover); }
.dropdown-item:disabled { opacity: 0.45; cursor: not-allowed; }
.dropdown-item.danger { color: var(--danger); }
.dropdown-item.danger:hover:not(:disabled) { background: var(--danger-soft); }
.item-icon { flex-shrink: 0; color: var(--text-3); }
.dropdown-item.danger .item-icon { color: var(--danger); }

.dropdown-divider {
  height: 1px;
  background: var(--border-soft);
  margin: 5px 8px;
}
</style>
