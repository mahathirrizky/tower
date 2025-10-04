import { defineStore } from 'pinia';

export const useUiStore = defineStore('ui', {
  state: () => ({
    isDarkMode: localStorage.getItem('isDarkMode') === 'true',
    isDetailSidebarVisible: false,
    detailItem: null,
  }),
  actions: {
    toggleDarkMode() {
      this.isDarkMode = !this.isDarkMode;
      localStorage.setItem('isDarkMode', this.isDarkMode);
      if (this.isDarkMode) {
        document.documentElement.classList.add('my-app-dark');
      } else {
        document.documentElement.classList.remove('my-app-dark');
      }
    },
    initializeDarkMode() {
      if (this.isDarkMode) {
        document.documentElement.classList.add('my-app-dark');
      } else {
        document.documentElement.classList.remove('my-app-dark');
      }
    },
    openDetailSidebar(item) {
      this.isDetailSidebarVisible = true;
      this.detailItem = item;
    },
    closeDetailSidebar() {
      this.isDetailSidebarVisible = false;
      this.detailItem = null;
    },
  }
});
