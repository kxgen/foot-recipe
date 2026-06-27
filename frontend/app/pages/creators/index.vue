<!--  -->
<!-- all creators page -->
<template>
  <div class="py-12">
    <div class="page-container">
      <h1 class="text-3xl font-semibold text-heading">Creators</h1>
      <p class="mt-2 text-sm text-heading-muted">Discover recipes by creator.</p>

      <p v-if="loadError" class="mt-4 error-message">
        {{ loadError }}
      </p>

      <div v-if="loading" class="mt-8 text-sm text-heading-muted"><InlineLoader/></div>

      <div v-else class="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <CreatorCard
          v-for="creator in creators"
          :key="creator.id"
          :username="creator.name"
          :slug="creator.slug"
          :avatar-url="creator.avatarUrl"
          :recipe-count="creator.recipes"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { GET_CREATORS } from "~/graphql/queries.js"

useHead({
  title: "Creators — Share Plate",
});

const { $apollo: apollo } = useNuxtApp();

const creators = ref([]);
const loading = ref(true);
const loadError = ref("");

onMounted(async () => {
  try {
    const { data } = await apollo.query({
      query: GET_CREATORS,
      fetchPolicy: "network-only",
    });
    creators.value = (data?.users ?? [])
      .filter((u) => (u.recipes_aggregate?.aggregate?.count ?? 0) > 0)
      .map((u) => ({
        id: u.id,
        name: u.username,
        slug: u.slug,
        avatarUrl: u.avatar_url,
        recipes: u.recipes_aggregate?.aggregate?.count ?? 0,
      }));
  } catch (err) {
    loadError.value = err?.message || "Failed to load creators.";
  } finally {
    loading.value = false;
  }
});
</script>