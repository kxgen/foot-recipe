<!-- xo -->
<!-- specific creator page-->
<script setup>
import { GET_CREATOR_BY_SLUG } from "~/graphql/queries.js"
import { mapRecipeCard } from "~/utils/recipe-mappers.js"

function mapRecipes(list) {
  return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

const route = useRoute();
// holds the route's creators/[slug], slug portion as variable
const creatorSlug = computed(() => route.params.slug);
const { $apollo: apollo } = useNuxtApp();

const creatorName = ref("Creator");
const creatorAvatar = ref("");
const creatorBio = ref("");
const creatorRecipes = ref([]);
const loading = ref(true);
const loadError = ref("");

useHead(() => ({
  title: `${creatorName.value} — Share Plate`,
}));

onMounted(loadCreator);
// when i get creatorslug, run loadCreator()
watch(creatorSlug, loadCreator);

async function loadCreator() {
  if (!creatorSlug.value) return;

  loading.value = true;
  loadError.value = "";

  try {
    const { data } = await apollo.query({
      query: GET_CREATOR_BY_SLUG,
      variables: { slug: creatorSlug.value },
      fetchPolicy: "network-only",
    });

    const creator = data?.users?.[0];
    creatorName.value = creator?.username ?? "Unknown creator";
    creatorBio.value = creator?.bio ?? "";
    creatorAvatar.value = creator?.avatar_url ?? "";
    creatorRecipes.value = mapRecipes(creator?.recipes);
  } catch (err) {
    loadError.value = err?.message || "Failed to load creator.";
    creatorRecipes.value = [];
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="py-12">
    <div class="page-container">
      <NuxtLink to="/creators" class="text-sm text-primary hover:underline">
        ← All creators
      </NuxtLink>

      <div class="mt-6 flex items-center gap-5">
        <div v-if="creatorAvatar" class="size-25 overflow-hidden rounded-full border border-primary-subtle/20 bg-white shadow-sm shrink-0">
          <img :src="creatorAvatar" :alt="creatorName" class="h-full w-full object-cover" />
        </div>
        <div v-else class="flex size-25 items-center justify-center rounded-full border border-primary/20 bg-white text-2xl font-bold text-heading shrink-0">
          {{ creatorName.charAt(0).toUpperCase() }}
        </div>
        <h1 class="text-3xl font-semibold text-heading leading-tight">{{ creatorName }}</h1>
      </div>
      <p v-if="creatorBio" class="mt-3 text-sm leading-relaxed text-heading-muted max-w-2xl">
        {{ creatorBio }}
      </p>
      <p v-else class="mt-2 text-sm text-heading-muted">Recipes shared by this creator.</p>

      <p v-if="loadError" class="mt-4 error-message">
        {{ loadError }}
      </p>

      <div v-if="loading" class="mt-8 text-sm text-heading-muted"><InlineLoader/></div>

      <div v-else-if="creatorRecipes.length" class="mt-8 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
        <RecipeCard
          v-for="recipe in creatorRecipes"
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
        This creator has not shared any recipes yet.
      </div>
    </div>
  </div>
</template>
