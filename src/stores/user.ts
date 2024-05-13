import { ref, reactive } from 'vue'
import { defineStore } from 'pinia'

interface IUser {
  id: number
  userid: number
  username: string
  receivePublicKey: string
  sendPublicKey: string
}

export const useUserStore = defineStore('user', () => {
  const isAuth = ref(false)
  const user = reactive<IUser>({
    id: 0,
    userid: 0,
    username: '',
    receivePublicKey: '',
    sendPublicKey: '',
  })

  return { isAuth, user }
})
