<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800">
    <div class="max-w-2xl mx-auto p-8 text-center">
      <h1 class="text-5xl font-bold text-white mb-4">
        Welcome!
      </h1>
      <p class="text-xl text-slate-300 mb-8">
        Connect your Spotify account to get started with your music trends.
      </p>
      <NuxtLink
        to="/usetrmnl/spotify/authorize"
        class="inline-block bg-gradient-to-r from-green-500 to-green-600 text-white px-8 py-4 rounded-full font-semibold text-lg hover:from-green-600 hover:to-green-700 transition-all transform hover:scale-105 shadow-lg hover:shadow-xl"
      >
        Get Started →
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
// SEO
useHead({
  title: 'Home - Spotify Trends',
  meta: [
    { name: 'description', content: 'Connect your Spotify account to view your music trends and top tracks' }
  ]
})
</script>
