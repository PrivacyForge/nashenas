<script lang="ts" setup>
import { ref } from 'vue'

import axios from '@/plugins/axios'
import { decrypt, encrypt } from '@/cryptography'

import ReplyIcon from '@/components/icons/Reply.vue'
import Button from '@/components/Button.vue'
import Time from '@/components/Time.vue'

const props = defineProps<{
  id: number
  text: string
  time: string
  owner: boolean
  mark: boolean
  quote?: {
    id: number
    content: string
  }
}>()

const replaying = ref(false)
const replayMessage = ref('')
const replaySent = ref(false)

const vDecrypt = {
  mounted: async (el: HTMLParagraphElement) => {
    try {
      let privateKey, decryptedMsg
      if (props.owner) {
        privateKey = localStorage.getItem('receive_private_key')
        decryptedMsg = await decrypt(
          el.innerText,
          privateKey!,
          el.getAttribute('quote') == 'true'
        )
      } else {
        privateKey = localStorage.getItem('send_private_key')
        decryptedMsg = await decrypt(
          el.innerText,
          privateKey!,
          el.getAttribute('quote') == 'true'
        )
      }

      el.innerText = decryptedMsg
    } catch (error) {
      el.innerText = 'Decryption Error.'
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

  axios.get(`/get-key/${props.id}`).then(async ({ data: key }) => {
    let publicKey, encryptedMsg
    if (props.owner) {
      publicKey = localStorage.getItem('receive_public_key')
      encryptedMsg = await encrypt(replayMessage.value, key, publicKey!)
    } else {
      publicKey = localStorage.getItem('send_public_key')
      encryptedMsg = await encrypt(replayMessage.value, key, publicKey!)
    }

    axios
      .post('/replay-message', {
        message_id: props.id,
        message: encryptedMsg,
      })
      .then(() => {
        replaying.value = false

        replaySent.value = true

        setTimeout(() => (replaySent.value = false), 1500)
      })
  })
}
</script>
<template>
  <div
    class="flex flex-col bg-[#ffffff] px-4 pt-3 pb-4 rounded-lg shadow-sm"
    :class="mark && ['border-2 border-[#119af5]']"
  >
    <Time :value="time" class="text-gray-400 text-end"></Time>
    <p
      v-if="quote?.content"
      class="border-l-4 rounded-md border-l-blue-500 pl-2 py-2 mt-2"
      style="background-color: rgba(137, 207, 240, 0.3)"
      quote="true"
      v-decrypt
    >
      {{ quote.content }}
    </p>
    <p class="break-words pt-2" dir="auto" v-decrypt>{{ text }}</p>
    <div v-if="!replaying" class="flex justify-end text-gray-400 text-end">
      <div class="flex items-center" @click="replaying = true">
        <span class="mr-1">Reply</span> <ReplyIcon size="20" color="#9CA38F" />
      </div>
    </div>
    <div v-else class="flex flex-col mt-4">
      <textarea
        v-model="replayMessage"
        class="textarea textarea-bordered focus:outline-[#119af5]"
        placeholder="Write..."
        v-focus
      ></textarea>
      <Button :block="true" class="mt-4" @click="Submit"
        >Send Replay Message</Button
      >
      <p
        class="text-center pt-4 text-[#119af5] font-bold"
        @click="replaying = false"
      >
        Cancel
      </p>
    </div>
    <p v-if="replaySent" class="text-center text-[#119af5]">Replay Sent.</p>
  </div>
</template>
