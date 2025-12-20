export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  const getSpotifyAuthUrl = async () => {
    try {
      const response = await $fetch<{ authUrl: string }>(`${apiBase}/usetrmnl/spotify/authorize`)
      return response
    } catch (error) {
      console.error('Failed to fetch auth URL:', error)
      throw error
    }
  }

  const exchangeSpotifyCode = async (code: string, state?: string) => {
    try {
      const response = await $fetch<{
        success: boolean
        refreshToken?: string
        accessToken?: string
        expiresIn?: number
        error?: string
      }>(`${apiBase}/usetrmnl/spotify/callback`, {
        method: 'POST',
        body: { code, state }
      })
      return response
    } catch (error) {
      console.error('Failed to exchange code:', error)
      throw error
    }
  }

  return {
    getSpotifyAuthUrl,
    exchangeSpotifyCode
  }
}
