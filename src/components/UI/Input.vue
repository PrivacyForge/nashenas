<script lang="ts" setup>
const value = defineModel<string>()

defineProps<{
  placeholder?: string
}>()
</script>
<template>
  <input
    v-model="value"
    type="text"
    class="input input-bordered w-full focus:outline-[#119af5] outline-[#119af5] border rounded-md py-3"
    :placeholder="placeholder"
  />
</template>
