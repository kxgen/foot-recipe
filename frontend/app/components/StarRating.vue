<!--  -->
<script setup>
const props = defineProps({
  modelValue: {
    type: Number,
    default: 0,
  },
  readonly: {
    type: Boolean,
    default: false,
  },
  size: {
    type: String,
    default: "md",
  },
});

const emit = defineEmits(["update:modelValue"]);

const hoverScore = ref(0);

const displayScore = computed(() => {
  if (props.readonly) return props.modelValue;
  return hoverScore.value || props.modelValue;
});

const sizeClass = computed(() =>
  props.size === "sm" ? "text-base gap-0.5" : "text-xl gap-1",
);

function setScore(score) {
  if (props.readonly) return;
  emit("update:modelValue", score);
}

function onMouseLeave() {
  hoverScore.value = 0;
}
</script>

<template>
  <div
    class="inline-flex items-center"
    :class="sizeClass"
    role="group"
    :aria-label="readonly ? `${modelValue} out of 5 stars` : 'Rate from 1 to 5 stars'"
    @mouseleave="onMouseLeave"
  >
    <button
      v-for="star in 5"
      :key="star"
      type="button"
      class="leading-none transition"
      :class="[
        readonly ? 'cursor-default' : 'cursor-pointer hover:scale-110',
        star <= displayScore ? 'text-star' : 'text-primary/30',
      ]"
      :disabled="readonly"
      :aria-label="`${star} star${star === 1 ? '' : 's'}`"
      @click="setScore(star)"
      @mouseenter="!readonly && (hoverScore = star)"
    >
      ★
    </button>
  </div>
</template>
