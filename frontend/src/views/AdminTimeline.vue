<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Tower Timeline</h2>
    <Card>
      <template #content>
        <div class="flex flex-col gap-4">
          <FloatLabel variant="on">

            <Select
            v-model="selectedTower"
            :options="towers"
            optionLabel="ID"

            @change="onTowerSelect"
            fluid
            >
            <template #option="slotProps">
              <span>Tower ID: {{ slotProps.option.ID }}</span>
            </template>
          </Select>
          <label for="latitude">Pilih Tower</label>
        </FloatLabel>

          <div v-if="isLoading" class="text-center">
            <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
            <p>Loading timeline...</p>
          </div>

          <div v-else-if="timelineEvents.length > 0">
            <TowerTimeline :events="timelineEvents" />
          </div>

          <div v-else-if="selectedTower" class="text-center text-gray-500">
            <p>No timeline events found for this tower.</p>
          </div>

          <div v-else class="text-center text-gray-500">
            <p>Please select a tower to view its timeline.</p>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useDataStore } from '../stores/data';
import { storeToRefs } from 'pinia';
import Select from 'primevue/select';
import Card from 'primevue/card';
import TowerTimeline from '../components/TowerTimeline.vue';
import FloatLabel from 'primevue/floatlabel';
const dataStore = useDataStore();
const { towers } = storeToRefs(dataStore);
const { fetchTowers, fetchTowerHistory } = dataStore;

const selectedTower = ref(null);
const timelineEvents = ref([]);
const isLoading = ref(false);

onMounted(() => {
  if (towers.value.length === 0) {
    fetchTowers();
  }
});

const onTowerSelect = async (event) => {
  if (event.value) {
    isLoading.value = true;
    timelineEvents.value = [];
    try {
      timelineEvents.value = await fetchTowerHistory(event.value.ID);
    } catch (error) {
      console.error("Failed to fetch tower history:", error);
    } finally {
      isLoading.value = false;
    }
  }
};
</script>