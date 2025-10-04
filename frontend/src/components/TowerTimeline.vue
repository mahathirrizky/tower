<template>
  <div class="timeline-container">
    <Timeline :value="events" align="alternate" class="customized-timeline">
      <template #marker="slotProps">
        <span class="custom-marker shadow-2" :style="{ backgroundColor: getColor(slotProps.item.event_type) }">
          <i :class="getIcon(slotProps.item.event_type)"></i>
        </span>
      </template>
      <template #content="slotProps">
        <Card>
          <template #title>
            {{ formatTitle(slotProps.item.event_type) }}
          </template>
          <template #subtitle>
            {{ new Date(slotProps.item.timestamp).toLocaleString() }}
          </template>
          <template #content>
            <p>{{ slotProps.item.description }}</p>
          </template>
        </Card>
      </template>
    </Timeline>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import Timeline from 'primevue/timeline';
import Card from 'primevue/card';

const props = defineProps({
  events: {
    type: Array,
    required: true,
    default: () => []
  }
});

const getIcon = (eventType) => {
  switch (eventType) {
    case 'Created':
      return 'pi pi-plus';
    case 'OwnershipChange':
      return 'pi pi-users';
    case 'Relocation':
      return 'pi pi-map-marker';
    case 'Dismantled':
      return 'pi pi-trash';
    case 'PhotoUpdate':
      return 'pi pi-image';
    default:
      return 'pi pi-info-circle';
  }
};

const getColor = (eventType) => {
  switch (eventType) {
    case 'Created':
      return '#4CAF50'; // Green
    case 'OwnershipChange':
      return '#2196F3'; // Blue
    case 'Relocation':
      return '#FFC107'; // Amber
    case 'Dismantled':
      return '#F44336'; // Red
    case 'PhotoUpdate':
      return '#9C27B0'; // Purple
    default:
      return '#607D8B'; // Blue Grey
  }
};

const formatTitle = (eventType) => {
  // Add spaces before uppercase letters to format the title
  return eventType.replace(/([A-Z])/g, ' $1').trim();
};
</script>

<style scoped>
.timeline-container {
  padding: 1rem;
}

.custom-marker {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  border-radius: 50%;
  width: 2rem;
  height: 2rem;
}

.customized-timeline .p-timeline-event:nth-child(even) {
  flex-direction: row-reverse;
}

.customized-timeline .p-timeline-event-content,
.customized-timeline .p-timeline-event-opposite {
  line-height: 1;
}
</style>
