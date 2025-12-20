<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800 p-4">
    <div class="max-w-md w-full">
      <!-- Card -->
      <div class="bg-white rounded-2xl shadow-2xl p-8">
        <div class="text-center mb-8">
          <div class="text-6xl mb-4">🎵</div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            Connect Spotify
          </h1>
          <p class="text-gray-600">
            Link your Spotify account to view your music trends and top tracks
          </p>
        </div>

        <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-red-800 text-sm">{{ error }}</p>
        </div>

        <div v-if="loading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-gray-200 border-t-green-600"></div>
          <p class="mt-4 text-gray-600">Loading...</p>
        </div>

        <div v-else>
          <a
            v-if="authUrl"
            :href="authUrl"
            class="block w-full bg-gradient-to-r from-green-500 to-green-600 text-white text-center px-6 py-4 rounded-full font-semibold text-lg hover:from-green-600 hover:to-green-700 transition-all transform hover:scale-105 shadow-lg hover:shadow-xl"
          >
            🔗 Link to Spotify Account
          </a>
        </div>

        <div class="mt-6 text-center">
          <NuxtLink
            to="/"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            ← Back to Home
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authUrl = ref<string>('')
const loading = ref(true)
const error = ref<string>('')

const api = useApi()

// SEO
useHead({
  title: 'Authorize Spotify - Spotify Trends',
  meta: [
    { name: 'description', content: 'Authorize your Spotify account to access your music trends' }
  ]
})

onMounted(async () => {
  try {
    const response = await api.getSpotifyAuthUrl()
    authUrl.value = response.authUrl
  } catch (err) {
    console.error('Failed to load auth URL:', err)
    error.value = 'Failed to load authorization URL. Please try again.'
  } finally {
    loading.value = false
  }
})
</script>
