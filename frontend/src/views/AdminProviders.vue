<template>
  <div>
    
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Manage Providers</h2>
      <div class="flex items-center gap-2">
        <Button label="New Provider" icon="pi pi-plus" class="p-button-success mr-2" @click="openNewProviderDialog" />
        <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
        <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
        <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
      </div>
    </div>
    <Card class="w-full">
      <template #content>
        <DataTable :value="providers" :paginator="true" :rows="5">
          <Column header="Nomor">
            <template #body="slotProps">
              {{ slotProps.index + 1 }}
            </template>
          </Column>
          <Column field="name" header="Name"></Column>
          <Column field="address" header="Address"></Column>
          <Column header="Actions">
            <template #body="slotProps">
              <Button v-tooltip.top="'Edit Provider'" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" @click="editProvider(slotProps.data)" />
              <Button v-tooltip.top="'Delete Provider'" icon="pi pi-trash" class="p-button-rounded p-button-danger" @click="confirmDeleteProvider(slotProps.data)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <Dialog v-model:visible="providerDialog" :style="{ width: '450px' }" header="Provider Details" :modal="true" class="p-fluid">
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="name" v-model.trim="provider.name" required="true" autofocus :class="{ 'p-invalid': submitted && !provider.name }" fluid />
          <label for="name">Name</label>
        </FloatLabel>
        <small class="p-error" v-if="submitted && !provider.name">Name is required.</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="address" v-model.trim="provider.address" fluid />
          <label for="address">Address</label>
        </FloatLabel>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideProviderDialog" />
        <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveProvider" />
      </template>
    </Dialog>

    <Dialog v-model:visible="deleteProviderDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="provider">Are you sure you want to delete provider <b>{{ provider.name }}</b>?</span>
      </div>
      <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="deleteProviderDialog = false" />
        <Button label="Yes" icon="pi pi-check" class="p-button-text" @click="deleteProvider" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useDataStore } from '../stores/data';
import { useUiStore } from '../stores/ui';
import { useToast } from 'primevue/usetoast';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Card from 'primevue/card';
import ToggleSwitch from 'primevue/toggleswitch';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';

import FloatLabel from 'primevue/floatlabel';

// Stores
const dataStore = useDataStore();
const uiStore = useUiStore();

// Store State
const { providers } = storeToRefs(dataStore);

// Store Actions
const { fetchProviders, createProvider, updateProvider, deleteProvider: deleteProviderAction } = dataStore;
const { toggleDarkMode } = uiStore;

// UI & State
const toast = useToast();
const providerDialog = ref(false);
const deleteProviderDialog = ref(false);
const provider = ref({});
const submitted = ref(false);

const isDarkMode = computed(() => uiStore.isDarkMode);

// Lifecycle
onMounted(() => {
  fetchProviders();
});

// Methods
const openNewProviderDialog = () => {
  provider.value = { name: '', address: '' };
  submitted.value = false;
  providerDialog.value = true;
};

const hideProviderDialog = () => {
  providerDialog.value = false;
  submitted.value = false;
};

const saveProvider = async () => {
  submitted.value = true;
  if (provider.value.name) {
    try {
      if (provider.value.ID) {
        await updateProvider(provider.value.ID, provider.value);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Provider Updated', life: 3000 });
      } else {
        await createProvider(provider.value);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Provider Created', life: 3000 });
      }
      providerDialog.value = false;
      provider.value = {};
      fetchProviders(); // Refresh the list
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Operation failed', life: 3000 });
    }
  }
};

const editProvider = (p) => {
  provider.value = { ...p };
  providerDialog.value = true;
};

const confirmDeleteProvider = (p) => {
  provider.value = p;
  deleteProviderDialog.value = true;
};

const deleteProvider = async () => {
  try {
    await deleteProviderAction(provider.value.ID);
    deleteProviderDialog.value = false;
    provider.value = {};
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Provider Deleted', life: 3000 });
    fetchProviders(); // Refresh the list
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Deletion failed', life: 3000 });
  }
};
</script>
