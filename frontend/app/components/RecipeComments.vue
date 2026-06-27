<!-- Recipes/:[id].vue -->
<script setup>
import { GET_RECIPE_COMMENTS } from "~/graphql/queries.js"
import { UPSERT_RECIPE_RATING, INSERT_RECIPE_COMMENT, UPDATE_RECIPE_COMMENT, DELETE_RECIPE_COMMENT } from "~/graphql/mutations.js"
import { getUserIdFromToken } from "~/utils/jwt";
import { useNotifications } from "~/composables/useNotifications";

const props = defineProps({
    recipeId: {
        type: Number,
        required: true,
    },
    avgRating: {
        type: [Number, String],
        default: 0,
    },
    ratingCount: {
        type: Number,
        default: 0,
    },
});

const { $apollo: apollo } = useNuxtApp();
const { token, isLoggedIn } = useAuth();
const { notifySuccess, notifyError } = useNotifications();

const emit = defineEmits(["refresh"]);

const comments = ref([]);
const ratingsByUserId = ref({});
const loading = ref(true);
const error = ref("");
const body = ref("");
const rating = ref(5);
const submitting = ref(false);

const editingCommentId = ref(null);
const editingBody = ref("");
const updatingComment = ref(false);

const currentUserId = computed(() => getUserIdFromToken(token.value));

const ratingDistribution = computed(() => {
    const dist = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
    Object.values(ratingsByUserId.value).forEach((score) => {
        if (dist[score] !== undefined) dist[score]++;
    });
    return dist;
});

const totalRatings = computed(() => Object.keys(ratingsByUserId.value).length);

function mapRatings(ratings) {
    const map = {};
    for (const entry of ratings ?? []) {
        map[entry.user_id] = entry.score;
    }
    return map;
}

const isInitialLoad = ref(true);

async function loadComments() {
    loading.value = true;
    error.value = "";

    try {
        const { data } = await apollo.query({
            query: GET_RECIPE_COMMENTS,
            variables: { recipeId: props.recipeId },
            fetchPolicy: "network-only",
        });
        comments.value = data?.recipe_comments ?? [];
        ratingsByUserId.value = mapRatings(data?.recipe_ratings);

        if (currentUserId.value && ratingsByUserId.value[currentUserId.value]) {
            rating.value = ratingsByUserId.value[currentUserId.value];
        }
        
        // Wait a tick for the watcher to acknowledge the data assignment silently
        await nextTick();
        isInitialLoad.value = false;
        
    } catch (err) {
        error.value = err?.message || "Could not load comments.";
    } finally {
        loading.value = false;
    }
}

async function submitRating(newScore) {
    if (!isLoggedIn.value || !currentUserId.value) {
        await navigateTo("/auth/login");
        return;
    }

    try {
        await apollo.mutate({
            mutation: UPSERT_RECIPE_RATING,
            variables: {
                recipeId: props.recipeId,
                score: Number(newScore),
            },
        });

        ratingsByUserId.value = {
            ...ratingsByUserId.value,
            [currentUserId.value]: Number(newScore),
        };

        notifySuccess("Rating updated.");
        emit("refresh");
    } catch (err) {
        notifyError(err?.message || "Could not update rating.");
    }
}

async function submitComment() {
    const trimmed = body.value.trim();
    if (!trimmed) return;

    if (!isLoggedIn.value || !currentUserId.value) {
        await navigateTo("/auth/login");
        return;
    }

    submitting.value = true;
    error.value = "";

    try {
        const { data } = await apollo.mutate({
            mutation: INSERT_RECIPE_COMMENT,
            variables: {
                recipeId: props.recipeId,
                body: trimmed,
            },
        });

        const created = data?.insert_recipe_comments_one;
        if (created) {
            comments.value = [created, ...comments.value];
            body.value = "";
            notifySuccess("Comment posted successfully.");
        } else {
            await loadComments();
        }
    } catch (err) {
        notifyError(err?.message || "Could not post comment.");
    } finally {
        submitting.value = false;
    }
}

watch(rating, (newVal, oldVal) => {
    if (isInitialLoad.value || loading.value) return;

    if (newVal !== oldVal && newVal > 0) {
        submitRating(newVal);
    }
});

function startEditing(comment) {
    editingCommentId.value = comment.id;
    editingBody.value = comment.body;
}

function cancelEditing() {
    editingCommentId.value = null;
    editingBody.value = "";
}

async function updateComment() {
    const trimmed = editingBody.value.trim();
    if (!trimmed) return;

    updatingComment.value = true;

    try {
        await apollo.mutate({
            mutation: UPDATE_RECIPE_COMMENT,
            variables: {
                id: editingCommentId.value,
                body: trimmed,
            },
        });

        // Replaces the item by mapping to a new object copy
        comments.value = comments.value.map((c) =>
            c.id === editingCommentId.value ? { ...c, body: trimmed } : c,
        );

        editingCommentId.value = null;
        editingBody.value = "";
        notifySuccess("Comment updated.");
    } catch (err) {
        notifyError(err?.message || "Could not update comment.");
    } finally {
        updatingComment.value = false;
    }
}

async function deleteComment(id) {
    if (!confirm("Are you sure you want to delete this comment?")) return;

    try {
        await apollo.mutate({
            mutation: DELETE_RECIPE_COMMENT,
            variables: { id },
        });

        comments.value = comments.value.filter((c) => c.id !== id);
        notifySuccess("Comment deleted.");
    } catch (err) {
        notifyError(err?.message || "Could not delete comment.");
    }
}

function formatDate(value) {
    if (!value) return "";
    return new Date(value).toLocaleString();
}

function ratingForUser(userId) {
    return ratingsByUserId.value[userId] ?? 0;
}

watch(
    () => props.recipeId,
    () => {
        if (props.recipeId) loadComments();
    },
    { immediate: true },
);
</script>

<template>
    <section class="mt-16 border-t border-primary/10 pt-10">
        <div class="flex flex-wrap items-end justify-between gap-4">
            <div>
                <h2 class="text-2xl font-bold text-heading">Community Reviews</h2>
                <p class="mt-1 text-sm text-subtle">
                    Share your thoughts and rate this recipe
                </p>
            </div>
            <div
                v-if="ratingCount > 0"
                class="flex items-center gap-3 rounded-2xl bg-primary-subtle/10 px-4 py-2"
            >
                <span class="text-2xl font-bold text-primary"
                    >⭐ {{ Number(avgRating).toFixed(1) }}</span
                >
                <div class="h-8 w-px bg-primary/20"></div>
                <span class="text-sm font-medium text-subtle"
                    >{{ ratingCount }} reviews</span
                >
            </div>
        </div>

        <!-- Rating Summary & Form Grid -->
        <div class="mt-8 grid gap-8 lg:grid-cols-12">
            <!-- Left: Form -->
            <div class="lg:col-span-7">
                <form class="card-surface p-6" @submit.prevent="submitComment">
                    <h3 class="mb-4 text-lg font-semibold text-heading">
                        Write a Review
                    </h3>
                    <div class="space-y-5">
                        <div>
                            <label
                                class="mb-2 block text-sm font-medium text-heading"
                                >Your Rating</label
                            >
                            <StarRating
                                v-model="rating"
                                :readonly="!isLoggedIn || submitting"
                            />
                        </div>

                        <div>
                            <label
                                class="mb-2 block text-sm font-medium text-heading"
                                >Your Comment</label
                            >
                            <textarea
                                v-model="body"
                                rows="4"
                                class="input-field min-h-30 resize-y text-sm"
                                :placeholder="
                                    isLoggedIn
                                        ? 'What did you think of this recipe? Any tips or changes?'
                                        : 'Please log in to leave a review...'
                                "
                                :disabled="!isLoggedIn || submitting"
                            />
                        </div>

                        <div class="flex items-center justify-between">
                            <button
                                v-if="isLoggedIn"
                                type="submit"
                                class="btn-primary px-8"
                                :disabled="submitting || !body.trim()"
                            >
                                {{
                                    submitting ? "Posting..." : "Submit Review"
                                }}
                            </button>
                            <NuxtLink
                                v-else
                                to="/auth/login"
                                class="btn-primary px-8"
                            >
                                Log in to Review
                            </NuxtLink>
                        </div>
                    </div>
                </form>
            </div>

            <!-- Right: Stats -->
            <div class="lg:col-span-5">
                <div class="card-surface p-6">
                    <h3 class="mb-4 text-lg font-semibold text-heading">
                        Review Summary
                    </h3>
                    <div class="space-y-3">
                        <div
                            v-for="star in [5, 4, 3, 2, 1]"
                            :key="star"
                            class="flex items-center gap-3"
                        >
                            <span
                                class="w-12 text-sm font-medium text-subtle"
                                >{{ star }} stars</span
                            >
                            <div
                                class="relative h-2 flex-1 overflow-hidden rounded-full bg-primary/10"
                            >
                                <div
                                    class="absolute inset-y-0 left-0 bg-star transition-all duration-500"
                                    :style="{
                                        width:
                                            totalRatings > 0
                                                ? `${(ratingDistribution[star] / totalRatings) * 100}%`
                                                : '0%',
                                    }"
                                ></div>
                            </div>
                            <span
                                class="w-8 text-right text-xs font-medium text-subtle"
                            >
                                {{ ratingDistribution[star] }}
                            </span>
                        </div>
                    </div>
                    <p class="mt-6 text-center text-xs text-subtle">
                        Showing ratings from {{ totalRatings }} community
                        members
                    </p>
                </div>
            </div>
        </div>

        <!-- Comments List -->
        <div class="mt-12">
            <h3 class="mb-6 text-xl font-bold text-heading">
                Latest Comments
                <span
                    v-if="comments.length"
                    class="ml-2 text-sm font-normal text-subtle"
                >
                    ({{ comments.length }})
                </span>
            </h3>

            <div
                v-if="loading"
                class="flex flex-col items-center py-12 text-subtle"
            >
                <InlineLoader />
                <p class="mt-4 text-sm">Loading reviews...</p>
            </div>

            <p
                v-else-if="error"
                class="rounded-2xl bg-danger-bg p-4 text-center text-sm text-danger"
            >
                {{ error }}
            </p>

            <div v-else-if="comments.length" class="space-y-6">
                <div
                    v-for="comment in comments"
                    :key="comment.id"
                    class="card-surface p-6 transition-shadow hover:shadow-md"
                >
                    <div class="flex items-start gap-4">
                        <div
                            class="flex size-12 shrink-0 items-center justify-center rounded-2xl bg-primary-subtle/20 text-lg font-bold text-primary"
                        >
                            {{
                                comment.user?.username
                                    ?.charAt(0)
                                    ?.toUpperCase() || "?"
                            }}
                        </div>
                        <div class="min-w-0 flex-1">
                            <div
                                class="flex flex-wrap items-center justify-between gap-2"
                            >
                                <div class="flex items-center gap-3">
                                    <span class="font-bold text-heading">
                                        {{
                                            comment.user?.username ||
                                            "A community member"
                                        }}
                                    </span>
                                    <StarRating
                                        v-if="ratingForUser(comment.user?.id)"
                                        :model-value="
                                            ratingForUser(comment.user?.id)
                                        "
                                        readonly
                                        size="sm"
                                    />
                                </div>
                                <div class="flex items-center gap-3">
                                    <time class="text-xs text-subtle">
                                        {{ formatDate(comment.created_at) }}
                                    </time>

                                    <!-- Edit/Delete Actions -->
                                    <div
                                        v-if="
                                            currentUserId ===
                                                comment.user?.id &&
                                            editingCommentId !== comment.id
                                        "
                                        class="flex items-center gap-2 border-l border-primary/15 pl-3"
                                    >
                                        <button
                                            class="text-xs font-semibold text-primary hover:underline"
                                            @click="startEditing(comment)"
                                        >
                                            Edit
                                        </button>
                                        <button
                                            class="text-xs font-semibold text-primary hover:underline"
                                            @click="deleteComment(comment.id)"
                                        >
                                            Delete
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <!-- Comment Content / Edit Form -->
                            <div
                                v-if="editingCommentId === comment.id"
                                class="mt-3 space-y-3"
                            >
                                <textarea
                                    v-model="editingBody"
                                    rows="3"
                                    class="input-field min-h-20 resize-y text-sm"
                                    :disabled="updatingComment"
                                />
                                <div class="flex items-center gap-2">
                                    <button
                                        class="btn-primary py-1.5 px-4 text-xs"
                                        :disabled="
                                            updatingComment ||
                                            !editingBody.trim()
                                        "
                                        @click="updateComment"
                                    >
                                        {{
                                            updatingComment
                                                ? "Saving..."
                                                : "Save"
                                        }}
                                    </button>
                                    <button
                                        class="btn-secondary py-1.5 px-4 text-xs"
                                        :disabled="updatingComment"
                                        @click="cancelEditing"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </div>
                            <div
                                v-else
                                class="mt-3 text-sm leading-relaxed text-subtle whitespace-pre-wrap"
                            >
                                {{ comment.body }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div
                v-else
                class="rounded-3xl border-2 border-dashed border-primary/20 py-16 text-center"
            >
                <div class="text-4xl mb-4">💬</div>
                <p class="text-subtle">
                    No reviews yet. Be the first to share your experience!
                </p>
            </div>
        </div>
    </section>
</template>
