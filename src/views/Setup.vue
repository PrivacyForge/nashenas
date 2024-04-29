<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import forge from 'node-forge'

const route = useRoute()

const loading = ref(false)
const FileInput = ref<any>()

const state = ref<
  'set-username' | 'key-question' | 'key-generation' | 'key-upload'
>('key-question')

async function generateKeyPair() {
  const keyPair = await forge.pki.rsa.generateKeyPair({
    bits: 512,
    workers: 2,
  })

  const publicKey = forge.pki.publicKeyToPem(keyPair.publicKey)
  const privateKey = forge.pki.privateKeyToPem(keyPair.privateKey)

  localStorage.setItem('private_key', privateKey)
  localStorage.setItem('public_key', publicKey)

  loading.value = true

  setTimeout(() => {
    loading.value = false
  }, 2000)

  state.value = 'key-generation'
}

function exportKeys() {
  const link = document.createElement('a')

  const content = `${localStorage.getItem(
    'private_key'
  )}divide\n${localStorage.getItem('public_key')}`

  const file = new Blob([content], { type: 'text/plain' })

  link.href = URL.createObjectURL(file)

  link.download = 'keys.txt'

  link.click()
  URL.revokeObjectURL(link.href)
  state.value = 'key-generation'
}

function importKeys(event: Event) {
  // @ts-ignore
  const file = event.target!.files[0]
  const reader = new FileReader()
  reader.onload = (e) => {
    const rawData = e.target!.result as string
    const privateKey = rawData.split('divide')[0]
    const publicKey = rawData.split('divide')[1].slice(1) // send to server

    localStorage.setItem('private_key', privateKey)
  }
  reader.readAsText(file)
  state.value = 'key-upload'
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
            type="text"
            class="input input-bordered w-full focus:outline-[#119af5] outline-[#119af5] mt-3 pl-6"
            placeholder="username"
          />
          <span class="absolute left-2 top-12">@</span>
        </div>
        <button
          class="btn mt-4 bg-[#119af5] text-white"
          @click="state = 'key-question'"
        >
          Next
        </button>
      </template>
      <template v-if="state === 'key-question'">
        <p class="pb-5 pt-3 text-center">Do you have any public/private key?</p>
        <div class="grid grid-cols-1 gap-y-2">
          <button
            class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
            @click="generateKeyPair"
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
        <button class="btn mt-6 bg-[#119af5] text-white">Continue</button>
      </template>

      <template v-if="state === 'key-generation'">
        <p class="text-center text-green-600 mt-2">
          Your keys have been generated.
        </p>
        <button class="btn mt-6 bg-[#119af5] text-white">Continue</button>
        <button
          class="text-[#119af5] pt-4 pb-2 rounded-md font-semibold"
          @click="exportKeys"
        >
          Export Keys
        </button>
      </template>
    </template>
  </div>
</template>
