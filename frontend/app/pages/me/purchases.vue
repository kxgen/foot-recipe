<template>
  <div class="page-container py-14">
    <h1 class="text-2xl font-semibold text-heading">Purchases</h1>
    <p class="mt-2 text-sm text-heading-muted">
      Recipes you have purchased.
    </p>

    <p v-if="!isLoggedIn" class="mt-6 text-sm text-heading-muted">
      <NuxtLink to="/auth/login" class="font-semibold text-primary hover:underline">
        Log in
      </NuxtLink>
      to view purchases.
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
      />
    </div>

    <div v-else class="mt-8 card-surface p-6 text-center text-sm text-heading-muted">
      No purchased recipes yet.
    </div>
  </div>
</template>

<script setup>
import { gql } from "@apollo/client/core";
import { getUserIdFromToken } from "~/utils/jwt";
import { GET_MY_PURCHASES } from "~/graphql/queries.js"
import { mapRecipeCard } from "~/utils/recipe-mappers.js"

const { $apollo: apollo } = useNuxtApp();
const { token, isLoggedIn } = useAuth();
const query = (options) => apollo.query({ fetchPolicy: "network-only", ...options });

function mapRecipes(list) {
  return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

useHead({ title: "Purchases — Share Plate" });

const recipes = ref([]);
const loading = ref(false);
const loadError = ref("");

onMounted(loadPurchases);
watch(() => token.value, loadPurchases);

async function loadPurchases() {
  const userId = getUserIdFromToken(token.value);
  if (!userId) {
    recipes.value = [];
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const { data } = await query({
      query: GET_MY_PURCHASES,
      variables: { userId },
    });
    recipes.value = (data?.recipe_purchases ?? [])
      .map((row) => mapRecipes([row.recipe])[0])
      .filter(Boolean);
  } catch (err) {
    loadError.value = err?.message || "Failed to load purchases.";
  } finally {
    loading.value = false;
  }
}
</script>

