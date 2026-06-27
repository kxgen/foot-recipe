<script setup>
import { ref, computed } from "vue";
import { onClickOutside } from "@vueuse/core";
import { TIME_BUCKETS, POPULAR_INGREDIENTS } from "~/config/search.js";

const props = defineProps({
    allIngredients: {
        type: Array,
        default: () => [],
    },
});

const selectedIngredients = defineModel("selectedIngredients", {
    default: () => [],
});
const selectedTimeBucketIds = defineModel("selectedTimeBucketIds", {
    default: () => [],
});

const emit = defineEmits(["change", "clear"]);

const ingredientSearchQuery = ref("");
const isIngredientDropdownOpen = ref(false);
const ingredientDropdownRef = ref(null);

const filteredIngredientSuggestions = computed(() => {
    const search = ingredientSearchQuery.value.toLowerCase().trim();
    if (!search) return [];

    return props.allIngredients
        .filter(
            (ing) =>
                ing.toLowerCase().includes(search) &&
                !selectedIngredients.value.includes(ing),
        )
        .slice(0, 8);
});

const hasActiveFilters = computed(
    () =>
        selectedIngredients.value.length > 0 ||
        selectedTimeBucketIds.value.length > 0,
);

function toggleInArray(list, value) {
    const index = list.indexOf(value);
    if (index === -1) list.push(value);
    else list.splice(index, 1);
}

function toggleIngredient(ing) {
    toggleInArray(selectedIngredients.value, ing);
    emit("change");
}

function addIngredientFromSearch(ing) {
    if (!selectedIngredients.value.includes(ing)) {
        selectedIngredients.value.push(ing);
        emit("change");
    }
    ingredientSearchQuery.value = "";
    isIngredientDropdownOpen.value = false;
}

function toggleTimeBucket(id) {
    toggleInArray(selectedTimeBucketIds.value, id);
    emit("change");
}

function clearAllFilters() {
    selectedIngredients.value = [];
    selectedTimeBucketIds.value = [];
    emit("clear");
}

onClickOutside(ingredientDropdownRef, () => {
    isIngredientDropdownOpen.value = false;
});
</script>

<template>
    <div
        class="sticky top-24 rounded-2xl border border-primary/15 bg-pure p-5 backdrop-blur"
    >
        <div class="flex items-center justify-between mb-5">
            <h2 class="text-sm font-semibold uppercase tracking-wider text-heading">
                Filters
            </h2>
            <button
                v-if="hasActiveFilters"
                @click="clearAllFilters"
                class="text-xs font-medium text-primary hover:underline"
            >
                Clear all
            </button>
        </div>

        <div class="space-y-8">
            <!-- PREP TIME FILTER -->
            <section>
                <h3 class="text-sm font-semibold text-heading">Preparation time</h3>
                <div class="mt-3 space-y-2">
                    <label
                        v-for="b in TIME_BUCKETS"
                        :key="b.id"
                        class="flex cursor-pointer items-center justify-between gap-3 rounded-xl border border-primary/15 bg-canvas/60 px-3 py-2 transition hover:bg-layer"
                    >
                        <span class="text-xs font-medium text-subtle">{{
                            b.label
                        }}</span>
                        <input
                            type="checkbox"
                            class="size-4 rounded border-primary/30 text-primary focus:ring-primary/40"
                            :checked="selectedTimeBucketIds.includes(b.id)"
                            @change="toggleTimeBucket(b.id)"
                        />
                    </label>
                </div>
            </section>

            <!-- INGREDIENTS FILTER -->
            <section>
                <h3 class="text-sm font-semibold text-heading">Ingredients</h3>

                <!-- Selected Ingredient Chips -->
                <div
                    v-if="selectedIngredients.length"
                    class="mt-3 flex flex-wrap gap-2 mb-4"
                >
                    <button
                        v-for="ing in selectedIngredients"
                        :key="ing"
                        @click="toggleIngredient(ing)"
                        class="inline-flex items-center gap-1 rounded-full bg-accent-subtle border border-primary/30 px-3 py-1 text-xs font-medium text-body transition hover:bg-cta hover:text-inverted group"
                    >
                        {{ ing }}
                        <span
                            class="text-[10px] opacity-60 group-hover:opacity-100"
                            >✕</span
                        >
                    </button>
                </div>

                <!-- Ingredient Autocomplete Search -->
                <div class="relative mt-3" ref="ingredientDropdownRef">
                    <input
                        v-model="ingredientSearchQuery"
                        type="text"
                        placeholder="Add ingredient..."
                        class="input-field py-2 text-xs"
                        @focus="isIngredientDropdownOpen = true"
                    />

                    <!-- Dropdown suggestions -->
                    <div
                        v-if="
                            isIngredientDropdownOpen &&
                            filteredIngredientSuggestions.length
                        "
                        class="absolute left-0 right-0 z-50 mt-1 max-h-48 overflow-y-auto rounded-xl border border-primary/15 bg-canvas shadow-lg"
                    >
                        <button
                            v-for="ing in filteredIngredientSuggestions"
                            :key="ing"
                            @click="addIngredientFromSearch(ing)"
                            class="block w-full px-4 py-2 text-left text-xs text-subtle transition hover:bg-accent-subtle hover:text-body"
                        >
                            {{ ing }}
                        </button>
                    </div>
                </div>

                <!-- Popular Ingredients -->
                <div class="mt-6">
                    <h4
                        class="text-[11px] font-bold uppercase tracking-widest text-subtle/70"
                    >
                        Popular
                    </h4>
                    <div class="mt-3 flex flex-wrap gap-2">
                        <button
                            v-for="ing in POPULAR_INGREDIENTS"
                            :key="ing"
                            type="button"
                            class="rounded-full border px-3 py-1 text-[11px] transition"
                            :class="
                                selectedIngredients.includes(ing)
                                    ? 'border-cta bg-accent-subtle text-body'
                                    : 'border-primary/25 bg-canvas/60 text-subtle hover:bg-layer'
                            "
                            @click="toggleIngredient(ing)"
                        >
                            {{ ing }}
                        </button>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>
