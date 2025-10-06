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
        <Tabs value="0">
          <TabList>
            <Tab value="0">Active</Tab>
            <Tab value="1">Dismantled</Tab>
          </TabList>
          <TabPanels>
            <TabPanel value="0">
              <DataTable :value="activeTowers" :paginator="true" :rows="5" :key="activeTowers.length">
                <Column header="Nomor">
                  <template #body="slotProps">
                    {{ slotProps.index + 1 }}
                  </template>
                </Column>
                <Column field="latitude" header="Latitude"></Column>
                <Column field="longitude" header="Longitude"></Column>
                <Column field="kecamatan" header="Kecamatan"></Column>
                <Column field="kelurahan" header="Kelurahan"></Column>
                <Column field="address" header="Alamat"></Column>
                <Column field="status" header="Status"></Column>
                <Column header="Photo">
                  <template #body="slotProps">
                    <div v-if="slotProps.data">
                      <img 
                        v-if="slotProps.data.photo_url" 
                        :src="slotProps.data.photo_url" 
                        :alt="slotProps.data.ID" 
                        class="w-10 h-10 object-cover rounded-full cursor-pointer hover:opacity-75 transition-opacity" 
                        @click="showImageDialog(slotProps.data.photo_url)" 
                      />
                      <span v-else>No Photo</span>
                    </div>
                  </template>
                </Column>
                <Column header="Providers">
                  <template #body="slotProps">
                    <div v-if="slotProps.data">
                      <span v-if="slotProps.data.providers && slotProps.data.providers.length > 0">
                        <span v-for="provider in slotProps.data.providers" :key="provider.ID" class="mr-1">{{ provider.name }}</span>
                      </span>
                      <span v-else>N/A</span>
                    </div>
                  </template>
                </Column>
                <Column header="Actions">
                  <template #body="slotProps">
                    <div v-if="slotProps.data">
                      <Button v-tooltip.top="'Edit Tower'" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" @click="editTower(slotProps.data)" />
                      <Button v-tooltip.top="'Relocate Tower'" icon="pi pi-map-marker" class="p-button-rounded p-button-info mr-2" @click="confirmRelocateTower(slotProps.data)" />
                      <Button v-tooltip.top="'Change Ownership'" icon="pi pi-users" class="p-button-rounded p-button-secondary mr-2" @click="confirmChangeOwnership(slotProps.data)" />
                      <Button v-tooltip.top="'Dismantle Tower'" icon="pi pi-power-off" class="p-button-rounded p-button-warning" @click="confirmDismantleTower(slotProps.data)" />
                    </div>
                  </template>
                </Column>
              </DataTable>
            </TabPanel>
            <TabPanel value="1">
              <DataTable :value="dismantledTowers" :paginator="true" :rows="5" :key="dismantledTowers.length">
                <Column header="Nomor">
                  <template #body="slotProps">
                    {{ slotProps.index + 1 }}
                  </template>
                </Column>
                <Column field="latitude" header="Latitude"></Column>
                <Column field="longitude" header="Longitude"></Column>
                <Column field="kecamatan" header="Kecamatan"></Column>
                <Column field="kelurahan" header="Kelurahan"></Column>
                <Column field="address" header="Alamat"></Column>
                <Column field="status" header="Status"></Column>
                <Column header="Photo">
                  <template #body="slotProps">
                    <div v-if="slotProps.data">
                      <img 
                        v-if="slotProps.data.photo_url" 
                        :src="slotProps.data.photo_url" 
                        :alt="slotProps.data.ID" 
                        class="w-10 h-10 object-cover rounded-full cursor-pointer hover:opacity-75 transition-opacity" 
                        @click="showImageDialog(slotProps.data.photo_url)" 
                      />
                      <span v-else>No Photo</span>
                    </div>
                  </template>
                </Column>
                <Column header="Actions">
                  <template #body="slotProps">
                    <div v-if="slotProps.data">
                      <Button v-tooltip.top="'Delete Permanently'" icon="pi pi-trash" class="p-button-rounded p-button-danger" @click="confirmPermanentDelete(slotProps.data)" />
                    </div>
                  </template>
                </Column>
              </DataTable>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </template>
    </Card>

    <!-- Dialogs -->
    <Dialog v-model:visible="imageDialogVisible" modal header="Tower Photo" :style="{ width: '50vw' }" :breakpoints="{ '960px': '75vw', '640px': '90vw' }" >
      <img :src="displayedImageUrl" alt="Tower Photo" class="w-full h-auto" />
    </Dialog>

    <Dialog v-model:visible="towerDialog" :style="{ width: '600px' }" header="Tower Details" :modal="true" class="p-fluid" @show="initDialogMap" @hide="destroyDialogMap">
      <div class="field my-4">
        <FloatLabel  variant="on">
          <InputText id="latitude" v-model="tower.latitude" type="number" step="any" :invalid="!!errors.latitude" fluid />
          <label for="latitude">Latitude</label>
        </FloatLabel>
        <small v-if="errors.latitude" class="p-error">{{ errors.latitude }}</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="longitude" v-model="tower.longitude" type="number" step="any" :invalid="!!errors.longitude" fluid />
          <label for="longitude">Longitude</label>
        </FloatLabel>
        <small v-if="errors.longitude" class="p-error">{{ errors.longitude }}</small>
      </div>
      <div class="field my-4">
        <label>Location on Map</label>
        <div id="dialogMap" style="height: 300px; width: 100%;"></div>
      </div>
       <div class="field my-4">
        <FloatLabel variant="on" >
          <InputText id="kecamatan" v-model="tower.kecamatan" :invalid="!!errors.kecamatan" fluid />
          <label for="kecamatan">Kecamatan</label>
        </FloatLabel>
        <small v-if="errors.kecamatan" class="p-error">{{ errors.kecamatan }}</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <InputText id="kelurahan" v-model="tower.kelurahan" :invalid="!!errors.kelurahan" fluid />
          <label for="kelurahan">Kelurahan</label>
        </FloatLabel>
        <small v-if="errors.kelurahan" class="p-error">{{ errors.kelurahan }}</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <InputText id="address" v-model="tower.address" fluid />
          <label for="address">Alamat</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="tipe" v-model="tower.tipe" fluid />
          <label for="tipe">Tipe Tower</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <InputNumber id="tinggi" v-model="tower.tinggi" mode="decimal" :minFractionDigits="1" :maxFractionDigits="2" fluid/>
          <label for="tinggi">Tinggi Tower (meter)</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <label for="photo">Photo</label>
        <FileUpload mode="basic" name="photo" accept="image/*" @select="onFileSelect" fluid />
        <img v-if="tower.photo_url" :src="tower.photo_url" :alt="tower.ID" class="mt-2 w-20 h-20 object-cover" />
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <MultiSelect id="providers" v-model="selectedProviders" :options="providers" optionLabel="name" display="chip" dataKey="ID" fluid />
          <label for="providers">Providers</label>
        </FloatLabel>
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideTowerDialog" />
        <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveTower" />
      </template>
    </Dialog>

    <Dialog v-model:visible="relocateTowerDialog" :style="{ width: '450px' }" header="Relocate Tower" :modal="true" class="p-fluid" @show="initRelocateMap" @hide="destroyRelocateMap">
      <div class="confirmation-content">
        <i class="pi pi-map-marker mr-3" style="font-size: 2rem" />
        <span v-if="towerToRelocate">Relocate tower <b>{{ towerToRelocate.kelurahan }}</b>?</span>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <InputText id="relocateLatitude" v-model="relocateInput.latitude" type="number" step="any" fluid />
          <label for="relocateLatitude">New Latitude</label>
        </FloatLabel>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="relocateLongitude" v-model="relocateInput.longitude" type="number" step="any" fluid />
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

    <Dialog v-model:visible="changeOwnershipDialog" :style="{ width: '450px' }" header="Change Ownership" :modal="true" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-users mr-3" style="font-size: 2rem" />
        <span v-if="towerToChangeOwnership">Change ownership for tower <b>{{ towerToChangeOwnership.ID }}</b>?</span>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on" >
          <Dropdown id="newProvider" v-model="changeOwnershipInput.newProviderId" :options="providers" optionLabel="name" optionValue="ID" placeholder="Select New Provider" fluid />
          <label for="newProvider">New Provider</label>
        </FloatLabel>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="changeOwnershipDialog = false" />
        <Button label="Change" icon="pi pi-check" class="p-button-text" @click="changeOwnership" />
      </template>
    </Dialog>

    <Dialog v-model:visible="dismantleTowerDialog" :style="{ width: '450px' }" header="Confirm Dismantle" :modal="true">
      <div class="flex items-center">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="towerToDismantle">Are you sure you want to change status of tower <b>{{ towerToDismantle.ID }}</b> to dismantled?</span>
      </div>
      <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="dismantleTowerDialog = false" />
        <Button label="Yes, Dismantle" icon="pi pi-check" class="p-button-text" @click="dismantleTower" />
      </template>
    </Dialog>

    <Dialog v-model:visible="permanentDeleteDialog" :style="{ width: '450px' }" header="Confirm Permanent Delete" :modal="true">
      <div class="flex items-center">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="towerToDelete">Are you sure you want to permanently delete tower <b>{{ towerToDelete.ID }}</b>? This action cannot be undone.</span>
      </div>
      <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="permanentDeleteDialog = false" />
        <Button label="Yes, Delete Permanently" icon="pi pi-trash" class="p-button-danger" @click="executePermanentDelete" />
      </template>
    </Dialog>

  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useDataStore } from '../stores/data';
import { useUiStore } from '../stores/ui';
import { useToast } from 'primevue/usetoast';
import { z } from 'zod';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Card from 'primevue/card';
import ToggleSwitch from 'primevue/toggleswitch';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import FileUpload from 'primevue/fileupload';
import MultiSelect from 'primevue/multiselect';
import Dropdown from 'primevue/dropdown';
import FloatLabel from 'primevue/floatlabel';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Fix for default marker icon in Leaflet with Vite
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';
import iconRetina from 'leaflet/dist/images/marker-icon-2x.png';

let DefaultIcon = L.icon({
    iconUrl: icon,
    iconRetinaUrl: iconRetina,
    shadowUrl: iconShadow,
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    tooltipAnchor: [16, -28],
    shadowSize: [41, 41]
});

L.Marker.prototype.options.icon = DefaultIcon;

// Zod Schema for Tower Validation
const towerSchema = z.object({
  latitude: z.coerce.number().min(-90, 'Invalid latitude').max(90, 'Invalid latitude'),
  longitude: z.coerce.number().min(-180, 'Invalid longitude').max(180, 'Invalid longitude'),
  kecamatan: z.string().min(1, 'Kecamatan is required'),
  kelurahan: z.string().min(1, 'Kelurahan is required'),
});

// Stores
const dataStore = useDataStore();
const uiStore = useUiStore();

// Store State
const { towers, providers } = storeToRefs(dataStore);

// Computed properties for separating towers by status
const activeTowers = computed(() => 
  Array.isArray(towers.value) ? towers.value.filter(t => t && t.status !== 'dismantled') : []
);

const dismantledTowers = computed(() => 
  Array.isArray(towers.value) ? towers.value.filter(t => t && t.status === 'dismantled') : []
);

// Store Actions
const { fetchTowers, createTower, updateTower, dismantleTower: dismantleTowerAction, permanentDeleteTower, changeTowerOwnership: changeTowerOwnershipAction, relocateTower: relocateTowerAction, fetchProviders } = dataStore;
const { toggleDarkMode } = uiStore;

// UI & State
const toast = useToast();
const towerDialog = ref(false);
const relocateTowerDialog = ref(false);
const changeOwnershipDialog = ref(false);
const dismantleTowerDialog = ref(false);
const permanentDeleteDialog = ref(false);
const tower = ref({});
const errors = ref({});
const towerToRelocate = ref(null);
const towerToChangeOwnership = ref(null);
const towerToDismantle = ref(null);
const towerToDelete = ref(null);
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
let debounceTimer = null;

// Lifecycle
onMounted(() => {
  fetchTowers();
  fetchProviders();
});

// Watchers for Map Interaction
watch([() => tower.value.latitude, () => tower.value.longitude], ([newLat, newLon]) => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    const lat = parseFloat(newLat);
    const lon = parseFloat(newLon);
    if (!isNaN(lat) && !isNaN(lon) && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180) {
      updateDialogMarker();
    }
  }, 500); // 500ms debounce
});

// Methods for Tower CRUD (Create/Edit)
const openNewTowerDialog = () => {
  tower.value = { latitude: null, longitude: null, photo_url: '', Status: 'active', kelurahan: '', kecamatan: '', address: '', tipe: '', tinggi: null };
  selectedProviders.value = [];
  photoFile.value = null;
  submitted.value = false;
  errors.value = {};
  towerDialog.value = true;
};

const hideTowerDialog = () => {
  towerDialog.value = false;
  submitted.value = false;
};

const saveTower = async () => {
  submitted.value = true;
  errors.value = {}; // Clear previous errors

  const result = towerSchema.safeParse(tower.value);

  if (!result.success) {
    errors.value = result.error.flatten().fieldErrors;
    toast.add({ severity: 'error', summary: 'Validation Error', detail: 'Please check the required fields.', life: 3000 });
    return;
  }

  const isUpdate = !!tower.value.ID;
  const formData = new FormData();

  // Append fields for both create and update
  formData.append('kelurahan', tower.value.kelurahan || '');
  formData.append('kecamatan', tower.value.kecamatan || '');
  formData.append('address', tower.value.address || '');
  formData.append('tipe', tower.value.tipe || '');
  formData.append('tinggi', tower.value.tinggi || 0);
  if (photoFile.value) {
    formData.append('photo', photoFile.value);
  }

  // Append fields only for create
  if (!isUpdate) {
    formData.append('latitude', tower.value.latitude);
    formData.append('longitude', tower.value.longitude);
    selectedProviders.value.forEach(provider => {
      formData.append('provider_ids', provider.ID);
    });
  }

  try {
    if (isUpdate) {
      await updateTower(tower.value.ID, formData);
      toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Details Updated', life: 3000 });
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
};

const editTower = (t) => {
  tower.value = { ...t };
  errors.value = {};
  selectedProviders.value = []; 
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

// Methods for Permanent Delete
const confirmPermanentDelete = (t) => {
  towerToDelete.value = { ...t };
  permanentDeleteDialog.value = true;
};

const executePermanentDelete = async () => {
  try {
    await permanentDeleteTower(towerToDelete.value.ID);
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Tower Permanently Deleted', life: 3000 });
    permanentDeleteDialog.value = false;
    towerToDelete.value = null;
    fetchTowers();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Permanent delete failed', life: 3000 });
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
      // Watcher will handle the marker update
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
    const lat = parseFloat(tower.value.latitude);
    const lon = parseFloat(tower.value.longitude);
    if (isNaN(lat) || isNaN(lon)) return;

    const latlng = [lat, lon];
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
    relocateMapInstance.invalidateSize(); // Force map to re-calculate its size
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