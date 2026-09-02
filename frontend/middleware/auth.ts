import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async (to) => {
  const nuxtApp = useNuxtApp()
  const auth = useAuthStore(nuxtApp.$pinia as any)

  // Ensure session is initialized on first navigation
  if (!auth.isInitialized) {
    try {
      await auth.restoreSession()
    } catch (e) {
      console.warn('Middleware auth restore warning:', e)
    }
  }

  if (!auth.isAuthenticated) {
    return navigateTo({
      path: '/login',
      query: { redirect: to.fullPath !== '/login' ? to.fullPath : undefined }
    })
  }
})
