<!-- Header -->
<!-- Default Layout -->

<script setup>
const route = useRoute();
const { isLoggedIn, logout } = useAuth();

const authLinks = computed(() => {
    if (!isLoggedIn.value) {
        return [
            { label: "Log in", to: "/auth/login", class: "btn-secondary" },
            { label: "Sign up", to: "/auth/register", class: "btn-primary" }
        ];
    } else {
        return [
            { label: "Me", to: "/me", class: "btn-primary" },
            { label: "Logout", action: logout, class: "btn-secondary" }
        ];
    }
});

const navLinks = [
    { label: "Home", to: "/" },
    { label: "Recipes", to: "/recipes" },
    { label: "Categories", to: "/categories" },
    { label: "Creators", to: "/creators" },
];

const hamburgerMenuOpen = ref(false);

function isActive(path) {
    if (path === "/")
        return route.path === "/";
    return route.path.startsWith(path);
}

watch(
    () => route.path,
    () => {
        hamburgerMenuOpen.value = false;
    },
);
</script>

<template>
    <header class="sticky top-0 z-50 border-b border-primary/15 bg-layer/70 backdrop-blur-md">
        <div class="page-container flex h-16 items-center justify-between gap-4">
            <NuxtLink to="/" class="group flex items-center shrink-0">
                <span
                    class="font-brand text-4xl leading-none text-primary  transition-all duration-300 drop-shadow-sm group-hover:-rotate-2">Share</span>
                <span
                    class="font-brand text-4xl leading-none text-accent  transition-all duration-300 -ml-1 group-hover:-translate-y-0.5">Plate</span>
            </NuxtLink>

            <nav class="hidden items-center gap-1 md:flex" aria-label="Main">
                <NuxtLink v-for="link in navLinks" :key="link.to" :to="link.to"
                    class="rounded-lg px-3 py-2 text-sm font-medium transition" :class="isActive(link.to)
                        ? 'bg-accent-subtle text-heading'
                        : 'text-subtle hover:bg-layer hover:text-body'
                        ">
                    {{ link.label }}
                </NuxtLink>
            </nav>

            <div class="hidden items-center gap-2 md:flex">


                <NuxtLink v-for="btn in authLinks" :key="btn.label" :to="btn.to"
                    :class="[btn.class, 'px-4 py-2 text-sm']" @click="btn.action ? btn.action() : null">
                    {{ btn.label }}
                </NuxtLink>
            </div>

            <button type="button"
                class="inline-flex size-10 items-center justify-center rounded-lg border border-primary/20 text-body md:hidden"
                :aria-expanded="hamburgerMenuOpen" aria-controls="mobile-nav" aria-label="Toggle menu"
                @click="hamburgerMenuOpen = !hamburgerMenuOpen">
                <svg v-if="!hamburgerMenuOpen" xmlns="http://www.w3.org/2000/svg" class="size-5" fill="none"
                    viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" d="M4 7h16M4 12h16M4 17h16" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="size-5" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" d="M6 6l12 12M18 6L6 18" />
                </svg>
            </button>
        </div>

        <div v-show="hamburgerMenuOpen" id="mobile-nav" class="border-t border-primary/15 bg-canvas md:hidden">
            <nav class="page-container flex flex-col gap-1 py-4" aria-label="Mobile">
                <NuxtLink v-for="link in navLinks" :key="link.to" :to="link.to"
                    class="rounded-lg px-3 py-2.5 text-sm font-medium" :class="isActive(link.to)
                        ? 'bg-accent-subtle text-heading'
                        : 'text-subtle hover:bg-layer'
                        ">
                    {{ link.label }}
                </NuxtLink>

                <div class="my-2 border-t border-primary/10"></div>

                <div class="flex flex-col gap-2 px-3">


                    <NuxtLink v-for="btn in authLinks" :key="btn.label" :to="btn.to"
                        :class="[btn.class, 'w-full py-2.5 text-center text-sm']"
                        @click="btn.action ? btn.action() : null">
                        {{ btn.label }}
                    </NuxtLink>
                </div>
            </nav>
        </div>
    </header>
</template>