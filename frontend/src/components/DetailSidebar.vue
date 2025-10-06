<template>
  <Drawer
    v-model:visible="isVisible"
    position="right"
    class="md:w-25rem lg:w-35rem"
    @hide="$emit('close')"
  >
    <template #header>
      <div class="flex items-center gap-2">
        <i :class="headerIcon" class="text-xl"></i>
        <span class="font-bold text-xl">{{ headerTitle }}</span>
      </div>
    </template>

    <div v-if="item">
      <Tabs value="0">
        <TabList>
          <Tab value="0">Details</Tab>
          <Tab value="1" v-if="isTower">History</Tab>
        </TabList>
        <TabPanels>
          <TabPanel value="0">
            <div class="flex flex-col gap-4 p-3">
            <!-- Photo/Icon display at the top -->
            <div v-if="isTower" class="flex flex-col items-center justify-center bg-gray-100 dark:bg-gray-800 p-4 rounded-md h-48">
              <img 
                v-if="item.photo_url" 
                :src="item.photo_url" 
                alt="Tower Photo" 
                class="w-full h-full object-cover rounded-md cursor-pointer hover:opacity-75 transition-opacity" 
                @click="showImageDialog(item.photo_url)"
              />
              <template v-else>
                <i class="pi pi-image text-6xl text-gray-400"></i>
                <span class="text-gray-500 mt-2">No Photo Available</span>
              </template>
            </div>

            <!-- Image Preview Dialog -->
            <Dialog v-model:visible="imageDialogVisible" modal header="Tower Photo" :style="{ width: '50vw' }" :breakpoints="{ '960px': '75vw', '640px': '90vw' }" >
              <img :src="displayedImageUrl" alt="Tower Photo" class="w-full h-auto" />
            </Dialog>

              <template v-if="isTower">
                <div class="grid grid-cols-2 gap-4">
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Latitude</span>
                    <span class="text-lg">{{ item.latitude || 'N/A' }}</span>
                  </div>
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Longitude</span>
                    <span class="text-lg">{{ item.longitude || 'N/A' }}</span>
                  </div>
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Kecamatan</span>
                    <span class="text-lg">{{ item.kecamatan || 'N/A' }}</span>
                  </div>
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Kelurahan</span>
                    <span class="text-lg">{{ item.kelurahan || 'N/A' }}</span>
                  </div>
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Tipe Tower</span>
                    <span class="text-lg">{{ item.tipe || 'N/A' }}</span>
                  </div>
                  <div class="flex flex-col">
                    <span class="font-semibold text-sm text-gray-500">Tinggi</span>
                    <span class="text-lg">{{ item.tinggi ? item.tinggi + ' m' : 'N/A' }}</span>
                  </div>
                </div>
                <div class="flex flex-col">
                  <span class="font-semibold text-sm text-gray-500">Status</span>
                  <span class="text-lg">{{ item.status || 'N/A' }}</span>
                </div>
                <div class="flex flex-col">
                  <span class="font-semibold text-sm text-gray-500">Providers</span>
                  <span v-if="item.providers && item.providers.length" class="text-lg">
                      {{ item.providers.map(p => p.name).join(', ') }}
                  </span>
                  <span v-else class="text-lg text-gray-400">
                      No provider information.
                  </span>
                </div>
              </template>

              <template v-else-if="isBlankspot">
                <div v-for="field in fields" :key="field.key" class="flex flex-col">
                  <span class="font-semibold text-sm text-gray-500">{{ field.label }}</span>
                  <div v-if="field.key === 'color'" class="flex items-center gap-2">
                    <span :style="{ backgroundColor: item[field.key] }" class="inline-block w-4 h-4 rounded-full"></span>
                    <span>{{ item[field.key] }}</span>
                  </div>
                  <span v-else class="text-lg">{{ item[field.key] || 'N/A' }}</span>
                </div>
              </template>
            </div>
          </TabPanel>
          <TabPanel value="1" v-if="isTower">
            <div class="p-3">
              <Timeline :value="formattedTowerHistory" align="alternate" class="custom-timeline">
                <template #marker="slotProps">
                  <span class="flex w-2rem h-2rem align-items-center justify-content-center text-white border-circle z-1 shadow-1" :style="{ backgroundColor: getEventColor(slotProps.item.eventType) }">
                    <i :class="getEventIcon(slotProps.item.eventType)"></i>
                  </span>
                </template>
                <template #content="slotProps">
                  <Card>
                    <template #title>{{ slotProps.item.description }}</template>
                    <template #subtitle>{{ formatDate(slotProps.item.timestamp) }}</template>
                    <template #content>
                      <p class="m-0">{{ slotProps.item.details }}</p>
                    </template>
                  </Card>
                </template>
              </Timeline>
              <p v-if="formattedTowerHistory.length === 0" class="text-center text-gray-500">No history available for this tower.</p>
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </div>
    <div v-else>
      <p class="p-3">No item selected.</p>
    </div>
  </Drawer>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useDataStore } from '../stores/data'; // Import data store
import Drawer from 'primevue/drawer';
import Panel from 'primevue/panel'; // Keep Panel for now, might remove later if not used
import Tabs from 'primevue/tabs';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Timeline from 'primevue/timeline'; // New import
import Card from 'primevue/card'; // New import for timeline items
import Dialog from 'primevue/dialog';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  item: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['close']);

// State for Image Preview Dialog
const imageDialogVisible = ref(false);
const displayedImageUrl = ref('');

const showImageDialog = (url) => {
  displayedImageUrl.value = url;
  imageDialogVisible.value = true;
};

const dataStore = useDataStore(); // Initialize data store

const isVisible = ref(props.visible);
const towerHistory = ref([]); // To store fetched tower history

watch(() => props.visible, (newVal) => {
  isVisible.value = newVal;
  if (newVal && isTower.value && props.item.ID) {
    fetchTowerHistory(props.item.ID);
  } else if (!newVal) {
    towerHistory.value = []; // Clear history when sidebar closes
  }
});

watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('close');
  }
});

const isTower = computed(() => props.item && props.item.hasOwnProperty('latitude'));
const isBlankspot = computed(() => props.item && props.item.hasOwnProperty('kelurahan'));

const headerIcon = computed(() => {
  if (isTower.value) return 'pi pi-signal';
  if (isBlankspot.value) return 'pi pi-map-marker';
  return 'pi pi-info-circle';
});

const headerTitle = computed(() => {
  if (isTower.value) return `Tower ${props.item.ID}`;
  if (isBlankspot.value) return props.item.name;
  return 'Details';
});

const fields = computed(() => {
  if (isTower.value) {
    return [
      { key: 'latitude', label: 'Latitude' },
      { key: 'longitude', label: 'Longitude' },
      { key: 'status', label: 'Status' },
      { key: 'providers', label: 'Providers' }
    ];
  }
  if (isBlankspot.value) {
    return [
      { key: 'name', label: 'Area Name' },
      { key: 'type', label: 'Type' },
      { key: 'kelurahan', label: 'Kelurahan' },
      { key: 'color', label: 'Color' }
    ];
  }
  return [];
});

// Fetch tower history
const fetchTowerHistory = async (towerId) => {
  try {
    towerHistory.value = await dataStore.fetchTowerHistory(towerId);
  } catch (error) {
    console.error('Failed to fetch tower history:', error);
    // Optionally show a toast error
  }
};

// Format history for Timeline component
const formattedTowerHistory = computed(() => {
  return towerHistory.value.map(event => {
    let details = '';
    try {
      const oldData = event.OldData ? JSON.parse(event.OldData) : {};
      const newData = event.NewData ? JSON.parse(event.NewData) : {};

      switch (event.EventType) {
        case 'Created':
          details = `Initial creation at Lat: ${newData.latitude}, Long: ${newData.longitude}. Status: ${newData.status}.`;
          break;
        case 'OwnershipChange':
          const oldProviders = oldData.providers ? oldData.providers.map(p => p.name).join(', ') : 'N/A';
          const newProviders = newData.providers ? newData.providers.map(p => p.name).join(', ') : 'N/A';
          details = `Providers changed from ${oldProviders} to ${newProviders}.`;
          break;
        case 'Relocation':
          details = `Moved from Lat: ${oldData.latitude}, Long: ${oldData.longitude} to Lat: ${newData.latitude}, Long: ${newData.longitude}.`;
          break;
        case 'Dismantled':
          details = `Tower was dismantled. Old status: ${oldData.status}.`;
          break;
        case 'PhotoUpdate':
          details = `A new photo was uploaded for the tower.`;
          break;
        default:
          details = event.Description;
      }
    } catch (e) {
      console.error('Error parsing event data:', e);
      details = event.Description; // Fallback to description
    }
    return {
      ...event,
      details: details,
    };
  });
});

// Helper to get icon based on event type
const getEventIcon = (eventType) => {
  switch (eventType) {
    case 'Created': return 'pi pi-plus';
    case 'OwnershipChange': return 'pi pi-users';
    case 'Relocation': return 'pi pi-map-marker';
    case 'Dismantled': return 'pi pi-trash';
    case 'PhotoUpdate': return 'pi pi-image';
    default: return 'pi pi-info-circle';
  }
};

// Helper to get color based on event type
const getEventColor = (eventType) => {
  switch (eventType) {
    case 'Created': return '#2196F3'; // Blue
    case 'OwnershipChange': return '#FFC107'; // Amber
    case 'Relocation': return '#4CAF50'; // Green
    case 'Dismantled': return '#F44336'; // Red
    case 'PhotoUpdate': return '#9C27B0'; // Purple
    default: return '#607D8B'; // Blue Grey
  }
};

// Helper to format date
const formatDate = (date) => {
  return new Date(date).toLocaleString();
};
</script>

<style scoped>
/* Custom styles for timeline if needed */
</style>