<script lang="ts" setup>
import { ref } from 'vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateKeyPair } from '@/cryptography/RSA'
import { generateKeysTemplate, extractKeys, exportKeys } from '@/utils'

const userStore = useUserStore()

const username = ref(userStore.user.username)
const usernameErr = ref('')

const loading = ref(false)
const FileInput = ref<any>()

const state = ref<
  'set-username' | 'key-question' | 'key-generation' | 'key-upload'
>('set-username')

async function generateKeys() {
  loading.value = true

  let receivePublicKey: string, sendPublicKey: string

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    localStorage.setItem('receive_private_key', privateKey)
    localStorage.setItem('receive_public_key', publicKey)
    receivePublicKey = publicKey
  })

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    localStorage.setItem('send_private_key', privateKey)
    localStorage.setItem('send_public_key', publicKey)
    sendPublicKey = publicKey
  })

  setTimeout(() => {
    axios
      .post('/set-key', {
        receive_public_key: receivePublicKey,
        send_public_key: sendPublicKey,
      })
      .then(({ data }) => {
        state.value = 'key-generation'
        userStore.user.receivePublicKey = data.receive_public_key
        userStore.user.sendPublicKey = data.send_public_key
      })
      .finally(() => {
        loading.value = false
      })
  }, 1500)
}

function exportHandler() {
  exportKeys()
  state.value = 'key-generation'
}

function importKeys(event: Event) {
  // @ts-ignore
  const file = event.target!.files[0]
  const reader = new FileReader()
  reader.onload = (e) => {
    const rawData = e.target!.result as string

    const {
      receivePrivateKey,
      receivePublicKey,
      sendPrivateKey,
      sendPublicKey,
    } = extractKeys(rawData)

    localStorage.setItem('receive_private_key', receivePrivateKey)
    localStorage.setItem('receive_public_key', receivePublicKey)

    localStorage.setItem('send_private_key', sendPrivateKey)
    localStorage.setItem('send_public_key', sendPublicKey)

    axios
      .post('/set-key', {
        receive_public_key: receivePublicKey,
        send_public_key: sendPublicKey,
      })
      .then(({ data }) => {
        state.value = 'key-upload'
        userStore.user.receivePublicKey = data.receive_public_key
        userStore.user.sendPublicKey = data.send_public_key
      })
  }
  reader.readAsText(file)
}

function usernameSubmit() {
  axios
    .post('/set-username', { username: username.value })
    .then(() => {
      state.value = 'key-question'
      userStore.user.username = username.value
    })
    .catch((error) => {
      usernameErr.value = error.response.data.message
    })
}
</script>
<template>
  <div class="grid grid-cols-1 m-4 bg-white shadow-sm p-5 rounded-lg">
    <div v-if="loading" class="flex justify-center items-center py-6">
      <p class="text-[#119af5] font-semibold">Generating your keys...</p>
      <span class="loading loading-infinity loading-lg mx-2 text-[#119af5]" />
    </div>
    <template v-else>
      <template v-if="state === 'set-username'">
        <div class="relative">
          <label class="font-semibold"> Choose a stylish username: </label>
          <input
            v-model="username"
            type="text"
            class="input input-bordered w-full focus:outline-[#119af5] outline-[#119af5] mt-3 pl-6"
            placeholder="username"
          />
          <span class="absolute left-2 top-12">@</span>
        </div>
        <p v-if="usernameErr" class="text-red-500 mt-2" v-text="usernameErr" />
        <button
          class="btn mt-4 bg-[#119af5] text-white"
          @click="usernameSubmit()"
        >
          Next
        </button>
      </template>

      <template v-if="state === 'key-question'">
        <p class="pb-5 pt-3 text-center">Do you have any public/private key?</p>
        <div class="grid grid-cols-1 gap-y-2">
          <button
            class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
            @click="generateKeys"
          >
            No, Please Generate.
          </button>
          <button
            class="text-[#119af5] py-2 rounded-md font-semibold"
            @click="FileInput.click()"
          >
            Yes, I can upload them.
          </button>
          <input
            ref="FileInput"
            type="file"
            class="hidden"
            @change="importKeys"
          />
        </div>
      </template>

      <template v-if="state === 'key-upload'">
        <p class="text-center text-green-600 mt-2">
          Your keys have been set successfully.
        </p>
        <button
          class="btn mt-6 bg-[#119af5] text-white"
          @click="$router.push({ name: 'inbox' })"
        >
          Continue
        </button>
      </template>

      <template v-if="state === 'key-generation'">
        <p class="text-center text-green-600 mt-2">
          Your keys have been generated.
        </p>
        <button
          class="btn mt-6 bg-[#119af5] text-white"
          @click="$router.push({ name: 'inbox' })"
        >
          Continue
        </button>
        <button
          class="text-[#119af5] pt-4 pb-2 rounded-md font-semibold"
          @click="exportHandler"
        >
          Export Keys
        </button>
      </template>
    </template>
  </div>
</template>
