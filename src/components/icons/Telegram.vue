<script lang="ts" setup>
import { useAttrs } from 'vue'

defineProps<{
  color: string
  size: string
}>()

const attrs = useAttrs()
</script>

<template>
  <svg
    v-bind="attrs"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M23.9303 4.62262L20.3089 19.5664C20.0357 20.6211 19.3232 20.8836 18.3107 20.3867L12.7928 16.8289L10.1303 19.0695C9.83569 19.3273 9.58926 19.5429 9.02141 19.5429L9.41784 14.6257L19.6446 6.53981C20.0893 6.19294 19.5482 6.00075 18.9536 6.34762L6.31069 13.3132L0.867837 11.8226C-0.316092 11.4992 -0.33752 10.7867 1.11427 10.2898L22.4036 3.11325C23.3893 2.78981 24.2518 3.30544 23.9303 4.62262Z"
      :fill="color"
    />
  </svg>
</template>
