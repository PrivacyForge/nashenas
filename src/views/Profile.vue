<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import forge from 'node-forge'
import axios from '@/plugins/axios'

import TelegramIcon from '@/components/icons/Telegram.vue'
import GithubIcon from '@/components/icons/Github.vue'

const route = useRoute()

const message = ref('')
const username = route.params.username

const publicKey = ref('')
const id = ref<number>()

axios.get(`/get-profile?username=${username}`).then((response) => {
  publicKey.value = response.data.PublicKey
  id.value = response.data.ID
})

function submit() {
  const encryptedMsg = forge.pki
    .publicKeyFromPem(publicKey.value!)
    .encrypt(message.value)

  axios
    .post(`/send-message?id=${id.value}`, {
      message: forge.util.bytesToHex(encryptedMsg),
    })
    .then(() => {
      message.value = ''
    })
}
</script>

<template>
  <div class="grid grid-cols-1 gap-4 m-4 bg-[#ffffff] p-5 rounded-lg shadow-sm">
    <p>Write something for Yasha...</p>
    <textarea
      v-model="message"
      class="textarea textarea-bordered w-full focus:outline-[#119af5]"
      placeholder="Message..."
      dir="auto"
    ></textarea>
    <button class="btn bg-[#119af5] text-white btn-block" @click="submit">
      Send to @{{ $route.params.username }}
    </button>
  </div>
  <div class="fixed bottom-2 left-[50%] w-max translate-x-[-50%] text-gray-500">
    <div class="flex justify-center items-center pb-2">
      <GithubIcon class="mr-1" size="26" color="gray" />
      <TelegramIcon class="ml-1" size="24" color="gray" />
    </div>
    <p>Open Source, E2E Encrption, Safe.</p>
  </div>
</template>
