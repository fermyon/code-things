<script setup lang="ts">
import { UserCircleIcon } from "@heroicons/vue/24/outline";
import { useAuth0 } from "@auth0/auth0-vue";
const {
    user,
} = useAuth0();
</script>

<template>
<img
    v-if="user && user.picture"
    class="h-8 w-8 rounded-full"
    :src="user.picture"
    alt=""
    />
<UserCircleIcon
    v-else
    class="h-8 w-8 rounded-full"
    />
</template>