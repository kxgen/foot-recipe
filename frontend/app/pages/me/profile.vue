<script setup>
import { gql } from "@apollo/client/core";
import { getUserIdFromToken } from "~/utils/jwt";
import { useNotifications } from "~/composables/useNotifications";
import { GET_USER_PROFILE } from "~/graphql/queries.js"
import { UPDATE_USER_PROFILE, UPDATE_PASSWORD } from "~/graphql/mutations.js"

const { $apollo: apollo } = useNuxtApp();
const query = (options) => apollo.query({ fetchPolicy: "network-only", ...options });
const mutate = (options) => apollo.mutate(options);
const { notifySuccess, notifyError } = useNotifications();
const { uploadAvatar } = useFileUpload();

const { token } = useAuth();
const userId = computed(() => getUserIdFromToken(token.value));

const form = reactive({
    username: "",
    slug: "",
    email: "",
    password: "",
    bio: "",
});

const loadError = ref("");
const fileInput = ref(null);
const avatarPreview = ref("");

const editing = reactive({
    username: false,
    password: false,
    bio: false,
});

function toggleEdit(field) {
    editing[field] = !editing[field];
}

async function saveField(field) {
    if (!userId.value) {
        await navigateTo("/auth/login");
        return;
    }

    const set = { [field]: form[field] };
    if (field === "password" && !form.password) return;

    try {
        if (field === "password") {
            await mutate({
                mutation: UPDATE_PASSWORD,
                variables: {
                    newPassword: form.password,
                },
            });
        } else {
            await mutate({
                mutation: UPDATE_USER_PROFILE,
                variables: {
                    id: userId.value,
                    set,
                },
            });
        }

        editing[field] = false;
        if (field === "password") form.password = "";
        notifySuccess(`${field.charAt(0).toUpperCase() + field.slice(1)} updated successfully.`);
    } catch (error) {
        notifyError(error?.message || `Could not save ${field}.`);
        console.error(error);
    }
}

function triggerAvatarSelect() {
    fileInput.value?.click();
}

async function onAvatarFileChange(event) {
    const file = event.target.files[0];
    if (!file) return;

    try {
        // Preview
        const reader = new FileReader();
        reader.onload = (e) => {
            avatarPreview.value = e.target.result;
        };
        reader.readAsDataURL(file);

        const imageUrl = await uploadAvatar(file, { 
            userId: userId.value,
            slug: form.slug
        });
        
        avatarPreview.value = imageUrl;
        notifySuccess("Profile picture updated successfully.");
    } catch (error) {
        notifyError(error?.message || "Could not update profile picture.");
        console.error(error);
    }
}

const initials = computed(() => {
    return form.username?.charAt(0)?.toUpperCase() || "U";
});

onMounted(loadProfile);
watch(userId, loadProfile);

async function loadProfile() {
    if (!userId.value) {
        await navigateTo("/auth/login");
        return;
    }

    loadError.value = "";

    try {
        const { data } = await query({
            query: GET_USER_PROFILE,
            variables: { id: userId.value },
        });
        const user = data?.users_by_pk;
        if (!user) return;
        form.username = user.username ?? "";
        form.slug = user.slug ?? "";
        form.email = user.email ?? "";
        form.bio = user.bio ?? "";
        avatarPreview.value = user.avatar_url ?? "";
    } catch (error) {
        loadError.value = error?.message || "Failed to load profile.";
        console.error(error);
    }
}
</script>

<template>
    <div class="py-10">
        <div class="page-container max-w-4xl">
            <!-- Header -->
            <section
                class="relative overflow-hidden rounded-3xl border border-primary/15 bg-surface px-6 py-8 sm:px-8"
            >
                <div
                    class="pointer-events-none absolute inset-0 -z-10 bg-gradient-to-br from-lime/20 via-transparent to-sage/10"
                    aria-hidden="true"
                />

                <div class="flex items-center justify-between gap-5">
                    <div class="flex items-center gap-5">
                        <div
                            class="group relative flex size-24 items-center justify-center overflow-hidden rounded-3xl border border-primary/15 bg-white text-3xl font-semibold text-primary cursor-pointer"
                            @click="triggerAvatarSelect"
                        >
                            <img
                                v-if="avatarPreview"
                                :src="avatarPreview"
                                alt="Avatar"
                                class="h-full w-full object-cover transition-opacity group-hover:opacity-75"
                            />

                            <span v-else>
                                {{ initials }}
                            </span>

                            <!-- Overlay on hover -->
                            <div class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition-opacity group-hover:opacity-100">
                                <svg xmlns="http://www.w3.org/2000/svg" class="size-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                                </svg>
                            </div>
                        </div>

                        <div>
                            <h1 class="text-2xl font-semibold text-heading">
                                Your Profile
                            </h1>

                            <p class="mt-1 text-sm text-heading-muted">
                                Manage your account information and public profile.
                            </p>
                        </div>
                    </div>

                    <!-- Email on the rightmost side -->
                    <div class="hidden sm:block text-right">
                        <p class="text-sm font-semibold text-heading-muted">Email - {{ form.email }}</p>
                    </div>
                </div>

                <!-- Hidden input -->
                <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="onAvatarFileChange" />
            </section>

            <p
                v-if="updateError"
                class="mt-4 error-message"
            >
                {{ loadError }}
            </p>

            <!-- Form -->
            <section class="mt-8 space-y-5">
                <!-- Username -->
                <div class="card-surface p-5">
                    <div class="mb-3 flex items-center justify-between gap-4">
                        <div>
                            <h2 class="text-sm font-semibold text-heading">
                                Username
                            </h2>

                            <p class="text-xs text-heading-muted">
                                Your public username.
                            </p>
                        </div>

                        <div class="flex items-center gap-3">
                            <button
                                v-if="editing.username"
                                class="btn-primary text-sm"
                                @click="saveField('username')"
                            >
                                Save
                            </button>

                            <button
                                class="text-sm font-medium text-primary hover:underline"
                                @click="toggleEdit('username')"
                            >
                                {{ editing.username ? "Cancel" : "Edit" }}
                            </button>
                        </div>
                    </div>

                    <input
                        v-model="form.username"
                        type="text"
                        :disabled="!editing.username"
                        placeholder="Username"
                        class="w-full rounded-xl border border-primary/15 bg-surface-muted/40 px-4 py-3 text-sm text-heading outline-none transition placeholder:text-heading-muted focus:border-primary-dark/40 focus:bg-white disabled:cursor-not-allowed disabled:opacity-60"
                    />
                </div>

                <!-- Email (Read-only) -->
                <div class="card-surface p-5 text-heading-muted/80">
                    <div class="mb-3">
                        <h2 class="text-sm font-semibold text-heading">
                            Email
                        </h2>

                        <p class="text-xs text-heading-muted">
                            Used for authentication and notifications.
                        </p>
                    </div>

                    <input
                        v-model="form.email"
                        type="email"
                        disabled
                        placeholder="Email address"
                        class="w-full rounded-xl border border-primary/15 bg-surface-muted/20 px-4 py-3 text-sm text-heading-muted outline-none transition cursor-not-allowed opacity-60"
                    />
                </div>

                <!-- Password -->
                <div class="card-surface p-5">
                    <div class="mb-3 flex items-center justify-between gap-4">
                        <div>
                            <h2 class="text-sm font-semibold text-heading">
                                Password
                            </h2>

                            <p class="text-xs text-heading-muted">
                                Update your account password.
                            </p>
                        </div>

                        <div class="flex items-center gap-3">
                            <button
                                v-if="editing.password"
                                class="btn-primary text-sm"
                                @click="saveField('password')"
                            >
                                Save
                            </button>

                            <button
                                class="text-sm font-medium text-primary hover:underline"
                                @click="toggleEdit('password')"
                            >
                                {{ editing.password ? "Cancel" : "Edit" }}
                            </button>
                        </div>
                    </div>

                    <input
                        v-model="form.password"
                        type="password"
                        :disabled="!editing.password"
                        placeholder="New password"
                        class="w-full rounded-xl border border-primary/15 bg-surface-muted/40 px-4 py-3 text-sm text-heading outline-none transition placeholder:text-heading-muted focus:border-primary-dark/40 focus:bg-white disabled:cursor-not-allowed disabled:opacity-60"
                    />
                </div>

                <!-- Bio -->
                <div class="card-surface p-5">
                    <div class="mb-3 flex items-center justify-between gap-4">
                        <div>
                            <h2 class="text-sm font-semibold text-heading">Bio</h2>

                            <p class="text-xs text-heading-muted">
                                Tell people a little about yourself.
                            </p>
                        </div>

                        <div class="flex items-center gap-3">
                            <button
                                v-if="editing.bio"
                                class="btn-primary text-sm"
                                @click="saveField('bio')"
                            >
                                Save
                            </button>

                            <button
                                class="text-sm font-medium text-primary hover:underline"
                                @click="toggleEdit('bio')"
                            >
                                {{ editing.bio ? "Cancel" : "Edit" }}
                            </button>
                        </div>
                    </div>

                    <textarea
                        v-model="form.bio"
                        rows="5"
                        :disabled="!editing.bio"
                        placeholder="Write your bio..."
                        class="w-full resize-none rounded-xl border border-primary/15 bg-surface-muted/40 px-4 py-3 text-sm text-heading outline-none transition placeholder:text-heading-muted focus:border-primary-dark/40 focus:bg-white disabled:cursor-not-allowed disabled:opacity-60"
                    />
                </div>

            </section>
        </div>
    </div>
</template>
