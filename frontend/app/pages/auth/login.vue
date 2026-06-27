<script setup>
import { useForm } from "vee-validate";
import { useNotifications } from "~/composables/useNotifications";
import { LOGIN_USER } from "~/graphql/mutations.js";

const { $apollo: apollo } = useNuxtApp();
const { login } = useAuth();
const { notifySuccess, notifyError } = useNotifications();

definePageMeta({
    layout: "default",
});

useHead({
    title: "Login — Share Plate",
});

/*
username, password
1. Initialize the form brain with simple inline validation rules
*/
const { errors, defineField, handleSubmit, isSubmitting } = useForm({
    validationSchema: {
        username: (val) => (val?.trim() ? true : "Username is required"),
        password: (val) => (val && val.length >= 6 ? true : "Password must be at least 6 characters"),
    },
});

// Define the form fields and their binding attributes
const [username, usernameAttrs] = defineField("username");
const [password, passwordAttrs] = defineField("password");

// Wrap submission function inside Vee-Validate's handleSubmit gatekeeper
const onSubmit = handleSubmit(async (values) => {
    try {
        const result = await apollo.mutate({
            mutation: LOGIN_USER,
            variables: {
                username: values.username,
                password: values.password,
            },
        });

        const token = result.data.loginUser.token;

        if (token) {
            notifySuccess("Welcome back! You have successfully logged in.");
            login(token);
            console.log(token);
        } else {
            notifyError("Login failed. Please check your credentials.");
        }
    } catch (err) {
        notifyError(
            err?.message || "An unexpected error occurred during login.",
        );
        console.error(err);
    }
});
</script>

<template>
    <div class="relative overflow-hidden py-12 sm:py-16 lg:py-20">
        <div
            class="pointer-events-none absolute inset-0 -z-10"
            aria-hidden="true"
        />


        <div class="page-container">
            <div class="mx-auto max-w-md">
                <div class="text-center">
                    <p
                        class="text-sm font-semibold uppercase tracking-wider text-primary"
                    >
                        Welcome back
                    </p>
                    <h1 class="mt-2 text-3xl font-semibold sm:text-4xl">
                        Log in to Share Plate
                    </h1>
                    <p class="mt-2 text-sm text-heading-muted">
                        Save bookmarks, rate recipes, and share your own
                        creations.
                    </p>
                </div>

                <form
                    class="card-surface mt-8 space-y-5 p-6 sm:p-8"
                    @submit.prevent="onSubmit"
                >
                    <div>
                        <label
                            for="username"
                            class="mb-1.5 block text-sm font-medium text-heading"
                        >
                            Username
                        </label>
                        <input
                            id="username"
                            v-model="username"
                            v-bind="usernameAttrs"
                            type="text"
                            autocomplete="username"
                            placeholder="username"
                            class="input-field"
                            :class="{ 'border-danger focus:ring-danger': errors.username }"
                        />
                        <!-- Error Message Display -->
                        <p v-if="errors.username" class="mt-1 text-xs text-danger">
                            {{ errors.username }}
                        </p>
                    </div>

                    <div>
                        <div class="mb-1.5 flex items-center justify-between">
                            <label
                                for="password"
                                class="text-sm font-medium text-heading"
                            >
                                Password
                            </label>
                        </div>
                        <input
                            id="password"
                            v-model="password"
                            v-bind="passwordAttrs"
                            type="password"
                            autocomplete="current-password"
                            placeholder="••••••••"
                            class="input-field"
                            :class="{ 'border-danger focus:ring-danger': errors.password }"
                        />
                        <!-- Error Message Display -->
                        <p v-if="errors.password" class="mt-1 text-xs text-danger">
                            {{ errors.password }}
                        </p>
                    </div>

                    <!-- Button disables itself during API execution automatically -->
                    <button 
                        type="submit" 
                        :disabled="isSubmitting"
                        class="btn-primary w-full py-3 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {{ isSubmitting ? 'Logging in...' : 'Log in' }}
                    </button>
                </form>

                <p class="mt-6 text-center text-sm text-heading-muted">
                    Don&apos;t have an account?
                    <NuxtLink
                        to="/auth/register"
                        class="font-semibold text-primary hover:underline"
                    >
                        Sign up free
                    </NuxtLink>
                </p>
            </div>
        </div>
    </div>
</template>