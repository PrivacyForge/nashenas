<script lang="ts" setup>
import { ref } from 'vue'

import Message from '@/components/Message.vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user';

const userStore = useUserStore()

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
    userStore.user.receivePublicKey = data.receive_public_key
    userStore.user.sendPublicKey = data.send_public_key

    userStore.isAuth = true
  })
</script>

<template>
  <div class="grid grid-cols-1 gap-y-3 my-4">
    <template v-if="messages?.length">
      <Message v-for="(message, i) in messages" :key="i" :id="message.id" :time="message.time" :owner="message.owner"
        :quote="message.quote" :text="message.content" :canReplay="message.can_replay" :mark="false"
        :sender_public_key="message.sender_public_key" />
    </template>
    <p class="text-center text-gray-400" v-else>
      فعلا هیچ پیامی دریافت نکردی :)
    </p>
  </div>
</template>
