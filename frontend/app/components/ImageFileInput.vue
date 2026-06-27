
<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    default: "Add image",
  },
  inputId: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:modelValue"]);

const { uploadRecipeImage } = useFileUpload();
const error = ref("");
const loading = ref(false);
const cacheBuster = ref(Date.now());

const displayUrl = computed(() => {
  if (!props.modelValue) return "";
  const separator = props.modelValue.includes("?") ? "&" : "?";
  return `${props.modelValue}${separator}t=${cacheBuster.value}`;
});

async function onFileChange(event) {
  const file = event.target.files?.[0];
  event.target.value = "";
  if (!file) return;

  error.value = "";
  loading.value = true;

  try {
    const serverUrl = await uploadRecipeImage(file);
    cacheBuster.value = Date.now();
    emit("update:modelValue", serverUrl);
  } catch (err) {
    error.value = err?.message || "Could not upload image.";
  } finally {
    loading.value = false;
  }
}

function removeImage() {
  emit("update:modelValue", "");
}
</script>

<template>
  <div class="space-y-2">
    <div v-if="modelValue" class="relative inline-block">
      <img
        :src="displayUrl"
        alt="Uploaded preview"
        class="h-24 w-24 rounded-xl border border-primary/20 object-cover"
      />
      <button
        type="button"
        class="absolute -right-2 -top-2 flex size-6 items-center justify-center rounded-full bg-danger text-xs font-bold text-inverted shadow"
        aria-label="Remove image"
        @click="removeImage"
      >
        ×
      </button>
    </div>

    <label
      :for="inputId"
      class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-dashed border-primary/35 bg-white/60 px-3 py-2 text-xs font-semibold text-primary transition hover:border-primary hover:bg-accent-subtle/30"
    >
      <span>{{ loading ? "Uploading..." : label }}</span>
      <input
        :id="inputId"
        type="file"
        accept="image/*"
        class="sr-only"
        :disabled="loading"
        @change="onFileChange"
      />
    </label>

    <p v-if="error" class="text-xs text-danger">{{ error }}</p>
  </div>
</template>
