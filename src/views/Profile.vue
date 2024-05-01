<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import forge from 'node-forge'
import axios from '@/plugins/axios'

import TelegramIcon from '@/components/icons/Telegram.vue'
import GithubIcon from '@/components/icons/Github.vue'
import Button from '@/components/Button.vue'

const route = useRoute()

const message = ref('')
const username = route.params.username

const loading = ref(true)
const notFoundUser = ref<boolean>()
const errorMessage = ref('')

const submitLoading = ref(false)

const sent = ref(false)

const user = reactive<{
  id: number | null
  publicKey: string | null
}>({
  id: null,
  publicKey: null,
})

function submit() {
  if (message.value === '') return

  const encryptedMsg = forge.pki
    .publicKeyFromPem(user.publicKey!)
    .encrypt(message.value)

  submitLoading.value = true

  setTimeout(() => {
    axios
      .post(`/send-message?id=${user.id}`, {
        message: forge.util.bytesToHex(encryptedMsg),
      })
      .then(() => {
        message.value = ''
        sent.value = true
      })
      .finally(() => {
        submitLoading.value = false
      })
  }, 2000)
}

onMounted(() => {
  setTimeout(() => {
    axios
      .get(`/get-profile?username=${username}`)
      .then((response) => {
        user.publicKey = response.data.PublicKey
        user.id = response.data.ID
      })
      .catch((error) => {
        notFoundUser.value = true
        errorMessage.value = error.response.data ?? 'Error.'
      })
      .finally(() => {
        loading.value = false
      })
  }, 1000)
})
</script>

<template>
  <div class="grid grid-cols-1 gap-4 m-4 bg-[#ffffff] p-5 rounded-lg shadow-sm">
    <div v-if="loading" class="flex justify-center items-center py-6">
      <p class="text-[#119af5] font-semibold">Waiting for server response...</p>
      <span class="loading loading-infinity loading-lg mx-2 text-[#119af5]" />
    </div>
    <template v-else>
      <template v-if="!sent">
        <template v-if="notFoundUser">
          <p
            class="text-center text-[#119af5] font-semibold"
            v-text="errorMessage"
          />
        </template>
        <template v-else>
          <p>Write something for Yasha...</p>
          <textarea
            v-model="message"
            class="textarea textarea-bordered w-full focus:outline-[#119af5]"
            placeholder="Message..."
            dir="auto"
          ></textarea>
          <Button :loading="submitLoading" @click="submit"
            >Send to @{{ $route.params.username }}</Button
          >
        </template>
      </template>
      <template v-else>
        <p class="text-center text-[#119af5] font-semibold">
          Your encrypted message has been sent.
        </p>
        <Button @click="sent = false">Send another message</Button>
      </template>
    </template>
  </div>
  <div class="fixed bottom-2 left-[50%] w-max translate-x-[-50%] text-gray-500">
    <div class="flex justify-center items-center pb-2">
      <GithubIcon class="mr-1" size="26" color="gray" />
      <TelegramIcon class="ml-1" size="24" color="gray" />
    </div>
    <p>Open Source, E2E Encrption, Safe.</p>
  </div>
</template>
