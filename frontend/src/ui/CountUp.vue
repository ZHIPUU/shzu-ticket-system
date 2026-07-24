<template>
  <span class="count-up tnum">{{ display }}</span>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const props = defineProps({
  value: { type: [Number, String], default: 0 },
  duration: { type: Number, default: 700 },
})

const display = ref(0)
let raf = null

const animate = (from, to) => {
  if (typeof to !== 'number' || isNaN(to)) { display.value = to; return }
  if (raf) cancelAnimationFrame(raf)
  const start = performance.now()
  const tick = (now) => {
    const t = Math.min(1, (now - start) / props.duration)
    const eased = 1 - Math.pow(1 - t, 3)
    display.value = Math.round(from + (to - from) * eased)
    if (t < 1) raf = requestAnimationFrame(tick)
  }
  raf = requestAnimationFrame(tick)
}

watch(() => props.value, (nv, ov) => {
  animate(typeof ov === 'number' ? ov : 0, nv)
})
onMounted(() => animate(0, props.value))
</script>
