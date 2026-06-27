<script setup>
const sortBy = defineModel("sortBy", { default: "created_at" });
const sortOrder = defineModel("sortOrder", { default: "desc" });

const emit = defineEmits(["change"]);

const handleChange = () => {
    emit("change");
};
</script>

<template>
    <div class="flex gap-2">
        <div class="flex flex-col">
            <select
                id="sort-by"
                v-model="sortBy"
                @change="handleChange"
                class="rounded-xl border border-primary/20 bg-white/60 px-3 py-2 text-sm text-heading focus:border-primary focus:outline-none focus:ring-2 focus:ring-accent/40"
            >
                <option value="created_at">Newest</option>
                <option value="avg_rating">Rating</option>
                <option value="total_time_minutes">Prep Time</option>
                <option value="price">Price</option>
            </select>
        </div>

        <div class="flex flex-col">
            <select
                id="sort-order"
                v-model="sortOrder"
                @change="handleChange"
                class="rounded-xl border border-primary/20 bg-white/60 px-3 py-2 text-sm text-heading focus:border-primary focus:outline-none focus:ring-2 focus:ring-accent/40"
            >
                <option value="desc">Desc</option>
                <option value="asc">Asc</option>
            </select>
        </div>
    </div>
</template>
