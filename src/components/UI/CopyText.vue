<script lang="ts" setup>
import { ref, useAttrs } from 'vue'

const props = defineProps<{
  text: string
  copy: string
}>()

const attrs = useAttrs()

const delay = ref(false)

function copy() {
  try {
    navigator.clipboard.writeText(props.copy)
  } catch (error) {
    alert(error)
  }

  delay.value = true
  setTimeout(() => {
    delay.value = false
  }, 2000)
}
</script>
<template>
  <button v-if="!delay" v-bind="attrs" v-text="text" @click="copy" />
  <div class="text-center" v-else>کپی شد!</div>
</template>
