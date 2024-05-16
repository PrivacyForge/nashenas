<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'

import Card from '@/components/UI/Card.vue'

const userStore = useUserStore()

const route = useRoute()
const router = useRouter()

const errorMessage = ref('')
const code = route.params.code

setTimeout(() => {
  axios
    .get(`/confirm/${code}`)
    .then(({ data }) => {
      localStorage.setItem('receive_public_key', data.receive_public_key)
      localStorage.setItem('send_public_key', data.send_public_key)
      localStorage.setItem('token', data.token)

      userStore.user.id = data.id
      userStore.user.receivePublicKey = data.receive_public_key
      userStore.user.sendPublicKey = data.send_public_key
      userStore.user.username = data.username
      userStore.user.userid = data.userid

      userStore.isAuth = true

      if (!data.receive_public_key || !data.send_public_key) {
        router.push({ name: 'setup' })
      } else {
        router.push({ name: 'inbox' })
      }
    })
    .catch((error) => {
      errorMessage.value = error.response.data.message
    })
    .finally(() => {})
}, 1000)
</script>

<template>
  <Card class="grid grid-cols-1 m-4">
    <div v-if="!errorMessage" class="flex justify-center items-center py-5">
      <p class="text-[#119af5] font-semibold">Waiting for server response...</p>
      <span class="loading loading-infinity loading-lg mx-2 text-[#119af5]" />
    </div>
    <template v-else>
      <p class="text-red-500 text-center" v-text="errorMessage"></p>
    </template>
  </Card>
</template>
