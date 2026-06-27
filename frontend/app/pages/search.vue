<template>
    <div class="pb-10 pt-2 sm:pb-14 sm:pt-4">
        <div class="page-container">
            <div class="mx-auto max-w-5xl">
                <div class="flex flex-col gap-4 sm:flex-row sm:items-end">
                    <div class="grow">
                        <RecipeSearchBar :initial-query="q" />
                    </div>

                    <SearchSort
                        v-model:sortBy="sortBy"
                        v-model:sortOrder="sortOrder"
                        @change="runSearch"
                    />
                </div>
                <div>
                    <p v-if="q" class="mt-2 text-sm text-subtle">
                        Showing results for
                        <span class="font-semibold text-body"
                            >&quot;{{ q }}&quot;</span
                        >
                    </p>
                </div>
            </div>

            <div
                class="mt-8 h-px bg-linear-to-r from-primary/0 via-primary/35 to-primary/0"
                role="separator"
                aria-hidden="true"
            />

            <div class="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-12">
                <!-- FILTERS ASIDE -->
                <aside class="lg:col-span-3 lg:order-1">
                    <SearchFilters
                        v-model:selectedIngredients="selectedIngredients"
                        v-model:selectedTimeBucketIds="selectedTimeBucketIds"
                        :all-ingredients="allIngredients"
                        @change="runSearch"
                        @clear="runSearch"
                    />
                </aside>

                <!-- RESULTS SECTION -->
                <section class="lg:col-span-9 lg:order-2">
                    <div
                        v-if="q || hasActiveFilters"
                        class="mb-6 flex items-center justify-between gap-4"
                    >
                        <h2 class="text-xl font-semibold">
                            Recipe results
                        </h2>
                        <div class="text-sm text-subtle">
                            <span v-if="loading">
                                <InlineLoader />
                            </span>
                            <span v-else>
                                {{ results.length }} recipe{{
                                    results.length === 1 ? "" : "s"
                                }}
                            </span>
                        </div>
                    </div>

                    <p
                        v-if="loadError"
                        class="mb-4 error-message"
                    >
                        {{ loadError }}
                    </p>

                    <RecipesSection
                        v-if="results.length"
                        :title="''"
                        :recipes="results"
                    />

                    <SearchEmptyState
                        v-else-if="!loading && (q || hasActiveFilters)"
                        type="no-results"
                    />

                    <!-- BLANK STATE WHEN NO SEARCH -->
                    <SearchEmptyState v-else-if="!loading" type="initial" />
                </section>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import { GET_RECIPES, GET_FILTER_OPTIONS } from "~/graphql/queries.js";
import { mapRecipeCard } from "~/utils/recipe-mappers.js";
import { TIME_BUCKETS } from "~/config/search.js";

const { $apollo: apollo } = useNuxtApp();
const query = (options) =>
    apollo.query({ fetchPolicy: "network-only", ...options });

function buildRecipesWhere({
    titleQuery,
    categoryIds,
    userIds,
    ingredientNames,
    prepTimeMin,
    prepTimeMax,
    timeBucketIds,
}) {
    const and = [];

    const trimmedTitle = titleQuery?.trim();
    if (trimmedTitle) {
        and.push({ title: { _ilike: `%${trimmedTitle}%` } });
    }

    if (categoryIds?.length) {
        and.push({ category_id: { _in: categoryIds } });
    }

    if (userIds?.length) {
        and.push({ user_id: { _in: userIds } });
    }

    if (ingredientNames?.length) {
        and.push({
            recipe_ingredients: { name: { _in: ingredientNames } },
        });
    }

    if (timeBucketIds?.length) {
        const orBuckets = TIME_BUCKETS.filter((b) =>
            timeBucketIds.includes(b.id),
        ).map((b) => {
            if (b.max == null) {
                return { total_time_minutes: { _gte: b.min } };
            }
            return {
                total_time_minutes: { _gte: b.min, _lte: b.max },
            };
        });
        if (orBuckets.length === 1) and.push(orBuckets[0]);
        else if (orBuckets.length) and.push({ _or: orBuckets });
    } else if (prepTimeMin != null || prepTimeMax != null) {
        const prep = {};
        if (prepTimeMin != null) prep._gte = prepTimeMin;
        if (prepTimeMax != null) prep._lte = prepTimeMax;
        and.push({ total_time_minutes: prep });
    }

    if (!and.length) return {};
    if (and.length === 1) return and[0];
    return { _and: and };
}

function prepRangeFromBucketIds(bucketIds) {
    if (!bucketIds?.length) return { mins: [], maxs: [] };

    const selected = TIME_BUCKETS.filter((b) => bucketIds.includes(b.id));
    return {
        mins: selected.map((b) => b.min),
        maxs: selected.map((b) => (b.max == null ? 9999 : b.max)),
    };
}

function mapRecipes(list) {
    return (list ?? []).map(mapRecipeCard).filter(Boolean);
}

async function searchRecipes({
    titleQuery,
    categoryIds,
    userIds,
    ingredientNames,
    timeBucketIds,
    limit,
    orderBy,
}) {
    const { mins, maxs } = prepRangeFromBucketIds(timeBucketIds);

    const activeOrderBy = orderBy ?? [{ created_at: "desc" }];

    if (timeBucketIds?.length) {
        const merged = new Map();

        for (let i = 0; i < mins.length; i++) {
            const where = buildRecipesWhere({
                titleQuery,
                categoryIds,
                userIds,
                ingredientNames,
                prepTimeMin: mins[i],
                prepTimeMax: maxs[i],
            });

            const { data } = await query({
                query: GET_RECIPES,
                variables: {
                    where,
                    limit: limit ?? 50,
                    orderBy: activeOrderBy,
                },
            });

            for (const recipe of data?.recipes ?? []) {
                merged.set(recipe.id, recipe);
            }
        }

        return mapRecipes([...merged.values()]);
    }

    const where = buildRecipesWhere({
        titleQuery,
        categoryIds,
        userIds,
        ingredientNames,
    });

    const { data } = await query({
        query: GET_RECIPES,
        variables: {
            where,
            limit: limit ?? 50,
            orderBy: activeOrderBy,
        },
    });

    return mapRecipes(data?.recipes);
}

const route = useRoute();

const q = computed(() => {
    const raw = route.query.q;
    if (typeof raw === "string") return raw;
    if (Array.isArray(raw)) return raw[0] ?? "";
    return "";
});

useHead(() => ({
    title: q.value
        ? `Search: ${q.value} — Share Plate`
        : "Search — Share Plate",
}));

const allIngredients = ref([]);
const selectedIngredients = ref([]);
const selectedTimeBucketIds = ref([]);

const sortBy = ref("created_at");
const sortOrder = ref("desc");

const results = ref([]);
const loading = ref(false);
const loadError = ref("");

const hasActiveFilters = computed(
    () =>
        selectedIngredients.value.length > 0 ||
        selectedTimeBucketIds.value.length > 0,
);

watch(q, () => {
    runSearch();
});

async function runSearch() {
    const title = q.value.trim();
    const active = title || hasActiveFilters.value;

    if (!active) {
        results.value = [];
        loading.value = false;
        return;
    }

    loading.value = true;
    loadError.value = "";

    try {
        results.value = await searchRecipes({
            titleQuery: title,
            ingredientNames: selectedIngredients.value,
            timeBucketIds: selectedTimeBucketIds.value,
            orderBy: [{ [sortBy.value]: sortOrder.value }],
        });
    } catch (err) {
        loadError.value = err?.message || "Search failed.";
        results.value = [];
    } finally {
        loading.value = false;
    }
}

onMounted(async () => {
    try {
        const { data } = await query({ query: GET_FILTER_OPTIONS });

        allIngredients.value = (data?.recipe_ingredients ?? [])
            .map((i) => i.name)
            .filter(Boolean);
    } catch (err) {
        loadError.value = err?.message || "Failed to load filters.";
    }

    await runSearch();
});
</script>
