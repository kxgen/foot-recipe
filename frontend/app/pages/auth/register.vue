<script setup>
import { useForm } from "vee-validate";
import { useNotifications } from "~/composables/useNotifications";
import { REGISTER_USER } from "~/graphql/mutations.js"

definePageMeta({
    layout: "default",
});

useHead({
    title: "Sign up — Share Plate",
});

const { $apollo: apollo } = useNuxtApp();
const { login: authLogin } = useAuth();
const { notifySuccess, notifyError } = useNotifications();
const submitError = ref("");
const submitSuccess = ref("");

const { errors, defineField, handleSubmit, isSubmitting, resetForm, values } = useForm({
    initialValues: {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
    },
    validationSchema: computed(() => ({
        username: (value) =>
            value && value.trim().length >= 3
                ? true
                : "Username must be at least 3 characters.",
        email: (value) =>
            /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value || "")
                ? true
                : "Please enter a valid email.",
        password: (value) =>
            value && value.length >= 8
                ? true
                : "Password must be at least 8 characters.",
        confirmPassword: (value) => {
            if (!value) return "Please confirm your password.";
            if (value !== values.password) return "Passwords do not match.";
            return true;
        },
    })),
});

const [username, usernameAttrs] = defineField("username");
const [email, emailAttrs] = defineField("email");
const [password, passwordAttrs] = defineField("password");
const [confirmPassword, confirmPasswordAttrs] = defineField("confirmPassword");

const onSubmit = handleSubmit(async (formValues) => {
    submitError.value = "";
    submitSuccess.value = "";

    if (formValues.password !== formValues.confirmPassword) {
        submitError.value = "Passwords do not match.";
        return;
    }

    if (!apollo) {
        notifyError("Apollo client is not configured.");
        return;
    }

    try {
        const { data } = await apollo.mutate({
            mutation: REGISTER_USER,
            variables: {
                username: formValues.username.trim(),
                email: formValues.email.trim().toLowerCase(),
                password: formValues.password,
            },
        });

        if (data?.registerUser?.token) {
            notifySuccess("Account created successfully!");
            authLogin(data.registerUser.token);
        } else {
            notifySuccess("Account created successfully! You can now log in.");
            await navigateTo("/auth/login");
        }
    } catch (error) {
        notifyError(
            error?.message || "Failed to create account. Please try again.",
        );
    }
});
</script>

<template>
    <div class="relative overflow-hidden py-12 sm:py-16 lg:py-20">
        <div
            class="pointer-events-none absolute inset-0 -z-10"
            aria-hidden="true"
        />
        <div
            class="pointer-events-none absolute right-0 top-1/3 size-64 rounded-full bg-accent-subtle/70 blur-3xl"
            aria-hidden="true"
        />

        <div class="page-container">
            <div class="mx-auto max-w-md">
                <div class="text-center">
                    <p
                        class="text-sm font-semibold uppercase tracking-wider text-primary"
                    >
                        Join the kitchen
                    </p>
                    <h1 class="mt-2 text-3xl font-semibold sm:text-4xl">
                        Create your account
                    </h1>
                    <p class="mt-2 text-sm text-subtle">
                        Share recipes, bookmark favorites, and connect with
                        fellow food lovers.
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
                        />
                        <p
                            v-if="errors.username"
                            class="mt-1 text-xs text-danger"
                        >
                            {{ errors.username }}
                        </p>
                    </div>

                    <div>
                        <label
                            for="email"
                            class="mb-1.5 block text-sm font-medium text-heading"
                        >
                            Email
                        </label>
                        <input
                            id="email"
                            v-model="email"
                            v-bind="emailAttrs"
                            type="email"
                            autocomplete="email"
                            placeholder="yourname@company.com"
                            class="input-field"
                        />
                        <p
                            v-if="errors.email"
                            class="mt-1 text-xs text-danger"
                        >
                            {{ errors.email }}
                        </p>
                    </div>

                    <div>
                        <label
                            for="password"
                            class="mb-1.5 block text-sm font-medium text-heading"
                        >
                            New Password
                        </label>
                        <input
                            id="password"
                            v-model="password"
                            v-bind="passwordAttrs"
                            type="password"
                            autocomplete="new-password"
                            placeholder="At least 8 characters"
                            class="input-field"
                        />
                        <p
                            v-if="errors.password"
                            class="mt-1 text-xs text-danger"
                        >
                            {{ errors.password }}
                        </p>
                    </div>

                    <div>
                        <label
                            for="confirmPassword"
                            class="mb-1.5 block text-sm font-medium text-heading"
                        >
                            Confirm Password
                        </label>
                        <input
                            id="confirmPassword"
                            v-model="confirmPassword"
                            v-bind="confirmPasswordAttrs"
                            type="password"
                            autocomplete="new-password"
                            placeholder="Confirm your password"
                            class="input-field"
                        />
                        <p
                            v-if="errors.confirmPassword"
                            class="mt-1 text-xs text-danger"
                        >
                            {{ errors.confirmPassword }}
                        </p>
                    </div>

                    <p
                        v-if="submitError"
                        class="error-message"
                    >
                        {{ submitError }}
                    </p>
                    <p
                        v-if="submitSuccess"
                        class="success-message"
                    >
                        {{ submitSuccess }}
                    </p>

                    <button
                        type="submit"
                        class="btn-primary w-full py-3"
                        :disabled="isSubmitting"
                        :class="{
                            'cursor-not-allowed opacity-60': isSubmitting,
                        }"
                    >
                        {{
                            isSubmitting
                                ? "Creating account..."
                                : "Create account"
                        }}
                    </button>
                </form>

                <p class="mt-6 text-center text-sm text-subtle">
                    Already have an account?
                    <NuxtLink
                        to="/auth/login"
                        class="font-semibold text-primary hover:underline"
                    >
                        Log in
                    </NuxtLink>
                </p>
            </div>
        </div>
    </div>
</template>
