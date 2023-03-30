<script setup lang="ts">
import { Disclosure, DisclosureButton, DisclosurePanel } from "@headlessui/vue";
import { getPermalinkPreview } from "@/utils";
import { ref, watch, computed } from "vue";

const props = defineProps({
  type: {
    type: String,
    required: true,
    validator(value: string) {
      return ["permalink-range"].includes(value);
    },
  },
  data: {
    type: String,
  },
});

const previewContent = ref<string | null>();
const previewContentHead = ref<string | null>();
const previewContentTail = ref<string | null>();

const computeHeadTail = (value: string | null | undefined) => {
  const splitAt = 5;
  const lines = value?.split("\n");
  previewContentHead.value = lines?.slice(0, splitAt).join("\n");
  previewContentTail.value = lines?.slice(5).join("\n");
};

watch(
  () => props.data,
  async (value) => {
    if (value && props.type == "permalink-range") {
      previewContent.value = await getPermalinkPreview(value);
    }
  },
  {
    immediate: true,
  }
);
watch(previewContent, computeHeadTail);
</script>

<template>
  <Disclosure as="pre"
        class="rounded-md bg-gray-800 px-6 py-4 text-white cursor-pointer">
    <DisclosureButton as="pre">{{ previewContentHead }}</DisclosureButton>
    <transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-out"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <DisclosurePanel
        as="pre"
        >{{ previewContentTail }}</DisclosurePanel
      >
    </transition>
  </Disclosure>
</template>
