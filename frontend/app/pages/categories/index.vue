<!-- xo -->
<template>
  <div class="py-12">
    <div class="page-container">
      <h1 class="text-3xl font-semibold text-heading">Categories</h1>
      <p class="mt-2 text-sm text-heading-muted">Browse recipes by food type.</p>

      <p v-if="loadError" class="mt-4 error-message">
        {{ loadError }}
      </p>

      <div v-if="loading" class="mt-8 text-sm text-heading-muted"><InlineLoader/></div>

      <div v-else class="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <CategoryCard
          v-for="category in categories"
          :key="category.id"
          :id="category.id"
          :name="category.name"
          :slug="category.slug"
          :count="category.count"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { GET_CATEGORIES } from "~/graphql/queries.js"

useHead({
  title: "Categories — Share Plate",
});

const { $apollo: apollo } = useNuxtApp();

const categories = ref([]);
const loading = ref(true);
const loadError = ref("");

onMounted(async () => {
  try {
    const { data } = await apollo.query({
      query: GET_CATEGORIES,
      fetchPolicy: "network-only",
    });
    categories.value = (data?.categories ?? []).map((c) => ({
      id: c.id,
      name: c.name,
      slug: c.slug ,
      count: c.recipes_aggregate?.aggregate?.count ?? 0,
    }));
  } catch (err) {
    loadError.value = err?.message || "Failed to load categories.";
  } finally {
    loading.value = false;
  }
});
</script>
