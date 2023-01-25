import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedState from 'pinia-plugin-persistedstate';

import App from "./App.vue";
import router from "./router";

import "./assets/main.css";

const app = createApp(App);

app.use(createPinia()
    .use(piniaPluginPersistedState));
app.use(router);

app.mount("#app");
