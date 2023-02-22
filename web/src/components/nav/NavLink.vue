<template>
  <RouterLink
    v-slot="{ href, isActive, isExactActive, navigate }"
    :to="href"
    custom
  >
    <a
      v-bind="$attrs"
      :href="href"
      :class="isActive ? activeClass : inactiveClass"
      :aria-current="isExactActive ? ariaCurrentValue : undefined"
      @click="navigate"
      >
      <slot />
    </a>
  </RouterLink>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  href: string;
  activeClass?: string;
  inactiveClass?: string;
  ariaCurrentValue?: 'page' | 'step' | 'location' | 'date' | 'time' | 'true' | 'false';
}>(), {
  ariaCurrentValue: 'page',
});
</script>
