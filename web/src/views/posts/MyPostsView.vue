<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useAuth0 } from "@auth0/auth0-vue";
import CodePreview from "@/components/CodePreview.vue";

const { getAccessTokenSilently } = useAuth0();

const posts = ref<
  [
    {
      id: number;
      author_id: string;
      content: string;
      type: string;
      data: string;
      visibility: string;
    }
  ]
>();

onMounted(async () => {
  // posts.value = [
  //   {
  //     id: 10,
  //     author_id: "github|3060890",
  //     content: "Test Content",
  //     type: "permalink-range",
  //     data: "https://github.com/fermyon/spin/blob/9095fe60a54775672141bbb12f3d3716b333ce9e/src/commands/build.rs#L31-L53",
  //     visibility: "public",
  //   },
  //   {
  //     id: 11,
  //     author_id: "github|3060890",
  //     content: "Test Content",
  //     type: "permalink-range",
  //     data: "https://github.com/fermyon/spin/blob/9095fe60a54775672141bbb12f3d3716b333ce9e/src/commands/build.rs#L31-L53",
  //     visibility: "public",
  //   },
  //   {
  //     id: 12,
  //     author_id: "github|3060890",
  //     content: "Test Content",
  //     type: "permalink-range",
  //     data: "https://github.com/fermyon/spin/blob/9095fe60a54775672141bbb12f3d3716b333ce9e/src/commands/build.rs#L31-L53",
  //     visibility: "public",
  //   },
  //   {
  //     id: 13,
  //     author_id: "github|3060890",
  //     content: "Test Content",
  //     type: "permalink-range",
  //     data: "https://github.com/fermyon/spin/blob/9095fe60a54775672141bbb12f3d3716b333ce9e/src/commands/build.rs#L31-L53",
  //     visibility: "public",
  //   },
  //   {
  //     id: 14,
  //     author_id: "github|3060890",
  //     content: "Test Content",
  //     type: "permalink-range",
  //     data: "https://github.com/fermyon/spin/blob/9095fe60a54775672141bbb12f3d3716b333ce9e/src/commands/build.rs#L31-L53",
  //     visibility: "public",
  //   },
  // ];
  // return;
  const limit = 5;
  const offset = 0;
  const accessToken = await getAccessTokenSilently();
  const response = await fetch(`/api/post?limit=${limit}&offset=${offset}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Content-Type": "application/json",
    },
  });
  posts.value = await response.json();
});
</script>

<template>
  <header class="bg-white shadow">
    <div class="mx-auto max-w-7xl py-6 px-4 sm:px-6 lg:px-8">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">My Posts</h1>
    </div>
  </header>
  <main>
    <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div class="px-6 py-4 sm:px-0" v-for="post in posts">
        <div class="rounded-lg bg-white px-6 py-4">
          <div class="pb-4">
            {{ post.content }}
          </div>
          <CodePreview :type="post.type" :data="post.data" />
          <div class="pt-4">
            <a :href="post.data" target="_blank" class="hover:text-oxfordblue-600"
              >View source</a
            >
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
