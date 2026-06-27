<script setup>
defineProps({
    type: {
        type: String,
        required: true,
        validator: (value) => ["initial", "no-results"].includes(value),
    },
});
</script>

<template>
    <div class="flex flex-col items-center justify-center py-20 text-center">
        <!-- ICON -->
        <div
            class="mb-4 flex size-16 items-center justify-center rounded-full bg-accent-subtle text-primary"
            :class="type === 'no-results' ? 'text-primary/60' : ''"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-8"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
                />
                <path
                    v-if="type === 'no-results'"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M6 18 18 6M6 6l12 12"
                />
            </svg>
        </div>

        <!-- TEXT -->
        <template v-if="type === 'no-results'">
            <h3 class="text-lg font-semibold text-heading">No recipes found</h3>
            <p class="mt-2 max-w-sm text-sm text-subtle">
                Try changing your filters or searching a different title to find
                what you're looking for.
            </p>
        </template>
        <template v-else>
            <h3 class="text-lg font-semibold text-heading">Ready to cook?</h3>
            <p class="mt-2 max-w-sm text-sm text-subtle">
                Search for recipes by title, ingredients, or use filters to find
                exactly what you're looking for.
            </p>
        </template>
    </div>
</template>
