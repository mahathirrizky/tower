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
        const response = await apiClient.get('/towers');
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
        this.towers.push(response.data.data);
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
        const index = this.towers.findIndex(tower => tower.ID === id);
        if (index !== -1) {
          this.towers[index] = response.data.data;
        }
        return response.data;
      } catch (error) {
        console.error('Error updating tower:', error);
        throw error;
      }
    },
    async deleteTower(id) {
      try {
        await apiClient.delete(`/towers/${id}`);
        this.towers = this.towers.filter(tower => tower.ID !== id);
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
        this.providers.push(response.data.data);
        return response.data;
      } catch (error) {
        console.error('Error creating provider:', error);
        throw error;
      }
    },
    async updateProvider(id, providerData) {
      try {
        const response = await apiClient.put(`/providers/${id}`, providerData);
        const index = this.providers.findIndex(provider => provider.ID === id);
        if (index !== -1) {
          this.providers[index] = response.data.data;
        }
        return response.data;
      } catch (error) {
        console.error('Error updating provider:', error);
        throw error;
      }
    },
    async deleteProvider(id) {
      try {
        await apiClient.delete(`/providers/${id}`);
        this.providers = this.providers.filter(provider => provider.ID !== id);
      } catch (error) {
        console.error('Error deleting provider:', error);
        throw error;
      }
    },

    async changeTowerOwnership(towerId, newProviderId) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/ownership`, { new_provider_id: newProviderId });
        // Update the specific tower in the state
        const index = this.towers.findIndex(tower => tower.ID === towerId);
        if (index !== -1) {
          this.towers[index] = response.data.data; // Assuming backend returns updated tower
        }
        return response.data;
      } catch (error) {
        console.error('Error changing tower ownership:', error);
        throw error;
      }
    },
    async relocateTower(towerId, latitude, longitude) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/relocate`, { latitude, longitude });
        // Update the specific tower in the state
        const index = this.towers.findIndex(tower => tower.ID === towerId);
        if (index !== -1) {
          this.towers[index] = response.data.data; // Assuming backend returns updated tower
        }
        return response.data;
      } catch (error) {
        console.error('Error relocating tower:', error);
        throw error;
      }
    },
    async dismantleTower(towerId) {
      try {
        const response = await apiClient.put(`/towers/${towerId}/dismantle`);
        // Update the specific tower in the state (status to dismantled)
        const index = this.towers.findIndex(tower => tower.ID === towerId);
        if (index !== -1) {
          this.towers[index] = response.data.data; // Assuming backend returns updated tower with status
        }
        // Optionally remove from active towers list if UI only shows active
        // this.towers = this.towers.filter(tower => tower.ID !== towerId);
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
