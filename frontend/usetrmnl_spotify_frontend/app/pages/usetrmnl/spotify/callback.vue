<template>
  <div class="min-h-screen flex items-center justify-center bg-linear-to-br from-slate-900 to-slate-800 p-4">
    <div class="max-w-2xl w-full">
      <!-- Loading State -->
      <Card v-if="loading" class="shadow-2xl">
        <CardContent class="p-12 text-center">
          <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-gray-200 border-t-green-600 mb-4"></div>
          <h2 class="text-2xl font-semibold text-gray-900">Processing...</h2>
          <p class="text-gray-600 mt-2">Connecting your Spotify account</p>
        </CardContent>
      </Card>

      <!-- Success State -->
      <div
        v-else-if="success"
        class="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl shadow-2xl p-10 text-white"
      >
        <Card class="bg-transparent border-white/20 text-white">
          <CardHeader class="text-center">
            <div class="text-6xl mb-4">✅</div>
            <CardTitle class="text-4xl mb-4 text-white">Successfully Connected!</CardTitle>
            <CardDescription class="text-green-50 text-lg">
              Your Spotify account has been linked. Below is your refresh token.
            </CardDescription>
          </CardHeader>

          <CardContent>
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-6 mb-6">
              <label class="block text-green-50 text-sm font-medium mb-3">
                Refresh Token:
              </label>
              <Textarea
                id="refreshTokenDisplay"
                :model-value="refreshToken"
                readonly
                class="w-full bg-black/20 text-white border-white/20 rounded-lg font-mono text-sm resize-none min-h-[120px] focus:ring-2 focus:ring-white/50"
              />

              <Button
                @click="copyToClipboard"
                size="lg"
                :class="[
                  'mt-4 w-full font-semibold text-lg transition-all transform hover:scale-105 shadow-lg',
                  copied
                    ? 'bg-white text-green-600'
                    : 'bg-white text-green-600 hover:bg-green-50'
                ]"
              >
                {{ copied ? '✓ Copied!' : '📋 Copy to Clipboard' }}
              </Button>
            </div>

            <Alert class="bg-white/10 border-white/20">
              <AlertDescription class="text-green-50 text-sm">
                Save this refresh token securely. You'll need it to access your Spotify data.
              </AlertDescription>
            </Alert>

            <div class="text-center mt-6">
              <NuxtLink
                to="/"
                class="inline-block text-white hover:text-green-100 font-medium transition-colors"
              >
                ← Back to Home
              </NuxtLink>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Error State -->
      <div
        v-else-if="errorMessage"
        class="bg-linear-to-br from-gray-500 to-gray-600 rounded-2xl shadow-2xl p-10 text-white"
      >
        <Card class="bg-transparent border-white/20 text-white">
          <CardHeader class="text-center">
            <div class="text-6xl mb-4">❌</div>
            <CardTitle class="text-4xl mb-4 text-white">Connection Failed</CardTitle>
            <CardDescription class="text-red-50 text-lg mb-8">
              {{ errorMessage }}
            </CardDescription>
          </CardHeader>

          <CardContent>
            <Button
              as-child
              size="lg"
              class="w-full bg-white text-red-600 hover:bg-red-50 font-semibold text-lg transition-all transform hover:scale-105 shadow-lg mb-6"
            >
              <NuxtLink to="/usetrmnl/spotify/authorize">
                🔄 Try Again
              </NuxtLink>
            </Button>

            <div class="text-center">
              <NuxtLink
                to="/"
                class="text-white hover:text-red-100 font-medium transition-colors"
              >
                ← Back to Home
              </NuxtLink>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useExchangeSpotifyCode } from '~/composables/useApi'

const route = useRoute()

const success = ref(false)
const errorMessage = ref('')
const refreshToken = ref('')
const accessToken = ref('')
const expiresIn = ref(0)
const copied = ref(false)

// Use Vue Query mutation for exchanging code
const { mutate, isPending } = useExchangeSpotifyCode()

// Computed for loading state
const loading = computed(() => isPending.value)

// SEO
useHead({
  title: 'Spotify Callback - Spotify Trends',
  meta: [
    { name: 'description', content: 'Processing your Spotify authorization' }
  ]
})

const copyToClipboard = async () => {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(refreshToken.value)
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    } else {
      // Fallback for older browsers
      const textarea = document.getElementById('refreshTokenDisplay') as HTMLTextAreaElement
      textarea.select()
      document.execCommand('copy')
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    }
  } catch (err) {
    console.error('Failed to copy:', err)
    alert('Please manually copy the token')
  }
}

onMounted(() => {
  const code = route.query.code as string
  const error = route.query.error as string
  const state = route.query.state as string

  // Check if user denied authorization
  if (error) {
    errorMessage.value = 'Authorization was denied. Please try again.'
    return
  }

  // Validate code parameter
  if (!code) {
    errorMessage.value = 'Missing authorization code. Please start the authorization process again.'
    return
  }

  // Execute mutation
  mutate(
    { code, state },
    {
      onSuccess: (response) => {
        if (response.success) {
          success.value = true
          refreshToken.value = response.refreshToken || ''
          accessToken.value = response.accessToken || ''
          expiresIn.value = response.expiresIn || 0
        } else {
          errorMessage.value = response.error || 'Failed to exchange authorization code for tokens'
        }
      },
      onError: (err) => {
        console.error('Exchange failed:', err)
        errorMessage.value = 'Failed to exchange authorization code for tokens. Please try again.'
      }
    }
  )
})
</script>
