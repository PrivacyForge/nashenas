import { createApp } from 'vue'
import { createPinia } from 'pinia'
import type { Telegram } from '@twa-dev/types'

import App from './App.vue'
import router from './router'
import './assets/styles/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

declare global {
  interface Window {
    Telegram: Telegram
  }
}

app.mount('#app')
