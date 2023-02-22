<script setup lang="ts">
import { UserCircleIcon } from "@heroicons/vue/24/outline";
import { ref } from "vue";
const avatar = ref("");
</script>

<template>
<img
    v-if="avatar"
    class="h-8 w-8 rounded-full"
    :src="avatar"
    alt=""
    />
<UserCircleIcon
    v-else
    class="h-8 w-8 rounded-full"
    />
</template>