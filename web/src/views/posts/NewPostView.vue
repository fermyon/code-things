<script setup lang="ts">
import { Tab, TabGroup, TabList, TabPanels, TabPanel } from "@headlessui/vue";
import CodePreview from "@/components/CodePreview.vue";
import { ref } from "vue";
import { useAuth0 } from "@auth0/auth0-vue";

const { user, getAccessTokenSilently } = useAuth0();

const postType = ref("permalink-range");
const permalink = ref("");
const content = ref("");
const visibility = ref("public");

async function submit() {
  const accessToken = await getAccessTokenSilently();
  const post = {
    author_id: user.value.sub,
    content: content.value,
    type: postType.value,
    data: permalink.value,
    visibility: visibility.value,
  };
  const response = await fetch("/api/post", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(post),
  });
  if (response.ok) {
    //TODO: navigate back?
  } else {
    const message = await response.text();
    console.error(`Create Post API Error: ${response.statusText} ${message}`);
  }
}
function postTypeChanged(index: number) {
  postType.value = ["permalink-range", "code"][index];
}
</script>

<template>
  <header class="bg-white shadow">
    <div class="mx-auto max-w-7xl py-6 px-4 sm:px-6 lg:px-8">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">
        Create a Post
      </h1>
    </div>
  </header>
  <main>
    <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div class="px-6 py-4 sm:px-0">
        <div class="rounded-lg bg-white">
          <TabGroup @change="postTypeChanged">
            <TabList class="border-b-[1px]">
              <Tab class="ui-selected:bg-seagreen-500 px-6 py-4 text-gray-700"
                >Link</Tab
              >
              <Tab class="cursor-not-allowed px-6 py-4 text-gray-200" disabled
                >Code</Tab
              >
            </TabList>
            <TabPanels class="px-6 py-4">
              <TabPanel>
                <form @submit.prevent="submit">
                  <div class="pb-4">
                    <label
                      for="link"
                      class="block text-sm font-medium text-gray-700"
                      >Permalink (must be a public GitHub permalink with hash
                      for line range)</label
                    >
                    <input
                      v-model="permalink"
                      type="text"
                      name="link"
                      id="link"
                      autocomplete="url"
                      class="focus:border-indigo-500 focus:ring-indigo-500 mt-1 block w-full rounded-md border-gray-300 shadow-sm sm:text-sm"
                    />
                  </div>
                  <div class="pb-4">
                    <label
                      for="preview"
                      class="block text-sm font-medium text-gray-700"
                      >Preview</label
                    >
                    <CodePreview :type="postType" :data="permalink" />
                  </div>
                  <div class="pb-4">
                    <label
                      for="content"
                      class="block text-sm font-medium text-gray-700"
                      >Content</label
                    >
                    <textarea
                      v-model="content"
                      name="content"
                      id="content"
                      class="focus:border-indigo-500 focus:ring-indigo-500 mt-1 block w-full rounded-md border-gray-300 shadow-sm sm:text-sm"
                    />
                  </div>
                  <div class="text-right">
                    <button
                      type="submit"
                      class="inline-flex justify-center rounded-md border border-transparent bg-seagreen-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-seagreen-700 focus:outline-none focus:ring-2 focus:ring-seagreen-500 focus:ring-offset-2"
                    >
                      Post
                    </button>
                  </div>
                </form>
              </TabPanel>
              <TabPanel> Code Panel </TabPanel>
            </TabPanels>
          </TabGroup>
        </div>
      </div>
    </div>
  </main>
</template>
