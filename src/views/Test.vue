<script lang="ts" setup>
import { onMounted } from 'vue';

// function store(id: string, value: string) {
//     window.Telegram.WebApp.CloudStorage.setItem("sessions-1", "23", (error, result) => {

//     })

// }

onMounted(() => {
    // store("1233", "odj3h28fg8923hf823gtrh23r8g23r")
    window.Telegram.WebApp.CloudStorage.setItem("x", "y")

    window.Telegram.WebApp.CloudStorage.getKeys((error, result) => {
        alert(error)
        alert(result)
    })
})
</script>

<template>
    <input type="file">
</template>