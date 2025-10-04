import axios from 'axios';
import { useAuthStore } from '../stores/auth.js';

const apiClient = axios.create({
  baseURL: '/api'
});

apiClient.interceptors.request.use(config => {
  const authStore = useAuthStore();
  const token = authStore.token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default apiClient;
