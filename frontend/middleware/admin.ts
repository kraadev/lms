import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()

  if (!auth.isInitialized) {
    await auth.restoreSession()
  }

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (auth.user?.role !== 'admin') {
    return navigateTo('/dashboard')
  }
})
