import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth.js';

const Towers = () => import('../views/Towers.vue');
const Login = () => import('../views/Login.vue');
const Providers = () => import('../views/Providers.vue');
const Admin = () => import('../views/Admin.vue');
const AdminTowers = () => import('../views/AdminTowers.vue');
const AdminProviders = () => import('../views/AdminProviders.vue');
const AdminBlankspots = () => import('../views/AdminBlankspots.vue');
const AdminSettings = () => import('../views/AdminSettings.vue');
const AdminTimeline = () => import('../views/AdminTimeline.vue');

const routes = [
  {
    path: '/',
    name: 'Towers',
    component: Towers
  },
  {
    path: '/providers',
    name: 'Providers',
    component: Providers
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'AdminRedirect',
        redirect: '/admin/towers'
      },
      {
        path: 'towers',
        name: 'AdminTowers',
        component: AdminTowers
      },
      {
        path: 'providers',
        name: 'AdminProviders',
        component: AdminProviders
      },
      {
        path: 'blankspots',
        name: 'AdminBlankspots',
        component: AdminBlankspots
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: AdminSettings
      },
      {
        path: 'timeline',
        name: 'AdminTimeline',
        component: AdminTimeline
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else {
    next();
  }
});

export default router;