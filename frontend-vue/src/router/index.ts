import { createRouter, createWebHistory } from 'vue-router'
import { isAuthenticated } from '@/api.ts'

const authedRoutes = {

}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '',
      redirect: "/entry",
    },
    {
      path: '/entry',
      name: 'entry',
      component: () => {},
    },
    {
      path: '/incident-detail/:id',
      name: 'incident-detail',
      component: () => import('@/views/IncidentDetailView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/incidents',
      name: 'incidents',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/IncidentsView.vue'),
    },
  ],
})

router.beforeEach(async (to) => {
  if (to.name == "entry") {
    const ok = await isAuthenticated()
    return ok == true ? "incidents" : "login"
  }
})

export default router
