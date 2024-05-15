<script lang="ts" setup>
import { computed, ref } from 'vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateKeyPair } from '@/cryptography/RSA'
import { exportKeys, importKeysFromFile } from '@/utils'

import Modal from '@/components/UI/Modal.vue'
import SettingsIcon from '@/components/icons/Settings.vue'

const userStore = useUserStore()

const visible = ref(false)
const loading = ref(false)

const FileInput = ref<any>()

const hasKeys = computed(
  () =>
    localStorage.getItem('receive_private_key') &&
    localStorage.getItem('receive_public_key') &&
    localStorage.getItem('send_private_key') &&
    localStorage.getItem('send_public_key')
)

async function generateKeys() {
  loading.value = true

  let receivePublicKey: string, sendPublicKey: string

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    receivePublicKey = publicKey
    localStorage.setItem('receive_private_key', privateKey)
  })

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    sendPublicKey = publicKey
    localStorage.setItem('send_private_key', privateKey)
  })

  axios
    .post('/set-key', {
      send_public_key: sendPublicKey!,
      receive_public_key: receivePublicKey!,
    })
    .then(({ data }) => {
      userStore.user.sendPublicKey = data.send_public_key
      userStore.user.receivePublicKey = data.receive_public_key

      localStorage.setItem('send_public_key', data.send_public_key)
      localStorage.setItem('receive_public_key', data.receive_public_key)
    })
    .finally(() => {
      loading.value = false
    })
}

function importHandler(event: Event) {
  // @ts-ignore
  const file = event.target!.files[0] as File

  importKeysFromFile(file).then((keys) => {
    localStorage.setItem('send_private_key', keys.sendPrivateKey)
    localStorage.setItem('receive_private_key', keys.receivePrivateKey)
    axios
      .post('/set-key', {
        send_public_key: keys.sendPublicKey,
        receive_public_key: keys.receivePublicKey,
      })
      .then(({ data }) => {
        userStore.user.sendPublicKey = data.send_public_key
        userStore.user.receivePublicKey = data.receive_public_key

        localStorage.setItem('send_public_key', data.send_public_key)
        localStorage.setItem('receive_public_key', data.receive_public_key)
      })
  })
}
</script>

<template>
  <SettingsIcon size="24" color="black" @click="visible = true" />

  <Modal v-model="visible">
    <template #header>
      <h3 class="font-bold text-lg">Security Keys</h3>
    </template>
    <template #body>
      <template v-if="!loading">
        <div v-if="hasKeys">
          <p class="pt-5 pb-10 text-center text-green-600">
            Your keys are set and you are safe.
          </p>
          <div class="grid grid-cols-1 gap-y-2">
            <div class="grid grid-cols-2 gap-x-2">
              <button
                class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
                @click="exportKeys"
              >
                Export
              </button>
              <button
                class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
                @click="FileInput.click()"
              >
                Import
              </button>
              <input
                ref="FileInput"
                type="file"
                class="hidden"
                @change="importHandler"
              />
            </div>
            <button
              class="text-[#119af5] py-2 rounded-md font-semibold"
              @click="generateKeys"
            >
              I want to regenerate my keys.
            </button>
          </div>
        </div>
        <template v-else>
          <p class="py-10 text-center text-red-600">
            You currently have no keys!
          </p>
          <div class="grid grid-cols-1 gap-y-2">
            <button
              class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
              @click="generateKeys"
            >
              Generate
            </button>
            <button class="text-[#119af5] py-2 rounded-md font-semibold">
              I Want to import my keys.
            </button>
          </div>
        </template>
      </template>
      <div v-else>
        <p class="flex items-center justify-center py-10">
          Generating your keys...
          <span class="loading loading-infinity loading-md mx-2"></span>
        </p>
      </div>
    </template>
  </Modal>
</template>
