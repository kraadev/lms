import { api } from './api'
import type { ChatMessage } from '~/types'

export const messagesService = {
  getByClass(classId: number | string, params?: { limit?: number; offset?: number }): Promise<ChatMessage[]> {
    return api.get<ChatMessage[]>(`/classes/${classId}/messages`, params)
  },

  send(classId: number | string, message: string): Promise<ChatMessage> {
    return api.post<ChatMessage>(`/classes/${classId}/messages`, { message })
  }
}
