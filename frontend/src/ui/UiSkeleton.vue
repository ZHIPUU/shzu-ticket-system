<template>
  <div class="ui-skeleton" :class="`t-${type}`">
    <div v-if="type === 'card'" class="sk-card" />
    <template v-else-if="type === 'rows'">
      <div v-for="i in count" :key="i" class="sk-row" :style="{ '--i': i }">
        <div class="sk-block sk-w-sm" />
        <div class="sk-block sk-w-lg" />
        <div class="sk-block sk-w-md" />
      </div>
    </template>
    <div v-else v-for="i in count" :key="i" class="sk-block sk-line" :style="{ width: i === count ? '60%' : '100%' }" />
  </div>
</template>

<script setup>
defineProps({
  type: { type: String, default: 'rows' }, // rows | lines | card
  count: { type: Number, default: 5 },
})
</script>

<style scoped>
.sk-block, .sk-card {
  border-radius: var(--r-sm);
  background: linear-gradient(
    90deg,
    var(--bg-hover) 25%,
    var(--bg-active) 50%,
    var(--bg-hover) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s ease infinite;
}

.sk-card { height: 120px; border-radius: var(--r-lg); }

.sk-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 15px 4px;
}
.sk-w-sm { width: 90px; height: 14px; flex-shrink: 0; }
.sk-w-lg { flex: 1; height: 14px; }
.sk-w-md { width: 70px; height: 14px; flex-shrink: 0; }

.sk-line { height: 14px; margin-bottom: 12px; }
</style>
