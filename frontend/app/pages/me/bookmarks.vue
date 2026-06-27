<template>
  <div class="page-container py-14">
    <h1 class="text-2xl font-semibold text-heading">Bookmarks</h1>
    <p class="mt-2 text-sm text-heading-muted">
      Recipes you saved to cook later.
    </p>

    <p v-if="!isLoggedIn" class="mt-6 text-sm text-heading-muted">
      <NuxtLink to="/auth/login" class="font-semibold text-primary hover:underline">
        Log in
      </NuxtLink>
      to view bookmarks.
    </p>

    <p v-else-if="loadError" class="mt-6 error-message">
      {{ loadError }}
    </p>

    <div v-else-if="loading" class="mt-6 text-sm text-heading-muted"><InlineLoader/></div>

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
        :recipe-link="recipe.recipeLink"
      />
    </div>

    <div v-else class="mt-8 card-surface p-6 text-center text-sm text-heading-muted">
      No bookmarked recipes yet. Save recipes from any recipe page.
    </div>
  </div>
</template>

<script setup>
import { getUserIdFromToken } from "~/utils/jwt";
import { GET_MY_BOOKMARKS } from "~/graphql/queries.js"
import { mapRecipeCard } from "~/utils/recipe-mappers.js"

useHead({ title: "Bookmarks — Share Plate" });

const { token, isLoggedIn } = useAuth();
const { $apollo: apollo } = useNuxtApp();

const recipes = ref([]);
const loading = ref(false);
const loadError = ref("");

onMounted(loadBookmarks);
watch(() => token.value, loadBookmarks);

async function loadBookmarks() {
  const userId = getUserIdFromToken(token.value);
  if (!userId) {
    recipes.value = [];
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const { data } = await apollo.query({
      query: GET_MY_BOOKMARKS,
      variables: { userId },
      fetchPolicy: "network-only",
    });
    recipes.value = (data?.recipe_bookmarks ?? [])
      .map((row) => mapRecipeCard(row.recipe))
      .filter(Boolean);
  } catch (err) {
    loadError.value = err?.message || "Failed to load bookmarks.";
  } finally {
    loading.value = false;
  }
}
</script>

