<script lang="ts" setup>
import { ref } from 'vue'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { extractKeys, exportKeys } from '@/utils'
import { generateKeyPair } from '@/cryptography/RSA'

import Card from '@/components/UI/Card.vue'
import Input from '@/components/UI/Input.vue'
import Button from '@/components/UI/Button.vue'
import LoadingIcon from '@/components/icons/Loading.vue'

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
    .then(({ data }) => {
      state.value = 'key-question'
      userStore.user.username = data.username
    })
    .catch((error) => {
      usernameErr.value = error.response.data.message
    })
}
</script>
<template>
  <Card class="grid grid-cols-1 m-4">
    <div v-if="loading" class="flex justify-center items-center py-6">
      <p class="text-[#119af5] font-semibold">درحال ساخت جفت کلید...</p>
      <LoadingIcon color="#119af5" size="26px" />
    </div>

    <template v-else>
      <template v-if="state === 'set-username'">
        <div class="relative">
          <label class="font-semibold">
            یک نام کاربر برای خود انتخاب کنید:
          </label>
          <Input
            v-model="username"
            class="pl-7 my-4"
            placeholder="Username..."
            dir="ltr"
          />
          <span class="absolute left-2 top-[54px] font-bold">@</span>
        </div>
        <p v-if="usernameErr" class="text-red-500 mt-2" v-text="usernameErr" />
        <Button @click="usernameSubmit()"> بعدی </Button>
      </template>

      <template v-if="state === 'key-question'">
        <p class="pb-5 pt-3 text-center">
          آیا از قبل جفت کلید عمومی و خصوصی دارید؟
        </p>
        <div class="grid grid-cols-1 gap-y-2">
          <button
            class="bg-[#119af5] text-white py-2 rounded-md font-semibold"
            @click="generateKeys"
          >
            ندارم لطفا بساز
          </button>
          <button
            class="text-[#119af5] py-2 rounded-md font-semibold"
            @click="FileInput.click()"
          >
            بله آپلود می‌کنم
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
        <p class="text-center text-green-600 mt-2 mb-4">
          کلیدهای شما با موفقیت پیکربندی شدند.
        </p>
        <Button @click="$router.push({ name: 'inbox' })"> ادامه </Button>
      </template>

      <template v-if="state === 'key-generation'">
        <p class="text-center text-green-600 mt-2 mb-4">
          کلیدهای شما ساخته شد.
        </p>
        <Button @click="$router.push({ name: 'inbox' })"> ادامه </Button>
        <p
          class="text-center mt-4 text-[#119af5] font-semibold"
          @click="exportHandler"
        >
          دانلود کلیدها
        </p>
      </template>
    </template>
  </Card>
</template>
