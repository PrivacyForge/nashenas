<script lang="ts" setup>
import { decrypt } from '@/cryptography'

import ReplyIcon from '@/components/icons/Reply.vue'

defineProps<{
  text: string
  time: string
  mark: boolean
}>()

const vDecrypt = {
  mounted: async (el: HTMLParagraphElement) => {
    try {
      const privateKey = localStorage.getItem('private_key')

      const decrypted = await decrypt(el.innerText, privateKey!)

      el.innerText = decrypted
    } catch (error) {
      el.innerText = 'Decryption Error.'
    }
  },
}
</script>
<template>
  <div
    class="flex flex-col bg-[#ffffff] p-4 rounded-lg shadow-sm"
    :class="mark && ['border-2 border-[#119af5]']"
  >
    <p class="text-gray-400 text-end">{{ time.split('.')[0] }}</p>
    <p class="break-words pt-2" dir="auto" v-decrypt>{{ text }}</p>
    <div class="flex justify-end text-gray-400 text-end">
      <span class="mr-1">Reply</span> <ReplyIcon size="20" color="#9CA38F" />
    </div>
  </div>
</template>
