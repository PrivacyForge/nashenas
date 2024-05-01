<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Message from '@/components/Message.vue'
import axios from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

if (!authStore.user.PublicKey) router.push({ name: 'setup' })

const messages = ref<any[]>([])

axios.get('/get-messages').then((response) => {
  messages.value = response.data
})
</script>

<template>
  <div class="grid grid-cols-1 gap-y-3 m-4">
    <template v-if="messages.length">
      <Message
        v-for="(m, i) in messages"
        :key="i"
        :text="m.Message"
        :time="m.Time"
        :mark="false"
      />
    </template>
    <p class="text-center text-gray-400" v-else>You have no message :)</p>
  </div>
</template>
