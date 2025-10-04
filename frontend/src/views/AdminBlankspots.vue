<template>
  <div>
    
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Manage Blankspot Areas</h2>
      <div class="flex items-center gap-2">
        <Button label="New Blankspot Area" icon="pi pi-plus" class="p-button-success mr-2" @click="openNewBlankspotAreaDialog" />
        <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
        <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
        <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
      </div>
    </div>
    <Card class="w-full">
      <template #content>
        <DataTable :value="blankspotAreas" :paginator="true" :rows="5">
          <Column header="Nomor">
            <template #body="slotProps">
              {{ slotProps.index + 1 }}
            </template>
          </Column>
          <Column field="name" header="Name"></Column>
          <Column field="kelurahan" header="Kelurahan"></Column>
          <Column field="type" header="Type"></Column>
          <Column field="color" header="Color">
            <template #body="slotProps">
              <span :style="{ backgroundColor: slotProps.data.color }" class="inline-block w-4 h-4 rounded-full mr-2"></span>
              {{ slotProps.data.color }}
            </template>
          </Column>
          <Column header="Actions">
            <template #body="slotProps">
              <Button v-tooltip.top="'Edit Blankspot'" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" @click="editBlankspotArea(slotProps.data)" />
              <Button v-tooltip.top="'Delete Blankspot'" icon="pi pi-trash" class="p-button-rounded p-button-danger" @click="confirmDeleteBlankspotArea(slotProps.data)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <Dialog v-model:visible="blankspotAreaDialog" :style="{ width: '800px' }" header="Blankspot Area Details" :modal="true" class="p-fluid" @show="initMap" @hide="destroyMap">
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="name" v-model.trim="blankspotArea.name" required="true" autofocus :class="{ 'p-invalid': submitted && !blankspotArea.name }" fluid />
          <label for="name">Name</label>
        </FloatLabel>
        <small class="p-error" v-if="submitted && !blankspotArea.name">Name is required.</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <InputText id="kelurahan" v-model.trim="blankspotArea.kelurahan" required="true" :class="{ 'p-invalid': submitted && !blankspotArea.kelurahan }" fluid />
          <label for="kelurahan">Kelurahan</label>
        </FloatLabel>
        <small class="p-error" v-if="submitted && !blankspotArea.kelurahan">Kelurahan is required.</small>
      </div>
      <div class="field my-4">
        <FloatLabel variant="on">
          <Select id="type" fluid v-model="blankspotArea.type" :options="blankspotTypes" placeholder="Select a Type" required="true" :class="{ 'p-invalid': submitted && !blankspotArea.type }" @change="updateDefaultColor" />
          <label for="type">Type</label>
        </FloatLabel>
        <small class="p-error" v-if="submitted && !blankspotArea.type">Type is required.</small>
      </div>
      <!-- Removed ColorPicker section -->
      <div class="field my-4">
        <label>Draw Blankspot Area on Map</label>
        <div id="blankspotMap" style="height: 300px; width: 100%;"></div>
      </div>
      <div class="field my-4">
        <label>Coordinates (JSON)</label>
        <InputText id="coordinates" v-model="blankspotArea.coordinates" disabled fluid />
      </div>

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideBlankspotAreaDialog" />
        <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveBlankspotArea" />
      </template>
    </Dialog>

    <Dialog v-model:visible="deleteBlankspotAreaDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="blankspotArea">Are you sure you want to delete blankspot area <b>{{ blankspotArea.name }}</b>?</span>
      </div>
      <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="deleteBlankspotAreaDialog = false" />
        <Button label="Yes" icon="pi pi-check" class="p-button-text" @click="deleteBlankspotArea" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useBlankspotStore } from '../stores/blankspot';
import { useUiStore } from '../stores/ui';
import { useToast } from 'primevue/usetoast';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Card from 'primevue/card';
import ToggleSwitch from 'primevue/toggleswitch';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';

import FloatLabel from 'primevue/floatlabel';
import 'leaflet/dist/leaflet.css';
import 'leaflet-draw/dist/leaflet.draw.css'; // Import Leaflet.draw CSS
import L from 'leaflet';
import 'leaflet-draw'; // Import Leaflet.draw JS

// Stores
const blankspotStore = useBlankspotStore();
const uiStore = useUiStore();

// Store State
const { blankspotAreas } = storeToRefs(blankspotStore);

// Store Actions
const { fetchBlankspotAreas, createBlankspotArea, updateBlankspotArea, deleteBlankspotArea: deleteBlankspotAreaAction } = blankspotStore;
const { toggleDarkMode } = uiStore;

// UI & State
const toast = useToast();
const blankspotAreaDialog = ref(false);
const deleteBlankspotAreaDialog = ref(false);
const blankspotArea = ref({});
const submitted = ref(false);
const blankspotTypes = ref(['Blankspot', 'Weak Signal']); // Options for the Type dropdown

const isDarkMode = computed(() => uiStore.isDarkMode);

// Map instance
let mapInstance = null;
let editableLayer = null; // Layer to hold drawn items
let drawControl = null; // Leaflet.draw control

// Lifecycle
onMounted(() => {
  fetchBlankspotAreas();
});

// Watch for blankspotArea.type changes to update default color
watch(() => blankspotArea.value.type, (newType) => {
  updateDefaultColor(newType);
});

// Methods
const updateDefaultColor = (type) => {
  if (type === 'Blankspot') {
    blankspotArea.value.color = '#FF0000'; // Red
  } else if (type === 'Weak Signal') {
    blankspotArea.value.color = '#FFFF00'; // Yellow
  } else {
    blankspotArea.value.color = '#FF0000'; // Default to red
  }
  // Update color of current polygon if drawing
  if (editableLayer && editableLayer.getLayers().length > 0) {
    editableLayer.eachLayer(layer => {
      if (layer instanceof L.Polygon) {
        layer.setStyle({ color: blankspotArea.value.color, fillColor: blankspotArea.value.color });
      }
    });
  }
};

const openNewBlankspotAreaDialog = () => {
  blankspotArea.value = { name: '', kelurahan: '', coordinates: '[]', type: 'Blankspot', color: '#FF0000' }; // Default type and color
  submitted.value = false;
  blankspotAreaDialog.value = true;
};

const hideBlankspotAreaDialog = () => {
  blankspotAreaDialog.value = false;
  submitted.value = false;
  // Clear drawn items when dialog is hidden
  if (editableLayer) {
    editableLayer.clearLayers();
  }
  // Remove draw control
  if (drawControl) {
    mapInstance.removeControl(drawControl);
    drawControl = null;
  }
};

const saveBlankspotArea = async () => {
  submitted.value = true;
  if (blankspotArea.value.name && blankspotArea.value.kelurahan && blankspotArea.value.coordinates && blankspotArea.value.type) {
    try {
      if (blankspotArea.value.ID) {
        await updateBlankspotArea(blankspotArea.value.ID, blankspotArea.value);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Blankspot Area Updated', life: 3000 });
      } else {
        await createBlankspotArea(blankspotArea.value);
        toast.add({ severity: 'success', summary: 'Successful', detail: 'Blankspot Area Created', life: 3000 });
      }
      blankspotAreaDialog.value = false;
      blankspotArea.value = {};
      fetchBlankspotAreas(); // Refresh the list
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Operation failed', life: 3000 });
    }
  }
};

const editBlankspotArea = (area) => {
  blankspotArea.value = { ...area };
  blankspotAreaDialog.value = true;
};

const confirmDeleteBlankspotArea = (area) => {
  blankspotArea.value = area;
  deleteBlankspotAreaDialog.value = true;
};

const deleteBlankspotArea = async () => {
  try {
    await deleteBlankspotAreaAction(blankspotArea.value.ID);
    deleteBlankspotAreaDialog.value = false;
    blankspotArea.value = {};
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Blankspot Area Deleted', life: 3000 });
    fetchBlankspotAreas(); // Refresh the list
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Deletion failed', life: 3000 });
  }
};

// Map Methods
const initMap = () => {
  nextTick(() => {
    if (mapInstance) {
      mapInstance.remove();
    }
    mapInstance = L.map('blankspotMap').setView([-8.4601, 118.7267], 13); // Centered on Bima
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(mapInstance);

    editableLayer = new L.FeatureGroup();
    mapInstance.addLayer(editableLayer);

    // Initialize draw control
    drawControl = new L.Control.Draw({
      edit: {
        featureGroup: editableLayer,
        poly: {
          allowIntersection: false
        }
      },
      draw: {
        polygon: {
          allowIntersection: false,
          showArea: true,
          shapeOptions: {
            color: blankspotArea.value.color,
            fillColor: blankspotArea.value.color,
            fillOpacity: 0.3,
            opacity: 0.5
          }
        },
        polyline: false,
        rectangle: false,
        circle: false,
        marker: false,
        circlemarker: false
      }
    });
    mapInstance.addControl(drawControl);

    // Load existing polygon if editing
    if (blankspotArea.value.coordinates && blankspotArea.value.coordinates !== '[]') {
      try {
        const coords = JSON.parse(blankspotArea.value.coordinates);
        const polygon = L.polygon(coords, { color: blankspotArea.value.color, fillColor: blankspotArea.value.color, fillOpacity: 0.3, opacity: 0.5 });
        editableLayer.addLayer(polygon);
        mapInstance.fitBounds(polygon.getBounds());
      } catch (e) {
        console.error("Error parsing coordinates:", e);
      }
    }

    // Handle draw events
    mapInstance.on(L.Draw.Event.CREATED, (e) => {
      const layer = e.layer;
      editableLayer.addLayer(layer);
      blankspotArea.value.coordinates = JSON.stringify(layer.getLatLngs()[0].map(ll => [ll.lat, ll.lng]));
    });

    mapInstance.on(L.Draw.Event.EDITED, (e) => {
      e.layers.eachLayer(layer => {
        if (layer instanceof L.Polygon) {
          blankspotArea.value.coordinates = JSON.stringify(layer.getLatLngs()[0].map(ll => [ll.lat, ll.lng]));
        }
      });
    });

    mapInstance.on(L.Draw.Event.DELETED, () => {
      blankspotArea.value.coordinates = '[]';
    });
  });
};

const destroyMap = () => {
  if (mapInstance) {
    mapInstance.remove();
    mapInstance = null;
    editableLayer = null;
    drawControl = null;
  }
};
</script>
