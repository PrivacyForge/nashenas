<script lang="ts" setup>
import { computed ,onMounted} from 'vue'
import { RouterView } from 'vue-router'

import { useUserStore } from '@/stores/user'

import Settings from '@/components/Settings.vue'
import CopyText from '@/components/UI/CopyText.vue'

const userStore = useUserStore()

const myLink = computed(() => `https://t.me/E2EChatbot?start=${userStore.user.username}`)
onMounted(() => {
  window.Telegram.WebApp.expand()
  window.Telegram.WebApp.disableVerticalSwipes()
}),
</script>

<template>
  <div class="lg:w-5/12 mx-auto px-4">
    <nav
      class="flex justify-between bg-[#ffffff] px-5 py-4 my-4 rounded-lg shadow-sm"
    >
      <div class="cursor-pointer">
        <Settings />
      </div>
      <CopyText
        text="کپی لینک"
        :copy="myLink"
        class="text-[#119af5] test- font-semibold text-end"
      />
    </nav>
    <hr class="my-1" />
    <RouterView />
  </div>
</template>
