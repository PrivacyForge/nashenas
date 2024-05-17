<script lang="ts" setup>
const value = defineModel<string>()

defineProps<{
  placeholder?: string
}>()
</script>

<template>
  <textarea v-model="value" :placeholder="placeholder" dir="rtl"></textarea>
</template>

<style scoped>
textarea {
  min-height: 3rem /* 48px */;
  flex-shrink: 1;
  padding-left: 1rem /* 16px */;
  padding-right: 1rem /* 16px */;
  padding-top: 0.5rem /* 8px */;
  padding-bottom: 0.5rem /* 8px */;
  font-size: 1rem /* 14px */;
  line-height: 1.25rem /* 20px */;
  line-height: 2;
  border-radius: var(--rounded-btn, 0.5rem /* 8px */);
  border-width: 1px;
  --tw-bg-opacity: 1;
  background-color: var(--fallback-b1, oklch(var(--b1) / var(--tw-bg-opacity)));
}
textarea:focus {
  box-shadow: none;
  outline-style: solid;
  outline-width: 2px;
  outline-offset: 2px;
  outline-color: var(--fallback-bc, oklch(var(--bc) / 0.2));
  outline-color: #119af5;
}
textarea-disabled,
textarea:disabled,
textarea[disabled] {
  cursor: not-allowed;
  --tw-border-opacity: 1;
  border-color: var(--fallback-b2, oklch(var(--b2) / var(--tw-border-opacity)));
  --tw-bg-opacity: 1;
  background-color: var(--fallback-b2, oklch(var(--b2) / var(--tw-bg-opacity)));
  color: var(--fallback-bc, oklch(var(--bc) / 0.4));
}
textarea-disabled,
textarea:disabled,
textarea[disabled] {
  cursor: not-allowed;
  --tw-border-opacity: 1;
  border-color: var(--fallback-b2, oklch(var(--b2) / var(--tw-border-opacity)));
  --tw-bg-opacity: 1;
  background-color: var(--fallback-b2, oklch(var(--b2) / var(--tw-bg-opacity)));
  color: var(--fallback-bc, oklch(var(--bc) / 0.4));
}
textarea-disabled::placeholder,
textarea:disabled::placeholder,
textarea[disabled]::placeholder {
  color: var(--fallback-bc, oklch(var(--bc) / var(--tw-placeholder-opacity)));
  --tw-placeholder-opacity: 0.2;
}
textarea-disabled::placeholder,
textarea:disabled::placeholder,
textarea[disabled]::placeholder {
  color: var(--fallback-bc, oklch(var(--bc) / var(--tw-placeholder-opacity)));
  --tw-placeholder-opacity: 0.2;
}
</style>
