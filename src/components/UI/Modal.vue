<script lang="ts" setup>
// copy-pasted from official Vue examples
const model = defineModel()

function close() {
  model.value = false
}
</script>

<template>
  <Transition name="modal">
    <dialog v-if="model" class="modal-mask" @click.self="close()">
      <div class="modal-container rounded-xl">
        <div class="flex justify-between items-center">
          <slot name="header">header-text</slot>
          <button @click="close()">✕</button>
        </div>

        <div class="modal-body">
          <slot name="body">body-text</slot>
        </div>

        <div class="modal-footer">
          <slot name="footer" />
        </div>
      </div>
    </dialog>
  </Transition>
</template>

<style>
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  transition: opacity 0.3s ease;
}

.modal-container {
  width: 90%;
  max-width: 600px;
  margin: auto;
  padding: 20px 20px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-body {
  margin: 20px 0;
}

.modal-default-button {
  float: right;
}

.modal-enter-from {
  opacity: 0;
}

.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  -webkit-transform: scale(1.1);
  transform: scale(1.1);
}
</style>
