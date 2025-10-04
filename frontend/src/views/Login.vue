<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-white dark:bg-gray-800 rounded-xl shadow-lg">
      <div class="flex justify-end">
        <div class="flex items-center gap-2">
          <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
          <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
          <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
        </div>
      </div>
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
          Sign in to your account
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="space-y-4">
          <div class="field pt-4">
            <FloatLabel variant="on">
              <InputText id="email-address" name="email" type="email" autocomplete="email" required v-model="email"
                class="w-full dark:bg-gray-700 dark:text-white" />
              <label for="email-address">Email address</label>
            </FloatLabel>
          </div>
          <div class="field pt-4">
            <FloatLabel variant="on">
              <Password id="password" v-model="password" toggleMask required inputClass="w-full dark:bg-gray-700 dark:text-white" autocomplete="current-password" fluid></Password>
              <label for="password">Password</label>
            </FloatLabel>
          </div>
        </div>

        <div class="pt-4">
          <Button type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            label="Sign in" />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth.js';
import { useUiStore } from '../stores/ui';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import ToggleSwitch from 'primevue/toggleswitch';
import FloatLabel from 'primevue/floatlabel';
import Password from 'primevue/password';

// Stores & Router
const authStore = useAuthStore();
const uiStore = useUiStore();
const router = useRouter();

// State
const email = ref('');
const password = ref('');

// Computed
const isDarkMode = computed(() => uiStore.isDarkMode);

// Actions
const { toggleDarkMode } = uiStore;

// Methods
const handleLogin = async () => {
  try {
    const credentials = {
      email: email.value,
      password: password.value
    };
    await authStore.login(credentials);
    router.push('/admin');
  } catch (error) {
    console.error('Error logging in:', error);
    // Optionally, show a toast message on error
  }
};
</script>
