<script setup lang="ts">
import Avatar from "@/components/Avatar.vue";
</script>

<template>
  <main>
    <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="rounded-lg bg-white py-4 px-4">
          <div class="flex items-center gap-2">
            <Avatar class="h-8 w-8" />
            <RouterLink
              to="/posts/new"
              class="w-full rounded-md py-1.5 px-1.5 text-gray-400 ring-1 ring-inset ring-gray-300"
              >Create Post</RouterLink
            >
          </div>
        </div>
      </div>
    </div>
    <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="h-96 rounded-lg border-4 border-dashed border-gray-200" />
      </div>
    </div>
  </main>
</template>
