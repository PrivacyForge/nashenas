<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '@/plugins/axios'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const errorMessage = ref('')
const code = route.params.code

onMounted(() => {
  setTimeout(() => {
    axios
      .get(`/confirm?code=${code}`)
      .then((response) => {
        localStorage.setItem('token', response.data)
        router.push({ name: 'setup' })
      })
      .catch((error) => {
        errorMessage.value = error.response.data
      })
      .finally(() => {
        loading.value = false
      })
  }, 1000)
})
</script>

<template>
  <div class="grid grid-cols-1 m-4 bg-white shadow-sm p-5 rounded-lg">
    <div v-if="loading" class="flex justify-center items-center py-6">
      <p class="text-[#119af5] font-semibold">Waiting for server response...</p>
      <span class="loading loading-infinity loading-lg mx-2 text-[#119af5]" />
    </div>
    <div v-else>
      <p class="text-red-500 text-center" v-text="errorMessage"></p>
    </div>
  </div>
</template>
