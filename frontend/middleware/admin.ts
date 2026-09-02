import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const nuxtApp = useNuxtApp()
  const auth = useAuthStore(nuxtApp.$pinia as any)

  if (!auth.isInitialized) {
    try {
      await auth.restoreSession()
    } catch (e) {
      console.warn('Middleware admin restore warning:', e)
    }
  }

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (auth.user?.role !== 'admin') {
    return navigateTo('/dashboard')
  }
})
