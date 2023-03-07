<script setup lang="ts">
import { Tab, TabGroup, TabList, TabPanels, TabPanel } from "@headlessui/vue";
import { ref, watch } from "vue";
import { useAuth0 } from "@auth0/auth0-vue";
const { user, getAccessTokenSilently } = useAuth0();

// regular expression to validate permalink with
const permalinkRegex =
  /https:\/\/github\.com\/[a-zA-Z0-9-_\.]+\/[a-zA-Z0-9-_\.]+\/blob\/[a-z0-9]{40}(\/[a-zA-Z0-9-_\.]+)+#L[0-9]+-L[0-9]+/;

const postType = ref("permalink-range");
const permalink = ref("");
const permalinkPreview = ref("");
const content = ref("");

// function to get the permalink
const getPermalinkPreview = async () => {
  try {
    // parse the permalink
    const permalinkUrl = new URL(permalink.value);

    // get the range start/end from the hash
    const [rangeStart, rangeEnd] = permalinkUrl.hash
      .slice(1) // remove the '#'
      .split("-") // separate start/end
      .map((v) => parseInt(v.slice(1))); // remove the 'L' from start/end & parse as int
    permalinkUrl.hash = "";

    // change the host from github.com to raw.githubusercontent.com
    permalinkUrl.host = "raw.githubusercontent.com";

    // remove the /blob segment from the url
    permalinkUrl.pathname = permalinkUrl.pathname
      .split("/")
      .filter((part) => part != "blob")
      .join("/");

    const response = await fetch(permalinkUrl);
    const contents = await response.text();
    const contentRange = contents
      .split(/\r\n|\n|\r/)
      .slice(rangeStart - 1, rangeEnd)
      .join("\n");
    permalinkPreview.value = contentRange;
  } catch (e: any) {
    //TODO: better error handling
    console.error("Failed to get the code preview", e);
    permalinkPreview.value = "";
  }
};
watch(
  permalink,
  (value) => {
    if (permalinkRegex.test(value)) {
      getPermalinkPreview();
    }
  },
  {
    immediate: true,
  }
);

async function submit() {
  const accessToken = await getAccessTokenSilently();
  const post = {
    author_id: user.value.sub,
    content: content.value,
    type: postType.value,
    data: permalink.value,
  };
  const response = await fetch("/api/post", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
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
              <Tab class="px-6 py-4 text-gray-700 ui-selected:bg-seagreen-500"
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
                    <pre
                      name="preview"
                      class="rounded-md bg-gray-800 px-6 py-4 text-white"
                      >{{ permalinkPreview }}</pre
                    >
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
