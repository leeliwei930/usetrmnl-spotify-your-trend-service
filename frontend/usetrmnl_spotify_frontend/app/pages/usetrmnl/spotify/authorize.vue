<template>
  <div class="min-h-screen flex items-center justify-center bg-linear-to-br from-slate-900 to-slate-800 p-4">
    <div class="max-w-md w-full">
      <Card class="shadow-2xl">
        <CardHeader class="text-center">
          <div class="text-6xl mb-4">🎵</div>
          <CardTitle class="text-3xl mb-2">Connect Spotify</CardTitle>
          <CardDescription class="text-base">
            Link your Spotify account to view your music trends and top tracks
          </CardDescription>
        </CardHeader>

        <CardContent>
          <Alert v-if="error" variant="destructive" class="mb-6">
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>

          <div v-if="loading" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-gray-200 border-t-green-600"></div>
            <p class="mt-4 text-gray-600">Loading...</p>
          </div>

          <div v-else>
            <Button
              v-if="authUrl"
              as-child
              size="lg"
              class="w-full bg-linear-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white font-semibold text-lg transition-all transform hover:scale-105 shadow-lg hover:shadow-xl"
            >
              <a :href="authUrl">
                🔗 Link to Spotify Account
              </a>
            </Button>
          </div>

          <div class="mt-6 text-center">
            <NuxtLink
              to="/"
              class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
            >
              ← Back to Home
            </NuxtLink>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSpotifyAuthUrl } from '~/composables/useApi'

// Use Vue Query to fetch auth URL
const { data, isLoading, isError } = useSpotifyAuthUrl()

// Computed properties for template
const authUrl = computed(() => data.value?.authUrl || '')
const loading = computed(() => isLoading.value)
const error = computed(() => isError.value ? 'Failed to load authorization URL. Please try again.' : '')

// SEO
useHead({
  title: 'Authorize Spotify - Spotify Trends',
  meta: [
    { name: 'description', content: 'Authorize your Spotify account to access your music trends' }
  ]
})
</script>
