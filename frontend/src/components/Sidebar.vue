<template>
  <div class="w-16 bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200 flex flex-col items-center py-4">
    <router-link v-tooltip.right="'Towers'" to="/" class="p-4 rounded-md transition-colors duration-200" :class="{ 'bg-blue-500 text-white': $route.path === '/', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/' }">
      <i class="pi pi-sitemap" style="font-size: 1.5rem"></i>
    </router-link>
    <router-link v-tooltip.right="'Providers'" to="/providers" class="p-4 rounded-md transition-colors duration-200 mt-4" :class="{ 'bg-blue-500 text-white': $route.path === '/providers', 'hover:bg-gray-200 dark:hover:bg-gray-700': $route.path !== '/providers' }">
      <i class="pi pi-building" style="font-size: 1.5rem"></i>
    </router-link>
  </div>
</template>

<script setup>
</script>