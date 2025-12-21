import { useQuery, useMutation } from '@tanstack/vue-query'

// Query for getting Spotify auth URL
export const useSpotifyAuthUrl = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  return useQuery({
    queryKey: ['spotify', 'authUrl'],
    queryFn: async () => {
      const response = await $fetch<{ authUrl: string }>(`${apiBase}/usetrmnl/spotify/authorize`)
      return response
    },
    staleTime: 1000 * 60 * 5, // Consider the auth URL fresh for 5 minutes
    retry: 2
  })
}

// Mutation for exchanging Spotify code
export const useExchangeSpotifyCode = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  return useMutation({
    mutationFn: async (params: { code: string; state?: string }) => {
      const response = await $fetch<{
        success: boolean
        refreshToken?: string
        accessToken?: string
        expiresIn?: number
        error?: string
      }>(`${apiBase}/usetrmnl/spotify/callback`, {
        method: 'POST',
        body: { code: params.code, state: params.state }
      })
      return response
    }
  })
}
