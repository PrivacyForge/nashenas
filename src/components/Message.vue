<script lang="ts" setup>
import { ref } from 'vue'

import axios from '@/plugins/axios'
import { decryptE2EPacket, createE2EPacket } from '@/cryptography/DiffieHellman'

import Time from '@/components/UI/Time.vue'
import Button from '@/components/UI/Button.vue'
import Textarea from '@/components/UI/Textarea.vue'

const props = defineProps<{
  id: number
  text: string
  time: string
  owner: boolean
  mark: boolean
  canReply: boolean
  sender_public_key: string
  quote?: {
    id: number
    content: string
  }
}>()

const replying = ref(false)
const replyMessage = ref('')
const replySent = ref(false)

const vDecrypt = {
  mounted: async (el: HTMLParagraphElement) => {
    try {
      const isQuote = !!el.getAttribute('quote')
      if (props.owner) {
        window.Telegram.WebApp.CloudStorage.getItem(
          'receive_private_key',
          async (error, privateKey) => {
            const decryptedMsg = await decryptE2EPacket(
              privateKey!,
              props.sender_public_key,
              el.innerText,
            )
            el.innerText = decryptedMsg!
          },
        )
      } else {
        window.Telegram.WebApp.CloudStorage.getItem(
          'send_private_key',
          async (error, privateKey) => {
            const decryptedMsg = await decryptE2EPacket(
              privateKey!,
              props.sender_public_key,
              el.innerText,
            )
            el.innerText = decryptedMsg!
          },
        )
      }
    } catch (error) {
      alert(error)
      el.innerText = 'خطا در رمزگشایی!'
    }
  },
}

const vFocus = {
  mounted: (el: HTMLParagraphElement) => {
    el.focus()
  },
}

function Submit() {
  if (!replyMessage.value) return

  axios.get(`/get-key/${props.id}`).then(async ({ data: key }) => {
    window.Telegram.WebApp.CloudStorage.getItem(
      props.owner ? 'receive_private_key' : 'send_private_key',
      async (error, privateKey) => {
        const encryptedMsg = await createE2EPacket(
          key,
          privateKey!,
          replyMessage.value,
        )
        axios
          .post('/reply-message', {
            message_id: props.id,
            message: encryptedMsg,
          })
          .then(() => {
            replying.value = false
            replyMessage.value = ''
            replySent.value = true

            setTimeout(() => (replySent.value = false), 1500)
          })
      },
    )
  })
}
</script>
<template>
  <div
    class="flex flex-col bg-[#ffffff] px-4 pt-3 pb-4 rounded-lg shadow-sm"
    :class="mark && ['border-2 border-[#119af5]']"
  >
    <Time :value="time" class="text-gray-400 text-end text-sm"></Time>
    <p
      v-if="quote?.content"
      class="border-r-4 rounded-md border-r-blue-500 pr-2 py-2 mt-2 truncate w-full"
      style="background-color: rgba(137, 207, 240, 0.3)"
      quote="true"
      v-decrypt
    >
      {{ quote.content }}
    </p>

    <p class="break-words py-2" dir="auto" v-decrypt>{{ text }}</p>

    <template v-if="canReply">
      <div v-if="!replying" class="flex justify-end text-gray-400 text-end">
        <div class="flex items-center cursor-pointer" @click="replying = true">
          <span class="ml-1 text-sm">پاسخ</span>
          <ReplyIcon size="20" color="#9CA38F" />
        </div>
      </div>

      <div v-else class="flex flex-col mt-4">
        <Textarea
          v-model="replyMessage"
          placeholder="پاسخ شما..."
          v-focus
        ></Textarea>
        <Button :block="true" class="mt-4" @click="Submit">ارسال</Button>
        <p
          class="text-center pt-4 text-[#119af5] font-bold cursor-pointer"
          @click="replying = false"
        >
          بیخیال
        </p>
      </div>
      <p v-if="replySent" class="text-center text-[#119af5]">
        پاسخ شما ارسال شد.
      </p>
    </template>
  </div>
</template>
