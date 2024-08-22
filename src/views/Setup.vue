<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import { generateKeyPair } from '@/cryptography/RSA'

import Card from '@/components/UI/Card.vue'
import Input from '@/components/UI/Input.vue'
import Button from '@/components/UI/Button.vue'
import LoadingIcon from '@/components/icons/Loading.vue'

const userStore = useUserStore()
const route = useRoute()
const router = useRouter()

const username = ref(userStore.user.username)
const usernameErr = ref('')

const loading = ref(false)
const delay = ref(false)

const state = ref<
  'set-username' | 'key-question' | 'key-generation' | 'key-upload'
>('set-username')

async function generateKeyPairs() {
  loading.value = true

  generateKeyPair().then(({ privateKey, publicKey }) => {
    window.Telegram.WebApp.CloudStorage.setItem("private_key", privateKey)
    window.Telegram.WebApp.CloudStorage.setItem("public_key", publicKey)
    setTimeout(() => {
      axios
        .post('/set-key', {
          public_key: publicKey,
        })
        .then(({ data }) => {
          state.value = 'key-generation'
          userStore.user.publicKey = data.public_key
        })
        .finally(() => {
          loading.value = false
        })
    }, 1500)
  })
}

function exportHandler() {
  window.Telegram.WebApp.CloudStorage.getItem("private_key", async (error, privateKey) => {
    window.Telegram.WebApp.CloudStorage.getItem("public_key", async (error, publicKey) => {
      navigator.clipboard.writeText(`${privateKey}\n\n\n${publicKey}`)
      delay.value = true
      setTimeout(() => {
        delay.value = false
      }, 2000);
    })
  })
  state.value = 'key-generation'
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

function done() {
  if (route.query.next)
    router.push({
      name: 'profile',
      params: { usernameWithHash: route.query.next as string },
    })
  else {
    router.push({ name: 'inbox' })
  }
}

onMounted(() => {
  window.Telegram.WebApp.expand()
})
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
            یه نام کاربر برای خودت انتخاب کن:
          </label>
          <Input v-model="username" class="pl-7 my-4" placeholder="Username..." dir="ltr" />
          <span class="absolute left-2 top-[54px] font-bold">@</span>
        </div>
        <p v-if="usernameErr" class="text-red-500 my-2" v-text="usernameErr" />
        <Button @click="usernameSubmit()"> بعدی </Button>
        <p class="mt-5">
          اگه قصد داری بروزرسانی ناشناس باشه یه مقدار تصادفی رو وارد کن.
        </p>
      </template>

      <template v-if="state === 'key-question'">
        <p class="pb-5 pt-3 text-center">
          توی این مرحله، ما برات یه جفت کلید RSA برای رمزنگاری پیام‌های ارسالی و
          دریافتی می‌سازیم. در آینده امکان اینکه خودت کلیدت رو بسازی و وارد کنی
          هم اضافه میشه. اینم یادت باشه که کلید private تو هیچوقت سمت هیچ سروری
          ارسال نمیشه و تمام فرایند رمزنگاری سمت تلگرام انجام میشه.
        </p>
        <div class="grid grid-cols-1 gap-y-2">
          <button class="bg-[#119af5] text-white py-2 rounded-md font-semibold" @click="generateKeyPairs">
            متوجه شدم
          </button>
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
          کلیدهای رمزنگاری تو با موفقیت ساخته شد.
        </p>
        <Button @click="done()"> ادامه </Button>
        <p class="text-center mt-4 text-[#119af5] font-semibold cursor-pointer" @click="exportHandler">
          {{ !delay ? 'کپی کردن' : 'کپی شد' }}
        </p>
      </template>
    </template>
  </Card>
</template>
