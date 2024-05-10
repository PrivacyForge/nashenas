import axios from '@/plugins/axios'
import { useUserStore } from '@/stores/user'
import type { NavigationGuard } from 'vue-router'

const middleware: NavigationGuard = async (to, from, next) => {
  const userStore = useUserStore()
  if (userStore.isAuth) return next()

  await axios
    .get('/me')
    .then(({ data }) => {
      userStore.user.id = data.id
      userStore.user.userid = data.userid
      userStore.user.username = data.username
      userStore.user.publicKey = data.public_key
      userStore.isAuth = true

      data.public_key ? next() : next({ name: 'setup' })
    })
    .catch(() => {
      next({ name: 'auth' })
    })
}

export default middleware
