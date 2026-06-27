// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2026-06-11",
  devtools: { enabled: true },
  css: ["~/assets/css/main.css"],

  runtimeConfig: {
    public: {
      hasuraUrl: "",
      imageBase: "",
    },
  },

  app: {
    head: {
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "",
        },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Rubik+Spray+Paint&display=swap",
        },
      ],
    },
  },

  vite: {
    plugins: [tailwindcss()],
    optimizeDeps: {
      include: [
        "@apollo/client",
        "@apollo/client/core",
        "@apollo/client/link/context",
        "@apollo/client/link/error",
        "@lucide/vue",
        "@vue/devtools-core",
        "@vue/devtools-kit",
        "@vueuse/core",
        "vee-validate",
      ],
    },
  },

  modules: ["@pinia/nuxt"],
});
