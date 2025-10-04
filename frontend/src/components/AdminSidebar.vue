<template>
  <div class="w-16 bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200 flex flex-col items-center py-4">
    <router-link v-tooltip.right="'Towers'" to="/admin/towers" class="p-4 rounded-md transition-colors duration-200" :class="{ 'bg-blue-500 text-white': $route.path === '/admin/towers', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/admin/towers' }">
      <i class="pi pi-sitemap" style="font-size: 1.5rem"></i>
    </router-link>
    <router-link v-tooltip.right="'Providers'" to="/admin/providers" class="p-4 rounded-md transition-colors duration-200 mt-4" :class="{ 'bg-blue-500 text-white': $route.path === '/admin/providers', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/admin/providers' }">
      <i class="pi pi-building" style="font-size: 1.5rem"></i>
    </router-link>
    <router-link v-tooltip.right="'Blankspots'" to="/admin/blankspots" class="p-4 rounded-md transition-colors duration-200 mt-4" :class="{ 'bg-blue-500 text-white': $route.path === '/admin/blankspots', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/admin/blankspots' }">
      <i class="pi pi-map" style="font-size: 1.5rem"></i> <!-- Using pi-map for blankspots -->
    </router-link>

    <router-link v-tooltip.right="'Timeline'" to="/admin/timeline" class="p-4 rounded-md transition-colors duration-200 mt-4" :class="{ 'bg-blue-500 text-white': $route.path === '/admin/timeline', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/admin/timeline' }">
      <i class="pi pi-history" style="font-size: 1.5rem"></i>
    </router-link>

    <router-link v-tooltip.right="'Settings'" to="/admin/settings" class="p-4 rounded-md transition-colors duration-200 mt-4" :class="{ 'bg-blue-500 text-white': $route.path === '/admin/settings', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/admin/settings' }">
      <i class="pi pi-cog" style="font-size: 1.5rem"></i> <!-- Settings icon -->
    </router-link>

    <div class="mt-auto flex flex-col items-center gap-4">
      <a v-tooltip.right="'Logout'" @click="handleLogout" class="p-4 rounded-md transition-colors duration-200 hover:bg-gray-200 dark:hover:bg-gray-700 cursor-pointer block">
        <i class="pi pi-sign-out" style="font-size: 1.5rem"></i>
      </a>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};
</script>
