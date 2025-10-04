<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-3xl font-bold">Providers</h1>
      <div class="flex items-center gap-2">
        <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
        <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
        <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
      </div>
    </div>
    <DataTable :value="filteredProviders" :paginator="true" :rows="5">
      <template #header>
        <div class="flex justify-end">
          <span class="p-input-icon-left">
            <i class="pi pi-search" />
            <InputText v-model="searchQuery" placeholder="Search Providers" />
          </span>
        </div>
      </template>
      <Column field="name" header="Name"></Column>
      <Column field="address" header="Address"></Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useDataStore } from '../stores/data';
import { useUiStore } from '../stores/ui';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import ToggleSwitch from 'primevue/toggleswitch';

// Stores
const dataStore = useDataStore();
const uiStore = useUiStore();

// Store State
const { providers } = storeToRefs(dataStore);

// Store Actions
const { fetchProviders } = dataStore;
const { toggleDarkMode } = uiStore;

// State
const searchQuery = ref('');

// Computed
const isDarkMode = computed(() => uiStore.isDarkMode);

const filteredProviders = computed(() => {
  if (!searchQuery.value) {
    return providers.value;
  }
  const query = searchQuery.value.toLowerCase();
  return providers.value.filter(provider =>
    (provider.name && provider.name.toLowerCase().includes(query)) ||
    (provider.address && provider.address.toLowerCase().includes(query))
  );
});

// Lifecycle
onMounted(() => {
  fetchProviders();
});
</script>