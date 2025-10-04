<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Manage Towers</h2>
      <div class="flex items-center gap-2">
        <Button label="New Tower" icon="pi pi-plus" class="p-button-success mr-2" @click="openNewTowerDialog" />
        <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
        <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
        <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
      </div>
    </div>
    <Card class="w-full">
      <template #content>
        <DataTable :value="towers" :paginator="true" :rows="5">
          <Column header="Nomor">
            <template #body="slotProps">
              {{ slotProps.index + 1 }}
            </template>
          </Column>
          <Column field="latitude" header="Latitude"></Column>
          <Column field="longitude" header="Longitude"></Column>
          <Column field="status" header="Status"></Column>
          <Column header="Photo">
            <template #body="slotProps">
              <img 
                v-if="slotProps.data.photo_url" 
                :src="slotProps.data.photo_url" 
                :alt="slotProps.data.ID" 
                class="w-10 h-10 object-cover rounded-full cursor-pointer hover:opacity-75 transition-opacity" 
                @click="showImageDialog(slotProps.data.photo_url)" 
              />
              <span v-else>No Photo</span>
            </template>
          </Column>
          <Column header="Providers">
            <template #body="slotProps">
              <span v-if="slotProps.data.providers && slotProps.data.providers.length > 0">
                <span v-for="provider in slotProps.data.providers" :key="provider.ID" class="mr-1">{{ provider.name }}</span>
              </span>
              <span v-else>N/A</span>
            </template>
          </Column>
          <Column header="Actions">
            <template #body="slotProps">
              <Button v-tooltip.top="'Edit Tower'" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" @click="editTower(slotProps.data)" />
              <Button v-tooltip.top="'Relocate Tower'" icon="pi pi-map-marker" class="p-button-rounded p-button-info mr-2" @click="confirmRelocateTower(slotProps.data)" />
              <Button v-tooltip.top="'Change Ownership'" icon="pi pi-users" class="p-button-rounded p-button-secondary mr-2" @click="confirmChangeOwnership(slotProps.data)" />
              <Button v-tooltip.top="'Dismantle Tower'" icon="pi pi-trash" class="p-button-rounded p-button-danger" @click="confirmDismantleTower(slotProps.data)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Image Preview Dialog -->
    <Dialog v-model:visible="imageDialogVisible" modal header="Tower Photo" :style="{ width: '50vw' }" :breakpoints="{ '960px': '75vw', '640px': '90vw' }" >
      <img :src="displayedImageUrl" alt="Tower Photo" class="w-full h-auto" />
    </Dialog>

    <!-- Tower Details Dialog (Create/Edit) -->
    <Dialog v-model:visible="towerDialog" :style="{ width: '600px' }" header="Tower Details" :modal="true" class="p-fluid" @show="initDialogMap" @hide="destroyDialogMap">
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputNumber id="latitude" v-model="tower.latitude" mode="decimal" :minFractionDigits="2" :maxFractionDigits="6" fluid />
          <label for="latitude">Latitude</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputNumber id="longitude" v-model="tower.longitude" mode="decimal" :minFractionDigits="2" :maxFractionDigits="6" fluid />
          <label for="longitude">Longitude</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <label>Location on Map</label>
        <div id="dialogMap" style="height: 300px; width: 100%;"></div>
      </div>
      <div class="field my-4">
        <label for="photo">Photo</label>
        <FileUpload mode="basic" name="photo" accept="image/*" @select="onFileSelect" fluid />
        <img v-if="tower.photo_url" :src="tower.photo_url" :alt="tower.ID" class="mt-2 w-20 h-20 object-cover" />
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <MultiSelect id="providers" v-model="selectedProviders" :options="providers" optionLabel="name" display="chip" dataKey="ID" fluid />
          <label for="providers">Providers</label>
        </FloatLabel>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideTowerDialog" />
        <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveTower" />
      </template>
    </Dialog>

    <!-- Relocate Tower Dialog -->
    <Dialog v-model:visible="relocateTowerDialog" :style="{ width: '450px' }" header="Relocate Tower" :modal="true" class="p-fluid" @show="initRelocateMap" @hide="destroyRelocateMap">
      <div class="confirmation-content">
        <i class="pi pi-map-marker mr-3" style="font-size: 2rem" />
        <span v-if="towerToRelocate">Relocate tower <b>{{ towerToRelocate.ID }}</b>?</span>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputNumber id="relocateLatitude" v-model="relocateInput.latitude" mode="decimal" :minFractionDigits="2" :maxFractionDigits="6" fluid />
          <label for="relocateLatitude">New Latitude</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputNumber id="relocateLongitude" v-model="relocateInput.longitude" mode="decimal" :minFractionDigits="2" :maxFractionDigits="6" fluid />
          <label for="relocateLongitude">New Longitude</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <label>New Location on Map</label>
        <div id="relocateMap" style="height: 300px; width: 100%;"></div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="relocateTowerDialog = false" />
        <Button label="Relocate" icon="pi pi-check" class="p-button-text" @click="relocateTower" />
      </template>
    </Dialog>

    <!-- Change Ownership Dialog -->
    <Dialog v-model:visible="changeOwnershipDialog" :style="{ width: '450px' }" header="Change Ownership" :modal="true" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-users mr-3" style="font-size: 2rem" />
        <span v-if="towerToChangeOwnership">Change ownership for tower <b>{{ towerToChangeOwnership.ID }}</b>?</span>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <Dropdown id="newProvider" v-model="changeOwnershipInput.newProviderId" :options="providers" optionLabel="name" optionValue="ID" placeholder="Select New Provider" fluid />
          <label for="newProvider">New Provider</label>
        </FloatLabel>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="changeOwnershipDialog = false" />
        <Button label="Change" icon="pi pi-check" class="p-button-text" @click="changeOwnership" />
      </template>
    </Dialog>

    <!-- Dismantle Tower Dialog -->
    <Dialog v-model:visible="dismantleTowerDialog" :style="{ width: '450px' }" header="Confirm Dismantle" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="towerToDismantle">Are you sure you want to dismantle tower <b>{{ towerToDismantle.ID }}</b>? This action cannot be undone.</span>
      </div>
      <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="dismantleTowerDialog = false" />
        <Button label="Yes, Dismantle" icon="pi pi-check" class="p-button-text" @click="dismantleTower" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick } from 'vue';
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
import InputNumber from 'primevue/inputnumber';
import FileUpload from 'primevue/fileupload';
import MultiSelect from 'primevue/multiselect';
import Dropdown from 'primevue/dropdown';
import FloatLabel from 'primevue/floatlabel';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Stores
const dataStore = useDataStore();
const uiStore = useUiStore();

// Store State
const { towers, providers } = storeToRefs(dataStore);

// Store Actions
const { fetchTowers, createTower, updateTower, dismantleTower: dismantleTowerAction, changeTowerOwnership: changeTowerOwnershipAction, relocateTower: relocateTowerAction, fetchProviders } = dataStore;
const { toggleDarkMode } = uiStore;

// UI & State
const toast = useToast();
const towerDialog = ref(false);
const relocateTowerDialog = ref(false);
const changeOwnershipDialog = ref(false);
const dismantleTowerDialog = ref(false);
const tower = ref({});
const towerToRelocate = ref(null);
const towerToChangeOwnership = ref(null);
const towerToDismantle = ref(null);
const selectedProviders = ref([]);
const submitted = ref(false);
const photoFile = ref(null);

const relocateInput = ref({ latitude: null, longitude: null });
const changeOwnershipInput = ref({ newProviderId: null });

const isDarkMode = computed(() => uiStore.isDarkMode);

// State for Image Preview Dialog
const imageDialogVisible = ref(false);
const displayedImageUrl = ref('');

const showImageDialog = (url) => {
  displayedImageUrl.value = url;
  imageDialogVisible.value = true;
};

// Map instances
let dialogMapInstance = null;
let dialogMarker = null;
let relocateMapInstance = null;
let relocateMarker = null;

// Lifecycle

// Map instances
onMounted(() => {
  fetchTowers();
  fetchProviders();
});

// Methods for Tower CRUD (Create/Edit)
const openNewTowerDialog = () => {
  tower.value = { latitude: null, longitude: null, photo_url: '', Status: 'active' };
  selectedProviders.value = [];
  photoFile.value = null;
  submitted.value = false;
  towerDialog.value = true;
};

const hideTowerDialog = () => {
  towerDialog.value = false;
  submitted.value = false;
};

const saveTower = async () => {
  submitted.value = true;
  if (tower.value.latitude && tower.value.longitude && (photoFile.value || tower.value.photo_url)) {
    const formData = new FormData();
    formData.append('latitude', tower.value.latitude);
    formData.append('longitude', tower.value.longitude);
    formData.append('status', tower.value.Status);
    if (photoFile.value) {
      formData.append('photo', photoFile.value);
    }
    selectedProviders.value.forEach(provider => {
      formData.append('provider_ids', provider.ID);
    });

    try {
      if (tower.value.ID) {
        await updateTower(tower.value.ID, formData);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Updated', life: 3000 });
      } else {
        await createTower(formData);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Created', life: 3000 });
      }
      towerDialog.value = false;
      tower.value = {};
      selectedProviders.value = [];
      photoFile.value = null;
      fetchTowers();
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Operation failed', life: 3000 });
    }
  }
};

const editTower = (t) => {
  tower.value = { ...t };
  selectedProviders.value = t.providers ? t.providers.map(p => ({ ID: p.ID, name: p.name })) : [];
  towerDialog.value = true;
};

const onFileSelect = (event) => {
  photoFile.value = event.files[0];
  const reader = new FileReader();
  reader.onload = (e) => {
    tower.value.photo_url = e.target.result;
  };
  reader.readAsDataURL(photoFile.value);
};

// Methods for Relocate Tower
const confirmRelocateTower = (t) => {
  towerToRelocate.value = { ...t };
  relocateInput.value.latitude = t.latitude;
  relocateInput.value.longitude = t.longitude;
  relocateTowerDialog.value = true;
};

const relocateTower = async () => {
  try {
    await relocateTowerAction(towerToRelocate.value.ID, relocateInput.value.latitude, relocateInput.value.longitude);
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Relocated', life: 3000 });
    relocateTowerDialog.value = false;
    towerToRelocate.value = null;
    fetchTowers();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Relocation failed', life: 3000 });
  }
};

// Methods for Change Ownership
const confirmChangeOwnership = (t) => {
  towerToChangeOwnership.value = { ...t };
  changeOwnershipInput.value.newProviderId = t.providers && t.providers.length > 0 ? t.providers[0].ID : null;
  changeOwnershipDialog.value = true;
};

const changeOwnership = async () => {
  try {
    await changeTowerOwnershipAction(towerToChangeOwnership.value.ID, changeOwnershipInput.value.newProviderId);
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Ownership Changed', life: 3000 });
    changeOwnershipDialog.value = false;
    towerToChangeOwnership.value = null;
    fetchTowers();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Ownership change failed', life: 3000 });
  }
};

// Methods for Dismantle Tower
const confirmDismantleTower = (t) => {
  towerToDismantle.value = { ...t };
  dismantleTowerDialog.value = true;
};

const dismantleTower = async () => {
  try {
    await dismantleTowerAction(towerToDismantle.value.ID);
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Dismantled', life: 3000 });
    dismantleTowerDialog.value = false;
    towerToDismantle.value = null;
    fetchTowers();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Dismantling failed', life: 3000 });
  }
};

// Map Methods for Tower Details Dialog
const initDialogMap = () => {
  nextTick(() => {
    if (dialogMapInstance) {
      dialogMapInstance.remove();
    }
    dialogMapInstance = L.map('dialogMap').setView([-8.4601, 118.7267], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(dialogMapInstance);

    dialogMapInstance.on('click', (e) => {
      tower.value.latitude = e.latlng.lat;
      tower.value.longitude = e.latlng.lng;
      updateDialogMarker();
    });

    updateDialogMarker();
  });
};

const destroyDialogMap = () => {
  if (dialogMapInstance) {
    dialogMapInstance.remove();
    dialogMapInstance = null;
    dialogMarker = null;
  }
};

const updateDialogMarker = () => {
  if (dialogMapInstance && tower.value.latitude && tower.value.longitude) {
    const latlng = [tower.value.latitude, tower.value.longitude];
    if (dialogMarker) {
      dialogMarker.setLatLng(latlng);
    } else {
      dialogMarker = L.marker(latlng).addTo(dialogMapInstance);
    }
    dialogMapInstance.setView(latlng, dialogMapInstance.getZoom());
  } else if (dialogMarker) {
    dialogMapInstance.removeLayer(dialogMarker);
    dialogMarker = null;
  }
};

// Map Methods for Relocate Tower Dialog
const initRelocateMap = () => {
  nextTick(() => {
    if (relocateMapInstance) {
      relocateMapInstance.remove();
    }
    relocateMapInstance = L.map('relocateMap').setView([relocateInput.value.latitude || -8.4601, relocateInput.value.longitude || 118.7267], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(relocateMapInstance);

    relocateMapInstance.on('click', (e) => {
      relocateInput.value.latitude = e.latlng.lat;
      relocateInput.value.longitude = e.latlng.lng;
      updateRelocateMarker();
    });

    updateRelocateMarker();
  });
};

const destroyRelocateMap = () => {
  if (relocateMapInstance) {
    relocateMapInstance.remove();
    relocateMapInstance = null;
    relocateMarker = null;
  }
};

const updateRelocateMarker = () => {
  if (relocateMapInstance && relocateInput.value.latitude && relocateInput.value.longitude) {
    const latlng = [relocateInput.value.latitude, relocateInput.value.longitude];
    if (relocateMarker) {
      relocateMarker.setLatLng(latlng);
    } else {
      relocateMarker = L.marker(latlng).addTo(relocateMapInstance);
    }
    relocateMapInstance.setView(latlng, relocateMapInstance.getZoom());
  } else if (relocateMarker) {
    relocateMapInstance.removeLayer(relocateMarker);
    relocateMarker = null;
  }
};
</script>