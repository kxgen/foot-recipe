<script setup>
import { getUserIdFromToken } from "~/utils/jwt";
import { useNotifications } from "~/composables/useNotifications";
import {
    GET_RECIPE_BY_ID,
    GET_RECIPE_LIKE_COUNT,
    CHECK_USER_RECIPE_FLAGS,
} from "~/graphql/queries.js";
import {
    INSERT_RECIPE_LIKE,
    DELETE_RECIPE_LIKE,
    INSERT_RECIPE_BOOKMARK,
    DELETE_RECIPE_BOOKMARK,
    INITIATE_PURCHASE,
} from "~/graphql/mutations.js";

const route = useRoute();
const { $apollo: apollo } = useNuxtApp();
const query = (options) =>
    apollo.query({ fetchPolicy: "network-only", ...options });
const mutate = (options) => apollo.mutate(options);
const { token, isLoggedIn } = useAuth();
const { notifySuccess, notifyError } = useNotifications();

const recipeId = computed(() => Number(route.params.id));
const currentUserId = computed(() => getUserIdFromToken(token.value));

const recipe = ref(null);
const loading = ref(true);
const loadError = ref("");
const liked = ref(false);
const bookmarked = ref(false);
const purchased = ref(false);
const purchasing = ref(false);
const total_likes = ref(0);
const actionError = ref("");

const featuredImage = computed(() => {
    const images = recipe.value?.recipe_images ?? [];
    return images.find((img) => img.is_featured)?.url ?? images[0]?.url ?? null;
});

const canEdit = computed(
    () => currentUserId.value && recipe.value?.user_id === currentUserId.value,
);

const isPaidRecipe = computed(() => Number(recipe.value?.price) > 0);

const canPurchase = computed(
    () => !canEdit.value && !purchased.value && isPaidRecipe.value,
);

const purchasePrice = computed(() => {
    const price = Number(recipe.value?.price);
    return Number.isFinite(price) && price > 0 ? price : null;
});

useHead(() => ({
    title: recipe.value?.title
        ? `${recipe.value.title} — Share Plate`
        : "Recipe — Share Plate",
}));

const router = useRouter();
const config = useRuntimeConfig();
const apiBase = config.public.apiBase || "http://localhost:8000";

onMounted(() => {
    loadRecipe();

    // Check for payment status from Chapa redirect
    if (route.query.payment === "success") {
        notifySuccess("Payment successful! You can now view the recipe.");
        // Remove query param without refreshing
        router.replace({ query: { ...route.query, payment: undefined } });
    } else if (route.query.payment === "failed") {
        notifyError("Payment failed. Please try again.");
        router.replace({ query: { ...route.query, payment: undefined } });
    }
});

watch(recipeId, () => loadRecipe());

async function loadRecipe(silent = false) {
    if (!recipeId.value || Number.isNaN(recipeId.value)) return;

    if (!silent) loading.value = true;
    loadError.value = "";

    try {
        const { data } = await query({
            query: GET_RECIPE_BY_ID,
            variables: { id: recipeId.value },
        });

        if (data?.recipes_by_pk) {
            recipe.value = { ...data.recipes_by_pk };
        } else {
            recipe.value = null;
        }

        if (isLoggedIn.value && currentUserId.value) {
            await loadUserFlags();
        }
    } catch (err) {
        if (!silent) {
            loadError.value = err?.message || "Failed to load recipe.";
            recipe.value = null;
        }
    } finally {
        if (!silent) loading.value = false;
    }
}

async function loadUserFlags() {
    const { data } = await query({
        query: CHECK_USER_RECIPE_FLAGS,
        variables: {
            recipeId: recipeId.value,
            userId: currentUserId.value,
        },
    });

    liked.value = (data?.recipe_likes?.length ?? 0) > 0;
    bookmarked.value = (data?.recipe_bookmarks?.length ?? 0) > 0;
    purchased.value = (data?.recipe_purchases?.length ?? 0) > 0;
}

async function toggleLike() {
    if (!requireAuth()) return;
    actionError.value = "";

    try {
        if (liked.value) {
            const { data } = await mutate({
                mutation: DELETE_RECIPE_LIKE,
                variables: {
                    recipeId: recipeId.value,
                    userId: currentUserId.value,
                },
            });

            // Extract the new count from the deletion returning array
            const newCount =
                data?.delete_recipe_likes?.returning[0]?.recipe?.like_count;
            if (typeof newCount === "number") {
                recipe.value.like_count = newCount;
            }

            liked.value = false;
            notifySuccess("Recipe unliked.");
        } else {
            const { data } = await mutate({
                mutation: INSERT_RECIPE_LIKE,
                variables: {
                    recipeId: recipeId.value,
                },
            });

            // Extract the new count from the single object insert
            const newCount = data?.insert_recipe_likes_one?.recipe?.like_count;
            if (typeof newCount === "number") {
                recipe.value.like_count = newCount;
            }

            liked.value = true;
            notifySuccess("Recipe liked!");
        }
    } catch (err) {
        notifyError(err?.message || "Could not update like.");
    }
}

async function toggleBookmark() {
    if (!requireAuth()) return;
    actionError.value = "";

    try {
        if (bookmarked.value) {
            await mutate({
                mutation: DELETE_RECIPE_BOOKMARK,
                variables: {
                    recipeId: recipeId.value,
                    userId: currentUserId.value,
                },
            });
            bookmarked.value = false;
            notifySuccess("Bookmark removed.");
        } else {
            await mutate({
                mutation: INSERT_RECIPE_BOOKMARK,
                variables: {
                    recipeId: recipeId.value,
                },
            });
            bookmarked.value = true;
            notifySuccess("Recipe bookmarked!");
        }
    } catch (err) {
        notifyError(err?.message || "Could not update bookmark.");
    }
}

async function purchaseRecipe() {
    if (!requireAuth()) return;
    if (purchasing.value) return;

    actionError.value = "";
    purchasing.value = true;

    try {
        const { data } = await mutate({
            mutation: INITIATE_PURCHASE,
            variables: {
                recipeId: recipeId.value,
                userId: currentUserId.value,
            },
        });

        const checkoutUrl = data?.initiatePurchase?.checkoutUrl;

        if (checkoutUrl) {
            // Redirect to Chapa payment page
            window.location.href = checkoutUrl;
        } else {
            notifyError(
                "Failed to initialize payment. No checkout URL returned.",
            );
            purchasing.value = false;
        }
    } catch (err) {
        notifyError(err?.message || "Purchase failed.");
        purchasing.value = false;
    }
}

function requireAuth() {
    if (!isLoggedIn.value || !currentUserId.value) {
        navigateTo("/auth/login");
        return false;
    }
    return true;
}
</script>

<template>
    <div class="py-12">
        <div class="page-container">
            <div class="pb-10 pt-2 sm:pb-14 sm:pt-4">
                <div class="mb-6">
                    <NuxtLink
                        to="/recipes"
                        class="text-sm text-primary hover:underline"
                    >
                        ← Back to recipes
                    </NuxtLink>
                </div>

                <div v-if="loading" class="text-sm text-subtle">
                    <InlineLoader />
                </div>

                <p v-if="purchaseError" class="error-message">
                    {{ purchaseError }}
                </p>

                <template v-else-if="recipe">
                    <article class="card-surface overflow-hidden">
                        <div
                            class="relative flex h-80 w-full items-center justify-center overflow-hidden rounded-xl bg-black/10"
                        >
                            <span
                                v-if="!featuredImage"
                                class="text-xs font-medium uppercase tracking-wider text-subtle/40"
                            >
                                No feature image
                            </span>

                            <template v-else>
                                <img
                                    :src="featuredImage"
                                    alt=""
                                    aria-hidden="true"
                                    class="absolute inset-0 h-full w-full object-cover opacity-40 blur-[2px]"
                                />

                                <div
                                    class="absolute inset-0 bg-black/20 z-0 pointer-events-none"
                                    aria-hidden="true"
                                />

                                <img
                                    :src="featuredImage"
                                    alt="Featured Recipe Image"
                                    class="relative z-10 h-full w-full object-contain object-center mask-[radial-gradient(circle,black_60%,transparent_100%)]"
                                />
                            </template>
                        </div>
                        <div class="p-6 sm:p-8">
                            <div
                                class="flex flex-wrap items-start justify-between gap-4"
                            >
                                <div>
                                    <h1
                                        class="text-3xl font-semibold text-heading"
                                    >
                                        {{ recipe.title }}
                                    </h1>
                                    <p class="mt-1 text-sm text-subtle">
                                        By
                                        <NuxtLink
                                            v-if="recipe.user?.slug"
                                            :to="`/creators/${recipe.user.slug}`"
                                            class="font-semibold text-primary hover:underline"
                                        >
                                            {{ recipe.user?.username }}
                                        </NuxtLink>
                                        <span v-else>Unknown</span>
                                        <span v-if="recipe.category?.name">
                                            ·
                                            <NuxtLink
                                                :to="`/categories/${recipe.category.slug}`"
                                                class="hover:underline"
                                            >
                                                {{ recipe.category.name }}
                                            </NuxtLink>
                                        </span>
                                    </p>
                                </div>
                            </div>

                            <p
                                class="mt-5 max-w-3xl text-sm leading-relaxed text-subtle"
                            >
                                {{ recipe.description }}
                            </p>

                            <div
                                class="mt-5 flex flex-wrap gap-4 text-sm text-subtle"
                            >
                                <span class="flex items-center gap-1">
                                    ❤️ {{ recipe.like_count || 0 }}
                                </span>
                                <span v-if="recipe.total_time_minutes"
                                    >🕒
                                    {{ recipe.total_time_minutes }} mins</span
                                >
                                <span v-if="recipe.difficulty"
                                    >Difficulty: {{ recipe.difficulty }}</span
                                >
                            </div>

                            <div class="mt-8 flex flex-wrap items-center gap-4">
                                <template v-if="isLoggedIn">
                                    <button
                                        type="button"
                                        :class="
                                            liked
                                                ? 'btn-primary'
                                                : 'btn-secondary'
                                        "
                                        class="min-w-25"
                                        @click="toggleLike"
                                    >
                                        {{ liked ? "Liked" : "Like" }}
                                    </button>
                                    <button
                                        type="button"
                                        :class="
                                            bookmarked
                                                ? 'btn-primary'
                                                : 'btn-secondary'
                                        "
                                        class="min-w-25"
                                        @click="toggleBookmark"
                                    >
                                        {{ bookmarked ? "Saved" : "Bookmark" }}
                                    </button>
                                </template>

                                <NuxtLink
                                    v-if="canEdit"
                                    :to="`/recipes/${recipe.id}/edit`"
                                    class="btn-secondary"
                                >
                                    Edit recipe
                                </NuxtLink>



                                <span
                                    v-if="purchased && !canEdit"
                                    class="rounded-xl bg-primary-subtle/25 px-4 py-2 text-sm font-medium text-heading"
                                >
                                    Purchased
                                </span>
                                <span
                                    v-else-if="!canEdit && !isPaidRecipe"
                                    class="rounded-xl bg-accent/30 px-4 py-2 text-sm font-medium text-heading"
                                >
                                    Free recipe
                                </span>
                            </div>

                            <p
                                v-if="actionError"
                                class="mt-4 text-center text-sm text-danger"
                            >
                                {{ actionError }}
                            </p>
                            <!-- Recipe Content Preview and Faded Overlay Paywall -->
                            <div class="relative mt-10" :class="{ 'min-h-[450px]': canPurchase }">
                                <!-- Blurred Content Preview (mimics ingredients & steps hidden under gradient) -->
                                <div :class="{ 'max-h-80 overflow-hidden pointer-events-none select-none filter blur-[1.5px] opacity-60': canPurchase }" class="space-y-10">
                                    <section v-if="recipe.recipe_ingredients?.length">
                                        <h2 class="text-xl font-semibold text-heading">Ingredients</h2>
                                        <ul class="mt-4 list-disc space-y-2 pl-5 text-sm text-subtle">
                                            <li v-for="ing in recipe.recipe_ingredients" :key="ing.id">
                                                <span class="text-heading">{{ ing.name }}</span>
                                                <span v-if="ing.quantity || ing.unit"> — {{ ing.quantity }} {{ ing.unit }} </span>
                                            </li>
                                        </ul>
                                    </section>

                                    <section v-if="recipe.recipe_steps?.length">
                                        <h2 class="text-xl font-semibold text-heading">Preparation Steps</h2>
                                        <ol class="mt-4 space-y-6">
                                            <li v-for="step in recipe.recipe_steps" :key="step.id" class="flex flex-col items-center text-center rounded-2xl border border-primary/15 bg-white/60 p-6 sm:p-8 gap-6">
                                                
                                                <div class="flex-1 w-full">
                                                    <div class="flex items-center justify-center gap-3">
                                                        <span class="inline-flex items-center justify-center bg-primary/10 text-primary rounded-full px-3 py-1 text-sm font-semibold">Step {{ step.step_number }}</span>
                                                        <span v-if="step.duration_minutes" class="text-sm font-medium text-subtle">🕒 {{ step.duration_minutes }} mins</span>
                                                    </div>
                                                    <p class="mt-4 text-base leading-relaxed text-subtle whitespace-pre-line max-w-2xl mx-auto">{{ step.instruction }}</p>
                                                </div>
                                                
                                                <div v-if="step.image_url" class="w-full max-w-xl aspect-auto overflow-hidden rounded-md bg-black/5 flex items-center justify-center shadow-sm">
                                                    <img :src="step.image_url" :alt="'Step ' + step.step_number + ' illustration'" class="h-full w-full object-cover" />
                                                </div>
                                                
                                            </li>
                                        </ol>
                                    </section>
                                </div>

                                <!-- Faded Overlay Paywall (New York Times style) -->
                                <div v-if="canPurchase" class="absolute bottom-0 left-0 right-0 top-0 flex flex-col justify-end bg-gradient-to-t from-white via-white/95 to-transparent">
                                    <div class="pb-12 pt-36 text-center bg-gradient-to-t from-white via-white to-transparent">
                                        <div class="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 text-xl text-primary mb-4 shadow-sm">🔒</div>
                                        <h3 class="text-2xl font-bold text-heading tracking-tight">Unlock Full Recipe</h3>
                                        <p class="mt-2 text-subtle max-w-sm mx-auto text-sm leading-relaxed px-4">
                                            This premium recipe includes the complete ingredients list and step-by-step instructions.
                                        </p>
                                        <button
                                            type="button"
                                            class="btn-primary mt-6 px-10 py-3.5 text-sm font-semibold shadow-md hover:shadow-lg transition duration-200"
                                            :disabled="purchasing"
                                            @click="purchaseRecipe"
                                        >
                                            <span v-if="purchasing">Initializing...</span>
                                            <span v-else>Buy Recipe for ${{ purchasePrice?.toFixed(2) }}</span>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </article>

                    <RecipeComments
                        :recipe-id="recipe.id"
                        :avg-rating="recipe.avg_rating"
                        :rating-count="recipe.rating_count"
                        @refresh="loadRecipe(true)"
                    />
                </template>
            </div>
        </div>
    </div>
</template>
