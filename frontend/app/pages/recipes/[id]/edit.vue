<script setup>
import { getUserIdFromToken } from "~/utils/jwt";
import { useNotifications } from "~/composables/useNotifications";
import { gql } from "@apollo/client/core";
import { GET_CATEGORIES, GET_RECIPE_BY_ID } from "~/graphql/queries.js"
import { UPDATE_RECIPE_WITH_RELATIONS, DELETE_RECIPE } from "~/graphql/mutations.js"

const route = useRoute();
const router = useRouter();
const { token } = useAuth();
const { $apollo: apollo } = useNuxtApp();
const { notifySuccess, notifyError } = useNotifications();

const recipeId = computed(() => Number(route.params.id));
const currentUserId = computed(() => getUserIdFromToken(token.value));

const categories = ref([]);
const loadError = ref("");
const submitError = ref("");
const loading = ref(true);
const submitting = ref(false);
let stepKeyCounter = 0;

function createStep(instruction = "", durationMinutes = 5, imageUrl = "") {
    stepKeyCounter += 1;
    return {
        key: stepKeyCounter,
        instruction,
        durationMinutes,
        imageUrl,
    };
}

function createIngredient(quantity = "", name = "", unit = "") {
    return { quantity, name, unit };
}

const form = reactive({
    title: "",
    description: "",
    categoryId: null,
    isFree: true,
    price: null,
    featuredImage: "",
    steps: [],
    ingredients: [],
});

const totalPrepTime = computed(() => {
    return form.steps.reduce((acc, step) => acc + (Number(step.durationMinutes) || 0), 0);
});

useHead({
    title: "Edit recipe — Share Plate",
});

onMounted(async () => {
    if (!currentUserId.value) {
        await navigateTo("/auth/login");
        return;
    }

    try {
        const [recipeRes, categoriesRes] = await Promise.all([
            apollo.query({
                query: GET_RECIPE_BY_ID,
                variables: { id: recipeId.value },
                fetchPolicy: "network-only",
            }),
            apollo.query({
                query: GET_CATEGORIES,
                fetchPolicy: "network-only",
            }),
        ]);

        categories.value = categoriesRes.data?.categories ?? [];
        const recipe = recipeRes.data?.recipes_by_pk;

        if (!recipe) {
            loadError.value = "Recipe not found.";
            return;
        }

        if (recipe.user_id !== currentUserId.value) {
            loadError.value = "You can only edit your own recipes.";
            return;
        }

        form.title = recipe.title ?? "";
        form.description = recipe.description ?? "";
        form.categoryId = recipe.category_id;

        // Load pricing
        form.isFree = recipe.price === 0 || recipe.price === null;
        form.price = form.isFree ? null : recipe.price;

        // Load featured image
        const featured =
            recipe.recipe_images?.find((img) => img.is_featured) ??
            recipe.recipe_images?.[0];
        form.featuredImage = featured?.url ?? "";

        // Load steps
        form.steps = (recipe.recipe_steps ?? []).map((step) =>
            createStep(step.instruction, step.duration_minutes ?? 5, step.image_url),
        );
        if (form.steps.length === 0) {
            form.steps.push(createStep());
        }

        // Load ingredients
        form.ingredients = (recipe.recipe_ingredients ?? []).map((ing) =>
            createIngredient(
                ing.quantity !== null ? String(ing.quantity) : "",
                ing.name ?? "",
                ing.unit ?? "",
            ),
        );
        if (form.ingredients.length === 0) {
            form.ingredients.push(createIngredient());
        }
    } catch (err) {
        loadError.value = err?.message || "Failed to load recipe.";
    } finally {
        loading.value = false;
    }
});

function addStep() {
    form.steps.push(createStep());
}

function removeStep(index) {
    if (form.steps.length === 1) return;
    form.steps.splice(index, 1);
}

function moveStepUp(index) {
    if (index <= 0) return;
    const item = form.steps[index];
    form.steps.splice(index, 1);
    form.steps.splice(index - 1, 0, item);
}

function moveStepDown(index) {
    if (index >= form.steps.length - 1) return;
    const item = form.steps[index];
    form.steps.splice(index, 1);
    form.steps.splice(index + 1, 0, item);
}

function addIngredient() {
    form.ingredients.push(createIngredient());
}

function removeIngredient(index) {
    if (form.ingredients.length === 1) return;
    form.ingredients.splice(index, 1);
}

async function updateRecipe() {
    if (!form.isFree && (!form.price || form.price <= 0)) {
        notifyError("Please enter a valid price for paid recipes.");
        return;
    }

    submitting.value = true;
    submitError.value = "";

    const stepsData = form.steps
        .map((step, index) => ({
            recipe_id: recipeId.value,
            step_number: index + 1,
            instruction: step.instruction.trim(),
            duration_minutes: Number(step.durationMinutes) || 0,
            image_url: step.imageUrl || null,
        }))
        .filter((step) => step.instruction);

    const ingredientsData = form.ingredients
        .map((ingredient, index) => ({
            recipe_id: recipeId.value,
            name: ingredient.name.trim(),
            quantity: ingredient.quantity.trim()
                ? Number(ingredient.quantity)
                : null,
            unit: ingredient.unit.trim() || null,
            sort_order: index,
        }))
        .filter((ingredient) => ingredient.name);

    const imagesData = form.featuredImage
        ? [
              {
                  recipe_id: recipeId.value,
                  url: form.featuredImage,
                  sort_order: 0,
                  is_featured: true,
              },
          ]
        : [];

    try {
        await apollo.mutate({
            mutation: UPDATE_RECIPE_WITH_RELATIONS,
            variables: {
                id: recipeId.value,
                set: {
                    title: form.title.trim(),
                    description: form.description.trim(),
                    category_id: form.categoryId,
                    price: form.isFree ? 0 : parseFloat(Number(form.price).toFixed(2)),
                },
                steps: stepsData,
                ingredients: ingredientsData,
                images: imagesData,
            },
        });
        notifySuccess("Recipe updated successfully!");
        await router.push(`/recipes/${recipeId.value}`);
    } catch (err) {
        submitError.value = err?.message || "Failed to update recipe.";
    } finally {
        submitting.value = false;
    }
}

const showDeleteModal = ref(false);

function triggerDelete() {
    showDeleteModal.value = true;
}

async function deleteRecipe() {
    showDeleteModal.value = false;
    submitting.value = true;
    submitError.value = "";

    try {
        await apollo.mutate({
            mutation: DELETE_RECIPE,
            variables: { id: recipeId.value },
        });

        notifySuccess("Recipe deleted successfully!");
        await router.push("/me/recipes");
    } catch (err) {
        submitError.value = err?.message || "Failed to delete recipe.";
    } finally {
        submitting.value = false;
    }
}
</script>

<template>
    <div class="py-12">
        <div class="page-container max-w-4xl">
            <h1 class="text-3xl font-semibold text-heading">Edit recipe</h1>
            <p class="mt-2 text-sm text-subtle">
                Update recipe details, upload a featured image, attach images to
                each step, and set ingredients with quantities.
            </p>

            <p
                v-if="loadError"
                class="mt-4 error-message"
            >
                {{ loadError }}
            </p>
            <p
                v-if="submitError"
                class="mt-4 error-message"
            >
                {{ submitError }}
            </p>

            <div v-if="loading" class="mt-8 text-sm text-subtle">
                <InlineLoader />
            </div>

            <form
                v-else-if="!loadError"
                class="card-surface mt-8 space-y-8 p-6 sm:p-8"
                @submit.prevent="updateRecipe"
            >
                <div>
                    <label class="mb-1.5 block text-sm font-medium text-heading"
                        >Title</label
                    >
                    <input v-model="form.title" class="input-field" required />
                </div>

                <div>
                    <label class="mb-1.5 block text-sm font-medium text-heading"
                        >Description</label
                    >
                    <textarea
                        v-model="form.description"
                        class="input-field min-h-24"
                        required
                    />
                </div>

                <div class="grid gap-4 sm:grid-cols-2">
                    <div>
                        <label class="mb-1.5 block text-sm font-medium text-heading"
                            >Category</label
                        >
                        <select
                            v-model="form.categoryId"
                            class="input-field"
                            required
                        >
                            <option disabled :value="null">
                                Choose category
                            </option>
                            <option
                                v-for="c in categories"
                                :key="c.id"
                                :value="c.id"
                            >
                                {{ c.name }}
                            </option>
                        </select>
                    </div>
                    <div>
                        <label class="mb-1.5 block text-sm font-medium text-heading"
                            >Total Preparation Time (mins)</label
                        >
                        <div class="input-field bg-primary/5 flex items-center font-semibold text-primary">
                            {{ totalPrepTime }} mins
                        </div>
                    </div>
                </div>

                <div
                    class="rounded-2xl border border-primary/15 bg-surface-muted/40 p-5"
                >
                    <h2 class="text-sm font-semibold text-heading">Pricing</h2>
                    <p class="mt-1 text-xs text-subtle">
                        Choose whether this recipe is free or paid.
                    </p>

                    <div class="mt-4 flex flex-wrap gap-4">
                        <label
                            class="inline-flex cursor-pointer items-center gap-2 text-sm text-heading"
                        >
                            <input
                                v-model="form.isFree"
                                type="radio"
                                :value="true"
                                class="accent-primary"
                            />
                            Free
                        </label>
                        <label
                            class="inline-flex cursor-pointer items-center gap-2 text-sm text-heading"
                        >
                            <input
                                v-model="form.isFree"
                                type="radio"
                                :value="false"
                                class="accent-primary"
                            />
                            Paid
                        </label>
                    </div>

                    <div v-if="!form.isFree" class="mt-4 max-w-xs">
                        <label class="mb-1.5 block text-sm font-medium text-heading"
                            >Price (ETB)</label
                        >
                        <input
                            v-model.number="form.price"
                            type="number"
                            min="0.01"
                            step="0.01"
                            class="input-field"
                            required
                            placeholder="200"
                        />
                    </div>
                </div>

                <div
                    class="rounded-2xl border border-primary/15 bg-surface-muted/40 p-5"
                >
                    <h2 class="text-sm font-semibold text-heading">
                        Featured image
                    </h2>
                    <p class="mt-1 text-xs text-subtle">
                        This image appears on recipe cards and at the top of the
                        recipe page.
                    </p>
                    <div class="mt-4">
                        <ImageFileInput
                            v-model="form.featuredImage"
                            input-id="featured-image"
                            label="Upload featured image"
                        />
                    </div>
                </div>

                <div>
                    <div class="mb-3">
                        <div>
                            <label class="text-sm font-medium text-heading"
                                >Steps</label
                            >
                            <p class="text-xs text-subtle">
                                Move steps up or down to reorder. Add an image
                                to each step if you like.
                            </p>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <div
                            v-for="(step, i) in form.steps"
                            :key="step.key"
                            class="rounded-2xl border border-primary/15 bg-white/70 p-4"
                        >
                            <div class="mb-3 flex items-center gap-2">
                                <span
                                    class="select-none rounded-lg bg-primary-subtle/25 px-2 py-1 text-xs font-semibold text-primary"
                                >
                                    Step {{ i + 1 }}
                                </span>

                                <div class="flex items-center gap-1 ml-2">
                                    <button
                                        type="button"
                                        class="p-1 rounded-md hover:bg-primary/10 disabled:opacity-30 disabled:cursor-not-allowed"
                                        :disabled="i === 0"
                                        title="Move up"
                                        @click="moveStepUp(i)"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            width="16"
                                            height="16"
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        >
                                            <path d="m18 15-6-6-6 6" />
                                        </svg>
                                    </button>
                                    <button
                                        type="button"
                                        class="p-1 rounded-md hover:bg-primary/10 disabled:opacity-30 disabled:cursor-not-allowed"
                                        :disabled="i === form.steps.length - 1"
                                        title="Move down"
                                        @click="moveStepDown(i)"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            width="16"
                                            height="16"
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        >
                                            <path d="m6 9 6 6 6-6" />
                                        </svg>
                                    </button>
                                </div>

                                <button
                                    type="button"
                                    class="ml-auto btn-secondary px-3 py-1.5 text-xs"
                                    :disabled="form.steps.length === 1"
                                    @click="removeStep(i)"
                                >
                                    Remove
                                </button>
                            </div>

                            <div class="mb-3 grid gap-4 sm:grid-cols-[1fr_120px]">
                                <div>
                                    <label class="mb-1 block text-xs font-medium text-subtle">Instruction</label>
                                    <textarea
                                        v-model="step.instruction"
                                        class="input-field min-h-20"
                                        :placeholder="`Describe step ${i + 1}`"
                                        required
                                    />
                                </div>
                                <div>
                                    <label class="mb-1 block text-xs font-medium text-subtle">Time (mins)</label>
                                    <input
                                        v-model.number="step.durationMinutes"
                                        type="number"
                                        min="0"
                                        class="input-field"
                                        required
                                    />
                                </div>
                            </div>

                            <div class="mt-3">
                                <ImageFileInput
                                    v-model="step.imageUrl"
                                    :input-id="`step-image-${step.key}`"
                                    label="Add step image"
                                />
                            </div>
                        </div>
                    </div>

                    <div class="mt-4 flex justify-end">
                        <button
                            type="button"
                            class="text-xs font-semibold text-primary hover:underline"
                            @click="addStep"
                        >
                            + Add step
                        </button>
                    </div>
                </div>

                <div>
                    <div class="mb-3">
                        <div>
                            <label class="text-sm font-medium text-heading"
                                >Ingredients</label
                            >
                            <p class="text-xs text-subtle">
                                Quantity, name, and unit for each ingredient.
                            </p>
                        </div>
                    </div>

                    <div class="space-y-2">
                        <div
                            v-for="(ingredient, i) in form.ingredients"
                            :key="`ingredient-${i}`"
                            class="grid gap-2 sm:grid-cols-[minmax(0,1fr)_minmax(0,2fr)_minmax(0,1fr)_auto]"
                        >
                            <input
                                v-model="ingredient.quantity"
                                type="text"
                                class="input-field"
                                placeholder="Qty"
                            />
                            <input
                                v-model="ingredient.name"
                                class="input-field"
                                :placeholder="`Ingredient ${i + 1}`"
                            />
                            <input
                                v-model="ingredient.unit"
                                class="input-field"
                                placeholder="Unit"
                            />
                            <button
                                type="button"
                                class="btn-secondary px-3"
                                @click="removeIngredient(i)"
                            >
                                ×
                            </button>
                        </div>
                    </div>

                    <div class="mt-4 flex justify-end">
                        <button
                            type="button"
                            class="text-xs font-semibold text-primary hover:underline"
                            @click="addIngredient"
                        >
                            + Add ingredient
                        </button>
                    </div>
                </div>

                <div class="flex flex-wrap gap-3 pt-2">
                    <button
                        type="submit"
                        class="btn-primary"
                        :disabled="submitting"
                    >
                        {{ submitting ? "Saving..." : "Update recipe" }}
                    </button>
                    <button
                        type="button"
                        class="btn-secondary border-danger/30 text-danger hover:bg-danger-bg/50"
                        :disabled="submitting"
                        @click="triggerDelete"
                    >
                        Delete recipe
                    </button>
                </div>
            </form>
        </div>

        <ConfirmModal
            :is-open="showDeleteModal"
            title="Delete Recipe"
            message="Are you sure you want to permanently delete this recipe? This action cannot be undone."
            confirm-text="Delete"
            cancel-text="Cancel"
            @confirm="deleteRecipe"
            @close="showDeleteModal = false"
        />
    </div>
</template>
