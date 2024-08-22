<script lang="ts" setup>
import { ref } from 'vue'

import axios from '@/plugins/axios'
import { decryptE2EPacket, createE2EPacket } from '@/cryptography/DiffieHellman'
import * as RSA from '@/cryptography/RSA'
import * as AES from '@/cryptography/AES'

import Time from '@/components/UI/Time.vue'
import Button from '@/components/UI/Button.vue'
import Textarea from '@/components/UI/Textarea.vue'

const props = defineProps<{
  message: {
    id: number
    content: string
    time: string
    owner: boolean
    mark: boolean
    can_replay: boolean
    session_id: number
    session_key?: string
    sender_public_key: string
    quote?: {
      id: number
      content: string
    }
  }
}>()

const replaying = ref(false)
const replayMessage = ref('')
const replaySent = ref(false)

const vDecrypt = {
  mounted: async (el: HTMLParagraphElement) => {
    try {
      window.Telegram.WebApp.CloudStorage.getItem(
        String(props.message.session_id),
        async (error, sessionKey) => {
          if (!sessionKey) {
            window.Telegram.WebApp.CloudStorage.getItem("private_key", async (error, value) => {
              const decryptedSessionKey = await RSA.decrypt(props.message.session_key!, value!)
              window.Telegram.WebApp.CloudStorage.setItem(String(props.message.session_id), decryptedSessionKey)

              try {
                const decryptedMsg = await AES.decrypt(props.message.content, decryptedSessionKey!)
                el.innerText = decryptedMsg!
              } catch (error) {
                alert(error)
                el.innerText = 'خطا در رمزگشایی!'
              }
            })
          } else {
            try {
              const decryptedMsg = await AES.decrypt(props.message.content, sessionKey!)
              el.innerText = decryptedMsg!
            } catch (error) {
              alert(error)
              el.innerText = 'خطا در رمزگشایی!'
            }
          }

        },
      )
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
  if (!replayMessage.value) return

  alert(replayMessage.value)

  window.Telegram.WebApp.CloudStorage.getItem(String(props.message.session_id), (error, sessionKey) => {
    alert(sessionKey)
    const encryptedMsg = AES.encrypt(replayMessage.value, sessionKey!)
    axios
      .post('/replay-message', {
        message_id: props.message.id,
        message: encryptedMsg,
      })
      .then(() => {
        replaying.value = false
        replayMessage.value = ''
        replaySent.value = true

        setTimeout(() => (replaySent.value = false), 1500)
      })
  })
}
</script>
<template>
  <div class="flex flex-col bg-[#ffffff] px-4 pt-3 pb-4 rounded-lg shadow-sm">
    <Time :value="message.time" class="text-gray-400 text-end text-sm"></Time>
    <p v-if="message.quote?.content" class="border-r-4 rounded-md border-r-blue-500 pr-2 py-2 mt-2 truncate w-full"
      style="background-color: rgba(137, 207, 240, 0.3)" quote="true" v-decrypt>
      {{ message.quote.content }}
    </p>

    <p class="break-words py-2" dir="auto" v-decrypt>{{ message.content }}</p>

    <template v-if="message.can_replay">
      <div v-if="!replaying" class="flex justify-end text-gray-400 text-end">
        <div class="flex items-center cursor-pointer" @click="replaying = true">
          <span class="ml-1 text-sm">پاسخ</span>
          <ReplyIcon size="20" color="#9CA38F" />
        </div>
      </div>

      <div v-else class="flex flex-col mt-4">
        <Textarea v-model="replayMessage" placeholder="پاسخ شما..." v-focus></Textarea>
        <Button :block="true" class="mt-4" @click="Submit">ارسال</Button>
        <p class="text-center pt-4 text-[#119af5] font-bold cursor-pointer" @click="replaying = false">
          بیخیال
        </p>
      </div>
      <p v-if="replaySent" class="text-center text-[#119af5]">
        پاسخ شما ارسال شد.
      </p>
    </template>
  </div>
</template>
