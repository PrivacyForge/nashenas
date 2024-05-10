<script lang="ts" setup>
import TelegramIcon from '@/components/icons/Telegram.vue'
</script>

<template>
  <div class="flex justify-center items-center h-screen w-screen">
    <div>
      <a
        class="flex bg-[#119af5] text-white font-semibold py-4 px-8 rounded-lg text-lg"
        href="https://telegram.me/nashenase2ebot?start=otp"
        target="_blank"
      >
        Connect to telegram
        <TelegramIcon class="ml-2" size="24" color="white" />
      </a>
    </div>
  </div>
</template>
