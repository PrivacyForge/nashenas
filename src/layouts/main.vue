<script lang="ts" setup>
import { computed, ref } from 'vue'
import { RouterView } from 'vue-router'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateKeyPair } from '@/cryptography/RSA'

import CopyText from '@/components/CopyText.vue'
import SettingsIcon from '@/components/icons/Settings.vue'
import { extractKeys } from '@/utils'

const userStore = useUserStore()

const FileInput = ref<any>()

const hasKeys = computed(
  () =>
    localStorage.getItem('receive_private_key') &&
    localStorage.getItem('receive_public_key') &&
    localStorage.getItem('send_private_key') &&
    localStorage.getItem('send_public_key')
)

const loading = ref(false)

const myLink = computed(() => `${location.origin}/@${userStore.user.username}`)

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

  axios
    .post('/set-key', {
      receive_public_key: receivePublicKey!,
      send_public_key: sendPublicKey!,
    })
    .then(({ data }) => {
      userStore.user.receivePublicKey = data.receive_public_key
      userStore.user.sendPublicKey = data.send_public_key
    })
    .finally(() => {
      loading.value = false
    })
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
        userStore.user.receivePublicKey = data.receive_public_key
        userStore.user.sendPublicKey = data.send_public_key
      })
  }
  reader.readAsText(file)
}

function exportKeys() {
  const link = document.createElement('a')

  const content = `${localStorage.getItem(
    'private_key'
  )}divide\n${localStorage.getItem('public_key')}`

  console.log(content.split('divide'))

  const file = new Blob([content], { type: 'text/plain' })

  link.href = URL.createObjectURL(file)

  link.download = 'keys.txt'

  link.click()
  URL.revokeObjectURL(link.href)
}
</script>

<template>
  <div class="lg:w-5/12 mx-auto">
    <nav class="flex justify-between bg-[#ffffff] p-5 m-4 rounded-lg shadow-sm">
      <div class="grid grid-cols-2 gap-x-3">
        <SettingsIcon size="24" color="black" onclick="modal_1.showModal()" />
      </div>
      <CopyText
        text="My Link"
        :copy="myLink"
        class="text-[#119af5] test- font-semibold text-end"
      />
    </nav>
    <hr class="my-1 mx-4" />
    <RouterView />
  </div>

  <dialog id="modal_1" class="modal">
    <div class="modal-box">
      <form method="dialog">
        <button
          class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2 outline-none"
        >
          ✕
        </button>
      </form>
      <h3 class="font-bold text-lg">Security Keys</h3>
      <template v-if="!loading">
        <div v-if="hasKeys">
          <p class="py-10 text-center text-green-600">
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
                @change="importKeys"
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
    </div>
  </dialog>
</template>
