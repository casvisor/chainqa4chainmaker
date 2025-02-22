import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UploadView from '@/views/upload/UploadView.vue'
import DownloadView from '@/views/download/DownloadView.vue'
import QueryView from '@/views/query/QueryView.vue'
import LogView from '@/views/log/LogView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/upload',
      name: 'upload',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: UploadView,
    },
    {
      path: '/download',
      name: 'download',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: DownloadView,
    },
    {
      path: '/query',
      name: 'query',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: QueryView,
    },
    {
      path: '/log',
      name: 'log',
      component: LogView,
    }
  ],
})

export default router
