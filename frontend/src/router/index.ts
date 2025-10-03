import { createRouter, createWebHistory } from 'vue-router'

import Home from '@/views/HomeView.vue'
import Player from '@/views/PlayerView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },

    {
      path: '/player',
      name: 'player',
      component: Player,
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
