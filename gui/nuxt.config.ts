// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  extends: ['@nuxt/ui-pro'],
  modules: ['@nuxt/ui', '@pinia/nuxt'],
  nitro: {
    experimental: {
      websocket: true
    }
  },
  pinia: {
    storesDirs: ['stores/**'],
  },
  colorMode: {
    preference: 'light'
  },
  ui: {
    safelistColors: ['primary', 'red', 'orange', 'green', 'blue', 'indigo', 'purple', 'pink', 'gray', 'yellow', 'teal', 'cyan', 'white', 'slate'],
  },
  ssr: false,
  devtools: { enabled: false }
})