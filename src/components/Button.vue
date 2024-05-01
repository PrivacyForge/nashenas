<script lang="ts" setup>
import { useAttrs } from 'vue'

defineProps<{
  loading?: boolean
}>()

const attrs = useAttrs()
</script>
<template>
  <button
    v-bind="attrs"
    class="btn bg-[#119af5] text-white btn-block disabled:bg-[#119af5] disabled:text-white"
    :class="loading && ['opacity-70']"
    :disabled="loading"
  >
    <!-- Send to @{{ $route.params.username }} -->
    <slot />
    <span v-if="loading" class="loading loading-infinity loading-md"></span>
  </button>
</template>
