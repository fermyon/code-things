<template>
  <nav class="bg-gray-800">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-between">
        <!-- sticks to left side -->
        <div class="flex items-center">
          <!-- brand -->
          <div class="flex-shrink-0">
            <slot name="brand"></slot>
          </div>
          <div class="ml-10 text-white">
            <slot name="default"></slot>
          </div>
        </div>
        <div class="text-white">
          <slot name="nav-menu"></slot>
        </div>
      </div>
    </div>
  </nav>
</template>
