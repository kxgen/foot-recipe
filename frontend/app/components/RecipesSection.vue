<!-- Search.vue -->
<script setup>
defineProps({
  title: { type: String, default: "" },
  recipes: { type: Array, default: () => [] },
  muted: { type: Boolean, default: false },
  viewMoreLink: { type: String, default: "" },
})
</script>

<template>
  <section class="py-14" :class="muted ? 'bg-surface-muted/35' : ''">
    <div class="page-container">
      <div v-if="title" class="mb-8 flex items-center justify-between gap-4">
        <h2 class="text-2xl font-semibold text-heading">{{ title }}</h2>
        <NuxtLink
          v-if="viewMoreLink"
          :to="viewMoreLink"
          class="text-sm font-semibold text-heading-muted hover:underline"
        >
          View more →
        </NuxtLink>
      </div>

      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <RecipeCard
          v-for="recipe in recipes"
          :key="recipe.id ?? recipe.title"
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
    </div>
  </section>
</template>


