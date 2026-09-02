import { useAuthStore } from '~/stores/auth'
import { useNotificationsStore } from '~/stores/notifications'
import { useWebSocket } from '~/composables/useWebSocket'
import { useTheme } from '~/composables/useTheme'

export default defineNuxtPlugin(async () => {
  const auth = useAuthStore()
  const notifStore = useNotificationsStore()
  const { initTheme } = useTheme()
  const { connect } = useWebSocket()

  // Initialize theme
  initTheme()

  // Initialize and restore session
  await auth.restoreSession()

  if (auth.isAuthenticated) {
    // Fetch initial notifications
    notifStore.fetchNotifications()
    // Connect WebSocket
    connect()
  }
})
