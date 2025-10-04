import { defineStore } from 'pinia';
import apiClient from '../api/axios';

export const useBlankspotStore = defineStore('blankspot', {
  state: () => ({
    blankspotAreas: [],
  }),
  actions: {
    async fetchBlankspotAreas() {
      try {
        const response = await apiClient.get('/blankspots');
        this.blankspotAreas = response.data.data;
      } catch (error) {
        console.error('Error fetching blankspot areas:', error);
        throw error;
      }
    },
    async createBlankspotArea(blankspotData) {
      try {
        const response = await apiClient.post('/blankspots', blankspotData);
        this.blankspotAreas.push(response.data.data);
        return response.data;
      } catch (error) {
        console.error('Error creating blankspot area:', error);
        throw error;
      }
    },
    async updateBlankspotArea(id, blankspotData) {
      try {
        const response = await apiClient.put(`/blankspots/${id}`, blankspotData);
        const index = this.blankspotAreas.findIndex(area => area.ID === id);
        if (index !== -1) {
          this.blankspotAreas[index] = response.data.data;
        }
        return response.data;
      } catch (error) {
        console.error('Error updating blankspot area:', error);
        throw error;
      }
    },
    async deleteBlankspotArea(id) {
      try {
        await apiClient.delete(`/blankspots/${id}`);
        this.blankspotAreas = this.blankspotAreas.filter(area => area.ID !== id);
      } catch (error) {
        console.error('Error deleting blankspot area:', error);
        throw error;
      }
    },
  },
});
