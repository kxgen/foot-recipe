<!-- xo -->
<!-- Default Layout -->
<script setup>
const year = new Date().getFullYear();

// object -> title with list of links
const footerLinks = [
    {
        title: "Explore",
        links: [
            { label: "Recipes", to: "/recipes" },
            { label: "Categories", to: "/categories" },
            { label: "Creators", to: "/creators" },
            { label: "Search", to: "/search" },
        ],
    },
    {
        title: "Account",
        links: [
            { label: "Log in", to: "/auth/login" },
            { label: "Sign up", to: "/auth/register" },
            { label: "My recipes", to: "/me/recipes" },
            { label: "Bookmarks", to: "/me/bookmarks" },
        ],
    },
    {
        title: "About",
        links: [
            { label: "About us", to: "/about" },
            { label: "Create recipe", to: "/recipes/create" },
        ],
    },
];
</script>

<template>
    <footer
        class="mt-auto border-t border-primary/15 bg-accent/10"
    >
        <div class="page-container py-12 lg:py-14">
            <div class="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
                <div class="sm:col-span-2 lg:col-span-1">
                    <NuxtLink to="/" class="inline-flex items-center gap-2">
                        <span
                            class="font-brand text-2xl font-normal text-heading leading-none"
                            >Share Plate</span
                        >
                    </NuxtLink>
                    <p
                        class="mt-3 max-w-xs text-sm leading-relaxed text-heading-muted"
                    >
                        Discover, share, and savor recipes from home cooks
                        around the world.
                    </p>
                </div>

                <div v-for="group in footerLinks" :key="group.title">
                    <h3
                        class="text-sm font-semibold uppercase tracking-wide text-heading"
                    >
                        {{ group.title }}
                    </h3>
                    <ul class="mt-3 space-y-2">
                        <li v-for="link in group.links" :key="link.to">
                            <NuxtLink
                                :to="link.to"
                                class="text-sm text-heading-muted hover:text-primary"
                            >
                                {{ link.label }}
                            </NuxtLink>
                        </li>
                    </ul>
                </div>
            </div>

            <div
                class="mt-10 flex flex-col items-center justify-between gap-4 border-t border-primary/15 pt-8 sm:flex-row"
            >
                <p class="text-sm text-heading-muted">
                    &copy; {{ year }} Share Plate. All rights reserved.
                </p>
                <p class="text-xs text-heading-muted/80">
                    Made with fresh ingredients &amp; love.
                </p>
            </div>
        </div>
    </footer>
</template>
