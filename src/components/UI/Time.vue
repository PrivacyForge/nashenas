<script lang="ts" setup>
import { ref, useAttrs } from 'vue'

const props = defineProps<{
  value: string
}>()

const attrs = useAttrs()

const time = ref<string>('')

function calculate() {
  const then = new Date(props.value)
  const now = new Date()

  const since = Math.trunc((now.getTime() - then.getTime()) / 1000)

  if (since <= 60) {
    time.value = `${since} ثانیه پیش`
  } else if (since <= 3600) {
    time.value = `${Math.trunc(since / 60)} دقیقه پیش`
  } else if (since <= 86400) {
    time.value = `${Math.trunc(since / 60 / 60)} ساعت پیش`
  } else if (since <= 604800) {
    time.value = `${Math.trunc(since / 60 / 60 / 24)} روز قبل`
  } else if (since > 604800) {
    time.value = `${Math.trunc(since / 60 / 60 / 24 / 7)} هفته قبل`
  }
}

calculate()

setInterval(() => {
  calculate()
}, 1000)
</script>

<template>
  <p v-bind="attrs" v-text="time" />
</template>
