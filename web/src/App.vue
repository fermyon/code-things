<script setup lang="ts">
import { computed, ref } from "vue";
import { RouterView } from "vue-router";

import { NavBar, NavLink } from "@/components/nav";
import DropdownMenu from "@/components/DropdownMenu.vue";
import UserProfileImage from "@/components/Avatar.vue";

const isAuthenticated = ref(false);

const userNavItems = computed(() =>
  isAuthenticated.value
    ? [{ text: "Profile", route: "/profile" }, { text: "Log Out", click: logoutHandler }]
    : [{ text: "Log In", click: loginHandler }]
);


function loginHandler() {
  isAuthenticated.value = true;
}

function logoutHandler() {
  isAuthenticated.value = false;
}

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
