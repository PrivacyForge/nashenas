<script lang="ts" setup>
import { ref } from 'vue'

import Message from '@/components/Message.vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user';
import { useRouter } from 'vue-router';

const userStore = useUserStore()
const router = useRouter()

const messages = ref<any[]>([])

axios.get('/get-messages').then((response) => {
  messages.value = response.data

  messages.value.reverse()
})

axios
  .get('/me')
  .then(({ data }) => {
    userStore.user.id = data.id
    userStore.user.userid = data.userid
    userStore.user.username = data.username
    userStore.user.publicKey = data.public_key

    userStore.isAuth = true
  })
  .catch((error) => {
    alert(error)
    router.push({ name: 'error' })
  })
</script>

<template>
  <div class="grid grid-cols-1 gap-y-3 my-4">
    <template v-if="messages?.length">
      {{ messages }}
      <Message v-for="(message, i) in messages" :key="i" :message="message" />
    </template>
    <p class="text-center text-gray-400" v-else>
      فعلا هیچ پیامی دریافت نکردی :)
    </p>
  </div>
</template>
