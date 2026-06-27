<!-- Home | Search.vue -->
<script setup>
const props = defineProps({
  placeholder: String,
  initialQuery: String,
})

const emit = defineEmits(["search"])

const route = useRoute()
const router = useRouter()

const qFromRoute = computed(() => {
  const raw = route.query.q
  if (typeof raw === "string") return raw
  if (Array.isArray(raw)) return raw[0] ?? ""
  return ""
})

const query = ref(props.initialQuery ?? qFromRoute.value)

watch(
  () => qFromRoute.value,
  (next) => {
    if (props.initialQuery) return
    query.value = next
  },
)

function submit() {
  const q = query.value.trim()
  if (!q) return
  emit("search", q)
  router.push(`/search?q=${encodeURIComponent(q)}`)
}
</script>

<template>
  <form
    class="flex items-center gap-3 rounded-2xl border border-primary/20 bg-pure px-3 py-2 shadow-sm backdrop-blur-md"
    @submit.prevent="submit"
  >
    <svg
      class="size-5 text-primary/80"
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <path
        d="M10.5 18a7.5 7.5 0 1 1 0-15 7.5 7.5 0 0 1 0 15Z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <path
        d="M16.2 16.2 21 21"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>

    <input
      v-model="query"
      :placeholder="placeholder ?? 'Search recipes...'"
      class="w-full bg-transparent text-sm text-body placeholder:text-subtle/70 outline-none"
      type="search"
      name="q"
      autocomplete="off"
      aria-label="Search recipes"
    />

    <button type="submit" class="btn-primary">
      Search
    </button>
  </form>
</template>