<template>
    <div class="flex min-h-screen flex-col bg-linear-to-b pure to-accent/20">
        <AppHeader />
        <AppToast />
        <main class="flex-1">
            <slot />
        </main>
        <AppFooter />
    </div>
</template>
