<template>
  <Transition name="fade">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
    >
      <div
        class="w-full max-w-sm rounded-3xl border border-primary/15 bg-white p-6 shadow-2xl transition-all"
      >
        <h3 class="text-lg font-bold text-heading">
          {{ title }}
        </h3>
        <p class="mt-2 text-sm text-subtle leading-relaxed">
          {{ message }}
        </p>

        <div class="mt-6 flex justify-end gap-3">
          <button
            type="button"
            class="rounded-xl bg-surface-muted px-4 py-2.5 text-xs font-semibold text-heading transition hover:bg-primary/10"
            @click="cancel"
          >
            {{ cancelText }}
          </button>
          <button
            type="button"
            class="rounded-xl bg-danger px-4 py-2.5 text-xs font-semibold text-white transition hover:bg-danger/90"
            @click="confirm"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Are you sure?",
  },
  message: {
    type: String,
    default: "This action cannot be undone.",
  },
  confirmText: {
    type: String,
    default: "Confirm",
  },
  cancelText: {
    type: String,
    default: "Cancel",
  },
});

const emit = defineEmits(["confirm", "close"]);

function confirm() {
  emit("confirm");
}

function cancel() {
  emit("close");
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
