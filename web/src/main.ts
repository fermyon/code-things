import { createApp } from "vue";
import { createPinia } from "pinia";
import { createAuth0 } from "@auth0/auth0-vue";
import piniaPluginPersistedState from 'pinia-plugin-persistedstate';

import "@/assets/main.css";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);

// setup application state
const pinia = createPinia();
pinia.use(piniaPluginPersistedState)
app.use(pinia);

// router must come after store
app.use(router);

// debug configuration
const domain = import.meta.env.VITE_AUTH0_DOMAIN;
const clientId = import.meta.env.VITE_AUTH0_CLIENT_ID;
const audience = import.meta.env.VITE_AUTH0_AUDIENCE;

const auth0 = createAuth0({
    domain,
    clientId,
    authorizationParams: {
        audience,
        redirect_uri: window.location.origin,
    },
    cacheLocation: "localstorage",
});
app.use(auth0);

app.mount("#app");
