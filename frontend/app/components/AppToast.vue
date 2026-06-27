<!-- Toast -->
<!-- Default Layout -->
<script setup>
import { useNotifications } from '~/composables/useNotifications';

const { notifications, removeNotification } = useNotifications();

const typeClasses = {
  success: 'bg-success-bg border-success/20 text-success',
  error: 'bg-danger-bg border-danger/20 text-danger',
  info: 'bg-info-bg border-info/20 text-info',
  warning: 'bg-warning-bg border-warning/20 text-warning',
};

const iconPaths = {
  success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
  error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
  info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  warning: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
};

</script>

<template>
  <div
    class="pointer-events-none fixed inset-x-0 top-0 z-[100] flex flex-col items-center gap-2 p-4 sm:items-end"
  >
    <TransitionGroup
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="translate-y-[-100%] opacity-0 sm:translate-x-[100%] sm:translate-y-0"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="pointer-events-auto flex w-full max-w-sm items-start gap-3 rounded-xl border p-4 shadow-lg"
        :class="typeClasses[notification.type]"
      >
        <div class="shrink-0 pt-0.5">
          <svg
            class="size-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path stroke-linecap="round" stroke-linejoin="round" :d="iconPaths[notification.type]" />
          </svg>
        </div>

        <div class="flex-1 text-sm font-medium leading-5">
          {{ notification.message }}
        </div>

        <button
          type="button"
          class="shrink-0 rounded-lg p-0.5 opacity-60 transition hover:opacity-100 focus:outline-none"
          @click="removeNotification(notification.id)"
        >
          <span class="sr-only">Close</span>
          <svg class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
