<script lang="ts" setup>
import { computed, ref } from 'vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateKeyPair } from '@/cryptography/RSA'
import { exportKeys, importKeysFromFile } from '@/utils'

import Modal from '@/components/UI/Modal.vue'
import LoadingIcon from '@/components/icons/Loading.vue'

import LockOpenIcon from '@/components/icons/LockOpen.vue'
import LockCloseIcon from '@/components/icons/LockClose.vue'

const userStore = useUserStore()

const visible = ref(false)
const loading = ref(false)

const FileInput = ref<any>()

const hasKeys = true

async function generateKeys() {
  loading.value = true

  let receivePublicKey: string, sendPublicKey: string

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    receivePublicKey = publicKey
    window.Telegram.WebApp.CloudStorage.setItem("receive_private_key", privateKey)

  })

  await generateKeyPair().then(({ privateKey, publicKey }) => {
    sendPublicKey = publicKey
    window.Telegram.WebApp.CloudStorage.setItem("send_private_key", privateKey)
  })

  axios
    .post('/set-key', {
      send_public_key: sendPublicKey!,
      receive_public_key: receivePublicKey!,
    })
    .then(({ data }) => {
      userStore.user.publicKey = data.public_key

      window.Telegram.WebApp.CloudStorage.setItem("send_public_key", data.send_public_key)
      window.Telegram.WebApp.CloudStorage.setItem("receive_public_key", data.receive_public_key)
    })
    .finally(() => {
      loading.value = false
    })
}

function importHandler(event: Event) {
  // @ts-ignore
  const file = event.target!.files[0] as File

  importKeysFromFile(file).then((keys) => {
    window.Telegram.WebApp.CloudStorage.setItem("send_private_key", keys.sendPrivateKey)
    window.Telegram.WebApp.CloudStorage.setItem("receive_private_key", keys.receivePrivateKey)

    axios
      .post('/set-key', {
        send_public_key: keys.sendPublicKey,
        receive_public_key: keys.receivePublicKey,
      })
      .then(({ data }) => {
        userStore.user.publicKey = data.public_key

        window.Telegram.WebApp.CloudStorage.setItem("send_public_key", data.send_public_key)
        window.Telegram.WebApp.CloudStorage.setItem("receive_public_key", data.receive_public_key)
      })
  })
}
</script>

<template>
  <LockCloseIcon v-if="hasKeys" size="30" color="#4BB543" />
  <LockOpenIcon v-else size="30" color="#FF5733" />

  <Modal v-model="visible">
    <template #header>
      <h3 class="font-bold text-lg">کلیدهای امنیتی</h3>
    </template>
    <template #body>
      <input ref="FileInput" type="file" class="hidden" @change="importHandler" />

      <template v-if="!loading">
        <div v-if="hasKeys">
          <p class="pt-5 pb-10 text-center text-green-600">
            کلیدهای شما پیکربندی شده است.
          </p>
          <div class="grid grid-cols-1 gap-y-2">
            <div class="grid grid-cols-2 gap-x-2">
              <button class="bg-[#119af5] text-white py-2 rounded-md font-semibold" @click="exportKeys">
                دانلود
              </button>
              <button class="bg-[#119af5] text-white py-2 rounded-md font-semibold" @click="FileInput.click()">
                آپلود
              </button>
            </div>
            <button class="text-[#119af5] py-2 rounded-md font-semibold" @click="generateKeys">
              می‌خواهم مجددا جفت کلید بسازم.
            </button>
          </div>
        </div>
        <template v-else>
          <p class="py-10 text-center text-red-600">
            درحال حاضر هیچ کلیدی ندارید.
          </p>
          <div class="grid grid-cols-1 gap-y-2">
            <button class="bg-[#119af5] text-white py-2 rounded-md font-semibold" @click="generateKeys">
              ساختن
            </button>
            <button class="text-[#119af5] py-2 rounded-md font-semibold" @click="FileInput.click()">
              آپلود می‌کنم.
            </button>
          </div>
        </template>
      </template>
      <div v-else>
        <p class="flex items-center justify-center py-10">
          درحال ساخت جفت کلید...
          <LoadingIcon color="#119af5" size="26px" />
        </p>
      </div>
    </template>
  </Modal>
</template>
