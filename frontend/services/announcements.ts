import { api } from './api'
import type { Announcement } from '~/types'

export const announcementsService = {
  getByClass(classId: number | string): Promise<Announcement[]> {
    return api.get<Announcement[]>(`/classes/${classId}/announcements`)
  },

  create(classId: number | string, payload: { title: string; content: string }): Promise<Announcement> {
    return api.post<Announcement>(`/classes/${classId}/announcements`, payload)
  },

  delete(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/announcements/${id}`)
  }
}
