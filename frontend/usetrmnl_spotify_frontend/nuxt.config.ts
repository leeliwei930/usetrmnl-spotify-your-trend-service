import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: process.env.NODE_ENV !== 'production' },
  ssr: false,
  modules: ['@nuxt/eslint', '@nuxt/test-utils', 'shadcn-nuxt'],
  css: ['~/assets/css/tailwind.css'],
  plugins: ['~/plugins/vue-query.ts'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  shadcn: {
    prefix: '',
    componentDir: '~/components/ui'
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8009/api'
    }
  }
})
