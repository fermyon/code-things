<script setup lang="ts">
import { computed } from "vue";
import { RouterView } from "vue-router";
import { useAuth0 } from "@auth0/auth0-vue";

import { NavBar, NavLink } from "@/components/nav";
import DropdownMenu from "@/components/DropdownMenu.vue";
import UserProfileImage from "@/components/UserProfileImage.vue";

const {
  isAuthenticated,
  loginWithRedirect,
  logout,
} = useAuth0();

function loginHandler() {
  loginWithRedirect({
    appState: {
      target: "/",
    },
  });
}

function logoutHandler() {
  logout();
}

const userNavItems = computed(() =>
  isAuthenticated.value
    ? [{ text: "Profile", route: "/profile" }, { text: "Log Out", click: logoutHandler }]
    : [{ text: "Log In", click: loginHandler }]
);
</script>

<template>
  <NavBar>
    <template #brand>
      <NavLink href="/">
        <img
          class="h-8 w-auto lg:block"
          src="@/assets/logo.svg"
          alt="Code Things"
        />
      </NavLink>
    </template>
    <template #default></template>
    <template #nav-menu>
      <DropdownMenu :items="userNavItems">
        <template #button>
          <UserProfileImage />
        </template>
      </DropdownMenu>
    </template>
  </NavBar>

  <RouterView />
</template>
