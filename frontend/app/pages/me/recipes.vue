<template>
  <div class="page-container py-14">
    <div class="mb-8 flex flex-wrap items-end justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-heading">My recipes</h1>
        <p class="mt-2 text-sm text-heading-muted">
          Recipes you have created.
        </p>
      </div>
      <NuxtLink to="/recipes/create" class="btn-primary">Create recipe</NuxtLink>
    </div>

    <p v-if="!isLoggedIn" class="text-sm text-heading-muted">
      <NuxtLink to="/auth/login" class="font-semibold text-primary hover:underline">
        Log in
      </NuxtLink>
      to see your recipes.
    </p>

    <p v-else-if="loadError" class="error-message">
      {{ loadError }}
    </p>

    <div v-else-if="loading" class="text-sm text-heading-muted"><InlineLoader/></div>

    <div v-else-if="recipes.length" class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
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

    <div v-else class="card-surface p-6 text-center text-sm text-heading-muted">
      You have not created any recipes yet.
    </div>
  </div>
</template>

<script setup>
import { gql } from "@apollo/client/core";
import { getUserIdFromToken } from "~/utils/jwt";
import { GET_MY_RECIPES } from "~/graphql/queries.js"
import { mapRecipeCard } from "~/utils/recipe-mappers.js"

const { $apollo: apollo } = useNuxtApp();
const { token, isLoggedIn } = useAuth();
const query = (options) => apollo.query({ fetchPolicy: "network-only", ...options });

function mapRecipes(list) {
  return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

useHead({ title: "My recipes — Share Plate" });

const recipes = ref([]);
const loading = ref(false);
const loadError = ref("");

onMounted(loadMyRecipes);
watch(() => token.value, loadMyRecipes);

async function loadMyRecipes() {
  const userId = getUserIdFromToken(token.value);
  if (!userId) {
    recipes.value = [];
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const { data } = await query({
      query: GET_MY_RECIPES,
      variables: { userId },
    });
    recipes.value = mapRecipes(data?.recipes);
  } catch (err) {
    loadError.value = err?.message || "Failed to load your recipes.";
  } finally {
    loading.value = false;
  }
}

</script>





