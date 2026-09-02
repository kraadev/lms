import { api } from './api'
import type { Notification } from '~/types'

export const notificationsService = {
  getAll(): Promise<{ notifications: Notification[]; unread_count: number }> {
    return api.get<{ notifications: Notification[]; unread_count: number }>('/notifications')
  },

  markAsRead(id: number | string): Promise<{ success: boolean }> {
    return api.patch<{ success: boolean }>(`/notifications/${id}/read`)
  },

  markAllAsRead(): Promise<{ success: boolean }> {
    return api.post<{ success: boolean }>('/notifications/read-all')
  }
}
