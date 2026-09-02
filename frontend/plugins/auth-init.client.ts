import { useAuthStore } from '~/stores/auth'
import { useNotificationsStore } from '~/stores/notifications'
import { useWebSocket } from '~/composables/useWebSocket'
import { useTheme } from '~/composables/useTheme'

export default defineNuxtPlugin(async (nuxtApp) => {
  const pinia = nuxtApp.$pinia as any
  const auth = useAuthStore(pinia)
  const notifStore = useNotificationsStore(pinia)
  const { initTheme } = useTheme()
  const { connect } = useWebSocket()

  // Initialize theme
  initTheme()

  // Initialize and restore session safely
  try {
    await auth.restoreSession()

    if (auth.isAuthenticated) {
      // Fetch initial notifications
      try {
        notifStore.fetchNotifications()
      } catch (e) {
        console.warn('Initial notifications fetch warning:', e)
      }

      // Connect WebSocket
      try {
        connect()
      } catch (e) {
        console.warn('Initial WebSocket connect warning:', e)
      }
    }
  } catch (err) {
    console.warn('Auth initialization warning:', err)
  }
})
