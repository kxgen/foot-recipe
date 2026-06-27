<!-- xo -->
<!-- Creators.vue -->
<script setup>
const props = defineProps({
    username: String,
    slug: String,
    avatarUrl: String,
    recipeCount: Number,
});
</script>

<template>
    <NuxtLink
        :to="`/creators/${slug}`"
        class="group card-surface flex flex-col items-center p-4 transition hover:bg-accent-subtle/40 hover:border-primary/30 hover:-translate-y-0.5 text-center"
    >
        <!-- avatar -->
        <div class="mb-4">
            <img
                v-if="avatarUrl"
                :src="avatarUrl"
                :alt="username"
                class="size-30 rounded-full object-cover border-2 border-primary-subtle/20 transition group-hover:scale-105"
            />
            <div
                v-else
                class="flex size-30 items-center justify-center rounded-full bg-canvas border border-primary/20 text-2xl font-semibold text-body transition group-hover:scale-105"
            >
                {{ username?.charAt(0).toUpperCase() }}
            </div>
        </div>

        <!-- username -->
        <div class="min-w-0">
            <div class="text-base font-semibold text-heading truncate">
                {{ username }}
            </div>

            <div class="mt-1 text-xs text-subtle">
                {{ recipeCount }} recipe{{ recipeCount === 1 ? "" : "s" }}
            </div>
        </div>
    </NuxtLink>
</template>
