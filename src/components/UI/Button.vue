<script lang="ts" setup>
import { useAttrs } from 'vue'

defineProps<{
  loading?: boolean
  disabled?: boolean
  block?: boolean
}>()

const attrs = useAttrs()
</script>
<template>
  <button
    v-bind="attrs"
    class="bg-[#119af5] text-white disabled:bg-[#119af5] py-3 font-bold rounded-md active:scale-95"
    :class="[(loading || disabled) && 'opacity-70', block && 'btn-block']"
    :disabled="loading || disabled"
  >
    <!-- Send to @{{ $route.params.username }} -->
    <slot />
    <span v-if="loading" class="loading loading-infinity loading-md"></span>
  </button>
</template>
