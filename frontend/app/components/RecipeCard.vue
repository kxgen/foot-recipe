<!-- Recipes.vue -->
<script setup>
import { CHECK_USER_BOOKMARK } from "~/graphql/queries.js"
import { INSERT_RECIPE_BOOKMARK, DELETE_RECIPE_BOOKMARK } from "~/graphql/mutations.js"
import { getUserIdFromToken } from "~/utils/jwt";
import { useNotifications } from "~/composables/useNotifications";

const props = defineProps({
  id: Number,
  title: { type: String, required: true },
  creator: { type: String, required: true },
  ratingScore: { type: Number, required: true },
  ratingCount: { type: Number, required: true },
  timeMinutes: { type: Number, required: true },
  imageSrc: String,
  linkable: { type: Boolean, default: true },
  recipeLink: String,
});

const recipeLink = computed(() => {
  if (props.recipeLink) return props.recipeLink;
  return `/recipes/${props.id || "1"}`;
});

const imgSrc = computed(() => props.imageSrc);

const imgFailed = ref(false);
const onImgError = (e) => {
  imgFailed.value = true;
};

const { $apollo: apollo } = useNuxtApp();
const { token, isLoggedIn } = useAuth();
const { notifySuccess, notifyError } = useNotifications();
const currentUserId = computed(() => getUserIdFromToken(token.value));

const saved = ref(false);
const toggling = ref(false);

onMounted(loadBookmarkStatus);
watch(() => token.value, loadBookmarkStatus);

async function loadBookmarkStatus() {
  if (!isLoggedIn.value || !currentUserId.value || !props.id) {
    saved.value = false;
    return;
  }

  try {
    const { data } = await apollo.query({
      query: CHECK_USER_BOOKMARK,
      variables: {
        recipeId: props.id,
        userId: currentUserId.value,
      },
      fetchPolicy: "network-only",
    });
    saved.value = (data?.recipe_bookmarks?.length ?? 0) > 0;
  } catch (err) {
    console.error("Failed to check bookmark status:", err);
  }
}

async function toggleBookmark() {
  if (!requireAuth() || toggling.value || !props.id) return;

  toggling.value = true;
  const previousState = saved.value;

  saved.value = !previousState;

  try {
    if (previousState) {
      await apollo.mutate({
        mutation: DELETE_RECIPE_BOOKMARK,
        variables: {
          recipeId: props.id,
          userId: currentUserId.value,
        },
      });
      notifySuccess("Bookmark removed.");
    } else {
      await apollo.mutate({
        mutation: INSERT_RECIPE_BOOKMARK,
        variables: {
          recipeId: props.id,
        },
      });
      notifySuccess("Recipe bookmarked!");
    }
  } catch (err) {
    saved.value = previousState;
    notifyError(err?.message || "Could not update bookmark.");
  } finally {
    toggling.value = false;
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
  <article
    class="group relative flex h-full flex-col overflow-hidden card-surface transition hover:-translate-y-0.5 hover:shadow-lg"
  >
    <div class="relative bg-canvas">
      <div class="relative aspect-video bg-canvas overflow-hidden border-b border-primary/10">
        <img
          v-if="imgSrc && !imgFailed"
          class="h-full w-full object-cover"
          :src="imgSrc"
          :alt="title"
          loading="lazy"
          @error="onImgError"
        />
        <div
          v-else
          class="flex h-full w-full items-center justify-center bg-canvas"
        >
          <span class="text-xs font-medium tracking-wider text-subtle/40 uppercase"
            >No feature image</span
          >
        </div>
      </div>
      <div
        v-if="imgSrc && !imgFailed"
        class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/35 via-black/0 to-black/0"
        aria-hidden="true"
      />
    </div>

    <div class="flex flex-1 flex-col p-5">
      <h3 class="text-lg font-semibold text-heading group-hover:text-primary">
        <NuxtLink
          v-if="linkable"
          :to="recipeLink"
          class="before:absolute before:inset-0 before:z-10"
        >
          {{ title }}
        </NuxtLink>
        <span v-else>{{ title }}</span>
      </h3>

      <p class="mt-1 text-xs text-subtle">
        By <span class="font-semibold text-body">{{ creator }}</span>
      </p>

      <div class="mt-3 flex items-center justify-between gap-3 pb-4">
        <div class="flex items-center gap-1.5">
          <div
            class="flex items-center gap-0.5 font-medium text-body"
            aria-label="Rating score"
          >
            <span>{{ ratingScore }}</span>
            <svg
              class="size-4 text-star"
              viewBox="0 0 24 24"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                d="M12 17.27 18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"
              />
            </svg>
          </div>
          <span class="text-xs text-subtle">
            ({{ ratingCount.toLocaleString() }})
          </span>
        </div>

        <div class="flex items-center gap-2 text-sm text-subtle">
          <svg
            class="size-4 text-primary/80"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            aria-hidden="true"
          >
            <path
              d="M12 7v5l3 2"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
              stroke="currentColor"
              stroke-width="2"
            />
          </svg>
          <span>{{ timeMinutes }} mins</span>
        </div>
      </div>

      <button
        type="button"
        class="relative z-20 mt-auto w-full rounded-xl border px-4 py-2.5 text-sm font-semibold shadow-sm transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/30"
        :class="
          saved
            ? 'bg-cta border-cta text-inverted hover:bg-primary'
            : 'bg-canvas border-primary/25 text-body hover:bg-accent-subtle/40'
        "
        :disabled="toggling"
        @click="toggleBookmark"
      >
        <span class="inline-flex items-center justify-center gap-2">
          <span>{{ saved ? "Saved" : "Save recipe" }}</span>
          <svg
            v-if="saved"
            class="size-4"
            viewBox="0 0 24 24"
            fill="currentColor"
            aria-hidden="true"
          >
            <path d="M5 5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16l-7-3.5L5 21V5Z" />
          </svg>
          <svg
            v-else
            class="size-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <path d="M5 5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16l-7-3.5L5 21V5Z" />
          </svg>
        </span>
      </button>
    </div>

    <div
      class="pointer-events-none absolute inset-0 -z-10 bg-linear-to-br from-accent/10 via-accent-subtle/10 to-primary-subtle/10 opacity-0 transition group-hover:opacity-100"
      aria-hidden="true"
    />
  </article>
</template>
