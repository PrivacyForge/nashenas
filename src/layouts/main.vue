<script lang="ts" setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div>
    <nav class="flex justify-between bg-[#ffffff] p-5 m-4 rounded-lg">
        <img src="@/assets/lock.svg" alt="">
        <RouterLink to="" class="text-[#119af5] test- font-semibold text-end">My Profile</RouterLink>
    </nav>
    <hr class="my-1 mx-4">
    <RouterView />
  </div>
</template>
