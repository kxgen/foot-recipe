<template>
    <div>
        <section class="py-14">
            <div class="page-container">
                <h1 class="text-3xl font-semibold text-heading">
                    Welcome back, {{ username || "Chef" }}!
                </h1>
                <p class="mt-2 text-sm text-heading-muted">
                    Quick access to your recipes, bookmarks, purchases, and
                    profile settings.
                </p>
            </div>
        </section>

        <section class="py-14">
            <div class="page-container">
                <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
                    <NuxtLink
                        to="/recipes/create"
                        class="card-surface border-2 border-primary-dark/25 bg-accent-subtle/20 p-6 hover:bg-accent-subtle/40 transition group"
                    >
                        <div
                            class="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary-dark/10 text-primary group-hover:scale-105 transition"
                        >
                            <PlusCircle class="size-7" />
                        </div>
                        <div class="text-sm font-semibold text-primary">
                            Create recipe
                        </div>
                        <div class="mt-1 text-xs text-heading-muted">
                            Add a new recipe with steps, ingredients, and
                            images.
                        </div>
                    </NuxtLink>
                    <NuxtLink
                        to="/me/recipes"
                        class="card-surface p-6 hover:bg-accent-subtle/25 transition group"
                    >
                        <div
                            class="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary-subtle/20 text-heading group-hover:scale-105 transition"
                        >
                            <ChefHat class="size-7" />
                        </div>
                        <div class="text-sm font-semibold text-heading">
                            My recipes
                        </div>
                        <div class="mt-1 text-xs text-heading-muted">
                            Create, edit, and manage your recipes.
                        </div>
                    </NuxtLink>
                    <NuxtLink
                        to="/me/bookmarks"
                        class="card-surface p-6 hover:bg-accent-subtle/25 transition group"
                    >
                        <div
                            class="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary-subtle/20 text-heading group-hover:scale-105 transition"
                        >
                            <Bookmark class="size-7" />
                        </div>
                        <div class="text-sm font-semibold text-heading">
                            Bookmarks
                        </div>
                        <div class="mt-1 text-xs text-heading-muted">
                            Saved recipes you want to cook later.
                        </div>
                    </NuxtLink>
                    <NuxtLink
                        to="/me/purchases"
                        class="card-surface p-6 hover:bg-accent-subtle/25 transition group"
                    >
                        <div
                            class="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary-subtle/20 text-heading group-hover:scale-105 transition"
                        >
                            <CreditCard class="size-7" />
                        </div>
                        <div class="text-sm font-semibold text-heading">
                            Purchases
                        </div>
                        <div class="mt-1 text-xs text-heading-muted">
                            Recipes you’ve purchased.
                        </div>
                    </NuxtLink>
                    <NuxtLink
                        to="/me/profile"
                        class="card-surface p-6 hover:bg-accent-subtle/25 transition group"
                    >
                        <div
                            class="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary-subtle/20 text-heading group-hover:scale-105 transition"
                        >
                            <User class="size-7" />
                        </div>
                        <div class="text-sm font-semibold text-heading">
                            Profile
                        </div>
                        <div class="mt-1 text-xs text-heading-muted">
                            Update avatar, name, password, and more.
                        </div>
                    </NuxtLink>
                </div>
            </div>
        </section>
    </div>
</template>

<script setup>
import { PlusCircle, ChefHat, Bookmark, CreditCard, User } from "@lucide/vue";
import { getUserIdFromToken } from "~/utils/jwt";
import { GET_USER_PROFILE } from "~/graphql/queries.js"

definePageMeta({
    middleware: ["auth"],
});
useHead({
    title: "Me — Share Plate",
});


const { $apollo: apollo } = useNuxtApp();
const { token } = useAuth();
const userId = computed(() => getUserIdFromToken(token.value));
const username = ref("");

onMounted(async () => {
    if (!userId.value) return;
    try {
        const { data } = await apollo.query({
            query: GET_USER_PROFILE,
            variables: { id: userId.value },
        });
        if (data?.users_by_pk?.username) {
            username.value = data.users_by_pk.username;
        }
    } catch (err) {
        console.error("Failed to fetch user profile", err);
    }
});
</script>
