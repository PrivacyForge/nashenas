<script setup lang="ts">
import { RouterView } from 'vue-router'
function copy(value: string) {
  navigator.clipboard.writeText(value)
}
</script>

<template>
  <span @click="copy($route.fullPath)">click</span>
  <RouterView />
</template>
