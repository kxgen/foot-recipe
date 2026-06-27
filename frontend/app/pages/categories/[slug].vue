<!-- extra data -->
<script setup>
import { GET_CATEGORY_BY_SLUG } from "~/graphql/queries.js"
import { mapRecipeCard } from "~/utils/recipe-mappers.js"

function mapRecipes(list) {
  return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

const route = useRoute();
const categorySlug = computed(() => route.params.slug);
const { $apollo: apollo } = useNuxtApp();

const categoryName = ref("Category");
const recipes = ref([]);
const loading = ref(true);
const loadError = ref("");

useHead(() => ({
  title: `${categoryName.value} — Share Plate`,
}));

onMounted(loadCategory);
watch(categorySlug, loadCategory);

async function loadCategory() {
  if (!categorySlug.value) return;

  loading.value = true;
  loadError.value = "";

  try {
    const { data } = await apollo.query({
      query: GET_CATEGORY_BY_SLUG,
      variables: { slug: categorySlug.value },
      fetchPolicy: "network-only",
    });

    const category = data?.categories?.[0];
    categoryName.value = category?.name ?? "Category";
    recipes.value = mapRecipes(category?.recipes);
  } catch (err) {
    loadError.value = err?.message || "Failed to load category.";
    recipes.value = [];
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="py-12">
    <div class="page-container">
      <NuxtLink to="/categories" class="text-sm text-primary hover:underline">
        ← All categories
      </NuxtLink>

      <h1 class="mt-4 text-3xl font-semibold text-heading">{{ categoryName }} recipes</h1>
      <p class="mt-2 text-sm text-heading-muted">
        Recipes in the {{ categoryName }} category.
      </p>

      <p v-if="loadError" class="mt-4 error-message">
        {{ loadError }}
      </p>

      <div v-if="loading" class="mt-8 text-sm text-heading-muted"><InlineLoader/></div>

      <div v-else-if="recipes.length" class="mt-8 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
        <RecipeCard
          v-for="recipe in recipes"
          :key="recipe.id"
          :id="recipe.id"
          :title="recipe.title"
          :creator="recipe.creator"
          :rating-score="recipe.avgRating"
          :rating-count="recipe.ratingCount"
          :time-minutes="recipe.timeMinutes"
          :image-src="recipe.imageSrc"
        />
      </div>

      <div v-else class="mt-8 card-surface p-6 text-center text-sm text-heading-muted">
        No recipes in this category yet.
      </div>
    </div>
  </div>
</template>
