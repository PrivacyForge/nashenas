import { createRouter, createWebHistory } from "vue-router";
import AuthMiddleware from "@/middlewares/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: { name: "inbox" },
      component: () => import("@/layouts/main.vue"),
      children: [
        {
          path: "/inbox",
          name: "inbox",
          component: () => import("@/views/Inbox.vue"),
          beforeEnter: AuthMiddleware,
        },
      ],
    },
    {
      path: "/auth",
      name: "auth",
      component: () => import("@/views/Auth.vue"),
    },
    {
      path: "/confirm/:code",
      name: "confirm",
      component: () => import("@/views/Confirm.vue"),
    },
    {
      path: "/setup",
      name: "setup",
      component: () => import("@/views/Setup.vue"),
      beforeEnter: AuthMiddleware,
    },
    {
      path: "/@:username([a-zA-Z]{1}[a-zA-Z0-9]{4,})",
      name: "profile",
      component: () => import("@/views/Profile.vue"),
    },
  ],
});

export default router;
