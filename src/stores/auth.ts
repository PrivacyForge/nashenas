import { ref, computed, reactive } from 'vue'
import { defineStore } from 'pinia'

interface IUser {
  ID: number
  Userid: number
  Username: string
  PublicKey: string
}

export const useAuthStore = defineStore('auth', () => {
  const isAuth = ref(false)

  const user = reactive<IUser>({
    ID: 0,
    Userid: 0,
    Username: '',
    PublicKey: '',
  })

  return { isAuth, user }
})
