<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateRandomKey as generateRandomAESKey } from '@/cryptography/AES'
import * as RSA from '@/cryptography/RSA'
import * as AES from '@/cryptography/AES'

import Card from '@/components/UI/Card.vue'
import Button from '@/components/UI/Button.vue'
import Textarea from '@/components/UI/Textarea.vue'
import GithubIcon from '@/components/icons/Github.vue'
import LoadingIcon from '@/components/icons/Loading.vue'
import TelegramIcon from '@/components/icons/Telegram.vue'
import { bufferToHex } from '@/utils'

const userStore = useUserStore()

const route = useRoute()
const router = useRouter()

const user = reactive<{
  id: number | null
  publicKey: string | null
  username: string | null
}>({
  id: null,
  username: null,
  publicKey: null,
})

const message = ref('')

const loading = ref(true)
const notFoundUser = ref<boolean>()
const errorMessage = ref('')

const submitLoading = ref(false)

const sent = ref(false)

async function submit() {
  alert(message.value)
  if (message.value === '') return

  try {
    const sessionKey = generateRandomAESKey()

    alert(sessionKey)

    const encryptedMsg = await AES.encrypt(message.value, sessionKey)
    alert(encryptedMsg)
    alert(user.publicKey)
    const encryptedKey = await RSA.encrypt(sessionKey, user.publicKey!)
    alert(encryptedKey)

    submitLoading.value = true
    setTimeout(() => {
      axios
        .post(`/send-message`, {
          message: encryptedMsg,
          session_key: encryptedKey,
          id: user.id!
        })
        .then(({ data }) => {
          message.value = ''
          sent.value = true

          window.Telegram.WebApp.CloudStorage.setItem(data.session_id, sessionKey)
        })
        .finally(() => {
          submitLoading.value = false
        })
    }, 2000)
  } catch (error) {
    alert(error)
  }
}

onMounted(async () => {
  const words = (route.params.usernameWithHash as string).split("-")

  if (words.length !== 2) router.push({ name: 'error' })

  const username = words[0].slice(1)
  const hash = words[1]

  alert(route.params.usernameWithHash)
  alert(username)
  alert(hash)

  loading.value = true
  axios
    .get('/me')
    .then(({ data }) => {
      userStore.user.id = data.id
      userStore.user.userid = data.userid
      userStore.user.username = data.username
      userStore.user.publicKey = data.public_key
      userStore.isAuth = true

      if (!userStore.user.username) router.push({ name: "setup", query: { next: hash } })

      axios
        .get(`/profile/${username}`)
        .then(async (response) => {
          user.id = response.data.id
          user.publicKey = response.data.public_key
          user.username = response.data.username
          alert(response.data.public_key)
          alert(user.publicKey)

          const encoder = new TextEncoder();
          const data = encoder.encode(user.publicKey!);

          const hashBuffer = await crypto.subtle.digest("SHA-256", data);

          const publicKeyHash = bufferToHex(hashBuffer)

          if (publicKeyHash !== hash) {
            router.push({ name: "error" })
          }
        })
        .catch((err: any) => {
          alert(err)
        })
        .finally(() => {
          loading.value = false
        })
    })
})
</script>

<template>
  <Card class="grid grid-cols-1 lg:w-4/12 lg:mx-auto gap-4 m-4">
    <div v-if="loading" class="flex justify-center items-center py-4">
      <p class="text-[#119af5] font-semibold">درحال بارگذاری...</p>
      <LoadingIcon color="#119af5" size="26px" />
    </div>

    <template v-else>
      <template v-if="!sent">
        <template v-if="notFoundUser">
          <p class="text-center text-[#119af5] font-semibold" v-text="errorMessage" />
        </template>
        <template v-else>
          <p>برای کاربر مقصدت یه پیام ناشناس بنویس...</p>
          <Textarea v-model="message" placeholder="متن..."></Textarea>
          <Button :disabled="message.length < 3" :loading="submitLoading" @click="submit">ارسال پیام</Button>
        </template>
      </template>
      <template v-else>
        <p class="text-center font-semibold">
          پیام شما رمزنگاری و به کاربر ارسال شد.
        </p>
        <router-link class="text-[#119af5] text-center" to="/inbox">صفحه اصلی</router-link>
      </template>
    </template>
  </Card>

  <div class="fixed bottom-2 left-[50%] w-max translate-x-[-50%] text-gray-500">
    <div class="flex justify-center items-center pb-2">
      <a href="https://github.com/PrivacyForge/nashenas" target="_blank">
        <GithubIcon class="ml-1" size="26" color="gray" />
      </a>
      <a href="https://nashenase2ebot.t.me" target="_blank">
        <TelegramIcon class="mr-1" size="24" color="gray" />
      </a>
    </div>
    <p dir="ltr">Open Source + E2E Encryption + Safe</p>
  </div>
</template>
