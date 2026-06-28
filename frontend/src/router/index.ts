import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Contact from '@/views/Contact.vue'
import Search from '@/views/Search.vue'
import Root from '@/views/Root.vue'
import Nothing from '@/views/Nothing.vue'

const routes = [
  { path: '/', component: Root },
  { path: '/Home', component: Home },
  { path: '/Contact', component: Contact },
  { path: '/Search', component: Search },
  { path: '/:pathMatch(.*)*', component: Nothing, meta: { title: '404 Not Found' } },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
