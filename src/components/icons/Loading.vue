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
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    style="shape-rendering: auto"
    viewBox="0 0 100 100"
    preserveAspectRatio="xMidYMid"
  >
    <path
      fill="none"
      :stroke="color"
      stroke-width="10"
      stroke-dasharray="205.271142578125 51.317785644531256"
      d="M24.3 30C11.4 30 5 43.3 5 50s6.4 20 19.3 20c19.3 0 32.1-40 51.4-40 C88.6 30 95 43.3 95 50s-6.4 20-19.3 20C56.4 70 43.6 30 24.3 30z"
      stroke-linecap="round"
      style="transform: scale(0.8); transform-origin: 50px 50px"
    >
      <animate
        attributeName="stroke-dashoffset"
        repeatCount="indefinite"
        dur="2s"
        keyTimes="0;1"
        values="0;256.58892822265625"
      ></animate>
    </path>
  </svg>
</template>
