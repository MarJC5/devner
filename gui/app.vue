<script setup lang="ts">
import { useStatsStore } from '~/stores/stats';

const statsStore = useStatsStore();
const colorMode = useColorMode()

const color = computed(() => {
  return colorMode.value === 'dark' ? '#000000' : '#ffffff'
});

useHead({
  meta: [
    { charset: 'utf-8' },
    { name: 'viewport', content: 'width=device-width, initial-scale=1' },
    { key: 'theme-color', name: 'theme-color', content: color }
  ],
  link: [
    { rel: 'icon', href: '/favicon.ico' }
  ],
  htmlAttrs: {
    lang: 'en'
  }
})

const title = 'Devner - The ultimate development tool'
const description = 'Devner is a tool that helps you manage your development environment with ease. It is the ultimate tool for developers who want to focus on coding and not on managing their development environment'

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
})

onMounted(() => {
  statsStore.startTracking();
})

onBeforeUnmount(() => {
  statsStore.stopTracking();
});
</script>


<template>
  <div>
    <NuxtLoadingIndicator />

    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    
    <UNotifications />
  </div>
</template>