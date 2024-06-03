<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { onMounted, reactive, ref } from 'vue'

import axios from '@/plugins/axios'
import { encrypt } from '@/cryptography'
import { useUserStore } from '@/stores/user'

import Card from '@/components/UI/Card.vue'
import GithubIcon from '@/components/icons/Github.vue'
import LoadingIcon from '@/components/icons/Loading.vue'
import TelegramIcon from '@/components/icons/Telegram.vue'
import { ArrowRightIcon } from '@heroicons/vue/24/outline'
import { PaperAirplaneIcon } from '@heroicons/vue/24/solid'

import StickerImg from '@/assets/images/sticker.webp'

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
  <div class="fixed top-0 right-0 left-0 h-16 bg-white shadow z-20">
    <div class="max-w-lg mx-auto flex items-center h-full">
      <RouterLink class="w-5 h-5" :to="{ name: 'inbox' }">
        <ArrowRightIcon class="w-full h-full" />
      </RouterLink>
      <p class="font-bold mt-0.5 mr-4">{{ username }}</p>

      <a
        class="mr-auto"
        href="https://github.com/PrivacyForge/nashenas"
        target="_blank"
      >
        <GithubIcon class="ml-1" size="26" color="gray" />
      </a>
      <a href="https://nashenase2ebot.t.me" target="_blank">
        <TelegramIcon class="mr-1" size="24" color="gray" />
      </a>
    </div>
  </div>

  <div class="pt-16">
    <Card v-if="loading" class="grid grid-cols-1 max-w-lg mx-auto gap-4 m-4">
      <div class="flex justify-center items-center py-4">
        <p class="text-[#119af5] font-semibold">درحال بارگذاری...</p>
        <LoadingIcon color="#119af5" size="26px" />
      </div>
    </Card>

    <div v-else class="max-w-lg mx-auto pt-10">
      <template v-if="!sent">
        <template v-if="notFoundUser">
          <p
            class="text-center text-[#119af5] font-semibold"
            v-text="errorMessage"
          />
        </template>
        <template v-else>
          <img class="w-1/2 mx-auto" :src="StickerImg" alt="sticker" />
          <p class="text-center mt-2 text-gray-600 text-sm">
            <span>درحال ارسال پیام ناشناس به </span>
            <span>{{ username }}</span>
            <span> هستی...</span>
          </p>
        </template>
      </template>
      <template v-else>
        <Card>
          <p class="text-center text-[#119af5] font-semibold">
            پیام شما رمزنگاری و به کاربر ارسال شد.
          </p>
        </Card>
      </template>
    </div>
  </div>

  <div class="fixed bottom-0 right-0 left-0 h-16 bg-white z-20 shadow-lg">
    <form
      @submit.prevent="submit"
      class="max-w-lg h-full mx-auto flex items-center"
    >
      <input
        v-model="message"
        class="h-full grow !outline-none"
        placeholder="پیام شما..."
      />

      <button
        :disabled="message.length < 3 || submitLoading"
        class="w-5 h-5"
        type="submit"
      >
        <PaperAirplaneIcon
          v-if="!submitLoading"
          class="w-full h-full rotate-180"
          :class="message.length < 3 ? 'text-gray-400' : 'text-blue-500'"
        />
        <LoadingIcon v-else color="rgb(59, 130, 246)" size="20px" />
      </button>
    </form>
  </div>
</template>
