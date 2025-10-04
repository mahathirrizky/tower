<template>
  <div class="flex h-screen">
    <Toast />
    <Sidebar v-if="$route.path !== '/login' && !$route.path.startsWith('/admin')" />
    <AdminSidebar v-if="$route.path.startsWith('/admin')" />
    <div class="flex-1 p-4 overflow-y-auto">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import Sidebar from './components/Sidebar.vue';
import AdminSidebar from './components/AdminSidebar.vue';
import Toast from 'primevue/toast';
</script>