import { defineStore } from 'pinia'
import { notificationsService } from '~/services/notifications'
import type { Notification } from '~/types'

export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref<Notification[]>([])
  const unreadCount = ref(0)
  const isLoading = ref(false)

  async function fetchNotifications() {
    isLoading.value = true
    try {
      const data = await notificationsService.getAll()
      notifications.value = data.notifications || []
      unreadCount.value = data.unread_count || 0
    } catch (err) {
      console.error('Failed to fetch notifications:', err)
    } finally {
      isLoading.value = false
    }
  }

  async function markAsRead(id: number) {
    const notif = notifications.value.find(n => n.id === id)
    if (notif && !notif.is_read) {
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
      try {
        await notificationsService.markAsRead(id)
      } catch (err) {
        console.error('Failed to mark notification as read:', err)
      }
    }
  }

  async function markAllAsRead() {
    notifications.value.forEach(n => { n.is_read = true })
    unreadCount.value = 0
    try {
      await notificationsService.markAllAsRead()
    } catch (err) {
      console.error('Failed to mark all notifications as read:', err)
    }
  }

  function addNotification(notif: Notification) {
    notifications.value.unshift(notif)
    if (!notif.is_read) {
      unreadCount.value += 1
    }
  }

  return {
    notifications,
    unreadCount,
    isLoading,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    addNotification
  }
})
