<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'

import axios from '@/plugins/axios'
import { encrypt } from '@/cryptography'
import { useUserStore } from '@/stores/user'

import Card from '@/components/UI/Card.vue'
import Button from '@/components/UI/Button.vue'
import Textarea from '@/components/UI/Textarea.vue'
import GithubIcon from '@/components/icons/Github.vue'
import LoadingIcon from '@/components/icons/Loading.vue'
import TelegramIcon from '@/components/icons/Telegram.vue'

const userStore = useUserStore()

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

async function submit() {
  if (message.value === '') return

  let myPublicKey
  if (userStore.isAuth) {
    myPublicKey = localStorage.getItem('send_public_key')
  }

  const encryptedMsg = await encrypt(
    message.value,
    user.publicKey!,
    myPublicKey!,
  )

  submitLoading.value = true

  setTimeout(() => {
    axios
      .post(`/send-message`, {
        message: encryptedMsg,
        id: user.id,
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
  setTimeout(async () => {
    await axios
      .get(`/profile/${username}`)
      .then(async (response) => {
        user.publicKey = response.data.public_key
        user.id = response.data.id

        await axios
          .get('/me')
          .then(({ data }) => {
            userStore.user.id = data.id
            userStore.user.userid = data.userid
            userStore.user.username = data.username
            userStore.user.receivePublicKey = data.receive_public_key
            userStore.user.sendPublicKey = data.send_public_key

            userStore.isAuth = true
          })
          .catch(() => {})
      })
      .catch(() => {
        notFoundUser.value = true
        errorMessage.value = 'Not found user.'
      })
      .finally(() => {
        loading.value = false
      })
  }, 1000)
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
          <p
            class="text-center text-[#119af5] font-semibold"
            v-text="errorMessage"
          />
        </template>
        <template v-else>
          <p>برای کاربر {{ $route.params.username }} یه پیام ناشناس بنویس...</p>
          <Textarea v-model="message" placeholder="متن..."></Textarea>
          <Button
            :disabled="message.length < 3"
            :loading="submitLoading"
            @click="submit"
            >ارسال پیام</Button
          >
        </template>
      </template>
      <template v-else>
        <p class="text-center text-[#119af5] font-semibold">
          پیام شما رمزنگاری و به کاربر ارسال شد.
        </p>
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
