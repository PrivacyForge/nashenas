<script lang="ts" setup>
import Message from '@/components/Message.vue'
import { ref } from 'vue'

const messages = ref<any>()

messages.value = JSON.parse(localStorage.getItem('messages')!)
</script>

<template>
  <div class="grid grid-cols-1 gap-y-3 m-4">
    <Message
      v-for="(m, i) in messages"
      :key="i"
      :text="m.text"
      :time="m.time"
    />
  </div>
</template>
