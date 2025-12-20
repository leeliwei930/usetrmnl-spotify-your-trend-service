<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800 p-4">
    <div class="max-w-2xl w-full">
      <!-- Loading State -->
      <div v-if="loading" class="bg-white rounded-2xl shadow-2xl p-12 text-center">
        <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-gray-200 border-t-green-600 mb-4"></div>
        <h2 class="text-2xl font-semibold text-gray-900">Processing...</h2>
        <p class="text-gray-600 mt-2">Connecting your Spotify account</p>
      </div>

      <!-- Success State -->
      <div
        v-else-if="success"
        class="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl shadow-2xl p-10 text-white"
      >
        <div class="text-center mb-8">
          <div class="text-6xl mb-4">✅</div>
          <h1 class="text-4xl font-bold mb-4">
            Successfully Connected!
          </h1>
          <p class="text-green-50 text-lg">
            Your Spotify account has been linked. Below is your refresh token.
          </p>
        </div>

        <div class="bg-white/10 backdrop-blur-sm rounded-xl p-6 mb-6">
          <label class="block text-green-50 text-sm font-medium mb-3">
            Refresh Token:
          </label>
          <textarea
            id="refreshTokenDisplay"
            :value="refreshToken"
            readonly
            class="w-full bg-black/20 text-white border border-white/20 rounded-lg p-4 font-mono text-sm resize-none min-h-[120px] focus:outline-none focus:ring-2 focus:ring-white/50"
          />

          <button
            @click="copyToClipboard"
            :class="[
              'mt-4 w-full px-6 py-3 rounded-full font-semibold text-lg transition-all transform hover:scale-105 shadow-lg',
              copied
                ? 'bg-white text-green-600'
                : 'bg-white text-green-600 hover:bg-green-50'
            ]"
          >
            {{ copied ? '✓ Copied!' : '📋 Copy to Clipboard' }}
          </button>
        </div>

        <div class="border-t border-white/20 pt-6 text-center">
          <p class="text-green-50 text-sm">
            Save this refresh token securely. You'll need it to access your Spotify data.
          </p>
          <NuxtLink
            to="/"
            class="inline-block mt-4 text-white hover:text-green-100 font-medium transition-colors"
          >
            ← Back to Home
          </NuxtLink>
        </div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="errorMessage"
        class="bg-gradient-to-br from-red-500 to-red-600 rounded-2xl shadow-2xl p-10 text-white"
      >
        <div class="text-center mb-8">
          <div class="text-6xl mb-4">❌</div>
          <h1 class="text-4xl font-bold mb-4">
            Connection Failed
          </h1>
          <p class="text-red-50 text-lg mb-8">
            {{ errorMessage }}
          </p>

          <NuxtLink
            to="/spotify/authorize"
            class="inline-block bg-white text-red-600 px-12 py-4 rounded-full font-semibold text-lg hover:bg-red-50 transition-all transform hover:scale-105 shadow-lg"
          >
            🔄 Try Again
          </NuxtLink>
        </div>

        <div class="text-center mt-6">
          <NuxtLink
            to="/"
            class="text-white hover:text-red-100 font-medium transition-colors"
          >
            ← Back to Home
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()

const loading = ref(true)
const success = ref(false)
const errorMessage = ref('')
const refreshToken = ref('')
const accessToken = ref('')
const expiresIn = ref(0)
const copied = ref(false)

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

onMounted(async () => {
  const code = route.query.code as string
  const error = route.query.error as string
  const state = route.query.state as string

  // Check if user denied authorization
  if (error) {
    errorMessage.value = 'Authorization was denied. Please try again.'
    loading.value = false
    return
  }

  // Validate code parameter
  if (!code) {
    errorMessage.value = 'Missing authorization code. Please start the authorization process again.'
    loading.value = false
    return
  }

  try {
    const response = await api.exchangeSpotifyCode(code, state)
    
    if (response.success) {
      success.value = true
      refreshToken.value = response.refreshToken || ''
      accessToken.value = response.accessToken || ''
      expiresIn.value = response.expiresIn || 0
    } else {
      errorMessage.value = response.error || 'Failed to exchange authorization code for tokens'
    }
  } catch (err) {
    console.error('Exchange failed:', err)
    errorMessage.value = 'Failed to exchange authorization code for tokens. Please try again.'
  } finally {
    loading.value = false
  }
})
</script>
