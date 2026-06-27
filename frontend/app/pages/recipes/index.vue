<template>
    <div class="py-12">
        <div class="page-container">
            <div class="mb-8 flex items-end justify-between gap-4">
                <div>
                    <h1 class="text-3xl font-semibold text-heading">
                        All Recipes
                    </h1>
                    <p class="mt-2 text-sm text-heading-muted">
                        Explore free and premium recipes from the community.
                    </p>
                </div>
                <NuxtLink to="/recipes/create" class="btn-secondary"
                    >Create recipe</NuxtLink
                >
            </div>

            <p v-if="loadError" class="mb-4 error-message">
                {{ loadError }}
            </p>

            <div v-if="loading" class="text-sm text-heading-muted">
                <InlineLoader />
            </div>

            <div v-else class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
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
        </div>
    </div>
</template>

<script setup>
import { gql } from "@apollo/client/core";
import { GET_RECIPES } from "~/graphql/queries.js";
import { mapRecipeCard } from "~/utils/recipe-mappers.js";

const { $apollo: apollo } = useNuxtApp();
const query = (options) =>
    apollo.query({ fetchPolicy: "network-only", ...options });

function mapRecipes(list) {
    return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

useHead({
    title: "Recipes — Share Plate",
});

const recipes = ref([]);
const loading = ref(true);
const loadError = ref("");

onMounted(async () => {
    try {
        const { data } = await query({
            query: GET_RECIPES,
            variables: {
                orderBy: [{ created_at: "desc" }],
                limit: 60,
            },
        });
        recipes.value = mapRecipes(data?.recipes);
    } catch (err) {
        loadError.value = err?.message || "Failed to load recipes.";
    } finally {
        loading.value = false;
    }
});
</script>
