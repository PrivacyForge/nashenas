<script lang="ts" setup>
import forge from 'node-forge'

const props = defineProps<{
  text: string
  time: string
  mark: boolean
}>()

const vDecrypt = {
  mounted: (el: HTMLParagraphElement) => {
    try {
      // const publicKey = localStorage.getItem('public_key')
      const privateKey = localStorage.getItem('private_key')

      // const encryptedMsg = forge.pki
      //   .publicKeyFromPem(publicKey!)
      //   .encrypt("hello Yasha!")

      const decrypted = forge.pki
        .privateKeyFromPem(privateKey!)
        .decrypt(forge.util.hexToBytes(el.innerText))

      el.innerText = decrypted
    } catch (error) {
      el.innerText = "Decryption Error."
    }
  },
}
</script>
<template>
  <div
    class="flex flex-col bg-[#ffffff] p-5 rounded-lg"
    :class="mark && ['border-2 border-[#119af5]']"
  >
    <p class="mb-2 text-gray-400 text-end">{{ time }}</p>
    <p class="break-words" dir="auto" v-decrypt>{{ text }}</p>
  </div>
</template>
