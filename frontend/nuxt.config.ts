// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: false },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],

  css: [
    '~/assets/css/main.css'
  ],

  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api',
      wsUrl: process.env.NUXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'
    }
  },

  app: {
    head: {
      title: 'LMS - Learning Management System',
      titleTemplate: '%s | LMS',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Modern, intuitive, and fast Learning Management System for students, teachers, and administrators.' },
        { name: 'theme-color', content: '#4f46e5' }
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap'
        }
      ]
    }
  },

  devServer: {
    port: 3000,
    host: '0.0.0.0'
  },

  typescript: {
    strict: true,
    shim: false
  }
})
