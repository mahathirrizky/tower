import { createApp } from "vue";
import { createPinia } from 'pinia';
import App from "./App.vue";
import PrimeVue from "primevue/config";
import Aura from '@primeuix/themes/aura';
import router from "./router/index.js";
import 'primeicons/primeicons.css';
import './style.css';

import { useUiStore } from './stores/ui'; // Import the UI store

import ToastService from 'primevue/toastservice';
import Tooltip from 'primevue/tooltip';
const app = createApp(App);
const pinia = createPinia();


app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: '.my-app-dark',
            cssLayer:{
              name:'primevue',
              order:'theme, base,primevue'
            }
        }
    }
 });

app.directive('tooltip', Tooltip);
app.use(pinia);
app.use(router);
app.use(ToastService);

// Initialize dark mode after Pinia is installed
const uiStore = useUiStore();
uiStore.initializeDarkMode();

app.mount("#app");