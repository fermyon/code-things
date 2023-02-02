<script setup lang="ts">
import { ref } from "vue";
import { useAuth0 } from "@auth0/auth0-vue";
const { user, getAccessTokenSilently } = useAuth0();

const handle = ref(user.value.nickname);
const avatar = ref(user.value.picture);

async function submit() {
  try {
    const accessToken = await getAccessTokenSilently();
    const profile = {
      handle: handle.value,
      avatar: avatar.value,
    };
    const response = await fetch("/api/profile", {
      method: "POST",
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
      body: JSON.stringify(profile),
    });
    if (response.ok) {
      const updatedProfile = await response.json();
      handle.value = updatedProfile.handle;
      avatar.value = updatedProfile.avatar;
    } else {
      const message = await response.text();
      console.error(`Profile API Error: ${response.statusText} ${message}`);
    }
  } catch (ex) {
    console.log(ex);
  }
}
</script>

<template>
  <main class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
    <div>
      <div class="md:grid md:grid-cols-3 md:gap-6">
        <div class="md:col-span-1">
          <div class="px-4 sm:px-0">
            <h3 class="text-lg font-medium leading-6 text-gray-900">Profile</h3>
            <p class="mt-1 text-sm text-gray-600">
              This information will be displayed publicly so be careful what you
              share.
            </p>
          </div>
        </div>
        <div class="mt-5 md:col-span-2 md:mt-0">
          <form @submit.prevent="submit">
            <div class="shadow sm:overflow-hidden sm:rounded-md">
              <div class="space-y-6 bg-white px-4 py-5 sm:p-6">
                <div class="grid grid-cols-6 gap-6">
                  <div class="col-span-6 sm:col-span-3">
                    <label
                      for="handle"
                      class="block text-sm font-medium text-gray-700"
                      >Handle</label
                    >
                    <input
                      v-model="handle"
                      type="text"
                      name="handle"
                      id="handle"
                      autocomplete="handle"
                      class="focus:border-indigo-500 focus:ring-indigo-500 mt-1 block w-full rounded-md border-gray-300 shadow-sm sm:text-sm"
                    />
                  </div>

                  <div class="col-span-6 sm:col-span-4">
                    <label
                      for="avatar"
                      class="block text-sm font-medium text-gray-700"
                      >Avatar</label
                    >
                    <div class="mt-1 flex items-center content-center">
                      <img
                        class="h-10 w-10 mr-2 rounded-full"
                        :src="avatar" />
                      <input
                        v-model="avatar"
                        type="text"
                        name="avatar"
                        id="avatar"
                        autocomplete="email"
                        class="focus:border-indigo-500 focus:ring-indigo-500 block w-full rounded-md border-gray-300 shadow-sm sm:text-sm"
                      />
                    </div>
                  </div>
                </div>
              </div>
              <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
                <button
                  type="submit"
                  class="inline-flex justify-center rounded-md border border-transparent bg-seagreen-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-seagreen-700 focus:outline-none focus:ring-2 focus:ring-seagreen-500 focus:ring-offset-2"
                >
                  Save
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </main>
</template>
