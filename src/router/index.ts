import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: 'inbox',
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/views/Auth.vue'),
    },
    // {
    //   path: '/@:username',
    //   name: 'profile',
    //   component: () => import('@/views/Profile.vue'),
    // },
    {
      path: '/',
      component: () => import('@/layouts/main.vue'),
      children: [
        {
          path: '/inbox',
          name: 'inbox',
          component: () => import('@/views/Inbox.vue'),
        },
      ],
    },
  ],
})

export default router
