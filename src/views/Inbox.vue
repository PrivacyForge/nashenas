<script lang="ts" setup>
import { ref } from 'vue'
import axios from '@/plugins/axios'
import Message from '@/components/Message.vue'

const messages = ref<any[]>([])

axios.get('/get-messages').then((response) => {
  messages.value = response.data

  messages.value.reverse()
})
</script>

<template>
  <div class="grid grid-cols-1 gap-y-3 m-4">
    <template v-if="messages.length">
      <Message
        v-for="(m, i) in messages"
        :key="i"
        :text="m.Content"
        :time="m.Time"
        :mark="false"
      />
    </template>
    <p class="text-center text-gray-400" v-else>You have no message :)</p>
  </div>
</template>
