import axios from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth'
import type { NavigationGuard } from 'vue-router'

const middleware: NavigationGuard = async (to, from, next) => {
  const authStore = useAuthStore()
  if (authStore.isAuth) return next()

  await axios
    .get('/me')
    .then(({ data }) => {
      authStore.user.ID = data.ID
      authStore.user.Userid = data.Userid
      authStore.user.Username = data.Username
      authStore.user.PublicKey = data.PublicKey
      authStore.isAuth = true

      next()
    })
    .catch(() => {
      next({ name: 'auth' })
    })
}

export default middleware
