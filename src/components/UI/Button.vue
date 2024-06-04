<script lang="ts" setup>
import { useAttrs } from 'vue'
import LoadingIcon from '@/components/icons/Loading.vue'

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
    class="flex justify-center bg-[#119af5] text-white disabled:bg-[#119af5] py-3 font-bold rounded-md active:scale-95"
    :class="[(loading || disabled) && 'opacity-70', block && 'btn-block']"
    :disabled="loading || disabled"
  >
    <slot />
    <LoadingIcon v-if="loading" color="white" size="26px" />
  </button>
</template>
