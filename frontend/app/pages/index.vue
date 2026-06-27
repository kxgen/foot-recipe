<!-- xo -->
<template>
    <div>
        <section class="relative overflow-hidden pt-12 sm:pt-16 pb-6">
            <div class="pointer-events-none absolute inset-0 -z-10" />

            <div class="page-container">
                <div class="mx-auto max-w-2xl text-center">
                    <p
                        class="text-sm font-semibold uppercase tracking-wider text-primary"
                    >
                        <span class="font-brand text-primary text-3xl"
                            >Share Plate</span
                        >
                    </p>
                    <h1 class="mt-2 text-3xl font-semibold sm:text-5xl">
                        Find recipes you&apos;ll love to cook
                    </h1>
                    <p class="mt-3 text-sm text-subtle sm:text-base">
                        Browse categories, discover creators, and save the ones
                        you want to make next.
                    </p>

                    <div class="mt-7">
                        <RecipeSearchBar />
                    </div>
                </div>
            </div>
        </section>

        <section class="py-14">
            <div class="page-container">
                <div class="mb-8 flex items-end justify-between gap-4">
                    <h2 class="text-2xl font-semibold">Latest Recipes</h2>
                    <NuxtLink
                        to="/recipes"
                        class="text-sm font-semibold text-subtle hover:underline"
                    >
                        View more →
                    </NuxtLink>
                </div>

                <div v-if="recipesLoading" class="text-sm text-subtle">
                    <InlineLoader />
                </div>

                <div v-else class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
                    <RecipeCard
                        v-for="recipe in latestRecipes"
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
        </section>

        <section class="py-14">
            <div class="page-container">
                <div class="mb-8 flex items-end justify-between gap-4">
                    <h2 class="text-2xl font-semibold">Categories</h2>
                    <NuxtLink
                        to="/categories"
                        class="text-sm font-semibold text-subtle hover:underline"
                    >
                        View more →
                    </NuxtLink>
                </div>

                <div v-if="categoriesLoading" class="text-sm text-subtle">
                    <InlineLoader />
                </div>

                <div v-else class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
                    <CategoryCard
                        v-for="c in categories"
                        :key="c.id"
                        :id="c.id"
                        :name="c.name"
                        :slug="c.slug"
                        :count="c.count"
                    />
                </div>
            </div>
        </section>

        <section class="py-14">
            <div class="page-container">
                <div class="mb-8 flex items-end justify-between gap-4">
                    <h2 class="text-2xl font-semibold">Creators</h2>
                    <NuxtLink
                        to="/creators"
                        class="text-sm font-semibold text-subtle hover:underline"
                    >
                        View more →
                    </NuxtLink>
                </div>

                <div v-if="creatorsLoading" class="text-sm text-subtle">
                    <InlineLoader />
                </div>

                <div v-else class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
                    <CreatorCard
                        v-for="cr in creators"
                        :key="cr.id"
                        :username="cr.name"
                        :slug="cr.slug"
                        :avatar-url="cr.avatarUrl"
                        :recipe-count="cr.recipes"
                    />
                </div>
            </div>
        </section>
    </div>
</template>

<script setup>
import { GET_LATEST_RECIPES, GET_CATEGORIES } from "~/graphql/queries.js";
import { GET_CREATORS } from "~/graphql/queries.js";
import { mapRecipeCard } from "~/utils/recipe-mappers.js";

const { $apollo: apollo } = useNuxtApp();
const query = (options) =>
    apollo.query({ fetchPolicy: "network-only", ...options });

function mapRecipes(list) {
    return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

useHead({
    title: "Share Plate — Food recipes",
});

const latestRecipes = ref([]);
const categories = ref([]);
const creators = ref([]);
const loadError = ref("");
const recipesLoading = ref(true);
const categoriesLoading = ref(true);
const creatorsLoading = ref(true);

onMounted(async () => {
    loadError.value = "";

    try {
        const [latestRes, categoriesRes, creatorsRes] = await Promise.all([
            query({ query: GET_LATEST_RECIPES, variables: { limit: 6 } }),
            query({ query: GET_CATEGORIES, variables: { limit: 10 } }),
            query({ query: GET_CREATORS }),
        ]);

        latestRecipes.value = mapRecipes(latestRes.data?.recipes);
        // [c:{}, c:{}]
        categories.value = (categoriesRes.data?.categories ?? []).map((c) => ({
            id: c.id,
            name: c.name,
            // Deriving a slug icon from lucid
            slug: c.slug || c.name?.toLowerCase().replace(/\s+/g, "-"),
            count: c.recipes_aggregate?.aggregate?.count ?? 0,
        }));

        // map data to creator card
        creators.value = (creatorsRes.data?.users ?? [])
            .filter((u) => (u.recipes_aggregate?.aggregate?.count ?? 0) > 0)
            .slice(0, 4)
            .map((u) => ({
                id: u.id,
                name: u.username,
                slug: u.slug,
                avatarUrl: u.avatar_url,
                recipes: u.recipes_aggregate?.aggregate?.count ?? 0,
            }));
    } catch (err) {
        loadError.value = err?.message || "Failed to load home page data.";
    } finally {
        recipesLoading.value = false;
        categoriesLoading.value = false;
        creatorsLoading.value = false;
    }
});
</script>
