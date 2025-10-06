import { defineStore } from 'pinia';
import apiClient from '../api/axios';

export const useDataStore = defineStore('data', {
  state: () => ({
    towers: [],
    providers: [],
  }),
  actions: {
    async fetchTowers() {
      try {
        const response = await apiClient.get(`/towers?_=${new Date().getTime()}`);
        this.towers = response.data.data;
      } catch (error) {
        console.error('Error fetching towers:', error);
        throw error;
      }
    },
    async createTower(towerData) {
      try {
        const response = await apiClient.post('/towers', towerData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        // State is updated by calling fetchTowers() in the component
        return response.data;
      } catch (error) {
        console.error('Error creating tower:', error);
        throw error;
      }
    },
    async updateTower(id, towerData) {
      try {
        const response = await apiClient.put(`/towers/${id}`, towerData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        // State is updated by calling fetchTowers() in the component
        return response.data;
      } catch (error) {
        console.error('Error updating tower:', error);
        throw error;
      }
    },
    async deleteTower(id) {
      try {
        await apiClient.delete(`/towers/${id}`);
        // State is updated by calling fetchTowers() in the component
      } catch (error) {
        console.error('Error deleting tower:', error);
        throw error;
      }
    },

    async fetchProviders() {
      try {
        const response = await apiClient.get('/providers');
        this.providers = response.data.data;
      } catch (error) {
        console.error('Error fetching providers:', error);
        throw error;
      }
    },
    async createProvider(providerData) {
      try {
        const response = await apiClient.post('/providers', providerData);
        // For simplicity, we can refetch providers in the component as well
        return response.data;
      } catch (error) {
        console.error('Error creating provider:', error);
        throw error;
      }
    },
    async updateProvider(id, providerData) {
      try {
        const response = await apiClient.put(`/providers/${id}`, providerData);
        // For simplicity, we can refetch providers in the component as well
        return response.data;
      } catch (error) {
        console.error('Error updating provider:', error);
        throw error;
      }
    },
    async deleteProvider(id) {
      try {
        await apiClient.delete(`/providers/${id}`);
        // For simplicity, we can refetch providers in the component as well
      } catch (error) {
        console.error('Error deleting provider:', error);
        throw error;
      }
    },

    async permanentDeleteTower(towerId) {
      try {
        await apiClient.delete(`/towers/${towerId}`);
        // State is updated by calling fetchTowers() in the component
      } catch (error) {
        console.error('Error permanently deleting tower:', error);
        throw error;
      }
    },

    async changeTowerOwnership(towerId, newProviderId) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/ownership`, { new_provider_id: newProviderId });
        // State is updated by calling fetchTowers() in the component
        return response.data;
      } catch (error) {
        console.error('Error changing tower ownership:', error);
        throw error;
      }
    },
    async relocateTower(towerId, latitude, longitude) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/relocate`, { latitude, longitude });
        // State is updated by calling fetchTowers() in the component
        return response.data;
      } catch (error) {
        console.error('Error relocating tower:', error);
        throw error;
      }
    },
    async dismantleTower(towerId) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/dismantle`);
        // State is updated by calling fetchTowers() in the component
        return response.data;
      } catch (error) {
        console.error('Error dismantling tower:', error);
        throw error;
      }
    },
    async fetchTowerHistory(towerId) {
      try {
        const response = await apiClient.get(`/towers/${towerId}/history`);
        return response.data.data; // Returns array of TowerEvent
      } catch (error) {
        console.error('Error fetching tower history:', error);
        throw error;
      }
    },
  },
});