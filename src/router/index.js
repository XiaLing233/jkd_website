// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import Contact from '../views/Contact.vue';
import Search from '../views/Search.vue';

const routes = [
  { path: '/Home', component: Home },
  { path: '/Contact', component: Contact },
  { path: '/Search', component: Search},
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
