<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-3xl font-bold">Tower Locations</h1>
      <div class="flex items-center gap-2">
        <i class="pi pi-sun" :class="{ 'text-gray-400': isDarkMode }"></i>
        <ToggleSwitch :modelValue="isDarkMode" @update:modelValue="toggleDarkMode" />
        <i class="pi pi-moon" :class="{ 'text-gray-400': !isDarkMode }"></i>
      </div>
    </div>
    <div id="map" class="h-[calc(100vh-10rem)] w-full"></div>

    <DetailSidebar :visible="isDetailSidebarVisible" :item="detailItem" @close="closeDetailSidebar" />
  </div>
</template>

<script setup>
import { onMounted, computed, watch, nextTick } from 'vue';
import { storeToRefs } from 'pinia';
import { useDataStore } from '../stores/data';
import { useUiStore } from '../stores/ui';
import { useBlankspotStore } from '../stores/blankspot'; // Import blankspot store
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import ToggleSwitch from 'primevue/toggleswitch';
import DetailSidebar from '../components/DetailSidebar.vue'; // Import DetailSidebar

// Stores
const dataStore = useDataStore();
const uiStore = useUiStore();
const blankspotStore = useBlankspotStore(); // Initialize blankspot store

// Store State
const { towers } = storeToRefs(dataStore);
const { blankspotAreas } = storeToRefs(blankspotStore); // Get blankspot areas
const { isDetailSidebarVisible, detailItem } = storeToRefs(uiStore); // Get sidebar state from UI store

// Store Actions
const { fetchTowers } = dataStore;
const { toggleDarkMode, openDetailSidebar, closeDetailSidebar } = uiStore; // Get sidebar actions from UI store
const { fetchBlankspotAreas } = blankspotStore; // Get fetch action for blankspots

// Map instance & layers
let map = null;
let markerLayerGroup = null;
let blankspotLayerGroup = null; // New layer group for blankspots

// Custom Icon
const towerIcon = L.icon({
  iconUrl: '/tower-icon.svg',
  iconSize: [48, 48],
  iconAnchor: [24, 48],
  popupAnchor: [0, -48]
});


// Computed
const isDarkMode = computed(() => uiStore.isDarkMode);

const initMap = () => {
  map = L.map('map').setView([-8.4601, 118.7267], 13); // Centered on Bima
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);
  markerLayerGroup = L.layerGroup().addTo(map);
  blankspotLayerGroup = L.layerGroup().addTo(map); // Add blankspot layer group
};

const addMarkers = (newTowers) => {
  if (!markerLayerGroup) return;
  markerLayerGroup.clearLayers(); // Clear existing markers
  newTowers.forEach(tower => {
    const marker = L.marker([tower.latitude, tower.longitude], { icon: towerIcon })
      .addTo(markerLayerGroup);
    marker.on('click', () => {
      openDetailSidebar(tower);
    });
  });
};

const addBlankspotAreasToMap = (areas) => {
  if (!blankspotLayerGroup) return;
  blankspotLayerGroup.clearLayers(); // Clear existing blankspot areas
  areas.forEach(area => {
    try {
      const coordinates = JSON.parse(area.coordinates);
      if (coordinates && coordinates.length > 0) {
        const polygon = L.polygon(coordinates, {
          color: area.color,
          fillColor: area.color,
          fillOpacity: 0.3,
          weight: 2
        }).addTo(blankspotLayerGroup);
        polygon.on('click', () => {
          openDetailSidebar(area);
        });
      } else {
        console.warn(`Blankspot area ${area.name} has no valid coordinates.`);
      }
    } catch (e) {
      console.error("Error parsing blankspot coordinates:", e);
    }
  });
};

// Lifecycle
onMounted(() => {
  nextTick(() => {
    initMap();
    fetchTowers();
    fetchBlankspotAreas(); // Fetch blankspot areas
  });
});

// Watch for tower data changes
watch(towers, (newTowers) => {
  addMarkers(newTowers);
});

// Watch for blankspot area data changes
watch(blankspotAreas, (newAreas) => {
  addBlankspotAreasToMap(newAreas);
});
</script>
