import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  // Ensure session is initialized on first navigation
  if (!auth.isInitialized) {
    await auth.restoreSession()
  }

  if (!auth.isAuthenticated) {
    return navigateTo({
      path: '/login',
      query: { redirect: to.fullPath !== '/login' ? to.fullPath : undefined }
    })
  }
})
