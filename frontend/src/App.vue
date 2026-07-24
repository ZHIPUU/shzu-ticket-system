<template>
  <router-view />
  <ToastHost />
  <ConfirmHost />
</template>

<script setup>
import ToastHost from './ui/ToastHost.vue'
import ConfirmHost from './ui/ConfirmHost.vue'
</script>
