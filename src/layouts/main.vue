<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'

import Settings from '@/components/Settings.vue'
import CopyText from '@/components/UI/CopyText.vue'

import { bufferToHex } from '@/utils'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const delay = ref(false)

async function copy() {
  try {
    const encoder = new TextEncoder();
    const data = encoder.encode(userStore.user.publicKey);

    const hashBuffer = await crypto.subtle.digest("SHA-1", data);
    const hash = bufferToHex(hashBuffer)
    navigator.clipboard.writeText(`https://t.me/Nashenase2ebot?start=${userStore.user}-${hash}`)
  } catch (error) {
    alert(error)
  }

  delay.value = true
  setTimeout(() => {
    delay.value = false
  }, 2000)
}

onMounted(() => {
  window.Telegram.WebApp.expand()
  window.Telegram.WebApp.disableVerticalSwipes()
})
</script>

<template>
  <div class="lg:w-5/12 mx-auto px-4">
    <nav class="flex justify-between bg-[#ffffff] px-5 py-4 my-4 rounded-lg shadow-sm">
      <div class="cursor-pointer">
        <Settings />
      </div>
      <button v-if="!delay" class="text-[#119af5] font-semibold text-end" @click="copy">کپی لینک</button>
      <div class="text-center" v-else>کپی شد!</div>
    </nav>
    <hr class="my-1" />
    <RouterView />
  </div>
</template>
