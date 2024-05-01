<script lang="ts" setup>
import { computed, ref } from 'vue'
import { RouterView } from 'vue-router'
import forge from 'node-forge'
import { useAuthStore } from '@/stores/auth'

import CopyText from '@/components/CopyText.vue'
import SettingsIcon from '@/components/icons/Settings.vue'

const authStore = useAuthStore()

const hasKeys = ref(
  !!localStorage.getItem('private_key') && !!localStorage.getItem('public_key')
)
const loading = ref(false)

const myLink = computed(() => `${location.origin}/@${authStore.user.Username}`)

async function generateKeyPair() {
  const keyPair = await forge.pki.rsa.generateKeyPair({
    bits: 512,
    workers: 2,
  })

  const publicKey = forge.pki.publicKeyToPem(keyPair.publicKey)
  const privateKey = forge.pki.privateKeyToPem(keyPair.privateKey)

  localStorage.setItem('private_key', privateKey)
  localStorage.setItem('public_key', publicKey)

  hasKeys.value = true

  loading.value = true

  setTimeout(() => {
    loading.value = false
  }, 2000)
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
  <div>
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
            <button
              class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
              @click="exportKeys"
            >
              Export
            </button>
            <button
              class="text-[#119af5] py-2 rounded-md font-semibold"
              @click="generateKeyPair"
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
              @click="generateKeyPair"
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
