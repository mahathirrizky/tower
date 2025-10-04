import { defineStore } from 'pinia';
import apiClient from '../api/axios'; // Import apiClient

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    setToken(token) {
      this.token = token;
      localStorage.setItem('token', token);
    },
    async login(credentials) {
      try {
        const response = await apiClient.post('/auth/login', credentials);
        
        this.setToken(response.data.data.token);
        return response.data; // Return data if needed
      } catch (error) {
        this.token = null; // Clear token on login failure
        localStorage.removeItem('token');
        throw error; // Re-throw to allow component to handle
      }
    },
    logout() {
      this.token = null;
      localStorage.removeItem('token');
    },
    async changePassword(currentPassword, newPassword) {
      try {
        const response = await apiClient.put('/auth/change-password', {
          current_password: currentPassword,
          new_password: newPassword
        });
        return response.data;
      } catch (error) {
        throw error.response?.data?.message || 'Failed to change password';
      }
    }
  }
});
