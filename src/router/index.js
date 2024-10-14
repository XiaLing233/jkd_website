// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import Contact from '../views/Contact.vue';
import Search from '../views/Search.vue';
import Root from '../views/Root.vue';

const routes = [
  { path: '/', component: Root },
  { path: '/Home', component: Home },
  { path: '/Contact', component: Contact },
  { path: '/Search', component: Search},
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
