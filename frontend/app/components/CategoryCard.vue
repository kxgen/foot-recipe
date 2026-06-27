<!-- xo -->
<!-- Catagory.vue -->
<script setup>
import { computed } from "vue";
import { categoryIcons } from "~/config/category-ui";

// props from parent
const props = defineProps({
    id: Number,
    name: String,
    slug: String,
    count: Number,
});

// icon lookup
const icon = computed(() => {
    return categoryIcons[props.slug] || null;
});
</script>

<template>
    <NuxtLink
        :to="`/categories/${slug}`"
        class="group card-surface flex items-center gap-4 p-5 transition hover:bg-accent-subtle/40 hover:border-primary/30 hover:-translate-y-0.5"
    >
        <!-- icon lucid image -->
        <div
            class="flex size-10 items-center justify-center rounded-xl bg-canvas border border-primary/20 text-body transition group-hover:scale-105"
            aria-hidden="true"
        >
            <component v-if="icon" :is="icon" class="w-5 h-5" />
            <span v-else>🍳</span>
        </div>

        <!-- catagory name -->
        <div class="min-w-0">
            <div class="text-sm font-semibold text-heading truncate">
                {{ name }}
            </div>

            <div class="mt-1 text-xs text-subtle">
                {{ count }} recipe{{ count === 1 ? "" : "s" }}
            </div>
        </div>
    </NuxtLink>
</template>
