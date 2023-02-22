<script setup lang="ts">
import { useRouter } from "vue-router";
import { Menu, MenuItem, MenuItems, MenuButton } from "@headlessui/vue";

export interface DropdownMenuItem {
  text: string;
  route?: string;
  click?: (payload: MouseEvent) => void;
}

defineProps<{
  items: DropdownMenuItem[];
}>();

const router = useRouter();
function menuClick(payload: MouseEvent, item: DropdownMenuItem) {
  if (item.click) {
    item.click(payload);
    return;
  }
  if (item.route) {
    router.push(item.route);
    return;
  }
  throw new Error(
    "Developer Error: DropdownMenuItem must have one of click or route defined."
  );
}
</script>

<template>
  <Menu as="div" class="relative ml-3">
    <div>
      <MenuButton
        class="flex max-w-xs items-center rounded-full bg-gray-800 text-sm hover:outline-none hover:ring-2 hover:ring-white hover:ring-offset-2 hover:ring-offset-gray-800"
      >
        <span class="sr-only">Open user menu</span>
        <slot name="button"></slot>
      </MenuButton>
    </div>
    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <MenuItems
        class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
      >
        <MenuItem v-for="item in items" :key="item.text" v-slot="{ active }">
          <a
            @click="(e) => menuClick(e, item)"
            :class="[
              active ? 'bg-gray-100' : '',
              'block px-4 py-2 text-sm text-gray-700',
            ]"
          >
            {{ item.text }}
          </a>
        </MenuItem>
      </MenuItems>
    </transition>
  </Menu>
</template>
